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
	CICD           Category = "CI/CD"
	Database       Category = "Database"
	Monitoring     Category = "Monitoring"
	Security       Category = "Security"
)

type Diagnosis struct {
	Category Category
	Level    Level

	// 감점
	ScoreImpact int

	Title   string
	Message string

	Reason string
	Fix    string
}

type Result struct {
	Score      int
	Diagnoses  []Diagnosis
}