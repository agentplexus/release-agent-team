---
name: check
description: Run validation checks on the current repository
dependencies: [atrelease]
process:
  - Detect project language(s)
  - Run language-specific build checks
  - Run test suite
  - Run linter
  - Check code formatting
  - Report pass/fail status for each check
---

Run all validation checks for the current repository without making any changes.

## Go Projects

- `go build ./...` - Compilation
- `go test ./...` - Unit tests
- `golangci-lint run` - Linting
- `gofmt` - Formatting
- `go mod tidy` - Dependencies

## TypeScript/JavaScript Projects

- `npm test` or `yarn test` - Tests
- `eslint` - Linting
- `prettier --check` - Formatting
- `tsc --noEmit` - Type checking

## Usage

```
/agent-team-release:check
```

## Examples

Run checks:

```
/agent-team-release:check
```

Executes: `atrelease check --verbose`
