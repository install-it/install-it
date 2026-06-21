package matching

import (
	"testing"

	"install-it/pkg/storage"
)

// fakeHardwareQuerier returns a fixed hardware map for testing.
type fakeHardwareQuerier struct {
	hw map[storage.RuleSource][]string
}

func (f fakeHardwareQuerier) HardwareMap() (map[storage.RuleSource][]string, error) {
	return f.hw, nil
}

// fakeRuleSetReader returns a fixed list of rule sets for testing.
type fakeRuleSetReader struct {
	ruleSets []storage.RuleSet
}

func (f fakeRuleSetReader) All() ([]storage.RuleSet, error) {
	return f.ruleSets, nil
}

// ==================== testRule ====================

func TestTestRule_Contain(t *testing.T) {
	rule := storage.Rule{
		Source:   storage.Cpu,
		Operator: storage.Contain,
		Values:   []string{"Intel"},
	}
	if !testRule(rule, "Intel Core i7") {
		t.Error("Contain should match substring")
	}
	if testRule(rule, "AMD Ryzen") {
		t.Error("Contain should not match non-substring")
	}
}

func TestTestRule_NotContain(t *testing.T) {
	rule := storage.Rule{
		Source:   storage.Cpu,
		Operator: storage.NotContain,
		Values:   []string{"Intel"},
	}
	if testRule(rule, "Intel Core i7") {
		t.Error("NotContain should not match when substring present")
	}
	if !testRule(rule, "AMD Ryzen") {
		t.Error("NotContain should match when substring absent")
	}
}

func TestTestRule_Equal(t *testing.T) {
	rule := storage.Rule{
		Source:   storage.Cpu,
		Operator: storage.Equal,
		Values:   []string{"Intel Core i7"},
	}
	if !testRule(rule, "Intel Core i7") {
		t.Error("Equal should match exact string")
	}
	if testRule(rule, "Intel Core i7-9700K") {
		t.Error("Equal should not match partial")
	}
}

func TestTestRule_NotEqual(t *testing.T) {
	rule := storage.Rule{
		Source:   storage.Cpu,
		Operator: storage.NotEqual,
		Values:   []string{"Intel Core i7"},
	}
	if !testRule(rule, "AMD Ryzen") {
		t.Error("NotEqual should match different string")
	}
	if testRule(rule, "Intel Core i7") {
		t.Error("NotEqual should not match same string")
	}
}

func TestTestRule_Regex(t *testing.T) {
	rule := storage.Rule{
		Source:   storage.Cpu,
		Operator: storage.Regex,
		Values:   []string{`^Intel.*`},
	}
	if !testRule(rule, "Intel Core i7") {
		t.Error("Regex should match pattern")
	}
	if testRule(rule, "AMD Ryzen") {
		t.Error("Regex should not match non-matching pattern")
	}
}

func TestTestRule_CaseInsensitive(t *testing.T) {
	rule := storage.Rule{
		Source:          storage.Cpu,
		Operator:        storage.Contain,
		IsCaseSensitive: false,
		Values:          []string{"intel"},
	}
	if !testRule(rule, "INTEL Core i7") {
		t.Error("Case-insensitive Contain should match regardless of case")
	}
}

func TestTestRule_CaseSensitive(t *testing.T) {
	rule := storage.Rule{
		Source:          storage.Cpu,
		Operator:        storage.Contain,
		IsCaseSensitive: true,
		Values:          []string{"intel"},
	}
	if testRule(rule, "INTEL Core i7") {
		t.Error("Case-sensitive Contain should not match different case")
	}
	if !testRule(rule, "intel Core i7") {
		t.Error("Case-sensitive Contain should match same case")
	}
}

func TestTestRule_ShouldHitAll_Values(t *testing.T) {
	rule := storage.Rule{
		Source:       storage.Cpu,
		Operator:     storage.Contain,
		ShouldHitAll: true,
		Values:       []string{"Intel", "Core"},
	}
	if !testRule(rule, "Intel Core i7") {
		t.Error("ShouldHitAll with all values present should match")
	}
	if testRule(rule, "Intel i7") {
		t.Error("ShouldHitAll with missing value should not match")
	}
}

func TestTestRule_ShouldHitAny_Values(t *testing.T) {
	rule := storage.Rule{
		Source:       storage.Cpu,
		Operator:     storage.Contain,
		ShouldHitAll: false,
		Values:       []string{"Intel", "Ryzen"},
	}
	if !testRule(rule, "Intel Core i7") {
		t.Error("ShouldHitAny with one value present should match")
	}
	if !testRule(rule, "AMD Ryzen") {
		t.Error("ShouldHitAny with one value present should match")
	}
	if testRule(rule, "AMD Athlon") {
		t.Error("ShouldHitAny with no values present should not match")
	}
}

// ==================== anyMatchesRule ====================

func TestAnyMatchesRule_AnyInputMatches(t *testing.T) {
	rule := storage.Rule{
		Source:   storage.Cpu,
		Operator: storage.Contain,
		Values:   []string{"Intel"},
	}
	inputs := []string{"AMD Ryzen", "Intel Core i7", "ARM"}
	if !anyMatchesRule(rule, inputs) {
		t.Error("should match when at least one input matches")
	}
}

func TestAnyMatchesRule_NoInputMatches(t *testing.T) {
	rule := storage.Rule{
		Source:   storage.Cpu,
		Operator: storage.Contain,
		Values:   []string{"Intel"},
	}
	inputs := []string{"AMD Ryzen", "ARM"}
	if anyMatchesRule(rule, inputs) {
		t.Error("should not match when no input matches")
	}
}

