package cmd

import (
	"github.com/greenstatic/kga/pkg/app"
	"github.com/greenstatic/kga/pkg/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var Generate = &cobra.Command{
	Use:   "generate <app path>",
	Short: "Generate the YAML manifests for a kga app",
	Long: `Generate the YAML manifests for a kga app.

Use the environment variable HELM to specify an alternative path to helm,
otherwise we will use helm and hope it is in your path.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appPath := args[0]

		if !filepath.IsAbs(appPath) {
			wd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			appPath = filepath.Join(wd, appPath)
		}

		log.FatalOnError(app.Generate(appPath))
		log.Info("Successfully generated kga app")
	},
}
