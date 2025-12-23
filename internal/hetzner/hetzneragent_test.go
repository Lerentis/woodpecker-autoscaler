package hetzner

import (
	"strings"
	"testing"
	"time"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
	"github.com/hetznercloud/hcloud-go/hcloud"
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
		WoodpeckerAgentVersion:  "latest",
		HcloudToken:             "EvenMoreGeheim1!",
		HcloudInstanceType:      "cpx21",
		HcloudLocation:          "fsn1",
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
          - WOODPECKER_AGENT_SECRET=Geheim1!
          - WOODPECKER_FILTER_LABELS=uploadfilter24.eu/instance-role=WoodpeckerTest
          - WOODPECKER_GRPC_SECURE=true
          - WOODPECKER_HOSTNAME=test-instance
          - WOODPECKER_MAX_WORKFLOWS=4
          - WOODPECKER_SERVER=grpc-test.woodpecker.test.tld:443
  path: /root/docker-compose.yml
runcmd:
- [ sh, -xc, "cd /root; docker run --rm --privileged multiarch/qemu-user-static --reset -p yes; docker compose up -d" ]
`
	got, err := generateConfig(&cfg, "test-instance", "Geheim1!")
	if err != nil {
		t.Errorf("Error in generating Config: %v", err)
	}
	if wanted != got {
		t.Errorf("got:\n%v\n, wanted:\n%v", got, wanted)
	}
}

func TestGenerateUserData_MultipleCases(t *testing.T) {
	base := config.Config{
		WoodpeckerGrpc:          "grpc-test.woodpecker.test.tld:443",
		WoodpeckerLabelSelector: "uploadfilter24.eu/instance-role=WoodpeckerTest",
		WoodpeckerAgentVersion:  "latest",
	}

	cases := []struct {
		name         string
		cfg          config.Config
		agentName    string
		agentToken   string
		wantContains []string
	}{
		{
			name:       "basic",
			cfg:        base,
			agentName:  "test-instance",
			agentToken: "Geheim1!",
			wantContains: []string{
				"image: woodpeckerci/woodpecker-agent:latest",
				"- WOODPECKER_AGENT_SECRET=Geheim1!",
				"- WOODPECKER_FILTER_LABELS=uploadfilter24.eu/instance-role=WoodpeckerTest",
				"- WOODPECKER_SERVER=grpc-test.woodpecker.test.tld:443",
			},
		},
		{
			name:       "empty token",
			cfg:        base,
			agentName:  "no-token",
			agentToken: "",
			wantContains: []string{
				"image: woodpeckerci/woodpecker-agent:latest",
				"- WOODPECKER_AGENT_SECRET=",
				"- WOODPECKER_HOSTNAME=no-token",
			},
		},
	}

	for _, tc := range cases {
		got, err := generateConfig(&tc.cfg, tc.agentName, tc.agentToken)
		if err != nil {
			t.Fatalf("%s: generateConfig returned error: %v", tc.name, err)
		}
		for _, want := range tc.wantContains {
			if !strings.Contains(got, want) {
				t.Errorf("%s: expected generated userdata to contain %q, got:\n%s", tc.name, want, got)
			}
		}
	}
}

func TestCheckRuntime_MockedRefresh(t *testing.T) {
	// Mock refreshNodeInfo to return a server with a known Created time
	orig := refreshNodeInfo
	defer func() { refreshNodeInfo = orig }()

	created := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	refreshNodeInfo = func(cfg *config.Config, serverID int) (*hcloud.Server, error) {
		return &hcloud.Server{Created: created}, nil
	}

	cfg := config.Config{}
	// Capture minute before call to avoid flakiness across minute boundary
	minute := time.Now().Minute()
	got, err := CheckRuntime(&cfg, &hcloud.Server{ID: 123})
	if err != nil {
		t.Fatalf("CheckRuntime returned error: %v", err)
	}
	want := created.Add(time.Duration(minute))
	if !got.Equal(want) {
		t.Fatalf("unexpected runtime: got %v, want %v", got, want)
	}
}
