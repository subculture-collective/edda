package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellDraconicSpirit is one source-listed row or nested block.
type SpellDraconicSpirit struct {
	Name                string             `json:"name"`
	Size                string             `json:"size"`
	CreatureType        string             `json:"creature_type"`
	Alignment           string             `json:"alignment"`
	ArmorClassBase      int                `json:"armor_class_base"`
	ArmorClassAdds      string             `json:"armor_class_adds"`
	HitPointsBase       int                `json:"hit_points_base"`
	HitPointsPerLevel   int                `json:"hit_points_per_level"`
	HitPointsAboveLevel int                `json:"hit_points_above_level"`
	SpeedFeet           int                `json:"speed_feet"`
	FlySpeedFeet        int                `json:"fly_speed_feet"`
	SwimSpeedFeet       int                `json:"swim_speed_feet"`
	Abilities           []SpellStatAbility `json:"abilities"`
	DamageResistances   []string           `json:"damage_resistances"`
	ConditionImmunities []string           `json:"condition_immunities"`
	BlindsightFeet      int                `json:"blindsight_feet"`
	DarkvisionFeet      int                `json:"darkvision_feet"`
	PassivePerception   int                `json:"passive_perception"`
	Languages           []string           `json:"languages"`
	ChallengeRating     string             `json:"challenge_rating"`
	XP                  int                `json:"xp"`
	ProficiencyBonus    string             `json:"proficiency_bonus"`
	SharedResistance    string             `json:"shared_resistance"`
	MultiattackRend     string             `json:"multiattack_rend"`
	MultiattackAlso     string             `json:"multiattack_also"`
	RendAttackBonus     string             `json:"rend_attack_bonus"`
	RendReachFeet       int                `json:"rend_reach_feet"`
	RendDiceCount       int                `json:"rend_dice_count"`
	RendDieSides        int                `json:"rend_die_sides"`
	RendBonus           int                `json:"rend_bonus"`
	RendBonusAdds       string             `json:"rend_bonus_adds"`
	RendDamageType      string             `json:"rend_damage_type"`
	BreathSave          string             `json:"breath_save"`
	BreathSaveDC        string             `json:"breath_save_dc"`
	BreathConeFeet      int                `json:"breath_cone_feet"`
	BreathDiceCount     int                `json:"breath_dice_count"`
	BreathDieSides      int                `json:"breath_die_sides"`
	BreathDamageType    string             `json:"breath_damage_type"`
	BreathSuccess       string             `json:"breath_success"`
}

// SpellSummonDragon retains typed SRD source facts without implementing them.
type SpellSummonDragon struct {
	Summons        string              `json:"summons"`
	Placement      string              `json:"placement"`
	DisappearsWhen []string            `json:"disappears_when"`
	CombatRole     string              `json:"combat_role"`
	Initiative     string              `json:"initiative"`
	TurnOrder      string              `json:"turn_order"`
	Commands       string              `json:"commands"`
	NoCommand      string              `json:"no_command"`
	SlotLevelUse   string              `json:"slot_level_use"`
	StatBlock      SpellDraconicSpirit `json:"stat_block"`
}

// SpellSunbeam retains typed SRD source facts without implementing them.
type SpellSunbeam struct {
	LineWidthFeet     int    `json:"line_width_feet"`
	LineLengthFeet    int    `json:"line_length_feet"`
	Save              string `json:"save"`
	DamageDiceCount   int    `json:"damage_dice_count"`
	DamageDieSides    int    `json:"damage_die_sides"`
	DamageType        string `json:"damage_type"`
	FailureCondition  string `json:"failure_condition"`
	ConditionEnds     string `json:"condition_ends"`
	SaveSuccess       string `json:"save_success"`
	NewLineAction     string `json:"new_line_action"`
	MoteLocation      string `json:"mote_location"`
	BrightLightFeet   int    `json:"bright_light_feet"`
	DimLightExtraFeet int    `json:"dim_light_extra_feet"`
	LightKind         string `json:"light_kind"`
}

// SpellSunburst retains typed SRD source facts without implementing them.
type SpellSunburst struct {
	RadiusFeet       int    `json:"radius_feet"`
	Save             string `json:"save"`
	DamageDiceCount  int    `json:"damage_dice_count"`
	DamageDieSides   int    `json:"damage_die_sides"`
	DamageType       string `json:"damage_type"`
	FailureCondition string `json:"failure_condition"`
	ConditionSeconds int    `json:"condition_seconds"`
	SaveSuccess      string `json:"save_success"`
	RepeatSave       string `json:"repeat_save"`
	RepeatTiming     string `json:"repeat_timing"`
	RepeatSuccess    string `json:"repeat_success"`
	Dispels          string `json:"dispels"`
}

