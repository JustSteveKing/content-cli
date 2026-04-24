package cmd

import (
	"os"
	"os/exec"

	"github.com/juststeveking/content-cli/internal/content"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <slug>",
	Short: "Open a content file in $EDITOR",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := scanCollections("")
		if err != nil {
			return err
		}

		f, err := content.FindBySlug(files, args[0])
		if err != nil {
			return err
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = os.Getenv("VISUAL")
		}
		if editor == "" {
			editor = "vi"
		}

		c := exec.Command(editor, f.Path)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}
