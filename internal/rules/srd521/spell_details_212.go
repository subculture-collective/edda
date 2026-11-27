package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellLocateObject retains typed SRD source facts without implementing them.
type SpellLocateObject struct {
	Query                 string   `json:"query"`
	Senses                string   `json:"senses"`
	SenseRadiusFeet       int      `json:"sense_radius_feet"`
	MovingTarget          string   `json:"moving_target"`
	SpecificRequirement   string   `json:"specific_requirement"`
	CloseDistanceFeet     int      `json:"close_distance_feet"`
	CloseSightingsMinimum int      `json:"close_sightings_minimum"`
	Alternative           string   `json:"alternative"`
	KindExamples          []string `json:"kind_examples"`
	LeadBlocker           string   `json:"lead_blocker"`
}

// SpellLongstrider retains typed SRD source facts without implementing them.
type SpellLongstrider struct {
	Target                   string `json:"target"`
	SpeedIncreaseFeet        int    `json:"speed_increase_feet"`
	BaseTargets              int    `json:"base_targets"`
	AdditionalTargetsPerSlot int    `json:"additional_targets_per_slot"`
}

// SpellMageArmor retains typed SRD source facts without implementing them.
type SpellMageArmor struct {
	Target             string `json:"target"`
	BaseArmorClass     int    `json:"base_armor_class"`
	ArmorClassModifier string `json:"armor_class_modifier"`
	EarlyEnd           string `json:"early_end"`
}

// SpellMageHand retains typed SRD source facts without implementing them.
type SpellMageHand struct {
	Hand               string   `json:"hand"`
	VanishDistanceFeet int      `json:"vanish_distance_feet"`
	VanishTriggers     []string `json:"vanish_triggers"`
	Uses               []string `json:"uses"`
	LaterAction        string   `json:"later_action"`
	LaterMoveFeet      int      `json:"later_move_feet"`
	Forbidden          []string `json:"forbidden"`
	CarryLimitPounds   int      `json:"carry_limit_pounds"`
}

// SpellMagicCircle retains typed SRD source facts without implementing them.
type SpellMagicCircle struct {
	SourcePages                      []int    `json:"source_pages"`
	CylinderRadiusFeet               int      `json:"cylinder_radius_feet"`
	CylinderHeightFeet               int      `json:"cylinder_height_feet"`
	Center                           string   `json:"center"`
	Runes                            string   `json:"runes"`
	CreatureTypes                    []string `json:"creature_types"`
	TypeChoice                       string   `json:"type_choice"`
	EntryRule                        string   `json:"entry_rule"`
	TravelMethods                    []string `json:"travel_methods"`
	TravelSave                       string   `json:"travel_save"`
	TravelSaveResult                 string   `json:"travel_save_result"`
	AttackRule                       string   `json:"attack_rule"`
	Protections                      []string `json:"protections"`
	ReverseOption                    string   `json:"reverse_option"`
	ReverseEffect                    string   `json:"reverse_effect"`
	AdditionalDurationSecondsPerSlot int      `json:"additional_duration_seconds_per_slot"`
}

// SpellMagicJar retains typed SRD source facts without implementing them.
type SpellMagicJar struct {
	BodyState                string   `json:"body_state"`
	SoulDestination          string   `json:"soul_destination"`
	ContainerAwareness       string   `json:"container_awareness"`
	ContainerRestrictions    []string `json:"container_restrictions"`
	ContainerOnlyAction      string   `json:"container_only_action"`
	ProjectionDistanceFeet   int      `json:"projection_distance_feet"`
	ProjectionChoices        []string `json:"projection_choices"`
	PossessionTarget         string   `json:"possession_target"`
	PossessionDistanceFeet   int      `json:"possession_distance_feet"`
	PossessionWards          []string `json:"possession_wards"`
	PossessionSave           string   `json:"possession_save"`
	PossessionFailure        string   `json:"possession_failure"`
	PossessionSuccess        string   `json:"possession_success"`
	PossessionRetryHours     int      `json:"possession_retry_hours"`
	HostControl              string   `json:"host_control"`
	ReplacedStatistics       []string `json:"replaced_statistics"`
	OtherStatistics          string   `json:"other_statistics"`
	TrappedSoul              []string `json:"trapped_soul"`
	ReturnAction             string   `json:"return_action"`
	ReturnDistanceFeet       int      `json:"return_distance_feet"`
	ReturnResult             string   `json:"return_result"`
	HostDeath                string   `json:"host_death"`
	HostDeathSaveDC          string   `json:"host_death_save_dc"`
	HostDeathSuccess         string   `json:"host_death_success"`
	HostDeathOtherwise       string   `json:"host_death_otherwise"`
	ContainerEndTriggers     []string `json:"container_end_triggers"`
	ContainerEndResult       string   `json:"container_end_result"`
	CasterDeathConditions    []string `json:"caster_death_conditions"`
	TrappedSoulOnDestruction string   `json:"trapped_soul_on_destruction"`
	TrappedSoulOtherwise     string   `json:"trapped_soul_otherwise"`
	SpellEndContainer        string   `json:"spell_end_container"`
}

