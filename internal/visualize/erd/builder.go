package erd

import "github.com/Team-Shell-We/infra-doctor/internal/project"

// Build : project.SchemaInfo를 ERD Diagram으로 변환한다. 양쪽에서 각각 선언된 관계
// (예: FacilityReview.comments @OneToMany(mappedBy="facilityReview")와 그 짝이 되는
// 선언)는 하나의 Relation으로 병합한다 — analyzer는 선언된 사실을 그대로 모으기만 하고,
// 중복 제거 같은 표현상의 판단은 시각화 쪽 책임이라는 원칙을 따른다.
func Build(info project.Info) Diagram {

	diagram := Diagram{Title: "Entity-Relationship Diagram"}

	for _, entity := range info.Schema.Entities {
		diagram.Entities = append(diagram.Entities, buildEntity(entity, info.Schema.Relationships))
	}

	diagram.Relations = mergeRelationships(info.Schema.Relationships)

	return diagram
}

func buildEntity(entity project.EntityInfo, relationships []project.RelationshipInfo) Entity {

	result := Entity{Name: entity.Name}

	for _, column := range entity.Columns {
		result.Columns = append(result.Columns, Column{
			Name:   column.Column,
			Type:   column.Type,
			PK:     column.IsID,
			Unique: column.Unique,
		})
	}

	// 관계 필드 자체는 entity.Columns에 없으므로(별도 Relationships로 관리됨),
	// 소유 측(@JoinColumn 보유)만 FK 컬럼으로 합성해 보여준다.
	for _, rel := range relationships {
		if rel.From != entity.Name || rel.JoinColumn == "" {
			continue
		}
		result.Columns = append(result.Columns, Column{
			Name: rel.JoinColumn,
			Type: rel.To,
			FK:   true,
		})
	}

	return result
}

// mergeRelationships : 양방향으로 짝이 맞는 관계 선언 둘을 하나의 Relation으로 합친다.
// 짝을 못 찾으면(단방향 선언) 그대로 하나만 만든다.
func mergeRelationships(relationships []project.RelationshipInfo) []Relation {

	consumed := make([]bool, len(relationships))
	var result []Relation

	for i, rel := range relationships {

		if consumed[i] {
			continue
		}
		consumed[i] = true

		pairIdx := -1
		for j := i + 1; j < len(relationships); j++ {
			if consumed[j] {
				continue
			}
			if isPair(rel, relationships[j]) {
				pairIdx = j
				break
			}
		}

		if pairIdx == -1 {
			result = append(result, Relation{
				From:        rel.From,
				To:          rel.To,
				Label:       rel.Field,
				Cardinality: cardinality(rel.Type),
			})
			continue
		}

		consumed[pairIdx] = true
		result = append(result, mergedRelation(rel, relationships[pairIdx]))
	}

	return result
}

// isPair : 두 관계 선언이 같은 물리적 관계의 양쪽인지 — 서로 엔티티가 맞물리고,
// 한쪽의 mappedBy가 다른 쪽의 필드명을 가리키면 짝으로 본다
func isPair(a, b project.RelationshipInfo) bool {
	if a.From != b.To || a.To != b.From {
		return false
	}
	return a.MappedBy == b.Field || b.MappedBy == a.Field
}

// mergedRelation : 짝이 맞는 두 선언 중 소유 측(mappedBy가 없는 쪽, 보통 @JoinColumn을 가짐)을
// 기준으로 Relation을 만든다. 소유 측의 Type이 카디널리티를 그대로 결정한다.
func mergedRelation(a, b project.RelationshipInfo) Relation {

	owner := a
	if a.MappedBy != "" {
		owner = b
	}

	return Relation{
		From:        owner.From,
		To:          owner.To,
		Label:       owner.Field,
		Cardinality: cardinality(owner.Type),
	}
}

func cardinality(relType string) string {
	switch relType {
	case "OneToOne":
		return "one-to-one"
	case "OneToMany":
		return "one-to-many"
	case "ManyToMany":
		return "many-to-many"
	default: // ManyToOne
		return "many-to-one"
	}
}
