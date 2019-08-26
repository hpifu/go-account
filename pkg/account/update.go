package account

import "github.com/hpifu/go-account/internal/c"

type UpdateReq struct {
	Token       string   `json:"token" http:"body"`
	Field       string   `json:"field,omitempty" http:"body"`
	Email       string   `json:"email,omitempty" http:"body"`
	Phone       string   `json:"phone,omitempty" http:"body"`
	FirstName   string   `json:"firstName,omitempty" http:"body"`
	LastName    string   `json:"lastName,omitempty" http:"body"`
	Birthday    string   `json:"birthday,omitempty" http:"body"`
	Password    string   `json:"password,omitempty" http:"body"`
	OldPassword string   `json:"oldPassword,omitempty" http:"body"`
	Gender      c.Gender `json:"gender,omitempty" http:"body"`
	Avatar      string   `json:"avatar,omitempty" http:"body"`
}

type UpdateRes struct {
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}
