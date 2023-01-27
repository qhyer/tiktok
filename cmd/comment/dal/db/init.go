package db

import (
	"tiktok/pkg/constants"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormOpentracing "gorm.io/plugin/opentracing"
)

var DB *gorm.DB

// Init init DB
func Init() {
	var err error
	DB, err = gorm.Open(mysql.Open(constants.MySQLDefaultDSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}

	if err = DB.Use(gormOpentracing.New()); err != nil {
		panic(err)
	}
}
