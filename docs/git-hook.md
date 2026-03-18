# Git Hook

`git-msg` can be installed as a `prepare-commit-msg` hook so that commit
message generation happens automatically every time you run `git commit`.

---

## Installing the hook

### Per-repository (recommended)

Run inside a git repository:

```sh
git-msg hook install
```

This writes an executable script to `.git/hooks/prepare-commit-msg`.

### Globally (all repositories)

Install once and have it apply to every git repo on your machine:

```sh
# First, configure a global hooks directory
git config --global core.hooksPath ~/.config/git/hooks
mkdir -p ~/.config/git/hooks

# Then install the hook globally
git-msg hook install --global
```

`git-msg` will read `core.hooksPath` from your git config and install there.

---

## What the hook script looks like

```sh
#!/bin/sh
# Managed by git-msg. Do not edit manually.
# Skip generation when stdin is not an interactive terminal (e.g. CI,
# editor integrations, or any non-interactive shell environment).
if ! [ -t 0 ]; then
  exit 0
fi
exec git-msg generate --hook-mode --hook-msg-file "$1" --hook-source "${2:-}"
```

Key behaviours:

- **No TTY, no generation** — the script checks whether `stdin` is an
  interactive terminal (`-t 0`). The review TUI requires an interactive stdin
  to read user input, so the hook exits silently when one is not available
  (CI pipelines, editor integrations, non-interactive shells).
- **Hook mode** — `git-msg` writes the generated message to the commit message
  file that git passes as `$1`, instead of calling `git commit -m` itself.
- **Source-aware** — the `$2` argument from git tells `git-msg` what triggered
  the commit. Only `normal` and `message` sources trigger generation; merges,
  squashes, and amends are passed through unchanged.

---

## How commit source dispatch works

`prepare-commit-msg` receives a second argument indicating why the commit is
happening:

| Source | Example | `git-msg` behaviour |
| --- | --- | --- |
| *(empty)* | `git commit` | Generates message |
| `message` | `git commit -m "..."` | Generates message (overwrites) |
| `merge` | Merge commit | Passes through — no generation |
| `squash` | Squash commit | Passes through — no generation |
| `template` | `commit.template` set | Passes through — no generation |
| `commit` | `git commit --amend` | Passes through — no generation |

---

## Using the hook

Once installed, just run `git commit` as normal:

```sh
git add src/
git commit
```

The review TUI appears:

```
  Generated commit message:
  > Use as-is
    Edit inline
    Open $EDITOR
    Abort

feat(auth): add JWT refresh token rotation
```

Select an option. The chosen message is written to the commit and git
proceeds normally.

---

## Checking hook status

```sh
# Check if the local hook is installed
ls -la .git/hooks/prepare-commit-msg

# Check if a global hook is installed
git config --get core.hooksPath
ls -la "$(git config --get core.hooksPath)/prepare-commit-msg"
```

---

## Reinstalling after updates

If you update `git-msg` and want to ensure the latest hook script is in
place (e.g. after a bug fix to the script), reinstall:

```sh
git-msg hook install       # local
git-msg hook install --global  # global
```

Install is idempotent — running it again replaces the script silently.

---

## Removing the hook

```sh
git-msg hook uninstall              # local
git-msg hook uninstall --global     # global
```

Uninstall is also idempotent — exits 0 if no hook is present.
