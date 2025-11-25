package storage

import (
	"testing"
)

// ==================== DriverGroup Tests ====================

func TestDriverGroup_GetId(t *testing.T) {
	dg := DriverGroup{Id: "test-id", Name: "Test Group"}
	if dg.GetId() != "test-id" {
		t.Errorf("expected 'test-id', got '%s'", dg.GetId())
	}
}

func TestDriverGroup_SetId(t *testing.T) {
	dg := DriverGroup{}
	dg.SetId("new-id")
	if dg.Id != "new-id" {
		t.Errorf("expected 'new-id', got '%s'", dg.Id)
	}
}

// ==================== Driver Tests ====================

func TestDriver_GetId(t *testing.T) {
	d := Driver{Id: "driver-1", Name: "Test Driver"}
	if d.GetId() != "driver-1" {
		t.Errorf("expected 'driver-1', got '%s'", d.GetId())
	}
}

func TestDriver_SetId(t *testing.T) {
	d := Driver{}
	d.SetId("new-driver-id")
	if d.Id != "new-driver-id" {
		t.Errorf("expected 'new-driver-id', got '%s'", d.Id)
	}
}

// ==================== DriverGroupStorage Tests ====================

func TestDriverGroupStorage_All_EmptyStore(t *testing.T) {
	store := &MemoryStore{}
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)

	groups, err := dgs.All()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestDriverGroupStorage_All_ExistingStore(t *testing.T) {
	existingData := []*DriverGroup{
		{
			Id:   "group1",
			Name: "Network Drivers",
			Type: Network,
			Drivers: []*Driver{
				{Id: "driver1", Name: "Driver 1", Type: Network},
			},
		},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)

	groups, err := dgs.All()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Id != "group1" {
		t.Errorf("expected 'group1', got '%s'", groups[0].Id)
	}
}

func TestDriverGroupStorage_Get_NotFound(t *testing.T) {
	store := &MemoryStore{}
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = []*DriverGroup{}

	_, err := dgs.Get("non-existent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDriverGroupStorage_Add_SingleGroup(t *testing.T) {
	store := &MemoryStore{}
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)

	group := DriverGroup{
		Name: "Network Drivers",
		Type: Network,
		Drivers: []*Driver{
			{Name: "Driver 1", Type: Network},
			{Name: "Driver 2", Type: Network},
		},
	}

	id, err := dgs.Add(group)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id == "" {
		t.Error("expected non-empty id")
	}

	if len(dgs.data) != 1 {
		t.Errorf("expected 1 group in storage, got %d", len(dgs.data))
	}

	if len(dgs.data[0].Drivers) != 2 {
		t.Errorf("expected 1 driver, got %d", len(dgs.data[0].Drivers))
	}

	if dgs.data[0].Drivers[0].Id == "" {
		t.Error("expected driver ID to be generated")
	}

	// Check all drivers have unique IDs
	driverIds := make(map[string]bool)
	for _, d := range dgs.data[0].Drivers {
		if driverIds[d.Id] {
			t.Errorf("duplicate driver ID: %s", d.Id)
		}
		driverIds[d.Id] = true
	}
}

