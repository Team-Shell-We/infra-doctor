package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type UpdateInfo struct {
	CurrentVersion  string
	LatestVersion   string
	UpdateAvailable bool
}

type semanticVersion struct {
	Major int
	Minor int
	Patch int
}

func GetUpdateInfo() (UpdateInfo, error) {
	return getUpdateInfo(context.Background(), Version(), NewGitHubClient())
}

func getUpdateInfo(
	ctx context.Context,
	currentVersion string,
	fetcher ReleaseFetcher,
) (UpdateInfo, error) {
	release, err := fetcher.LatestRelease(ctx)
	if err != nil {
		return UpdateInfo{}, fmt.Errorf("fetch latest GitHub release: %w", err)
	}

	current := normalizeVersion(currentVersion)
	latest := normalizeVersion(release.TagName)

	currentSemVer, err := parseVersion(current)
	// go run 등 개발 환경이거나, go install ...@브랜치처럼 태그가 아닌 방식으로
	// 설치해 Go pseudo-version(예: v0.3.1-0.20260101120000-abcdef123456)이나
	// pre-release 태그가 나온 경우엔 버전 비교가 불가능하므로 "dev"로 취급한다.
	if current == "dev" || err != nil {
		return UpdateInfo{
			CurrentVersion:  "dev",
			LatestVersion:   latest,
			UpdateAvailable: false,
		}, nil
	}

	latestSemVer, err := parseVersion(latest)
	if err != nil {
		return UpdateInfo{}, fmt.Errorf("invalid latest version %q: %w", latest, err)
	}

	return UpdateInfo{
		CurrentVersion:  current,
		LatestVersion:   latest,
		UpdateAvailable: isUpdateAvailable(currentSemVer, latestSemVer),
	}, nil
}

// current가 latest보다 낮을 때만 true를 반환합니다.
func isUpdateAvailable(current, latest semanticVersion) bool {
	if current.Major != latest.Major {
		return current.Major < latest.Major
	}

	if current.Minor != latest.Minor {
		return current.Minor < latest.Minor
	}

	return current.Patch < latest.Patch
}

// "v1.2.3"을 Major=1, Minor=2, Patch=3으로 변환합니다.
func parseVersion(version string) (semanticVersion, error) {
	version = normalizeVersion(version)
	version = strings.TrimPrefix(version, "v")

	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return semanticVersion{}, fmt.Errorf(
			"version must use major.minor.patch format",
		)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil || major < 0 {
		return semanticVersion{}, fmt.Errorf(
			"major version must be a non-negative integer",
		)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil || minor < 0 {
		return semanticVersion{}, fmt.Errorf(
			"minor version must be a non-negative integer",
		)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil || patch < 0 {
		return semanticVersion{}, fmt.Errorf(
			"patch version must be a non-negative integer",
		)
	}

	return semanticVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func normalizeVersion(version string) string {
	version = strings.TrimSpace(version)

	if version == "" || version == "dev" || version == "development" || version == "unknown" || version == "(devel)" {
		return "dev"
	}

	if strings.HasPrefix(version, "v") {
		return version
	}

	return "v" + version
}
