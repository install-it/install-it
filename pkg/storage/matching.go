package storage

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"install-it/pkg/sysinfo"
)

type RuleSource string

const (
	Cpu         RuleSource = "cpu"
	Motherboard RuleSource = "motherboard"
	Gpu         RuleSource = "gpu"
	Memory      RuleSource = "memory"
	Nic         RuleSource = "nic"
	Storage     RuleSource = "storage"
)

type RuleOperator string

const (
	Contain    RuleOperator = "contain"
	NotContain RuleOperator = "not_contain"
	Equal      RuleOperator = "equal"
	NotEqual   RuleOperator = "not_equal"
	Regex      RuleOperator = "regex"
)

type Rule struct {
	Source          RuleSource   `json:"source"`
	Operator        RuleOperator `json:"operator"`
	IsCaseSensitive bool         `json:"is_case_sensitive"`
	ShouldHitAll    bool         `json:"should_hit_all"`
	Values          []string     `json:"values"`
}

type RuleSet struct {
	Id             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string         `json:"name"`
	Rules          []Rule         `json:"rules" gorm:"serializer:json"`
	ShouldHitAll   bool           `json:"should_hit_all"`
	DriverGroups   []*DriverGroup `json:"-" gorm:"many2many:rule_set_driver_groups;constraint:OnDelete:CASCADE"`
	DriverGroupIds []uint         `json:"driver_group_ids" gorm:"-"`
}

func populateDriverGroupIds(rs *RuleSet) {
	rs.DriverGroupIds = make([]uint, len(rs.DriverGroups))
	for i, dg := range rs.DriverGroups {
		rs.DriverGroupIds[i] = dg.Id
	}
}

type MatchRuleStorage struct {
	DB *gorm.DB
}

func NewMatchRuleStorage(db *gorm.DB) *MatchRuleStorage {
	return &MatchRuleStorage{DB: db}
}

func (s *MatchRuleStorage) All() ([]RuleSet, error) {
	var ruleSets []*RuleSet
	if err := s.DB.Preload("DriverGroups").Find(&ruleSets).Error; err != nil {
		return nil, err
	}
	result := make([]RuleSet, len(ruleSets))
	for i, rs := range ruleSets {
		populateDriverGroupIds(rs)
		result[i] = *rs
	}
	return result, nil
}

func (s *MatchRuleStorage) Get(id uint) (RuleSet, error) {
	var rs RuleSet
	result := s.DB.Preload("DriverGroups").First(&rs, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return RuleSet{}, errors.New("store: no item with the same ID was found")
		}
		return RuleSet{}, result.Error
	}
	populateDriverGroupIds(&rs)
	return rs, nil
}

func (s *MatchRuleStorage) Add(ruleSet RuleSet) error {
	ruleSet.DriverGroups = idsToDriverGroups(ruleSet.DriverGroupIds)
	return s.DB.Omit("DriverGroups.*").Create(&ruleSet).Error
}

func (s *MatchRuleStorage) Update(ruleSet RuleSet) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("DriverGroups").Save(&ruleSet).Error; err != nil {
			return err
		}
		groups := idsToDriverGroups(ruleSet.DriverGroupIds)
		return tx.Model(&ruleSet).Association("DriverGroups").Replace(groups)
	})
}

func (s *MatchRuleStorage) Remove(id uint) error {
	result := s.DB.Delete(&RuleSet{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("store: no item with the same ID was found")
	}
	return nil
}

func (s *MatchRuleStorage) Clone(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var original RuleSet
		if err := tx.Preload("DriverGroups").First(&original, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("store: no item with the same ID was found")
			}
			return err
		}

		newRS := RuleSet{
			Name:         original.Name + " (copy)",
			Rules:        original.Rules,
			ShouldHitAll: original.ShouldHitAll,
			DriverGroups: original.DriverGroups,
		}
		return tx.Omit("DriverGroups.*").Create(&newRS).Error
	})
}

// MatchedGroupIds evaluates all RuleSets against live hardware info and returns
// the IDs of DriverGroups that should be selected.
func (s *MatchRuleStorage) MatchedGroupIds() ([]uint, error) {
	hw, err := buildHardwareMap()
	if err != nil {
		return nil, err
	}

	var ruleSets []*RuleSet
	if err := s.DB.Preload("DriverGroups").Find(&ruleSets).Error; err != nil {
		return nil, err
	}

	seen := make(map[uint]bool)
	matched := []uint{}

	for _, rs := range ruleSets {
		results := make([]bool, len(rs.Rules))
		for i, rule := range rs.Rules {
			inputs := hw[rule.Source]
			results[i] = anyMatchesRule(rule, inputs)
		}

		var hit bool
		if rs.ShouldHitAll {
			hit = allTrue(results)
		} else {
			hit = anyTrue(results)
		}
		if hit {
			for _, dg := range rs.DriverGroups {
				if !seen[dg.Id] {
					matched = append(matched, dg.Id)
					seen[dg.Id] = true
				}
			}
		}
	}

	return matched, nil
}

