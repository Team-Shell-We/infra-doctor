# 🩺 Infra Doctor

[English](README_EN.md) | **한국어**

> **Spring Boot 프로젝트를 위한 AI 기반 인프라 분석 CLI**

Infra Doctor는 Spring Boot 프로젝트의 현재 인프라 구성을 분석해서, 개발자가 배포 환경을 이해하고 진단하고 시각화하고 개선할 수 있도록 돕는 CLI 도구입니다.

프로젝트 설정을 스캔하고 인프라 구성 요소를 감지해서, 실제로 적용 가능한 인사이트를 제공합니다.

---

## ✨ 기능

| 분류 | 명령어 | 설명 | 상태 |
| :------ | :------ | :---------- | :----: |
| 🔍 분석 | `scan` | 프로젝트를 스캔해 프레임워크, 데이터베이스, Docker, CI/CD, 프로파일, 인프라 정보를 수집합니다. | ✅ |
| 🩺 진단 | `doctor` | 규칙 기반 체크리스트로 배포 준비도(0~100%)를 진단하고 개선 방향을 제안합니다. | ✅ |
| 🤖 AI | `explain <topic>` | 특정 인프라 기술이 현재 프로젝트에서 어떻게 쓰이는지 설명합니다. | ✅ |
| 🤖 AI | `recommend` | 감지된 인프라 복잡도를 기반으로 배포 전략(Docker Compose vs Kubernetes)을 추천합니다. | ✅ |
| 🏗 시각화 | `visualize architecture` | 현재 인프라의 아키텍처 다이어그램을 생성합니다. | ✅ |
| 🏗 시각화 | `visualize flow` | 빌드·배포 워크플로를 시각화합니다. | ✅ |
| ⚙️ 생성 | `generate <target>` | 스캔 결과를 바탕으로 인프라 설정 파일을 생성합니다. | ✅ |
| 📄 내보내기 | `export` | 전체 분석 결과(리포트, 다이어그램, 생성 파일)를 디렉터리 하나로 내보냅니다. | ✅ |

> **상태**
>
> - ✅ 완료
> - 🚧 진행 중
> - 📅 예정

`explain`/`recommend`는 먼저 `login`이 필요합니다. AI가 생성하는 내용은 항상 스캐너가 이미 결정론적으로 확인한 사실에 근거합니다 — AI는 어떤 파일이 있는지 지어내지 않고, 왜 그게 중요한지만 설명합니다 (전체 설계 근거는 [docs/기능명세서.md](docs/기능명세서.md) 참고).

Gradle(`build.gradle`/`build.gradle.kts`)과 Maven(`pom.xml`) 프로젝트 모두 빌드 도구로 감지됩니다. **알려진 한계:** Maven 프로젝트는 현재 프레임워크/의존성/데이터베이스 감지만 동작합니다 — Docker/Compose/Kubernetes/CI/프로파일 감지는 아직 Gradle 전용이라, 잘 갖춰진 Maven 프로젝트도 `doctor`/`recommend`/`export`가 준비도를 실제보다 낮게 보고할 수 있습니다. 수정 예정인 버그로 추적 중입니다. 자세한 내용은 [docs/기능명세서.md](docs/기능명세서.md) 참고.

---

## 🔧 유틸리티

| 명령어 | 설명 | 상태 |
| :------ | :---------- | :----: |
| `init` | 현재 프로젝트에 Infra Doctor 설정을 초기화합니다. | ✅ |
| `config` | CLI 설정(LLM 제공자, 출력 언어, 출력 형식, 자동 내보내기)을 보거나 변경합니다. | ✅ |
| `login` | OpenAI API 키로 인증해 AI 기능을 활성화합니다. | ✅ |
| `update` | 최신 CLI 버전을 확인합니다. | ✅ |
| `version` | 현재 CLI 버전을 표시합니다. | ✅ |
| `help` | 사용 가능한 명령어와 사용법을 보여줍니다. | ✅ |
| `about` | 프로젝트 정보와 크레딧을 표시합니다. | ✅ |
| `donate` | Infra Doctor 프로젝트를 후원합니다. | ✅ |

