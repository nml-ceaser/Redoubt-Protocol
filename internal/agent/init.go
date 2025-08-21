package agent

import (
	"log"
	"os"
)

func Init() error {
	// Ensure log directory
	if err := os.MkdirAll("_logs", 0755); err != nil {
		return err
	}

	// Ensure incident directory
	if err := os.MkdirAll("_incidents", 0755); err != nil {
		return err
	}

	// Create log file
	logFile, err := os.OpenFile("_logs/redoubt.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	log.SetOutput(logFile)
	log.Println("Agent initialized.")
	return nil
}
