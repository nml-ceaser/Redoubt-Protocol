package core

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Policy struct {
	Mode       string `yaml:"mode"`
	LogPath    string `yaml:"log_path"`
	ListenAddr string `yaml:"listen_addr"`
}

func LoadPolicy(path string) (*Policy, error) {
	b, err := os.ReadFile(path)
	if err != nil { return nil, err }
	var p Policy
	if err := yaml.Unmarshal(b, &p); err != nil { return nil, err }
	return &p, nil
}
