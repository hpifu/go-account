package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/rule"
	"github.com/sirupsen/logrus"
)

type VerifyAccountReq struct {
	Field string `json:"field,omitempty" form:"field"`
	Value string `json:"value,omitempty" form:"value"`
}

type VerifyAccountRes string

func (s *Service) VerifyAccount(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &VerifyAccountReq{}
	var res VerifyAccountRes
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

	if err := c.Bind(req); err != nil {
		err = fmt.Errorf("bind json failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = s.checkGETAccountVerifyReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.verifyAccount(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("verifyAccount failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusOK
	c.String(status, string(res))
}

func (s *Service) checkGETAccountVerifyReqBody(req *VerifyAccountReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Field: {rule.Required, rule.In(map[interface{}]struct{}{"phone": {}, "email": {}, "username": {}})},
		req.Value: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) verifyAccount(req *VerifyAccountReq) (VerifyAccountRes, error) {
	if req.Field == "phone" {
		account, err := s.db.SelectAccountByPhone(req.Value)
		if err != nil {
			return "", err
		}
		if account == nil {
			return "", nil
		}
		return "电话号码已存在", nil
	}

	if req.Field == "email" {
		account, err := s.db.SelectAccountByEmail(req.Value)
		if err != nil {
			return "", err
		}
		if account == nil {
			return "", nil
		}
		return "邮箱已存在", nil
	}

	if req.Field == "username" {
		account, err := s.db.SelectAccountByPhoneOrEmail(req.Value)
		if err != nil {
			return "", err
		}
		if account == nil {
			return "账号不存在", nil
		}
		return "", nil
	}

	return VerifyAccountRes(fmt.Sprintf("未知字段 [%v]", req.Field)), nil
}
