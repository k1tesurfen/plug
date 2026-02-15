package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Initialize configuration and check dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setting up Plug...")

		// 1. Check Prerequisites (gh cli)
		path, err := exec.LookPath("gh")
		if err != nil {
			fmt.Println("❌ Error: 'gh' CLI not found. Please install GitHub CLI first.")
			fmt.Println("   Mac: brew install gh")
			fmt.Println("   Win: winget install GitHub.cli")
			return // Don't exit(1), just return so user sees the message
		}
		fmt.Printf("Found 'gh' CLI at: %s\n", path)

		// 2. Determine Config Directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Could not determine user home directory: %v\n", err)
			return
		}

		plugDir := filepath.Join(home, ".config", "plug")
		templatesDir := filepath.Join(plugDir, "templates")

		// 3. Create Directories
		if err := os.MkdirAll(templatesDir, 0755); err != nil {
			fmt.Printf("❌ Failed to create directories: %v\n", err)
			return
		}
		fmt.Printf("✅ Verified directory: %s\n", plugDir)
		fmt.Printf("✅ Verified directory: %s\n", templatesDir)

		// 4. Create Default Config if missing
		configFile := filepath.Join(plugDir, "config.toml")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			defaultConfig := `# Default configuration for Plug CLI
github_org = "YourGithubOrgHere"

[starters]
  [starters.default]
  # Files or Folders inside your templates/ directory
  items = [] 
`
			if err := os.WriteFile(configFile, []byte(defaultConfig), 0644); err != nil {
				fmt.Printf("❌ Failed to create config file: %v\n", err)
				return
			}
			fmt.Printf("Created default config: %s\n", configFile)
			fmt.Println("ACTION REQUIRED: Edit config.toml and set your 'github_org'")
		} else {
			fmt.Printf("Config file already exists: %s\n", configFile)
		}

		fmt.Println("\nSetup complete! You can now populate your 'templates' folder.")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
