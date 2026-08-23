package generate

import (
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/doctor"
	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func TestBuildContextUsesScanAndDoctorResults(t *testing.T) {
	info := project.Info{
		Framework: project.FrameworkInfo{
			BuildTool: project.BuildToolInfo{Type: "Gradle"},
			Java:      project.JavaInfo{Version: "21"},
		},
		Database: project.DatabaseInfo{
			Primary: project.Database{Type: "PostgreSQL"},
			Redis:   &project.RedisInfo{Enabled: true},
		},
	}
	diagnosis := doctor.Analyze(&info)

	ctx, _ := BuildContext(info, diagnosis, Config{}, "en")

	if len(ctx.Databases) != 1 || ctx.Databases[0].Image != "postgres:16" {
		t.Fatalf("databases = %+v, want detected PostgreSQL service", ctx.Databases)
	}
	if !ctx.Redis {
		t.Fatal("Redis = false, want true from scan result")
	}
	if !ctx.NeedsCompose || !ctx.NeedsDocker || !ctx.NeedsHealthCheck || !ctx.NeedsNginx {
		t.Fatalf("doctor flags were not applied: %+v", ctx)
	}
}

// 회귀 테스트: Actuator가 없는 프로젝트에 "/actuator/health"를 기본값으로
// 넣으면 404만 반환해 HEALTHCHECK/probe가 항상 실패한다.
func TestBuildContextHealthPathReflectsActuator(t *testing.T) {
	withActuator := project.Info{Dependencies: project.DependencyInfo{
		Actuator: project.ActuatorInfo{Enabled: true},
	}}
	ctx, _ := BuildContext(withActuator, nil, Config{}, "en")
	if ctx.HealthPath != "/actuator/health" {
		t.Errorf("HealthPath = %q, want /actuator/health when Actuator is detected", ctx.HealthPath)
	}

	withoutActuator := project.Info{}
	ctx, _ = BuildContext(withoutActuator, nil, Config{}, "en")
	if ctx.HealthPath != "/" {
		t.Errorf("HealthPath = %q, want / when Actuator is not detected", ctx.HealthPath)
	}

	explicitOverride := project.Info{}
	ctx, _ = BuildContext(explicitOverride, nil, Config{HealthPath: "/health"}, "en")
	if ctx.HealthPath != "/health" {
		t.Errorf("HealthPath = %q, want explicit config override /health", ctx.HealthPath)
	}
}
