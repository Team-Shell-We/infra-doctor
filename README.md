# 🩺 Infra Doctor

> **AI-powered Infrastructure Analysis CLI for Spring Boot Projects**

Infra Doctor is a CLI tool that analyzes the current infrastructure of Spring Boot projects and helps developers understand, diagnose, visualize, and improve their deployment environment.

It scans project configurations, detects infrastructure components, and provides actionable insights for modern DevOps practices.

---

## ✨ Features

| Category | Command | Description | Status |
| :------ | :------ | :---------- | :----: |
| 🔍 Analysis | `scan` | Scan the project and collect framework, database, Docker, CI/CD, profile, and infrastructure information. | ✅ |
| 🩺 Diagnosis | `doctor` | Diagnose deployment readiness (0–100%) against a rule-based checklist and suggest infrastructure improvements. | ✅ |
| 🤖 AI | `explain <topic>` | Explain how a specific infrastructure technology is used within the current project. | ✅ |
| 🤖 AI | `recommend` | Recommend a deployment strategy (Docker Compose vs. Kubernetes) based on detected infrastructure complexity. | ✅ |
| 🏗 Visualization | `visualize architecture` | Generate an architecture diagram of the current infrastructure. | ✅ |
| 🏗 Visualization | `visualize flow` | Visualize the build and deployment workflow. | ✅ |
| ⚙️ Generator | `generate <target>` | Generate infrastructure configuration files from scan results. | ✅ |
| 📄 Export | `export` | Export the full analysis (report, diagrams, generated configs) into one directory. | ✅ |

> **Status**
>
> - ✅ Available
> - 🚧 In Progress
> - 📅 Planned

`explain`/`recommend` require `login` first. Their AI-generated sections are always grounded in facts the scanner already verified deterministically — the AI never invents which files exist, it only explains why they matter (see [docs/기능명세서.md](docs/기능명세서.md) for the full design rationale).

Both Gradle (`build.gradle`/`build.gradle.kts`) and Maven (`pom.xml`) projects are detected as a build tool. **Known limitation:** for Maven projects, only framework/dependency/database detection currently runs — Docker/Compose/Kubernetes/CI/profile detection is Gradle-only right now, so `doctor`/`recommend`/`export` will under-report readiness on an otherwise well-configured Maven project. Tracked as a bug to fix; see [docs/기능명세서.md](docs/기능명세서.md) for details.

### Using `doctor` as a CI gate

`doctor --json` prints the same result as machine-readable JSON instead of the box UI, and `--fail-under <score>` exits non-zero when the readiness score is below the threshold — combine them to fail a CI pipeline on infrastructure regressions:

```bash
infra-doctor doctor --fail-under 70
infra-doctor doctor --json > doctor-result.json
```

---

## 🔧 Utilities

| Command | Description | Status |
| :------ | :---------- | :----: |
| `init` | Initialize Infra Doctor configuration for the current project. | ✅ |
| `config` | View or update CLI configuration (LLM provider, output language, output format, auto-export). | ✅ |
| `login` | Authenticate with an OpenAI API key to enable AI-powered features. | ✅ |
| `update` | Check for the latest CLI version. | ✅ |
| `version` | Display the current CLI version. | ✅ |
| `help` | Show available commands and usage information. | ✅ |
| `about` | Display project information and credits. | ✅ |
| `donate` | Support the Infra Doctor project. | ✅ |

> `config --lang ko` switches every command's fixed CLI text (labels, headers, messages) to Korean, and `explain`/`recommend` generate their AI text in that language too. The one thing still English-only is `doctor`'s diagnostic rule text (the `message`/`reason`/`fix` in `internal/doctor/rules/*.yaml`) — translating those needs a bilingual YAML field design that hasn't been done yet.

---

## 🚀 Installation

### Install via Go

```bash
go install github.com/Team-Shell-We/infra-doctor@latest
```

### Build from Source

```bash
git clone https://github.com/Team-Shell-We/infra-doctor.git

cd infra-doctor

go build -o infra-doctor
```

---

## 📖 Usage

### Typical workflow

A new Spring Boot project usually goes through these commands in order:

