# PRD: `git-msg`

> Go CLI for AI-assisted git commit message generation.

## Status: Draft

## Background

Shell-script approaches (`git-gemini-commit.sh`, Harper Reed's `prepare-commit-msg` hook) solve the problem of bad commit messages but are fragile:
- Single LLM provider, tightly coupled
- No configuration management — prompts hardcoded or require manual file setup
- No git hook lifecycle management
- No context injection (branch, recent log, ticket refs)
- Poor error handling, no retry
- Not composable or testable

**Reference implementations:**
- https://harper.blog/2024/03/11/use-an-llm-to-automagically-generate-meaningful-git-commit-messages/
- `git-gemini-commit.sh` in this repo

---

## Goals

1. Multi-provider LLM support (OpenAI, Anthropic, Gemini, Ollama) via pluggable backends
2. Customizable Jinja2-style prompt templates with variable interpolation
3. Interactive TUI review with inline edit and `$EDITOR` fallback
4. Git hook lifecycle management (install/uninstall, global or per-repo)
5. Secure API key storage via system keychain
6. Clean, idiomatic Go — testable business logic fully decoupled from CLI wiring

---

## Non-Goals (v1)

- No diff chunking for very large diffs
- No `git add -p` staging integration
- No team/shared prompt server
- No Windows support (keychain assumption)
- No i18n

---

## CLI Surface

```
git-msg
├── generate              # Generate commit message from staged diff (default command)
│   ├── --dry-run         # Print only, don't commit
│   ├── --provider NAME   # Override configured provider for this run
│   └── --template NAME   # Override prompt template for this run
│
├── config
│   ├── set KEY VALUE     # Set a config value
│   ├── get KEY           # Read a config value
│   └── show              # Print full resolved config
│
├── prompt
│   ├── list              # List available templates
│   ├── show NAME         # Print a template
│   ├── edit NAME         # Open template in $EDITOR
│   └── reset NAME        # Restore template to built-in default
│
└── hook
    ├── install           # Install as prepare-commit-msg hook
    │   └── --global      # Install globally via git core.hooksPath
    └── uninstall         # Remove the hook
        └── --global
```

---

## Architecture

### Cobra Double-File Pattern

Every command is split into two files:

- `cmd/generate.go` — pure business logic function `Run(ctx, opts) error`. No Cobra dependency. Directly testable.
- `cmd/generate_cobra.go` — `NewGenerateCmd() *cobra.Command`. Binds flags, constructs options struct, calls `Run`.

This keeps `cmd/` as thin glue. All real logic lives in `internal/` packages and the non-cobra `cmd/*.go` files.

### Package Structure

```
git-msg/
├── main.go
│
├── cmd/
│   ├── root.go              # Version, persistent flags, context setup
│   ├── root_cobra.go        # Execute() entry point
│   ├── generate.go          # Run(ctx, GenerateOptions) error
│   ├── generate_cobra.go    # NewGenerateCmd() *cobra.Command
│   ├── config.go            # SetConfig, GetConfig, ShowConfig
│   ├── config_cobra.go      # NewConfigCmd() *cobra.Command
│   ├── prompt.go            # ListPrompts, ShowPrompt, EditPrompt, ResetPrompt
│   ├── prompt_cobra.go      # NewPromptCmd() *cobra.Command
│   ├── hook.go              # InstallHook, UninstallHook
│   └── hook_cobra.go        # NewHookCmd() *cobra.Command
│
└── internal/
    ├── config/
    │   ├── config.go        # Config struct + defaults
    │   ├── store.go         # Store interface: Load(), Save()
    │   └── file.go          # TOML file impl (~/.config/git-msg/config.toml)
    │
    ├── secret/
    │   ├── store.go         # SecretStore interface: Get, Set, Delete
    │   └── keychain.go      # 99designs/keyring impl + env var fallback
    │
    ├── prompt/
    │   ├── template.go      # Template struct {Name, SystemText, UserText}
    │   ├── store.go         # TemplateStore interface: Get, List, Save, Delete
    │   ├── file.go          # File impl (~/.config/git-msg/prompts/)
    │   ├── defaults.go      # Embedded defaults via //go:embed
    │   └── renderer.go      # Wraps ason Engine.Render(); injects TemplateVars
    │
    ├── llm/
    │   ├── provider.go      # Provider interface: Generate(ctx, system, user) (string, error)
    │   ├── openai.go        # OpenAI REST impl
    │   ├── anthropic.go     # Anthropic REST impl
    │   ├── gemini.go        # Gemini REST impl
    │   ├── ollama.go        # Ollama REST impl
    │   └── factory.go       # NewProvider(cfg, secrets) (Provider, error)
    │
    ├── git/
    │   ├── client.go        # Client struct wrapping os/exec
    │   ├── diff.go          # StagedDiff() (string, error)
    │   ├── branch.go        # CurrentBranch() (string, error)
    │   ├── log.go           # RecentLog(n int) (string, error)
    │   └── commit.go        # Commit(msg string) error
    │
    ├── ui/
    │   ├── review.go        # huh form: confirm | edit inline | open $EDITOR
    │   ├── spinner.go       # Spinner(label) StopFunc
    │   └── setup.go         # First-run wizard (huh): provider, model, api key
    │
    └── hook/
        ├── source.go        # SourceType enum + ShouldGenerate(SourceType) bool
        ├── manager.go       # Manager interface: Install, Uninstall, IsInstalled
        └── install.go       # Hook script generation + git config wiring
```

---

## Data Flow: `git-msg generate`

```
cmd/generate.go: Run(ctx, opts)
  │
  ├─ config.Store.Load()                 → Config
  │   └─ first-run: ui.Setup() wizard if no config exists
  │
  ├─ secret.Store.Get(provider)          → api key
  │
  ├─ git.Client.StagedDiff()             → diff string
  │   └─ error if no staged changes
  │
  ├─ git.Client.CurrentBranch()          → branch string
  ├─ git.Client.RecentLog(5)             → log string
  │
  ├─ prompt.FileStore.Get(name)          → Template{SystemText, UserText}
  ├─ prompt.Renderer.Render(t, vars)     → rendered system + user strings
  │     └─ ason Engine.Render(str, map{
  │              "diff": diff,
  │              "branch": branch,
  │              "log": log })
  │
  ├─ ui.Spinner("Generating...")
  ├─ llm.Provider.Generate(ctx, system, user) → raw message string
  ├─ spinner.Stop()
  │
  ├─ if --dry-run: print and return
  │
  ├─ ui.Review(message) → (finalMessage, error)
  │     ├─ [confirm]  use as-is
  │     ├─ [edit]     huh textarea inline edit
  │     ├─ [editor]   write tmp file → open $EDITOR → read back
  │     └─ [abort]    return error, exit 0
  │
  └─ if --hook-mode:
  │     write finalMessage to opts.HookMsgFile
  └─ else:
        git.Client.Commit(finalMessage)
```

---

## Hook Mode

### Installed Script

The hook installed into `.git/hooks/prepare-commit-msg` (or global hooks dir) is a minimal self-referencing shell script:

```sh
#!/bin/sh
# Managed by git-msg. Do not edit manually.
exec git-msg generate --hook-mode --hook-msg-file "$1" --hook-source "${2:-}"
```

### Source Type Dispatch

`prepare-commit-msg` receives `$2` indicating commit source. v1 implements `normal` and `message`; all others are stubbed to skip.

```go
// internal/hook/source.go

type SourceType string

const (
    SourceNormal   SourceType = ""         // plain git commit
    SourceMessage  SourceType = "message"  // git commit -m "..." given
    SourceTemplate SourceType = "template" // commit.template used
    SourceMerge    SourceType = "merge"    // merge commit
    SourceSquash   SourceType = "squash"   // squash commit
    SourceCommit   SourceType = "commit"   // amend
)

// ShouldGenerate returns true for source types where generation is implemented.
// template, merge, squash, commit are stubbed — behavior not yet defined.
func ShouldGenerate(s SourceType) bool {
    switch s {
    case SourceNormal, SourceMessage:
        return true
    default:
        return false
    }
}
```

---

## Prompt Templates

### Storage

- Default templates embedded at build time via `//go:embed` in `internal/prompt/defaults.go`
- User templates stored at `~/.config/git-msg/prompts/<name>.toml`
- User template takes precedence over embedded default with same name

### Template File Format

```toml
name = "conventional"
description = "Conventional Commits with branch and log context"

system = """
You are a git commit message generator.
Follow the Conventional Commits specification (type(scope): subject).
Output only the commit message. No explanation, no markdown fencing.
"""

user = """
Branch: {{ branch }}

Recent commits:
{{ log }}

Staged diff:
{{ diff }}
"""
```

### Template Variables

| Variable  | Source                         |
| --------- | ------------------------------ |
| `diff`    | `git diff --cached`            |
| `branch`  | `git rev-parse --abbrev-ref HEAD` |
| `log`     | `git log --oneline -5`         |

### Rendering

Uses `github.com/madstone-tech/ason/pkg` Engine directly:

```go
engine := pkg.NewDefaultEngine()
rendered, err := engine.Render(templateStr, map[string]interface{}{
    "diff":   diff,
    "branch": branch,
    "log":    log,
})
```

---

## Configuration

### Schema (`~/.config/git-msg/config.toml`)

```toml
[provider]
name  = "anthropic"         # openai | anthropic | gemini | ollama
model = "claude-haiku-4-5"
# api_key stored in system keychain, not here

[ollama]
host = "http://localhost:11434"

[prompt]
default = "conventional"

[hook]
global = false
```

### API Key Storage

Keys stored in system keychain via `github.com/99designs/keyring`.
Service name: `git-msg`. Account name: provider name (e.g. `anthropic`).

Lookup order:
1. System keychain
2. Environment variable: `GIT_MSG_<PROVIDER>_API_KEY` (e.g. `GIT_MSG_ANTHROPIC_API_KEY`)
3. Error with instructions

---

## First-Run Setup Wizard

Triggered when `~/.config/git-msg/config.toml` does not exist.
Implemented as a `huh` form sequence in `internal/ui/setup.go`:

1. Select provider (`openai` / `anthropic` / `gemini` / `ollama`)
2. Enter model name (prefilled with sensible default per provider)
3. Enter API key → stored via keyring (skipped for ollama)
4. Confirm → write `config.toml`

Default models per provider:

| Provider    | Default model            |
| ----------- | ------------------------ |
| `anthropic` | `claude-haiku-4-5`       |
| `openai`    | `gpt-4o-mini`            |
| `gemini`    | `gemini-1.5-flash`       |
| `ollama`    | `llama3`                 |

---

## LLM Providers

All providers implement the same interface using stdlib `net/http` only — no vendor SDKs.

```go
// internal/llm/provider.go
type Provider interface {
    Generate(ctx context.Context, system, user string) (string, error)
}
```

| Provider    | API base URL                                          |
| ----------- | ----------------------------------------------------- |
| `anthropic` | `https://api.anthropic.com/v1/messages`               |
| `openai`    | `https://api.openai.com/v1/chat/completions`          |
| `gemini`    | `https://generativelanguage.googleapis.com/v1beta/...`|
| `ollama`    | `http://localhost:11434/api/chat`                     |

---

## External Dependencies

| Package                              | Role                                      |
| ------------------------------------ | ----------------------------------------- |
| `github.com/spf13/cobra`             | CLI framework                             |
| `github.com/charmbracelet/huh`       | TUI forms (review, setup wizard, spinner) |
| `github.com/99designs/keyring`       | System keychain (note: package is `keyring`) |
| `github.com/madstone-tech/ason/pkg`  | Jinja2-style template rendering           |
| `github.com/pelletier/go-toml/v2`   | TOML config and template file read/write  |
| stdlib `net/http`                    | All LLM API calls                         |
| stdlib `os/exec`                     | Git subprocess calls                      |

---

## Testing Strategy

- `internal/` packages tested directly with no CLI overhead
- `cmd/*.go` (non-cobra files) tested via `Run(ctx, opts)` with fakes injected
- LLM providers: interface-based fakes, no live API calls in unit tests
- Git client: fake via interface, fixture diffs in `testdata/`
- Hook source dispatch: pure function, table-driven tests

---

## Open Questions

- Should `--dry-run` still invoke `ui.Review` (for copy-paste) or skip it entirely and just print?
- Should `prompt edit` create the file from the embedded default if no user override exists yet?
- Max diff size before truncating to avoid LLM context limits — configurable or hardcoded?
