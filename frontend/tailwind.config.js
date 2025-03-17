/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#2DD4BF',
          50: '#DEFAF7',
          100: '#C8F6F1',
          200: '#9EEFE5',
          300: '#74E8D9',
          400: '#4AE1CD',
          500: '#2DD4BF',
          600: '#22A898',
          700: '#187C71',
          800: '#0E504A',
          900: '#042423',
        },
      },
    },
  },
  plugins: [],
  darkMode: 'class',
}
