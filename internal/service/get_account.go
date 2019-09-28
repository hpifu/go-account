package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
)

type GETAccountReq struct {
	Token string `uri:"token" json:"token,omitempty"`
}

type GETAccountRes Account

func (s *Service) GETAccount(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &GETAccountReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := s.validGETAccount(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	account, err := s.cache.GetAccount(req.Token)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("redis get account failed. err: [%v]", err)
	}
	if account == nil {
		return req, nil, http.StatusNoContent, nil
	}

	return req, &GETAccountRes{
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

func (s *Service) validGETAccount(req *GETAccountReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}
