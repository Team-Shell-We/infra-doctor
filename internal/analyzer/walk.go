package analyzer

// excludedDirs : infrastructure.go/github.go/profile.go의 filepath.Walk가
// 내려가지 않을 디렉터리 이름. cmd/init.go가 예전에 examples/ 폴백에서
// 쓰던 목록에 target(Maven 빌드 산출물), .infra-doctor(자체 생성 디렉터리)를 추가.
var excludedDirs = map[string]bool{
	".git":          true,
	".gradle":       true,
	".idea":         true,
	"build":         true,
	"node_modules":  true,
	"target":        true,
	".infra-doctor": true,
}

// shouldSkipDir : 이 이름의 디렉터리는 재귀 탐색에서 제외해야 하는지.
func shouldSkipDir(name string) bool {
	return excludedDirs[name]
}
