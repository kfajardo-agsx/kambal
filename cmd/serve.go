package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	restAPI := createRestAPI()
	if err := restAPI.Run(); err != nil {
		log.WithError(err).Error("REST API terminated")
	}
}
