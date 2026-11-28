package srd521

import (
	"fmt"
	"strings"
)

// validateConditionalSpeedNote accepts retained form-only speeds. A movement
// mode the note names cannot also be typed as an unconditional speed.
func (m Monster) validateConditionalSpeedNote() error {
	if m.ConditionalSpeedNote == "" {
		return nil
	}
	for mode := range m.Speeds {
		if strings.Contains(m.ConditionalSpeedNote, monsterModeName(mode)+" ") {
			return fmt.Errorf("conditional speed note replaces the typed %s speed", mode)
		}
	}
	return m.validateSourceNote(m.ConditionalSpeedNote, "movement/conditional-speed")
}

// monsterModeName is the capitalised movement mode as a Speed line prints it.
func monsterModeName(mode string) string {
	if mode == "" {
		return mode
	}
	return strings.ToUpper(mode[:1]) + mode[1:]
}
