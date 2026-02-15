package gh

import (
	"fmt"
	"os"
	"plug/internal/shell"
)

func CreateRepo(dir, org, repoName string) error {
	fullName := fmt.Sprintf("%s/%s", org, repoName)

	// 1. Create the Repo on GitHub (without pushing yet)
	// We use --source=. so it configures the 'origin' remote for us automatically.
	err := shell.Run(dir, "gh", "repo", "create", fullName, "--private", "--source=.")
	if err != nil {
		return fmt.Errorf("failed to create repo on GitHub: %w", err)
	}

	// 2. Push manually
	// This helps us catch if the issue is purely a git-push error
	// -u origin HEAD pushes the current branch (main/master) to origin
	err = shell.Run(dir, "git", "push", "-u", "origin", "HEAD")
	if err != nil {
		return fmt.Errorf("failed to push initial commit: %w", err)
	}

	return nil
}

// Cleanup removes the local directory if things fail
func Cleanup(path string) {
	fmt.Printf("Cleaning up: removing directory %s...\n", path)
	_ = os.RemoveAll(path)
}
