package url

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"url-shortener-go/config"
	"url-shortener-go/entity"

	"gorm.io/gorm"
)

// type UrlRepository interface {
// 	GetAll() ([]UrlEntity, error)
// }

func GetAll() ([]entity.Url, error) {
	var urls []entity.Url

	if err := config.DBConn.Find(&urls).Error; err != nil {
		log.Fatalf("Failed to connect DB: %v", err)

		return nil, err
	}

	return urls, nil
}

func InsertNewUrl(shortUrl string, longUrl string) (*entity.Url, error) {

	urlEntity := entity.Url{
		UserId:   1,
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
		Status:   "active",
	}

	if err := config.DBConn.Create(&urlEntity).Error; err != nil {
		return nil, err
	}

	return &urlEntity, nil
}

func FindRedirection(shortUrl string) (*entity.Url, error) {
	url, err := config.CacheRemember(shortUrl, 3600, func() interface{} {
		urlEntity := entity.Url{ShortUrl: shortUrl}

		err := config.DBConn.First(&urlEntity).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Fatalf("Failed to connect DB: %v", err)
			}

			return nil
		}

		return &urlEntity
	})

	if err != nil {
		return nil, err
	}

	urlEntity := entity.Url{}

	if len(url) > 0 {
		err = json.Unmarshal([]byte(url), &urlEntity)

		if err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
			return nil, err
		}
	}

	return &urlEntity, nil
}

func IncementHitCount(key string) error {

	if err := config.CacheIncrement("hit-count-temp", key); err != nil {
		return err
	}

	return nil
}

func SyncronizeHitCount(inputData map[string]interface{}) error {

	urlEntities := []entity.Url{}
	shortUrlInput := []string{}

	for key := range inputData {
		shortUrlInput = append(shortUrlInput, key)
	}

	if err := config.DBConn.Where(shortUrlInput).Find(&urlEntities).Error; err != nil {
		return err
	}

	for _, urlEntity := range urlEntities {

		for key, input := range inputData {

			if key == urlEntity.ShortUrl {
				hitCount, err := strconv.Atoi(input.(string))
				if err != nil {
					return err
				}

				urlEntity.HitCount += hitCount
				if err := config.DBConn.Save(&urlEntity).Error; err != nil {
					return err
				}
			}
		}

	}

	return nil
}

func BatchInsertUrls(urlEntities []entity.Url) error {

	if err := config.DBConn.Create(&urlEntities).Error; err != nil {
		return err
	}

	return nil
}
