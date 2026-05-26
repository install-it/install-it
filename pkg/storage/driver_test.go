package storage

import (
	"path/filepath"
	"testing"

	"gorm.io/gorm"
)

// openTestDB creates an isolated SQLite database in the test's temp directory
// and runs all migrations, giving each test a clean schema.
func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := OpenDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("openTestDB: %v", err)
	}
	if err := RunMigrations(db); err != nil {
		t.Fatalf("openTestDB RunMigrations: %v", err)
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

// ==================== DriverGroup / Driver Struct Tests ====================

func TestDriverGroup_BeforeCreate_GeneratesId(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	id, err := dgs.Add(DriverGroup{Name: "Auto ID Group", Type: Network})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if len(id) != 8 {
		t.Errorf("generated group ID %q should be 8 chars", id)
	}
}

func TestDriver_BeforeCreate_GeneratesId(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	groupId, err := dgs.Add(DriverGroup{
		Name: "Driver ID Group",
		Type: Network,
		Drivers: []*Driver{
			{Name: "D1"},
			{Name: "D2"},
		},
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	group, err := dgs.Get(groupId)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(group.Drivers) != 2 {
		t.Fatalf("expected 2 drivers, got %d", len(group.Drivers))
	}
	ids := make(map[string]bool)
	for _, d := range group.Drivers {
		if len(d.Id) != 8 {
			t.Errorf("driver ID %q should be 8 chars", d.Id)
		}
		if ids[d.Id] {
			t.Errorf("duplicate driver ID: %s", d.Id)
		}
		ids[d.Id] = true
	}
}

// ==================== DriverGroupStorage CRUD ====================

func TestDriverGroupStorage_All_Empty(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	groups, err := dgs.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestDriverGroupStorage_Get_NotFound(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	if _, err := dgs.Get("nonexistent"); err == nil {
		t.Fatal("expected error for nonexistent group, got nil")
	}
}

func TestDriverGroupStorage_Add_And_All(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	id1, _ := dgs.Add(DriverGroup{Name: "G1", Type: Network})
	id2, _ := dgs.Add(DriverGroup{Name: "G2", Type: Display})

	all, err := dgs.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(all))
	}
	// Verify ordering by position (insertion order)
	if all[0].Id != id1 || all[1].Id != id2 {
		t.Errorf("unexpected order: [%s, %s]", all[0].Id, all[1].Id)
	}
}

func TestDriverGroupStorage_Update_ScalarFields(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	id, _ := dgs.Add(DriverGroup{Name: "Original", Type: Network, MutuallyExclusive: false})

	updated, err := dgs.Update(DriverGroup{
		Id:                id,
		Name:              "Updated",
		Type:              Display,
		MutuallyExclusive: true,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "Updated" || updated.Type != Display || !updated.MutuallyExclusive {
		t.Errorf("unexpected updated values: %+v", updated)
	}
}

func TestDriverGroupStorage_Update_RemoveDriverCleansIncompatibles(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	groupId, _ := dgs.Add(DriverGroup{
		Name: "Incompat Group",
		Type: Network,
		Drivers: []*Driver{
			{Name: "D1"},
			{Name: "D2"},
		},
	})

	group, _ := dgs.Get(groupId)
	d1Id := group.Drivers[0].Id
	d2Id := group.Drivers[1].Id

	// Set D1 incompatible with D2
	group.Drivers[0].Incompatibles = []string{d2Id}
	group.Drivers[1].Incompatibles = []string{d1Id}
	dgs.Update(group)

	// Remove D2 — D1's incompatibles should be cleaned
	group.Drivers = group.Drivers[:1]
	updated, err := dgs.Update(group)
	if err != nil {
		t.Fatalf("Update removing D2: %v", err)
	}
	if len(updated.Drivers) != 1 {
		t.Fatalf("expected 1 driver, got %d", len(updated.Drivers))
	}
	if len(updated.Drivers[0].Incompatibles) != 0 {
		t.Errorf("expected empty incompatibles after D2 removal, got %v",
			updated.Drivers[0].Incompatibles)
	}
}

func TestDriverGroupStorage_Update_AddNewDriver(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	groupId, _ := dgs.Add(DriverGroup{
		Name: "G",
		Type: Network,
		Drivers: []*Driver{{Name: "D1"}},
	})

	group, _ := dgs.Get(groupId)
	group.Drivers = append(group.Drivers, &Driver{Name: "D2"})
	updated, err := dgs.Update(group)
	if err != nil {
		t.Fatalf("Update with new driver: %v", err)
	}
	if len(updated.Drivers) != 2 {
		t.Fatalf("expected 2 drivers, got %d", len(updated.Drivers))
	}
	if updated.Drivers[1].Id == "" {
		t.Error("new driver should have a generated ID")
	}
}

func TestDriverGroupStorage_Remove_DeletesGroup(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	id1, _ := dgs.Add(DriverGroup{Name: "G1", Type: Network})
	id2, _ := dgs.Add(DriverGroup{Name: "G2", Type: Display})

	if err := dgs.Remove(id1); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	all, _ := dgs.All()
	if len(all) != 1 || all[0].Id != id2 {
		t.Errorf("expected only G2 remaining, got %v", all)
	}
}

func TestDriverGroupStorage_Remove_NotFound(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	if err := dgs.Remove("nonexistent"); err == nil {
		t.Fatal("expected error removing nonexistent group, got nil")
	}
}

func TestDriverGroupStorage_Remove_CascadeIncompatibles(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	g1Id, _ := dgs.Add(DriverGroup{
		Name:    "G1",
		Type:    Network,
		Drivers: []*Driver{{Name: "D1"}},
	})
	g2Id, _ := dgs.Add(DriverGroup{
		Name:    "G2",
		Type:    Display,
		Drivers: []*Driver{{Name: "D2"}},
	})

	g1, _ := dgs.Get(g1Id)
	g2, _ := dgs.Get(g2Id)
	d1Id := g1.Drivers[0].Id
	d2Id := g2.Drivers[0].Id

	// D2 marks D1 as incompatible
	g2.Drivers[0].Incompatibles = []string{d1Id}
	dgs.Update(g2)

	// Remove G1 — D2's incompatibles must be cleaned
	if err := dgs.Remove(g1Id); err != nil {
		t.Fatalf("Remove G1: %v", err)
	}

	g2After, err := dgs.Get(g2Id)
	if err != nil {
		t.Fatalf("Get G2 after removal: %v", err)
	}
	for _, d := range g2After.Drivers {
		if d.Id == d2Id {
			for _, inc := range d.Incompatibles {
				if inc == d1Id {
					t.Errorf("d1Id still in D2 Incompatibles after G1 removed")
				}
			}
		}
	}
}

func TestDriverGroupStorage_IndexOf(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	id1, _ := dgs.Add(DriverGroup{Name: "G1", Type: Network})
	id2, _ := dgs.Add(DriverGroup{Name: "G2", Type: Display})
	id3, _ := dgs.Add(DriverGroup{Name: "G3", Type: Miscellaneous})

	idx, err := dgs.IndexOf(id2)
	if err != nil || idx != 1 {
		t.Errorf("IndexOf G2: expected 1, got %d (err=%v)", idx, err)
	}
	idx, err = dgs.IndexOf(id1)
	if err != nil || idx != 0 {
		t.Errorf("IndexOf G1: expected 0, got %d (err=%v)", idx, err)
	}
	idx, err = dgs.IndexOf(id3)
	if err != nil || idx != 2 {
		t.Errorf("IndexOf G3: expected 2, got %d (err=%v)", idx, err)
	}
}

func TestDriverGroupStorage_IndexOf_NotFound(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	if _, err := dgs.IndexOf("nonexistent"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDriverGroupStorage_MoveBehind_Forward(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	ids := make([]string, 4)
	for i, name := range []string{"A", "B", "C", "D"} {
		ids[i], _ = dgs.Add(DriverGroup{Name: name, Type: Network})
	}

	// MoveBehind A (index 0) behind index 2 → A lands after C: [B, C, D, A]
	// but the existing algorithm bubbles: [A,B,C,D] → [B,A,C,D] → [B,C,A,D] → [B,C,D,A]
	result, err := dgs.MoveBehind(ids[0], 2)
	if err != nil {
		t.Fatalf("MoveBehind: %v", err)
	}
	if len(result) != 4 {
		t.Fatalf("expected 4 groups, got %d", len(result))
	}
	if result[3].Id != ids[0] {
		t.Errorf("expected A at index 3, got %s", result[3].Id)
	}
	if result[0].Id != ids[1] {
		t.Errorf("expected B at index 0, got %s", result[0].Id)
	}
}

func TestDriverGroupStorage_MoveBehind_OutOfBounds(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	id, _ := dgs.Add(DriverGroup{Name: "G1", Type: Network})
	dgs.Add(DriverGroup{Name: "G2", Type: Display})

	if _, err := dgs.MoveBehind(id, 5); err == nil {
		t.Fatal("expected error for out-of-bounds index")
	}
}

func TestDriverGroupStorage_MoveBehind_NegativeIndex(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	id, _ := dgs.Add(DriverGroup{Name: "G1", Type: Network})
	dgs.Add(DriverGroup{Name: "G2", Type: Display})

	result, err := dgs.MoveBehind(id, -1)
	if err != nil {
		t.Fatalf("MoveBehind(-1): %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 groups unchanged, got %d", len(result))
	}
}

func TestDriverGroupStorage_CompleteWorkflow(t *testing.T) {
	db := openTestDB(t)
	dgs := NewDriverGroupStorage(db)

	id1, _ := dgs.Add(DriverGroup{Name: "Network", Type: Network,
		Drivers: []*Driver{{Name: "Net Driver"}}})
	id2, _ := dgs.Add(DriverGroup{Name: "Display", Type: Display,
		Drivers: []*Driver{{Name: "GPU Driver"}}})

	all, _ := dgs.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(all))
	}

	g, _ := dgs.Get(id1)
	g.Name = "Updated Network"
	dgs.Update(g)

	g, _ = dgs.Get(id1)
	if g.Name != "Updated Network" {
		t.Errorf("expected updated name, got %s", g.Name)
	}

	dgs.MoveBehind(id1, 0)

	all, _ = dgs.All()
	if all[1].Id != id1 {
		t.Errorf("expected id1 at index 1 after MoveBehind, got %s", all[1].Id)
	}

	dgs.Remove(id2)
	all, _ = dgs.All()
	if len(all) != 1 {
		t.Errorf("expected 1 group after remove, got %d", len(all))
	}
}
