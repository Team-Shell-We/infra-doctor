// 현재 디렉토리가 spring boot인지 dectect!
package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//BuildTool: project가 사용하는 빌드 도구의 종류
type BuildTool string

const (
	BuildToolUnknown BuildTool = "unknown"	//인식 불가
	BuildToolMaven   BuildTool = "maven"	//Maven
	BuildToolGradle  BuildTool = "gradle"	//Gradle or Gradle Kotlin DSL
)

//DetectionResult: spring boot dectect 결과
type DetectionResult struct {
	Detected bool // spring boot 여부
	BuildTool BuildTool // 감지된 빌드 도구
	BuildFile string // 감지에 사용된 빌드 파일 이름
}

func DetectSpringBoot(dir string) (*DetectionResult, error) {
	//디렉토리가 empty 일 때 오류 반환
	if dir == "" {
		return nil, errors.New("project directory is empty")
	}

	//전달받은 경로가 실제로 있는지 확인
	info, err := os.Stat(dir)
	if err != nil {
		return nil, errors.New("project directory does not exist")
		
	}

	//전달받은 경로가 파일이고 dir가 아닌 경우 오류 반환
	if !info.IsDir() {
		return &DetectionResult{}, fmt.Errorf(
			"project directory %q is not a directory",
			dir,
		)
	}
	
	//gradle -> build.gradle, build.gradle.kts
	//maven -> pom.xml
	candidateFiles := []struct{
		FileName string
		BuildTool BuildTool
		Markers []string
	}{
		{
			FileName:  "build.gradle",
			BuildTool: BuildToolGradle,
			Markers: []string{
				"org.springframework.boot",
				"spring-boot-starter",
			},
		},
		{
			FileName:  "build.gradle.kts",
			BuildTool: BuildToolGradle,
			Markers: []string{
				"org.springframework.boot",
				"spring-boot-starter",
				},
		},
		{
			FileName:  "pom.xml",
			BuildTool: BuildToolMaven,
			Markers: []string{
				"org.springframework.boot",
				"spring-boot-starter",
				"spring-boot-maven-plugin",
			},
		},
	}

	//후보 빌드 파일들을 검사
	for _, candidate := range candidateFiles {
		buildFilePath := filepath.Join(dir, candidate.FileName)

		content, err := os.ReadFile(buildFilePath)
		if err != nil {
			if os.IsNotExist(err) { //해당 빌드 파일이 존재하지 않는 경우, 오류 반환
				continue
			}
			// 파일은 존재하지만 권한 등의 이유로 읽지 못한 경우 오류 반환
			return nil, fmt.Errorf(
				"cannot read build file %q: %w",
				buildFilePath,
				err,
			)
		}

		//빌드 파일에 spring boot 관련 문자열(Markers)이 포함되어 있는지 확인
		if containsAnyMarker(string(content), candidate.Markers) {
			return &DetectionResult{
				Detected: true,
				BuildTool: candidate.BuildTool,
				BuildFile: buildFilePath,
			}, nil
		}

	}
	
	//spring boot 프로젝트가 아닌 것으로 판단 -> 오류 반환
	return &DetectionResult{
		Detected: false,
		BuildTool: BuildToolUnknown,
	}, nil
}
	


func containsAnyMarker(content string, markers []string) bool {
	lowerContent := strings.ToLower(content) 

	//spring boot 관련 문자열을 하나씩 검사
	for _, marker := range markers {
		if strings.Contains(lowerContent, strings.ToLower(marker)) {
			return true //spring boot 프로젝트라고 판단
		}
	}

	return false //spring boot 관련 문자열이 하나도 없는 경우 -> 오류 반환
}

