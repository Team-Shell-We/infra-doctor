package analyzer

import (
	"os"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
	"gopkg.in/yaml.v3"
)

type workflowYAML struct {
	Name string `yaml:"name"`

	On yaml.Node `yaml:"on"`

	Jobs map[string]struct {
		RunsOn string `yaml:"runs-on"`

		Steps []struct {
			Name string            `yaml:"name"`
			Uses string            `yaml:"uses"`
			Run  string            `yaml:"run"`
			With map[string]string `yaml:"with"`
			Env  map[string]string `yaml:"env"`
		} `yaml:"steps"`
	} `yaml:"jobs"`
}

func ParseWorkflow(path string, fileName string) (*project.WorkflowInfo, error) {

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var wf workflowYAML

	if err := yaml.Unmarshal(content, &wf); err != nil {
		return nil, err
	}

	workflow := &project.WorkflowInfo{
		Name: wf.Name,
		File: fileName,
		Path: path,
	}

	switch wf.On.Kind {

	// on: push
	case yaml.ScalarNode:
		workflow.Triggers = append(workflow.Triggers, project.TriggerInfo{Event: wf.On.Value})

	// on: [push, pull_request]
	case yaml.SequenceNode:
		for _, event := range wf.On.Content {
			if event.Kind == yaml.ScalarNode {
				workflow.Triggers = append(workflow.Triggers, project.TriggerInfo{Event: event.Value})
			}
		}

	// on: { push: { branches: [...] }, pull_request: {} }
	case yaml.MappingNode:
		for i := 0; i < len(wf.On.Content); i += 2 {

			key := wf.On.Content[i]
			value := wf.On.Content[i+1]

			trigger := project.TriggerInfo{
				Event: key.Value,
			}

			if value.Kind == yaml.MappingNode {

				for j := 0; j < len(value.Content); j += 2 {

					if value.Content[j].Value != "branches" {
						continue
					}

					branches := value.Content[j+1]

					if branches.Kind == yaml.SequenceNode {

						for _, branch := range branches.Content {
							trigger.Branches = append(trigger.Branches, branch.Value)
						}
					}
				}
			}

			workflow.Triggers = append(workflow.Triggers, trigger)
		}
	}

	for jobName, job := range wf.Jobs {

		jobInfo := project.JobInfo{
			Name: jobName,
		}

		for _, step := range job.Steps {

			jobInfo.Steps = append(jobInfo.Steps,
				project.StepInfo{
					Name: step.Name,
					Uses: step.Uses,
					Run:  step.Run,
					With: step.With,
					Env:  step.Env,
				},
			)
		}

		workflow.Jobs = append(workflow.Jobs, jobInfo)
	}

	return workflow, nil
}
