package generate

type Result struct {
	Target Target
	DryRun bool

	Detections []string
	Warnings   []string

	Planned     []string
	Created     []string
	Overwritten []string
	Skipped     []string
}

// StatusOf : path 하나의 생성 결과 상태("created"/"planned"/"skipped"/"overwritten")
func (r Result) StatusOf(path string) string {

	for _, skipped := range r.Skipped {
		if path == skipped {
			return "skipped"
		}
	}

	for _, overwritten := range r.Overwritten {
		if path == overwritten {
			return "overwritten"
		}
	}

	if r.DryRun {
		return "planned"
	}

	return "created"
}
