package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellPrismaticWallLayer is one source-listed row or nested block.
type SpellPrismaticWallLayer struct {
	Order             int    `json:"order"`
	Color             string `json:"color"`
	Failure           string `json:"failure"`
	DiceCount         int    `json:"dice_count"`
	DieSides          int    `json:"die_sides"`
	DamageType        string `json:"damage_type"`
	AdditionalEffect  string `json:"additional_effect"`
	DestroyDamage     int    `json:"destroy_damage"`
	DestroyDamageType string `json:"destroy_damage_type"`
}

// SpellPrismaticWall retains typed SRD source facts without implementing them.
type SpellPrismaticWall struct {
	Appearance              string                    `json:"appearance"`
	WallLengthFeet          int                       `json:"wall_length_feet"`
	WallHeightFeet          int                       `json:"wall_height_feet"`
	WallThicknessInches     int                       `json:"wall_thickness_inches"`
	WallOrientation         string                    `json:"wall_orientation"`
	GlobeDiameterFeet       int                       `json:"globe_diameter_feet"`
	ShapeChoices            []string                  `json:"shape_choices"`
	OccupiedSpace           string                    `json:"occupied_space"`
	BrightLightFeet         int                       `json:"bright_light_feet"`
	DimLightExtraFeet       int                       `json:"dim_light_extra_feet"`
	SafeCreatures           []string                  `json:"safe_creatures"`
	BlindnessDistanceFeet   int                       `json:"blindness_distance_feet"`
	BlindnessTriggers       []string                  `json:"blindness_triggers"`
	BlindnessSave           string                    `json:"blindness_save"`
	BlindnessSeconds        int                       `json:"blindness_seconds"`
	LayerCount              int                       `json:"layer_count"`
	Traversal               string                    `json:"traversal"`
	LayerSave               string                    `json:"layer_save"`
	ArmorClass              int                       `json:"armor_class"`
	DestructionOrder        string                    `json:"destruction_order"`
	DestroyedLayer          string                    `json:"destroyed_layer"`
	AntimagicField          string                    `json:"antimagic_field"`
	DispelMagic             string                    `json:"dispel_magic"`
	Layers                  []SpellPrismaticWallLayer `json:"layers"`
	IndigoCondition         string                    `json:"indigo_condition"`
	IndigoSave              string                    `json:"indigo_save"`
	IndigoSaveTiming        string                    `json:"indigo_save_timing"`
	IndigoSuccessesToEnd    int                       `json:"indigo_successes_to_end"`
	IndigoFailuresToPetrify int                       `json:"indigo_failures_to_petrify"`
	IndigoPetrifiedUntil    string                    `json:"indigo_petrified_until"`
	IndigoTracking          string                    `json:"indigo_tracking"`
	VioletCondition         string                    `json:"violet_condition"`
	VioletSave              string                    `json:"violet_save"`
	VioletSaveTiming        string                    `json:"violet_save_timing"`
	VioletSuccess           string                    `json:"violet_success"`
	VioletFailure           string                    `json:"violet_failure"`
}

// SpellPrivateSanctum retains typed SRD source facts without implementing them.
type SpellPrivateSanctum struct {
	MinCubeSideFeet           int      `json:"min_cube_side_feet"`
	MaxCubeSideFeet           int      `json:"max_cube_side_feet"`
	PropertyChoiceTiming      string   `json:"property_choice_timing"`
	ChoiceCount               string   `json:"choice_count"`
	Properties                []string `json:"properties"`
	PermanentCastingDays      int      `json:"permanent_casting_days"`
	PermanentRequirement      string   `json:"permanent_requirement"`
	PermanentResult           string   `json:"permanent_result"`
	AdditionalCubeFeetPerSlot int      `json:"additional_cube_feet_per_slot"`
}

// SpellProduceFlame retains typed SRD source facts without implementing them.
type SpellProduceFlame struct {
	FlameLocation     string                   `json:"flame_location"`
	FlameHeat         string                   `json:"flame_heat"`
	BrightLightFeet   int                      `json:"bright_light_feet"`
	DimLightExtraFeet int                      `json:"dim_light_extra_feet"`
	Recast            string                   `json:"recast"`
	HurlAction        string                   `json:"hurl_action"`
	HurlRangeFeet     int                      `json:"hurl_range_feet"`
	HurlTargets       []string                 `json:"hurl_targets"`
	Attack            string                   `json:"attack"`
	Upgrades          []SpellCantripUpgradeRow `json:"upgrades"`
}

