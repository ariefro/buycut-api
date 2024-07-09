package cronjobs

import (
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func Trigger() {
	// Attempt to load the Jakarta time zone location
	jakartaTime, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// If it fails, log the error and use the default UTC time zone
		log.Printf("Failed to load Jakarta time zone: %v. Falling back to UTC.", err)
		jakartaTime = time.UTC
	}

	// Initialize the scheduler with the time zone location
	scd := gocron.NewScheduler(jakartaTime)
	_, _ = scd.Every(5).Minute().Do(func() {
		log.Println("OK")
	})
	scd.StartAsync()
}
