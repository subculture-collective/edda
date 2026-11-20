package srd521

import (
	"fmt"
	"slices"
	"strings"
)

// TelepathyNote retains the parenthetical the source prints after a telepathy
// range, such as the Dretch's "(works only with creatures that understand
// Abyssal)". Like every stat-block field it requires a dependency key, so a
// record that carries one stays unavailable.

// validateTelepathyNote accepts a note only on a record that lists telepathy.
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
