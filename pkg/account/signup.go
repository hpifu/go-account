package account

import "github.com/hpifu/go-account/internal/c"

type SignUpReq struct {
	FirstName string   `json:"firstName,omitempty" http:"body"`
	LastName  string   `json:"lastName,omitempty" http:"body"`
	Phone     string   `json:"phone,omitempty" http:"body"`
	Email     string   `json:"email,omitempty" http:"body"`
	Password  string   `json:"password,omitempty" http:"body"`
	Birthday  string   `json:"birthday,omitempty" http:"body"`
	Gender    c.Gender `json:"gender,omitempty" http:"body"`
}

type SignUpRes struct {
	Success bool `json:"success,omitempty"`
}
