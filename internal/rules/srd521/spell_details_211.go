package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellInstantSummons retains typed SRD source facts without implementing them.
type SpellInstantSummons struct {
	ObjectMaxWeightPounds         int      `json:"object_max_weight_pounds"`
	ObjectMaxLongestDimensionFeet int      `json:"object_max_longest_dimension_feet"`
	ObjectMark                    string   `json:"object_mark"`
	SapphireInscription           string   `json:"sapphire_inscription"`
	SapphireReuse                 string   `json:"sapphire_reuse"`
	ActivationAction              string   `json:"activation_action"`
	ActivationSteps               []string `json:"activation_steps"`
	Result                        string   `json:"result"`
	ResultDistance                string   `json:"result_distance"`
	AfterResult                   string   `json:"after_result"`
	HeldByAnother                 string   `json:"held_by_another"`
	HeldKnowledge                 []string `json:"held_knowledge"`
}

// SpellInvisibility retains typed SRD source facts without implementing them.
type SpellInvisibility struct {
	Target                   string   `json:"target"`
	Condition                string   `json:"condition"`
	EarlyEndTriggers         []string `json:"early_end_triggers"`
	EarlyEndTiming           string   `json:"early_end_timing"`
	BaseTargets              int      `json:"base_targets"`
	AdditionalTargetsPerSlot int      `json:"additional_targets_per_slot"`
}

// SpellIrresistibleDance retains typed SRD source facts without implementing them.
type SpellIrresistibleDance struct {
	Save             string   `json:"save"`
	SuccessEffect    string   `json:"success_effect"`
	SuccessMovement  string   `json:"success_movement"`
	FailureCondition string   `json:"failure_condition"`
	CharmedEffects   []string `json:"charmed_effects"`
	RepeatAction     string   `json:"repeat_action"`
	RepeatSave       string   `json:"repeat_save"`
	RepeatSuccess    string   `json:"repeat_success"`
}

// SpellJump retains typed SRD source facts without implementing them.
type SpellJump struct {
	Target                   string `json:"target"`
	JumpMaxFeet              int    `json:"jump_max_feet"`
	MovementCostFeet         int    `json:"movement_cost_feet"`
	UsesPerTurn              int    `json:"uses_per_turn"`
	BaseTargets              int    `json:"base_targets"`
	AdditionalTargetsPerSlot int    `json:"additional_targets_per_slot"`
}

// SpellKnock retains typed SRD source facts without implementing them.
type SpellKnock struct {
	SourcePages                  []int    `json:"source_pages"`
	ObjectExamples               []string `json:"object_examples"`
	MundaneStates                []string `json:"mundane_states"`
	MundaneResults               []string `json:"mundane_results"`
	LocksUnlockedIfMultiple      int      `json:"locks_unlocked_if_multiple"`
	ArcaneLockSuppressionMinutes int      `json:"arcane_lock_suppression_minutes"`
	ArcaneLockResult             string   `json:"arcane_lock_result"`
	Sound                        string   `json:"sound"`
	SoundAudibleFeet             int      `json:"sound_audible_feet"`
}

// SpellLegendLore retains typed SRD source facts without implementing them.
type SpellLegendLore struct {
	IncenseMinCostCP     int      `json:"incense_min_cost_cp"`
	IncenseConsumed      bool     `json:"incense_consumed"`
	IvoryStripCount      int      `json:"ivory_strip_count"`
	IvoryStripMinCostCP  int      `json:"ivory_strip_min_cost_cp"`
	IvoryStripsConsumed  bool     `json:"ivory_strips_consumed"`
	SubjectInput         string   `json:"subject_input"`
	Subjects             []string `json:"subjects"`
	Result               string   `json:"result"`
	DescribedBy          string   `json:"described_by"`
	PossibleContent      []string `json:"possible_content"`
	DetailScaling        string   `json:"detail_scaling"`
	Accuracy             string   `json:"accuracy"`
	AccuracyDeterminedBy string   `json:"accuracy_determined_by"`
	NotFamous            string   `json:"not_famous"`
}

// SpellLesserRestoration retains typed SRD source facts without implementing them.
type SpellLesserRestoration struct {
	Target               string   `json:"target"`
	ConditionsEndedCount int      `json:"conditions_ended_count"`
	ConditionChoices     []string `json:"condition_choices"`
}

// SpellLevitate retains typed SRD source facts without implementing them.
type SpellLevitate struct {
	RiseMaxFeet             int      `json:"rise_max_feet"`
	Suspension              string   `json:"suspension"`
	ObjectMaxWeightPounds   int      `json:"object_max_weight_pounds"`
	UnwillingSave           string   `json:"unwilling_save"`
	UnwillingSuccess        string   `json:"unwilling_success"`
	Movement                string   `json:"movement"`
	MovementSurfaceExamples []string `json:"movement_surface_examples"`
	MovementMode            string   `json:"movement_mode"`
	AltitudeChangeMaxFeet   int      `json:"altitude_change_max_feet"`
	AltitudeDirections      []string `json:"altitude_directions"`
	SelfAltitude            string   `json:"self_altitude"`
	OtherAltitudeAction     string   `json:"other_altitude_action"`
	OtherAltitudeLimit      string   `json:"other_altitude_limit"`
	Ending                  string   `json:"ending"`
}

