// Package storage_test provides external black-box tests for RuleSetStorage.
package storage_test

import (
	"testing"

	"install-it/pkg/storage"
)

// TestRuleSetStorage_AllEmpty verifies that an empty store returns a non-nil
// empty slice.
func TestRuleSetStorage_AllEmpty(t *testing.T) {
	db := openExternalTestDB(t)
	rss := storage.NewRuleSetStorage(db)

	results, err := rss.All()
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

// TestRuleSetStorage_AddGetUpdate performs a full CRUD round-trip.
func TestRuleSetStorage_AddGetUpdate(t *testing.T) {
	db := openExternalTestDB(t)
	rss := storage.NewRuleSetStorage(db)

	if err := rss.Add(storage.RuleSet{
		Name:         "CRUD Test",
		ShouldHitAll: true,
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	all, err := rss.All()
	if err != nil {
		t.Fatalf("All after Add: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 ruleset, got %d", len(all))
	}
	id := all[0].Id
	if id == 0 {
		t.Fatal("expected non-zero autoincrement ID")
	}

	got, err := rss.Get(id)
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
	if err := rss.Update(got); err != nil {
		t.Fatalf("Update: %v", err)
	}

	final, err := rss.Get(id)
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

// TestRuleSetStorage_AddPreservesRuleFields verifies all Rule fields survive
// the Add → Get cycle.
func TestRuleSetStorage_AddPreservesRuleFields(t *testing.T) {
	db := openExternalTestDB(t)
	rss := storage.NewRuleSetStorage(db)

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

	if err := rss.Add(storage.RuleSet{
		Name:         "Fields Test",
		Rules:        rules,
		ShouldHitAll: true,
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	all, _ := rss.All()
	got, err := rss.Get(all[0].Id)
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

// TestRuleSetStorage_RemovePurgesOrphanedGroupIds verifies that removing a
// DriverGroup also removes it from RuleSet associations via FK CASCADE.
func TestRuleSetStorage_RemovePurgesOrphanedGroupIds(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)
	rss := storage.NewRuleSetStorage(db)

	groupID := addTestGroup(t, dgs, storage.DriverGroup{
		Name:    "Temp Group",
		Type:    storage.Network,
		Drivers: []*storage.Driver{{Name: "D1"}},
	})
	keepID := addTestGroup(t, dgs, storage.DriverGroup{
		Name: "Keep Group",
		Type: storage.Display,
	})

	if err := rss.Add(storage.RuleSet{
		Name:           "Purge Test",
		DriverGroupIds: []uint{groupID, keepID},
	}); err != nil {
		t.Fatalf("rss.Add: %v", err)
	}

	if err := dgs.Remove(groupID); err != nil {
		t.Fatalf("dgs.Remove: %v", err)
	}

	all, err := rss.All()
	if err != nil {
		t.Fatalf("rss.All: %v", err)
	}
	rs := all[0]
	if containsUint(rs.DriverGroupIds, groupID) {
		t.Errorf("groupID %d should have been purged, still in %v", groupID, rs.DriverGroupIds)
	}
	if !containsUint(rs.DriverGroupIds, keepID) {
		t.Errorf("keepID %d should still be present, not in %v", keepID, rs.DriverGroupIds)
	}
}

// TestRuleSetStorage_UpdateRuleSet verifies Update persists rule changes.
func TestRuleSetStorage_UpdateRuleSet(t *testing.T) {
	db := openExternalTestDB(t)
	rss := storage.NewRuleSetStorage(db)

	if err := rss.Add(storage.RuleSet{
		Name:  "Before Update",
		Rules: []storage.Rule{{Source: storage.Memory, Operator: storage.Equal, Values: []string{"8GB"}}},
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	all, _ := rss.All()
	original, err := rss.Get(all[0].Id)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	original.Name = "After Update"
	original.Rules = []storage.Rule{
		{Source: storage.Storage, Operator: storage.NotContain, Values: []string{"USB"}},
		{Source: storage.Motherboard, Operator: storage.Regex, Values: []string{`^ASUS.*`}},
	}

	if err := rss.Update(original); err != nil {
		t.Fatalf("Update: %v", err)
	}

	persisted, err := rss.Get(all[0].Id)
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

// TestRuleSetStorage_RemoveNonExistent verifies that removing a non-existent
// ID returns an error.
func TestRuleSetStorage_RemoveNonExistent(t *testing.T) {
	db := openExternalTestDB(t)
	rss := storage.NewRuleSetStorage(db)

	if err := rss.Remove(9999); err == nil {
		t.Error("expected error when removing non-existent id, got nil")
	}
}

// TestRuleSetStorage_AllReturnsAllAdded verifies All() returns every added ruleset.
func TestRuleSetStorage_AllReturnsAllAdded(t *testing.T) {
	db := openExternalTestDB(t)
	rss := storage.NewRuleSetStorage(db)

	names := []string{"Alpha", "Beta", "Gamma"}
	for _, name := range names {
		if err := rss.Add(storage.RuleSet{Name: name}); err != nil {
			t.Fatalf("Add %s: %v", name, err)
		}
	}

	all, err := rss.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(all) != len(names) {
		t.Errorf("All() count: got %d, want %d", len(all), len(names))
	}
}

// TestRuleSetStorage_DriverGroupIds_RoundTrip verifies that DriverGroupIds
// are persisted and returned correctly via the M2M join table.
func TestRuleSetStorage_DriverGroupIds_RoundTrip(t *testing.T) {
	db := openExternalTestDB(t)
	dgs := storage.NewDriverGroupStorage(db)
	rss := storage.NewRuleSetStorage(db)

	g1 := addTestGroup(t, dgs, storage.DriverGroup{Name: "Net", Type: storage.Network})
	g2 := addTestGroup(t, dgs, storage.DriverGroup{Name: "Disp", Type: storage.Display})

	if err := rss.Add(storage.RuleSet{
		Name:           "IDs Test",
		DriverGroupIds: []uint{g1, g2},
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	all, _ := rss.All()
	rs := all[0]
	if !containsUint(rs.DriverGroupIds, g1) || !containsUint(rs.DriverGroupIds, g2) {
		t.Errorf("DriverGroupIds %v missing g1=%d or g2=%d", rs.DriverGroupIds, g1, g2)
	}
}
