package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"redoubt-protocol/internal/agent"
	"redoubt-protocol/internal/policy"
)

func main() {
	// Load policy early so Init can use simulate flag if needed
	pol, err := policy.LoadPolicy("configs/policy.yaml")
	if err != nil {
		log.Printf("warning: failed to load policy (using defaults): %v", err)
		pol = policy.DefaultPolicy()
	}

	// Initialize agent (logging + incidents)
	if err := agent.Init(pol); err != nil {
		log.Fatalf("failed to initialize agent: %v", err)
	}

	// Register routes
	http.HandleFunc("/healthz", agent.HealthHandler)
	http.HandleFunc("/trigger", agent.TriggerHandler)

	// Listen in a goroutine so we can gracefully handle signals
	go func() {
		addr := ":8686"
		log.Printf("Redoubt Protocol Agent listening on %s (simulate=%v)", addr, pol.Simulate)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for termination signals to exit cleanly
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down agent.")
}
