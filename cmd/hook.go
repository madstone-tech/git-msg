package cmd

import (
	"context"
	"fmt"

	"github.com/madstone0-0/git-msg/internal/hook"
)

// HookOptions holds inputs for hook install/uninstall commands.
type HookOptions struct {
	Global  bool
	Manager hook.Manager
}

// InstallHook installs the prepare-commit-msg hook.
func InstallHook(ctx context.Context, opts HookOptions) error {
	if err := opts.Manager.Install(opts.Global); err != nil {
		return fmt.Errorf("Error: hook install failed\n  → %w", err)
	}
	scope := "local"
	if opts.Global {
		scope = "global"
	}
	fmt.Printf("Hook installed (%s).\n", scope)
	return nil
}

// UninstallHook removes the prepare-commit-msg hook.
func UninstallHook(ctx context.Context, opts HookOptions) error {
	if err := opts.Manager.Uninstall(opts.Global); err != nil {
		return fmt.Errorf("Error: hook uninstall failed\n  → %w", err)
	}
	fmt.Println("Hook uninstalled.")
	return nil
}
