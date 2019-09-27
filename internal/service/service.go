package service

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/mail"
	"github.com/hpifu/go-account/internal/mysqldb"
	"github.com/hpifu/go-account/internal/rediscache"
	"github.com/sirupsen/logrus"
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

func Decorator(inner func(*gin.Context) (interface{}, interface{}, int, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		rid := c.DefaultQuery("rid", NewToken())
		req, res, status, err := inner(c)
		if err != nil {
			c.String(status, err.Error())
			WarnLog.WithField("@rid", rid).WithField("err", err).Warn()
		} else if res == nil {
			c.Status(status)
		} else {
			switch res.(type) {
			case string:
				c.String(status, res.(string))
			default:
				c.JSON(status, res)
			}
		}

		AccessLog.WithFields(logrus.Fields{
			"host":   c.Request.Host,
			"url":    c.Request.URL.String(),
			"req":    req,
			"res":    res,
			"rid":    rid,
			"err":    err,
			"status": status,
		}).Info()
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
