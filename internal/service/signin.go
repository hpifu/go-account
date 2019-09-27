package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/rediscache"
	"github.com/hpifu/go-account/internal/rule"
)

type SignInReq struct {
	Username string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

type SignInRes string

func (s *Service) ProcessSignIn(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &SignInReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validSignIn(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	account, err := s.db.SelectAccountByPhoneOrEmail(req.Username)
	if err != nil {
		return req, nil, http.StatusInternalServerError, err
	}

	if account == nil || account.Password != req.Password {
		return req, "密码错误", http.StatusOK, nil
	}

	token := NewToken()
	if err := s.cache.SetAccount(token, rediscache.NewAccount(account)); err != nil {
		return req, nil, http.StatusInternalServerError, err
	}

	// 最后一个参数是 httponly 需要设置成 false 否则 axios 不能访问到 cookie
	c.SetCookie("token", token, 7*24*3600, "/", s.domain, s.secure, false)
	return req, nil, http.StatusOK, nil
}

func (s *Service) validSignIn(req *SignInReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Username: {rule.Required},
		req.Password: {rule.Required, rule.AtLeast8Characters},
	}); err != nil {
		return err
	}

	return nil
}
