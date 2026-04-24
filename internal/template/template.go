package template

import "os"

func Read(path string) string {
	if path == "" {
		return "\n"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "\n"
	}
	return string(data)
}
