package gitops

import (
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func OpenRepo(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

func GitAdd(repo *git.Repository, filePath string) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(filePath)
	return err
}

func GitCommit(repo *git.Repository, message string) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name: "Yuvraj1811",
			Email: "ys2029472@gmail.com",
			When: time.Now(),

		},
	})
	return err
}

func GitPush(repo *git.Repository) error {
	err := repo.Push(&git.PushOptions{})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}