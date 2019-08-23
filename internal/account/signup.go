package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/c"
	"github.com/hpifu/go-account/internal/mysqldb"
	"github.com/hpifu/go-account/internal/rule"
	api "github.com/hpifu/go-account/pkg/account"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (s *Service) SignUp(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &api.SignUpReqBody{}
	var res *api.SignUpResBody
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

	if err := c.BindJSON(req); err != nil {
		err = fmt.Errorf("bind json failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = s.checkSignUpReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. req: [%v], err: [%v]", req, err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.signUp(req)
	if err != nil {
		err = fmt.Errorf("signUp failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusOK
	c.JSON(status, res)
}

func (s *Service) checkSignUpReqBody(req *api.SignUpReqBody) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Password: {rule.Required, rule.AtLeast8Characters},
		req.Birthday: {rule.Required, rule.ValidBirthday},
		req.Gender: {rule.In(map[interface{}]struct{}{
			c.GenderUnknown: {}, c.Male: {}, c.Famale: {},
		})},
		req.FirstName: {rule.Required, rule.AtMost32Characters},
		req.LastName:  {rule.Required, rule.AtMost32Characters},
	}); err != nil {
		return err
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

func (s *Service) signUp(req *api.SignUpReqBody) (*api.SignUpResBody, error) {
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

	return &api.SignUpResBody{Success: ok}, err
}
