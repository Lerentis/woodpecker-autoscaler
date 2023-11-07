package hetzner

import (
	"testing"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
)

func TestGenerateUserData(t *testing.T) {
	cfg := config.Config{
		LogLevel:                "Info",
		CheckInterval:           5,
		DryRun:                  false,
		WoodpeckerLabelSelector: "uploadfilter24.eu/instance-role=WoodpeckerTest",
		WoodpeckerInstance:      "http://woodpecker.test.tld",
		WoodpeckerGrpc:          "grpc-test.woodpecker.test.tld:443",
		WoodpeckerAgentSecret:   "Geheim1!",
		WoodpeckerApiToken:      "VeryGeheim1!",
		HcloudToken:             "EvenMoreGeheim1!",
		HcloudInstanceType:      "cpx21",
		HcloudRegion:            "eu-central",
		HcloudDatacenter:        "fsn1-dc14",
		HcloudSSHKeys:           "test-key",
	}
	wanted := `
#cloud-config
write_files:
- content: |
    # docker-compose.yml
    version: '3'
    services:
      woodpecker-agent:
        image: woodpeckerci/woodpecker-agent:latest
        command: agent
        restart: always
        volumes:
          - /var/run/docker.sock:/var/run/docker.sock
        environment:
          - WOODPECKER_AGENT_SECRET="Geheim1!"
          - WOODPECKER_FILTER_LABELS="uploadfilter24.eu/instance-role=WoodpeckerTest"
          - WOODPECKER_GRPC_SECURE=true
          - WOODPECKER_HOSTNAME="test-instance"
          - WOODPECKER_SERVER="grpc-test.woodpecker.test.tld:443"
  path: /root/docker-compose.yml
runcmd:
- [ sh, -xc, "cd /root; docker run --rm --privileged multiarch/qemu-user-static --reset -p yes; docker compose up -d" ]
`
	got, err := generateConfig(&cfg, "test-instance")
	if err != nil {
		t.Errorf("Error in generating Config: %v", err)
	}
	if wanted != got {
		t.Errorf("got:\n%v\n, wanted:\n%v", got, wanted)
	}
}
