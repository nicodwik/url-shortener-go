package config

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitDB() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("root:root@tcp(%v:%v)/url_shortener?parseTime=True", dbHost, dbPort)
	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		return nil, err
	}

	DBConn = db

	return db, nil
}
