package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd returns a cobra command that prints the binary version and build info.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of git-msg",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("git-msg %s\n", Version)
			fmt.Printf("commit: %s\n", Commit)
			fmt.Printf("date:   %s\n", Date)
			fmt.Printf("built by: %s\n", BuiltBy)
		},
	}
}
