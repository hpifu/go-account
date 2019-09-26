package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/rule"
	"github.com/sirupsen/logrus"
	"net/http"
)

type GETAccountReq struct {
	Token string `uri:"token" json:"token,omitempty"`
}

type GETAccountRes Account

func (s *Service) GETAccount(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &GETAccountReq{}
	var res *GETAccountRes
	var err error
	status := http.StatusOK

	defer func() {
		AccessLog.WithFields(logrus.Fields{
			"host":   c.Request.Host,
			"url":    c.Request.URL.String(),
			"req":    req,
			"res":    res,
			"rid":    rid,
			"err":    err,
			"status": status,
		}).Info()
	}()

	if err := c.BindUri(req); err != nil {
		err = fmt.Errorf("bind failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = s.validGETAccountReq(req); err != nil {
		err = fmt.Errorf("check request body failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.getAccount(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("getAccount failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	if res == nil {
		status = http.StatusNoContent
		c.Status(status)
		return
	}

	status = http.StatusOK
	c.JSON(status, res)
}

func (s *Service) validGETAccountReq(req *GETAccountReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) getAccount(req *GETAccountReq) (*GETAccountRes, error) {
	account, err := s.cache.GetAccount(req.Token)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, nil
	}

	return &GETAccountRes{
		ID:        account.ID,
		Email:     account.Email,
		Phone:     account.Phone,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Birthday:  account.Birthday,
		Gender:    account.Gender,
		Avatar:    account.Avatar,
	}, nil
}
