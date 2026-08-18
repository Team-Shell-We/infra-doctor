package generate

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/doctor"
	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

// 전체 파이프라인 실행
type Analyzer func(string) (*project.Info, error)
type Diagnoser func(*project.Info) *doctor.Result

type Options struct {
	Root      string
	OutputDir string
	Overwrite bool
	DryRun    bool
	Config    Config
}

// Config contains the optional values used to customize generated files.
// It lives in the generate package so callers do not need an application layer.
type Config struct {
	ProjectName     string
	ApplicationPort int
	HealthPath      string
	ServiceName     string
	DockerImage     string
	Namespace       string
	Replicas        int
}

type Service struct {
	Analyze    Analyzer
	Diagnose   Diagnoser
	Generators Registry
	Writer     Writer
}

func (s Service) Generate(
	target Target,
	options Options,
) (Result, error) {
	info, err := s.Analyze(options.Root)
	if err != nil {
		return Result{}, err
	}

	var diagnosis *doctor.Result
	if s.Diagnose != nil {
		diagnosis = s.Diagnose(info)
	}

	ctx, warnings := BuildContext(
		*info,
		diagnosis,
		options.Config,
	)

	generator, found := s.Generators[target]
	if !found {
		return Result{}, fmt.Errorf(
			"generator for %q is not registered",
			target,
		)
	}

	files, err := generator.Plan(ctx)
	if err != nil {
		return Result{}, err
	}

	result, err := s.Writer.Write(
		options.OutputDir,
		files,
		WriteOptions{
			Overwrite: options.Overwrite,
			DryRun:    options.DryRun,
		},
	)

	result.Target = target
	result.Warnings = warnings
	result.Detections = buildDetections(*info)

	return result, err
}
