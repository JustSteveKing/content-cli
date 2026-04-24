package cmd

import (
	"fmt"
	"strconv"

	"github.com/juststeveking/content-cli/internal/content"
	"github.com/juststeveking/content-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	listDrafts     bool
	listPublished  bool
	listSort       string
	listCollection string
	listJSON       bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List content files",
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := scanCollections(listCollection)
		if err != nil {
			return err
		}

		var filtered []*content.ContentFile
		for _, f := range files {
			isDraft := f.FM.GetBool("draft")
			if listDrafts && !isDraft {
				continue
			}
			if listPublished && isDraft {
				continue
			}
			filtered = append(filtered, f)
		}

		content.SortBy(filtered, listSort)

		if listJSON {
			type item struct {
				Slug        string `json:"slug"`
				Title       string `json:"title"`
				Date        string `json:"date"`
				Draft       bool   `json:"draft"`
				Words       int    `json:"words"`
				ReadingTime int    `json:"reading_time"`
				Collection  string `json:"collection"`
			}
			result := make([]item, 0, len(filtered))
			for _, f := range filtered {
				result = append(result, item{
					Slug:        f.Slug(),
					Title:       f.FM.GetString("title"),
					Date:        f.FM.GetString("date"),
					Draft:       f.FM.GetBool("draft"),
					Words:       f.WordCount(),
					ReadingTime: f.ReadingTime(),
					Collection:  f.Collection,
				})
			}
			output.JSON(result)
			return nil
		}

		showCollection := listCollection == ""
		headers := []string{"Slug", "Title", "Date", "Draft", "Words", "Reading Time"}
		if showCollection {
			headers = append([]string{"Collection"}, headers...)
		}

		rows := make([][]string, 0, len(filtered))
		for _, f := range filtered {
			row := []string{
				f.Slug(),
				f.FM.GetString("title"),
				f.FM.GetString("date"),
				strconv.FormatBool(f.FM.GetBool("draft")),
				strconv.Itoa(f.WordCount()),
				fmt.Sprintf("%d min", f.ReadingTime()),
			}
			if showCollection {
				row = append([]string{f.Collection}, row...)
			}
			rows = append(rows, row)
		}

		output.Table(headers, rows)
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVar(&listDrafts, "drafts", false, "Show only drafts")
	listCmd.Flags().BoolVar(&listPublished, "published", false, "Show only published")
	listCmd.Flags().StringVar(&listSort, "sort", "date", "Sort by: date, title, words")
	listCmd.Flags().StringVarP(&listCollection, "collection", "c", "", "Filter to a specific collection")
	listCmd.Flags().BoolVarP(&listJSON, "json", "j", false, "Output as JSON")
}
