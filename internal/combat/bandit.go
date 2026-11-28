package combat

import (
	"fmt"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
)

// BanditAttack validates the selected named action against symmetric known facts.
// Long range, unseen targets and close ranged attacks need unsupported rulings.
func BanditAttack(profile *srd521.BanditProfile, space *Space) error {
	if profile == nil || space == nil {
		return fmt.Errorf("explicit Bandit and spatial facts required")
	}
	_, action, err := profile.Record()
	if err != nil {
		return err
	}
	if err := space.Validate(); err != nil {
		return err
	}
	if action.Attack.Kind == "melee" {
		return space.Melee()
	}
	if !*space.LineOfSight || space.Cover == "total" || *space.DistanceFeet <= 5 || *space.DistanceFeet > action.Attack.Range {
		return fmt.Errorf("Bandit ranged attack requires visible target beyond 5 feet within normal range")
	}
	return nil
}

// BanditTurn plans the active Bandit's turn from the declared distance. Beyond
// its Light Crossbow's normal range it first walks toward the player. It then
// attacks with the one weapon that is legal without Disadvantage: the Scimitar
// within 5 feet, the Light Crossbow beyond 5 feet and within normal range. It
// returns the feet closed and that weapon's profile, or nil when no attack is
// legal. Unresolved facts or an unsupported profile return an error.
func (e *Encounter) BanditTurn() (int, *srd521.BanditProfile, error) {
	if e.Bandit == nil {
		return 0, nil, fmt.Errorf("explicit Bandit required")
	}
	m, _, err := e.Bandit.Record()
	if err != nil {
		return 0, nil, err
	}
	closed := 0
	for _, a := range m.Actions {
		if a.ID == "light-crossbow" {
			if closed, _, err = e.Approach(a.Attack.Range); err != nil {
				return 0, nil, err
			}
		}
	}
	for _, id := range []string{"scimitar", "light-crossbow"} {
		profile := *e.Bandit
		profile.ActionID = id
		if BanditAttack(&profile, e.Space) == nil {
			return closed, &profile, nil
		}
	}
	return closed, nil, nil
}
