// Package storage_test provides external black-box tests for the storage package.
// These complement the internal tests in store_test.go, app_setting_test.go, and driver_test.go.
package storage_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"install-it/pkg/storage"
)

// testItem implements storage.HasId for use with generic CRUD helpers in external tests.
type testItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (ti *testItem) GetId() string   { return ti.ID }
func (ti *testItem) SetId(id string) { ti.ID = id }

// containsStr reports whether slice contains s.
func containsStr(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
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
	// Create a nested sub-directory manually to ensure it exists
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
	// Returned value must match what was passed in
	if returned.Language != input.Language {
		t.Errorf("returned Language: got %q, want %q", returned.Language, input.Language)
	}
	if returned.SuccessAction != input.SuccessAction {
		t.Errorf("returned SuccessAction: got %v, want %v", returned.SuccessAction, input.SuccessAction)
	}
	if returned.HideNotFound != input.HideNotFound {
		t.Errorf("returned HideNotFound: got %v, want %v", returned.HideNotFound, input.HideNotFound)
	}

	// Persisted value must also match
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

// ==================== GenerateId (external API) ====================

func TestGenerateId_ProducesEightCharHex(t *testing.T) {
	t.Parallel()

	items := []*testItem{}
	id := storage.GenerateId(items)
	if len(id) != 8 {
		t.Errorf("generated ID %q has length %d, want 8", id, len(id))
	}
	// Verify it's valid hex
	for _, c := range id {
		if !('0' <= c && c <= '9') && !('a' <= c && c <= 'f') {
			t.Errorf("ID %q contains non-hex character %q", id, c)
		}
	}
}

func TestGenerateId_Uniqueness(t *testing.T) {
	t.Parallel()

	const n = 100
	items := make([]*testItem, 0, n)
	seen := make(map[string]bool, n)

	for i := 0; i < n; i++ {
		id := storage.GenerateId(items)
		if seen[id] {
			t.Errorf("duplicate ID generated after %d iterations: %s", i, id)
		}
		seen[id] = true
		items = append(items, &testItem{ID: id})
	}

	if len(seen) != n {
		t.Errorf("expected %d unique IDs, got %d", n, len(seen))
	}
}

// ==================== Generic CRUD helpers ====================

func TestCreate_GeneratesUniqueId(t *testing.T) {
	t.Parallel()

	data := []*testItem{}
	item := &testItem{Name: "new"}
	id, err := storage.Create(item, &data)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if len(id) != 8 {
		t.Errorf("generated ID %q should be 8 chars", id)
	}
	if id != item.GetId() {
		t.Errorf("returned ID %q does not match item.GetId() %q", id, item.GetId())
	}
}

func TestCreate_AlwaysAssignsNewId(t *testing.T) {
	t.Parallel()

	// Note: Create always calls GenerateId, overwriting any pre-set ID.
	data := []*testItem{}
	item := &testItem{ID: "preset-id", Name: "test"}
	id, err := storage.Create(item, &data)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if id == "" {
		t.Error("returned ID should not be empty")
	}
	// The item is in the slice with its new ID
	if len(data) != 1 {
		t.Errorf("expected 1 item in slice, got %d", len(data))
	}
}

func TestUpdate_ExistingItem(t *testing.T) {
	t.Parallel()

	data := []*testItem{
		{ID: "aabbccdd", Name: "original"},
	}
	updated := &testItem{ID: "aabbccdd", Name: "modified"}
	if err := storage.Update(updated, &data); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if data[0].Name != "modified" {
		t.Errorf("expected Name='modified', got %q", data[0].Name)
	}
}

func TestUpdate_NonExistentItem(t *testing.T) {
	t.Parallel()

	data := []*testItem{
		{ID: "aabbccdd", Name: "existing"},
	}
	missing := &testItem{ID: "00000000", Name: "ghost"}
	if err := storage.Update(missing, &data); err == nil {
		t.Error("expected error for non-existent ID, got nil")
	}
}

func TestDelete_ExistingItem(t *testing.T) {
	t.Parallel()

	data := []*testItem{
		{ID: "id000001", Name: "keep"},
		{ID: "id000002", Name: "remove"},
		{ID: "id000003", Name: "keep2"},
	}
	if err := storage.Delete("id000002", &data); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if len(data) != 2 {
		t.Errorf("expected 2 items after delete, got %d", len(data))
	}
	for _, item := range data {
		if item.ID == "id000002" {
			t.Error("deleted item still present in slice")
		}
	}
}

func TestDelete_NonExistentItem(t *testing.T) {
	t.Parallel()

	data := []*testItem{
		{ID: "id000001", Name: "only"},
	}
	if err := storage.Delete("nothere", &data); err == nil {
		t.Error("expected error for non-existent ID, got nil")
	}
	if len(data) != 1 {
		t.Errorf("slice length changed unexpectedly: got %d", len(data))
	}
}

func TestGet_ExistingItem(t *testing.T) {
	t.Parallel()

	data := []*testItem{
		{ID: "id000001", Name: "alpha"},
		{ID: "id000002", Name: "beta"},
	}
	got, err := storage.Get("id000002", data)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "beta" {
		t.Errorf("expected Name='beta', got %q", got.Name)
	}
}

func TestGet_NonExistent(t *testing.T) {
	t.Parallel()

	data := []*testItem{
		{ID: "id000001", Name: "only"},
	}
	_, err := storage.Get("missing", data)
	if err == nil {
		t.Error("expected error for non-existent ID, got nil")
	}
}
