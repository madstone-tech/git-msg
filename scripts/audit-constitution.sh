#!/usr/bin/env bash
# audit-constitution.sh
# Enforces git-msg Constitution v1.1.0 Principle II and Principle VII.
# Exit 0 = all checks pass. Exit 1 = violations found.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VIOLATIONS=0

red() { printf '\033[0;31m%s\033[0m\n' "$*"; }
green() { printf '\033[0;32m%s\033[0m\n' "$*"; }
info() { printf '  %s\n' "$*"; }

fail() {
	red "❌ $1"
	shift
	while [[ $# -gt 0 ]]; do
		info "$1"
		shift
	done
	VIOLATIONS=$((VIOLATIONS + 1))
}

# pass() prints a success line only when no violations were added in the
# preceding block. Callers pass the violation count before the block as $2.
check_pass() {
	local label="$1"
	local before="$2"
	if [[ "$VIOLATIONS" -eq "$before" ]]; then
		green "✅ $label"
	fi
}

echo ""
echo "git-msg Constitution Audit (v1.1.0)"
echo "======================================"
echo ""

# ---------------------------------------------------------------------------
# Principle II — Cobra Double-File Separation
# cmd/*.go (non-cobra) must not import github.com/spf13/cobra
# ---------------------------------------------------------------------------
echo "Principle II: Cobra double-file separation"
before=$VIOLATIONS

while IFS= read -r -d '' file; do
	if [[ "$file" == *"_cobra.go" ]]; then continue; fi
	if grep -q '"github.com/spf13/cobra"' "$file"; then
		fail "Cobra import in logic file: $file" \
			"Non-_cobra.go files must not import cobra. Move wiring to $(basename "${file%.go}")_cobra.go"
	fi
done < <(find "$ROOT/cmd" -name "*.go" -not -name "*_test.go" -print0)

check_pass "cmd/*.go (non-cobra) files have no cobra imports" "$before"

# ---------------------------------------------------------------------------
# Principle VII — Clean Architecture: internal/ui must not import internal/*
# Uses a broad pattern — any internal/ import is a violation, regardless of
# which package, so future packages are caught automatically.
# ---------------------------------------------------------------------------
echo ""
echo "Principle VII: internal/ui has no internal/* imports"
before=$VIOLATIONS

while IFS= read -r -d '' file; do
	if grep -qE '"github.com/madstone-tech/git-msg/internal/' "$file"; then
		IMPORT=$(grep -E '"github.com/madstone-tech/git-msg/internal/' "$file" | head -1 | xargs)
		fail "internal/ui imports another internal package: $file" \
			"Found: $IMPORT" \
			"ui must return plain values. Persistence belongs in cmd/root.go:EnsureConfig."
	fi
done < <(find "$ROOT/internal/ui" -name "*.go" -not -name "*_test.go" -print0)

check_pass "internal/ui has no internal/* imports" "$before"

# ---------------------------------------------------------------------------
# Principle VII — internal/llm must not import internal/config or internal/secret
# ---------------------------------------------------------------------------
echo ""
echo "Principle VII: internal/llm has no config/secret imports"
before=$VIOLATIONS

while IFS= read -r -d '' file; do
	if grep -qE '"github.com/madstone-tech/git-msg/internal/(config|secret)"' "$file"; then
		IMPORT=$(grep -E '"github.com/madstone-tech/git-msg/internal/(config|secret)"' "$file" | head -1 | xargs)
		fail "internal/llm imports config or secret: $file" \
			"Found: $IMPORT" \
			"Provider construction belongs in cmd/llm.go:NewLLMProvider, not in the llm package."
	fi
done < <(find "$ROOT/internal/llm" -name "*.go" -not -name "*_test.go" -print0)

check_pass "internal/llm has no config/secret imports" "$before"

# ---------------------------------------------------------------------------
# Principle VII — internal/hook must not call exec.Command directly
# ---------------------------------------------------------------------------
echo ""
echo "Principle VII: internal/hook has no raw exec.Command calls"
before=$VIOLATIONS

while IFS= read -r -d '' file; do
	if grep -q 'exec\.Command\b' "$file"; then
		fail "Raw exec.Command in internal/hook: $file" \
			"hook.FileManager must receive a GitConfigReader interface." \
			"Inject git.ExecClient from cmd/hook_cobra.go instead."
	fi
done < <(find "$ROOT/internal/hook" -name "*.go" -not -name "*_test.go" -print0)

check_pass "internal/hook has no raw exec.Command calls" "$before"

# ---------------------------------------------------------------------------
# Paths — only internal/dirs may call os.UserConfigDir()
# ---------------------------------------------------------------------------
echo ""
echo "Paths: only internal/dirs calls os.UserConfigDir()"
before=$VIOLATIONS

while IFS= read -r -d '' file; do
	if [[ "$file" == "$ROOT/internal/dirs/"* ]]; then continue; fi
	if grep -q 'os\.UserConfigDir()' "$file"; then
		fail "os.UserConfigDir() outside internal/dirs: $file" \
			"All config path resolution must go through internal/dirs.ConfigRoot()," \
			"internal/dirs.ConfigFile(), or internal/dirs.PromptsDir()."
	fi
done < <(find "$ROOT" -name "*.go" -not -name "*_test.go" -not -path "*/vendor/*" -print0)

check_pass "os.UserConfigDir() confined to internal/dirs" "$before"

# ---------------------------------------------------------------------------
# Paths — cmd/ must not import go-toml
# ---------------------------------------------------------------------------
echo ""
echo "Paths: cmd/ does not import go-toml directly"
before=$VIOLATIONS

while IFS= read -r -d '' file; do
	if grep -q '"github.com/pelletier/go-toml' "$file"; then
		fail "go-toml imported in cmd/: $file" \
			"TOML serialisation must stay in internal/config or internal/prompt." \
			"Use store.Format(), store.Marshal(), or store.Unmarshal() instead."
	fi
done < <(find "$ROOT/cmd" -name "*.go" -not -name "*_test.go" -print0)

check_pass "cmd/ has no go-toml imports" "$before"

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo ""
echo "======================================"
if [[ $VIOLATIONS -eq 0 ]]; then
	green "All constitution checks passed (0 violations)"
	exit 0
else
	red "$VIOLATIONS violation(s) found"
	exit 1
fi
