# File Generation Package

This package is responsible for generating all static files, including Markdown and HTML for gists and the main index pages.

## Responsibilities

-   **Markdown Generation**: Creates composite Markdown files for each gist, combining file contents with metadata.
-   **HTML Conversion**: Converts the generated Markdown files into static HTML pages.
-   **Index File Creation**: Generates `index.md` and `index.html` to list all archived gists.
-   **Utilities**: Provides helper functions, such as identifying whether a file is a Markdown file.
