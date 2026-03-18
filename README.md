# git-msg

> AI-assisted git commit message generator with multi-provider LLM support.

[![CI](https://github.com/madstone-tech/git-msg/actions/workflows/ci.yml/badge.svg)](https://github.com/madstone-tech/git-msg/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/madstone-tech/git-msg)](https://goreportcard.com/report/github.com/madstone-tech/git-msg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`git-msg` watches your staged diff, calls an LLM, and drops a commit message
into an interactive TUI where you can confirm, edit, or abort — before anything
is committed.

```
git add .
git-msg generate
```

```
  Generated commit message:
  > Use as-is
    Edit inline
    Open $EDITOR
    Abort

feat(config): add XDG-compliant path resolution via internal/dirs
```

---

## Features

- **Multi-provider** — OpenAI, Anthropic (Claude), Google Gemini, Ollama (local)
- **Interactive TUI review** — confirm, edit inline, open `$EDITOR`, or abort
- **Git hook integration** — install as `prepare-commit-msg` for automatic generation on every `git commit`
- **Customisable prompts** — Jinja2-style templates with `{{ diff }}`, `{{ branch }}`, `{{ log }}` variables
- **Secure credentials** — API keys stored in the system keychain, never in config files
- **XDG config** — all config lives under `~/.config/mdstn/git-msg/`

---

## Installation

### Homebrew (macOS / Linux)

```sh
brew tap madstone-tech/tap
brew install git-msg
```

### Go install

```sh
go install github.com/madstone-tech/git-msg@latest
```

### Binary releases

Download the latest binary for your platform from the
[Releases](https://github.com/madstone-tech/git-msg/releases) page,
extract, and move to a directory on your `PATH`:

```sh
tar -xzf git-msg_Darwin_arm64.tar.gz
sudo mv git-msg /usr/local/bin/
```

---

## Quick start

**1. Run any `git-msg` command** — the first-run wizard launches automatically
when no config file exists:

```
  git-msg — first-run setup
  ─────────────────────────

  Select LLM provider
  > Anthropic (Claude)
    OpenAI (GPT)
    Google Gemini
    Ollama (local)
```

For Ollama, `git-msg` queries `ollama list` and presents a model picker.
For cloud providers, you enter a model name and API key (stored in the keychain).

**2. Stage some changes and generate:**

```sh
git add src/
git-msg generate
```

**3. (Optional) Install the git hook** so generation happens automatically on
every `git commit`:

```sh
git-msg hook install
```

See [Getting Started](docs/getting-started.md) for a full walkthrough.

---

## Supported providers

| Provider    | Default model      | Needs API key |
| ----------- | ------------------ | ------------- |
| `anthropic` | `claude-haiku-4-5` | Yes           |
| `openai`    | `gpt-4o-mini`      | Yes           |
| `gemini`    | `gemini-1.5-flash` | Yes           |
| `ollama`    | *(from ollama list)* | No          |

Override the provider for a single run:

```sh
git-msg generate --provider ollama
```

See [Providers](docs/providers.md) for setup instructions for each provider.

---

## Configuration

Config file: `~/.config/mdstn/git-msg/config.toml`

```toml
[provider]
name  = "anthropic"
model = "claude-haiku-4-5"

[ollama]
host = "http://localhost:11434"

[prompt]
default = "conventional"

[hook]
global = false
```

```sh
git-msg config set provider.name ollama
git-msg config set provider.model llama3
git-msg config show
```

See [Configuration](docs/configuration.md) for all keys and credential management.

---

## Prompt templates

`git-msg` ships with a built-in `conventional` template that follows the
[Conventional Commits](https://www.conventionalcommits.org) specification.
Create your own templates and edit them in `$EDITOR`:

```sh
git-msg prompt list
git-msg prompt show conventional
git-msg prompt edit conventional
git-msg prompt reset conventional
```

See [Prompt Templates](docs/prompt-templates.md) for the template format and
variable reference.

---

## Documentation

| Document | Description |
| --- | --- |
| [Getting Started](docs/getting-started.md) | Installation, first run, and daily workflow |
| [Configuration](docs/configuration.md) | Config file, all keys, credential storage |
| [Providers](docs/providers.md) | Setup guide for each LLM provider |
| [Prompt Templates](docs/prompt-templates.md) | Template format, variables, and customisation |
| [Git Hook](docs/git-hook.md) | Automatic generation via `prepare-commit-msg` |
| [CLI Reference](docs/cli-reference.md) | Complete command and flag reference |

---

## Development

```sh
git clone https://github.com/madstone-tech/git-msg
cd git-msg
task build          # compile to bin/git-msg
task test           # run test suite
task ci             # fmt + vet + lint + constitution + test + build
```

**Requirements**: Go 1.26+, [Task](https://taskfile.dev)

---

## License

[MIT](LICENSE)
