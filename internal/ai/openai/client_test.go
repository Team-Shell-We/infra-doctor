package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
)

func newTestClient(baseURL string) *Client {
	return &Client{apiKey: "test-key", baseURL: baseURL, httpClient: http.DefaultClient}
}

func TestCompleteReturnsContentFromFirstChoice(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body chatRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		if body.Model != defaultModel {
			t.Errorf("Model = %q, want default %q when req.Model is empty", body.Model, defaultModel)
		}

		if len(body.Messages) != 2 || body.Messages[0].Role != "system" || body.Messages[1].Role != "user" {
			t.Errorf("Messages = %+v, want [system, user]", body.Messages)
		}

		json.NewEncoder(w).Encode(chatResponse{
			Choices: []chatChoice{{Message: chatMessage{Role: "assistant", Content: "hello"}}},
		})
	}))
	defer server.Close()

	resp, err := newTestClient(server.URL).Complete(context.Background(), ai.CompletionRequest{
		SystemPrompt: "sys", UserPrompt: "user",
	})
	if err != nil {
		t.Fatalf("Complete returned error: %v", err)
	}

	if resp.Content != "hello" {
		t.Errorf("Content = %q, want %q", resp.Content, "hello")
	}
}

func TestCompleteSetsJSONResponseFormatWhenRequested(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body chatRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.ResponseFormat == nil || body.ResponseFormat.Type != "json_object" {
			t.Errorf("ResponseFormat = %+v, want type json_object", body.ResponseFormat)
		}

		json.NewEncoder(w).Encode(chatResponse{Choices: []chatChoice{{Message: chatMessage{Content: "{}"}}}})
	}))
	defer server.Close()

	_, err := newTestClient(server.URL).Complete(context.Background(), ai.CompletionRequest{
		SystemPrompt: "sys", UserPrompt: "user", JSONMode: true,
	})
	if err != nil {
		t.Fatalf("Complete returned error: %v", err)
	}
}

func TestCompleteReturnsErrorOnEmptyChoices(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(chatResponse{Choices: []chatChoice{}})
	}))
	defer server.Close()

	_, err := newTestClient(server.URL).Complete(context.Background(), ai.CompletionRequest{})
	if err == nil {
		t.Fatal("Complete with empty choices: want error, got nil")
	}
}

func TestCompleteReturnsErrInvalidAPIKeyOn401(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	_, err := newTestClient(server.URL).Complete(context.Background(), ai.CompletionRequest{})
	if err != ai.ErrInvalidAPIKey {
		t.Errorf("Complete on 401: err = %v, want %v", err, ai.ErrInvalidAPIKey)
	}
}

func TestCompleteReturnsErrorOnNon200Status(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("boom"))
	}))
	defer server.Close()

	_, err := newTestClient(server.URL).Complete(context.Background(), ai.CompletionRequest{})
	if err == nil {
		t.Fatal("Complete on 500: want error, got nil")
	}
}

func TestVerifyCredentialsSucceedsOn200(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization header = %q, want %q", got, "Bearer test-key")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	if err := newTestClient(server.URL).VerifyCredentials(context.Background()); err != nil {
		t.Errorf("VerifyCredentials returned error: %v", err)
	}
}

func TestVerifyCredentialsReturnsErrInvalidAPIKeyOn401(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	err := newTestClient(server.URL).VerifyCredentials(context.Background())
	if err != ai.ErrInvalidAPIKey {
		t.Errorf("VerifyCredentials on 401: err = %v, want %v", err, ai.ErrInvalidAPIKey)
	}
}

func TestVerifyCredentialsReturnsErrorOnUnexpectedStatus(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	if err := newTestClient(server.URL).VerifyCredentials(context.Background()); err == nil {
		t.Fatal("VerifyCredentials on 503: want error, got nil")
	}
}

func TestNewSetsDefaultBaseURL(t *testing.T) {

	c := New("key")

	if c.baseURL != defaultBaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL, defaultBaseURL)
	}
}
