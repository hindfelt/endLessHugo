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
            
            // Add new posts - exclude hero post
            const newPosts = doc.querySelectorAll('.post-card:not(.hero)');
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
    
        try {
            // Fetch the index to get the next post
            const response = await fetch('/index.json');
            const data = await response.json();
            
            // Find current post index
            const currentIndex = data.posts.findIndex(post => 
                post.permalink.endsWith(currentUrl) || currentUrl.includes(post.permalink)
            );
            
            console.log('Current index:', currentIndex); // Debug
    
            if (currentIndex === -1 || currentIndex === data.posts.length - 1) {
                console.log('No next post available');
                return;
            }
    
            // Get next post
            const nextPost = data.posts[currentIndex + 1];
            if (!nextPost?.permalink) {
                console.log('Next post URL not found');
                return;
            }
    
            console.log('Loading next post:', nextPost.permalink); // Debug
    
            // Fetch and display next post
            const postResponse = await fetch(nextPost.permalink);
            const html = await postResponse.text();
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');
            
            const nextPostContent = doc.querySelector('.single-post');
            if (nextPostContent) {
                // Create wrapper for the new post
                const postWrapper = document.createElement('article');
                postWrapper.className = 'single-post';
                postWrapper.innerHTML = nextPostContent.innerHTML;
                
                // Add separator
                const separator = document.createElement('div');
                separator.className = 'post-separator';
                container.appendChild(separator);
                
                // Add the new post
                container.appendChild(postWrapper);
                
                // Update URL without page reload
                window.history.pushState({}, '', nextPost.permalink);
                
                // Update document title
                document.title = nextPost.title;
                
                // Observe the new last article
                this.observer.observe(postWrapper);
            }
    
        } catch (error) {
            console.error('Error loading next post:', error);
        } finally {
            this.loading = false;
        }
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