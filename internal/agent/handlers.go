package agent

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"redoubt-protocol/internal/isolation"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func TriggerHandler(w http.ResponseWriter, r *http.Request) {
	incident := map[string]interface{}{
		"time":    time.Now().UTC().Format(time.RFC3339),
		"trigger": "manual",
	}

	// Write incident file
	filename := filepath.Join("_incidents",
		time.Now().UTC().Format("20060102-150405")+".json")
	file, _ := os.Create(filename)
	defer file.Close()
	_ = json.NewEncoder(file).Encode(incident)

	// Log event
	log.Println("Incident triggered:", incident)

	// Call isolation (stubbed for Milestone 2)
	isolation.Execute()

	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte("incident triggered"))
}
