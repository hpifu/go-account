package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis"
	"github.com/hpifu/go-account/internal/mail"
	"github.com/hpifu/go-account/internal/mysql"
	"github.com/hpifu/go-account/internal/redis"
	"github.com/hpifu/go-account/internal/service"
	godtoken "github.com/hpifu/go-godtoken/api"
	"github.com/hpifu/go-kit/hhttp"
	"github.com/hpifu/go-kit/logger"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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

	service.InfoLog = infoLog
	service.WarnLog = warnLog
	service.AccessLog = accessLog

	// init mysql
	db, err := mysql.NewMysql(config.GetString("mysql.uri"))
	if err != nil {
		panic(err)
	}
	infoLog.Infof("init mysql success. uri [%v]", config.GetString("mysql.uri"))

	// init redis cache
	option := &goredis.Options{
		Addr:         config.GetString("redis.addr"),
		DialTimeout:  config.GetDuration("redis.dialTimeout"),
		ReadTimeout:  config.GetDuration("redis.readTimeout"),
		WriteTimeout: config.GetDuration("redis.writeTimeout"),
		MaxRetries:   config.GetInt("redis.maxRetries"),
		PoolSize:     config.GetInt("redis.poolSize"),
		Password:     config.GetString("redis.password"),
		DB:           config.GetInt("redis.db"),
	}
	cache, err := redis.NewRedis(
		option,
		config.GetDuration("authCodeExpiration"),
		config.GetDuration("tokenExpiration"),
	)
	if err != nil {
		panic(err)
	}
	infoLog.Infof("init redis success. option [%#v]", option)

	// init godtoken client
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}
	conn, err := grpc.Dial(
		config.GetString("godtoken.address"),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
	)
	if err != nil {
		panic(err)
	}
	godtokenCli := godtoken.NewServiceClient(conn)
	infoLog.Infof("init godtoken client success. address: [%v]", config.GetString("godtoken.address"))

	// init mail client
	mc := &mail.MailClient{}
	if err := config.Sub("mail").Unmarshal(mc); err != nil {
		panic(err)
	}
	infoLog.Infof("init mail client success. mailclient [%#v]", mc)

	secure := config.GetBool("service.cookieSecure")
	domain := config.GetString("service.cookieDomain")
	origins := config.GetStringSlice("service.allowOrigins")
	// init services
	svc := service.NewService(db, cache, mc, godtokenCli, secure, domain)

	// init gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"PUT", "POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
	}))

	// set handler
	d := hhttp.NewGinHttpDecorator(infoLog, warnLog, accessLog)
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	r.POST("/account", d.Decorate(svc.POSTAccount))
	r.GET("/accounts", d.Decorate(svc.GETAccounts))
	r.GET("/account/token/:token", d.Decorate(svc.GETAccountToken))
	r.PUT("/account/token/:token/:field", d.Decorate(svc.PUTAccountToken))
	r.POST("/authcode/:type", d.Decorate(svc.POSTAuthCode))
	r.GET("/verify/account", d.Decorate(svc.VerifyAccount))
	r.GET("/verify/authcode/:type", d.Decorate(svc.VerifyAuthCode))
	r.POST("/signin", d.Decorate(svc.SignIn))
	r.GET("/signout/:token", d.Decorate(svc.SignOut))

	infoLog.Infof("%v init success, port [%v]", os.Args[0], config.GetString("service.port"))

	// run server
	server := &http.Server{
		Addr:    config.GetString("service.port"),
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// graceful quit
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	infoLog.Infof("%v shutdown ...", os.Args[0])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		warnLog.Errorf("%v shutdown fail or timeout", os.Args[0])
		return
	}
	warnLog.Out.(*rotatelogs.RotateLogs).Close()
	accessLog.Out.(*rotatelogs.RotateLogs).Close()
	infoLog.Errorf("%v shutdown success", os.Args[0])
	infoLog.Out.(*rotatelogs.RotateLogs).Close()
}
