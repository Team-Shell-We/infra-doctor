package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/doctor"
	"github.com/Team-Shell-We/infra-doctor/internal/generate"
	"github.com/Team-Shell-We/infra-doctor/internal/generate/architecture"
	"github.com/Team-Shell-We/infra-doctor/internal/generate/ci"
	"github.com/Team-Shell-We/infra-doctor/internal/generate/compose"
	"github.com/Team-Shell-We/infra-doctor/internal/generate/docker"
	"github.com/Team-Shell-We/infra-doctor/internal/generate/k8s"
	"github.com/Team-Shell-We/infra-doctor/internal/generate/nginx"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type GenerateRunner interface {
	Generate(generate.Target, generate.Options) (generate.Result, error)
}

type generateRequest struct {
	Target, Root, OutputDir, ConfigPath string
	Force, DryRun                       bool
}

func generateCommand(runner GenerateRunner) *cobra.Command {
	request := generateRequest{}
	command := &cobra.Command{
		Use: "generate <target> [path]", Short: "Generate infrastructure configuration files",
		Args: cobra.RangeArgs(1, 2), SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			request.Target, request.Root = args[0], "."
			if len(args) == 2 {
				request.Root = args[1]
			}
			target, err := generate.ParseTarget(request.Target)
			if err != nil {
				return err
			}
			loaded, err := loadGenerateConfig(request.Root, request.ConfigPath)
			if err != nil {
				return err
			}
			outputDir := request.OutputDir
			if outputDir == "" {
				outputDir = loaded.OutputDir
			}
			if outputDir == "" {
				outputDir = request.Root
			}
			if !filepath.IsAbs(outputDir) && outputDir != request.Root {
				outputDir = filepath.Join(request.Root, outputDir)
			}
			overwrite := loaded.Overwrite || request.Force
			result, err := runner.Generate(target, generate.Options{Root: request.Root, OutputDir: outputDir, Overwrite: overwrite, DryRun: request.DryRun, Config: loaded.Generate})
			if err != nil {
				return err
			}
			return printGenerateResult(cmd, result)
		},
	}
	command.Flags().BoolVarP(&request.Force, "force", "f", false, "overwrite existing files")
	command.Flags().BoolVar(&request.DryRun, "dry-run", false, "show planned files without writing")
	command.Flags().StringVarP(&request.OutputDir, "output-dir", "o", "", "directory for generated files")
	command.Flags().StringVar(&request.ConfigPath, "config", "", "infra-doctor configuration file")
	return command
}

func newGenerateService() generate.Service {
	components := []generate.Generator{docker.Generator{}, compose.Generator{}, nginx.Generator{}, ci.Generator{}, k8s.Generator{}}
	all := append([]generate.Generator{}, components...)
	all = append(all, architecture.Generator{Components: components})
	return generate.Service{Analyze: analyzer.AnalyzeProject, Diagnose: doctor.Analyze, Generators: generate.NewRegistry(all...), Writer: generate.Writer{}}
}

type generateConfigFile struct {
	Project struct {
		Name string `yaml:"name"`
	} `yaml:"project"`
	Generate struct {
		ApplicationPort int    `yaml:"applicationPort"`
		HealthPath      string `yaml:"healthPath"`
		ServiceName     string `yaml:"serviceName"`
		DockerImage     string `yaml:"dockerImage"`
		Namespace       string `yaml:"namespace"`
		Replicas        int    `yaml:"replicas"`
	} `yaml:"generate"`
	Output struct {
		Directory string `yaml:"directory"`
		Overwrite bool   `yaml:"overwrite"`
	} `yaml:"output"`
}

type loadedGenerateConfig struct {
	Generate  generate.Config
	OutputDir string
	Overwrite bool
}

func loadGenerateConfig(root, explicitPath string) (loadedGenerateConfig, error) {
	path := explicitPath
	if path == "" {
		path = filepath.Join(root, ".infra-doctor", "config.yaml")
	}
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) && explicitPath == "" {
		return loadedGenerateConfig{}, nil
	}
	if err != nil {
		return loadedGenerateConfig{}, fmt.Errorf("read config: %w", err)
	}
	var file generateConfigFile
	if err := yaml.Unmarshal(content, &file); err != nil {
		return loadedGenerateConfig{}, fmt.Errorf("parse config: %w", err)
	}
	return loadedGenerateConfig{
		Generate:  generate.Config{ProjectName: file.Project.Name, ApplicationPort: file.Generate.ApplicationPort, HealthPath: file.Generate.HealthPath, ServiceName: file.Generate.ServiceName, DockerImage: file.Generate.DockerImage, Namespace: file.Generate.Namespace, Replicas: file.Generate.Replicas},
		OutputDir: file.Output.Directory,
		Overwrite: file.Output.Overwrite,
	}, nil
}

func printGenerateResult(cmd *cobra.Command, result generate.Result) error {
	for _, warning := range result.Warnings {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "warning: %s\n", warning); err != nil {
			return err
		}
	}
	for _, path := range result.Planned {
		status := "created"
		if result.DryRun {
			status = "planned"
		}
		for _, skipped := range result.Skipped {
			if path == skipped {
				status = "skipped"
				break
			}
		}
		for _, overwritten := range result.Overwritten {
			if path == overwritten {
				status = "overwritten"
				break
			}
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", status, path); err != nil {
			return err
		}
	}
	return nil
}

func init() { rootCmd.AddCommand(generateCommand(newGenerateService())) }
