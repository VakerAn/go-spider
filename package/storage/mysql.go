package storage

import (
	"fmt"
	"go-spider/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

func InitMysqlDB() {
	var err error
	dbs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.ConfData.MysqlDB.User,
		config.ConfData.MysqlDB.Password,
		config.ConfData.MysqlDB.Host,
		config.ConfData.MysqlDB.Port,
		config.ConfData.MysqlDB.Scheme,
	)
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dbs,
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}))

	if err != nil {
		log.Fatalf(" gorm.Open.err: %v", err)
	}
	// Get generic database object sql.DB to use its functions
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf(" gorm.Open.err: %v", err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
}
