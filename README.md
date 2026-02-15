# plug CLI

A highly opinionated, lightning-fast scaffolding tool for WordPress Plugins.
It generates boilerplates, sets up tests, and initializes private GitHub repositories in seconds.

## Features

- **Scaffolding**: Generate WordPress plugins from custom templates.
- **Additive Configuration**: Combine different template "starters" (e.g., Base + React + API).
- **GitHub Integration**: Automatically creates a private repo and pushes code using the `gh` CLI.
- **Zero-Config Logic**: Uses your local file system for templates (`~/.config/plug/templates`).

## Prerequisites

- [Go](https://go.dev/) (to build)
- [GitHub CLI (gh)](https://cli.github.com/) (must be authenticated)

## Installation

```bash
# Clone the repo
git clone [https://github.com/k1tesurfen/plug.git](https://github.com/k1tesurfen/plug.git)

# Build the binary
cd plug
go build -o plug main.go

# Move to path (optional)
sudo mv plug /usr/local/bin/


```

## Example Use

```bash
# 1. Generate necessary config files (do this once).
plug setup

# 2. Create new plugin and create a remote repository on github.com
# For this to work you need gh cli installed and authorized,
# as well as .config/plug/config.toml having the github username/org set.
plug init my-new-plugin --gh

```

```

```
