# pivot

> Switch Python, PHP, Node.js, and Go versions on the fly.

## Install

**Linux / macOS**

```bash
curl -fsSL https://raw.githubusercontent.com/zidan-herlangga/pivot/main/scripts/install.sh | sh
```

**Windows (PowerShell)**

```powershell
iwr -Uri https://raw.githubusercontent.com/zidan-herlangga/pivot/main/scripts/install.ps1 | iex
```

**Manual (Go required)**

```bash
go install github.com/zidan-herlangga/pivot@latest
```

## Quick Start

```bash
pivot                       # Interactive menu
pivot use node 22           # Switch Node.js to v22
pivot install python        # Download latest Python
pivot list                  # Show all available versions
pivot env                   # Print PATH setup for your shell
```

Add `~/.pivot/bin` to your PATH, then run `pivot use <runtime> <version>` to switch anytime.

## Commands

| Command                    | Description                       |
| -------------------------- | --------------------------------- |
| `pivot`                    | Interactive TUI menu              |
| `pivot list`               | Show installed versions           |
| `pivot use <rt> <ver>`     | Activate a version                |
| `pivot install <rt> [ver]` | Download a portable runtime       |
| `pivot run <rt> <ver> <cmd>` | Run a command with a version    |
| `pivot doctor`             | Diagnose system and PATH          |
| `pivot upgrade`            | Upgrade pivot to latest version   |
| `pivot clean`              | Remove unused runtime versions    |
| `pivot update`             | Check latest upstream versions    |
| `pivot init`               | Create `.pivotrc` in current dir  |
| `pivot env`                | Print PATH setup for shell config |

Version aliases: `latest`, `system`, `lts` (e.g. `pivot use node lts`).

### Project Scaffolding

```bash
pivot create laravel myapp
pivot create react myapp
pivot create nextjs myapp
```

Supported: Laravel, CodeIgniter 4, Symfony, WordPress (Bedrock), React, Next.js, Vue, AdonisJS.

### Profiles

```bash
pivot profile save backend     # Save current versions
pivot profile load backend     # Restore saved versions
pivot profile list             # List all profiles
pivot profile delete backend   # Delete a profile
```

## .pivotrc

Place a `.pivotrc` file in any directory to auto-apply runtimes when you `cd` into it:

```ini
python=3.12.0
node=22.0.0
```

pivot walks up the directory tree and picks the nearest `.pivotrc`.

## Language

pivot automatically detects your system language — **English** and **Bahasa Indonesia** are supported.

## Directory Structure

```
~/.pivot/
├── bin/            # Symlinks to active runtime binaries (add to PATH)
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
