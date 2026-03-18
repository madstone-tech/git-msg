package cmd

import (
	"github.com/madstone0-0/git-msg/internal/prompt"
	"github.com/spf13/cobra"
)

// NewPromptCmd returns the cobra command for prompt management.
func NewPromptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prompt",
		Short: "Manage prompt templates",
	}
	cmd.AddCommand(
		newPromptListCmd(),
		newPromptShowCmd(),
		newPromptEditCmd(),
		newPromptResetCmd(),
	)
	return cmd
}

func newPromptListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available prompt templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := newPromptStore()
			if err != nil {
				return err
			}
			return ListPrompts(cmd.Context(), store)
		},
	}
}

func newPromptShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show [name]",
		Short: "Show a prompt template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := newPromptStore()
			if err != nil {
				return err
			}
			return ShowPrompt(cmd.Context(), store, args[0])
		},
	}
}

func newPromptEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit [name]",
		Short: "Edit a prompt template in $EDITOR",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := newPromptStore()
			if err != nil {
				return err
			}
			return EditPrompt(cmd.Context(), store, args[0])
		},
	}
}

func newPromptResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset [name]",
		Short: "Reset a template to its embedded default",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := newPromptStore()
			if err != nil {
				return err
			}
			return ResetPrompt(cmd.Context(), store, args[0])
		},
	}
}

func newPromptStore() (*prompt.FileStore, error) {
	return prompt.NewFileStore()
}
