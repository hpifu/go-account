package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/rule"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SignOutReq struct {
	Token string `json:"token,omitempty" uri:"token"`
}

type SignOutRes struct{}

func (s *Service) SignOut(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &SignOutReq{}
	var err error
	var buf []byte
	status := http.StatusOK

	defer func() {
		AccessLog.WithFields(logrus.Fields{
			"host":   c.Request.Host,
			"body":   string(buf),
			"url":    c.Request.URL.String(),
			"req":    req,
			"res":    nil,
			"rid":    rid,
			"err":    err,
			"status": status,
		}).Info()
	}()

	if err := c.BindUri(req); err != nil {
		err = fmt.Errorf("bind uri failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = s.checkSignOutReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	err = s.signOut(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("signOut failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusAccepted
	c.Status(status)
}

func (s *Service) checkSignOutReqBody(req *SignOutReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) signOut(req *SignOutReq) error {
	return s.cache.DelAccount(req.Token)
}
