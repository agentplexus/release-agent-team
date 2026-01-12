// Package checks provides pre-push checks for various languages.
package checks

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Result represents the result of a check.
type Result struct {
	Name    string
	Passed  bool
	Output  string
	Error   error
	Skipped bool
	Reason  string
	Warning bool // Soft check: reported but doesn't fail the build
}

// Checker is the interface for language-specific checks.
type Checker interface {
	Name() string
	Check(dir string, opts Options) []Result
}

// Options configures which checks to run.
type Options struct {
	Test     bool
	Lint     bool
	Format   bool
	Coverage bool
	Verbose  bool

	// Language-specific options
	GoExcludeCoverage string // directories to exclude from coverage (e.g., "cmd")
}

// DefaultOptions returns the default check options.
func DefaultOptions() Options {
	return Options{
		Test:              true,
		Lint:              true,
		Format:            true,
		Coverage:          false,
		Verbose:           false,
		GoExcludeCoverage: "cmd",
	}
}

// RunCommand executes a command and returns the result.
func RunCommand(name string, dir string, command string, args ...string) Result {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()

	return Result{
		Name:   name,
		Passed: err == nil,
		Output: strings.TrimSpace(string(output)),
		Error:  err,
	}
}

// CommandExists checks if a command is available in PATH.
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// PrintResults prints check results to stdout.
// Returns counts: passed, failed, skipped, warnings
func PrintResults(results []Result, verbose bool) (passed int, failed int, skipped int, warnings int) {
	for _, r := range results {
		if r.Skipped {
			fmt.Printf("⊘ %s (skipped: %s)\n", r.Name, r.Reason)
			skipped++
			continue
		}

		if r.Warning {
			// Soft check: show warning but count as passed
			if r.Passed {
				fmt.Printf("✓ %s\n", r.Name)
			} else {
				fmt.Printf("⚠ %s (warning)\n", r.Name)
				warnings++
			}
			// Always show output for warnings
			if r.Output != "" {
				lines := strings.Split(r.Output, "\n")
				for _, line := range lines {
					fmt.Printf("  %s\n", line)
				}
			}
			if r.Passed {
				passed++
			}
			continue
		}

		if r.Passed {
			fmt.Printf("✓ %s\n", r.Name)
			passed++
		} else {
			fmt.Printf("✗ %s\n", r.Name)
			failed++
		}

		if verbose || !r.Passed {
			if r.Output != "" {
				// Indent output
				lines := strings.Split(r.Output, "\n")
				for _, line := range lines {
					fmt.Printf("  %s\n", line)
				}
			}
			if r.Error != nil && r.Output == "" {
				fmt.Printf("  Error: %v\n", r.Error)
			}
		}
	}

	return passed, failed, skipped, warnings
}

// FileExists checks if a file exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ValidationStatus represents a Go/No-Go status for a check.
type ValidationStatus struct {
	Name   string
	Status string // "GO" or "NO-GO"
	Icon   string // UTF-8 icon
	Detail string // Optional detail message
}

// GoNoGoIcons defines the UTF-8 icons for Go/No-Go status.
const (
	IconGo      = "🟢" // Green circle for GO
	IconNoGo    = "🔴" // Red circle for NO-GO
	IconSkipped = "⚪" // White circle for skipped
	IconWarning = "🟡" // Yellow circle for warning
)

// PrintGoNoGoReport prints results in NASA-style Go/No-Go format.
// Returns true if all required checks pass (GO), false otherwise (NO-GO).
func PrintGoNoGoReport(results []Result, verbose bool) bool {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              RELEASE VALIDATION - GO/NO-GO                   ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")

	allGo := true
	var statuses []ValidationStatus

	for _, r := range results {
		var status ValidationStatus
		status.Name = r.Name

		if r.Skipped {
			status.Status = "SKIP"
			status.Icon = IconSkipped
			status.Detail = r.Reason
		} else if r.Warning && !r.Passed {
			status.Status = "WARN"
			status.Icon = IconWarning
			status.Detail = "Warning (non-blocking)"
		} else if r.Passed {
			status.Status = "GO"
			status.Icon = IconGo
		} else {
			status.Status = "NO-GO"
			status.Icon = IconNoGo
			allGo = false
			if r.Output != "" {
				// Truncate long output for summary
				lines := strings.Split(r.Output, "\n")
				if len(lines) > 0 {
					status.Detail = lines[0]
					if len(status.Detail) > 40 {
						status.Detail = status.Detail[:40] + "..."
					}
				}
			}
		}

		statuses = append(statuses, status)
	}

	// Print each status
	for _, s := range statuses {
		// Format: Icon STATUS Name
		line := fmt.Sprintf("║ %s %-6s %-50s ║", s.Icon, s.Status, s.Name)
		fmt.Println(line)
		if s.Detail != "" && verbose {
			detailLine := fmt.Sprintf("║          └─ %-49s ║", s.Detail)
			fmt.Println(detailLine)
		}
	}

	fmt.Println("╠══════════════════════════════════════════════════════════════╣")

	// Final verdict
	if allGo {
		fmt.Println("║                    🚀 ALL SYSTEMS GO 🚀                      ║")
	} else {
		fmt.Println("║                  🛑 NO-GO FOR RELEASE 🛑                     ║")
	}
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	return allGo
}

// PrintCompactGoNoGo prints a compact Go/No-Go summary.
func PrintCompactGoNoGo(results []Result) bool {
	allGo := true
	var goCount, noGoCount, warnCount, skipCount int

	for _, r := range results {
		switch {
		case r.Skipped:
			skipCount++
		case r.Warning && !r.Passed:
			warnCount++
		case r.Passed:
			goCount++
		default:
			noGoCount++
			allGo = false
		}
	}

	fmt.Printf("\n━━━ Validation Summary ━━━\n")
	fmt.Printf("%s GO: %d  %s NO-GO: %d  %s WARN: %d  %s SKIP: %d\n",
		IconGo, goCount, IconNoGo, noGoCount, IconWarning, warnCount, IconSkipped, skipCount)

	if allGo {
		fmt.Printf("\n%s RELEASE VALIDATION: GO\n", IconGo)
	} else {
		fmt.Printf("\n%s RELEASE VALIDATION: NO-GO\n", IconNoGo)
	}

	return allGo
}
