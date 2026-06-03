package storage

import (
	"encoding/json"
	"os"
	"testing"
)

// ==================== Helpers ====================

func createTempFile(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()
	return tmpFile.Name()
}

func cleanupTempFile(t *testing.T, path string) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to cleanup temp file: %v", err)
	}
}

// ==================== FileStore Tests ====================

func TestFileStore_ReadNonExistentFile(t *testing.T) {
	store := &FileStore{Path: "/nonexistent/path/file.json"}
	var data struct{ V int }

	err := store.Read(&data)
	if err != nil {
		t.Errorf("expected nil error for non-existent file, got %v", err)
	}
}

func TestFileStore_ReadValidJSON(t *testing.T) {
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	type payload struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Val  int    `json:"val"`
	}
	testObj := payload{ID: "123", Name: "Test", Val: 42}
	bytes, _ := json.Marshal(testObj)
	os.WriteFile(path, bytes, os.ModePerm)

	store := &FileStore{Path: path}
	var result payload
	err := store.Read(&result)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if result.ID != "123" || result.Name != "Test" || result.Val != 42 {
		t.Errorf("unexpected read data: %+v", result)
	}
}

func TestFileStore_ReadInvalidJSON(t *testing.T) {
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	os.WriteFile(path, []byte("invalid json {"), os.ModePerm)

	store := &FileStore{Path: path}
	var result struct{ V int }
	err := store.Read(&result)

	if err == nil {
		t.Errorf("expected error for invalid JSON, got nil")
	}
}

func TestFileStore_WriteValidData(t *testing.T) {
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	type payload struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Val  int    `json:"val"`
	}
	testObj := payload{ID: "456", Name: "WriteTest", Val: 100}
	store := &FileStore{Path: path}

	err := store.Write(testObj)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	bytes, _ := os.ReadFile(path)
	var result payload
	json.Unmarshal(bytes, &result)

	if result.ID != "456" || result.Name != "WriteTest" || result.Val != 100 {
		t.Errorf("written data mismatch: %+v", result)
	}
}

func TestFileStore_WriteUnmarshalableData(t *testing.T) {
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	unmarshalable := make(chan int)
	store := &FileStore{Path: path}

	err := store.Write(unmarshalable)
	if err == nil {
		t.Errorf("expected error for unmarshalable data, got nil")
	}
}

func TestFileStore_ExistWhenFileExists(t *testing.T) {
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	store := &FileStore{Path: path}
	if !store.Exist() {
		t.Errorf("expected Exist() to return true for existing file")
	}
}

func TestFileStore_ExistWhenFileNotExists(t *testing.T) {
	store := &FileStore{Path: "/nonexistent/path/file.json"}
	if store.Exist() {
		t.Errorf("expected Exist() to return false for non-existent file")
	}
}

// ==================== MemoryStore Tests ====================

func TestMemoryStore_ReadEmpty(t *testing.T) {
	m := &MemoryStore{}
	var v map[string]any
	if err := m.Read(&v); err != nil {
		t.Errorf("Read on empty MemoryStore: expected nil error, got %v", err)
	}
	if v != nil {
		t.Errorf("expected nil value for empty store, got %v", v)
	}
}

func TestMemoryStore_WriteAndRead(t *testing.T) {
	type payload struct {
		X int    `json:"x"`
		Y string `json:"y"`
	}
	m := &MemoryStore{}
	want := payload{X: 7, Y: "hello"}
	if err := m.Write(want); err != nil {
		t.Fatalf("Write: %v", err)
	}
	var got payload
	if err := m.Read(&got); err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got != want {
		t.Errorf("roundtrip mismatch: got %+v, want %+v", got, want)
	}
}

func TestMemoryStore_ExistAfterWrite(t *testing.T) {
	m := &MemoryStore{}
	if m.Exist() {
		t.Error("expected Exist()=false before Write")
	}
	m.Write(struct{}{})
	if !m.Exist() {
		t.Error("expected Exist()=true after Write")
	}
}
