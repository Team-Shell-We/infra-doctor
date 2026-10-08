package erd

import (
	"strings"
)

var mermaidEscaper = strings.NewReplacer(`"`, `#quot;`, "\n", " ")

// RenderMermaid : Mermaid의 erDiagram 문법으로 렌더링한다. internal/visualize/mermaid.go의
// flowchart TD 문법과는 완전히 다른 블록 구조라 그쪽 렌더러를 재사용할 수 없다.
func RenderMermaid(d Diagram) string {

	var b strings.Builder
	b.WriteString("erDiagram\n")

	for _, entity := range d.Entities {
		if len(entity.Columns) == 0 {
			b.WriteString("    " + entity.Name + "\n")
			continue
		}

		b.WriteString("    " + entity.Name + " {\n")
		for _, column := range entity.Columns {
			b.WriteString("        " + mermaidAttribute(column) + "\n")
		}
		b.WriteString("    }\n")
	}

	for _, relation := range d.Relations {
		b.WriteString("    " + mermaidRelation(relation) + "\n")
	}

	return b.String()
}

func mermaidAttribute(c Column) string {
	line := typeOrDefault(c.Type) + " " + c.Name
	if key := keyLabel(c); key != "" {
		line += " " + key
	}
	return line
}

func mermaidRelation(r Relation) string {
	symbol := crowsFoot(r.Cardinality)
	label := mermaidEscaper.Replace(r.Label)
	return r.From + " " + symbol + " " + r.To + " : \"" + label + "\""
}

// crowsFoot : 카디널리티를 Mermaid의 crow's-foot 표기(예: "||--o{")로 변환한다.
// From쪽 기호가 From의 많음/적음을, To쪽 기호가 To의 많음/적음을 나타낸다.
func crowsFoot(cardinality string) string {
	switch cardinality {
	case "one-to-one":
		return "||--||"
	case "one-to-many":
		return "||--o{"
	case "many-to-many":
		return "}o--o{"
	default: // many-to-one
		return "}o--||"
	}
}

func typeOrDefault(t string) string {
	if strings.TrimSpace(t) == "" {
		return "string"
	}
	return t
}
