package visualize

import "fmt"

type Format string

const (
	ASCII    Format = "ascii"
	Mermaid  Format = "mermaid"
	Markdown Format = "markdown"
)

func Render(d Diagram, f Format) (string, error) {
	switch f {
	case ASCII:
		return RenderASCII(d), nil
	case Mermaid:
		return RenderMermaid(d), nil
	case Markdown:
		return RenderMarkdown(d), nil
	default:
		return "", fmt.Errorf("unsupported output format %q", f)
	}
}
