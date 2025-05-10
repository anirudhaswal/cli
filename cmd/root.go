/*
Copyright © 2025 SuprSend
*/
package cmd

import (
	_ "fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.szostok.io/version/extension"
	"suprsend-cli/util"
)

var (
	cfgFile       string
	outputType    string
	verbosity     string
	serviceToken  string
	noColorOutput bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "suprsend",
	Short: "CLI to interact with SuprSend, a Notification Infrastructure",
	Long: heredoc.Doc(`SuprSend is a robust notification infrastructure that helps you deploy multi-channel product notifications effortlessly and take care of user experience.

	This CLI lets you interact with your SuprSend workspace and do actions like fetching/modifying template, workflows etc.`),
	Version: "0.0.1",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		// 1. Register the 'version' command
		extension.NewVersionCobraCmd(
			// 2. Explicitly enable upgrade notice
			extension.WithUpgradeNotice("suprsend", "cli"),
		),
	)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := SetUpLogs(verbosity, outputType); err != nil {
			return err
		}

		util.InitSDK(viper.GetString("service_token"), viper.GetBool("debug"))
		return nil
	}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.suprsend.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputType, "output", "o", "pretty", "Output Tyle (pretty, yaml, json)")
	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", "info", "Log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringVarP(&serviceToken, "service-token", "s", "", "Service token (default: $SUPRSEND_SERVICE_TOKEN)")
	rootCmd.PersistentFlags().BoolVarP(&noColorOutput, "no-color", "n", false, "Disable color output")

	viper.BindPFlag("service_token", rootCmd.PersistentFlags().Lookup("service-token"))
	viper.BindPFlag("NO_COLOR", rootCmd.PersistentFlags().Lookup("no-color"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".suprsend" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".suprsend")
	}

	viper.AutomaticEnv() // read in environment variables that match
	// load configs from env
	viper.BindEnv("debug", "DEBUG")
	// if NO_COLOR is set, disable color output
	if viper.GetBool("NO_COLOR") {
		color.NoColor = true
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	}
}

// setUpLogs set the log output ans the log level
func SetUpLogs(level string, outputType string) error {
	if outputType == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	if viper.GetBool("debug") {
		verbosity = "debug"
	}
	lvl, err := log.ParseLevel(verbosity)
	if err != nil {
		return errors.Wrap(err, "parsing log level")
	}
	log.SetOutput(os.Stderr)
	log.SetLevel(lvl)
	return nil
}
