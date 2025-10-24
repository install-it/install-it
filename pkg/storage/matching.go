package storage

import (
	"install-it/pkg/utils"
	"reflect"
	"slices"
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

type RuleType string

const (
	Contain    RuleType = "contain"
	NotContain RuleType = "not_contain"
	Equal      RuleType = "equal"
	NotEqual   RuleType = "not_equal"
	Regex      RuleType = "regex"
)

type Rule struct {
	Source          RuleSource `json:"source"`
	Type            RuleType   `json:"type"`
	IsCaseSensitive bool       `json:"is_case_sensitive"`
	Values          []string   `json:"values"`
}

type RuleSet struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Rules          []Rule   `json:"rules"`
	DriverGroupIds []string `json:"driver_group_ids"`
}

func (r RuleSet) GetId() string { return r.Id }

func (r *RuleSet) SetId(id string) { r.Id = id }

type MatchRuleStorage struct {
	Store    Store
	EventBus *DeleteEventBus
	data     []*RuleSet
}

func NewMatchRuleStorage(store Store, eventBus *DeleteEventBus) *MatchRuleStorage {
	m := &MatchRuleStorage{Store: store, EventBus: eventBus}
	m.RegisterEventHandlers()
	return m
}

func (s *MatchRuleStorage) RegisterEventHandlers() {
	s.EventBus.Subscribe(reflect.TypeFor[DriverGroup]().Name(),
		func(deletedIds []string) error {
			for _, ruleSet := range s.data {
				ruleSet.DriverGroupIds = slices.DeleteFunc(ruleSet.DriverGroupIds, func(id string) bool {
					return slices.Contains(deletedIds, id)
				})
			}
			return s.Store.Write(s.data)
		})
}

func (s *MatchRuleStorage) All() ([]RuleSet, error) {
	if !s.Store.Exist() {
		s.data = []*RuleSet{}
		s.Store.Write(s.data)
	} else {
		s.Store.Read(&s.data)
	}
	return s.copyOfAll(), nil
}

func (s MatchRuleStorage) Get(id string) (RuleSet, error) {
	if ruleSet, err := Get(id, s.data); err != nil {
		return RuleSet{}, err
	} else {
		return *ruleSet, nil
	}
}

func (s *MatchRuleStorage) Add(ruleSet RuleSet) (string, error) {
	if id, err := Create(&ruleSet, &s.data); err != nil {
		return "", err
	} else {
		return id, s.Store.Write(s.data)
	}
}

func (s *MatchRuleStorage) Update(ruleSet RuleSet) (RuleSet, error) {
	if err := Update(&ruleSet, &s.data); err != nil {
		return RuleSet{}, err
	}
	return ruleSet, s.Store.Write(s.data)
}

func (s *MatchRuleStorage) Remove(id string) error {
	if err := Delete(id, &s.data); err != nil {
		return err
	}
	return s.Store.Write(s.data)
}

func (s MatchRuleStorage) copyOfAll() []RuleSet {
	return utils.Map(s.data, func(g *RuleSet) RuleSet { return *g })
}
