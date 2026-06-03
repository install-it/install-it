// Package execute_test provides external black-box tests for the execute package.
//
// Design notes:
//   - CommandExecutor.RunAndOutput does not touch ce.ctx or ce.commands, so it
//     works safely without calling SetContext first.
//   - CommandExecutor.Run and Abort require SetContext (which initialises the
//     internal xsync.MapOf). Passing context.Background() is safe: the Wails
//     runtime.EventsEmit called inside dispatch is a no-op when the context
//     carries no Wails app state.
//   - All process-spawning commands use Windows-native paths (cmd.exe) because
//     this project targets Windows exclusively.
package execute_test

import (
	"context"
	"strings"
	"testing"

	"install-it/pkg/execute"
)

// ==================== RunAndOutput ====================

func TestCommandExecutor_RunAndOutput_Success(t *testing.T) {
	t.Parallel()

	var ce execute.CommandExecutor
	result := ce.RunAndOutput("cmd", []string{"/c", "echo hello"}, false)

	if result.Error != "" {
		t.Fatalf("unexpected error field: %s", result.Error)
	}
	if result.ExitCode != 0 {
		t.Errorf("exit code: got %d, want 0", result.ExitCode)
	}
	if !strings.Contains(result.Stdout, "hello") {
		t.Errorf("stdout %q should contain 'hello'", result.Stdout)
	}
}

func TestCommandExecutor_RunAndOutput_NonZeroExit(t *testing.T) {
	t.Parallel()

	var ce execute.CommandExecutor
	result := ce.RunAndOutput("cmd", []string{"/c", "exit 1"}, false)

	if result.ExitCode != 1 {
		t.Errorf("exit code: got %d, want 1", result.ExitCode)
	}
}

func TestCommandExecutor_RunAndOutput_InvalidProgram(t *testing.T) {
	t.Parallel()

	var ce execute.CommandExecutor
	result := ce.RunAndOutput("__nonexistent_binary_xyz__", []string{}, false)

	if result.Error == "" {
		t.Error("expected non-empty Error field for invalid program, got empty string")
	}
}

func TestCommandExecutor_RunAndOutput_Stderr(t *testing.T) {
	t.Parallel()

	var ce execute.CommandExecutor
	// Redirect echo output to stderr inside cmd.exe (1>&2 = stdout → stderr).
	result := ce.RunAndOutput("cmd", []string{"/c", "echo stderr_content 1>&2"}, false)

	if result.Stderr == "" {
		t.Errorf("expected non-empty Stderr, got empty string (stdout=%q, error=%q)",
			result.Stdout, result.Error)
	}
}

func TestCommandResult_Lapse(t *testing.T) {
	t.Parallel()

	var ce execute.CommandExecutor
	result := ce.RunAndOutput("cmd", []string{"/c", "echo lapse_test"}, false)

	// Lapse is in seconds; must be >= 0 for a command that ran successfully.
	if result.Lapse < 0 {
		t.Errorf("Lapse should be >= 0 for a completed command, got %f", result.Lapse)
	}
	// Sanity upper bound: the echo command should finish within 30 seconds.
	if result.Lapse > 30 {
		t.Errorf("Lapse %f seconds seems unreasonably large for echo", result.Lapse)
	}
}

// ==================== Abort ====================

func TestCommandExecutor_Abort_NonExistentId(t *testing.T) {
	t.Parallel()

	var ce execute.CommandExecutor
	ce.SetContext(context.Background())

	if err := ce.Abort("nonexistent_id"); err == nil {
		t.Error("expected error when aborting non-existent id, got nil")
	}
}

func TestCommandExecutor_Abort_RunningCommand(t *testing.T) {
	// Skipped: CommandExecutor.Run starts a goroutine that calls dispatch(),
	// which calls runtime.EventsEmit(ce.ctx, ...).  With context.Background()
	// Wails calls log.Fatal("invalid context"), crashing the entire test binary.
	// This path requires a real Wails app context and is covered by integration tests.
	t.Skip("Run() goroutine calls runtime.EventsEmit which requires a real Wails context — use integration environment")
}

// ==================== Run (event emission skipped — requires Wails runtime) ====================

func TestCommandExecutor_Run_ReturnsNonEmptyId(t *testing.T) {
	t.Skip("Run() dispatches via Wails EventsEmit; requires integration environment with Wails runtime context")
}

// TestCommandExecutor_RunWailsEvents skips tests that require a real Wails context.
func TestCommandExecutor_RunWailsEvents(t *testing.T) {
	t.Skip("Wails runtime.EventsEmit requires a real Wails app context — tested via integration tests only")
}
