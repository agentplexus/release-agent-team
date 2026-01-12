// Command releaseagent is an autonomous release preparation tool.
//
// It auto-detects languages based on project files (go.mod, package.json, etc.)
// and runs appropriate checks for each language found. It can also generate
// changelogs, update documentation, and manage the release lifecycle.
//
// Usage:
//
//	releaseagent                    # Run validation checks (default)
//	releaseagent check              # Run validation checks
//	releaseagent check /path/to/repo
//	releaseagent check --verbose    # Show detailed output
//	releaseagent check --no-test    # Skip tests
//	releaseagent version            # Show version information
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
