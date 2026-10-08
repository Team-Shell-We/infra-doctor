package doctor

import (
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func baseRelationalInfo() *project.Info {
	return &project.Info{
		Database: project.DatabaseInfo{Primary: project.Database{Type: "PostgreSQL"}},
		Dependencies: project.DependencyInfo{
			Migration: project.MigrationInfo{Enabled: true, Tool: "Flyway"},
		},
	}
}

func TestDatabaseRuleFlagsUnindexedForeignKey(t *testing.T) {

	info := baseRelationalInfo()
	info.Schema = project.SchemaInfo{
		Entities: []project.EntityInfo{{Name: "Review"}}, // 인덱스 없음
		Relationships: []project.RelationshipInfo{
			{From: "Review", To: "Facility", Type: "ManyToOne", JoinColumn: "facility_id", Fetch: "LAZY"},
		},
	}

	diagnoses := DatabaseRule{}.Check(info)

	if !containsID(diagnoses, "unindexed_foreign_key") {
		t.Error("expected unindexed_foreign_key diagnosis, got none")
	}
}

func TestDatabaseRuleSkipsIndexedForeignKey(t *testing.T) {

	info := baseRelationalInfo()
	info.Schema = project.SchemaInfo{
		Entities: []project.EntityInfo{{
			Name:    "Review",
			Indexes: []project.IndexInfo{{Columns: []string{"facility_id"}}},
		}},
		Relationships: []project.RelationshipInfo{
			{From: "Review", To: "Facility", Type: "ManyToOne", JoinColumn: "facility_id", Fetch: "LAZY"},
		},
	}

	diagnoses := DatabaseRule{}.Check(info)

	if containsID(diagnoses, "unindexed_foreign_key") {
		t.Error("did not expect unindexed_foreign_key diagnosis when the column is indexed")
	}
}

func TestDatabaseRuleFlagsEagerFetch(t *testing.T) {

	info := baseRelationalInfo()
	info.Schema = project.SchemaInfo{
		Relationships: []project.RelationshipInfo{
			{From: "Review", To: "Facility", Type: "ManyToOne"}, // Fetch 미지정 -> 기본 EAGER
		},
	}

	diagnoses := DatabaseRule{}.Check(info)

	if !containsID(diagnoses, "eager_fetch_risk") {
		t.Error("expected eager_fetch_risk diagnosis, got none")
	}
}

func TestDatabaseRuleIgnoresOneToManyFetchDefault(t *testing.T) {

	info := baseRelationalInfo()
	info.Schema = project.SchemaInfo{
		Relationships: []project.RelationshipInfo{
			{From: "Facility", To: "Review", Type: "OneToMany"}, // 기본값이 이미 LAZY
		},
	}

	diagnoses := DatabaseRule{}.Check(info)

	if containsID(diagnoses, "eager_fetch_risk") {
		t.Error("OneToMany without explicit fetch should not trigger eager_fetch_risk (default is already LAZY)")
	}
}

func TestDatabaseRuleFlagsCascadeAllOnManyToOne(t *testing.T) {

	info := baseRelationalInfo()
	info.Schema = project.SchemaInfo{
		Relationships: []project.RelationshipInfo{
			{From: "Review", To: "Facility", Type: "ManyToOne", Fetch: "LAZY", Cascade: "ALL"},
		},
	}

	diagnoses := DatabaseRule{}.Check(info)

	if !containsID(diagnoses, "cascade_all_on_many_to_one") {
		t.Error("expected cascade_all_on_many_to_one diagnosis, got none")
	}
}

func TestDatabaseRuleIgnoresCascadeAllOnOneToMany(t *testing.T) {

	info := baseRelationalInfo()
	info.Schema = project.SchemaInfo{
		Relationships: []project.RelationshipInfo{
			{From: "Facility", To: "Review", Type: "OneToMany", Cascade: "ALL"}, // 부모->자식 방향은 정상 패턴
		},
	}

	diagnoses := DatabaseRule{}.Check(info)

	if containsID(diagnoses, "cascade_all_on_many_to_one") {
		t.Error("CascadeType.ALL on the OneToMany (parent) side should not be flagged")
	}
}

func TestDatabaseRuleFlagsMissingMigrationTool(t *testing.T) {

	info := baseRelationalInfo()
	info.Dependencies.Migration = project.MigrationInfo{}

	diagnoses := DatabaseRule{}.Check(info)

	if !containsID(diagnoses, "no_migration_tool") {
		t.Error("expected no_migration_tool diagnosis, got none")
	}
}

func TestDatabaseRuleSkipsMigrationCheckForNonRelationalDatabase(t *testing.T) {

	info := baseRelationalInfo()
	info.Database.Primary.Type = "Unknown"
	info.Dependencies.Migration = project.MigrationInfo{}

	diagnoses := DatabaseRule{}.Check(info)

	if containsID(diagnoses, "no_migration_tool") {
		t.Error("no_migration_tool should not fire when there is no relational database in use")
	}
}

func containsID(diagnoses []Diagnosis, id string) bool {
	registry, err := LoadRules()
	if err != nil {
		return false
	}
	want, err := registry.DatabaseRule(id)
	if err != nil {
		return false
	}
	for _, d := range diagnoses {
		if d.Title == want.Title {
			return true
		}
	}
	return false
}
