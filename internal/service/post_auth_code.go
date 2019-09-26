package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/mail"
	"github.com/hpifu/go-account/internal/rule"
	"github.com/sirupsen/logrus"
	"net/http"
)

type POSTAuthCodeReq struct {
	Type      string `json:"type" uri:"type"`
	FirstName string `json:"firstName,omitempty" form:"firstName"`
	LastName  string `json:"lastName,omitempty" form:"lastName"`
	Email     string `json:"email,omitempty" form:"email"`
	Phone     string `json:"phone,omitempty" form:"phone"`
}

type POSTAuthCodeRes struct{}

func (s *Service) POSTAuthCode(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &POSTAuthCodeReq{}
	var err error
	status := http.StatusOK

	defer func() {
		AccessLog.WithFields(logrus.Fields{
			"host":   c.Request.Host,
			"url":    c.Request.URL.String(),
			"req":    req,
			"res":    nil,
			"rid":    rid,
			"err":    err,
			"status": status,
		}).Info()
	}()

	if err := c.BindUri(req); err != nil {
		err = fmt.Errorf("bind failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err := c.Bind(req); err != nil {
		err = fmt.Errorf("bind json failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = s.checkPOSTAuthCodeReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	err = s.postAuthCode(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("postAuthCode failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusCreated
	c.Status(status)
}

func (s *Service) checkPOSTAuthCodeReqBody(req *POSTAuthCodeReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Type: {rule.Required, rule.In(map[interface{}]struct{}{
			"email": {}, "phone": {},
		})},
		req.FirstName: {rule.Required, rule.AtMost32Characters},
		req.LastName:  {rule.Required, rule.AtMost32Characters},
	}); err != nil {
		return err
	}

	switch req.Type {
	case "phone":
		return rule.Check(map[interface{}][]rule.Rule{
			req.Phone: {rule.Required, rule.ValidPhone},
		})
	case "email":
		return rule.Check(map[interface{}][]rule.Rule{
			req.Email: {rule.Required, rule.ValidEmail, rule.AtMost64Characters},
		})
	}

	return nil
}

func (s *Service) postAuthCode(req *POSTAuthCodeReq) error {
	code, err := s.cache.GetAuthCode(req.Email)
	if err != nil {
		return err
	}
	if code == "" {
		code = NewCode()
	}
	if req.Type == "email" {
		if err := s.cache.SetAuthCode(req.Email, code); err != nil {
			return err
		}
		if err := s.mc.Send(req.Email, "hpifu 账号验证", mail.NewAuthCodeTpl(req.FirstName, req.LastName, code)); err != nil {
			return err
		}
	}

	if req.Type == "phone" {
		if err := s.cache.SetAuthCode(req.Phone, code); err != nil {
			return err
		}
	}

	return nil
}
