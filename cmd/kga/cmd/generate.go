package cmd

import (
	"github.com/greenstatic/kga/internal/generate"
	"github.com/greenstatic/kga/internal/layout"
	"github.com/greenstatic/kga/pkg/config"
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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appPath := args[0]

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		appPath = filepath.Join(wd, appPath)

		// Verify we have a kga app
		isKgaApp, err := layout.PathIsKgaApp(appPath)
		if err != nil {
			log.Fatal(err)
		}
		if !isKgaApp {
			log.Fatal("App path is not a kga app")
		}

		// Verify kga.yaml file
		if err := config.VerifyKgaFile(filepath.Join(appPath, "kga.yaml")); err != nil {
			log.Error(err)
			log.Fatal("Bad configuration!")
		}

		// Read kga.yaml
		log.Info("Reading kga.yaml")
		kgaConfig, err := config.ParseFile(filepath.Join(appPath, "kga.yaml"))
		if err != nil {
			log.Fatal(err)
		}

		// Delete base directory if one exists
		log.Info("Removing " + filepath.Join(appPath, "base"))
		if err := layout.RemoveAppBaseDir(appPath); err != nil {
			log.Fatal(err)
		}

		if kgaConfig.AppType() == config.AppTypeHelm {
			log.Info("Running Helm manifest generation")
			generate.CreateHelmChartManifests(appPath, kgaConfig.Spec.Namespace, kgaConfig.Spec.Helm)

		} else {
			log.Info("Running URL manifest generation")
			generate.DownloadManifestFiles(appPath, kgaConfig.Spec.Manifest)
		}

		if kgaConfig.HasExcludeSpec() {
			log.Info("Running kga manifest exclusion spec")
			if err := generate.RemoveExcludedBaseManifests(appPath, kgaConfig.Spec.Exclude); err != nil {
				log.Error(err)
				log.Fatal("Failed to run kga exclusion spec")
			}
		}

		log.Info("Generating base/kustomization.yaml")
		if err := layout.CreateBaseKustomization(appPath); err != nil {
			log.Fatal(err)
		}

		overlayExists, err := layout.OverlayExists(appPath)
		if err != nil {
			log.Fatal(err)
		}

		if overlayExists {
			log.Info("Overlay exists, leaving it alone")
		} else {
			log.Info("Creating overlay directory structure")

			if err := layout.OverlayCreateGeneralLayout(appPath, kgaConfig.Spec.Namespace); err != nil {
				log.Fatal(err)
			}

			if err := layout.CreateMainKustomizationFile(appPath); err != nil {
				log.Fatal(err)
			}
		}

		log.Info("Successfully generated kga app")
	},
}
