/** @type {import('tailwindcss').Config} */
export const content = [
  "./views/*.templ",
]

export const theme = {
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
      pastel: {
        blue: '#AEC6CF',
        purple: '#CBAACB',
        pink: '#FFB6C1',
        orange: '#FFDAB9',
        green: '#C4E17F',
        yellow: '#FFDD94',
        gray: '#E5E5E5'
      }
    }
  }
}

export const darkMode = 'selector'

export const plugins = [require("@tailwindcss/forms"), require("@tailwindcss/typography")]
