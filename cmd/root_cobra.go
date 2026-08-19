package cmd

import (
	"github.com/madstone-tech/git-msg/internal/config"
	"github.com/madstone-tech/git-msg/internal/secret"
	"github.com/spf13/cobra"
)

// Execute is the entry point called from main.go.
func Execute() error {
	cfgStore, err := config.NewFileStore()
	if err != nil {
		return err
	}
	secrets := secret.KeychainStore{}

	root := newRootCmd(cfgStore, secrets)
	return root.Execute()
}

// newRootCmd builds the full command tree with PersistentPreRunE wired in.
func newRootCmd(cfgStore config.Store, secrets secret.SecretStore) *cobra.Command {
	root := &cobra.Command{
		Use:   "git-msg",
		Short: "AI-assisted git commit message generator",
		// Version is set via -ldflags; expose it here so --version works.
		Version: Version,
		// PersistentPreRunE runs before every subcommand's RunE.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "help" || cmd.Name() == "version" {
				return nil
			}
			cfg, err := EnsureConfig(cmd.Context(), cfgStore, secrets)
			if err != nil {
				return err
			}
			ctx := WithConfig(cmd.Context(), cfg)
			ctx = WithSecrets(ctx, secrets)
			cmd.SetContext(ctx)
			return nil
		},
	}

	root.AddCommand(
		NewGenerateCmd(),
		NewConfigCmd(cfgStore),
		NewPromptCmd(),
		NewHookCmd(),
		NewVersionCmd(),
	)

	return root
}
