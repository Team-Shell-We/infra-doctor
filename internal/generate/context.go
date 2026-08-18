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
}

func BuildContext(info project.Info, diagnosis *doctor.Result, config Config) (Context, []string) {
	ctx := Context{
		ProjectName:     valueOr(config.ProjectName, "application"),
		Framework:       "spring-boot",
		Runtime:         "java",
		RuntimeVersion:  valueOr(info.Framework.Java.Version, "21"),
		BuildTool:       strings.ToLower(info.Framework.BuildTool.Type),
		ApplicationPort: intOr(config.ApplicationPort, 8080),
		HealthPath:      valueOr(config.HealthPath, "/actuator/health"),
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

func databaseFor(name string) (Database, bool) {
	switch strings.ToLower(name) {
	case "postgresql":
		return Database{Name: "postgres", ServiceName: "postgres", Image: "postgres:16", Port: 5432, DataPath: "/var/lib/postgresql/data"}, true
	case "mysql":
		return Database{Name: "mysql", ServiceName: "mysql", Image: "mysql:8", Port: 3306, DataPath: "/var/lib/mysql"}, true
	case "mariadb":
		return Database{Name: "mariadb", ServiceName: "mariadb", Image: "mariadb:11", Port: 3306, DataPath: "/var/lib/mysql"}, true
	default:
		return Database{}, false
	}
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
