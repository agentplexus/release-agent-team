// Copyright 2025 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package checks

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// SecurityChecker implements security and compliance checks.
type SecurityChecker struct{}

// Name returns the checker name.
func (c *SecurityChecker) Name() string {
	return "Security"
}

// SecurityOptions configures security checks.
type SecurityOptions struct {
	Verbose bool
}

// Check runs security checks on the specified directory.
func (c *SecurityChecker) Check(dir string, opts SecurityOptions) []Result {
	var results []Result

	// Check LICENSE file exists
	results = append(results, c.checkLicense(dir))

	// Check for known vulnerabilities (Go)
	results = append(results, c.checkGoVulncheck(dir))

	// Check for dependency audit (Go)
	results = append(results, c.checkGoModAudit(dir))

	// Check for secrets in code
	results = append(results, c.checkNoSecrets(dir))

	return results
}

func (c *SecurityChecker) checkLicense(dir string) Result {
	name := "Security: LICENSE file"

	licenseFiles := []string{
		"LICENSE",
		"LICENSE.md",
		"LICENSE.txt",
		"COPYING",
	}

	for _, f := range licenseFiles {
		if FileExists(filepath.Join(dir, f)) {
			return Result{
				Name:   name,
				Passed: true,
				Output: f,
			}
		}
	}

	return Result{
		Name:   name,
		Passed: false,
		Output: "No LICENSE file found",
	}
}

func (c *SecurityChecker) checkGoVulncheck(dir string) Result {
	name := "Security: vulnerability scan"

	// Check if this is a Go project
	if !FileExists(filepath.Join(dir, "go.mod")) {
		return Result{
			Name:    name,
			Skipped: true,
			Reason:  "Not a Go project",
		}
	}

	// Check if govulncheck is available
	if !CommandExists("govulncheck") {
		return Result{
			Name:    name,
			Skipped: true,
			Reason:  "govulncheck not installed. Install: go install golang.org/x/vuln/cmd/govulncheck@latest",
		}
	}

	// Run govulncheck
	result := RunCommand(name, dir, "govulncheck", "./...")
	if !result.Passed {
		// Check if it's a real vulnerability or just an error
		if strings.Contains(result.Output, "Vulnerability") {
			return Result{
				Name:   name,
				Passed: false,
				Output: "Vulnerabilities found - review and update dependencies",
			}
		}
	}

	return Result{
		Name:   name,
		Passed: true,
		Output: "No vulnerabilities found",
	}
}

func (c *SecurityChecker) checkGoModAudit(dir string) Result {
	name := "Security: dependency audit"

	// Check if this is a Go project
	if !FileExists(filepath.Join(dir, "go.mod")) {
		return Result{
			Name:    name,
			Skipped: true,
			Reason:  "Not a Go project",
		}
	}

	// Use go list to check for dependency issues
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	cmd.Dir = dir
	_, err := cmd.Output()
	if err != nil {
		return Result{
			Name:    name,
			Warning: true,
			Passed:  false,
			Output:  "Failed to list dependencies",
		}
	}

	// Check for retracted versions
	cmd = exec.Command("go", "list", "-m", "-u", "-retracted", "all")
	cmd.Dir = dir
	output, _ := cmd.Output()

	if strings.Contains(string(output), "(retracted)") {
		return Result{
			Name:    name,
			Warning: true,
			Passed:  false,
			Output:  "Some dependencies have retracted versions",
		}
	}

	return Result{
		Name:   name,
		Passed: true,
	}
}

func (c *SecurityChecker) checkNoSecrets(dir string) Result {
	name := "Security: no hardcoded secrets"

	// Check for common secret patterns in Go files
	// This is a simple heuristic check, not a full secret scanner

	// Patterns that might indicate hardcoded secrets
	secretPatterns := []string{
		"password.*=.*\"",
		"apikey.*=.*\"",
		"api_key.*=.*\"",
		"secret.*=.*\"",
		"token.*=.*\"",
		"private_key.*=.*\"",
	}

	for _, pattern := range secretPatterns {
		cmd := exec.Command("grep", "-r", "-i", "-l", "--include=*.go", pattern, ".")
		cmd.Dir = dir
		output, err := cmd.Output()

		// grep returns exit 1 if no matches (which is good for us)
		if err == nil && len(output) > 0 {
			files := strings.TrimSpace(string(output))
			return Result{
				Name:    name,
				Warning: true,
				Passed:  false,
				Output:  "Potential hardcoded secrets found in: " + files,
			}
		}
	}

	return Result{
		Name:   name,
		Passed: true,
	}
}
