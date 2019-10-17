package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
	"net/http"
)

type GETAccountTokenReq struct {
	Token string `uri:"token" json:"token,omitempty"`
}

type GETAccountTokenRes Account

func (s *Service) GETAccountToken(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &GETAccountTokenReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := s.validGETAccountToken(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	account, err := s.redis.GetAccount(req.Token)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("redis get account failed. err: [%v]", err)
	}
	if account == nil {
		return req, nil, http.StatusUnauthorized, nil
	}

	return req, &GETAccountTokenRes{
		ID:        account.ID,
		Email:     account.Email,
		Phone:     account.Phone,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Birthday:  account.Birthday,
		Gender:    account.Gender,
		Avatar:    account.Avatar,
	}, http.StatusOK, nil
}

func (s *Service) validGETAccountToken(req *GETAccountTokenReq) error {
	if err := rule.Check([][3]interface{}{
		{"token", req.Token, []rule.Rule{rule.Required}},
	}); err != nil {
		return err
	}

	return nil
}
