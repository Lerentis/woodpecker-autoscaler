package woodpecker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/models"

	log "github.com/sirupsen/logrus"
)

func QueueInfo(cfg *config.Config, target interface{}) error {
	apiRoute := fmt.Sprintf("%s/api/queue/info", cfg.WoodpeckerInstance)
	req, err := http.NewRequest("GET", apiRoute, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create queue request: %s", err.Error()))
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.WoodpeckerApiToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not query queue info: %s", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Error from queue info api: %s", err.Error()))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func CheckPending(cfg *config.Config) (bool, error) {
	expectedKV := strings.Split(cfg.WoodpeckerLabelSelector, "=")
	queueInfo := new(models.QueueInfo)
	err := QueueInfo(cfg, queueInfo)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Error from QueueInfo: %s", err.Error()))
	}
	if queueInfo.Stats.PendingCount > 0 {
		if queueInfo.Pending != nil {
			for _, pendingJobs := range queueInfo.Pending {
				val, exists := pendingJobs.Labels[expectedKV[0]]
				if exists && val == expectedKV[1] {
					log.WithFields(log.Fields{
						"Caller": "CheckPending",
					}).Info("Found pending job for us")
					return true, nil
				} else {
					log.WithFields(log.Fields{
						"Caller": "CheckPending",
					}).Info("No Jobs for us in Queue")
					return false, nil
				}
			}
		}
	}
	return false, nil
}

func CheckRunning(cfg *config.Config) (bool, error) {
	expectedKV := strings.Split(cfg.WoodpeckerLabelSelector, "=")
	queueInfo := new(models.QueueInfo)
	err := QueueInfo(cfg, queueInfo)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Error from QueueInfo: %s", err.Error()))
	}
	if queueInfo.Stats.RunningCount > 0 {
		for _, runningJobs := range queueInfo.Running {
			val, exists := runningJobs.Labels[expectedKV[0]]
			if exists && val == expectedKV[1] {
				log.WithFields(log.Fields{
					"Caller": "CheckRunning",
				}).Info("Found running job for us")
				return true, nil
			} else {
				log.WithFields(log.Fields{
					"Caller": "CheckRunning",
				}).Info("No running job for us")
				return false, nil
			}
		}
	}
	return false, nil
}
