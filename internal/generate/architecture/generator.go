package architecture

import "github.com/Team-Shell-We/infra-doctor/internal/generate"

type Generator struct {
	Components []generate.Generator
}

func (Generator) Target() generate.Target {
	return generate.TargetArchitecture
}

func (g Generator) Plan(
	ctx generate.Context,
) ([]generate.File, error) {
	var files []generate.File

	for _, component := range g.Components {
		if !shouldGenerate(component.Target(), ctx) {
			continue
		}

		part, err := component.Plan(ctx)
		if err != nil {
			return nil, err
		}

		files = append(files, part...)
	}

	return files, nil
}

func shouldGenerate(
	target generate.Target,
	ctx generate.Context,
) bool {
	switch target {
	case generate.TargetCompose:
		return ctx.NeedsCompose && (len(ctx.Databases) > 0 || ctx.Redis)

	case generate.TargetDocker:
		return ctx.NeedsDocker && ctx.Runtime != ""

	case generate.TargetNginx:
		return ctx.NeedsNginx

	case generate.TargetCI:
		return ctx.NeedsCI

	default:
		return true
	}
}
