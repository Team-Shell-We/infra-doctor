package utils

//read cli version

import (
	"fmt"
	"runtime/debug"
	"strings"
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

func VersionString() string {
	return fmt.Sprintf(
		"Infra Doctor\n\nVersion : %s\n\nGo : %s",
		Version(),
		GoVersion(),
	)
}
