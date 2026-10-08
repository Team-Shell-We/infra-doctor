package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeEntitiesParsesColumnsAndID(t *testing.T) {

	root := t.TempDir()

	src := `package com.example.demo;

import jakarta.persistence.*;

@Entity
public class Product {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long productId;

    @Column(nullable = false, unique = true)
    private String sku;

    @Column(name = "display_name")
    private String name;
}
`
	if err := os.WriteFile(filepath.Join(root, "Product.java"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	schema, err := AnalyzeEntities(root)
	if err != nil {
		t.Fatalf("AnalyzeEntities failed: %v", err)
	}

	if len(schema.Entities) != 1 {
		t.Fatalf("len(Entities) = %d, want 1", len(schema.Entities))
	}

	entity := schema.Entities[0]
	if entity.Name != "Product" {
		t.Errorf("Name = %q, want %q", entity.Name, "Product")
	}
	if len(entity.Columns) != 3 {
		t.Fatalf("len(Columns) = %d, want 3", len(entity.Columns))
	}

	byField := map[string]int{}
	for i, c := range entity.Columns {
		byField[c.Field] = i
	}

	id := entity.Columns[byField["productId"]]
	if !id.IsID {
		t.Errorf("productId.IsID = false, want true")
	}
	if id.Type != "Long" {
		t.Errorf("productId.Type = %q, want %q", id.Type, "Long")
	}

	sku := entity.Columns[byField["sku"]]
	if sku.Nullable {
		t.Errorf("sku.Nullable = true, want false")
	}
	if !sku.Unique {
		t.Errorf("sku.Unique = false, want true")
	}

	name := entity.Columns[byField["name"]]
	if name.Column != "display_name" {
		t.Errorf("name.Column = %q, want %q", name.Column, "display_name")
	}
}

func TestAnalyzeEntitiesParsesTableIndexes(t *testing.T) {

	root := t.TempDir()

	src := `package com.example.demo;

import jakarta.persistence.*;

@Entity
@Table(indexes = {
    @Index(name = "idx_order_customer", columnList = "customer_id, status")
})
public class Order {

    @Id
    private Long orderId;
}
`
	if err := os.WriteFile(filepath.Join(root, "Order.java"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	schema, err := AnalyzeEntities(root)
	if err != nil {
		t.Fatalf("AnalyzeEntities failed: %v", err)
	}

	if len(schema.Entities) != 1 {
		t.Fatalf("len(Entities) = %d, want 1", len(schema.Entities))
	}

	indexes := schema.Entities[0].Indexes
	if len(indexes) != 1 {
		t.Fatalf("len(Indexes) = %d, want 1", len(indexes))
	}
	if indexes[0].Name != "idx_order_customer" {
		t.Errorf("Indexes[0].Name = %q, want %q", indexes[0].Name, "idx_order_customer")
	}
	if len(indexes[0].Columns) != 2 || indexes[0].Columns[0] != "customer_id" || indexes[0].Columns[1] != "status" {
		t.Errorf("Indexes[0].Columns = %v, want [customer_id status]", indexes[0].Columns)
	}
}

func TestAnalyzeEntitiesParsesRelationships(t *testing.T) {

	root := t.TempDir()

	src := `package com.example.demo;

import jakarta.persistence.*;
import java.util.List;

@Entity
public class Review {

    @Id
    private Long reviewId;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "facility_id", nullable = false)
    private Facility facility;

    @OneToMany(mappedBy = "review", cascade = CascadeType.ALL, orphanRemoval = true, fetch = FetchType.LAZY)
    private List<Comment> comments;
}
`
	if err := os.WriteFile(filepath.Join(root, "Review.java"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	schema, err := AnalyzeEntities(root)
	if err != nil {
		t.Fatalf("AnalyzeEntities failed: %v", err)
	}

	if len(schema.Relationships) != 2 {
		t.Fatalf("len(Relationships) = %d, want 2", len(schema.Relationships))
	}

	var foundManyToOne, foundOneToMany bool
	for _, rel := range schema.Relationships {
		switch rel.Field {
		case "facility":
			foundManyToOne = true
			if rel.Type != "ManyToOne" {
				t.Errorf("facility.Type = %q, want ManyToOne", rel.Type)
			}
			if rel.JoinColumn != "facility_id" {
				t.Errorf("facility.JoinColumn = %q, want facility_id", rel.JoinColumn)
			}
			if rel.Fetch != "LAZY" {
				t.Errorf("facility.Fetch = %q, want LAZY", rel.Fetch)
			}
			if rel.To != "Facility" {
				t.Errorf("facility.To = %q, want Facility", rel.To)
			}
		case "comments":
			foundOneToMany = true
			if rel.Type != "OneToMany" {
				t.Errorf("comments.Type = %q, want OneToMany", rel.Type)
			}
			if rel.MappedBy != "review" {
				t.Errorf("comments.MappedBy = %q, want review", rel.MappedBy)
			}
			if rel.Cascade != "ALL" {
				t.Errorf("comments.Cascade = %q, want ALL", rel.Cascade)
			}
			if rel.To != "Comment" {
				t.Errorf("comments.To = %q, want Comment (List<> wrapper stripped)", rel.To)
			}
		}
	}

	if !foundManyToOne {
		t.Error("no relationship found for field 'facility'")
	}
	if !foundOneToMany {
		t.Error("no relationship found for field 'comments'")
	}
}

func TestAnalyzeEntitiesMergesMappedSuperclassFields(t *testing.T) {

	root := t.TempDir()

	base := `package com.example.demo;

import jakarta.persistence.MappedSuperclass;
import java.time.LocalDateTime;

@MappedSuperclass
public class BaseEntity {

    @Column(nullable = false)
    private LocalDateTime createdAt;
}
`
	child := `package com.example.demo;

import jakarta.persistence.*;

@Entity
public class Account extends BaseEntity {

    @Id
    private Long accountId;
}
`
	if err := os.WriteFile(filepath.Join(root, "BaseEntity.java"), []byte(base), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "Account.java"), []byte(child), 0644); err != nil {
		t.Fatal(err)
	}

	schema, err := AnalyzeEntities(root)
	if err != nil {
		t.Fatalf("AnalyzeEntities failed: %v", err)
	}

	if len(schema.Entities) != 1 {
		t.Fatalf("len(Entities) = %d, want 1 (BaseEntity must not be treated as an entity)", len(schema.Entities))
	}

	account := schema.Entities[0]
	if len(account.Columns) != 2 {
		t.Fatalf("len(Account.Columns) = %d, want 2 (own accountId + inherited createdAt)", len(account.Columns))
	}

	var foundInherited bool
	for _, c := range account.Columns {
		if c.Field == "createdAt" {
			foundInherited = true
			if !c.Inherited {
				t.Error("createdAt.Inherited = false, want true")
			}
		}
	}
	if !foundInherited {
		t.Error("createdAt was not merged in from BaseEntity")
	}
}

func TestAnalyzeEntitiesIgnoresPlainClasses(t *testing.T) {

	root := t.TempDir()

	src := `package com.example.demo;

public class RequestDto {
    private String name;
}
`
	if err := os.WriteFile(filepath.Join(root, "RequestDto.java"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	schema, err := AnalyzeEntities(root)
	if err != nil {
		t.Fatalf("AnalyzeEntities failed: %v", err)
	}

	if len(schema.Entities) != 0 {
		t.Errorf("len(Entities) = %d, want 0", len(schema.Entities))
	}
}
