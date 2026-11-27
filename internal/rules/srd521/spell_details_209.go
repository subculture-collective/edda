package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellHallucinatoryTerrain retains typed SRD source facts without implementing them.
type SpellHallucinatoryTerrain struct {
	CubeFeet               int      `json:"cube_feet"`
	AffectedTerrain        string   `json:"affected_terrain"`
	DisguisedSenses        []string `json:"disguised_senses"`
	Disguise               string   `json:"disguise"`
	Examples               []string `json:"examples"`
	Unchanged              []string `json:"unchanged"`
	TactileCharacteristics string   `json:"tactile_characteristics"`
	TouchResult            string   `json:"touch_result"`
	DisbelieveCondition    string   `json:"disbelieve_condition"`
	DisbelieveAction       string   `json:"disbelieve_action"`
	DisbelieveCheck        string   `json:"disbelieve_check"`
	DisbelieveDC           string   `json:"disbelieve_dc"`
	Discerned              string   `json:"discerned"`
}

// SpellHarm retains typed SRD source facts without implementing them.
type SpellHarm struct {
	Target                 string `json:"target"`
	Save                   string `json:"save"`
	DamageDiceCount        int    `json:"damage_dice_count"`
	DamageDieSides         int    `json:"damage_die_sides"`
	DamageType             string `json:"damage_type"`
	FailureHitPointMaximum string `json:"failure_hit_point_maximum"`
	Success                string `json:"success"`
	HitPointMaximumFloor   int    `json:"hit_point_maximum_floor"`
}

// SpellHaste retains typed SRD source facts without implementing them.
type SpellHaste struct {
	Target                   string   `json:"target"`
	SpeedMultiplier          int      `json:"speed_multiplier"`
	ArmorClassBonus          int      `json:"armor_class_bonus"`
	Advantage                string   `json:"advantage"`
	AdditionalActionsPerTurn int      `json:"additional_actions_per_turn"`
	AllowedActions           []string `json:"allowed_actions"`
	EndCondition             string   `json:"end_condition"`
	EndSpeedFeet             int      `json:"end_speed_feet"`
	EndDuration              string   `json:"end_duration"`
	EndDescription           string   `json:"end_description"`
}

// SpellHeal retains typed SRD source facts without implementing them.
type SpellHeal struct {
	Target                     string   `json:"target"`
	HealingHitPoints           int      `json:"healing_hit_points"`
	ConditionsEnded            []string `json:"conditions_ended"`
	AdditionalHitPointsPerSlot int      `json:"additional_hit_points_per_slot"`
}

// SpellHealingWord retains typed SRD source facts without implementing them.
type SpellHealingWord struct {
	Target           string `json:"target"`
	HealingDiceCount int    `json:"healing_dice_count"`
	HealingDieSides  int    `json:"healing_die_sides"`
	HealingAdds      string `json:"healing_adds"`
	SlotDiceCount    int    `json:"slot_dice_count"`
	SlotDieSides     int    `json:"slot_die_sides"`
}

// SpellHeatMetal retains typed SRD source facts without implementing them.
type SpellHeatMetal struct {
	SourcePages         []int    `json:"source_pages"`
	Object              string   `json:"object"`
	ObjectExamples      []string `json:"object_examples"`
	ObjectEffect        string   `json:"object_effect"`
	Damaged             string   `json:"damaged"`
	DamageDiceCount     int      `json:"damage_dice_count"`
	DamageDieSides      int      `json:"damage_die_sides"`
	DamageType          string   `json:"damage_type"`
	InitialTiming       string   `json:"initial_timing"`
	RepeatAction        string   `json:"repeat_action"`
	RepeatRequirement   string   `json:"repeat_requirement"`
	SaveTrigger         string   `json:"save_trigger"`
	Save                string   `json:"save"`
	SaveFailure         string   `json:"save_failure"`
	KeptPenalty         []string `json:"kept_penalty"`
	KeptPenaltyDuration string   `json:"kept_penalty_duration"`
	SlotDiceCount       int      `json:"slot_dice_count"`
	SlotDieSides        int      `json:"slot_die_sides"`
}

// SpellHellishRebuke retains typed SRD source facts without implementing them.
type SpellHellishRebuke struct {
	Target              string `json:"target"`
	TriggerDistanceFeet int    `json:"trigger_distance_feet"`
	TriggerRequirement  string `json:"trigger_requirement"`
	Appearance          string `json:"appearance"`
	Save                string `json:"save"`
	DamageDiceCount     int    `json:"damage_dice_count"`
	DamageDieSides      int    `json:"damage_die_sides"`
	DamageType          string `json:"damage_type"`
	Success             string `json:"success"`
	SlotDiceCount       int    `json:"slot_dice_count"`
	SlotDieSides        int    `json:"slot_die_sides"`
}

// SpellHeroesFeast retains typed SRD source facts without implementing them.
type SpellHeroesFeast struct {
	Location                 string   `json:"location"`
	CubeFeet                 int      `json:"cube_feet"`
	ConsumeHours             int      `json:"consume_hours"`
	FeastDisappears          string   `json:"feast_disappears"`
	BenefitsBegin            string   `json:"benefits_begin"`
	MaxCreatures             int      `json:"max_creatures"`
	BenefitHours             int      `json:"benefit_hours"`
	DamageResistance         string   `json:"damage_resistance"`
	ConditionImmunities      []string `json:"condition_immunities"`
	HitPointMaximumDiceCount int      `json:"hit_point_maximum_dice_count"`
	HitPointMaximumDieSides  int      `json:"hit_point_maximum_die_sides"`
	HitPointsGained          string   `json:"hit_points_gained"`
}

