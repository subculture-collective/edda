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
		return "", nil, ErrUnparseableChoices
	}
	return cleaned, choices, nil
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
