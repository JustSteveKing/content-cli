package cmd

import (
	"fmt"
	"os"

	"github.com/juststeveking/content-cli/internal/output"
	"github.com/juststeveking/content-cli/internal/schema"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Validate .content.yaml against the config schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(".content.yaml")
		if err != nil {
			return fmt.Errorf("reading .content.yaml: %w", err)
		}

		var raw map[string]any
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return fmt.Errorf("parsing .content.yaml: %w", err)
		}

		errs, err := schema.ValidateConfig(raw)
		if err != nil {
			return fmt.Errorf("running schema validation: %w", err)
		}

		if len(errs) == 0 {
			output.Success(".content.yaml is valid")
			return nil
		}

		for _, e := range errs {
			output.Warn(e.String())
		}
		fmt.Printf("\nFound %d schema violation(s)\n", len(errs))
		os.Exit(1)
		return nil
	},
}
