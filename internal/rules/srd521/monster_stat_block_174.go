package srd521

import (
	"fmt"
	"slices"
)

// validateSwarm accepts the member size of a printed swarm type line on a
// record with one size. The record keeps the members' creature type.
func (m Monster) validateSwarm() error {
	if m.SwarmOf == "" {
		return nil
	}
	if m.SwarmOf != "tiny" || len(m.Sizes) != 1 || m.Sizes[0] == m.SwarmOf {
		return fmt.Errorf("invalid swarm member size %q", m.SwarmOf)
	}
	if !slices.Contains(m.RequiredMechanics, "monster/swarm") {
		return fmt.Errorf("stat-block field requires mechanic %q", "monster/swarm")
	}
	return nil
}

// validateSpeedNote accepts a retained open speed choice.
func (m Monster) validateSpeedNote() error {
	if m.SpeedNote == "" {
		return nil
	}
	return m.validateSourceNote(m.SpeedNote, "movement/chosen-mode")
}
