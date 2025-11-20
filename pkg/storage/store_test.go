package storage

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"
)

// MockObject implements HasId interface for testing
type MockObject struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (m *MockObject) GetId() string {
	if m == nil {
		return ""
	}
	return m.ID
}

func (m *MockObject) SetId(id string) {
	if m != nil {
		m.ID = id
	}
}

// Helper function to create a temporary file for testing
func createTempFile(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()
	return tmpFile.Name()
}

// Helper function to clean up temp file
func cleanupTempFile(t *testing.T, path string) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to cleanup temp file: %v", err)
	}
}

// ==================== FileStore Tests ====================

// func TestFileStore_ReadNonExistentFile(t *testing.T) {
// 	store := &FileStore{Path: "/nonexistent/path/file.json"}
// 	var data MockObject

// 	err := store.Read(&data)
// 	if err != nil {
// 		t.Errorf("expected nil error for non-existent file, got %v", err)
// 	}
// }

func TestFileStore_ReadValidJSON(t *testing.T) {
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	// Write test data
	testObj := MockObject{ID: "123", Name: "Test", Value: 42}
	bytes, _ := json.Marshal(testObj)
	os.WriteFile(path, bytes, os.ModePerm)

	// Read test data
	store := &FileStore{Path: path}
	var result MockObject
	err := store.Read(&result)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if result.ID != "123" || result.Name != "Test" || result.Value != 42 {
		t.Errorf("unexpected read data: %+v", result)
	}
}

func TestFileStore_ReadInvalidJSON(t *testing.T) {
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	// Write invalid JSON
	os.WriteFile(path, []byte("invalid json {"), os.ModePerm)

	store := &FileStore{Path: path}
	var result MockObject
	err := store.Read(&result)

	if err == nil {
		t.Errorf("expected error for invalid JSON, got nil")
	}
}

func TestFileStore_WriteValidData(t *testing.T) {
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	testObj := MockObject{ID: "456", Name: "WriteTest", Value: 100}
	store := &FileStore{Path: path}

	err := store.Write(testObj)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Verify written file
	bytes, _ := os.ReadFile(path)
	var result MockObject
	json.Unmarshal(bytes, &result)

	if result.ID != "456" || result.Name != "WriteTest" || result.Value != 100 {
		t.Errorf("written data mismatch: %+v", result)
	}

	// Verify stat is updated
	if store.stat == nil {
		t.Errorf("expected stat to be set after write")
	}
}

func TestFileStore_WriteUnmarshalableData(t *testing.T) {
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	// Create a value that can't be marshaled (e.g., channel)
	unmarshalable := make(chan int)
	store := &FileStore{Path: path}

	err := store.Write(unmarshalable)
	if err == nil {
		t.Errorf("expected error for unmarshalable data, got nil")
	}
}

func TestFileStore_ModifiedWhenStatIsNil(t *testing.T) {
	store := &FileStore{Path: "/some/path", stat: nil}

	if !store.Modified() {
		t.Errorf("expected Modified() to return true when stat is nil")
	}
}

func TestFileStore_ModifiedWhenFileNotExists(t *testing.T) {
	// Create a temporary file first to get a valid FileInfo
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	oldStat, _ := os.Stat(path)

	// Now delete it and test
	os.Remove(path)

	store := &FileStore{Path: path, stat: oldStat}

	if store.Modified() {
		t.Errorf("expected Modified() to return false when file doesn't exist")
	}
}

func TestFileStore_ModifiedWhenFileChanged(t *testing.T) {
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	// Write initial file
	os.WriteFile(path, []byte("initial"), os.ModePerm)
	store := &FileStore{Path: path}
	store.stat, _ = os.Stat(path)

	// Wait a bit and modify file
	time.Sleep(10 * time.Millisecond)
	os.WriteFile(path, []byte("modified"), os.ModePerm)

	if !store.Modified() {
		t.Errorf("expected Modified() to return true after file modification")
	}
}

