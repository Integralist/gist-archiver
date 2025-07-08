document.addEventListener('DOMContentLoaded', () => {
  const alertTypes = {
    '[!NOTE]': 'note',
    '[!TIP]': 'tip',
    '[!IMPORTANT]': 'important',
    '[!WARNING]': 'warning',
    '[!CAUTION]': 'caution',
  };

  // The CSS is scoped with .markdown-content, so we search within that container.
  const contentContainer = document.querySelector('.markdown-content');
  if (!contentContainer) {
    return; // Exit if no markdown container is found on the page.
  }

  const blockquotes = contentContainer.querySelectorAll('blockquote');

  blockquotes.forEach(quote => {
    const firstParagraph = quote.querySelector('p:first-child');
    if (!firstParagraph) return;

    const paragraphHtml = firstParagraph.innerHTML.trim();

    for (const [key, value] of Object.entries(alertTypes)) {
      if (paragraphHtml.startsWith(key)) {
        quote.classList.add('markdown-alert', `markdown-alert-${value}`);
        
        // Remove the alert keyword (e.g., "[!TIP]") from the paragraph.
        firstParagraph.innerHTML = paragraphHtml.substring(key.length).trim();
        
        break; // Move to the next blockquote.
      }
    }
  });
});
