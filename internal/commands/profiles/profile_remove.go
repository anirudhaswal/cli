package profiles

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	removeName string
)

var profileRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a profile",
	Run: func(cmd *cobra.Command, args []string) {
		if removeName == "" {
			log.Error("You must specify the name as --name")
			return
		}

		path, err := cmd.Flags().GetString("config")
		if err != nil {
			log.WithError(err).Error("Couldn't find the path")
			return
		}

		cfg, path, err := EnsureConfig(path)
		if err != nil {
			log.WithError(err).Error("Failed to load config")
			return
		}

		if _, exists := cfg.Profiles[removeName]; !exists {
			log.Infof("Profile %q does not exist. Use the command 'suprsend profiles list' to see all profiles.", removeName)
			return
		}

		delete(cfg.Profiles, removeName)
		log.Printf("Profile %q has been removed.", removeName)

		if cfg.ActiveProfile == removeName {
			if _, hasDefault := cfg.Profiles["default"]; hasDefault {
				cfg.ActiveProfile = "default"
				log.Warn("Active profile was removed. Falling back to 'default'.")
			} else {
				cfg.ActiveProfile = ""
				log.Warn("Active profile was remove. No active profile set.")
			}
		}

		if err := SaveConfig(cfg, path); err != nil {
			log.WithError(err).Error("Failed to save")
			return
		}
	},
}

func init() {
	profileRemoveCmd.Flags().StringVar(&removeName, "name", "", "Name of the profile to remove")
	ProfilesCmd.AddCommand(profileRemoveCmd)
}
