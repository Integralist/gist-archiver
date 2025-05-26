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
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const (
	apiBaseURL     = "https://api.github.com"
	gistsPerPage   = 100
	githubUsername = "integralist" // Replace if needed

	// Output structure constants
	webDeployDir    = "." // MODIFIED: Root for deployable web content is current directory
	assetsSubDir    = "assets"
	cssSubDir       = "css"
	gistsHtmlSubDir = "gists"

	markdownOutputDir = "markdown_files" // Separate directory for generated markdown (e.g., ./markdown_files/)
	gistsMdSubDir     = "gists"

	cssFileName          = "styles.css"
	maxConcurrentFetches = 10
)

const cssContent = `
/* General body and container styles */
body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
    margin: 20px;
    line-height: 1.6;
    background-color: #f4f4f4;
    color: #333;
}

.container {
    max-width: 800px;
    margin: auto;
    background: #fff;
    padding: 20px 30px;
    border-radius: 8px;
    box-shadow: 0 2px 15px rgba(0,0,0,0.1);
}

/* Headings */
h1, h2 {
    color: #333;
    border-bottom: 1px solid #eaecef;
    padding-bottom: 0.3em;
    margin-top: 24px;
    margin-bottom: 16px;
}

h1:first-child, h2:first-child {
    margin-top: 0;
}

h1 {
    font-size: 2em;
}

h2 { /* For filenames */
    font-size: 1.5em;
}

/* Links */
a {
    color: #0366d6;
    text-decoration: none;
}

a:hover {
    text-decoration: underline;
}

/* Index Page List */
.container > ul {
    list-style-type: none;
    padding: 0;
}

.container > ul > li {
    padding: 8px 0;
    border-bottom: 1px solid #eee;
}
.container > ul > li:last-child {
    border-bottom: none;
}
.container > ul > li > small.date { /* Style for the date */
    display: block;
    font-size: 0.8em;
    color: #666;
    margin-top: 4px;
}


/* Gist page specific "back to index" link */
a.back-link {
    display: inline-block;
    margin-bottom: 20px;
    font-size: 0.9em;
}

/* Code blocks */
pre {
    background-color: #f6f8fa;
    padding: 16px;
    overflow: auto;
    font-size: 85%;
    line-height: 1.45;
    border-radius: 6px;
    border: 1px solid #ddd;
    margin-bottom: 16px;
}

code {
    font-family: 'SFMono-Regular', Consolas, "Liberation Mono", Menlo, Courier, monospace;
    font-size: inherit;
}

:not(pre) > code {
  padding: .2em .4em;
  margin: 0 2px;
  font-size: 85%;
  background-color: rgba(27,31,35,0.07);
  border-radius: 3px;
}
`

var (
	httpClient = &http.Client{Timeout: 30 * time.Second}
	linkRegex  = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)
)

// --- Struct definitions ---
type GistListEntry struct {
	ID          string                  `json:"id"`
	Description string                  `json:"description"`
	Public      bool                    `json:"public"`
	Owner       GistOwner               `json:"owner"`
	Files       map[string]GistListFile `json:"files"`
}

type GistListFile struct {
	Filename string `json:"filename"`
}

