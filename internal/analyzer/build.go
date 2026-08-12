package analyzer

import "os"

func HasGradle(root string) bool {

	// os.Stat() : 파일이 존재하는지 확인하는 함수

	if _, err := os.Stat(root + "/build.gradle"); err == nil {
		return true
	}

	if _, err := os.Stat(root + "/build.gradle.kts"); err == nil {
		return true
	}

	return false
}