# pivot

> ALL-IN-ONE runtime version switcher & project scaffolder.
> Python, PHP, Node.js, Go, Deno, Bun, Java, Rust — switch, install, scaffold.

## Install

**Linux / macOS**
```bash
curl -fsSL https://raw.githubusercontent.com/zidan-herlangga/pivot/main/scripts/install.sh | sh
```

**Windows (PowerShell)**
```powershell
iwr -UseBasicParsing -Uri https://raw.githubusercontent.com/zidan-herlangga/pivot/main/scripts/install.ps1 | iex
```

**From source**
```bash
git clone https://github.com/zidan-herlangga/pivot.git
cd pivot
go build -o pivot .
```

Add `~/.pivot/bin` to your PATH once, then `pivot use` handles the rest.

## Quick Start

```bash
pivot                          # Interactive TUI menu
pivot use node 22              # Activate Node.js v22
pivot list                     # Show all versions (system + portable)
pivot install python           # Download latest Python
pivot run python 3 --version   # Run with specific version
pivot doctor --fix             # Auto-fix PATH issues
pivot create react myapp       # Scaffold a React project
```

## Commands

| Command | Description |
|---|---|
| `pivot` | Interactive TUI menu |
| `pivot list` | List all installed versions |
| `pivot use <rt> <ver>` | Activate a runtime version |
| `pivot install <rt> [ver]` | Download a portable runtime |
| `pivot run <rt> <ver> <cmd>` | Run command with specific version |
| `pivot shell <rt> <ver>` | Spawn subshell with specific version |
| `pivot create <fw> <name>` | Scaffold a framework project |
| `pivot profile <op> <name>` | Save/load/list/delete profiles |
| `pivot doctor [--fix]` | Diagnose & auto-fix issues |
| `pivot clean` | Remove unused runtime versions |
| `pivot upgrade` | Self-update to latest pivot |
| `pivot update` | Check latest upstream versions |
| `pivot init` / `pin` | Create `.pivotrc` in current dir |
| `pivot env` | Print PATH setup for shell config |
| `pivot hook` | Print shell hook for auto-apply `.pivotrc` |
| `pivot completion bash\|zsh\|fish` | Generate tab completion |

**Version aliases:** `latest`, `system`, `lts` — e.g. `pivot use node lts`.

## Supported Runtimes

| Runtime | System Detection | Portable Download |
|---|---|---|
| **Python** | ✓ | Windows, Linux, macOS |
| **PHP** | ✓ | Windows |
| **Node.js** | ✓ | Windows, Linux, macOS |
| **Go** | ✓ | Windows, Linux, macOS |
| **Deno** | ✓ | Windows, Linux, macOS |
| **Bun** | ✓ | Windows, Linux, macOS |
| **Java** | ✓ | Windows, Linux, macOS (Adoptium) |
| **Rust** | ✓ | via rustup |

## Project Scaffolding

```bash
pivot create laravel myapp
pivot create react myapp
pivot create django myapp
pivot create rails myapp
```

20 frameworks supported:

**PHP:** Laravel, CodeIgniter 4, Symfony, WordPress (Bedrock)
**JS/TS:** React (Vite), Next.js, Vue (Vite), AdonisJS, Svelte (Vite), Nuxt, Solid (Vite)
**Python:** Django, Flask, FastAPI
**Go:** Gin, Echo, Fiber
**Ruby:** Ruby on Rails
**Java:** Spring Boot

## Shell Hook

Auto-apply `.pivotrc` whenever you `cd` into a directory:

```bash
eval "$(pivot hook)"        # bash/zsh
pivot hook | source         # fish
```

## Tab Completion

```bash
eval "$(pivot completion bash)"    # bash
eval "$(pivot completion zsh)"     # zsh
pivot completion fish >> ~/.config/fish/config.fish  # fish
```

## Profiles

Save and restore your entire runtime configuration:

```bash
pivot profile save backend    # saves all 8 runtimes
pivot profile load backend
pivot profile list
pivot profile delete backend
```

## .pivotrc

Place `.pivotrc` in any project directory:

```ini
python=3.12.0
node=22.0.0
go=1.22.5
```

pivot walks up the directory tree and auto-applies the nearest `.pivotrc`.

## Language

Auto-detects `LANG` / `LC_ALL` / `LC_MESSAGES` — supports **English** and **Bahasa Indonesia**.

## Directory Structure

```
~/.pivot/
├── bin/            # Active binary symlinks (add this to PATH)
├── runtimes/       # Downloaded portable runtimes
│   ├── python/
│   ├── node/
│   ├── go/
│   ├── java/
│   └── ...
├── profiles/       # Saved version profiles
├── config.json     # Current active versions
└── update-cache.json
```

## Bug Reports

Found a bug? Open an issue at [github.com/zidan-herlangga/pivot/issues](https://github.com/zidan-herlangga/pivot/issues)

## MIT License

Copyright (c) 2026 zidan-herlangga
