package visualize

import (
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

// 회귀 테스트: DB 미감지 시 Primary.Type이 "Unknown"으로 채워지는데,
// 이게 실제 데이터베이스처럼 다이어그램에 노드로 나오면 안 됨.
func TestBuildOmitsUnknownDatabaseNode(t *testing.T) {

	info := project.Info{}
	info.Database.Primary.Type = "Unknown"

	diagram := Build(info)

	for _, node := range diagram.Nodes {
		if node.Kind == Database {
			t.Errorf("expected no Database node when Primary.Type is Unknown, got %+v", node)
		}
	}
}

func TestBuildIncludesRealDatabaseNode(t *testing.T) {

	info := project.Info{}
	info.Database.Primary.Type = "PostgreSQL"

	diagram := Build(info)

	found := false
	for _, node := range diagram.Nodes {
		if node.Kind == Database && node.Label == "PostgreSQL" {
			found = true
		}
	}
	if !found {
		t.Error("expected a PostgreSQL Database node")
	}
}
