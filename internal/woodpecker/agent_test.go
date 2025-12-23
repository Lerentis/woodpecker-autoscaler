package woodpecker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/config"
	"git.uploadfilter24.eu/covidnetes/woodpecker-autoscaler/internal/models"
)

func TestCreateAndGetAndDeleteAgent(t *testing.T) {
	// prepare a fake agent to return
	createdAgent := models.Agent{
		ID:    42,
		Name:  "woodpecker-autoscaler-agent-abcde",
		Token: "tok",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/agents", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// ensure content-type
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				t.Fatalf("expected json content-type, got %s", ct)
			}
			body, _ := io.ReadAll(r.Body)
			defer r.Body.Close()
			if !strings.Contains(string(body), "woodpecker-autoscaler-agent-") {
				t.Fatalf("unexpected agent request body: %s", string(body))
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(createdAgent)
			return
		}
		// For GET listing, return an AgentList
		w.WriteHeader(http.StatusOK)
		list := models.AgentList{Agents: []models.Agent{createdAgent}}
		_ = json.NewEncoder(w).Encode(list)
	})

	mux.HandleFunc("/api/agents?page=1&perPage=100", func(w http.ResponseWriter, r *http.Request) {
		// return list in expected format for GetAgentIdByName
		w.WriteHeader(http.StatusOK)
		// GetAgentIdByName expects a models.AgentList; encode accordingly
		list := models.AgentList{Agents: []models.Agent{createdAgent}}
		_ = json.NewEncoder(w).Encode(list)
	})

	// handle delete
	mux.HandleFunc(fmt.Sprintf("/api/agents/%d", createdAgent.ID), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	cfg := config.Config{
		WoodpeckerInstance: srv.URL,
		WoodpeckerApiToken: "testtoken",
	}

	// Test CreateWoodpeckerAgent
	a, err := CreateWoodpeckerAgent(&cfg)
	if err != nil {
		t.Fatalf("CreateWoodpeckerAgent failed: %v", err)
	}
	if a == nil || !strings.HasPrefix(a.Name, "woodpecker-autoscaler-agent-") {
		t.Fatalf("unexpected agent returned: %#v", a)
	}

	// Test GetAgentIdByName
	id, err := GetAgentIdByName(&cfg, a.Name)
	if err != nil {
		t.Fatalf("GetAgentIdByName failed: %v", err)
	}
	if id != int(a.ID) {
		t.Fatalf("unexpected id: got %d want %d", id, a.ID)
	}

	// Test DecomAgent
	if err := DecomAgent(&cfg, a.ID); err != nil {
		t.Fatalf("DecomAgent failed: %v", err)
	}
}

func TestGetAgentIdByName_NotFound(t *testing.T) {
	// server returns empty list
	mux := http.NewServeMux()
	mux.HandleFunc("/api/agents?page=1&perPage=100", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		list := models.AgentList{Agents: []models.Agent{{ID: 1, Name: "other"}}}
		_ = json.NewEncoder(w).Encode(list)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cfg := config.Config{WoodpeckerInstance: srv.URL, WoodpeckerApiToken: "t"}
	_, err := GetAgentIdByName(&cfg, "nonexistent")
	if err == nil {
		t.Fatalf("expected error for unknown agent name")
	}
}
