package logging

import (
	"os"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
	log "github.com/sirupsen/logrus"
)

func ConfigureLogger(cfg *config.Config) {

	switch cfg.LogLevel {
	case "Debug":
		log.SetLevel(log.DebugLevel)
	case "Info":
		log.SetLevel(log.InfoLevel)
	case "Warn":
		log.SetLevel(log.WarnLevel)
	case "Error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
		log.Warnf("Home: invalid log level supplied: '%s'", cfg.LogLevel)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}
