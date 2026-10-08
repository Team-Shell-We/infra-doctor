package export

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

// Category : export가 만들 수 있는 산출물 묶음 하나. ID는 i18n 키("export.select.category."+ID)와
// 테스트에서 그대로 쓰는 안정적인 식별자다.
type Category struct {
	ID    string
	Files []generate.File
}

// categoryOrder : 체크리스트에 보여주는 고정 순서. buildFiles()가 실제로 만드는
// 산출물과 1:1로 대응한다 — export는 항상 새로 계산해서 쓰기 때문에, 디스크를
// 스캔하는 대신 이 고정 목록으로 "만들 수 있는 것"을 보여준다.
var categoryOrder = []string{
	"report", "architecture", "erd", "flow", "recommendations", "docker", "kubernetes", "github",
}

// Categorize : buildFiles()가 만든 평평한 파일 목록을 8개 고정 카테고리로 묶는다.
// 매칭되는 파일이 하나도 없는 카테고리(예: ERD 빌드가 실패해 erd 파일이 없는 경우)는
// 목록에서 빠진다.
func Categorize(files []generate.File) []Category {

	buckets := make(map[string][]generate.File, len(categoryOrder))
	for _, file := range files {
		id := categoryID(file.Path)
		if id == "" {
			continue
		}
		buckets[id] = append(buckets[id], file)
	}

	var categories []Category
	for _, id := range categoryOrder {
		if matched, ok := buckets[id]; ok {
			categories = append(categories, Category{ID: id, Files: matched})
		}
	}

	return categories
}

func categoryID(path string) string {

	slash := filepath.ToSlash(path)

	switch slash {
	case "report.md":
		return "report"
	case "architecture.md", "architecture.mmd":
		return "architecture"
	case "erd.md", "erd.mmd":
		return "erd"
	case "deployment-flow.md":
		return "flow"
	case "recommendations.md":
		return "recommendations"
	}

	top := strings.SplitN(slash, "/", 2)[0]
	switch top {
	case "docker", "kubernetes", "github":
		return top
	}

	return ""
}

// PromptSelection : 카테고리를 번호 목록으로 보여주고, 쉼표로 구분된 선택(또는 "all")을
// 한 번 읽어 해당 카테고리의 파일만 모아 반환한다. cmd/login.go의 번호 선택 프롬프트와
// 같은 스타일 — 재시도 루프 없이 1회만 읽고, 잘못된 입력은 에러로 끝낸다.
func PromptSelection(in io.Reader, out io.Writer, categories []Category, lang string) ([]generate.File, error) {

	reader := bufio.NewReader(in)

	fmt.Fprintln(out, i18n.Get(lang, "export.select.title"))
	fmt.Fprintln(out)
	for i, category := range categories {
		fmt.Fprintf(out, "  %d. %s\n", i+1, i18n.Get(lang, "export.select.category."+category.ID))
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, i18n.Get(lang, "export.select.instructions"))
	fmt.Fprint(out, i18n.Get(lang, "export.select.prompt"))

	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	if strings.EqualFold(line, "all") {
		var files []generate.File
		for _, category := range categories {
			files = append(files, category.Files...)
		}
		return files, nil
	}

	seen := make(map[int]bool)
	var selected []generate.File

	for _, token := range strings.Split(line, ",") {

		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		n, err := strconv.Atoi(token)
		if err != nil || n < 1 || n > len(categories) {
			return nil, fmt.Errorf("%s", i18n.Get(lang, "export.select.invalidChoice"))
		}

		if seen[n] {
			continue
		}
		seen[n] = true

		selected = append(selected, categories[n-1].Files...)
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("%s", i18n.Get(lang, "export.select.empty"))
	}

	return selected, nil
}
