import { defineConfig } from 'cspell';

export default defineConfig({
  version: "0.2",
  allowCompoundWords: true,
  words: ["dropdx"],
  import: [
    "@spellbookx/cspell-config"
  ],
  dictionaryDefinitions: [
    {
      name: "devenv",
      path: "./.config/cspell/dicts/devenv.txt",
      addWords: true
    },
    {
      name: "misc",
      path: "./.config/cspell/dicts/misc.txt",
      addWords: true
    }
  ],
  dictionaries: ["misc", "devenv"],
  languageSettings: [
    { languageId: "*", dictionaries: ["misc", "devenv"] }
  ]
})
