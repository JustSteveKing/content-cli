package cmd

import (
	"fmt"

	"github.com/juststeveking/content-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	statusCollection string
	statusJSON       bool
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show content health summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := scanCollections(statusCollection)
		if err != nil {
			return err
		}

		type stats struct {
			Total      int `json:"total"`
			Published  int `json:"published"`
			Drafts     int `json:"drafts"`
			TotalWords int `json:"total_words"`
			AvgWords   int `json:"avg_words"`
			Stale      int `json:"stale"`
		}

		byCollection := map[string]*stats{}
		for _, f := range files {
			if _, ok := byCollection[f.Collection]; !ok {
				byCollection[f.Collection] = &stats{}
			}
			s := byCollection[f.Collection]
			s.Total++
			if f.FM.GetBool("draft") {
				s.Drafts++
			} else {
				s.Published++
			}
			s.TotalWords += f.WordCount()
			if f.IsStale() {
				s.Stale++
			}
		}

		for _, s := range byCollection {
			if s.Total > 0 {
				s.AvgWords = s.TotalWords / s.Total
			}
		}

		if statusJSON {
			type payload struct {
				Collections map[string]*stats `json:"collections"`
				Totals      *stats            `json:"totals,omitempty"`
			}
			p := payload{Collections: byCollection}
			if len(byCollection) > 1 {
				total := &stats{}
				for _, s := range byCollection {
					total.Total += s.Total
					total.Published += s.Published
					total.Drafts += s.Drafts
					total.TotalWords += s.TotalWords
					total.Stale += s.Stale
				}
				if total.Total > 0 {
					total.AvgWords = total.TotalWords / total.Total
				}
				p.Totals = total
			}
			output.JSON(p)
			return nil
		}

		printStats := func(label string, s *stats) {
			output.Info(fmt.Sprintf("%s", label))
			output.Info(fmt.Sprintf("  Total files:    %d", s.Total))
			output.Info(fmt.Sprintf("  Published:      %d", s.Published))
			output.Info(fmt.Sprintf("  Drafts:         %d", s.Drafts))
			output.Info(fmt.Sprintf("  Total words:    %d", s.TotalWords))
			output.Info(fmt.Sprintf("  Avg words/post: %d", s.AvgWords))
			output.Info(fmt.Sprintf("  Stale (>30d):   %d", s.Stale))
		}

		if statusCollection != "" {
			s := byCollection[statusCollection]
			if s == nil {
				s = &stats{}
			}
			printStats(statusCollection, s)
			return nil
		}

		total := &stats{}
		for name, s := range byCollection {
			printStats(name+":", s)
			fmt.Println()
			total.Total += s.Total
			total.Published += s.Published
			total.Drafts += s.Drafts
			total.TotalWords += s.TotalWords
			total.Stale += s.Stale
		}

		if len(byCollection) > 1 {
			if total.Total > 0 {
				total.AvgWords = total.TotalWords / total.Total
			}
			printStats("total:", total)
		}

		return nil
	},
}

func init() {
	statusCmd.Flags().StringVarP(&statusCollection, "collection", "c", "", "Show stats for a specific collection only")
	statusCmd.Flags().BoolVarP(&statusJSON, "json", "j", false, "Output as JSON")
}
