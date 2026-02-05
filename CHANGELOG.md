# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Markdown table support via goldmark Table extension
- Test coverage for tables (simple, alignment, inline formatting)
- Total test count: 29 (was 26)

### Fixed
- Tables now render as proper HTML `<table>` elements instead of plain text

## [0.1.0] - 2026-02-05

### Added
- Initial release of fizzy-md! 🎉
- Transparent Markdown→HTML conversion for Fizzy CLI
- Support for `--description`, `--body`, `--description_file`, `--body_file` flags
- Smart file detection (`.md` → convert, `.html` → passthrough, no extension → assume Markdown)
- Version flag support (`--version` / `-v`)
- Comprehensive test suite (26 unit tests)
- Benchmark tests (performance: ~10μs conversion time)
- Multi-platform binaries (macOS, Linux, Windows - amd64 + arm64)
- Homebrew tap support (`brew tap zainfathoni/fizzy`)
- GitHub Actions workflows (test + release)
- GoReleaser configuration for automated releases

### Features
- **Markdown parser:** goldmark (CommonMark 0.31.2 compliant)
- **Performance:** 10,000x faster than the 100ms requirement
- **Platforms:** macOS (Intel + Apple Silicon), Linux (x86_64 + ARM64), Windows (x86_64 + ARM64)
- **Installation:** Homebrew, pre-built binaries, `go install`, build from source

### Tested
- All Markdown elements (headers, lists, code blocks, links, blockquotes, etc.)
- Special character escaping (`&`, `<`, `>`, quotes)
- Emoji preservation (✅ 🚀 🔧)
- Complex nested structures
- Edge cases (empty lines, line breaks)
- UI rendering verification in Fizzy
- Real-world usage (dogfooded on Fizzy Card #92)

[0.1.0]: https://github.com/zainfathoni/fizzy-md/releases/tag/v0.1.0
