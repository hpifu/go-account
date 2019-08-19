package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/account/internal/rule"
	"github.com/sirupsen/logrus"
	"net/http"
)

type VerifyReqBody struct {
	Field string `json:"field,omitempty"`
	Value string `json:"value,omitempty"`
}

type VerifyResBody struct {
	OK  bool   `json:"ok"`
	Tip string `json:"tip"`
}

func (s *Service) Verify(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &VerifyReqBody{
		Field: c.DefaultQuery("field", ""),
		Value: c.DefaultQuery("value", ""),
	}
	var res *VerifyResBody
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

	if err = s.checkVerifyReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.verify(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("verify failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusOK
	c.JSON(status, res)
}

func (s *Service) checkVerifyReqBody(req *VerifyReqBody) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Field: {rule.Required, rule.In(map[interface{}]struct{}{"phone": {}, "email": {}, "username": {}})},
		req.Value: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) verify(req *VerifyReqBody) (*VerifyResBody, error) {
	if req.Field == "phone" {
		account, err := s.db.SelectAccountByPhone(req.Value)
		if err != nil {
			return nil, err
		}
		if account == nil {
			return &VerifyResBody{OK: true}, nil
		}
		return &VerifyResBody{OK: false, Tip: "电话号码已存在"}, nil
	}

	if req.Field == "email" {
		account, err := s.db.SelectAccountByEmail(req.Value)
		if err != nil {
			return nil, err
		}
		if account == nil {
			return &VerifyResBody{OK: true}, nil
		}
		return &VerifyResBody{OK: false, Tip: "邮箱已存在"}, nil
	}

	if req.Field == "username" {
		account, err := s.db.SelectAccountByPhoneOrEmail(req.Value)
		if err != nil {
			return nil, err
		}
		if account == nil {
			return &VerifyResBody{OK: false, Tip: "账号不存在"}, nil
		}
		return &VerifyResBody{OK: true}, nil
	}

	return &VerifyResBody{OK: false, Tip: fmt.Sprintf("未知字段 [%v]", req.Field)}, nil
}
