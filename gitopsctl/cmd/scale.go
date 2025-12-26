package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitopsctl/internal/gitops"
)

var (
	scaleEnv   string
	scaleRep   int
	scaleDry   bool
)

var scaleCmd = &cobra.Command{
	Use:   "scale <service>",
	Short: "Scale replicas via GitOps",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]

		// 1️⃣ validation
		if scaleRep <= 0 {
			return fmt.Errorf("--replicas must be > 0")
		}
		if scaleEnv == "" {
			return fmt.Errorf("--env is required")
		}

		// 2️⃣ update replicas
		patchPath, err := gitops.UpdateReplicas(
			service,
			scaleEnv,
			scaleRep,
			scaleDry,
		)
		if err != nil {
			return err
		}

		// 3️⃣ git (COMMON)
		commitMsg := fmt.Sprintf(
			"scale(%s): replicas=%d env=%s",
			service,
			scaleRep,
			scaleEnv,
		)

		if err := gitops.CommitAndPush(
			scaleDry,
			commitMsg,
			patchPath,
		); err != nil {
			return err
		}

		fmt.Println("Scale completed successfully")
		return nil
	},
}

func init() {
	scaleCmd.Flags().StringVar(&scaleEnv, "env", "", "Environment (dev|staging|prod)")
	scaleCmd.Flags().IntVar(&scaleRep, "replicas", 0, "Number of replicas")
	scaleCmd.Flags().BoolVar(&scaleDry, "dry-run", false, "Show changes without committing")

	scaleCmd.MarkFlagRequired("env")
	scaleCmd.MarkFlagRequired("replicas")

	rootCmd.AddCommand(scaleCmd)
}
