package cmd

import (
	"github.com/madstone0-0/git-msg/internal/git"
	"github.com/madstone0-0/git-msg/internal/hook"
	"github.com/spf13/cobra"
)

// NewHookCmd returns the cobra command for hook management.
func NewHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook",
		Short: "Manage git hook lifecycle",
	}
	cmd.AddCommand(newHookInstallCmd(), newHookUninstallCmd())
	return cmd
}

func newHookInstallCmd() *cobra.Command {
	var global bool
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install prepare-commit-msg hook",
		RunE: func(cmd *cobra.Command, args []string) error {
			gitClient := git.ExecClient{}
			root, err := gitClient.RepoRoot(cmd.Context())
			if err != nil {
				return err
			}
			return InstallHook(cmd.Context(), HookOptions{
				Global:  global,
				Manager: hook.NewFileManager(root, gitClient),
			})
		},
	}
	cmd.Flags().BoolVar(&global, "global", false, "Install hook globally via git core.hooksPath")
	return cmd
}

func newHookUninstallCmd() *cobra.Command {
	var global bool
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove prepare-commit-msg hook",
		RunE: func(cmd *cobra.Command, args []string) error {
			gitClient := git.ExecClient{}
			root, err := gitClient.RepoRoot(cmd.Context())
			if err != nil {
				return err
			}
			return UninstallHook(cmd.Context(), HookOptions{
				Global:  global,
				Manager: hook.NewFileManager(root, gitClient),
			})
		},
	}
	cmd.Flags().BoolVar(&global, "global", false, "Uninstall hook globally")
	return cmd
}
