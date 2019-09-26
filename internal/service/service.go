package service

import (
	"encoding/hex"
	"fmt"
	"github.com/hpifu/go-account/internal/mail"
	"github.com/hpifu/go-account/internal/mysqldb"
	"github.com/hpifu/go-account/internal/rediscache"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

var InfoLog *logrus.Logger
var WarnLog *logrus.Logger
var AccessLog *logrus.Logger

func init() {
	InfoLog = logrus.New()
	WarnLog = logrus.New()
	AccessLog = logrus.New()
}

type Service struct {
	db     *mysqldb.MysqlDB
	cache  *rediscache.RedisCache
	mc     *mail.MailClient
	secure bool
	domain string
}

func NewService(
	db *mysqldb.MysqlDB,
	cache *rediscache.RedisCache,
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

func NewToken() string {
	buf := make([]byte, 32)
	token := make([]byte, 16)
	rand.New(rand.NewSource(time.Now().UnixNano())).Read(token)
	hex.Encode(buf, token)
	return string(buf)
}

func NewCode() string {
	return fmt.Sprintf("%06d", int(rand.NewSource(time.Now().UnixNano()).(rand.Source64).Uint64()%1000000))
}
