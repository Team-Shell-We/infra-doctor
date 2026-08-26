package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	infraDoctorDirectoryName = ".infra-doctor"
	configFileName           = "config.yaml"
	gitignoreFileName        = ".gitignore"
)

var ErrAlreadyInitialized = errors.New("Infra Doctor is already initialized")

type Result struct {
	CreatedPaths []string
}

func Initialize(projectDir string) (*Result, error) {
	if projectDir == "" {
		return nil, errors.New("project directory is empty")
	}

	infraDoctorDir := filepath.Join(
		projectDir,
		infraDoctorDirectoryName,
	)

	exists, err := pathExists(infraDoctorDir)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to check initialization directory: %w",
			err,
		)
	}

	if exists {
		return nil, ErrAlreadyInitialized
	}

	if err := os.Mkdir(infraDoctorDir, 0755); err != nil {
		return nil, fmt.Errorf(
			"failed to create directory %s: %w",
			infraDoctorDir,
			err,
		)
	}

	initializationSucceeded := false
	defer func() {
		if !initializationSucceeded {
			_ = os.RemoveAll(infraDoctorDir)
		}
	}()

	configPath := filepath.Join(
		infraDoctorDir,
		configFileName,
	)

	if err := writeFile(configPath, defaultConfig, 0644); err != nil {
		return nil, fmt.Errorf(
			"failed to create config file %s: %w",
			configPath,
			err,
		)
	}

	gitignorePath := filepath.Join(
		infraDoctorDir,
		gitignoreFileName,
	)

	if err := writeFile(gitignorePath, defaultGitignore, 0644); err != nil {
		return nil, fmt.Errorf(
			"failed to create gitignore file %s: %w",
			gitignorePath,
			err,
		)
	}

	initializationSucceeded = true

	return &Result{
		CreatedPaths: []string{
			".infra-doctor",
			".infra-doctor/config.yaml",
			".infra-doctor/.gitignore",
		},
	}, nil
}

func writeFile(path string, content string, permission os.FileMode) error {
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		permission)

	if err != nil {
		return err
	}

	if _, err = file.WriteString(content); err != nil {
		_ = file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)

	switch {
	case err == nil:
		return true, nil

	case os.IsNotExist(err):
		return false, nil

	default:
		return false, err
	}

}
