package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultReleaseURL = "https://api.github.com/repos/Team-Shell-We/infra-doctor/releases/latest"

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

type ReleaseFetcher interface {
	LatestRelease(context.Context) (GitHubRelease, error)
}

type GitHubClient struct {
	httpClient *http.Client
	releaseURL string
}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		releaseURL: defaultReleaseURL,
	}
}

// 단위 테스트에서 httptest.Server를 주입하기 위한 생성자입니다.
func NewGitHubClientWithHTTPClient(
	httpClient *http.Client,
	releaseURL string,
) *GitHubClient {
	return &GitHubClient{
		httpClient: httpClient,
		releaseURL: releaseURL,
	}
}

func (c *GitHubClient) LatestRelease(
	ctx context.Context,
) (GitHubRelease, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.releaseURL,
		nil,
	)
	if err != nil {
		return GitHubRelease{}, fmt.Errorf(
			"create GitHub release request: %w",
			err,
		)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "infra-doctor")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GitHubRelease{}, fmt.Errorf(
			"request GitHub release: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GitHubRelease{}, fmt.Errorf(
			"GitHub Releases API returned %s",
			resp.Status,
		)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return GitHubRelease{}, fmt.Errorf(
			"decode GitHub release response: %w",
			err,
		)
	}

	if release.TagName == "" {
		return GitHubRelease{}, fmt.Errorf(
			"GitHub release response has no tag_name",
		)
	}

	return release, nil
}
