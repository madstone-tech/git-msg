# CLI Reference

Complete reference for all `git-msg` commands and flags.

---

## Global flags

These flags are available on every command:

| Flag | Description |
| --- | --- |
| `-h`, `--help` | Show help for the command |
| `-v`, `--version` | Print the version string and exit |

---

## `git-msg generate`

Generate a commit message from the current staged diff.

```sh
git-msg generate [flags]
```

### Flags

| Flag | Type | Description |
| --- | --- | --- |
| `--dry-run` | bool | Print the generated message to stdout without committing or showing the review TUI |
| `--provider NAME` | string | Override the configured provider for this run only (`anthropic`, `openai`, `gemini`, `ollama`) |
| `--template NAME` | string | Override the configured default template for this run only |

### Behaviour

1. Loads config from `~/.config/mdstn/git-msg/config.toml`
2. Collects `git diff --cached`, current branch, and last 5 commits
3. Renders the prompt template with the collected context
4. Calls the configured LLM provider
5. Shows the review TUI (unless `--dry-run`)
6. On confirm: runs `git commit -m <message>`

### Examples

```sh
# Standard generation with review
git-msg generate

# Preview without committing
git-msg generate --dry-run

# Use Ollama for this commit only
git-msg generate --provider ollama

# Use a custom template for this commit only
git-msg generate --template jira
```

### Exit codes

| Code | Meaning |
| --- | --- |
| `0` | Success, or user aborted (no commit was made) |
| `1` | Error (no staged changes, provider failure, config missing) |

---

## `git-msg config`

Manage `git-msg` configuration.

### `git-msg config set KEY VALUE`

Set a configuration value.

```sh
git-msg config set KEY VALUE
```

**Supported keys:**

| Key | Example value |
| --- | --- |
| `provider.name` | `ollama` |
| `provider.model` | `fast-coder:latest` |
| `ollama.host` | `http://localhost:11434` |
| `prompt.default` | `conventional` |
| `hook.global` | `true` |

```sh
git-msg config set provider.name anthropic
git-msg config set provider.model claude-sonnet-4-5
```

### `git-msg config get KEY`

Print the current value of a config key.

```sh
git-msg config get provider.name
# anthropic
```

### `git-msg config show`

Print the full resolved configuration as TOML.

```sh
git-msg config show
```

```toml
[provider]
  name = 'anthropic'
  model = 'claude-haiku-4-5'

[ollama]
  host = 'http://localhost:11434'

[prompt]
  default = 'conventional'

[hook]
  global = false
```

API keys are **never** included in this output.

---

## `git-msg prompt`

Manage prompt templates.

### `git-msg prompt list`

List all available templates and their source.

```sh
git-msg prompt list
```

```
  conventional         [embedded]
  jira                 [user]
```

- `[embedded]` — built-in template shipped with `git-msg`
- `[user]` — user-defined template in `~/.config/mdstn/git-msg/prompts/`

### `git-msg prompt show NAME`

Print the full content of a template.

```sh
git-msg prompt show conventional
```

### `git-msg prompt edit NAME`

Open a template in `$EDITOR` for editing. If no user override exists, the
embedded default is pre-loaded.

```sh
git-msg prompt edit conventional
```

Requires `$EDITOR` to be set. On save, the result is written to
`~/.config/mdstn/git-msg/prompts/<name>.toml`.

### `git-msg prompt reset NAME`

Delete the user override for a template, restoring the embedded default.
Exits 0 if no override exists.

```sh
git-msg prompt reset conventional
```

---

## `git-msg hook`

Manage the `prepare-commit-msg` git hook lifecycle.

### `git-msg hook install`

Install the hook into the current repository's `.git/hooks/` directory.

```sh
git-msg hook install [--global]
```

| Flag | Description |
| --- | --- |
| `--global` | Install into the directory specified by `git config core.hooksPath` instead |

The installed script exits silently when no TTY is available, so it does not
block non-interactive flows.

### `git-msg hook uninstall`

Remove the `prepare-commit-msg` hook.

```sh
git-msg hook uninstall [--global]
```

| Flag | Description |
| --- | --- |
| `--global` | Remove from the global hooks path |

Exits 0 if the hook file does not exist.

---

## `git-msg completion`

Generate shell completion scripts.

```sh
# bash
git-msg completion bash > /etc/bash_completion.d/git-msg

# zsh
git-msg completion zsh > "${fpath[1]}/_git-msg"

# fish
git-msg completion fish > ~/.config/fish/completions/git-msg.fish
```

---

## Environment variables

| Variable | Description |
| --- | --- |
| `XDG_CONFIG_HOME` | Override the base config directory (default: `~/.config`) |
| `EDITOR` | Editor used by `prompt edit` and the "Open $EDITOR" review option |
| `GIT_MSG_ANTHROPIC_API_KEY` | Anthropic API key (fallback if not in keychain) |
| `GIT_MSG_OPENAI_API_KEY` | OpenAI API key (fallback if not in keychain) |
| `GIT_MSG_GEMINI_API_KEY` | Google Gemini API key (fallback if not in keychain) |
