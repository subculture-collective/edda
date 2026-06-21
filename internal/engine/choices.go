package engine

import (
	"regexp"
	"strings"
)

const rightSingleQuote = "\u2019"

var numberedChoicePattern = regexp.MustCompile(`^\s*(\d+)[.)]\s+(.*\S)\s*$`)

func hasChoicesMarker(narrative string) bool {
	lower := strings.ToLower(strings.ReplaceAll(narrative, rightSingleQuote, "'"))
	markers := []string{
		"**choices:**",
		"choices:",
		"options:",
		"what do you do?",
		"what will you do?",
	}
	for _, marker := range markers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func extractChoicesStrict(narrative string) (string, []Choice, error) {
	cleaned, choices := extractChoices(narrative)
	if len(choices) == 0 && hasChoicesMarker(narrative) {
		return stripDanglingChoicesMarker(narrative), nil, nil
	}
	return cleaned, choices, nil
}

// CleanNarrativeText removes backend-generated choice scaffolding that should
// not be stored or exposed as narrative prose. It is intentionally tolerant so
// older session logs with dangling markers can still be cleaned at API edges.
func CleanNarrativeText(narrative string) string {
	cleaned, _, err := extractChoicesStrict(narrative)
	if err != nil {
		return stripDanglingChoicesMarker(narrative)
	}
	return cleaned
}

func stripDanglingChoicesMarker(narrative string) string {
	lines := strings.Split(narrative, "\n")
	for len(lines) > 0 {
		trimmed := strings.TrimSpace(lines[len(lines)-1])
		if trimmed == "" || isChoiceMarkerLine(trimmed) {
			lines = lines[:len(lines)-1]
			continue
		}
		break
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func isChoiceMarkerLine(line string) bool {
	lower := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(line), rightSingleQuote, "'"))
	switch lower {
	case "**choices:**", "choices:", "options:", "what do you do?", "what will you do?":
		return true
	default:
		return false
	}
}

func extractChoices(narrative string) (string, []Choice) {
	lines := strings.Split(narrative, "\n")
	if len(lines) == 0 {
		return narrative, nil
	}

	var (
		narrativeLines []string
		choices        []Choice
		inChoices      bool
	)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" && !inChoices {
			narrativeLines = append(narrativeLines, line)
			continue
		}

		if matches := numberedChoicePattern.FindStringSubmatch(line); matches != nil {
			inChoices = true
			id := matches[1]
			text := strings.TrimSpace(matches[2])
			choices = append(choices, Choice{ID: id, Text: text})
			continue
		}

		if inChoices {
			lower := strings.ToLower(strings.ReplaceAll(trimmed, rightSingleQuote, "'"))
			if strings.HasPrefix(lower, "or describe what you'd like to do") {
				continue
			}
			if trimmed == "" {
				continue
			}
			narrativeLines = append(narrativeLines, line)
			inChoices = false
			continue
		}

		narrativeLines = append(narrativeLines, line)
	}

	cleaned := strings.TrimSpace(strings.Join(narrativeLines, "\n"))
	if len(choices) == 0 {
		return narrative, nil
	}
	return cleaned, choices
}