// SpellSymbolEffect is one source-listed row or nested block.
type SpellSymbolEffect struct {
	Effect           string `json:"effect"`
	Save             string `json:"save"`
	Result           string `json:"result"`
	DiceCount        int    `json:"dice_count"`
	DieSides         int    `json:"die_sides"`
	DamageType       string `json:"damage_type"`
	ConditionSeconds int    `json:"condition_seconds"`
}

// SpellSymbol retains typed SRD source facts without implementing them.
type SpellSymbol struct {
	Placements           []string            `json:"placements"`
	MaxGlyphDiameterFeet int                 `json:"max_glyph_diameter_feet"`
	MaxObjectMoveFeet    int                 `json:"max_object_move_feet"`
	MovedResult          string              `json:"moved_result"`
	NoticeCheck          string              `json:"notice_check"`
	NoticeDC             string              `json:"notice_dc"`
	CastChoices          []string            `json:"cast_choices"`
	SurfaceTriggers      []string            `json:"surface_triggers"`
	ObjectTriggers       []string            `json:"object_triggers"`
	TriggerRefinements   []string            `json:"trigger_refinements"`
	RefinementExample    string              `json:"refinement_example"`
	ActivationRadiusFeet int                 `json:"activation_radius_feet"`
	ActivationLight      string              `json:"activation_light"`
	ActiveMinutes        int                 `json:"active_minutes"`
	ActiveEnd            string              `json:"active_end"`
	Targeted             []string            `json:"targeted"`
	TargetsPerTurn       int                 `json:"targets_per_turn"`
	Effects              []SpellSymbolEffect `json:"effects"`
	FearMinMoveFeet      int                 `json:"fear_min_move_feet"`
}

// SpellTelekinesis retains typed SRD source facts without implementing them.
type SpellTelekinesis struct {
	SourcePages          []int    `json:"source_pages"`
	UseTiming            string   `json:"use_timing"`
	TargetKinds          []string `json:"target_kinds"`
	Retargeting          string   `json:"retargeting"`
	MaxCreatureSize      string   `json:"max_creature_size"`
	CreatureSave         string   `json:"creature_save"`
	CreatureMoveFeet     int      `json:"creature_move_feet"`
	MoveLimit            string   `json:"move_limit"`
	CreatureCondition    string   `json:"creature_condition"`
	ConditionEnds        string   `json:"condition_ends"`
	Lifted               string   `json:"lifted"`
	Fall                 string   `json:"fall"`
	MaxObjectSize        string   `json:"max_object_size"`
	UnattendedObject     string   `json:"unattended_object"`
	ObjectMoveFeet       int      `json:"object_move_feet"`
	CarriedObjectSave    string   `json:"carried_object_save"`
	CarriedObjectFailure string   `json:"carried_object_failure"`
	FineControl          string   `json:"fine_control"`
	FineControlExamples  []string `json:"fine_control_examples"`
	SourceTextEnd        string   `json:"source_text_end"`
}

// SpellTelepathicBond retains typed SRD source facts without implementing them.
type SpellTelepathicBond struct {
	MaxCreatures  int    `json:"max_creatures"`
	Link          string `json:"link"`
	Excluded      string `json:"excluded"`
	Communication string `json:"communication"`
	Distance      string `json:"distance"`
	PlanarLimit   string `json:"planar_limit"`
}

// SpellTeleportOutcomeRange is one source-listed row or nested block.
type SpellTeleportOutcomeRange struct {
	Outcome string `json:"outcome"`
	Low     int    `json:"low"`
	High    int    `json:"high"`
}

// SpellTeleportFamiliarity is one source-listed row or nested block.
type SpellTeleportFamiliarity struct {
	Familiarity string                      `json:"familiarity"`
	Outcomes    []SpellTeleportOutcomeRange `json:"outcomes"`
}

// SpellTeleportFamiliarityMeaning is one source-listed row or nested block.
type SpellTeleportFamiliarityMeaning struct {
	Familiarity string `json:"familiarity"`
	Meaning     string `json:"meaning"`
}

