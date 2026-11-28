package combat

import (
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
)

// BeastFacts validates what a Wolf or Brown Bear encounter needs to start: the
// pinned revision and action, a declared target size, spatial rulings and, for
// Pack Tactics, a resolved ally fact. It does not require reach.
func BeastFacts(profile *srd521.BeastProfile, space *Space) error {
	if profile == nil || space == nil {
		return fmt.Errorf("explicit beast and spatial facts required")
	}
	if _, _, err := profile.Record(); err != nil {
		return err
	}
	if err := space.Validate(); err != nil {
		return err
	}
	if profile.PackTactics() {
		if _, err := space.PackAlly.Eligible(); err != nil {
			return err
		}
	}
	return nil
}

// BeastAttack validates the pinned action against the encounter's facts before
// any roll: BeastFacts, and every attack needs visible five-foot reach.
func BeastAttack(profile *srd521.BeastProfile, space *Space) error {
	if err := BeastFacts(profile, space); err != nil {
		return err
	}
	return space.Melee()
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
	// Every beast attack is a melee attack within 5 feet, so a Prone or
	// Unconscious target only ever gives Advantage here.
	target, err := e.IncomingAttack()
	if err != nil {
		return nil, err
	}
	return append(reasons, target.Advantage...), nil
}
