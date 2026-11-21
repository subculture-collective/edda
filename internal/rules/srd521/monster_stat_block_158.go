package srd521

import (
	"fmt"
	"slices"
	"strings"
)

// validateTelepathyNote accepts a note only on a record that lists telepathy
// and requires language/telepathy-condition.
func (m Monster) validateTelepathyNote() error {
	if m.TelepathyNote == "" {
		return nil
	}
	if m.Telepathy <= 0 || m.TelepathyNote != strings.TrimSpace(m.TelepathyNote) {
		return fmt.Errorf("telepathy note needs a telepathy range and no surrounding space")
	}
	if !slices.Contains(m.RequiredMechanics, "language/telepathy-condition") {
		return fmt.Errorf("stat-block field requires mechanic %q", "language/telepathy-condition")
	}
	return nil
}
