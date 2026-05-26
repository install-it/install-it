package storage

import (
	"errors"

	"gorm.io/gorm"
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
	Id             string   `json:"id" gorm:"primaryKey"`
	Name           string   `json:"name"`
	Rules          []Rule   `json:"rules" gorm:"serializer:json"`
	ShouldHitAll   bool     `json:"should_hit_all"`
	DriverGroupIds []string `json:"driver_group_ids" gorm:"serializer:json"`
}

func (r *RuleSet) BeforeCreate(tx *gorm.DB) error {
	if r.Id == "" {
		r.Id = generateHexId()
	}
	return nil
}

type MatchRuleStorage struct {
	DB *gorm.DB
}

func NewMatchRuleStorage(db *gorm.DB) *MatchRuleStorage {
	return &MatchRuleStorage{DB: db}
}

func (s *MatchRuleStorage) All() ([]RuleSet, error) {
	var ruleSets []*RuleSet
	if err := s.DB.Find(&ruleSets).Error; err != nil {
		return nil, err
	}
	result := make([]RuleSet, len(ruleSets))
	for i, rs := range ruleSets {
		result[i] = *rs
	}
	return result, nil
}

func (s *MatchRuleStorage) Get(id string) (RuleSet, error) {
	var rs RuleSet
	result := s.DB.First(&rs, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return RuleSet{}, errors.New("store: no item with the same ID was found")
		}
		return RuleSet{}, result.Error
	}
	return rs, nil
}

func (s *MatchRuleStorage) Add(ruleSet RuleSet) (string, error) {
	if err := s.DB.Create(&ruleSet).Error; err != nil {
		return "", err
	}
	return ruleSet.Id, nil
}

func (s *MatchRuleStorage) Update(ruleSet RuleSet) (RuleSet, error) {
	if err := s.DB.Save(&ruleSet).Error; err != nil {
		return RuleSet{}, err
	}
	return ruleSet, nil
}

func (s *MatchRuleStorage) Remove(id string) error {
	result := s.DB.Delete(&RuleSet{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("store: no item with the same ID was found")
	}
	return nil
}
