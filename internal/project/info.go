package project

type Info struct {
	Framework FrameworkInfo
}

// FrameworkInfo : 프로젝트에서 사용하는 프레임워크 정보
type FrameworkInfo struct {
	BuildTool     string
	SpringBoot    bool
	SpringVersion string
	JavaVersion   string
}