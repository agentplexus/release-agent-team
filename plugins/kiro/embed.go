// Package kiro provides embedded Kiro CLI agent configurations.
package kiro

import "embed"

// AgentFiles contains embedded Kiro agent JSON files.
//
//go:embed agents/*.json
var AgentFiles embed.FS

// SteeringFiles contains embedded Kiro steering markdown files.
//
//go:embed steering/*.md
var SteeringFiles embed.FS
