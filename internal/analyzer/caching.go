package analyzer

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

var (
	cacheableAnnotationRegex   = regexp.MustCompile(`@(?:Cacheable|CachePut|CacheEvict)\b`)
	cacheConfigAnnotationRegex = regexp.MustCompile(`@CacheConfig\b`)
	tokenCandidateNameRegex    = regexp.MustCompile(`(?i)(RefreshToken|Session)`)
)

// AnalyzeCaching : .java 소스에서 캐시 어노테이션 사용과 토큰/세션 저장 패턴을 탐지한다.
// AnalyzeAPI(api.go)와 동일한 filepath.Walk + shouldSkipDir 패턴을 쓴다.
func AnalyzeCaching(root string) (*project.CachingInfo, error) {

	info := &project.CachingInfo{}

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

		content := string(data)

		info.CacheableCount += len(cacheableAnnotationRegex.FindAllString(content, -1))

		if cacheConfigAnnotationRegex.MatchString(content) {
			info.HasCacheConfig = true
		}

		name := strings.TrimSuffix(filepath.Base(path), ".java")
		if tokenCandidateNameRegex.MatchString(name) {
			info.TokenCandidates = append(info.TokenCandidates, name)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return info, nil
}
