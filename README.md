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
| 📄 Export | `export` | Export analysis results as Markdown, Mermaid, or report files. | 📅 |

> **Status**
>
> - ✅ Available
> - 🚧 In Progress
> - 📅 Planned

`explain`/`recommend` require `login` first. Their AI-generated sections are always grounded in facts the scanner already verified deterministically — the AI never invents which files exist, it only explains why they matter (see [docs/기능명세서.md](docs/기능명세서.md) for the full design rationale).

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

> `config --lang ko` currently switches only `config`'s own output to Korean — translating the rest of the CLI's output is still in progress.

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

Scan the current project (or a specific path).

```bash
infra-doctor scan
infra-doctor scan ~/workspace/my-project
```

Diagnose deployment readiness.

```bash
infra-doctor doctor ~/workspace/my-project
```

Log in once to enable AI-powered commands (OpenAI API key).

```bash
infra-doctor login
```

Explain a technology in the context of your project.

```bash
infra-doctor explain docker ~/workspace/my-project
```

Get a deployment strategy recommendation.

```bash
infra-doctor recommend ~/workspace/my-project
```

View or change CLI configuration.

```bash
infra-doctor config
infra-doctor config --lang ko
```

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
- [x] Framework Detection
- [x] Database Detection
- [x] Docker Detection
- [x] Profile Detection
- [x] GitHub Actions Analysis
- [x] Docker Compose Analysis
- [x] Deployment Doctor
- [x] AI-powered Explain
- [x] AI-powered Recommendation
- [ ] Infrastructure Visualization
- [ ] Configuration Generator
- [ ] Export Reports
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