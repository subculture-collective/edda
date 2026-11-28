package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellSpeakWithDead retains typed SRD source facts without implementing them.
type SpellSpeakWithDead struct {
	CorpseRequirement string   `json:"corpse_requirement"`
	FailureConditions []string `json:"failure_conditions"`
	RepeatBlockDays   int      `json:"repeat_block_days"`
	MaxQuestions      int      `json:"max_questions"`
	Knowledge         string   `json:"knowledge"`
	AnswerStyle       []string `json:"answer_style"`
	Truthfulness      string   `json:"truthfulness"`
	Soul              string   `json:"soul"`
	Limits            []string `json:"limits"`
}

// SpellSpeakWithPlants retains typed SRD source facts without implementing them.
type SpellSpeakWithPlants struct {
	EmanationFeet      int      `json:"emanation_feet"`
	EmanationMobility  string   `json:"emanation_mobility"`
	Grants             string   `json:"grants"`
	QuestionTopics     []string `json:"question_topics"`
	QuestionWindowDays int      `json:"question_window_days"`
	TerrainOptions     []string `json:"terrain_options"`
	MovementLimit      string   `json:"movement_limit"`
	PlantCreatures     string   `json:"plant_creatures"`
}

// SpellSpiderClimb retains typed SRD source facts without implementing them.
type SpellSpiderClimb struct {
	Surfaces                 []string `json:"surfaces"`
	Hands                    string   `json:"hands"`
	ClimbSpeed               string   `json:"climb_speed"`
	BaseTargets              int      `json:"base_targets"`
	AdditionalTargetsPerSlot int      `json:"additional_targets_per_slot"`
}

// SpellSpikeGrowth retains typed SRD source facts without implementing them.
type SpellSpikeGrowth struct {
	RadiusFeet             int    `json:"radius_feet"`
	Terrain                string `json:"terrain"`
	DamageIntervalFeet     int    `json:"damage_interval_feet"`
	DamageDiceCount        int    `json:"damage_dice_count"`
	DamageDieSides         int    `json:"damage_die_sides"`
	Camouflage             string `json:"camouflage"`
	RecognitionRequirement string `json:"recognition_requirement"`
	RecognitionAction      string `json:"recognition_action"`
	RecognitionCheck       string `json:"recognition_check"`
	RecognitionDC          string `json:"recognition_dc"`
	RecognitionTiming      string `json:"recognition_timing"`
}

// SpellSpiritGuardiansForm is one source-listed row or nested block.
type SpellSpiritGuardiansForm struct {
	Alignment         string   `json:"alignment"`
	AppearanceChoices []string `json:"appearance_choices"`
	DamageType        string   `json:"damage_type"`
}

// SpellSpiritGuardians retains typed SRD source facts without implementing them.
type SpellSpiritGuardians struct {
	SourcePages     []int                      `json:"source_pages"`
	EmanationFeet   int                        `json:"emanation_feet"`
	Forms           []SpellSpiritGuardiansForm `json:"forms"`
	Exemptions      string                     `json:"exemptions"`
	SpeedDivisor    int                        `json:"speed_divisor"`
	SaveTriggers    []string                   `json:"save_triggers"`
	Save            string                     `json:"save"`
	DamageDiceCount int                        `json:"damage_dice_count"`
	DamageDieSides  int                        `json:"damage_die_sides"`
	SaveSuccess     string                     `json:"save_success"`
	SavesPerTurn    int                        `json:"saves_per_turn"`
	SlotDiceCount   int                        `json:"slot_dice_count"`
	SlotDieSides    int                        `json:"slot_die_sides"`
}

// SpellSpiritualWeapon retains typed SRD source facts without implementing them.
type SpellSpiritualWeapon struct {
	Appearance      string `json:"appearance"`
	Attack          string `json:"attack"`
	AttackReachFeet int    `json:"attack_reach_feet"`
	AttackTargets   int    `json:"attack_targets"`
	DamageDiceCount int    `json:"damage_dice_count"`
	DamageDieSides  int    `json:"damage_die_sides"`
	DamageAdds      string `json:"damage_adds"`
	DamageType      string `json:"damage_type"`
	LaterAction     string `json:"later_action"`
	MoveMaxFeet     int    `json:"move_max_feet"`
	RepeatAttack    string `json:"repeat_attack"`
	SlotDiceCount   int    `json:"slot_dice_count"`
	SlotDieSides    int    `json:"slot_die_sides"`
}

// SpellStarryWisp retains typed SRD source facts without implementing them.
type SpellStarryWisp struct {
	Attack       string                   `json:"attack"`
	TargetKinds  []string                 `json:"target_kinds"`
	DimLightFeet int                      `json:"dim_light_feet"`
	Invisibility string                   `json:"invisibility"`
	EffectEnds   string                   `json:"effect_ends"`
	Upgrades     []SpellCantripUpgradeRow `json:"upgrades"`
}

// SpellStinkingCloud retains typed SRD source facts without implementing them.
type SpellStinkingCloud struct {
	RadiusFeet    int    `json:"radius_feet"`
	Appearance    string `json:"appearance"`
	Obscurement   string `json:"obscurement"`
	Lingers       string `json:"lingers"`
	WindExample   string `json:"wind_example"`
	Save          string `json:"save"`
	SaveTrigger   string `json:"save_trigger"`
	Condition     string `json:"condition"`
	ConditionEnds string `json:"condition_ends"`
	PoisonedLimit string `json:"poisoned_limit"`
}

