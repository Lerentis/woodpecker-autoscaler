package woodpecker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

func CheckPending(cfg *config.Config) error {
	queueInfo := new(models.QueueInfo)
	err := QueueInfo(cfg, queueInfo)
	if err != nil {
		return errors.New(fmt.Sprintf("Error from QueueInfo: %s", err.Error()))
	}
	if queueInfo.Stats.PendingCount > 0 {
		for _, pendingJobs := range queueInfo.Pending {
			// TODO: separate key and value from LabelSelector and compare them deeply
			_, exists := pendingJobs.Labels[cfg.LabelSelector]
			if exists {
				log.WithFields(log.Fields{
					"Caller": "CheckPending",
				}).Info("Found pending job for us. Requesting new Agent")
			} else {
				log.WithFields(log.Fields{
					"Caller": "CheckPending",
				}).Info("No Jobs for us in Queue")
			}
		}
	}
	return nil
}
