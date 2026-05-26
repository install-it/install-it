package storage

import (
	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenDB(path string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
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
	})
	return m.Migrate()
}