// SpellMagicMissile retains typed SRD source facts without implementing them.
type SpellMagicMissile struct {
	DartCount              int      `json:"dart_count"`
	DartDiceCount          int      `json:"dart_dice_count"`
	DartDieSides           int      `json:"dart_die_sides"`
	DartBonus              int      `json:"dart_bonus"`
	DamageType             string   `json:"damage_type"`
	StrikeTiming           string   `json:"strike_timing"`
	Distribution           []string `json:"distribution"`
	AdditionalDartsPerSlot int      `json:"additional_darts_per_slot"`
}

// SpellMagicMouth retains typed SRD source facts without implementing them.
type SpellMagicMouth struct {
	MessageMaxWords           int      `json:"message_max_words"`
	MessageDeliveryMaxMinutes int      `json:"message_delivery_max_minutes"`
	TriggerDefinition         string   `json:"trigger_definition"`
	Mouth                     string   `json:"mouth"`
	Recitation                string   `json:"recitation"`
	MouthPlacement            string   `json:"mouth_placement"`
	MouthExample              string   `json:"mouth_example"`
	AfterDeliveryChoices      []string `json:"after_delivery_choices"`
	TriggerDetail             string   `json:"trigger_detail"`
	TriggerBasis              []string `json:"trigger_basis"`
	TriggerRangeFeet          int      `json:"trigger_range_feet"`
	TriggerExamples           []string `json:"trigger_examples"`
}

// SpellMagicWeaponSlotBonus is one source-listed row.
type SpellMagicWeaponSlotBonus struct {
	MinSlot int `json:"min_slot"`
	MaxSlot int `json:"max_slot"`
	Bonus   int `json:"bonus"`
}

// SpellMagicWeapon retains typed SRD source facts without implementing them.
type SpellMagicWeapon struct {
	SourcePages     []int                       `json:"source_pages"`
	Target          string                      `json:"target"`
	Becomes         string                      `json:"becomes"`
	AttackRollBonus int                         `json:"attack_roll_bonus"`
	DamageRollBonus int                         `json:"damage_roll_bonus"`
	EarlyEnd        string                      `json:"early_end"`
	SlotBonuses     []SpellMagicWeaponSlotBonus `json:"slot_bonuses"`
}

// SpellMagnificentMansion retains typed SRD source facts without implementing them.
type SpellMagnificentMansion struct {
	DoorAppearance          string   `json:"door_appearance"`
	DoorWidthFeet           int      `json:"door_width_feet"`
	DoorHeightFeet          int      `json:"door_height_feet"`
	DoorDestination         string   `json:"door_destination"`
	Entrants                string   `json:"entrants"`
	EntryRequirement        string   `json:"entry_requirement"`
	DoorControlAction       string   `json:"door_control_action"`
	DoorControlDistanceFeet int      `json:"door_control_distance_feet"`
	ClosedDoor              string   `json:"closed_door"`
	Foyer                   string   `json:"foyer"`
	Atmosphere              []string `json:"atmosphere"`
	FloorPlan               string   `json:"floor_plan"`
	MaxContiguousCubes      int      `json:"max_contiguous_cubes"`
	CubeFeet                int      `json:"cube_feet"`
	Furnishing              string   `json:"furnishing"`
	BanquetCourses          int      `json:"banquet_courses"`
	BanquetPeople           int      `json:"banquet_people"`
	RemovedObjects          string   `json:"removed_objects"`
	ServantCount            int      `json:"servant_count"`
	ServantAppearance       string   `json:"servant_appearance"`
	ServantTraits           []string `json:"servant_traits"`
	ServantTasks            string   `json:"servant_tasks"`
	ServantTaskExamples     []string `json:"servant_task_examples"`
	ServantLimits           []string `json:"servant_limits"`
	EndExpulsion            string   `json:"end_expulsion"`
}

// SpellMajorImage retains typed SRD source facts without implementing them.
type SpellMajorImage struct {
	ImageSubjects               []string `json:"image_subjects"`
	MaxCubeFeet                 int      `json:"max_cube_feet"`
	Qualities                   []string `json:"qualities"`
	Limits                      []string `json:"limits"`
	MoveRequirement             string   `json:"move_requirement"`
	MoveAction                  string   `json:"move_action"`
	MoveDestination             string   `json:"move_destination"`
	MovementAppearance          string   `json:"movement_appearance"`
	SoundChanges                string   `json:"sound_changes"`
	PhysicalInteraction         string   `json:"physical_interaction"`
	DiscernAction               string   `json:"discern_action"`
	DiscernCheck                string   `json:"discern_check"`
	DiscernDC                   string   `json:"discern_dc"`
	Discerned                   []string `json:"discerned"`
	UntilDispelledMinSlot       int      `json:"until_dispelled_min_slot"`
	UntilDispelledConcentration string   `json:"until_dispelled_concentration"`
}

