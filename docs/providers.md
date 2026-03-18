# Providers

`git-msg` supports four LLM providers. All use stdlib `net/http` — no vendor
SDKs are bundled.

---

## Anthropic (Claude)

**API reference**: [docs.anthropic.com](https://docs.anthropic.com)

### Setup

1. Create an account at [console.anthropic.com](https://console.anthropic.com)
2. Generate an API key under **API Keys**
3. Run the setup wizard:

```sh
rm ~/.config/mdstn/git-msg/config.toml
git-msg generate --dry-run
# Select: Anthropic (Claude)
# Enter your key when prompted
```

Or configure manually:

```sh
git-msg config set provider.name anthropic
git-msg config set provider.model claude-haiku-4-5
# Key is stored interactively via wizard, or:
export GIT_MSG_ANTHROPIC_API_KEY=sk-ant-...
```

### Recommended models

| Model | Speed | Quality | Notes |
| --- | --- | --- | --- |
| `claude-haiku-4-5` | Fast | Good | Default — best cost/speed ratio |
| `claude-sonnet-4-5` | Medium | Great | Better for complex diffs |
| `claude-opus-4-5` | Slow | Best | Large context, detailed commits |

### API endpoint

```
https://api.anthropic.com/v1/messages
```

---

## OpenAI (GPT)

**API reference**: [platform.openai.com/docs](https://platform.openai.com/docs)

### Setup

1. Create an account at [platform.openai.com](https://platform.openai.com)
2. Generate an API key under **API Keys**
3. Configure:

```sh
git-msg config set provider.name openai
git-msg config set provider.model gpt-4o-mini
export GIT_MSG_OPENAI_API_KEY=sk-...
```

### Recommended models

| Model | Speed | Quality | Notes |
| --- | --- | --- | --- |
| `gpt-4o-mini` | Fast | Good | Default — cheap and capable |
| `gpt-4o` | Medium | Great | Better reasoning |
| `o1-mini` | Slow | Best | Highest quality, higher cost |

### API endpoint

```
https://api.openai.com/v1/chat/completions
```

---

## Google Gemini

**API reference**: [ai.google.dev](https://ai.google.dev)

### Setup

1. Go to [aistudio.google.com](https://aistudio.google.com) and create an API key
2. Configure:

```sh
git-msg config set provider.name gemini
git-msg config set provider.model gemini-1.5-flash
export GIT_MSG_GEMINI_API_KEY=AIza...
```

### Recommended models

| Model | Speed | Quality | Notes |
| --- | --- | --- | --- |
| `gemini-1.5-flash` | Fast | Good | Default — generous free tier |
| `gemini-1.5-pro` | Medium | Great | Larger context window |
| `gemini-2.0-flash` | Fast | Great | Latest Flash generation |

### API endpoint

```
https://generativelanguage.googleapis.com/v1beta/models/<model>:generateContent
```

The API key is passed as a query parameter (`?key=...`), not a header.

---

## Ollama (local)

**Documentation**: [ollama.ai](https://ollama.ai)

Ollama runs models locally on your machine — no API key required, no data
leaves your network.

### Setup

1. Install Ollama from [ollama.ai/download](https://ollama.ai/download)
2. Pull a model:

```sh
ollama pull llama3
ollama pull mistral
ollama pull codellama
```

3. Configure `git-msg`:

```sh
git-msg config set provider.name ollama
git-msg config set provider.model llama3
```

Or run the setup wizard — it will query `ollama list` and show a picker of
your installed models.

### Listing available models

```sh
ollama list
```

`git-msg` queries this automatically during the setup wizard.

### Custom host

If Ollama is running on a different machine or port:

```sh
git-msg config set ollama.host http://192.168.1.50:11434
```

### Recommended models for commit messages

| Model | Size | Notes |
| --- | --- | --- |
| `llama3` | 4.7GB | Good all-rounder |
| `mistral` | 4.1GB | Fast, concise output |
| `codellama` | 3.8GB | Code-focused |
| `phi4` | 9.1GB | High quality, larger |

### API endpoint

```
http://localhost:11434/api/chat
```

No `Authorization` header is sent.

---

## Switching providers

For a single run, use `--provider`:

```sh
git-msg generate --provider ollama
git-msg generate --provider anthropic
```

To change the default permanently:

```sh
git-msg config set provider.name gemini
git-msg config set provider.model gemini-1.5-flash
```

---

## Timeouts

All providers use a 30-second HTTP timeout. If the LLM is slow to respond
(large diff, slow network), the request will fail with a timeout error.
Reduce diff size by staging fewer files, or use a faster model.
