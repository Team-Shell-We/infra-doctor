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

// Semantic Version의 각 숫자를 저장합니다.
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

	// go run 등 개발 환경에서는 현재 버전을 비교할 수 없습니다.
	if current == "dev" {
		return UpdateInfo{
			CurrentVersion:  current,
			LatestVersion:   latest,
			UpdateAvailable: false,
		}, nil
	}

	updateAvailable, err := isUpdateAvailable(current, latest)
	if err != nil {
		return UpdateInfo{}, err
	}

	return UpdateInfo{
		CurrentVersion:  current,
		LatestVersion:   latest,
		UpdateAvailable: updateAvailable,
	}, nil
}

// current가 latest보다 낮을 때만 true를 반환합니다.
func isUpdateAvailable(current, latest string) (bool, error) {
	currentSemVer, err := parseVersion(current)
	if err != nil {
		return false, fmt.Errorf(
			"invalid current version %q: %w",
			current,
			err,
		)
	}

	latestSemVer, err := parseVersion(latest)
	if err != nil {
		return false, fmt.Errorf(
			"invalid latest version %q: %w",
			latest,
			err,
		)
	}

	if currentSemVer.Major != latestSemVer.Major {
		return currentSemVer.Major < latestSemVer.Major, nil
	}

	if currentSemVer.Minor != latestSemVer.Minor {
		return currentSemVer.Minor < latestSemVer.Minor, nil
	}

	return currentSemVer.Patch < latestSemVer.Patch, nil
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
