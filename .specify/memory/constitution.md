<!--
SYNC IMPACT REPORT
==================
Version change: 1.0.0 → 1.1.0
Bump type: MINOR — new principle added (VII. Clean Architecture & Dependency Direction),
  multiple sections materially expanded with concrete layer rules, interface inventory
  updated, and new packages documented.

Modified principles:
  - I.  Multi-Provider Abstraction: factory moved from internal/llm to cmd/llm.go
  - II. Cobra Double-File: tightened — opts structs MUST have all deps injected,
        no construction or I/O inside Run()
  - III. Interface-Driven: interface inventory updated (added Marshal/Unmarshal,
         Format, RepoRoot, GitConfigReader); narrow interfaces documented

Added sections:
  - Principle VII: Clean Architecture & Dependency Direction (new)
  - Layer Map (new subsection under VII)
  - Package Responsibility Rules (new subsection under VII)
  - internal/dirs as canonical path package (new subsection under Configuration)
  - UI layer contract (new subsection under VII)
  - Narrow interface rule (new subsection under VII)

Removed:
  - Reference to internal/llm/factory.go (factory now in cmd/llm.go)

Templates reviewed:
  - .specify/templates/plan-template.md     ✅ Constitution Check gate updated
  - .specify/templates/spec-template.md     ✅ no changes required
  - .specify/templates/tasks-template.md    ✅ no changes required

Follow-up TODOs:
  - None
-->

# git-msg Constitution

## Core Principles

### I. Multi-Provider Abstraction (MUST)

All LLM integrations MUST be implemented behind the `llm.Provider` interface
(`internal/llm/provider.go`). No command or business-logic layer may import
a provider SDK directly. Supported providers: OpenAI, Anthropic, Gemini, Ollama.

Provider construction (the factory) lives at the application layer in
`cmd/llm.go:NewLLMProvider` — not inside the `internal/llm` package. The
`llm` infrastructure package exports only the interface, four concrete HTTP
implementations, `ErrProviderRequest`, and `FakeProvider`.

**Rationale**: Prevents lock-in to a single LLM vendor, enables seamless
per-run provider overrides, and keeps peer infrastructure packages decoupled
from each other.

### II. Cobra Double-File Separation (MUST)

Every CLI command MUST be split into two files:

- `cmd/<command>.go` — pure orchestration function `Run(ctx, opts) error`.
  Zero Cobra imports. Zero construction logic. Zero direct I/O (no
  `os.Create`, `exec.Command`, spinner calls). All dependencies arrive via
  the options struct; `Run` only coordinates them.
- `cmd/<command>_cobra.go` — `New<Command>Cmd() *cobra.Command`. Binds
  flags, constructs all infrastructure dependencies, calls `Run`. No logic
  beyond wiring.

Injectable I/O behaviours in opts structs:
- `SpinnerFunc func(label string) ui.StopFunc` — nil falls back to `ui.Spinner`
- `ReviewFunc func(message string) (ui.ReviewResult, error)` — nil falls back
  to `ui.Review`

**Rationale**: Every `Run` function is fully testable by calling it directly
with fakes, without cobra, without a TTY, without any subprocess.

### III. Interface-Driven Internal Packages (MUST)

Every internal subsystem MUST expose its behaviour through a Go interface.
Concrete implementations are injected at the call site.

**Current interface inventory**:

| Interface | File | Key methods |
| --- | --- | --- |
| `config.Store` | `internal/config/store.go` | `Load`, `Save`, `Format` |
| `secret.SecretStore` | `internal/secret/store.go` | `Get`, `Set`, `Delete` |
| `prompt.TemplateStore` | `internal/prompt/store.go` | `Get`, `List`, `Save`, `Delete`, `Marshal`, `Unmarshal` |
| `llm.Provider` | `internal/llm/provider.go` | `Generate(ctx, system, user)` |
| `git.Client` | `internal/git/client.go` | `StagedDiff`, `CurrentBranch`, `RecentLog`, `Commit`, `RunConfig`, `RepoRoot` |
| `hook.Manager` | `internal/hook/manager.go` | `Install`, `Uninstall`, `IsInstalled` |
| `hook.GitConfigReader` | `internal/hook/install.go` | `RunConfig` (narrow interface) |

`Marshal`/`Unmarshal` on `TemplateStore` and `Format` on `config.Store`
exist specifically to keep serialisation format knowledge (TOML) inside
the infrastructure layer and out of `cmd/`.

**Rationale**: Enables fake injection in tests, enforces clear boundaries,
and allows future alternative implementations without changing callers.

### IV. Test-First for Business Logic (MUST)

