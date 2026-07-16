package project

type Info struct {
	Framework FrameworkInfo
	Database  DatabaseInfo
	Docker    DockerInfo
	Github    GithubInfo
	Profiles  []ProfileInfo
}

// FrameworkInfo : 프로젝트에서 사용하는 프레임워크 정보
type FrameworkInfo struct {
	BuildTool  BuildToolInfo
	SpringBoot SpringBootInfo
	Java       JavaInfo
}

// BuildToolInfo : 프로젝트에서 사용하는 빌드 도구 정보
type BuildToolInfo struct {
	Type string
	File string
	Path string
}

// SpringBootInfo : 프로젝트에서 사용하는 Spring Boot 정보
type SpringBootInfo struct {
	Enabled bool
	Version string
}

// JavaInfo : 프로젝트에서 사용하는 Java 정보
type JavaInfo struct {
	Version string
}

// DatabaseInfo : 프로젝트에서 사용하는 데이터베이스 정보
type DatabaseInfo struct {
	Primary Database
	Redis   *RedisInfo
}

// Database : 프로젝트에서 사용하는 메인 데이터베이스 정보
type Database struct {
	Type string
}

// RedisInfo : 프로젝트에서 사용하는 Redis 정보
type RedisInfo struct {
	Enabled bool
}

// DockerInfo : 프로젝트에서 사용하는 도커 정보
type DockerInfo struct {
	Dockerfiles []DockerfileInfo
	Compose     []ComposeInfo
}

// DockerfileInfo : 프로젝트에서 사용하는 Dockerfile 정보
type DockerfileInfo struct {
	File string
	Path string
}

// ComposeInfo : 프로젝트에서 사용하는 Docker Compose 정보
type ComposeInfo struct {
	File string
	Path string
}

// GithubInfo : 프로젝트에서 사용하는 GitHub 정보
type GithubInfo struct {
	Workflows []WorkflowInfo
}

// WorkflowInfo : 프로젝트에서 사용하는 GitHub Workflow 정보
type WorkflowInfo struct {
	Name string
	File string
	Path string

	Triggers []TriggerInfo
	Jobs      []JobInfo
}

// TriggerInfo : GitHub Workflow Trigger 정보
type TriggerInfo struct {
	Event    string
	Branches []string
}

// JobInfo : GitHub Workflow Job 정보
type JobInfo struct {
	Name  string
	Steps []StepInfo
}

// StepInfo : GitHub Workflow Step 정보
type StepInfo struct {
	Name string

	Uses string
	Run  string

	With map[string]string
	Env  map[string]string
}

// ProfileInfo : 프로젝트에서 사용하는 프로파일 정보
type ProfileInfo struct {
	Name string

	// Path : application-prod.yml 내용 분석 시 재사용
	Path string // 내부 분석용
	File string // 외부 출력용
}