# 🩺 Infra Doctor

> **AI-powered Infrastructure Analysis CLI for Spring Boot Projects**

Infra Doctor is a CLI tool that analyzes the current infrastructure of Spring Boot projects and helps developers understand, diagnose, visualize, and improve their deployment environment.

It scans project configurations, detects infrastructure components, and provides actionable insights for modern DevOps practices.

---

## ✨ Features

| Category | Command | Description | Status |
| :------ | :------ | :---------- | :----: |
| 🔍 Analysis | `scan` | Scan the project and collect framework, database, Docker, CI/CD, profile, and infrastructure information. | 🚧 |
| 🩺 Diagnosis | `doctor` | Diagnose deployment readiness and suggest infrastructure improvements. | 📅 |
| 🏗 Visualization | `visualize architecture` | Generate an architecture diagram of the current infrastructure. | 📅 |
| 🏗 Visualization | `visualize flow` | Visualize the build and deployment workflow. | 📅 |
| 🤖 AI | `explain <topic>` | Explain how a specific infrastructure technology is used within the current project. | 📅 |
| 🤖 AI | `recommend` | Recommend an infrastructure strategy based on project analysis. | 📅 |
| ⚙️ Generator | `generate <target>` | Generate infrastructure configuration files from scan results. | 📅 |
| 📄 Export | `export` | Export analysis results as Markdown, Mermaid, or report files. | 📅 |

> **Status**
>
> - ✅ Available
> - 🚧 In Progress
> - 📅 Planned

---

## 🔧 Utilities

| Command | Description | Status |
| :------ | :---------- | :----: |
| `init` | Initialize Infra Doctor configuration for the current project. | 📅 |
| `config` | View or update CLI configuration. | 📅 |
| `login` | Authenticate to enable AI-powered features. | 📅 |
| `update` | Check for the latest CLI version. | 📅 |
| `version` | Display the current CLI version. | 🚧 |
| `help` | Show available commands and usage information. | ✅ |
| `about` | Display project information and credits. | 📅 |
| `donate` | Support the Infra Doctor project. | 📅 |

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

Analyze the current project.

```bash
infra-doctor scan
```

Analyze a specific project.

```bash
infra-doctor scan <project-path>
```

Example

```bash
infra-doctor scan ~/workspace/my-project
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

The following directories can be excluded:

```text
.git/
.gradle/
.idea/
build/
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
- [ ] GitHub Actions Analysis
- [ ] Docker Compose Analysis
- [ ] Infrastructure Visualization
- [ ] SSH / EC2 Analysis
- [ ] Kubernetes Analysis
- [ ] Deployment Doctor
- [ ] Configuration Generator
- [ ] AI Recommendation

---

## 🤝 Contributing

Contributions are always welcome.

If you'd like to contribute, please read the contribution guide before opening an Issue or Pull Request.

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
```

The test project is intended **only for local development and testing**. Do not include sensitive source code or production configuration files in the repository.

---

## 📄 License

This project is licensed under the MIT License.