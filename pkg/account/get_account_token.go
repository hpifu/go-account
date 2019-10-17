package account

import (
	"fmt"
	"net/http"

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

func (c *Client) GETAccountToken(rid string, token string) (*Account, error) {
	result := c.client.GET("http://"+c.address+"/account/token/"+token, nil, map[string]interface{}{
		"rid": rid,
	}, nil)
	if result.Err != nil {
		return nil, result.Err
	}

	if result.Status == http.StatusOK {
		res := &Account{}
		if err := result.Interface(res); err != nil {
			return nil, err
		}

		return res, nil
	}

	if result.Status == http.StatusNoContent || result.Status == http.StatusUnauthorized {
		return nil, nil
	}

	return nil, fmt.Errorf("GETAccountToken failed. res [%v]", string(result.Res))
}
