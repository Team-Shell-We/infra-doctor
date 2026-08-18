package ci

import (
	"embed"
	"errors"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
)

//go:embed templates/*
var templates embed.FS

type Generator struct{}

func (Generator) Target() generate.Target {
	return generate.TargetCI
}

func (Generator) Plan(
	ctx generate.Context,
) ([]generate.File, error) {
	if ctx.BuildCommand == "" {
		return nil, errors.New(
			"cannot generate CI without a build command",
		)
	}

	content, err := generate.RenderTemplate(
		templates,
		"templates/github-actions.yml.tmpl",
		ctx,
	)
	if err != nil {
		return nil, err
	}

	return []generate.File{
		{
			Path:    ".github/workflows/ci.yml",
			Content: content,
		},
	}, nil
}
