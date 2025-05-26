package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const (
	apiBaseURL      = "https://api.github.com"
	gistsHtmlSubDir = "gists"
	gistsMdSubDir   = "gists"
	gistsPerPage    = 100
	githubUsername  = "integralist"
	htmlBaseDir     = "output/html"
	markdownBaseDir = "output/markdown"
)

var (
	httpClient = &http.Client{Timeout: 20 * time.Second}
	linkRegex  = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)
)

// GistListEntry is for the paginated list of gists (metadata only)
type GistListEntry struct {
	ID          string                  `json:"id"`
	Description string                  `json:"description"`
	Public      bool                    `json:"public"`
	Owner       GistOwner               `json:"owner"`
	Files       map[string]GistListFile `json:"files"` // to get a fallback title if description is empty
}

// GistListFile provides minimal file info from the list endpoint
type GistListFile struct {
	Filename string `json:"filename"`
}

// GistDetail represents a single Gist with full file content
type GistDetail struct {
	ID          string              `json:"id"`
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Owner       GistOwner           `json:"owner"`
	Files       map[string]GistFile `json:"files"`
}

// GistFile represents a file within a Gist, including its content
type GistFile struct {
	Filename  string `json:"filename"`
	Language  string `json:"language"`
	Content   string `json:"content"`
	RawURL    string `json:"raw_url"` // optional fallback if content is missing
	Size      int    `json:"size"`
	Truncated bool   `json:"truncated"` // indicates if content is too large for direct inclusion
}

// GistOwner represents the owner of a Gist
type GistOwner struct {
	Login string `json:"login"`
}

// IndexEntry holds info for the index file
type IndexEntry struct {
	Title        string
	MarkdownPath string // Relative to markdownBaseDir
	HTMLPath     string // Relative to htmlBaseDir
}

