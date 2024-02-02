package config

import (
	"fmt"
	"url-shortener-go/helpers"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitDB() (*gorm.DB, error) {

	dbHost := helpers.Env("DB_HOST", "host.docker.internal")
	dbPort := helpers.Env("DB_PORT", "3306")

	dsn := fmt.Sprintf("root:root@tcp(%v:%v)/url_shortener?parseTime=True", dbHost, dbPort)
	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		return nil, err
	}

	DBConn = db

	return db, nil
}
