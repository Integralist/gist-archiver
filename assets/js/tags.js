document.addEventListener('DOMContentLoaded', () => {
  const tagContainerSelector = 'small.date';
  const dynamicTagLinkClass = 'dynamic-tag-link';
  const filterInputId = 'filterInput';
  const gistListSelector = '#gistList li';

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

      const firstGistLi = document.querySelector(gistListSelector);
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
});