```bash
cd ~/workspace/my-project

infra-doctor init                 # create .infra-doctor/config.yaml
infra-doctor scan                 # see what's detected
infra-doctor doctor                # see what's missing, and get a readiness score
infra-doctor login                 # one-time, needed for explain/recommend
infra-doctor recommend             # Docker Compose vs Kubernetes, with reasons
infra-doctor generate compose      # write the files recommend pointed to
infra-doctor doctor                # re-run to confirm the score went up
infra-doctor visualize architecture
```

`[path]` is optional on every command below — it defaults to the current directory.

### `scan [path]`

Prints detected framework, dependencies, database, infrastructure, CI/CD, and profiles. No flags, no login required.

```bash
infra-doctor scan
infra-doctor scan ~/workspace/my-project
```

### `doctor [path]`

Diagnoses deployment readiness as a 0–100% score plus a fix-it list. No login required.

| Flag | Default | Description |
| :--- | :--- | :--- |
| `--json` | `false` | Print the result as JSON instead of the box UI |
| `--fail-under <score>` | *(unset)* | Exit with code 1 if the score is below this threshold |

```bash
infra-doctor doctor ~/workspace/my-project

# CI gate: fail the build if readiness drops below 70
infra-doctor doctor --fail-under 70

# machine-readable output for other tooling
infra-doctor doctor --json > doctor-result.json
```

### `login`

Interactive prompt to store an OpenAI API key at `~/.infra-doctor/config.json` (required once, before `explain` or `recommend`).

```bash
infra-doctor login
```

### `explain <topic> [path]`

Explains what a technology means for *this* project — grounded in what `scan` actually found, not a generic tutorial. Requires `login`.

Valid topics: `compose`, `container`, `docker`, `github-actions`, `image`, `k8s`, `nginx`, `postgres`, `rds`, `redis`.

```bash
infra-doctor explain docker ~/workspace/my-project
infra-doctor explain k8s        # works even if you haven't adopted k8s yet — says so explicitly
```

### `recommend [path]`

Recommends Docker Compose vs. Kubernetes based on detected complexity signals (Kafka, relational DB, Redis, CI workflow count, API endpoint count, multi-module Gradle), and suggests the `generate` command(s) to close the gap. Requires `login`.

```bash
infra-doctor recommend ~/workspace/my-project
```

### `generate <target> [path]`

Writes real infrastructure config files based on the scan + `doctor` diagnosis, not generic boilerplate — port, health path, and database env vars are filled in from what was actually detected.

| Target | Produces |
| :--- | :--- |
| `docker` | `Dockerfile` |
| `compose` | `docker-compose.yml` |
| `nginx` | `nginx.conf` |
| `ci` | `.github/workflows/*.yml` |
| `k8s` | Kubernetes deployment/service/configmap manifests |
| `architecture` | Whichever of the above `doctor` says is currently missing, generated together |

| Flag | Default | Description |
| :--- | :--- | :--- |
| `-f`, `--force` | `false` | Overwrite files that already exist |
| `--dry-run` | `false` | Print what would be written without writing it |
| `-o`, `--output-dir` | *(project root)* | Directory to write generated files into |
| `--config` | `.infra-doctor/config.yaml` | Path to an alternate config file |

```bash
infra-doctor generate docker ~/workspace/my-project
infra-doctor generate compose --dry-run
infra-doctor generate k8s --force
```

Customize the generated content by adding a `generate:` section to `.infra-doctor/config.yaml` (created by `init`):

```yaml
project:
  name: my-service

generate:
  applicationPort: 8080
  healthPath: /actuator/health
  serviceName: my-service
  dockerImage: my-service:latest
  namespace: default
  replicas: 2

output:
  directory: .infra-doctor/generated   # set by `init`; omit to write into the project root instead
  overwrite: false
```

### `visualize architecture [path]` / `visualize flow [path]`

`architecture` diagrams the current runtime infrastructure; `flow` diagrams the build-to-deploy pipeline (GitHub Actions workflow is searched for in the current or a parent directory). Neither needs login — both are drawn straight from `scan` results.

| Flag | Default | Description |
| :--- | :--- | :--- |
| `--format` | `ascii` | `ascii`, `mermaid`, or `markdown` |
| `--output` | *(stdout)* | Write to a file instead of printing |

```bash
infra-doctor visualize architecture
infra-doctor visualize flow --format mermaid --output flow.md
```

### `config [--lang en\|ko]`

