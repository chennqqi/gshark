package models

import (
	"regexp"

	"github.com/chennqqi/gshark/vars"
)

type FilterRule struct {
	Id        int64
	RuleType  int // 0: blacklist rule, 1: whitelist rule
	RuleKey   string
	RuleValue string `xorm:"text"`
}

type ExcludeFilter struct {
	FilterRule
	exp *regexp.Regexp
}

func (r *ExcludeFilter) compile() error {
	var err error
	r.exp, err = regexp.Compile(r.RuleValue)
	return err
}

func (r *ExcludeFilter) Exclude(name string) bool {
	return r.exp.MatchString(name)
}

func NewFilterRule(ruleType int, ruleKey, ruleValue string) *FilterRule {
	return &FilterRule{RuleType: ruleType, RuleKey: ruleKey, RuleValue: ruleValue}
}

func (r *FilterRule) Insert() (err error) {
	_, err = Engine.Insert(r)
	return err
}

func GetFilterRules() ([]ExcludeFilter, error) {
	var rules []ExcludeFilter
	err := Engine.Table("filter_rule").Where(`rule_key!=?`, "name").Find(&rules)
	for i := 0; i < len(rules); i++ {
		rule = &rules[i]
		err = rule.compile()
		if err != nil {
			return nil, err
		}
	}
	return rules, err
}

func GetExcludeNames() ([]FilterRule, error) {
	rules := make([]FilterRule, 0)
	err := Engine.Table("filter_rule").Where(`rule_key=?`, "name").Find(&rules)
	return rules, err
}

func GetFilterRulesPage(page int) ([]FilterRule, int, error) {
	rules := make([]FilterRule, 0)
	totalPages, err := Engine.Table("filter_rule").Count()

	var pages int
	pageSize := int64(vars.PAGE_SIZE)

	if totalPages%pageSize == 0 {
		pages = int(totalPages / pageSize)
	} else {
		pages = int(totalPages/pageSize) + 1
	}

	if page >= pages {
		page = pages
	}

	if page < 1 {
		page = 1
	}

	err = Engine.Table("filter_rule").Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).
		Desc("rule_type").Find(&rules)
	return rules, pages, err
}

func GetFilterRuleById(id int64) (*FilterRule, bool, error) {
	rule := new(FilterRule)
	has, err := Engine.ID(id).Get(rule)
	return rule, has, err
}

func EditFilterRuleById(id int64, ruleType int, ruleKey, ruleValue string) error {
	rule := new(FilterRule)
	_, has, err := GetFilterRuleById(id)
	if err == nil && has {
		rule.RuleKey = ruleKey
		rule.RuleType = ruleType
		rule.RuleValue = ruleValue
		_, err = Engine.ID(id).Update(rule)
	}
	return err
}

func DeleteFilterRuleById(id int64) (err error) {
	rule := new(FilterRule)
	_, err = Engine.Id(id).Delete(rule)
	return err
}

func ConvertRuleType(id int64) (err error) {
	rule := new(FilterRule)
	has, err := Engine.Id(id).Get(rule)
	if err == nil && has {
		rule.RuleType = rule.RuleType ^ 1
		_, err = Engine.Id(id).Cols("rule_type").Update(rule)
	}
	return err
}
