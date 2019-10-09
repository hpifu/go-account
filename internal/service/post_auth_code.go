package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/mail"
	"github.com/hpifu/go-kit/rule"
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

func (s *Service) POSTAuthCode(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &POSTAuthCodeReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validPOSTAuthCode(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	code, err := s.cache.GetAuthCode(req.Email)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("redis get auth code failed. err: [%v]", err)
	}
	if code == "" {
		code = NewCode()
	}

	if req.Type == "email" {
		if err := s.cache.SetAuthCode(req.Email, code); err != nil {
			return req, nil, http.StatusInternalServerError, fmt.Errorf("redis set auth code failed. err: [%v]", err)
		}
		if err := s.mc.Send(req.Email, "hpifu 账号验证", mail.NewAuthCodeTpl(req.FirstName, req.LastName, code)); err != nil {
			return req, nil, http.StatusInternalServerError, fmt.Errorf("mail send auth code failed. err: [%v]", err)
		}
	}
	if req.Type == "phone" {
		if err := s.cache.SetAuthCode(req.Phone, code); err != nil {
			return req, nil, http.StatusInternalServerError, fmt.Errorf("redis set auth code failed. err: [%v]", err)
		}
		// todo send auth code to phone
	}

	return req, nil, http.StatusCreated, nil
}

func (s *Service) validPOSTAuthCode(req *POSTAuthCodeReq) error {
	if err := rule.Check([][3]interface{}{
		{"type", req.Type, []rule.Rule{rule.Required, rule.In("email", "phone")}},
		{"firstName", req.FirstName, []rule.Rule{rule.Required, rule.AtMost32Characters}},
		{"lastName", req.LastName, []rule.Rule{rule.Required, rule.AtMost32Characters}},
	}); err != nil {
		return err
	}

	switch req.Type {
	case "phone":
		return rule.Check([][3]interface{}{{"phone", req.Phone, []rule.Rule{rule.Required, rule.ValidPhone}}})
	case "email":
		return rule.Check([][3]interface{}{{"email", req.Email, []rule.Rule{rule.Required, rule.ValidEmail, rule.AtMost64Characters}}})
	}

	return nil
}
