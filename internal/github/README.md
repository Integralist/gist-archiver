# GitHub Client Package

This package provides a client for interacting with the GitHub Gists API.

## Responsibilities

-   **API Requests**: Handles making authenticated or unauthenticated requests to the GitHub API.
-   **Gist Fetching**: Fetches the list of all public gists for a specified user.
-   **Detail Fetching**: Fetches the detailed content and metadata for each individual gist.
-   **Concurrency**: Manages concurrent API requests to improve performance while respecting rate limits.
