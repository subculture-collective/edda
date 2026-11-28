package combat

import "fmt"

// Space records the declared, symmetric spatial facts for the two-actor fixture.
// Nil/empty facts mean unresolved, never a guessed distance or clear sight line.
type Space struct {
	DistanceFeet *int   `json:"distance_feet"`
	LineOfSight  *bool  `json:"line_of_sight"`
	Cover        string `json:"cover"`
	// Light and Obscurement are optional declared rulings at the target's
	// space (#145). Legacy facts omit them; Sight then fails closed.
	Light       string `json:"light,omitempty"`
	Obscurement string `json:"obscurement,omitempty"`
	// Door optionally names one known door between the actor and the target
	// (#146). The other facts describe the path with that door open.
	Door     string `json:"door,omitempty"`
	Revision string `json:"revision"`
	// PackAlly is the declared state of the NPC's nearest ally relative to the
	// player. It is read only by Pack Tactics; nil means unresolved.
	PackAlly *PackAlly `json:"pack_ally,omitempty"`
}

// PackAlly declares whether an ally of the NPC is present, its exact distance
// from the player and whether it has the Incapacitated condition. The ally is
// a declared fact, not an encounter actor.
type PackAlly struct {
	Present       *bool `json:"present"`
	DistanceFeet  *int  `json:"distance_feet,omitempty"`
	Incapacitated *bool `json:"incapacitated,omitempty"`
}

// Eligible reports whether the ally grants Pack Tactics. Unknown facts fail.
func (a *PackAlly) Eligible() (bool, error) {
	if a == nil || a.Present == nil {
		return false, fmt.Errorf("Clarify whether an ally of the attacker is near the target.")
	}
	if !*a.Present {
		return false, nil
	}
	if a.DistanceFeet == nil || *a.DistanceFeet < 0 || a.Incapacitated == nil {
		return false, fmt.Errorf("Clarify the ally's exact distance from the target and whether it is Incapacitated.")
	}
	return *a.DistanceFeet <= 5 && !*a.Incapacitated, nil
}

// ThroughDoor applies the declared door's state. A closed door is a solid
// obstacle: total cover and no line of sight. An open door leaves the declared
// facts. An unresolved door clears line of sight so every ruling fails closed.
func (s Space) ThroughDoor(state string) Space {
	switch state {
	case "open":
	case "closed":
		sight := false
		s.LineOfSight, s.Cover = &sight, "total"
	default:
		s.LineOfSight = nil
	}
	return s
}

// Validate requires an exact distance, sight ruling and recognized cover ruling.
func (s Space) Validate() error {
	if s.DistanceFeet == nil || *s.DistanceFeet < 0 || s.LineOfSight == nil {
		return fmt.Errorf("Clarify the exact distance and line of sight before rolling.")
	}
	switch s.Cover {
	case "none", "half", "three-quarters", "total":
		return nil
	default:
		return fmt.Errorf("Clarify cover before rolling.")
	}
}

// Melee validates ordinary five-foot reach. Unseen attacks remain unsupported.
func (s Space) Melee() error {
	if err := s.Validate(); err != nil {
		return err
	}
	if !*s.LineOfSight {
		return fmt.Errorf("No clear line of sight; unseen attacks require adjudication.")
	}
	if s.Cover == "total" {
		return fmt.Errorf("Total cover blocks this attack.")
	}
	if *s.DistanceFeet > 5 {
		return fmt.Errorf("Out of reach: this melee attack requires 5 feet or less.")
	}
	return nil
}

// CoverBonus is the symmetric fixture's declared bonus to Armor Class.
func (s Space) CoverBonus() int {
	switch s.Cover {
	case "half":
		return 2
	case "three-quarters":
		return 5
	}
	return 0
}

// Move changes an exact distance by five feet without inventing a grid position.
func (s *Space) Move(kind string) error {
	if err := s.Validate(); err != nil {
		return err
	}
	if !*s.LineOfSight || s.Cover != "none" {
		return fmt.Errorf("Clarify an unobstructed path before moving.")
	}
	distance := *s.DistanceFeet
	switch kind {
	case EncounterApproach:
		if distance <= 5 {
			return fmt.Errorf("Already within 5 feet; cannot approach further.")
		}
		distance -= 5
	case EncounterRetreat:
		distance += 5
	default:
		return fmt.Errorf("Unsupported spatial movement.")
	}
	s.DistanceFeet = &distance
	return nil
}
