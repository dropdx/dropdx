import { defineConfig } from 'eslint/config';
import spellbookx from 'eslint-plugin-spellbookx';

export default defineConfig([
  ...spellbookx.configs.recommended,
  {
    ignores: [
      '**/node_modules/**',
      '**/dist/**',
      '**/.turbo/**',
      'dropdx',
      '*.exe',
      '*.test',
      '*.out',
      'coverage.*',
      'go.mod',
      'go.sum',
      '**/*.go',
    ],
  },
]);
