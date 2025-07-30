package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var genDocsCmd = &cobra.Command{
	Use:    "gendocs [dir]",
	Hidden: true,
	Short:  "Generate CLI documentation in Markdown",
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		return doc.GenMarkdownTree(rootCmd, dir)
	},
}

func init() {
	genDocsCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		command.Flags().MarkHidden("workspace")
		command.Flags().MarkHidden("service-token")
		command.Flags().MarkHidden("output")
		command.Flags().MarkHidden("verbosity")
		command.Flags().MarkHidden("no-color")
		command.Flags().MarkHidden("config")
		command.Parent().HelpFunc()(command, strings)
	})

	rootCmd.AddCommand(genDocsCmd)
}
