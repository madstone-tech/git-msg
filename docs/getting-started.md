# Getting Started

This guide walks you from zero to a working `git-msg` setup in under five
minutes.

## Prerequisites

- macOS or Linux
- Git
- One of:
  - An API key for [Anthropic](https://console.anthropic.com),
    [OpenAI](https://platform.openai.com), or
    [Google Gemini](https://aistudio.google.com)
  - [Ollama](https://ollama.ai) running locally with at least one model pulled

---

## Installation

### Homebrew (recommended)

```sh
brew tap madstone-tech/tap
brew install git-msg
```

### Go install

```sh
go install github.com/madstone-tech/git-msg@latest
```

### Binary releases

Download the binary for your platform from
[GitHub Releases](https://github.com/madstone-tech/git-msg/releases),
extract, and move to your `PATH`:

```sh
# macOS (Apple Silicon)
tar -xzf git-msg_Darwin_arm64.tar.gz
sudo mv git-msg /usr/local/bin/

# macOS (Intel)
tar -xzf git-msg_Darwin_x86_64.tar.gz
sudo mv git-msg /usr/local/bin/

# Linux (x86_64)
tar -xzf git-msg_Linux_x86_64.tar.gz
sudo mv git-msg /usr/local/bin/
```

Verify:

```sh
git-msg --version
```

---

## First run

The first time you run any `git-msg` command, the setup wizard launches
automatically if no config file exists.

```sh
git-msg generate --dry-run
```

### Step 1 — Select a provider

```
  git-msg — first-run setup
  ─────────────────────────

  Select LLM provider
  > Anthropic (Claude)
    OpenAI (GPT)
    Google Gemini
    Ollama (local)
```

Use arrow keys to move, Enter to select.

### Step 2 — Choose a model

**For Ollama**: `git-msg` runs `ollama list` and shows all locally available
models. Select the one you want.

**For cloud providers**: a text input appears pre-filled with a sensible
default. Press Enter to accept or type a different model name.

| Provider    | Default         |
| ----------- | --------------- |
| `anthropic` | `claude-haiku-4-5` |
| `openai`    | `gpt-4o-mini`   |
| `gemini`    | `gemini-1.5-flash` |

### Step 3 — API key (cloud providers only)

Enter your API key. It is stored in the system keychain — never written to
disk. This step is skipped entirely for Ollama.

### Completion

```
  Config saved. Provider: anthropic  Model: claude-haiku-4-5
```

Your config file is now at `~/.config/mdstn/git-msg/config.toml`.

---

## Daily workflow

### Generate a commit message

```sh
# Stage your changes as normal
git add src/auth.go

# Generate a message
git-msg generate
```

The tool collects your staged diff, current branch name, and five most recent
commits, then sends them to the LLM. A spinner shows progress.

When the message is ready, the review TUI appears:

```
  Generated commit message:
  > Use as-is
    Edit inline
    Open $EDITOR
    Abort

feat(auth): add JWT refresh token rotation
```

| Option | What happens |
| --- | --- |
| **Use as-is** | Commits with the generated message |
| **Edit inline** | Opens a textarea in the terminal to edit the message |
| **Open $EDITOR** | Opens `$EDITOR` (e.g. `vim`, `nvim`) with the message in a temp file |
| **Abort** | Exits with code 0; your staged changes are untouched |

### Dry run (no commit)

Print the generated message without committing or showing the review TUI:

```sh
git-msg generate --dry-run
```

Useful for previewing output or piping into other tools.

### Override provider for one run

```sh
git-msg generate --provider ollama
git-msg generate --provider openai
```

The override applies only to the current invocation; your saved config is
unchanged.

### Override template for one run

```sh
git-msg generate --template conventional
```

---

## Automatic generation via git hook

Install `git-msg` as a `prepare-commit-msg` hook so it fires automatically
on every `git commit`:

```sh
git-msg hook install
```

From that point on:

```sh
git add .
git commit          # review TUI appears automatically
```

To remove the hook:

```sh
git-msg hook uninstall
```

See [Git Hook](git-hook.md) for global hook installation and hook behaviour
details.

---

## Where config lives

| Path | Purpose |
| --- | --- |
| `~/.config/mdstn/git-msg/config.toml` | Provider, model, defaults |
| `~/.config/mdstn/git-msg/prompts/` | User prompt template overrides |
| System keychain (service `git-msg`) | API keys |

The `XDG_CONFIG_HOME` environment variable is respected if set.

---

## Next steps

- [Configuration](configuration.md) — all config keys, credential management
- [Providers](providers.md) — detailed setup for each LLM provider
- [Prompt Templates](prompt-templates.md) — customise what gets sent to the LLM
- [Git Hook](git-hook.md) — automate generation on every commit
- [CLI Reference](cli-reference.md) — full command and flag reference
