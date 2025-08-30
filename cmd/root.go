package cmd

import (
	"fmt"

	"github.com/seqra/seqra/internal/globals"
	"github.com/seqra/seqra/internal/utils"
	"github.com/seqra/seqra/internal/utils/log"
	"github.com/seqra/seqra/internal/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var toolVersion bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "seqra",
	Short: "Seqra Analyzer",
	Long:  `Seqra is a CLI tool that analyzes Java projects to find vulnerabilities`,
	Run: func(cmd *cobra.Command, args []string) {
		if toolVersion {
			fmt.Printf("seqra version %s\n", version.Version)
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.seqra/config.yaml)")

	rootCmd.Flags().BoolVarP(&toolVersion, "version", "v", false, "Print the version information")

	rootCmd.PersistentFlags().StringVar(&globals.VerboseLevel, "verbosity", logrus.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
	_ = viper.BindPFlag("verbosity", rootCmd.PersistentFlags().Lookup("verbosity"))

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Set up logging to both console and file
		logFile, logPath, err := log.OpenLogFile()
		globals.LogPath = logPath
		cobra.CheckErr(err)

		if err := log.SetUpLogs(logFile, globals.VerboseLevel); err != nil {
			return fmt.Errorf("failed to set up logging: %w", err)
		}

		if viper.ConfigFileUsed() != "" {
			logrus.Infof("Using config file: %v", viper.ConfigFileUsed())
		}
		return nil
	}

	rootCmd.PersistentFlags().BoolVarP(&globals.Quiet, "quiet", "q", false, "Suppress interactive console output. (default: false)")
	_ = viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))

	rootCmd.PersistentFlags().StringVar(&globals.AnalyzerVersion, "analyzer-version", globals.AnalyzerBindVersion, "Version of seqra analyzer")
	_ = rootCmd.PersistentFlags().MarkHidden("analyzer-version")

	rootCmd.PersistentFlags().StringVar(&globals.AutobuilderVersion, "autobuilder-version", globals.AutobuilderBindVersion, "Version of seqra autobuilder")
	_ = rootCmd.PersistentFlags().MarkHidden("autobuilder-version")

	rootCmd.PersistentFlags().StringVar(&globals.GithubToken, "github-token", "", "Token for docker image pull from ghcr.io")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		seqraHomePath, err := utils.GetSeqraHome()
		cobra.CheckErr(err)
		viper.AddConfigPath(seqraHomePath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
		_ = viper.SafeWriteConfig()
	}

	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig()
	if err != nil {
		// Only error if it's not a "config file not found" error
		_, notFound := err.(viper.ConfigFileNotFoundError)
		if !notFound {
			cobra.CheckErr(err)
		}
	}
}
