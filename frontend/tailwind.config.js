/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	darkMode: 'class',
	theme: {
		extend: {
			colors: {
				// Cosmic space background colors
				space: {
					950: '#050510',
					900: '#0a0a1a',
					800: '#0d1025',
					700: '#121830',
					600: '#1a2040',
					500: '#252850',
				},
				// Primary - Electric Cyan/Blue (main glow color)
				primary: {
					50: '#ecfeff',
					100: '#cffafe',
					200: '#a5f3fc',
					300: '#67e8f9',
					400: '#22d3ee',
					500: '#00d4ff',
					600: '#00b8e6',
					700: '#0891b2',
					800: '#0e7490',
					900: '#155e75',
					950: '#083344'
				},
				// Accent - Hot Pink/Magenta (alert highlights)
				accent: {
					50: '#fdf2f8',
					100: '#fce7f3',
					200: '#fbcfe8',
					300: '#f9a8d4',
					400: '#f472b6',
					500: '#ff0080',
					600: '#db2777',
					700: '#be185d',
					800: '#9d174d',
					900: '#831843',
					950: '#500724'
				},
				// Neon colors for status indicators
				neon: {
					cyan: '#00ffff',
					blue: '#00d4ff',
					pink: '#ff0080',
					purple: '#8b5cf6',
					green: '#00ff88',
					yellow: '#ffdd00',
					orange: '#ff6b35',
					red: '#ff3366'
				}
			},
			backgroundImage: {
				'space-gradient': 'radial-gradient(ellipse at center, #1a2040 0%, #0a0a1a 50%, #050510 100%)',
				'glow-cyan': 'radial-gradient(ellipse at center, rgba(0, 212, 255, 0.15) 0%, transparent 70%)',
				'glow-pink': 'radial-gradient(ellipse at center, rgba(255, 0, 128, 0.15) 0%, transparent 70%)',
			},
			boxShadow: {
				'neon-cyan': '0 0 20px rgba(0, 212, 255, 0.5), 0 0 40px rgba(0, 212, 255, 0.3)',
				'neon-pink': '0 0 20px rgba(255, 0, 128, 0.5), 0 0 40px rgba(255, 0, 128, 0.3)',
				'neon-green': '0 0 20px rgba(0, 255, 136, 0.5), 0 0 40px rgba(0, 255, 136, 0.3)',
				'neon-purple': '0 0 20px rgba(139, 92, 246, 0.5), 0 0 40px rgba(139, 92, 246, 0.3)',
			},
			animation: {
				'pulse-glow': 'pulse-glow 2s ease-in-out infinite',
			},
			keyframes: {
				'pulse-glow': {
					'0%, 100%': { opacity: '1' },
					'50%': { opacity: '0.5' },
				}
			}
		}
	},
	plugins: []
};
