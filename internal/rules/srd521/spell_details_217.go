package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellRayOfSickness retains typed SRD source facts without implementing them.
type SpellRayOfSickness struct {
	Appearance      string `json:"appearance"`
	Attack          string `json:"attack"`
	Condition       string `json:"condition"`
	ConditionEnds   string `json:"condition_ends"`
	DamageDiceCount int    `json:"damage_dice_count"`
	DamageDieSides  int    `json:"damage_die_sides"`
	DamageType      string `json:"damage_type"`
	SlotDiceCount   int    `json:"slot_dice_count"`
	SlotDieSides    int    `json:"slot_die_sides"`
}

// SpellRegenerate retains typed SRD source facts without implementing them.
type SpellRegenerate struct {
	InitialDiceCount int    `json:"initial_dice_count"`
	InitialDieSides  int    `json:"initial_die_sides"`
	InitialBonus     int    `json:"initial_bonus"`
	TurnHealing      int    `json:"turn_healing"`
	TurnTiming       string `json:"turn_timing"`
	SeveredParts     string `json:"severed_parts"`
	RegrowthMinutes  int    `json:"regrowth_minutes"`
}

// SpellReincarnateSpecies is one source-listed row or nested block.
type SpellReincarnateSpecies struct {
	Roll    int    `json:"roll"`
	Species string `json:"species"`
}

// SpellReincarnate retains typed SRD source facts without implementing them.
type SpellReincarnate struct {
	TargetForms    []string                  `json:"target_forms"`
	MaxDaysDead    int                       `json:"max_days_dead"`
	Result         string                    `json:"result"`
	TableDieSides  int                       `json:"table_die_sides"`
	Species        []SpellReincarnateSpecies `json:"species"`
	GMAlternative  string                    `json:"gm_alternative"`
	SpeciesChoices string                    `json:"species_choices"`
	Memory         string                    `json:"memory"`
	Retained       string                    `json:"retained"`
	TraitChange    string                    `json:"trait_change"`
}

// SpellRemoveCurse retains typed SRD source facts without implementing them.
type SpellRemoveCurse struct {
	CursesEnded      string `json:"curses_ended"`
	CursedMagicItem  string `json:"cursed_magic_item"`
	AttunementResult string `json:"attunement_result"`
}

// SpellResilientSphere retains typed SRD source facts without implementing them.
type SpellResilientSphere struct {
	MaxTargetSize    string   `json:"max_target_size"`
	TargetKinds      []string `json:"target_kinds"`
	Save             string   `json:"save"`
	SaveApplies      string   `json:"save_applies"`
	Passage          string   `json:"passage"`
	Blocked          []string `json:"blocked"`
	Breathing        string   `json:"breathing"`
	Immunity         string   `json:"immunity"`
	InsideProtection string   `json:"inside_protection"`
	InsideLimit      string   `json:"inside_limit"`
	Size             string   `json:"size"`
	PushAction       string   `json:"push_action"`
	RollSpeedDivisor int      `json:"roll_speed_divisor"`
	OutsideMovement  string   `json:"outside_movement"`
	DestroyedBy      string   `json:"destroyed_by"`
}

// SpellResistance retains typed SRD source facts without implementing them.
type SpellResistance struct {
	DamageTypeChoices  []string `json:"damage_type_choices"`
	ChooseCount        int      `json:"choose_count"`
	ReductionDiceCount int      `json:"reduction_dice_count"`
	ReductionDieSides  int      `json:"reduction_die_sides"`
	UsesPerTurn        int      `json:"uses_per_turn"`
}

// SpellResurrection retains typed SRD source facts without implementing them.
type SpellResurrection struct {
	SourcePages                 []int    `json:"source_pages"`
	MaxYearsDead                int      `json:"max_years_dead"`
	DeathRequirements           []string `json:"death_requirements"`
	ReturnHitPoints             string   `json:"return_hit_points"`
	PoisonBenefit               string   `json:"poison_benefit"`
	Wounds                      string   `json:"wounds"`
	MissingParts                string   `json:"missing_parts"`
	D20TestPenalty              int      `json:"d20_test_penalty"`
	PenaltyReductionPerLongRest int      `json:"penalty_reduction_per_long_rest"`
	PenaltyEndsAt               int      `json:"penalty_ends_at"`
	TaxingDaysDead              int      `json:"taxing_days_dead"`
	TaxingPenalties             []string `json:"taxing_penalties"`
	TaxingEnds                  string   `json:"taxing_ends"`
}

