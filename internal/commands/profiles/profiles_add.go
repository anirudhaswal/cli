package profiles

import (
	"fmt"

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
			return
		}

		if addName != "" && addServiceToken != "" {
			cfg.Profiles[addName] = Profile{
				BaseUrl:      addBaseUrl,
				MgmntUrl:     addMgmntUrl,
				ServiceToken: addServiceToken,
			}

			err := SaveConfig(cfg, path)
			if err != nil {
				log.WithError(err).Error("Failed to save config")
				return
			}

			log.Infof("Profile %s added successfully", addName)
		} else {
			runAddInteractive(cfg, path)
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

	var questions []cobra_ui.Question

	if addName == "" {
		questions = append(questions, cobra_ui.Question{
			Text: "Name: ",
			Handler: func(s string) error {
				if s == "" {
					return fmt.Errorf("profile name cannot be empty")
				}
				if _, exists := cfg.Profiles[s]; exists {
					return fmt.Errorf("profile '%s' already exists", s)
				}
				addName = s
				return nil
			},
		})
	}

	if addServiceToken == "" {
		questions = append(questions, cobra_ui.Question{
			Text: "Service Token: ",
			Handler: func(s string) error {
				if s == "" {
					return fmt.Errorf("service token cannot be empty")
				}
				addServiceToken = s
				return nil
			},
		})
	}

	if addBaseUrl == "" {
		questions = append(questions, cobra_ui.Question{
			Text: "Base URL (optional, press Enter for default): ",
			Handler: func(s string) error {
				addBaseUrl = s
				return nil
			},
		})
	}

	if addMgmntUrl == "" {
		questions = append(questions, cobra_ui.Question{
			Text: "Management URL (optional, press Enter for default): ",
			Handler: func(s string) error {
				addMgmntUrl = s
				return nil
			},
		})
	}

	if len(questions) > 0 {
		ui.SetQuestions(questions)
		ui.RunInteractiveUI()
	}

	cfg.Profiles[addName] = Profile{
		BaseUrl:      addBaseUrl,
		MgmntUrl:     addMgmntUrl,
		ServiceToken: addServiceToken,
	}

	err := SaveConfig(cfg, path)
	if err != nil {
		log.WithError(err).Error("Failed to save config")
		return
	}

	log.Infof("Profile %s added successfully", addName)
}
