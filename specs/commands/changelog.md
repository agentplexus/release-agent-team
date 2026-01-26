---
name: changelog
description: Generate changelog entries for the specified version
arguments: [version]
dependencies: [schangelog, git]
process:
  - Get current version from latest git tag
  - Parse commits since that tag using schangelog
  - Classify commits by type (feat, fix, docs, etc.)
  - Generate changelog entries in CHANGELOG.json
  - Regenerate CHANGELOG.md
---

Generate changelog entries for the specified version based on commits since the last tag.

Uses schangelog to parse conventional commits and generate structured changelog entries.

## Usage

```
/agent-team-release:changelog <version>
```

## Arguments

- **version** (required): Version to generate changelog for (e.g., v1.2.3)

## Examples

Generate changelog:

```
/agent-team-release:changelog v0.9.0
```

Updates CHANGELOG.json and CHANGELOG.md with entries for the specified version.
