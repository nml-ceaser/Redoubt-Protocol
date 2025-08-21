package isolation

import (
	"log"

	"redoubt-protocol/internal/policy"
)

// HandleAndroid currently only logs (Termux without root cannot change firewall).
func HandleAndroid(actions []string, pol *policy.Policy) {
	for _, a := range actions {
		log.Printf("[Isolation][android] action=%s (simulate=%v)", a, pol.Simulate)
		if a == "log_only" {
			log.Println("[Isolation][android] log_only: no enforcement performed on Android/Termux by default")
		} else {
			// Unknown action: in the future, implement root-based iptables or termux-api calls
			log.Printf("[Isolation][android] unknown action: %s (no-op)", a)
		}
	}
}
