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
	option := &rediscache.Option{}
	if err := config.Sub("rediscache").Unmarshal(option); err != nil {
		panic(err)
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

	// init services
	service := account.NewService(db, cache, mc)

	// init gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://account.hatlonely.com")
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

	infoLog.Infof("%v init success, port [%v]", os.Args[0], config.GetString("service.port"))

	// run server
	if err := r.Run(config.GetString("service.port")); err != nil {
		panic(err)
	}
}
