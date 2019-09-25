package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/internal/account"
	"github.com/hpifu/go-account/internal/logger"
	"github.com/hpifu/go-account/internal/mail"
	"github.com/hpifu/go-account/internal/mysqldb"
	"github.com/hpifu/go-account/internal/rediscache"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// AppVersion name
var AppVersion = "unknown"

func main() {
	version := flag.Bool("v", false, "print current version")
	configfile := flag.String("c", "configs/account.json", "config file path")
	flag.Parse()
	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	// load config
	config := viper.New()
	config.SetEnvPrefix("account")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()
	config.SetConfigType("json")
	fp, err := os.Open(*configfile)
	if err != nil {
		panic(err)
	}
	err = config.ReadConfig(fp)
	if err != nil {
		panic(err)
	}

	// init logger
	infoLog, err := logger.NewTextLoggerWithViper(config.Sub("logger.infoLog"))
	if err != nil {
		panic(err)
	}
	warnLog, err := logger.NewTextLoggerWithViper(config.Sub("logger.warnLog"))
	if err != nil {
		panic(err)
	}
	accessLog, err := logger.NewJsonLoggerWithViper(config.Sub("logger.accessLog"))
	if err != nil {
		panic(err)
	}
	account.InfoLog = infoLog
	account.WarnLog = warnLog
	account.AccessLog = accessLog

	// init mysqldb
	db, err := mysqldb.NewMysqlDB(config.GetString("mysqldb.uri"))
	if err != nil {
		panic(err)
	}
	infoLog.Infof("init mysqldb success. uri [%v]", config.GetString("mysqldb.uri"))

	// init redis cache
	option := &rediscache.Option{
		Address:            config.GetString("rediscache.address"),
		Timeout:            config.GetDuration("rediscache.timeout"),
		Retries:            config.GetInt("rediscache.retries"),
		PoolSize:           config.GetInt("rediscache.poolSize"),
		Password:           config.GetString("rediscache.password"),
		DB:                 config.GetInt("rediscache.0"),
		TokenExpiration:    config.GetDuration("rediscache.tokenExpiration"),
		AuthCodeExpiration: config.GetDuration("rediscache.authCodeExpiration"),
	}
	cache, err := rediscache.NewRedisCache(option)
	if err != nil {
		panic(err)
	}
	infoLog.Infof("init redis cache success. option [%#v]", option)

	// init mail client
	mc := &mail.MailClient{}
	if err := config.Sub("mail").Unmarshal(mc); err != nil {
		panic(err)
	}
	infoLog.Infof("init mail client success. mailclient [%#v]", mc)

	secure := config.GetBool("service.cookieSecure")
	domain := config.GetString("service.cookieDomain")
	origin := config.GetString("service.allowOrigin")
	// init services
	service := account.NewService(db, cache, mc, secure, domain)

	// init gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// set handler
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	r.GET("/verify", service.Verify)
	r.GET("/getaccount", service.GetAccount)
	r.GET("/signout", service.SignOut)
	r.GET("/verifyauthcode", service.VerifyAuthCode)
	r.POST("/genauthcode", service.GenAuthCode)
	r.POST("/update", service.Update)
	r.POST("/signin", service.SignIn)
	r.POST("/signup", service.SignUp)
	r.GET("/account/:token", service.GETAccountV2)
	r.PUT("/account", service.POSTAccount)

	infoLog.Infof("%v init success, port [%v]", os.Args[0], config.GetString("service.port"))

	// run server
	if err := r.Run(config.GetString("service.port")); err != nil {
		panic(err)
	}
}
