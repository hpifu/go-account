package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/c"
	"github.com/hpifu/go-account/internal/mysqldb"
	"github.com/hpifu/go-account/internal/rule"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (s *Service) POSTAccount(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &Account{}
	var err error
	status := http.StatusOK

	defer func() {
		AccessLog.WithFields(logrus.Fields{
			"host":   c.Request.Host,
			"url":    c.Request.URL.String(),
			"req":    req,
			"res":    nil,
			"rid":    rid,
			"err":    err,
			"status": status,
		}).Info()
	}()

	if err := c.Bind(req); err != nil {
		err = fmt.Errorf("bind json failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = s.checkPOSTAccountReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. req: [%v], err: [%v]", req, err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	ok, err := s.postAccount(req)
	if err != nil {
		err = fmt.Errorf("postAccount failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	if !ok {
		status = http.StatusNotModified
	} else {
		status = http.StatusCreated
	}
	c.Status(status)
}

func (s *Service) checkPOSTAccountReqBody(req *Account) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Password: {rule.Required, rule.AtLeast8Characters},
		req.Gender: {rule.In(map[interface{}]struct{}{
			c.GenderUnknown: {}, c.Male: {}, c.Famale: {},
		})},
	}); err != nil {
		return err
	}

	if req.FirstName != "" {
		if err := rule.AtMost32Characters(req.FirstName); err != nil {
			return err
		}
	}
	if req.LastName != "" {
		if err := rule.AtMost32Characters(req.LastName); err != nil {
			return err
		}
	}
	if req.Birthday != "" {
		if err := rule.ValidBirthday(req.Birthday); err != nil {
			return err
		}
	} else {
		req.Birthday = "1970-01-02"
	}

	if req.Phone == "" && req.Email == "" {
		return fmt.Errorf("电话和邮箱不可同时为空")
	}
	if req.Phone != "" {
		if err := rule.Check(map[interface{}][]rule.Rule{
			req.Phone: {rule.Required, rule.ValidPhone},
		}); err != nil {
			return err
		}
	}
	if req.Email != "" {
		if err := rule.Check(map[interface{}][]rule.Rule{
			req.Email: {rule.Required, rule.ValidEmail, rule.AtMost64Characters},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) postAccount(req *Account) (bool, error) {
	birthday, _ := time.Parse("2006-01-02", req.Birthday)
	ok, err := s.db.InsertAccount(&mysqldb.Account{
		Phone:     req.Phone,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Birthday:  birthday,
		Gender:    req.Gender,
	})

	return ok, err
}
