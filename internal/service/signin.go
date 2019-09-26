package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/rediscache"
	"github.com/hpifu/go-account/internal/rule"
	"github.com/sirupsen/logrus"
)

type SignInReq struct {
	Username string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

type SignInRes string

func (s *Service) SignIn(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &SignInReq{}
	var res SignInRes
	var err error
	status := http.StatusOK

	defer func() {
		AccessLog.WithFields(logrus.Fields{
			"host":   c.Request.Host,
			"url":    c.Request.URL.String(),
			"req":    req,
			"res":    res,
			"rid":    rid,
			"err":    err,
			"status": status,
		}).Info()
	}()

	if err := c.Bind(req); err != nil {
		err = fmt.Errorf("bind json failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = s.checkSignInReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.signIn(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("signIn failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	if res != "" {
		// 最后一个参数是 httponly 需要设置成 false 否则 axios 不能访问到 cookie
		c.SetCookie("token", string(res), 7*24*3600, "/", s.domain, s.secure, false)
		status = http.StatusOK
	} else {
		status = http.StatusForbidden
	}
	c.Status(status)
}

func (s *Service) checkSignInReqBody(req *SignInReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Username: {rule.Required},
		req.Password: {rule.Required, rule.AtLeast8Characters},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) signIn(req *SignInReq) (SignInRes, error) {
	account, err := s.db.SelectAccountByPhoneOrEmail(req.Username)
	if err != nil {
		return "", err
	}

	if account == nil {
		return "", nil
	}

	if account.Password != req.Password {
		return "", nil
	}

	token := NewToken()
	if err := s.cache.SetAccount(token, rediscache.NewAccount(account)); err != nil {
		return "", err
	}

	return SignInRes(token), nil
}
