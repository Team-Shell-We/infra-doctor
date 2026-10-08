package analyzer

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

var (
	entityAnnoFileRegex           = regexp.MustCompile(`@Entity\b`)
	mappedSuperclassAnnoFileRegex = regexp.MustCompile(`@MappedSuperclass\b`)
	classDeclRegex                = regexp.MustCompile(`\bclass\s+(\w+)(?:\s+extends\s+(\w+))?`)

	// fieldBlockRegex : 필드 바로 위에 붙은 어노테이션 묶음(1개 이상, 줄바꿈 포함) + 필드 선언을 한 번에 캡처한다.
	// 어노테이션이 하나도 없는 필드는 DB 분석 대상이 아니므로 의도적으로 매칭하지 않는다.
	fieldBlockRegex = regexp.MustCompile(`(?s)((?:@\w+(?:\([^)]*\))?\s*)+)(?:private|protected|public)\s+([\w.<>\[\],\s]+?)\s+(\w+)\s*(?:=[^;]*)?;`)

	idAnnoRegex         = regexp.MustCompile(`@Id\b`)
	columnAnnoRegex     = regexp.MustCompile(`@Column\b(?:\(([^)]*)\))?`)
	joinColumnAnnoRegex = regexp.MustCompile(`@JoinColumn\b(?:\(([^)]*)\))?`)
	oneToOneAnnoRegex   = regexp.MustCompile(`@OneToOne\b(?:\(([^)]*)\))?`)
	oneToManyAnnoRegex  = regexp.MustCompile(`@OneToMany\b(?:\(([^)]*)\))?`)
	manyToOneAnnoRegex  = regexp.MustCompile(`@ManyToOne\b(?:\(([^)]*)\))?`)
	manyToManyAnnoRegex = regexp.MustCompile(`@ManyToMany\b(?:\(([^)]*)\))?`)

	nameAttrRegex       = regexp.MustCompile(`\bname\s*=\s*"([^"]+)"`)
	nullableAttrRegex   = regexp.MustCompile(`\bnullable\s*=\s*(true|false)`)
	uniqueAttrRegex     = regexp.MustCompile(`\bunique\s*=\s*(true|false)`)
	mappedByAttrRegex   = regexp.MustCompile(`\bmappedBy\s*=\s*"([^"]+)"`)
	fetchAttrRegex      = regexp.MustCompile(`\bfetch\s*=\s*FetchType\.(\w+)`)
	cascadeValueRegex   = regexp.MustCompile(`CascadeType\.(\w+)`)
	columnListAttrRegex = regexp.MustCompile(`\bcolumnList\s*=\s*"([^"]+)"`)
	indexEntryRegex     = regexp.MustCompile(`@Index\s*\(([^)]*)\)`)
	collectionWrapRegex = regexp.MustCompile(`^(?:List|Set|Collection)\s*<\s*(\w+)\s*>$`)
)

// javaFile : AnalyzeEntities 1차 스캔에서 모은, 파일 하나에 대한 원시 정보
type javaFile struct {
	path      string
	content   string
	className string
	extends   string
	isEntity  bool
	isMapped  bool
}

