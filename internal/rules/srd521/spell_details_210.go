package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellHoldPerson retains typed SRD source facts without implementing them.
type SpellHoldPerson struct {
	Target                   string `json:"target"`
	Save                     string `json:"save"`
	FailureCondition         string `json:"failure_condition"`
	RepeatTiming             string `json:"repeat_timing"`
	RepeatSuccess            string `json:"repeat_success"`
	BaseTargets              int    `json:"base_targets"`
	AdditionalTargetsPerSlot int    `json:"additional_targets_per_slot"`
	AdditionalTargetType     string `json:"additional_target_type"`
}

// SpellHolyAura retains typed SRD source facts without implementing them.
type SpellHolyAura struct {
	EmanationFeet        int      `json:"emanation_feet"`
	Beneficiaries        string   `json:"beneficiaries"`
	SaveBenefit          string   `json:"save_benefit"`
	AttackPenalty        string   `json:"attack_penalty"`
	RetaliationTypes     []string `json:"retaliation_types"`
	RetaliationTrigger   string   `json:"retaliation_trigger"`
	RetaliationSave      string   `json:"retaliation_save"`
	RetaliationCondition string   `json:"retaliation_condition"`
	RetaliationDuration  string   `json:"retaliation_duration"`
}

// SpellHuntersMarkSlotDuration is one source-listed row.
type SpellHuntersMarkSlotDuration struct {
	MinSlot int `json:"min_slot"`
	MaxSlot int `json:"max_slot"`
	Seconds int `json:"seconds"`
}

// SpellHuntersMark retains typed SRD source facts without implementing them.
type SpellHuntersMark struct {
	Quarry               string                         `json:"quarry"`
	ExtraDamageDiceCount int                            `json:"extra_damage_dice_count"`
	ExtraDamageDieSides  int                            `json:"extra_damage_die_sides"`
	DamageType           string                         `json:"damage_type"`
	DamageTrigger        string                         `json:"damage_trigger"`
	FindChecks           []string                       `json:"find_checks"`
	FindBenefit          string                         `json:"find_benefit"`
	RetargetCondition    string                         `json:"retarget_condition"`
	RetargetAction       string                         `json:"retarget_action"`
	RetargetResult       string                         `json:"retarget_result"`
	SlotDurations        []SpellHuntersMarkSlotDuration `json:"slot_durations"`
}

// SpellHypnoticPattern retains typed SRD source facts without implementing them.
type SpellHypnoticPattern struct {
	CubeFeet         int      `json:"cube_feet"`
	Pattern          string   `json:"pattern"`
	Affected         string   `json:"affected"`
	Save             string   `json:"save"`
	FailureCondition string   `json:"failure_condition"`
	CharmedCondition string   `json:"charmed_condition"`
	CharmedSpeedFeet int      `json:"charmed_speed_feet"`
	EndsForCreature  []string `json:"ends_for_creature"`
}

// SpellIceKnife retains typed SRD source facts without implementing them.
type SpellIceKnife struct {
	Attack              string `json:"attack"`
	PiercingDiceCount   int    `json:"piercing_dice_count"`
	PiercingDieSides    int    `json:"piercing_die_sides"`
	ExplosionTiming     string `json:"explosion_timing"`
	ExplosionTargets    string `json:"explosion_targets"`
	ExplosionRadiusFeet int    `json:"explosion_radius_feet"`
	Save                string `json:"save"`
	ColdDiceCount       int    `json:"cold_dice_count"`
	ColdDieSides        int    `json:"cold_die_sides"`
	SaveSuccess         string `json:"save_success"`
	SlotColdDiceCount   int    `json:"slot_cold_dice_count"`
	SlotColdDieSides    int    `json:"slot_cold_die_sides"`
}

