package analyzer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func FindProfiles(root string) ([]project.ProfileInfo, error) {

	var profiles []project.ProfileInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			if shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		name := info.Name()

		// application*.yml 또는 application*.yaml만 찾는다.
		// TODO: .properties 등 다른 형식도 지원
		if strings.HasPrefix(name, "application") &&
			(strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml")) {

			profile := "default"

			// application-local.yml -> local 추출
			if strings.Contains(name, "-") {
				start := strings.Index(name, "-") + 1
				end := strings.LastIndex(name, ".")
				profile = name[start:end]
			}

			profiles = append(profiles, project.ProfileInfo{
				Name: profile,
				File: name,
				Path: path,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return profiles, nil
}
