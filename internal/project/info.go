package project

type Info struct {
	Framework      FrameworkInfo
	Dependencies   DependencyInfo
	Database       DatabaseInfo
	Infrastructure InfrastructureInfo
	Github         GithubInfo
	Profiles       []ProfileInfo
	API            APIInfo
	Schema         SchemaInfo
}

// SchemaInfo : JPA 엔티티/관계 정적 분석 결과
type SchemaInfo struct {
	Entities      []EntityInfo
	Relationships []RelationshipInfo
}

// EntityInfo : @Entity 클래스 하나에 대한 정보
type EntityInfo struct {
	Name    string // 클래스명, 예: "FacilityReview"
	Table   string // @Table(name=...), 없으면 ""
	File    string // 표시용 상대 경로
	Columns []ColumnInfo
	Indexes []IndexInfo
}

type ColumnInfo struct {
	Field     string
	Column    string // @Column(name=...), 없으면 Field와 동일
	Type      string // Java 필드 타입, 예: "Long", "String"
	Nullable  bool
	Unique    bool
	IsID      bool
	Inherited bool // @MappedSuperclass 부모 클래스에서 병합된 필드인지
}

// IndexInfo : @Table(indexes = {@Index(...)}) 항목 하나
type IndexInfo struct {
	Name    string
	Columns []string
}

// RelationshipInfo : @OneToOne/@OneToMany/@ManyToOne/@ManyToMany 필드 하나에 대한 선언
type RelationshipInfo struct {
	From       string // 선언한 엔티티
	To         string // 대상 엔티티(컬렉션 래퍼 제거한 타입)
	Type       string // "OneToOne" | "OneToMany" | "ManyToOne" | "ManyToMany"
	Field      string
	JoinColumn string // 소유 측 @JoinColumn(name=...), 없으면 ""
	MappedBy   string // 역방향 측 mappedBy 값, 없으면 ""
	Fetch      string // "LAZY" | "EAGER" | "" (미지정)
	Cascade    string // 원문 그대로, 예: "ALL" (CascadeType.ALL 포함 여부 판단용)
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
	Security  SecurityInfo
	JPA       JPAInfo
	Kafka     KafkaInfo
	AWS       AWSInfo
	Lombok    LombokInfo
	Actuator  ActuatorInfo
	OpenAPI   OpenAPIInfo
	Migration MigrationInfo
}

// MigrationInfo : 스키마 마이그레이션 도구(Flyway/Liquibase) 사용 여부
type MigrationInfo struct {
	Enabled bool
	Tool    string // "Flyway" | "Liquibase"
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
