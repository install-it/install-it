// Package storage_test provides external black-box tests for DriverGroupStorage.
// The internal tests in driver_test.go cover individual CRUD operations; this file
// focuses on cross-storage integration (cascade delete) and ordering.
package storage_test

import (
	"testing"

	"install-it/pkg/storage"
)

// TestDriverGroupStorage_CascadeDeleteFromMatchRule verifies that
// DriverGroupStorage.Remove cleans up DriverGroupIds in all RuleSets.
//
// The driver IDs belonging to the removed group are stripped from every
// RuleSet.DriverGroupIds as part of the same transaction.
func TestDriverGroupStorage_CascadeDeleteFromMatchRule(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)
	mrs := storage.NewMatchRuleStorage(db)

	groupID, err := dgs.Add(storage.DriverGroup{
		Name: "Cascade Group",
		Type: storage.Network,
		Drivers: []*storage.Driver{
			{Name: "Net Driver", Type: storage.Network},
		},
	})
	if err != nil {
		t.Fatalf("dgs.Add: %v", err)
	}

	ruleSetID, err := mrs.Add(storage.RuleSet{
		Name:           "Cascade RuleSet",
		DriverGroupIds: []string{groupID, "unrelated-id"},
	})
	if err != nil {
		t.Fatalf("mrs.Add: %v", err)
	}

	rs, err := mrs.Get(ruleSetID)
	if err != nil {
		t.Fatalf("mrs.Get before remove: %v", err)
	}
	if !containsStr(rs.DriverGroupIds, groupID) {
		t.Fatalf("pre-condition: groupID %s not in DriverGroupIds", groupID)
	}

	if err := dgs.Remove(groupID); err != nil {
		t.Fatalf("dgs.Remove: %v", err)
	}

	rs, err = mrs.Get(ruleSetID)
	if err != nil {
		t.Fatalf("mrs.Get after cascade: %v", err)
	}
	if containsStr(rs.DriverGroupIds, groupID) {
		t.Errorf("cascade delete failed: groupID %s still present in DriverGroupIds %v",
			groupID, rs.DriverGroupIds)
	}
	if !containsStr(rs.DriverGroupIds, "unrelated-id") {
		t.Errorf("cascade delete incorrectly removed 'unrelated-id' from DriverGroupIds")
	}
}

// TestDriverGroupStorage_RemoveGroup_NotInAllAfterRemove verifies a removed group
// is absent from the All() result.
func TestDriverGroupStorage_RemoveGroup_NotInAllAfterRemove(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)

	id, err := dgs.Add(storage.DriverGroup{Name: "Temporary", Type: storage.Display})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if err := dgs.Remove(id); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	groups, err := dgs.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	for _, g := range groups {
		if g.Id == id {
			t.Errorf("removed group %s is still returned by All()", id)
		}
	}
}

// TestDriverGroupStorage_AddAssignsIdsToDrivers checks that drivers without
// pre-set IDs receive generated 8-char hex IDs after Add.
func TestDriverGroupStorage_AddAssignsIdsToDrivers(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)

	id, err := dgs.Add(storage.DriverGroup{
		Name: "ID Test Group",
		Type: storage.Miscellaneous,
		Drivers: []*storage.Driver{
			{Name: "Driver A"},
			{Name: "Driver B"},
		},
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	group, err := dgs.Get(id)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if len(group.Drivers) != 2 {
		t.Fatalf("expected 2 drivers, got %d", len(group.Drivers))
	}

	seen := make(map[string]bool)
	for _, d := range group.Drivers {
		if len(d.Id) != 8 {
			t.Errorf("driver %q has ID %q (len=%d), want 8-char hex", d.Name, d.Id, len(d.Id))
		}
		if seen[d.Id] {
			t.Errorf("duplicate driver ID %s", d.Id)
		}
		seen[d.Id] = true
	}
}

// TestDriverGroupStorage_MoveBehind_OrderChanges verifies that MoveBehind
// correctly reorders the groups list.
func TestDriverGroupStorage_MoveBehind_OrderChanges(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)

	var ids [4]string
	for i, name := range []string{"A", "B", "C", "D"} {
		id, err := dgs.Add(storage.DriverGroup{Name: name, Type: storage.Network})
		if err != nil {
			t.Fatalf("Add %s: %v", name, err)
		}
		ids[i] = id
	}

	// MoveBehind A (index 0) behind index 2 → bubbles: [A,B,C,D] → [B,C,D,A]
	result, err := dgs.MoveBehind(ids[0], 2)
	if err != nil {
		t.Fatalf("MoveBehind: %v", err)
	}
	if len(result) != 4 {
		t.Fatalf("expected 4 groups, got %d", len(result))
	}
	if result[3].Id != ids[0] {
		t.Errorf("after MoveBehind(A, 2): index 3 has %s, want %s (A)", result[3].Id, ids[0])
	}
	if result[0].Id != ids[1] {
		t.Errorf("after MoveBehind(A, 2): index 0 has %s, want %s (B)", result[0].Id, ids[1])
	}
}

// TestDriverGroupStorage_UpdateDriverIncompatibles verifies that removing a
// driver from the group via Update cleans up incompatible references.
func TestDriverGroupStorage_UpdateDriverIncompatibles(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)

	groupID, err := dgs.Add(storage.DriverGroup{
		Name: "Incompat Group",
		Type: storage.Network,
		Drivers: []*storage.Driver{
			{Name: "D1"},
			{Name: "D2"},
		},
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	group, err := dgs.Get(groupID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	d1ID := group.Drivers[0].Id
	d2ID := group.Drivers[1].Id

	group.Drivers[0].Incompatibles = []string{d2ID}
	group.Drivers[1].Incompatibles = []string{d1ID}
	if _, err := dgs.Update(group); err != nil {
		t.Fatalf("Update with incompatibles: %v", err)
	}

	group.Drivers = group.Drivers[:1]
	if _, err := dgs.Update(group); err != nil {
		t.Fatalf("Update removing D2: %v", err)
	}

	updated, err := dgs.Get(groupID)
	if err != nil {
		t.Fatalf("Get after update: %v", err)
	}
	if len(updated.Drivers) != 1 {
		t.Fatalf("expected 1 driver, got %d", len(updated.Drivers))
	}
	if len(updated.Drivers[0].Incompatibles) != 0 {
		t.Errorf("expected empty incompatibles after removing D2, got %v",
			updated.Drivers[0].Incompatibles)
	}
}
