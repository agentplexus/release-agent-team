// Copyright 2025 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DocChecker implements documentation checks.
type DocChecker struct{}

// Name returns the checker name.
func (c *DocChecker) Name() string {
	return "Documentation"
}

// DocOptions configures documentation checks.
type DocOptions struct {
	Version string // Target release version (e.g., "v0.2.0")
	Verbose bool
}

// Check runs documentation checks on the specified directory.
func (c *DocChecker) Check(dir string, opts DocOptions) []Result {
	var results []Result

	// Check README.md exists
	results = append(results, c.checkReadme(dir))

	// Check PRD.md (if exists, verify it's not empty)
	results = append(results, c.checkOptionalDoc(dir, "PRD.md", "Product Requirements Document"))

	// Check TRD.md (if exists, verify it's not empty)
	results = append(results, c.checkOptionalDoc(dir, "TRD.md", "Technical Requirements Document"))

	// Check for MkDocs site
	results = append(results, c.checkMkDocs(dir))

	// Check for release notes
	results = append(results, c.checkReleaseNotes(dir, opts.Version))

	// Check CHANGELOG.md exists
	results = append(results, c.checkChangelog(dir))

	return results
}

func (c *DocChecker) checkReadme(dir string) Result {
	name := "Docs: README.md"
	readmePath := filepath.Join(dir, "README.md")

	if !FileExists(readmePath) {
		return Result{
			Name:   name,
			Passed: false,
			Output: "README.md not found - create one to document the project",
		}
	}

	// Check if README has reasonable content (not empty or just a title)
	info, err := os.Stat(readmePath)
	if err != nil {
		return Result{
			Name:   name,
			Passed: false,
			Error:  err,
		}
	}

	if info.Size() < 100 {
		return Result{
			Name:   name,
			Passed: false,
			Output: "README.md exists but appears too short (< 100 bytes)",
		}
	}

	return Result{
		Name:   name,
		Passed: true,
	}
}

func (c *DocChecker) checkOptionalDoc(dir, filename, description string) Result {
	name := fmt.Sprintf("Docs: %s", filename)
	docPath := filepath.Join(dir, filename)

	if !FileExists(docPath) {
		return Result{
			Name:    name,
			Skipped: true,
			Reason:  fmt.Sprintf("%s not found (optional)", description),
		}
	}

	// If it exists, check it has content
	info, err := os.Stat(docPath)
	if err != nil {
		return Result{
			Name:   name,
			Passed: false,
			Error:  err,
		}
	}

	if info.Size() < 50 {
		return Result{
			Name:    name,
			Warning: true,
			Passed:  false,
			Output:  fmt.Sprintf("%s exists but appears too short", filename),
		}
	}

	return Result{
		Name:   name,
		Passed: true,
	}
}

func (c *DocChecker) checkMkDocs(dir string) Result {
	name := "Docs: MkDocs site"

	// Check for docs directory
	docsPath := filepath.Join(dir, "docs")
	if !FileExists(docsPath) {
		return Result{
			Name:    name,
			Skipped: true,
			Reason:  "docs/ directory not found (optional)",
		}
	}

	// Check for mkdocs.yml
	mkdocsPath := filepath.Join(docsPath, "mkdocs.yml")
	if !FileExists(mkdocsPath) {
		mkdocsPath = filepath.Join(dir, "mkdocs.yml")
		if !FileExists(mkdocsPath) {
			return Result{
				Name:    name,
				Warning: true,
				Passed:  false,
				Output:  "docs/ exists but mkdocs.yml not found",
			}
		}
	}

	return Result{
		Name:   name,
		Passed: true,
	}
}

func (c *DocChecker) checkReleaseNotes(dir string, version string) Result {
	name := "Docs: Release Notes"

	if version == "" {
		return Result{
			Name:    name,
			Skipped: true,
			Reason:  "No version specified for release notes check",
		}
	}

	// Clean version string
	ver := strings.TrimPrefix(version, "v")
	versionWithV := "v" + ver

	// Check for docs/releases/vX.Y.Z.md first
	docsPath := filepath.Join(dir, "docs", "releases", versionWithV+".md")
	if FileExists(docsPath) {
		return Result{
			Name:   name,
			Passed: true,
			Output: fmt.Sprintf("Found: docs/releases/%s.md", versionWithV),
		}
	}

	// Check for ./RELEASE_NOTES_vX.Y.Z.md
	releaseNotesPath := filepath.Join(dir, fmt.Sprintf("RELEASE_NOTES_%s.md", versionWithV))
	if FileExists(releaseNotesPath) {
		return Result{
			Name:   name,
			Passed: true,
			Output: fmt.Sprintf("Found: RELEASE_NOTES_%s.md", versionWithV),
		}
	}

	// Determine expected location
	docsDir := filepath.Join(dir, "docs")
	var expectedPath string
	if FileExists(docsDir) {
		expectedPath = fmt.Sprintf("docs/releases/%s.md", versionWithV)
	} else {
		expectedPath = fmt.Sprintf("RELEASE_NOTES_%s.md", versionWithV)
	}

	return Result{
		Name:   name,
		Passed: false,
		Output: fmt.Sprintf("Release notes not found. Expected: %s", expectedPath),
	}
}

func (c *DocChecker) checkChangelog(dir string) Result {
	name := "Docs: CHANGELOG.md"
	changelogPath := filepath.Join(dir, "CHANGELOG.md")

	if !FileExists(changelogPath) {
		// Check for CHANGELOG.json as alternative
		changelogJSON := filepath.Join(dir, "CHANGELOG.json")
		if FileExists(changelogJSON) {
			return Result{
				Name:    name,
				Warning: true,
				Passed:  false,
				Output:  "CHANGELOG.json exists but CHANGELOG.md not generated. Run: schangelog generate",
			}
		}
		return Result{
			Name:   name,
			Passed: false,
			Output: "CHANGELOG.md not found",
		}
	}

	return Result{
		Name:   name,
		Passed: true,
	}
}
