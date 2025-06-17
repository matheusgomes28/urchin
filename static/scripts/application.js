function initDropdowns() {
    const dropdownButtons = document.querySelectorAll("[id^=dropdown-button]");

    // We assume that every dropdown button will have a matching
    // dropdown list!
    for (const buttonElem of dropdownButtons) {
        const listId = buttonElem.id.replace("dropdown-button", "dropdown-list");
        const listElem = document.getElementById(listId);
        const arrowElem = buttonElem.querySelector("[id^=dropdown-arrow");

        // Switches the arrow element's direction
        const toggleArrowDir = () => {
            arrowElem.classList.toggle("icon-caret-down");
            arrowElem.classList.toggle("icon-caret-up");
        };

        buttonElem.addEventListener('click', function () {
            const isHidden = listElem.classList.contains('hidden');

            if (isHidden) {
                listElem.classList.remove('hidden');
                buttonElem.setAttribute('aria-expanded', 'true');
                toggleArrowDir();
            } else {
                listElem.classList.add('hidden');
                buttonElem.setAttribute('aria-expanded', 'false');
                toggleArrowDir();
            }
        });

        // close when clicking outside
        // Lazy - should probably do all of the elements in one go
        document.addEventListener('click', function (event) {
            if (!buttonElem.contains(event.target) && !listElem.contains(event.target)) {
                listElem.classList.add('hidden');
                buttonElem.setAttribute('aria-expanded', 'false');
                arrowElem.classList.remove("icon-caret-up");
                arrowElem.classList.add("icon-caret-down");
            }
        });
    }
}

function initThemeToggles() {
    const themeToggle = document.getElementById('theme-toggle');
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
}

function initMenuToggle() {
    document.getElementById('menu-toggle').addEventListener('click', function () {
        const menu = document.getElementById('mobile-menu');
        menu.classList.toggle('hidden')
    });
}

// Main entrypoint
function init() {
    initThemeToggles();
    initDropdowns();
    initMenuToggle();
}

// Initialize when DOM is loaded
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}

export { initDropdowns, initThemeToggles, init }

// This is not working
// document.getElementById('demo-form').addEventListener('submit', function() {
//   const event = new Event('verified');
//   const elem = document.querySelector("#demo-form");
//   elem.dispatchEvent(event);
// })
