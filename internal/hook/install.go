package hook

import (
	"context"
	"errors"
	"os"
	"path/filepath"
)

const hookScript = `#!/bin/sh
# Managed by git-msg. Do not edit manually.
# Skip generation when there is no interactive terminal (e.g. --no-verify,
# non-interactive shells, or CI environments).
if ! [ -t 1 ]; then
  exit 0
fi
exec git-msg generate --hook-mode --hook-msg-file "$1" --hook-source "${2:-}"
`

// GitConfigReader is a narrow interface for reading a single git config value.
// Satisfied by git.Client so the full client can be injected, but the hook
// package only declares the slice it needs — avoiding a dependency on the
// git package.
type GitConfigReader interface {
	RunConfig(ctx context.Context, key string) (string, error)
}

// FileManager implements Manager using the filesystem.
type FileManager struct {
	repoRoot  string
	gitConfig GitConfigReader
}

// NewFileManager returns a FileManager.
// gitConfig is used to resolve core.hooksPath for global installs.
func NewFileManager(repoRoot string, gitConfig GitConfigReader) *FileManager {
	return &FileManager{repoRoot: repoRoot, gitConfig: gitConfig}
}

func (m *FileManager) hookDir(global bool) (string, error) {
	if global {
		val, err := m.gitConfig.RunConfig(context.Background(), "core.hooksPath")
		if err != nil || val == "" {
			return "", errors.New("git config core.hooksPath is not set; configure it first")
		}
		return val, nil
	}
	return filepath.Join(m.repoRoot, ".git", "hooks"), nil
}

func (m *FileManager) Install(global bool) error {
	dir, err := m.hookDir(global)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, "prepare-commit-msg")
	if err := os.WriteFile(path, []byte(hookScript), 0755); err != nil {
		return err
	}
	return os.Chmod(path, 0755)
}

func (m *FileManager) Uninstall(global bool) error {
	dir, err := m.hookDir(global)
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(dir, "prepare-commit-msg"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (m *FileManager) IsInstalled(global bool) (bool, error) {
	dir, err := m.hookDir(global)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(filepath.Join(dir, "prepare-commit-msg"))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return err == nil, err
}