func TestFileStore_ModifiedWhenFileNotChanged(t *testing.T) {
	path := createTempFile(t)
	defer cleanupTempFile(t, path)

	os.WriteFile(path, []byte("content"), os.ModePerm)
	store := &FileStore{Path: path}
	store.stat, _ = os.Stat(path)

	if store.Modified() {
		t.Errorf("expected Modified() to return false when file not modified")
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

// ==================== GenerateId Tests ====================

func TestGenerateId_UniqueIds(t *testing.T) {
	data := []*MockObject{
		{ID: "id1", Name: "obj1"},
		{ID: "id2", Name: "obj2"},
		{ID: "id3", Name: "obj3"},
	}

	id := GenerateId(data)

	// Check that id is not in the data
	for _, obj := range data {
		if obj.ID == id {
			t.Errorf("generated ID %s already exists in data", id)
		}
	}

	// Check that id has valid format (hex string with length 8, from 4 bytes)
	if len(id) != 8 {
		t.Errorf("expected ID length 8, got %d", len(id))
	}
}

func TestGenerateId_EmptyData(t *testing.T) {
	data := []*MockObject{}

	id := GenerateId(data)

	if id == "" {
		t.Errorf("expected non-empty ID for empty data")
	}
	if len(id) != 8 {
		t.Errorf("expected ID length 8, got %d", len(id))
	}
}

// ==================== IndexOf Tests ====================

func TestIndexOf_ExistingId(t *testing.T) {
	data := []*MockObject{
		{ID: "id1", Name: "obj1"},
		{ID: "id2", Name: "obj2"},
		{ID: "id3", Name: "obj3"},
	}

	index, err := IndexOf("id2", data)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if index != 1 {
		t.Errorf("expected index 1, got %d", index)
	}
}

func TestIndexOf_NonExistingId(t *testing.T) {
	data := []*MockObject{
		{ID: "id1", Name: "obj1"},
		{ID: "id2", Name: "obj2"},
	}

	index, err := IndexOf("nonexistent", data)

	if err == nil {
		t.Errorf("expected error for non-existent ID, got nil")
	}
	if index != -1 {
		t.Errorf("expected index -1, got %d", index)
	}
}

func TestIndexOf_EmptyData(t *testing.T) {
	data := []*MockObject{}

	index, err := IndexOf("id1", data)

	if err == nil {
		t.Errorf("expected error for empty data, got nil")
	}
	if index != -1 {
		t.Errorf("expected index -1, got %d", index)
	}
}

// ==================== Create Tests ====================

func TestCreate_AddsItemToSlice(t *testing.T) {
	data := []*MockObject{
		{ID: "id1", Name: "obj1"},
	}

	newObj := &MockObject{Name: "newobj", Value: 99}
	id, err := Create(newObj, &data)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(data) != 2 {
		t.Errorf("expected slice length 2, got %d", len(data))
	}
	if data[1].ID == "" {
		t.Errorf("expected ID to be set on new object")
	}
	if data[0].ID == data[1].ID {
		t.Errorf("expected different IDs, got same: %s", data[1].ID)
	}
	if id != data[1].ID {
		t.Errorf("expected returned ID %s to match object ID %s", id, data[1].ID)
	}
}

func TestCreate_EmptySlice(t *testing.T) {
	data := []*MockObject{}

	newObj := &MockObject{Name: "first"}
	id, err := Create(newObj, &data)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(data) != 1 {
		t.Errorf("expected slice length 1, got %d", len(data))
	}
	if id == "" {
		t.Errorf("expected non-empty ID")
	}
}

// ==================== Update Tests ====================

func TestUpdate_UpdatesExistingItem(t *testing.T) {
	data := []*MockObject{
		{ID: "id1", Name: "obj1", Value: 10},
		{ID: "id2", Name: "obj2", Value: 20},
		{ID: "id3", Name: "obj3", Value: 30},
	}

	updated := &MockObject{ID: "id2", Name: "updated", Value: 200}
	err := Update(updated, &data)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if data[1].Name != "updated" || data[1].Value != 200 {
		t.Errorf("expected object to be updated: %+v", data[1])
	}
	if data[0].Name != "obj1" || data[2].Name != "obj3" {
		t.Errorf("expected other items to remain unchanged")
	}
}

func TestUpdate_NonExistingItem(t *testing.T) {
	data := []*MockObject{
		{ID: "id1", Name: "obj1"},
	}

	updated := &MockObject{ID: "nonexistent", Name: "updated"}
	err := Update(updated, &data)

	if err == nil {
		t.Errorf("expected error for non-existent ID, got nil")
	}
}

func TestUpdate_EmptySlice(t *testing.T) {
	data := []*MockObject{}

	updated := &MockObject{ID: "id1", Name: "obj1"}
	err := Update(updated, &data)

	if err == nil {
		t.Errorf("expected error for empty slice, got nil")
	}
}

// ==================== Delete Tests ====================

func TestDelete_RemovesExistingItem(t *testing.T) {
	data := []*MockObject{
		{ID: "id1", Name: "obj1"},
		{ID: "id2", Name: "obj2"},
		{ID: "id3", Name: "obj3"},
	}

	err := Delete("id2", &data)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(data) != 2 {
		t.Errorf("expected slice length 2, got %d", len(data))
	}
	if data[1].ID != "id3" {
		t.Errorf("expected id3 to be at index 1 after deletion")
	}
}

func TestDelete_NonExistingItem(t *testing.T) {
	data := []*MockObject{
		{ID: "id1", Name: "obj1"},
	}

	err := Delete("nonexistent", &data)

	if err == nil {
		t.Errorf("expected error for non-existent ID, got nil")
	}
	if len(data) != 1 {
		t.Errorf("expected slice length unchanged at 1, got %d", len(data))
	}
}

func TestDelete_EmptySlice(t *testing.T) {
	data := []*MockObject{}

	err := Delete("id1", &data)

	if err == nil {
		t.Errorf("expected error for empty slice, got nil")
	}
}

// ==================== Get Tests ====================

func TestGet_RetrievesExistingItem(t *testing.T) {
	data := []*MockObject{
		{ID: "id1", Name: "obj1", Value: 10},
		{ID: "id2", Name: "obj2", Value: 20},
	}

	result, err := Get("id1", data)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if result.ID != "id1" || result.Name != "obj1" || result.Value != 10 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestGet_NonExistingItem(t *testing.T) {
	data := []*MockObject{
		{ID: "id1", Name: "obj1"},
	}

	result, err := Get("nonexistent", data)

	if err == nil {
		t.Errorf("expected error for non-existent ID, got nil")
	}
	if result != nil {
		t.Errorf("expected nil for error case, got %+v", result)
	}
}

func TestGet_EmptySlice(t *testing.T) {
	data := []*MockObject{}

	result, err := Get("id1", data)

	if err == nil {
		t.Errorf("expected error for empty slice, got nil")
	}
	if result != nil {
		t.Errorf("expected nil for error case")
	}
}

// ==================== DeleteEventBus Tests ====================

func TestNewEventBus_InitializesCorrectly(t *testing.T) {
	bus := NewEventBus()

	if bus == nil {
		t.Errorf("expected non-nil event bus")
	}
	if bus.subscribers == nil {
		t.Errorf("expected subscribers map to be initialized")
	}
}

func TestDeleteEventBus_SubscribeAddsHandler(t *testing.T) {
	bus := NewEventBus()

	handler := func(ids []string) error {
		return nil
	}

	bus.Subscribe("storage1", handler)

	if len(bus.subscribers["storage1"]) != 1 {
		t.Errorf("expected 1 handler for storage1, got %d", len(bus.subscribers["storage1"]))
	}
}

func TestDeleteEventBus_SubscribeMultipleHandlers(t *testing.T) {
	bus := NewEventBus()

	handler1 := func(ids []string) error {
		return nil
	}
	handler2 := func(ids []string) error {
		return nil
	}

	bus.Subscribe("storage1", handler1)
	bus.Subscribe("storage1", handler2)

	if len(bus.subscribers["storage1"]) != 2 {
		t.Errorf("expected 2 handlers for storage1, got %d", len(bus.subscribers["storage1"]))
	}
}

func TestDeleteEventBus_SubscribeMultipleStorages(t *testing.T) {
	bus := NewEventBus()

	handler := func(ids []string) error {
		return nil
	}

	bus.Subscribe("storage1", handler)
	bus.Subscribe("storage2", handler)

	if len(bus.subscribers["storage1"]) != 1 {
		t.Errorf("expected 1 handler for storage1, got %d", len(bus.subscribers["storage1"]))
	}
	if len(bus.subscribers["storage2"]) != 1 {
		t.Errorf("expected 1 handler for storage2, got %d", len(bus.subscribers["storage2"]))
	}
}

func TestDeleteEventBus_PublishCallsHandlers(t *testing.T) {
	bus := NewEventBus()
	called := false

	handler := func(ids []string) error {
		called = true
		return nil
	}

	bus.Subscribe("storage1", handler)
	bus.Publish("storage1", []string{"id1"})

	if !called {
		t.Errorf("expected handler to be called")
	}
}

func TestDeleteEventBus_PublishPassesIds(t *testing.T) {
	bus := NewEventBus()
	var receivedIds []string

	handler := func(ids []string) error {
		receivedIds = ids
		return nil
	}

	bus.Subscribe("storage1", handler)
	bus.Publish("storage1", []string{"id1", "id2", "id3"})

	if len(receivedIds) != 3 || receivedIds[0] != "id1" || receivedIds[1] != "id2" || receivedIds[2] != "id3" {
		t.Errorf("expected ids to be passed correctly, got %v", receivedIds)
	}
}

func TestDeleteEventBus_PublishMultipleHandlers(t *testing.T) {
	bus := NewEventBus()
	count := 0

	handler1 := func(ids []string) error {
		count++
		return nil
	}
	handler2 := func(ids []string) error {
		count++
		return nil
	}

	bus.Subscribe("storage1", handler1)
	bus.Subscribe("storage1", handler2)
	bus.Publish("storage1", []string{})

	if count != 2 {
		t.Errorf("expected both handlers to be called, count: %d", count)
	}
}

func TestDeleteEventBus_PublishNoHandlers(t *testing.T) {
	bus := NewEventBus()

	err := bus.Publish("storage1", []string{})

	if err != nil {
		t.Errorf("expected nil error when publishing with no handlers, got %v", err)
	}
}

func TestDeleteEventBus_PublishHandlerError(t *testing.T) {
	bus := NewEventBus()

	handler := func(ids []string) error {
		return errors.New("handler failed")
	}

	bus.Subscribe("storage1", handler)
	err := bus.Publish("storage1", []string{})

	if err == nil {
		t.Errorf("expected error from handler, got nil")
	}
}

func TestDeleteEventBus_PublishWrongStorage(t *testing.T) {
	bus := NewEventBus()
	called := false

	handler := func(ids []string) error {
		called = true
		return nil
	}

	bus.Subscribe("storage1", handler)
	bus.Publish("storage2", []string{})

	if called {
		t.Errorf("expected handler not to be called for different storage")
	}
}

func TestDeleteEventBus_ConcurrentOperations(t *testing.T) {
	bus := NewEventBus()

	// Test concurrent subscribe and publish operations
	handler := func(ids []string) error {
		return nil
	}

	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 10; i++ {
			bus.Subscribe("storage1", handler)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			bus.Publish("storage1", []string{})
		}
		done <- true
	}()

	<-done
	<-done

	if len(bus.subscribers["storage1"]) != 10 {
		t.Errorf("expected 10 handlers, got %d", len(bus.subscribers["storage1"]))
	}
}
