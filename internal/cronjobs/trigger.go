package cronjobs

import (
	"context"
	"time"

	"github.com/ariefro/buycut-api/internal/company"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func Trigger(companyController company.Controller) {
	// Attempt to load the Jakarta time zone location
	jakartaTime, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// If it fails, log the error and use the default UTC time zone
		log.Printf("Failed to load Jakarta time zone: %v. Falling back to UTC.", err)
		jakartaTime = time.UTC
	}

	// Initialize the scheduler with the time zone location
	ctx := context.Background()
	scd := gocron.NewScheduler(jakartaTime)
	_, _ = scd.Every(3).Minute().Do(func() {
		err = companyController.FindOneDummy(ctx)
		if err != nil {
			log.Println("Error running cron job:", err)
		}

		log.Println("OK")
	})
	scd.StartAsync()
}
