package hetzner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/models"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/utils"
	"github.com/hetznercloud/hcloud-go/hcloud"

	log "github.com/sirupsen/logrus"
)

var USER_DATA_TEMPLATE = `
#cloud-config
write_files:
- content: |
    # docker-compose.yml
    version: '3'
    services:
      woodpecker-agent:
        image: {{ .Image }}
        command: agent
        restart: always
        volumes:
          - /var/run/docker.sock:/var/run/docker.sock
        environment:
		{{- range $key, $val := .EnvConfig }}
          - {{ $key }}={{ $val }}
	  	{{- end }}
  path: /root/docker-compose.yml
runcmd:
- [ sh, -xc, "cd /root; docker run --rm --privileged multiarch/qemu-user-static --reset -p yes; docker compose up -d" ]
`

type UserDataConfig struct {
	Image     string
	EnvConfig map[string]interface{}
}

func generateConfig(cfg *config.Config, name string, agentToken string) (string, error) {
	envConfig := map[string]interface{}{
		"WOODPECKER_SERVER":        fmt.Sprintf("%s", cfg.WoodpeckerGrpc),
		"WOODPECKER_GRPC_SECURE":   true,
		"WOODPECKER_AGENT_SECRET":  fmt.Sprintf("%s", agentToken),
		"WOODPECKER_FILTER_LABELS": fmt.Sprintf("%s", cfg.WoodpeckerLabelSelector),
		"WOODPECKER_HOSTNAME":      fmt.Sprintf("%s", name),
		"WOODPECKER_MAX_WORKFLOWS": 4,
	}
	config := UserDataConfig{
		Image:     fmt.Sprintf("woodpeckerci/woodpecker-agent:%s", cfg.WoodpeckerAgentVersion),
		EnvConfig: envConfig,
	}
	tmpl, err := template.New("userdata").Parse(USER_DATA_TEMPLATE)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Errors in userdata template: %s", err.Error()))
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, &config)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not render userdata template: %s", err.Error()))
	}
	return buf.String(), nil
}

func CreateNewAgent(cfg *config.Config, woodpeckerAgent *models.Agent) (*hcloud.Server, error) {
	client := hcloud.NewClient(hcloud.WithToken(cfg.HcloudToken))
	userdata, err := generateConfig(cfg, woodpeckerAgent.Name, woodpeckerAgent.Token)
	keys := []*hcloud.SSHKey{}
	for _, keyName := range strings.Split(cfg.HcloudSSHKeys, ",") {
		key, _, err := client.SSHKey.GetByName(context.Background(), keyName)
		if err != nil {
			log.WithFields(log.Fields{
				"Caller": "CreateNewAgent",
			}).Warnf("Failed to look up ssh key %s: %s", keyName, err.Error())
			continue
		}
		keys = append(keys, key)
	}
	img, _, err := client.Image.GetByNameAndArchitecture(context.Background(), "docker-ce", "x86")
	utils.CheckError(err, "GetImageByNameAndArchitecture")
	loc, _, err := client.Location.GetByName(context.Background(), cfg.HcloudLocation)
	utils.CheckError(err, "GetRegionByName")
	pln, _, err := client.ServerType.GetByName(context.Background(), cfg.HcloudInstanceType)
	utils.CheckError(err, "GetServerTypeByName")
	labels := map[string]string{}
	labels["Role"] = "WoodpeckerAgent"
	labels["ControledBy"] = "WoodpeckerAutoscaler"
	labels["ID"] = fmt.Sprintf("%d", woodpeckerAgent.ID)

	networkConf := hcloud.ServerCreatePublicNet{
		EnableIPv4: !cfg.HcloudIPv6Only,
		EnableIPv6: true,
	}

	res, _, err := client.Server.Create(context.Background(), hcloud.ServerCreateOpts{
		Name:             woodpeckerAgent.Name,
		ServerType:       pln,
		Image:            img,
		SSHKeys:          keys,
		Location:         loc,
		UserData:         userdata,
		StartAfterCreate: utils.BoolPointer(true),
		Labels:           labels,
		PublicNet:        &networkConf,
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create new Agent: %s", err.Error()))
	}

	log.WithFields(log.Fields{
		"Caller": "CreateNewAgent",
	}).Infof("Created new Build Agent %s", res.Server.Name)

	return res.Server, nil
}

func ListAgents(cfg *config.Config) ([]hcloud.Server, error) {
	client := hcloud.NewClient(hcloud.WithToken(cfg.HcloudToken))
	allServers, err := client.Server.All(context.Background())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not query Server list: %s", err.Error()))
	}
	myServers := []hcloud.Server{}
	for _, server := range allServers {
		val, exists := server.Labels["ControledBy"]
		if exists && val == "WoodpeckerAutoscaler" {
			myServers = append(myServers, *server)
			log.WithFields(log.Fields{
				"Caller": "ListAgents",
			}).Debugf("Owning %s Hetzner node", server.Name)
		}
	}
	return myServers, nil
}

func DecomNode(cfg *config.Config, server *hcloud.Server) (int64, error) {
	client := hcloud.NewClient(hcloud.WithToken(cfg.HcloudToken))
	var woodpeckerAgentID int64
	val, exists := server.Labels["ID"]
	if exists {
		log.WithFields(log.Fields{
			"Caller": "DecomNode",
		}).Debugf("Found woodpecker agent id: %s", val)
		woodpeckerAgentID, _ = strconv.ParseInt(val, 10, 64)
	} else {
		log.WithFields(log.Fields{
			"Caller": "DecomNode",
		}).Warnf("Did not find woodpecker agent id for node %s", server.Name)
	}
	log.WithFields(log.Fields{
		"Caller": "DecomNode",
	}).Debugf("Deleting %s node", server.Name)
	_, _, err := client.Server.DeleteWithResult(context.Background(), server)
	if err != nil {
		return woodpeckerAgentID, errors.New(fmt.Sprintf("Could not delete Agent: %s", err.Error()))
	}
	return woodpeckerAgentID, nil
}

func RefreshNodeInfo(cfg *config.Config, serverID int) (*hcloud.Server, error) {
	client := hcloud.NewClient(hcloud.WithToken(cfg.HcloudToken))
	server, _, err := client.Server.GetByID(context.Background(), serverID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not refresh server info: %s", err.Error()))
	}
	return server, nil
}

func CheckRuntime(cfg *config.Config, server *hcloud.Server) (time.Time, error) {
	server, err := RefreshNodeInfo(cfg, server.ID)
	now := time.Now()
	if err != nil {
		return time.Time{}, errors.New(fmt.Sprintf("Could not check Runtime: %s", err.Error()))
	}
	return server.Created.Add(time.Duration(now.Minute())), nil
}
