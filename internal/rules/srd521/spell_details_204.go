package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellDominatePersonSlotDuration is one source-listed row.
type SpellDominatePersonSlotDuration struct {
	MinSlot int `json:"min_slot"`
	MaxSlot int `json:"max_slot"`
	Seconds int `json:"seconds"`
}

// SpellDominatePerson retains typed SRD source facts without implementing them.
type SpellDominatePerson struct {
	Target          string                            `json:"target"`
	Save            string                            `json:"save"`
	Condition       string                            `json:"condition"`
	SaveAdvantage   string                            `json:"save_advantage"`
	RepeatSave      string                            `json:"repeat_save"`
	RepeatSuccess   string                            `json:"repeat_success"`
	LinkRequirement string                            `json:"link_requirement"`
	CommandAction   string                            `json:"command_action"`
	CommandTiming   string                            `json:"command_timing"`
	CommandExamples []string                          `json:"command_examples"`
	Obedience       string                            `json:"obedience"`
	NoDirection     string                            `json:"no_direction"`
	ReactionCommand string                            `json:"reaction_command"`
	SlotDurations   []SpellDominatePersonSlotDuration `json:"slot_durations"`
}

// SpellDragonsBreath retains typed SRD source facts without implementing them.
type SpellDragonsBreath struct {
	Target          string   `json:"target"`
	DamageTypes     []string `json:"damage_types"`
	ExhaleAction    string   `json:"exhale_action"`
	ExhaleTiming    string   `json:"exhale_timing"`
	AreaShape       string   `json:"area_shape"`
	AreaFeet        int      `json:"area_feet"`
	AreaTargets     string   `json:"area_targets"`
	Save            string   `json:"save"`
	DamageDiceCount int      `json:"damage_dice_count"`
	DamageDieSides  int      `json:"damage_die_sides"`
	FailDamage      string   `json:"fail_damage"`
	SuccessDamage   string   `json:"success_damage"`
	SlotDiceCount   int      `json:"slot_dice_count"`
	SlotDieSides    int      `json:"slot_die_sides"`
}

// SpellDream retains typed SRD source facts without implementing them.
type SpellDream struct {
	Target                 string   `json:"target"`
	MessengerChoices       []string `json:"messenger_choices"`
	MessengerState         string   `json:"messenger_state"`
	MessengerConditions    []string `json:"messenger_conditions"`
	MessengerSpeed         int      `json:"messenger_speed"`
	AsleepEffect           string   `json:"asleep_effect"`
	DreamShaping           string   `json:"dream_shaping"`
	Emerge                 string   `json:"emerge"`
	Recall                 string   `json:"recall"`
	AwakeKnowledge         string   `json:"awake_knowledge"`
	AwakeOptions           []string `json:"awake_options"`
	TerrifyOption          string   `json:"terrify_option"`
	TerrifyMessageMaxWords int      `json:"terrify_message_max_words"`
	TerrifySave            string   `json:"terrify_save"`
	TerrifyFailure         []string `json:"terrify_failure"`
	TerrifyDamageDiceCount int      `json:"terrify_damage_dice_count"`
	TerrifyDamageDieSides  int      `json:"terrify_damage_die_sides"`
	TerrifyDamageType      string   `json:"terrify_damage_type"`
}

// SpellDruidcraft retains typed SRD source facts without implementing them.
type SpellDruidcraft struct {
	Options                     []string `json:"options"`
	WeatherSensorSize           string   `json:"weather_sensor_size"`
	WeatherSensorPredicts       string   `json:"weather_sensor_predicts"`
	WeatherSensorHours          int      `json:"weather_sensor_hours"`
	WeatherSensorPersistsRounds int      `json:"weather_sensor_persists_rounds"`
	WeatherSensorExamples       []string `json:"weather_sensor_examples"`
	BloomTargets                []string `json:"bloom_targets"`
	BloomTiming                 string   `json:"bloom_timing"`
	SensoryExamples             []string `json:"sensory_examples"`
	SensoryCubeFeet             int      `json:"sensory_cube_feet"`
	FirePlayTargets             []string `json:"fire_play_targets"`
	FirePlayActions             []string `json:"fire_play_actions"`
}

