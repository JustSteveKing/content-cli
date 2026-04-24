package content_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/juststeveking/content-cli/internal/content"
)

func writeTestFile(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestScan(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "post-one.mdx"), "---\ntitle: Post One\ndraft: false\n---\nBody here.")
	writeTestFile(t, filepath.Join(dir, "post-two.mdx"), "---\ntitle: Post Two\ndraft: true\n---\n")
	writeTestFile(t, filepath.Join(dir, "ignored.md"), "---\ntitle: Wrong ext\n---\n")

	files, err := content.Scan(dir, "mdx")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestScanNoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "plain.mdx"), "Just plain text, no frontmatter.")

	files, err := content.Scan(dir, "mdx")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0].Body != "Just plain text, no frontmatter." {
		t.Errorf("unexpected body: %q", files[0].Body)
	}
}

func TestFindBySlug(t *testing.T) {
	files := []*content.ContentFile{
		{Path: "src/content/blog/hello-world.mdx"},
		{Path: "src/content/blog/hello-go.mdx"},
		{Path: "src/content/blog/getting-started.mdx"},
	}

	t.Run("exact match", func(t *testing.T) {
		f, err := content.FindBySlug(files, "getting-started")
		if err != nil {
			t.Fatal(err)
		}
		if f.Slug() != "getting-started" {
			t.Errorf("got slug %q, want %q", f.Slug(), "getting-started")
		}
	})

	t.Run("prefix match — single result", func(t *testing.T) {
		f, err := content.FindBySlug(files, "getting")
		if err != nil {
			t.Fatal(err)
		}
		if f.Slug() != "getting-started" {
			t.Errorf("got slug %q, want %q", f.Slug(), "getting-started")
		}
	})

	t.Run("ambiguous prefix returns error", func(t *testing.T) {
		_, err := content.FindBySlug(files, "hello")
		if err == nil {
			t.Error("expected error for ambiguous slug")
		}
	})

	t.Run("no match returns error", func(t *testing.T) {
		_, err := content.FindBySlug(files, "nonexistent")
		if err == nil {
			t.Error("expected error for unknown slug")
		}
	})

	t.Run("empty list returns error", func(t *testing.T) {
		_, err := content.FindBySlug([]*content.ContentFile{}, "anything")
		if err == nil {
			t.Error("expected error on empty file list")
		}
	})
}

func TestSortBy(t *testing.T) {
	files := []*content.ContentFile{
		{Path: "c.mdx", FM: content.Frontmatter{"title": "Zebra", "date": "2024-03-01"}, Body: "word"},
		{Path: "a.mdx", FM: content.Frontmatter{"title": "Apple", "date": "2024-01-01"}, Body: "word word word"},
		{Path: "b.mdx", FM: content.Frontmatter{"title": "Mango", "date": "2024-02-01"}, Body: "word word"},
	}

	t.Run("sort by title", func(t *testing.T) {
		cp := copyFiles(files)
		content.SortBy(cp, "title")
		titles := []string{cp[0].FM.GetString("title"), cp[1].FM.GetString("title"), cp[2].FM.GetString("title")}
		want := []string{"Apple", "Mango", "Zebra"}
		for i, tt := range want {
			if titles[i] != tt {
				t.Errorf("position %d: got %q, want %q", i, titles[i], tt)
			}
		}
	})

	t.Run("sort by words descending", func(t *testing.T) {
		cp := copyFiles(files)
		content.SortBy(cp, "words")
		if cp[0].WordCount() < cp[1].WordCount() {
			t.Error("expected descending word count order")
		}
	})

	t.Run("sort by date descending", func(t *testing.T) {
		cp := copyFiles(files)
		content.SortBy(cp, "date")
		if cp[0].FM.GetString("date") < cp[1].FM.GetString("date") {
			t.Error("expected descending date order")
		}
	})
}

func copyFiles(files []*content.ContentFile) []*content.ContentFile {
	out := make([]*content.ContentFile, len(files))
	for i, f := range files {
		cp := *f
		out[i] = &cp
	}
	return out
}
