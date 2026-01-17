// Command releaseagent is an autonomous release preparation tool.
//
// It auto-detects languages based on project files (go.mod, package.json, etc.)
// and runs appropriate checks for each language found. It can also generate
// changelogs, update documentation, and manage the release lifecycle.
//
// Usage:
//
//	releaseagent                    # Run validation checks (default)
//	release-agent-team check              # Run validation checks
//	release-agent-team check /path/to/repo
//	release-agent-team check --verbose    # Show detailed output
//	release-agent-team check --no-test    # Skip tests
//	release-agent-team version            # Show version information
//
// Configuration:
//
// Create a .releaseagent.yaml file to customize behavior:
//
//	verbose: true
//	languages:
//	  go:
//	    lint: true
//	    test: true
//	    coverage: true
//	    exclude_coverage: "cmd"
//	  typescript:
//	    enabled: true
//	    paths: ["frontend/"]
package main

func main() {
	Execute()
}
