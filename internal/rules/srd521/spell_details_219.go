package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellShillelaghDieUpgrade is one source-listed row or nested block.
type SpellShillelaghDieUpgrade struct {
	Level     int `json:"level"`
	DiceCount int `json:"dice_count"`
	DieSides  int `json:"die_sides"`
}

// SpellShillelagh retains typed SRD source facts without implementing them.
type SpellShillelagh struct {
	WeaponChoices       []string                    `json:"weapon_choices"`
	WeaponRequirement   string                      `json:"weapon_requirement"`
	AbilitySubstitution string                      `json:"ability_substitution"`
	Rolls               []string                    `json:"rolls"`
	DamageDiceCount     int                         `json:"damage_dice_count"`
	DamageDieSides      int                         `json:"damage_die_sides"`
	DamageTypeChoices   []string                    `json:"damage_type_choices"`
	EarlyEnd            []string                    `json:"early_end"`
	DieUpgrades         []SpellShillelaghDieUpgrade `json:"die_upgrades"`
}

// SpellShiningSmite retains typed SRD source facts without implementing them.
type SpellShiningSmite struct {
	ExtraDamageDiceCount int    `json:"extra_damage_dice_count"`
	ExtraDamageDieSides  int    `json:"extra_damage_die_sides"`
	DamageType           string `json:"damage_type"`
	BrightLightFeet      int    `json:"bright_light_feet"`
	AttackBenefit        string `json:"attack_benefit"`
	Invisibility         string `json:"invisibility"`
	SlotDiceCount        int    `json:"slot_dice_count"`
	SlotDieSides         int    `json:"slot_die_sides"`
}

// SpellShockingGrasp retains typed SRD source facts without implementing them.
type SpellShockingGrasp struct {
	Attack          string                   `json:"attack"`
	Restriction     string                   `json:"restriction"`
	RestrictionEnds string                   `json:"restriction_ends"`
	Upgrades        []SpellCantripUpgradeRow `json:"upgrades"`
}

// SpellSilence retains typed SRD source facts without implementing them.
type SpellSilence struct {
	RadiusFeet         int      `json:"radius_feet"`
	Sound              string   `json:"sound"`
	ImmunityTargets    []string `json:"immunity_targets"`
	ImmunityDamageType string   `json:"immunity_damage_type"`
	InsideRequirement  string   `json:"inside_requirement"`
	Condition          string   `json:"condition"`
	VerbalCasting      string   `json:"verbal_casting"`
}

// SpellSilentImage retains typed SRD source facts without implementing them.
type SpellSilentImage struct {
	Phenomena           []string `json:"phenomena"`
	MaxCubeFeet         int      `json:"max_cube_feet"`
	SensoryLimit        string   `json:"sensory_limit"`
	MoveAction          string   `json:"move_action"`
	MoveDestination     string   `json:"move_destination"`
	MoveAppearance      string   `json:"move_appearance"`
	Example             string   `json:"example"`
	PhysicalInteraction string   `json:"physical_interaction"`
	ExamineAction       string   `json:"examine_action"`
	ExamineCheck        string   `json:"examine_check"`
	ExamineDC           string   `json:"examine_dc"`
	Discerned           string   `json:"discerned"`
}

// SpellSimulacrum retains typed SRD source facts without implementing them.
type SpellSimulacrum struct {
	SourcePages            []int    `json:"source_pages"`
	CreatureTypes          []string `json:"creature_types"`
	CastingDistanceFeet    int      `json:"casting_distance_feet"`
	Completion             string   `json:"completion"`
	ResultKind             string   `json:"result_kind"`
	Statistics             string   `json:"statistics"`
	ResultCreatureType     string   `json:"result_creature_type"`
	HitPointMaximumDivisor int      `json:"hit_point_maximum_divisor"`
	ForbiddenCasting       string   `json:"forbidden_casting"`
	Attitude               string   `json:"attitude"`
	Control                string   `json:"control"`
	Limits                 []string `json:"limits"`
	RepairRequirement      string   `json:"repair_requirement"`
	RepairGPPerHitPoint    int      `json:"repair_gp_per_hit_point"`
	RepairDistanceFeet     int      `json:"repair_distance_feet"`
	Destruction            string   `json:"destruction"`
	Recast                 string   `json:"recast"`
}