func main() {
	log.Printf("Starting gist archival for user: %s", githubUsername)

	// Create base output directories
	if err := os.MkdirAll(filepath.Join(markdownBaseDir, gistsMdSubDir), 0755); err != nil {
		log.Fatalf("Failed to create markdown output directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(htmlBaseDir, gistsHtmlSubDir), 0755); err != nil {
		log.Fatalf("Failed to create HTML output directory: %v", err)
	}

	// 1. Fetch all Gist details
	allGists, err := fetchAllGistDetails(githubUsername)
	if err != nil {
		log.Fatalf("Error fetching gists: %v", err)
	}

	if len(allGists) == 0 {
		log.Println("No public gists found for this user.")
		// Create empty index files if no gists
		emptyEntries := []IndexEntry{}
		_, err = generateIndexMarkdown(emptyEntries)
		if err != nil {
			log.Printf("Error generating empty markdown index: %v", err)
		} else {
			err = convertMarkdownToHTML(filepath.Join(markdownBaseDir, "index.md"), filepath.Join(htmlBaseDir, "index.html"))
			if err != nil {
				log.Printf("Error converting empty markdown index to HTML: %v", err)
			}
		}
		return
	}

	// 2. Generate individual Markdown files for each Gist
	var indexEntries []IndexEntry
	for _, gist := range allGists {
		if !gist.Public { // Double-check, though API should only return public ones
			continue
		}
		relativeMdPath, err := generateGistMarkdownFile(gist)
		if err != nil {
			log.Printf("Error generating markdown for gist %s: %v. Skipping.", gist.ID, err)
			continue
		}
		title := getGistTitle(gist)
		indexEntries = append(indexEntries, IndexEntry{
			Title:        title,
			MarkdownPath: filepath.ToSlash(relativeMdPath), // Ensure forward slashes for web/md links
			HTMLPath:     filepath.ToSlash(strings.Replace(relativeMdPath, ".md", ".html", 1)),
		})
	}

	// 3. Generate Index Markdown file
	relativeIndexMdPath, err := generateIndexMarkdown(indexEntries)
	if err != nil {
		log.Fatalf("Error generating index.md: %v", err)
	}

	// 4. Convert Markdown to HTML
	// Convert index.md to index.html
	fullIndexHtmlPath := filepath.Join(htmlBaseDir, strings.Replace(relativeIndexMdPath, ".md", ".html", 1))

	var htmlIndexContent strings.Builder
	htmlIndexContent.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n")
	htmlIndexContent.WriteString("<title>My Gists Archive</title>\n")
	// Basic styling
	htmlIndexContent.WriteString("<style>\n")
	htmlIndexContent.WriteString("body { font-family: sans-serif; margin: 20px; line-height: 1.6; background-color: #f4f4f4; color: #333; }\n")
	htmlIndexContent.WriteString("h1, h2 { color: #333; }\n")
	htmlIndexContent.WriteString(".container { max-width: 800px; margin: auto; background: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); }\n")
	htmlIndexContent.WriteString("ul {}\n")
	htmlIndexContent.WriteString("li a { text-decoration: none; color: #007bff; }\n")
	htmlIndexContent.WriteString("li a:hover { text-decoration: underline; }\n")
	htmlIndexContent.WriteString("pre { background: #eee; padding: 10px; border-radius: 4px; overflow-x: auto; }\n")
	htmlIndexContent.WriteString("code { font-family: monospace; }\n")
	htmlIndexContent.WriteString("</style>\n")
	htmlIndexContent.WriteString("</head>\n<body>\n<div class=\"container\">\n")
	htmlIndexContent.WriteString("<h1>My Gists Archive</h1>\n<ul>\n")

	for _, entry := range indexEntries {
		htmlIndexContent.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a></li>\n", entry.HTMLPath, entry.Title))
	}
	htmlIndexContent.WriteString("</ul>\n</div>\n</body>\n</html>")

	if err := os.WriteFile(fullIndexHtmlPath, []byte(htmlIndexContent.String()), 0644); err != nil {
		log.Fatalf("Failed to write HTML index file %s: %v", fullIndexHtmlPath, err)
	}
	log.Printf("Generated HTML Index: %s", fullIndexHtmlPath)

	// Convert individual gist markdown files to HTML
	for _, entry := range indexEntries {
		fullGistMdPath := filepath.Join(markdownBaseDir, entry.MarkdownPath)
		fullGistHtmlPath := filepath.Join(htmlBaseDir, entry.HTMLPath)

		mdInput, err := os.ReadFile(fullGistMdPath)
		if err != nil {
			log.Printf("Error reading gist markdown %s: %v", fullGistMdPath, err)
			continue
		}

		// Generate HTML content for individual gist page
		var gistHtmlBuffer bytes.Buffer
		gistHtmlBuffer.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n")
		gistHtmlBuffer.WriteString(fmt.Sprintf("<title>%s</title>\n", entry.Title))
		gistHtmlBuffer.WriteString("<style>\n") // Same style as index
		gistHtmlBuffer.WriteString("body { font-family: sans-serif; margin: 20px; line-height: 1.6; background-color: #f4f4f4; color: #333; }\n")
		gistHtmlBuffer.WriteString(".container { max-width: 800px; margin: auto; background: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); }\n")
		gistHtmlBuffer.WriteString("pre { background: #f0f0f0; padding: 1em; border-radius: 5px; overflow-x: auto; border: 1px solid #ddd; }\n")
		gistHtmlBuffer.WriteString("code { font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace; font-size: 0.9em;}\n")
		gistHtmlBuffer.WriteString("h1, h2 { border-bottom: 1px solid #eee; padding-bottom: 0.3em;}\n")
		gistHtmlBuffer.WriteString("a.back-link { display: inline-block; margin-bottom: 20px; color: #007bff; text-decoration: none; }\n")
		gistHtmlBuffer.WriteString("a.back-link:hover { text-decoration: underline; }\n")
		gistHtmlBuffer.WriteString("</style>\n")
		gistHtmlBuffer.WriteString("</head>\n<body>\n<div class=\"container\">\n")
		gistHtmlBuffer.WriteString("<a href=\"../index.html\" class=\"back-link\">&laquo; Back to Index</a>\n")

		p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.FencedCode)
		doc := p.Parse(mdInput)
		htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.UseXHTML // UseXHTML for self-closing tags if any
		opts := html.RendererOptions{Flags: htmlFlags, Title: entry.Title}
		renderer := html.NewRenderer(opts)
		renderedMarkdown := md.Render(doc, renderer)

		gistHtmlBuffer.Write(renderedMarkdown)
		gistHtmlBuffer.WriteString("\n</div>\n</body>\n</html>")

		if err := os.MkdirAll(filepath.Dir(fullGistHtmlPath), 0755); err != nil {
			log.Printf("Failed to create directory for HTML gist file %s: %v", fullGistHtmlPath, err)
			continue
		}
		if err := os.WriteFile(fullGistHtmlPath, gistHtmlBuffer.Bytes(), 0644); err != nil {
			log.Printf("Error converting gist markdown %s to HTML: %v", fullGistMdPath, err)
			continue
		}
		log.Printf("Converted %s to %s", fullGistMdPath, fullGistHtmlPath)
	}

	log.Println("-----------------------------------------------------")
	log.Println("Processing complete!")
	log.Printf("Markdown files are in: %s", filepath.Join(markdownBaseDir))
	log.Printf("HTML files are in: %s", filepath.Join(htmlBaseDir))
	log.Printf("Open %s in your browser to view the archive.", filepath.Join(htmlBaseDir, "index.html"))
	log.Println("-----------------------------------------------------")
}

