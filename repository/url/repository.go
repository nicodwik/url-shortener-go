package url

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"url-shortener-go/config"
	"url-shortener-go/entity"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PaginationResponse struct {
	Items interface{}    `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	CurrentPage  int   `json:"current_page,omitempty"`
	Limit        int   `json:"limit,omitempty"`
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	ItemsPerPage int   `json:"items_per_page"`
}

func Paginate(c echo.Context, model interface{}, paginationMeta *PaginationMeta, db *gorm.DB) func(db *gorm.DB) *gorm.DB {

	var totalItems int64
	db.Model(&model).Count(&totalItems)

	paginationMeta.TotalItems = totalItems

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}

	paginationMeta.Limit = limit
	paginationMeta.CurrentPage = page
	paginationMeta.TotalPages = (int(totalItems) / limit) + 1

	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func GetAll(c echo.Context) (*PaginationResponse, error) {
	var urls []entity.Url
	paginationResponse := PaginationResponse{}
	paginationMeta := PaginationMeta{}

	if err := config.DBConn.Scopes(Paginate(c, urls, &paginationMeta, config.DBConn)).Order("created_at DESC").Find(&urls).Error; err != nil {
		log.Fatalf("Failed to connect DB: %v", err)

		return &paginationResponse, err
	}

	paginationMeta.ItemsPerPage = len(urls)

	paginationResponse.Items = &urls
	paginationResponse.Meta = paginationMeta

	return &paginationResponse, nil
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