// SpellReverseGravity retains typed SRD source facts without implementing them.
type SpellReverseGravity struct {
	RadiusFeet      int    `json:"radius_feet"`
	HeightFeet      int    `json:"height_feet"`
	Affected        string `json:"affected"`
	Fall            string `json:"fall"`
	GrabSave        string `json:"grab_save"`
	GrabResult      string `json:"grab_result"`
	Collision       string `json:"collision"`
	UnobstructedTop string `json:"unobstructed_top"`
	Ending          string `json:"ending"`
}

// SpellRevivify retains typed SRD source facts without implementing them.
type SpellRevivify struct {
	MaxSecondsDead  int    `json:"max_seconds_dead"`
	ReturnHitPoints int    `json:"return_hit_points"`
	CannotRevive    string `json:"cannot_revive"`
	MissingParts    string `json:"missing_parts"`
}

// SpellRopeTrick retains typed SRD source facts without implementing them.
type SpellRopeTrick struct {
	RopeMotion       string `json:"rope_motion"`
	PortalVisibility string `json:"portal_visibility"`
	PortalWidthFeet  int    `json:"portal_width_feet"`
	PortalLengthFeet int    `json:"portal_length_feet"`
	PortalLocation   string `json:"portal_location"`
	SpaceKind        string `json:"space_kind"`
	Access           string `json:"access"`
	RopeHandling     string `json:"rope_handling"`
	MaxCreatures     int    `json:"max_creatures"`
	MaxCreatureSize  string `json:"max_creature_size"`
	Isolation        string `json:"isolation"`
	Sight            string `json:"sight"`
	Ending           string `json:"ending"`
}

// SpellSacredFlame retains typed SRD source facts without implementing them.
type SpellSacredFlame struct {
	Appearance   string                   `json:"appearance"`
	Save         string                   `json:"save"`
	IgnoredCover []string                 `json:"ignored_cover"`
	Upgrades     []SpellCantripUpgradeRow `json:"upgrades"`
}

// SpellSanctuary retains typed SRD source facts without implementing them.
type SpellSanctuary struct {
	Triggers               []string `json:"triggers"`
	Save                   string   `json:"save"`
	FailureChoices         []string `json:"failure_choices"`
	AreaEffects            string   `json:"area_effects"`
	EndsWhenWardedCreature []string `json:"ends_when_warded_creature"`
}

func spell217Templates() map[string]any {
	return map[string]any{
		"ray_of_sickness":  SpellRayOfSickness{},
		"regenerate":       SpellRegenerate{},
		"reincarnate":      SpellReincarnate{},
		"remove_curse":     SpellRemoveCurse{},
		"resilient_sphere": SpellResilientSphere{},
		"resistance":       SpellResistance{},
		"resurrection":     SpellResurrection{},
		"reverse_gravity":  SpellReverseGravity{},
		"revivify":         SpellRevivify{},
		"rope_trick":       SpellRopeTrick{},
		"sacred_flame":     SpellSacredFlame{},
		"sanctuary":        SpellSanctuary{},
	}
}

func validSpell217Target(t SpellTarget) bool {
	switch t.Shape {
	case "single":
		return t.Count == 1 && slices.Contains([]string{"dead-humanoid-or-piece-of-one-touch", "creature-or-object-touch", "large-or-smaller-creature-or-object-in-range", "rope-touch"}, t.Eligibility)
	case "cylinder":
		return t.Count == 0 && t.Eligibility == "creatures-and-objects-in-reversed-gravity"
	}
	return false
}

// spell217CantripUpgrade reports whether this batch types the cantrip dice rows.
func spell217CantripUpgrade(d *SpellSourceDetails) bool {
	return d.SacredFlame != nil && len(d.SacredFlame.Upgrades) == 3
}

// has217Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has217Payload() bool {
	return d != nil && (d.RayOfSickness != nil || d.Regenerate != nil || d.Reincarnate != nil || d.RemoveCurse != nil || d.ResilientSphere != nil || d.Resistance != nil ||
		d.Resurrection != nil || d.ReverseGravity != nil || d.Revivify != nil || d.RopeTrick != nil || d.SacredFlame != nil || d.Sanctuary != nil)
}

// validateSpell217 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell217() (bool, error) {
	pin, ok := spell217Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has217Payload() {
			return true, fmt.Errorf("mismatched batch 217 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 217 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 217 record differs from the printed source")
	}
	return true, nil
}
