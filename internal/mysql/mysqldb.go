package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

// Mysql db
type Mysql struct {
	db *gorm.DB
}

// NewMysql create a db
func NewMysql(uri string) (*Mysql, error) {
	db, err := gorm.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	if !db.HasTable(&Account{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Account{}).Error; err != nil {
			panic(err)
		}
	}

	// 服务器主动断开连接，报 "invalid connection" 错误
	// 临时解决方案，设置短一点连接时间，主动重连
	// https://github.com/jinzhu/gorm/issues/1822
	db.DB().SetConnMaxLifetime(60 * time.Second)

	return &Mysql{
		db: db,
	}, nil
}
