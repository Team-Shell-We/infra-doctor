package analyzer

import (
	"os"
	"path/filepath"
	"regexp"
)

var (
	includeLineRegex  = regexp.MustCompile(`(?m)^\s*include\b.*$`)
	quotedModuleRegex = regexp.MustCompile(`['"]([^'"]+)['"]`)
)

// AnalyzeGradleModules : root의 settings.gradle(.kts)에서 include된
// 서브모듈 개수를 센다. 파일이 없거나 include가 없으면 0.
func AnalyzeGradleModules(root string) (int, error) {

	for _, name := range []string{"settings.gradle", "settings.gradle.kts"} {

		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return 0, err
		}

		count := 0
		for _, line := range includeLineRegex.FindAllString(string(data), -1) {
			count += len(quotedModuleRegex.FindAllString(line, -1))
		}

		return count, nil
	}

	return 0, nil
}
