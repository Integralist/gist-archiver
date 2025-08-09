document.addEventListener('DOMContentLoaded', () => {
  const alertTypes = {
    '[!NOTE]': 'note',
    '[!TIP]': 'tip',
    '[!IMPORTANT]': 'important',
    '[!WARNING]': 'warning',
    '[!CAUTION]': 'caution',
  };

  /**
   * Checks if a paragraph's content starts with a known alert keyword.
   * @param {string} html - The innerHTML of the paragraph element.
   * @returns {{type: string, keyword: string}|null} The alert type and keyword, or null if no match.
   */
  function findAlertType(html) {
    for (const [keyword, type] of Object.entries(alertTypes)) {
      if (html.startsWith(keyword)) {
        return { type, keyword };
      }
    }
    return null;
  }

  const contentContainer = document.querySelector('.markdown-content');
  if (!contentContainer) {
    return;
  }

  const blockquotes = contentContainer.querySelectorAll('blockquote');

  blockquotes.forEach(quote => {
    // If the blockquote is already styled (e.g. by a previous run), skip it.
    if (quote.classList.contains('markdown-alert')) {
      return;
    }

    // A single blockquote from the original markdown can contain multiple "alert" sections.
    // This script iterates through the blockquote's children and moves them into
    // new, properly styled blockquote elements when an alert keyword is found.
    const children = Array.from(quote.children);
    let currentAlertBlock = null;

    children.forEach(child => {
      // An alert can only be triggered by a paragraph. If we encounter other elements
      // (like lists or nested blockquotes) while inside an alert, append them.
      if (child.tagName !== 'P') {
        if (currentAlertBlock) {
          currentAlertBlock.appendChild(child);
        }
        return;
      }

      const p_html = child.innerHTML.trim();
      const alertInfo = findAlertType(p_html);

      if (alertInfo) {
        // This paragraph starts a new alert. Create a new, styled blockquote for it.
        // This new block will become the active container for this and subsequent elements.
        currentAlertBlock = document.createElement('blockquote');
        currentAlertBlock.classList.add('markdown-alert', `markdown-alert-${alertInfo.type}`);
        
        // Insert the new alert block before the original, unstyled blockquote.
        quote.parentNode.insertBefore(currentAlertBlock, quote);
        
        // Clean the alert keyword from the paragraph's content.
        child.innerHTML = p_html.substring(alertInfo.keyword.length).trim();
      }

      if (currentAlertBlock) {
        // If we are in an active alert block, move the current child element into it.
        // This moves elements from the original blockquote to the new styled one.
        currentAlertBlock.appendChild(child);
      }
    });

    // If the original blockquote is now empty, it means all its content
    // was moved into new alert blocks. Remove the empty original.
    if (quote.children.length === 0) {
      quote.parentNode.removeChild(quote);
    }
  });
});