func TestDriverGroupStorage_Update_Success(t *testing.T) {
	existingData := []*DriverGroup{
		{
			Id:   "group1",
			Name: "Original Name",
			Type: Network,
			Drivers: []*Driver{
				{Id: "driver1", Name: "Driver 1", Type: Network},
			},
		},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = existingData

	updatedGroup := DriverGroup{
		Id:   "group1",
		Name: "Updated Name",
		Type: Display,
		Drivers: []*Driver{
			{Id: "driver1", Name: "Driver 1", Type: Display},
		},
	}

	result, err := dgs.Update(updatedGroup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "Updated Name" {
		t.Errorf("expected 'Updated Name', got '%s'", result.Name)
	}

	if dgs.data[0].Type != Display {
		t.Errorf("expected Display type, got %v", dgs.data[0].Type)
	}
}

func TestDriverGroupStorage_Update_RemoveDrivers(t *testing.T) {
	existingData := []*DriverGroup{
		{
			Id:   "group1",
			Name: "Network Drivers",
			Type: Network,
			Drivers: []*Driver{
				{Id: "driver1", Name: "Driver 1", Type: Network, Incompatibles: []string{"driver2"}},
				{Id: "driver2", Name: "Driver 2", Type: Network},
			},
		},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = existingData

	// Remove driver2 from incompatibles
	existingData[0].Drivers[0].Incompatibles = []string{"driver2"}

	updatedGroup := DriverGroup{
		Id:   "group1",
		Name: "Network Drivers",
		Type: Network,
		Drivers: []*Driver{
			{Id: "driver1", Name: "Driver 1", Type: Network, Incompatibles: []string{"driver2"}},
		},
	}

	_, err := dgs.Update(updatedGroup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if driver2 was removed from incompatibles
	if len(dgs.data[0].Drivers[0].Incompatibles) != 0 {
		t.Errorf("expected incompatibles to be empty, got %v", dgs.data[0].Drivers[0].Incompatibles)
	}
}

func TestDriverGroupStorage_Remove_Success(t *testing.T) {
	existingData := []*DriverGroup{
		{
			Id:   "group1",
			Name: "Group 1",
			Type: Network,
			Drivers: []*Driver{
				{Id: "driver1", Name: "Driver 1", Type: Network},
			},
		},
		{
			Id:   "group2",
			Name: "Group 2",
			Type: Display,
			Drivers: []*Driver{
				{Id: "driver2", Name: "Driver 2", Type: Display, Incompatibles: []string{"driver1"}},
			},
		},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = existingData

	err := dgs.Remove("group1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dgs.data) != 1 {
		t.Errorf("expected 1 group, got %d", len(dgs.data))
	}

	if dgs.data[0].Id != "group2" {
		t.Errorf("expected 'group2', got '%s'", dgs.data[0].Id)
	}

	// Check if driver1 was removed from incompatibles
	if len(dgs.data[0].Drivers[0].Incompatibles) != 0 {
		t.Errorf("expected incompatibles to be empty, got %v", dgs.data[0].Drivers[0].Incompatibles)
	}
}

func TestDriverGroupStorage_Remove_NotFound(t *testing.T) {
	store := &MemoryStore{}
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = []*DriverGroup{}

	err := dgs.Remove("non-existent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDriverGroupStorage_Remove_PublishesEvent(t *testing.T) {
	existingData := []*DriverGroup{
		{
			Id:   "group1",
			Name: "Group 1",
			Type: Network,
			Drivers: []*Driver{
				{Id: "driver1", Name: "Driver 1", Type: Network},
				{Id: "driver2", Name: "Driver 2", Type: Network},
			},
		},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()

	eventPublished := false
	var publishedIds []string
	eventBus.Subscribe("DriverGroup", func(ids []string) error {
		eventPublished = true
		publishedIds = ids
		return nil
	})

	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = existingData

	err := dgs.Remove("group1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !eventPublished {
		t.Error("expected event to be published")
	}

	if len(publishedIds) != 2 {
		t.Errorf("expected 2 driver IDs, got %d", len(publishedIds))
	}
}

func TestDriverGroupStorage_IndexOf_Success(t *testing.T) {
	existingData := []*DriverGroup{
		{Id: "group1", Name: "Group 1", Type: Network},
		{Id: "group2", Name: "Group 2", Type: Display},
		{Id: "group3", Name: "Group 3", Type: Miscellaneous},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = existingData

	index, err := dgs.IndexOf("group2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if index != 1 {
		t.Errorf("expected index 1, got %d", index)
	}
}

func TestDriverGroupStorage_IndexOf_NotFound(t *testing.T) {
	store := &MemoryStore{}
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = []*DriverGroup{}

	_, err := dgs.IndexOf("non-existent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDriverGroupStorage_MoveBehind_ValidMove(t *testing.T) {
	existingData := []*DriverGroup{
		{Id: "group1", Name: "Group 1", Type: Network},
		{Id: "group2", Name: "Group 2", Type: Display},
		{Id: "group3", Name: "Group 3", Type: Miscellaneous},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = existingData

	result, err := dgs.MoveBehind("group1", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result[2].Id != "group1" {
		t.Errorf("expected 'group1' at index 2, got '%s'", result[2].Id)
	}
}

func TestDriverGroupStorage_MoveBehind_OutOfBounds(t *testing.T) {
	existingData := []*DriverGroup{
		{Id: "group1", Name: "Group 1", Type: Network},
		{Id: "group2", Name: "Group 2", Type: Display},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = existingData

	_, err := dgs.MoveBehind("group1", 5)
	if err == nil {
		t.Fatal("expected error for out of bounds index")
	}
}

func TestDriverGroupStorage_MoveBehind_InvalidSourceId(t *testing.T) {
	existingData := []*DriverGroup{
		{Id: "group1", Name: "Group 1", Type: Network},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = existingData

	_, err := dgs.MoveBehind("non-existent", 0)
	if err == nil {
		t.Fatal("expected error for non-existent group")
	}
}

func TestDriverGroupStorage_MoveBehind_NoChange(t *testing.T) {
	existingData := []*DriverGroup{
		{Id: "group1", Name: "Group 1", Type: Network},
		{Id: "group2", Name: "Group 2", Type: Display},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = existingData

	result, err := dgs.MoveBehind("group2", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 groups, got %d", len(result))
	}
}

func TestDriverGroupStorage_MoveBehind_NegativeIndex(t *testing.T) {
	existingData := []*DriverGroup{
		{Id: "group1", Name: "Group 1", Type: Network},
		{Id: "group2", Name: "Group 2", Type: Display},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = existingData

	result, err := dgs.MoveBehind("group1", -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return current state unchanged when index is -1
	if len(result) != 2 {
		t.Errorf("expected 2 groups, got %d", len(result))
	}
}

func TestDriverGroupStorage_CopyOfAll_Isolation(t *testing.T) {
	existingData := []*DriverGroup{
		{
			Id:   "group1",
			Name: "Group 1",
			Type: Network,
			Drivers: []*Driver{
				{Id: "driver1", Name: "Driver 1", Type: Network},
			},
		},
	}

	store := &MemoryStore{}
	store.Write(existingData)
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)
	dgs.data = existingData

	copy := dgs.copyOfAll()

	// Modify the copy
	copy[0].Name = "Modified Name"

	// Original should be unchanged
	if dgs.data[0].Name != "Group 1" {
		t.Errorf("expected 'Group 1', got '%s'", dgs.data[0].Name)
	}
}

// ==================== Integration Tests ====================

func TestDriverGroupStorage_CompleteWorkflow(t *testing.T) {
	store := &MemoryStore{}
	eventBus := NewEventBus()
	dgs := NewDriverGroupStorage(store, eventBus)

	// Add first group
	group1 := DriverGroup{
		Name: "Network Drivers",
		Type: Network,
		Drivers: []*Driver{
			{Name: "Network Driver 1", Type: Network},
		},
	}

	id1, err := dgs.Add(group1)
	if err != nil {
		t.Fatalf("failed to add group1: %v", err)
	}

	// Add second group
	group2 := DriverGroup{
		Name: "Display Drivers",
		Type: Display,
		Drivers: []*Driver{
			{Name: "Display Driver 1", Type: Display},
		},
	}

	id2, err := dgs.Add(group2)
	if err != nil {
		t.Fatalf("failed to add group2: %v", err)
	}

	// Verify both groups are in storage
	if len(dgs.data) != 2 {
		t.Errorf("expected 2 groups in storage, got %d", len(dgs.data))
	}

	// Get specific group
	retrieved, err := dgs.Get(id1)
	if err != nil {
		t.Fatalf("failed to get group: %v", err)
	}

	if retrieved.Name != "Network Drivers" {
		t.Errorf("expected 'Network Drivers', got '%s'", retrieved.Name)
	}

	// Update group
	retrieved.Name = "Updated Network Drivers"
	_, err = dgs.Update(retrieved)
	if err != nil {
		t.Fatalf("failed to update group: %v", err)
	}

	// Verify update
	updated, err := dgs.Get(id1)
	if err != nil {
		t.Fatalf("failed to get updated group: %v", err)
	}

	if updated.Name != "Updated Network Drivers" {
		t.Errorf("expected 'Updated Network Drivers', got '%s'", updated.Name)
	}

	// Move group
	_, err = dgs.MoveBehind(id1, 0)
	if err != nil {
		t.Fatalf("failed to move group: %v", err)
	}

	// Remove group
	err = dgs.Remove(id2)
	if err != nil {
		t.Fatalf("failed to remove group: %v", err)
	}

	// Verify removal
	if len(dgs.data) != 1 {
		t.Errorf("expected 1 group, got %d", len(dgs.data))
	}
}
