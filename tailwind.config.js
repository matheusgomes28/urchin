/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/*.templ"],
  darkMode: "class",
  theme: {
    container: {
      center: true,
      padding: {
        DEFAULT: "1rem",
        mobile: "1rem",
        tablet: "1rem",
        desktop: "1rem",
      },
    },
    extend: {
      colors: {
        "urchin-bg-dark": "var(--color-urchin-bg-dark)",
        "urchin-bg": "var(--color-urchin-bg)",
        "urchin-bg-light": "var(--color-urchin-bg-light)",
        "urchin-text": "var(--color-urchin-text)",
        "urchin-text-muted": "var(--color-urchin-text-muted)",
        "urchin-highlight": "var(--color-urchin-highlight)",
        "urchin-border": "var(--color-urchin-border)",
        "urchin-border-muted": "var(--color-urchin-border-muted)",
        "urchin-primary": "var(--color-urchin-primary)",
        "urchin-primary-highlight": "var(--color-urchin-primary-highlight)",
        "urchin-secondary": "var(--color-urchin-secondary)",
        "urchin-secondary-highlight": "var(--color-urchin-secondary-highlight)",
        "urchin-danger": "var(--color-urchin-danger)",
        "urchin-warning": "var(--color-urchin-warning)",
        "urchin-success": "var(--color-urchin-success)",
        "urchin-info": "var(--color-urchin-info)",
      },
    },
  },
  plugins: [require("@tailwindcss/forms"), require("@tailwindcss/typography")],
};
