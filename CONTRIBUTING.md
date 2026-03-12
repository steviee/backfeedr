# Contributing to backfeedr

First off, thank you for considering contributing to backfeedr! 👋

## Getting Started

### Setting up your environment

1. Fork the repository on GitHub
2. Clone your fork locally
3. Ensure you have Go 1.22+ installed
4. Run `make build` to verify everything works

### Running tests

```bash
make test
```

## How to Contribute

We use [git-issues](https://github.com/steviee/git-issues) for task tracking. Each issue is a markdown file in `.issues/`.

### Finding work

1. Check open issues in `.issues/`
2. Comment on an issue if you want to work on it
3. We'll assign it to you

### Making changes

1. Create a branch: `git checkout -b feature/description`
2. Make your changes
3. Update the relevant issue file (set status to `in-progress`)
4. Commit with clear messages
5. Push and open a Pull Request

### Pull Request Process

1. Update the README.md if needed
2. Ensure tests pass
3. Update relevant issue to `closed`
4. Request review from maintainers
5. Address feedback
6. Squash merge when approved

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Keep functions focused and small
- Add comments for exported functions
- Write tests for new functionality

## Communication

- Be respectful and constructive
- Ask questions if something is unclear
- Share ideas in issues before big changes
- We're all learning together!

## Questions?

Open an issue or reach out. We're happy to help.
