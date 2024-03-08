package config

import (
	"fmt"
	"log"
	"time"
	"url-shortener-go/entity"
	"url-shortener-go/helpers"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-faker/faker/v4"
)

var DBConn *gorm.DB

type UrlSeeder struct {
	ShortUrl string `faker:"username,len=8"`
	LongUrl  string `faker:"url"`
}

func InitDB() (*gorm.DB, error) {

	dbHost := helpers.Env("DB_HOST", "host.docker.internal")
	dbPort := helpers.Env("DB_PORT", "3306")
	dbPassword := helpers.Env("DB_PASSWORD", "root")

	dsn := fmt.Sprintf("root:%v@tcp(%v:%v)/url_shortener?parseTime=True", dbPassword, dbHost, dbPort)
	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		return nil, err
	}

	DBConn = db

	RunUrlSeeder(100)

	return db, nil
}

func RunUrlSeeder(count int) ([]entity.Redirection, error) {

	urlEntities := []entity.Redirection{}
	userEntity := entity.User{}

	if err := DBConn.First(&userEntity).Error; err != nil {
		log.Println("Error DB: ", err.Error())
		return nil, err
	}

	fmt.Println(&userEntity)

	for i := 0; i < count; i++ {
		urlSeeder := UrlSeeder{}

		if err := faker.FakeData(&urlSeeder); err != nil {
			return urlEntities, err
		}

		urlEntity := entity.Redirection{
			ShortUrl:  urlSeeder.ShortUrl,
			LongUrl:   urlSeeder.LongUrl,
			UserId:    userEntity.Id,
			Status:    "active",
			HitCount:  0,
			CreatedAt: time.Now(),
		}

		urlEntities = append(urlEntities, urlEntity)

	}

	return urlEntities, nil
}
