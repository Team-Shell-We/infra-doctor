package analyzer

// excludedDirs: filepath.Walk 재귀 탐색에서 건너뛸 디렉터리. target은 Maven 빌드 산출물, .infra-doctor는 자체 생성 디렉터리
var excludedDirs = map[string]bool{
	".git":          true,
	".gradle":       true,
	".idea":         true,
	"build":         true,
	"node_modules":  true,
	"target":        true,
	".infra-doctor": true,
}

func shouldSkipDir(name string) bool {
	return excludedDirs[name]
}
