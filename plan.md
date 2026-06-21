# Refactor Plan: pkg/storage

## Context

The `pkg/storage` package has become overly complex and is approaching an unmaintainable state. This plan documents the agreed-upon redesign to make it simple, clean, maintainable, and production-ready while following idiomatic Go conventions.

## Design Philosophy (from user)

- Avoid premature abstraction
- Simplicity over complexity
- Implicit over explicit
- Follow Go conventions / idiomatic Go
- Avoid overengineering

## Key Decisions (confirmed with user)

1. **Scope**: No limit on what files can change, as long as it makes sense. Breaking public API is acceptable.
2. **Database**: SQLite-only. No need for database-agnostic architecture.
3. **Repository pattern**: No heavy repository interfaces. Concrete structs for data access. In-memory SQLite for tests instead of mocks. Interfaces only at integration boundaries (HardwareQuerier).
4. **Matching logic**: Extract to separate `pkg/matching` package. Storage should not know about hardware enumeration.
5. **Transactions**: Consistent per-method transactions (Option B). Every mutation wraps in `db.Transaction(...)`. No Unit of Work or transaction manager.
6. **AppSetting**: Keep file-based (JSON), but inline the file I/O directly. Delete the `Store` interface, `FileStore`, and `MemoryStore` entirely.
7. **Package structure**: Flat packages, canonical Go approach. No sub-packages within storage, no separate models/domain package.
8. **Testing**: In-memory SQLite (`:memory:`) for storage tests. Fake `HardwareQuerier` for matching tests. Pure unit tests for rule evaluation functions.
9. **Error handling**: Sentinel `ErrNotFound` + `fmt.Errorf` wrapping. No typed error package.
10. **Hardware errors**: Silent degradation is intended behavior (buildHardwareMap swallows errors). Preserve this.
11. **MemoryStore**: Only used for tests. Remove from production code entirely. Tests can define local helpers if needed.
12. **HardwareQuerier interface**: Returns formatted `map[storage.RuleSource][]string` (not raw WMI types) so the matching package doesn't import sysinfo in its interface. WMI implementation lives in `pkg/matching/hardware.go`.
13. **WMI refactor**: User has plans to refactor WMI queries to work without drivers. Out of scope for this refactor, but the HardwareQuerier interface design should not block it.

## Current State Summary

### Files in pkg/storage (non-test, ~661 lines total)
- `db.go` (53 lines) — `OpenDB(path) (*gorm.DB, error)`, `RunMigrations(db *gorm.DB) error`
- `store.go` (93 lines) — `Store` interface, `FileStore`, `MemoryStore` (to be deleted)
- `driver.go` (254 lines) — `DriverGroup`, `Driver` models + `DriverGroupStorage` (CRUD + Clone + MoveBehind)
- `matching.go` (310 lines) — `RuleSet`, `Rule` models + `MatchRuleStorage` (CRUD + Clone + MatchedGroupIds + rule evaluation + WMI hardware queries)
- `app_setting.go` (50 lines) — `AppSetting` model + `AppSettingStorage` (uses Store interface)

### Key Problems Being Fixed
1. **No layering**: matching.go mixes data access, domain logic, and WMI infrastructure
2. **No interfaces for testability**: DriverGroupStorage and MatchRuleStorage embed `*gorm.DB` directly
3. **Exported mutable DB fields**: `groupStorage.DB = db` in main.go OnAfterImport hook (fragile)
4. **Inconsistent transactions**: Add methods not transactional (race on position calc), Remove in MatchRuleStorage not transactional
5. **Inconsistent error handling**: Duplicated error strings, no wrapping, some errors silently ignored
6. **MatchedGroupIds untestable**: Hardcodes `sysinfo.SysInfo{}` with no injection point
7. **Clone methods untested**: Complex transactional logic with zero test coverage
8. **Store abstraction overengineered**: Serves single consumer (AppSettingStorage)
9. **FileStore.Write uses 0777**: Security issue
10. **testRule recompiles regex every call**: Performance issue

### Permission Fix (commit df3e317, must preserve)
The fix added `OnBeforeBackup`/`OnAfterImport` hooks to Porter struct in `pkg/porter/porter.go`. In `main.go`, `OnAfterImport` reopens the DB and reassigns `groupStorage.DB = db` / `matchRuleStorage.DB = db`. The refactor improves this: repositories hold `*Database` (not `*gorm.DB`), and `Database.Reopen()` swaps the internal pointer — no field mutation needed.

