// @ts-check
import { createRequire } from 'module'
const require = createRequire(import.meta.url)

// 尝试从 @nuxt/eslint 获取配置
let nuxtEslintConfig = null
try {
  const configPath = require.resolve('@nuxt/eslint/config')
  const { default: config } = await import(configPath)
  nuxtEslintConfig = config
} catch (e) {
  console.log('Using fallback ESLint config')
}

// Fallback 简化配置
export default nuxtEslintConfig || {
  ignores: ['.nuxt/**', 'node_modules/**', 'dist/**', '.output/**', 'public/**', 'eslint.config.js'],
  rules: {
    'vue/multi-word-component-names': 'off',
    'no-undef': 'off',
    'no-unused-vars': 'off'
  }
}
