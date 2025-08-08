document.addEventListener('DOMContentLoaded', function () {
  const searchInput = document.getElementById('filterInput');
  const gistListContainer = document.getElementById('gistList');

  if (!searchInput || !gistListContainer) {
    console.warn('Gist Search: Search input or gist list container not found. Search disabled.');
    return;
  }

  if (typeof lunr === 'undefined') {
    console.error('Gist Search: Lunr.js is not loaded. Search disabled.');
    searchInput.disabled = true;
    searchInput.placeholder = 'Search not available';
    return;
  }

  if (!window.SEARCH_DATA) {
    console.error('Gist Search: Search data not found. Search disabled.');
    searchInput.disabled = true;
    searchInput.placeholder = 'Search not available';
    return;
  }

  const gists = Array.from(gistListContainer.getElementsByClassName('gist-entry'));
  let lunrIndex;

  function showAllGists() {
    gists.forEach(gist => {
      gist.style.display = ''; // Use empty string to revert to stylesheet's display property
    });
  }

  // Build the Lunr index immediately with the embedded data
  try {
    lunrIndex = lunr(function () {
      this.ref('id');
      this.field('title', { boost: 10 });
      this.field('tags', { boost: 5 });
      this.field('content');

      window.SEARCH_DATA.forEach(item => {
        this.add(item);
      });
    });
  } catch (e) {
    console.error('Gist Search: Error building Lunr index:', e);
    searchInput.disabled = true;
    searchInput.placeholder = 'Search not available';
    return;
  }

  searchInput.addEventListener('input', function (event) {
    const query = event.target.value.trim();

    if (!query) {
      showAllGists();
      return;
    }

    try {
      // Perform a search that supports both partial matching (wildcard) and stemming.
      const searchResults = lunrIndex.search(query + '* ' + query);
      const resultIds = new Set(searchResults.map(r => `gist-${r.ref}`));

      gists.forEach(gist => {
        if (resultIds.has(gist.id)) {
          gist.style.display = '';
        } else {
          gist.style.display = 'none';
        }
      });
    } catch (e) {
      // Lunr can throw errors for invalid query syntax (e.g., a trailing colon)
      if (e instanceof lunr.QueryParseError) {
        // This is expected for certain query inputs, do nothing.
      } else {
        console.warn('Gist Search: Error during search:', e);
      }
    }
  });
});
