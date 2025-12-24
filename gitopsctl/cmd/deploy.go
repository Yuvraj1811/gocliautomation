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
	env     string
	image   string
	dryRun  bool
	replicas int
)

var deployCmd = &cobra.Command{
	Use:   "deploy <app>",
	Short: "Deploy an application by updating GitOps repo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := args[0]

		// validation
		if err := validateDeployFlags(app); err != nil {
			return err
		}

		// Locate overlay kustomization.yaml
		path := filepath.Join("apps", app, "overlays", env, "kustomization.yaml")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", path)
		}

		// Update image in YAML (map-based)
		updatedData, err := gitops.UpdateImageTagInFile(path, image, dryRun)
		if err != nil {
			return err
		}

		// Dry-run output
		if dryRun {
			originalData, _ := os.ReadFile(path)
			fmt.Println("----- DRY RUN -----")
			fmt.Println("File:", path)
			fmt.Println("----- BEFORE -----")
			fmt.Println(string(originalData))
			fmt.Println("----- AFTER ------")
			fmt.Println(string(updatedData))
			return nil
		}

		fmt.Println("YAML updated successfully")

		// Git automation
		repo, err := gitops.OpenRepo(".")
		if err != nil {
			return err
		}

		if err := gitops.GitAdd(repo, path); err != nil {
			return err
		}

		commitMsg := fmt.Sprintf("deploy(%s): update image to %s in %s", app, image, env)
		if err := gitops.GitCommit(repo, commitMsg); err != nil {
			return err
		}

		if err := gitops.GitPush(repo); err != nil {
			return err
		}

		fmt.Println("Deployment committed to Git")
		return nil
	},
}

func init() {
	deployCmd.Flags().StringVar(&env, "env", "", "Environment (dev|staging|prod)")
	deployCmd.Flags().StringVar(&image, "image", "", "Container image (required)")
	deployCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show changes without committing")
	deployCmd.Flags().IntVar(&replicas, "replicas", 0, "Override replicas")
	rootCmd.AddCommand(deployCmd)
}

func validateDeployFlags(app string) error {
	if app == "" {
		return errors.New("app name is required")
	}
	if env == "" {
		return errors.New("--env is required")
	}
	if env != "dev" && env != "staging" && env != "prod" {
		return fmt.Errorf("invalid env: %s", env)
	}
	if image == "" {
		return errors.New("--image is required")
	}
	return nil
}
