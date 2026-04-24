package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/juststeveking/content-cli/internal/config"
	"github.com/spf13/cobra"
)

var Version = "dev"

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:     "content",
	Short:   "Manage Markdown/MDX content files",
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "init" || cmd.Name() == "config" {
			return nil
		}
		if !config.Exists() {
			return fmt.Errorf("no .content.yaml found — run 'content init' first")
		}
		loaded, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		cfg = loaded
		return nil
	},
}

func Execute() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(tagsCmd)
	rootCmd.AddCommand(serveCmd)
}
