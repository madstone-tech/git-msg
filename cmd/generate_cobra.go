package cmd

import (
	"github.com/madstone-tech/git-msg/internal/git"
	"github.com/madstone-tech/git-msg/internal/prompt"
	"github.com/spf13/cobra"
)

// NewGenerateCmd returns the cobra command for generate. Cobra wiring only.
func NewGenerateCmd() *cobra.Command {
	var opts GenerateOptions

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a commit message from staged changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Pull config and secrets injected by PersistentPreRunE.
			if cfg, ok := ConfigFromContext(ctx); ok {
				opts.Cfg = &cfg
			}
			if secrets, ok := SecretsFromContext(ctx); ok {
				opts.Secrets = secrets
			}

			// Construct infrastructure dependencies here, not in Run().
			tmplStore, err := prompt.NewFileStore()
			if err != nil {
				return err
			}
			opts.Templates = tmplStore
			opts.Git = git.ExecClient{}
			// SpinnerFunc and ReviewFunc left nil → Run() uses ui.Spinner / ui.Review.

			return Run(ctx, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Print generated message without committing")
	cmd.Flags().StringVar(&opts.Provider, "provider", "", "Override LLM provider for this run")
	cmd.Flags().StringVar(&opts.Template, "template", "", "Override prompt template for this run")
	cmd.Flags().BoolVar(&opts.HookMode, "hook-mode", false, "")
	cmd.Flags().StringVar(&opts.HookMsgFile, "hook-msg-file", "", "")
	cmd.Flags().StringVar(&opts.HookSource, "hook-source", "", "")
	_ = cmd.Flags().MarkHidden("hook-mode")
	_ = cmd.Flags().MarkHidden("hook-msg-file")
	_ = cmd.Flags().MarkHidden("hook-source")

	return cmd
}
