package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/configor"
)

type Config = struct {
	LogLevel                string `default:"Info" env:"WOODPECKER_AUTOSCALER_LOGLEVEL"`
	CheckInterval           int    `default:"15" env:"WOODPECKER_AUTOSCALER_CHECK_INTERVAL"`
	DryRun                  bool   `default:"false" env:"WOODPECKER_AUTOSCALER_DRY_RUN"`
	CostOptimizedMode       bool   `default:"false" env:"WOODPECKER_AUTOSCALER_COST_OPTIMIZED"`
	WoodpeckerLabelSelector string `default:"uploadfilter24.eu/instance-role=Woodpecker" env:"WOODPECKER_AUTOSCALER_WOODPECKER_LABEL_SELECTOR"`
	WoodpeckerInstance      string `default:"" env:"WOODPECKER_AUTOSCALER_WOODPECKER_INSTANCE"`
	WoodpeckerGrpc          string `default:"" env:"WOODPECKER_AUTOSCALER_WOODPECKER_GRPC"`
	WoodpeckerAgentSecret   string `default:"" env:"WOODPECKER_AUTOSCALER_WOODPECKER_AGENT_SECRET"`
	WoodpeckerApiToken      string `default:"" env:"WOODPECKER_AUTOSCALER_WOODPECKER_API_TOKEN"`
	WoodpeckerAgentVersion  string `default:"latest" env:"WOODPECKER_AUTOSCALER_WOODPECKER_AGENT_VERSION"`
	HcloudToken             string `default:"" env:"WOODPECKER_AUTOSCALER_HCLOUD_TOKEN"`
	HcloudInstanceType      string `default:"cpx21" env:"WOODPECKER_AUTOSCALER_HCLOUD_INSTANCE_TYPE"`
	HcloudLocation          string `default:"" env:"WOODPECKER_AUTOSCALER_HCLOUD_LOCATION"`
	HcloudSSHKeys           string `default:"" env:"WOODPECKER_AUTOSCALER_HCLOUD_SSH_KEYS"`
	HcloudIPv6Only          bool   `default:"false" env:"WOODPECKER_AUTOSCALER_HCLOUD_IPV6_ONLY"`
}

func GenConfig() (cfg *Config, err error) {

	cfg = &Config{}

	err = configor.New(&configor.Config{
		ENVPrefix:          "WOODPECKER_AUTOSCALER",
		AutoReload:         true,
		Silent:             true,
		AutoReloadInterval: time.Minute}).Load(cfg, "config.json")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error generating Config: %s", err.Error()))
	}
	return cfg, nil
}
