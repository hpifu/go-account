package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/mysql"
	"github.com/hpifu/go-kit/rule"
	"github.com/hpifu/pb-constant/c"
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

type POSTAccountReq Account

type POSTAccountRes struct{}

func (s *Service) POSTAccount(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &POSTAccountReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validPOSTAccount(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	birthday, _ := time.Parse("2006-01-02", req.Birthday)
	ok, err := s.db.InsertAccount(&mysql.Account{
		Phone:     req.Phone,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Birthday:  birthday,
		Gender:    req.Gender,
	})

	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql insert account failed. err: [%v]", err)
	}

	if !ok {
		return req, nil, http.StatusNotModified, nil
	}
	return req, nil, http.StatusCreated, nil
}

func (s *Service) validPOSTAccount(req *POSTAccountReq) error {
	if err := rule.Check([][3]interface{}{
		{"password", req.Password, []rule.Rule{rule.Required, rule.AtLeast8Characters}},
		{"gender", req.Gender, []rule.Rule{rule.In(c.Gender_Null, c.Gender_Unknown, c.Gender_Male, c.Gender_Famale)}},
	}); err != nil {
		return err
	}

	if req.FirstName != "" {
		if err := rule.Check([][3]interface{}{{"firstName", req.FirstName, []rule.Rule{rule.AtMost32Characters}}}); err != nil {
			return err
		}
	}
	if req.LastName != "" {
		if err := rule.Check([][3]interface{}{{"lastName", req.LastName, []rule.Rule{rule.AtMost32Characters}}}); err != nil {
			return err
		}
	}
	if req.Birthday != "" {
		if err := rule.Check([][3]interface{}{{"birthday", req.Birthday, []rule.Rule{rule.ValidBirthday}}}); err != nil {
			return err
		}
	} else {
		req.Birthday = "1970-01-02"
	}

	if req.Phone == "" && req.Email == "" {
		return fmt.Errorf("电话和邮箱不可同时为空")
	}
	if req.Phone != "" {
		if err := rule.Check([][3]interface{}{{"phone", req.Phone, []rule.Rule{rule.ValidPhone}}}); err != nil {
			return err
		}
	}
	if req.Email != "" {
		if err := rule.Check([][3]interface{}{{"email", req.Email, []rule.Rule{rule.ValidEmail, rule.AtMost64Characters}}}); err != nil {
			return err
		}
	}

	return nil
}
