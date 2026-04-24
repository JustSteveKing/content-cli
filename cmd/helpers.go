package cmd

import (
	"fmt"

	"github.com/juststeveking/content-cli/internal/content"
	"github.com/juststeveking/content-cli/internal/output"
)

// scanCollections scans one named collection, or all collections when name is empty.
// Each file's Collection field is set to the collection name.
func scanCollections(name string) ([]*content.ContentFile, error) {
	if name != "" {
		col, ok := cfg.Collections[name]
		if !ok {
			return nil, fmt.Errorf("unknown collection %q — check .content.yaml", name)
		}
		files, err := content.Scan(col.Dir, col.Format)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			f.Collection = name
		}
		return files, nil
	}

	var all []*content.ContentFile
	for cname, col := range cfg.Collections {
		files, err := content.Scan(col.Dir, col.Format)
		if err != nil {
			output.Warn(fmt.Sprintf("scanning collection %q: %v", cname, err))
			continue
		}
		for _, f := range files {
			f.Collection = cname
		}
		all = append(all, files...)
	}
	return all, nil
}
