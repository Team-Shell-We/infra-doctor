package doctor

import (
	"embed"
	"fmt"
	"strings"
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
	loadErr  error
	once     sync.Once
)

// LoadRules : 모든 룰 파일을 딱 한 번만 로드. once.Do는 첫 호출에서만
// 실행되므로, 로드 실패 시의 err도 loadErr에 저장해 이후 호출에서도
// 계속 반환해야 한다 — 로컬 var err에만 담으면 두 번째 호출부터는
// nil로 돌아가 실패를 숨긴다.
func LoadRules() (*RuleRegistry, error) {

	once.Do(func() {

		registry = &RuleRegistry{
			Deployment: make(map[string]Diagnosis),
			Production: make(map[string]Diagnosis),
			LocalDev:   make(map[string]Diagnosis),
		}

		if loadErr = loadDeploymentRules(registry); loadErr != nil {
			return
		}

		if loadErr = loadProductionRules(registry); loadErr != nil {
			return
		}

		loadErr = loadLocalDevRules(registry)
	})

	return registry, loadErr
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
		Level:    Level(strings.ToUpper(rule.Level)),

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
