package profiles

import (
	"fmt"
	"strings"

	"github.com/sabouaram/cobra_ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	modifyName         string
	modifyBaseUrl      string
	modifyMgmntUrl     string
	modifyServiceToken string
)

var profilesModifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify a profile",
	Long:  "Modify a profile in the configs",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("config")

		cfg, path, err := EnsureConfig(path)
		if err != nil {
			log.WithError(err).Error("Failed to load or create config")
			return
		}
		if modifyName != "" {
			if _, exists := cfg.Profiles[modifyName]; !exists {
				log.Errorf("Profile %s not found", modifyName)
				return
			}
		}

		if modifyName != "" && modifyServiceToken != "" {
			selectedProfile := cfg.Profiles[modifyName]

			if modifyBaseUrl != "" {
				selectedProfile.BaseUrl = modifyBaseUrl
			}
			if modifyMgmntUrl != "" {
				selectedProfile.MgmntUrl = modifyMgmntUrl
			}
			selectedProfile.ServiceToken = modifyServiceToken

			cfg.Profiles[modifyName] = selectedProfile

			err := SaveConfig(cfg, path)
			if err != nil {
				log.WithError(err).Error("Failed to save config")
				return
			}

			log.Infof("Profile %s modified successfully", modifyName)
		} else {
			runModifyInteractive(cfg, path)
		}
	},
}

func init() {
	profilesModifyCmd.Flags().StringVar(&modifyName, "name", "", "Name of the profile to modify")
	profilesModifyCmd.Flags().StringVar(&modifyBaseUrl, "base-url", "", "Base URL")
	profilesModifyCmd.Flags().StringVar(&modifyMgmntUrl, "mgmnt-url", "", "Management URL")
	profilesModifyCmd.Flags().StringVar(&modifyServiceToken, "token", "", "Service Token")
	ProfilesCmd.AddCommand(profilesModifyCmd)
}

func runModifyInteractive(cfg *Config, path string) {
	ui := cobra_ui.New()

	var profileNames []string
	for name := range cfg.Profiles {
		profileNames = append(profileNames, name)
	}

	if len(profileNames) == 0 {
		log.Info("No profiles found. Use the command 'suprsend profiles add' to add a profile.")
		return
	}

	if modifyName == "" {
		ui.SetQuestions([]cobra_ui.Question{
			{
				Text: "Select a profile to modify: ",
				Handler: func(s string) error {
					modifyName = s
					return nil
				},
				Options: profileNames,
			},
		})

		ui.RunInteractiveUI()
	}

	selectedProfile, exists := cfg.Profiles[modifyName]
	if !exists {
		log.Errorf("Profile %s not found", modifyName)
		return
	}

	ui2 := cobra_ui.New()

	currentBaseURL := selectedProfile.BaseUrl
	if currentBaseURL == "" {
		currentBaseURL = "https://hub.suprsend.com"
	}
	currentMgmntURL := selectedProfile.MgmntUrl
	if currentMgmntURL == "" {
		currentMgmntURL = "https://api.suprsend.com"
	}
	currentToken := selectedProfile.ServiceToken
	if currentToken != "" {
		// Mask the token for security
		if len(currentToken) > 4 {
			currentToken = strings.Repeat("*", len(currentToken)-4) + currentToken[len(currentToken)-4:]
		} else {
			currentToken = "****"
		}
	} else {
		currentToken = "not set"
	}

	var questions []cobra_ui.Question

	if modifyBaseUrl == "" {
		questions = append(questions, cobra_ui.Question{
			Text: fmt.Sprintf("Base URL (current: %s, press Enter to keep): ", currentBaseURL),
			Handler: func(s string) error {
				if s != "" {
					modifyBaseUrl = s
				} else {
					modifyBaseUrl = selectedProfile.BaseUrl
				}
				return nil
			},
		})
	}

	if modifyMgmntUrl == "" {
		questions = append(questions, cobra_ui.Question{
			Text: fmt.Sprintf("Management URL (current: %s, press Enter to keep): ", currentMgmntURL),
			Handler: func(s string) error {
				if s != "" {
					modifyMgmntUrl = s
				} else {
					modifyMgmntUrl = selectedProfile.MgmntUrl
				}
				return nil
			},
		})
	}

	if modifyServiceToken == "" {
		questions = append(questions, cobra_ui.Question{
			Text: fmt.Sprintf("Service Token (current: %s, press Enter to keep): ", currentToken),
			Handler: func(s string) error {
				if s != "" {
					modifyServiceToken = s
				} else {
					modifyServiceToken = selectedProfile.ServiceToken
				}
				return nil
			},
		})
	}

	if len(questions) > 0 {
		ui2.SetQuestions(questions)
		ui2.RunInteractiveUI()
	}

	if modifyName == "" {
		log.Error("Profile name is required")
		return
	}
	if modifyServiceToken == "" {
		log.Error("Service token is required")
		return
	}

	updatedProfile := Profile{
		BaseUrl:      modifyBaseUrl,
		MgmntUrl:     modifyMgmntUrl,
		ServiceToken: modifyServiceToken,
	}

	cfg.Profiles[modifyName] = updatedProfile

	err := SaveConfig(cfg, path)
	if err != nil {
		log.WithError(err).Error("Failed to save config")
		return
	}

	log.Infof("Profile %s modified successfully", modifyName)
}
