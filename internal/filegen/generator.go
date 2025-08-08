package filegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/integralist/gist-archiver/internal/config"
	"github.com/integralist/gist-archiver/internal/model"
)

// generateGistMarkdownFile creates the content for a single gist's markdown file.
func generateGistMarkdownFile(gist model.GistDetail) (string, error) {
	var mdPageContent strings.Builder
	title := gist.GetTitle()
	tagsStart := strings.Index(title, "#")
	if tagsStart >= 0 {
		title = title[:tagsStart]
	}
	mdPageContent.WriteString(fmt.Sprintf("# %s\n\n", title))

	if gist.HTMLURL != "" {
		mdPageContent.WriteString(fmt.Sprintf("[View original Gist on GitHub](%s)\n\n", gist.HTMLURL))
	}
	if gist.Tags != "" {
		mdPageContent.WriteString(fmt.Sprintf("**Tags:** %s\n\n", gist.Tags))
	}

	if len(gist.Files) == 0 {
		mdPageContent.WriteString("This gist has no files.\n\n")
	} else {
		var filenames []string
		for fname := range gist.Files {
			filenames = append(filenames, fname)
		}
		sort.Strings(filenames)

		for _, filename := range filenames {
			file := gist.Files[filename]
			isMD := isMarkdownFile(file)

			fileContent := file.Content
			if file.Truncated && fileContent == "" {
				if isMD {
					fileContent = fmt.Sprintf("\n\n> **Note:** Content of this file is truncated. Full content available at %s\n\n", file.RawURL)
				} else {
					fileContent = fmt.Sprintf("*Content of this file is truncated. Full content available at %s*", file.RawURL)
				}
			}

			mdPageContent.WriteString(fmt.Sprintf("## %s\n\n", filename))

			if isMD {
				mdPageContent.WriteString(fileContent)
				mdPageContent.WriteString("\n\n")
			} else {
				lang := file.Language
				if lang == "" {
					ext := strings.TrimPrefix(filepath.Ext(filename), ".")
					if ext != "" {
						lang = ext
					} else {
						lang = "text"
					}
				}
				mdPageContent.WriteString(fmt.Sprintf("```%s\n", strings.ToLower(lang)))
				mdPageContent.WriteString(fileContent)
				mdPageContent.WriteString("\n```\n\n")
			}
		}
	}

	mdFilenameForTempStore := gist.ID + ".md"
	filePath := filepath.Join(config.MarkdownOutputDir, config.GistsMdSubDir, mdFilenameForTempStore)

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", filepath.Dir(filePath), err)
	}
	if err := os.WriteFile(filePath, []byte(mdPageContent.String()), 0o644); err != nil {
		return "", fmt.Errorf("failed to write composite markdown file %s: %w", filePath, err)
	}
	return filepath.Join(config.GistsMdSubDir, mdFilenameForTempStore), nil
}

