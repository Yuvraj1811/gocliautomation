package gitops

import "fmt"

func CommitAndPush(dryRun bool, commitMsg string, files ...string) error {
	if dryRun {
		fmt.Println("DRY-RUN: git commit & push skipped")
		return nil
	}

	repo, err := OpenRepo(".")
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := GitAdd(repo, f); err != nil {
			return err
		}
	}

	if err := GitCommit(repo, commitMsg); err != nil {
		return err
	}

	if err := GitPush(repo); err != nil {
		return err
	}

	return nil
}
