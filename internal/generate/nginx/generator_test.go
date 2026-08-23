package nginx

import (
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
)

func TestPlanIncludesForwardedProtoHeader(t *testing.T) {
	ctx := generate.Context{ServiceName: "application", ApplicationPort: 8080, Lang: "en"}

	files, err := (Generator{}).Plan(ctx)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}
	content := string(files[0].Content)

	// Spring Boot의 forwarded-header 처리가 기대하는 헤더 중 하나
	if !strings.Contains(content, "X-Forwarded-Proto") {
		t.Error("nginx.conf must forward X-Forwarded-Proto")
	}

	for _, want := range []string{"events {}", "http {"} {
		if !strings.Contains(content, want) {
			t.Errorf("nginx.conf must wrap the config in %q to be valid on its own", want)
		}
	}
}
