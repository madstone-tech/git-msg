package hook_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/madstone0-0/git-msg/internal/hook"
)

// noopGitConfig satisfies hook.GitConfigReader for tests that only exercise
// local (non-global) hook paths.
type noopGitConfig struct{}

func (noopGitConfig) RunConfig(_ context.Context, _ string) (string, error) { return "", nil }

func TestFileManager_Install(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, ".git", "hooks")
	os.MkdirAll(hooksDir, 0755)

	mgr := hook.NewFileManager(dir, noopGitConfig{})
	if err := mgr.Install(false); err != nil {
		t.Fatal(err)
	}

	hookPath := filepath.Join(hooksDir, "prepare-commit-msg")
	info, err := os.Stat(hookPath)
	if err != nil {
		t.Fatalf("hook file not created: %v", err)
	}
	if info.Mode()&0100 == 0 {
		t.Error("hook file is not executable")
	}
}

func TestFileManager_Install_Idempotent(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git", "hooks"), 0755)
	mgr := hook.NewFileManager(dir, noopGitConfig{})
	if err := mgr.Install(false); err != nil {
		t.Fatal(err)
	}
	if err := mgr.Install(false); err != nil {
		t.Fatalf("second install failed: %v", err)
	}
}

func TestFileManager_Uninstall(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git", "hooks"), 0755)
	mgr := hook.NewFileManager(dir, noopGitConfig{})
	_ = mgr.Install(false)
	if err := mgr.Uninstall(false); err != nil {
		t.Fatal(err)
	}
	_, err := os.Stat(filepath.Join(dir, ".git", "hooks", "prepare-commit-msg"))
	if !os.IsNotExist(err) {
		t.Error("hook file still exists after uninstall")
	}
}

func TestFileManager_Uninstall_NotPresent(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git", "hooks"), 0755)
	mgr := hook.NewFileManager(dir, noopGitConfig{})
	if err := mgr.Uninstall(false); err != nil {
		t.Fatalf("uninstall on absent hook failed: %v", err)
	}
}