// SpellEarthquake retains typed SRD source facts without implementing them.
type SpellEarthquake struct {
	Point                      string   `json:"point"`
	RadiusFeet                 int      `json:"radius_feet"`
	AreaShape                  string   `json:"area_shape"`
	Terrain                    string   `json:"terrain"`
	Save                       string   `json:"save"`
	SaveTiming                 []string `json:"save_timing"`
	SaveTargets                string   `json:"save_targets"`
	FailureConditions          []string `json:"failure_conditions"`
	FailureBreaksConcentration bool     `json:"failure_breaks_concentration"`
	FissureCountDiceCount      int      `json:"fissure_count_dice_count"`
	FissureCountDieSides       int      `json:"fissure_count_die_sides"`
	FissureTiming              string   `json:"fissure_timing"`
	FissureLocations           string   `json:"fissure_locations"`
	FissureDepthDiceCount      int      `json:"fissure_depth_dice_count"`
	FissureDepthDieSides       int      `json:"fissure_depth_die_sides"`
	FissureDepthMultiplierFeet int      `json:"fissure_depth_multiplier_feet"`
	FissureWidthFeet           int      `json:"fissure_width_feet"`
	FissureExtent              string   `json:"fissure_extent"`
	FissureFallTrigger         string   `json:"fissure_fall_trigger"`
	FissureSave                string   `json:"fissure_save"`
	FissureFailure             string   `json:"fissure_failure"`
	FissureSuccess             string   `json:"fissure_success"`
	StructureDamage            int      `json:"structure_damage"`
	StructureDamageType        string   `json:"structure_damage_type"`
	StructureContact           string   `json:"structure_contact"`
	StructureTiming            []string `json:"structure_timing"`
	StructureCollapse          string   `json:"structure_collapse"`
	CollapseDistance           string   `json:"collapse_distance"`
	CollapseSave               string   `json:"collapse_save"`
	CollapseDamageDiceCount    int      `json:"collapse_damage_dice_count"`
	CollapseDamageDieSides     int      `json:"collapse_damage_die_sides"`
	CollapseDamageType         string   `json:"collapse_damage_type"`
	CollapseFailureConditions  []string `json:"collapse_failure_conditions"`
	CollapseEscapeAction       string   `json:"collapse_escape_action"`
	CollapseEscapeCheck        string   `json:"collapse_escape_check"`
	CollapseEscapeDC           int      `json:"collapse_escape_dc"`
	CollapseSuccess            string   `json:"collapse_success"`
}

// SpellBlastBeamStep is one source-listed row.
type SpellBlastBeamStep struct {
	CharacterLevel int `json:"character_level"`
	Beams          int `json:"beams"`
}

// SpellEldritchBlast retains typed SRD source facts without implementing them.
type SpellEldritchBlast struct {
	Attack          string               `json:"attack"`
	Targets         string               `json:"targets"`
	DamageDiceCount int                  `json:"damage_dice_count"`
	DamageDieSides  int                  `json:"damage_die_sides"`
	DamageType      string               `json:"damage_type"`
	BaseBeams       int                  `json:"base_beams"`
	BeamSteps       []SpellBlastBeamStep `json:"beam_steps"`
	BeamTargeting   string               `json:"beam_targeting"`
	AttackRolls     string               `json:"attack_rolls"`
}

