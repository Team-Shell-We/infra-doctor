package i18n

// messages : 언어별 문구 데이터. 로직(Get/Supported/IsSupported)은
// i18n.go에 분리돼 있음. 프로퍼 노운(Docker, PostgreSQL, 파일명 등)이나
// AI가 생성하는 문장은 여기 포함되지 않음 — 전자는 번역 대상이 아니고,
// 후자는 프롬프트에 언어 지시를 넣는 별도 메커니즘을 쓴다
// (internal/ai/explain, internal/ai/recommend의 prompt.go 참고).
var messages = map[string]map[string]string{
	English: {
		// config
		"config.llm":             "LLM",
		"config.language":        "Language",
		"config.output":          "Output",
		"config.autoExport":      "Auto Export",
		"config.enabled":         "Enabled",
		"config.disabled":        "Disabled",
		"config.unconfigured":    "Not configured",
		"config.langFlagDesc":    "Set output language (en, ko)",
		"config.unsupportedLang": "Unsupported language %q (supported: %v)",

		// scan
		"scan.title":          "Project Scan",
		"scan.framework":      "Framework",
		"scan.dependencies":   "Dependencies",
		"scan.database":       "Database",
		"scan.infrastructure": "Infrastructure",
		"scan.cicd":           "CI/CD",
		"scan.trigger":        "Trigger",
		"scan.branch":         "Branch",
		"scan.jobs":           "Jobs",
		"scan.profiles":       "Profiles",

		// doctor
		"doctor.title":               "Infrastructure Doctor",
		"doctor.readiness":           "Deployment Readiness",
		"doctor.infraCheck":          "Infrastructure Check",
		"doctor.recommendation":      "Recommendation",
		"doctor.noIssues":            "No issues found.",
		"doctor.check.docker":        "Docker",
		"doctor.check.dockerCompose": "Docker Compose",
		"doctor.check.healthCheck":   "Health Check",
		"doctor.check.reverseProxy":  "Reverse Proxy",
		"doctor.check.monitoring":    "Monitoring",
		"doctor.check.logRotation":   "Log Rotation",
		"doctor.check.dbBackup":      "DB Backup",

		// login
		"login.title":              "Login",
		"login.selectMethod":       "Select Login Method",
		"login.option1":            "1. OpenAI API Key",
		"login.option2":            "2. Infra Doctor Account",
		"login.choose":             "Choose (1-2): ",
		"login.accountUnavailable": "Infra Doctor Account login is not available yet. Please use an OpenAI API Key for now.",
		"login.invalidChoice":      "Invalid choice. Run 'infra-doctor login' again and choose 1 or 2.",
		"login.apiKeyPrompt":       "OpenAI API Key",
		"login.apiKeyEmpty":        "API Key cannot be empty.",
		"login.invalidKey":         "Invalid OpenAI API Key.",
		"login.checkKey":           "Please check your API Key.",
		"login.verifyFailed":       "Failed to verify API Key: %v",
		"login.saveFailed":         "Failed to save credentials: %v",
		"login.verified":           "OpenAI API Key verified.",
		"login.completed":          "Login completed.",
		"login.provider":           "Provider",

		// explain
		"explain.suffix":                       "Explained",
		"explain.notAdopted":                   "Not adopted in this project yet — the sections below describe what adopting it would look like.",
		"explain.currentProject":               "Current Project",
		"explain.buildFlow":                    "Build Flow",
		"explain.whyTopic":                     "Why %s?",
		"explain.currentStatus":                "Current Status",
		"explain.status.dockerfile":            "Dockerfile",
		"explain.status.dockerCompose":         "Docker Compose",
		"explain.status.healthCheck":           "Health Check",
		"explain.status.dockerComposeYml":      "docker-compose.yml",
		"explain.status.githubActionsWorkflow": "GitHub Actions workflow",
		"explain.status.kubernetesManifests":   "Kubernetes manifests",
		"explain.status.nginxConfig":           "Nginx configuration",
		"explain.status.postgresql":            "PostgreSQL",
		"explain.status.awsSdk":                "AWS SDK dependency",
		"explain.status.relationalDb":          "Relational database (PostgreSQL/MySQL)",
		"explain.status.redis":                 "Redis",

		// recommend
		"recommend.title":        "Deployment Recommendation",
		"recommend.currentStack": "Current Stack",
		"recommend.recommended":  "Recommended",
		"recommend.fit.yes":      "Recommended",
		"recommend.fit.no":       "Not Recommended",
		"recommend.reason":       "Reason",
		"recommend.nextStep":     "Next Step",

		// init
		"init.springBootProject": "Spring Boot project:",
		"init.created":           "Created:",

		// update
		"update.current":  "Current",
		"update.latest":   "Latest",
		"update.devMode":  "Update status cannot be checked in development mode.",
		"update.run":      "Run",
		"update.upToDate": "Infra Doctor is already up to date.",

		// version
		"version.unknown":     "unknown",
		"version.development": "development",
		"version.label":       "Version :",
		"version.goLabel":     "Go :",

		// help
		"help.usage":             "Usage",
		"help.coreCommands":      "Core Commands",
		"help.doctorDesc":        "Analyze project readiness",
		"help.visualizeArchDesc": "Show infrastructure architecture",
		"help.visualizeFlowDesc": "Show deployment flow",
		"help.explainDesc":       "Explain infrastructure concepts",
		"help.recommendDesc":     "Recommend deployment strategy",
		"help.generateDesc":      "Generate configuration files",
		"help.exportDesc":        "Export project report",
		"help.run":               "Run",
		"help.forDetails":        "for detailed information.",

		// about
		"about.tagline":  "AI-powered Infrastructure Analysis CLI",
		"about.builtFor": "Built for backend developers.",
		"about.github":   "GitHub",
		"about.license":  "License",

		// donate
		"donate.thanks": "Thank you for supporting Infra Doctor!",

		// common (shared across commands)
		"common.notLoggedIn":  "You're not logged in. Run 'infra-doctor login' to set up your OpenAI API Key first.",
		"common.openaiFailed": "Failed to reach OpenAI: %v",
		"common.error":        "Error:",
	},

	Korean: {
		// config
		"config.llm":             "LLM",
		"config.language":        "언어",
		"config.output":          "출력 형식",
		"config.autoExport":      "자동 내보내기",
		"config.enabled":         "사용",
		"config.disabled":        "미사용",
		"config.unconfigured":    "설정 안 됨",
		"config.langFlagDesc":    "출력 언어 설정 (en, ko)",
		"config.unsupportedLang": "지원하지 않는 언어 %q (지원: %v)",

		// scan
		"scan.title":          "프로젝트 스캔",
		"scan.framework":      "프레임워크",
		"scan.dependencies":   "의존성",
		"scan.database":       "데이터베이스",
		"scan.infrastructure": "인프라",
		"scan.cicd":           "CI/CD",
		"scan.trigger":        "트리거",
		"scan.branch":         "브랜치",
		"scan.jobs":           "작업",
		"scan.profiles":       "프로파일",

		// doctor
		"doctor.title":               "인프라 진단",
		"doctor.readiness":           "배포 준비도",
		"doctor.infraCheck":          "인프라 점검",
		"doctor.recommendation":      "권장 사항",
		"doctor.noIssues":            "발견된 문제가 없습니다.",
		"doctor.check.docker":        "Docker",
		"doctor.check.dockerCompose": "Docker Compose",
		"doctor.check.healthCheck":   "헬스 체크",
		"doctor.check.reverseProxy":  "리버스 프록시",
		"doctor.check.monitoring":    "모니터링",
		"doctor.check.logRotation":   "로그 로테이션",
		"doctor.check.dbBackup":      "DB 백업",

		// login
		"login.title":              "로그인",
		"login.selectMethod":       "로그인 방법 선택",
		"login.option1":            "1. OpenAI API 키",
		"login.option2":            "2. Infra Doctor 계정",
		"login.choose":             "선택 (1-2): ",
		"login.accountUnavailable": "Infra Doctor 계정 로그인은 아직 지원하지 않습니다. 지금은 OpenAI API 키를 사용해주세요.",
		"login.invalidChoice":      "잘못된 선택입니다. 'infra-doctor login'을 다시 실행해서 1 또는 2를 선택해주세요.",
		"login.apiKeyPrompt":       "OpenAI API 키",
		"login.apiKeyEmpty":        "API 키를 비워둘 수 없습니다.",
		"login.invalidKey":         "유효하지 않은 OpenAI API 키입니다.",
		"login.checkKey":           "API 키를 확인해주세요.",
		"login.verifyFailed":       "API 키 검증 실패: %v",
		"login.saveFailed":         "자격증명 저장 실패: %v",
		"login.verified":           "OpenAI API 키가 확인되었습니다.",
		"login.completed":          "로그인이 완료되었습니다.",
		"login.provider":           "제공자",

		// explain
		"explain.suffix":                       "설명",
		"explain.notAdopted":                   "아직 이 프로젝트에 도입되지 않았습니다 — 아래는 도입할 경우의 시나리오입니다.",
		"explain.currentProject":               "현재 프로젝트",
		"explain.buildFlow":                    "빌드 흐름",
		"explain.whyTopic":                     "왜 %s인가?",
		"explain.currentStatus":                "현재 상태",
		"explain.status.dockerfile":            "Dockerfile",
		"explain.status.dockerCompose":         "Docker Compose",
		"explain.status.healthCheck":           "헬스 체크",
		"explain.status.dockerComposeYml":      "docker-compose.yml",
		"explain.status.githubActionsWorkflow": "GitHub Actions 워크플로",
		"explain.status.kubernetesManifests":   "Kubernetes 매니페스트",
		"explain.status.nginxConfig":           "Nginx 설정",
		"explain.status.postgresql":            "PostgreSQL",
		"explain.status.awsSdk":                "AWS SDK 의존성",
		"explain.status.relationalDb":          "관계형 데이터베이스 (PostgreSQL/MySQL)",
		"explain.status.redis":                 "Redis",

		// recommend
		"recommend.title":        "배포 추천",
		"recommend.currentStack": "현재 스택",
		"recommend.recommended":  "추천",
		"recommend.fit.yes":      "권장",
		"recommend.fit.no":       "비권장",
		"recommend.reason":       "이유",
		"recommend.nextStep":     "다음 단계",

		// init
		"init.springBootProject": "Spring Boot 프로젝트:",
		"init.created":           "생성됨:",

		// update
		"update.current":  "현재",
		"update.latest":   "최신",
		"update.devMode":  "개발 모드에서는 업데이트 상태를 확인할 수 없습니다.",
		"update.run":      "실행",
		"update.upToDate": "Infra Doctor가 이미 최신 버전입니다.",

		// version
		"version.unknown":     "알 수 없음",
		"version.development": "개발 버전",
		"version.label":       "버전 :",
		"version.goLabel":     "Go :",

		// help
		"help.usage":             "사용법",
		"help.coreCommands":      "핵심 명령어",
		"help.doctorDesc":        "프로젝트 준비도 분석",
		"help.visualizeArchDesc": "인프라 아키텍처 표시",
		"help.visualizeFlowDesc": "배포 흐름 표시",
		"help.explainDesc":       "인프라 개념 설명",
		"help.recommendDesc":     "배포 전략 추천",
		"help.generateDesc":      "설정 파일 생성",
		"help.exportDesc":        "프로젝트 리포트 내보내기",
		"help.run":               "실행",
		"help.forDetails":        "자세한 정보를 확인하세요.",

		// about
		"about.tagline":  "AI 기반 인프라 분석 CLI",
		"about.builtFor": "백엔드 개발자를 위해 만들었습니다.",
		"about.github":   "GitHub",
		"about.license":  "라이선스",

		// donate
		"donate.thanks": "Infra Doctor를 후원해주셔서 감사합니다!",

		// common (shared across commands)
		"common.notLoggedIn":  "로그인이 필요합니다. 'infra-doctor login'으로 먼저 OpenAI API 키를 설정해주세요.",
		"common.openaiFailed": "OpenAI 호출 실패: %v",
		"common.error":        "오류:",
	},
}
