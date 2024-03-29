package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/hpifu/go-account/internal/mysql"
	"github.com/hpifu/pb-constant/c"
)

type Redis struct {
	client             *redis.Client
	option             *redis.Options
	authCodeExpiration time.Duration
	tokenExpiration    time.Duration
}

type Account struct {
	ID        int      `json:"id,omitempty"`
	Email     string   `json:"email,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	FirstName string   `json:"firstName,omitempty"`
	LastName  string   `json:"lastName,omitempty"`
	Birthday  string   `json:"birthday,omitempty"`
	Password  string   `json:"password,omitempty"`
	Gender    c.Gender `json:"gender"`
	Avatar    string   `json:"avatar,omitempty"`
}

func NewAccount(account *mysql.Account) *Account {
	return &Account{
		ID:        account.ID,
		Email:     account.Email,
		Phone:     account.Phone,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Password:  account.Password,
		Birthday:  account.Birthday.Format("2006-01-02"),
		Gender:    account.Gender,
		Avatar:    account.Avatar,
	}
}

func NewRedis(option *redis.Options, authCodeExpiration, tokenExpiration time.Duration) (*Redis, error) {
	client := redis.NewClient(option)

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &Redis{
		client:             client,
		option:             option,
		authCodeExpiration: authCodeExpiration,
		tokenExpiration:    tokenExpiration,
	}, nil
}

func (rc *Redis) SetAuthCode(key string, code string) error {
	return rc.client.Set("ac_"+key, code, rc.authCodeExpiration).Err()
}

func (rc *Redis) DelAuthCode(key string) error {
	return rc.client.Del("ac_" + key).Err()
}

func (rc *Redis) GetAuthCode(key string) (string, error) {
	buf, err := rc.client.Get("ac_" + key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (rc *Redis) SetAccount(token string, account *Account) error {
	buf, err := json.Marshal(account)
	if err != nil {
		return err
	}

	return rc.client.Set(token, buf, rc.tokenExpiration).Err()
}

func (rc *Redis) DelAccount(token string) error {
	return rc.client.Del(token).Err()
}

func (rc *Redis) GetAccount(token string) (*Account, error) {
	buf, err := rc.client.Get(token).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	account := &Account{}
	if err := json.Unmarshal([]byte(buf), account); err != nil {
		return nil, err
	}

	return account, nil
}
