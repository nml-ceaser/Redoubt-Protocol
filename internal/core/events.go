package core

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Event struct {
	Time   time.Time         `json:"time"`
	Kind   string            `json:"kind"`
	Reason string            `json:"reason"`
	Meta   map[string]any    `json:"meta"`
}

func PersistIncident(ev Event) error {
	dir := "_incidents"
	if err := os.MkdirAll(dir, 0o755); err != nil { return err }
	fn := filepath.Join(dir, time.Now().UTC().Format("20060102T150405Z")+".json")
	f, err := os.Create(fn)
	if err != nil { return err }
	defer f.Close()
	return json.NewEncoder(f).Encode(ev)
}

func AppendLog(path string, line string) {
	if path == "" { return }
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil { log.Printf("log open: %v", err); return }
	defer f.Close()
	_, _ = f.WriteString(line+"\n")
}
