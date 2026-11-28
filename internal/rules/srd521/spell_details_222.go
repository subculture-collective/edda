package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellTongues retains typed SRD source facts without implementing them.
type SpellTongues struct {
	SourcePages           []int  `json:"source_pages"`
	Understands           string `json:"understands"`
	UnderstoodWhen        string `json:"understood_when"`
	UnderstoodBy          string `json:"understood_by"`
	UnderstoodRequirement string `json:"understood_requirement"`
}

// SpellTransportViaPlants retains typed SRD source facts without implementing them.
type SpellTransportViaPlants struct {
	MinimumPlantSize       string `json:"minimum_plant_size"`
	PlantKind              string `json:"plant_kind"`
	Destination            string `json:"destination"`
	DestinationRequirement string `json:"destination_requirement"`
	Travelers              string `json:"travelers"`
	MovementCostFeet       int    `json:"movement_cost_feet"`
}

// SpellTreeStride retains typed SRD source facts without implementing them.
type SpellTreeStride struct {
	DestinationTree        string   `json:"destination_tree"`
	RangeFeet              int      `json:"range_feet"`
	TreeRequirement        string   `json:"tree_requirement"`
	EnterMovementFeet      int      `json:"enter_movement_feet"`
	Awareness              string   `json:"awareness"`
	AwarenessFeet          int      `json:"awareness_feet"`
	EntryChoices           []string `json:"entry_choices"`
	ExitDistanceFeet       int      `json:"exit_distance_feet"`
	ExitMovementFeet       int      `json:"exit_movement_feet"`
	NoMovement             string   `json:"no_movement"`
	NoMovementDistanceFeet int      `json:"no_movement_distance_feet"`
	UsesPerTurn            int      `json:"uses_per_turn"`
	TurnEnd                string   `json:"turn_end"`
}

// SpellTruePolymorphCreatureForm is one source-listed row or nested block.
type SpellTruePolymorphCreatureForm struct {
	FormChoice            string   `json:"form_choice"`
	FormLimit             string   `json:"form_limit"`
	Statistics            string   `json:"statistics"`
	Retained              []string `json:"retained"`
	TemporaryHitPoints    string   `json:"temporary_hit_points"`
	TemporaryHitPointsEnd string   `json:"temporary_hit_points_end"`
	Limits                []string `json:"limits"`
	Gear                  string   `json:"gear"`
	GearUse               string   `json:"gear_use"`
}

// SpellTruePolymorphObjectCreature is one source-listed row or nested block.
type SpellTruePolymorphObjectCreature struct {
	CreatureKind         string `json:"creature_kind"`
	SizeLimit            string `json:"size_limit"`
	MaxChallengeRating   int    `json:"max_challenge_rating"`
	Attitude             string `json:"attitude"`
	TurnOrder            string `json:"turn_order"`
	Commands             string `json:"commands"`
	ControlSeconds       int    `json:"control_seconds"`
	AfterControl         string `json:"after_control"`
	AfterControlAttitude string `json:"after_control_attitude"`
}

// SpellTruePolymorphObjectForm is one source-listed row or nested block.
type SpellTruePolymorphObjectForm struct {
	Gear       string `json:"gear"`
	SizeLimit  string `json:"size_limit"`
	Statistics string `json:"statistics"`
	Memory     string `json:"memory"`
}

// SpellTruePolymorph retains typed SRD source facts without implementing them.
type SpellTruePolymorph struct {
	TargetChoices            []string                         `json:"target_choices"`
	CreatureResults          []string                         `json:"creature_results"`
	ObjectResult             string                           `json:"object_result"`
	ObjectRequirement        string                           `json:"object_requirement"`
	EndsEarly                []string                         `json:"ends_early"`
	FullConcentrationSeconds int                              `json:"full_concentration_seconds"`
	FullConcentrationResult  string                           `json:"full_concentration_result"`
	UnwillingSave            string                           `json:"unwilling_save"`
	UnwillingSaveSuccess     string                           `json:"unwilling_save_success"`
	CreatureIntoCreature     SpellTruePolymorphCreatureForm   `json:"creature_into_creature"`
	ObjectIntoCreature       SpellTruePolymorphObjectCreature `json:"object_into_creature"`
	CreatureIntoObject       SpellTruePolymorphObjectForm     `json:"creature_into_object"`
}

