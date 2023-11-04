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

func DecomAgent(cfg *config.Config, agentId int) error {
	apiRoute := fmt.Sprintf("%s/api/agents/%d", cfg.WoodpeckerInstance, agentId)
	req, err := http.NewRequest("DELETE", apiRoute, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create delete request: %s", err.Error()))
	}
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.WoodpeckerApiToken))

	log.WithFields(log.Fields{
		"Caller": "DecomAgent",
	}).Debugf("Deleting %d agent from woodpecker", agentId)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not delete agent: %s", err.Error()))
	}
	defer resp.Body.Close()
	return nil
}

func GetAgentIdByName(cfg *config.Config, name string) (int, error) {
	apiRoute := fmt.Sprintf("%s/api/agents?page=1&perPage=100", cfg.WoodpeckerInstance)
	req, err := http.NewRequest("GET", apiRoute, nil)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Could not create agent query request: %s", err.Error()))
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.WoodpeckerApiToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Could not query agent list: %s", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(fmt.Sprintf("Invalid status code from API: %d", resp.StatusCode))
	}
	agentList := new(models.AgentList)
	err = json.NewDecoder(resp.Body).Decode(agentList)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Could not unmarshal api response: %s", err.Error()))
	}

	for _, agent := range agentList.Agents {
		if agent.Name == name {
			log.WithFields(log.Fields{
				"Caller": "GetAgentIdByName",
			}).Debugf("Found ID %d for Agent %s", agent.ID, name)
			return int(agent.ID), nil
		}
	}
	return 0, errors.New(fmt.Sprintf("Agent with name %s is not in server", name))
}
