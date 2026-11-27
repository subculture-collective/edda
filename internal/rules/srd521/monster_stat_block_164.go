package srd521

import (
	"fmt"
	"slices"
	"strings"
)

// validateResistanceNote accepts a note only in place of a typed resistance
// list, on a record that requires defense/chosen-damage-resistance.
func (m Monster) validateResistanceNote() error {
	if m.ResistanceNote == "" {
		return nil
	}
	if m.Resistances != nil || m.ResistanceNote != strings.TrimSpace(m.ResistanceNote) {
		return fmt.Errorf("resistance note replaces typed resistances and has no surrounding space")
	}
	if !slices.Contains(m.RequiredMechanics, "defense/chosen-damage-resistance") {
		return fmt.Errorf("stat-block field requires mechanic %q", "defense/chosen-damage-resistance")
	}
	return nil
}
