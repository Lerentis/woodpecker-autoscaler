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
		ownedNodes, err := hetzner.ListAgents(cfg)
		if err != nil {
			log.WithFields(log.Fields{
				"Caller": "Main",
			}).Fatal(fmt.Sprintf("Error collecting owned hetzner nodes: %s", err.Error()))
		}
		log.WithFields(log.Fields{
			"Caller": "Main",
		}).Infof("Currently owning %d Agents", len(ownedNodes))
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
				}).Infof("%s is in status %s", server.Name, server.Status)
				time.Sleep(30 * time.Second)
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
				for _, server := range ownedNodes {
					hetzner.DecomNode(cfg, &server)
					agentId, err := woodpecker.GetAgentIdByName(cfg, server.Name)
					if err != nil {
						log.WithFields(log.Fields{
							"Caller": "Main",
						}).Warnf("Could not find agent %s in woodpecker. Assuming it was never added", server.Name)
					} else {
						woodpecker.DecomAgent(cfg, agentId)
					}
				}
			}
		}
		time.Sleep(time.Duration(cfg.CheckInterval) * time.Minute)
	}
}
