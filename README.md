# dropdx

A cross-platform CLI tool to sync and update Personal Access Tokens (PATs) and configurations across different machines.

## Purpose

`dropdx` manages your local configuration files (like `.npmrc`, `.env`, etc.) by using a local directory (default: `~/.dropdx`) as a "Single Source of Truth". This directory is intended to be a Git repository, allowing you to synchronize your tokens and templates across multiple machines securely.

## Features

- **Cross-platform**: Works on Linux (bash/zsh) and Windows (PowerShell).
- **Template Engine**: Inject your secret tokens into templates using Go's standard templating engine.
- **Git Sync**: Easily pull and push your configuration changes using built-in git commands.
- **Obfuscated List**: View your configured tokens and providers without exposing full secrets.

## Getting Started

### Installation

```bash
# Clone the repository
git clone https://github.com/dcdavidev/dropdx.git
cd dropdx

# Build the binary
go build -o dropdx ./cmd/dropdx/main.go

# (Optional) Move to your PATH
mv dropdx /usr/local/bin/
```

### Initialization

Initialize your `dropdx` home directory:

```bash
dropdx init
```

This creates `~/.dropdx` with a `config.yaml` and a `templates/` folder.

### Configuration

Edit `~/.dropdx/config.yaml` to add your tokens and define providers:

```yaml
tokens:
  npm_token: "npm_..."
  github_token: "ghp_..."

providers:
  npm:
    template: "templates/.npmrc.tmpl"
    target: "~/.npmrc"
```

Create a template in `~/.dropdx/templates/.npmrc.tmpl`:

```text
//registry.npmjs.org/:_authToken={{.npm_token}}
```

### Usage

**Apply configurations:**
```bash
dropdx apply npm
# or apply all providers
dropdx apply
```

**Sync with Git:**
```bash
# First, initialize git in ~/.dropdx if you haven't
cd ~/.dropdx
git init
git remote add origin <your-repo-url>

# Then sync anytime
dropdx sync
```

**List status:**
```bash
dropdx list
```

## Technical Details

- **Language**: Go
- **CLI Framework**: spf13/cobra
- **Configuration**: spf13/viper (YAML)
- **Home Directory**: Prioritizes `DROPDX_HOME` environment variable, then defaults to `~/.dropdx`.

## License

MIT
