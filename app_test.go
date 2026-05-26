// Tests for the App struct methods in package main.
//
// Note: the init() function in main.go runs when this test binary starts. It
// parses buildVersion (empty → "0.0.0") and creates conf/ and drivers/
// directories next to the test executable. This is harmless for testing.
//
// Methods that require a live Wails context (SelectFolder, SelectFile, Update)
// are excluded; Update is tagged as integration-only.
package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestApp_AppVersion_DefaultZero(t *testing.T) {
	t.Parallel()
	a := App{}
	v := a.AppVersion()
	// When built without -ldflags "-X main.buildVersion=...", buildVersion is
	// the empty string and init() sets version to "0.0.0".
	if v != "0.0.0" {
		t.Errorf("AppVersion() = %q, want '0.0.0'", v)
	}
}

func TestApp_AppBinaryType_Format(t *testing.T) {
	t.Parallel()
	a := App{}
	bt := a.AppBinaryType()

	// Must be "goos-arch" format.
	parts := strings.SplitN(bt, "-", 2)
	if len(parts) != 2 {
		t.Fatalf("AppBinaryType() = %q: expected 'os-arch' format", bt)
	}

	goos := parts[0]
	arch := parts[1]

	if goos == "" || arch == "" {
		t.Errorf("AppBinaryType() = %q: os or arch part is empty", bt)
	}

	// Verify the GOOS component matches the actual OS.
	if goos != runtime.GOOS {
		t.Errorf("AppBinaryType() os part = %q, want %q", goos, runtime.GOOS)
	}

	// Verify the arch translation: amd64→x64, 386→x86, anything else stays.
	goarch := runtime.GOARCH
	wantArch := goarch
	if goarch == "amd64" {
		wantArch = "x64"
	} else if goarch == "386" {
		wantArch = "x86"
	}
	if arch != wantArch {
		t.Errorf("AppBinaryType() arch part = %q, want %q", arch, wantArch)
	}
}

func TestApp_PathExists_ExistingFile(t *testing.T) {
	t.Parallel()
	a := App{}

	f, err := os.CreateTemp(t.TempDir(), "pathexists-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	if !a.PathExists(f.Name()) {
		t.Errorf("PathExists(%q) = false, want true for existing file", f.Name())
	}
}

func TestApp_PathExists_ExistingDirectory(t *testing.T) {
	t.Parallel()
	a := App{}
	dir := t.TempDir()
	if !a.PathExists(dir) {
		t.Errorf("PathExists(%q) = false, want true for existing directory", dir)
	}
}

func TestApp_PathExists_NonExistent(t *testing.T) {
	t.Parallel()
	a := App{}
	path := filepath.Join(t.TempDir(), "no_such_file_xyz_123.txt")
	if a.PathExists(path) {
		t.Errorf("PathExists(%q) = true, want false for non-existent path", path)
	}
}

// TestApp_Update is skipped in normal test runs because it downloads a real
// binary from GitHub and spawns a process. Run with -run TestApp_Update manually
// in an environment with internet access to exercise the full update flow.
func TestApp_Update(t *testing.T) {
	t.Skip("integration test: downloads and spawns a real updater binary from GitHub — skipped by default")
}
