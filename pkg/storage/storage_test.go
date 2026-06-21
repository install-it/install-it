// Package storage_test provides external black-box tests for the storage package.
package storage_test

import (
	"testing"

	"install-it/pkg/storage"
)

// openExternalTestDB creates an in-memory SQLite database and runs all migrations.
func openExternalTestDB(t *testing.T) *storage.Database {
	t.Helper()
	db, err := storage.Open(":memory:")
	if err != nil {
		t.Fatalf("openExternalTestDB: %v", err)
	}
	if err := db.Migrate(); err != nil {
		t.Fatalf("openExternalTestDB Migrate: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close database: %v", err)
		}
	})
	return db
}

// containsUint reports whether slice contains v.
func containsUint(slice []uint, v uint) bool {
	for _, u := range slice {
		if u == v {
			return true
		}
	}
	return false
}

// addTestGroup adds a group and returns its autoincrement ID via All().
func addTestGroup(t *testing.T, dgs *storage.DriverGroupStorage, group storage.DriverGroup) uint {
	t.Helper()
	if err := dgs.Add(group); err != nil {
		t.Fatalf("addTestGroup Add: %v", err)
	}
	all, err := dgs.All()
	if err != nil {
		t.Fatalf("addTestGroup All: %v", err)
	}
	return all[len(all)-1].Id
}
