# CLI 출력 다국어(i18n) 가이드

`infra-doctor`의 CLI 출력은 두 가지 서로 다른 메커니즘으로 다국어를 지원합니다. 새 명령어(`generate`/`export`/`visualize` 등)를 추가하거나 기존 명령어에 문구를 추가할 때는 그 문구가 아래 둘 중 어디에 해당하는지부터 판단하세요.

## 두 가지 메커니즘

| | 대상 | 처리 방식 |
|---|---|---|
| **A. 고정 문구** | 헤더, 라벨, 에러 메시지 등 CLI 코드가 직접 작성한 문장 | `internal/i18n`의 key-lookup (`i18n.Get(lang, key)`) |
| **B. AI 생성 문장** | `explain`/`recommend`의 LLM 응답 | 정적 번역 대상이 아님 — 프롬프트에 "Respond entirely in {{.Language}}." 지시를 넣어 모델이 직접 해당 언어로 생성 |

B는 모델 출력이라 사전에 번역 문구를 준비할 수 없습니다. `internal/ai/explain/prompt.go`, `internal/ai/recommend/prompt.go`의 `BuildRequest(..., lang string)` 패턴을 그대로 따르면 됩니다 — 새 AI 생성 명령어를 추가할 때만 해당.

이 문서는 A(고정 문구)를 추가하는 방법을 다룹니다.

## 파일 구조

- `i18n.go` — 로직만. `Get(lang, key)`(폴백: lang에 없으면 영어 → 그것도 없으면 key 그대로), `Supported()`, `IsSupported(lang)`.
- `messages.go` — 데이터만. `map[string]map[string]string`, `<command>.<label>` 형식의 key로 명령어별 섹션 주석 아래 정리돼 있음.

## 새 문구 추가하는 법

1. `messages.go`의 `English` 맵과 `Korean` 맵 **양쪽 모두**에 같은 key를 추가합니다. 명령어별 섹션(`// scan`, `// doctor` 등) 아래 알파벳/등장 순서에 맞춰 넣으세요. 새 명령어라면 새 섹션 주석을 추가합니다.
2. key 네이밍은 `<command>.<label>` (예: `generate.title`, `generate.noTemplate`). 여러 명령어에서 완전히 동일한 문구(로그인 안내, 공통 에러 등)는 `common.*`으로 통합합니다.
3. 호출부(`cmd/*.go`)에서 `lang := currentLang()`(`cmd/lang.go`)로 언어를 한 번 읽고 `i18n.Get(lang, "your.key")`로 조회합니다.
4. `go test ./internal/i18n/...`를 실행해 `TestLanguagesHaveMatchingKeys`가 통과하는지 확인합니다 — en/ko 두 맵의 key 집합이 다르면 이 테스트가 실패합니다. 한쪽 언어만 추가하고 잊어버리는 실수를 여기서 잡습니다.

## 번역하면 안 되는 것

- 고유명사: Docker, PostgreSQL, Redis, Kubernetes, GitHub Actions 등
- 파일명/경로: `docker-compose.yml`, `Dockerfile` 등
- 동적 값: 버전 번호, 실제 파일 경로, 셸 명령어
- `doctor`의 진단 룰(`internal/doctor/rules/*.yaml`)의 `message`/`reason`/`fix` 문장 — 룰마다 자유 서술문이라 이중 언어 YAML 필드 설계가 필요한 별도 작업입니다(현재 미지원, `internal/doctor/rules/README.md` 참고).

## 박스 UI와 폭 계산

`internal/ui/box.go`의 `Header`/`Line`/`Blank`/`Wrap`/`center`는 `DisplayWidth()`로 텍스트 폭을 계산합니다. 한글/한자/가나 등 East Asian Wide 문자는 터미널에서 2칸을 차지하므로, `len()`(바이트 길이)이나 `utf8.RuneCountInString()`(글자 개수)을 직접 쓰면 한글이 섞인 문구에서 박스 정렬이 깨집니다. 새 라벨을 박스에 출력할 때는 반드시 `ui.Line`/`ui.Header`/`ui.Wrap`을 통해서만 출력하세요 — 직접 `fmt.Printf("%-Ns", ...)` 같은 서식을 쓰지 마세요.

## PR 리뷰 체크리스트

- [ ] `messages.go`의 `English`/`Korean` 양쪽에 key를 추가했는지
- [ ] key가 `<command>.<label>` 형식이고 다른 key와 안 겹치는지
- [ ] 고유명사/파일명/동적 값을 번역하지 않았는지
- [ ] 박스에 출력되는 문구라면 `ui.Line`/`ui.Header`/`ui.Wrap`을 거쳐 출력하는지 (직접 폭 계산 금지)
- [ ] `go test ./internal/i18n/... ./internal/ui/...` 통과 확인했는지

## 로컬 검증

```bash
go test ./internal/i18n/... ./internal/ui/...

# 실제 출력 확인
go run . config --lang ko
go run . <command>
go run . config --lang en   # 확인 후 원래대로 되돌리기
```

이 테스트들은 CI(`.github/workflows/ci.yml`)에서 모든 PR마다 자동으로도 실행됩니다.
