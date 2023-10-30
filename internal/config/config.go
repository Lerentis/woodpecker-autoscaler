package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/configor"
)

type Config = struct {
	LogLevel              string `default:"Info" env:"WOODPECKER_AUTOSCALER_LOGLEVEL"`
	CheckInterval         int    `default:"15" env:"WOODPECKER_AUTOSCALER_CHECK_INTERVAL"`
	LabelSelector         string `default:"uploadfilter24.eu/instance-role=Woodpecker" env:"WOODPECKER_AUTOSCALER_LABELSELECTOR"`
	WoodpeckerInstance    string `default:"" env:"WOODPECKER_AUTOSCALER_WOODPECKER_INSTANCE"`
	WoodpeckerAgentSecret string `default:"" env:"WOODPECKER_AUTOSCALER_WOODPECKER_AGENT_SECRET"`
	WoodpeckerApiToken    string `default:"" env:"WOODPECKER_AUTOSCALER_WOODPECKER_API_TOKEN"`
	Protocol              string `default:"http" env:"WOODPECKER_AUTOSCALER_PROTOCOL"`
	HcloudToken           string `default:"" env:"WOODPECKER_AUTOSCALER_HCLOUD_TOKEN"`
	InstanceType          string `default:"" env:"WOODPECKER_AUTOSCALER_INSTANCE_TYPE"`
	Zone                  string `default:"" env:"WOODPECKER_AUTOSCALER_ZONE"`
	DryRun                bool   `default:"false" env:"WOODPECKER_AUTOSCALER_DRY_RUN"`
	SSHKey                string `default:"" env:"WOODPECKER_AUTOSCALER_SSH_KEY"`
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
