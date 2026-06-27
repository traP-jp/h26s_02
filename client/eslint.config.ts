import { defineConfig, globalIgnores } from 'eslint/config'
import pluginVue from 'eslint-plugin-vue'
import tseslint from 'typescript-eslint'
import eslintConfigPrettier from 'eslint-config-prettier'
import js from '@eslint/js'

// エディタ上の Lint と npx eslint . による Lint の両方に対応

export default defineConfig([
  globalIgnores(['dist']),

  {
    files: ['src/**/*.{ts,vue}'],
    extends: [
      js.configs.recommended,
      tseslint.configs.recommendedTypeChecked,
      tseslint.configs.stylisticTypeChecked, // TypeScript 用のベース設定
    ],
    languageOptions: {
      parserOptions: {
        parser: tseslint.parser,
        project: ['./tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
        extraFileExtensions: ['.vue'],
      },
    },
  },

  {
    files: ['src/**/*.vue'],
    extends: [pluginVue.configs['flat/recommended']],
    rules: {
      'vue/component-name-in-template-casing': ['warn', 'PascalCase'],
      'vue/no-template-target-blank': ['error', { enforceDynamicLinks: 'always' }],
      'vue/no-v-html': 'error',
      // 'vue/multi-word-component-names': 'off',
    },
  },
  {
    files: ['src/**/*.{ts,vue}'],
    rules: {
      '@typescript-eslint/consistent-type-definitions': 'off', // interface 推奨を無効化
      '@typescript-eslint/no-unnecessary-condition': 'warn', // 不要な条件分岐の警告
      'no-undef': 'off', // TypeScript で型チェックされるので不要
      '@typescript-eslint/no-unused-vars': [
        'warn',
        {
          vars: 'all',
          args: 'after-used', // 使用された引数の後ろに続く引数をチェック
          argsIgnorePattern: '^_',
          ignoreRestSiblings: true, // 不要なキーを取り除く操作を許容
        },
      ],
    },
  },

  eslintConfigPrettier,
  // Prettier のルールと衝突する ESLint の設定を無効化するだけ。常に最後に置く
])
