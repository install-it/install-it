// Package storage_test provides external black-box tests for DriverGroupStorage.
// The internal tests in driver_test.go cover individual CRUD operations; this file
// focuses on cross-storage integration (cascade delete via EventBus) and ordering.
package storage_test

import (
	"testing"

	"install-it/pkg/storage"
)

// TestDriverGroupStorage_CascadeDeleteFromMatchRule verifies the EventBus wiring
// between DriverGroupStorage and MatchRuleStorage.
//
// Implementation note: DriverGroupStorage.Remove publishes the IDs of the *drivers*
// within the removed group (not the group ID itself) under the "DriverGroup" EventBus
// key. MatchRuleStorage subscribes to that key and removes the received IDs from
// any RuleSet.DriverGroupIds that contains them.
//
// This test exercises that end-to-end wiring by placing the driver IDs into
// DriverGroupIds and verifying they are removed after the group is deleted.
func TestDriverGroupStorage_CascadeDeleteFromMatchRule(t *testing.T) {
	t.Parallel()

	eventBus := storage.NewEventBus()
	dgs := storage.NewDriverGroupStorage(&storage.MemoryStore{}, eventBus)
	mrs := storage.NewMatchRuleStorage(&storage.MemoryStore{}, eventBus)

	// Prime both storages so their in-memory slices are initialised.
	if _, err := dgs.All(); err != nil {
		t.Fatalf("dgs.All: %v", err)
	}
	if _, err := mrs.All(); err != nil {
		t.Fatalf("mrs.All: %v", err)
	}

	// Add a driver group with one driver.
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

	// Retrieve the generated driver ID.
	group, err := dgs.Get(groupID)
	if err != nil {
		t.Fatalf("dgs.Get: %v", err)
	}
	if len(group.Drivers) == 0 {
		t.Fatal("group has no drivers")
	}
	driverID := group.Drivers[0].Id

	// Add a ruleset where DriverGroupIds holds the driver ID.
	// (Remove publishes driver IDs; the subscriber removes them from DriverGroupIds.)
	ruleSetID, err := mrs.Add(storage.RuleSet{
		Name:           "Cascade RuleSet",
		DriverGroupIds: []string{driverID, "unrelated-id"},
	})
	if err != nil {
		t.Fatalf("mrs.Add: %v", err)
	}

	// Pre-condition: driverID is present.
	rs, err := mrs.Get(ruleSetID)
	if err != nil {
		t.Fatalf("mrs.Get before remove: %v", err)
	}
	if !containsStr(rs.DriverGroupIds, driverID) {
		t.Fatalf("pre-condition: driverID %s not in DriverGroupIds", driverID)
	}

	// Remove the driver group — publishes the driver ID via EventBus.
	if err := dgs.Remove(groupID); err != nil {
		t.Fatalf("dgs.Remove: %v", err)
	}

	// The EventBus subscriber should have removed driverID from DriverGroupIds.
	rs, err = mrs.Get(ruleSetID)
	if err != nil {
		t.Fatalf("mrs.Get after cascade: %v", err)
	}
	if containsStr(rs.DriverGroupIds, driverID) {
		t.Errorf("cascade delete failed: driverID %s still present in DriverGroupIds %v",
			driverID, rs.DriverGroupIds)
	}
	// Unrelated ID must be untouched.
	if !containsStr(rs.DriverGroupIds, "unrelated-id") {
		t.Errorf("cascade delete incorrectly removed 'unrelated-id' from DriverGroupIds")
	}
}

// TestDriverGroupStorage_RemoveGroup_NotInAllAfterRemove verifies a removed group
// is absent from the All() result — using the storage public API only.
func TestDriverGroupStorage_RemoveGroup_NotInAllAfterRemove(t *testing.T) {
	t.Parallel()

	eventBus := storage.NewEventBus()
	dgs := storage.NewDriverGroupStorage(&storage.MemoryStore{}, eventBus)

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
	t.Parallel()

	dgs := storage.NewDriverGroupStorage(&storage.MemoryStore{}, storage.NewEventBus())

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
//
// MoveBehind(id, index) moves the element with the given id so that it sits
// directly after the element currently at position index (i.e. at index+1 after
// bubbling through all intermediate positions).
func TestDriverGroupStorage_MoveBehind_OrderChanges(t *testing.T) {
	t.Parallel()

	dgs := storage.NewDriverGroupStorage(&storage.MemoryStore{}, storage.NewEventBus())

	// Add 4 groups in order A, B, C, D.
	var ids [4]string
	names := []string{"A", "B", "C", "D"}
	for i, name := range names {
		id, err := dgs.Add(storage.DriverGroup{Name: name, Type: storage.Network})
		if err != nil {
			t.Fatalf("Add %s: %v", name, err)
		}
		ids[i] = id
	}

	// MoveBehind(A, 2): move A (currently at index 0) behind index 2.
	// The algorithm bubbles A rightward: [A,B,C,D] → [B,A,C,D] → [B,C,A,D] → [B,C,D,A]
	// Result: A lands at the last position (index 3 in a 4-element array).
	result, err := dgs.MoveBehind(ids[0], 2)
	if err != nil {
		t.Fatalf("MoveBehind: %v", err)
	}

	if len(result) != 4 {
		t.Fatalf("expected 4 groups, got %d", len(result))
	}

	// ids[0] (A) should now be at index 3 (last position).
	if result[3].Id != ids[0] {
		t.Errorf("after MoveBehind(A, 2): index 3 has %s, want %s (A)",
			result[3].Id, ids[0])
	}
	// ids[1] (B) should be first.
	if result[0].Id != ids[1] {
		t.Errorf("after MoveBehind(A, 2): index 0 has %s, want %s (B)",
			result[0].Id, ids[1])
	}
}

// TestDriverGroupStorage_UpdateDriverIncompatibles verifies that removing a
// driver from the group via Update cleans up incompatible references within
// the same group.
func TestDriverGroupStorage_UpdateDriverIncompatibles(t *testing.T) {
	t.Parallel()

	dgs := storage.NewDriverGroupStorage(&storage.MemoryStore{}, storage.NewEventBus())

	// Bootstrap with two drivers, each referencing the other as incompatible.
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

	// Set D1's incompatibles to [D2] and D2's incompatibles to [D1].
	group.Drivers[0].Incompatibles = []string{d2ID}
	group.Drivers[1].Incompatibles = []string{d1ID}
	if _, err := dgs.Update(group); err != nil {
		t.Fatalf("Update with incompatibles: %v", err)
	}

	// Now update the group removing D2 — D1's incompatibles should be cleared.
	group.Drivers = group.Drivers[:1] // keep only D1
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
