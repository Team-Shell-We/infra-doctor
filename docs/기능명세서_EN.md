# Infra Doctor Feature Specification

**English** | [한국어](기능명세서.md)

> Last updated against: `develop` (PR #45 full i18n coverage, #46 explain topic-not-adopted clarification, #49 visualize, #50/#51 generate, #53 recommend copy cleanup, #54 recommend complexity signal expansion, #55 Next Step layer, #56 doctor CI gate, #64 Maven support, #66 update release API, #67 Gradle Kotlin DSL detection fix, #70 Maven infrastructure detection fix, #77 help/timeout/recommend label/visualize DB node cleanup, #81 scaffold/dead code cleanup, #82 v0.2.0 release, #84 duplicate logic consolidation, #85 comment style unification, #79 CONTRIBUTING split, #94 donate Ko-fi link, #95 update pseudo-version fallback, #96 GitHub Actions shorthand trigger parsing, #97 ui.Wrap force-break, #98 test coverage)

## Project Description

An AI-powered CLI for operational diagnosis and infrastructure guidance. An open-source tool that analyzes a Spring Boot project, visualizes its current operational structure, diagnoses what's missing for production, and auto-generates Docker/Kubernetes/CI-CD config files plus improvement guidance.

## Pain Point

- Backend developers can build applications but struggle when it comes to setting up the operational environment (when to use Docker, where Nginx fits, why Redis is needed, when to adopt Kubernetes, etc.)
- There's plenty of information online, but getting an answer specific to *your* project usually means pasting your whole codebase into an AI and asking — a lot of manual effort.
- Infra Doctor aims to analyze the project itself, explain the current structure from an operations standpoint, and provide both a direction for improvement and config files you can actually apply.

## Core Design Principle

Applies to every AI command (`explain`, `recommend`).

> **AI never judges facts — it only explains facts the code has already determined.**

- `doctor` uses no AI at all — it is 100% rule-based (YAML rules).
- `explain`'s "Current Status" (whether related files/settings exist) is decided deterministically by Go code based on `project.Info`; the AI only explains "why it matters" in natural language. This keeps the AI from inventing content the scan never found.
- `recommend`'s deployment strategy recommendation (Docker Compose vs. Kubernetes) is likewise decided deterministically by code, with the AI only writing out the reasoning (Reason) in prose.
- Because of this principle, the AI commands have reproducible core information (checklists, recommended values, Next Step, etc.) guaranteed by code, while the AI's role is scoped to explanation.

---

## Usage Scenario

An example command sequence for onboarding a new Spring Boot project. See "Detailed Command Specification" below for each command's full spec.

```bash
cd my-project

infra-doctor init                  # create .infra-doctor/config.yaml (for generate customization; works fine without it)
infra-doctor scan                  # see what's detected first
infra-doctor doctor                 # readiness score + what's missing
infra-doctor login                  # one-time, needed for explain/recommend

infra-doctor recommend              # Docker Compose vs Kubernetes recommendation + reasons
infra-doctor explain docker         # why it's needed, explained in project context

infra-doctor generate compose       # actually generate the files recommend pointed to
infra-doctor doctor                 # re-check whether the score went up

infra-doctor visualize architecture # view the finished infrastructure as a diagram
```

In a CI pipeline, `doctor --fail-under <score>` can gate infrastructure regressions (issue #42, see the `doctor` entry below).

Instead of running the above commands one at a time, `infra-doctor export` can collect the report + diagrams + generated files into one `infra-doctor/` directory in a single run (issue #39, see the `export` entry below).

`generate` also works fine without creating `.infra-doctor/config.yaml` first — in that case it falls back to scan-based defaults for port/service name/image, and writes generated files directly to the project root. The default config `init` creates sets `output.directory: .infra-doctor/generated`, so after `init`, generated files go into that subdirectory instead of the root.

---

## Detailed Command Specification

### scan

| Item | Content |
|---|---|
| Command | `infra-doctor scan [path]` |
| Status | ✅ Done (MVP) |
| Uses AI | N |
| Owner | gimye |

**Description**: Analyzes the project and detects the technology stack and infrastructure components currently in use. Omitting `path` analyzes the current directory.

**What it analyzes**: Spring Boot version / Gradle·Maven / Java version, Spring Security·JPA·Kafka·AWS SDK·Lombok·Actuator·OpenAPI dependencies, PostgreSQL·MySQL·MariaDB·Redis, Docker·Docker Compose·Kubernetes, GitHub Actions (including triggers/branches/jobs), profiles (dev/prod/test/default).

**Example output**
```
🔍 Project Scan

 Framework
   ✓ Spring Boot 3.5.7
   ✓ Gradle
   ✓ Java 17

 Dependencies
   ✓ Spring Security
   ✓ Spring Data JPA
   ...

 Database
   ✓ PostgreSQL
   ✓ Redis

 Infrastructure
   ✓ Docker
      └─ Dockerfile
   ✗ Docker Compose
   ✗ Kubernetes

 CI/CD
   ✗ GitHub Actions

 Profiles
   ✓ dev
      └─ application-dev.yml
   ...
```

Build file (`build.gradle`/`pom.xml`) discovery only happens in the given `path` directory, while Docker/Compose/K8s/GitHub Actions/Profile file discovery recurses into subdirectories, excluding `.git`/`.gradle`/`.idea`/`build`/`node_modules`/`target`/`.infra-doctor` (`internal/analyzer/walk.go`). This applies identically to Gradle and Maven — however, dependency/database/framework detection itself only reads the single root build file for both build tools, so a multi-module project whose real dependencies are declared in a child module can have its detection missed.

---

### doctor

| Item | Content |
|---|---|
| Command | `infra-doctor doctor [path] [--json] [--fail-under <score>]` |
| Status | ✅ Done |
| Uses AI | **N (originally spec'd as Y, switched to 100% rule-based during implementation)** |
| Owner | gimye |

**Description**: Diagnoses the project's deployment readiness as a 0–100% score. Diagnostic rules fall into three categories, all managed as `internal/doctor/rules/*.yaml` with automated schema validation (tests + CI).

**CI gate (issue #42)**: `--json` prints `doctor.Result` as JSON instead of the box UI (the default human-readable output is unchanged when the flag isn't passed). `--fail-under <score>` calls `os.Exit(1)` when the score is below the threshold — without the flag, it always exits 0 as before. The two flags are independent, so you can use either or both.

- **deployment**: Docker/GitHub Actions combination (`no_deployment` -40 Critical, `no_docker` -20 Warning, `no_github_actions` -15 Warning)
- **production**: production hardening (`no_health_check`/`no_reverse_proxy`/`no_monitoring` each -10~-15 Warning, `no_log_rotation` -5 Info, `no_backup` -20 Critical)
- **localdev**: local development convenience (`no_compose_with_dependencies` -10 Warning, `no_dev_profile` -5 Info)

**Infrastructure Check checklist (7 items, deterministic)**: Docker, Docker Compose, Health Check, Reverse Proxy, Monitoring, Log Rotation, DB Backup

**Example output**
```
🩺 Infrastructure Doctor

Deployment Readiness
 ████░░░░░░░░░░░░░░░░░░░░░░░░░░ 15%

Infrastructure Check
 ✓ Docker
 ✗ Docker Compose
 ✗ Health Check
 ✗ Reverse Proxy
 ✗ Monitoring
 ✗ Log Rotation
 ✗ DB Backup

Recommendation
 • Configure GitHub Actions.
 • Add a docker-compose.yml that starts local database/Redis dependencies with one command.
 ...
```

**Contributor guide**: `internal/doctor/rules/README.md` documents the rule-authoring conventions (id naming, category/level enums, score ranges, PR checklist).

---

### login

| Item | Content |
|---|---|
| Command | `infra-doctor login` |
| Status | ✅ Done |
| Uses AI | N (prerequisite for other AI commands) |
| Owner | gimye |

**Description**: Login for using the AI commands (`explain`, `recommend`). Offers a choice of method.

1. **OpenAI API Key** — ✅ Implemented. Masked input (`golang.org/x/term`), verifies the key against the real OpenAI API, then saves to `~/.infra-doctor/config.json` (`{"provider":"openai","apiKey":"...","login":true}`, file permission `0600`/directory `0700`).
2. **Infra Doctor Account** — ❌ No backend. Selectable in the menu, but only prints a "not yet supported" message.

**Note**: The `~/.infra-doctor/config.json` saved here (home directory, AI credentials + global CLI settings) is an entirely different file from the `.infra-doctor/config.yaml` (analysis settings) that `init` creates in the project directory.

**Example output**
```
🔐 Login

Select Login Method
  1. OpenAI API Key
  2. Infra Doctor Account
Choose (1-2): 1

OpenAI API Key
> ********

✅ OpenAI API Key verified.
✅ Login completed.

Provider
OpenAI
```

---

### explain <topic>

| Item | Content |
|---|---|
| Command | `infra-doctor explain <topic> [path]` |
| Status | ✅ Done |
| Uses AI | Y (some sections only) |
| Argument | compose, container, docker, github-actions, image, k8s, nginx, postgres, rds, redis |
| Owner | gimye |

**Description**: Explains what role a given technology plays **in the context of the currently scanned project**. Requires `login`.

**How each section is generated**:
| Section | Generated by |
|---|---|
| Current Project | AI (based on the scan summary, fixed prompt) |
| Build Flow | AI |
| Why \<Topic\>? | AI |
| Current Status | **Code (deterministic)** — decided by a hardcoded mapping of which files/settings to check per topic (`internal/ai/explain/status.go`); the AI has no involvement |

When explaining a topic not yet present in the project (e.g. `explain k8s` without k8s), "Current Status" is shown before the AI's narration, with a code-decided "not yet adopted" banner at the top (`explain.TopicPresent`). The AI prompt also instructs it to phrase any fact marked absent in the subjunctive.

**Example output**
```
💡 Docker Explained

Current Project
 ✓ Uses Spring Boot 3.5.7 with Gradle, leveraging Java 17.
 ✓ Relies on PostgreSQL and Redis for data management.
 ...

Build Flow
 Dockerfile defines the build instructions for the application.
 ...

Why Docker?
 • Docker ensures consistent environment setups across development,
   testing, and production.
 ...

Current Status
 ✓ Dockerfile
 ✗ Docker Compose
 ✗ Health Check
```

**When not logged in**: prints `You're not logged in. Run 'infra-doctor login' to set up your OpenAI API Key first.` and exits.

---

### recommend

| Item | Content |
|---|---|
| Command | `infra-doctor recommend [path]` |
| Status | ✅ Done |
| Uses AI | Y (some sections only) |
| Owner | gimye |

**Description**: Recommends a deployment strategy (Docker Compose vs. Kubernetes) based on static complexity signals detected in the project (Kafka/relational DB/Redis usage, CI workflow/API endpoint/module counts). True "scale" indicators like traffic, team size, or microservice count aren't observable from a single-repo scan, so the command deals only with signals that can be checked statically in code. Requires `login`.

**6 signals (`complexityThreshold = 4`)**: Kafka usage, relational DB usage, Redis usage, 2+ CI workflows, more than 20 API endpoints (`internal/analyzer/api.go` counts `@RestController`/`@GetMapping`-style annotations in `.java` sources), multi-module structure (for Gradle, `internal/analyzer/modules.go` counts `include` statements in `settings.gradle`; for Maven, `AnalyzeMavenModules` counts `<modules>` entries in `pom.xml`). A project already using Kubernetes is kept regardless of these signals, but if replicas are 2 or more, an "already scaling" rationale is added to the Reason (`internal/analyzer/infrastructure.go` also parses the `replicas:` value from manifests).

**How each section is generated**:
| Section | Generated by |
|---|---|
| Current Stack | Code (reuses `ai.Summary`) |
| Recommended | **Code** — keeps Kubernetes if manifests already exist, otherwise decides Kubernetes vs. Docker Compose by the count of complexity signals (Kafka/DB/Redis/multiple CI workflows) |
| Kubernetes Fit | Code |
| Reason | AI — expands the code-decided recommendation reasons (labels) into natural language. **The system prompt forces it not to change the recommendation itself** |
| Next Step | Code — suggests `infra-doctor generate ...` commands based on what's missing (compose/k8s manifest/CI) |

**Example output**
```
🚀 Deployment Recommendation

Current Stack
 ✓ Spring Boot 3.5.7 (Gradle), Java 17
 ✓ PostgreSQL
 ✓ Redis

Recommended
 ⭐ Docker Compose

Kubernetes
 ✗ Not Recommended

Reason
 • Your project utilizes a single Spring Boot API server, making
   Docker Compose an ideal fit...
 ...

Next Step
 infra-doctor generate docker
 infra-doctor generate ci
```

---

### visualize architecture

| Item | Content |
|---|---|
| Command | `infra-doctor visualize architecture [path] --format ascii\|mermaid\|markdown --output <file>` |
| Status | ✅ Done (PR #49, issue #38) |
| Uses AI | N |
| Output | ASCII diagram, Mermaid diagram, or a Markdown file |
| Owner | jade |

**Description**: Analyzes the project and automatically visualizes the current operational architecture. No AI needed — it can be drawn from `scan` results alone. `internal/visualize` converts `project.Info` into a node/edge graph and renders it in one of three formats.

---

### visualize flow

| Item | Content |
|---|---|
| Command | `infra-doctor visualize flow [path] --format ascii\|mermaid\|markdown --output <file>` |
| Status | ✅ Done (PR #49, issue #38) |
| Uses AI | N |
| Output | Build & deploy flow diagram |
| Owner | jade |

**Description**: Visualizes the process from build to deployment as a flow diagram. The GitHub Actions workflow is searched for in the current or a nearby parent directory.

---

### generate <target>

| Item | Content |
|---|---|
| Command | `infra-doctor generate <target> [path] [--force] [--dry-run] [--output-dir] [--config]` |
| Status | ✅ Done (PR #50/#51, issue #37) |
| Uses AI | N |
| Argument | docker, compose, nginx, k8s, ci, architecture (= everything needed, at once) |
| Output | Dockerfile, docker-compose.yml, nginx.conf, GitHub Actions workflow, Kubernetes manifests (deployment/service/configmap) |
| Owner | jade & gimye |

**Description**: Generates config files at the project root based on project analysis + `doctor` diagnosis results (`internal/generate`, a `Generator` interface + `Registry` pattern, `.tmpl` + `go:embed`). It reads the `output`/`generate` settings in `.infra-doctor/config.yaml` for real, letting you customize port/service name/image/namespace, etc. Refuses to overwrite an existing file without `--force`; `--dry-run` previews what would be written. `recommend`'s Next Step points here.

The `architecture` target has a different meaning from `visualize architecture` (a diagram) — it picks whatever `doctor`'s diagnosis says is missing (docker/compose/nginx/ci) and generates all of it at once.

Generated files use exactly the filenames/content `internal/analyzer` recognizes, so `doctor`/`scan` immediately detect them as "present" right after generation (`Dockerfile` gets a `HEALTHCHECK`, `docker-compose.yml` gets the required env vars per DB + a `healthcheck` + `depends_on: condition: service_healthy`, k8s manifests get `resources` requests/limits, CI gets dependency caching + minimal `permissions`). Every generated file includes a "next steps" banner localized via `config --lang` (en/ko).

---

### export

| Item | Content |
|---|---|
| Command | `infra-doctor export [path] [--force] [--dry-run]` |
| Status | ✅ Done (PR #61/#62, issue #39) |
| Uses AI | N |
| Output | An `infra-doctor/` directory (report.md, architecture.md/.mmd, deployment-flow.md, recommendations.md, docker/, kubernetes/, github/) |
| Owner | jade |

**Description**: Saves the combined output of `scan`+`doctor`+`visualize architecture`+`visualize flow`+`generate` into a single `infra-doctor/` directory at the project root in one run, without needing to run each command separately (`internal/export`). No login required — everything exported is deterministic (no AI-generated prose).

- `report.md`: combines the scan results (Framework/Dependencies/Database/Infrastructure/CI-CD/Profiles) and the doctor results (score + Infrastructure Check) into one report
- `architecture.md`/`architecture.mmd`, `deployment-flow.md`: `internal/visualize` rendered as Markdown/Mermaid
- `recommendations.md`: the doctor diagnoses' `Fix` text (or `Message` if `Fix` is empty), sorted and listed
- `docker/`, `kubernetes/`, `github/`: generated by calling `internal/generate`'s docker/compose/k8s/ci Generators directly (filenames are rearranged to fit export's directory layout — e.g. `k8s/deployment.yml` → `kubernetes/deployment.yaml`)

`--force`/`--dry-run` mean the same as in `generate`. Works identically for both Gradle and Maven projects.

---

### about / version / update / donate

| Command | Status | Notes |
|---|---|---|
| `infra-doctor about` | ✅ Done | Matches the spec |
| `infra-doctor version` | ✅ Done | Based on `runtime/debug.ReadBuildInfo` |
| `infra-doctor update` | ✅ Done | Real semver comparison logic + GitHub Releases API integration (issue #41, PR #66, `internal/utils/github.go`). Since the `v0.2.0` official release (PR #82, develop → main), `GET /releases/latest` responds correctly — confirmed `Latest: v0.2.0` in practice. A Go pseudo-version or pre-release tag (e.g. from `go install ...@branch` instead of a tag) now falls back to "dev" the same way `(devel)`/`unknown` do (issue #87, PR #95). |
| `infra-doctor donate` | ✅ Done | Shows a Ko-fi (`https://ko-fi.com/shellwe`) link. Chose Ko-fi over Buy Me a Coffee since the latter doesn't support South Korea (issue #90, PR #94) |

---

### init

| Item | Content |
|---|---|
| Command | `infra-doctor init` |
| Status | ✅ Done |
| Uses AI | N |

**Description**: Confirms the current directory is a Spring Boot project and creates `.infra-doctor/config.yaml` + `.infra-doctor/.gitignore`. Spring Boot detection uses `internal/project/detector.go`'s `DetectSpringBoot`; if the current directory isn't a Spring Boot project, it fails with a clear error.

---

### help / error

| Item | Content |
|---|---|
| Command | `infra-doctor help`, (on an invalid command) |
| Status | ✅ Done |
| Uses AI | N |

**Description**: Suggesting help on an invalid command is handled by Cobra by default, so there's no separate implementation for it. `help` is registered via `rootCmd.SetHelpCommand` (`Use: "help"` is a reserved name for Cobra's own built-in help command, so registering it via `AddCommand` would collide).

---

### config

| Item | Content |
|---|---|
| Command | `infra-doctor config [--lang en\|ko]` |
| Status | ✅ Done (language switching is only the mechanism so far — see below) |
| Uses AI | N |

**Description**: Displays the current LLM/Language/Output/Auto Export settings. Actually reads and writes `~/.infra-doctor/config.json` (the same file `login` uses; `Language`/`OutputFormat`/`AutoExport` fields were added to `Credentials` in `internal/ai/credentials.go`). The `--lang` flag changes and persists the language. `Output`/`Auto Export` are currently read-only — there's no flag yet to change their values (issue #89).

**Scope of language switching**: The `internal/i18n` package (map-based key lookup, no external dependency) is applied to the fixed text of every command — `scan`/`doctor`/`login`/`explain`/`recommend`/`init`/`update`/`version`/`help`/`about`/`donate`/`config`/`generate`. `explain`/`recommend`'s AI responses are also genuinely generated in Korean via a prompt language instruction. The box UI (`internal/ui/box.go`) also handles East Asian Width via `DisplayWidth`, since Korean characters take up 2 columns in a terminal. `doctor`'s rule YAML (`message`/`reason`/`fix`) isn't translated yet — that needs a separate bilingual YAML field design (issue #91).

---

### Next Step

| Item | Content |
|---|---|
| Status | ✅ Done (doctor + recommend, issue #40) |
| Uses AI | N |

**Description**: Unified as `internal/nextstep.Suggest(info, wantsKubernetes)`. It looks at what's missing among Dockerfile/compose/k8s/CI workflow and suggests an `infra-doctor generate <target>` command; if nothing is missing, it suggests `infra-doctor doctor`. `doctor` shows it as a Next Step section below the existing Recommendation (per-rule Fix text); `recommend` also calls this same package for its Next Step section. Extending this to `scan`/`explain` is out of scope.

---

Bug reports and other TODOs are tracked on [GitHub Issues](https://github.com/Team-Shell-We/infra-doctor/issues).
