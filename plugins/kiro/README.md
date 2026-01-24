# release-agent-team - Kiro CLI Plugin

Multi-agent team for automating software release workflows including versioning, changelog generation, CI verification, security scanning, and Git tagging

## Agents

| Agent | Description |
|-------|-------------|
| `release-coordinator` | Orchestrates software releases including semantic versioning, changelog generation, CI verification, and Git tagging. Use when preparing a new release, automating release workflows, or managing version bumps. |

## Steering Files

Copy steering files to `.kiro/steering/` for automatic context loading:

```bash
mkdir -p .kiro/steering
cp steering/*.md .kiro/steering/
```

