package compose

import (
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
)

func TestPlanIncludesHealthCheckAndDependsOnCondition(t *testing.T) {
	ctx := generate.Context{
		ApplicationPort: 8080,
		HealthPath:      "/",
		Lang:            "en",
		Databases: []generate.Database{
			{
				Name: "postgres", ServiceName: "postgres", Image: "postgres:16",
				Port: 5432, DataPath: "/var/lib/postgresql/data",
				EnvVars:         map[string]string{"POSTGRES_PASSWORD": "changeme"},
				HealthCheckTest: `["CMD-SHELL", "pg_isready -U postgres"]`,
			},
		},
	}

	files, err := (Generator{}).Plan(ctx)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}
	content := string(files[0].Content)

	// plain depends_on만으로는 컨테이너 "시작"만 확인되고 DB가 연결을 받을
	// 준비가 됐는지는 보장 안 됨
	// condition: service_healthy가 있어야 하고, 그러려면 DB 쪽에도
	// healthcheck이 있어야 함
	if !strings.Contains(content, "condition: service_healthy") {
		t.Error("depends_on should wait for service_healthy")
	}
	if !strings.Contains(content, "POSTGRES_PASSWORD: changeme") {
		t.Error("postgres service must have its required password env var, or the official image refuses to start")
	}
	if !strings.Contains(content, `pg_isready`) {
		t.Error("postgres service should have its own healthcheck")
	}
	if !strings.Contains(content, "healthcheck:") {
		t.Error("app service must have a healthcheck block")
	}
}

func TestPlanOmitsDependsOnWhenNoDependencies(t *testing.T) {
	ctx := generate.Context{ApplicationPort: 8080, HealthPath: "/"}

	files, err := (Generator{}).Plan(ctx)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}
	content := string(files[0].Content)

	if strings.Contains(content, "depends_on:") {
		t.Error("no databases/redis detected, should not emit depends_on")
	}
}
