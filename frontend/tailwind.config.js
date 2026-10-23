/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // 主色（背景）60%
        dark: {
          50: '#F3F7FC',
          100: '#E6EDF7',
          200: '#CBD6E4',
          300: '#A7B4C7',
          400: '#7C8BA1',
          500: '#5E6E86',
          600: '#46566E',
          700: '#2E3D56',
          800: '#1D2A47',
          900: '#16203A',
          950: '#0B1220',
        },
        // 强调色（唯一 CTA）10%
        accent: {
          50: '#EFFEF9',
          100: '#D9FBF2',
          200: '#B3F5E4',
          300: '#7DEBD0',
          400: '#4CE0BC',
          500: '#22D3A6',
          600: '#17B890',
          700: '#0E9A78',
          800: '#0B7D62',
          900: '#07604B',
        },
        // 语义色（仅状态提示，不与强调色抢层级）
        warn: '#F5B54A',
        danger: '#FF6B6B',
        neon: {
          blue: '#22D3A6',
          purple: '#0E9A78',
          green: '#22D3A6',
          orange: '#F5B54A',
          pink: '#FF6B6B',
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
      animation: {
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'slide-up': 'slideUp 0.32s cubic-bezier(0.22, 1, 0.36, 1)',
        'slide-down': 'slideDown 0.32s cubic-bezier(0.22, 1, 0.36, 1)',
        'fade-in': 'fadeIn 0.24s ease-out',
        'card-in': 'cardIn 0.34s cubic-bezier(0.22, 1, 0.36, 1)',
        'card-out': 'cardOut 0.2s ease-in',
      },
      keyframes: {
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        slideDown: {
          '0%': { transform: 'translateY(-10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        cardIn: {
          '0%': { transform: 'translateX(24px) scale(0.98)', opacity: '0' },
          '100%': { transform: 'translateX(0) scale(1)', opacity: '1' },
        },
        cardOut: {
          '0%': { transform: 'translateX(0) scale(1)', opacity: '1' },
          '100%': { transform: 'translateX(-24px) scale(0.98)', opacity: '0' },
        },
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'mesh-gradient': 'linear-gradient(135deg, #0B1220 0%, #16203A 60%, #1D2A47 100%)',
      },
    },
  },
  plugins: [],
}
