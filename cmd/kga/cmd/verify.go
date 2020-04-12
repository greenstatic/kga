package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var Verify = &cobra.Command{
	Use:   "verify <app path or kga.yaml path>",
	Short: "Verify kga.yaml file or a kga app. file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		//appPath := args[0]
		//
		//appPathIsDir, err := isDir(appPath)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//
		//if appPathIsDir {
		//	appPath = filepath.Join(appPath, "kga.yaml")
		//}
		//
		//if err := config.VerifyKgaFile(appPath); err != nil {
		//	log.Error(err)
		//	log.Fatal("Bad configuration!")
		//} else {
		//	log.Info("Configuration is valid!")
		//}
	},
}

func isDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	mode := fi.Mode()
	return mode.IsDir(), nil
}
