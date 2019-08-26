package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/c"
	"github.com/hpifu/go-account/internal/rule"
	api "github.com/hpifu/go-account/pkg/account"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (s *Service) Update(c *gin.Context) {
	rid := c.DefaultQuery("rid", NewToken())
	req := &api.UpdateReq{}
	var res *api.UpdateRes
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

	if err := c.BindJSON(req); err != nil {
		err = fmt.Errorf("bind json failed. err: [%v]", err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	if err = s.checkUpdateReqBody(req); err != nil {
		err = fmt.Errorf("check request body failed. body: [%v], err: [%v]", string(buf), err)
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		status = http.StatusBadRequest
		c.String(status, err.Error())
		return
	}

	res, err = s.update(req)
	if err != nil {
		WarnLog.WithField("@rid", rid).WithField("err", err).Warn("update failed")
		status = http.StatusInternalServerError
		c.String(status, err.Error())
		return
	}

	status = http.StatusOK
	c.JSON(status, res)
}

func (s *Service) checkUpdateReqBody(req *api.UpdateReq) error {
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

	return nil
}

func (s *Service) update(req *api.UpdateReq) (*api.UpdateRes, error) {
	account, err := s.cache.GetAccount(req.Token)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return &api.UpdateRes{OK: false, Err: "会话已过期，请重新登录"}, nil
	}

	var ok bool
	switch req.Field {
	case "phone":
		ok, err = s.db.UpdateAccountPhone(account.ID, req.Phone)
		account.Phone = req.Phone
	case "email":
		ok, err = s.db.UpdateAccountEmail(account.ID, req.Email)
		account.Email = req.Email
	case "password":
		if req.OldPassword != account.Password {
			return &api.UpdateRes{OK: false, Err: fmt.Sprintf("密码错误")}, nil
		}
		ok, err = s.db.UpdateAccountPassword(account.ID, req.Password)
	case "gender":
		ok, err = s.db.UpdateAccountGender(account.ID, req.Gender)
		account.Gender = req.Gender
	case "name":
		ok, err = s.db.UpdateAccountName(account.ID, req.FirstName, req.LastName)
		account.FirstName = req.FirstName
		account.LastName = req.LastName
	case "birthday":
		birthday, _ := time.Parse("2006-01-02", req.Birthday)
		ok, err = s.db.UpdateAccountBirthday(account.ID, birthday)
		account.Birthday = req.Birthday
	case "avatar":
		ok, err = s.db.UpdateAccountAvatar(account.ID, req.Avatar)
		account.Avatar = req.Avatar
	default:
		return &api.UpdateRes{OK: false, Err: fmt.Sprintf("未知字段 [%v]", req.Field)}, nil
	}

	if err != nil {
		return nil, err
	}

	if err := s.cache.SetAccount(req.Token, account); err != nil {
		return nil, err
	}

	return &api.UpdateRes{OK: ok}, nil
}
