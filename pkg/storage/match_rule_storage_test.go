// Package storage_test provides external black-box tests for MatchRuleStorage.
package storage_test

import (
	"testing"

	"install-it/pkg/storage"
)

// TestMatchRuleStorage_AllEmpty verifies that an empty store returns a non-nil
// empty slice.
func TestMatchRuleStorage_AllEmpty(t *testing.T) {
	db := openExternalTestDB(t)
	mrs := storage.NewMatchRuleStorage(db)

	results, err := mrs.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if results == nil {
		t.Error("expected empty (non-nil) slice, got nil")
	}
	if len(results) != 0 {
		t.Errorf("expected 0 items, got %d", len(results))
	}
}

// TestMatchRuleStorage_AddGetUpdate performs a full CRUD round-trip.
func TestMatchRuleStorage_AddGetUpdate(t *testing.T) {
	db := openExternalTestDB(t)
	mrs := storage.NewMatchRuleStorage(db)

	id, err := mrs.Add(storage.RuleSet{
		Name:           "CRUD Test",
		ShouldHitAll:   true,
		DriverGroupIds: []string{"group-aabbccdd"},
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if len(id) != 8 {
		t.Errorf("expected 8-char ID, got %q", id)
	}

	got, err := mrs.Get(id)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "CRUD Test" {
		t.Errorf("Name: got %q, want 'CRUD Test'", got.Name)
	}
	if !got.ShouldHitAll {
		t.Error("ShouldHitAll: got false, want true")
	}

	got.Name = "Updated CRUD Test"
	got.ShouldHitAll = false
	updated, err := mrs.Update(got)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "Updated CRUD Test" {
		t.Errorf("updated Name: got %q, want 'Updated CRUD Test'", updated.Name)
	}

	final, err := mrs.Get(id)
	if err != nil {
		t.Fatalf("Get after update: %v", err)
	}
	if final.Name != "Updated CRUD Test" {
		t.Errorf("persisted Name: got %q", final.Name)
	}
	if final.ShouldHitAll {
		t.Error("persisted ShouldHitAll should be false")
	}
}

// TestMatchRuleStorage_AddPreservesRuleFields verifies all Rule fields survive
// the Add → Get cycle.
func TestMatchRuleStorage_AddPreservesRuleFields(t *testing.T) {
	db := openExternalTestDB(t)
	mrs := storage.NewMatchRuleStorage(db)

	rules := []storage.Rule{
		{
			Source:          storage.Cpu,
			Operator:        storage.Contain,
			IsCaseSensitive: true,
			ShouldHitAll:    false,
			Values:          []string{"Intel", "AMD"},
		},
		{
			Source:          storage.Nic,
			Operator:        storage.Regex,
			IsCaseSensitive: false,
			ShouldHitAll:    true,
			Values:          []string{`^Realtek.*`},
		},
		{
			Source:   storage.Gpu,
			Operator: storage.NotEqual,
			Values:   []string{"Basic Display Adapter"},
		},
	}

	id, err := mrs.Add(storage.RuleSet{
		Name:         "Fields Test",
		Rules:        rules,
		ShouldHitAll: true,
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	got, err := mrs.Get(id)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(got.Rules) != len(rules) {
		t.Fatalf("rule count: got %d, want %d", len(got.Rules), len(rules))
	}

	for i, want := range rules {
		r := got.Rules[i]
		if r.Source != want.Source {
			t.Errorf("rule[%d].Source: got %v, want %v", i, r.Source, want.Source)
		}
		if r.Operator != want.Operator {
			t.Errorf("rule[%d].Operator: got %v, want %v", i, r.Operator, want.Operator)
		}
		if r.IsCaseSensitive != want.IsCaseSensitive {
			t.Errorf("rule[%d].IsCaseSensitive: got %v, want %v", i, r.IsCaseSensitive, want.IsCaseSensitive)
		}
		if r.ShouldHitAll != want.ShouldHitAll {
			t.Errorf("rule[%d].ShouldHitAll: got %v, want %v", i, r.ShouldHitAll, want.ShouldHitAll)
		}
		if len(r.Values) != len(want.Values) {
			t.Errorf("rule[%d] Values len: got %d, want %d", i, len(r.Values), len(want.Values))
			continue
		}
		for j, v := range want.Values {
			if r.Values[j] != v {
				t.Errorf("rule[%d].Values[%d]: got %q, want %q", i, j, r.Values[j], v)
			}
		}
	}
}

// TestMatchRuleStorage_RemovePurgesOrphanedGroupIds verifies that
// DriverGroupStorage.Remove also cleans up driver IDs from DriverGroupIds
// in all RuleSets.
func TestMatchRuleStorage_RemovePurgesOrphanedGroupIds(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)
	mrs := storage.NewMatchRuleStorage(db)

	groupID, err := dgs.Add(storage.DriverGroup{
		Name:    "Temp Group",
		Type:    storage.Network,
		Drivers: []*storage.Driver{{Name: "D1"}},
	})
	if err != nil {
		t.Fatalf("dgs.Add: %v", err)
	}

	idToKeep := "ccdd3344"
	rsID, err := mrs.Add(storage.RuleSet{
		Name:           "Purge Test",
		DriverGroupIds: []string{groupID, idToKeep},
	})
	if err != nil {
		t.Fatalf("mrs.Add: %v", err)
	}

	if err := dgs.Remove(groupID); err != nil {
		t.Fatalf("dgs.Remove: %v", err)
	}

	rs, err := mrs.Get(rsID)
	if err != nil {
		t.Fatalf("mrs.Get: %v", err)
	}
	if containsStr(rs.DriverGroupIds, groupID) {
		t.Errorf("groupID %q should have been purged, still in %v", groupID, rs.DriverGroupIds)
	}
	if !containsStr(rs.DriverGroupIds, idToKeep) {
		t.Errorf("idToKeep %q should still be present, not in %v", idToKeep, rs.DriverGroupIds)
	}
}

// TestMatchRuleStorage_UpdateRuleSet verifies Update persists rule changes.
func TestMatchRuleStorage_UpdateRuleSet(t *testing.T) {
	db := openExternalTestDB(t)
	mrs := storage.NewMatchRuleStorage(db)

	id, err := mrs.Add(storage.RuleSet{
		Name:  "Before Update",
		Rules: []storage.Rule{{Source: storage.Memory, Operator: storage.Equal, Values: []string{"8GB"}}},
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	original, err := mrs.Get(id)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	original.Name = "After Update"
	original.Rules = []storage.Rule{
		{Source: storage.Storage, Operator: storage.NotContain, Values: []string{"USB"}},
		{Source: storage.Motherboard, Operator: storage.Regex, Values: []string{`^ASUS.*`}},
	}

	updated, err := mrs.Update(original)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "After Update" {
		t.Errorf("updated.Name: got %q", updated.Name)
	}
	if len(updated.Rules) != 2 {
		t.Errorf("updated.Rules len: got %d, want 2", len(updated.Rules))
	}

	persisted, err := mrs.Get(id)
	if err != nil {
		t.Fatalf("Get after update: %v", err)
	}
	if persisted.Name != "After Update" {
		t.Errorf("persisted.Name: got %q", persisted.Name)
	}
	if len(persisted.Rules) != 2 {
		t.Errorf("persisted.Rules len: got %d", len(persisted.Rules))
	}
}

// TestMatchRuleStorage_RemoveNonExistent verifies that removing a non-existent
// ID returns an error.
func TestMatchRuleStorage_RemoveNonExistent(t *testing.T) {
	db := openExternalTestDB(t)
	mrs := storage.NewMatchRuleStorage(db)

	if err := mrs.Remove("notexist"); err == nil {
		t.Error("expected error when removing non-existent id, got nil")
	}
}

// TestMatchRuleStorage_AllReturnsAllAdded verifies All() returns every added ruleset.
func TestMatchRuleStorage_AllReturnsAllAdded(t *testing.T) {
	db := openExternalTestDB(t)
	mrs := storage.NewMatchRuleStorage(db)

	names := []string{"Alpha", "Beta", "Gamma"}
	for _, name := range names {
		if _, err := mrs.Add(storage.RuleSet{Name: name}); err != nil {
			t.Fatalf("Add %s: %v", name, err)
		}
	}

	all, err := mrs.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(all) != len(names) {
		t.Errorf("All() count: got %d, want %d", len(all), len(names))
	}
}
