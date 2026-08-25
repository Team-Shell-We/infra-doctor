package k8s

import (
	"embed"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
)

//go:embed templates/*
var templates embed.FS

type Generator struct{}

func (Generator) Target() generate.Target {
	return generate.TargetK8s
}

func (Generator) Plan(
	ctx generate.Context,
) ([]generate.File, error) {
	ctx.Header = generate.BuildHeader(ctx.Lang, "generate.k8s.nextSteps")
	ctx.ImageNote = generate.NoteBlock(ctx.Lang, "generate.k8s.imageNote")
	ctx.ResourcesNote = generate.NoteBlock(ctx.Lang, "generate.k8s.resourcesNote")

	items := []struct {
		template string
		output   string
	}{
		{
			"templates/deployment.yml.tmpl",
			"k8s/deployment.yml",
		},
		{
			"templates/service.yml.tmpl",
			"k8s/service.yml",
		},
		{
			"templates/configmap.yml.tmpl",
			"k8s/configmap.yml",
		},
	}

	files := make([]generate.File, 0, len(items))

	for _, item := range items {
		content, err := generate.RenderTemplate(
			templates,
			item.template,
			ctx,
		)
		if err != nil {
			return nil, err
		}

		files = append(files, generate.File{
			Path:    item.output,
			Content: content,
		})
	}

	return files, nil
}
