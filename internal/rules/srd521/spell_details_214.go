package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellMirrorImage retains typed SRD source facts without implementing them.
type SpellMirrorImage struct {
	DuplicateCount            int      `json:"duplicate_count"`
	Location                  string   `json:"location"`
	Behavior                  []string `json:"behavior"`
	RedirectTrigger           string   `json:"redirect_trigger"`
	RedirectDieSides          int      `json:"redirect_die_sides"`
	DicePerRemainingDuplicate int      `json:"dice_per_remaining_duplicate"`
	RedirectThreshold         int      `json:"redirect_threshold"`
	RedirectResult            string   `json:"redirect_result"`
	DuplicateIgnores          string   `json:"duplicate_ignores"`
	EndCondition              string   `json:"end_condition"`
	UnaffectedCreatures       []string `json:"unaffected_creatures"`
}

// SpellMislead retains typed SRD source facts without implementing them.
type SpellMislead struct {
	InvisibilityTiming      string   `json:"invisibility_timing"`
	DoubleLocation          string   `json:"double_location"`
	DoubleDuration          string   `json:"double_duration"`
	InvisibilityEndTriggers []string `json:"invisibility_end_triggers"`
	MoveAction              string   `json:"move_action"`
	MoveSpeedMultiplier     int      `json:"move_speed_multiplier"`
	DoubleBehaviors         []string `json:"double_behaviors"`
	DoubleTraits            []string `json:"double_traits"`
	Senses                  string   `json:"senses"`
}

// SpellMistyStep retains typed SRD source facts without implementing them.
type SpellMistyStep struct {
	Appearance      string `json:"appearance"`
	TeleportMaxFeet int    `json:"teleport_max_feet"`
	Destination     string `json:"destination"`
}

// SpellModifyMemorySlotWindow is one source-listed row.
type SpellModifyMemorySlotWindow struct {
	Slot int `json:"slot"`
	Days int `json:"days"`
}

// SpellModifyMemory retains typed SRD source facts without implementing them.
type SpellModifyMemory struct {
	Save                    string                        `json:"save"`
	SaveAdvantage           string                        `json:"save_advantage"`
	FailureCondition        string                        `json:"failure_condition"`
	CharmedCondition        string                        `json:"charmed_condition"`
	CharmedAwareness        string                        `json:"charmed_awareness"`
	InterruptTriggers       []string                      `json:"interrupt_triggers"`
	InterruptResult         string                        `json:"interrupt_result"`
	EventWindowHours        int                           `json:"event_window_hours"`
	EventMaxMinutes         int                           `json:"event_max_minutes"`
	MemoryOptions           []string                      `json:"memory_options"`
	DescriptionRequirements []string                      `json:"description_requirements"`
	Gaps                    string                        `json:"gaps"`
	UnfinishedDescription   string                        `json:"unfinished_description"`
	TakesHold               string                        `json:"takes_hold"`
	BehaviorEffect          string                        `json:"behavior_effect"`
	Contradictions          []string                      `json:"contradictions"`
	IllogicalMemory         string                        `json:"illogical_memory"`
	IllogicalExample        string                        `json:"illogical_example"`
	GMDiscretion            string                        `json:"gm_discretion"`
	RestoredBy              []string                      `json:"restored_by"`
	SlotWindows             []SpellModifyMemorySlotWindow `json:"slot_windows"`
	AnyPastTimeSlot         int                           `json:"any_past_time_slot"`
}

// SpellMoonbeam retains typed SRD source facts without implementing them.
type SpellMoonbeam struct {
	SourcePages         []int    `json:"source_pages"`
	Beam                string   `json:"beam"`
	CylinderRadiusFeet  int      `json:"cylinder_radius_feet"`
	CylinderHeightFeet  int      `json:"cylinder_height_feet"`
	Center              string   `json:"center"`
	Light               string   `json:"light"`
	MoveAction          string   `json:"move_action"`
	MoveMaxFeet         int      `json:"move_max_feet"`
	Save                string   `json:"save"`
	DamageDiceCount     int      `json:"damage_dice_count"`
	DamageDieSides      int      `json:"damage_die_sides"`
	DamageType          string   `json:"damage_type"`
	ShapeShiftedFailure string   `json:"shape_shifted_failure"`
	ShapeShiftExample   string   `json:"shape_shift_example"`
	SaveSuccess         string   `json:"save_success"`
	SaveTriggers        []string `json:"save_triggers"`
	SavesPerTurn        int      `json:"saves_per_turn"`
	SlotDiceCount       int      `json:"slot_dice_count"`
	SlotDieSides        int      `json:"slot_die_sides"`
}

