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
	ctx.Header = generate.BuildHeader(ctx.Lang, "generate.compose.nextSteps")

	if len(ctx.Databases) > 0 {
		ctx.CredentialsNote = generate.IndentLines(
			generate.NoteBlock(ctx.Lang, "generate.compose.credentialsNote"),
			2,
		)
	}

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
