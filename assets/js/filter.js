document.addEventListener('DOMContentLoaded', function() {
    const filterInput = document.getElementById('filterInput');
    const gistList = document.getElementById('gistList');

    // Check if essential elements exist to prevent errors
    if (!filterInput || !gistList) {
        console.error("Filter input or gist list not found on the page.");
        return;
    }

    const listItems = gistList.getElementsByTagName('li');

    filterInput.addEventListener('keydown', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault();
            const filterText = filterInput.value.toLowerCase().trim();

            for (let i = 0; i < listItems.length; i++) {
                const item = listItems[i];
                const itemText = item.textContent || item.innerText || '';

                if (filterText === '') {
                    item.classList.remove('hidden');
                } else {
                    if (itemText.toLowerCase().includes(filterText)) {
                        item.classList.remove('hidden');
                    } else {
                        item.classList.add('hidden');
                    }
                }
            }
        }
    });
});