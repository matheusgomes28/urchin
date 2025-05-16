// MARK: Images Modal

/**
 * Shows the modal image popup
 * @param {string} title title of the image
 */
function showImageModal(title, excerpt, src) {
  document.getElementById("modal-title").innerHTML = title;
  document.getElementById("modal-excerpt").innerHTML = excerpt;
  document.getElementById("modal-image").src = src;
  modal.showModal();
}

document.getElementById('menu-toggle').addEventListener('click', function () {
  const menu = document.getElementById('mobile-menu');
  menu.classList.toggle('hidden')
});

const themeToggle = document.getElementById('theme-toggle');
const lightIcon = document.getElementById('light-icon');
const darkIcon = document.getElementById('dark-icon');

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

document.getElementById('demo-form').addEventListener('submit', function() {
  const event = new Event('verified');
  const elem = document.querySelector("#demo-form");
  elem.dispatchEvent(event);
})