// SpellIceStorm retains typed SRD source facts without implementing them.
type SpellIceStorm struct {
	SourcePages              []int  `json:"source_pages"`
	CylinderRadiusFeet       int    `json:"cylinder_radius_feet"`
	CylinderHeightFeet       int    `json:"cylinder_height_feet"`
	Center                   string `json:"center"`
	Save                     string `json:"save"`
	BludgeoningDiceCount     int    `json:"bludgeoning_dice_count"`
	BludgeoningDieSides      int    `json:"bludgeoning_die_sides"`
	ColdDiceCount            int    `json:"cold_dice_count"`
	ColdDieSides             int    `json:"cold_die_sides"`
	SaveSuccess              string `json:"save_success"`
	Ground                   string `json:"ground"`
	GroundDuration           string `json:"ground_duration"`
	SlotBludgeoningDiceCount int    `json:"slot_bludgeoning_dice_count"`
	SlotBludgeoningDieSides  int    `json:"slot_bludgeoning_die_sides"`
}

// SpellIdentify retains typed SRD source facts without implementing them.
type SpellIdentify struct {
	Touch             string   `json:"touch"`
	ObjectEligibility string   `json:"object_eligibility"`
	ObjectReveals     []string `json:"object_reveals"`
	CreatureReveals   string   `json:"creature_reveals"`
}

// SpellIllusoryScript retains typed SRD source facts without implementing them.
type SpellIllusoryScript struct {
	Materials           []string `json:"materials"`
	DesignatedReaders   string   `json:"designated_readers"`
	DesignatedView      []string `json:"designated_view"`
	OthersView          string   `json:"others_view"`
	AlternativeAlters   []string `json:"alternative_alters"`
	AlternativeLanguage string   `json:"alternative_language"`
	Dispelled           string   `json:"dispelled"`
	Truesight           string   `json:"truesight"`
}

// SpellImprisonmentEffect is one source-listed row.
type SpellImprisonmentEffect struct {
	Name  string   `json:"name"`
	Rules []string `json:"rules"`
}

// SpellImprisonment retains typed SRD source facts without implementing them.
type SpellImprisonment struct {
	Save                 string                    `json:"save"`
	SaveSuccess          string                    `json:"save_success"`
	SuccessImmunityHours int                       `json:"success_immunity_hours"`
	Unneeded             []string                  `json:"unneeded"`
	Ages                 bool                      `json:"ages"`
	Divination           string                    `json:"divination"`
	Teleport             string                    `json:"teleport"`
	EffectChoiceCount    int                       `json:"effect_choice_count"`
	Effects              []SpellImprisonmentEffect `json:"effects"`
	MinimusHeightInches  int                       `json:"minimus_height_inches"`
	TriggerSpecified     string                    `json:"trigger_specified"`
	TriggerComplexity    string                    `json:"trigger_complexity"`
	TriggerAgreement     string                    `json:"trigger_agreement"`
	TriggerWindowYears   int                       `json:"trigger_window_years"`
	TriggerRequirement   string                    `json:"trigger_requirement"`
	TriggerExamples      []string                  `json:"trigger_examples"`
	DispelMagicSlotLevel int                       `json:"dispel_magic_slot_level"`
	DispelMagicTargets   []string                  `json:"dispel_magic_targets"`
}

// SpellIncendiaryCloud retains typed SRD source facts without implementing them.
type SpellIncendiaryCloud struct {
	SourcePages      []int    `json:"source_pages"`
	SphereRadiusFeet int      `json:"sphere_radius_feet"`
	Contents         []string `json:"contents"`
	Center           string   `json:"center"`
	Obscurement      string   `json:"obscurement"`
	DispersedBy      string   `json:"dispersed_by"`
	WindExample      string   `json:"wind_example"`
	Save             string   `json:"save"`
	DamageDiceCount  int      `json:"damage_dice_count"`
	DamageDieSides   int      `json:"damage_die_sides"`
	DamageType       string   `json:"damage_type"`
	SaveSuccess      string   `json:"save_success"`
	SaveTriggers     []string `json:"save_triggers"`
	SavesPerTurn     int      `json:"saves_per_turn"`
	DriftFeet        int      `json:"drift_feet"`
	DriftDirection   string   `json:"drift_direction"`
	DriftTiming      string   `json:"drift_timing"`
}

