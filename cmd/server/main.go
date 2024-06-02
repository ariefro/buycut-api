package main

import (
	"github.com/ariefro/buycut-api/internal/initializer"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func main() {
	initializer.InitializedServer()
}
