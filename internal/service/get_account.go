package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	godtoken "github.com/hpifu/go-godtoken/api"
	"github.com/hpifu/go-kit/rule"
	"net/http"
	"time"
)

type GETAccountReq struct {
	Token string `uri:"token" json:"token,omitempty"`
}

type GETAccountRes Account

func (s *Service) GETAccount(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &GETAccountReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := s.validGETAccount(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	account, err := s.redis.GetAccount(req.Token)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("redis get account failed. err: [%v]", err)
	}
	if account == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		res, err := s.godtokenCli.Verify(ctx, &godtoken.VerifyReq{Rid: rid, Token: req.Token})
		if err != nil {
			return req, nil, http.StatusInternalServerError, fmt.Errorf("godtoken verify failed. err: [%v]", err)
		}
		if !res.Ok {
			return req, nil, http.StatusUnauthorized, nil
		}
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
	if err := rule.Check([][3]interface{}{
		{"token", req.Token, []rule.Rule{rule.Required}},
	}); err != nil {
		return err
	}

	return nil
}
