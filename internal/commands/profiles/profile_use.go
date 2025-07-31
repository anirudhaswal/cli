package profiles

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	useName string
)

var profileUseCmd = &cobra.Command{
	Use:   "use",
	Short: "Set the active profile",
	Run: func(cmd *cobra.Command, args []string) {
		if useName == "" {
			log.Error("You must specify --name")
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

		if _, exists := cfg.Profiles[useName]; !exists {
			log.Info("Create a profile first with 'suprsend profiles add'")
			return
		}

		cfg.ActiveProfile = useName

		if err := SaveConfig(cfg, path); err != nil {
			log.WithError(err).Error("Failed to save config")
			return
		}

		log.Infof("Active profile set to %q.", useName)
	},
}

func init() {
	profileUseCmd.Flags().StringVarP(&useName, "name", "", "", "Profile name to set as active.")
	ProfilesCmd.AddCommand(profileUseCmd)
}
