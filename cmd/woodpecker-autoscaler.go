package main

import (
	"fmt"
	"time"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/health"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/hetzner"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/logging"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/woodpecker"
	"github.com/hetznercloud/hcloud-go/hcloud"
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
		pendingTasks, err := woodpecker.CheckPending(cfg)
		if err != nil {
			log.WithFields(log.Fields{
				"Caller": "Main",
			}).Fatal(fmt.Sprintf("Error checking woodpecker queue: %s", err.Error()))
		}
		if pendingTasks {
			server, err := hetzner.CreateNewAgent(cfg)
			if err != nil {
				log.WithFields(log.Fields{
					"Caller": "Main",
				}).Fatal(fmt.Sprintf("Error spawning new agent: %s", err.Error()))
			}
			for {
				if server.Status == hcloud.ServerStatusRunning {
					log.WithFields(log.Fields{
						"Caller": "Main",
					}).Infof("%s started!", server.Name)
					break
				}
				log.WithFields(log.Fields{
					"Caller": "Main",
				}).Infof("Waiting for agent %s to start", server.Name)
			}
		} else {
			log.WithFields(log.Fields{
				"Caller": "Main",
			}).Info("Checking if agents can be removed")
			runningTasks, err := woodpecker.CheckRunning(cfg)
			if err != nil {
				log.WithFields(log.Fields{
					"Caller": "Main",
				}).Fatal(fmt.Sprintf("Error checking woodpecker queue: %s", err.Error()))
			}
			if runningTasks {
				log.WithFields(log.Fields{
					"Caller": "Main",
				}).Info("Still found running tasks. No agent to be removed")
			} else {
				log.WithFields(log.Fields{
					"Caller": "Main",
				}).Info("No tasks running. Will remove agents")
				// TODO: iterate over agents and remove ours
				// agent name should match hetzner name
			}
		}
		time.Sleep(time.Duration(cfg.CheckInterval) * time.Minute)
	}
}
