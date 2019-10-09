package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
)

type VerifyAccountReq struct {
	Field string `json:"field,omitempty" form:"field"`
	Value string `json:"value,omitempty"  form:"value"`
}

type VerifyAccountRes string

func (s *Service) VerifyAccount(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &VerifyAccountReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validVerifyAccount(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	if req.Field == "phone" {
		account, err := s.db.SelectAccountByPhone(req.Value)
		if err != nil {
			return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select account failed. err: [%v]", err)
		}
		if account == nil {
			return req, nil, http.StatusOK, nil
		}
		return req, "电话号码已存在", http.StatusOK, nil
	}

	if req.Field == "email" {
		account, err := s.db.SelectAccountByEmail(req.Value)
		if err != nil {
			return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select account failed. err: [%v]", err)
		}
		if account == nil {
			return req, nil, http.StatusOK, nil
		}
		return req, "邮箱已存在", http.StatusOK, nil
	}

	if req.Field == "username" {
		account, err := s.db.SelectAccountByPhoneOrEmail(req.Value)
		if err != nil {
			return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select account failed. err: [%v]", err)
		}
		if account != nil {
			return req, nil, http.StatusOK, nil
		}
		return req, "账号不存在", http.StatusOK, nil
	}

	return req, VerifyAccountRes(fmt.Sprintf("未知字段 [%v]", req.Field)), http.StatusOK, nil
}

func (s *Service) validVerifyAccount(req *VerifyAccountReq) error {
	if err := rule.Check([][3]interface{}{
		{"field", req.Field, []rule.Rule{rule.Required, rule.In(map[interface{}]struct{}{"phone": {}, "email": {}, "username": {}})}},
		{"value", req.Value, []rule.Rule{rule.Required}},
	}); err != nil {
		return err
	}

	return nil
}
