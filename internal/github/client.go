package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/integralist/gist-archiver/internal/model"
)

var (
	httpClient = &http.Client{Timeout: 30 * time.Second}
	linkRegex  = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)
)

// Client handles communication with the GitHub API.
type Client struct {
	BaseURL              string
	Username             string
	Token                string
	GistsPerPage         int
	MaxConcurrentFetches int
}

// NewClient creates a new GitHub API client.
func NewClient(baseURL, username, token string, gistsPerPage, maxConcurrentFetches int) *Client {
	if token != "" {
		log.Printf("🎉 Using GitHub API Token")
	} else {
		log.Printf("⚠️ No GitHub API Token found")
	}
	return &Client{
		BaseURL:              baseURL,
		Username:             username,
		Token:                token,
		GistsPerPage:         gistsPerPage,
		MaxConcurrentFetches: maxConcurrentFetches,
	}
}

// doAPIRequest is a helper that handles the common logic for making a GitHub API request.
func (c *Client) doAPIRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "go-gist-archiver/1.6")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to %s: %w", url, err)
	}

	if resp.StatusCode == http.StatusForbidden && strings.Contains(resp.Header.Get("X-RateLimit-Remaining"), "0") {
		log.Printf("Rate limit hit for URL %s. Status: %s", url, resp.Status)
		return nil, fmt.Errorf("rate limit hit for %s. X-RateLimit-Remaining: %s, X-RateLimit-Reset: %s", url, resp.Header.Get("X-RateLimit-Remaining"), resp.Header.Get("X-RateLimit-Reset"))
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("api request to %s failed with status %s: %s", url, resp.Status, string(bodyBytes))
	}

	return resp, nil
}

func (c *Client) makeAPIRequest(url string, target any) error {
	resp, err := c.doAPIRequest(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode JSON response from %s: %w", url, err)
		}
	}
	return nil
}

func (c *Client) makeAPIRequestForGistList(url string, target any) (string, error) {
	resp, err := c.doAPIRequest(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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

// FetchAndProcessGistDetails fetches all public gists for a user and sends them to a channel.
func (c *Client) FetchAndProcessGistDetails(gistDetailChan chan<- model.GistDetail, fetchWg *sync.WaitGroup) {
	defer close(gistDetailChan)

	nextURL := fmt.Sprintf("%s/users/%s/gists?per_page=%d", c.BaseURL, c.Username, c.GistsPerPage)
	page := 1
	fetchSemaphore := make(chan struct{}, c.MaxConcurrentFetches)

	for nextURL != "" {
		log.Printf("Fetching page %d of gist list for user %s from %s", page, c.Username, nextURL)
		var gistListEntries []model.GistListEntry
		var err error
		nextURL, err = c.makeAPIRequestForGistList(nextURL, &gistListEntries)
		if err != nil {
			log.Printf("Error: Failed to fetch gist list page %d: %v. Stopping further list fetching.", page, err)
			return
		}

		for _, entry := range gistListEntries {
			if !entry.Public {
				continue
			}

			fetchWg.Add(1)
			go func(gistEntry model.GistListEntry) {
				defer fetchWg.Done()
				fetchSemaphore <- struct{}{}
				defer func() { <-fetchSemaphore }()

				tempGistDetail := model.GistDetail{Description: gistEntry.Description, Files: convertListFilesToDetailFiles(gistEntry.Files), ID: gistEntry.ID}
				log.Printf("Fetching details for gist ID: %s (%s)", gistEntry.ID, tempGistDetail.GetTitle())

				var gistDetail model.GistDetail
				detailURL := fmt.Sprintf("%s/gists/%s", c.BaseURL, gistEntry.ID)
				if err := c.makeAPIRequest(detailURL, &gistDetail); err != nil {
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

func convertListFilesToDetailFiles(listFiles map[string]model.GistListFile) map[string]model.GistFile {
	detailFiles := make(map[string]model.GistFile)
	for k, v := range listFiles {
		detailFiles[k] = model.GistFile{Filename: v.Filename}
	}
	return detailFiles
}
