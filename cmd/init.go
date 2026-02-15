package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"plug/internal/config"
	"plug/internal/gh"
	"plug/internal/scaffold"
	"plug/internal/shell"
	"strings"

	"github.com/spf13/cobra"
)

var (
	starterName string
	useGithub   bool
	isPublic    bool
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
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		// 2. Validate Starter
		starter, exists := cfg.Starters[starterName]
		if !exists {
			fmt.Printf("Starter '%s' not defined in config.toml\n", starterName)
			return
		}

		// 3. Get Template Directory
		tmplDir, err := config.GetTemplateDir()
		if err != nil {
			fmt.Printf("Error finding template dir: %v\n", err)
			return
		}

		// --- PHASE 1: SCAFFOLDING ---
		fmt.Printf("Scaffolding '%s' using starter '%s'...\n", pluginName, starterName)
		if err := scaffold.Generate(targetDir, starter.Items, starter.Scripts, tmplDir); err != nil {
			fmt.Printf("❌ Error generating files: %v\n", err)
			return
		}

		if len(starter.Scripts) > 0 {
			fmt.Println("Running setup scripts...")
			for _, script := range starter.Scripts {
				// Split "npm install" into ["npm", "install"]
				parts := strings.Fields(script)
				if len(parts) == 0 {
					continue
				}

				cmdName := parts[0]
				cmdArgs := parts[1:]

				fmt.Printf("   > Running: %s %s\n", cmdName, strings.Join(cmdArgs, " "))

				if err := shell.Run(targetDir, cmdName, cmdArgs...); err != nil {
					// We warn but don't abort, in case it's just a minor script failure
					fmt.Printf("Script failed: %s (Error: %v)\n", script, err)
				}
			}
		}

		// --- PHASE 2: LOCAL GIT (THIS WAS MISSING) ---
		fmt.Println("g  Initializing Git...")

		// A. Git Init
		if err := shell.Run(targetDir, "git", "init"); err != nil {
			fmt.Printf("Git init failed: %v\n", err)
			// Clean up if git fails
			gh.Cleanup(targetDir)
			return
		}

		// B. Git Add
		// We add "." to stage all files we just scaffolded
		if err := shell.Run(targetDir, "git", "add", "."); err != nil {
			fmt.Printf("Git add failed: %v\n", err)
			gh.Cleanup(targetDir)
			return
		}

		// C. Git Commit
		// gh requires a commit to exist before it can push
		if err := shell.Run(targetDir, "git", "commit", "-m", "Initial commit via Plug CLI"); err != nil {
			fmt.Printf("Git commit failed: %v\n", err)
			gh.Cleanup(targetDir)
			return
		}

		// --- PHASE 3: GITHUB REPO ---
		if useGithub {
			visibilityLabel := "private"
			if isPublic {
				visibilityLabel = "PUBLIC"
			}

			fmt.Printf("Creating %s repository in %s...\n", visibilityLabel, cfg.GithubOrg)

			if err := gh.CreateRepo(targetDir, cfg.GithubOrg, pluginName, isPublic); err != nil {
				fmt.Printf("GitHub creation failed: %v\n", err)
			}
		} else {
			if isPublic {
				fmt.Println("Warning: --public ignored because --gh was not set.")
			}
			fmt.Println("Skipping GitHub repository creation.")
		}
		fmt.Println("\nDone! Plugin created successfully.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&starterName, "starter", "s", "default", "The starter template config to use")
	initCmd.Flags().BoolVar(&useGithub, "gh", false, "Create a private GitHub repository")
}