// SpellProgrammedIllusion retains typed SRD source facts without implementing them.
type SpellProgrammedIllusion struct {
	Phenomena             []string `json:"phenomena"`
	DormantState          string   `json:"dormant_state"`
	MaxCubeFeet           int      `json:"max_cube_feet"`
	ScriptChoices         []string `json:"script_choices"`
	ScriptTiming          string   `json:"script_timing"`
	MaxPerformanceMinutes int      `json:"max_performance_minutes"`
	Activation            string   `json:"activation"`
	AfterPerformance      string   `json:"after_performance"`
	ResetMinutes          int      `json:"reset_minutes"`
	TriggerDetail         string   `json:"trigger_detail"`
	TriggerBasis          []string `json:"trigger_basis"`
	TriggerRangeFeet      int      `json:"trigger_range_feet"`
	TriggerExample        string   `json:"trigger_example"`
	PhysicalInteraction   string   `json:"physical_interaction"`
	ExamineAction         string   `json:"examine_action"`
	ExamineCheck          string   `json:"examine_check"`
	ExamineDC             string   `json:"examine_dc"`
	Discerned             string   `json:"discerned"`
}

// SpellProjectImage retains typed SRD source facts without implementing them.
type SpellProjectImage struct {
	SourcePages         []int    `json:"source_pages"`
	CopyOf              string   `json:"copy_of"`
	RangeMiles          int      `json:"range_miles"`
	Placement           string   `json:"placement"`
	Likeness            string   `json:"likeness"`
	Tangibility         string   `json:"tangibility"`
	DamageResult        string   `json:"damage_result"`
	Senses              string   `json:"senses"`
	MoveAction          string   `json:"move_action"`
	MoveMaxFeet         int      `json:"move_max_feet"`
	Behaviors           []string `json:"behaviors"`
	Mannerisms          string   `json:"mannerisms"`
	PhysicalInteraction string   `json:"physical_interaction"`
	ExamineAction       string   `json:"examine_action"`
	ExamineCheck        string   `json:"examine_check"`
	ExamineDC           string   `json:"examine_dc"`
	Discerned           string   `json:"discerned"`
}

// SpellProtectionFromEnergy retains typed SRD source facts without implementing them.
type SpellProtectionFromEnergy struct {
	Benefit           string   `json:"benefit"`
	DamageTypeChoices []string `json:"damage_type_choices"`
	ChooseCount       int      `json:"choose_count"`
}

// SpellProtectionFromEvilAndGood retains typed SRD source facts without implementing them.
type SpellProtectionFromEvilAndGood struct {
	ProtectedAgainst         []string `json:"protected_against"`
	AttackPenalty            string   `json:"attack_penalty"`
	ImmunitiesFromThoseTypes []string `json:"immunities_from_those_types"`
	ExistingEffectBenefit    string   `json:"existing_effect_benefit"`
}

// SpellProtectionFromPoison retains typed SRD source facts without implementing them.
type SpellProtectionFromPoison struct {
	ConditionEnded       string `json:"condition_ended"`
	SaveBenefit          string `json:"save_benefit"`
	ResistanceDamageType string `json:"resistance_damage_type"`
}

// SpellPurifyFoodAndDrink retains typed SRD source facts without implementing them.
type SpellPurifyFoodAndDrink struct {
	Removes    []string `json:"removes"`
	Eligible   string   `json:"eligible"`
	RadiusFeet int      `json:"radius_feet"`
	Center     string   `json:"center"`
}

// SpellRaiseDead retains typed SRD source facts without implementing them.
type SpellRaiseDead struct {
	MaxDaysDead                 int    `json:"max_days_dead"`
	DeathRequirement            string `json:"death_requirement"`
	ReturnHitPoints             int    `json:"return_hit_points"`
	PoisonBenefit               string `json:"poison_benefit"`
	Wounds                      string `json:"wounds"`
	MissingParts                string `json:"missing_parts"`
	AutomaticFailure            string `json:"automatic_failure"`
	FailureExample              string `json:"failure_example"`
	D20TestPenalty              int    `json:"d20_test_penalty"`
	PenaltyReductionPerLongRest int    `json:"penalty_reduction_per_long_rest"`
	PenaltyEndsAt               int    `json:"penalty_ends_at"`
}

