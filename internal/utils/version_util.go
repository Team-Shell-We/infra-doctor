package utils

//read cli version

import (
	"runtime/debug"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

func Version() string { // 현재 CLI 버전을 읽음
	info, ok := debug.ReadBuildInfo() // 빌드 정보를 읽고 빌드 정보를 정상적으로 읽었는지의 여부를 저장
	/*
		info 안에 담길 정보
		GoVersion: go1.24.3
		Main.Path: github.com/example/infra-doctor
		Main.Version: v1.2.0 : 출력되어야 할 정보
	*/

	if !ok {
		return "unknown"
	}

	if info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "development"
	}

	return info.Main.Version
}

func GoVersion() string { // CLI를 빌드할 때 사용한 GO 버전을 읽음
	info, ok := debug.ReadBuildInfo() //ex: info.GoVersion → go1.24.3
	if !ok {
		return "unknown"
	}

	version := strings.TrimPrefix(info.GoVersion, "go")
	parts := strings.Split(version, ".")

	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}

	return version // 1.24
}

// LocalizeVersion : Version()/GoVersion()이 반환하는 "unknown"/
// "development" 폴백 문자열만 언어에 맞게 바꾼다. 실제 버전 번호(예:
// "v0.1.0", "1.26")는 그대로 둔다. cmd/version.go와
// internal/utils/help_util.go가 공유해서 쓴다.
func LocalizeVersion(lang, version string) string {

	switch version {
	case "unknown":
		return i18n.Get(lang, "version.unknown")
	case "development":
		return i18n.Get(lang, "version.development")
	default:
		return version
	}
}
