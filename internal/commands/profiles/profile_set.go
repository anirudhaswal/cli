package profiles

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	setName         string
	setBaseURL      string
	setMgmntURL     string
	setServiceToken string
)

var setProfileCmd = &cobra.Command{
	Use:   "set",
	Short: "Create or update a profile",
	Long:  "Create or update a profile and set it as active",
	Run: func(cmd *cobra.Command, args []string) {
		if setName == "" {
			log.Error("You must specify --name")
		}

		path, err := cmd.Flags().GetString("config")
		if err != nil {
			log.WithError(err).Error("Couldn't find the path")
		}
		cfg, path, err := EnsureConfig(path)
		if err != nil {
			log.WithError(err).Error("Failed to load or create config")
		}

		profile, exists := cfg.Profiles[setName]
		if exists {
			log.Infof("Updating existing profile %q", setName)
		} else {
			profile = Profile{}
			log.Infof("Creating new profile %q", setName)
		}

		if cmd.Flags().Changed("base-url") {
			profile.BaseUrl = setBaseURL
		}
		if cmd.Flags().Changed("mgmnt-url") {
			profile.MgmntUrl = setMgmntURL
		}
		if cmd.Flags().Changed("token") {
			profile.ServiceToken = setServiceToken
		}

		cfg.Profiles[setName] = profile
		cfg.ActiveProfile = setName

		if err := SaveConfig(cfg, path); err != nil {
			log.WithError(err).Error("Failed to save config")
		}

		log.Infof("Profile %q is now active.", setName)
	},
}

func init() {
	setProfileCmd.Flags().StringVar(&setName, "name", "", "Name of the profile (required)")
	setProfileCmd.Flags().StringVar(&setBaseURL, "base-url", "", "Base URL")
	setProfileCmd.Flags().StringVar(&setMgmntURL, "mgmnt-url", "", "Management URL")
	setProfileCmd.Flags().StringVar(&setServiceToken, "token", "", "Service token")
	ProfilesCmd.AddCommand(setProfileCmd)
}
