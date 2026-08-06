package ai

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Credentials is stored at ~/.infra-doctor/config.json — the user's home
// directory, holding AI provider credentials. This is unrelated to the
// per-project .infra-doctor/config.yaml that `infra-doctor init` creates
// (internal/utils/init_util.go): that one holds analysis settings for a
// single project, this one holds global AI login state shared across every
// project the CLI is run against.
type Credentials struct {
	Provider string `json:"provider"`
	APIKey   string `json:"apiKey"`
	Login    bool   `json:"login"`
}

const (
	credentialsDirName  = ".infra-doctor"
	credentialsFileName = "config.json"
)

// Path returns ~/.infra-doctor/config.json for the current user.
func Path() (string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, credentialsDirName, credentialsFileName), nil
}

// Load reads the current user's stored credentials, returning
// ErrNotLoggedIn if `infra-doctor login` has never been run successfully.
func Load() (*Credentials, error) {

	path, err := Path()
	if err != nil {
		return nil, err
	}

	return LoadFrom(path)
}

// Save persists creds to ~/.infra-doctor/config.json.
func Save(creds Credentials) error {

	path, err := Path()
	if err != nil {
		return err
	}

	return SaveTo(path, creds)
}

// LoadFrom/SaveTo take an explicit path so credential persistence can be
// unit tested against a temp directory instead of the real home directory.

func LoadFrom(path string) (*Credentials, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotLoggedIn
		}
		return nil, err
	}

	var creds Credentials

	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	if !creds.Login {
		return nil, ErrNotLoggedIn
	}

	return &creds, nil
}

func SaveTo(path string, creds Credentials) error {

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return err
	}

	return nil
}
