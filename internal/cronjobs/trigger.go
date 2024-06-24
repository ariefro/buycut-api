package cronjobs

import (
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func Trigger() {
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scd := gocron.NewScheduler(jakartaTime)
	_, _ = scd.Every(10).Minute().Do(func() {
		log.Println("OK")
	})
	scd.StartAsync()
}
