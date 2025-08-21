package isolation

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Policy struct {
	Isolation map[string][]string `yaml:"isolation"`
}

func loadPolicy() (*Policy, error) {
	data, err := os.ReadFile("configs/policy.yaml")
	if err != nil {
		return nil, err
	}
	var p Policy
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func Execute() {
	osType := DetectOS()
	log.Printf("[Isolation] Detected OS: %s", osType)

	policy, err := loadPolicy()
	if err != nil {
		log.Printf("[Isolation] Failed to load policy: %v", err)
		return
	}

	actions, ok := policy.Isolation[osType]
	if !ok {
		log.Printf("[Isolation] No rules for OS: %s", osType)
		return
	}

	for _, action := range actions {
		// For now, just log instead of enforcing
		log.Printf("[Isolation] (%s) would %s", osType, action)
		fmt.Printf("[Isolation] (%s) would %s\n", osType, action)
	}
}