> `config --lang ko`는 모든 명령어의 고정 CLI 텍스트(라벨, 헤더, 메시지)를 한글로 바꾸고, `explain`/`recommend`의 AI 텍스트도 같은 언어로 생성합니다. 아직 영어로만 나오는 건 `doctor`의 진단 룰 텍스트(`internal/doctor/rules/*.yaml`의 `message`/`reason`/`fix`)뿐입니다 — 번역하려면 이중 언어 YAML 필드 설계가 필요한데 아직 안 돼 있습니다.

---

## 🚀 설치

### Go로 설치

```bash
go install github.com/Team-Shell-We/infra-doctor@latest
```

### 소스에서 빌드

```bash
git clone https://github.com/Team-Shell-We/infra-doctor.git

cd infra-doctor

go build -o infra-doctor
```

---

## 📖 사용법

### 일반적인 흐름

신규 Spring Boot 프로젝트는 보통 이 순서로 명령어를 씁니다.

```bash
cd ~/workspace/my-project

infra-doctor init                 # .infra-doctor/config.yaml 생성
infra-doctor scan                 # 뭐가 감지됐는지 확인
infra-doctor doctor                # 부족한 것 확인 + 준비도 점수
infra-doctor login                 # explain/recommend 쓰려면 최초 1회 필요
infra-doctor recommend             # Docker Compose vs Kubernetes 추천 + 이유
infra-doctor generate compose      # recommend가 안내한 파일 생성
infra-doctor doctor                # 점수가 올랐는지 재확인
infra-doctor visualize architecture
```

아래 모든 명령어에서 `[path]`는 생략 가능하며, 생략하면 현재 디렉터리를 씁니다.

| 명령어 | 주요 플래그 | 설명 |
| :--- | :--- | :--- |
| `scan [path]` | – | 감지된 프레임워크/의존성/DB/인프라/CI/프로파일 요약 |
| `doctor [path]` | `--json`, `--fail-under <score>` | 준비도 점수 + 개선 목록. CI 게이트로 사용 가능 |
| `login` | – | OpenAI API 키 등록 (`explain`/`recommend` 사전 조건) |
| `explain <topic> [path]` | – | topic: compose·container·docker·github-actions·image·k8s·nginx·postgres·rds·redis |
| `recommend [path]` | – | Docker Compose vs Kubernetes 추천 + 이유 |
| `generate <target> [path]` | `-f/--force`, `--dry-run`, `-o/--output-dir`, `--config` | target: docker·compose·nginx·ci·k8s·architecture |
| `visualize architecture\|flow [path]` | `--format`, `--output` | ascii/mermaid/markdown 다이어그램 |
| `export [path]` | `-f/--force`, `--dry-run` | 위 결과 전체를 `infra-doctor/` 디렉터리 하나로 |
| `config [--lang en\|ko]` | – | CLI 설정 조회/변경 |
| `init` | – | `.infra-doctor/config.yaml` 생성 |

`version`, `update`, `help`, `about`, `donate`는 인자·로그인 없이 바로 씁니다. 각 명령어의 상세 동작과 출력 예시, `generate`/`init`의 설정 파일 스키마는 [docs/기능명세서.md](docs/기능명세서.md)를, 플래그 전체 목록은 `infra-doctor <command> --help`를 참고하세요.

---

## 📂 프로젝트 구조

```text
infra-doctor
├── .github
├── cmd
├── docs
├── internal
├── scripts
├── CONTRIBUTING.md
├── go.mod
├── main.go
└── README.md
```

---

## 🛣️ 로드맵

- [x] 프로젝트 스캐너
- [x] 프레임워크 감지 (Gradle; Maven은 일부만 — 위 안내 참고)
- [x] 데이터베이스 감지
- [x] Docker 감지
- [x] 프로파일 감지
- [x] GitHub Actions 분석
- [x] Docker Compose 분석
- [x] 배포 진단(Doctor)
- [x] AI 기반 Explain
- [x] AI 기반 추천
- [x] 인프라 시각화
- [x] 설정 파일 생성기
- [x] 리포트 내보내기
- [ ] Kubernetes 심층 분석
- [ ] SSH / EC2 분석
- [ ] 전 명령어 다국어 출력

---

## 🤝 기여하기

기여는 언제나 환영합니다. 개발 환경 설정, 변경 사항 테스트 방법, 진단 룰·번역 추가 방법은 [CONTRIBUTING.md](CONTRIBUTING.md)를 참고하세요.

---

## 📄 라이선스

이 프로젝트는 MIT License를 따릅니다.
