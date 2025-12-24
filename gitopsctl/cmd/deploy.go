package cmd 

import (
	"errors"
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"gitopsctl/internal/gitops"
)

var (
	env string
	image string
	dryRun bool
	replicas int
)

var deployCmd = &cobra.Command{
	Use: "deploy <app>",
	Short: "Deploy an application by updating Gitops repo",
	Args: cobra.ExactArgs(1), // app name mandatory 
	RunE: func(cmd *cobra.Command, args []string) error {  //RunE
		app := args[0]

		//validation
		if err := validateDeployFlags(app); err != nil {
			return err
		}

		path := gitops.KustomizationPath(app, env)


		// read original file
		originalData, err := os.ReadFile(path)
		if err != nil {
			return err
		}


		k, err := gitops.LoadKustomization(path)
		if err != nil {
			return err
		}

		if err := gitops.UpdateImageTag(k, app, image); err != nil {
			return err
		}

		if err := gitops.SaveKustomization(path, k); err != nil {
			return err
		}

		updatedData, err := gitops.RenderYAML(k)
		if err != nil {
			return err
		}
				
		if dryRun {
			fmt.Println("----- DRY RUN -----")
			fmt.Println("File:", path)
			fmt.Println("----- BEFORE -----")
			fmt.Println(string(originalData))
			fmt.Println("----- AFTER ------")
			fmt.Println(string(updatedData))
			return nil
		}

		fmt.Println("YAML updated successfully")

		repo, err := gitops.OpenRepo(".")
		if err != nil {
			return err
		}

		if err := gitops.GitAdd(repo, path); err != nil {
			return err
		}

		commitMsg := fmt.Sprintf(
			"deploy(%s): update image to %s in %s",
			app, image, env,
		)

		if err := gitops.GitCommit(repo, commitMsg); err != nil {
	        return err
        }

		if !dryRun {
			if err := gitops.GitPush(repo); err != nil {
				return err
			}
		}


		fmt.Println("Deploying app:", app)
		fmt.Println("Env:", env)
		fmt.Println("Image:", image)
		fmt.Println("Dry-run:", dryRun)

		return nil

	},
}

func init(){
	deployCmd.Flags().StringVar(&env, "env", "", "Environment (dev|staging|prod)")
	deployCmd.Flags().StringVar(&image, "image", "", "Container image (required)")
	deployCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show changes without comitting")
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