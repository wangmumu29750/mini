import forms from '@tailwindcss/forms'
import type { Config } from 'tailwindcss'

export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        rail: {
          ink: '#18212f',
          line: '#d9e2ec',
          surface: '#f7fafc',
          primary: '#0f766e',
          accent: '#c2410c',
        },
      },
      boxShadow: {
        subtle: '0 1px 2px rgb(15 23 42 / 0.08)',
      },
    },
  },
  plugins: [forms],
} satisfies Config

