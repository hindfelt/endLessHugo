document.addEventListener('DOMContentLoaded', function() {
    // Prevent default image click behavior
    document.querySelectorAll('.post-content img').forEach(img => {
        img.addEventListener('click', function(e) {
            e.preventDefault();
            showImageOverlay(this.src);
        });
        
        // Prevent drag and right-click default behaviors
        img.addEventListener('dragstart', e => e.preventDefault());
        img.addEventListener('contextmenu', e => e.preventDefault());
        
        // Make images look clickable
        img.style.cursor = 'zoom-in';
    });

    // Create and append overlay
    const overlay = document.createElement('div');
    overlay.className = 'image-overlay';
    overlay.innerHTML = `
        <img class="fullsize-image">
        <div class="image-overlay-instructions">Click image to open full size in new tab • ESC to close</div>
    `;
    document.body.appendChild(overlay);

    const fullsizeImage = overlay.querySelector('.fullsize-image');

    // Add click handler to open full size in new tab
    fullsizeImage.addEventListener('click', function(e) {
        e.stopPropagation(); // Prevent overlay from closing
        window.open(this.src, '_blank');
    });

    function showImageOverlay(src) {
        fullsizeImage.src = src;
        overlay.style.display = 'flex';
        document.body.style.overflow = 'hidden';
        fullsizeImage.style.cursor = 'zoom-in'; // Show it's clickable
    }

    function hideOverlay() {
        overlay.style.display = 'none';
        document.body.style.overflow = '';
    }

    // Close on overlay click
    overlay.addEventListener('click', function(e) {
        if (e.target === overlay) {
            hideOverlay();
        }
    });

    // Close on ESC
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') {
            hideOverlay();
        }
    });
});