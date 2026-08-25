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
