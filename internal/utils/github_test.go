package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLatestReleaseReturnsTagName(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"tag_name":"v0.4.0"}`))
	}))
	defer server.Close()

	client := NewGitHubClientWithHTTPClient(server.Client(), server.URL)

	release, err := client.LatestRelease(context.Background())
	if err != nil {
		t.Fatalf("LatestRelease failed: %v", err)
	}

	if release.TagName != "v0.4.0" {
		t.Errorf("TagName = %q, want %q", release.TagName, "v0.4.0")
	}
}

func TestLatestReleaseReturnsErrorOnNon200Status(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewGitHubClientWithHTTPClient(server.Client(), server.URL)

	if _, err := client.LatestRelease(context.Background()); err == nil {
		t.Fatal("LatestRelease on 404: want error, got nil")
	}
}

func TestLatestReleaseReturnsErrorOnMalformedJSON(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not json`))
	}))
	defer server.Close()

	client := NewGitHubClientWithHTTPClient(server.Client(), server.URL)

	if _, err := client.LatestRelease(context.Background()); err == nil {
		t.Fatal("LatestRelease on malformed JSON: want error, got nil")
	}
}

func TestLatestReleaseReturnsErrorOnEmptyTagName(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"tag_name":""}`))
	}))
	defer server.Close()

	client := NewGitHubClientWithHTTPClient(server.Client(), server.URL)

	if _, err := client.LatestRelease(context.Background()); err == nil {
		t.Fatal("LatestRelease with empty tag_name: want error, got nil")
	}
}

func TestNewGitHubClientSetsDefaultReleaseURL(t *testing.T) {

	client := NewGitHubClient()

	if client.releaseURL != defaultReleaseURL {
		t.Errorf("releaseURL = %q, want %q", client.releaseURL, defaultReleaseURL)
	}
}
