package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/account/internal/rule"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SignOutReqBody struct {
	Token string `json:"token,omitempty"`
}

type SignOutResBody struct {
	OK bool `json:"ok"`
}

func (s *Service) SignOut(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &SignOutReqBody{
		Token: c.DefaultQuery("token", ""),
	}
	var res *SignOutResBody
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

	if err = s.checkSignOutReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.signOut(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("signOut failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusOK
	c.JSON(status, res)
}

func (s *Service) checkSignOutReqBody(req *SignOutReqBody) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) signOut(req *SignOutReqBody) (*SignOutResBody, error) {
	err := s.cache.DelAccount(req.Token)
	if err != nil {
		return nil, err
	}
	return &SignOutResBody{OK: true}, nil
}
