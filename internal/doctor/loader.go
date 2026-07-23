package doctor

import (
	"embed"
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed rules/*.yaml
var rulesFS embed.FS

type RuleDefinition struct {
	Category string `yaml:"category"`
	Level    string `yaml:"level"`
	Score    int    `yaml:"score"`

	Title   string `yaml:"title"`
	Message string `yaml:"message"`
	Reason  string `yaml:"reason"`
	Fix     string `yaml:"fix"`
}

type deploymentRuleFile struct {
	Deployment map[string]RuleDefinition `yaml:"deployment"`
}

type productionRuleFile struct {
	Production map[string]RuleDefinition `yaml:"production"`
}

type localDevRuleFile struct {
	LocalDev map[string]RuleDefinition `yaml:"localdev"`
}

type RuleRegistry struct {
	Deployment map[string]Diagnosis
	Production map[string]Diagnosis
	LocalDev   map[string]Diagnosis
}

var (
	registry *RuleRegistry
	once     sync.Once
)

// LoadRules loads every rule file only once.
func LoadRules() (*RuleRegistry, error) {

	var err error

	once.Do(func() {

		registry = &RuleRegistry{
			Deployment: make(map[string]Diagnosis),
			Production: make(map[string]Diagnosis),
			LocalDev:   make(map[string]Diagnosis),
		}

		if err = loadDeploymentRules(registry); err != nil {
			return
		}

		if err = loadProductionRules(registry); err != nil {
			return
		}

		err = loadLocalDevRules(registry)
	})

	return registry, err
}

func loadDeploymentRules(registry *RuleRegistry) error {

	data, err := rulesFS.ReadFile("rules/deployment.yaml")
	if err != nil {
		return err
	}

	var file deploymentRuleFile

	if err := yaml.Unmarshal(data, &file); err != nil {
		return err
	}

	for id, rule := range file.Deployment {
		registry.Deployment[id] = toDiagnosis(rule)
	}

	return nil
}

func loadProductionRules(registry *RuleRegistry) error {

	data, err := rulesFS.ReadFile("rules/production.yaml")
	if err != nil {
		return err
	}

	var file productionRuleFile

	if err := yaml.Unmarshal(data, &file); err != nil {
		return err
	}

	for id, rule := range file.Production {
		registry.Production[id] = toDiagnosis(rule)
	}

	return nil
}

func loadLocalDevRules(registry *RuleRegistry) error {

	data, err := rulesFS.ReadFile("rules/localdev.yaml")
	if err != nil {
		return err
	}

	var file localDevRuleFile

	if err := yaml.Unmarshal(data, &file); err != nil {
		return err
	}

	for id, rule := range file.LocalDev {
		registry.LocalDev[id] = toDiagnosis(rule)
	}

	return nil
}

func toDiagnosis(rule RuleDefinition) Diagnosis {

	return Diagnosis{
		Category: Category(rule.Category),
		Level:    Level(rule.Level),

		ScoreImpact: rule.Score,

		Title:   rule.Title,
		Message: rule.Message,
		Reason:  rule.Reason,
		Fix:     rule.Fix,
	}
}

func (r *RuleRegistry) DeploymentRule(id string) (Diagnosis, error) {

	rule, ok := r.Deployment[id]
	if !ok {
		return Diagnosis{}, fmt.Errorf("deployment rule '%s' not found", id)
	}

	return rule, nil
}

func (r *RuleRegistry) ProductionRule(id string) (Diagnosis, error) {

	rule, ok := r.Production[id]
	if !ok {
		return Diagnosis{}, fmt.Errorf("production rule '%s' not found", id)
	}

	return rule, nil
}

func (r *RuleRegistry) LocalDevRule(id string) (Diagnosis, error) {

	rule, ok := r.LocalDev[id]
	if !ok {
		return Diagnosis{}, fmt.Errorf("localdev rule '%s' not found", id)
	}

	return rule, nil
}
