package generate

import (
	"fmt"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

// commentBlock: 각 줄 앞에 "# "를 붙인다. generate가 만드는 Dockerfile/YAML/nginx.conf가 전부 "#"를 줄 주석으로 쓴다
func commentBlock(text string) string {

	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")

	for i, line := range lines {
		if line == "" {
			lines[i] = "#"
		} else {
			lines[i] = "# " + line
		}
	}

	return strings.Join(lines, "\n")
}

// BuildHeader: 생성 파일 맨 위 배너를 target별 "다음 할 일" 목록으로 lang(en/ko)에 맞춰 조립한다
func BuildHeader(lang, nextStepsKey string) string {

	var b strings.Builder

	b.WriteString(commentBlock(i18n.Get(lang, "generate.header")))
	b.WriteString("\n#\n")
	b.WriteString(commentBlock(i18n.Get(lang, "generate.nextStepsLabel")))
	b.WriteString("\n")
	b.WriteString(commentBlock(i18n.Get(lang, nextStepsKey)))
	b.WriteString("\n")

	return b.String()
}

// NoteBlock: i18n key 문구를 찾아 필요하면 fmt.Sprintf로 값을 채운 뒤 "# " 주석으로 감싼다
func NoteBlock(lang, key string, args ...any) string {

	text := i18n.Get(lang, key)

	if len(args) > 0 {
		text = fmt.Sprintf(text, args...)
	}

	return commentBlock(text)
}

// IndentLines: 각 줄 앞에 공백 n칸을 붙인다. 템플릿 안 들여쓴 위치의 노트를 주변 코드와 맞추는 데 쓴다
func IndentLines(text string, n int) string {

	prefix := strings.Repeat(" ", n)
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		lines[i] = prefix + line
	}

	return strings.Join(lines, "\n")
}
