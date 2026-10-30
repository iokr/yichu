import type { Config } from 'tailwindcss'

export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        primary: { DEFAULT: '#A67C89', hover: '#94697A', soft: '#F5EDEF' },
        accent: '#7E9A93',
        bg: '#FAFAF8',
        surface: { DEFAULT: '#FFFFFF', 2: '#F4F3F1' },
        line: '#EAE8E5',
        ink: { DEFAULT: '#1F1D1B', 2: '#6B6661', 3: '#A6A19B' },
        success: '#6E9C7C',
        warning: '#D1A055',
        danger: '#C06C6C',
        info: '#7C93AE',
      },
      borderRadius: { md: '12px', lg: '16px', xl: '24px' },
      boxShadow: {
        sm: '0 1px 2px rgba(31,29,27,.04)',
        md: '0 4px 16px rgba(31,29,27,.06)',
        lg: '0 12px 32px rgba(31,29,27,.10)',
      },
      aspectRatio: { photo: '3 / 4' },
      fontFamily: {
        sans: ['Inter', 'PingFang SC', 'HarmonyOS Sans SC', 'Microsoft YaHei', 'system-ui', 'sans-serif'],
      },
      maxWidth: { page: '1440px' },
      transitionTimingFunction: {
        smooth: 'cubic-bezier(.4,0,.2,1)',
        drawer: 'cubic-bezier(.32,.72,0,1)',
      },
      animation: {
        'fade-up': 'fadeUp .24s cubic-bezier(.4,0,.2,1) both',
      },
      keyframes: {
        fadeUp: {
          '0%': { opacity: '0', transform: 'translateY(8px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
} satisfies Config
