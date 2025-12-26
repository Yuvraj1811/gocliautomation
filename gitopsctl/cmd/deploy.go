package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gitopsctl/internal/gitops"
)

var (
	deployEnv   string
	deployImage string
	deployDry   bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy <service>",
	Short: "Deploy service by updating image in GitOps repo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]

		// 1️⃣ validation
		if err := validateDeployFlags(service); err != nil {
			return err
		}

		// 2️⃣ locate kustomization.yaml
		kustPath := filepath.Join(
			"apps", service, "overlays", deployEnv, "kustomization.yaml",
		)

		if _, err := os.Stat(kustPath); err != nil {
			return fmt.Errorf("kustomization not found: %s", kustPath)
		}

		// 3️⃣ update image
		_, err := gitops.UpdateImageTagInFile(
			kustPath,
			deployImage,
			deployDry,
		)
		if err != nil {
			return err
		}

		// 4️⃣ git (COMMON)
		commitMsg := fmt.Sprintf(
			"deploy(%s): image=%s env=%s",
			service,
			deployImage,
			deployEnv,
		)

		if err := gitops.CommitAndPush(
			deployDry,
			commitMsg,
			kustPath,
		); err != nil {
			return err
		}

		fmt.Println("Deploy completed successfully")
		return nil
	},
}

func init() {
	deployCmd.Flags().StringVar(&deployEnv, "env", "", "Environment (dev|staging|prod)")
	deployCmd.Flags().StringVar(&deployImage, "image", "", "Container image (name:tag)")
	deployCmd.Flags().BoolVar(&deployDry, "dry-run", false, "Show changes without committing")

	deployCmd.MarkFlagRequired("env")
	deployCmd.MarkFlagRequired("image")

	rootCmd.AddCommand(deployCmd)
}

func validateDeployFlags(service string) error {
	if service == "" {
		return errors.New("service name is required")
	}
	if deployEnv != "dev" && deployEnv != "staging" && deployEnv != "prod" {
		return fmt.Errorf("invalid env: %s", deployEnv)
	}
	if deployImage == "" {
		return errors.New("--image is required")
	}
	return nil
}