### Wails Binding Constraint
This is a Wails app. `main.go` binds concrete struct instances to the frontend via `Bind: []interface{}{...}`. The frontend calls Go methods through auto-generated TypeScript stubs (`wailsjs/go/storage/*.js`). The generated `models.ts` mirrors Go struct JSON tags exactly. Constraints:
- JSON tags on all models MUST be preserved exactly
- Wails-bound method signatures must stay stable or be coordinated with `wails generate`
- `MatchedGroupIds` moves from `MatchRuleStorage` to `matching.Matcher` — frontend call site in `pages/index.vue:42` changes from `matchRuleStorage.MatchedGroupIds()` to `matcher.MatchedGroupIds()`
- After refactor, run `wails generate` to regenerate TypeScript bindings

### Frontend Usage (for context)
- `main.go:118-151` binds: AppSettingStorage, groupStorage, matchRuleStorage, porter.Porter, sysinfo.SysInfo
- Frontend Pinia stores: `useAppSettingStore.ts`, `useDriverGroupStore.ts`, `useMatchRuleStore.ts`
- App.vue:26-50 loads all data on startup via Promise.all
- Driver CRUD in `DriverFormComponent.vue:69-91`, list in `pages/drivers/index.vue`
- Match rule CRUD in `MatchRuleFormComponent.vue:49-58`, list in `pages/match-rules/index.vue`
- Auto-match in `pages/index.vue:42-60` calls `matchRuleStorage.MatchedGroupIds()`
- Settings in `pages/settings.vue:33-45`

## Target Architecture

### Package Structure
```
pkg/
  storage/
    db.go           — Database struct (Open, Close, Reopen, Migrate)
    driver.go       — DriverGroup, Driver models + DriverGroupStorage
    ruleset.go      — RuleSet, Rule, RuleSource, RuleOperator models + RuleSetStorage (CRUD only)
    appsetting.go   — AppSetting + AppSettingStorage (inline file I/O, no Store interface)
    errors.go       — ErrNotFound sentinel
  matching/
    matching.go     — Matcher struct, MatchedGroupIds, rule evaluation (testRule, anyMatchesRule, allTrue, anyTrue)
    hardware.go     — HardwareQuerier interface + WMIHardwareQuerier implementation
```

### Dependency Graph
```
pkg/matching  →  pkg/storage   (for RuleSet, RuleSource, DriverGroup types)
pkg/matching  →  pkg/sysinfo    (only in WMIHardwareQuerier implementation)
pkg/storage   →  (nothing new — no longer imports sysinfo)
main.go       →  pkg/storage, pkg/matching  (wires everything)
```
No circular dependencies. Storage is pure data access. Matching is domain logic + hardware querying.

### Key Types

#### pkg/storage/db.go
```go
type Database struct {
    db   *gorm.DB
    path string
}

func Open(path string) (*Database, error)       // gorm.Open + PRAGMA foreign_keys = ON
func (d *Database) DB() *gorm.DB                 // accessor, not exported field
func (d *Database) Close() error                // closes underlying sqlDB
func (d *Database) Reopen() error               // close + reopen at same path + PRAGMA
func (d *Database) Migrate() error               // runs gormigrate migrations (current logic from RunMigrations)
```
Repositories hold `*Database`, not `*gorm.DB`. `Reopen()` swaps internal `*gorm.DB` pointer — all repositories using `db.DB()` see the new connection automatically.

#### pkg/storage/driver.go
```go
// Models unchanged (preserve all JSON tags and GORM tags exactly)
type DriverType string  // Network, Display, Miscellaneous
type DriverGroup struct { ... }  // same fields, same tags
type Driver struct { ... }       // same fields, same tags

// Storage holds *Database instead of *gorm.DB
type DriverGroupStorage struct {
    db *Database
}
func NewDriverGroupStorage(db *Database) *DriverGroupStorage

// Methods unchanged in signature (Wails binding contract):
// All, Get, Add, Update, Remove, Clone, MoveBehind
// Changes: all use s.db.DB() instead of s.DB
// Add wraps in transaction (fixes race on position calc)
// All mutations wrap in transaction consistently
// Error handling: use ErrNotFound sentinel + fmt.Errorf wrapping
```

#### pkg/storage/ruleset.go (renamed from matching.go, CRUD only)
```go
// Models unchanged (preserve all JSON tags and GORM tags exactly)
type RuleSource string   // Cpu, Motherboard, Gpu, Memory, Nic, Storage
type RuleOperator string // Contain, NotContain, Equal, NotEqual, Regex
type Rule struct { ... }       // same fields, same tags
type RuleSet struct { ... }    // same fields, same tags

// Storage holds *Database
type RuleSetStorage struct {
    db *Database
}
func NewRuleSetStorage(db *Database) *RuleSetStorage

// Methods (CRUD only, no MatchedGroupIds):
// All, Get, Add, Update, Remove, Clone
// Changes: all use s.db.DB() instead of s.DB
// Add wraps in transaction
// Remove wraps in transaction
// Error handling: use ErrNotFound sentinel + fmt.Errorf wrapping
```
Note: Type renamed from `MatchRuleStorage` to `RuleSetStorage` for clarity. Frontend binding name changes — coordinate with `wails generate`.