func makeAPIRequest(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	// GitHub API requires a User-Agent header
	req.Header.Set("User-Agent", "go-gist-archiver/1.0")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api request to %s failed with status %s: %s", url, resp.Status, string(bodyBytes))
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode JSON response from %s: %w", url, err)
		}
	}
	return nil
}

func makeAPIRequestWithPagination(url string, target any) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "go-gist-archiver/1.0")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("api request to %s failed with status %s: %s", url, resp.Status, string(bodyBytes))
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return "", fmt.Errorf("failed to decode JSON response from %s: %w", url, err)
		}
	}

	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		matches := linkRegex.FindStringSubmatch(linkHeader)
		if len(matches) > 1 {
			return matches[1], nil // Return the next URL
		}
	}
	return "", nil // No next page
}

// fetchAllGistDetails fetches all public gists for a user
func fetchAllGistDetails(username string) ([]GistDetail, error) {
	var allGistDetails []GistDetail
	nextURL := fmt.Sprintf("%s/users/%s/gists?per_page=%d", apiBaseURL, username, gistsPerPage)
	page := 1

	for nextURL != "" {
		log.Printf("Fetching page %d of gists for user %s from %s", page, username, nextURL)
		var gistListEntries []GistListEntry
		var err error
		nextURL, err = makeAPIRequestWithPagination(nextURL, &gistListEntries)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch gist list page %d: %w", page, err)
		}

		for _, entry := range gistListEntries {
			if !entry.Public { // Ensure we only process public gists, though the endpoint should only return them for unauth user
				continue
			}
			log.Printf("Fetching details for gist ID: %s", entry.ID)
			var gistDetail GistDetail
			detailURL := fmt.Sprintf("%s/gists/%s", apiBaseURL, entry.ID)
			if err := makeAPIRequest(detailURL, &gistDetail); err != nil {
				log.Printf("Warning: Failed to fetch detail for gist %s: %v. Skipping.", entry.ID, err)
				continue // Skip this gist if detail fetch fails
			}
			allGistDetails = append(allGistDetails, gistDetail)
			time.Sleep(200 * time.Millisecond) // Be polite to the API, ~5 req/sec
		}
		page++
		if nextURL != "" {
			time.Sleep(500 * time.Millisecond) // Pause between fetching pages
		}
	}
	log.Printf("Fetched details for %d public gists.", len(allGistDetails))
	return allGistDetails, nil
}

