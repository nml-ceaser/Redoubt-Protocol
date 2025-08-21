package isolation

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"syscall"

	"redoubt-protocol/internal/policy"
)

// HandleLinux executes actions for linux as per policy.
// It respects pol.Simulate: when simulate==true it only logs what it *would* do.
func HandleLinux(actions []string, pol *policy.Policy) {
	for _, a := range actions {
		log.Printf("[Isolation][linux] action=%s (simulate=%v)", a, pol.Simulate)
		switch a {
		case "block_network":
			if pol.Simulate {
				log.Println("[Isolation][linux] (simulate) would create iptables drop chain")
			} else {
				if err := applyIptablesDrop(); err != nil {
					log.Printf("[Isolation][linux] applyIptablesDrop error: %v", err)
				} else {
					log.Println("[Isolation][linux] iptables drop applied")
				}
			}
		case "freeze_process":
			// check control flags
			controls := pol.Controls["linux"]
			targetPid := 0
			if v, ok := controls["target_pid"]; ok {
				switch t := v.(type) {
				case int:
					targetPid = t
				case int64:
					targetPid = int(t)
				case float64:
					targetPid = int(t)
				}
			}
			pauseSelf := false
			if v, ok := controls["pause_self"]; ok {
				if b, ok2 := v.(bool); ok2 {
					pauseSelf = b
				}
			}

			if pol.Simulate {
				if targetPid > 0 {
					log.Printf("[Isolation][linux] (simulate) would SIGSTOP pid=%d", targetPid)
				} else if pauseSelf {
					log.Printf("[Isolation][linux] (simulate) would SIGSTOP self")
				} else {
					log.Printf("[Isolation][linux] (simulate) would SIGSTOP target_pid (none set); no-op")
				}
			} else {
				if targetPid > 0 {
					if err := syscall.Kill(targetPid, syscall.SIGSTOP); err != nil {
						log.Printf("[Isolation][linux] failed SIGSTOP pid=%d: %v", targetPid, err)
					} else {
						log.Printf("[Isolation][linux] SIGSTOP pid=%d", targetPid)
					}
				} else if pauseSelf {
					pid := syscall.Getpid()
					log.Printf("[Isolation][linux] pausing self pid=%d", pid)
					if err := syscall.Kill(pid, syscall.SIGSTOP); err != nil {
						log.Printf("[Isolation][linux] failed to SIGSTOP self: %v", err)
					}
				} else {
					log.Printf("[Isolation][linux] no target_pid and pause_self==false; skip freeze_process")
				}
			}
		default:
			log.Printf("[Isolation][linux] unknown action: %s", a)
		}
	}
}

func applyIptablesDrop() error {
	// Creates a dedicated chain REDOUBT, insert jump rules for OUTPUT and INPUT if not present.
	// This implementation is conservative and idempotent.
	cmds := [][]string{
		{"iptables", "-N", "REDOUBT"},                 // may fail if exists
		{"iptables", "-F", "REDOUBT"},                 // flush chain
		{"iptables", "-C", "OUTPUT", "-j", "REDOUBT"}, // check rule
	}
	// run first commands ignoring specific errors
	for _, c := range cmds[:2] {
		if err := exec.Command(c[0], c[1:]...).Run(); err != nil {
			// ignore errors for create/flush (e.g., already exists)
		}
	}
	// ensure OUTPUT jumps to REDOUBT
	if err := exec.Command("iptables", "-C", "OUTPUT", "-j", "REDOUBT").Run(); err != nil {
		// not present → insert
		if err2 := exec.Command("iptables", "-I", "OUTPUT", "1", "-j", "REDOUBT").Run(); err2 != nil {
			return fmt.Errorf("failed to insert OUTPUT REDOUBT: %w", err2)
		}
	}

	// In chain: accept loopback then drop all
	// Clear existing chain then append accept/drop
	_ = exec.Command("iptables", "-F", "REDOUBT").Run()
	if err := exec.Command("iptables", "-A", "REDOUBT", "-o", "lo", "-j", "ACCEPT").Run(); err != nil {
		return err
	}
	if err := exec.Command("iptables", "-A", "REDOUBT", "-j", "DROP").Run(); err != nil {
		return err
	}
	return nil
}
