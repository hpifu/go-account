package account

import "github.com/hpifu/go-account/internal/c"

type SignUpReq struct {
	FirstName string   `json:"firstName,omitempty"`
	LastName  string   `json:"lastName,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	Email     string   `json:"email,omitempty"`
	Password  string   `json:"password,omitempty"`
	Birthday  string   `json:"birthday,omitempty"`
	Gender    c.Gender `json:"gender,omitempty"`
}

type SignUpRes struct {
	Success bool `json:"success,omitempty"`
}
