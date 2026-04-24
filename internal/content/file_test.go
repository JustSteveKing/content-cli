package content_test

import (
	"os"
	"testing"
	"time"

	"github.com/juststeveking/content-cli/internal/content"
)

type mockFileInfo struct {
	modTime time.Time
}

func (m mockFileInfo) Name() string      { return "" }
func (m mockFileInfo) Size() int64       { return 0 }
func (m mockFileInfo) Mode() os.FileMode { return 0 }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool       { return false }
func (m mockFileInfo) Sys() any          { return nil }

func TestSlug(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"src/content/blog/hello-world.mdx", "hello-world"},
		{"src/content/blog/getting-started.md", "getting-started"},
		{"post.mdx", "post"},
		{"deep/nested/path/my-file.mdx", "my-file"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f := &content.ContentFile{Path: tt.path}
			if got := f.Slug(); got != tt.want {
				t.Errorf("Slug() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWordCount(t *testing.T) {
	tests := []struct {
		body string
		want int
	}{
		{"one two three", 3},
		{"  spaces   everywhere  ", 2},
		{"", 0},
		{"single", 1},
		{"line one\nline two\nline three", 6},
	}
	for _, tt := range tests {
		t.Run(tt.body, func(t *testing.T) {
			f := &content.ContentFile{Body: tt.body}
			if got := f.WordCount(); got != tt.want {
				t.Errorf("WordCount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestReadingTime(t *testing.T) {
	tests := []struct {
		words int
		want  int
	}{
		{0, 1},
		{100, 1},
		{199, 1},
		{200, 1},
		{400, 2},
		{1000, 5},
	}
	for _, tt := range tests {
		body := makeWords(tt.words)
		f := &content.ContentFile{Body: body}
		if got := f.ReadingTime(); got != tt.want {
			t.Errorf("ReadingTime() for %d words = %d, want %d", tt.words, got, tt.want)
		}
	}
}

func TestIsStale(t *testing.T) {
	fresh := &content.ContentFile{Info: mockFileInfo{modTime: time.Now()}}
	stale := &content.ContentFile{Info: mockFileInfo{modTime: time.Now().Add(-31 * 24 * time.Hour)}}
	boundary := &content.ContentFile{Info: mockFileInfo{modTime: time.Now().Add(-29 * 24 * time.Hour)}}

	if fresh.IsStale() {
		t.Error("fresh file should not be stale")
	}
	if !stale.IsStale() {
		t.Error("31-day-old file should be stale")
	}
	if boundary.IsStale() {
		t.Error("29-day-old file should not be stale")
	}
}

func makeWords(n int) string {
	if n == 0 {
		return ""
	}
	word := "word "
	result := make([]byte, 0, n*5)
	for i := 0; i < n; i++ {
		result = append(result, []byte(word)...)
	}
	return string(result)
}
