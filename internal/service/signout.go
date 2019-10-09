package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
)

type SignOutReq struct {
	Token string `json:"token,omitempty" uri:"token"`
}

type SignOutRes struct{}

func (s *Service) SignOut(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &SignOutReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := s.validSignOut(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	if err := s.cache.DelAccount(req.Token); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("redis del account failed. err: [%v]", err)
	}

	return req, nil, http.StatusAccepted, nil
}

func (s *Service) validSignOut(req *SignOutReq) error {
	if err := rule.Check([][3]interface{}{{"token", req.Token, []rule.Rule{rule.Required}}}); err != nil {
		return err
	}

	return nil
}
