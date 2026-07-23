# doctor 진단 룰 작성 가이드

`infra-doctor doctor`가 사용하는 진단 룰은 Go 코드가 아니라 이 디렉터리의 YAML 파일로 정의됩니다. 새로운 진단 항목을 추가하거나 기존 항목의 문구/점수를 수정하려면 Go 코드를 몰라도 이 디렉터리의 YAML만 수정하면 됩니다.

## 파일 구조

카테고리(체크 그룹)마다 파일이 하나씩 있습니다. 최상위 YAML 키가 카테고리명이고, 그 아래 `id: {필드...}` 형태로 룰을 나열합니다.

| 파일 | 최상위 키 | 담당 |
|---|---|---|
| `deployment.yaml` | `deployment` | Docker/GitHub Actions 조합 진단 |
| `production.yaml` | `production` | 운영 환경 하드닝(Health Check, Reverse Proxy, Monitoring, Log Rotation, DB Backup) |
| `localdev.yaml` | `localdev` | 로컬 개발 편의성(Compose, dev 프로파일) |

기존 파일에 룰을 추가/수정하는 건 YAML만 고치면 되지만, 완전히 새로운 카테고리(새 파일)를 추가하려면 `internal/doctor/loader.go`에 로드 로직을 먼저 연결해야 합니다 — 이 경우는 Go 코드 변경이 필요합니다.

## 룰 필드 규칙

```yaml
production:
  no_health_check:            # id: snake_case, "no_" 접두사 권장

    category: Infrastructure  # Infrastructure | CICD | Database | Monitoring | Security 중 하나
    level: Warning             # Info | Warning | Critical 중 하나 (대소문자는 검증 시 정규화되지만 이 표기를 따를 것)
    score: -10                 # 반드시 음수. 레벨별 권장 범위는 아래 참고

    title: Health check is not configured        # 한 줄 요약

    message: >                                     # 무엇이 감지/미감지됐는지 사실 서술
      No Docker HEALTHCHECK instruction or Compose healthcheck block was found.

    reason: >                                      # 왜 문제인지
      Without a health check, orchestrators and load balancers cannot detect
      an unhealthy instance and route traffic away from it.

    fix: >                                         # 구체적인 실행 한 가지 (명령형, 짧게)
      Add a HEALTHCHECK instruction to your Dockerfile or a healthcheck
      block to docker-compose.yml.
```

- **id**: `^[a-z][a-z0-9_]*$` (snake_case). 다른 파일의 id와도 겹치지 않게 지어주세요.
- **category / level**: 오타나 다른 표기(예: `CI/CD`, `warning`)는 `internal/doctor/rules_test.go`의 스키마 테스트에서 걸러집니다.
- **score**: 레벨별 권장 범위 (기존 룰 기준) — Info `-5`, Warning `-10 ~ -15`, Critical `-20 ~ -40`.
- **fix**: CLI의 60자 박스 UI 안에서 자동 줄바꿈되어 출력됩니다. 한 문장, 한 가지 행동으로 짧게 쓰세요 (대략 100자 이내 권장).

## PR 리뷰 체크리스트

- [ ] id가 snake_case이고 다른 룰과 겹치지 않는지
- [ ] category/level 값이 위 표기와 정확히 일치하는지 (오타 없는지)
- [ ] score가 음수이고, level 대비 과하거나 약하지 않은지
- [ ] fix가 "개선하세요" 같은 모호한 표현이 아니라 구체적인 한 가지 액션인지
- [ ] 로컬에서 `go test ./internal/doctor/...` 통과 확인했는지

## 로컬 검증

```bash
go test ./internal/doctor/...
```

이 테스트는 CI(`.github/workflows/ci.yml`)에서 모든 PR마다 자동으로도 실행됩니다.
