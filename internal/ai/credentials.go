package ai

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Credentials : ~/.infra-doctor/config.json(사용자 홈 디렉터리)에 저장되는 AI provider 자격증명. `infra-doctor init`이 프로젝트 디렉터리에 만드는
// .infra-doctor/config.yaml(internal/utils/init_util.go, 프로젝트별 분석 설정)과는 별개(어느 프로젝트에서 실행하든 공유되는 전역 로그인 상태)
type Credentials struct {
	Provider string `json:"provider"`
	APIKey   string `json:"apiKey"`
	Login    bool   `json:"login"`
}

const (
	credentialsDirName  = ".infra-doctor"
	credentialsFileName = "config.json"
)

// Path : 현재 사용자의 ~/.infra-doctor/config.json 경로를 반환
func Path() (string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, credentialsDirName, credentialsFileName), nil
}

// Load : 저장된 자격증명을 읽음. 
// `infra-doctor login`을 성공적으로 실행한 적이 없으면 ErrNotLoggedIn을 반환
func Load() (*Credentials, error) {

	path, err := Path()
	if err != nil {
		return nil, err
	}

	return LoadFrom(path)
}

// Save : creds를 ~/.infra-doctor/config.json에 저장
func Save(creds Credentials) error {

	path, err := Path()
	if err != nil {
		return err
	}

	return SaveTo(path, creds)
}

// LoadFrom/SaveTo : 경로를 직접 받아, 실제 홈 디렉터리 대신 임시 디렉터리로 단위 테스트할 수 있게 해줌

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
