package gh

import (
	"fmt"
	"os"
	"plug/internal/shell"
)

func CreateRepo(dir, org, repoName string, isPublic bool) error {
	fullName := fmt.Sprintf("%s/%s", org, repoName)

	// Determine visibility flag
	visibility := "--private"
	if isPublic {
		visibility = "--public"
	}

	// Pass the visibility flag to the gh command
	return shell.Run(dir, "gh", "repo", "create", fullName, visibility, "--source=.", "--push")
}

// Cleanup removes the local directory if things fail
func Cleanup(path string) {
	fmt.Printf("Cleaning up: removing directory %s...\n", path)
	_ = os.RemoveAll(path)
}
