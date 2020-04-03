package cmd

import (
	"github.com/greenstatic/kga/internal/layout"
	"github.com/greenstatic/kga/pkg/config"
	"github.com/greenstatic/kga/pkg/log"
	"github.com/spf13/cobra"
	"path/filepath"
)

var Create = &cobra.Command{
	Use:   "create <app name>",
	Short: "Create a kga app",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		appType := config.AppType(cmd.Flag("type").Value.String())
		if appType != config.AppTypeHelm && appType != config.AppTypeManifest {
			log.Fatal("Unrecognized type, type can only be: helm or manifest")
		}

		wd, err := filepath.Abs(cmd.Flag("output").Value.String())
		if err != nil {
			log.Fatal(err)
		}

		appPath := filepath.Join(wd, appName)
		log.Infof("Creating app in: %s", appPath)

		if err := layout.CreateAppDir(wd, appName); err != nil {
			log.Fatal(err)
		}

		log.Infof("Creating %s/kga.yaml", appPath)
		if err := config.KgaYamlCreate(appType, appName, appPath); err != nil {
			log.Fatalf("Failed to create kga.yaml file, error: %s", err)
		}

		if appType == config.AppTypeHelm {
			// Create helm_values.yaml file
			log.Infof("Creating %s/helm_values.yaml", appPath)
			if err := layout.CreateHelmValuesFile(appPath); err != nil {
				log.Fatal("Failed to create helm_values.yaml file")
			}
		}

		log.Info("Successfully created kga app, now fill out the kga.yaml file and then run `kga generate <name>`")
	},
}

func init() {
	Create.Flags().StringP("type", "t", "helm", "Type of app [helm/manifest]")
	Create.Flags().StringP("output", "o", "./", "Directory where to create the app")
}