// SpellStoneShape retains typed SRD source facts without implementing them.
type SpellStoneShape struct {
	MaxObjectSize            string   `json:"max_object_size"`
	MaxSectionDimensionFeet  int      `json:"max_section_dimension_feet"`
	Result                   string   `json:"result"`
	Examples                 []string `json:"examples"`
	ExampleWallThicknessFeet int      `json:"example_wall_thickness_feet"`
	MaxHinges                int      `json:"max_hinges"`
	MaxLatches               int      `json:"max_latches"`
	DetailLimit              string   `json:"detail_limit"`
}

// SpellStoneskin retains typed SRD source facts without implementing them.
type SpellStoneskin struct {
	ResistanceDamageTypes []string `json:"resistance_damage_types"`
}

// SpellStormOfVengeanceTurn is one source-listed row or nested block.
type SpellStormOfVengeanceTurn struct {
	FirstTurn  int    `json:"first_turn"`
	LastTurn   int    `json:"last_turn"`
	Effect     string `json:"effect"`
	DiceCount  int    `json:"dice_count"`
	DieSides   int    `json:"die_sides"`
	DamageType string `json:"damage_type"`
	Save       string `json:"save"`
}

// SpellStormOfVengeance retains typed SRD source facts without implementing them.
type SpellStormOfVengeance struct {
	SourcePages              []int                       `json:"source_pages"`
	RangeMiles               int                         `json:"range_miles"`
	CloudRadiusFeet          int                         `json:"cloud_radius_feet"`
	InitialSave              string                      `json:"initial_save"`
	InitialDiceCount         int                         `json:"initial_dice_count"`
	InitialDieSides          int                         `json:"initial_die_sides"`
	InitialDamageType        string                      `json:"initial_damage_type"`
	InitialCondition         string                      `json:"initial_condition"`
	InitialConditionDuration string                      `json:"initial_condition_duration"`
	TurnTiming               string                      `json:"turn_timing"`
	Turns                    []SpellStormOfVengeanceTurn `json:"turns"`
	LightningBolts           int                         `json:"lightning_bolts"`
	LightningTargets         string                      `json:"lightning_targets"`
	GustEffects              []string                    `json:"gust_effects"`
	GustEffectsEnd           string                      `json:"gust_effects_end"`
}

// SpellSuggestion retains typed SRD source facts without implementing them.
type SpellSuggestion struct {
	MaxWords         int      `json:"max_words"`
	Requirements     []string `json:"requirements"`
	Examples         []string `json:"examples"`
	Save             string   `json:"save"`
	FailureCondition string   `json:"failure_condition"`
	ConditionEnds    []string `json:"condition_ends"`
	Pursuit          string   `json:"pursuit"`
	ActivityDuration string   `json:"activity_duration"`
	Completion       string   `json:"completion"`
}

func spell220Templates() map[string]any {
	return map[string]any{
		"speak_with_dead":    SpellSpeakWithDead{},
		"speak_with_plants":  SpellSpeakWithPlants{},
		"spider_climb":       SpellSpiderClimb{},
		"spike_growth":       SpellSpikeGrowth{},
		"spirit_guardians":   SpellSpiritGuardians{},
		"spiritual_weapon":   SpellSpiritualWeapon{},
		"starry_wisp":        SpellStarryWisp{},
		"stinking_cloud":     SpellStinkingCloud{},
		"stone_shape":        SpellStoneShape{},
		"stoneskin":          SpellStoneskin{},
		"storm_of_vengeance": SpellStormOfVengeance{},
		"suggestion":         SpellSuggestion{},
	}
}

func validSpell220Target(t SpellTarget) bool {
	switch t.Shape {
	case "single":
		return t.Count == 1 && slices.Contains([]string{"corpse-in-range", "stone-object-or-section-touch", "visible-hearing-understanding-creature-in-range"}, t.Eligibility)
	case "emanation":
		return t.Count == 0 && slices.Contains([]string{"plants-in-immobile-emanation", "undesignated-creatures-in-emanation"}, t.Eligibility)
	case "sphere":
		return t.Count == 0 && slices.Contains([]string{"ground-in-sphere-on-point-in-range", "area-centered-on-point-in-range", "creatures-and-objects-under-storm-cloud"}, t.Eligibility)
	case "created-effect":
		return t.Count == 1 && t.Eligibility == "chosen-space-in-range"
	}
	return false
}

// spell220TargetScaling returns this batch's typed extra-target increments.
func spell220TargetScaling(d *SpellSourceDetails) (base, additional int, ok bool) {
	if d.SpiderClimb != nil {
		return d.SpiderClimb.BaseTargets, d.SpiderClimb.AdditionalTargetsPerSlot, true
	}
	return 0, 0, false
}

// spell220CantripUpgrade reports whether this batch types the cantrip dice rows.
func spell220CantripUpgrade(d *SpellSourceDetails) bool {
	return d.StarryWisp != nil && len(d.StarryWisp.Upgrades) == 3
}

// has220Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has220Payload() bool {
	return d != nil && (d.SpeakWithDead != nil || d.SpeakWithPlants != nil || d.SpiderClimb != nil || d.SpikeGrowth != nil || d.SpiritGuardians != nil || d.SpiritualWeapon != nil ||
		d.StarryWisp != nil || d.StinkingCloud != nil || d.StoneShape != nil || d.Stoneskin != nil || d.StormOfVengeance != nil || d.Suggestion != nil)
}

// validateSpell220 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell220() (bool, error) {
	pin, ok := spell220Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has220Payload() {
			return true, fmt.Errorf("mismatched batch 220 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 220 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 220 record differs from the printed source")
	}
	return true, nil
}