// SpellMoveEarth retains typed SRD source facts without implementing them.
type SpellMoveEarth struct {
	MaxSideFeet            int      `json:"max_side_feet"`
	Materials              []string `json:"materials"`
	Changes                []string `json:"changes"`
	ExtentLimit            string   `json:"extent_limit"`
	ExampleSideFeet        int      `json:"example_side_feet"`
	ExampleMaxChangeFeet   int      `json:"example_max_change_feet"`
	ChangeMinutes          int      `json:"change_minutes"`
	CreatureSafety         string   `json:"creature_safety"`
	NewAreaIntervalMinutes int      `json:"new_area_interval_minutes"`
	NewArea                string   `json:"new_area"`
	Excluded               []string `json:"excluded"`
	RocksAndStructures     string   `json:"rocks_and_structures"`
	UnstableStructures     string   `json:"unstable_structures"`
	Plants                 string   `json:"plants"`
}

// SpellNondetection retains typed SRD source facts without implementing them.
type SpellNondetection struct {
	HiddenFrom       string   `json:"hidden_from"`
	TargetChoices    []string `json:"target_choices"`
	MaxDimensionFeet int      `json:"max_dimension_feet"`
	Protections      []string `json:"protections"`
}

// SpellPassWithoutTrace retains typed SRD source facts without implementing them.
type SpellPassWithoutTrace struct {
	Aura          string `json:"aura"`
	EmanationFeet int    `json:"emanation_feet"`
	Beneficiaries string `json:"beneficiaries"`
	Check         string `json:"check"`
	CheckBonus    int    `json:"check_bonus"`
	Tracks        string `json:"tracks"`
}

// SpellPasswall retains typed SRD source facts without implementing them.
type SpellPasswall struct {
	Surfaces         []string `json:"surfaces"`
	SurfaceExamples  []string `json:"surface_examples"`
	DimensionsChoice string   `json:"dimensions_choice"`
	MaxWidthFeet     int      `json:"max_width_feet"`
	MaxHeightFeet    int      `json:"max_height_feet"`
	MaxDepthFeet     int      `json:"max_depth_feet"`
	Stability        string   `json:"stability"`
	Ejection         string   `json:"ejection"`
}

// SpellPhantasmalForce retains typed SRD source facts without implementing them.
type SpellPhantasmalForce struct {
	SourcePages            []int    `json:"source_pages"`
	Save                   string   `json:"save"`
	Phenomena              []string `json:"phenomena"`
	MaxCubeFeet            int      `json:"max_cube_feet"`
	PerceivableBy          string   `json:"perceivable_by"`
	Stimuli                []string `json:"stimuli"`
	ExamineAction          string   `json:"examine_action"`
	ExamineCheck           string   `json:"examine_check"`
	ExamineDC              string   `json:"examine_dc"`
	ExamineSuccess         string   `json:"examine_success"`
	Belief                 string   `json:"belief"`
	RationalizationExample string   `json:"rationalization_example"`
	DamageRequirement      string   `json:"damage_requirement"`
	DamageDiceCount        int      `json:"damage_dice_count"`
	DamageDieSides         int      `json:"damage_die_sides"`
	DamageType             string   `json:"damage_type"`
	DamageTiming           string   `json:"damage_timing"`
	DamageReachFeet        int      `json:"damage_reach_feet"`
	DamagePosition         string   `json:"damage_position"`
	DamagePerception       string   `json:"damage_perception"`
}

