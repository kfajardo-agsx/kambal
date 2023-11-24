package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serve = &cobra.Command{
	Use:   "serve",
	Short: "run the server",
	Run:   runServe,
}

func init() {
	root.AddCommand(serve)
}

func runServe(c *cobra.Command, args []string) {
	runMigrate(c, args)
	log.Info("========================================")
	log.Info("File Service")
	log.Info("========================================")

	setLogLevel()

	restAPI := createRestAPI()

	log.Info("========================================")
	log.Info("Starting API Server")
	log.Info("========================================")

	if err := restAPI.Run(); err != nil {
		log.WithError(err).Error("REST API terminated")
	}
}

func setLogLevel() {
	logLevel := viper.GetString("log.level")
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithError(err).Errorf("log level %s is invalid", logLevel)
	}
	log.SetLevel(level)
	log.Infof("Log Level Set to: %s", level)
}
