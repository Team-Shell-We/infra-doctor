package export

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/analyzer"
	"github.com/Team-Shell-We/infra-doctor/internal/doctor"
	"github.com/Team-Shell-We/infra-doctor/internal/generate"
	"github.com/Team-Shell-We/infra-doctor/internal/generate/ci"
	"github.com/Team-Shell-We/infra-doctor/internal/generate/compose"
	"github.com/Team-Shell-We/infra-doctor/internal/generate/docker"
	"github.com/Team-Shell-We/infra-doctor/internal/generate/k8s"
	"github.com/Team-Shell-We/infra-doctor/internal/generate/nginx"
	"github.com/Team-Shell-We/infra-doctor/internal/project"
	"github.com/Team-Shell-We/infra-doctor/internal/visualize"
)

const OutputDirectory = "infra-doctor"

type Request struct {
	Root   string
	Force  bool
	DryRun bool
	Lang   string
}

type Application struct {
	Analyze  func(string) (*project.Info, error)
	Diagnose func(*project.Info) *doctor.Result
	Writer   generate.Writer
}

func NewApplication() *Application {
	return &Application{Analyze: analyzer.AnalyzeProject, Diagnose: doctor.Analyze, Writer: generate.Writer{}}
}

func (a *Application) Run(_ context.Context, request Request, output io.Writer) error {
	root := request.Root
	if root == "" {
		root = "."
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return err
	}

	fmt.Fprintln(output, "Analyzing project...")
	info, err := a.Analyze(root)
	if err != nil {
		return err
	}
	diagnosis := a.Diagnose(info)
	files, warnings, err := buildFiles(root, *info, diagnosis, request.Lang)
	if err != nil {
		return err
	}

	result, err := a.Writer.Write(filepath.Join(root, OutputDirectory), files, generate.WriteOptions{Overwrite: request.Force, DryRun: request.DryRun})
	if err != nil {
		return err
	}
	result.Warnings = warnings
	return printResult(output, result)
}

func buildFiles(root string, info project.Info, diagnosis *doctor.Result, lang string) ([]generate.File, []string, error) {
	architecture := visualize.Build(info)
	flow, err := visualize.BuildDeploymentFlow(root, info)
	if err != nil {
		return nil, nil, err
	}
	architectureMarkdown, err := visualize.Render(architecture, visualize.Markdown)
	if err != nil {
		return nil, nil, err
	}
	architectureMermaid, err := visualize.Render(architecture, visualize.Mermaid)
	if err != nil {
		return nil, nil, err
	}
	flowMarkdown, err := visualize.Render(flow, visualize.Markdown)
	if err != nil {
		return nil, nil, err
	}

	files := []generate.File{
		textFile("report.md", renderReport(info, diagnosis)),
		textFile("architecture.md", architectureMarkdown),
		textFile("architecture.mmd", architectureMermaid),
		textFile("deployment-flow.md", flowMarkdown),
		textFile("recommendations.md", renderRecommendations(diagnosis)),
	}
	ctx, warnings := generate.BuildContext(info, diagnosis, generate.Config{}, lang)
	generated, err := collectGeneratedFiles(ctx)
	if err != nil {
		return nil, nil, err
	}
	return append(files, generated...), warnings, nil
}

type generatorCall struct {
	generator generate.Generator
	directory string
	paths     map[string]string
}

func collectGeneratedFiles(ctx generate.Context) ([]generate.File, error) {
	calls := []generatorCall{
		{generator: docker.Generator{}, directory: "docker"},
		{generator: compose.Generator{}, directory: "docker"},
		{generator: nginx.Generator{}, directory: "docker"},
		{generator: k8s.Generator{}, directory: "kubernetes", paths: map[string]string{
			"k8s/deployment.yml": "deployment.yaml", "k8s/service.yml": "service.yaml", "k8s/configmap.yml": "configmap.yaml",
		}},
		{generator: ci.Generator{}, directory: "github", paths: map[string]string{".github/workflows/ci.yml": "ci.yml"}},
	}
	var files []generate.File
	for _, call := range calls {
		generated, err := call.generator.Plan(ctx)
		if err != nil {
			return nil, fmt.Errorf("export %s: %w", call.generator.Target(), err)
		}
		for _, file := range generated {
			name := filepath.Base(file.Path)
			if mapped, ok := call.paths[filepath.ToSlash(file.Path)]; ok {
				name = mapped
			}
			file.Path = filepath.Join(call.directory, name)
			files = append(files, file)
		}
	}
	return files, nil
}

func textFile(path, content string) generate.File {
	return generate.File{Path: path, Content: []byte(content), Mode: 0o644}
}

