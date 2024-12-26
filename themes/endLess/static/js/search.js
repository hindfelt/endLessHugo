document.addEventListener('DOMContentLoaded', function() {
    const searchModal = document.getElementById('search-modal');
    const searchInput = document.getElementById('search-input');
    const searchResults = document.getElementById('search-results');
    const closeSearch = document.getElementById('close-search');
    const searchTrigger = document.querySelector('.search-trigger');
    let posts = [];

    // Debug check for elements
    console.log({
        searchModal: !!searchModal,
        searchInput: !!searchInput,
        searchResults: !!searchResults,
        closeSearch: !!closeSearch
    });

    // Fetch posts data
    fetch('/index.json')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log('Data received:', data); // Debug log
            if (data && Array.isArray(data)) {
                posts = data;
            } else if (data && Array.isArray(data.posts)) {
                posts = data.posts;
            } else {
                console.error('Unexpected data structure:', data);
                posts = [];
            }
        })
        .catch(error => console.error('Error loading posts:', error));

    // Open search with Ctrl+K or Cmd+K
    document.addEventListener('keydown', function(e) {
        if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
            e.preventDefault();
            openSearch();
        }
        if (e.key === 'Escape' && searchModal.classList.contains('active')) {
            closeSearchModal();
        }
        if (searchTrigger) {
            searchTrigger.addEventListener('click', openSearch);
        }
    });

    // Close search modal
    if (closeSearch) {
        closeSearch.addEventListener('click', closeSearchModal);
    }
    
    if (searchModal) {
        searchModal.addEventListener('click', function(e) {
            if (e.target === searchModal) {
                closeSearchModal();
            }
        });
    }

    // Handle search input
    if (searchInput) {
        searchInput.addEventListener('input', function() {
            const query = this.value.toLowerCase();
            
            if (!Array.isArray(posts)) {
                console.error('Posts is not an array:', posts);
                return;
            }
            
            if (query.length < 2) {
                searchResults.innerHTML = '<div class="search-result-item">Type at least 2 characters to search...</div>';
                return;
            }

            const results = posts.filter(post => 
                (post.title && post.title.toLowerCase().includes(query)) || 
                (post.content && post.content.toLowerCase().includes(query))
            );

            displayResults(results);
        });
    }

    function displayResults(results) {
        if (!searchResults) return;
        
        if (results.length === 0) {
            searchResults.innerHTML = '<div class="search-result-item">No results found</div>';
            return;
        }

        searchResults.innerHTML = results
            .slice(0, 10)
            .map(post => `
                <div class="search-result-item" onclick="window.location.href='${post.permalink}'">
                    <h3>${post.title}</h3>
                    <time>${post.date}</time>
                </div>
            `)
            .join('');
    }

    function openSearch() {
        if (searchModal) {
            searchModal.classList.add('active');
            if (searchInput) {
                searchInput.focus();
            }
        }
    }

    function closeSearchModal() {
        if (searchModal) {
            searchModal.classList.remove('active');
            if (searchInput) {
                searchInput.value = '';
            }
            if (searchResults) {
                searchResults.innerHTML = '';
            }
        }
    }
});