// ProcessSingleGist handles markdown/HTML generation for one gist and sends an entry for the index.
func ProcessSingleGist(gist model.GistDetail, indexEntryChan chan<- model.IndexEntry, searchIndexEntryChan chan<- model.SearchIndexEntry) {
	title := gist.GetTitle()
	tagsStart := strings.Index(title, "#")
	if tagsStart >= 0 {
		gist.Tags = title[tagsStart:]
		title = title[:tagsStart]
	}
	relativeMdPath, err := generateGistMarkdownFile(gist)
	if err != nil {
		log.Printf("Error generating composite markdown for gist %s (%s): %v. Skipping.", gist.ID, title, err)
		return
	}

	htmlFilenameInGistsDir := strings.Replace(filepath.Base(relativeMdPath), ".md", ".html", 1)
	relativeHTMLPath := filepath.Join(config.GistsHtmlSubDir, htmlFilenameInGistsDir)
	fullGistHTMLPath := filepath.Join(config.WebDeployDir, relativeHTMLPath)

	fullGistCompositeMdPath := filepath.Join(config.MarkdownOutputDir, relativeMdPath)
	mdInput, err := os.ReadFile(fullGistCompositeMdPath)
	if err != nil {
		log.Printf("Error reading composite gist markdown for HTML conversion %s: %v", fullGistCompositeMdPath, err)
		return
	}

	var gistHTMLBuffer bytes.Buffer
	gistHTMLBuffer.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n")
	gistHTMLBuffer.WriteString(fmt.Sprintf("<title>%s</title>\n", title))
	gistHTMLBuffer.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1.0">`)
	gistHTMLBuffer.WriteString(fmt.Sprintf("<link rel=\"stylesheet\" type=\"text/css\" href=\"../%s/%s/%s\">\n",
		config.AssetsSubDir, config.CssSubDir, config.CssFileName))
	gistHTMLBuffer.WriteString("</head>\n<body>\n<div class=\"container\">\n")
	gistHTMLBuffer.WriteString("<a href=\"../index.html\" class=\"back-link\">&laquo; Back to Index</a>\n")

	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.FencedCode)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.UseXHTML | html.Smartypants
	opts := html.RendererOptions{Flags: htmlFlags, Title: title}
	renderer := html.NewRenderer(opts)
	renderedMarkdown := md.Render(p.Parse(mdInput), renderer)

	gistHTMLBuffer.WriteString("<div class=\"markdown-content\">\n")
	gistHTMLBuffer.Write(renderedMarkdown)
	gistHTMLBuffer.WriteString("</div>\n")

	gistHTMLBuffer.WriteString(fmt.Sprintf("<link rel=\"stylesheet\" href=\"../%s/%s/highlights/github.min.css\">\n", config.AssetsSubDir, config.CssSubDir))
	gistHTMLBuffer.WriteString(fmt.Sprintf("<script src=\"../%s/%s/%s\"></script>\n", config.AssetsSubDir, config.JsSubDir, config.JsHighlightFileName))
	gistHTMLBuffer.WriteString(fmt.Sprintf("<script src=\"../%s/%s/%s\" defer></script>\n", config.AssetsSubDir, config.JsSubDir, config.JsNotesFileName))
	gistHTMLBuffer.WriteString("<script>hljs.highlightAll();</script>\n")

	gistHTMLBuffer.WriteString("\n</div>\n</body>\n</html>")

	if err := os.MkdirAll(filepath.Dir(fullGistHTMLPath), 0o755); err != nil {
		log.Printf("Failed to create directory for HTML gist file %s: %v", fullGistHTMLPath, err)
		return
	}
	if err := os.WriteFile(fullGistHTMLPath, gistHTMLBuffer.Bytes(), 0o644); err != nil {
		log.Printf("Error writing HTML for gist %s: %v", fullGistHTMLPath, err)
		return
	}
	log.Printf("Processed gist: %s (%s) -> HTML: %s", gist.ID, title, fullGistHTMLPath)

	var gistDate time.Time
	if gist.CreatedAt != "" {
		parsedDate, err := time.Parse(time.RFC3339, gist.CreatedAt)
		if err != nil {
			log.Printf("Warning: Could not parse date string '%s' for gist %s: %v", gist.CreatedAt, gist.ID, err)
		} else {
			gistDate = parsedDate
		}
	}

	indexEntryChan <- model.IndexEntry{
		GistID:       gist.ID,
		Title:        title,
		MarkdownPath: filepath.ToSlash(relativeMdPath),
		HTMLPath:     filepath.ToSlash(relativeHTMLPath),
		Date:         gistDate,
		Tags:         gist.Tags,
	}

	var contentBuilder strings.Builder
	for _, file := range gist.Files {
		contentBuilder.WriteString(file.Content)
		contentBuilder.WriteString("\n")
	}

	searchIndexEntryChan <- model.SearchIndexEntry{
		ID:      gist.ID,
		Title:   title,
		Content: contentBuilder.String(),
		Tags:    gist.Tags,
	}
}

