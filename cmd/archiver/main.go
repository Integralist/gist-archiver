package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/integralist/gist-archiver/internal/config"
	"github.com/integralist/gist-archiver/internal/filegen"
	"github.com/integralist/gist-archiver/internal/github"
	"github.com/integralist/gist-archiver/internal/model"
)

func main() {
	startTime := time.Now()
	log.Printf("Starting gist archival for user: %s", config.GithubUsername)

	if err := os.MkdirAll(filepath.Join(config.MarkdownOutputDir, config.GistsMdSubDir), 0755); err != nil {
		log.Fatalf("Failed to create markdown gists directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(config.WebDeployDir, config.GistsHtmlSubDir), 0755); err != nil {
		log.Fatalf("Failed to create HTML gists directory: %v", err)
	}

	gistDetailChan := make(chan model.GistDetail, config.MaxConcurrentFetches)
	indexEntryChan := make(chan model.IndexEntry, runtime.NumCPU()*2)
	searchIndexEntryChan := make(chan model.SearchIndexEntry, runtime.NumCPU()*2)

	var fetchWg sync.WaitGroup
	var processWg sync.WaitGroup

	token := os.Getenv("API_TOKEN_GITHUB")
	client := github.NewClient(config.APIBaseURL, config.GithubUsername, token, config.GistsPerPage, config.MaxConcurrentFetches)

	go client.FetchAndProcessGistDetails(gistDetailChan, &fetchWg)

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
				filegen.ProcessSingleGist(gist, indexEntryChan, searchIndexEntryChan)
			}
		}(i)
	}

	var collectedIndexEntries []model.IndexEntry
	var searchIndexEntries []model.SearchIndexEntry
	var indexWg sync.WaitGroup
	indexWg.Add(1)
	go func() {
		defer indexWg.Done()
		for {
			select {
			case entry, ok := <-indexEntryChan:
				if !ok {
					indexEntryChan = nil
				} else {
					collectedIndexEntries = append(collectedIndexEntries, entry)
				}
			case searchEntry, ok := <-searchIndexEntryChan:
				if !ok {
					searchIndexEntryChan = nil
				} else {
					searchIndexEntries = append(searchIndexEntries, searchEntry)
				}
			}
			if indexEntryChan == nil && searchIndexEntryChan == nil {
				break
			}
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
		filegen.GenerateIndexFiles(collectedIndexEntries, searchIndexEntries)
	}()

	fetchWg.Wait()
	log.Println("All gist detail fetching goroutines launched and channel closed by fetcher.")

	processWg.Wait()
	log.Println("All gist processing workers have completed.")

	close(indexEntryChan)
	close(searchIndexEntryChan)
	log.Println("Index entry channel closed.")

	indexWg.Wait()
	log.Println("Index file generation complete.")

	log.Println("-----------------------------------------------------")
	log.Printf("Processing complete in %v!", time.Since(startTime))
	log.Printf("Web content (HTML, assets) in: current directory (./)")
	log.Printf("  - Index: ./index.html")
	log.Printf("  - Gists: ./%s/", config.GistsHtmlSubDir)
	log.Printf("  - Assets: ./%s/ (css & js files)", config.AssetsSubDir)
	log.Printf("Markdown files (source for HTML) in: ./%s", config.MarkdownOutputDir)
	log.Printf("Open ./index.html in your browser to view the archive.")
	log.Println("-----------------------------------------------------")
}
