/** @type {import('tailwindcss').Config} */
const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
    content: [
        "../../node/web/index.html",
        "../wasm/templates/**/*.html"
    ],
    darkMode: ['class', "[class='dark']"],
    theme: {
        extend: {
            colors: {
                'os': {
                    'bg': '#f8fafc',
                    'taskbar': '#f1f5f9',
                    'accent': '#e2e8f0',
                    'border': '#cbd5e1',
                    'text': '#475569'
                },
                'terminal': {
                    'bg': '#0a1322',
                    'taskbar': '#0f172a',
                    'accent': '#1e293b',
                    'border': '#00ff9580',
                    'text': '#00ff95',
                    'glow': '#00ff9520'
                }
            },
            fontFamily: {
                sans: ['InterVariable', ...defaultTheme.fontFamily.sans],
            },
            fontSize: {
                'xxs': '.6rem'
            },
            screens: {
                ...defaultTheme.screens,
                'sm': '640px',
                'md': '768px',
                'lg': '1024px',
                'xl': '1280px',
                '2xl': '1536px'
            }
        },
    },
    plugins: [],
};