#### pkg/storage/appsetting.go (simplified, no Store interface)
```go
// Models unchanged (preserve all JSON tags)
type AppSetting struct { ... }       // same fields, same tags
type SuccessAction string            // Nothing, Shutdown, Reboot, Firmware

type AppSettingStorage struct {
    Path    string
    setting AppSetting
}

func (s *AppSettingStorage) All() (AppSetting, error) {
    // os.ReadFile(s.Path); if os.ErrNotExist, write defaults
    // json.Unmarshal directly
    // Check Read error (current code silently ignores it — fix this)
}

func (s *AppSettingStorage) Update(v AppSetting) (AppSetting, error) {
    s.setting = v
    // json.Marshal + os.WriteFile(s.Path, bytes, 0644) — fix 0777 to 0644
}
```

#### pkg/storage/errors.go
```go
var ErrNotFound = errors.New("storage: not found")
```
Replace all 5 instances of `errors.New("store: no item with the same ID was found")` with `ErrNotFound`. Use `fmt.Errorf("driver group: %w", err)` for wrapping.

#### pkg/matching/matching.go
```go
type RuleSetReader interface {
    All() ([]storage.RuleSet, error)
}

type Matcher struct {
    rules    RuleSetReader
    hardware HardwareQuerier
}

func NewMatcher(rules RuleSetReader, hw HardwareQuerier) *Matcher
func (m *Matcher) MatchedGroupIds() ([]uint, error)

// Pure functions moved from storage/matching.go (currently unexported, make them unexported here too):
// anyMatchesRule, testRule, allTrue, anyTrue
// These are pure functions — test directly with controlled inputs
```

#### pkg/matching/hardware.go
```go
type HardwareQuerier interface {
    HardwareMap() (map[storage.RuleSource][]string, error)
}

type WMIHardwareQuerier struct{}

func (WMIHardwareQuerier) HardwareMap() (map[storage.RuleSource][]string, error) {
    // Current buildHardwareMap logic moves here verbatim
    // Imports sysinfo, formats hardware data into strings
    // Silent error degradation preserved (user confirmed this is intended)
}
```
`sysinfo.SysInfo` does NOT directly satisfy `HardwareQuerier` — `WMIHardwareQuerier` wraps it. This is intentional: the interface returns formatted strings, not raw WMI types, so it's stable when WMI is refactored later.

### main.go Changes

```go
// Package-level vars change:
var (
    db            *storage.Database
    groupStorage  *storage.DriverGroupStorage
    ruleSetStorage *storage.RuleSetStorage   // renamed from matchRuleStorage
    matcher       *matching.Matcher          // NEW
)

// init() unchanged

// main():
db, err = storage.Open(filepath.Join(dirConf, "data.db"))
if err != nil { panic(err) }
if err := db.Migrate(); err != nil { panic(err) }

groupStorage = storage.NewDriverGroupStorage(db)
ruleSetStorage = storage.NewRuleSetStorage(db)
matcher = matching.NewMatcher(ruleSetStorage, matching.WMIHardwareQuerier{})

// Bind array changes:
// - AppSettingStorage: &storage.AppSettingStorage{Path: filepath.Join(dirConf, "setting.json")}
// - groupStorage (same)
// - ruleSetStorage (replaces matchRuleStorage)
// - matcher (NEW — for MatchedGroupIds)
// - Porter hooks simplified:

&porter.Porter{
    ...
    OnBeforeBackup: func() error {
        return db.Close()  // was: sqlDB, err := db.DB(); ...; sqlDB.Close()
    },
    OnAfterImport: func() error {
        if err := db.Reopen(); err != nil { return err }
        return db.Migrate()
        // No more: groupStorage.DB = db; matchRuleStorage.DB = db
    },
},

// EnumBind: RuleSource and RuleOperator stay in storage package (they're model types)
```

### Frontend Changes (after wails generate)
- `matchRuleStorage` → `ruleSetStorage` (or whatever name Wails generates from the bound struct)
- `matchRuleStorage.MatchedGroupIds()` → `matcher.MatchedGroupIds()` (in `pages/index.vue:42`)
- TypeScript types in `models.ts` regenerate automatically
- Pinia store `useMatchRuleStore.ts` may need rename for consistency (optional)

