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
	"runtime" // Added for sorting filenames
	"sort"
	"strings"
	"sync"
	"time"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Constants (apiBaseURL, gistsPerPage, githubUsername, etc.)
// (Remain the same as your last provided version)
const (
	apiBaseURL     = "https://api.github.com"
	gistsPerPage   = 100
	githubUsername = "integralist" // Replace if needed

	webDeployDir    = "."
	assetsSubDir    = "assets"
	cssSubDir       = "css"
	gistsHtmlSubDir = "gists"

	markdownOutputDir = "markdown_files"
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

/* Specific styling for H2 when used as a file section header */
.gist-file-header {
    margin-top: 30px; /* More space before a new file section */
    font-size: 1.3em; /* Slightly smaller than primary H2 if needed */
    color: #1f2328;
}


h1:first-child, h2.gist-file-header:first-child { /* Adjust if .gist-file-header is the first */
    margin-top: 0;
}

h1 { /* Page title / Gist title */
    font-size: 2em;
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
    margin-top: 8px; /* Add some space above pre if it follows a heading */
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

/* Styles for rendered Markdown content within a gist page */
.markdown-content {
    margin-top: 8px; /* Space between file header and rendered markdown */
}
.markdown-content h1, .markdown-content h2, .markdown-content h3, .markdown-content h4, .markdown-content h5, .markdown-content h6 {
    /* Reset gist-file-header specific margins if needed or add specific styling */
    /* For example, make them less prominent than the file header */
    margin-top: 1em;
    margin-bottom: 0.5em;
    border-bottom: none; /* Remove double border if h2 is used in MD */
}
.markdown-content h1 { font-size: 1.6em; }
.markdown-content h2 { font-size: 1.4em; }
.markdown-content h3 { font-size: 1.2em; }

`

var (
	httpClient = &http.Client{Timeout: 30 * time.Second}
	linkRegex  = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)
)

// Struct definitions (GistListEntry, GistListFile, GistDetail, GistFile, GistOwner, IndexEntry)
// (Remain unchanged from the previous version with date support)
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
	Type      string `json:"type"` // The API often provides a MIME type here
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

// API fetching functions (makeAPIRequest, makeAPIRequestForGistList, fetchAndProcessGistDetails, convertListFilesToDetailFiles)
// (Remain unchanged)
func makeAPIRequest(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "go-gist-archiver/1.6") // Version bump
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
	req.Header.Set("User-Agent", "go-gist-archiver/1.6")
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

// isMarkdownFile checks if a Gist file should be treated as Markdown.
func isMarkdownFile(file GistFile) bool {
	// Prioritize the Language field from the API
	if strings.ToLower(file.Language) == "markdown" {
		return true
	}
	// Fallback to checking file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return ext == ".md" || ext == ".markdown"
}

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
	var mdPageContent strings.Builder // This will be the content for the *entire Gist page*
	pageTitle := getGistTitle(gist)
	mdPageContent.WriteString(fmt.Sprintf("# %s\n\n", pageTitle)) // Main title for the Gist (becomes <h1>)

	if len(gist.Files) == 0 {
		mdPageContent.WriteString("This gist has no files.\n\n")
	} else {
		var filenames []string
		for fname := range gist.Files {
			filenames = append(filenames, fname)
		}
		sort.Strings(filenames) // Alphabetical sort for consistent file order

		for _, filename := range filenames {
			file := gist.Files[filename]
			isMD := isMarkdownFile(file)
			// You can uncomment this log for debugging file identification:
			// log.Printf("File: '%s', Language: '%s', Type: '%s', IsMD: %t", filename, file.Language, file.Type, isMD)

			fileContent := file.Content
			if file.Truncated && fileContent == "" {
				// Make the truncated message stand out a bit more if it's raw Markdown
				if isMD {
					fileContent = fmt.Sprintf("\n\n> **Note:** Content of this file is truncated. Full content available at %s\n\n", file.RawURL)
				} else {
					fileContent = fmt.Sprintf("*Content of this file is truncated. Full content available at %s*", file.RawURL)
				}
			}

			// Add filename as a Markdown H2 heading
			mdPageContent.WriteString(fmt.Sprintf("## %s\n\n", filename)) // Becomes <h2>

			if isMD {
				// For Markdown files, append their content directly.
				// This content is already Markdown and will be parsed as such.
				mdPageContent.WriteString(fileContent)
				mdPageContent.WriteString("\n\n") // Ensure separation
			} else {
				// For other files (code, text, etc.), wrap in a Markdown code block
				lang := file.Language
				if lang == "" { // Default language for code block if not specified
					ext := strings.TrimPrefix(filepath.Ext(filename), ".")
					if ext != "" {
						lang = ext // Guess from extension
					} else {
						lang = "text" // Ultimate fallback
					}
				}
				mdPageContent.WriteString(fmt.Sprintf("```%s\n", strings.ToLower(lang)))
				mdPageContent.WriteString(fileContent)
				mdPageContent.WriteString("\n```\n\n")
			}
		}
	}

	// This mdPageContent is now a "pure" Markdown string representing the whole Gist page.
	// It will be saved to a temporary .md file, and then processSingleGist will convert this .md to HTML.
	mdFilenameForTempStore := fmt.Sprintf("%s.md", gist.ID)
	filePath := filepath.Join(markdownOutputDir, gistsMdSubDir, mdFilenameForTempStore)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", filepath.Dir(filePath), err)
	}
	if err := os.WriteFile(filePath, []byte(mdPageContent.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to write composite markdown file %s: %w", filePath, err)
	}
	// Return the relative path of this composite markdown file (relative to markdownOutputDir)
	return filepath.Join(gistsMdSubDir, mdFilenameForTempStore), nil
}