func renderReport(info project.Info, diagnosis *doctor.Result) string {
	var b strings.Builder
	b.WriteString("# Infrastructure Analysis Report\n\n")

	// scan command result
	b.WriteString("## Scan Result\n\n### Framework\n\n")
	fmt.Fprintf(&b, "- Spring Boot: %s\n- Java: %s\n- Build tool: %s\n", value(info.Framework.SpringBoot.Version), value(info.Framework.Java.Version), value(info.Framework.BuildTool.Type))

	b.WriteString("\n### Dependencies\n\n")
	writeDetected(&b, "Spring Security", info.Dependencies.Security.Enabled)
	writeDetected(&b, "Spring Data JPA", info.Dependencies.JPA.Enabled)
	writeDetected(&b, "Kafka", info.Dependencies.Kafka.Enabled)
	writeDetected(&b, "AWS SDK", info.Dependencies.AWS.Enabled)
	writeDetected(&b, "Lombok", info.Dependencies.Lombok.Enabled)
	writeDetected(&b, "Spring Boot Actuator", info.Dependencies.Actuator.Enabled)
	writeDetected(&b, "SpringDoc OpenAPI", info.Dependencies.OpenAPI.Enabled)

	b.WriteString("\n### Database\n\n")
	fmt.Fprintf(&b, "- Primary: %s\n", value(info.Database.Primary.Type))
	writeDetected(&b, "Redis", info.Database.Redis != nil && info.Database.Redis.Enabled)

	b.WriteString("\n### Infrastructure\n\n")
	writeDetected(&b, "Docker", info.Infrastructure.Docker.Enabled)
	for _, file := range info.Infrastructure.Docker.Dockerfiles {
		fmt.Fprintf(&b, "  - `%s`\n", file.File)
	}
	writeDetected(&b, "Docker Compose", len(info.Infrastructure.Docker.Compose) > 0)
	for _, file := range info.Infrastructure.Docker.Compose {
		fmt.Fprintf(&b, "  - `%s`\n", file.File)
	}
	writeDetected(&b, "Kubernetes", info.Infrastructure.Kubernetes.Enabled)
	for _, file := range info.Infrastructure.Kubernetes.Files {
		fmt.Fprintf(&b, "  - `%s`\n", file.File)
	}
	writeDetected(&b, "Nginx", info.Infrastructure.Nginx.Enabled)
	writeDetected(&b, "Terraform", info.Infrastructure.Terraform.Enabled)

	b.WriteString("\n### CI/CD\n\n")
	writeDetected(&b, "GitHub Actions", len(info.Github.Workflows) > 0)
	for _, workflow := range info.Github.Workflows {
		fmt.Fprintf(&b, "- Workflow: %s (`%s`)\n", value(workflow.Name), workflow.File)
		for _, trigger := range workflow.Triggers {
			fmt.Fprintf(&b, "  - Trigger: %s", trigger.Event)
			if len(trigger.Branches) > 0 {
				fmt.Fprintf(&b, " (%s)", strings.Join(trigger.Branches, ", "))
			}
			b.WriteByte('\n')
		}
		for _, job := range workflow.Jobs {
			fmt.Fprintf(&b, "  - Job: %s\n", job.Name)
		}
	}

	b.WriteString("\n### Profiles\n\n")
	if len(info.Profiles) == 0 {
		b.WriteString("- Not detected\n")
	}
	for _, profile := range info.Profiles {
		fmt.Fprintf(&b, "- %s (`%s`)\n", profile.Name, profile.File)
	}

	// doctor command result
	fmt.Fprintf(&b, "\n## Doctor Result\n\n### Readiness\n\n- Score: %d%%\n\n### Infrastructure Checks\n\n", diagnosis.Score)
	for _, check := range diagnosis.Checks {
		mark := "❌"
		if check.Passed {
			mark = "✅"
		}
		fmt.Fprintf(&b, "- %s %s\n", mark, check.Name)
	}
	return b.String()
}

func writeDetected(b *strings.Builder, name string, detected bool) {
	mark := "❌"
	if detected {
		mark = "✅"
	}
	fmt.Fprintf(b, "- %s %s\n", mark, name)
}

func renderRecommendations(diagnosis *doctor.Result) string {
	items := make([]string, 0, len(diagnosis.Diagnoses))
	for _, item := range diagnosis.Diagnoses {
		text := strings.TrimSpace(item.Fix)
		if text == "" {
			text = strings.TrimSpace(item.Message)
		}
		if text != "" {
			items = append(items, text)
		}
	}
	if len(items) == 0 {
		items = append(items, "No immediate infrastructure recommendations were detected.")
	}
	sort.Strings(items)
	return "# Recommendations\n\n- " + strings.Join(items, "\n- ") + "\n"
}

func value(value string) string {
	if strings.TrimSpace(value) == "" {
		return "Not detected"
	}
	return value
}

func printResult(output io.Writer, result generate.Result) error {
	for _, warning := range result.Warnings {
		if _, err := fmt.Fprintf(output, "warning: %s\n", warning); err != nil {
			return err
		}
	}
	for _, path := range result.Planned {
		status := "created"
		if result.DryRun {
			status = "planned"
		}
		if contains(result.Skipped, path) {
			status = "skipped"
		} else if contains(result.Overwritten, path) {
			status = "overwritten"
		}
		if _, err := fmt.Fprintf(output, "%s: %s\n", status, filepath.Join(OutputDirectory, path)); err != nil {
			return err
		}
	}
	return nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
