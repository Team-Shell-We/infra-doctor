package nginx

import (
	"embed"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
)

//go:embed templates/*
var templates embed.FS

type Generator struct{}

func (Generator) Target() generate.Target {
	return generate.TargetNginx
}

func (Generator) Plan(
	ctx generate.Context,
) ([]generate.File, error) {
	ctx.Header = generate.BuildHeader(ctx.Lang, "generate.nginx.nextSteps")

	content, err := generate.RenderTemplate(
		templates,
		"templates/nginx.conf.tmpl",
		ctx,
	)
	if err != nil {
		return nil, err
	}

	return []generate.File{
		{
			Path:    "nginx.conf",
			Content: content,
		},
	}, nil
}
