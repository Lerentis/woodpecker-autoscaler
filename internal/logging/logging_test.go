package logging

import (
	"testing"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
	log "github.com/sirupsen/logrus"
)

func TestLoggingDebug(t *testing.T) {
	cfg := config.Config{LogLevel: "Debug"}
	ConfigureLogger(&cfg)
	if log.GetLevel() != log.DebugLevel {
		t.Fatalf("expected DebugLevel, got %v", log.GetLevel())
	}
}

func TestLoggingInfo(t *testing.T) {
	cfg := config.Config{LogLevel: "Info"}
	ConfigureLogger(&cfg)
	if log.GetLevel() != log.InfoLevel {
		t.Fatalf("expected InfoLevel, got %v", log.GetLevel())
	}
}

func TestLoggingWarning(t *testing.T) {
	cfg := config.Config{LogLevel: "Warn"}
	ConfigureLogger(&cfg)
	if log.GetLevel() != log.WarnLevel {
		t.Fatalf("expected WarnLevel, got %v", log.GetLevel())
	}
}

func TestLoggingError(t *testing.T) {
	cfg := config.Config{LogLevel: "Error"}
	ConfigureLogger(&cfg)
	if log.GetLevel() != log.ErrorLevel {
		t.Fatalf("expected ErrorLevel, got %v", log.GetLevel())
	}
}
