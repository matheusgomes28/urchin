// MARK: Enhanced Images Modal with Improved UX
class ImageModalManager {
  constructor() {
    this.pageImages = [];
    this.currImage = 0;
    this.modal = null;
    this.isLoading = false;
    this.touchStartX = 0;
    this.touchEndX = 0;
    this.preloadedImages = new Map();
    
    window.addEventListener("load", () => {
      this.refreshImages(window.imagesData);
      this.setupEventListeners();
      this.modal = document.getElementById("modal");
    });
  }

  setupEventListeners() {
    // Keyboard navigation
    window.addEventListener("keydown", (e) => {
      if (!this.modal?.hasAttribute("open")) return;
      
      switch(e.key) {
        case "ArrowLeft":
          e.preventDefault();
          this.prevImage();
          break;
        case "ArrowRight":
          e.preventDefault();
          this.nextImage();
          break;
        case "Escape":
          e.preventDefault();
          this.closeModal();
          break;
        case " ": // Spacebar
          e.preventDefault();
          this.nextImage();
          break;
      }
    });

    // Touch/swipe support
    if (this.modal) {
      this.modal.addEventListener("touchstart", (e) => {
        this.touchStartX = e.changedTouches[0].screenX;
      }, { passive: true });

      this.modal.addEventListener("touchend", (e) => {
        this.touchEndX = e.changedTouches[0].screenX;
        this.handleSwipe();
      }, { passive: true });
    }

    // Click outside modal to close
    window.addEventListener("click", (e) => {
      if (e.target === this.modal) {
        this.closeModal();
      }
    });
  }

  handleSwipe() {
    const swipeThreshold = 50;
    const diff = this.touchStartX - this.touchEndX;
    
    if (Math.abs(diff) > swipeThreshold) {
      if (diff > 0) {
        this.nextImage(); // Swipe left = next
      } else {
        this.prevImage(); // Swipe right = previous
      }
    }
  }

  // Preload adjacent images for smoother experience
  preloadAdjacentImages() {
    const preloadIndices = [
      (this.currImage + 1) % this.pageImages.length,
      (this.currImage - 1 + this.pageImages.length) % this.pageImages.length
    ];

    preloadIndices.forEach(index => {
      const image = this.pageImages[index];
      if (image && !this.preloadedImages.has(image.filepath)) {
        const img = new Image();
        img.src = image.filepath;
        this.preloadedImages.set(image.filepath, img);
      }
    });
  }

  showLoadingState() {
    const modalImage = document.getElementById("modal-image");
    const loadingSpinner = document.getElementById("modal-loading");
    
    if (modalImage) modalImage.style.opacity = "0.5";
    if (loadingSpinner) loadingSpinner.style.display = "block";
    this.isLoading = true;
  }

  hideLoadingState() {
    const modalImage = document.getElementById("modal-image");
    const loadingSpinner = document.getElementById("modal-loading");
    
    if (modalImage) modalImage.style.opacity = "1";
    if (loadingSpinner) loadingSpinner.style.display = "none";
    this.isLoading = false;
  }

