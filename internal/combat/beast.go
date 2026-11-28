package combat

import (
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
)

// BeastAttack validates the pinned Wolf or Brown Bear action against the
// encounter's facts before any roll: every attack needs visible five-foot
// reach, a declared target size and, for Pack Tactics, a resolved ally fact.
func BeastAttack(profile *srd521.BeastProfile, space *Space) error {
	if profile == nil || space == nil {
		return fmt.Errorf("explicit beast and spatial facts required")
	}
	if _, _, err := profile.Record(); err != nil {
		return err
	}
	if err := space.Melee(); err != nil {
		return err
	}
	if profile.PackTactics() {
		if _, err := space.PackAlly.Eligible(); err != nil {
			return err
		}
	}
	return nil
}

// BeastAdvantage names the durable facts that give the beast's melee attacks
// Advantage. Both reasons together still grant a single Advantage.
func (e *Encounter) BeastAdvantage() ([]string, error) {
	if err := BeastAttack(e.Beast, e.Space); err != nil {
		return nil, err
	}
	reasons := []string{}
	if e.Beast.PackTactics() {
		if ok, _ := e.Space.PackAlly.Eligible(); ok {
			reasons = append(reasons, "Pack Tactics: a non-Incapacitated ally is within 5 feet of the target")
		}
	}
	for _, a := range e.Actors {
		if a.Player && a.Prone {
			// Every beast attack is a melee attack within 5 feet.
			reasons = append(reasons, "target is Prone within 5 feet")
		}
	}
	return reasons, nil
}
