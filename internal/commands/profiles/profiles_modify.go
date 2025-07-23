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
		// add a check here if the passed profile name is valid
		path, _ := cmd.Flags().GetString("config")

		cfg, path, err := EnsureConfig(path)
		if err != nil {
			log.WithError(err).Error("Failed to load or create config")
		}

		if modifyName != "" {
			if _, exists := cfg.Profiles[modifyName]; !exists {
				log.Errorf("Profile %s not found", modifyName)
				return
			}
		}

		if modifyName == "" || modifyServiceToken == "" {
			runModifyInteractive(cfg, path)
		} else {
			cfg.Profiles[modifyName] = Profile{
				BaseUrl:      modifyBaseUrl,
				MgmntUrl:     modifyMgmntUrl,
				ServiceToken: modifyServiceToken,
			}

			err := SaveConfig(cfg, path)
			if err != nil {
				log.WithError(err).Error("Failed to save config")
			}
		}

	},
}

func init() {
	profilesModifyCmd.Flags().StringVar(&modifyName, "name", "", "Name of the profile (required)")
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
		log.Error("No profiles found")
		return
	}

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

	ui2.SetQuestions([]cobra_ui.Question{
		{
			Text: fmt.Sprintf("Base URL (current: %s): ", currentBaseURL),
			Handler: func(s string) error {
				if s != "" {
					modifyBaseUrl = s
				} else {
					modifyBaseUrl = selectedProfile.BaseUrl
				}
				return nil
			},
		},
		{
			Text: fmt.Sprintf("Management URL (current: %s): ", currentMgmntURL),
			Handler: func(s string) error {
				if s != "" {
					modifyMgmntUrl = s
				} else {
					modifyMgmntUrl = selectedProfile.MgmntUrl
				}
				return nil
			},
		},
		{
			Text: fmt.Sprintf("Service Token (current: %s): ", currentToken),
			Handler: func(s string) error {
				if s != "" {
					modifyServiceToken = s
				} else {
					modifyServiceToken = selectedProfile.ServiceToken
				}
				return nil
			},
		},
	})

	ui2.RunInteractiveUI()

	cfg.Profiles[modifyName] = Profile{
		BaseUrl:      modifyBaseUrl,
		MgmntUrl:     modifyMgmntUrl,
		ServiceToken: modifyServiceToken,
	}

	err := SaveConfig(cfg, path)
	if err != nil {
		log.WithError(err).Error("Failed to save config")
	}

	log.Infof("Profile %s modified successfully", modifyName)
}
