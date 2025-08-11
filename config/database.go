package config

import (
	"errors"
	"fmt"
	"log"
	"time"
	"url-shortener-go/entity"
	"url-shortener-go/helpers"

	"gorm.io/driver/postgres"
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
	dbPort := helpers.Env("DB_PORT", "5432")
	dbPassword := helpers.Env("DB_PASSWORD", "root")

	// dsn := fmt.Sprintf("postgres:%v@tcp(%v:%v)/url_shortener?parseTime=True", dbPassword, dbHost, dbPort)
	dsn := fmt.Sprintf("host=%s user=postgres password=%s dbname=url_shortener port=%s", dbHost, dbPassword, dbPort)
	db, err := gorm.Open(postgres.Open(dsn))

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

	if err := DBConn.Where(entity.User{Email: "superuser@nerdproject.id"}).First(&userEntity).Error; err != nil {
		log.Println("Error DB: ", err.Error())
		return nil, err
	}

	if userEntity.Id == "" {
		return nil, errors.New("Superuser not found")
	}

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
