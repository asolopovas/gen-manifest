package genWebmanifest

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

func Run() {
	if err := newRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

func newRootCmd() *cobra.Command {
	var (
		showVersion  bool
		icon         string
		iconsDir     string
		configPath   string
		manifestName string
		prefix       string
	)

	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	dirPath := filepath.Dir(currentFilePath)
	versionFilePath := filepath.Join(dirPath, "/../version")

	ver, err := os.ReadFile(versionFilePath)
	ErrChk(err)

	var rootCmd = &cobra.Command{
		Use:   "gen-manifest",
		Short: "Tool that help generate and resize web manifest for progressive web apps: " + string(ver),
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Println(string(ver))
				return
			}

			if configPath == "webmanifest.config.json" {
				configPath, err = filepath.Abs(configPath)
				if err != nil {
					log.Fatalf("Failed to get absolute path: %v", err)
				}
			}

			if !PathExist(configPath) {
				fmt.Println("No config found, generating `webmanifest.config.json` in current directory ")
				GenConfig(configPath)
				return
			}

			conf, err := GetJsonConfig(configPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config file '%s': %v\n", configPath, err)
				cmd.Help()
				return
			}

			if icon == "" {
				fmt.Println("Please povide icon path argument")
				return
			}

			GenWebmanifest(conf, icon, iconsDir, manifestName)
			fmt.Println("Webmanifest generated successfully")

		},
	}

	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Get Version")
	rootCmd.Flags().StringVarP(&icon, "icon", "i", "", "Path to app icon image")
	rootCmd.Flags().StringVarP(&iconsDir, "icons-dir", "d", "app-icons", "Resized icons destination directory")
	rootCmd.Flags().StringVarP(&manifestName, "manifest-name", "m", "manifest.webmanifest", "Name for webmanifest file")
	rootCmd.Flags().StringVarP(&prefix, "prefix", "p", "manifest.webmanifest", "Name for webmanifest file")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "webmanifest.config.json", "Custom config path")

	rootCmd.AddCommand(newCompletionCmd())

	return rootCmd
}

func newCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "Generate fish completion script",
		Run:   generateFishCompletion,
	}
}

func generateFishCompletion(cmd *cobra.Command, args []string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get user home directory: %v", err)
	}

	fishCompletionDir := filepath.Join(homeDir, ".config", "fish", "completions")
	if err := os.MkdirAll(fishCompletionDir, os.ModePerm); err != nil {
		log.Fatalf("failed to create fish completions directory: %v", err)
	}

	fishCompletionFile := filepath.Join(fishCompletionDir, "gen-webmanifest.fish")
	f, err := os.Create(fishCompletionFile)
	if err != nil {
		log.Fatalf("failed to create fish completion file: %v", err)
	}
	defer f.Close()

	if err := cmd.Root().GenFishCompletion(f, true); err != nil {
		log.Fatalf("failed to generate fish completion script: %v", err)
	}

	fmt.Printf("Fish completion script generated at: %s\n", fishCompletionFile)
}
