package main

import (
	"fmt"
	"time"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/health"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/logging"
	log "github.com/sirupsen/logrus"
)

func main() {

	cfg, err := config.GenConfig()
	logging.ConfigureLogger(cfg)

	if err != nil {
		log.WithFields(log.Fields{
			"Caller": "Main",
		}).Fatal(fmt.Sprintf("Error generating Config: %s", err.Error()))
	}

	go func() {
		log.WithFields(log.Fields{
			"Caller": "Main",
		}).Info("Starting Health Endpoint")
		health.StartHealthEndpoint()
	}()

	log.WithFields(log.Fields{
		"Caller": "Main",
	}).Info("Entering main event loop")

	for {
		time.Sleep(time.Duration(cfg.CheckInterval) * time.Minute)
	}
}
