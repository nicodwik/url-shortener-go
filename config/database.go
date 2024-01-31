package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitDB() (*gorm.DB, error) {
	dsn := "root:root@tcp(host.docker.internal:3306)/url_shortener?parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		return nil, err
	}

	DBConn = db

	return db, nil
}
