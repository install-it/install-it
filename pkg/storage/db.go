package storage

import (
	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	db.Exec("PRAGMA foreign_keys = ON")
	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "2026052601_baseline",
			Migrate: func(tx *gorm.DB) error {
				// Guard for legacy installations that already have raw tables:
				// skip AutoMigrate but let gormigrate record this migration as done.
				if tx.Migrator().HasTable(&DriverGroup{}) {
					return nil
				}
				return tx.AutoMigrate(&DriverGroup{}, &Driver{}, &RuleSet{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&RuleSet{}, &Driver{}, &DriverGroup{})
			},
		},
		{
			ID: "2026052901_m2m_and_uint_pks",
			Migrate: func(tx *gorm.DB) error {
				tx.Exec("PRAGMA foreign_keys = ON")
				// Drop old string-PK schema; no user data to preserve (pre-release)
				tx.Migrator().DropTable("rule_sets", "drivers", "driver_groups")
				return tx.AutoMigrate(&DriverGroup{}, &Driver{}, &RuleSet{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"rule_set_driver_groups", "driver_incompatibles",
					"rule_sets", "drivers", "driver_groups")
			},
		},
	})
	return m.Migrate()
}