// AnalyzeEntities : .java 소스에서 JPA 엔티티/관계 어노테이션을 파싱해 project.SchemaInfo를 만든다.
// AnalyzeAPI(api.go)와 동일한 filepath.Walk + shouldSkipDir 패턴을 쓰며, AST가 아닌 정규식 기반이라
// 중첩된 메서드 호출이 어노테이션 인자로 들어가는 경우나 한 파일에 top-level class가 여러 개인
// 경우는 다루지 않는다.
func AnalyzeEntities(root string) (*project.SchemaInfo, error) {

	var files []javaFile

	err := filepath.Walk(root, func(path string, fileInfo os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			if shouldSkipDir(fileInfo.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".java" {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		content := string(data)

		match := classDeclRegex.FindStringSubmatch(content)
		if match == nil {
			return nil
		}

		files = append(files, javaFile{
			path:      path,
			content:   content,
			className: match[1],
			extends:   match[2],
			isEntity:  entityAnnoFileRegex.MatchString(content),
			isMapped:  mappedSuperclassAnnoFileRegex.MatchString(content),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	superclassFields := make(map[string][]project.ColumnInfo)
	for _, file := range files {
		if file.isMapped {
			superclassFields[file.className] = parseColumns(file.content)
		}
	}

	schema := &project.SchemaInfo{}

	for _, file := range files {
		if !file.isEntity {
			continue
		}

		relPath, relErr := filepath.Rel(root, file.path)
		if relErr != nil {
			relPath = file.path
		}

		entity := project.EntityInfo{
			Name:    file.className,
			File:    filepath.ToSlash(relPath),
			Columns: parseColumns(file.content),
		}

		if inherited, ok := superclassFields[file.extends]; ok {
			for i := range inherited {
				inherited[i].Inherited = true
			}
			entity.Columns = append(entity.Columns, inherited...)
		}

		entity.Table, entity.Indexes = parseTable(file.content)

		schema.Entities = append(schema.Entities, entity)
		schema.Relationships = append(schema.Relationships, parseRelationships(file.className, file.content)...)
	}

	return schema, nil
}

// parseColumns : 파일 본문에서 @Id/@Column이 붙은 필드만 ColumnInfo로 변환한다
func parseColumns(content string) []project.ColumnInfo {

	var columns []project.ColumnInfo

	for _, block := range fieldBlockRegex.FindAllStringSubmatch(content, -1) {

		annotations, fieldType, field := block[1], strings.TrimSpace(block[2]), block[3]

		isID := idAnnoRegex.MatchString(annotations)
		columnMatch := columnAnnoRegex.FindStringSubmatch(annotations)

		if !isID && columnMatch == nil {
			continue
		}

		column := project.ColumnInfo{
			Field:    field,
			Column:   field,
			Type:     fieldType,
			Nullable: true,
			IsID:     isID,
		}

		if columnMatch != nil {
			attrs := columnMatch[1]
			if name := nameAttrRegex.FindStringSubmatch(attrs); name != nil {
				column.Column = name[1]
			}
			if nullable := nullableAttrRegex.FindStringSubmatch(attrs); nullable != nil {
				column.Nullable = nullable[1] == "true"
			}
			if unique := uniqueAttrRegex.FindStringSubmatch(attrs); unique != nil {
				column.Unique = unique[1] == "true"
			}
		}

		columns = append(columns, column)
	}

	return columns
}

// parseTable : @Table(name=..., indexes={@Index(columnList=...)})를 파싱한다.
// indexes는 중첩 괄호를 포함해 regexp 하나로 뜰 수 없으므로 괄호 카운팅으로 블록을 먼저 잘라낸다.
func parseTable(content string) (string, []project.IndexInfo) {

	idx := strings.Index(content, "@Table")
	if idx == -1 {
		return "", nil
	}

	openParen := strings.IndexByte(content[idx:], '(')
	if openParen == -1 {
		return "", nil
	}

	block := extractBalanced(content, idx+openParen)

	var tableName string
	if name := nameAttrRegex.FindStringSubmatch(block); name != nil {
		tableName = name[1]
	}

	var indexes []project.IndexInfo
	for _, entry := range indexEntryRegex.FindAllStringSubmatch(block, -1) {
		columnList := columnListAttrRegex.FindStringSubmatch(entry[1])
		if columnList == nil {
			continue
		}

		index := project.IndexInfo{}
		if name := nameAttrRegex.FindStringSubmatch(entry[1]); name != nil {
			index.Name = name[1]
		}
		for _, column := range strings.Split(columnList[1], ",") {
			index.Columns = append(index.Columns, strings.TrimSpace(column))
		}
		indexes = append(indexes, index)
	}

	return tableName, indexes
}

// parseRelationships : @OneToOne/@OneToMany/@ManyToOne/@ManyToMany가 붙은 필드마다 RelationshipInfo를 만든다.
// 양쪽에 다 선언돼 있어도 선언된 그대로 각각 저장한다 — 중복 제거는 하지 않는다(시각화 단계의 책임).
func parseRelationships(entityName, content string) []project.RelationshipInfo {

	var relationships []project.RelationshipInfo

	for _, block := range fieldBlockRegex.FindAllStringSubmatch(content, -1) {

		annotations, fieldType, field := block[1], strings.TrimSpace(block[2]), block[3]

		relType, attrs := "", ""
		switch {
		case oneToOneAnnoRegex.MatchString(annotations):
			relType, attrs = "OneToOne", oneToOneAnnoRegex.FindStringSubmatch(annotations)[1]
		case oneToManyAnnoRegex.MatchString(annotations):
			relType, attrs = "OneToMany", oneToManyAnnoRegex.FindStringSubmatch(annotations)[1]
		case manyToOneAnnoRegex.MatchString(annotations):
			relType, attrs = "ManyToOne", manyToOneAnnoRegex.FindStringSubmatch(annotations)[1]
		case manyToManyAnnoRegex.MatchString(annotations):
			relType, attrs = "ManyToMany", manyToManyAnnoRegex.FindStringSubmatch(annotations)[1]
		default:
			continue
		}

		relationship := project.RelationshipInfo{
			From:  entityName,
			To:    relationTargetType(fieldType),
			Type:  relType,
			Field: field,
		}

		if fetch := fetchAttrRegex.FindStringSubmatch(attrs); fetch != nil {
			relationship.Fetch = fetch[1]
		}

		if mappedBy := mappedByAttrRegex.FindStringSubmatch(attrs); mappedBy != nil {
			relationship.MappedBy = mappedBy[1]
		}

		cascades := cascadeValueRegex.FindAllStringSubmatch(attrs, -1)
		for _, cascade := range cascades {
			if cascade[1] == "ALL" {
				relationship.Cascade = "ALL"
				break
			}
			if relationship.Cascade == "" {
				relationship.Cascade = cascade[1]
			}
		}

		if joinColumn := joinColumnAnnoRegex.FindStringSubmatch(annotations); joinColumn != nil {
			if name := nameAttrRegex.FindStringSubmatch(joinColumn[1]); name != nil {
				relationship.JoinColumn = name[1]
			}
		}

		relationships = append(relationships, relationship)
	}

	return relationships
}

// relationTargetType : List<Foo>/Set<Foo>/Collection<Foo> 래퍼를 벗겨 대상 타입명만 남긴다
func relationTargetType(javaType string) string {
	if match := collectionWrapRegex.FindStringSubmatch(javaType); match != nil {
		return match[1]
	}
	return javaType
}

// extractBalanced : text[openParenIdx]가 '('일 때, 그와 짝이 맞는 ')'까지의 내부 문자열을 반환한다.
// @Table(indexes = {@Index(...)})처럼 중첩 괄호가 있는 어노테이션 인자를 regexp 하나로 뜰 수 없어 쓴다.
func extractBalanced(text string, openParenIdx int) string {

	depth := 0
	start := -1

	for i := openParenIdx; i < len(text); i++ {
		switch text[i] {
		case '(':
			if depth == 0 {
				start = i + 1
			}
			depth++
		case ')':
			depth--
			if depth == 0 {
				return text[start:i]
			}
		}
	}

	return ""
}
