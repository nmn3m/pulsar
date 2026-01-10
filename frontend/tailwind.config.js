/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        // Catppuccin Latte inspired colors (pgvoyager light theme)
        base: {
          DEFAULT: '#eff1f5',
          50: '#eff1f5',
          100: '#e6e9ef',
          200: '#dce0e8',
          300: '#bcc0cc',
          400: '#9ca0b0',
          500: '#7c7f93',
          600: '#6c6f85',
          700: '#5c5f77',
          800: '#4c4f69',
          900: '#3c3f52',
        },
        // Primary - Blue (main accent color)
        primary: {
          50: '#f0f5ff',
          100: '#e0ebff',
          200: '#c7d9fe',
          300: '#a5c0fc',
          400: '#7287fd',
          500: '#1e66f5',
          600: '#1a5ce0',
          700: '#1652cc',
          800: '#1248b8',
          900: '#0e3ea4',
        },
        // Secondary - Lavender
        secondary: {
          DEFAULT: '#7287fd',
          light: '#8899ff',
          dark: '#5c6cdb',
        },
        // Surface colors
        surface: {
          DEFAULT: '#e6e9ef',
          hover: '#dce0e8',
          active: '#ccd0da',
        },
        // Status colors
        success: {
          DEFAULT: '#40a02b',
          light: '#4db835',
          dark: '#358f24',
        },
        warning: {
          DEFAULT: '#df8e1d',
          light: '#e9a033',
          dark: '#c97d19',
        },
        error: {
          DEFAULT: '#d20f39',
          light: '#e5334d',
          dark: '#b80d32',
        },
        info: {
          DEFAULT: '#04a5e5',
          light: '#1fb8f5',
          dark: '#0394cc',
        },
        // Text colors
        text: {
          DEFAULT: '#4c4f69',
          muted: '#6c6f85',
          dim: '#8c8fa1',
        },
        // Border color
        border: {
          DEFAULT: '#ccd0da',
        },
      },
      backgroundImage: {
        'base-gradient': 'linear-gradient(135deg, #eff1f5 0%, #e6e9ef 100%)',
        'glow-primary':
          'radial-gradient(ellipse at center, rgba(30, 102, 245, 0.1) 0%, transparent 70%)',
      },
      boxShadow: {
        'glow-primary': '0 0 20px rgba(30, 102, 245, 0.3), 0 0 40px rgba(30, 102, 245, 0.15)',
        'glow-success': '0 0 20px rgba(64, 160, 43, 0.3), 0 0 40px rgba(64, 160, 43, 0.15)',
        'glow-error': '0 0 20px rgba(210, 15, 57, 0.3), 0 0 40px rgba(210, 15, 57, 0.15)',
      },
      animation: {
        'pulse-glow': 'pulse-glow 2s ease-in-out infinite',
      },
      keyframes: {
        'pulse-glow': {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.5' },
        },
      },
    },
  },
  plugins: [],
};
