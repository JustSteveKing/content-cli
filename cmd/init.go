package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/juststeveking/content-cli/internal/config"
	"github.com/juststeveking/content-cli/internal/output"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a .content.yaml config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.Exists() {
			output.Warn(".content.yaml already exists — delete it first to re-initialise")
			return nil
		}

		scanner := bufio.NewScanner(os.Stdin)
		prompt := func(question, defaultVal string) string {
			if defaultVal != "" {
				fmt.Printf("%s [%s]: ", question, defaultVal)
			} else {
				fmt.Printf("%s: ", question)
			}
			scanner.Scan()
			text := strings.TrimSpace(scanner.Text())
			if text == "" {
				return defaultVal
			}
			return text
		}

		collections := map[string]config.CollectionConfig{}
		var firstName string

		for {
			fmt.Println()
			name := prompt("Collection name (blank to finish)", "")
			if name == "" {
				if len(collections) == 0 {
					fmt.Println("  At least one collection is required.")
					continue
				}
				break
			}
			if firstName == "" {
				firstName = name
			}

			dir := prompt("  Directory", "src/content/"+name)
			format := prompt("  Format (md/mdx)", "mdx")
			requiredRaw := prompt("  Required fields (comma-separated)", "title,date,draft")
			optionalRaw := prompt("  Optional fields (comma-separated)", "description,tags")
			slugStyle := prompt("  Slug style (kebab/snake/raw)", "kebab")
			author := prompt("  Default author", "")
			draftDefault := prompt("  Draft by default? (true/false)", "true")
			template := prompt("  Template file path (blank to skip)", "")

			collections[name] = config.CollectionConfig{
				Dir:            dir,
				Format:         format,
				Template:       template,
				RequiredFields: splitTrim(requiredRaw),
				OptionalFields: splitTrim(optionalRaw),
				Slug:           slugStyle,
				Defaults: config.DefaultsConfig{
					Author: author,
					Draft:  draftDefault == "true",
				},
			}
		}

		defaultCollection := firstName
		if len(collections) > 1 {
			defaultCollection = prompt("Default collection", firstName)
		}

		cfg := &config.Config{
			DefaultCollection: defaultCollection,
			Collections:       collections,
		}

		if err := config.Write(cfg); err != nil {
			return fmt.Errorf("writing config: %w", err)
		}

		output.Success("Created .content.yaml")
		return nil
	},
}

func splitTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
