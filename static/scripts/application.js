document.getElementById('menu-toggle').addEventListener('click', function () {
  const menu = document.getElementById('mobile-menu');
  menu.classList.toggle('hidden')
});

const themeToggle = document.getElementById('theme-toggle');
const lightIcon = document.getElementById('light-icon');
const darkIcon = document.getElementById('dark-icon');
const galleryDropdownButton = document.getElementById('gallery-dropdown-button');
const galleryDropdownList = document.getElementById('gallery-dropdown-list');

const html = document.documentElement;

if (localStorage.getItem('color-theme') === 'dark') {
  html.classList.add('dark');
} else {
  html.classList.remove('dark');
}

themeToggle.addEventListener('click', function () {
  // Toggle dark class on HTML element
  html.classList.toggle('dark');

  if (html.classList.contains('dark')) {
    localStorage.setItem('color-theme', 'dark');
  } else {
    localStorage.setItem('color-theme', 'light');
  }
});

// Toggle dropdown
galleryDropdownButton.addEventListener('click', function() {
    const isHidden = galleryDropdownList.classList.contains('hidden');
    
    if (isHidden) {
        galleryDropdownList.classList.remove('hidden');
        galleryDropdownButton.setAttribute('aria-expanded', 'true');
    } else {
        galleryDropdownList.classList.add('hidden');
        galleryDropdownButton.setAttribute('aria-expanded', 'false');
    }
});
    
// Close when clicking outside
document.addEventListener('click', function(event) {
    if (!galleryDropdownButton.contains(event.target) && !galleryDropdownList.contains(event.target)) {
        galleryDropdownList.classList.add('hidden');
        galleryDropdownButton.setAttribute('aria-expanded', 'false');
    }
});

// This is not working
// document.getElementById('demo-form').addEventListener('submit', function() {
//   const event = new Event('verified');
//   const elem = document.querySelector("#demo-form");
//   elem.dispatchEvent(event);
// })
