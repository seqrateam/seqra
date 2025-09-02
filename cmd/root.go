package cmd

import (
	"fmt"
	"strings"

	"github.com/seqrateam/seqra/internal/globals"
	"github.com/seqrateam/seqra/internal/utils/log"
	"github.com/seqrateam/seqra/internal/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var toolVersion bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "seqra",
	Short: "Seqra Analyzer",
	Long:  `Seqra is a CLI tool that analyzes Java projects to find vulnerabilities`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Set up logging to both console and file
		logFile, logPath, err := log.OpenLogFile()
		globals.LogPath = logPath
		cobra.CheckErr(err)

		if err := log.SetUpLogs(logFile, globals.Config.Log.Verbosity); err != nil {
			return fmt.Errorf("failed to set up logging: %w", err)
		}

		if cmd.Annotations != nil && cmd.Annotations["PrintConfig"] == "true" {
			logrus.Infof("=== Config ===")
			logrus.Infof("Log level: %s", globals.Config.Log.Verbosity)
			if viper.ConfigFileUsed() != "" {
				logrus.Infof("Using config file: %v", viper.ConfigFileUsed())
			}
			logrus.Infof("Logging to file: %s", globals.LogPath)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if toolVersion {
			fmt.Printf("seqra version %s\n", version.Version)
		} else {
			_ = cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logrus.Fatalf("Unexpected error: %s", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&globals.ConfigFile, "config", "", "Path to a config file")

	rootCmd.Flags().BoolVarP(&toolVersion, "version", "v", false, "Print the version information")

	rootCmd.PersistentFlags().StringVar(&globals.Config.Log.Verbosity, "verbosity", logrus.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
	_ = viper.BindPFlag("log.verbosity", rootCmd.PersistentFlags().Lookup("verbosity"))

	rootCmd.PersistentFlags().BoolVarP(&globals.Config.Quiet, "quiet", "q", false, "Suppress interactive console output. (default: false)")
	_ = viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))

	rootCmd.PersistentFlags().StringVar(&globals.Config.Analyzer.Version, "analyzer-version", globals.AnalyzerBindVersion, "Version of seqra analyzer")
	_ = rootCmd.PersistentFlags().MarkHidden("analyzer-version")
	_ = viper.BindPFlag("analyzer.version", rootCmd.PersistentFlags().Lookup("analyzer-version"))

	rootCmd.PersistentFlags().StringVar(&globals.Config.Autobuilder.Version, "autobuilder-version", globals.AutobuilderBindVersion, "Version of seqra autobuilder")
	_ = rootCmd.PersistentFlags().MarkHidden("autobuilder-version")
	_ = viper.BindPFlag("autobuilder.version", rootCmd.PersistentFlags().Lookup("autobuilder-version"))

	rootCmd.PersistentFlags().StringVar(&globals.Config.Github.Token, "github-token", "", "Token for docker image pull from ghcr.io")
	_ = rootCmd.PersistentFlags().MarkHidden("github-token")
	_ = viper.BindPFlag("github.token", rootCmd.PersistentFlags().Lookup("github-token"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if globals.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(globals.ConfigFile)
	}

	viper.SetEnvPrefix("seqra")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	_ = viper.ReadInConfig()
	_ = viper.Unmarshal(&globals.Config)
}

func bindCompileTypeFlag(cmd *cobra.Command) {
	_ = viper.BindPFlag("compile.type", cmd.Flags().Lookup("compile-type"))
}