// SpellSleep retains typed SRD source facts without implementing them.
type SpellSleep struct {
	RadiusFeet             int      `json:"radius_feet"`
	Save                   string   `json:"save"`
	FirstCondition         string   `json:"first_condition"`
	FirstConditionEnds     string   `json:"first_condition_ends"`
	RepeatTiming           string   `json:"repeat_timing"`
	SecondFailureCondition string   `json:"second_failure_condition"`
	EndsOnTarget           []string `json:"ends_on_target"`
	WakeDistanceFeet       int      `json:"wake_distance_feet"`
	AutomaticSuccess       []string `json:"automatic_success"`
	NoSleepExample         string   `json:"no_sleep_example"`
}

// SpellSleetStorm retains typed SRD source facts without implementing them.
type SpellSleetStorm struct {
	HeightFeet   int      `json:"height_feet"`
	RadiusFeet   int      `json:"radius_feet"`
	Obscurement  string   `json:"obscurement"`
	Flames       string   `json:"flames"`
	Ground       string   `json:"ground"`
	SaveTriggers []string `json:"save_triggers"`
	Save         string   `json:"save"`
	Failure      []string `json:"failure"`
}

// SpellSlow retains typed SRD source facts without implementing them.
type SpellSlow struct {
	MaxTargets            int    `json:"max_targets"`
	CubeFeet              int    `json:"cube_feet"`
	Save                  string `json:"save"`
	SpeedDivisor          int    `json:"speed_divisor"`
	ArmorClassPenalty     int    `json:"armor_class_penalty"`
	DexteritySavePenalty  int    `json:"dexterity_save_penalty"`
	Reactions             string `json:"reactions"`
	TurnLimit             string `json:"turn_limit"`
	AttackLimit           int    `json:"attack_limit"`
	SomaticFailurePercent int    `json:"somatic_failure_percent"`
	SomaticFailureCause   string `json:"somatic_failure_cause"`
	RepeatTiming          string `json:"repeat_timing"`
	RepeatSuccess         string `json:"repeat_success"`
}

// SpellSorcerousBurst retains typed SRD source facts without implementing them.
type SpellSorcerousBurst struct {
	Attack              string                   `json:"attack"`
	TargetKinds         []string                 `json:"target_kinds"`
	DamageTypeChoices   []string                 `json:"damage_type_choices"`
	ChooseCount         int                      `json:"choose_count"`
	ExtraDieTriggerRoll int                      `json:"extra_die_trigger_roll"`
	ExtraDieSides       int                      `json:"extra_die_sides"`
	ExtraDieLimit       string                   `json:"extra_die_limit"`
	Upgrades            []SpellCantripUpgradeRow `json:"upgrades"`
}

// SpellSpareTheDyingRangeUpgrade is one source-listed row or nested block.
type SpellSpareTheDyingRangeUpgrade struct {
	Level int `json:"level"`
	Feet  int `json:"feet"`
}

// SpellSpareTheDying retains typed SRD source facts without implementing them.
type SpellSpareTheDying struct {
	RequiredHitPoints int                              `json:"required_hit_points"`
	Requirement       string                           `json:"requirement"`
	Result            string                           `json:"result"`
	BaseRangeFeet     int                              `json:"base_range_feet"`
	RangeUpgrades     []SpellSpareTheDyingRangeUpgrade `json:"range_upgrades"`
}

