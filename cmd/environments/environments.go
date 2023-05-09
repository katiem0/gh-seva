package environments

import (
	createCmd "github.com/katiem0/gh-seva/cmd/environments/create"
	exportCmd "github.com/katiem0/gh-seva/cmd/environments/export"
	"github.com/spf13/cobra"
)

func NewCmdEnvs() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "envs <command>",
		Short: "Export and Create environments and metadata.",
		Long:  "Export and Create environments and metadata for repos in an organization or single repository.",
	}
	cmd.Flags().Bool("help", false, "Show help for command")

	cmd.AddCommand(exportCmd.NewCmdExport())
	cmd.AddCommand(createCmd.NewCmdCreate())

	return cmd
}
