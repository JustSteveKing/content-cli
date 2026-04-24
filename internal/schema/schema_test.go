package schema_test

import (
	"testing"

	"github.com/juststeveking/content-cli/internal/schema"
)

func validConfig() map[string]any {
	return map[string]any{
		"collections": map[string]any{
			"blog": map[string]any{
				"dir":    "src/content/blog",
				"format": "mdx",
			},
		},
	}
}

func TestValidateConfig_valid(t *testing.T) {
	errs, err := schema.ValidateConfig(validConfig())
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid config, got: %v", errs)
	}
}

func TestValidateConfig_withAllFields(t *testing.T) {
	data := map[string]any{
		"default_collection": "blog",
		"collections": map[string]any{
			"blog": map[string]any{
				"dir":             "src/content/blog",
				"format":          "mdx",
				"template":        "templates/blog.mdx",
				"required_fields": []any{"title", "date", "draft"},
				"optional_fields": []any{"description", "tags"},
				"slug":            "kebab",
				"defaults": map[string]any{
					"author": "Steve",
					"draft":  true,
				},
			},
		},
	}
	errs, err := schema.ValidateConfig(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func TestValidateConfig_missingCollections(t *testing.T) {
	errs, err := schema.ValidateConfig(map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) == 0 {
		t.Error("expected errors for missing collections key")
	}
}

func TestValidateConfig_invalidFormat(t *testing.T) {
	data := validConfig()
	data["collections"].(map[string]any)["blog"].(map[string]any)["format"] = "xml"
	errs, err := schema.ValidateConfig(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) == 0 {
		t.Error("expected error for invalid format value 'xml'")
	}
}

func TestValidateConfig_invalidSlugStyle(t *testing.T) {
	data := validConfig()
	data["collections"].(map[string]any)["blog"].(map[string]any)["slug"] = "camelCase"
	errs, err := schema.ValidateConfig(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) == 0 {
		t.Error("expected error for invalid slug style 'camelCase'")
	}
}

func TestValidateConfig_unknownProperty(t *testing.T) {
	data := validConfig()
	data["collections"].(map[string]any)["blog"].(map[string]any)["unknown_field"] = "oops"
	errs, err := schema.ValidateConfig(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) == 0 {
		t.Error("expected error for unknown property")
	}
}

func TestValidateConfig_emptyCollections(t *testing.T) {
	data := map[string]any{
		"collections": map[string]any{},
	}
	errs, err := schema.ValidateConfig(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) == 0 {
		t.Error("expected error for empty collections map")
	}
}

func TestValidationError_String(t *testing.T) {
	tests := []struct {
		e    schema.ValidationError
		want string
	}{
		{schema.ValidationError{Location: "/format", Message: "invalid value"}, "/format: invalid value"},
		{schema.ValidationError{Location: "", Message: "root error"}, "root error"},
		{schema.ValidationError{Location: "/", Message: "root error"}, "root error"},
	}
	for _, tt := range tests {
		if got := tt.e.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}
