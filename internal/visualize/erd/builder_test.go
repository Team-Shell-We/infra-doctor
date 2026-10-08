package erd

import (
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func TestBuildMergesBidirectionalRelationship(t *testing.T) {

	info := project.Info{
		Schema: project.SchemaInfo{
			Entities: []project.EntityInfo{
				{Name: "Facility"},
				{Name: "FacilityReview"},
			},
			Relationships: []project.RelationshipInfo{
				{From: "FacilityReview", To: "Facility", Type: "ManyToOne", Field: "facility", JoinColumn: "facility_id", Fetch: "LAZY"},
				{From: "Facility", To: "FacilityReview", Type: "OneToMany", Field: "reviews", MappedBy: "facility"},
			},
		},
	}

	diagram := Build(info)

	if len(diagram.Relations) != 1 {
		t.Fatalf("len(Relations) = %d, want 1 (bidirectional pair must merge into one)", len(diagram.Relations))
	}

	rel := diagram.Relations[0]
	if rel.From != "FacilityReview" || rel.To != "Facility" {
		t.Errorf("merged relation = %s -> %s, want FacilityReview -> Facility (anchored on owning/JoinColumn side)", rel.From, rel.To)
	}
	if rel.Cardinality != "many-to-one" {
		t.Errorf("Cardinality = %q, want many-to-one", rel.Cardinality)
	}
	if rel.Label != "facility" {
		t.Errorf("Label = %q, want facility", rel.Label)
	}
}

func TestBuildKeepsUnidirectionalRelationshipAlone(t *testing.T) {

	info := project.Info{
		Schema: project.SchemaInfo{
			Entities: []project.EntityInfo{
				{Name: "Order"},
				{Name: "Customer"},
			},
			Relationships: []project.RelationshipInfo{
				{From: "Order", To: "Customer", Type: "ManyToOne", Field: "customer", JoinColumn: "customer_id"},
			},
		},
	}

	diagram := Build(info)

	if len(diagram.Relations) != 1 {
		t.Fatalf("len(Relations) = %d, want 1", len(diagram.Relations))
	}
	if diagram.Relations[0].Cardinality != "many-to-one" {
		t.Errorf("Cardinality = %q, want many-to-one", diagram.Relations[0].Cardinality)
	}
}

func TestBuildAddsSyntheticForeignKeyColumn(t *testing.T) {

	info := project.Info{
		Schema: project.SchemaInfo{
			Entities: []project.EntityInfo{
				{
					Name: "Order",
					Columns: []project.ColumnInfo{
						{Field: "orderId", Column: "orderId", Type: "Long", IsID: true},
					},
				},
				{Name: "Customer"},
			},
			Relationships: []project.RelationshipInfo{
				{From: "Order", To: "Customer", Type: "ManyToOne", Field: "customer", JoinColumn: "customer_id"},
			},
		},
	}

	diagram := Build(info)

	var order Entity
	for _, e := range diagram.Entities {
		if e.Name == "Order" {
			order = e
		}
	}

	if len(order.Columns) != 2 {
		t.Fatalf("len(Order.Columns) = %d, want 2 (own PK + synthesized FK)", len(order.Columns))
	}

	var foundFK bool
	for _, c := range order.Columns {
		if c.Name == "customer_id" {
			foundFK = true
			if !c.FK {
				t.Error("customer_id.FK = false, want true")
			}
		}
	}
	if !foundFK {
		t.Error("customer_id FK column was not synthesized from the relationship's JoinColumn")
	}
}
