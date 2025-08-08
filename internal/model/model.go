package model

import "time"

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
	HTMLURL     string              `json:"html_url"`
	Public      bool                `json:"public"`
	Owner       GistOwner           `json:"owner"`
	Files       map[string]GistFile `json:"files"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
	Tags        string              // custom for how I typically add #tag1 #tag2 to my gist description
}

// GetTitle extracts the title from a gist's description or filename.
func (g *GistDetail) GetTitle() string {
	if g.Description != "" {
		return g.Description
	}
	for _, file := range g.Files {
		return file.Filename
	}
	return g.ID
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
	GistID       string
	Title        string
	MarkdownPath string
	HTMLPath     string
	Date         time.Time
	Tags         string
}

type SearchIndexEntry struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
}
