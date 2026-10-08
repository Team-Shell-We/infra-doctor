package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

type DatabaseRule struct{}

func (r DatabaseRule) Check(info *project.Info) []Diagnosis {

	var diagnoses []Diagnosis

	registry, err := LoadRules()
	if err != nil {
		return diagnoses
	}

	checks := []struct {
		id      string
		missing bool
	}{
		{"unindexed_foreign_key", hasUnindexedForeignKey(info)},
		{"eager_fetch_risk", hasEagerFetchRisk(info)},
		{"cascade_all_on_many_to_one", hasCascadeAllOnOwningSide(info)},
		{"no_migration_tool", isRelationalDatabase(info) && !info.Dependencies.Migration.Enabled},
	}

	for _, c := range checks {

		if !c.missing {
			continue
		}

		if rule, err := registry.DatabaseRule(c.id); err == nil {
			diagnoses = append(diagnoses, rule)
		}
	}

	return diagnoses
}

// hasUnindexedForeignKey : 소유 측(@JoinColumn 보유) 관계 중, 그 컬럼이 자신이 속한
// 엔티티의 @Table(indexes=...)에 하나도 없는 경우가 있는지. PostgreSQL 등 일부 DB는
// FK 컬럼을 자동으로 인덱싱하지 않으므로 조인/삭제 성능에 실제로 영향을 준다.
func hasUnindexedForeignKey(info *project.Info) bool {

	indexedColumns := make(map[string]map[string]bool) // entity -> column -> true

	for _, entity := range info.Schema.Entities {
		columns := make(map[string]bool)
		for _, index := range entity.Indexes {
			for _, column := range index.Columns {
				columns[column] = true
			}
		}
		indexedColumns[entity.Name] = columns
	}

	for _, rel := range info.Schema.Relationships {
		if rel.JoinColumn == "" {
			continue
		}
		if !indexedColumns[rel.From][rel.JoinColumn] {
			return true
		}
	}

	return false
}

// hasEagerFetchRisk : @ManyToOne/@OneToOne에 fetch=LAZY를 명시하지 않은 경우가 있는지.
// 이 둘은 JPA 기본값이 EAGER라 명시하지 않으면 매번 연관 엔티티를 즉시 로딩한다
// (@OneToMany/@ManyToMany는 기본값이 이미 LAZY라 대상이 아니다).
func hasEagerFetchRisk(info *project.Info) bool {

	for _, rel := range info.Schema.Relationships {
		if (rel.Type == "ManyToOne" || rel.Type == "OneToOne") && rel.Fetch != "LAZY" {
			return true
		}
	}

	return false
}

// hasCascadeAllOnOwningSide : @ManyToOne/@OneToOne(자식 -> 부모 방향)에
// CascadeType.ALL이 붙어 있는 경우가 있는지. 이 방향의 ALL은 자식 저장/삭제가
// 의도치 않게 부모까지 전파될 수 있어 위험하다 — 반대로 @OneToMany/@ManyToMany
// (부모가 자식을 관리하는 방향)의 ALL은 일반적인 패턴이라 대상이 아니다.
func hasCascadeAllOnOwningSide(info *project.Info) bool {

	for _, rel := range info.Schema.Relationships {
		if (rel.Type == "ManyToOne" || rel.Type == "OneToOne") && rel.Cascade == "ALL" {
			return true
		}
	}

	return false
}

func isRelationalDatabase(info *project.Info) bool {
	switch info.Database.Primary.Type {
	case "PostgreSQL", "MySQL", "MariaDB":
		return true
	default:
		return false
	}
}
