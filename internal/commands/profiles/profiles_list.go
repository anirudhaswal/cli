package profiles

import (
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/suprsend/cli/internal/utils"
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

		outputType, _ := cmd.Flags().GetString("output")

		// Check if any profile has BaseUrl, MgmntUrl, or ServiceToken values
		hasBaseUrl := false
		hasMgmntUrl := false
		hasServiceToken := false

		for _, name := range names {
			profile := cfg.Profiles[name]
			if profile.BaseUrl != "" {
				hasBaseUrl = true
			}
			if profile.MgmntUrl != "" {
				hasMgmntUrl = true
			}
			if profile.ServiceToken != "" {
				hasServiceToken = true
			}
		}

		if hasBaseUrl || hasMgmntUrl || hasServiceToken {
			var profileData []ProfileListItem

			for _, name := range names {
				profile := cfg.Profiles[name]
				isActive := "no"
				if name == cfg.ActiveProfile {
					isActive = "yes"
				}

				profileData = append(profileData, ProfileListItem{
					Name:         name,
					Active:       isActive,
					BaseUrl:      profile.BaseUrl,
					MgmntUrl:     profile.MgmntUrl,
					ServiceToken: profile.ServiceToken,
				})
			}

			utils.OutputData(profileData, outputType)
		} else {
			var profileData []SimpleProfileListItem

			for _, name := range names {
				isActive := "no"
				if name == cfg.ActiveProfile {
					isActive = "yes"
				}

				profileData = append(profileData, SimpleProfileListItem{
					Name:   name,
					Active: isActive,
				})
			}

			utils.OutputData(profileData, outputType)
		}
	},
}

func init() {
	ProfilesCmd.AddCommand(listProfilesCmd)
}
