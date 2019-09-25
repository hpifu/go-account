package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/c"
	"github.com/hpifu/go-account/internal/rule"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Account struct {
	ID        int      `form:"id" json:"id,omitempty"`
	Email     string   `form:"email" json:"email,omitempty"`
	Phone     string   `form:"phone" json:"phone,omitempty"`
	FirstName string   `form:"firstName" json:"firstName,omitempty"`
	LastName  string   `form:"lastName" json:"lastName,omitempty"`
	Birthday  string   `form:"birthday" json:"birthday,omitempty"`
	Password  string   `form:"password" json:"password,omitempty"`
	Gender    c.Gender `form:"gender" json:"gender"`
	Avatar    string   `form:"avatar" json:"avatar"`
}

type GETAccountV2Req struct {
	Token string `uri:"token" json:"token,omitempty"`
}

func (s *Service) GETAccountV2(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &GETAccountV2Req{}
	var res *Account
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

	if err = s.checkGETAccountV2ReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.getAccountV2(req)
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

func (s *Service) checkGETAccountV2ReqBody(req *GETAccountV2Req) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) getAccountV2(req *GETAccountV2Req) (*Account, error) {
	account, err := s.cache.GetAccount(req.Token)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, nil
	}

	return &Account{
		ID:        account.ID,
		Email:     account.Email,
		Phone:     account.Phone,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Birthday:  account.Birthday,
		Password:  account.Password,
		Gender:    account.Gender,
		Avatar:    account.Avatar,
	}, nil
}
