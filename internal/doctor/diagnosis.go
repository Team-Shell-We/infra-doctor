package doctor

type Level string

const (
	Info     Level = "INFO"
	Warning  Level = "WARNING"
	Critical Level = "CRITICAL"
)

type Category string

const (
	Infrastructure Category = "Infrastructure"
	CICD           Category = "CICD"
	Database       Category = "Database"
	Monitoring     Category = "Monitoring"
	Security       Category = "Security"
)

type Diagnosis struct {
	Category Category `json:"category"`
	Level    Level    `json:"level"`

	// 감점
	ScoreImpact int `json:"scoreImpact"`

	Title   string `json:"title"`
	Message string `json:"message"`

	Reason string `json:"reason"`
	Fix    string `json:"fix"`
}

type Result struct {
	Score     int         `json:"score"`
	Diagnoses []Diagnosis `json:"diagnoses"`
	Checks    []Check     `json:"checks"`
}