package account

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/account/internal/mail"
	"github.com/hpifu/account/internal/rule"
	"github.com/sirupsen/logrus"
	"net/http"
)

type GenAuthCodeReqBody struct {
	Type      string `json:"type"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type GenAuthCodeResBody struct {
	OK bool `json:"ok"`
}

func (s *Service) GenAuthCode(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &GenAuthCodeReqBody{}
	var res *GenAuthCodeResBody
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

	buf, err = c.GetRawData()
	if err != nil {
		err = fmt.Errorf("get raw data failed, err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = json.Unmarshal(buf, req); err != nil {
		err = fmt.Errorf("json unmarshal body failed. body: [%v], err: [%v]", string(buf), err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = s.checkGenAuthCodeReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.genAuthCode(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("genAuthCode failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusOK
	c.JSON(status, res)
}

func (s *Service) checkGenAuthCodeReqBody(req *GenAuthCodeReqBody) error {
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

func (s *Service) genAuthCode(req *GenAuthCodeReqBody) (*GenAuthCodeResBody, error) {
	code, err := s.cache.GetAuthCode(req.Email)
	if err != nil {
		return nil, err
	}
	if code == "" {
		code = NewCode()
	}
	if req.Type == "email" {
		if err := s.cache.SetAuthCode(req.Email, code); err != nil {
			return nil, err
		}
		if err := s.mc.Send(req.Email, "hpifu 账号验证", mail.NewAuthCodeTpl(req.FirstName, req.LastName, code)); err != nil {
			return nil, err
		}
	}

	if req.Type == "phone" {
		if err := s.cache.SetAuthCode(req.Phone, code); err != nil {
			return nil, err
		}
	}

	return &GenAuthCodeResBody{OK: true}, nil
}
