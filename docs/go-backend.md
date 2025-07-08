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
