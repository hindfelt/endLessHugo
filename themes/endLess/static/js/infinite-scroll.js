class InfiniteScroll {
    constructor() {
        this.loading = false;
        this.observer = null;
        this.currentPage = 1;
        this.isListView = document.getElementById('post-grid') !== null;
        
        this.init();
    }

    init() {
        // Initialize based on view type
        if (this.isListView) {
            this.initListView();
        } else {
            this.initSingleView();
        }
    }

    initListView() {
        const loadMoreDiv = document.getElementById('page-loader');
        if (!loadMoreDiv) return;

        this.observer = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting) {
                this.loadMorePosts();
            }
        }, {
            rootMargin: '300px'
        });

        this.observer.observe(loadMoreDiv);
    }

    initSingleView() {
        const articles = document.querySelectorAll('article');
        if (articles.length === 0) return;

        this.observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    this.loadNextPost();
                }
            });
        }, {
            rootMargin: '200px',
            threshold: 0.1
        });

        // Observe the last article
        this.observer.observe(articles[articles.length - 1]);
    }

    async loadMorePosts() {
        if (this.loading) return;
        
        this.loading = true;
        const loadMoreDiv = document.getElementById('page-loader');
        const gridContainer = document.getElementById('post-grid');
        const nextPage = loadMoreDiv.dataset.nextPage;

        try {
            const response = await fetch(nextPage);
            const html = await response.text();
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');
            
            // Add new posts
            const newPosts = doc.querySelectorAll('.post-card');
            newPosts.forEach(post => {
                gridContainer.appendChild(post.cloneNode(true));
            });
            
            // Update loader with new next page
            const newLoader = doc.getElementById('page-loader');
            if (newLoader) {
                loadMoreDiv.dataset.nextPage = newLoader.dataset.nextPage;
                this.currentPage++;
            } else {
                loadMoreDiv.remove();
                this.observer.disconnect();
            }

        } catch (error) {
            console.error('Error loading more posts:', error);
        } finally {
            this.loading = false;
        }
    }

    async loadNextPost() {
        if (this.loading) return;
        
        this.loading = true;
        const container = document.querySelector('.next-posts-container');
        const currentUrl = window.location.pathname;
        const nextPostUrl = this.getNextPostUrl(currentUrl);

        try {
            const response = await fetch(nextPostUrl);
            const html = await response.text();
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');
            
            // Get the next post content
            const nextPost = doc.querySelector('.single-post');
            if (nextPost) {
                // Create wrapper for the new post
                const postWrapper = document.createElement('article');
                postWrapper.className = 'single-post';
                postWrapper.innerHTML = nextPost.innerHTML;
                
                // Add separator
                const separator = document.createElement('div');
                separator.className = 'post-separator';
                container.appendChild(separator);
                
                // Add the new post
                container.appendChild(postWrapper);
                
                // Update URL without page reload
                const newUrl = doc.querySelector('link[rel="canonical"]')?.href || nextPostUrl;
                window.history.pushState({}, '', newUrl);
                
                // Update document title
                document.title = doc.title;
                
                // Observe the new last article
                this.observer.observe(postWrapper);
            }

        } catch (error) {
            console.error('Error loading next post:', error);
        } finally {
            this.loading = false;
        }
    }

    getNextPostUrl(currentUrl) {
        // Remove trailing slash if present
        currentUrl = currentUrl.replace(/\/$/, '');
        
        // Split the URL into parts
        const parts = currentUrl.split('/');
        
        // If we're in a numbered post, increment the number
        const lastPart = parts[parts.length - 1];
        if (/^\d+$/.test(lastPart)) {
            const nextNum = parseInt(lastPart) + 1;
            parts[parts.length - 1] = nextNum.toString();
        } else {
            // If not numbered, append /2 (or increment existing number)
            const match = currentUrl.match(/\/(\d+)$/);
            if (match) {
                const nextNum = parseInt(match[1]) + 1;
                return currentUrl.replace(/\/\d+$/, `/${nextNum}`);
            } else {
                return `${currentUrl}/2`;
            }
        }
        
        return parts.join('/');
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new InfiniteScroll();
});

// Add some CSS for the separator
const style = document.createElement('style');
style.textContent = `
    .post-separator {
        height: 1px;
        background: #eee;
        margin: 4rem 0;
        width: 100%;
    }
`;
document.head.appendChild(style);