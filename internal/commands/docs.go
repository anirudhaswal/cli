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
	rootCmd.AddCommand(genDocsCmd)
}
