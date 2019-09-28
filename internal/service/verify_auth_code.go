package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
	"net/http"
)

type VerifyAuthCodeReq struct {
	Type  string `json:"type,omitempty" uri:"type"`
	Phone string `json:"phone,omitempty" form:"phone"`
	Email string `json:"email,omitempty" form:"email"`
	Code  string `json:"code,omitempty" form:"code"`
}

type VerifyAuthCodeRes string

func (s *Service) VerifyAuthCode(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &VerifyAuthCodeReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validVerifyAuthCode(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	if req.Type == "phone" {
		code, err := s.cache.GetAuthCode(req.Phone)
		if err != nil {
			return req, nil, http.StatusInternalServerError, fmt.Errorf("redis get auth code failed. err: [%v]", err)
		}
		if code == "" {
			return req, "验证码不存在", http.StatusOK, nil
		}
		if code != req.Code {
			return req, "验证失败", http.StatusOK, nil
		}
		return req, nil, http.StatusOK, nil
	}

	if req.Type == "email" {
		code, err := s.cache.GetAuthCode(req.Email)
		if err != nil {
			return req, nil, http.StatusInternalServerError, fmt.Errorf("redis get auth code failed. err: [%v]", err)
		}
		if code == "" {
			return req, "验证码不存在", http.StatusOK, nil
		}
		if code != req.Code {
			return req, "验证失败", http.StatusOK, nil
		}
		return req, nil, http.StatusOK, nil
	}

	return req, VerifyAccountRes(fmt.Sprintf("未知字段 [%v]", req.Type)), http.StatusOK, nil
}

func (s *Service) validVerifyAuthCode(req *VerifyAuthCodeReq) error {
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
