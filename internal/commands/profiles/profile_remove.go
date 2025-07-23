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
		}

		path, err := cmd.Flags().GetString("config")
		if err != nil {
			log.WithError(err).Error("Couldn't find the path")
		}

		cfg, path, err := EnsureConfig(path)
		if err != nil {
			log.WithError(err).Error("Failed to load config")
		}

		if _, exists := cfg.Profiles[removeName]; !exists {
			log.Errorf("Profile %q does not exist.", removeName)
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
		}
	},
}

func init() {
	profileRemoveCmd.Flags().StringVar(&removeName, "name", "", "Name of the profile to remove")
	ProfilesCmd.AddCommand(profileRemoveCmd)
}
