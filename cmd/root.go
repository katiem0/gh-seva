package cmd

import (
	secretsCmd "github.com/katiem0/gh-seva/cmd/secrets"
	variablesCmd "github.com/katiem0/gh-seva/cmd/variables"
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {

	cmdRoot := &cobra.Command{
		Use:   "seva <command> <subcommand> [flags]",
		Short: "Export and Create secrets and variables.",
		Long:  "Export and Create secrets and variables for an organization and/or repositories.",
	}
	cmdRoot.PersistentFlags().Bool("help", false, "Show help for command")

	cmdRoot.AddCommand(secretsCmd.NewCmdSecrets())
	cmdRoot.AddCommand(variablesCmd.NewCmdVariables())
	cmdRoot.CompletionOptions.DisableDefaultCmd = true
	cmdRoot.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	return cmdRoot
}
