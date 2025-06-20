package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var genConfigCmd = &cobra.Command{
	Use:   "generate-config [path]",
	Short: "Generate config file in yaml format (default: $HOME/suprsend.yaml)",
	RunE: func(cmd *cobra.Command, args []string) error {
		// if no args are provided, ask for filepath input with default value $HOME/suprsend.yaml
		path := ""
		serviceToken := ""
		if len(args) == 0 {
			fmt.Println("Enter the path to save the config file (default: $HOME/suprsend.yaml)")
			fmt.Scanln(&path)
			if path == "" {
				path = os.Getenv("HOME") + "/suprsend.yaml"
			}
			// check if the file exists
			if _, err := os.Stat(path); err == nil {
				fmt.Println("File already exists")
				return nil
			}
			// ask for service token
			fmt.Println("Enter the service token")
			fmt.Scanln(&serviceToken)
		}
		// create the file with the provided service token
		filecontent := fmt.Sprintf(`service_token: %s`, serviceToken)
		os.WriteFile(path, []byte(filecontent), 0o644)
		fmt.Println("Config file created at", path)
		return nil
	},
}

func init() {
	genConfigCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		// Hide flag for this command
		command.Flags().MarkHidden("config")
		command.Flags().MarkHidden("output")
		command.Flags().MarkHidden("service-token")
		// Call parent help func
		command.Parent().HelpFunc()(command, strings)
	})

	rootCmd.AddCommand(genConfigCmd)
}