All business logic in `internal/` packages and non-Cobra `cmd/*.go` files
MUST have unit tests written before or alongside implementation. Tests MUST
use interface fakes — no live API calls, no live git processes, no live
keychain access in unit tests.

`go test ./... -race -count=1` MUST pass fully offline.

Fixture diffs live in `testdata/diffs/`. Hook source dispatch
(`ShouldGenerate`) MUST be covered by table-driven tests over all six
`SourceType` values.

**Rationale**: Maintains a fast, reliable test suite and enforces
interface-driven design by demanding fakes exist for every interface.

### V. Secure Credential Handling (MUST)

API keys MUST NOT appear in config files, source code, logs, or CLI output.
The lookup order is strictly:
1. System keychain (via `github.com/99designs/keyring`, service `git-msg`)
2. Environment variable `GIT_MSG_<PROVIDER>_API_KEY`
3. Structured error with user-facing remediation instructions

No other lookup path is permitted. Keychain writes MUST be gated behind
explicit user consent (setup wizard or `config set`). Error messages MUST
be verified to contain no key values.

**Rationale**: Protects credentials from accidental exposure in dotfiles,
shell history, or version control.

### VI. Simplicity & Explicit Non-Goals (MUST respect)

The v1 scope is fixed. The following are explicitly out of scope and MUST
NOT be implemented or designed for in v1:

- Diff chunking for large diffs
- `git add -p` staging integration
- Team/shared prompt server
- Windows support (keychain assumption)
- Internationalization (i18n)

Any proposal to extend scope MUST update this constitution first.

**Rationale**: Keeps the codebase small, idiomatic, and shippable. YAGNI
applies — defer until there is concrete user demand.

### VII. Clean Architecture & Dependency Direction (MUST)

Dependencies flow strictly inward. Outer layers depend on inner layers;
inner layers MUST NOT depend on outer layers or on peer infrastructure
packages.

#### Layer Map

```
main.go
  └─ cmd/                        application layer
       ├─ *_cobra.go              CLI wiring: flag binding, construction, context
       ├─ *.go                    orchestration: Run(ctx, opts), pure logic
       └─ llm.go                  application factory: NewLLMProvider
            └─ internal/          infrastructure & domain
                 ├─ dirs/         XDG path resolution (no deps on other internal/)
                 ├─ config/       domain type + TOML persistence
                 ├─ secret/       credential abstraction + keychain
                 ├─ git/          git subprocess abstraction
                 ├─ llm/          HTTP provider implementations (no config/secret imports)
                 ├─ prompt/       template domain + file store + renderer
                 ├─ hook/         hook lifecycle + source classification
                 └─ ui/           terminal presentation (huh TUI, spinner, editor)
```

#### Package Responsibility Rules

- `internal/llm` MUST NOT import `internal/config` or `internal/secret`.
  Provider construction belongs in `cmd/llm.go`.
- `internal/ui` MUST NOT import `internal/config` or `internal/secret`.
  The wizard returns a plain `SetupResult` struct; persistence is the
  application layer's (`cmd/root.go:EnsureConfig`) responsibility.
- `internal/hook` MUST NOT call `exec.Command("git", ...)` directly.
  It receives a `GitConfigReader` interface injected from `cmd/`.
- `cmd/*.go` (non-cobra logic files) MUST NOT construct infrastructure
  types. Construction belongs exclusively in `*_cobra.go` files.
- No `cmd/` file may import `github.com/spf13/cobra` except `*_cobra.go`.

#### UI Layer Contract

`internal/ui` is a pure presentation layer. It:
- Accepts plain Go values (strings, slices, enums).
- Returns plain Go values or typed error sentinels (`ErrWizardAborted`).
- MUST NOT perform persistence, keychain access, or network calls.
- MUST NOT import any other `internal/` package.

#### Narrow Interface Rule

When a package needs only a subset of another package's behaviour, define
a narrow local interface rather than importing the full type. Example:
`hook.GitConfigReader` exposes only `RunConfig` — `hook` does not depend on
the full `git.Client`.

**Rationale**: Separates concerns, makes each layer independently testable,
prevents circular dependencies, and keeps infrastructure packages reusable
in isolation.

---

## Technical Constraints & Non-Goals

### Language & Runtime

- Implementation language: Go (idiomatic, modern)
- All LLM API calls MUST use stdlib `net/http` only — no vendor SDKs
- All git operations MUST use `os/exec` — no git library bindings
- Config and template files MUST use TOML format
  (`github.com/pelletier/go-toml/v2`)
