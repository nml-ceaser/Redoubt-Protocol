package agent

import (
	"log"
	"os"
	"path/filepath"

	"redoubt-protocol/internal/policy"
)

var pol *policy.Policy

func Init(policyObj *policy.Policy) error {
	// keep global policy
	pol = policyObj
	if pol == nil {
		pol = policy.DefaultPolicy()
	}

	// Ensure directories
	if err := os.MkdirAll("_logs", 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll("_incidents", 0o755); err != nil {
		return err
	}

	// Open log file and set global logger
	logPath := filepath.Join("_logs", "redoubt.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	log.SetOutput(f)
	log.Printf("Agent initialized. simulate=%v", pol.Simulate)
	return nil
}
