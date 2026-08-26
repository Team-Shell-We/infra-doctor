package utils

import (
	"context"
	"testing"
)

type fakeReleaseFetcher struct {
	release GitHubRelease
	err     error
}

func (f fakeReleaseFetcher) LatestRelease(context.Context) (GitHubRelease, error) {
	return f.release, f.err
}

func TestGetUpdateInfoFallsBackToDevForPseudoVersion(t *testing.T) {

	// go install ...@브랜치처럼 태그가 아닌 방식으로 설치했을 때 나오는 형태
	pseudoVersion := "v0.3.1-0.20260101120000-abcdef123456"

	fetcher := fakeReleaseFetcher{release: GitHubRelease{TagName: "v0.4.0"}}

	info, err := getUpdateInfo(context.Background(), pseudoVersion, fetcher)
	if err != nil {
		t.Fatalf("getUpdateInfo failed: %v", err)
	}

	if info.CurrentVersion != "dev" {
		t.Errorf("CurrentVersion = %q, want dev", info.CurrentVersion)
	}
	if info.UpdateAvailable {
		t.Error("UpdateAvailable = true, want false for an unparseable current version")
	}
}

func TestGetUpdateInfoFallsBackToDevForPrereleaseTag(t *testing.T) {

	fetcher := fakeReleaseFetcher{release: GitHubRelease{TagName: "v0.4.0"}}

	info, err := getUpdateInfo(context.Background(), "v1.4.0-rc1", fetcher)
	if err != nil {
		t.Fatalf("getUpdateInfo failed: %v", err)
	}

	if info.CurrentVersion != "dev" {
		t.Errorf("CurrentVersion = %q, want dev", info.CurrentVersion)
	}
}

func TestGetUpdateInfoDetectsAvailableUpdate(t *testing.T) {

	fetcher := fakeReleaseFetcher{release: GitHubRelease{TagName: "v0.3.0"}}

	info, err := getUpdateInfo(context.Background(), "v0.2.0", fetcher)
	if err != nil {
		t.Fatalf("getUpdateInfo failed: %v", err)
	}

	if !info.UpdateAvailable {
		t.Error("UpdateAvailable = false, want true (v0.2.0 -> v0.3.0)")
	}
}

func TestGetUpdateInfoNoUpdateWhenUpToDate(t *testing.T) {

	fetcher := fakeReleaseFetcher{release: GitHubRelease{TagName: "v0.2.0"}}

	info, err := getUpdateInfo(context.Background(), "v0.2.0", fetcher)
	if err != nil {
		t.Fatalf("getUpdateInfo failed: %v", err)
	}

	if info.UpdateAvailable {
		t.Error("UpdateAvailable = true, want false when already on the latest version")
	}
}

func TestGetUpdateInfoDevBuildNeverReportsUpdate(t *testing.T) {

	fetcher := fakeReleaseFetcher{release: GitHubRelease{TagName: "v0.2.0"}}

	info, err := getUpdateInfo(context.Background(), "(devel)", fetcher)
	if err != nil {
		t.Fatalf("getUpdateInfo failed: %v", err)
	}

	if info.CurrentVersion != "dev" || info.UpdateAvailable {
		t.Errorf("got %+v, want CurrentVersion=dev UpdateAvailable=false", info)
	}
}
