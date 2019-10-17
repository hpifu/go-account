package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	godtoken "github.com/hpifu/go-godtoken/api"
	"github.com/hpifu/go-kit/rule"
)

type GETAccountsReq struct {
	Token string `json:"token,omitempty"`
	IDs   []int  `form:"ids" json:"ids,omitempty"`
}

type GETAccountsRes []*Account

func (s *Service) GETAccounts(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &GETAccountsReq{
		Token: c.GetHeader("Authorization"),
	}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := s.validGETAccounts(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	res, err := s.godtokenCli.Verify(ctx, &godtoken.VerifyReq{Rid: rid, Token: req.Token})
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("godtoken verify failed. err: [%v]", err)
	}
	if !res.Ok {
		return req, nil, http.StatusUnauthorized, nil
	}

	accounts, err := s.db.SelectAccountByIDs(req.IDs)
	if err != nil {
		return req, nil, http.StatusInternalServerError, err
	}
	var as GETAccountsRes
	for _, account := range accounts {
		as = append(as, &Account{
			ID:        account.ID,
			Email:     account.Email,
			Phone:     account.Phone,
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Birthday:  account.Birthday.Format("2006-01-02"),
			Gender:    account.Gender,
			Avatar:    account.Avatar,
		})
	}

	return req, as, http.StatusOK, nil
}

func (s *Service) validGETAccounts(req *GETAccountsReq) error {
	if err := rule.Check([][3]interface{}{
		{"token", req.Token, []rule.Rule{rule.Required}},
	}); err != nil {
		return err
	}

	return nil
}
