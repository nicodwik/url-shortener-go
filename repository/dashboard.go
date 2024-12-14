package dashboard

import (
	"errors"
	"sync"
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

type DashboardResponseChan struct {
	name string
	data interface{}
}

func InitDasboardRepository(db *gorm.DB) DashboardContract {
	return &connection{db}
}

func (conn *connection) GetDashboardData(userId string) (*DashboardResponse, error) {
	totalExpectedData := 4
	dashboardResponse := DashboardResponse{}

	c := make(chan DashboardResponseChan, totalExpectedData)
	var wg sync.WaitGroup

	wg.Add(totalExpectedData)
	go conn.getTotalActive(userId, &wg, c)
	go conn.getMostVisited(userId, &wg, c)
	go conn.getMostNotVisited(userId, &wg, c)
	go conn.getLastCreated(userId, &wg, c)

	wg.Wait()
	close(c)

	for d := range c {
		if val, ok := d.data.(int64); ok {
			dashboardResponse.TotalActive = val
		}

		if val, ok := d.data.(entity.Redirection); ok {
			switch d.name {
			case "mv":
				dashboardResponse.MostVisited = &val
			case "mnv":
				dashboardResponse.MostNotVisited = &val
			case "lc":
				dashboardResponse.LastCreated = &val
			}
		}
	}

	return &dashboardResponse, nil
}

// get total active
func (conn *connection) getTotalActive(userId string, wg *sync.WaitGroup, c chan<- DashboardResponseChan) {
	defer wg.Done()

	var totalRedirections int64
	if err := conn.db.Model(&entity.Redirection{}).Where(&entity.Redirection{UserId: userId, Status: "active"}).Count(&totalRedirections).Error; err != nil {
		return
	}
	c <- DashboardResponseChan{name: "ta", data: totalRedirections}
}

// get most visited
func (conn *connection) getMostVisited(userId string, wg *sync.WaitGroup, c chan<- DashboardResponseChan) {
	defer wg.Done()

	var mostVisitedRedirection entity.Redirection
	if err := conn.db.Where(&entity.Redirection{UserId: userId}).Group("hit_count").Having("MAX(hit_count)").Order("hit_count DESC").First(&mostVisitedRedirection).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	}
	c <- DashboardResponseChan{name: "mv", data: mostVisitedRedirection}
}

// get most not visited
func (conn *connection) getMostNotVisited(userId string, wg *sync.WaitGroup, c chan<- DashboardResponseChan) {
	defer wg.Done()

	var mostNotVisitedRedirection entity.Redirection
	if err := conn.db.Where(&entity.Redirection{UserId: userId}).Group("hit_count").Having("MIN(hit_count) >= 0").Order("hit_count ASC").First(&mostNotVisitedRedirection).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	}
	c <- DashboardResponseChan{name: "mnv", data: mostNotVisitedRedirection}
}

// get last created
func (conn *connection) getLastCreated(userId string, wg *sync.WaitGroup, c chan<- DashboardResponseChan) {
	defer wg.Done()

	var lastCreatedRedirection entity.Redirection
	if err := conn.db.Where(&entity.Redirection{UserId: userId}).Group("created_at").Having("MIN(created_at)").Order("created_at DESC").First(&lastCreatedRedirection).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	}
	c <- DashboardResponseChan{name: "lc", data: lastCreatedRedirection}
}