type GistDetail struct {
	ID          string              `json:"id"`
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Owner       GistOwner           `json:"owner"`
	Files       map[string]GistFile `json:"files"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}

type GistFile struct {
	Filename  string `json:"filename"`
	Language  string `json:"language"`
	Content   string `json:"content"`
	RawURL    string `json:"raw_url"`
	Size      int    `json:"size"`
	Truncated bool   `json:"truncated"`
}

type GistOwner struct {
	Login string `json:"login"`
}

type IndexEntry struct {
	Title        string
	MarkdownPath string
	HTMLPath     string
	Date         time.Time
}

// --- API fetching functions ---
func makeAPIRequest(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "go-gist-archiver/1.5") // Version bump
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden && strings.Contains(resp.Header.Get("X-RateLimit-Remaining"), "0") {
		log.Printf("Rate limit hit for URL %s. Status: %s", url, resp.Status)
		return fmt.Errorf("rate limit hit for %s. X-RateLimit-Remaining: %s, X-RateLimit-Reset: %s", url, resp.Header.Get("X-RateLimit-Remaining"), resp.Header.Get("X-RateLimit-Reset"))
	}

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

func makeAPIRequestForGistList(url string, target any) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "go-gist-archiver/1.5")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden && strings.Contains(resp.Header.Get("X-RateLimit-Remaining"), "0") {
		log.Printf("Rate limit hit for URL %s. Status: %s", url, resp.Status)
		return "", fmt.Errorf("rate limit hit for %s. X-RateLimit-Remaining: %s, X-RateLimit-Reset: %s", url, resp.Header.Get("X-RateLimit-Remaining"), resp.Header.Get("X-RateLimit-Reset"))
	}

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
			return matches[1], nil
		}
	}
	return "", nil
}

func fetchAndProcessGistDetails(username string, gistDetailChan chan<- GistDetail, fetchWg *sync.WaitGroup) {
	defer close(gistDetailChan)

	nextURL := fmt.Sprintf("%s/users/%s/gists?per_page=%d", apiBaseURL, username, gistsPerPage)
	page := 1
	fetchSemaphore := make(chan struct{}, maxConcurrentFetches)

	for nextURL != "" {
		log.Printf("Fetching page %d of gist list for user %s from %s", page, username, nextURL)
		var gistListEntries []GistListEntry
		var err error
		nextURL, err = makeAPIRequestForGistList(nextURL, &gistListEntries)
		if err != nil {
			log.Printf("Error: Failed to fetch gist list page %d: %v. Stopping further list fetching.", page, err)
			return
		}

		for _, entry := range gistListEntries {
			if !entry.Public {
				continue
			}

			fetchWg.Add(1)
			go func(gistEntry GistListEntry) {
				defer fetchWg.Done()
				fetchSemaphore <- struct{}{}
				defer func() { <-fetchSemaphore }()

				log.Printf("Fetching details for gist ID: %s (%s)", gistEntry.ID, getGistTitle(GistDetail{Description: gistEntry.Description, Files: convertListFilesToDetailFiles(gistEntry.Files), ID: gistEntry.ID}))
				var gistDetail GistDetail
				detailURL := fmt.Sprintf("%s/gists/%s", apiBaseURL, gistEntry.ID)
				if err := makeAPIRequest(detailURL, &gistDetail); err != nil {
					log.Printf("Warning: Failed to fetch detail for gist %s: %v. Skipping.", gistEntry.ID, err)
					return
				}
				gistDetailChan <- gistDetail
			}(entry)
		}
		page++
		if nextURL != "" {
			time.Sleep(250 * time.Millisecond)
		}
	}
	fetchWg.Wait()
}

func convertListFilesToDetailFiles(listFiles map[string]GistListFile) map[string]GistFile {
	detailFiles := make(map[string]GistFile)
	for k, v := range listFiles {
		detailFiles[k] = GistFile{Filename: v.Filename}
	}
	return detailFiles
}

// --- Utility and Markdown/HTML generation functions ---
func getGistTitle(gist GistDetail) string {
	if gist.Description != "" {
		return gist.Description
	}
	for _, file := range gist.Files {
		return file.Filename
	}
	return gist.ID
}

func generateGistMarkdownFile(gist GistDetail) (string, error) {
	var mdContent strings.Builder
	title := getGistTitle(gist)
	mdContent.WriteString(fmt.Sprintf("# %s\n\n", title))

	if len(gist.Files) == 0 {
		mdContent.WriteString("This gist has no files.\n")
	}

	for filename, file := range gist.Files {
		lang := file.Language
		if lang == "" {
			lang = "text"
		}
		fileContent := file.Content
		if file.Truncated && fileContent == "" {
			fileContent = fmt.Sprintf("File content truncated. View original at %s", file.RawURL)
		}
		mdContent.WriteString(fmt.Sprintf("## %s\n\n```%s\n%s\n```\n\n", filename, strings.ToLower(lang), fileContent))
	}

	mdFilename := fmt.Sprintf("%s.md", gist.ID)
	filePath := filepath.Join(markdownOutputDir, gistsMdSubDir, mdFilename) // e.g. ./markdown_files/gists/id.md

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", filepath.Dir(filePath), err)
	}
	if err := os.WriteFile(filePath, []byte(mdContent.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to write markdown file %s: %w", filePath, err)
	}
	return filepath.Join(gistsMdSubDir, mdFilename), nil // Returns "gists/id.md"
}

func processSingleGist(gist GistDetail, indexEntryChan chan<- IndexEntry) {
	title := getGistTitle(gist)
	// relativeMdPath is "gists/gist_id.md" (relative to markdownOutputDir)
	relativeMdPath, err := generateGistMarkdownFile(gist)
	if err != nil {
		log.Printf("Error generating markdown for gist %s (%s): %v. Skipping.", gist.ID, title, err)
		return
	}

	htmlFilenameInGistsDir := strings.Replace(filepath.Base(relativeMdPath), ".md", ".html", 1)
	// relativeHtmlPath is "gists/gist_id.html" (relative to webDeployDir which is ".")
	relativeHtmlPath := filepath.Join(gistsHtmlSubDir, htmlFilenameInGistsDir)
	// fullGistHtmlPath is "./gists/gist_id.html"
	fullGistHtmlPath := filepath.Join(webDeployDir, relativeHtmlPath)

	fullGistMdPath := filepath.Join(markdownOutputDir, relativeMdPath)
	mdInput, err := os.ReadFile(fullGistMdPath)
	if err != nil {
		log.Printf("Error reading gist markdown for HTML conversion %s: %v", fullGistMdPath, err)
		return
	}

	var gistHtmlBuffer bytes.Buffer
	gistHtmlBuffer.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n")
	gistHtmlBuffer.WriteString(fmt.Sprintf("<title>%s</title>\n", title))
	// CSS link relative from ./gists/file.html to ./assets/css/styles.css is "../assets/css/styles.css"
	gistHtmlBuffer.WriteString(fmt.Sprintf("<link rel=\"stylesheet\" type=\"text/css\" href=\"../%s/%s/%s\">\n",
		assetsSubDir, cssSubDir, cssFileName))
	gistHtmlBuffer.WriteString("</head>\n<body>\n<div class=\"container\">\n")
	// Back link relative from ./gists/file.html to ./index.html is "../index.html"
	gistHtmlBuffer.WriteString("<a href=\"../index.html\" class=\"back-link\">&laquo; Back to Index</a>\n")

	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.FencedCode)
	doc := p.Parse(mdInput)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.UseXHTML
	opts := html.RendererOptions{Flags: htmlFlags, Title: title}
	renderer := html.NewRenderer(opts)
	renderedMarkdown := md.Render(doc, renderer)

	gistHtmlBuffer.Write(renderedMarkdown)
	gistHtmlBuffer.WriteString("\n</div>\n</body>\n</html>")

	if err := os.MkdirAll(filepath.Dir(fullGistHtmlPath), 0755); err != nil {
		log.Printf("Failed to create directory for HTML gist file %s: %v", fullGistHtmlPath, err)
		return
	}
	if err := os.WriteFile(fullGistHtmlPath, gistHtmlBuffer.Bytes(), 0644); err != nil {
		log.Printf("Error writing HTML for gist %s: %v", fullGistHtmlPath, err)
		return
	}
	log.Printf("Processed gist: %s (%s) -> HTML: %s", gist.ID, title, fullGistHtmlPath)

	var gistDate time.Time
	if gist.CreatedAt != "" {
		parsedDate, err := time.Parse(time.RFC3339, gist.CreatedAt)
		if err != nil {
			log.Printf("Warning: Could not parse date string '%s' for gist %s: %v", gist.CreatedAt, gist.ID, err)
		} else {
			gistDate = parsedDate
		}
	}

	indexEntryChan <- IndexEntry{
		Title:        title,
		MarkdownPath: filepath.ToSlash(relativeMdPath),
		HTMLPath:     filepath.ToSlash(relativeHtmlPath),
		Date:         gistDate,
	}
}

func generateIndexFiles(entries []IndexEntry) {
	// Generate Index Markdown (in ./markdown_files/index.md)
	var mdContent strings.Builder
	mdContent.WriteString("# Integralist's Gists Archive\n\n")
	if len(entries) == 0 {
		mdContent.WriteString("No gists found or processed.\n")
	} else {
		for _, entry := range entries {
			// entry.MarkdownPath is "gists/id.md", link should be "./gists/id.md" from markdown_files/index.md
			mdContent.WriteString(fmt.Sprintf("- [%s](./%s)\n", entry.Title, entry.MarkdownPath))
		}
	}
	indexMdFilePath := filepath.Join(markdownOutputDir, "index.md")
	if err := os.MkdirAll(filepath.Dir(indexMdFilePath), 0755); err != nil {
		log.Fatalf("Failed to create directory for markdown index: %v", err)
	}
	if err := os.WriteFile(indexMdFilePath, []byte(mdContent.String()), 0644); err != nil {
		log.Fatalf("Failed to write markdown index.md: %v", err)
	}
	log.Printf("Generated Markdown Index: %s", indexMdFilePath)

	// Generate Index HTML (in ./index.html)
	var htmlIndexContent strings.Builder
	htmlIndexContent.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n")
	htmlIndexContent.WriteString("<title>Integralist's Gists Archive</title>\n")
	// CSS link relative from ./index.html to ./assets/css/styles.css is "assets/css/styles.css"
	htmlIndexContent.WriteString(fmt.Sprintf("<link rel=\"stylesheet\" type=\"text/css\" href=\"%s/%s/%s\">\n",
		assetsSubDir, cssSubDir, cssFileName))
	htmlIndexContent.WriteString("</head>\n<body>\n<div class=\"container\">\n")
	htmlIndexContent.WriteString("<h1>Integralist's Gists Archive</h1>\n<ul>\n")

	if len(entries) == 0 {
		htmlIndexContent.WriteString("<li>No gists found or processed.</li>\n")
	} else {
		for _, entry := range entries {
			dateStr := ""
			if !entry.Date.IsZero() {
				dateStr = fmt.Sprintf("<br><small class=\"date\">Created: %s</small>", entry.Date.Format("January 2, 2006"))
			}
			// entry.HTMLPath is "gists/file.html", correct for link from ./index.html
			htmlIndexContent.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a>%s</li>\n", entry.HTMLPath, entry.Title, dateStr))
		}
	}
	htmlIndexContent.WriteString("</ul>\n</div>\n</body>\n</html>")

	indexHtmlFilePath := filepath.Join(webDeployDir, "index.html") // Results in "./index.html"
	if err := os.MkdirAll(filepath.Dir(indexHtmlFilePath), 0755); err != nil {
		// For webDeployDir = ".", filepath.Dir will be ".", MkdirAll for "." is a no-op and safe.
		log.Fatalf("Failed to create directory for HTML index: %v", err)
	}
	if err := os.WriteFile(indexHtmlFilePath, []byte(htmlIndexContent.String()), 0644); err != nil {
		log.Fatalf("Failed to write HTML index file %s: %v", indexHtmlFilePath, err)
	}
	log.Printf("Generated HTML Index: %s", indexHtmlFilePath)
}

func main() {
	startTime := time.Now()
	log.Printf("Starting gist archival for user: %s", githubUsername)

	// Create base output directories
	// Markdown files go into a subdirectory
	if err := os.MkdirAll(filepath.Join(markdownOutputDir, gistsMdSubDir), 0755); err != nil {
		log.Fatalf("Failed to create markdown gists directory: %v", err)
	}
	// HTML gists go into a subdirectory of the current dir
	if err := os.MkdirAll(filepath.Join(webDeployDir, gistsHtmlSubDir), 0755); err != nil {
		log.Fatalf("Failed to create HTML gists directory: %v", err)
	}

	// Assets directory (./assets/css/)
	cssDirPath := filepath.Join(webDeployDir, assetsSubDir, cssSubDir)
	if err := os.MkdirAll(cssDirPath, 0755); err != nil {
		log.Fatalf("Failed to create CSS assets directory %s: %v", cssDirPath, err)
	}
	cssFilePath := filepath.Join(cssDirPath, cssFileName)
	if err := os.WriteFile(cssFilePath, []byte(strings.TrimSpace(cssContent)), 0644); err != nil {
		log.Fatalf("Failed to write CSS file %s: %v", cssFilePath, err)
	}
	log.Printf("Generated CSS file: %s", cssFilePath)

	gistDetailChan := make(chan GistDetail, maxConcurrentFetches)
	indexEntryChan := make(chan IndexEntry, runtime.NumCPU()*2)

	var fetchWg sync.WaitGroup
	var processWg sync.WaitGroup

	numProcessors := runtime.NumCPU()
	if numProcessors < 1 {
		numProcessors = 1
	}
	log.Printf("Launching %d gist processing workers.", numProcessors)
	for i := 0; i < numProcessors; i++ {
		processWg.Add(1)
		go func(workerID int) {
			defer processWg.Done()
			for gist := range gistDetailChan {
				processSingleGist(gist, indexEntryChan)
			}
		}(i)
	}

	var collectedIndexEntries []IndexEntry
	var indexWg sync.WaitGroup
	indexWg.Add(1)
	go func() {
		defer indexWg.Done()
		for entry := range indexEntryChan {
			collectedIndexEntries = append(collectedIndexEntries, entry)
		}
		log.Printf("Collected %d index entries.", len(collectedIndexEntries))

		sort.Slice(collectedIndexEntries, func(i, j int) bool {
			if collectedIndexEntries[i].Date.IsZero() {
				return false
			}
			if collectedIndexEntries[j].Date.IsZero() {
				return true
			}
			return collectedIndexEntries[j].Date.Before(collectedIndexEntries[i].Date)
		})
		log.Println("Index entries sorted by date (newest first).")
		generateIndexFiles(collectedIndexEntries)
	}()

	fetchAndProcessGistDetails(githubUsername, gistDetailChan, &fetchWg)
	log.Println("All gist detail fetching goroutines launched and channel closed by fetcher.")

	processWg.Wait()
	log.Println("All gist processing workers have completed.")

	close(indexEntryChan)
	log.Println("Index entry channel closed.")

	indexWg.Wait()
	log.Println("Index file generation complete.")

	log.Println("-----------------------------------------------------")
	log.Printf("Processing complete in %v!", time.Since(startTime))
	log.Printf("Web content (HTML, assets) in: current directory (./)")
	log.Printf("  - Index: ./index.html")
	log.Printf("  - Gists: ./%s/", gistsHtmlSubDir)
	log.Printf("  - Assets: ./%s/", assetsSubDir)
	log.Printf("Markdown files in:             ./%s", markdownOutputDir)
	log.Printf("Open ./index.html in your browser to view the archive.")
	log.Println("-----------------------------------------------------")
}
