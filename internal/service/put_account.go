package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/c"
	"github.com/hpifu/go-account/internal/rule"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type PUTAccountReq struct {
	Token       string   `json:"token" uri:"token"`
	Field       string   `json:"field,omitempty" uri:"field"`
	Email       string   `json:"email,omitempty" form:"email"`
	Phone       string   `json:"phone,omitempty" form:"phone"`
	FirstName   string   `json:"firstName,omitempty" form:"firstName"`
	LastName    string   `json:"lastName,omitempty" form:"lastName"`
	Birthday    string   `json:"birthday,omitempty" form:"birthday"`
	Password    string   `json:"password,omitempty" form:"password"`
	OldPassword string   `json:"oldPassword,omitempty" form:"oldPassword"`
	Gender      c.Gender `json:"gender,omitempty" form:"gender"`
	Avatar      string   `json:"avatar,omitempty" form:"avatar"`
}

func (s *Service) PUTAccount(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &PUTAccountReq{}
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

	if err := c.BindUri(req); err != nil {
		err = fmt.Errorf("bind failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err := c.Bind(req); err != nil {
		err = fmt.Errorf("bind json failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = s.checkPUTAccountReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	var err1 error
	err1, err = s.putAccount(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("putAccount failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}
	if err1 != nil {
		status = http.StatusForbidden
		c.String(status, err1.Error())
		return
	}

	status = http.StatusAccepted
	c.Status(status)
}

func (s *Service) checkPUTAccountReqBody(req *PUTAccountReq) error {
	if err := rule.Check(map[interface{}][]rule.Rule{
		req.Token: {rule.Required},
		req.Field: {rule.Required, rule.In(map[interface{}]struct{}{
			"phone": {}, "email": {}, "name": {}, "birthday": {}, "gender": {}, "password": {}, "avatar": {},
		})},
	}); err != nil {
		return err
	}

	switch req.Field {
	case "phone":
		return rule.Check(map[interface{}][]rule.Rule{
			req.Phone: {rule.Required, rule.ValidPhone},
		})
	case "email":
		return rule.Check(map[interface{}][]rule.Rule{
			req.Email: {rule.Required, rule.ValidEmail, rule.AtMost64Characters},
		})
	case "name":
		return rule.Check(map[interface{}][]rule.Rule{
			req.FirstName: {rule.Required, rule.AtMost32Characters},
			req.LastName:  {rule.Required, rule.AtMost32Characters},
		})
	case "birthday":
		return rule.Check(map[interface{}][]rule.Rule{
			req.Birthday: {rule.Required, rule.ValidBirthday},
		})
	case "gender":
		return rule.Check(map[interface{}][]rule.Rule{
			req.Gender: {rule.In(map[interface{}]struct{}{
				c.GenderUnknown: {}, c.Male: {}, c.Famale: {},
			})},
		})
	case "password":
		return rule.Check(map[interface{}][]rule.Rule{
			req.Password:    {rule.Required, rule.AtLeast8Characters},
			req.OldPassword: {rule.Required, rule.AtLeast8Characters},
		})
	case "avatar":
		return rule.Check(map[interface{}][]rule.Rule{
			req.Avatar: {rule.Required},
		})
	default:
		return fmt.Errorf("未知字段 [%v]", req.Field)
	}
}

func (s *Service) putAccount(req *PUTAccountReq) (error, error) {
	account, err := s.cache.GetAccount(req.Token)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return fmt.Errorf("会话已过期，请重新登录"), nil
	}

	switch req.Field {
	case "phone":
		_, err = s.db.UpdateAccountPhone(account.ID, req.Phone)
		account.Phone = req.Phone
	case "email":
		_, err = s.db.UpdateAccountEmail(account.ID, req.Email)
		account.Email = req.Email
	case "password":
		if req.OldPassword != account.Password {
			return fmt.Errorf("密码错误"), nil
		}
		_, err = s.db.UpdateAccountPassword(account.ID, req.Password)
		account.Password = req.Password
	case "gender":
		_, err = s.db.UpdateAccountGender(account.ID, req.Gender)
		account.Gender = req.Gender
	case "name":
		_, err = s.db.UpdateAccountName(account.ID, req.FirstName, req.LastName)
		account.FirstName = req.FirstName
		account.LastName = req.LastName
	case "birthday":
		birthday, _ := time.Parse("2006-01-02", req.Birthday)
		_, err = s.db.UpdateAccountBirthday(account.ID, birthday)
		account.Birthday = req.Birthday
	case "avatar":
		_, err = s.db.UpdateAccountAvatar(account.ID, req.Avatar)
		account.Avatar = req.Avatar
	default:
		return fmt.Errorf("未知字段 [%v]", req.Field), nil
	}

	if err != nil {
		return nil, err
	}

	if err := s.cache.SetAccount(req.Token, account); err != nil {
		return nil, err
	}

	return nil, nil
}
