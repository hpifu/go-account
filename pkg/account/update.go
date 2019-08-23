package account

import "github.com/hpifu/go-account/internal/c"

type UpdateReq struct {
	Token       string   `json:"token"`
	Field       string   `json:"field,omitempty"`
	Email       string   `json:"email,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	FirstName   string   `json:"firstName,omitempty"`
	LastName    string   `json:"lastName,omitempty"`
	Birthday    string   `json:"birthday,omitempty"`
	Password    string   `json:"password,omitempty"`
	OldPassword string   `json:"oldPassword,omitempty"`
	Gender      c.Gender `json:"gender,omitempty"`
}

type UpdateRes struct {
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}
