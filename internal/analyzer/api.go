package analyzer

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

var (
	controllerAnnotationRegex = regexp.MustCompile(`@(?:Rest)?Controller\b`)
	endpointAnnotationRegex   = regexp.MustCompile(`@(?:Get|Post|Put|Delete|Patch)Mapping\b`)
)

// AnalyzeAPI : .java 소스에서 Controller/엔드포인트 어노테이션 개수를 셈
func AnalyzeAPI(root string) (*project.APIInfo, error) {

	info := &project.APIInfo{}

	err := filepath.Walk(root, func(path string, fileInfo os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			if shouldSkipDir(fileInfo.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".java" {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		info.ControllerCount += len(controllerAnnotationRegex.FindAll(data, -1))
		info.EndpointCount += len(endpointAnnotationRegex.FindAll(data, -1))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return info, nil
}