// processSingleGist - Ensure parser settings are as you had them (without the non-existent flags)
func processSingleGist(gist GistDetail, indexEntryChan chan<- IndexEntry) {
	title := getGistTitle(gist)
	relativeMdPath, err := generateGistMarkdownFile(gist) // This now generates a composite MD for the whole page
	if err != nil {
		log.Printf("Error generating composite markdown for gist %s (%s): %v. Skipping.", gist.ID, title, err)
		return
	}

	htmlFilenameInGistsDir := strings.Replace(filepath.Base(relativeMdPath), ".md", ".html", 1)
	relativeHtmlPath := filepath.Join(gistsHtmlSubDir, htmlFilenameInGistsDir)
	fullGistHtmlPath := filepath.Join(webDeployDir, relativeHtmlPath)

	// Read the composite Markdown file we just created
	fullGistCompositeMdPath := filepath.Join(markdownOutputDir, relativeMdPath)
	mdInput, err := os.ReadFile(fullGistCompositeMdPath)
	if err != nil {
		log.Printf("Error reading composite gist markdown for HTML conversion %s: %v", fullGistCompositeMdPath, err)
		return
	}

	var gistHtmlBuffer bytes.Buffer
	gistHtmlBuffer.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n")
	gistHtmlBuffer.WriteString(fmt.Sprintf("<title>%s</title>\n", title))
	gistHtmlBuffer.WriteString(fmt.Sprintf("<link rel=\"stylesheet\" type=\"text/css\" href=\"../%s/%s/%s\">\n",
		assetsSubDir, cssSubDir, cssFileName))
	gistHtmlBuffer.WriteString("</head>\n<body>\n<div class=\"container\">\n")
	gistHtmlBuffer.WriteString("<a href=\"../index.html\" class=\"back-link\">&laquo; Back to Index</a>\n")

	// Parse the entire composite Markdown.
	// The parser settings from your code (without the non-existent flags) are correct here.
	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.FencedCode)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.UseXHTML | html.Smartypants
	opts := html.RendererOptions{Flags: htmlFlags, Title: title}
	renderer := html.NewRenderer(opts)

	renderedMarkdown := md.Render(p.Parse(mdInput), renderer) // mdInput is the composite Markdown

	gistHtmlBuffer.Write(renderedMarkdown) // This now contains rendered HTML from all files
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
	// Generate Index Markdown
	var mdContent strings.Builder
	mdContent.WriteString("# Gists Archive\n\n")
	if len(entries) == 0 {
		mdContent.WriteString("No gists found or processed.\n")
	} else {
		for _, entry := range entries {
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

	// Generate Index HTML
	var htmlIndexContent strings.Builder
	htmlIndexContent.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n")
	htmlIndexContent.WriteString("<title>Gists Archive</title>\n")
	htmlIndexContent.WriteString(fmt.Sprintf("<link rel=\"stylesheet\" type=\"text/css\" href=\"%s/%s/%s\">\n",
		assetsSubDir, cssSubDir, cssFileName))
	htmlIndexContent.WriteString("</head>\n<body>\n<div class=\"container\">\n")
	htmlIndexContent.WriteString("<h1>Gists Archive</h1>\n<ul>\n")

	if len(entries) == 0 {
		htmlIndexContent.WriteString("<li>No gists found or processed.</li>\n")
	} else {
		for _, entry := range entries {
			dateStr := ""
			if !entry.Date.IsZero() {
				dateStr = fmt.Sprintf("<br><small class=\"date\">Created: %s</small>", entry.Date.Format("January 2, 2006"))
			}
			htmlIndexContent.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a>%s</li>\n", entry.HTMLPath, entry.Title, dateStr))
		}
	}
	htmlIndexContent.WriteString("</ul>\n</div>\n</body>\n</html>")

	indexHtmlFilePath := filepath.Join(webDeployDir, "index.html")
	if err := os.MkdirAll(filepath.Dir(indexHtmlFilePath), 0755); err != nil {
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

	if err := os.MkdirAll(filepath.Join(markdownOutputDir, gistsMdSubDir), 0755); err != nil {
		log.Fatalf("Failed to create markdown gists directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(webDeployDir, gistsHtmlSubDir), 0755); err != nil {
		log.Fatalf("Failed to create HTML gists directory: %v", err)
	}

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
	log.Printf("Markdown files (source for HTML) in: ./%s", markdownOutputDir) // Clarified this line
	log.Printf("Open ./index.html in your browser to view the archive.")
	log.Println("-----------------------------------------------------")
}
