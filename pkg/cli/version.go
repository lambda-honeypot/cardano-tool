package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Cardano Tool",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("cardano-tool:v0.1.0")
	},
}
