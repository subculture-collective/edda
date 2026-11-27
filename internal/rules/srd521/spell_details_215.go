package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellPlanarAllyPaymentRate is one source-listed row or nested block.
type SpellPlanarAllyPaymentRate struct {
	TaskUnit  string `json:"task_unit"`
	GPPerUnit int    `json:"gp_per_unit"`
}

// SpellPlanarAlly retains typed SRD source facts without implementing them.
type SpellPlanarAlly struct {
	EntityRequirement     string                       `json:"entity_requirement"`
	EntityExamples        []string                     `json:"entity_examples"`
	CreatureTypeChoices   []string                     `json:"creature_type_choices"`
	CreatureLoyalty       string                       `json:"creature_loyalty"`
	NamedRequest          string                       `json:"named_request"`
	Compulsion            string                       `json:"compulsion"`
	ServiceBasis          string                       `json:"service_basis"`
	TaskRange             string                       `json:"task_range"`
	TaskExamples          []string                     `json:"task_examples"`
	BargainingRequirement string                       `json:"bargaining_requirement"`
	PaymentExamples       []string                     `json:"payment_examples"`
	PaymentRates          []SpellPlanarAllyPaymentRate `json:"payment_rates"`
	MaxTaskDays           int                          `json:"max_task_days"`
	GMAdjustment          string                       `json:"gm_adjustment"`
	AlignedTaskPayment    string                       `json:"aligned_task_payment"`
	NonhazardousPayment   string                       `json:"nonhazardous_payment"`
	DangerousPayment      string                       `json:"dangerous_payment"`
	SuicidalTasks         string                       `json:"suicidal_tasks"`
	ReturnTriggers        []string                     `json:"return_triggers"`
	ReturnResult          string                       `json:"return_result"`
	NoAgreementResult     string                       `json:"no_agreement_result"`
}

// SpellPlanarBindingSlotDuration is one source-listed row or nested block.
type SpellPlanarBindingSlotDuration struct {
	Slot int `json:"slot"`
	Days int `json:"days"`
}

// SpellPlanarBinding retains typed SRD source facts without implementing them.
type SpellPlanarBinding struct {
	SourcePages          []int                            `json:"source_pages"`
	CreatureTypes        []string                         `json:"creature_types"`
	RangeRequirement     string                           `json:"range_requirement"`
	TypicalSetup         string                           `json:"typical_setup"`
	Save                 string                           `json:"save"`
	SaveTiming           string                           `json:"save_timing"`
	FailureResult        string                           `json:"failure_result"`
	SourceSpellExtension string                           `json:"source_spell_extension"`
	CommandObedience     string                           `json:"command_obedience"`
	CommandExamples      []string                         `json:"command_examples"`
	HostileBehavior      string                           `json:"hostile_behavior"`
	CompletionSamePlane  string                           `json:"completion_same_plane"`
	CompletionOtherPlane string                           `json:"completion_other_plane"`
	SlotDurations        []SpellPlanarBindingSlotDuration `json:"slot_durations"`
}

// SpellPlaneShift retains typed SRD source facts without implementing them.
type SpellPlaneShift struct {
	MaxAdditionalCreatures int      `json:"max_additional_creatures"`
	TravelerRequirement    string   `json:"traveler_requirement"`
	DestinationKind        string   `json:"destination_kind"`
	GeneralDestination     string   `json:"general_destination"`
	DestinationExamples    []string `json:"destination_examples"`
	CircleDestination      string   `json:"circle_destination"`
	SmallCircleOverflow    string   `json:"small_circle_overflow"`
}

// SpellPlantGrowthVariant is one source-listed row or nested block.
type SpellPlantGrowthVariant struct {
	Choice         string `json:"choice"`
	CastingTime    string `json:"casting_time"`
	CastingSeconds int    `json:"casting_seconds"`
}

// SpellPlantGrowth retains typed SRD source facts without implementing them.
type SpellPlantGrowth struct {
	Variants                  []SpellPlantGrowthVariant `json:"variants"`
	OvergrowthRadiusFeet      int                       `json:"overgrowth_radius_feet"`
	OvergrowthShape           string                    `json:"overgrowth_shape"`
	OvergrowthPlants          string                    `json:"overgrowth_plants"`
	MovementCostFeet          int                       `json:"movement_cost_feet"`
	MovementUnitFeet          int                       `json:"movement_unit_feet"`
	OvergrowthExclusions      string                    `json:"overgrowth_exclusions"`
	EnrichmentRadiusFeet      int                       `json:"enrichment_radius_feet"`
	EnrichmentDays            int                       `json:"enrichment_days"`
	EnrichmentYieldMultiplier int                       `json:"enrichment_yield_multiplier"`
	EnrichmentLimitPerYear    int                       `json:"enrichment_limit_per_year"`
}

