package compose

import (
	"embed"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
)

//go:embed templates/*
var templates embed.FS

type Generator struct{}

func (Generator) Target() generate.Target {
	return generate.TargetCompose
}

func (Generator) Plan(
	ctx generate.Context,
) ([]generate.File, error) {
	content, err := generate.RenderTemplate(
		templates,
		"templates/docker-compose.yml.tmpl",
		ctx,
	)
	if err != nil {
		return nil, err
	}

	return []generate.File{
		{
			Path:    "docker-compose.yml",
			Content: content,
		},
	}, nil
}
