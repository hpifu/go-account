package service

import (
	"github.com/hpifu/go-account/internal/mail"
	"github.com/hpifu/go-account/internal/mysql"
	"github.com/hpifu/go-account/internal/redis"
	godtoken "github.com/hpifu/go-godtoken/api"
	"github.com/sirupsen/logrus"
)

var InfoLog *logrus.Logger = logrus.New()
var WarnLog *logrus.Logger = logrus.New()
var AccessLog *logrus.Logger = logrus.New()

type Service struct {
	db          *mysql.Mysql
	cache       *redis.Redis
	mc          *mail.MailClient
	godtokenCli *godtoken.ServiceClient
	secure      bool
	domain      string
}

func NewService(
	db *mysql.Mysql,
	cache *redis.Redis,
	mc *mail.MailClient,
	godtokenCli *godtoken.ServiceClient,
	secure bool,
	domain string,
) *Service {
	return &Service{
		db:          db,
		cache:       cache,
		mc:          mc,
		godtokenCli: godtokenCli,
		secure:      secure,
		domain:      domain,
	}
}
