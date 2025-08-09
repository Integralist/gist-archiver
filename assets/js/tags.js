document.addEventListener('DOMContentLoaded', () => {
  const tagContainerSelector = 'small.date';
  const dynamicTagLinkClass = 'dynamic-tag-link';
  const filterInputId = 'filterInput';
  const gistListSelector = '#gistList li';

  /**
   * Finds plain text tags (e.g., #go) within designated containers
   * and converts them into clickable anchor tags for filtering.
   */
  function convertPlainTextTagsToDynamicLinks() {
    const tagContainers = document.querySelectorAll(tagContainerSelector);

    tagContainers.forEach(container => {
      if (container.classList.contains('tags-processed')) {
        return;
      }

      const boldTagElement = container.querySelector('b');
      if (boldTagElement && boldTagElement.textContent.trim().toLowerCase().startsWith('tags:')) {
        let htmlContent = container.innerHTML;
        htmlContent = htmlContent.replace(/(#\w+)/g, (match) => {
          const tagNameNoHash = match.substring(1);
          return `<a href="#" class="${dynamicTagLinkClass}" data-tag="${tagNameNoHash}">${match}</a>`;
        });
        container.innerHTML = htmlContent;
        container.classList.add('tags-processed');
      }
    });
  }

  /**
   * Programmatically sets the value of the search input and triggers
   * the 'input' event to activate the search functionality in search.js.
   * @param {string} query - The search query to trigger.
   */
  function triggerSearch(query) {
    const filterInputElement = document.getElementById(filterInputId);
    if (filterInputElement) {
      filterInputElement.value = query;
      filterInputElement.dispatchEvent(new Event('input', { bubbles: true }));
    }
  }

  // Main Logic
  convertPlainTextTagsToDynamicLinks();

  const tagLinks = document.querySelectorAll(`.${dynamicTagLinkClass}`);
  if (!tagLinks.length) {
    return;
  }

  tagLinks.forEach(link => {
    link.addEventListener('click', event => {
      event.preventDefault();
      const tagName = link.dataset.tag;
      if (!tagName) return;

      triggerSearch(`tags:${tagName}`);

      // After filtering, scroll the top of the gist list into view.
      const gistListElement = document.getElementById('gistList');
      if (gistListElement) {
        gistListElement.scrollIntoView({ behavior: 'smooth', block: 'start' });
      }
    });
  });
});