// SpellSpeakWithAnimals retains typed SRD source facts without implementing them.
type SpellSpeakWithAnimals struct {
	CreatureType         string   `json:"creature_type"`
	Communication        string   `json:"communication"`
	SkillOptions         string   `json:"skill_options"`
	TypicalTopics        string   `json:"typical_topics"`
	MinimumInformation   []string `json:"minimum_information"`
	PerceptionWindowDays int      `json:"perception_window_days"`
}

func spell219Templates() map[string]any {
	return map[string]any{
		"shillelagh":         SpellShillelagh{},
		"shining_smite":      SpellShiningSmite{},
		"shocking_grasp":     SpellShockingGrasp{},
		"silence":            SpellSilence{},
		"silent_image":       SpellSilentImage{},
		"simulacrum":         SpellSimulacrum{},
		"sleep":              SpellSleep{},
		"sleet_storm":        SpellSleetStorm{},
		"slow":               SpellSlow{},
		"sorcerous_burst":    SpellSorcerousBurst{},
		"spare_the_dying":    SpellSpareTheDying{},
		"speak_with_animals": SpellSpeakWithAnimals{},
	}
}

func validSpell219Target(t SpellTarget) bool {
	switch t.Shape {
	case "single":
		return t.Count == 1 && slices.Contains([]string{"held-club-or-quarterstaff", "beast-or-humanoid-within-10-feet-for-entire-casting", "creature-with-0-hit-points-not-dead"}, t.Eligibility)
	case "sphere":
		return t.Count == 0 && slices.Contains([]string{"area-centered-on-chosen-point", "chosen-creatures-in-sphere"}, t.Eligibility)
	case "created-effect":
		return t.Count == 1 && t.Eligibility == "visible-phenomenon-in-range"
	case "cylinder":
		return t.Count == 0 && t.Eligibility == "area-centered-on-chosen-point"
	case "up-to-count":
		return t.Count == 6 && t.Eligibility == "chosen-creatures-in-cube-in-range"
	}
	return false
}

// spell219CantripUpgrade reports whether this batch types the cantrip dice rows.
func spell219CantripUpgrade(d *SpellSourceDetails) bool {
	return (d.ShockingGrasp != nil && len(d.ShockingGrasp.Upgrades) == 3) || (d.SorcerousBurst != nil && len(d.SorcerousBurst.Upgrades) == 3)
}

// validSpell219Scaling covers Shillelagh's changing weapon die and Spare the
// Dying's growing range. Both are cantrip upgrades without added dice.
func (s SpellRecord) validSpell219Scaling(effects map[string]bool) bool {
	sc, d := s.Scaling, s.SourceDetails
	if d == nil || s.Level != 0 || sc.AboveLevel != 0 || sc.DiceCount != 0 || sc.DieSides != 0 || len(sc.EffectIDs) != 1 || !effects[sc.EffectIDs[0]] {
		return false
	}
	switch sc.Kind {
	case "cantrip-die-table":
		return d.Shillelagh != nil && len(d.Shillelagh.DieUpgrades) == 3 && slices.Contains(s.RequiredMechanics, "scaling.cantrip-die-upgrade")
	case "cantrip-range":
		return d.SpareTheDying != nil && len(d.SpareTheDying.RangeUpgrades) == 3 && slices.Contains(s.RequiredMechanics, "scaling.cantrip-range-upgrade")
	}
	return false
}

// has219Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has219Payload() bool {
	return d != nil && (d.Shillelagh != nil || d.ShiningSmite != nil || d.ShockingGrasp != nil || d.Silence != nil || d.SilentImage != nil || d.Simulacrum != nil ||
		d.Sleep != nil || d.SleetStorm != nil || d.Slow != nil || d.SorcerousBurst != nil || d.SpareTheDying != nil || d.SpeakWithAnimals != nil)
}

// validateSpell219 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell219() (bool, error) {
	pin, ok := spell219Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has219Payload() {
			return true, fmt.Errorf("mismatched batch 219 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 219 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 219 record differs from the printed source")
	}
	return true, nil
}
