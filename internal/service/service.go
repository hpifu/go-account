package service

import (
	"github.com/hpifu/go-account/internal/mail"
	"github.com/hpifu/go-account/internal/mysql"
	"github.com/hpifu/go-account/internal/redis"
	"github.com/sirupsen/logrus"
)

var InfoLog *logrus.Logger = logrus.New()
var WarnLog *logrus.Logger = logrus.New()
var AccessLog *logrus.Logger = logrus.New()

type Service struct {
	db     *mysql.Mysql
	cache  *redis.Redis
	mc     *mail.MailClient
	secure bool
	domain string
}

func NewService(
	db *mysql.Mysql,
	cache *redis.Redis,
	mc *mail.MailClient,
	secure bool,
	domain string,
) *Service {
	return &Service{
		db:     db,
		cache:  cache,
		mc:     mc,
		secure: secure,
		domain: domain,
	}
}
