package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/rule"
	"github.com/sirupsen/logrus"
	"net/http"
)

type POSTAuthCodeVertifyReq struct {
	Type  string `json:"type,omitempty" uri:"type"`
	Phone string `json:"phone,omitempty" form:"phone"`
	Email string `json:"email,omitempty" form:"email"`
	Code  string `json:"code,omitempty" form:"code"`
}

type POSTAuthCodeVertifyRes string

func (s *Service) POSTAuthCodeVerify(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &POSTAuthCodeVertifyReq{}
	var res POSTAuthCodeVertifyRes
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

	if err = s.validPOSTAuthCodeVerify(req); err != nil {
		err = fmt.Errorf("check request body failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.postAuthCodeVerify(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("postAuthCodeVerify failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}
	if res != "" {
		status = http.StatusForbidden
		c.String(status, string(res))
		return
	}

	status = http.StatusOK
	c.Status(status)
}

func (s *Service) validPOSTAuthCodeVerify(req *POSTAuthCodeVertifyReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Type: {rule.Required, rule.In(map[interface{}]struct{}{"phone": {}, "email": {}})},
		req.Code: {rule.Required, rule.ValidCode},
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

func (s *Service) postAuthCodeVerify(req *POSTAuthCodeVertifyReq) (POSTAuthCodeVertifyRes, error) {
	if req.Type == "phone" {
		code, err := s.cache.GetAuthCode(req.Phone)
		if err != nil {
			return "", err
		}
		if code == "" {
			return "验证码不存在", nil
		}
		if code != req.Code {
			return "验证失败", nil
		}
		return "", nil
	}

	if req.Type == "email" {
		code, err := s.cache.GetAuthCode(req.Email)
		if err != nil {
			return "", err
		}
		if code == "" {
			return "验证码不存在", nil
		}
		if code != req.Code {
			return "验证失败", nil
		}
		return "", nil
	}

	return POSTAuthCodeVertifyRes(fmt.Sprintf("未知字段 [%v]", req.Type)), nil
}