// SpellLight retains typed SRD source facts without implementing them.
type SpellLight struct {
	ObjectMaxSize          string `json:"object_max_size"`
	ObjectRequirement      string `json:"object_requirement"`
	BrightLightRadiusFeet  int    `json:"bright_light_radius_feet"`
	DimLightAdditionalFeet int    `json:"dim_light_additional_feet"`
	Color                  string `json:"color"`
	Blocking               string `json:"blocking"`
	Recast                 string `json:"recast"`
}

// SpellLightningBolt retains typed SRD source facts without implementing them.
type SpellLightningBolt struct {
	LineLengthFeet  int    `json:"line_length_feet"`
	LineWidthFeet   int    `json:"line_width_feet"`
	Direction       string `json:"direction"`
	Save            string `json:"save"`
	DamageDiceCount int    `json:"damage_dice_count"`
	DamageDieSides  int    `json:"damage_die_sides"`
	DamageType      string `json:"damage_type"`
	SaveSuccess     string `json:"save_success"`
	SlotDiceCount   int    `json:"slot_dice_count"`
	SlotDieSides    int    `json:"slot_die_sides"`
}

// SpellLocateAnimalsOrPlants retains typed SRD source facts without implementing them.
type SpellLocateAnimalsOrPlants struct {
	Query             string   `json:"query"`
	KindCategories    []string `json:"kind_categories"`
	Learns            []string `json:"learns"`
	Result            string   `json:"result"`
	SearchRadiusMiles int      `json:"search_radius_miles"`
	Requirement       string   `json:"requirement"`
}

// SpellLocateCreature retains typed SRD source facts without implementing them.
type SpellLocateCreature struct {
	Query                 string   `json:"query"`
	Senses                string   `json:"senses"`
	SenseRadiusFeet       int      `json:"sense_radius_feet"`
	MovingTarget          string   `json:"moving_target"`
	Targets               []string `json:"targets"`
	KindRequirement       string   `json:"kind_requirement"`
	KindExamples          []string `json:"kind_examples"`
	CloseDistanceFeet     int      `json:"close_distance_feet"`
	CloseSightingsMinimum int      `json:"close_sightings_minimum"`
	DifferentForm         string   `json:"different_form"`
	DifferentFormExamples []string `json:"different_form_examples"`
	LeadBlocker           string   `json:"lead_blocker"`
}

func spell211Templates() map[string]any {
	return map[string]any{
		"instant_summons":          SpellInstantSummons{},
		"invisibility":             SpellInvisibility{},
		"irresistible_dance":       SpellIrresistibleDance{},
		"jump":                     SpellJump{},
		"knock":                    SpellKnock{},
		"legend_lore":              SpellLegendLore{},
		"lesser_restoration":       SpellLesserRestoration{},
		"levitate":                 SpellLevitate{},
		"light":                    SpellLight{},
		"lightning_bolt":           SpellLightningBolt{},
		"locate_animals_or_plants": SpellLocateAnimalsOrPlants{},
		"locate_creature":          SpellLocateCreature{},
	}
}

func validSpell211Target(t SpellTarget) bool {
	switch t.Shape {
	case "single":
		return t.Count == 1 && slices.Contains([]string{"sapphire-and-eligible-object-touch", "visible-access-prevented-object-in-range", "visible-creature-or-loose-object-in-range", "unattended-large-or-smaller-object-touch"}, t.Eligibility)
	}
	return false
}

// spell211TargetScaling returns this batch's typed extra-target increments.
func spell211TargetScaling(d *SpellSourceDetails) (base, additional int, ok bool) {
	if d.Invisibility != nil {
		return d.Invisibility.BaseTargets, d.Invisibility.AdditionalTargetsPerSlot, true
	}
	if d.Jump != nil {
		return d.Jump.BaseTargets, d.Jump.AdditionalTargetsPerSlot, true
	}
	return 0, 0, false
}

// has211Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has211Payload() bool {
	return d != nil && (d.InstantSummons != nil || d.Invisibility != nil || d.IrresistibleDance != nil || d.Jump != nil || d.Knock != nil || d.LegendLore != nil ||
		d.LesserRestoration != nil || d.Levitate != nil || d.Light != nil || d.LightningBolt != nil || d.LocateAnimalsOrPlants != nil || d.LocateCreature != nil)
}

// validateSpell211 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell211() (bool, error) {
	pin, ok := spell211Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has211Payload() {
			return true, fmt.Errorf("mismatched batch 211 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 211 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 211 record differs from the printed source")
	}
	return true, nil
}
