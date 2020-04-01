package cmd

import (
	"github.com/greenstatic/kga/pkg/log"
	"github.com/greenstatic/kga/pkg/version"
	"github.com/spf13/cobra"
)

var Version = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Version: %s", version.String())
	},
}
