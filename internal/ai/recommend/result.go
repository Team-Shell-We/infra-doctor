package recommend

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Result: AI가 실제로 생성하는 유일한 값 — 배포 전략 추천 이유를 자연어로
// 풀어쓴 문장 목록. Recommended/KubernetesFit/NextSteps는 전부 decision.go에서
// 결정론적으로 계산되므로 여기 포함되지 않는다.
type Result struct {
	Reasons []string `json:"reasons"`
}

// Parse는 모델 응답을 디코딩하고 검증한다. reasons가 비어있으면 렌더링할
// 문구가 없다는 뜻이므로 실패로 처리한다.
func Parse(raw string) (*Result, error) {

	var result Result

	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &result); err != nil {
		return nil, fmt.Errorf("recommend: failed to parse AI response as JSON: %w", err)
	}

	if len(result.Reasons) == 0 {
		return nil, fmt.Errorf("recommend: AI response is missing required field: reasons")
	}

	return &result, nil
}
