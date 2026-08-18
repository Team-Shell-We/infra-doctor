package generate

import "io/fs"

type File struct {
	Path    string
	Content []byte
	Mode    fs.FileMode
}

type Generator interface {
	Target() Target
	Plan(Context) ([]File, error)
}

type Registry map[Target]Generator

func NewRegistry(
	generators ...Generator,
) Registry {
	registry := Registry{}

	for _, generator := range generators {
		registry[generator.Target()] = generator
	}

	return registry
}
