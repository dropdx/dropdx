# Changelog

## [0.6.0](https://github.com/dropdx/dropdx/compare/v0.5.0...v0.6.0) (2026-03-02)

### Features

* unify gh/github templates and improve migration ([fdfb0f5](https://github.com/dropdx/dropdx/commit/fdfb0f521c6d5981ff9d378ab46d65de6ea0d43c))

### Bug Fixes

* ensure apply command uses v2 list structure and fix config loading ([65cc8b4](https://github.com/dropdx/dropdx/commit/65cc8b467527a66b0eac0fbb4a116b13cb35168c))

## [0.5.0](https://github.com/dropdx/dropdx/compare/v0.4.0...v0.5.0) (2026-03-02)

### Features

* implement config v2 with universal token lists and migration ([eae46af](https://github.com/dropdx/dropdx/commit/eae46afbcde18993d09dc327b5b7fb7a01dd2821))

## [0.4.0](https://github.com/dropdx/dropdx/compare/v0.3.3...v0.4.0) (2026-02-27)

### Features

* add dropdx create machine command and machines configuration ([9a56cfd](https://github.com/dropdx/dropdx/commit/9a56cfdae4ba0c03a7303b817c2beefa6095ff11))
* add shell source advice and fix build errors ([c5c751a](https://github.com/dropdx/dropdx/commit/c5c751abbe2d9b9003d6da08801efebc953bbaca))
* add ssh-config and remotes management commands ([e4f1c2b](https://github.com/dropdx/dropdx/commit/e4f1c2bc65cce675133eb4dff25389c646fdb679))
* add sync ssh-keys command ([a1f538a](https://github.com/dropdx/dropdx/commit/a1f538a4673dc9f825251ca9352b00504dcba011))
* add version check to dropdx ([d1b6df6](https://github.com/dropdx/dropdx/commit/d1b6df6889052e8f83abffff599df0e6b92c8b11))
* auto-create .gitignore during dropdx init ([89cf3c8](https://github.com/dropdx/dropdx/commit/89cf3c8c273e954f40ebd09ee91c8f7c3bd594d6))
* improve apply interactive mode and pass nested tokens to templates ([a836d29](https://github.com/dropdx/dropdx/commit/a836d29203ad564f31bd0e9d9c1b3338c0725355))
* make all commands interactive ([d43b30a](https://github.com/dropdx/dropdx/commit/d43b30a0a8e8dc20ab98aa82832a89edf8e0a8cf))
* make set-token interactive and support multiple npm registries ([9a28508](https://github.com/dropdx/dropdx/commit/9a28508f93fc7b1dd2bb1f837b0489aeaa00b4ee))
* rename create ssh-config to create remote ([5fb6e77](https://github.com/dropdx/dropdx/commit/5fb6e77a82415f2bb7155faa72cf46f81b2428b8))
* restructure sync command and add repository subcommand with autostash ([4e12e03](https://github.com/dropdx/dropdx/commit/4e12e03c384c38233e5b606565f527a3c1cddb3f))
* support appending/updating shell config files and add github provider ([4a43062](https://github.com/dropdx/dropdx/commit/4a430629eb59baef2868595b18e44c5b807f1ed7))

### Bug Fixes

* ensure github provider and template are always available for apply ([504f451](https://github.com/dropdx/dropdx/commit/504f4513a32b874c573b3caea606ed558cc0b7dd))
* ensure github provider is always available in apply ([acc9d36](https://github.com/dropdx/dropdx/commit/acc9d368bf26089324a1b0fc2a52da42a7b98650))
* ensure provider name key is populated in templates for multi-registry ([50f6c08](https://github.com/dropdx/dropdx/commit/50f6c0844d9b82c9b010c0fb662e8f0ba36e9caf))
* force rebuild and add debug marker ([bf15ddc](https://github.com/dropdx/dropdx/commit/bf15ddcc7df4d5cba6fcca74c4055716ad3c7eb7))
* resolve build errors in apply and list commands ([0527837](https://github.com/dropdx/dropdx/commit/0527837ce85e073b83b1fc5b598c17cb30701ccf))
* resolve build errors in engine.go ([e927ede](https://github.com/dropdx/dropdx/commit/e927ededd0689ce6353596048af399a5110a2aaf))
* ultra-robust token mapping for multi-registry npm templates ([1318210](https://github.com/dropdx/dropdx/commit/1318210fa296adc3d012fa28305d8e2766ae10fd))
* use yaml.Unmarshal directly to fix map keys with dots ([732c46e](https://github.com/dropdx/dropdx/commit/732c46e96a23f7f7ec9ec3e5fd6ed8c4fc2d366f))

## [0.3.3](https://github.com/dropdx/dropdx/compare/v0.3.2...v0.3.3) (2026-02-27)

## [0.3.2](https://github.com/dropdx/dropdx/compare/v0.3.1...v0.3.2) (2026-02-27)

### Bug Fixes

* ensure binary is executable and add snapcraft config ([cadc5d8](https://github.com/dropdx/dropdx/commit/cadc5d8d3b9802af2e014daf0d8a492ddb3e2926))
* ensure binary is executable before spawning in apps/cli ([4ff87ab](https://github.com/dropdx/dropdx/commit/4ff87ab05654fb9c872008a0e6296b67a9211126))

## [0.3.1](https://github.com/dropdx/dropdx/compare/v0.3.0...v0.3.1) (2026-02-27)

## [0.3.0](https://github.com/dropdx/dropdx/compare/v0.2.1...v0.3.0) (2026-02-26)

## [0.2.1](https://github.com/dropdx/dropdx/compare/v0.2.0...v0.2.1) (2026-02-26)

### Features

* enhance CLI with colors, spinners and a colorful banner ([556caf1](https://github.com/dropdx/dropdx/commit/556caf11654484db80d811a72a4392f3b6c16a73))
* improve CLI branding and help (bastion style) ([5333a2b](https://github.com/dropdx/dropdx/commit/5333a2b6ca44d7d08a3bb93b04db62dc7bd83745))
* refactor CLI to use pterm and bastion-style architecture ([9b9cdd3](https://github.com/dropdx/dropdx/commit/9b9cdd3fe36bef53a6aaa60d7dfc4e91bcb18639))

## 0.2.0 (2026-02-26)

### Features

- add .goreleaser.yaml for multi-platform distribution ([1e7a32d](https://github.com/dropdx/dropdx/commit/1e7a32dc26d03ef79ef027884f9d9acec7e64c55))
- add default templates and configuration for PyPI and NPM in init command ([0f9fd94](https://github.com/dropdx/dropdx/commit/0f9fd94ede23f7186dfe56de8f8a6f1a70e3927b))
- add docker default template and configuration in init command ([15923a6](https://github.com/dropdx/dropdx/commit/15923a65a033143a967499dbc2268016f46a82e9))
- add dpdx alias and update snapcraft package name ([4151eba](https://github.com/dropdx/dropdx/commit/4151eba675ef4d2a1cbca6b40014215fdbd8ab95))
- add GitHub Actions pipeline for automated builds and releases ([7755fda](https://github.com/dropdx/dropdx/commit/7755fdae55eacbd367fa6ab19d7683bd879df38e))
- add JS wrapper for NPM publishing and rename cli package to dropdx ([1fca2f7](https://github.com/dropdx/dropdx/commit/1fca2f72169c890ea5c6397910e67b7cb1cdf0dd))
- add release-it configuration for automated versioning and changelog ([e8c6a35](https://github.com/dropdx/dropdx/commit/e8c6a35d95cdf15bc83f2b9d4e5e35bfd294a84e))
- add support for --exp false in set-token command ([b4ed002](https://github.com/dropdx/dropdx/commit/b4ed0025870b5920d224b8177a63b8e2d84c3abf))
- add version command, update README, and refactor engine to use Provider interface ([325db9e](https://github.com/dropdx/dropdx/commit/325db9efded39c707f978f982c7a55db7ce1fd03))
- implement configuration management and template engine for providers ([8da93d0](https://github.com/dropdx/dropdx/commit/8da93d048ff724d41ae6a3afe59143c62e63d3c7))
- implement init command for base directory setup ([a71602e](https://github.com/dropdx/dropdx/commit/a71602e67f84ec4307f47c389c48886b5aa6e11c))
- implement list command to view configured tokens and providers ([39e3d44](https://github.com/dropdx/dropdx/commit/39e3d44b140d796db5100c9bf8ca57af95eb888e))
- implement set-token command with metadata and secure prompt ([98626de](https://github.com/dropdx/dropdx/commit/98626de8e3b6511ad8ea65dafe97463140108f5f))
- implement sync command to synchronize configurations using git ([1c6cd93](https://github.com/dropdx/dropdx/commit/1c6cd93818296cec65754804dad0891db528f6ab))
- initial project structure with cobra and viper ([43f4cac](https://github.com/dropdx/dropdx/commit/43f4cac64d2a8eaad4432509d5c37ec1708db880))

### Bug Fixes

- correct .release-it.json syntax and clean working dir ([2449f1e](https://github.com/dropdx/dropdx/commit/2449f1ee0118e508fa04eb8ecc868a38370b5742))
- ignore go.mod in linting/formatting and fix package filter ([067d189](https://github.com/dropdx/dropdx/commit/067d189ba379b9eca167ac3d52bced6ea394da83))
