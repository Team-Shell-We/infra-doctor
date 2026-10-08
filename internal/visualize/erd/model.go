package erd

// Format : ERD 렌더링 출력 형식. internal/visualize.Format과 같은 세 가지 값을 쓰지만,
// Diagram 모델 자체가 달라 독립된 타입으로 둔다.
type Format string

const (
	ASCII    Format = "ascii"
	Mermaid  Format = "mermaid"
	Markdown Format = "markdown"
)

// Diagram : 엔티티-관계 다이어그램 전체
type Diagram struct {
	Title     string
	Entities  []Entity
	Relations []Relation
}

// Entity : ERD에서 박스 하나로 그려지는 엔티티(JPA @Entity 클래스) 하나
type Entity struct {
	Name    string
	Columns []Column
}

// Column : 엔티티 박스 안에 나열되는 컬럼 한 줄
type Column struct {
	Name   string
	Type   string
	PK     bool
	FK     bool
	Unique bool
}

// Relation : 엔티티 두 개를 잇는 관계 한 줄. 양방향으로 선언된 관계(@OneToMany의
// mappedBy와 짝이 되는 @ManyToOne)는 Build 단계에서 하나로 병합된다.
type Relation struct {
	From, To    string
	Label       string
	Cardinality string // "one-to-one" | "one-to-many" | "many-to-one" | "many-to-many"
}

// keyLabel : 세 renderer(ascii/mermaid/markdown)가 공통으로 쓰는 PK/FK/UK 표시
func keyLabel(c Column) string {
	switch {
	case c.PK:
		return "PK"
	case c.FK:
		return "FK"
	case c.Unique:
		return "UK"
	default:
		return ""
	}
}
