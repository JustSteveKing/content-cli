package cmd

import (
	"fmt"
	"os"

	"github.com/juststeveking/content-cli/internal/content"
	"github.com/juststeveking/content-cli/internal/output"
	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish <slug>",
	Short: "Publish a draft post by setting draft: false",
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

		if !f.FM.GetBool("draft") {
			output.Info(fmt.Sprintf("%s is already published", f.Path))
			return nil
		}

		f.FM["draft"] = false
		encoded := content.Encode(f.FM, f.Body)
		if err := os.WriteFile(f.Path, encoded, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", f.Path, err)
		}

		output.Success(fmt.Sprintf("Published %s", f.Path))
		return nil
	},
}