// SpellTeleport retains typed SRD source facts without implementing them.
type SpellTeleport struct {
	MaxAdditionalCreatures     int                               `json:"max_additional_creatures"`
	CreatureRequirement        string                            `json:"creature_requirement"`
	ObjectAlternative          string                            `json:"object_alternative"`
	MaxObjectSize              string                            `json:"max_object_size"`
	ObjectRestriction          string                            `json:"object_restriction"`
	DestinationRequirements    []string                          `json:"destination_requirements"`
	TableDieSides              int                               `json:"table_die_sides"`
	Table                      []SpellTeleportFamiliarity        `json:"table"`
	FamiliarityMeanings        []SpellTeleportFamiliarityMeaning `json:"familiarity_meanings"`
	LinkedObjectMonths         int                               `json:"linked_object_months"`
	LinkedObjectExample        string                            `json:"linked_object_example"`
	FalseDestinationExamples   []string                          `json:"false_destination_examples"`
	MishapDiceCount            int                               `json:"mishap_dice_count"`
	MishapDieSides             int                               `json:"mishap_die_sides"`
	MishapDamageType           string                            `json:"mishap_damage_type"`
	MishapTargets              string                            `json:"mishap_targets"`
	MishapReroll               string                            `json:"mishap_reroll"`
	SimilarArea                string                            `json:"similar_area"`
	SimilarAreaExample         string                            `json:"similar_area_example"`
	OffTargetDistanceDiceCount int                               `json:"off_target_distance_dice_count"`
	OffTargetDistanceDieSides  int                               `json:"off_target_distance_die_sides"`
	OffTargetDistanceUnit      string                            `json:"off_target_distance_unit"`
	OffTargetDirectionDieSides int                               `json:"off_target_direction_die_sides"`
	OffTargetDirections        []string                          `json:"off_target_directions"`
	OnTarget                   string                            `json:"on_target"`
}

// SpellTeleportationCircle retains typed SRD source facts without implementing them.
type SpellTeleportationCircle struct {
	SourcePages               []int    `json:"source_pages"`
	CircleRadiusFeet          int      `json:"circle_radius_feet"`
	Destination               string   `json:"destination"`
	PortalCloses              string   `json:"portal_closes"`
	ArrivalDistanceFeet       int      `json:"arrival_distance_feet"`
	OccupiedArrival           string   `json:"occupied_arrival"`
	CommonLocations           []string `json:"common_locations"`
	SigilSequence             string   `json:"sigil_sequence"`
	InitialDestinations       int      `json:"initial_destinations"`
	InitialDestinationPlane   string   `json:"initial_destination_plane"`
	InitialDestinationChooser string   `json:"initial_destination_chooser"`
	Learning                  string   `json:"learning"`
	MemorizeMinutes           int      `json:"memorize_minutes"`
	PermanentCastingDays      int      `json:"permanent_casting_days"`
	PermanentRequirement      string   `json:"permanent_requirement"`
}

// SpellThaumaturgyOption is one source-listed row or nested block.
type SpellThaumaturgyOption struct {
	Choice string `json:"choice"`
	Effect string `json:"effect"`
	Lasts  string `json:"lasts"`
}

// SpellThaumaturgy retains typed SRD source facts without implementing them.
type SpellThaumaturgy struct {
	Options               []SpellThaumaturgyOption `json:"options"`
	ChooseCount           int                      `json:"choose_count"`
	VoiceVolumeMultiplier int                      `json:"voice_volume_multiplier"`
	MaxActiveOneMinute    int                      `json:"max_active_one_minute"`
}

// SpellThunderwave retains typed SRD source facts without implementing them.
type SpellThunderwave struct {
	CubeFeet          int    `json:"cube_feet"`
	Save              string `json:"save"`
	DamageDiceCount   int    `json:"damage_dice_count"`
	DamageDieSides    int    `json:"damage_die_sides"`
	DamageType        string `json:"damage_type"`
	CreaturePushFeet  int    `json:"creature_push_feet"`
	SaveSuccess       string `json:"save_success"`
	ObjectRequirement string `json:"object_requirement"`
	ObjectPushFeet    int    `json:"object_push_feet"`
	AudibleFeet       int    `json:"audible_feet"`
	SlotDiceCount     int    `json:"slot_dice_count"`
	SlotDieSides      int    `json:"slot_die_sides"`
}

