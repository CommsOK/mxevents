# Contributing to mxevents

Thank you for your interest in contributing to mxevents! This document provides guidelines and instructions for contributing.

## How to Contribute

### Reporting Bugs

Before creating a bug report, please check existing issues to avoid duplicates. When creating a bug report, include:

- A clear, descriptive title
- Steps to reproduce the issue
- Expected behavior vs actual behavior
- Go version and operating system
- Relevant code snippets or error messages

### Suggesting Features

Feature suggestions are welcome! Please open an issue with:

- A clear description of the feature
- The problem it solves or use case it addresses
- Any implementation ideas you have

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Write tests** for any new functionality
3. **Ensure tests pass** by running `go test ./...`
4. **Follow Go conventions** - run `go fmt` and `go vet`
5. **Update documentation** if you're changing behavior
6. **Write clear commit messages** following conventional commits

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/mxevents.git
cd mxevents

# Add upstream remote
git remote add upstream https://github.com/commsok/mxevents.git

# Install dependencies
go mod download

# Run tests
go test ./...
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Code Style

This project follows standard Go conventions:

- Run `go fmt ./...` before committing
- Run `go vet ./...` to catch common issues
- Use meaningful variable and function names
- Add comments for exported types and functions

## Project Structure

```
mxevents/
├── events.go           # Core event types
├── reasons.go          # Reason taxonomy
├── facts.go            # EventFacts and ClassificationResult
├── classifier.go       # Classifier interface
├── enricher.go         # Enricher interface
├── vendors.go          # Vendor constants
├── classifier/         # Classifier implementations
│   └── bounce/         # Bounce classification
├── enrichment/         # Enricher implementations
└── toolkit/            # High-level API
```

## Making Changes to the Taxonomy

The event taxonomy is versioned. When making changes:

1. **Additive changes** (new events, new reasons) are generally safe
2. **Breaking changes** (removing or renaming) require a version bump
3. Update `CurrentTaxonomyVersion` in `events.go` for breaking changes
4. Document changes in the changelog

## Questions?

Feel free to open an issue for any questions about contributing.
