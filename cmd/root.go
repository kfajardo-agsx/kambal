package cmd

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string

	root = &cobra.Command{
		Use:   "file-service",
		Short: "File Service - save and get files",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	root.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ${HOME}/.file-service.yaml")
}

func initConfig() {
	if configFile != "" {
		// load config file if provided via args
		viper.SetConfigFile(configFile)
	} else {
		// else try to load it from the home directory
		home, err := os.UserHomeDir()
		if err != nil {
			log.WithError(err).Error("unable to locate home directory for configuration")
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".file-service.yaml")
	}

	// enable environment vars
	viper.SetEnvPrefix("FILE_SERVICE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Execute is the entrypoint for the application
func Execute(version string) {
	root.Version = version
	if err := root.Execute(); err != nil {
		log.WithError(err).Error("Cannot execute root command")
		os.Exit(1)
	}
}
