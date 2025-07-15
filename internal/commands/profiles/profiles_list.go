package profiles

import (
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listProfilesCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all profiles",
	Long:  "Lists all profiles from the config",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("config")
		if err != nil {
			log.WithError(err).Error("Couldn't find the path")
		}
		cfg, _, err := EnsureConfig(path)
		if err != nil {
			log.WithError(err)
		}

		var names []string
		for name := range cfg.Profiles {
			names = append(names, name)
		}
		sort.Strings(names)

		if len(names) == 0 {
			log.Warn("No profiles found.")
		}

		fmt.Println("Profiles:")
		for _, name := range names {
			active := ""
			if name == cfg.ActiveProfile {
				active = "(active)"
			}
			fmt.Printf(" - %s %s\n", name, active)
		}
	},
}

func init() {
	ProfilesCmd.AddCommand(listProfilesCmd)
}
