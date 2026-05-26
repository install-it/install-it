// Package porter_test provides external black-box tests for the porter package.
package porter_test

import (
	"archive/zip"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"install-it/pkg/porter"
	"install-it/pkg/status"
)

// ==================== Progress ====================

func TestProgress_WriteAccumulates(t *testing.T) {
	t.Parallel()

	p := &porter.Progress{}
	n1, err := p.Write([]byte("hello")) // 5 bytes
	if err != nil {
		t.Fatalf("Write 1: %v", err)
	}
	n2, err := p.Write([]byte("world!")) // 6 bytes
	if err != nil {
		t.Fatalf("Write 2: %v", err)
	}
	if n1 != 5 || n2 != 6 {
		t.Errorf("Write returned (%d, %d), want (5, 6)", n1, n2)
	}
	if p.Current != 11 {
		t.Errorf("Current after two writes: got %d, want 11", p.Current)
	}
}

func TestProgress_StartCompleteFlow(t *testing.T) {
	t.Parallel()

	p := &porter.Progress{}
	p.Start(256)

	if p.Status != status.Running {
		t.Errorf("Status after Start: got %v, want Running", p.Status)
	}
	if p.Total != 256 {
		t.Errorf("Total: got %d, want 256", p.Total)
	}
	if p.StartAt.IsZero() {
		t.Error("StartAt should be set after Start")
	}

	p.Complete()

	if p.Status != status.Completed {
		t.Errorf("Status after Complete: got %v, want Completed", p.Status)
	}
	if p.Current != p.Total {
		t.Errorf("Current should equal Total after Complete: got %d, want %d", p.Current, p.Total)
	}
}

func TestProgress_FailFlow(t *testing.T) {
	t.Parallel()

	p := &porter.Progress{}
	p.Start(100)

	sentinel := errors.New("disk error")
	p.Fail(sentinel)

	if p.Status != status.Failed {
		t.Errorf("Status after Fail(err): got %v, want Failed", p.Status)
	}
	if p.Error == nil {
		t.Error("Error field should be set after Fail")
	}
}

func TestProgress_FailWithContextCanceled_SetsAborted(t *testing.T) {
	t.Parallel()

	p := &porter.Progress{}
	p.Start(100)
	p.Fail(context.Canceled)

	if p.Status != status.Aborted {
		t.Errorf("Status after Fail(context.Canceled): got %v, want Aborted", p.Status)
	}
}

func TestProgress_AccumulateIncreasesCurrentOnly(t *testing.T) {
	t.Parallel()

	p := &porter.Progress{}
	p.Total = 100
	p.Accumulate(40)
	p.Accumulate(35)

	if p.Current != 75 {
		t.Errorf("Current after Accumulate: got %d, want 75", p.Current)
	}
	if p.Total != 100 {
		t.Error("Accumulate should not change Total")
	}
}

// ==================== Porter — idle state ====================

func TestPorter_StatusIdle(t *testing.T) {
	t.Parallel()

	p := &porter.Porter{}
	if p.Status() != status.Pending {
		t.Errorf("fresh Porter.Status(): got %v, want Pending", p.Status())
	}
}

func TestPorter_AbortNoJob(t *testing.T) {
	t.Parallel()

	p := &porter.Porter{}
	if err := p.Abort(); err == nil {
		t.Error("expected error from Abort with no running job, got nil")
	}
}

func TestPorter_ProgressNoJob(t *testing.T) {
	t.Parallel()

	p := &porter.Porter{}
	_, err := p.Progress()
	if err == nil {
		t.Error("expected error from Progress with no started job, got nil")
	}
}

// ==================== Porter — ImportFromFile ====================

func TestPorter_ImportFromFileInvalidPath(t *testing.T) {
	t.Parallel()

	msg := make(chan string, 512)
	p := &porter.Porter{
		DirRoot: t.TempDir(),
		Message: msg,
	}

	err := p.ImportFromFile(filepath.Join(t.TempDir(), "nonexistent_completely_random.zip"))
	if err == nil {
		t.Error("expected error for non-existent zip path, got nil")
	}
}

