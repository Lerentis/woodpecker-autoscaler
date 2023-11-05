package hetzner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
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
          - {{ $key }}="{{ $val }}"
	  	{{- end }}
  path: /root/docker-compose.yml
runcmd:
- [ sh, -xc, "cd /root; docker run --rm --privileged multiarch/qemu-user-static --reset -p yes; docker compose up -d" ]
`

type UserDataConfig struct {
	Image     string
	EnvConfig map[string]string
}

func generateConfig(cfg *config.Config, name string) (string, error) {
	envConfig := map[string]string{}
	envConfig["WOODPECKER_SERVER"] = cfg.WoodpeckerGrpc
	envConfig["WOODPECKER_GRPC_SECURE"] = "true" // TODO: should probably made configurable
	envConfig["WOODPECKER_AGENT_SECRET"] = cfg.WoodpeckerAgentSecret
	envConfig["WOODPECKER_FILTER_LABELS"] = cfg.WoodpeckerLabelSelector
	envConfig["WOODPECKER_HOSTNAME"] = name
	config := UserDataConfig{
		Image:     "woodpeckerci/woodpecker-agent:latest",
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

func CreateNewAgent(cfg *config.Config) (*hcloud.Server, error) {
	client := hcloud.NewClient(hcloud.WithToken(cfg.HcloudToken))
	name := fmt.Sprintf("woodpecker-autoscaler-agent-%s", utils.RandStringBytes(5))
	userdata, err := generateConfig(cfg, name)
	img, _, err := client.Image.GetByNameAndArchitecture(context.Background(), "docker-ce", "x86")
	loc, _, err := client.Location.GetByName(context.Background(), cfg.HcloudRegion)
	pln, _, err := client.ServerType.GetByName(context.Background(), cfg.HcloudInstanceType)
	key, _, err := client.SSHKey.GetByName(context.Background(), cfg.HcloudSSHKey)
	dc, _, err := client.Datacenter.GetByName(context.Background(), cfg.HcloudDatacenter)
	labels := map[string]string{}
	labels["Role"] = "WoodpeckerAgent"
	labels["ControledBy"] = "WoodpeckerAutoscaler"

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse agent spec: %s", err.Error()))
	}

	res, _, err := client.Server.Create(context.Background(), hcloud.ServerCreateOpts{
		Name:             name,
		ServerType:       pln,
		Image:            img,
		SSHKeys:          []*hcloud.SSHKey{key},
		Location:         loc,
		Datacenter:       dc,
		UserData:         userdata,
		StartAfterCreate: utils.BoolPointer(true),
		Labels:           labels,
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

func DecomNode(cfg *config.Config, server *hcloud.Server) error {
	client := hcloud.NewClient(hcloud.WithToken(cfg.HcloudToken))
	log.WithFields(log.Fields{
		"Caller": "DecomNode",
	}).Debugf("Deleting %s node", server.Name)
	_, _, err := client.Server.DeleteWithResult(context.Background(), server)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not delete Agent: %s", err.Error()))
	}
	return nil
}

func RefreshNodeInfo(cfg *config.Config, serverID int) (*hcloud.Server, error) {
	client := hcloud.NewClient(hcloud.WithToken(cfg.HcloudToken))
	server, _, err := client.Server.GetByID(context.Background(), serverID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not refresh server info: %s", err.Error()))
	}
	return server, nil
}
