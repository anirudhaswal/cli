package event

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
	"github.com/yarlson/pin"
)

var force bool

var eventPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull events from workspace to local directory",
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")

		dirPath := filepath.Join(".", "suprsend", "event")
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if force {
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					fmt.Fprintf(os.Stdout, "Error: Failed to create directory: %v\n", err)
					return
				}
				fmt.Fprintf(os.Stdout, "Succcess: Directory created at %s\n", dirPath)
			} else {
				fmt.Fprintf(os.Stdout, "Directory '%s' does not exist. Do you want to create it?(y/n): ", dirPath)
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.ToLower(strings.TrimSpace(input))

				if input == "y" || input == "yes" {
					err := os.MkdirAll(dirPath, 0755)
					if err != nil {
						fmt.Fprintf(os.Stdout, "Error: Failed to create directory: %v\n", err)
						return
					}
					fmt.Fprintf(os.Stdout, "Success: Directory created at %s\n", dirPath)
				} else {
					fmt.Fprintf(os.Stdout, "Error: Directory not created. Exiting\n")
					return
				}
			}
		} else {
			fmt.Fprintf(os.Stdout, "Directory already exists: %s\n", dirPath)
		}

		var p *pin.Pin
		if !utils.IsOutputPiped() {
			p = pin.New("Loading...",
				pin.WithSpinnerColor(pin.ColorCyan),
				pin.WithTextColor(pin.ColorYellow),
			)
			cancel := p.Start(context.Background())
			defer cancel()
		}

		mgmnt_client := utils.GetSuprSendMgmntClient()
		events_resp, err := mgmnt_client.GetEvents(workspace)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to get events: %v\n", err)
			return
		}
		if p != nil {
			p.Stop(fmt.Sprintf("Pulled %d events", len(events_resp.Results)))
		}

		_, err = WriteEventsToFiles(events_resp, dirPath)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to save events: %v\n", err)
			return
		}
	},
}

func init() {
	eventPullCmd.Flags().StringP("workspace", "w", "staging", "Workspace to pull events from")
	eventPullCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing directory")
	EventCmd.AddCommand(eventPullCmd)
}
