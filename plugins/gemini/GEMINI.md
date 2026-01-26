# Release Agent Team

This plugin provides a multi-agent team for automating software release workflows.

## Available Commands

- `/agent-team-release:release <version>` - Execute full release workflow
- `/agent-team-release:check` - Run validation checks
- `/agent-team-release:changelog <version>` - Generate changelog
- `/agent-team-release:version-next` - Suggest next version

## Agent Team

- PM - Product scope and versioning
- QA - Build, test, lint validation
- Documentation - Docs completeness
- Security - Vulnerability scanning
- Release - Git and CI validation
- Coordinator - Orchestrates all agents

## Supported Languages

- Go (build, test, golangci-lint, gofmt)
- TypeScript/JavaScript (ESLint, Prettier, npm test)

## Dependencies

- `atrelease` - Release automation CLI
- `schangelog` - Changelog generation
- `git` - Version control
- `gh` - GitHub CLI