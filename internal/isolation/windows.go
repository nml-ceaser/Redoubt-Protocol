package isolation

import (
	"log"
	"os/exec"

	"redoubt-protocol/internal/policy"
)

// HandleWindows executes Windows-specific actions via PowerShell/Netsh.
// This file assumes it's running on Windows with necessary privileges.
func HandleWindows(actions []string, pol *policy.Policy) {
	for _, a := range actions {
		log.Printf("[Isolation][windows] action=%s (simulate=%v)", a, pol.Simulate)
		switch a {
		case "disable_interface":
			if pol.Simulate {
				log.Println("[Isolation][windows] (simulate) would disable network interface (via netsh)")
			} else {
				// Attempt to use the interface name from controls if provided
				ifName := ""
				if controls, ok := pol.Controls["windows"]; ok {
					if v, ok2 := controls["interface_name"].(string); ok2 {
						ifName = v
					}
				}
				if ifName == "" {
					log.Println("[Isolation][windows] no interface_name configured; attempting to disable all non-loopback adapters (may require admin)")
					// A conservative approach: use PowerShell to disable adapters except loopback
					powershell := `Get-NetAdapter -Physical | Where-Object { $_.Status -eq 'Up' -and $_.InterfaceDescription -notlike '*Loopback*' } | Disable-NetAdapter -Confirm:$false`
					cmd := exec.Command("powershell", "-NoProfile", "-Command", powershell)
					if out, err := cmd.CombinedOutput(); err != nil {
						log.Printf("[Isolation][windows] error disabling adapters: %v, out=%s", err, string(out))
					} else {
						log.Printf("[Isolation][windows] adapters disabled, out=%s", string(out))
					}
				} else {
					cmd := exec.Command("netsh", "interface", "set", "interface", ifName, "admin=disable")
					if out, err := cmd.CombinedOutput(); err != nil {
						log.Printf("[Isolation][windows] netsh disable error: %v, out=%s", err, string(out))
					} else {
						log.Printf("[Isolation][windows] interface %s disabled: %s", ifName, string(out))
					}
				}
			}
		default:
			log.Printf("[Isolation][windows] unknown action: %s", a)
		}
	}
}
