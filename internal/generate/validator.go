package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// 생성 파일의 경로와 내용 검증
func ValidateFiles(files []File) error {
	seen := map[string]bool{}

	for _, file := range files {
		if file.Path == "" {
			return fmt.Errorf("generated path is empty")
		}

		if len(file.Content) == 0 {
			return fmt.Errorf(
				"generated file %q is empty",
				file.Path,
			)
		}

		if filepath.IsAbs(file.Path) {
			return fmt.Errorf(
				"generated path must be relative: %s",
				file.Path,
			)
		}

		clean := filepath.Clean(file.Path)

		if clean == ".." ||
			strings.HasPrefix(
				clean,
				".."+string(filepath.Separator),
			) {
			return fmt.Errorf(
				"generated path escapes output directory: %s",
				file.Path,
			)
		}

		if seen[clean] {
			return fmt.Errorf(
				"duplicate generated path %q",
				file.Path,
			)
		}

		seen[clean] = true
	}

	return nil
}
