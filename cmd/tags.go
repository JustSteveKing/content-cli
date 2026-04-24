package cmd

import (
	"sort"
	"strconv"

	"github.com/juststeveking/content-cli/internal/content"
	"github.com/juststeveking/content-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	tagsCollection string
	tagsJSON       bool
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags and their frequency across content files",
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := scanCollections(tagsCollection)
		if err != nil {
			return err
		}

		freq := map[string]int{}
		for _, f := range files {
			for _, tag := range extractTags(f.FM) {
				freq[tag]++
			}
		}

		type entry struct {
			Tag   string `json:"tag"`
			Count int    `json:"count"`
		}

		entries := make([]entry, 0, len(freq))
		for tag, count := range freq {
			entries = append(entries, entry{tag, count})
		}
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Count != entries[j].Count {
				return entries[i].Count > entries[j].Count
			}
			return entries[i].Tag < entries[j].Tag
		})

		if tagsJSON {
			if len(entries) == 0 {
				output.JSON([]entry{})
				return nil
			}
			output.JSON(entries)
			return nil
		}

		if len(entries) == 0 {
			output.Info("no tags found across content files")
			return nil
		}

		rows := make([][]string, len(entries))
		for i, e := range entries {
			rows[i] = []string{e.Tag, strconv.Itoa(e.Count)}
		}
		output.Table([]string{"Tag", "Count"}, rows)
		return nil
	},
}

func init() {
	tagsCmd.Flags().StringVarP(&tagsCollection, "collection", "c", "", "Filter to a specific collection")
	tagsCmd.Flags().BoolVarP(&tagsJSON, "json", "j", false, "Output as JSON")
}

func extractTags(fm content.Frontmatter) []string {
	v, ok := fm["tags"]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case []any:
		tags := make([]string, 0, len(t))
		for _, item := range t {
			if s, ok := item.(string); ok && s != "" {
				tags = append(tags, s)
			}
		}
		return tags
	case string:
		if t != "" {
			return []string{t}
		}
	}
	return nil
}
