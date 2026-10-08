package erd

import (
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/ui"
)

const minBoxWidth = 28

// RenderASCII : 엔티티마다 컬럼 목록이 담긴 박스를 세로로 나열하고, 관계는 그 아래
// 텍스트 목록으로 보여준다. 엔티티가 많아지면(수십 개) 2차원 레이아웃은 터미널에서
// 어차피 읽기 어려워지므로, 개수가 늘어도 안정적으로 읽을 수 있는 목록형으로 둔다.
func RenderASCII(d Diagram) string {

	var b strings.Builder
	b.WriteString(d.Title + "\n\n")

	for _, entity := range d.Entities {
		b.WriteString(entityBox(entity))
		b.WriteString("\n")
	}

	if len(d.Relations) > 0 {
		b.WriteString("Relations\n")
		for _, relation := range d.Relations {
			b.WriteString("  " + relationLine(relation) + "\n")
		}
	}

	return b.String()
}

func entityBox(e Entity) string {

	width := minBoxWidth
	if w := ui.DisplayWidth(e.Name); w > width {
		width = w
	}
	for _, column := range e.Columns {
		if w := ui.DisplayWidth(columnLine(column)); w > width {
			width = w
		}
	}

	border := "+" + strings.Repeat("-", width+2) + "+\n"

	var b strings.Builder
	b.WriteString(border)
	b.WriteString("| " + ui.PadRight(e.Name, width) + " |\n")
	b.WriteString(border)
	for _, column := range e.Columns {
		b.WriteString("| " + ui.PadRight(columnLine(column), width) + " |\n")
	}
	b.WriteString(border)

	return b.String()
}

func columnLine(c Column) string {
	line := typeOrDefault(c.Type) + " " + c.Name
	if key := keyLabel(c); key != "" {
		line += " " + key
	}
	return line
}

func relationLine(r Relation) string {
	line := r.From + " --(" + r.Cardinality + ")--> " + r.To
	if r.Label != "" {
		line += " : " + r.Label
	}
	return line
}
