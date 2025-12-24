package cmd 

import (
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use: "gitopsctl",
	Short: "GitOps CLI for managing Kubernetes via Git",
	Long: "gitopsctl update GitOps repositories. It never touches cluster directly.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)

	}
}