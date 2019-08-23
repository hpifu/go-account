package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/rule"
	api "github.com/hpifu/go-account/pkg/account"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Service) VerifyAuthCode(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &api.VerifyAuthCodeReq{
		Type:  c.DefaultQuery("type", ""),
		Phone: c.DefaultQuery("phone", ""),
		Email: c.DefaultQuery("email", ""),
		Code:  c.DefaultQuery("code", ""),
	}
	var res *api.VerifyAuthCodeRes
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

	if err = s.checkVerifyAuthCodeReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.verifyAuthCode(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("verifyAuthCode failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusOK
	c.JSON(status, res)
}

func (s *Service) checkVerifyAuthCodeReqBody(req *api.VerifyAuthCodeReq) error {
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

func (s *Service) verifyAuthCode(req *api.VerifyAuthCodeReq) (*api.VerifyAuthCodeRes, error) {
	if req.Type == "phone" {
		code, err := s.cache.GetAuthCode(req.Phone)
		if err != nil {
			return nil, err
		}
		if code == "" {
			return &api.VerifyAuthCodeRes{OK: false, Tip: "验证码不存在"}, nil
		}
		if code != req.Code {
			return &api.VerifyAuthCodeRes{OK: false, Tip: "验证失败"}, nil
		}
		return &api.VerifyAuthCodeRes{OK: true}, nil
	}

	if req.Type == "email" {
		code, err := s.cache.GetAuthCode(req.Email)
		if err != nil {
			return nil, err
		}
		if code == "" {
			return &api.VerifyAuthCodeRes{OK: false, Tip: "验证码不存在"}, nil
		}
		if code != req.Code {
			return &api.VerifyAuthCodeRes{OK: false, Tip: "验证失败"}, nil
		}
		return &api.VerifyAuthCodeRes{OK: true}, nil
	}

	return &api.VerifyAuthCodeRes{OK: false, Tip: fmt.Sprintf("未知字段 [%v]", req.Type)}, nil
}