// SpellPhantasmalKiller retains typed SRD source facts without implementing them.
type SpellPhantasmalKiller struct {
	Illusion         string   `json:"illusion"`
	Save             string   `json:"save"`
	DamageDiceCount  int      `json:"damage_dice_count"`
	DamageDieSides   int      `json:"damage_die_sides"`
	DamageType       string   `json:"damage_type"`
	FailurePenalties []string `json:"failure_penalties"`
	SaveSuccess      string   `json:"save_success"`
	RepeatTiming     string   `json:"repeat_timing"`
	RepeatFailure    string   `json:"repeat_failure"`
	RepeatSuccess    string   `json:"repeat_success"`
	SlotDiceCount    int      `json:"slot_dice_count"`
	SlotDieSides     int      `json:"slot_die_sides"`
}

// SpellPhantomSteed retains typed SRD source facts without implementing them.
type SpellPhantomSteed struct {
	CreatureSize             string   `json:"creature_size"`
	CreatureNature           string   `json:"creature_nature"`
	Appearance               string   `json:"appearance"`
	Equipment                []string `json:"equipment"`
	EquipmentVanishes        string   `json:"equipment_vanishes"`
	EquipmentMaxDistanceFeet int      `json:"equipment_max_distance_feet"`
	Rider                    string   `json:"rider"`
	StatBlock                string   `json:"stat_block"`
	SpeedFeet                int      `json:"speed_feet"`
	TravelMilesPerHour       int      `json:"travel_miles_per_hour"`
	Fade                     string   `json:"fade"`
	DismountMinutes          int      `json:"dismount_minutes"`
	EarlyEnd                 string   `json:"early_end"`
}

func spell214Templates() map[string]any {
	return map[string]any{
		"mirror_image":       SpellMirrorImage{},
		"mislead":            SpellMislead{},
		"misty_step":         SpellMistyStep{},
		"modify_memory":      SpellModifyMemory{},
		"moonbeam":           SpellMoonbeam{},
		"move_earth":         SpellMoveEarth{},
		"nondetection":       SpellNondetection{},
		"pass_without_trace": SpellPassWithoutTrace{},
		"passwall":           SpellPasswall{},
		"phantasmal_force":   SpellPhantasmalForce{},
		"phantasmal_killer":  SpellPhantasmalKiller{},
		"phantom_steed":      SpellPhantomSteed{},
	}
}

func validSpell214Target(t SpellTarget) bool {
	switch t.Shape {
	case "square":
		return t.Count == 0 && t.Eligibility == "terrain-area-in-range"
	case "choice":
		return t.Count == 1 && t.Eligibility == "willing-creature-place-or-object-touch"
	case "emanation":
		return t.Count == 0 && t.Eligibility == "caster-and-chosen-creatures-in-aura"
	case "created-effect":
		return t.Count == 1 && t.Eligibility == "visible-point-on-wood-plaster-or-stone-surface"
	}
	return false
}

// validSpell214Scaling covers how old a memory Modify Memory reaches by slot.
func (s SpellRecord) validSpell214Scaling(effects map[string]bool) bool {
	sc, d := s.Scaling, s.SourceDetails
	return d != nil && d.ModifyMemory != nil && sc.Kind == "slot-memory-window" && sc.AboveLevel == s.Level && sc.DiceCount == 0 && sc.DieSides == 0 &&
		len(sc.EffectIDs) == 1 && effects[sc.EffectIDs[0]] && len(d.ModifyMemory.SlotWindows) == 3 && slices.Contains(s.RequiredMechanics, "scaling.memory-window-by-slot")
}

// has214Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has214Payload() bool {
	return d != nil && (d.MirrorImage != nil || d.Mislead != nil || d.MistyStep != nil || d.ModifyMemory != nil || d.Moonbeam != nil || d.MoveEarth != nil ||
		d.Nondetection != nil || d.PassWithoutTrace != nil || d.Passwall != nil || d.PhantasmalForce != nil || d.PhantasmalKiller != nil || d.PhantomSteed != nil)
}

// validateSpell214 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell214() (bool, error) {
	pin, ok := spell214Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has214Payload() {
			return true, fmt.Errorf("mismatched batch 214 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 214 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 214 record differs from the printed source")
	}
	return true, nil
}
