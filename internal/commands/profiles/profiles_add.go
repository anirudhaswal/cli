package profiles

import (
	"github.com/sabouaram/cobra_ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	addName         string
	addBaseUrl      string
	addMgmntUrl     string
	addServiceToken string
)

var profilesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new profile",
	Long:  "Add a new profile to the configs",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("config")

		cfg, path, err := EnsureConfig(path)
		if err != nil {
			log.WithError(err).Error("Failed to load or create config")
		}

		if addName == "" || addServiceToken == "" {
			runAddInteractive(cfg, path)
		} else {
			cfg.Profiles[addName] = Profile{
				BaseUrl:      addBaseUrl,
				MgmntUrl:     addMgmntUrl,
				ServiceToken: addServiceToken,
			}

			err := SaveConfig(cfg, path)
			if err != nil {
				log.WithError(err).Error("Failed to save config")
			}
		}

	},
}

func init() {
	profilesAddCmd.Flags().StringVar(&addName, "name", "", "Name of the profile (required)")
	profilesAddCmd.Flags().StringVar(&addBaseUrl, "base-url", "", "Base URL")
	profilesAddCmd.Flags().StringVar(&addMgmntUrl, "mgmnt-url", "", "Management URL")
	profilesAddCmd.Flags().StringVar(&addServiceToken, "token", "", "Service token")
	ProfilesCmd.AddCommand(profilesAddCmd)
}

func runAddInteractive(cfg *Config, path string) {
	ui := cobra_ui.New()

	ui.SetQuestions([]cobra_ui.Question{
		{
			Text: "Name: ",
			Handler: func(s string) error {
				addName = s
				return nil
			},
		},
		{
			Text: "Base URL: ",
			Handler: func(s string) error {
				addBaseUrl = s
				return nil
			},
		},
		{
			Text: "Management URL: ",
			Handler: func(s string) error {
				addMgmntUrl = s
				return nil
			},
		},
		{
			Text: "Service Token: ",
			Handler: func(s string) error {
				addServiceToken = s
				return nil
			},
		},
	})

	ui.RunInteractiveUI()

	cfg.Profiles[addName] = Profile{
		BaseUrl:      addBaseUrl,
		MgmntUrl:     addMgmntUrl,
		ServiceToken: addServiceToken,
	}

	err := SaveConfig(cfg, path)
	if err != nil {
		log.WithError(err).Error("Failed to save config")
	}

	log.Infof("Profile %s added successfully", addName)
}
