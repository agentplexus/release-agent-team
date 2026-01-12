// Copyright 2025 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package checks

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ReleaseChecker implements release management checks.
type ReleaseChecker struct{}

// Name returns the checker name.
func (c *ReleaseChecker) Name() string {
	return "Release"
}

// ReleaseOptions configures release checks.
type ReleaseOptions struct {
	Version string // Target release version (e.g., "v0.2.0")
	Verbose bool
}

// Check runs release management checks on the specified directory.
func (c *ReleaseChecker) Check(dir string, opts ReleaseOptions) []Result {
	var results []Result

	// Check version format and availability
	results = append(results, c.checkVersionAvailable(dir, opts.Version))

	// Check git status (clean working directory for release)
	results = append(results, c.checkGitStatus(dir))

	// Check git remote is configured
	results = append(results, c.checkGitRemote(dir))

	// Check CHANGELOG.json exists and is valid
	results = append(results, c.checkChangelogJSON(dir))

	// Check for CI configuration
	results = append(results, c.checkCIConfig(dir))

	return results
}

func (c *ReleaseChecker) checkVersionAvailable(dir string, version string) Result {
	name := "Release: version available"

	if version == "" {
		return Result{
			Name:    name,
			Skipped: true,
			Reason:  "No version specified",
		}
	}

	// Ensure version has v prefix
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	// Check if tag already exists
	cmd := exec.Command("git", "tag", "-l", version)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return Result{
			Name:   name,
			Passed: false,
			Error:  err,
		}
	}

	if strings.TrimSpace(string(output)) != "" {
		return Result{
			Name:   name,
			Passed: false,
			Output: fmt.Sprintf("Tag %s already exists", version),
		}
	}

	return Result{
		Name:   name,
		Passed: true,
		Output: fmt.Sprintf("Tag %s is available", version),
	}
}

func (c *ReleaseChecker) checkGitStatus(dir string) Result {
	name := "Release: git working directory"

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return Result{
			Name:   name,
			Passed: false,
			Error:  err,
		}
	}

	status := strings.TrimSpace(string(output))
	if status != "" {
		lines := strings.Split(status, "\n")
		summary := fmt.Sprintf("%d uncommitted changes", len(lines))
		if len(lines) <= 5 {
			summary = status
		}
		return Result{
			Name:    name,
			Warning: true,
			Passed:  false,
			Output:  summary,
		}
	}

	return Result{
		Name:   name,
		Passed: true,
		Output: "Working directory is clean",
	}
}

func (c *ReleaseChecker) checkGitRemote(dir string) Result {
	name := "Release: git remote"

	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return Result{
			Name:   name,
			Passed: false,
			Output: "No 'origin' remote configured",
		}
	}

	remote := strings.TrimSpace(string(output))
	return Result{
		Name:   name,
		Passed: true,
		Output: remote,
	}
}

func (c *ReleaseChecker) checkChangelogJSON(dir string) Result {
	name := "Release: CHANGELOG.json"
	changelogPath := filepath.Join(dir, "CHANGELOG.json")

	if !FileExists(changelogPath) {
		return Result{
			Name:   name,
			Passed: false,
			Output: "CHANGELOG.json not found. Run: schangelog init",
		}
	}

	// Validate with schangelog if available
	if CommandExists("schangelog") {
		result := RunCommand("validate", dir, "schangelog", "validate", "CHANGELOG.json")
		if !result.Passed {
			return Result{
				Name:   name,
				Passed: false,
				Output: "CHANGELOG.json validation failed",
			}
		}
	}

	return Result{
		Name:   name,
		Passed: true,
	}
}

func (c *ReleaseChecker) checkCIConfig(dir string) Result {
	name := "Release: CI configuration"

	// Check for common CI configurations
	ciConfigs := []string{
		".github/workflows",
		".gitlab-ci.yml",
		".circleci/config.yml",
		"Jenkinsfile",
		".travis.yml",
		"azure-pipelines.yml",
	}

	for _, config := range ciConfigs {
		configPath := filepath.Join(dir, config)
		if FileExists(configPath) {
			return Result{
				Name:   name,
				Passed: true,
				Output: fmt.Sprintf("Found: %s", config),
			}
		}
	}

	return Result{
		Name:    name,
		Warning: true,
		Passed:  false,
		Output:  "No CI configuration found",
	}
}
