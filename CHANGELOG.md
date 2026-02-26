# Changelog

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
