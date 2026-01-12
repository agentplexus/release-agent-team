// Copyright 2025 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package report

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

const (
	// boxWidth is the inner width of the box (between the border characters).
	boxWidth = 76
)

// Renderer renders TeamStatusReport to various formats.
type Renderer struct {
	w io.Writer
}

// NewRenderer creates a new Renderer writing to w.
func NewRenderer(w io.Writer) *Renderer {
	return &Renderer{w: w}
}

// Render renders the report using the box template.
func (r *Renderer) Render(report *TeamStatusReport) error {
	return r.renderBox(report)
}

// renderBox renders the report in the box format.
func (r *Renderer) renderBox(report *TeamStatusReport) error {
	tmpl, err := template.New("report").Funcs(templateFuncs()).Parse(boxTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	return tmpl.Execute(r.w, report)
}

// templateFuncs returns the template function map.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"header":       header,
		"separator":    separator,
		"footer":       footer,
		"teamHeader":   teamHeader,
		"teamSep":      teamSep,
		"checkLine":    checkLine,
		"centerLine":   centerLine,
		"paddedLine":   paddedLine,
		"finalMessage": finalMessage,
	}
}

// header returns the top border of the box.
func header() string {
	return "╔" + strings.Repeat("═", boxWidth) + "╗"
}

// separator returns a separator line.
func separator() string {
	return "╠" + strings.Repeat("═", boxWidth) + "╣"
}

// footer returns the bottom border of the box.
func footer() string {
	return "╚" + strings.Repeat("═", boxWidth) + "╝"
}

// centerLine centers text within the box.
func centerLine(text string) string {
	// Account for emoji width (each emoji takes ~2 characters visually)
	visualLen := visualLength(text)
	padding := max(0, boxWidth-visualLen)
	left := padding / 2
	right := padding - left
	return "║" + strings.Repeat(" ", left) + text + strings.Repeat(" ", right) + "║"
}

// paddedLine left-aligns text with padding.
func paddedLine(text string) string {
	visualLen := visualLength(text)
	padding := max(0, boxWidth-visualLen-1) // -1 for left space
	return "║ " + text + strings.Repeat(" ", padding) + "║"
}

// teamHeader formats a team header line.
func teamHeader(team Team) string {
	// Format: "team-id (name)"
	text := fmt.Sprintf("%s (%s)", team.ID, team.Name)
	return paddedLine(text)
}

// teamSep returns a light separator for within teams.
func teamSep() string {
	return "╠" + strings.Repeat("═", boxWidth) + "╣"
}

// checkLine formats a single check result line.
func checkLine(check Check) string {
	// Format: "  id                     🟢 GO    detail"
	// ID is left-aligned in 24 chars, status icon + text in ~10, detail fills rest

	id := check.ID
	if len(id) > 24 {
		id = id[:21] + "..."
	}

	icon := check.Status.Icon()
	statusText := string(check.Status)

	detail := check.Detail
	// Calculate available space for detail
	// Layout: "  " + id(24) + " " + icon(~2) + " " + status(5) + "  " + detail
	// Total fixed: 2 + 24 + 1 + 2 + 1 + 5 + 2 = 37
	maxDetail := boxWidth - 40 // Leave some margin
	if len(detail) > maxDetail {
		detail = detail[:maxDetail-3] + "..."
	}

	line := fmt.Sprintf("  %-24s %s %-5s %s", id, icon, statusText, detail)
	return paddedLine(line)
}

// finalMessage formats the final status message line.
func finalMessage(report *TeamStatusReport) string {
	return centerLine(report.FinalMessage())
}

// visualLength calculates the visual length of a string,
// accounting for emoji characters that take 2 columns.
func visualLength(s string) int {
	length := 0
	for _, r := range s {
		// Emoji and other wide characters
		if r >= 0x1F300 && r <= 0x1FAFF {
			length += 2
		} else if r >= 0x2600 && r <= 0x27BF {
			length += 2
		} else {
			length++
		}
	}
	return length
}

// boxTemplate is the text/template for the box format.
const boxTemplate = `{{ header }}
{{ centerLine "TEAM STATUS REPORT" }}
{{ separator }}
{{ paddedLine (printf "Project: %s" .Project) }}
{{ paddedLine (printf "Target:  %s" .Target) }}
{{ separator }}
{{ paddedLine .Phase }}
{{- range .Teams }}
{{ separator }}
{{ teamHeader . }}
{{- range .Checks }}
{{ checkLine . }}
{{- end }}
{{- end }}
{{ separator }}
{{ finalMessage . }}
{{ footer }}
`
