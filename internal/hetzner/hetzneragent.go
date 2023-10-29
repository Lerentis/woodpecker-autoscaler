package hetzner

import (
	"os"
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

func generateConfig(cfg *config.Config) {
	envConfig := map[string]string{}
	envConfig["WOODPECKER_SERVER"] = cfg.WoodpeckerInstance
	envConfig["WOODPECKER_AGENT_SECRET"] = cfg.WoodpeckerAgentSecret
	envConfig["WOODPECKER_FILTER_LABELS"] = cfg.LabelSelector
	config := UserDataConfig{
		Image:     "woodpeckerci/woodpecker-agent:latest",
		EnvConfig: envConfig,
	}
	tmpl, err := template.New("test").Parse(USER_DATA_TEMPLATE)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, config)
	if err != nil {
		panic(err)
	}
}