// SpellHeroism retains typed SRD source facts without implementing them.
type SpellHeroism struct {
	Target                   string `json:"target"`
	ConditionImmunity        string `json:"condition_immunity"`
	TemporaryHitPoints       string `json:"temporary_hit_points"`
	TemporaryHitPointsTiming string `json:"temporary_hit_points_timing"`
	BaseTargets              int    `json:"base_targets"`
	AdditionalTargetsPerSlot int    `json:"additional_targets_per_slot"`
}

// SpellHexSlotDuration is one source-listed row.
type SpellHexSlotDuration struct {
	MinSlot int `json:"min_slot"`
	MaxSlot int `json:"max_slot"`
	Seconds int `json:"seconds"`
}

// SpellHex retains typed SRD source facts without implementing them.
type SpellHex struct {
	Target               string                 `json:"target"`
	ExtraDamageDiceCount int                    `json:"extra_damage_dice_count"`
	ExtraDamageDieSides  int                    `json:"extra_damage_die_sides"`
	DamageType           string                 `json:"damage_type"`
	DamageTrigger        string                 `json:"damage_trigger"`
	AbilityChoiceTiming  string                 `json:"ability_choice_timing"`
	AbilityChoiceCount   int                    `json:"ability_choice_count"`
	CheckPenalty         string                 `json:"check_penalty"`
	RetargetCondition    string                 `json:"retarget_condition"`
	RetargetAction       string                 `json:"retarget_action"`
	RetargetResult       string                 `json:"retarget_result"`
	SlotDurations        []SpellHexSlotDuration `json:"slot_durations"`
}

// SpellHideousLaughter retains typed SRD source facts without implementing them.
type SpellHideousLaughter struct {
	Target                   string   `json:"target"`
	Save                     string   `json:"save"`
	FailureConditions        []string `json:"failure_conditions"`
	Laughter                 string   `json:"laughter"`
	ProneRestriction         string   `json:"prone_restriction"`
	RepeatTriggers           []string `json:"repeat_triggers"`
	DamageSave               string   `json:"damage_save"`
	RepeatSuccess            string   `json:"repeat_success"`
	BaseTargets              int      `json:"base_targets"`
	AdditionalTargetsPerSlot int      `json:"additional_targets_per_slot"`
}

// SpellHoldMonster retains typed SRD source facts without implementing them.
type SpellHoldMonster struct {
	SourcePages              []int  `json:"source_pages"`
	Target                   string `json:"target"`
	Save                     string `json:"save"`
	FailureCondition         string `json:"failure_condition"`
	RepeatTiming             string `json:"repeat_timing"`
	RepeatSuccess            string `json:"repeat_success"`
	BaseTargets              int    `json:"base_targets"`
	AdditionalTargetsPerSlot int    `json:"additional_targets_per_slot"`
}

func spell209Templates() map[string]any {
	return map[string]any{
		"hallucinatory_terrain": SpellHallucinatoryTerrain{},
		"harm":                  SpellHarm{},
		"haste":                 SpellHaste{},
		"heal":                  SpellHeal{},
		"healing_word":          SpellHealingWord{},
		"heat_metal":            SpellHeatMetal{},
		"hellish_rebuke":        SpellHellishRebuke{},
		"heroes_feast":          SpellHeroesFeast{},
		"heroism":               SpellHeroism{},
		"hex":                   SpellHex{},
		"hideous_laughter":      SpellHideousLaughter{},
		"hold_monster":          SpellHoldMonster{},
	}
}

func validSpell209Target(t SpellTarget) bool {
	switch t.Shape {
	case "cube":
		return t.Count == 0 && t.Eligibility == "natural-terrain-in-cube"
	case "single":
		return t.Count == 1 && slices.Contains([]string{"visible-willing-creature-in-range", "visible-manufactured-metal-object-in-range", "creature-that-damaged-caster"}, t.Eligibility)
	case "created-effect":
		return t.Count == 1 && t.Eligibility == "surface-in-unoccupied-cube-next-to-caster"
	}
	return false
}

// spell209TargetScaling returns this batch's typed extra-target increments.
func spell209TargetScaling(d *SpellSourceDetails) (base, additional int, ok bool) {
	if d.Heroism != nil {
		return d.Heroism.BaseTargets, d.Heroism.AdditionalTargetsPerSlot, true
	}
	if d.HideousLaughter != nil {
		return d.HideousLaughter.BaseTargets, d.HideousLaughter.AdditionalTargetsPerSlot, true
	}
	if d.HoldMonster != nil {
		return d.HoldMonster.BaseTargets, d.HoldMonster.AdditionalTargetsPerSlot, true
	}
	return 0, 0, false
}

// spell209DurationScaling reports whether a typed slot table or increment backs
// this batch's slot-duration scaling.
func spell209DurationScaling(d *SpellSourceDetails, mechanics []string) bool {
	return (d.Hex != nil && len(d.Hex.SlotDurations) == 3 && slices.Contains(mechanics, "scaling.duration-by-slot"))
}

// has209Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has209Payload() bool {
	return d != nil && (d.HallucinatoryTerrain != nil || d.Harm != nil || d.Haste != nil || d.Heal != nil || d.HealingWord != nil || d.HeatMetal != nil ||
		d.HellishRebuke != nil || d.HeroesFeast != nil || d.Heroism != nil || d.Hex != nil || d.HideousLaughter != nil || d.HoldMonster != nil)
}

// validateSpell209 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell209() (bool, error) {
	pin, ok := spell209Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has209Payload() {
			return true, fmt.Errorf("mismatched batch 209 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 209 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 209 record differs from the printed source")
	}
	return true, nil
}
