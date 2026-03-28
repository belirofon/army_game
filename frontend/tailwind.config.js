/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Background colors
        'background-primary': '#1A1F16',
        'background-secondary': '#252B1E',
        
        // Surface colors
        'surface': '#2D3527',
        'surface-elevated': '#353D2E',
        
        // Accent colors
        'military-green': '#4A5D3F',
        'military-green-light': '#5A6D4F',
        'military-green-dark': '#3A4D2F',
        'accent-secondary': '#8B7355',
        'accent-tertiary': '#5D7A4A',
        
        // Text colors
        'text-primary': '#E8E4DC',
        'text-secondary': '#9A968E',
        'text-muted': '#6B665E',
        
        // Status colors
        'success': '#5D7A4A',
        'warning': '#B8860B',
        'danger': '#8B3A3A',
        'info': '#4A5D7F',
        
        // Border colors
        'border': '#3D4533',
        'border-light': '#4A5540',
        
        // Light gray for buttons
        'light-gray': '#E8E4DC',
      },
      fontFamily: {
        'roboto': ['Roboto', 'sans-serif'],
        'roboto-condensed': ['Roboto Condensed', 'sans-serif'],
        'roboto-mono': ['Roboto Mono', 'monospace'],
        'inter': ['Inter', 'sans-serif'],
      },
      spacing: {
        'xs': '4px',
        'sm': '8px',
        'md': '16px',
        'lg': '24px',
        'xl': '32px',
        'xxl': '48px',
      },
      borderRadius: {
        'sm': '4px',
        'md': '8px',
        'lg': '12px',
      },
      boxShadow: {
        'card': '0 2px 8px rgba(0,0,0,0.3)',
        'modal': '0 8px 24px rgba(0,0,0,0.5)',
        'button': '0 2px 4px rgba(0,0,0,0.2)',
        'elevated': '0 4px 12px rgba(0,0,0,0.4)',
      },
    },
  },
  plugins: [],
}