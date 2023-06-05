package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

var GlobalDB *gorm.DB

var (
	db  *gorm.DB
	err error
)

func GetGlobalDSN(conf DatabaseConf) (*gorm.DB, error) {
	mysqlPrefix := "mysql://"
	if strings.HasPrefix(conf.DSN, mysqlPrefix) {
		db, err = gorm.Open(mysql.Open(strings.ReplaceAll(conf.DSN, mysqlPrefix, "")), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error), //不打印慢SQL日志
		})
		if err != nil {
			return nil, err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		if err := sqlDB.Ping(); err != nil {
			return nil, err
		}
		if conf.MaxOpenConn > 0 {
			sqlDB.SetMaxOpenConns(conf.MaxOpenConn)
		}
		if conf.MaxIdleConn > 0 {
			sqlDB.SetMaxIdleConns(conf.MaxIdleConn)
		}
		if conf.ConnMaxLiftTime > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLiftTime) * time.Second)
		}
		return db, nil
	} else {
		return nil, gorm.ErrUnsupportedDriver
	}
}

func GetGlobalDB(conf DatabaseConf) error {
	db, err := GetGlobalDSN(conf)
	if err != nil {
		return err
	}
	GlobalDB = db
	return nil
}