// SpellElementalism retains typed SRD source facts without implementing them.
type SpellElementalism struct {
	Options                       []string `json:"options"`
	BeckonAirCubeFeet             int      `json:"beckon_air_cube_feet"`
	BeckonAirEffects              []string `json:"beckon_air_effects"`
	BeckonAirUnaffected           string   `json:"beckon_air_unaffected"`
	BeckonEarthAreaFeet           int      `json:"beckon_earth_area_feet"`
	BeckonEarthShroud             string   `json:"beckon_earth_shroud"`
	BeckonEarthAlternative        string   `json:"beckon_earth_alternative"`
	BeckonFireCubeFeet            int      `json:"beckon_fire_cube_feet"`
	BeckonFireCloud               string   `json:"beckon_fire_cloud"`
	BeckonFireChoices             []string `json:"beckon_fire_choices"`
	BeckonFireEmbers              string   `json:"beckon_fire_embers"`
	BeckonFireScentLingersSeconds int      `json:"beckon_fire_scent_lingers_seconds"`
	BeckonWaterCubeFeet           int      `json:"beckon_water_cube_feet"`
	BeckonWaterMist               string   `json:"beckon_water_mist"`
	BeckonWaterAlternativeCups    int      `json:"beckon_water_alternative_cups"`
	BeckonWaterAlternativePlaces  []string `json:"beckon_water_alternative_places"`
	BeckonWaterEvaporatesSeconds  int      `json:"beckon_water_evaporates_seconds"`
	SculptMaterials               []string `json:"sculpt_materials"`
	SculptCubeFeet                int      `json:"sculpt_cube_feet"`
	SculptShape                   string   `json:"sculpt_shape"`
	SculptSeconds                 int      `json:"sculpt_seconds"`
}

// SpellEnhanceAbility retains typed SRD source facts without implementing them.
type SpellEnhanceAbility struct {
	Target                   string   `json:"target"`
	Abilities                []string `json:"abilities"`
	Benefit                  string   `json:"benefit"`
	BaseTargets              int      `json:"base_targets"`
	AdditionalTargetsPerSlot int      `json:"additional_targets_per_slot"`
	AbilityPerTarget         string   `json:"ability_per_target"`
}

// SpellEnlargeReduce retains typed SRD source facts without implementing them.
type SpellEnlargeReduce struct {
	SourcePages            []int    `json:"source_pages"`
	Targets                string   `json:"targets"`
	ObjectCondition        string   `json:"object_condition"`
	UnwillingSave          string   `json:"unwilling_save"`
	UnwillingSuccess       string   `json:"unwilling_success"`
	CarriedItems           string   `json:"carried_items"`
	DroppedItems           string   `json:"dropped_items"`
	ThrownItems            string   `json:"thrown_items"`
	EnlargeSizeChange      string   `json:"enlarge_size_change"`
	EnlargeSizeExample     string   `json:"enlarge_size_example"`
	EnlargeAdvantage       []string `json:"enlarge_advantage"`
	EnlargeDamageDiceCount int      `json:"enlarge_damage_dice_count"`
	EnlargeDamageDieSides  int      `json:"enlarge_damage_die_sides"`
	EnlargeDamageSources   []string `json:"enlarge_damage_sources"`
	ReduceSizeChange       string   `json:"reduce_size_change"`
	ReduceSizeExample      string   `json:"reduce_size_example"`
	ReduceDisadvantage     []string `json:"reduce_disadvantage"`
	ReduceDamageDiceCount  int      `json:"reduce_damage_dice_count"`
	ReduceDamageDieSides   int      `json:"reduce_damage_die_sides"`
	ReduceDamageSources    []string `json:"reduce_damage_sources"`
	ReduceMinimumDamage    int      `json:"reduce_minimum_damage"`
}

// SpellEnsnaringStrike retains typed SRD source facts without implementing them.
type SpellEnsnaringStrike struct {
	Trigger          string   `json:"trigger"`
	Save             string   `json:"save"`
	LargeAdvantage   string   `json:"large_advantage"`
	FailureCondition string   `json:"failure_condition"`
	FailureDuration  string   `json:"failure_duration"`
	SuccessResult    string   `json:"success_result"`
	DamageDiceCount  int      `json:"damage_dice_count"`
	DamageDieSides   int      `json:"damage_die_sides"`
	DamageType       string   `json:"damage_type"`
	DamageTiming     string   `json:"damage_timing"`
	EscapeActors     []string `json:"escape_actors"`
	EscapeAction     string   `json:"escape_action"`
	EscapeCheck      string   `json:"escape_check"`
	EscapeDC         string   `json:"escape_dc"`
	EscapeSuccess    string   `json:"escape_success"`
	SlotDiceCount    int      `json:"slot_dice_count"`
	SlotDieSides     int      `json:"slot_die_sides"`
}

