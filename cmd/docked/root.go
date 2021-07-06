package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Build param: version
var version = "0.0.0"
// Build param: date
var date = "1970-01-01"
// Build param: commit
var commit = ""
// Build param: projectName
var projectName = "docked"

// Global persistent options
var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   projectName,
	Short: "Dockerfile linting.",
	Long: projectName + ` is a Dockerfile linting tool which aims to pull many
best practices and recommendations from multiple sources:

  * OWASP
  * Docker Official Documentation
  * Community recommendations
  * Package manager bug trackers
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}


func main() {
	formattedVersion := fmt.Sprintf("%s (%s) %s", version, commit, date)
	rootCmd.SetVersionTemplate(formattedVersion)
	rootCmd.Version = formattedVersion
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.docked.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	err := viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	cobra.CheckErr(err)
}

// initLogging initializes logging used by the tool.
func initLogging() {
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "info"
	}
	ll, err := logrus.ParseLevel(logLevel)
	if err != nil {
		ll = logrus.ErrorLevel
	}
	logrus.SetLevel(ll)
	logrus.SetOutput(os.Stderr)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath(fmt.Sprintf("/opt/%s", projectName))
		viper.AddConfigPath(fmt.Sprintf("/etc/%s", projectName))
		viper.AddConfigPath(fmt.Sprintf("%s/%s", home, projectName))
		viper.SetConfigType("yaml")
		viper.SetConfigName(".docked")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}
