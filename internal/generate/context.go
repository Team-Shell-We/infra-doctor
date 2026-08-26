package generate

import (
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/doctor"
	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

type Database struct {
	Name        string
	ServiceName string
	Image       string
	Port        int
	DataPath    string

	EnvVars         map[string]string
	HealthCheckTest string // 이미 YAML 배열 형태로 포맷된 문자열, 예: ["CMD", "redis-cli", "ping"]
}

type Context struct {
	ProjectName string

	Framework      string
	Runtime        string
	RuntimeVersion string

	BuildTool    string
	BuildCommand string
	RunCommand   string
	ArtifactPath string

	ApplicationPort int
	HealthPath      string

	ServiceName string
	DockerImage string

	Databases []Database
	Redis     bool

	Namespace string
	Replicas  int

	HasDockerfile bool

	NeedsDocker      bool
	NeedsCompose     bool
	NeedsHealthCheck bool
	NeedsNginx       bool
	NeedsCI          bool

	// Lang과 아래 배너 필드는 생성 파일 안내 주석용이다. 각 target의 generator.go가 자신이 쓰는 필드만 채워 RenderTemplate에 넘기고 나머지는 비워도 된다
	Lang            string
	Header          string
	PortNote        string
	CredentialsNote string
	ImageNote       string
	ResourcesNote   string
	BranchNote      string
}

func BuildContext(info project.Info, diagnosis *doctor.Result, config Config, lang string) (Context, []string) {
	ctx := Context{
		Lang:            lang,
		ProjectName:     valueOr(config.ProjectName, "application"),
		Framework:       "spring-boot",
		Runtime:         "java",
		RuntimeVersion:  valueOr(info.Framework.Java.Version, "21"),
		BuildTool:       strings.ToLower(info.Framework.BuildTool.Type),
		ApplicationPort: intOr(config.ApplicationPort, 8080),
		HealthPath:      valueOr(config.HealthPath, defaultHealthPath(info)),
		ServiceName:     valueOr(config.ServiceName, "application"),
		DockerImage:     valueOr(config.DockerImage, "application:latest"),
		Namespace:       valueOr(config.Namespace, "default"),
		Replicas:        intOr(config.Replicas, 1),
		Redis:           info.Database.Redis != nil && info.Database.Redis.Enabled,
		HasDockerfile:   info.Infrastructure.Docker.Enabled,
		NeedsCI:         len(info.Github.Workflows) == 0,
	}

	switch ctx.BuildTool {
	case "gradle":
		ctx.BuildCommand = "./gradlew clean bootJar"
		ctx.ArtifactPath = "build/libs/*.jar"
	case "maven":
		ctx.BuildCommand = "./mvnw clean package -DskipTests"
		ctx.ArtifactPath = "target/*.jar"
	}
	ctx.RunCommand = "java -jar app.jar"

	if db, ok := databaseFor(info.Database.Primary.Type); ok {
		ctx.Databases = append(ctx.Databases, db)
	}
	applyDoctorResult(&ctx, diagnosis)

	var warnings []string
	if info.Framework.Java.Version == "" {
		warnings = append(warnings, "Java version was not detected; using 21")
	}
	return ctx, warnings
}

func applyDoctorResult(ctx *Context, result *doctor.Result) {
	if result == nil {
		return
	}
	for _, check := range result.Checks {
		switch check.ID {
		case "docker":
			ctx.NeedsDocker = !check.Passed
		case "docker_compose":
			ctx.NeedsCompose = !check.Passed
		case "health_check":
			ctx.NeedsHealthCheck = !check.Passed
		case "reverse_proxy":
			ctx.NeedsNginx = !check.Passed
		}
	}
}

func buildDetections(info project.Info) []string {
	var detections []string
	if info.Framework.BuildTool.Type != "" {
		detections = append(detections, info.Framework.BuildTool.Type)
	}
	if info.Framework.Java.Version != "" {
		detections = append(detections, "Java "+info.Framework.Java.Version)
	}
	if info.Database.Primary.Type != "" && info.Database.Primary.Type != "Unknown" {
		detections = append(detections, info.Database.Primary.Type)
	}
	return detections
}

// databaseFor: DB별 최소 필수 환경변수를 채운다. 공식 postgres/mysql/mariadb 이미지는 이 값(또는 대체 값) 없이는 컨테이너가 바로 종료된다
func databaseFor(name string) (Database, bool) {
	switch strings.ToLower(name) {
	case "postgresql":
		return Database{
			Name: "postgres", ServiceName: "postgres", Image: "postgres:16", Port: 5432, DataPath: "/var/lib/postgresql/data",
			EnvVars:         map[string]string{"POSTGRES_PASSWORD": "changeme"},
			HealthCheckTest: `["CMD-SHELL", "pg_isready -U postgres"]`,
		}, true
	case "mysql":
		return Database{
			Name: "mysql", ServiceName: "mysql", Image: "mysql:8", Port: 3306, DataPath: "/var/lib/mysql",
			EnvVars:         map[string]string{"MYSQL_ROOT_PASSWORD": "changeme"},
			HealthCheckTest: `["CMD", "mysqladmin", "ping", "-h", "localhost"]`,
		}, true
	case "mariadb":
		return Database{
			Name: "mariadb", ServiceName: "mariadb", Image: "mariadb:11", Port: 3306, DataPath: "/var/lib/mysql",
			EnvVars:         map[string]string{"MARIADB_ROOT_PASSWORD": "changeme"},
			HealthCheckTest: `["CMD", "mysqladmin", "ping", "-h", "localhost"]`,
		}, true
	default:
		return Database{}, false
	}
}

// defaultHealthPath: Actuator 없는 프로젝트에 "/actuator/health"를 기본값으로 쓰면 404로 HEALTHCHECK/probe가 항상 실패한다
func defaultHealthPath(info project.Info) string {
	if info.Dependencies.Actuator.Enabled {
		return "/actuator/health"
	}
	return "/"
}

func valueOr(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func intOr(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}