// SpellRayOfEnfeeblement retains typed SRD source facts without implementing them.
type SpellRayOfEnfeeblement struct {
	Save                           string `json:"save"`
	SuccessResult                  string `json:"success_result"`
	SuccessWindow                  string `json:"success_window"`
	FailureResult                  string `json:"failure_result"`
	DamageRollSubtractionDiceCount int    `json:"damage_roll_subtraction_dice_count"`
	DamageRollSubtractionDieSides  int    `json:"damage_roll_subtraction_die_sides"`
	RepeatSaveTiming               string `json:"repeat_save_timing"`
	RepeatSuccess                  string `json:"repeat_success"`
}

// SpellRayOfFrost retains typed SRD source facts without implementing them.
type SpellRayOfFrost struct {
	Appearance         string                   `json:"appearance"`
	Attack             string                   `json:"attack"`
	SpeedReductionFeet int                      `json:"speed_reduction_feet"`
	SpeedReductionEnds string                   `json:"speed_reduction_ends"`
	Upgrades           []SpellCantripUpgradeRow `json:"upgrades"`
}

func spell216Templates() map[string]any {
	return map[string]any{
		"prismatic_wall":                SpellPrismaticWall{},
		"private_sanctum":               SpellPrivateSanctum{},
		"produce_flame":                 SpellProduceFlame{},
		"programmed_illusion":           SpellProgrammedIllusion{},
		"project_image":                 SpellProjectImage{},
		"protection_from_energy":        SpellProtectionFromEnergy{},
		"protection_from_evil_and_good": SpellProtectionFromEvilAndGood{},
		"protection_from_poison":        SpellProtectionFromPoison{},
		"purify_food_and_drink":         SpellPurifyFoodAndDrink{},
		"raise_dead":                    SpellRaiseDead{},
		"ray_of_enfeeblement":           SpellRayOfEnfeeblement{},
		"ray_of_frost":                  SpellRayOfFrost{},
	}
}

func validSpell216Target(t SpellTarget) bool {
	switch t.Shape {
	case "choice":
		return t.Count == 1 && t.Eligibility == "wall-or-globe-centered-in-range"
	case "cube":
		return t.Count == 0 && t.Eligibility == "warded-area-in-range"
	case "created-effect":
		return t.Count == 1 && slices.Contains([]string{"triggered-illusion-in-range", "previously-seen-location-in-range"}, t.Eligibility)
	case "sphere":
		return t.Count == 0 && t.Eligibility == "nonmagical-food-and-drink-in-sphere"
	case "single":
		return t.Count == 1 && t.Eligibility == "dead-creature-touch"
	}
	return false
}

// spell216CantripUpgrade reports whether this batch types the cantrip dice rows.
func spell216CantripUpgrade(d *SpellSourceDetails) bool {
	return (d.ProduceFlame != nil && len(d.ProduceFlame.Upgrades) == 3) || (d.RayOfFrost != nil && len(d.RayOfFrost.Upgrades) == 3)
}

// validSpell216Scaling covers Private Sanctum's larger Cube by slot.
func (s SpellRecord) validSpell216Scaling(effects map[string]bool) bool {
	sc, d := s.Scaling, s.SourceDetails
	return d != nil && d.PrivateSanctum != nil && sc.Kind == "slot-cube-size" && sc.AboveLevel == s.Level && sc.DiceCount == 0 && sc.DieSides == 0 &&
		len(sc.EffectIDs) == 1 && effects[sc.EffectIDs[0]] && d.PrivateSanctum.AdditionalCubeFeetPerSlot > 0 && slices.Contains(s.RequiredMechanics, "scaling.cube-size-per-slot")
}

// has216Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has216Payload() bool {
	return d != nil && (d.PrismaticWall != nil || d.PrivateSanctum != nil || d.ProduceFlame != nil || d.ProgrammedIllusion != nil || d.ProjectImage != nil || d.ProtectionFromEnergy != nil ||
		d.ProtectionFromEvilAndGood != nil || d.ProtectionFromPoison != nil || d.PurifyFoodAndDrink != nil || d.RaiseDead != nil || d.RayOfEnfeeblement != nil || d.RayOfFrost != nil)
}

// validateSpell216 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell216() (bool, error) {
	pin, ok := spell216Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has216Payload() {
			return true, fmt.Errorf("mismatched batch 216 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 216 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 216 record differs from the printed source")
	}
	return true, nil
}
