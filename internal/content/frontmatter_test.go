package content_test

import (
	"bytes"
	"testing"

	"github.com/juststeveking/content-cli/internal/content"
)

func TestParse(t *testing.T) {
	t.Run("valid frontmatter", func(t *testing.T) {
		input := []byte("---\ntitle: Hello World\ndraft: true\n---\nBody text here.")
		fm, body, err := content.Parse(input)
		if err != nil {
			t.Fatal(err)
		}
		if fm.GetString("title") != "Hello World" {
			t.Errorf("title = %q, want %q", fm.GetString("title"), "Hello World")
		}
		if !fm.GetBool("draft") {
			t.Error("expected draft = true")
		}
		if body != "Body text here." {
			t.Errorf("body = %q, want %q", body, "Body text here.")
		}
	})

	t.Run("no frontmatter returns full content as body", func(t *testing.T) {
		input := []byte("Just a plain body.")
		fm, body, err := content.Parse(input)
		if err != nil {
			t.Fatal(err)
		}
		if len(fm) != 0 {
			t.Errorf("expected empty frontmatter, got %v", fm)
		}
		if body != "Just a plain body." {
			t.Errorf("body = %q, want %q", body, "Just a plain body.")
		}
	})

	t.Run("unclosed delimiter returns full content as body", func(t *testing.T) {
		input := []byte("---\ntitle: Oops\n")
		fm, body, _ := content.Parse(input)
		if len(fm) != 0 {
			t.Errorf("expected empty frontmatter for unclosed delimiter")
		}
		_ = body
	})

	t.Run("empty input", func(t *testing.T) {
		fm, body, err := content.Parse([]byte{})
		if err != nil {
			t.Fatal(err)
		}
		if len(fm) != 0 {
			t.Errorf("expected empty frontmatter")
		}
		if body != "" {
			t.Errorf("expected empty body, got %q", body)
		}
	})
}

func TestEncode(t *testing.T) {
	fm := content.Frontmatter{
		"title": "Hello World",
		"draft": true,
	}
	body := "Some body content.\n"
	out := content.Encode(fm, body)

	if !bytes.Contains(out, []byte("---")) {
		t.Error("encoded output missing --- delimiter")
	}
	if !bytes.Contains(out, []byte("title: Hello World")) {
		t.Error("encoded output missing title field")
	}
	if !bytes.Contains(out, []byte("Some body content.")) {
		t.Error("encoded output missing body")
	}
}

func TestParseEncodeRoundtrip(t *testing.T) {
	original := []byte("---\ndate: \"2024-01-15\"\ndraft: false\ntitle: Roundtrip\n---\nBody.\n")
	fm, body, err := content.Parse(original)
	if err != nil {
		t.Fatal(err)
	}
	reencoded := content.Encode(fm, body)
	fm2, body2, err := content.Parse(reencoded)
	if err != nil {
		t.Fatalf("re-parse failed: %v", err)
	}
	if fm2.GetString("title") != "Roundtrip" {
		t.Errorf("title lost in roundtrip: %q", fm2.GetString("title"))
	}
	if body2 != body {
		t.Errorf("body changed in roundtrip: got %q, want %q", body2, body)
	}
}

func TestMissingKeys(t *testing.T) {
	fm := content.Frontmatter{
		"title": "Hello",
		"date":  "2024-01-01",
		"empty": "",
	}

	tests := []struct {
		name     string
		required []string
		want     []string
	}{
		{"all present", []string{"title", "date"}, nil},
		{"one missing", []string{"title", "date", "draft"}, []string{"draft"}},
		{"empty value counts as missing", []string{"empty"}, []string{"empty"}},
		{"no required fields", []string{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fm.MissingKeys(tt.required)
			if len(got) != len(tt.want) {
				t.Errorf("MissingKeys(%v) = %v, want %v", tt.required, got, tt.want)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	fm := content.Frontmatter{
		"yes":    true,
		"no":     false,
		"string": "true",
	}
	if !fm.GetBool("yes") {
		t.Error("expected true")
	}
	if fm.GetBool("no") {
		t.Error("expected false")
	}
	if fm.GetBool("missing") {
		t.Error("missing key should return false")
	}
	if fm.GetBool("string") {
		t.Error("string 'true' should not coerce to bool true")
	}
}
