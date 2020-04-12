package cmd

import (
	"github.com/greenstatic/kga/pkg/app"
	"github.com/greenstatic/kga/pkg/log"
	"github.com/spf13/cobra"
	"path/filepath"
)

var Init = &cobra.Command{
	Use:   "init <type> <directory>",
	Short: "Initialize a kga app",
	Long: `Initialize a kga app
<type> can be one of: basic, manifest or helm

basic: you write all the resources, kga will just create the app dir structure
manifest: we pull the manifests from the web and create the app dir structure
helm: we use helm to download the specified chart, use helm template to create 
      the base manifests and create the app dir structure`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appType := args[0]
		appName := filepath.Base(args[1])
		flagName := cmd.Flag("name").Value.String()
		if flagName != "" {
			appName = flagName
		}

		appPath := args[1]
		appPathAbsolute := app.PathAbsolute(appPath)

		log.FatalOnError(app.ValidateTypeString(appType))
		log.FatalOnError(app.ValidAppPathToInit(appPath))

		log.Infof("Initializing kga %s app named %s in: %s", appType, appName, appPathAbsolute)
		a := app.CreateType(appType)

		if err := app.Init(a, appPathAbsolute, appName); err != nil {
			log.Fatal(err)
		}

		log.Info("Successfully initialized kga app")
	},
}

func init() {
	Init.Flags().StringP("name", "n", "", "Name of app (default: last dir in app path)")
}