### Test Changes

#### Files to delete
- `store_test.go` (175 lines) — FileStore/MemoryStore tests, no longer needed
- `storage_test.go` (278 lines) — external black-box tests, will be replaced by per-storage test files

#### Files to update
- `driver_test.go` — change `openTestDB` to use `storage.Open(":memory:")` + `db.Migrate()`, pass `*Database` to `NewDriverGroupStorage`
- `driver_group_storage_test.go` — same DB setup change
- `match_rule_storage_test.go` → rename to `ruleset_test.go` — same DB setup change, update type names
- `app_setting_test.go` — remove `rawJSONStore` mock, use real temp file via `t.TempDir()`
- `fuzz_test.go` — update `FuzzRuleSet_JSONRoundtrip` for new package path, `FuzzFileStore_Read` removed (no more FileStore)

#### New test files
- `pkg/matching/matching_test.go` — unit tests for `testRule`, `anyMatchesRule`, `allTrue`, `anyTrue` (pure functions, currently untested)
- `pkg/matching/hardware_test.go` — test `MatchedGroupIds` with fake `HardwareQuerier` + fake `RuleSetReader`
- `pkg/storage/db_test.go` — test `Open`, `Close`, `Reopen`, `Migrate` error paths (currently untested)

#### Test helpers to consolidate
- Single `openTestDB(t)` helper using `:memory:` SQLite — either in a shared `testutil_test.go` or duplicated per package (Go idiom: duplicate is fine for test helpers)
- `addGroup` / `addTestGroup` helpers stay, consolidated

### Migration Steps (execution order)

1. **Create `pkg/storage/errors.go`** — `ErrNotFound` sentinel
2. **Create `pkg/storage/db.go`** — `Database` struct with Open/Close/Reopen/Migrate/DB()
3. **Update `pkg/storage/driver.go`** — `DriverGroupStorage` holds `*Database`, use `s.db.DB()`, wrap `Add` in transaction, use `ErrNotFound`
4. **Create `pkg/storage/ruleset.go`** — move `RuleSet`, `Rule`, `RuleSource`, `RuleOperator` + CRUD from `matching.go`, rename `MatchRuleStorage` → `RuleSetStorage`, hold `*Database`, wrap `Add`/`Remove` in transaction, use `ErrNotFound`
5. **Update `pkg/storage/appsetting.go`** — inline file I/O, remove Store dependency, fix 0777→0644, check Read error
6. **Delete `pkg/storage/store.go`** — Store interface, FileStore, MemoryStore all removed
7. **Delete `pkg/storage/matching.go`** — all content moved to ruleset.go (CRUD) and pkg/matching (evaluation)
8. **Create `pkg/matching/matching.go`** — `Matcher` struct, `MatchedGroupIds`, rule evaluation functions
9. **Create `pkg/matching/hardware.go`** — `HardwareQuerier` interface, `WMIHardwareQuerier` implementation
10. **Update `main.go`** — wire new types, simplify Porter hooks
11. **Update tests** — in-memory SQLite, new test files for matching, delete obsolete test files
12. **Run `go build ./...`** — verify compilation
13. **Run `go test ./...`** — verify all tests pass
14. **Run `wails generate`** — regenerate frontend bindings
15. **Update frontend call site** — `pages/index.vue:42` `matchRuleStorage.MatchedGroupIds()` → `matcher.MatchedGroupIds()`

### What Must NOT Change
- JSON tags on all model structs (frontend contract via Wails models.ts)
- GORM tags on all model structs (database schema compatibility)
- Migration IDs and logic (database compatibility with existing installations)
- Porter hook mechanism (OnBeforeBackup/OnAfterImport) — only the implementation in main.go changes
- Silent hardware error degradation in buildHardwareMap (user confirmed intended)
- AppSettingStorage file-based persistence (user confirmed keep file-based)

### Files Touched
- `pkg/storage/db.go` — rewrite
- `pkg/storage/driver.go` — modify (DB field, transactions, errors)
- `pkg/storage/matching.go` — delete (content split to ruleset.go + pkg/matching)
- `pkg/storage/store.go` — delete
- `pkg/storage/app_setting.go` — rewrite (inline file I/O)
- `pkg/storage/errors.go` — new
- `pkg/storage/ruleset.go` — new (from matching.go CRUD)
- `pkg/matching/matching.go` — new
- `pkg/matching/hardware.go` — new
- `main.go` — modify (wire new types, simplify hooks)
- Test files — update/delete/create as specified
- `frontend/src/pages/index.vue` — one line change (MatchedGroupIds call)
- `frontend/wailsjs/` — regenerated by `wails generate`