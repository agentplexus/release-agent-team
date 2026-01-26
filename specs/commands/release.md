---
name: release
description: Execute full release workflow for the specified version
arguments: [version]
dependencies: [atrelease, schangelog, git]
process:
  - Validate version format and check it doesn't exist
  - Check working directory is clean
  - Run validation checks (build, test, lint, format)
  - Generate changelog via schangelog
  - Update roadmap via sroadmap
  - Create release commit
  - Push to remote
  - Wait for CI to pass
  - Create and push release tag
---

Execute the complete release workflow for the specified version.

This command runs the full release process including validation, changelog generation, and git tagging.

## Usage

```
/agent-team-release:release <version>
```

## Arguments

- **version** (required): Semantic version for the release (e.g., v1.2.3)

## Examples

Create a release:

```
/agent-team-release:release v0.9.0
```

Executes: `atrelease release v0.9.0 --verbose`

Dry run:

```
atrelease release v1.0.0 --dry-run
```

Preview changes without executing.