// SpellTrueResurrection retains typed SRD source facts without implementing them.
type SpellTrueResurrection struct {
	SourcePages         []int    `json:"source_pages"`
	MaxDeadYears        int      `json:"max_dead_years"`
	ExcludedDeath       string   `json:"excluded_death"`
	RevivedHitPoints    string   `json:"revived_hit_points"`
	Restores            []string `json:"restores"`
	Undead              string   `json:"undead"`
	NewBody             string   `json:"new_body"`
	NewBodyRequirement  string   `json:"new_body_requirement"`
	NewBodySpace        string   `json:"new_body_space"`
	NewBodyDistanceFeet int      `json:"new_body_distance_feet"`
}

// SpellTrueSeeing retains typed SRD source facts without implementing them.
type SpellTrueSeeing struct {
	Recipient     string `json:"recipient"`
	Sense         string `json:"sense"`
	TruesightFeet int    `json:"truesight_feet"`
}

// SpellTrueStrike retains typed SRD source facts without implementing them.
type SpellTrueStrike struct {
	Weapon               string                   `json:"weapon"`
	AttackCount          int                      `json:"attack_count"`
	Ability              string                   `json:"ability"`
	AbilityAppliesTo     []string                 `json:"ability_applies_to"`
	DamageTypeChoices    []string                 `json:"damage_type_choices"`
	DamageTypeChoiceWhen string                   `json:"damage_type_choice_when"`
	UpgradeDamageType    string                   `json:"upgrade_damage_type"`
	UpgradeApplies       string                   `json:"upgrade_applies"`
	Upgrades             []SpellCantripUpgradeRow `json:"upgrades"`
}

// SpellTsunami retains typed SRD source facts without implementing them.
type SpellTsunami struct {
	MaxLengthFeet                 int      `json:"max_length_feet"`
	MaxHeightFeet                 int      `json:"max_height_feet"`
	MaxThicknessFeet              int      `json:"max_thickness_feet"`
	AppearSave                    string   `json:"appear_save"`
	AppearDiceCount               int      `json:"appear_dice_count"`
	AppearDieSides                int      `json:"appear_die_sides"`
	AppearDamageType              string   `json:"appear_damage_type"`
	AppearSuccess                 string   `json:"appear_success"`
	MoveTiming                    string   `json:"move_timing"`
	MoveFeet                      int      `json:"move_feet"`
	MoveDirection                 string   `json:"move_direction"`
	Carried                       string   `json:"carried"`
	MovingMaxCreatureSize         string   `json:"moving_max_creature_size"`
	MovingTriggers                []string `json:"moving_triggers"`
	MovingSave                    string   `json:"moving_save"`
	MovingDiceCount               int      `json:"moving_dice_count"`
	MovingDieSides                int      `json:"moving_die_sides"`
	MovingDamageType              string   `json:"moving_damage_type"`
	MovingDamageLimit             string   `json:"moving_damage_limit"`
	ReductionTiming               string   `json:"reduction_timing"`
	HeightReductionFeet           int      `json:"height_reduction_feet"`
	LaterDamageReductionDiceCount int      `json:"later_damage_reduction_dice_count"`
	LaterDamageReductionDieSides  int      `json:"later_damage_reduction_die_sides"`
	EndsAtHeightFeet              int      `json:"ends_at_height_feet"`
	Swimming                      string   `json:"swimming"`
	SwimCheck                     string   `json:"swim_check"`
	SwimCheckDC                   string   `json:"swim_check_dc"`
	SwimCheckFailure              string   `json:"swim_check_failure"`
	LeavingWall                   string   `json:"leaving_wall"`
}

// SpellUnseenServant retains typed SRD source facts without implementing them.
type SpellUnseenServant struct {
	Traits           []string `json:"traits"`
	Size             string   `json:"size"`
	Kind             string   `json:"kind"`
	Purpose          string   `json:"purpose"`
	ArmorClass       int      `json:"armor_class"`
	HitPoints        int      `json:"hit_points"`
	Strength         int      `json:"strength"`
	CanAttack        bool     `json:"can_attack"`
	CommandAction    string   `json:"command_action"`
	CommandFrequency string   `json:"command_frequency"`
	CommandKind      string   `json:"command_kind"`
	MoveFeet         int      `json:"move_feet"`
	CommandIncludes  string   `json:"command_includes"`
	TaskLimit        string   `json:"task_limit"`
	TaskExamples     []string `json:"task_examples"`
	AfterCommand     string   `json:"after_command"`
	MaxDistanceFeet  int      `json:"max_distance_feet"`
	EarlyEnd         []string `json:"early_end"`
}