func TestAnyMatchesRule_EmptyInputs(t *testing.T) {
	rule := storage.Rule{
		Source:   storage.Cpu,
		Operator: storage.Contain,
		Values:   []string{"Intel"},
	}
	if anyMatchesRule(rule, nil) {
		t.Error("should not match with empty inputs")
	}
}

// ==================== MatchedGroupIds ====================

func TestMatchedGroupIds_HitAll(t *testing.T) {
	hw := fakeHardwareQuerier{hw: map[storage.RuleSource][]string{
		storage.Cpu: {"Intel Core i7"},
		storage.Gpu: {"NVIDIA RTX 3080 (10GB)"},
	}}
	rules := fakeRuleSetReader{ruleSets: []storage.RuleSet{
		{
			Id:           1,
			Name:         "Intel + NVIDIA",
			ShouldHitAll: true,
			Rules: []storage.Rule{
				{Source: storage.Cpu, Operator: storage.Contain, Values: []string{"Intel"}},
				{Source: storage.Gpu, Operator: storage.Contain, Values: []string{"NVIDIA"}},
			},
			DriverGroups: []*storage.DriverGroup{{Id: 10}, {Id: 20}},
		},
	}}

	m := NewMatcher(rules, hw)
	ids, err := m.MatchedGroupIds()
	if err != nil {
		t.Fatalf("MatchedGroupIds: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 matched group IDs, got %d: %v", len(ids), ids)
	}
	if ids[0] != 10 || ids[1] != 20 {
		t.Errorf("expected [10, 20], got %v", ids)
	}
}

func TestMatchedGroupIds_HitAny(t *testing.T) {
	hw := fakeHardwareQuerier{hw: map[storage.RuleSource][]string{
		storage.Cpu: {"AMD Ryzen"},
	}}
	rules := fakeRuleSetReader{ruleSets: []storage.RuleSet{
		{
			Id:           1,
			Name:         "Intel or AMD",
			ShouldHitAll: false,
			Rules: []storage.Rule{
				{Source: storage.Cpu, Operator: storage.Contain, Values: []string{"Intel"}},
				{Source: storage.Cpu, Operator: storage.Contain, Values: []string{"AMD"}},
			},
			DriverGroups: []*storage.DriverGroup{{Id: 30}},
		},
	}}

	m := NewMatcher(rules, hw)
	ids, err := m.MatchedGroupIds()
	if err != nil {
		t.Fatalf("MatchedGroupIds: %v", err)
	}
	if len(ids) != 1 || ids[0] != 30 {
		t.Errorf("expected [30], got %v", ids)
	}
}

func TestMatchedGroupIds_NoMatch(t *testing.T) {
	hw := fakeHardwareQuerier{hw: map[storage.RuleSource][]string{
		storage.Cpu: {"ARM Cortex"},
	}}
	rules := fakeRuleSetReader{ruleSets: []storage.RuleSet{
		{
			Id:           1,
			Name:         "Intel only",
			ShouldHitAll: true,
			Rules: []storage.Rule{
				{Source: storage.Cpu, Operator: storage.Contain, Values: []string{"Intel"}},
			},
			DriverGroups: []*storage.DriverGroup{{Id: 40}},
		},
	}}

	m := NewMatcher(rules, hw)
	ids, err := m.MatchedGroupIds()
	if err != nil {
		t.Fatalf("MatchedGroupIds: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected 0 matched IDs, got %v", ids)
	}
}

func TestMatchedGroupIds_DeduplicatesGroupIds(t *testing.T) {
	hw := fakeHardwareQuerier{hw: map[storage.RuleSource][]string{
		storage.Cpu: {"Intel Core i7"},
	}}
	rules := fakeRuleSetReader{ruleSets: []storage.RuleSet{
		{
			Id:           1,
			Name:         "Rule 1",
			ShouldHitAll: false,
			Rules: []storage.Rule{
				{Source: storage.Cpu, Operator: storage.Contain, Values: []string{"Intel"}},
			},
			DriverGroups: []*storage.DriverGroup{{Id: 50}, {Id: 60}},
		},
		{
			Id:           2,
			Name:         "Rule 2",
			ShouldHitAll: false,
			Rules: []storage.Rule{
				{Source: storage.Cpu, Operator: storage.Contain, Values: []string{"Core"}},
			},
			DriverGroups: []*storage.DriverGroup{{Id: 60}, {Id: 70}},
		},
	}}

	m := NewMatcher(rules, hw)
	ids, err := m.MatchedGroupIds()
	if err != nil {
		t.Fatalf("MatchedGroupIds: %v", err)
	}
	// 50, 60, 70 — 60 should appear only once
	if len(ids) != 3 {
		t.Fatalf("expected 3 unique IDs, got %d: %v", len(ids), ids)
	}
	seen := make(map[uint]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("duplicate ID: %d", id)
		}
		seen[id] = true
	}
}

func TestMatchedGroupIds_EmptyRuleSets(t *testing.T) {
	hw := fakeHardwareQuerier{hw: map[storage.RuleSource][]string{}}
	rules := fakeRuleSetReader{ruleSets: nil}

	m := NewMatcher(rules, hw)
	ids, err := m.MatchedGroupIds()
	if err != nil {
		t.Fatalf("MatchedGroupIds: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected 0 IDs for empty rules, got %v", ids)
	}
}
