package account

import (
	"fmt"
	"net/http"
)

func (c *Client) GETAccounts(rid string, token string, ids []int) ([]*Account, error) {
	result := c.client.GET("http://"+c.address+"/accounts", map[string]string{
		"Authorization": token,
	}, map[string]interface{}{
		"rid": rid,
		"ids": ids,
	}, nil)
	if result.Err != nil {
		return nil, result.Err
	}

	if result.Status == http.StatusOK {
		var res []*Account
		if err := result.Interface(res); err != nil {
			return nil, err
		}

		return res, nil
	}

	if result.Status == http.StatusNoContent || result.Status == http.StatusUnauthorized {
		return nil, nil
	}

	return nil, fmt.Errorf("GET accounts failed. res [%v]", string(result.Res))
}
