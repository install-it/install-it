// Package storage_test provides external black-box tests for DriverGroupStorage.
// The internal tests in driver_test.go cover individual CRUD operations; this file
// focuses on cross-storage integration (cascade delete) and ordering.
package storage_test

import (
	"testing"

	"install-it/pkg/storage"
)

// TestDriverGroupStorage_CascadeDeleteFromMatchRule verifies that removing a
// DriverGroup also removes it from all RuleSet associations via FK CASCADE.
func TestDriverGroupStorage_CascadeDeleteFromMatchRule(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)
	mrs := storage.NewRuleSetStorage(db)

	groupID := addTestGroup(t, dgs, storage.DriverGroup{
		Name: "Cascade Group",
		Type: storage.Network,
		Drivers: []*storage.Driver{
			{Name: "Net Driver", Type: storage.Network},
		},
	})
	otherID := addTestGroup(t, dgs, storage.DriverGroup{
		Name: "Other Group",
		Type: storage.Display,
	})

	if err := mrs.Add(storage.RuleSet{
		Name:           "Cascade RuleSet",
		DriverGroupIds: []uint{groupID, otherID},
	}); err != nil {
		t.Fatalf("mrs.Add: %v", err)
	}

	all, err := mrs.All()
	if err != nil {
		t.Fatalf("mrs.All: %v", err)
	}
	rs := all[0]
	if !containsUint(rs.DriverGroupIds, groupID) {
		t.Fatalf("pre-condition: groupID %d not in DriverGroupIds %v", groupID, rs.DriverGroupIds)
	}

	if err := dgs.Remove(groupID); err != nil {
		t.Fatalf("dgs.Remove: %v", err)
	}

	all, err = mrs.All()
	if err != nil {
		t.Fatalf("mrs.All after cascade: %v", err)
	}
	rs = all[0]
	if containsUint(rs.DriverGroupIds, groupID) {
		t.Errorf("cascade delete failed: groupID %d still present in DriverGroupIds %v",
			groupID, rs.DriverGroupIds)
	}
	if !containsUint(rs.DriverGroupIds, otherID) {
		t.Errorf("cascade delete incorrectly removed otherID %d from DriverGroupIds", otherID)
	}
}

// TestDriverGroupStorage_RemoveGroup_NotInAllAfterRemove verifies a removed group
// is absent from the All() result.
func TestDriverGroupStorage_RemoveGroup_NotInAllAfterRemove(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)

	id := addTestGroup(t, dgs, storage.DriverGroup{Name: "Temporary", Type: storage.Display})

	if err := dgs.Remove(id); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	groups, err := dgs.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	for _, g := range groups {
		if g.Id == id {
			t.Errorf("removed group %d is still returned by All()", id)
		}
	}
}

// TestDriverGroupStorage_AddAssignsIdsToDrivers checks that drivers without
// pre-set IDs receive non-zero autoincrement IDs after Add.
func TestDriverGroupStorage_AddAssignsIdsToDrivers(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)

	id := addTestGroup(t, dgs, storage.DriverGroup{
		Name: "ID Test Group",
		Type: storage.Miscellaneous,
		Drivers: []*storage.Driver{
			{Name: "Driver A"},
			{Name: "Driver B"},
		},
	})

	group, err := dgs.Get(id)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if len(group.Drivers) != 2 {
		t.Fatalf("expected 2 drivers, got %d", len(group.Drivers))
	}

	seen := make(map[uint]bool)
	for _, d := range group.Drivers {
		if d.Id == 0 {
			t.Errorf("driver %q has zero ID, want non-zero uint", d.Name)
		}
		if seen[d.Id] {
			t.Errorf("duplicate driver ID %d", d.Id)
		}
		seen[d.Id] = true
	}
}

// TestDriverGroupStorage_MoveBehind_OrderChanges verifies that MoveBehind
// correctly reorders the groups list.
func TestDriverGroupStorage_MoveBehind_OrderChanges(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)

	var ids [4]uint
	for i, name := range []string{"A", "B", "C", "D"} {
		ids[i] = addTestGroup(t, dgs, storage.DriverGroup{Name: name, Type: storage.Network})
	}

	// Move A (position 0) to behind index 2 → [B, C, D, A]
	if err := dgs.MoveBehind(ids[0], 2); err != nil {
		t.Fatalf("MoveBehind: %v", err)
	}

	result, err := dgs.All()
	if err != nil {
		t.Fatalf("All after MoveBehind: %v", err)
	}
	if len(result) != 4 {
		t.Fatalf("expected 4 groups, got %d", len(result))
	}
	if result[3].Id != ids[0] {
		t.Errorf("after MoveBehind(A, 2): index 3 has %d, want %d (A)", result[3].Id, ids[0])
	}
	if result[0].Id != ids[1] {
		t.Errorf("after MoveBehind(A, 2): index 0 has %d, want %d (B)", result[0].Id, ids[1])
	}
}

// TestDriverGroupStorage_UpdateDriverIncompatibles verifies that removing a
// driver from the group via Update cleans up incompatible references.
func TestDriverGroupStorage_UpdateDriverIncompatibles(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)

	groupID := addTestGroup(t, dgs, storage.DriverGroup{
		Name: "Incompat Group",
		Type: storage.Network,
		Drivers: []*storage.Driver{
			{Name: "D1"},
			{Name: "D2"},
		},
	})

	group, err := dgs.Get(groupID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	d1ID := group.Drivers[0].Id
	d2ID := group.Drivers[1].Id

	group.Drivers[0].IncompatibleIds = []uint{d2ID}
	group.Drivers[1].IncompatibleIds = []uint{d1ID}
	if err := dgs.Update(group); err != nil {
		t.Fatalf("Update with incompatibles: %v", err)
	}

	// Simulate UI removing D2: also clear D2 from D1's incompatibles
	var d1 *storage.Driver
	for _, d := range group.Drivers {
		if d.Id == d1ID {
			d.IncompatibleIds = nil
			d1 = d
			break
		}
	}
	group.Drivers = []*storage.Driver{d1}
	if err := dgs.Update(group); err != nil {
		t.Fatalf("Update removing D2: %v", err)
	}

	updated, err := dgs.Get(groupID)
	if err != nil {
		t.Fatalf("Get after update: %v", err)
	}
	if len(updated.Drivers) != 1 {
		t.Fatalf("expected 1 driver, got %d", len(updated.Drivers))
	}
	if len(updated.Drivers[0].IncompatibleIds) != 0 {
		t.Errorf("expected empty incompatibles after removing D2, got %v",
			updated.Drivers[0].IncompatibleIds)
	}
}
