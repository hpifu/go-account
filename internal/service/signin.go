package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/redis"
	"github.com/hpifu/go-kit/rule"
)

type SignInReq struct {
	Username string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

type SignInRes string

func (s *Service) SignIn(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &SignInReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validSignIn(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	account, err := s.db.SelectAccountByPhoneOrEmail(req.Username)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select account failed. err: [%v]", err)
	}

	if account == nil || account.Password != req.Password {
		return req, "密码错误", http.StatusOK, nil
	}

	token := NewToken()
	if err := s.cache.SetAccount(token, redis.NewAccount(account)); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("redis set account falied. err: [%v]", err)
	}

	// 最后一个参数是 httponly 需要设置成 false 否则 axios 不能访问到 cookie
	c.SetCookie("token", token, 7*24*3600, "/", s.domain, s.secure, false)
	return req, nil, http.StatusOK, nil
}

func (s *Service) validSignIn(req *SignInReq) error {
	if err := rule.Check([][3]interface{}{
		{"username", req.Username, []rule.Rule{rule.Required}},
		{"password", req.Password, []rule.Rule{rule.Required, rule.AtLeast8Characters}},
	}); err != nil {
		return err
	}

	return nil
}
