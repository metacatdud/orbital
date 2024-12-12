/** @type {import('tailwindcss').Config} */
const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: [
    "../../node/web/index.html",
    "../wasm/templates/**/*.html",
  ],
  theme: {
    extend: {
      screens:{
        ...defaultTheme.screens,
        'sm': '640px',
        'md': '768px',
        'lg': '1024px',
        'xl': '1280px',
        '2xl': '1536px'
      },
      fontFamily: {
        sans: ['InterVariable', ...defaultTheme.fontFamily.sans],
      },
      fontSize: {
        'xxs': '.6rem'
      }
    },
  },
  plugins: [],
};