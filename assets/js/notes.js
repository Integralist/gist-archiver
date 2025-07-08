document.addEventListener('DOMContentLoaded', () => {
  const alertTypes = {
    '[!NOTE]': 'note',
    '[!TIP]': 'tip',
    '[!IMPORTANT]': 'important',
    '[!WARNING]': 'warning',
    '[!CAUTION]': 'caution',
  };

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

    const children = Array.from(quote.children);
    let currentAlertBlock = null;

    children.forEach(child => {
      // We only check for alert keywords in paragraphs.
      if (child.tagName !== 'P') {
        if (currentAlertBlock) {
          currentAlertBlock.appendChild(child);
        }
        return;
      }

      const p_html = child.innerHTML.trim();
      let foundAlertType = null;
      let keyword = null;

      for (const [key, value] of Object.entries(alertTypes)) {
        if (p_html.startsWith(key)) {
          foundAlertType = value;
          keyword = key;
          break;
        }
      }

      if (foundAlertType) {
        // This paragraph starts a new alert. Create a new blockquote for it.
        currentAlertBlock = document.createElement('blockquote');
        currentAlertBlock.classList.add('markdown-alert', `markdown-alert-${foundAlertType}`);
        
        // Insert the new, styled blockquote before the original one.
        quote.parentNode.insertBefore(currentAlertBlock, quote);
        
        // Clean the alert keyword from the paragraph's content.
        child.innerHTML = p_html.substring(keyword.length).trim();
      }

      if (currentAlertBlock) {
        // If we are processing an alert, move the current child element into it.
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
