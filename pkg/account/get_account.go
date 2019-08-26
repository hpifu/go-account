package account

import (
	"github.com/hpifu/go-account/internal/c"
)

type Account struct {
	ID        int      `json:"id,omitempty"`
	Email     string   `json:"email,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	FirstName string   `json:"firstName,omitempty"`
	LastName  string   `json:"lastName,omitempty"`
	Birthday  string   `json:"birthday,omitempty"`
	Password  string   `json:"password,omitempty"`
	Gender    c.Gender `json:"gender"`
}

type GetAccountReq struct {
	Token string `json:"token,omitempty" http:"param"`
}

type GetAccountRes struct {
	OK      bool     `json:"ok"`
	Account *Account `json:"account"`
}
