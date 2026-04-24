package content

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func Scan(dir, ext string) ([]*ContentFile, error) {
	var files []*ContentFile

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, "."+ext) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			fmt.Printf(" WARN could not stat %s: %v\n", path, err)
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf(" WARN could not read %s: %v\n", path, err)
			return nil
		}

		fm, body, err := Parse(data)
		if err != nil {
			fmt.Printf(" WARN could not parse %s: %v\n", path, err)
			fm = Frontmatter{}
			body = string(data)
		}

		files = append(files, &ContentFile{
			Path: path,
			Info: info,
			FM:   fm,
			Body: body,
		})
		return nil
	})

	return files, err
}

func SortBy(files []*ContentFile, field string) {
	sort.Slice(files, func(i, j int) bool {
		switch field {
		case "title":
			return files[i].FM.GetString("title") < files[j].FM.GetString("title")
		case "words":
			return files[i].WordCount() > files[j].WordCount()
		default:
			di := parseDate(files[i].FM.GetString("date"))
			dj := parseDate(files[j].FM.GetString("date"))
			return di.After(dj)
		}
	})
}

func FindBySlug(files []*ContentFile, slug string) (*ContentFile, error) {
	for _, f := range files {
		if f.Slug() == slug {
			return f, nil
		}
	}
	var matches []*ContentFile
	for _, f := range files {
		if strings.HasPrefix(f.Slug(), slug) {
			matches = append(matches, f)
		}
	}
	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return nil, fmt.Errorf("no content file matching slug %q", slug)
	default:
		paths := make([]string, len(matches))
		for i, m := range matches {
			paths[i] = m.Path
		}
		return nil, fmt.Errorf("ambiguous slug %q — matches: %s", slug, strings.Join(paths, ", "))
	}
}

func parseDate(s string) time.Time {
	formats := []string{"2006-01-02", time.RFC3339}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