// SpellInflictWounds retains typed SRD source facts without implementing them.
type SpellInflictWounds struct {
	Target          string `json:"target"`
	Save            string `json:"save"`
	DamageDiceCount int    `json:"damage_dice_count"`
	DamageDieSides  int    `json:"damage_die_sides"`
	DamageType      string `json:"damage_type"`
	SaveSuccess     string `json:"save_success"`
	SlotDiceCount   int    `json:"slot_dice_count"`
	SlotDieSides    int    `json:"slot_die_sides"`
}

// SpellInsectPlague retains typed SRD source facts without implementing them.
type SpellInsectPlague struct {
	SphereRadiusFeet int      `json:"sphere_radius_feet"`
	Contents         string   `json:"contents"`
	Center           string   `json:"center"`
	Obscurement      string   `json:"obscurement"`
	Terrain          string   `json:"terrain"`
	Save             string   `json:"save"`
	DamageDiceCount  int      `json:"damage_dice_count"`
	DamageDieSides   int      `json:"damage_die_sides"`
	DamageType       string   `json:"damage_type"`
	SaveSuccess      string   `json:"save_success"`
	SaveTriggers     []string `json:"save_triggers"`
	SavesPerTurn     int      `json:"saves_per_turn"`
	SlotDiceCount    int      `json:"slot_dice_count"`
	SlotDieSides     int      `json:"slot_die_sides"`
}

func spell210Templates() map[string]any {
	return map[string]any{
		"hold_person":      SpellHoldPerson{},
		"holy_aura":        SpellHolyAura{},
		"hunters_mark":     SpellHuntersMark{},
		"hypnotic_pattern": SpellHypnoticPattern{},
		"ice_knife":        SpellIceKnife{},
		"ice_storm":        SpellIceStorm{},
		"identify":         SpellIdentify{},
		"illusory_script":  SpellIllusoryScript{},
		"imprisonment":     SpellImprisonment{},
		"incendiary_cloud": SpellIncendiaryCloud{},
		"inflict_wounds":   SpellInflictWounds{},
		"insect_plague":    SpellInsectPlague{},
	}
}

func validSpell210Target(t SpellTarget) bool {
	switch t.Shape {
	case "emanation":
		return t.Count == 0 && t.Eligibility == "chosen-creatures-in-aura"
	case "cube":
		return t.Count == 0 && t.Eligibility == "creatures-seeing-pattern-in-cube"
	case "cylinder":
		return t.Count == 0 && t.Eligibility == "creatures-in-cylinder-centered-in-range"
	case "choice":
		return t.Count == 1 && t.Eligibility == "object-or-creature-touched-throughout-casting"
	case "single":
		return t.Count == 1 && t.Eligibility == "writing-material-touched"
	}
	return false
}

// spell210TargetScaling returns this batch's typed extra-target increments.
func spell210TargetScaling(d *SpellSourceDetails) (base, additional int, ok bool) {
	if d.HoldPerson != nil {
		return d.HoldPerson.BaseTargets, d.HoldPerson.AdditionalTargetsPerSlot, true
	}
	return 0, 0, false
}

// spell210DurationScaling reports whether a typed slot table or increment backs
// this batch's slot-duration scaling.
func spell210DurationScaling(d *SpellSourceDetails, mechanics []string) bool {
	return (d.HuntersMark != nil && len(d.HuntersMark.SlotDurations) == 2 && slices.Contains(mechanics, "scaling.duration-by-slot"))
}

// has210Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has210Payload() bool {
	return d != nil && (d.HoldPerson != nil || d.HolyAura != nil || d.HuntersMark != nil || d.HypnoticPattern != nil || d.IceKnife != nil || d.IceStorm != nil ||
		d.Identify != nil || d.IllusoryScript != nil || d.Imprisonment != nil || d.IncendiaryCloud != nil || d.InflictWounds != nil || d.InsectPlague != nil)
}

// validateSpell210 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell210() (bool, error) {
	pin, ok := spell210Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has210Payload() {
			return true, fmt.Errorf("mismatched batch 210 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 210 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 210 record differs from the printed source")
	}
	return true, nil
}
