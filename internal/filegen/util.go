package filegen

import (
	"path/filepath"
	"strings"

	"github.com/integralist/gist-archiver/internal/model"
)

// isMarkdownFile checks if a Gist file should be treated as Markdown.
func isMarkdownFile(file model.GistFile) bool {
	// Prioritize the Language field from the API
	if strings.ToLower(file.Language) == "markdown" {
		return true
	}
	// Fallback to checking file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return ext == ".md" || ext == ".markdown"
}
