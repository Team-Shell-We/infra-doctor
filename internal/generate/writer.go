package generate

import (
	"fmt"
	"os"
	"path/filepath"
)

// 실제 파일 저장
type WriteOptions struct {
	Overwrite bool
	DryRun    bool
}

type Writer struct{}

func (Writer) Write(
	root string,
	files []File,
	options WriteOptions,
) (Result, error) {
	result := Result{
		DryRun: options.DryRun,
	}

	if err := ValidateFiles(files); err != nil {
		return result, err
	}

	root, err := filepath.Abs(root)
	if err != nil {
		return result, err
	}

	for _, file := range files {
		result.Planned = append(
			result.Planned,
			file.Path,
		)

		destination := filepath.Join(
			root,
			filepath.Clean(file.Path),
		)

		_, statErr := os.Stat(destination)
		exists := statErr == nil

		if exists && !options.Overwrite {
			result.Skipped = append(
				result.Skipped,
				file.Path,
			)
			continue
		}

		if statErr != nil && !os.IsNotExist(statErr) {
			return result, statErr
		}

		if options.DryRun {
			continue
		}

		if err := os.MkdirAll(
			filepath.Dir(destination),
			0o755,
		); err != nil {
			return result, err
		}

		mode := file.Mode
		if mode == 0 {
			mode = 0o644
		}

		if err := os.WriteFile(
			destination,
			file.Content,
			mode,
		); err != nil {
			return result, fmt.Errorf(
				"write %s: %w",
				file.Path,
				err,
			)
		}

		if exists {
			result.Overwritten = append(result.Overwritten, file.Path)
		} else {
			result.Created = append(result.Created, file.Path)
		}
	}

	return result, nil
}
