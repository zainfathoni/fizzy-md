# Contributing to fizzy-md

Thanks for your interest in contributing! 🎉

## Quick Start

1. **Fork** the repository
2. **Clone** your fork:
   ```bash
   git clone https://github.com/YOUR-USERNAME/fizzy-md.git
   cd fizzy-md
   ```
3. **Install dependencies:**
   ```bash
   go mod download
   ```
4. **Run tests:**
   ```bash
   go test -v ./...
   ```

## Development

### Running Tests

```bash
# All tests
go test -v ./...

# With coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Benchmarks
go test -bench=. -benchmem ./...
```

### Building

```bash
# Development build
go build -o fizzy-md

# Build with version info
go build -ldflags "-X main.version=dev" -o fizzy-md

# Test the build
./fizzy-md --version
```

### Testing Manually

```bash
# Make sure fizzy-cli is installed
brew install robzolkos/fizzy/fizzy

# Test with inline Markdown
export PATH="/path/to/fizzy-cli:$PATH"
./fizzy-md card create --title "Test" --description "## Hello\n\n**Bold** text"

# Test with file
echo "## Test" > test.md
./fizzy-md card create --description_file test.md
```

## Pull Request Process

1. **Create a branch** for your feature:
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes** and add tests if needed

3. **Run tests** to make sure everything passes:
   ```bash
   go test -v ./...
   ```

4. **Commit** with a clear message:
   ```bash
   git commit -m "Add feature: description of what you did"
   ```

5. **Push** to your fork:
   ```bash
   git push origin feature/my-new-feature
   ```

6. **Open a Pull Request** on GitHub with:
   - Clear description of what you changed
   - Why the change is needed
   - Any relevant issue numbers

## Release Process

Releases are automated — maintainers only need to update the CHANGELOG:

1. **Update CHANGELOG.md** with the new version:
   ```markdown
   ## [X.Y.Z] - YYYY-MM-DD
   
   ### Added
   - New feature descriptions
   
   ### Changed
   - Changes to existing functionality
   
   ### Fixed
   - Bug fixes
   ```
   Follow [Keep a Changelog](https://keepachangelog.com/) format and [Semantic Versioning](https://semver.org/).

2. **Commit and push to master:**
   ```bash
   git add CHANGELOG.md
   git commit -m "Bump version to X.Y.Z"
   git push origin master
   ```

3. **Automation handles the rest:**
   - The `auto-tag` workflow detects the version change
   - Creates an annotated git tag (`vX.Y.Z`)
   - Pushes the tag to trigger the release workflow
   - Release workflow builds binaries and creates GitHub release
   - Homebrew tap gets automatically updated

**Manual fallback** (if automation fails):
```bash
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin vX.Y.Z
```

The release workflow will still run and create the GitHub release.

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Add tests for new functionality
- Keep functions focused and small
- Document exported functions

## Testing Guidelines

- Unit tests for all conversion logic
- Test both success and error cases
- Add benchmarks for performance-critical code
- Verify multi-platform compatibility when possible

## Reporting Issues

Found a bug? Have a feature request?

1. **Check existing issues** first
2. **Open a new issue** with:
   - Clear title and description
   - Steps to reproduce (for bugs)
   - Expected vs actual behavior
   - Your environment (OS, Go version, fizzy-md version)

## Feature Requests

We love new ideas! When suggesting a feature:

- Explain the problem it solves
- Show example usage
- Consider backward compatibility
- Think about impact on performance

## Code of Conduct

Be respectful, constructive, and kind. This is a community project built for agents and humans alike.

## Questions?

Open a [GitHub Discussion](https://github.com/zainfathoni/fizzy-md/discussions) or reach out to [@zainfathoni](https://github.com/zainfathoni).

---

**Happy coding!** 🔧
