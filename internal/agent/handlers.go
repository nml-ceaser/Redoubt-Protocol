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

type triggerReq struct {
	Reason string ` + "`json:\"reason\"`" + `
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func TriggerHandler(w http.ResponseWriter, r *http.Request) {
	// parse optional reason
	var req triggerReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Reason == "" {
		req.Reason = "manual"
	}

	incident := map[string]interface{}{
		"time":   time.Now().UTC().Format(time.RFC3339),
		"reason": req.Reason,
		"remote": r.RemoteAddr,
	}

	// persist incident
	fn := filepath.Join("_incidents", time.Now().UTC().Format("20060102-150405")+".json")
	f, err := os.Create(fn)
	if err != nil {
		log.Printf("failed to write incident file: %v", err)
	} else {
		_ = json.NewEncoder(f).Encode(incident)
		_ = f.Close()
	}

	// append a log entry
	log.Printf("Incident triggered: %v", incident)

	// Execute isolation actions (may be simulate)
	isolation.Execute(pol)

	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte("incident triggered"))
}
