package project

type Info struct {
	Framework FrameworkInfo
	Profiles []ProfileInfo
}

// FrameworkInfo : 프로젝트에서 사용하는 프레임워크 정보
type FrameworkInfo struct {
	BuildTool     string
	SpringBoot    bool
	SpringVersion string
	JavaVersion   string
}

// ProfileInfo : 프로젝트에서 사용하는 프로파일 정보
type ProfileInfo struct {
	Name string
	// Path : application-prod.yml 내용 분석 시 재사용
	Path string // 내부 분석용
	File string // 외부 출력용
}