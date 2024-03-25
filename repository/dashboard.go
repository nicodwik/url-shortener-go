package dashboard

import (
	"errors"
	"url-shortener-go/entity"

	"gorm.io/gorm"
)

type DashboardResponse struct {
	TotalActive    int64               `json:"total_active"`
	MostVisited    *entity.Redirection `json:"most_visited"`
	MostNotVisited *entity.Redirection `json:"most_not_visited"`
	LastCreated    *entity.Redirection `json:"last_created"`
}

type DashboardContract interface {
	GetDashboardData(userId string) (*DashboardResponse, error)
}

type connection struct {
	db *gorm.DB
}

func InitDasboardRepository(db *gorm.DB) *connection {
	return &connection{db}
}

func (conn *connection) GetDashboardData(userId string) (*DashboardResponse, error) {

	dashboardResponse := DashboardResponse{}
	var totalRedirections int64

	var redirection entity.Redirection
	var mostVisitedRedirection *entity.Redirection
	var mostNotVisitedRedirection entity.Redirection
	var lastCreatedRedirection entity.Redirection

	// get total active
	query1 := entity.Redirection{UserId: userId, Status: "active"}
	if err := conn.db.Model(&redirection).Where(query1).Count(&totalRedirections).Error; err != nil {
		return nil, err
	}

	// get most visited
	query2 := entity.Redirection{UserId: userId}
	if err := conn.db.Where(query2).Group("hit_count").Having("MAX(hit_count)").Order("hit_count DESC").First(&mostVisitedRedirection).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// get most not visited
	query3 := entity.Redirection{UserId: userId}
	if err := conn.db.Where(query3).Group("hit_count").Having("MIN(hit_count) >= 0").Order("hit_count ASC").First(&mostNotVisitedRedirection).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// get last created
	query4 := entity.Redirection{UserId: userId}
	if err := conn.db.Where(query4).Group("created_at").Having("MIN(created_at)").Order("created_at DESC").First(&lastCreatedRedirection).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	dashboardResponse.TotalActive = totalRedirections
	dashboardResponse.MostVisited = mostVisitedRedirection.NullSafe()
	dashboardResponse.MostNotVisited = mostNotVisitedRedirection.NullSafe()
	dashboardResponse.LastCreated = lastCreatedRedirection.NullSafe()

	return &dashboardResponse, nil
}
