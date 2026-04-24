package slug_test

import (
	"testing"

	"github.com/juststeveking/content-cli/internal/slug"
)

func TestFrom(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		style string
		want  string
	}{
		{"kebab basic", "Hello World", "kebab", "hello-world"},
		{"kebab punctuation", "Hello, World!", "kebab", "hello-world"},
		{"kebab extra spaces", "  Hello   World  ", "kebab", "hello-world"},
		{"kebab numbers", "Post 42", "kebab", "post-42"},
		{"kebab empty", "", "kebab", ""},
		{"snake basic", "Hello World", "snake", "hello_world"},
		{"snake punctuation", "Hello, World!", "snake", "hello_world"},
		{"raw preserves case", "Hello World", "raw", "Hello World"},
		{"raw trims only", "  Hello World  ", "raw", "Hello World"},
		{"unknown style defaults to kebab", "Hello World", "camel", "hello-world"},
		{"unicode stripped", "Héllo World", "kebab", "hllo-world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slug.From(tt.raw, tt.style)
			if got != tt.want {
				t.Errorf("From(%q, %q) = %q, want %q", tt.raw, tt.style, got, tt.want)
			}
		})
	}
}