Shows the current LLM provider, language, output format, and auto-export setting. `--lang` switches and persists the CLI's output language.

```bash
infra-doctor config
infra-doctor config --lang ko
```

### `export [path]`

Writes the full analysis — report, architecture/flow diagrams, and every `generate` target — into one `infra-doctor/` directory in the project, so you don't have to run each command separately. No login required (the report and diagrams are all deterministic; there's no AI-generated content in `export`'s output).

| Flag | Default | Description |
| :--- | :--- | :--- |
| `-f`, `--force` | `false` | Overwrite files that already exist |
| `--dry-run` | `false` | Print what would be written without writing it |

```bash
infra-doctor export ~/workspace/my-project
```

Produces:
```text
infra-doctor/
├── report.md              # scan + doctor results
├── architecture.md
├── architecture.mmd
├── deployment-flow.md
├── recommendations.md     # doctor's fix list, sorted
├── docker/
│   ├── Dockerfile
│   ├── .dockerignore
│   └── docker-compose.yml
├── kubernetes/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── configmap.yaml
└── github/
    └── ci.yml
```

### `init`

Detects that the current directory is a Spring Boot project and creates `.infra-doctor/config.yaml` + `.infra-doctor/.gitignore`. Run this once per project, before `generate`, if you want to customize output.

```bash
infra-doctor init
```

### Everything else

`version`, `update`, `help`, `about`, and `donate` take no arguments and need no login — run any of them with `--help` for details.

---

## 🧪 Testing

Example projects are **not included** in this repository.

To test Infra Doctor, prepare any Spring Boot project and run:

```bash
go run . scan <project-path>
```

Example

```bash
go run . scan ../spring-project
```

For richer analysis, the project may contain:

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

The following directories are automatically excluded from analysis:

```text
.git/
.gradle/
.idea/
build/
node_modules/
target/
.infra-doctor/
```

---

## 📂 Project Structure

```text
infra-doctor
├── .github
├── cmd
├── docs
├── internal
├── scripts
├── go.mod
├── main.go
└── README.md
```

---

## 🛣️ Roadmap

- [x] Project Scanner
- [x] Framework Detection (Gradle; Maven partially — see note above)
- [x] Database Detection
- [x] Docker Detection
- [x] Profile Detection
- [x] GitHub Actions Analysis
- [x] Docker Compose Analysis
- [x] Deployment Doctor
- [x] AI-powered Explain
- [x] AI-powered Recommendation
- [x] Infrastructure Visualization
- [x] Configuration Generator
- [x] Export Reports
- [ ] Kubernetes Deep Analysis
- [ ] SSH / EC2 Analysis
- [ ] Full Multi-language Output

---

## 🤝 Contributing

Contributions are always welcome.

If you'd like to contribute, please read the **Contributor Guide** section below before opening an Issue or Pull Request. For specific kinds of contributions, see:

- [internal/doctor/rules/README.md](internal/doctor/rules/README.md) — adding or changing a diagnostic rule
- [internal/i18n/README.md](internal/i18n/README.md) — adding or changing a translated string

---

## 👩‍💻 Contributor Guide

### 1. Configure Git

After cloning the repository, run the setup script once.

**Windows (PowerShell)**

```powershell
.\scripts\setup.ps1
```

This script automatically configures:

- Git Hooks (`.githooks`)
- Commit Message Template (`.gitmessage`)

You should see:

```text
Configuring Git...

Done!
[OK] Git Hooks configured
[OK] Commit template configured
```

---

### 2. Prepare a Test Project

For development, prepare a Spring Boot project under the `examples` directory.

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

Copy your Spring Boot project into `examples/spring-gradle`.

For example:

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

Before committing, remove the following directories from the copied project:

```text
.git/
.gradle/
.idea/
build/
node_modules/
target/
```

The test project is intended **only for local development and testing**. Do not include sensitive source code or production configuration files in the repository.

---

### 3. Adding a Diagnostic Rule

`doctor` checks are defined as YAML, not Go code. See [internal/doctor/rules/README.md](internal/doctor/rules/README.md) for the rule schema and PR checklist.

---

### 4. Adding a Translation

CLI output supports English/Korean via `internal/i18n`. See [internal/i18n/README.md](internal/i18n/README.md) for how to add a new string and which strings should stay untranslated.

---

## 📄 License

This project is licensed under the MIT License.