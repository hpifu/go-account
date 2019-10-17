package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
	"github.com/hpifu/pb-constant/c"
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

type PUTAccountRes string

func (s *Service) PUTAccount(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &PUTAccountReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validPUTAccount(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	account, err := s.cache.GetAccount(req.Token)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("redis get account failed. err: [%v]", err)
	}
	if account == nil {
		return req, "会话已过期，请重新登录", http.StatusOK, nil
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
			return req, "密码错误", http.StatusOK, nil
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
		return req, PUTAccountRes(fmt.Sprintf("未知字段 [%v]", req.Field)), http.StatusOK, nil
	}

	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql update account failed. err: [%v]", err)
	}

	if err := s.cache.SetAccount(req.Token, account); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("redis set account failed. err: [%v]", err)
	}

	return req, nil, http.StatusAccepted, nil
}

func (s *Service) validPUTAccount(req *PUTAccountReq) error {
	if err := rule.Check([][3]interface{}{
		{"token", req.Token, []rule.Rule{rule.Required}},
		{"field", req.Field, []rule.Rule{rule.Required, rule.In(
			"phone", "email", "name", "birthday", "gender", "password", "avatar",
		)}},
	}); err != nil {
		return err
	}

	switch req.Field {
	case "phone":
		return rule.Check([][3]interface{}{{"phone", req.Phone, []rule.Rule{rule.Required, rule.ValidPhone}}})
	case "email":
		return rule.Check([][3]interface{}{{"email", req.Email, []rule.Rule{rule.Required, rule.ValidEmail, rule.AtMost64Characters}}})
	case "name":
		return rule.Check([][3]interface{}{
			{"firstName", req.FirstName, []rule.Rule{rule.Required, rule.AtMost32Characters}},
			{"lastName", req.LastName, []rule.Rule{rule.Required, rule.AtMost32Characters}},
		})
	case "birthday":
		return rule.Check([][3]interface{}{{"birthday", req.Birthday, []rule.Rule{rule.Required, rule.ValidBirthday}}})
	case "gender":
		return rule.Check([][3]interface{}{
			{"gender", req.Gender, []rule.Rule{rule.In(c.Gender_Null, c.Gender_Unknown, c.Gender_Male, c.Gender_Famale)}},
		})
	case "password":
		return rule.Check([][3]interface{}{
			{"password", req.Password, []rule.Rule{rule.Required, rule.AtLeast8Characters}},
			{"oldPassword", req.OldPassword, []rule.Rule{rule.Required, rule.AtLeast8Characters}},
		})
	case "avatar":
		return rule.Check([][3]interface{}{{"avatar", req.Avatar, []rule.Rule{rule.Required}}})
	default:
		return fmt.Errorf("未知字段 [%v]", req.Field)
	}
}