// SpellTimeStop retains typed SRD source facts without implementing them.
type SpellTimeStop struct {
	TurnDiceCount   int      `json:"turn_dice_count"`
	TurnDieSides    int      `json:"turn_die_sides"`
	TurnBonus       int      `json:"turn_bonus"`
	Others          string   `json:"others"`
	TurnUse         string   `json:"turn_use"`
	EndingTriggers  []string `json:"ending_triggers"`
	MaxDistanceFeet int      `json:"max_distance_feet"`
}

// SpellTinyHut retains typed SRD source facts without implementing them.
type SpellTinyHut struct {
	EmanationFeet        int      `json:"emanation_feet"`
	Mobility             string   `json:"mobility"`
	CastingFailure       string   `json:"casting_failure"`
	FreePassage          string   `json:"free_passage"`
	Barred               string   `json:"barred"`
	MaxBlockedSpellLevel int      `json:"max_blocked_spell_level"`
	BlockedSpells        string   `json:"blocked_spells"`
	Atmosphere           string   `json:"atmosphere"`
	InteriorLightChoices []string `json:"interior_light_choices"`
	InteriorLightAction  string   `json:"interior_light_action"`
	Exterior             string   `json:"exterior"`
	InteriorView         string   `json:"interior_view"`
	EarlyEnd             []string `json:"early_end"`
}

func spell221Templates() map[string]any {
	return map[string]any{
		"summon_dragon":        SpellSummonDragon{},
		"sunbeam":              SpellSunbeam{},
		"sunburst":             SpellSunburst{},
		"symbol":               SpellSymbol{},
		"telekinesis":          SpellTelekinesis{},
		"telepathic_bond":      SpellTelepathicBond{},
		"teleport":             SpellTeleport{},
		"teleportation_circle": SpellTeleportationCircle{},
		"thaumaturgy":          SpellThaumaturgy{},
		"thunderwave":          SpellThunderwave{},
		"time_stop":            SpellTimeStop{},
		"tiny_hut":             SpellTinyHut{},
	}
}

func validSpell221Target(t SpellTarget) bool {
	switch t.Shape {
	case "line":
		return t.Count == 0 && t.Eligibility == "creatures-in-self-origin-line"
	case "sphere":
		return t.Count == 0 && t.Eligibility == "creatures-in-sphere-on-point-in-range"
	case "choice":
		return t.Count == 1 && slices.Contains([]string{"surface-or-closable-object-touch", "visible-creature-or-object-in-range", "willing-creatures-or-single-object-in-range", "one-chosen-minor-wonder-in-range"}, t.Eligibility)
	case "up-to-count":
		return t.Count == 8 && t.Eligibility == "chosen-willing-creatures-in-range"
	case "created-effect":
		return t.Count == 1 && t.Eligibility == "circle-on-ground-in-range"
	case "cube":
		return t.Count == 0 && t.Eligibility == "creatures-and-unsecured-objects-in-self-origin-cube"
	case "emanation":
		return t.Count == 0 && t.Eligibility == "stationary-emanation-around-caster"
	}
	return false
}

// validSpell221Scaling covers Summon Dragon's stat block using the slot level.
func (s SpellRecord) validSpell221Scaling(effects map[string]bool) bool {
	sc, d := s.Scaling, s.SourceDetails
	return d != nil && d.SummonDragon != nil && sc.Kind == "slot-stat-block" && sc.AboveLevel == s.Level && sc.DiceCount == 0 && sc.DieSides == 0 &&
		len(sc.EffectIDs) == 1 && effects[sc.EffectIDs[0]] && d.SummonDragon.SlotLevelUse != "" && slices.Contains(s.RequiredMechanics, "scaling.stat-block-by-slot")
}

// has221Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has221Payload() bool {
	return d != nil && (d.SummonDragon != nil || d.Sunbeam != nil || d.Sunburst != nil || d.Symbol != nil || d.Telekinesis != nil || d.TelepathicBond != nil ||
		d.Teleport != nil || d.TeleportationCircle != nil || d.Thaumaturgy != nil || d.Thunderwave != nil || d.TimeStop != nil || d.TinyHut != nil)
}

// validateSpell221 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell221() (bool, error) {
	pin, ok := spell221Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has221Payload() {
			return true, fmt.Errorf("mismatched batch 221 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 221 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 221 record differs from the printed source")
	}
	return true, nil
}
