package woodpecker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/models"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/utils"

	log "github.com/sirupsen/logrus"
)

func DecomAgent(cfg *config.Config, agentId int64) error {
	apiRoute := fmt.Sprintf("%s/api/agents/%d", cfg.WoodpeckerInstance, agentId)
	req, err := http.NewRequest("DELETE", apiRoute, nil)
	if err != nil {
		return fmt.Errorf("Could not create delete request: %s", err.Error())
	}
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.WoodpeckerApiToken))

	log.WithFields(log.Fields{
		"Caller": "DecomAgent",
	}).Debugf("Deleting agent with id %d from woodpecker", agentId)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not delete agent: %s", err.Error())
	}
	defer resp.Body.Close()
	return nil
}

func GetAgentIdByName(cfg *config.Config, name string) (int, error) {
	apiRoute := fmt.Sprintf("%s/api/agents?page=1&perPage=100", cfg.WoodpeckerInstance)
	req, err := http.NewRequest(http.MethodGet, apiRoute, nil)
	if err != nil {
		return 0, fmt.Errorf("Could not create agent query request: %s", err.Error())
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.WoodpeckerApiToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Could not query agent list: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Invalid status code from API: %d", resp.StatusCode)
	}
	agentList := new(models.AgentList)
	err = json.NewDecoder(resp.Body).Decode(agentList)
	if err != nil {
		return 0, fmt.Errorf("Could not unmarshal api response: %s", err.Error())
	}

	for _, agent := range agentList.Agents {
		if agent.Name == name {
			log.WithFields(log.Fields{
				"Caller": "GetAgentIdByName",
			}).Debugf("Found ID %d for Agent %s", agent.ID, name)
			return int(agent.ID), nil
		}
	}
	return 0, fmt.Errorf("Agent with name %s is not in server", name)
}

func ListAgents(cfg *config.Config) (*models.AgentList, error) {
	agentList := new(models.AgentList)
	apiRoute := fmt.Sprintf("%s/api/agents?page=1&perPage=100", cfg.WoodpeckerInstance)
	req, err := http.NewRequest(http.MethodGet, apiRoute, nil)
	if err != nil {
		return agentList, fmt.Errorf("Could not create agent query request: %s", err.Error())
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.WoodpeckerApiToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return agentList, fmt.Errorf("Could not query agent list: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return agentList, fmt.Errorf("Invalid status code from API: %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(agentList)
	if err != nil {
		return agentList, fmt.Errorf("Could not unmarshal api response: %s", err.Error())
	}
	return agentList, nil
}

func CreateWoodpeckerAgent(cfg *config.Config) (*models.Agent, error) {
	name := fmt.Sprintf("woodpecker-autoscaler-agent-%s", utils.RandStringBytes(5))
	agentRequest := models.AgentRequest{
		Name:       name,
		NoSchedule: false,
	}
	jsonBody, _ := json.Marshal(agentRequest)
	bodyReader := bytes.NewReader(jsonBody)

	apiRoute := fmt.Sprintf("%s/api/agents", cfg.WoodpeckerInstance)
	log.WithFields(log.Fields{
		"Caller": "CreateWoodpeckerAgent",
	}).Debugf("Sending the following data to %s: %s", apiRoute, jsonBody)
	req, err := http.NewRequest(http.MethodPost, apiRoute, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("Could not create agent request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.WoodpeckerApiToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not create new Agent: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code from API: %d", resp.StatusCode)
	}
	newAgent := new(models.Agent)
	err = json.NewDecoder(resp.Body).Decode(newAgent)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal api response: %s", err.Error())
	}
	return newAgent, nil

}
