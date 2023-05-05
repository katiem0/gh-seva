package variables

import (
	createCmd "github.com/katiem0/gh-seva/cmd/variables/create"
	exportCmd "github.com/katiem0/gh-seva/cmd/variables/export"
	"github.com/spf13/cobra"
)

func NewCmdVariables() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "variables <command>",
		Short: "Export and Create variables for an organization and/or repositories.",
		Long:  "Export and Create Actions variables for an organization and/or repositories.",
	}
	cmd.Flags().Bool("help", false, "Show help for command")

	cmd.AddCommand(exportCmd.NewCmdExport())
	cmd.AddCommand(createCmd.NewCmdCreate())

	return cmd
}
