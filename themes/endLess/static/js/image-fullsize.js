document.addEventListener('DOMContentLoaded', function() {
    console.log('Script loaded');
    
    // Find all images in post content
    const images = document.querySelectorAll('.post-content img');
    console.log('Found images:', images.length);
    
    images.forEach(img => {
        console.log('Original src:', img.src);
        let fullSizeUrl = img.src;
        
        // Remove -scaled from URL
        if (fullSizeUrl.includes('-scaled')) {
            fullSizeUrl = fullSizeUrl.replace('-scaled', '');
            console.log('Removed scaled:', fullSizeUrl);
        }
        
        // Update image src to full size
        img.src = fullSizeUrl;
        
        // Handle the wrapping anchor
        const parentAnchor = img.closest('a');
        if (parentAnchor) {
            console.log('Updating existing anchor:', parentAnchor.href);
            parentAnchor.href = fullSizeUrl;
        } else {
            console.log('Creating new anchor for:', fullSizeUrl);
            const wrapper = document.createElement('a');
            wrapper.href = fullSizeUrl;
            wrapper.style.cursor = 'zoom-in';
            img.parentNode.insertBefore(wrapper, img);
            wrapper.appendChild(img);
        }
    });
});