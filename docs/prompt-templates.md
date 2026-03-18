# Prompt Templates

`git-msg` uses Jinja2-style prompt templates to control what gets sent to
the LLM. Templates have a `system` section (instructions) and a `user`
section (the actual content with diff, branch, and log injected).

---

## Built-in templates

`git-msg` ships with one embedded template: `conventional`.

```sh
git-msg prompt list
#   conventional         [embedded]
```

The `conventional` template instructs the LLM to follow the
[Conventional Commits](https://www.conventionalcommits.org) specification
(`type(scope): subject`).

```sh
git-msg prompt show conventional
```

```
## System
You are a git commit message generator.
Follow the Conventional Commits specification (type(scope): subject).
Keep the subject line under 72 characters.
Output only the commit message. No explanation, no markdown fencing.

## User
Branch: {{ branch }}

Recent commits:
{{ log }}

Staged diff:
{{ diff }}
```

---

## Template variables

| Variable   | Source                             | Example                          |
| ---------- | ---------------------------------- | -------------------------------- |
| `{{ diff }}` | `git diff --cached`                | Full staged diff as a string     |
| `{{ branch }}` | `git rev-parse --abbrev-ref HEAD` | `feature/auth`                   |
| `{{ log }}`  | `git log --oneline -5`             | Last 5 commits, one per line     |

All three variables are always available. A variable that evaluates to an
empty string (e.g. no recent commits in a fresh repo) is injected as an
empty string.

---

## Template file format

Templates are TOML files:

```toml
name = "my-template"
description = "Short description shown in prompt list"

system = """
You are a git commit message generator.
Write a single-line commit message in imperative mood.
Output only the commit message.
"""

user = """
Branch: {{ branch }}

Recent commits:
{{ log }}

Staged diff:
{{ diff }}
"""
```

User templates are stored at:

```
~/.config/mdstn/git-msg/prompts/<name>.toml
```

A user template with the same name as an embedded template takes precedence.

---

## Managing templates

### List all templates

```sh
git-msg prompt list
```

```
  conventional         [embedded]
  my-template          [user]
```

### Show a template

```sh
git-msg prompt show conventional
git-msg prompt show my-template
```

### Edit a template in `$EDITOR`

```sh
git-msg prompt edit conventional
```

If the template has no user override yet, the embedded default is pre-loaded
into the editor. On save, the result is written to
`~/.config/mdstn/git-msg/prompts/conventional.toml` as a user override.

Make sure `$EDITOR` is set:

```sh
export EDITOR=nvim      # or vim, nano, code, etc.
```

### Reset a template to its embedded default

```sh
git-msg prompt reset conventional
```

Deletes the user override file. The embedded default is used again.
If no user override exists, this is a no-op (exits 0).

---

## Creating a custom template

Create a TOML file directly:

```sh
mkdir -p ~/.config/mdstn/git-msg/prompts/
cat > ~/.config/mdstn/git-msg/prompts/jira.toml << 'EOF'
name = "jira"
description = "Jira ticket prefix from branch name"

system = """
You are a git commit message generator.
Extract the Jira ticket number from the branch name (format: ABC-123).
Write a commit message in the format: ABC-123: imperative description.
Output only the commit message.
"""

user = """
Branch: {{ branch }}

Staged diff:
{{ diff }}
"""
EOF
```

Then set it as default:

```sh
git-msg config set prompt.default jira
```

Or use it for a single run:

```sh
git-msg generate --template jira
```

---

## Template design tips

**Be explicit about output format.** The LLM will follow your instructions
precisely. Include "Output only the commit message" to avoid explanatory text
or markdown fencing appearing in the commit.

**Use branch context for ticket references.** Branch names like
`ABC-123-add-login` can be parsed by the LLM to prefix commit messages with
a ticket number.

**Limit scope with the system prompt.** If you only want the subject line
(no body), say so explicitly: "Write a single line only."

**Include log for consistency.** The `{{ log }}` variable gives the LLM
context about your team's commit style from recent history.
