package profiles

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var profileUseCmd = &cobra.Command{
	Use:   "use",
	Short: "Set the active profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		useName := args[0]
		if useName == "" {
			log.Error("You must specify --name")
		}

		path, err := cmd.Flags().GetString("config")
		if err != nil {
			log.WithError(err).Error("Couldn't find the path")
		}

		cfg, path, err := EnsureConfig(path)
		if err != nil {
			log.WithError(err).Error("Failed to load config")
		}

		if _, exists := cfg.Profiles[useName]; !exists {
			log.Errorf("Profile %q does not exist.", useName)
		}

		cfg.ActiveProfile = useName

		if err := SaveConfig(cfg, path); err != nil {
			log.WithError(err).Error("Failed to save config")
		}

		log.Infof("Active profile set to %q.", useName)
	},
}

func init() {
	ProfilesCmd.AddCommand(profileUseCmd)
}
