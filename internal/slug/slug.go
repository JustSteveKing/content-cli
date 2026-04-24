package slug

import (
	"regexp"
	"strings"
)

var nonAlphanumSpace = regexp.MustCompile(`[^a-z0-9 ]+`)

func From(raw, style string) string {
	switch style {
	case "snake":
		return toSnake(raw)
	case "raw":
		return strings.TrimSpace(raw)
	default:
		return toKebab(raw)
	}
}

func toKebab(raw string) string {
	s := strings.ToLower(raw)
	s = nonAlphanumSpace.ReplaceAllString(s, "")
	return strings.Join(strings.Fields(s), "-")
}

func toSnake(raw string) string {
	s := strings.ToLower(raw)
	s = nonAlphanumSpace.ReplaceAllString(s, "")
	return strings.Join(strings.Fields(s), "_")
}
