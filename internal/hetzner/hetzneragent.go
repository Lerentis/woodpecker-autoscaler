package hetzner

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
)

var USER_DATA_TEMPLATE = `
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
			- {{ $key }}: {{ $val }}
	  	{{- end }}
  path: /root/docker-compose.yml
runcmd:
- [ sh, -xc, "cd /root; docker run --rm --privileged multiarch/qemu-user-static --reset -p yes; docker compose up -d" ]
`

type UserDataConfig struct {
	Image     string
	EnvConfig map[string]string
}

func generateConfig(cfg *config.Config) (string, error) {
	envConfig := map[string]string{}
	envConfig["WOODPECKER_SERVER"] = cfg.WoodpeckerInstance
	envConfig["WOODPECKER_AGENT_SECRET"] = cfg.WoodpeckerAgentSecret
	envConfig["WOODPECKER_FILTER_LABELS"] = cfg.LabelSelector
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
