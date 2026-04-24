package validate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/juststeveking/content-cli/internal/content"
	"github.com/juststeveking/content-cli/internal/validate"
)

func TestRequiredFields(t *testing.T) {
	tests := []struct {
		name       string
		fm         content.Frontmatter
		required   []string
		wantIssues int
	}{
		{
			name:       "all fields present",
			fm:         content.Frontmatter{"title": "Hello", "date": "2024-01-01", "draft": false},
			required:   []string{"title", "date", "draft"},
			wantIssues: 0,
		},
		{
			name:       "one field missing",
			fm:         content.Frontmatter{"title": "Hello", "date": "2024-01-01"},
			required:   []string{"title", "date", "draft"},
			wantIssues: 1,
		},
		{
			name:       "all fields missing",
			fm:         content.Frontmatter{},
			required:   []string{"title", "date", "draft"},
			wantIssues: 3,
		},
		{
			name:       "empty string counts as missing",
			fm:         content.Frontmatter{"title": ""},
			required:   []string{"title"},
			wantIssues: 1,
		},
		{
			name:       "no required fields",
			fm:         content.Frontmatter{},
			required:   []string{},
			wantIssues: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &content.ContentFile{Path: "test.mdx", FM: tt.fm}
			issues := validate.RequiredFields(f, tt.required)
			if len(issues) != tt.wantIssues {
				t.Errorf("got %d issues, want %d: %v", len(issues), tt.wantIssues, issues)
			}
		})
	}
}

func TestImagesExist(t *testing.T) {
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "photo.png")
	if err := os.WriteFile(imgPath, []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("existing image — no issue", func(t *testing.T) {
		f := &content.ContentFile{
			Path: filepath.Join(dir, "post.mdx"),
			Body: "![alt](photo.png)",
		}
		issues := validate.ImagesExist(f)
		if len(issues) != 0 {
			t.Errorf("expected no issues, got %v", issues)
		}
	})

	t.Run("missing image — issue reported", func(t *testing.T) {
		f := &content.ContentFile{
			Path: filepath.Join(dir, "post.mdx"),
			Body: "![alt](missing.png)",
		}
		issues := validate.ImagesExist(f)
		if len(issues) != 1 {
			t.Errorf("expected 1 issue, got %d", len(issues))
		}
	})

	t.Run("no images — no issues", func(t *testing.T) {
		f := &content.ContentFile{
			Path: filepath.Join(dir, "post.mdx"),
			Body: "Just text, no images.",
		}
		issues := validate.ImagesExist(f)
		if len(issues) != 0 {
			t.Errorf("expected no issues, got %v", issues)
		}
	})
}

func TestInternalLinks(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "other-post.mdx")
	if err := os.WriteFile(target, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("existing internal link — no issue", func(t *testing.T) {
		f := &content.ContentFile{
			Path: filepath.Join(dir, "post.mdx"),
			Body: "[see also](other-post.mdx)",
		}
		issues := validate.InternalLinks(f, dir)
		if len(issues) != 0 {
			t.Errorf("expected no issues, got %v", issues)
		}
	})

	t.Run("missing internal link — issue reported", func(t *testing.T) {
		f := &content.ContentFile{
			Path: filepath.Join(dir, "post.mdx"),
			Body: "[gone](vanished.mdx)",
		}
		issues := validate.InternalLinks(f, dir)
		if len(issues) != 1 {
			t.Errorf("expected 1 issue, got %d", len(issues))
		}
	})

	t.Run("external links skipped", func(t *testing.T) {
		f := &content.ContentFile{
			Path: filepath.Join(dir, "post.mdx"),
			Body: "[go](https://go.dev) and [here](http://example.com)",
		}
		issues := validate.InternalLinks(f, dir)
		if len(issues) != 0 {
			t.Errorf("external links should not be checked, got %v", issues)
		}
	})
}