// SpellMassCureWounds retains typed SRD source facts without implementing them.
type SpellMassCureWounds struct {
	Center           string `json:"center"`
	SphereRadiusFeet int    `json:"sphere_radius_feet"`
	MaxTargets       int    `json:"max_targets"`
	HealingDiceCount int    `json:"healing_dice_count"`
	HealingDieSides  int    `json:"healing_die_sides"`
	HealingAdds      string `json:"healing_adds"`
	SlotDiceCount    int    `json:"slot_dice_count"`
	SlotDieSides     int    `json:"slot_die_sides"`
}

func spell212Templates() map[string]any {
	return map[string]any{
		"locate_object":       SpellLocateObject{},
		"longstrider":         SpellLongstrider{},
		"mage_armor":          SpellMageArmor{},
		"mage_hand":           SpellMageHand{},
		"magic_circle":        SpellMagicCircle{},
		"magic_jar":           SpellMagicJar{},
		"magic_missile":       SpellMagicMissile{},
		"magic_mouth":         SpellMagicMouth{},
		"magic_weapon":        SpellMagicWeapon{},
		"magnificent_mansion": SpellMagnificentMansion{},
		"major_image":         SpellMajorImage{},
		"mass_cure_wounds":    SpellMassCureWounds{},
	}
}

func validSpell212Target(t SpellTarget) bool {
	switch t.Shape {
	case "single":
		return t.Count == 1 && slices.Contains([]string{"willing-unarmored-creature-touch", "visible-unworn-uncarried-object-in-range", "nonmagical-weapon-touch"}, t.Eligibility)
	case "cylinder":
		return t.Count == 0 && t.Eligibility == "circle-on-visible-ground-point-in-range"
	case "darts":
		return t.Count == 3 && t.Eligibility == "visible-creatures-in-range"
	case "created-effect":
		return t.Count == 1 && slices.Contains([]string{"door-in-range", "visible-spot-in-range"}, t.Eligibility)
	case "up-to-count":
		return t.Count == 6 && t.Eligibility == "creatures-in-sphere-on-visible-point"
	}
	return false
}

// spell212TargetScaling returns this batch's typed extra-target increments.
func spell212TargetScaling(d *SpellSourceDetails) (base, additional int, ok bool) {
	if d.Longstrider != nil {
		return d.Longstrider.BaseTargets, d.Longstrider.AdditionalTargetsPerSlot, true
	}
	return 0, 0, false
}

// spell212DurationScaling reports whether a typed slot table or increment backs
// this batch's slot-duration scaling.
func spell212DurationScaling(d *SpellSourceDetails, mechanics []string) bool {
	return (d.MagicCircle != nil && d.MagicCircle.AdditionalDurationSecondsPerSlot > 0 && slices.Contains(mechanics, "scaling.duration-per-slot"))
}

// validSpell212Scaling covers Magic Weapon's bonus by slot and Major Image's
// slot level at which it lasts until dispelled.
func (s SpellRecord) validSpell212Scaling(effects map[string]bool) bool {
	sc, d := s.Scaling, s.SourceDetails
	if d == nil || s.Level == 0 || sc.AboveLevel != s.Level || sc.DiceCount != 0 || sc.DieSides != 0 || len(sc.EffectIDs) != 1 || !effects[sc.EffectIDs[0]] {
		return false
	}
	switch sc.Kind {
	case "slot-bonus-table":
		return d.MagicWeapon != nil && len(d.MagicWeapon.SlotBonuses) == 2 && slices.Contains(s.RequiredMechanics, "scaling.bonus-by-slot")
	case "slot-until-dispelled":
		return d.MajorImage != nil && d.MajorImage.UntilDispelledMinSlot > s.Level && slices.Contains(s.RequiredMechanics, "scaling.until-dispelled-by-slot")
	}
	return false
}

// has212Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has212Payload() bool {
	return d != nil && (d.LocateObject != nil || d.Longstrider != nil || d.MageArmor != nil || d.MageHand != nil || d.MagicCircle != nil || d.MagicJar != nil ||
		d.MagicMissile != nil || d.MagicMouth != nil || d.MagicWeapon != nil || d.MagnificentMansion != nil || d.MajorImage != nil || d.MassCureWounds != nil)
}

// validateSpell212 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell212() (bool, error) {
	pin, ok := spell212Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has212Payload() {
			return true, fmt.Errorf("mismatched batch 212 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 212 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 212 record differs from the printed source")
	}
	return true, nil
}
