// setup user subcommand

package commands

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// rootCmd.AddCommand(userCmd)
}

// get user details
var getUserCmd = &cobra.Command{
	Use:   "get",
	Short: "Get user details",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Distinct ID is required")
		}
		distinctId := args[0]
		// Get flag values of distinct id and workspace
		workspace, _ := cmd.Flags().GetString("workspace")
		log.Debug("Getting user details for distinct ID: ", distinctId, " in workspace: ", workspace)
		mgmnt_client := utils.GetSuprSendMgmntClient()
		client, err := mgmnt_client.GetWorkspaceClient(workspace)
		if err != nil {
			log.Fatal(err)
		}
		user, err := client.Users.Get(context.Background(), distinctId)
		if err != nil {
			log.Fatal(err)
		}
		outputType, _ := cmd.Flags().GetString("output")
		utils.OutputData(user, outputType)
	},
}

func init() {
	getUserCmd.Flags().String("workspace", "", "Workspace name")
	getUserCmd.MarkFlagRequired("workspace")
	userCmd.AddCommand(getUserCmd)
}
