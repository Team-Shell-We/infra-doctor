// .infra-doctor 디렉토리와 파일을 생성

package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	infraDoctorDirectoryName = ".infra-doctor" //infra-doctor 디렉토리 이름
	configFileName           = "config.yaml"   //infra-doctor 설정 파일 이름
	gitignoreFileName        = ".gitignore"    //gitignore 파일 이름
)

//이미 .infra-doctor 디렉토리가 존재 -> 오류 반환
var ErrAlreadyInitialized = errors.New("Infra Doctor is already initialized")

type Result struct {
	CreatedPaths []string
}

//infra-doctor dir와 기본 설정 파일들을 생성
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

//주어진 경로에 새 파일을 만들고 내용 작성
func writeFile(path string, content string, permission os.FileMode) error {
	//os.O_CREATE: 파일이 없으면 새로 생성
	//os.O_WRONLY: 쓰기 전용으로 열기
	//os.O_TRUNC: 파일이 이미 존재하면 기존 내용을 지우고 새로 작성
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

//전달받은 경로가 존재하는지 확인
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
