package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/account/internal/rediscache"
	"github.com/hpifu/account/internal/rule"
	"github.com/sirupsen/logrus"
	"net/http"
)

type GetAccountReqBody struct {
	Token string `json:"token,omitempty"`
}

type GetAccountResBody struct {
	OK      bool                `json:"ok"`
	Account *rediscache.Account `json:"account"`
}

func (s *Service) GetAccount(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &GetAccountReqBody{
		Token: c.DefaultQuery("token", ""),
	}
	var res *GetAccountResBody
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

	if err = s.checkGetAccountReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
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

	status = http.StatusOK
	c.JSON(status, res)
}

func (s *Service) checkGetAccountReqBody(req *GetAccountReqBody) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) getAccount(req *GetAccountReqBody) (*GetAccountResBody, error) {
	account, err := s.cache.GetAccount(req.Token)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return &GetAccountResBody{OK: false}, nil
	}

	return &GetAccountResBody{
		OK:      true,
		Account: account,
	}, nil
}
