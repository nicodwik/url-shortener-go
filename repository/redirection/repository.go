package redirection

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"strconv"
	"url-shortener-go/config"
	"url-shortener-go/entity"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PaginationResponse struct {
	Items interface{}     `json:"items"`
	Meta  *PaginationMeta `json:"meta"`
}

type PaginationParam struct {
	Context        echo.Context
	Model          interface{}
	Query          interface{}
	PaginationMeta *PaginationMeta
	DB             *gorm.DB
}

type PaginationMeta struct {
	CurrentPage  int     `json:"current_page,omitempty"`
	Limit        int     `json:"limit,omitempty"`
	TotalItems   int     `json:"total_items"`
	TotalPages   float64 `json:"total_pages"`
	ItemsPerPage int     `json:"items_per_page"`
}

func Paginate(paginationParam PaginationParam) func(db *gorm.DB) *gorm.DB {

	var totalItems int64
	paginationParam.DB.Model(&paginationParam.Model).Where(paginationParam.Query).Count(&totalItems)

	paginationParam.PaginationMeta.TotalItems = int(totalItems)

	page, _ := strconv.Atoi(paginationParam.Context.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(paginationParam.Context.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}

	paginationParam.PaginationMeta.Limit = limit
	paginationParam.PaginationMeta.CurrentPage = page
	paginationParam.PaginationMeta.TotalPages = math.Ceil(float64(totalItems) / float64(limit))

	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func GetAll(c echo.Context) (*PaginationResponse, error) {
	var redirections []entity.Redirection
	paginationResponse := PaginationResponse{}
	paginationMeta := PaginationMeta{}
	paginationParam := PaginationParam{
		Context:        c,
		Model:          redirections,
		DB:             config.DBConn,
		PaginationMeta: &paginationMeta,
	}

	if err := config.DBConn.Scopes(Paginate(paginationParam)).Order("created_at DESC").Find(&redirections).Error; err != nil {
		log.Fatalf("Failed to connect DB: %v", err)

		return &paginationResponse, err
	}

	paginationMeta.ItemsPerPage = len(redirections)

	paginationResponse.Items = &redirections
	paginationResponse.Meta = &paginationMeta

	return &paginationResponse, nil
}

func GetAllByUserId(c echo.Context, userId string) (PaginationResponse, error) {
	var redirections []entity.Redirection
	query := entity.Redirection{UserId: userId}
	paginationResponse := PaginationResponse{}
	paginationMeta := PaginationMeta{}
	paginationParam := PaginationParam{
		Context:        c,
		Model:          redirections,
		DB:             config.DBConn,
		Query:          query,
		PaginationMeta: &paginationMeta,
	}

	if err := config.DBConn.Where(&query).Scopes(Paginate(paginationParam)).Order("created_at DESC").Find(&redirections).Error; err != nil {
		log.Fatalf("Failed to connect DB: %v", err)

		return paginationResponse, err
	}

	paginationMeta.ItemsPerPage = len(redirections)

	paginationResponse.Items = &redirections
	paginationResponse.Meta = &paginationMeta

	return paginationResponse, nil
}

func InsertNewUrl(shortUrl string, longUrl string, userId string) (*entity.Redirection, error) {

	redirectionEntity := entity.Redirection{
		// UserId:   "dbc587c2-98a8-4ab2-b306-4954d9e83dbf",
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
		Status:   "active",
	}

	if len(userId) > 0 {
		redirectionEntity.UserId = userId
	}

	if err := config.DBConn.Create(&redirectionEntity).Error; err != nil {
		return nil, err
	}

	return &redirectionEntity, nil
}

func FindRedirection(shortUrl string) (*entity.Redirection, error) {
	url, err := config.CacheRemember(shortUrl, 3600, func() interface{} {
		redirectionEntity := entity.Redirection{ShortUrl: shortUrl}

		err := config.DBConn.First(&redirectionEntity).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Fatalf("Failed to connect DB: %v", err)
			}

			return nil
		}

		return &redirectionEntity
	})

	if err != nil {
		return nil, err
	}

	redirectionEntity := entity.Redirection{}

	if len(url) > 0 {
		err = json.Unmarshal([]byte(url), &redirectionEntity)

		if err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
			return nil, err
		}
	}

	return &redirectionEntity, nil
}

func IncementHitCount(key string) error {

	if err := config.CacheIncrement("hit-count-temp", key); err != nil {
		return err
	}

	return nil
}

func SyncronizeHitCount(inputData map[string]interface{}) error {

	urlEntities := []entity.Redirection{}
	shortUrlInput := []string{}

	for key := range inputData {
		shortUrlInput = append(shortUrlInput, key)
	}

	if err := config.DBConn.Where(shortUrlInput).Find(&urlEntities).Error; err != nil {
		return err
	}

	for _, redirectionEntity := range urlEntities {

		for key, input := range inputData {

			if key == redirectionEntity.ShortUrl {
				hitCount, err := strconv.Atoi(input.(string))
				if err != nil {
					return err
				}

				redirectionEntity.HitCount += hitCount
				if err := config.DBConn.Save(&redirectionEntity).Error; err != nil {
					return err
				}
			}
		}

	}

	return nil
}

func BatchInsertredirections(urlEntities []entity.Redirection) error {

	if err := config.DBConn.Create(&urlEntities).Error; err != nil {
		return err
	}

	return nil
}
