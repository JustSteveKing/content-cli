package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/juststeveking/content-cli/internal/content"
	"github.com/juststeveking/content-cli/internal/output"
	"github.com/juststeveking/content-cli/internal/slug"
	tmpl "github.com/juststeveking/content-cli/internal/template"
)

var newCollection string

var newCmd = &cobra.Command{
	Use:   "new <title>",
	Short: "Scaffold a new content file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		col, ok := cfg.Collections[newCollection]
		if !ok {
			return fmt.Errorf("unknown collection %q — check .content.yaml", newCollection)
		}

		title := args[0]
		s := slug.From(title, col.Slug)
		filename := fmt.Sprintf("%s.%s", s, col.Format)
		destPath := filepath.Join(col.Dir, filename)

		fm := content.Frontmatter{
			"title": title,
			"date":  time.Now().Format("2006-01-02"),
			"draft": col.Defaults.Draft,
		}
		if col.Defaults.Author != "" {
			fm["author"] = col.Defaults.Author
		}
		for _, field := range col.RequiredFields {
			if _, exists := fm[field]; !exists {
				fm[field] = ""
			}
		}

		body := tmpl.Read(col.Template)

		if err := os.MkdirAll(col.Dir, 0755); err != nil {
			return fmt.Errorf("creating dir: %w", err)
		}
		encoded := content.Encode(fm, body)
		if err := os.WriteFile(destPath, encoded, 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}

		output.Success(fmt.Sprintf("Created %s  (slug: %s)", destPath, s))
		return nil
	},
}

func init() {
	newCmd.Flags().StringVarP(&newCollection, "collection", "c", "", "Collection to create content in (required)")
	_ = newCmd.MarkFlagRequired("collection")
}
