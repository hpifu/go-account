package service

import (
	"github.com/hpifu/go-account/internal/mail"
	"github.com/hpifu/go-account/internal/mysql"
	"github.com/hpifu/go-account/internal/redis"
	godtoken "github.com/hpifu/go-godtoken/api"
	"github.com/sirupsen/logrus"
)

type Service struct {
	mysql       *mysql.Mysql
	redis       *redis.Redis
	mc          *mail.MailClient
	godtokenCli godtoken.ServiceClient
	secure      bool
	domain      string
	infoLog     *logrus.Logger
	warnLog     *logrus.Logger
	accessLog   *logrus.Logger
}

func (s *Service) SetLogger(infoLog, warnLog, accessLog *logrus.Logger) {
	s.infoLog = infoLog
	s.warnLog = warnLog
	s.accessLog = accessLog
}

func NewService(
	db *mysql.Mysql,
	cache *redis.Redis,
	mc *mail.MailClient,
	godtokenCli godtoken.ServiceClient,
	secure bool,
	domain string,
) *Service {
	return &Service{
		mysql:       db,
		redis:       cache,
		mc:          mc,
		godtokenCli: godtokenCli,
		secure:      secure,
		domain:      domain,
		infoLog:     logrus.New(),
		warnLog:     logrus.New(),
		accessLog:   logrus.New(),
	}
}
