// Package storage_test provides external black-box tests for the storage package.
// These complement the internal tests in store_test.go, app_setting_test.go, and driver_test.go.
package storage_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gorm.io/gorm"

	"install-it/pkg/storage"
)

// openExternalTestDB creates an isolated SQLite database in the test's temp directory.
// Used by all external (package storage_test) tests that need a DB.
func openExternalTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := storage.OpenDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("openExternalTestDB: %v", err)
	}
	if err := storage.RunMigrations(db); err != nil {
		t.Fatalf("openExternalTestDB RunMigrations: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err != nil {
			t.Logf("failed to get underlying SQL DB: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			t.Logf("failed to close database: %v", err)
		}
	})
	return db
}

// containsUint reports whether slice contains v.
func containsUint(slice []uint, v uint) bool {
	for _, u := range slice {
		if u == v {
			return true
		}
	}
	return false
}

// addTestGroup adds a group and returns its autoincrement ID via All().
func addTestGroup(t *testing.T, dgs *storage.DriverGroupStorage, group storage.DriverGroup) uint {
	t.Helper()
	if err := dgs.Add(group); err != nil {
		t.Fatalf("addTestGroup Add: %v", err)
	}
	all, err := dgs.All()
	if err != nil {
		t.Fatalf("addTestGroup All: %v", err)
	}
	return all[len(all)-1].Id
}

// ==================== FileStore (external API) ====================

func TestFileStore_ReadWriteRoundtrip(t *testing.T) {
	t.Parallel()

	type payload struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}

	dir := t.TempDir()
	store := &storage.FileStore{Path: filepath.Join(dir, "data.json")}

	want := payload{Value: 42, Label: "roundtrip"}
	if err := store.Write(want); err != nil {
		t.Fatalf("Write: %v", err)
	}

	var got payload
	if err := store.Read(&got); err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got != want {
		t.Errorf("roundtrip mismatch: got %+v, want %+v", got, want)
	}
}

func TestFileStore_ReadNonExistent_ReturnsNilError(t *testing.T) {
	t.Parallel()

	store := &storage.FileStore{Path: filepath.Join(t.TempDir(), "nosuchfile.json")}
	var v map[string]any
	if err := store.Read(&v); err != nil {
		t.Errorf("Read on non-existent file: expected nil error, got %v", err)
	}
	if v != nil {
		t.Errorf("expected nil value for non-existent file, got %v", v)
	}
}

func TestFileStore_Exist_BeforeAndAfterWrite(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "newfile.json")
	store := &storage.FileStore{Path: path}

	if store.Exist() {
		t.Error("expected Exist()=false before any Write")
	}

	if err := store.Write(struct{}{}); err != nil {
		t.Fatalf("Write: %v", err)
	}

	if !store.Exist() {
		t.Error("expected Exist()=true after Write")
	}
}

func TestFileStore_WriteCreatesFileInExistingDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	subDir := filepath.Join(dir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(subDir, "output.json")

	store := &storage.FileStore{Path: path}
	if err := store.Write(map[string]int{"answer": 42}); err != nil {
		t.Fatalf("Write to nested path: %v", err)
	}
	if !store.Exist() {
		t.Error("expected file to exist after write to pre-created sub-directory")
	}
}

func TestFileStore_CorruptedJSON_ReturnsError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.json")
	if err := os.WriteFile(path, []byte("{not valid json"), 0644); err != nil {
		t.Fatal(err)
	}

	store := &storage.FileStore{Path: path}
	var v map[string]any
	if err := store.Read(&v); err == nil {
		t.Error("expected error when reading corrupted JSON, got nil")
	}
}

// ==================== AppSettingStorage (external API) ====================

func TestAppSettingStorage_AllReturnsDefault(t *testing.T) {
	t.Parallel()

	s := &storage.AppSettingStorage{Store: &storage.MemoryStore{}}
	result, err := s.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if result.ParallelInstall != true {
		t.Errorf("default ParallelInstall: got %v, want true", result.ParallelInstall)
	}
	if result.Language != "en" {
		t.Errorf("default Language: got %q, want 'en'", result.Language)
	}
	if result.SuccessAction != storage.Nothing {
		t.Errorf("default SuccessAction: got %v, want Nothing", result.SuccessAction)
	}
	if result.SuccessActionDelay != 5 {
		t.Errorf("default SuccessActionDelay: got %d, want 5", result.SuccessActionDelay)
	}
}

func TestAppSettingStorage_UpdatePersists(t *testing.T) {
	t.Parallel()

	s := &storage.AppSettingStorage{Store: &storage.MemoryStore{}}
	want := storage.AppSetting{
		Language:        "zh_Hant_HK",
		ParallelInstall: false,
		AutoCheckUpdate: false,
		SuccessAction:   storage.Reboot,
		Password:        "s3cr3t",
	}
	if _, err := s.Update(want); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err := s.All()
	if err != nil {
		t.Fatalf("All after Update: %v", err)
	}
	if got.Language != want.Language {
		t.Errorf("Language: got %q, want %q", got.Language, want.Language)
	}
	if got.ParallelInstall != want.ParallelInstall {
		t.Errorf("ParallelInstall: got %v, want %v", got.ParallelInstall, want.ParallelInstall)
	}
	if got.SuccessAction != want.SuccessAction {
		t.Errorf("SuccessAction: got %v, want %v", got.SuccessAction, want.SuccessAction)
	}
}

func TestAppSettingStorage_UpdateReturnsSaved(t *testing.T) {
	t.Parallel()

	s := &storage.AppSettingStorage{Store: &storage.MemoryStore{}}
	input := storage.AppSetting{
		Language:           "en",
		SuccessAction:      storage.Shutdown,
		SuccessActionDelay: 30,
		HideNotFound:       true,
	}
	returned, err := s.Update(input)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if returned.Language != input.Language {
		t.Errorf("returned Language: got %q, want %q", returned.Language, input.Language)
	}
	if returned.SuccessAction != input.SuccessAction {
		t.Errorf("returned SuccessAction: got %v, want %v", returned.SuccessAction, input.SuccessAction)
	}
	if returned.HideNotFound != input.HideNotFound {
		t.Errorf("returned HideNotFound: got %v, want %v", returned.HideNotFound, input.HideNotFound)
	}

	persisted, err := s.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if persisted.Language != input.Language {
		t.Errorf("persisted Language: got %q, want %q", persisted.Language, input.Language)
	}
}

// ==================== AppSetting JSON roundtrip ====================

func TestAppSetting_JSONRoundtrip(t *testing.T) {
	t.Parallel()

	original := storage.AppSetting{
		CreatePartition:    true,
		SetPassword:        true,
		Password:           "s3cr3t",
		ParallelInstall:    false,
		SuccessAction:      storage.Reboot,
		SuccessActionDelay: 10,
		FilterMiniportNic:  true,
		FilterMicrosoftNic: false,
		Language:           "zh_Hant_HK",
		DriverDownloadUrl:  "https://example.com/drivers",
		AutoCheckUpdate:    false,
		HideNotFound:       true,
	}

	b, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var decoded storage.AppSetting
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if decoded != original {
		t.Errorf("roundtrip mismatch:\n  got  %+v\n  want %+v", decoded, original)
	}
}
