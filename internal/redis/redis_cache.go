package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/hpifu/go-account/internal/c"
	"github.com/hpifu/go-account/internal/mysql"
	"time"
)

type Option struct {
	Address            string
	Timeout            time.Duration
	Retries            int
	PoolSize           int
	Password           string
	DB                 int
	TokenExpiration    time.Duration
	AuthCodeExpiration time.Duration
}

type Redis struct {
	client *redis.Client
	option *Option
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

func NewRedis(option *Option) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         option.Address,
		DialTimeout:  option.Timeout,
		ReadTimeout:  option.Timeout,
		WriteTimeout: option.Timeout,
		MaxRetries:   option.Retries,
		PoolSize:     option.PoolSize,
		Password:     option.Password,
		DB:           option.DB,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &Redis{
		client: client,
		option: option,
	}, nil
}

func (rc *Redis) SetAuthCode(key string, code string) error {
	return rc.client.Set("ac_"+key, code, rc.option.AuthCodeExpiration).Err()
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

	return rc.client.Set(token, buf, rc.option.TokenExpiration).Err()
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
