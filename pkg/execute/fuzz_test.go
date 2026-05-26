// Package execute_test provides fuzz tests for the execute package.
// Run with: go test -fuzz=FuzzCommandResult_JSONRoundtrip -fuzztime=30s ./pkg/execute/
package execute_test

import (
	"encoding/json"
	"testing"

	"install-it/pkg/execute"
)

// FuzzCommandResult_JSONRoundtrip verifies that CommandResult survives arbitrary
// JSON input (as might be received from the Wails frontend) without panicking,
// and that a valid struct survives a marshal → unmarshal cycle unchanged.
func FuzzCommandResult_JSONRoundtrip(f *testing.F) {
	// Seed corpus
	f.Add(`{"lapse":1.5,"exitCode":0,"stdout":"hello","stderr":"","error":"","aborted":false}`)
	f.Add(`{"lapse":-1,"exitCode":1,"stdout":"","stderr":"error output","error":"program not found","aborted":true}`)
	f.Add(`{}`)
	f.Add(`{"lapse":0,"exitCode":0,"stdout":"","stderr":"","error":"","aborted":false}`)
	f.Add(`{"lapse":3.14159,"exitCode":255,"stdout":"lots of output\nwith newlines","stderr":"err","error":"","aborted":false}`)
	f.Add(`not json at all`)
	f.Add(`null`)
	f.Add(`[]`)

	f.Fuzz(func(t *testing.T, s string) {
		var cr execute.CommandResult
		if err := json.Unmarshal([]byte(s), &cr); err != nil {
			// Not valid JSON for CommandResult — not a bug, skip the roundtrip check.
			return
		}

		b, err := json.Marshal(cr)
		if err != nil {
			t.Errorf("Marshal of valid CommandResult failed: %v", err)
			return
		}

		var cr2 execute.CommandResult
		if err := json.Unmarshal(b, &cr2); err != nil {
			t.Errorf("Unmarshal of re-marshaled CommandResult failed: %v", err)
			return
		}

		// All fields must survive the round-trip unchanged.
		if cr.ExitCode != cr2.ExitCode {
			t.Errorf("ExitCode: %d → %d", cr.ExitCode, cr2.ExitCode)
		}
		if cr.Error != cr2.Error {
			t.Errorf("Error: %q → %q", cr.Error, cr2.Error)
		}
		if cr.Stdout != cr2.Stdout {
			t.Errorf("Stdout: %q → %q", cr.Stdout, cr2.Stdout)
		}
		if cr.Stderr != cr2.Stderr {
			t.Errorf("Stderr: %q → %q", cr.Stderr, cr2.Stderr)
		}
		if cr.Aborted != cr2.Aborted {
			t.Errorf("Aborted: %v → %v", cr.Aborted, cr2.Aborted)
		}
		// Lapse is float32; re-encoding via float64 JSON intermediary may cause
		// tiny precision differences, so we accept a small tolerance.
		diff := cr.Lapse - cr2.Lapse
		if diff < -0.001 || diff > 0.001 {
			t.Errorf("Lapse: %f → %f (delta %f exceeds tolerance)", cr.Lapse, cr2.Lapse, diff)
		}
	})
}