func getGistTitle(gist GistDetail) string {
	if gist.Description != "" {
		return gist.Description
	}
	// Fallback to the first filename if description is empty
	for _, file := range gist.Files {
		return file.Filename // Returns the first one encountered
	}
	return gist.ID // Ultimate fallback
}

// generateGistMarkdownFile creates a markdown file for a single gist
func generateGistMarkdownFile(gist GistDetail) (string, error) {
	var mdContent strings.Builder
	title := getGistTitle(gist)
	mdContent.WriteString(fmt.Sprintf("# %s\n\n", title))

	if len(gist.Files) == 0 {
		mdContent.WriteString("This gist has no files.\n")
	}

	for _, file := range gist.Files {
		lang := file.Language
		if lang == "" { // Default to text if language is not specified
			lang = "text"
		}
		// Ensure content is present, especially if it could be truncated
		// For truncated files, GistFile.Content might be empty.
		// A more robust solution would fetch RawURL if file.Truncated is true.
		// This example assumes content is available or an empty block is acceptable.
		fileContent := file.Content
		if file.Truncated && fileContent == "" {
			fileContent = fmt.Sprintf("File content truncated. View original at %s", file.RawURL)
		}

		mdContent.WriteString(fmt.Sprintf("## %s\n\n", file.Filename))
		mdContent.WriteString(fmt.Sprintf("```%s\n", strings.ToLower(lang)))
		mdContent.WriteString(fileContent)
		mdContent.WriteString("\n```\n\n")
	}

	filename := fmt.Sprintf("%s.md", gist.ID)
	filePath := filepath.Join(markdownBaseDir, gistsMdSubDir, filename)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", filepath.Dir(filePath), err)
	}

	if err := os.WriteFile(filePath, []byte(mdContent.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to write markdown file %s: %w", filePath, err)
	}
	log.Printf("Generated Markdown: %s", filePath)
	return filepath.Join(gistsMdSubDir, filename), nil // Return relative path
}

// generateIndexMarkdown creates the main index.md file
func generateIndexMarkdown(entries []IndexEntry) (string, error) {
	var mdContent strings.Builder
	mdContent.WriteString("# My Gists Archive\n\n")

	if len(entries) == 0 {
		mdContent.WriteString("No gists found or processed.\n")
	}

	for _, entry := range entries {
		// Link relative to the index.md file itself
		mdContent.WriteString(fmt.Sprintf("- [%s](./%s)\n", entry.Title, entry.MarkdownPath))
	}

	filePath := filepath.Join(markdownBaseDir, "index.md")
	if err := os.WriteFile(filePath, []byte(mdContent.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to write index.md: %w", err)
	}
	log.Printf("Generated Markdown Index: %s", filePath)
	return "index.md", nil // Return relative path
}

// convertMarkdownToHTML converts a single markdown file to HTML
func convertMarkdownToHTML(mdFilePath, htmlFilePath string) error {
	mdInput, err := os.ReadFile(mdFilePath)
	if err != nil {
		return fmt.Errorf("failed to read markdown file %s: %w", mdFilePath, err)
	}

	// Setup gomarkdown parser and renderer
	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock)
	doc := p.Parse(mdInput)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank // Open external links in new tab
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	htmlOutput := md.Render(doc, renderer)

	if err := os.MkdirAll(filepath.Dir(htmlFilePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for html file %s: %w", htmlFilePath, err)
	}

	if err := os.WriteFile(htmlFilePath, htmlOutput, 0644); err != nil {
		return fmt.Errorf("failed to write html file %s: %w", htmlFilePath, err)
	}
	log.Printf("Converted %s to %s", mdFilePath, htmlFilePath)
	return nil
}
