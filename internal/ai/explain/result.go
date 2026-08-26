package explain

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Result : `explain` 응답이 채워야 하는 고정된 3섹션 구조
// 필드명/태그와 프롬프트 템플릿의 필드명은 항상 일치해야 함.
type Result struct {
	CurrentProject []string `json:"current_project"`
	BuildFlow      []string `json:"build_flow"`
	WhyTopic       []string `json:"why_topic"`
}

// Parse : 모델 응답을 디코딩하고 검증한다.
// 모든 섹션이 비어 있지 않아야 하며, 하나라도 비어 있으면 빈 박스 대신 렌더링 실패로 처리한다.
func Parse(raw string) (*Result, error) {

	var result Result

	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &result); err != nil {
		return nil, fmt.Errorf("explain: failed to parse AI response as JSON: %w", err)
	}

	var missing []string

	if len(result.CurrentProject) == 0 {
		missing = append(missing, "current_project")
	}

	if len(result.BuildFlow) == 0 {
		missing = append(missing, "build_flow")
	}

	if len(result.WhyTopic) == 0 {
		missing = append(missing, "why_topic")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("explain: AI response is missing required field(s): %s", strings.Join(missing, ", "))
	}

	return &result, nil
}
