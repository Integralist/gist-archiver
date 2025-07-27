# Go Backend (`main.go`)

The Go application is the core of the Gist Archiver. It's a command-line tool responsible for fetching all of a user's public Gists from the GitHub API, processing them, and generating static HTML and Markdown files.

## Key Responsibilities

1.  **Fetch Gists**: It retrieves a list of all public Gists for a specified GitHub user. It handles pagination to ensure all Gists are fetched.
2.  **Process Gists**: For each Gist, it fetches detailed information, including file contents.
3.  **Extract Tags**: It parses Gist descriptions to find hashtags (e.g., `#go`, `#javascript`) and associates them with the Gist.
4.  **Generate Markdown**: It creates a single composite Markdown file for each Gist, combining all of its files into one page. It also adds metadata like the Gist title, creation date, and a link back to the original Gist on GitHub.
5.  **Generate Index**: It creates a main `index.md` and `index.html` file that lists all the archived Gists, linking to their respective generated pages.

## How It Works

The application flow is managed in `main.go`:

- It uses goroutines and channels for concurrent fetching of Gist details to speed up the process (`fetchAndProcessGistDetails`). A semaphore (`fetchSemaphore`) is used to limit the number of concurrent API requests.
- The `GistDetail` struct holds the comprehensive data for a single Gist, including its files and extracted tags.
- `generateGistMarkdownFile` is responsible for creating the content of the Markdown file for a single Gist. It includes the Gist's title, content from all its files, and syntax highlighting hints for code blocks.
- `processSingleGist` orchestrates the processing of one Gist: extracting tags, generating the Markdown file, and sending an `IndexEntry` to a channel for later aggregation.
- Finally, `generateIndexFiles` takes all the `IndexEntry` items and builds the main index files.

## Project Structure Conventions

The code follows standard Go project layout conventions, although most of the logic is currently in `main.go`. For more complex versions, this logic would be broken out into `internal` packages as described in `CONVENTIONS.md`.
# Go Backend Architecture

This document details the Go application that fetches and processes Gists.

---

## 1. Overview

The Go backend is the core of the Gist Archiver. It is a command-line application responsible for:
- Fetching all public gists for a specific GitHub user.
- Processing each gist's content into both Markdown and HTML formats.
- Generating a static HTML website for browsing the archived content.

The application is built with concurrency in mind to efficiently handle multiple API requests and file processing tasks simultaneously.

---

## 2. Project Structure

The Go-related parts of the project follow standard layout conventions:

```
/cmd
  /archiver         # Main application entrypoint
/internal
  /config           # Application configuration
  /filegen          # Markdown and HTML file generation
  /github           # GitHub API client
  /model            # Data models (structs)
go.mod
go.sum
```

-   **/cmd/archiver**: Contains the `main` package, which is the entry point of the application. Its role is to initialize and coordinate the different components.
-   **/internal**: Contains the core application logic, separated into packages with distinct responsibilities. This code is not intended to be imported by other projects.

---

## 3. Package Details

### `internal/config`

This package centralizes all application-wide constants, such as API endpoints, directory paths, and concurrency settings. This makes configuration easy to manage from a single location.

### `internal/model`

This package defines all the primary data structures (structs) used throughout the application. These models represent entities like `GistDetail`, `GistFile`, and `IndexEntry`, which are used for processing data from the GitHub API and generating the final archive.

### `internal/github`

This package provides a client for interacting with the GitHub Gists API. Its responsibilities include:
-   Making authenticated API requests using a personal access token.
-   Fetching the list of all public gists for a user, handling API pagination.
-   Fetching the detailed content for each individual gist.
-   Managing concurrent API requests to improve performance.

### `internal/filegen`

This package is responsible for generating all static files. Its key functions are:
-   **Markdown Generation**: Creates a composite Markdown file for each gist, combining all of its file contents with metadata.
-   **HTML Conversion**: Converts the generated Markdown into static HTML pages, including syntax highlighting for code blocks via `gomarkdown`.
-   **Index File Creation**: Generates `index.md` and `index.html` to list all archived gists, sorted by date.

---

## 4. Execution Flow

The application's execution is orchestrated in `cmd/archiver/main.go` and relies heavily on goroutines and channels for concurrency.

1.  **Initialization**: The `main` function starts by creating the necessary output directories (`/gists` and `/markdown_files`).
2.  **Channels**: Two buffered channels are created to manage the flow of data between concurrent stages:
    -   `gistDetailChan`: To pass `GistDetail` objects from the API fetcher to the file processors.
    -   `indexEntryChan`: To pass `IndexEntry` objects from the processors to the final index generator.
3.  **Goroutines**: The application launches several goroutines to work in parallel:
    -   A single **fetcher goroutine** (`client.FetchAndProcessGistDetails`) fetches gist details from the GitHub API and sends them into `gistDetailChan`. Once all gists are fetched, it closes this channel.
    -   A pool of **processor goroutines** (`filegen.ProcessSingleGist`) reads from `gistDetailChan`. Each worker processes one gist at a time (generating its Markdown and HTML pages) and sends an `IndexEntry` into `indexEntryChan`. The number of workers is based on the number of available CPU cores.
    -   A single **indexer goroutine** collects all `IndexEntry` items from `indexEntryChan`. Once the channel is closed (signaling that all gists have been processed), it sorts the entries by date and generates the final `index.html` and `index.md` files.
4.  **Synchronization**: `sync.WaitGroup` is used to ensure that all steps complete in the correct order. The main function waits for fetching to complete, then for processing to complete, and finally for indexing to complete before exiting.

This concurrent model allows the application to overlap network I/O (fetching gists) with CPU-bound work (processing files), significantly speeding up the archival process.
