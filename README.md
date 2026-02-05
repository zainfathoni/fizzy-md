# fizzy-md

[![Test](https://github.com/zainfathoni/fizzy-md/actions/workflows/test.yml/badge.svg)](https://github.com/zainfathoni/fizzy-md/actions/workflows/test.yml)
[![Release](https://github.com/zainfathoni/fizzy-md/actions/workflows/release.yml/badge.svg)](https://github.com/zainfathoni/fizzy-md/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Transparent Markdown→HTML wrapper for [Fizzy CLI](https://github.com/robzolkos/fizzy-cli).

## Problem

AI agents naturally write in Markdown (`**bold**`, `## Headers`), but Fizzy CLI requires HTML (`<strong>bold</strong>`, `<h2>Headers</h2>`). This creates friction — agents either:

- Forget and use Markdown (poor rendering in Fizzy UI)
- Manually convert to HTML (slows them down, breaks flow)
- Avoid rich formatting entirely (less informative cards)

## Solution

`fizzy-md` is a transparent wrapper that accepts Markdown and converts it to HTML automatically. All `fizzy` commands work exactly as before — just write naturally!

## Features

✅ **Zero config** — just works when installed  
✅ **100% backward compatible** — all fizzy commands pass through  
✅ **Fast** — <100ms overhead, compiled Go binary  
✅ **Smart file detection** — `.md` files auto-convert, `.html` files pass through  
✅ **CommonMark compliant** — powered by [goldmark](https://github.com/yuin/goldmark)

## Installation

### Option 1: Homebrew (macOS/Linux - Recommended)

```bash
brew tap zainfathoni/fizzy
brew install fizzy-md
```

### Option 2: Pre-built binaries

Download the latest release for your platform from the [releases page](https://github.com/zainfathoni/fizzy-md/releases).

**macOS:**
```bash
# Intel
curl -L https://github.com/zainfathoni/fizzy-md/releases/latest/download/fizzy-md_Darwin_x86_64.tar.gz | tar xz
sudo mv fizzy-md /usr/local/bin/

# Apple Silicon (M1/M2/M3)
curl -L https://github.com/zainfathoni/fizzy-md/releases/latest/download/fizzy-md_Darwin_arm64.tar.gz | tar xz
sudo mv fizzy-md /usr/local/bin/
```

**Linux:**
```bash
# x86_64
curl -L https://github.com/zainfathoni/fizzy-md/releases/latest/download/fizzy-md_Linux_x86_64.tar.gz | tar xz
sudo mv fizzy-md /usr/local/bin/

# ARM64
curl -L https://github.com/zainfathoni/fizzy-md/releases/latest/download/fizzy-md_Linux_arm64.tar.gz | tar xz
sudo mv fizzy-md /usr/local/bin/
```

**Windows:**

Download `fizzy-md_Windows_x86_64.zip` from the [releases page](https://github.com/zainfathoni/fizzy-md/releases) and extract to a directory in your PATH.

### Option 3: `go install`

```bash
go install github.com/zainfathoni/fizzy-md@latest
```

Make sure `~/go/bin` is in your PATH.

### Option 4: Build from source

```bash
git clone https://github.com/zainfathoni/fizzy-md.git
cd fizzy-md
go build -o fizzy-md
sudo mv fizzy-md /usr/local/bin/
```

## Usage

fizzy-md is a **transparent wrapper** around fizzy CLI. You can use it in three ways:

### Option A: Direct Usage (Recommended for Scripts)

Use `fizzy-md` directly as a drop-in replacement for `fizzy`:

```bash
fizzy-md card create \
  --title "Test Card" \
  --description "## Overview

This is **important**.

- Item 1
- Item 2"
```

**Pros:**
- ✅ Explicit and clear
- ✅ No confusion about which tool is running
- ✅ Works reliably in scripts and automation

### Option B: Preprocessing (Cleanest for Automation)

Convert Markdown to HTML first, then pass to fizzy:

```bash
# Convert inline Markdown via stdin
DESC=$(echo -e "## Hello\n\n**Bold** text" | fizzy-md)
fizzy card create --title "My Card" --description "$DESC"

# Convert file via stdin
HTML=$(cat description.md | fizzy-md)
fizzy card create --title "My Card" --description "$HTML"

# Or store in variable first
MARKDOWN="## Overview

This is **important**.

- Item 1
- Item 2"

HTML=$(echo "$MARKDOWN" | fizzy-md)
fizzy card create --title "My Card" --description "$HTML"
```

**Pros:**
- ✅ Clear separation between conversion and CLI operations
- ✅ Easier to debug and compose with other tools
- ✅ Perfect for helper scripts and automation
- ✅ Can store/reuse HTML output

### Option C: Alias (Convenient for Interactive Use)

Make `fizzy` automatically use Markdown:

```bash
alias fizzy='fizzy-md'

# Now 'fizzy' supports Markdown
fizzy card create --title "Test" --description "## Hello"
```

**Pros:**
- ✅ Convenient for interactive terminal use
- ✅ No need to type `fizzy-md` every time

**Note:** This works because `fizzy-md` internally calls the real `fizzy` binary via `exec.LookPath()`, which finds the actual executable (not shell aliases). No recursion occurs.

### File-Based Conversion

Works with all three approaches:

```bash
# Create card.md with Markdown content
echo "## My Card\n\n- Item 1\n- Item 2" > card.md

# Direct usage
fizzy-md card create --title "My Card" --description_file card.md

# With alias
alias fizzy='fizzy-md'
fizzy card create --title "My Card" --description_file card.md
```

## Supported Flags

`fizzy-md` converts these flags automatically:

| Flag | Description |
|------|-------------|
| `--description "text"` | Inline card description |
| `--body "text"` | Inline comment body |
| `--description_file path` | Card description from file |
| `--body_file path` | Comment body from file |

**File detection:**
- `.md` files → convert to HTML
- `.html` files → pass through unchanged
- No extension → assume Markdown (agent-friendly default)

All other arguments pass through to `fizzy` unchanged.

## Examples

### Create a card with rich formatting

```bash
fizzy-md card create \
  --title "Bug Fix: Login Issue" \
  --description "## Problem

Users can't log in when password contains `&` character.

## Root Cause

- URL encoding not applied
- Special chars break form submission

## Fix

Added `encodeURIComponent()` to password field.

**Status:** ✅ Resolved"
```

### Add a comment with code block

```bash
fizzy-md comment create \
  --card 42 \
  --body "Suggested fix:

\`\`\`javascript
const password = encodeURIComponent(input.value);
\`\`\`

This handles all special characters correctly."
```

### Use a Markdown file

```bash
# notes.md
## Meeting Notes

- Discussed architecture
- Agreed on Go implementation
- Next: Build MVP

fizzy-md card create --title "Sprint Planning" --description_file notes.md
```

## How It Works

fizzy-md operates in two modes:

### Wrapper Mode (Default)
1. **Intercepts** your command-line arguments
2. **Detects** Markdown in `--description`, `--body`, and file flags
3. **Converts** Markdown → HTML using goldmark
4. **Passes through** to real `fizzy` CLI
5. **Preserves** all other args/flags exactly

No modifications to `fizzy-cli` needed. Just a transparent wrapper!

### Stdin Mode (Pipe Support)
When run with no arguments and piped input:
```bash
echo "## Hello" | fizzy-md
# Output: <h2>Hello</h2>
```

Perfect for preprocessing in scripts:
```bash
HTML=$(cat notes.md | fizzy-md)
fizzy card create --title "Notes" --description "$HTML"
```

## Requirements

- Go 1.23+ (for building)
- [fizzy-cli](https://github.com/robzolkos/fizzy-cli) installed and in PATH
- macOS or Linux (Windows untested but should work)

## Development

```bash
# Clone the repo
git clone https://github.com/zainfathoni/fizzy-md.git
cd fizzy-md

# Install dependencies
go mod download

# Build
go build -o fizzy-md

# Test
./fizzy-md card create --title "Test" --description "## Heading"
```

## Inspiration

This tool follows the "agent-first" design principle from Steve Yegge's [Software Survival 3.0](https://steve-yegge.medium.com/software-survival-3-0-97a2a6255f7b):

> "Agents always act like they're in a hurry, and if something appears to be failing for them, they will rapidly switch to trying workarounds. [...] Conversely, if you build the tool to their tastes, then agents will use the hell out of it."

**Design principle:** Minimize friction. Let agents write naturally (Markdown), handle the conversion transparently.

## License

MIT License - see [LICENSE](LICENSE) file.

## Credits

- Built with [goldmark](https://github.com/yuin/goldmark) - The best Markdown parser in Go
- Wraps [fizzy-cli](https://github.com/robzolkos/fizzy-cli) by Rob Zolkos
- Created for the [Autobots](https://github.com/zainfathoni) AI agent team

## Contributing

Contributions welcome! Please open an issue or PR.

**Future ideas:**
- Homebrew tap for easier installation
- Windows support testing
- Upstream integration into fizzy-cli (if Rob Zolkos is interested)

---

**Built with 🔧 by [Wheeljack](https://github.com/wheeljackz) for AI agents everywhere.**
