document.addEventListener('DOMContentLoaded', () => {
  // --- Configuration ---
  const gistEntrySelector = '#gistList li'; // Selector for individual gist list items
  const tagContainerSelector = 'small.date'; // Selector for the element containing "Tags: #tag1 #tag2"
  const dynamicTagLinkClass = 'dynamic-tag-link'; // Class for dynamically created tag links
  const filterInputId = 'filterInput'; // ID of the existing search input field

  // --- UI Elements for Active Filter Display and Clear Link ---
  let activeFilterDisplay = null;
  let clearFilterLink = null;

  function createFilterUIElements() {
    const filterInputElement = document.getElementById(filterInputId);
    if (!filterInputElement) {
      // console.warn('Gist Tag Filter: Filter input field not found. Cannot create UI elements for active filter.');
      return;
    }

    // Create span to display active filter
    activeFilterDisplay = document.createElement('span');
    activeFilterDisplay.id = 'activeTagFilterDisplay';
    activeFilterDisplay.style.display = 'none'; // Hidden by default
    activeFilterDisplay.style.marginLeft = '10px';
    activeFilterDisplay.style.fontSize = '0.9em';

    // Create "Clear filter" link
    clearFilterLink = document.createElement('a');
    clearFilterLink.id = 'clearTagFilterLink';
    clearFilterLink.href = '#';
    clearFilterLink.textContent = 'Clear active tag filter';
    clearFilterLink.style.display = 'none'; // Hidden by default
    clearFilterLink.style.marginLeft = '10px';
    clearFilterLink.style.fontSize = '0.9em';
    clearFilterLink.style.textDecoration = 'underline'; // Make it look more like a link

    // Insert after the filter input field
    filterInputElement.parentNode.insertBefore(activeFilterDisplay, filterInputElement.nextSibling);
    // Insert the clear link after the active filter display
    activeFilterDisplay.parentNode.insertBefore(clearFilterLink, activeFilterDisplay.nextSibling);

    // Add event listener for the clear link
    clearFilterLink.addEventListener('click', (event) => {
      event.preventDefault();
      performLocalFilter(''); // Pass empty string to clear filter
      if (filterInputElement) {
        filterInputElement.value = ''; // Also clear the text input field
      }
    });
  }

  // --- Helper: Convert plain text tags to clickable links ---
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

  // --- Helper Function to Prepare Gist Data with Parsed Tags ---
  let localGistData = [];

  function prepareLocalGistData() {
    const gistListItems = document.querySelectorAll(gistEntrySelector);
    if (!gistListItems.length) {
      return false;
    }

    localGistData = Array.from(gistListItems).map(listItem => {
      const itemTags = [];
      const smallElements = listItem.querySelectorAll('small.date');
      smallElements.forEach(small => {
        const smallText = small.textContent || "";
        const tagsPrefix = "tags:";
        const tagsIndex = smallText.toLowerCase().indexOf(tagsPrefix);
        if (tagsIndex !== -1) {
          const tagsString = smallText.substring(tagsIndex + tagsPrefix.length);
          const potentialTags = tagsString.trim().split(/\s+/);
          potentialTags.forEach(pt => {
            if (pt.startsWith('#') && pt.length > 1) {
              itemTags.push(pt.substring(1).toLowerCase());
            }
          });
        }
      });
      return {
        element: listItem,
        tagsArray: itemTags,
      };
    });
    return true;
  }

  // --- Helper Function to Perform Filtering Based on Exact Tag Match ---
  function performLocalFilter(clickedTagName) {
    if (!localGistData.length) {
      if (!prepareLocalGistData()) {
        console.warn('Gist Tag Filter: In performLocalFilter - No gist entries found or data could not be prepared.');
        return;
      }
    }

    const lowerCaseQuery = clickedTagName.toLowerCase().trim();

    if (!lowerCaseQuery) { // Clearing the filter
      localGistData.forEach(item => item.element.style.display = 'list-item');
      if (activeFilterDisplay) {
        activeFilterDisplay.innerHTML = ''; // Clear innerHTML
        activeFilterDisplay.style.display = 'none';
      }
      if (clearFilterLink) {
        clearFilterLink.style.display = 'none';
      }
      return;
    }

    // Applying a filter
    localGistData.forEach(item => {
      const hasExactTag = item.tagsArray.includes(lowerCaseQuery);
      item.element.style.display = hasExactTag ? 'list-item' : 'none';
    });

    if (activeFilterDisplay) {
      // Use innerHTML to allow for bold tag
      activeFilterDisplay.innerHTML = `<strong>Filtering by:</strong> #${lowerCaseQuery}`;
      activeFilterDisplay.style.display = 'inline'; // Show display
    }
    if (clearFilterLink) {
      clearFilterLink.style.display = 'inline'; // Show clear link
    }
  }

  // --- Main Logic ---
  // 0. Create UI elements for filter display and clear link
  createFilterUIElements();

  // 1. Convert plain text tags to links
  convertPlainTextTagsToDynamicLinks();

  // 2. Select all newly created dynamic tag links
  const tagLinks = document.querySelectorAll(`.${dynamicTagLinkClass}`);

  if (!tagLinks.length) {
    // console.log(`Gist Tag Filter: No dynamic tag links found.`);
    return;
  }

  // 3. Prepare the localGistData array.
  if (!prepareLocalGistData()) {
    console.warn('Gist Tag Filter: Initial preparation of local Gist data failed.');
  }

  // 4. Add click event listeners to each dynamic tag link
  tagLinks.forEach(link => {
    link.addEventListener('click', event => {
      event.preventDefault();
      const tagName = link.dataset.tag;
      if (!tagName) return;

      performLocalFilter(tagName);

      const firstGistLi = document.querySelector(gistEntrySelector);
      if (firstGistLi) {
        const gistListElement = document.getElementById('gistList');
        if (gistListElement) {
          gistListElement.scrollIntoView({ behavior: 'smooth', block: 'start' });
        } else {
          firstGistLi.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
      }
    });
  });

  // console.log('Gist Tag Filter: Script loaded with clear filter functionality and bold display.');
});
