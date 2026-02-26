# dropdx 🛡️

The official Node.js wrapper for **dropdx**, a cross-platform CLI to sync and update Personal Access Tokens (PATs) and configurations securely.

## Features

- **Cross-platform**: Works on Linux (bash/zsh) and Windows (PowerShell).
- **Template Engine**: Inject secret tokens into templates using Go's standard templating engine.
- **Git Sync**: Pull and push configuration changes using built-in Git commands.
- **Obfuscated List**: View configured tokens and providers without exposing full secrets.
- **Multiple Binaries**: Available as `dropdx` and the shorter `dpdx`.

## Installation

```bash
npm install -g dropdx
```

## Quick Start

### 1. Initialize

```bash
dropdx init
```
This creates `~/.dropdx` with a `config.yaml` and a `templates/` folder.

### 2. Set a Token

```bash
dropdx set-token github --exp 30d
```

### 3. Apply Configurations

```bash
dropdx apply
```

### 4. Sync with Git

```bash
dropdx sync
```

## CLI Reference 📖

### Setup & Sync

- `dropdx init`: Initialize the home directory and default configuration.
- `dropdx sync`: Perform git pull/push on the configuration directory.
- `dropdx version`: Print the version number.

### Token Management

- `dropdx set-token <provider>`: Store a new token with optional expiration and description.
- `dropdx list`: Display all configured tokens (obfuscated) and providers.

### Application

- `dropdx apply [provider]`: Inject tokens into templates for a specific provider or all.

## Global Flags

- `--config <path>`: Override the default config file path.
- `--help, -h`: Help for any command.

## Documentation

For full documentation and technical details, visit the [official repository](https://github.com/dropdx/dropdx).

## License

MIT © [dcdavidev](https://github.com/dcdavidev)
