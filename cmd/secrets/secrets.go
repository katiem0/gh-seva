package secrets

import (
	createCmd "github.com/katiem0/gh-seva/cmd/secrets/create"
	exportCmd "github.com/katiem0/gh-seva/cmd/secrets/export"
	"github.com/spf13/cobra"
)

func NewCmdSecrets() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "secrets <command> [flags]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Export and Create secrets for an organization and/or repositories.",
		Long:  "Export and Create Actions, Dependabot, and Codespaces secrets for an organization and/or repositories.",
	}
	cmd.Flags().Bool("help", false, "Show help for command")
	cmd.AddCommand(exportCmd.NewCmdExport())
	cmd.AddCommand(createCmd.NewCmdCreate())

	return cmd
}
