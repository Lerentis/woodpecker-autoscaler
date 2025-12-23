package woodpecker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/models"
)

func TestQueueInfoAndChecks(t *testing.T) {
	// Create queue info with one pending job matching label and one running matching
	qi := models.QueueInfo{
		Pending: []models.JobInformation{
			{ID: "1", Labels: map[string]string{"role": "worker"}},
		},
		Running: []models.JobInformation{
			{ID: "2", Labels: map[string]string{"role": "worker"}},
		},
		Stats: models.Stats{PendingCount: 1, RunningCount: 1},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/queue/info" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(qi)
	}))
	defer srv.Close()

	cfg := config.Config{
		WoodpeckerInstance:      srv.URL,
		WoodpeckerApiToken:      "t",
		WoodpeckerLabelSelector: "role=worker",
	}

	// Test QueueInfo
	var got models.QueueInfo
	if err := QueueInfo(&cfg, &got); err != nil {
		t.Fatalf("QueueInfo failed: %v", err)
	}
	if got.Stats.PendingCount != 1 || got.Stats.RunningCount != 1 {
		t.Fatalf("unexpected stats: %#v", got.Stats)
	}

	// Test CheckPending
	pending, err := CheckPending(&cfg)
	if err != nil {
		t.Fatalf("CheckPending error: %v", err)
	}
	if pending != 1 {
		t.Fatalf("expected 1 pending, got %d", pending)
	}

	// Test CheckRunning
	running, err := CheckRunning(&cfg)
	if err != nil {
		t.Fatalf("CheckRunning error: %v", err)
	}
	if running != 1 {
		t.Fatalf("expected 1 running, got %d", running)
	}
}
