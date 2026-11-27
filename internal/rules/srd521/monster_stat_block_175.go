package srd521

import "fmt"

// validateConditionImmunityNote accepts a note only in place of a typed
// condition immunity list.
func (m Monster) validateConditionImmunityNote() error {
	if m.ConditionImmunityNote == "" {
		return nil
	}
	if m.ConditionImmunities != nil {
		return fmt.Errorf("condition immunity note replaces typed condition immunities")
	}
	return m.validateSourceNote(m.ConditionImmunityNote, "defense/conditional-condition-immunity")
}
