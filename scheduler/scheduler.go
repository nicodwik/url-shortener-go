package scheduler

import (
	"fmt"
	"url-shortener-go/config"
	RedirectionRepository "url-shortener-go/repository/redirection"

	"github.com/jasonlvhit/gocron"
)

func Init() {
	scheduler := gocron.NewScheduler()

	fmt.Println("init scheduler")

	scheduler.Every(5).Minutes().Do(syncHitCountUrl)

	<-scheduler.Start()
}

func syncHitCountUrl() error {
	fmt.Println("sync [hit count] started")
	cacheData, err := config.CacheGetHitCountTemp()
	if err != nil {
		return err
	}

	fmt.Printf("[hit count] data: %v items \n", len(cacheData))
	if len(cacheData) > 0 {
		err = RedirectionRepository.SyncronizeHitCount(cacheData)
		if err != nil {
			return err
		}

		err = config.CacheResetHitCountTemp()
		if err != nil {
			return err
		}

		fmt.Println("sync [hit count] success")
	}

	fmt.Println("====================================")

	return nil
}
