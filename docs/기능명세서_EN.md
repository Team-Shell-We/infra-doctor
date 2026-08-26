# Infra Doctor Feature Specification

**English** | [한국어](기능명세서.md)

> Last updated against: `develop` (PR #45 full i18n coverage, #46 explain topic-not-adopted clarification, #49 visualize, #50/#51 generate, #53 recommend copy cleanup, #54 recommend complexity signal expansion, #55 Next Step layer, #56 doctor CI gate, #64 Maven support, #66 update release API, #67 Gradle Kotlin DSL detection fix, #70 Maven infrastructure detection fix, #77 help/timeout/recommend label/visualize DB node cleanup, #81 scaffold/dead code cleanup, #82 v0.2.0 release, #84 duplicate logic consolidation, #85 comment style unification, #79 CONTRIBUTING split)

## Project Description

An AI-powered CLI for operational diagnosis and infrastructure guidance. An open-source tool that analyzes a Spring Boot project, visualizes its current operational structure, diagnoses what's missing for production, and auto-generates Docker/Kubernetes/CI-CD config files plus improvement guidance.

## Pain Point

- Backend developers can build applications but struggle when it comes to setting up the operational environment (when to use Docker, where Nginx fits, why Redis is needed, when to adopt Kubernetes, etc.)
- There's plenty of information online, but getting an answer specific to *your* project usually means pasting your whole codebase into an AI and asking — a lot of manual effort.
- Infra Doctor aims to analyze the project itself, explain the current structure from an operations standpoint, and provide both a direction for improvement and config files you can actually apply.

## Core Design Principle (established during implementation)

Not in the original spec, but confirmed empirically while building `explain`/`doctor` and then applied to every subsequent AI command.

> **AI never judges facts — it only explains facts the code has already determined.**

- `doctor` was built without any AI at all, 100% rule-based (YAML rules) — noted again below wherever this differs from the original spec.
- `explain` initially let the AI freely generate "Current Status" (whether related files exist). In testing, we confirmed it would plausibly invent filenames it had never scanned (e.g. `Dockerfile.dev`). We then restructured it so "fact-finding" (what exists) is entirely decided by Go code based on `project.Info`, and the AI only handles the natural-language "why it matters" explanation.
- `recommend` applies the same principle from the start — the deployment strategy recommendation itself (Docker Compose vs. Kubernetes) is decided deterministically by code, and the AI only writes out the reasoning (Reason) in prose.
- Thanks to this principle, the AI commands have reproducible core information (checklists, recommended values, Next Step, etc.), and the AI's role is clearly scoped to "explaining things naturally."

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

Build file (`build.gradle`/`pom.xml`) discovery only happens in the given `path` directory, while Docker/Compose/K8s/GitHub Actions/Profile file discovery recurses into subdirectories, excluding `.git`/`.gradle`/`.idea`/`build`/`node_modules`/`target`/`.infra-doctor` (`internal/analyzer/walk.go`). This applies identically to Gradle and Maven (PR #70) — however, dependency/database/framework detection itself only reads the single root build file for both build tools, so a multi-module project whose real dependencies are declared in a child module can be missed (see "Known Issues" below).

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

**Description**: Recommends a deployment strategy (Docker Compose vs. Kubernetes) based on infrastructure complexity signals detected in the project (Kafka/relational DB/Redis/multiple CI workflows count). Requires `login`.

Previously described as "analyzing project scale," but there's actually no way to see true "scale" indicators like traffic, team size, or microservice count (only a single repo is scanned). Counting scannable static signals is what the code actually does, so the wording was updated to match (issue #44).

**6 signals (expanded in issue #52, `complexityThreshold = 4`)**: Kafka usage, relational DB usage, Redis usage, 2+ CI workflows, more than 20 API endpoints (`internal/analyzer/api.go` counts `@RestController`/`@GetMapping`-style annotations in `.java` sources), multi-module structure (for Gradle, `internal/analyzer/modules.go` counts `include` statements in `settings.gradle`; for Maven, `AnalyzeMavenModules` counts `<modules>` entries in `pom.xml` — PR #70). A project already using Kubernetes is kept regardless of these signals, but if replicas are 2 or more, an "already scaling" rationale is added to the Reason (`internal/analyzer/infrastructure.go` also parses the `replicas:` value from manifests).

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
| `infra-doctor update` | ✅ Done | Real semver comparison logic + GitHub Releases API integration (issue #41, PR #66, `internal/utils/github.go`). Since the `v0.2.0` official release (PR #82, develop → main), `GET /releases/latest` responds correctly — confirmed `Latest: v0.2.0` in practice. However, the version parser still hard-fails on the Go pseudo-version produced by installing via `go install ...@branch` instead of a tag (see "Known Issues" below). |
| `infra-doctor donate` | ✅ Done | Shows a Ko-fi (`https://ko-fi.com/shellwe`) link. Chose Ko-fi over Buy Me a Coffee since the latter doesn't support South Korea |

---

### init

| Item | Content |
|---|---|
| Command | `infra-doctor init` |
| Status | ✅ Done |
| Uses AI | N |

**Description**: Confirms the current directory is a Spring Boot project and creates `.infra-doctor/config.yaml` + `.infra-doctor/.gitignore`. Spring Boot detection uses `internal/project/detector.go`'s `DetectSpringBoot` (previously `cmd/init.go` had its own reimplementation plus a recursive `examples/` fallback; both were removed — if the current directory isn't a Spring Boot project, it now fails with a clear error).

---

### help / error

| Item | Content |
|---|---|
| Command | `infra-doctor help`, (on an invalid command) |
| Status | ✅ Done |
| Uses AI | N |

**Description**: `error` (and `internal/utils/error_util.go`) was removed entirely — Cobra provides "suggest help on an invalid command" by default. `help` is registered via `rootCmd.SetHelpCommand` instead of `rootCmd.AddCommand`, which also resolves a naming collision with Cobra's own built-in `help` (`Use: "help"` is a reserved name, so registering via `AddCommand` would duplicate it).

---

### config

| Item | Content |
|---|---|
| Command | `infra-doctor config [--lang en\|ko]` |
| Status | ✅ Done (language switching is only the mechanism so far — see below) |
| Uses AI | N |

**Description**: Displays the current LLM/Language/Output/Auto Export settings. Actually reads and writes `~/.infra-doctor/config.json` (the same file `login` uses; `Language`/`OutputFormat`/`AutoExport` fields were added to `Credentials` in `internal/ai/credentials.go`). The `--lang` flag changes and persists the language.

**Scope of language switching**: The `internal/i18n` package (map-based key lookup, no external dependency) is applied to the fixed text of every command — `scan`/`doctor`/`login`/`explain`/`recommend`/`init`/`update`/`version`/`help`/`about`/`donate`/`config`/`generate` (issue #36, PR #45). `explain`/`recommend`'s AI responses are also genuinely generated in Korean via a prompt language instruction. The box UI (`internal/ui/box.go`) also handles East Asian Width via `DisplayWidth`, since Korean characters take up 2 columns in a terminal. The only thing left untranslated is `doctor`'s rule YAML (`message`/`reason`/`fix`) — that needs a separate bilingual YAML field design.

---

### Next Step

| Item | Content |
|---|---|
| Status | ✅ Done (doctor + recommend, issue #40) |
| Uses AI | N |

**Description**: Unified as `internal/nextstep.Suggest(info, wantsKubernetes)`. It looks at what's missing among Dockerfile/compose/k8s/CI workflow and suggests an `infra-doctor generate <target>` command; if nothing is missing, it suggests `infra-doctor doctor`. `doctor` shows it as a new Next Step section below the existing Recommendation (per-rule Fix text); `recommend` replaced its existing Next Step section with a call into this same package. Extending this to `scan`/`explain` was left out of scope.

A pre-existing bug was also fixed during the migration — `recommend` used to suggest `generate docker` when it recommended Docker Compose and no compose file existed, but the target that actually creates `docker-compose.yml` is `compose` (`docker` only creates the Dockerfile). The case where the Dockerfile itself is missing, which wasn't checked at all before, was also added.

---

## Known Issues

### Resolved (for reference, issue #31)
- ~~`FindBuildFile` recursed into subdirectories and could pick up an unrelated project like `examples/`~~ → PR #30
- ~~`explain`/`recommend` were hardcoded to the current directory~~ → PR #30, added the `[path]` argument
- ~~`help`/`error` command name collision~~ → removed `error.go` and switched `help.go` to `SetHelpCommand` (also discovered and fixed the collision with Cobra's own built-in help)
- ~~`init`'s `examples/` fallback logic~~ → removed
- ~~Triple-duplicated Spring Boot detection logic~~ → `cmd/init.go` now uses `project.DetectSpringBoot` (previously unused code); `internal/analyzer/finder.go` serves a different purpose (parsing build files) and was kept separate.
- ~~`config` command's content didn't match the spec~~ → now shows LLM/Language/Output/Auto Export and `--lang` actually changes/persists a value. Only `--lang` actually changes a value; Output/Auto Export are still read-only.
- ~~`root.go` was still the raw Cobra scaffold~~ → replaced the Short/Long description, removed the unused `--toggle` flag
- ~~`scan`/`doctor`'s infrastructure file search walked the entire tree recursively~~ → `internal/analyzer/walk.go` now excludes `.git`/`.gradle`/`.idea`/`build`/`node_modules`/`target`/`.infra-doctor`

### Resolved (for reference, issues #36/#38/#37/#44)
- ~~Language switching only applied to the `config` command itself~~ → extended to every command, box UI now handles East Asian Width (PR #45)
- ~~`visualize architecture`/`visualize flow` unimplemented~~ → implemented as the `internal/visualize` package (PR #49)
- ~~`generate <target>` unimplemented~~ → implemented as `internal/generate` (docker/compose/nginx/k8s/ci + the unified `architecture` target), integrated with doctor's diagnosis, hardened generated output to a genuinely deployable level (HEALTHCHECK/depends_on condition/resources limits/CI caching, etc.) + localized banners (PR #50/#51)
- ~~`explain` described a not-yet-adopted topic as if it were already in use~~ → shows "Current Status" first + a not-yet-adopted banner + subjunctive-phrasing prompt instruction (PR #46)
- ~~`recommend`'s description overstated the actual logic (complexity signal counting) as "scale analysis"~~ → wording cleaned up (issue #44, final item)

### Resolved (for reference, issues #39/#41/#43)
- ~~`export` unimplemented~~ → implemented as `internal/export`, saves report+diagrams+generated files into `infra-doctor/` in one run (PR #61/#62)
- ~~Maven project analysis unsupported~~ → implemented pom.xml parsing in `internal/analyzer/maven.go` (Spring Boot version resolved via parent/dependency/plugin order, `${property}` substitution, Java version resolved via `sourceCompatibility`/`maven.compiler.*`/toolchain order) (PR #64)
- ~~`update`'s `latestVersion` was hardcoded~~ → integrated the GitHub Releases API (PR #66)

### Resolved (for reference, issues #65/#73/#74/#75/#76)
- ~~Gradle Kotlin DSL (`build.gradle.kts`) projects were misdetected as non-Spring-Boot~~ → the version regex now matches Kotlin DSL syntax too, added `sourceCompatibility`-based Java version detection (PR #67)
- ~~Infrastructure/CI/profile detection was always failing for Maven projects~~ → moved the `AnalyzeGitHub`/`FindProfiles`/`AnalyzeInfrastructure`/`AnalyzeAPI` calls outside the switch so both Gradle and Maven run them, added Maven multi-module counting and `dependencyManagement` (BOM) parsing, also fixed `export`'s missing nginx generation (PR #70)
- ~~`help <command>` always failed~~ → implemented so it actually shows the given command's help (issue #73, PR #77)
- ~~OpenAI requests were cut off at 30 seconds instead of 60~~ → removed the fixed timeout, now controlled solely by the context timeout (issue #74, PR #77)
- ~~`recommend` mislabeled Maven multi-module projects as "Gradle"~~ → removed the hardcoded build-tool name (issue #75, PR #77)
- ~~A fake "Unknown" node appeared in the architecture diagram for projects with no database~~ → added a filter in `visualize/builder.go` (issue #76, PR #77)
- ~~`update` always failed with a GitHub API 404~~ → the `v0.2.0` official release (PR #82, develop → main) made `/releases/latest` respond correctly. Resolved with no code change (the root cause was simply that no official release existed yet).
- ~~`main.go`'s cobra-cli scaffold placeholder, dead code (`AnalyzeDatabase`/`HasGradle`/`nodeMap`/`nodeLabel`)~~ → cleaned up (issue #80, PR #81)
- ~~Duplicated file-status logic in `generate`/`export`, duplicated language-mapping logic in `explain`/`recommend`~~ → consolidated into `generate.Result.StatusOf()`/`internal/ai/language.go` (PR #84)
- ~~Comments were a mix of English and Korean~~ → unified to Korean throughout, applied a clarity checklist, cleaned up 10 files that weren't gofmt-clean + added a `gofmt -l` check to CI (PR #85)
- ~~README had grown long enough that even someone who just wanted the contributing guide had to skim the whole thing~~ → split `CONTRIBUTING.md` out to the repo root (GitHub surfaces it automatically when opening an issue/PR), kept README focused on usage (issue #79)

### Still Open (in priority order)

1. **`update`'s version parser can't handle pseudo-versions** — hard-fails on the Go pseudo-version or pre-release tag you get from installing via `go install ...@branch` instead of an exact tag. Should be folded into the same graceful "dev" fallback as `(devel)`/`unknown`.
2. **Dependency/database detection only reads the root build file** — neither Gradle nor Maven descends into child modules. A multi-module project whose real dependencies live in a child module can have its JPA/Kafka/DB detection missed. A pre-existing limitation, symmetric across both build tools.
3. **GitHub Actions shorthand triggers (`on: push`, `on: [push, pull_request]`) parse as zero triggers** — only the multi-line mapping form is supported. Fails silently with no error.
4. **`config` displays settings it can't actually change** — shows "Output Format"/"Auto Export" as if they were configurable, but no such flag exists, so they can never change. Only `--lang` actually works.
5. **`doctor`'s rule YAML has no translation support** — needs a separate bilingual YAML field design.
6. **`ui.Wrap` doesn't force-break a single word longer than the box** — a long unbroken token in AI-generated text (e.g. a URL) can break the box border alignment.
7. **Test coverage gaps** — most of `cmd/*.go` (only doctor and generate have tests), `internal/utils`, and `internal/ai/openai` have no dedicated tests.

### Not Yet Implemented

- (None — `export` (#39), Maven support (#43), and `update`'s latest-version lookup (#41) are all implemented now. What remains is only the bugs/gaps listed under "Still Open" above.)