- Template rendering MUST use `github.com/madstone-tech/ason/pkg` Engine
- TOML serialisation knowledge is confined to `internal/config` and
  `internal/prompt`; `cmd/` MUST NOT import `go-toml` directly

### Configuration & Paths

All filesystem paths for configuration are resolved exclusively by
`internal/dirs`. No other package may call `os.UserConfigDir()` or
construct a config path directly.

| Path | Default | XDG override |
| --- | --- | --- |
| Config root | `~/.config/mdstn/git-msg/` | `$XDG_CONFIG_HOME/mdstn/git-msg/` |
| Config file | `~/.config/mdstn/git-msg/config.toml` | — |
| Prompt templates | `~/.config/mdstn/git-msg/prompts/` | — |

- User templates take precedence over embedded defaults of the same name.
- Default templates MUST be embedded at build time via `//go:embed`.
- `os.UserHomeDir()` is used (not `os.UserConfigDir()`) so the path is
  consistent across macOS and Linux.

### Hook Behaviour

- The installed hook script MUST be the exact constant defined in
  `internal/hook/install.go:hookScript` — a minimal shell wrapper
  delegating entirely to `git-msg generate --hook-mode`.
- `ShouldGenerate` MUST return `true` only for `SourceNormal` and
  `SourceMessage`; all other source types exit 0 silently.
- `hook.FileManager` receives a `GitConfigReader` for resolving
  `core.hooksPath`; it MUST NOT shell out independently.

### Platform (v1)

- Keychain integration targets macOS system keychain via `99designs/keyring`.
- Config paths use the XDG convention and work on both macOS and Linux.
- No Windows support in v1.

---

## Development Workflow & Quality Gates

### Before Starting Any Feature

1. Verify the feature does not conflict with v1 Non-Goals (Principle VI).
2. Confirm the relevant interface exists or create it before implementation.
3. Write unit tests (with fakes) that fail before writing production code.
4. Check Principle VII: identify which layer the new code belongs to and
   confirm its dependency direction is inward only.

### Pull Request Gates

- `go test ./... -race -count=1` MUST pass fully offline.
- `go vet ./...` MUST produce no output.
- No live external calls (API, git, keychain) in the unit test suite.
- No API keys or secrets in committed files.
- `cmd/<command>.go` `Run` functions MUST be testable by direct call with
  fakes — no cobra, no TTY, no subprocess.
- No `*_cobra.go` file may contain logic beyond flag binding, dependency
  construction, and calling `Run`.
- No `internal/` package may import a peer `internal/` package that is at
  the same or outer layer (see Principle VII layer map).

### Commit & Branch Conventions

- Conventional Commits format: `type(scope): subject`
- Branch names: `###-short-description` (numeric prefix optional)
- Squash or rebase before merge; linear history preferred

### Build

- Use `task build` (Taskfile.yml) — `clean` runs as a dependency before
  every build, ensuring no stale artifacts.
- Version is injected at build time:
  `-ldflags "-X main.Version=$(git describe --tags --always)"`

### Dependency Policy

Only the following external dependencies are permitted in v1:

| Package | Role |
| --- | --- |
| `github.com/spf13/cobra` | CLI framework (`*_cobra.go` files only) |
| `github.com/charmbracelet/huh` | TUI forms — review, wizard, inline edit |
| `github.com/99designs/keyring` | System keychain |
| `github.com/madstone-tech/ason/pkg` | Jinja2-style template rendering |
| `github.com/pelletier/go-toml/v2` | TOML read/write (`internal/` only) |
| stdlib `net/http` | All LLM API calls |
| stdlib `os/exec` | Git subprocess calls |

New dependencies MUST be approved via a constitution amendment.

---

## Governance

This constitution supersedes all other project conventions. Any practice
not explicitly permitted here is deferred pending an amendment.

**Amendment procedure**:
1. Propose the change in a PR that updates this file.
2. State the version bump type (MAJOR/MINOR/PATCH) with reasoning.
3. Update all affected template files in `.specify/templates/` in the same PR.
4. Obtain at least one reviewer approval before merging.

**Versioning policy** (semantic):
- MAJOR: Backward-incompatible governance change, principle removal, or
  redefinition that invalidates prior work.
- MINOR: New principle or section added, or materially expanded guidance.
- PATCH: Clarifications, wording fixes, non-semantic refinements.

**Compliance review**: Every plan's "Constitution Check" gate MUST be
completed before Phase 0 research begins and re-checked after Phase 1
design. See `.specify/templates/plan-template.md`.

**Version**: 1.1.0 | **Ratified**: 2026-03-18 | **Last Amended**: 2026-03-18