// SpellCantripUpgradeRow is one source-listed row or nested block.
type SpellCantripUpgradeRow struct {
	Level     int `json:"level"`
	DiceCount int `json:"dice_count"`
	DieSides  int `json:"die_sides"`
}

// SpellPoisonSpray retains typed SRD source facts without implementing them.
type SpellPoisonSpray struct {
	Appearance string                   `json:"appearance"`
	Attack     string                   `json:"attack"`
	Upgrades   []SpellCantripUpgradeRow `json:"upgrades"`
}

// SpellPolymorph retains typed SRD source facts without implementing them.
type SpellPolymorph struct {
	SourcePages            []int    `json:"source_pages"`
	Save                   string   `json:"save"`
	FormCreatureType       string   `json:"form_creature_type"`
	ChallengeLimit         string   `json:"challenge_limit"`
	NoChallengeRatingLimit string   `json:"no_challenge_rating_limit"`
	Statistics             string   `json:"statistics"`
	Retained               []string `json:"retained"`
	StatBlockReference     string   `json:"stat_block_reference"`
	TemporaryHitPoints     string   `json:"temporary_hit_points"`
	TemporaryHitPointsEnd  string   `json:"temporary_hit_points_end"`
	EarlyEnd               string   `json:"early_end"`
	ActionLimit            string   `json:"action_limit"`
	Forbidden              []string `json:"forbidden"`
	Gear                   string   `json:"gear"`
	GearUse                string   `json:"gear_use"`
}

// SpellPowerWordHeal retains typed SRD source facts without implementing them.
type SpellPowerWordHeal struct {
	HitPoints       string   `json:"hit_points"`
	ConditionsEnded []string `json:"conditions_ended"`
	ProneOption     string   `json:"prone_option"`
}

// SpellPowerWordKill retains typed SRD source facts without implementing them.
type SpellPowerWordKill struct {
	HitPointThreshold int    `json:"hit_point_threshold"`
	AtOrBelowResult   string `json:"at_or_below_result"`
	AboveDiceCount    int    `json:"above_dice_count"`
	AboveDieSides     int    `json:"above_die_sides"`
	AboveDamageType   string `json:"above_damage_type"`
}

// SpellPowerWordStun retains typed SRD source facts without implementing them.
type SpellPowerWordStun struct {
	HitPointThreshold  int    `json:"hit_point_threshold"`
	AtOrBelowCondition string `json:"at_or_below_condition"`
	AboveSpeedFeet     int    `json:"above_speed_feet"`
	AboveDuration      string `json:"above_duration"`
	RepeatSave         string `json:"repeat_save"`
	RepeatSaveTiming   string `json:"repeat_save_timing"`
	RepeatSuccess      string `json:"repeat_success"`
}

// SpellPrayerOfHealing retains typed SRD source facts without implementing them.
type SpellPrayerOfHealing struct {
	MaxCreatures        int    `json:"max_creatures"`
	PresenceRequirement string `json:"presence_requirement"`
	RestBenefit         string `json:"rest_benefit"`
	HealingDiceCount    int    `json:"healing_dice_count"`
	HealingDieSides     int    `json:"healing_die_sides"`
	ReuseLimit          string `json:"reuse_limit"`
	SlotDiceCount       int    `json:"slot_dice_count"`
	SlotDieSides        int    `json:"slot_die_sides"`
}

// SpellPrestidigitationOption is one source-listed row or nested block.
type SpellPrestidigitationOption struct {
	Choice  string   `json:"choice"`
	Effect  string   `json:"effect"`
	Details []string `json:"details"`
	Lasts   string   `json:"lasts"`
}

// SpellPrestidigitation retains typed SRD source facts without implementing them.
type SpellPrestidigitation struct {
	Options                   []SpellPrestidigitationOption `json:"options"`
	ChooseCount               int                           `json:"choose_count"`
	MaxActiveNonInstantaneous int                           `json:"max_active_non_instantaneous"`
}

// SpellPrismaticSprayDamageRay is one source-listed row or nested block.
type SpellPrismaticSprayDamageRay struct {
	Roll       int    `json:"roll"`
	Color      string `json:"color"`
	DiceCount  int    `json:"dice_count"`
	DieSides   int    `json:"die_sides"`
	DamageType string `json:"damage_type"`
	Success    string `json:"success"`
}

