package content

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Frontmatter map[string]any

func Parse(data []byte) (Frontmatter, string, error) {
	s := string(data)
	if !strings.HasPrefix(s, "---") {
		return Frontmatter{}, s, nil
	}

	rest := s[3:]
	idx := strings.Index(rest, "---")
	if idx == -1 {
		return Frontmatter{}, s, nil
	}

	yamlBlock := rest[:idx]
	body := rest[idx+3:]
	body = strings.TrimPrefix(body, "\n")

	fm := Frontmatter{}
	if err := yaml.Unmarshal([]byte(yamlBlock), &fm); err != nil {
		return nil, "", fmt.Errorf("parsing frontmatter: %w", err)
	}
	return fm, body, nil
}

func Encode(fm Frontmatter, body string) []byte {
	var buf bytes.Buffer
	buf.WriteString("---\n")
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	_ = enc.Encode(fm)
	_ = enc.Close()
	buf.WriteString("---\n")
	if body != "" {
		buf.WriteString(body)
	}
	return buf.Bytes()
}

func (fm Frontmatter) GetString(key string) string {
	if v, ok := fm[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func (fm Frontmatter) GetBool(key string) bool {
	if v, ok := fm[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func (fm Frontmatter) MissingKeys(required []string) []string {
	var missing []string
	for _, k := range required {
		v, ok := fm[k]
		if !ok || v == nil || fmt.Sprintf("%v", v) == "" {
			missing = append(missing, k)
		}
	}
	return missing
}