// GenerateIndexFiles creates the main index.md and index.html files.
func GenerateIndexFiles(entries []model.IndexEntry, searchEntries []model.SearchIndexEntry) {
	var mdContent strings.Builder
	mdContent.WriteString("# Gists Archive\n\n")
	if len(entries) == 0 {
		mdContent.WriteString("No gists found or processed.\n")
	} else {
		for _, entry := range entries {
			mdContent.WriteString(fmt.Sprintf("- [%s](./%s)\n", entry.Title, entry.MarkdownPath))
		}
	}
	indexMdFilePath := filepath.Join(config.MarkdownOutputDir, "index.md")
	if err := os.MkdirAll(filepath.Dir(indexMdFilePath), 0o755); err != nil {
		log.Fatalf("Failed to create directory for markdown index: %v", err)
	}
	if err := os.WriteFile(indexMdFilePath, []byte(mdContent.String()), 0o644); err != nil {
		log.Fatalf("Failed to write markdown index.md: %v", err)
	}
	log.Printf("Generated Markdown Index: %s", indexMdFilePath)

	var htmlIndexContent strings.Builder
	htmlIndexContent.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n")
	htmlIndexContent.WriteString("<title>Gists Archive</title>\n")
	htmlIndexContent.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1.0">`)
	htmlIndexContent.WriteString(fmt.Sprintf("<link rel=\"stylesheet\" type=\"text/css\" href=\"%s/%s/%s\">\n",
		config.AssetsSubDir, config.CssSubDir, config.CssFileName))
	htmlIndexContent.WriteString("</head>\n<body>\n<div class=\"container\">\n")
	htmlIndexContent.WriteString("<h1>Gists Archive</h1>\n")
	htmlIndexContent.WriteString("<input type=\"text\" id=\"filterInput\" placeholder=\"Search gists...\">\n")
	htmlIndexContent.WriteString("<ul id=\"gistList\">\n")

	if len(entries) == 0 {
		htmlIndexContent.WriteString("<li>No gists found or processed.</li>\n")
	} else {
		for _, entry := range entries {
			dateAndTags := ""
			if !entry.Date.IsZero() {
				dateAndTags += fmt.Sprintf("<br><small class=\"date\"><b>Created:</b> %s</small>", entry.Date.Format("January 2, 2006"))
			}
			if entry.Tags != "" {
				dateAndTags += fmt.Sprintf("<br><small class=\"date\"><b>Tags:</b> %s</small>", entry.Tags)
			}
			htmlIndexContent.WriteString(fmt.Sprintf("<li class=\"gist-entry\" id=\"gist-%s\"><a href=\"%s\">%s</a>%s</li>\n", entry.GistID, entry.HTMLPath, entry.Title, dateAndTags))
		}
	}
	htmlIndexContent.WriteString("</ul>\n")

	// --- JAVASCRIPT ---
	// Embed search data as a global variable before other scripts
	searchData, err := json.Marshal(searchEntries)
	if err != nil {
		log.Fatalf("Failed to marshal search index entries for embedding: %v", err)
	}
	htmlIndexContent.WriteString("<script>window.SEARCH_DATA = ")
	htmlIndexContent.Write(searchData)
	htmlIndexContent.WriteString(";</script>\n")

	htmlIndexContent.WriteString(fmt.Sprintf("<script src=\"%s/%s/lunr.min.js\"></script>\n", config.AssetsSubDir, config.JsSubDir))
	htmlIndexContent.WriteString(fmt.Sprintf("<script src=\"%s/%s/%s\" defer></script>\n", config.AssetsSubDir, config.JsSubDir, config.JsSearchFileName))
	htmlIndexContent.WriteString(fmt.Sprintf("<script src=\"%s/%s/%s\" defer></script>\n", config.AssetsSubDir, config.JsSubDir, config.JsTagsFileName))

	htmlIndexContent.WriteString("</div>\n</body>\n</html>")

	indexHTMLFilePath := filepath.Join(config.WebDeployDir, "index.html")
	if err := os.MkdirAll(filepath.Dir(indexHTMLFilePath), 0o755); err != nil {
		log.Fatalf("Failed to create directory for HTML index: %v", err)
	}
	if err := os.WriteFile(indexHTMLFilePath, []byte(htmlIndexContent.String()), 0o644); err != nil {
		log.Fatalf("Failed to write HTML index file %s: %v", indexHTMLFilePath, err)
	}
	log.Printf("Generated HTML Index: %s", indexHTMLFilePath)
}
