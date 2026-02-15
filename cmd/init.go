package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"plug/internal/config"
	"plug/internal/gh"
	"plug/internal/scaffold"
	"plug/internal/shell" // <--- Make sure this import is here

	"github.com/spf13/cobra"
)

var (
	starterName string
	useGithub   bool
)

var initCmd = &cobra.Command{
	Use:   "init [plugin-name]",
	Short: "Scaffold a new WordPress plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]

		// 0. Prepare Paths
		cwd, _ := os.Getwd()
		targetDir := filepath.Join(cwd, pluginName)

		// 1. Load Config
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("❌ Error loading config: %v\n", err)
			return
		}

		// 2. Validate Starter
		starter, exists := cfg.Starters[starterName]
		if !exists {
			fmt.Printf("❌ Starter '%s' not defined in config.toml\n", starterName)
			return
		}

		// 3. Get Template Directory
		tmplDir, err := config.GetTemplateDir()
		if err != nil {
			fmt.Printf("❌ Error finding template dir: %v\n", err)
			return
		}

		// --- PHASE 1: SCAFFOLDING ---
		fmt.Printf("🚀  Scaffolding '%s' using starter '%s'...\n", pluginName, starterName)
		if err := scaffold.Generate(targetDir, starter.Items, tmplDir); err != nil {
			fmt.Printf("❌ Error generating files: %v\n", err)
			return
		}

		// --- PHASE 2: LOCAL GIT (THIS WAS MISSING) ---
		fmt.Println("g  Initializing Git...")

		// A. Git Init
		if err := shell.Run(targetDir, "git", "init"); err != nil {
			fmt.Printf("❌ Git init failed: %v\n", err)
			// Clean up if git fails
			gh.Cleanup(targetDir)
			return
		}

		// B. Git Add
		// We add "." to stage all files we just scaffolded
		if err := shell.Run(targetDir, "git", "add", "."); err != nil {
			fmt.Printf("❌ Git add failed: %v\n", err)
			gh.Cleanup(targetDir)
			return
		}

		// C. Git Commit
		// gh requires a commit to exist before it can push
		if err := shell.Run(targetDir, "git", "commit", "-m", "Initial commit via Plug CLI"); err != nil {
			fmt.Printf("❌ Git commit failed: %v\n", err)
			gh.Cleanup(targetDir)
			return
		}

		// --- PHASE 3: GITHUB REPO ---
		if useGithub {
			fmt.Printf("Creating private repository in %s...\n", cfg.GithubOrg)

			if err := gh.CreateRepo(targetDir, cfg.GithubOrg, pluginName); err != nil {
				fmt.Printf("❌ GitHub creation failed: %v\n", err)
			}
		} else {
			fmt.Println("Skipping GitHub repository creation (use --gh to enable)")
		}
		fmt.Println("\nDone! Plugin created successfully.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&starterName, "starter", "s", "default", "The starter template config to use")
	initCmd.Flags().BoolVar(&useGithub, "gh", false, "Create a private GitHub repository")
}
