# Contributing to Infra Doctor

**English** | [한국어](CONTRIBUTING.md)

Thanks for considering a contribution. This guide covers everything you need to get set up, test your changes, and contribute a specific kind of change (a diagnostic rule or a translation).

---

## 1. Configure Git

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

## 2. Prepare a Test Project

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

## 3. Test Your Changes

Example projects are **not included** in this repository — use your own test project (see step 2) or any Spring Boot project on disk.

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

Run the full test suite and formatting check before opening a PR:

```bash
go build ./...
go vet ./...
gofmt -l .
go test ./...
```

---

## 4. Adding a Diagnostic Rule

`doctor` checks are defined as YAML, not Go code. See [internal/doctor/rules/README.md](internal/doctor/rules/README.md) for the rule schema and PR checklist.

---

## 5. Adding a Translation

CLI output supports English/Korean via `internal/i18n`. See [internal/i18n/README.md](internal/i18n/README.md) for how to add a new string and which strings should stay untranslated.

---

## Opening a Pull Request

Use the PR template — it asks for the linked issue, a summary, what you implemented, anything you got stuck on, and anything reviewers should know. Commit messages must match the `<type>(#issue): <summary>` format (`feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `ci`); the `commit-msg` hook from step 1 enforces this locally.