// TestPorter_ImportFromFile_BasicExtraction creates a minimal zip file by hand,
// imports it, and verifies the file lands in the correct location.
func TestPorter_ImportFromFile_BasicExtraction(t *testing.T) {
	// Not parallel — writes real files; isolated via TempDir.
	base := t.TempDir()
	importConf := filepath.Join(base, "conf")
	if err := os.MkdirAll(importConf, 0755); err != nil {
		t.Fatal(err)
	}

	// Craft a zip with a single entry: conf/data.json
	zipPath := filepath.Join(base, "test.zip")
	if err := createTestZip(t, zipPath, map[string][]byte{
		filepath.Join("conf", "data.json"): []byte(`{"ok":true}`),
	}); err != nil {
		t.Fatalf("createTestZip: %v", err)
	}

	msg := make(chan string, 512)
	p := &porter.Porter{
		DirRoot: base,
		Targets: []string{importConf}, // will be renamed → conf_old, then removed
		Message: msg,
	}

	if err := p.ImportFromFile(zipPath); err != nil {
		t.Fatalf("ImportFromFile: %v", err)
	}

	// The extracted file should exist in the import destination.
	extracted := filepath.Join(base, "conf", "data.json")
	data, err := os.ReadFile(extracted)
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if string(data) != `{"ok":true}` {
		t.Errorf("content mismatch: got %q, want '{\"ok\":true}'", string(data))
	}
}

// ==================== Porter — Export ====================

// testChdir temporarily changes the working directory and returns a cleanup func
// that restores it. It must be called before any t.TempDir() calls so that the
// cwd-restore cleanup runs FIRST in LIFO order (before TempDir removal), which
// is required on Windows where a directory cannot be deleted while it is the cwd.
func testChdirSetup(t *testing.T, dir string) {
	t.Helper()
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir to %s: %v", dir, err)
	}
	// Register AFTER Chdir so cleanup runs BEFORE any TempDir cleanups (LIFO).
	t.Cleanup(func() {
		if err := os.Chdir(origCwd); err != nil {
			t.Logf("warning: could not restore cwd to %s: %v", origCwd, err)
		}
	})
}

