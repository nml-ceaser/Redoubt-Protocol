package policy

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Policy struct {
	Simulate  bool                 `yaml:"simulate"`
	Allowlist []string             `yaml:"allowlist"`
	Isolation map[string][]string  `yaml:"isolation"`
	Controls  map[string]map[string]interface{} `yaml:"controls"`
}

func DefaultPolicy() *Policy {
	return &Policy{
		Simulate: true,
		Allowlist: []string{"127.0.0.1/32"},
		Isolation: map[string][]string{
			"linux":   {"block_network"},
			"windows": {"disable_interface"},
			"android": {"log_only"},
		},
		Controls: map[string]map[string]interface{}{},
	}
}

func LoadPolicy(path string) (*Policy, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Policy
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
