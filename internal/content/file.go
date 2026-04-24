package content

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ContentFile struct {
	Path       string
	Collection string
	Info       os.FileInfo
	FM         Frontmatter
	Body       string
}

func (f *ContentFile) Slug() string {
	base := filepath.Base(f.Path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func (f *ContentFile) WordCount() int {
	return len(strings.Fields(f.Body))
}

func (f *ContentFile) ReadingTime() int {
	wc := f.WordCount() / 200
	if wc < 1 {
		return 1
	}
	return wc
}

func (f *ContentFile) IsStale() bool {
	return time.Since(f.Info.ModTime()) > 30*24*time.Hour
}