func TestPorter_ExportCreatesZip(t *testing.T) {
	// Not parallel — temporarily changes the working directory.

	// Create base BEFORE registering the chdir cleanup so that the cwd-restore
	// runs FIRST (LIFO) and base can be removed afterwards.
	base, err := os.MkdirTemp("", "TestPorter_ExportCreatesZip-*")
	if err != nil {
		t.Fatal(err)
	}
	origCwd, err := os.Getwd()
	if err != nil {
		os.RemoveAll(base)
		t.Fatal(err)
	}
	if err := os.Chdir(base); err != nil {
		os.RemoveAll(base)
		t.Fatalf("Chdir to %s: %v", base, err)
	}
	t.Cleanup(func() {
		os.Chdir(origCwd)
		os.RemoveAll(base)
	})

	// Create source directory with a test file.
	confDir := filepath.Join(base, "conf")
	if err := os.MkdirAll(confDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(confDir, "setting.json"), []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	exportDest := t.TempDir()
	msg := make(chan string, 512)
	p := &porter.Porter{
		DirRoot: base,
		Targets: []string{confDir},
		Message: msg,
	}

	if err := p.Export(exportDest); err != nil {
		t.Fatalf("Export: %v", err)
	}

	zipPath := filepath.Join(exportDest, "install-it.zip")
	if _, err := os.Stat(zipPath); err != nil {
		t.Fatalf("install-it.zip not found after Export: %v", err)
	}

	// Verify at least one entry in the zip.
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer r.Close()
	if len(r.File) == 0 {
		t.Error("export produced an empty zip archive")
	}
}

func TestPorter_ExportCreatesProgressSteps(t *testing.T) {
	// Not parallel — temporarily changes the working directory.

	base, err := os.MkdirTemp("", "TestPorter_ExportProgress-*")
	if err != nil {
		t.Fatal(err)
	}
	origCwd, err := os.Getwd()
	if err != nil {
		os.RemoveAll(base)
		t.Fatal(err)
	}
	if err := os.Chdir(base); err != nil {
		os.RemoveAll(base)
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() {
		os.Chdir(origCwd)
		os.RemoveAll(base)
	})

	srcDir := filepath.Join(base, "drivers")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	exportDest := t.TempDir()
	msg := make(chan string, 512)
	p := &porter.Porter{
		DirRoot: base,
		Targets: []string{srcDir},
		Message: msg,
	}

	if err := p.Export(exportDest); err != nil {
		t.Fatalf("Export: %v", err)
	}

	progress, err := p.Progress()
	if err != nil {
		t.Fatalf("Progress: %v", err)
	}

	if len(progress.Progresses) != 2 {
		t.Fatalf("expected 2 progress steps, got %d", len(progress.Progresses))
	}

	names := make(map[string]bool)
	for _, pr := range progress.Progresses {
		names[pr.Name] = true
	}
	for _, want := range []string{"initialisation", "compression"} {
		if !names[want] {
			t.Errorf("missing progress step %q; got steps: %v", want, progress.Progresses)
		}
	}
}

// ==================== Full Export → Import roundtrip ====================

func TestPorter_ExportImportRoundtrip(t *testing.T) {
	// Not parallel — temporarily changes the working directory.

	base, err := os.MkdirTemp("", "TestPorter_Roundtrip-*")
	if err != nil {
		t.Fatal(err)
	}
	origCwd, err := os.Getwd()
	if err != nil {
		os.RemoveAll(base)
		t.Fatal(err)
	}
	if err := os.Chdir(base); err != nil {
		os.RemoveAll(base)
		t.Fatalf("Chdir to export base: %v", err)
	}
	t.Cleanup(func() {
		os.Chdir(origCwd)
		os.RemoveAll(base)
	})

	// ---- Export phase ----
	confDir := filepath.Join(base, "conf")
	if err := os.MkdirAll(confDir, 0755); err != nil {
		t.Fatal(err)
	}

	const wantContent = `{"version":"1.0"}`
	if err := os.WriteFile(filepath.Join(confDir, "config.json"), []byte(wantContent), 0644); err != nil {
		t.Fatal(err)
	}

	exportDest := t.TempDir()
	exportMsg := make(chan string, 512)
	exporter := &porter.Porter{
		DirRoot: base,
		Targets: []string{confDir},
		Message: exportMsg,
	}
	if err := exporter.Export(exportDest); err != nil {
		t.Fatalf("Export: %v", err)
	}

	zipPath := filepath.Join(exportDest, "install-it.zip")
	if _, err := os.Stat(zipPath); err != nil {
		t.Fatalf("zip not created: %v", err)
	}

	// ---- Import phase ----
	importBase := t.TempDir()
	importConf := filepath.Join(importBase, "conf")
	if err := os.MkdirAll(importConf, 0755); err != nil {
		t.Fatal(err)
	}

	importMsg := make(chan string, 512)
	importer := &porter.Porter{
		DirRoot: importBase,
		Targets: []string{importConf},
		Message: importMsg,
	}
	if err := importer.ImportFromFile(zipPath); err != nil {
		t.Fatalf("ImportFromFile: %v", err)
	}

	// Verify the extracted file
	extracted := filepath.Join(importBase, "conf", "config.json")
	got, err := os.ReadFile(extracted)
	if err != nil {
		t.Fatalf("read extracted config.json: %v", err)
	}
	if string(got) != wantContent {
		t.Errorf("content mismatch:\n  got  %q\n  want %q", string(got), wantContent)
	}
}

// ==================== helpers ====================

// createTestZip creates a zip archive at path containing the given files (name → content).
func createTestZip(t *testing.T, path string, files map[string][]byte) error {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	for name, content := range files {
		entry, err := w.Create(name)
		if err != nil {
			return err
		}
		if _, err := entry.Write(content); err != nil {
			return err
		}
	}
	return nil
}
