package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/juststeveking/content-cli/internal/content"
)

var (
	imageRegex = regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	linkRegex  = regexp.MustCompile(`\[.*?\]\(([^)]+)\)`)
)

type Issue struct {
	File    string
	Field   string
	Message string
}

func RequiredFields(file *content.ContentFile, required []string) []Issue {
	var issues []Issue
	for _, k := range file.FM.MissingKeys(required) {
		issues = append(issues, Issue{
			File:    file.Path,
			Field:   k,
			Message: fmt.Sprintf("required field '%s' is missing or empty", k),
		})
	}
	return issues
}

func ImagesExist(file *content.ContentFile) []Issue {
	var issues []Issue
	matches := imageRegex.FindAllStringSubmatch(file.Body, -1)
	dir := filepath.Dir(file.Path)
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		imgPath := m[1]
		if filepath.IsAbs(imgPath) {
			continue
		}
		full := filepath.Join(dir, imgPath)
		if _, err := os.Stat(full); err != nil {
			issues = append(issues, Issue{
				File:    file.Path,
				Field:   "image",
				Message: fmt.Sprintf("image not found: %s", imgPath),
			})
		}
	}
	return issues
}

func InternalLinks(file *content.ContentFile, contentDir string) []Issue {
	var issues []Issue
	matches := linkRegex.FindAllStringSubmatch(file.Body, -1)
	dir := filepath.Dir(file.Path)
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		target := m[1]
		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			continue
		}
		full := filepath.Join(dir, target)
		if _, err := os.Stat(full); err != nil {
			full2 := filepath.Join(contentDir, target)
			if _, err2 := os.Stat(full2); err2 != nil {
				issues = append(issues, Issue{
					File:    file.Path,
					Field:   "link",
					Message: fmt.Sprintf("internal link target not found: %s", target),
				})
			}
		}
	}
	return issues
}
