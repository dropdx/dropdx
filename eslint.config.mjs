import { defineConfig } from 'eslint/config';
import spellbookxPlugin from '@spellbookx/eslint-plugin';

export default defineConfig([
  spellbookxPlugin.configs.recommended
]);
