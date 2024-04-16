/** @type {import('tailwindcss').Config} */
export default {
  content: ["*.templ","./**/*.templ"], // this is where our templates are located
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
}