// SpellEntangle retains typed SRD source facts without implementing them.
type SpellEntangle struct {
	AreaShape        string `json:"area_shape"`
	AreaFeet         int    `json:"area_feet"`
	AreaOrigin       string `json:"area_origin"`
	Terrain          string `json:"terrain"`
	PlantsDisappear  string `json:"plants_disappear"`
	Targets          string `json:"targets"`
	Save             string `json:"save"`
	FailureCondition string `json:"failure_condition"`
	FailureDuration  string `json:"failure_duration"`
	EscapeAction     string `json:"escape_action"`
	EscapeCheck      string `json:"escape_check"`
	EscapeDC         string `json:"escape_dc"`
	EscapeSuccess    string `json:"escape_success"`
}

// SpellEnthrall retains typed SRD source facts without implementing them.
type SpellEnthrall struct {
	Targets          string   `json:"targets"`
	Save             string   `json:"save"`
	AutomaticSuccess string   `json:"automatic_success"`
	FailurePenalty   int      `json:"failure_penalty" srd:"signed"`
	PenalizedRolls   []string `json:"penalized_rolls"`
	PenaltyDuration  string   `json:"penalty_duration"`
}

func spell204Templates() map[string]any {
	return map[string]any{
		"dominate":     SpellDominatePerson{},
		"breath":       SpellDragonsBreath{},
		"dream":        SpellDream{},
		"druidcraft":   SpellDruidcraft{},
		"earthquake":   SpellEarthquake{},
		"blast":        SpellEldritchBlast{},
		"elementalism": SpellElementalism{},
		"enhance":      SpellEnhanceAbility{},
		"enlarge":      SpellEnlargeReduce{},
		"ensnaring":    SpellEnsnaringStrike{},
		"entangle":     SpellEntangle{},
		"enthrall":     SpellEnthrall{},
	}
}

func validSpell204Target(t SpellTarget) bool {
	switch t.Shape {
	case "single":
		return t.Count == 1 && slices.Contains([]string{"visible-humanoid-in-range", "known-creature-same-plane", "creature-or-object-in-range", "visible-creature-or-object-in-range"}, t.Eligibility)
	case "created-effect":
		return t.Count == 1 && t.Eligibility == "visible-ground-point-in-range"
	}
	return false
}

// validSpell204Scaling covers Eldritch Blast's beam count by character level.
func (s SpellRecord) validSpell204Scaling(effects map[string]bool) bool {
	sc, d := s.Scaling, s.SourceDetails
	return d != nil && d.Blast != nil && s.Level == 0 && sc.Kind == "cantrip-beams" && sc.AboveLevel == 0 && sc.DiceCount == 0 && sc.DieSides == 0 &&
		len(sc.EffectIDs) == 1 && effects[sc.EffectIDs[0]] && len(d.Blast.BeamSteps) == 3 && slices.Contains(s.RequiredMechanics, "scaling.cantrip-beams")
}

// has204Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has204Payload() bool {
	return d != nil && (d.Dominate != nil || d.Breath != nil || d.Dream != nil || d.Druidcraft != nil || d.Earthquake != nil || d.Blast != nil ||
		d.Elementalism != nil || d.Enhance != nil || d.Enlarge != nil || d.Ensnaring != nil || d.Entangle != nil || d.Enthrall != nil)
}

// validateSpell204 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell204() (bool, error) {
	pin, ok := spell204Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has204Payload() {
			return true, fmt.Errorf("mismatched batch 204 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 204 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 204 record differs from the printed source")
	}
	return true, nil
}
