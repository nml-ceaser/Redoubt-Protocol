package isolation

import (
	"log"

	"redoubt-protocol/internal/policy"
)

func Execute(pol *policy.Policy) {
	if pol == nil {
		log.Println("[Isolation] no policy provided; aborting isolation")
		return
	}

	osType := DetectOS()
	log.Printf("[Isolation] Detected OS: %s", osType)

	actions, ok := pol.Isolation[osType]
	if !ok || len(actions) == 0 {
		log.Printf("[Isolation] No actions configured for OS: %s", osType)
		return
	}

	// Dispatch to OS-specific handler
	switch osType {
	case "linux":
		HandleLinux(actions, pol)
	case "windows":
		HandleWindows(actions, pol)
	case "android":
		HandleAndroid(actions, pol)
	default:
		log.Printf("[Isolation] Unsupported OS: %s", osType)
	}
}   
