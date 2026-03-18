package cmd

import (
	"github.com/madstone-tech/git-msg/internal/config"
	"github.com/spf13/cobra"
)

// NewConfigCmd returns the cobra command for config management.
// cfgStore is the shared instance created in Execute() — no independent
// NewFileStore() calls are made here.
func NewConfigCmd(cfgStore config.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage git-msg configuration",
	}
	cmd.AddCommand(
		newConfigSetCmd(cfgStore),
		newConfigGetCmd(cfgStore),
		newConfigShowCmd(cfgStore),
	)
	return cmd
}

func newConfigSetCmd(cfgStore config.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "set KEY VALUE",
		Short: "Set a config value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return SetConfig(cmd.Context(), cfgStore, args[0], args[1])
		},
	}
}

func newConfigGetCmd(cfgStore config.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "get KEY",
		Short: "Get a config value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return GetConfig(cmd.Context(), cfgStore, args[0])
		},
	}
}

func newConfigShowCmd(cfgStore config.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Print full resolved configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ShowConfig(cmd.Context(), cfgStore)
		},
	}
}
