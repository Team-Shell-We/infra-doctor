package project

type Info struct {
	Framework      FrameworkInfo
	Dependencies   DependencyInfo
	Database       DatabaseInfo
	Infrastructure InfrastructureInfo
	Github         GithubInfo
	Profiles       []ProfileInfo
	API            APIInfo
}

type FrameworkInfo struct {
	BuildTool  BuildToolInfo
	SpringBoot SpringBootInfo
	Java       JavaInfo
	Modules    ModuleInfo
}

// ModuleInfo : Gradle 멀티모듈 구조 정보
type ModuleInfo struct {
	Count int
}

// APIInfo : 프로젝트의 Controller/엔드포인트 규모
type APIInfo struct {
	ControllerCount int
	EndpointCount   int
}

type DependencyInfo struct {
	Security SecurityInfo
	JPA      JPAInfo
	Kafka    KafkaInfo
	AWS      AWSInfo
	Lombok   LombokInfo
	Actuator ActuatorInfo
	OpenAPI  OpenAPIInfo
}

// SecurityInfo : 프로젝트에서 사용하는 Spring Security 정보
type SecurityInfo struct {
	Enabled bool
}

type JPAInfo struct {
	Enabled bool
}

type KafkaInfo struct {
	Enabled bool
}

type AWSInfo struct {
	Enabled bool
}

type LombokInfo struct {
	Enabled bool
}

// ActuatorInfo : 프로젝트에서 사용하는 Spring Boot Actuator 정보
type ActuatorInfo struct {
	Enabled bool
}

type OpenAPIInfo struct {
	Enabled bool
}

type BuildToolInfo struct {
	Type string
	File string
	Path string
}

type SpringBootInfo struct {
	Enabled bool
	Version string
}

type JavaInfo struct {
	Version string
}

type DatabaseInfo struct {
	Primary Database
	Redis   *RedisInfo
}

type Database struct {
	Type string
}

type RedisInfo struct {
	Enabled bool
}

type InfrastructureInfo struct {
	Docker      DockerInfo
	Kubernetes  KubernetesInfo
	Nginx       NginxInfo
	Terraform   TerraformInfo
	HealthCheck HealthCheckInfo
	Monitoring  MonitoringInfo
	LogRotation LogRotationInfo
	Backup      BackupInfo
}

type NginxInfo struct {
	Enabled bool
}

type TerraformInfo struct {
	Enabled bool
}

type HealthCheckInfo struct {
	Enabled bool
}

type MonitoringInfo struct {
	Enabled bool
}

type LogRotationInfo struct {
	Enabled bool
}

// BackupInfo : 프로젝트에서 사용하는 DB 백업 정보
type BackupInfo struct {
	Enabled bool
}

type DockerInfo struct {
	Enabled bool

	Dockerfiles []DockerfileInfo
	Compose     []ComposeInfo
}

type KubernetesInfo struct {
	Enabled bool

	Files []KubernetesFileInfo

	// Replicas : manifest에서 찾은 가장 큰 replicas 값(없으면 0)
	Replicas int
}

type KubernetesFileInfo struct {
	File string
	Path string
}

type DockerfileInfo struct {
	File string
	Path string
}

type ComposeInfo struct {
	File string
	Path string
}

type GithubInfo struct {
	Workflows []WorkflowInfo
}

type WorkflowInfo struct {
	Name string
	File string
	Path string

	Triggers []TriggerInfo
	Jobs     []JobInfo
}

type TriggerInfo struct {
	Event    string
	Branches []string
}

type JobInfo struct {
	Name  string
	Steps []StepInfo
}

type StepInfo struct {
	Name string

	Uses string
	Run  string

	With map[string]string
	Env  map[string]string
}

type ProfileInfo struct {
	Name string

	Path string // 내부 분석용
	File string // 외부 출력용
}
