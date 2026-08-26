# Infra Doctor 기여 가이드

[English](CONTRIBUTING_EN.md) | **한국어**

기여를 고려해주셔서 감사합니다. 이 가이드는 개발 환경 설정, 변경 사항 테스트 방법, 그리고 특정 종류의 기여(진단 룰 추가, 번역 추가)를 다룹니다.

---

## 1. Git 설정

저장소를 클론한 뒤, 설정 스크립트를 한 번 실행하세요.

**Windows (PowerShell)**

```powershell
.\scripts\setup.ps1
```

이 스크립트는 다음을 자동으로 설정합니다.

- Git Hooks (`.githooks`)
- 커밋 메시지 템플릿 (`.gitmessage`)

아래처럼 나오면 정상입니다.

```text
Configuring Git...

Done!
[OK] Git Hooks configured
[OK] Commit template configured
```

---

## 2. 테스트 프로젝트 준비

개발용으로 `examples` 디렉터리 아래에 Spring Boot 프로젝트를 준비하세요.

```text
infra-doctor
├── .github
├── cmd
├── docs
├── examples
│   └── spring-gradle
│       ├── <your-project>
│       └── ...
├── internal
└── ...
```

Spring Boot 프로젝트를 `examples/spring-gradle`에 복사합니다.

예:

```text
examples
└── spring-gradle
    └── .github
    └── <your-project-name>
        ├── gradle
        ├── src
        ├── build.gradle
        ├── settings.gradle
        ├── gradlew
        ├── gradlew.bat
        └── ...
```

커밋하기 전에 복사한 프로젝트에서 아래 디렉터리를 제거하세요.

```text
.git/
.gradle/
.idea/
build/
node_modules/
target/
```

테스트 프로젝트는 **로컬 개발·테스트 전용**입니다. 민감한 소스코드나 운영 설정 파일은 저장소에 포함하지 마세요.

---

## 3. 변경 사항 테스트

예제 프로젝트는 이 저장소에 **포함돼 있지 않습니다** — 위 2번에서 준비한 테스트 프로젝트나 로컬에 있는 아무 Spring Boot 프로젝트를 쓰면 됩니다.

```bash
go run . scan <project-path>
```

예시

```bash
go run . scan ../spring-project
```

더 풍부한 분석을 원한다면 프로젝트에 다음이 포함되면 좋습니다.

```text
build.gradle
pom.xml

Dockerfile
docker-compose.yml

.github/workflows/

application.yml
application-dev.yml
application-prod.yml
```

다음 디렉터리는 분석에서 자동으로 제외됩니다.

```text
.git/
.gradle/
.idea/
build/
node_modules/
target/
.infra-doctor/
```

PR을 올리기 전에 전체 테스트와 포맷 검사를 실행하세요.

```bash
go build ./...
go vet ./...
gofmt -l .
go test ./...
```

---

## 4. 진단 룰 추가

`doctor`의 체크 항목은 Go 코드가 아니라 YAML로 정의됩니다. 룰 스키마와 PR 체크리스트는 [internal/doctor/rules/README.md](internal/doctor/rules/README.md)를 참고하세요.

---

## 5. 번역 추가

CLI 출력은 `internal/i18n`을 통해 영어/한글을 지원합니다. 새 문구를 추가하는 법과 번역하면 안 되는 것들은 [internal/i18n/README.md](internal/i18n/README.md)를 참고하세요.

---

## Pull Request 올리기

PR 템플릿을 사용하세요 — 연관 이슈, 개요, 구현한 기능, 막혔던 부분, 리뷰어가 알아야 할 내용을 물어봅니다. 커밋 메시지는 `<type>(#이슈번호): <요약>` 형식(`feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `ci`)을 따라야 하며, 1번의 `commit-msg` 훅이 로컬에서 이를 강제합니다.
