package cmd

import (
	"fmt"
	"os"

	"github.com/juststeveking/content-cli/internal/output"
	"github.com/juststeveking/content-cli/internal/validate"
	"github.com/spf13/cobra"
)

var (
	checkCollection string
	checkJSON       bool
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate content files for issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := scanCollections(checkCollection)
		if err != nil {
			return err
		}

		var allIssues []validate.Issue
		for _, f := range files {
			col := cfg.Collections[f.Collection]
			allIssues = append(allIssues, validate.RequiredFields(f, col.RequiredFields)...)
			allIssues = append(allIssues, validate.ImagesExist(f)...)
			allIssues = append(allIssues, validate.InternalLinks(f, col.Dir)...)
		}

		if checkJSON {
			type result struct {
				OK      bool             `json:"ok"`
				Checked int              `json:"checked"`
				Issues  []validate.Issue `json:"issues"`
			}
			output.JSON(result{
				OK:      len(allIssues) == 0,
				Checked: len(files),
				Issues:  allIssues,
			})
			if len(allIssues) > 0 {
				os.Exit(1)
			}
			return nil
		}

		for _, issue := range allIssues {
			fmt.Printf(" WARN %s: %s: %s\n", issue.File, issue.Field, issue.Message)
		}

		if len(allIssues) == 0 {
			fmt.Printf("  OK  checked %d file(s), no issues found\n", len(files))
			return nil
		}

		fmt.Printf("\nFound %d issue(s) across %d file(s)\n", len(allIssues), len(files))
		os.Exit(1)
		return nil
	},
}

func init() {
	checkCmd.Flags().StringVarP(&checkCollection, "collection", "c", "", "Check a specific collection only")
	checkCmd.Flags().BoolVarP(&checkJSON, "json", "j", false, "Output as JSON")
}
