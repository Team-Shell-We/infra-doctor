package erd

import "fmt"

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
