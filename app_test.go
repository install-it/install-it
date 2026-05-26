// Tests for the App struct methods in package main.
//
// Note: the init() function in main.go runs when this test binary starts. It
// parses buildVersion (empty → "0.0.0") and creates conf/ and drivers/
// directories next to the test executable. This is harmless for testing.
//
// TriggerNativeUpdate cannot be tested end-to-end because it calls os.Exit(0);
// extractBinaryFromZip (its internal ZIP helper) is tested directly below.
package main

import (
	"archive/zip"
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
	if v != "0.0.0" {
		t.Errorf("AppVersion() = %q, want '0.0.0'", v)
	}
}

func TestApp_AppBinaryType_Format(t *testing.T) {
	t.Parallel()
	a := App{}
	bt := a.AppBinaryType()

	parts := strings.SplitN(bt, "-", 2)
	if len(parts) != 2 {
		t.Fatalf("AppBinaryType() = %q: expected 'os-arch' format", bt)
	}

	goos := parts[0]
	arch := parts[1]

	if goos == "" || arch == "" {
		t.Errorf("AppBinaryType() = %q: os or arch part is empty", bt)
	}

	if goos != runtime.GOOS {
		t.Errorf("AppBinaryType() os part = %q, want %q", goos, runtime.GOOS)
	}

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

// createTestZip writes a ZIP archive at zipPath containing entries from the
// provided map (name → content).
func createTestZip(t *testing.T, zipPath string, entries map[string][]byte) {
	t.Helper()
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("createTestZip Create: %v", err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	for name, content := range entries {
		fw, err := w.Create(name)
		if err != nil {
			t.Fatalf("createTestZip entry %q: %v", name, err)
		}
		if _, err := fw.Write(content); err != nil {
			t.Fatalf("createTestZip write %q: %v", name, err)
		}
	}
}

// TestApp_TriggerNativeUpdate_ExtractBinary verifies that extractBinaryFromZip
// correctly finds and extracts install-it.exe from a flat ZIP structure.
func TestApp_TriggerNativeUpdate_ExtractBinary(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	zipPath := filepath.Join(dir, "update.zip")
	wantContent := []byte("fake binary content v2")

	createTestZip(t, zipPath, map[string][]byte{
		"install-it.exe": wantContent,
	})

	destPath := filepath.Join(dir, "install-it.exe.new")
	if err := extractBinaryFromZip(zipPath, destPath); err != nil {
		t.Fatalf("extractBinaryFromZip: %v", err)
	}

	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if string(got) != string(wantContent) {
		t.Errorf("content mismatch: got %q, want %q", got, wantContent)
	}
}

// TestApp_TriggerNativeUpdate_ExtractBinary_NestedPath verifies that
// extractBinaryFromZip finds install-it.exe even when it is in a sub-directory.
func TestApp_TriggerNativeUpdate_ExtractBinary_NestedPath(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	zipPath := filepath.Join(dir, "update.zip")
	wantContent := []byte("nested binary")

	createTestZip(t, zipPath, map[string][]byte{
		"install-it/install-it.exe": wantContent,
	})

	destPath := filepath.Join(dir, "install-it.exe.new")
	if err := extractBinaryFromZip(zipPath, destPath); err != nil {
		t.Fatalf("extractBinaryFromZip nested: %v", err)
	}

	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if string(got) != string(wantContent) {
		t.Errorf("content mismatch: got %q, want %q", got, wantContent)
	}
}

// TestApp_TriggerNativeUpdate_ExtractBinary_NotFound verifies that
// extractBinaryFromZip returns an error when the ZIP contains no install-it.exe.
func TestApp_TriggerNativeUpdate_ExtractBinary_NotFound(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	zipPath := filepath.Join(dir, "update.zip")

	createTestZip(t, zipPath, map[string][]byte{
		"readme.txt":     []byte("docs"),
		"other-tool.exe": []byte("not it"),
	})

	destPath := filepath.Join(dir, "install-it.exe.new")
	if err := extractBinaryFromZip(zipPath, destPath); err == nil {
		t.Error("expected error when install-it.exe is absent from ZIP, got nil")
	}
}

// TestApp_TriggerNativeUpdate_ExtractBinary_CaseInsensitive verifies that
// extractBinaryFromZip matches Install-It.EXE (case-insensitive).
func TestApp_TriggerNativeUpdate_ExtractBinary_CaseInsensitive(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	zipPath := filepath.Join(dir, "update.zip")
	wantContent := []byte("uppercase binary")

	createTestZip(t, zipPath, map[string][]byte{
		"Install-It.EXE": wantContent,
	})

	destPath := filepath.Join(dir, "install-it.exe.new")
	if err := extractBinaryFromZip(zipPath, destPath); err != nil {
		t.Fatalf("extractBinaryFromZip case-insensitive: %v", err)
	}

	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if string(got) != string(wantContent) {
		t.Errorf("content mismatch: got %q, want %q", got, wantContent)
	}
}
