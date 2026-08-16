package explain

import (
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func TestBuildStatusDocker(t *testing.T) {

	info := &project.Info{}
	info.Infrastructure.Docker.Enabled = true

	status := BuildStatus("docker", info)

	if len(status) == 0 {
		t.Fatal("expected at least one status item for docker")
	}

	found := false
	for _, item := range status {
		if item.Label == "Dockerfile" {
			found = true
			if !item.Present {
				t.Error("Dockerfile should be reported present when Docker.Enabled is true")
			}
		}
	}
	if !found {
		t.Error("expected a Dockerfile status item for the docker topic")
	}
}

func TestBuildStatusNeverInventsBeyondItsMapping(t *testing.T) {
	info := &project.Info{}
	info.Infrastructure.Docker.Enabled = true

	status := BuildStatus("docker", info)

	want := map[string]bool{"Dockerfile": true, "Docker Compose": true, "Health Check": true}

	if len(status) != len(want) {
		t.Fatalf("BuildStatus(\"docker\", ...) returned %d items, want exactly %d", len(status), len(want))
	}

	for _, item := range status {
		if !want[item.Label] {
			t.Errorf("unexpected status label %q for docker topic", item.Label)
		}
	}
}

func TestBuildStatusRedisAbsent(t *testing.T) {

	status := BuildStatus("redis", &project.Info{})

	if len(status) != 1 || status[0].Label != "Redis" || status[0].Present {
		t.Errorf("expected exactly one absent Redis status item, got %+v", status)
	}
}

func TestBuildStatusRedisPresent(t *testing.T) {

	info := &project.Info{}
	info.Database.Redis = &project.RedisInfo{Enabled: true}

	status := BuildStatus("redis", info)

	if len(status) != 1 || !status[0].Present {
		t.Errorf("expected Redis to be reported present, got %+v", status)
	}
}

func TestBuildStatusEveryTopicIsHandled(t *testing.T) {

	info := &project.Info{}

	for _, topic := range Topics {
		if status := BuildStatus(topic, info); len(status) == 0 {
			t.Errorf("topic %q has no status mapping in BuildStatus — Current Status would render empty", topic)
		}
	}
}

func TestBuildStatusUnknownTopic(t *testing.T) {

	if status := BuildStatus("not-a-real-topic", &project.Info{}); status != nil {
		t.Errorf("expected nil status for an unrecognized topic, got %+v", status)
	}
}

// 회귀 테스트: explain k8s처럼 아직 도입 안 한 topic을 explain할 때
// cmd/explain.go가 "아직 도입되지 않음" 배너를 보여줄지 이 함수로 판단한다.
func TestTopicPresent(t *testing.T) {

	if TopicPresent(nil) {
		t.Error("TopicPresent(nil) should be false")
	}

	if TopicPresent([]StatusItem{{Label: "Kubernetes manifests", Present: false}}) {
		t.Error("TopicPresent should be false when the primary status item is absent")
	}

	if !TopicPresent([]StatusItem{{Label: "Dockerfile", Present: true}, {Label: "Docker Compose", Present: false}}) {
		t.Error("TopicPresent should be true when the primary (first) status item is present, even if later items are absent")
	}
}
