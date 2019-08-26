package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/rediscache"
	"github.com/hpifu/go-account/internal/rule"
	api "github.com/hpifu/go-account/pkg/account"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Service) SignIn(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &api.SignInReq{}
	var res *api.SignInRes
	var err error
	var buf []byte
	status := http.StatusOK

	defer func() {
		AccessLog.WithFields(logrus.Fields{
			"host":   c.Request.Host,
			"body":   string(buf),
			"url":    c.Request.URL.String(),
			"req":    req,
			"res":    res,
			"rid":    rid,
			"err":    err,
			"status": status,
		}).Info()
	}()

	if err := c.BindJSON(req); err != nil {
		err = fmt.Errorf("bind json failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = s.checkSignInReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
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

	if res.Valid {
		// 最后一个参数是 httponly 需要设置成 false 否则 axios 不能访问到 cookie
		c.SetCookie("token", res.Token, 7*24*3600, "/", "127.0.0.1", false, false)
	}

	status = http.StatusOK
	c.JSON(status, res)
}

func (s *Service) checkSignInReqBody(req *api.SignInReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Username: {rule.Required},
		req.Password: {rule.Required, rule.AtLeast8Characters},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) signIn(req *api.SignInReq) (*api.SignInRes, error) {
	account, err := s.db.SelectAccountByPhoneOrEmail(req.Username)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return &api.SignInRes{Valid: false}, nil
	}

	if account.Password != req.Password {
		return &api.SignInRes{Valid: false}, nil
	}

	token := NewToken()
	if err := s.cache.SetAccount(token, rediscache.NewAccount(account)); err != nil {
		return nil, err
	}

	return &api.SignInRes{Valid: true, Token: token}, nil
}
