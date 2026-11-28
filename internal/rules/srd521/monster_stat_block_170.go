package srd521

import (
	"fmt"
	"slices"
	"strings"
)

// validateSourceNote accepts a retained source phrase only without surrounding
// space and on a record that requires the mechanic that would have to apply it.
func (m Monster) validateSourceNote(note, key string) error {
	if note != strings.TrimSpace(note) {
		return fmt.Errorf("note for %q has surrounding space", key)
	}
	if !slices.Contains(m.RequiredMechanics, key) {
		return fmt.Errorf("stat-block field requires mechanic %q", key)
	}
	return nil
}

// validateVulnerabilityNote accepts a note only in place of a typed
// vulnerability list.
func (m Monster) validateVulnerabilityNote() error {
	if m.VulnerabilityNote == "" {
		return nil
	}
	if m.Vulnerabilities != nil {
		return fmt.Errorf("vulnerability note replaces typed vulnerabilities")
	}
	return m.validateSourceNote(m.VulnerabilityNote, "defense/conditional-damage-vulnerability")
}