  updatePaginationUI() {
    // Clear all active states
    for (let i = 0; i < this.pageImages.length; i++) {
      const paginationBtn = document.getElementById(`modal-pagination-${i}`);
      if (paginationBtn) {
        paginationBtn.classList.remove('text-urchin-text', 'bg-urchin-secondary-highlight');
        paginationBtn.setAttribute('aria-pressed', 'false');
      }
    }

    // Set active state for current image
    const currentBtn = document.getElementById(`modal-pagination-${this.currImage}`);
    if (currentBtn) {
      currentBtn.classList.add('text-urchin-text', 'bg-urchin-secondary-highlight');
      currentBtn.setAttribute('aria-pressed', 'true');
      currentBtn.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'center' });
    }

    // Update counter display
    const counter = document.getElementById("modal-counter");
    if (counter) {
      counter.textContent = `${this.currImage + 1} of ${this.pageImages.length}`;
    }

    // Update navigation button states
    this.updateNavigationButtons();
  }

  updateNavigationButtons() {
    const prevBtn = document.querySelector('[onclick="prevImage()"]');
    const nextBtn = document.querySelector('[onclick="nextImage()"]');
    
    if (this.pageImages.length <= 1) {
      if (prevBtn) prevBtn.style.display = 'none';
      if (nextBtn) nextBtn.style.display = 'none';
    } else {
      if (prevBtn) prevBtn.style.display = 'flex';
      if (nextBtn) nextBtn.style.display = 'flex';
    }
  }

  updateImageContent(image) {
    // Update with smooth transitions
    const elements = {
      title: document.getElementById("modal-title"),
      excerpt: document.getElementById("modal-excerpt"),
      image: document.getElementById("modal-image"),
      date: document.getElementById("modal-text-date"),
      location: document.getElementById("modal-text-name")
    };

    // Add fade transition class if not present
    Object.values(elements).forEach(el => {
      if (el && !el.classList.contains('transition-opacity')) {
        el.classList.add('transition-opacity', 'duration-200');
      }
    });

    if (elements.title) elements.title.textContent = image.name || 'Untitled';
    if (elements.excerpt) elements.excerpt.textContent = image.excerpt || image.name || 'No description available';
    if (elements.date) elements.date.textContent = image.date || 'Date not available';
    if (elements.location) elements.location.textContent = image.location?.name || 'Location not available';

    // Handle image loading with error handling
    if (elements.image) {
      this.showLoadingState();
      
      const newImage = new Image();
      newImage.onload = () => {
        elements.image.src = image.filepath;
        elements.image.alt = image.name || 'Image';
        this.hideLoadingState();
        this.preloadAdjacentImages();
      };
      
      newImage.onerror = () => {
        elements.image.src = '/static/images/placeholder.jpg'; // Fallback image
        elements.image.alt = 'Image failed to load';
        this.hideLoadingState();
        console.warn(`Failed to load image: ${image.filepath}`);
      };
      
      newImage.src = image.filepath;
    }
  }

  showImageModal(index) {
    if (this.isLoading || index < 0 || index >= this.pageImages.length) return;

    this.currImage = index;
    const image = this.pageImages[this.currImage];
    
    if (!image) {
      console.error(`No image found at index ${index}`);
      return;
    }

    this.updateImageContent(image);
    this.updatePaginationUI();
    
    if (this.modal) {
      this.modal.showModal();
      // Focus management for accessibility
      this.modal.focus();
      // Prevent body scroll when modal is open
      document.body.style.overflow = 'hidden';
    }
  }

  closeModal() {
    if (this.modal) {
      this.modal.close();
      document.body.style.overflow = '';
    }
  }

  nextImage() {
    if (this.pageImages.length <= 1) return;
    
    const nextIndex = (this.currImage + 1) % this.pageImages.length;
    this.showImageModal(nextIndex);
  }

  prevImage() {
    if (this.pageImages.length <= 1) return;
    
    const prevIndex = (this.currImage - 1 + this.pageImages.length) % this.pageImages.length;
    this.showImageModal(prevIndex);
  }

  // Enhanced validation with better error handling
  validateImageData(images) {
    if (!Array.isArray(images)) {
      console.error("Images data must be an array");
      return [];
    }

    return images.filter((img, index) => {
      const requiredFields = ['uuid', 'filename', 'name', 'date'];
      const missingFields = requiredFields.filter(field => !img.hasOwnProperty(field));
      
      if (missingFields.length > 0) {
        console.warn(`Image at index ${index} missing required fields: ${missingFields.join(', ')}`);
        return false;
      }

      // Validate location object
      if (img.location && typeof img.location === 'object') {
        const locationFields = ['name'];
        const missingLocationFields = locationFields.filter(field => !img.location.hasOwnProperty(field));
        if (missingLocationFields.length > 0) {
          console.warn(`Image at index ${index} has incomplete location data`);
        }
      }

      return true;
    });
  }

  sanitizeImageData(images) {
    return images.map(img => ({
      ...img,
      excerpt: img.excerpt || img.name || 'No description available',
      filepath: img.filepath || `/images/data/${img.filename}`,
      location: {
        name: img.location?.name || 'Unknown location',
        latitude: img.location?.latitude || null,
        longitude: img.location?.longitude || null
      }
    }));
  }

  refreshImages(images) {
    if (!images) {
      console.warn("No image data provided");
      return;
    }

    this.pageImages.length = 0; // Clear existing images
    
    const validImages = this.validateImageData(images);
    const sanitizedImages = this.sanitizeImageData(validImages);
    
    this.pageImages.push(...sanitizedImages);
    
    console.log(`Loaded ${this.pageImages.length} images`);
  }

  // Public API methods
  getCurrentImage() {
    return this.pageImages[this.currImage];
  }

  getTotalImages() {
    return this.pageImages.length;
  }

  goToImage(index) {
    this.showImageModal(index);
  }
}

// Initialize the image modal manager
const imageModal = new ImageModalManager();

// Global functions for template compatibility
function showImageModal(index) {
  imageModal.showImageModal(index);
}

function nextImage() {
  imageModal.nextImage();
}

function prevImage() {
  imageModal.prevImage();
}