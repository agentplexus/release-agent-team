---
name: version-next
description: Analyze commits and suggest next semantic version
dependencies: [git, schangelog]
process:
  - Get current version from latest git tag
  - Parse commits since that tag
  - Classify commits by type
  - Determine version bump (major/minor/patch)
  - Suggest next version with reasoning
---

Analyze git history since the last tag and suggest the next semantic version based on conventional commits.

Given `vMAJOR.MINOR.PATCH`:

- Breaking changes bump MAJOR
- New features (feat:) bump MINOR
- Bug fixes (fix:) bump PATCH

## Usage

```
/agent-team-release:version-next
```

## Examples

Analyze and suggest version:

```
/agent-team-release:version-next
```

Output: `Current: v0.8.0, Suggested: v0.9.0 (3 new features)`
