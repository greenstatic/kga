package main

import (
	"github.com/greenstatic/kga/cmd/kga/cmd"
	"github.com/greenstatic/kga/pkg/log"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd.AddCommand(cmd.Init)
	rootCmd.AddCommand(cmd.Generate)
	//rootCmd.AddCommand(cmd.Verify)
	rootCmd.AddCommand(cmd.Version)
	execute()
}

var rootCmd = &cobra.Command{
	Use:   "kga",
	Short: "Manage your Kubernetes GitOps apps",
	Long:  `A CLI tool to manage your Kubernetes GitOps apps`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
