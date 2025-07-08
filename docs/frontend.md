# Frontend

The frontend consists of HTML, CSS, and JavaScript files that create the user interface for browsing the archived Gists.

## Key Components

-   **`index.html`**: The main entry point that displays the list of all Gists.
-   **`assets/js/tags.js`**: A script that provides dynamic, client-side filtering of Gists by tags.
-   **`assets/js/highlight.min.js`**: A third-party library (highlight.js) used for syntax highlighting of code snippets within the Gist pages.
-   **`assets/css/style.css`**: Custom styles for the website.

## Tag Filtering (`tags.js`)

The `tags.js` script enhances the user experience by allowing users to filter the list of Gists without reloading the page.

### How It Works

1.  **Initialization**: When the page loads, the script scans the page for Gist entries and their associated tags.
2.  **UI Creation**: It dynamically creates a filter input field and a display area for the active filter.
3.  **Tag Extraction**: It finds tags within the Gist list items. The Go backend ensures tags from the Gist description (e.g., `#go #logs`) are present in the HTML. The script makes these tags clickable.
4.  **Local Filtering**: When a user clicks a tag or types in the filter box, the script performs a local search through the Gist data it has already collected from the DOM.
5.  **DOM Manipulation**: It shows or hides Gist entries in the list based on whether they match the filter query. This is done entirely on the client-side, making the filtering feel instantaneous.

This approach provides a fast and responsive way to navigate Gists, especially when the number of archived Gists is large.