// idsToDriverGroups converts a slice of IDs to DriverGroup pointer stubs for GORM associations.
func idsToDriverGroups(ids []uint) []*DriverGroup {
	groups := make([]*DriverGroup, len(ids))
	for i, id := range ids {
		groups[i] = &DriverGroup{Id: id}
	}
	return groups
}

// buildHardwareMap queries WMI and returns formatted strings per RuleSource,
// matching the formatting used in frontend/src/utils/index.ts getHardware().
func buildHardwareMap() (map[RuleSource][]string, error) {
	si := sysinfo.SysInfo{}
	hw := make(map[RuleSource][]string)

	if cpus, err := si.CpuInfo(); err == nil {
		for _, v := range cpus {
			hw[Cpu] = append(hw[Cpu], v.Name)
		}
	}

	if gpus, err := si.GpuInfo(); err == nil {
		for _, v := range gpus {
			gb := int(math.Round(float64(v.AdapterRAM) / math.Pow(1024, 3)))
			hw[Gpu] = append(hw[Gpu], fmt.Sprintf("%s (%dGB)", v.Name, gb))
		}
	}

	if mems, err := si.MemoryInfo(); err == nil {
		for _, v := range mems {
			gb := float64(v.Capacity) / math.Pow(1024, 3)
			gbStr := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.10g", gb), "0"), ".")
			hw[Memory] = append(hw[Memory], fmt.Sprintf("%s %s %sGB %dMHz",
				v.Manufacturer, strings.TrimSpace(v.PartNumber), gbStr, v.Speed))
		}
	}

	if boards, err := si.MotherboardInfo(); err == nil {
		for _, v := range boards {
			hw[Motherboard] = append(hw[Motherboard], fmt.Sprintf("%s %s", v.Manufacturer, v.Product))
		}
	}

	if nics, err := si.NicInfo(); err == nil {
		for _, v := range nics {
			hw[Nic] = append(hw[Nic], v.Name)
		}
	}

	if disks, err := si.DiskInfo(); err == nil {
		for _, v := range disks {
			gb := int(math.Round(float64(v.Size) / math.Pow(1024, 3)))
			hw[Storage] = append(hw[Storage], fmt.Sprintf("%s (%dGB)", v.Model, gb))
		}
	}

	return hw, nil
}

// anyMatchesRule returns true if any hardware input matches the rule.
func anyMatchesRule(rule Rule, inputs []string) bool {
	for _, input := range inputs {
		if testRule(rule, input) {
			return true
		}
	}
	return false
}

// testRule tests whether a single input string satisfies the rule,
// matching the logic of testMatchRule() in frontend/src/utils/index.ts.
func testRule(rule Rule, input string) bool {
	cmpInput := input
	values := make([]string, len(rule.Values))
	copy(values, rule.Values)

	if !rule.IsCaseSensitive {
		cmpInput = strings.ToLower(cmpInput)
		for i, v := range values {
			values[i] = strings.ToLower(v)
		}
	}

	hits := make([]bool, len(values))
	for i, v := range values {
		switch rule.Operator {
		case Contain:
			hits[i] = strings.Contains(cmpInput, v)
		case NotContain:
			hits[i] = !strings.Contains(cmpInput, v)
		case Equal:
			hits[i] = cmpInput == v
		case NotEqual:
			hits[i] = cmpInput != v
		case Regex:
			pattern := v // already lowercased if !IsCaseSensitive
			if !rule.IsCaseSensitive {
				pattern = "(?i)" + pattern
			}
			re, err := regexp.Compile(pattern)
			if err == nil {
				hits[i] = re.MatchString(cmpInput)
			}
		}
	}

	if rule.ShouldHitAll {
		return allTrue(hits)
	}
	return anyTrue(hits)
}

func allTrue(bs []bool) bool {
	for _, b := range bs {
		if !b {
			return false
		}
	}
	return true
}

func anyTrue(bs []bool) bool {
	for _, b := range bs {
		if b {
			return true
		}
	}
	return false
}
