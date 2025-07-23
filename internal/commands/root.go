/*
Copyright © 2025 SuprSend
*/
package commands

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/suprsend/cli/internal/commands/profiles"
	workflow "github.com/suprsend/cli/internal/commands/workflow"
	"github.com/suprsend/cli/internal/config"
	"github.com/suprsend/cli/internal/utils"
	"go.szostok.io/version/extension"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "suprsend",
	Short: "CLI to interact with SuprSend, a Notification Infrastructure",
	Long: heredoc.Doc(`SuprSend is a robust notification infrastructure that helps you deploy multi-channel product notifications effortlessly and take care of user experience.

	This CLI lets you interact with your SuprSend workspace and do actions like fetching/modifying template, workflows etc.`),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	conf := config.Cfg
	rootCmd.PersistentFlags().StringVarP(&conf.Workspace, "workspace", "w", "staging", "Workspace to use")
	rootCmd.PersistentFlags().StringVar(&conf.CfgFile, "config", "", "config file (default: $HOME/suprsend.yaml)")
	rootCmd.PersistentFlags().StringVarP(&conf.OutputType, "output", "o", "pretty", "Output Tyle (pretty, yaml, json)")
	rootCmd.PersistentFlags().StringVarP(&conf.Verbosity, "verbosity", "v", "info", "Log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringVarP(&conf.ServiceToken, "service-token", "s", "", "Service token (default: $SUPRSEND_SERVICE_TOKEN)")
	rootCmd.PersistentFlags().BoolVarP(&conf.NoColorOutput, "no-color", "n", false, "Disable color output (default: $NO_COLOR)")

	viper.BindPFlag("service_token", rootCmd.PersistentFlags().Lookup("service-token"))
	viper.BindPFlag("NO_COLOR", rootCmd.PersistentFlags().Lookup("no-color"))
	//
	cobra.OnInitialize(func() {
		config.InitConfig(conf.CfgFile)
	})
	rootCmd.AddCommand(
		// 1. Register the 'version' command
		extension.NewVersionCobraCmd(
			// 2. Explicitly enable upgrade notice
			extension.WithUpgradeNotice("suprsend", "cli"),
		),
	)

	rootCmd.AddCommand(profiles.ProfilesCmd)

	workflow.WorkflowCmd.PersistentFlags().IntP("limit", "l", 20, "Limit the number of workflows to list")
	workflow.WorkflowCmd.PersistentFlags().IntP("offset", "f", 0, "Offset the number of workflows to list (default: 0)")
	workflow.WorkflowCmd.PersistentFlags().StringP("mode", "m", "live", "Mode to list workflows (draft, live)")
	syncCmd.PersistentFlags().IntP("limit", "l", 20, "Limit the number of workflows to list")
	syncCmd.PersistentFlags().IntP("offset", "f", 0, "Offset the number of workflows to list (default: 0)")
	syncCmd.PersistentFlags().StringP("mode", "m", "live", "Mode to list workflows (draft, live)")

	rootCmd.AddCommand(workflow.WorkflowCmd)
	rootCmd.AddCommand(syncCmd)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := config.SetUpLogs(); err != nil {
			return err
		}
		// check the subcommand and return if it is generate-config or gendocs
		if cmd.Name() == "generate-config" || cmd.Name() == "gendocs" {
			return nil
		}

		cfg, _, err := profiles.EnsureConfig(conf.CfgFile)
		if err != nil {
			return err
		}

		activeProfile, exists := cfg.Profiles[cfg.ActiveProfile]
		if !exists {
			return fmt.Errorf("active profile '%s' not found", cfg.ActiveProfile)
		}

		// flag > profile > env
		serviceToken := viper.GetString("service_token")
		if serviceToken == "" {
			if activeProfile.ServiceToken != "" {
				serviceToken = activeProfile.ServiceToken
			} else {
				serviceToken = os.Getenv("SUPRSEND_SERVICE_TOKEN")
			}
		}

		conf.ServiceToken = serviceToken

		utils.InitSDKWithUrls(
			viper.GetString("service_token"),
			activeProfile.GetResolvedBaseUrl(),
			activeProfile.GetResolvedMgmntUrl(),
			viper.GetBool("debug"),
		)
		return nil
	}
}
