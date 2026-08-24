package docker

import (
	"embed"
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
)

//Dockerfile, .dockerignore 생성

//go:embed templates/*
var templates embed.FS

type Generator struct{}

func (Generator) Target() generate.Target {
	return generate.TargetDocker
}

func (Generator) Plan(
	ctx generate.Context,
) ([]generate.File, error) {
	if ctx.Runtime != "java" {
		return nil, fmt.Errorf(
			"docker generation does not support runtime %q",
			ctx.Runtime,
		)
	}

	ctx.Header = generate.BuildHeader(ctx.Lang, "generate.docker.nextSteps")

	dockerfile, err := generate.RenderTemplate(
		templates,
		"templates/Dockerfile.tmpl",
		ctx,
	)
	if err != nil {
		return nil, err
	}

	dockerignore, err := generate.RenderTemplate(
		templates,
		"templates/dockerignore.tmpl",
		ctx,
	)
	if err != nil {
		return nil, err
	}

	return []generate.File{
		{
			Path:    "Dockerfile",
			Content: dockerfile,
		},
		{
			Path:    ".dockerignore",
			Content: dockerignore,
		},
	}, nil
}
