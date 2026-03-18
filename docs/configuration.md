# Configuration

`git-msg` stores all configuration at:

```
~/.config/mdstn/git-msg/config.toml
```

If `XDG_CONFIG_HOME` is set, that path is used instead of `~/.config`.

The file is created by the [setup wizard](getting-started.md#first-run) on
first use. API keys are **never** written here — they live in the system
keychain.

---

## Config file reference

```toml
[provider]
name  = "anthropic"       # openai | anthropic | gemini | ollama
model = "claude-haiku-4-5"

[ollama]
host = "http://localhost:11434"

[prompt]
default = "conventional"

[hook]
global = false
```

### All keys

| Key | Type | Default | Description |
| --- | --- | --- | --- |
| `provider.name` | string | `anthropic` | LLM provider to use |
| `provider.model` | string | `claude-haiku-4-5` | Model name for the chosen provider |
| `ollama.host` | string | `http://localhost:11434` | Base URL of your Ollama instance |
| `prompt.default` | string | `conventional` | Template name used when `--template` is not given |
| `hook.global` | bool | `false` | Whether `hook install` targets the global hooks path |

---

## Reading and writing config

### Show the full resolved config

```sh
git-msg config show
```

### Get a single value

```sh
git-msg config get provider.name
git-msg config get provider.model
```

### Set a value

```sh
git-msg config set provider.name openai
git-msg config set provider.model gpt-4o
git-msg config set ollama.host http://192.168.1.10:11434
git-msg config set prompt.default conventional
git-msg config set hook.global true
```

Changes take effect immediately on the next `git-msg` invocation.

---

## API key storage

API keys are stored in the **system keychain** using service name `git-msg`
and the provider name as the account key (e.g. `anthropic`, `openai`).

### How lookup works

1. System keychain (`git-msg` / `<provider>`)
2. Environment variable `GIT_MSG_<PROVIDER>_API_KEY`
3. Error with remediation instructions

### Setting a key via environment variable

```sh
export GIT_MSG_ANTHROPIC_API_KEY=sk-ant-...
export GIT_MSG_OPENAI_API_KEY=sk-...
export GIT_MSG_GEMINI_API_KEY=AIza...
```

This is useful in CI, Docker, or any environment without an interactive
keychain.

### Updating a stored key

Re-run the setup wizard by removing your config and running any command:

```sh
rm ~/.config/mdstn/git-msg/config.toml
git-msg generate --dry-run
```

The wizard will prompt for a new key and overwrite the keychain entry.

---

## Switching providers

```sh
# Switch permanently
git-msg config set provider.name ollama
git-msg config set provider.model fast-coder:latest

# Switch for a single run only
git-msg generate --provider openai --provider-model gpt-4o
```

---

## Resetting to defaults

Delete the config file. The next command will trigger the setup wizard:

```sh
rm ~/.config/mdstn/git-msg/config.toml
```

To also remove all user prompt template overrides:

```sh
rm -rf ~/.config/mdstn/git-msg/
```

---

## Multiple environments

Use `XDG_CONFIG_HOME` to isolate configurations:

```sh
XDG_CONFIG_HOME=/path/to/work-config git-msg generate
XDG_CONFIG_HOME=/path/to/personal-config git-msg generate
```

Each will read from `$XDG_CONFIG_HOME/mdstn/git-msg/config.toml`.