// SpellVampiricTouch retains typed SRD source facts without implementing them.
type SpellVampiricTouch struct {
	Attack          string `json:"attack"`
	AttackTarget    string `json:"attack_target"`
	DamageDiceCount int    `json:"damage_dice_count"`
	DamageDieSides  int    `json:"damage_die_sides"`
	DamageType      string `json:"damage_type"`
	Healing         string `json:"healing"`
	HealingDivisor  int    `json:"healing_divisor"`
	RepeatAction    string `json:"repeat_action"`
	RepeatTiming    string `json:"repeat_timing"`
	RepeatTarget    string `json:"repeat_target"`
	SlotDiceCount   int    `json:"slot_dice_count"`
	SlotDieSides    int    `json:"slot_die_sides"`
}

// SpellViciousMockery retains typed SRD source facts without implementing them.
type SpellViciousMockery struct {
	SourcePages        []int                    `json:"source_pages"`
	Save               string                   `json:"save"`
	DamageDiceCount    int                      `json:"damage_dice_count"`
	DamageDieSides     int                      `json:"damage_die_sides"`
	DamageType         string                   `json:"damage_type"`
	DisadvantageOn     string                   `json:"disadvantage_on"`
	DisadvantageWindow string                   `json:"disadvantage_window"`
	Upgrades           []SpellCantripUpgradeRow `json:"upgrades"`
}

// SpellVitriolicSphere retains typed SRD source facts without implementing them.
type SpellVitriolicSphere struct {
	BallDiameterFeet int    `json:"ball_diameter_feet"`
	RadiusFeet       int    `json:"radius_feet"`
	Save             string `json:"save"`
	InitialDiceCount int    `json:"initial_dice_count"`
	InitialDieSides  int    `json:"initial_die_sides"`
	DelayedDiceCount int    `json:"delayed_dice_count"`
	DelayedDieSides  int    `json:"delayed_die_sides"`
	DamageType       string `json:"damage_type"`
	DelayedTiming    string `json:"delayed_timing"`
	SaveSuccess      string `json:"save_success"`
	SlotAppliesTo    string `json:"slot_applies_to"`
	SlotDiceCount    int    `json:"slot_dice_count"`
	SlotDieSides     int    `json:"slot_die_sides"`
}

func spell222Templates() map[string]any {
	return map[string]any{
		"tongues":              SpellTongues{},
		"transport_via_plants": SpellTransportViaPlants{},
		"tree_stride":          SpellTreeStride{},
		"true_polymorph":       SpellTruePolymorph{},
		"true_resurrection":    SpellTrueResurrection{},
		"true_seeing":          SpellTrueSeeing{},
		"true_strike":          SpellTrueStrike{},
		"tsunami":              SpellTsunami{},
		"unseen_servant":       SpellUnseenServant{},
		"vampiric_touch":       SpellVampiricTouch{},
		"vicious_mockery":      SpellViciousMockery{},
		"vitriolic_sphere":     SpellVitriolicSphere{},
	}
}

func validSpell222Target(t SpellTarget) bool {
	switch t.Shape {
	case "single":
		return t.Count == 1 && slices.Contains([]string{"large-or-larger-inanimate-plant-in-range", "creature-dead-no-longer-than-200-years-touch", "visible-or-heard-creature-in-range"}, t.Eligibility)
	case "choice":
		return t.Count == 1 && t.Eligibility == "visible-creature-or-nonmagical-object-in-range"
	case "wall":
		return t.Count == 0 && t.Eligibility == "creatures-in-water-wall-on-point-in-range"
	case "created-effect":
		return t.Count == 1 && t.Eligibility == "unoccupied-ground-space-in-range"
	case "sphere":
		return t.Count == 0 && t.Eligibility == "creatures-in-acid-sphere-on-point-in-range"
	}
	return false
}

// spell222CantripUpgrade reports whether this batch types the cantrip dice rows.
func spell222CantripUpgrade(d *SpellSourceDetails) bool {
	return (d.TrueStrike != nil && len(d.TrueStrike.Upgrades) == 3) || (d.ViciousMockery != nil && len(d.ViciousMockery.Upgrades) == 3)
}

// has222Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has222Payload() bool {
	return d != nil && (d.Tongues != nil || d.TransportViaPlants != nil || d.TreeStride != nil || d.TruePolymorph != nil || d.TrueResurrection != nil || d.TrueSeeing != nil || d.TrueStrike != nil || d.Tsunami != nil || d.UnseenServant != nil || d.VampiricTouch != nil || d.ViciousMockery != nil || d.VitriolicSphere != nil)
}

// validateSpell222 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell222() (bool, error) {
	pin, ok := spell222Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has222Payload() {
			return true, fmt.Errorf("mismatched batch 222 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 222 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 222 record differs from the printed source")
	}
	return true, nil
}
