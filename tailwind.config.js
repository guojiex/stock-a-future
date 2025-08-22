/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/static/**/*.{html,js,css}",
    "./web/static/*.html",
    "./web/static/js/**/*.js",
    "./web/static/js/modules/**/*.js"
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require("daisyui")
  ],
  daisyui: {
    themes: [
      "light",
      "dark",
      "corporate",
      "business"
    ],
    darkTheme: "dark",
    base: true,
    styled: true,
    utils: true,
    prefix: "",
    logs: true,
    themeRoot: ":root"
  }
}
