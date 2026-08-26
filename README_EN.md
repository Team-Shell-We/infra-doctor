# 🩺 Infra Doctor

**English** | [한국어](README.md)

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

`explain`/`recommend` require `login` first. Their AI-generated sections are always grounded in facts the scanner already verified deterministically — the AI never invents which files exist, it only explains why they matter (see [docs/기능명세서_EN.md](docs/기능명세서_EN.md) for the full design rationale).

Both Gradle (`build.gradle`/`build.gradle.kts`) and Maven (`pom.xml`) projects are detected as a build tool. **Known limitation:** for Maven projects, only framework/dependency/database detection currently runs — Docker/Compose/Kubernetes/CI/profile detection is Gradle-only right now, so `doctor`/`recommend`/`export` will under-report readiness on an otherwise well-configured Maven project. Tracked as a bug to fix; see [docs/기능명세서_EN.md](docs/기능명세서_EN.md) for details.

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

| Command | Key flags | Description |
| :--- | :--- | :--- |
| `scan [path]` | – | Detected framework/dependencies/database/infrastructure/CI/profiles |
| `doctor [path]` | `--json`, `--fail-under <score>` | Readiness score + fix-it list. Usable as a CI gate |
| `login` | – | Register an OpenAI API key (prerequisite for `explain`/`recommend`) |
| `explain <topic> [path]` | – | topic: compose·container·docker·github-actions·image·k8s·nginx·postgres·rds·redis |
| `recommend [path]` | – | Docker Compose vs. Kubernetes, with reasons |
| `generate <target> [path]` | `-f/--force`, `--dry-run`, `-o/--output-dir`, `--config` | target: docker·compose·nginx·ci·k8s·architecture |
| `visualize architecture\|flow [path]` | `--format`, `--output` | ascii/mermaid/markdown diagram |
| `export [path]` | `-f/--force`, `--dry-run` | All of the above into one `infra-doctor/` directory |
| `config [--lang en\|ko]` | – | View or change CLI configuration |
| `init` | – | Create `.infra-doctor/config.yaml` |

`version`, `update`, `help`, `about`, and `donate` take no arguments and need no login. For full command behavior, output examples, and the `generate`/`init` config file schema, see [docs/기능명세서_EN.md](docs/기능명세서_EN.md); for the complete flag list, run `infra-doctor <command> --help`.

---

## 📂 Project Structure

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

Contributions are always welcome. See [CONTRIBUTING_EN.md](CONTRIBUTING_EN.md) for setup, how to test your changes, and how to add a diagnostic rule or a translation.

---

## 📄 License

This project is licensed under the MIT License.
