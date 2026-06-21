# pivot

> Multi-runtime version switcher — Python, PHP, Node.js, Go, Deno, Bun, Java, Rust.

## Install

**Linux / macOS**
```bash
curl -fsSL https://raw.githubusercontent.com/zidan-herlangga/pivot/main/scripts/install.sh | sh
```

**Windows (PowerShell)**
```powershell
iwr -UseBasicParsing -Uri https://raw.githubusercontent.com/zidan-herlangga/pivot/main/scripts/install.ps1 | iex
```

**Manual (Go required)**
```bash
go install github.com/zidan-herlangga/pivot@latest
```

## Quick Start

```bash
pivot                # Interactive TUI menu
pivot use node 22    # Switch Node.js to v22
pivot list           # Show all available versions
pivot run python 3 --version
pivot doctor --fix   # Auto-diagnose and fix PATH issues
```

Add `~/.pivot/bin` to your PATH once, then `pivot use` handles the rest.

## Commands

| Command | Description |
|---|---|
| `pivot` | Interactive TUI menu |
| `pivot list` | List all installed versions |
| `pivot use <rt> <ver>` | Activate a runtime version |
| `pivot install <rt> [ver]` | Download a portable runtime |
| `pivot run <rt> <ver> <cmd>` | Run a command with specific version |
| `pivot shell <rt> <ver>` | Spawn a subshell with specific version |
| `pivot doctor [--fix]` | Diagnose system & auto-fix PATH |
| `pivot clean` | Remove unused runtime versions |
| `pivot upgrade` | Self-update to latest pivot |
| `pivot update` | Check latest upstream versions |
| `pivot init` / `pin` | Create `.pivotrc` in current dir |
| `pivot env` | Print PATH export for shell config |
| `pivot hook` | Print shell hook for auto-apply |
| `pivot completion bash\|zsh\|fish` | Generate tab completion |

**Version aliases:** `latest`, `system`, `lts` — e.g. `pivot use node lts`.

## Supported Runtimes

| Runtime | Detection | Portable Download |
|---|---|---|
| **Python** | System + Portable | Windows / Linux / macOS |
| **PHP** | System + Portable | Windows |
| **Node.js** | System + Portable | Windows / Linux / macOS |
| **Go** | System + Portable | Windows / Linux / macOS |
| **Deno** | System | — |
| **Bun** | System | — |
| **Java** | System | — |
| **Rust** | System | — |

## Run with Specific Version

Run any command without switching your active version:

```bash
pivot run python 3.12 my_script.py
pivot run node 22 npm test
pivot run go 1.26 go build .
```

Shortcut for version flags:

```bash
pivot run node 22 --version   # same as: node --version
```

## Shell Hook

Auto-apply `.pivotrc` whenever you `cd` into a directory.

```bash
eval "$(pivot hook)"        # bash/zsh
pivot hook | source         # fish
```

## Tab Completion

```bash
eval "$(pivot completion bash)"   # bash
eval "$(pivot completion zsh)"    # zsh
pivot completion fish >> ~/.config/fish/config.fish  # fish
```

## Project Scaffolding

```bash
pivot create laravel myapp
pivot create react myapp
pivot create nextjs myapp
```

Supported: Laravel, CodeIgniter 4, Symfony, WordPress (Bedrock), React, Next.js, Vue, AdonisJS.

## Profiles

Save and restore your runtime configurations:

```bash
pivot profile save backend
pivot profile load backend
pivot profile list
pivot profile delete backend
```

## .pivotrc

Place `.pivotrc` in any project directory:

```ini
python=3.12.0
node=22.0.0
```

pivot walks up the directory tree and auto-applies the nearest `.pivotrc` — no manual switching.

## Language

Auto-detects your system language — **English** and **Bahasa Indonesia** supported.

## Directory Structure

```
~/.pivot/
├── bin/            # Active runtime binary symlinks (add this to PATH)
├── runtimes/       # Downloaded portable runtimes
├── profiles/       # Saved version profiles
├── config.json     # Current active versions
└── update-cache.json
```

## Build from Source

```bash
git clone https://github.com/zidan-herlangga/pivot.git
cd pivot
go mod tidy
go build -o pivot .
```

## MIT License

Copyright (c) 2026 zidan-herlangga

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