// SpellPrismaticSpray retains typed SRD source facts without implementing them.
type SpellPrismaticSpray struct {
	SourcePages             []int                          `json:"source_pages"`
	ConeFeet                int                            `json:"cone_feet"`
	RayCount                int                            `json:"ray_count"`
	Save                    string                         `json:"save"`
	RayDieSides             int                            `json:"ray_die_sides"`
	DamageRays              []SpellPrismaticSprayDamageRay `json:"damage_rays"`
	IndigoRoll              int                            `json:"indigo_roll"`
	IndigoCondition         string                         `json:"indigo_condition"`
	IndigoSave              string                         `json:"indigo_save"`
	IndigoSaveTiming        string                         `json:"indigo_save_timing"`
	IndigoSuccessesToEnd    int                            `json:"indigo_successes_to_end"`
	IndigoFailuresToPetrify int                            `json:"indigo_failures_to_petrify"`
	IndigoPetrifiedUntil    string                         `json:"indigo_petrified_until"`
	IndigoTracking          string                         `json:"indigo_tracking"`
	VioletRoll              int                            `json:"violet_roll"`
	VioletCondition         string                         `json:"violet_condition"`
	VioletSave              string                         `json:"violet_save"`
	VioletSaveTiming        string                         `json:"violet_save_timing"`
	VioletSuccess           string                         `json:"violet_success"`
	VioletFailure           string                         `json:"violet_failure"`
	SpecialRoll             int                            `json:"special_roll"`
	SpecialResult           string                         `json:"special_result"`
	SpecialReroll           string                         `json:"special_reroll"`
}

func spell215Templates() map[string]any {
	return map[string]any{
		"planar_ally":       SpellPlanarAlly{},
		"planar_binding":    SpellPlanarBinding{},
		"plane_shift":       SpellPlaneShift{},
		"plant_growth":      SpellPlantGrowth{},
		"poison_spray":      SpellPoisonSpray{},
		"polymorph":         SpellPolymorph{},
		"power_word_heal":   SpellPowerWordHeal{},
		"power_word_kill":   SpellPowerWordKill{},
		"power_word_stun":   SpellPowerWordStun{},
		"prayer_of_healing": SpellPrayerOfHealing{},
		"prestidigitation":  SpellPrestidigitation{},
		"prismatic_spray":   SpellPrismaticSpray{},
	}
}

func validSpell215Target(t SpellTarget) bool {
	switch t.Shape {
	case "created-effect":
		return t.Count == 1 && t.Eligibility == "unoccupied-space-in-range-for-loyal-creature"
	case "single":
		return t.Count == 1 && t.Eligibility == "celestial-elemental-fey-or-fiend-in-range"
	case "selected-creatures":
		return t.Count == 0 && t.Eligibility == "caster-and-willing-creatures-linked-in-circle"
	case "choice":
		return t.Count == 1 && slices.Contains([]string{"overgrowth-or-enrichment-point-in-range", "one-chosen-minor-effect-in-range"}, t.Eligibility)
	case "up-to-count":
		return t.Count == 5 && t.Eligibility == "chosen-creatures-in-range-for-entire-casting"
	case "cone":
		return t.Count == 0 && t.Eligibility == "creatures-in-self-origin-cone"
	}
	return false
}

// spell215DurationScaling reports whether a typed slot table backs this batch's
// slot-duration scaling.
func spell215DurationScaling(d *SpellSourceDetails, mechanics []string) bool {
	return d.PlanarBinding != nil && len(d.PlanarBinding.SlotDurations) == 4 && slices.Contains(mechanics, "scaling.duration-by-slot")
}

// spell215CantripUpgrade reports whether this batch types the cantrip dice rows.
func spell215CantripUpgrade(d *SpellSourceDetails) bool {
	return d.PoisonSpray != nil && len(d.PoisonSpray.Upgrades) == 3
}

// has215Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has215Payload() bool {
	return d != nil && (d.PlanarAlly != nil || d.PlanarBinding != nil || d.PlaneShift != nil || d.PlantGrowth != nil || d.PoisonSpray != nil || d.Polymorph != nil ||
		d.PowerWordHeal != nil || d.PowerWordKill != nil || d.PowerWordStun != nil || d.PrayerOfHealing != nil || d.Prestidigitation != nil || d.PrismaticSpray != nil)
}

// validateSpell215 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell215() (bool, error) {
	pin, ok := spell215Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has215Payload() {
			return true, fmt.Errorf("mismatched batch 215 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 215 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 215 record differs from the printed source")
	}
	return true, nil
}
