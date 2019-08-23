package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/rule"
	api "github.com/hpifu/go-account/pkg/account"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Service) GetAccount(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	parm := &api.GetAccountReqParm{
		Token: c.DefaultQuery("token", ""),
	}
	req := &api.GetAccountReqBody{}
	var res *api.GetAccountResBody
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

	if err = s.checkGetAccountReqBody(parm, req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.getAccount(parm, req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("getAccount failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusOK
	c.JSON(status, res)
}

func (s *Service) checkGetAccountReqBody(parm *api.GetAccountReqParm, req *api.GetAccountReqBody) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		parm.Token: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) getAccount(parm *api.GetAccountReqParm, req *api.GetAccountReqBody) (*api.GetAccountResBody, error) {
	account, err := s.cache.GetAccount(parm.Token)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return &api.GetAccountResBody{OK: false}, nil
	}

	return &api.GetAccountResBody{
		OK: true,
		Account: &api.Account{
			ID:        account.ID,
			Email:     account.Email,
			Phone:     account.Phone,
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Birthday:  account.Birthday,
			Password:  account.Password,
			Gender:    account.Gender,
		},
	}, nil
}
