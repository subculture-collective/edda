package srd521

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// SpellAntipathyModesEntry retains typed SRD source facts without implementing them.
type SpellAntipathyModesEntry struct {
	Choice            string `json:"choice"`
	Condition         string `json:"condition"`
	Movement          string `json:"movement"`
	WithinFeet        int    `json:"within_feet"`
	WithinRestriction string `json:"within_restriction"`
	DamageSave        string `json:"damage_save"`
}

// SpellAntipathy retains typed SRD source facts without implementing them.
type SpellAntipathy struct {
	SourcePages         []int                      `json:"source_pages"`
	Modes               []SpellAntipathyModesEntry `json:"modes"`
	ChooseCount         int                        `json:"choose_count"`
	TargetChoices       []string                   `json:"target_choices"`
	MaximumSize         string                     `json:"maximum_size"`
	CreatureKind        string                     `json:"creature_kind"`
	KindExamples        []string                   `json:"kind_examples"`
	Save                string                     `json:"save"`
	SaveTrigger         string                     `json:"save_trigger"`
	SaveDistanceFeet    int                        `json:"save_distance_feet"`
	EndSaveTiming       string                     `json:"end_save_timing"`
	EndSaveDistanceFeet int                        `json:"end_save_distance_feet"`
	EndSaveSuccess      string                     `json:"end_save_success"`
	ImmunitySeconds     int                        `json:"immunity_seconds"`
	AfterImmunity       string                     `json:"after_immunity"`
}

// SpellEye retains typed SRD source facts without implementing them.
type SpellEye struct {
	Invisible            bool   `json:"invisible"`
	Invulnerable         bool   `json:"invulnerable"`
	Hovers               bool   `json:"hovers"`
	Vision               string `json:"vision"`
	DarkvisionFeet       int    `json:"darkvision_feet"`
	MoveAction           string `json:"move_action"`
	MoveFeet             int    `json:"move_feet"`
	MoveDirection        string `json:"move_direction"`
	Blocked              string `json:"blocked"`
	MinimumOpeningInches int    `json:"minimum_opening_inches"`
}

// SpellHandOptionsEntry retains typed SRD source facts without implementing them.
type SpellHandOptionsEntry struct {
	Choice              string `json:"choice"`
	ReachFeet           int    `json:"reach_feet"`
	MaximumSize         string `json:"maximum_size"`
	Resolution          string `json:"resolution"`
	DiceCount           int    `json:"dice_count"`
	DieSides            int    `json:"die_sides"`
	DamageType          string `json:"damage_type"`
	AddCastingModifier  bool   `json:"add_casting_modifier"`
	PushBaseFeet        int    `json:"push_base_feet"`
	PushFeetPerModifier int    `json:"push_feet_per_modifier"`
	Outcome             string `json:"outcome"`
	Followup            string `json:"followup"`
}

// SpellHandSlotDiceEntry retains typed SRD source facts without implementing them.
type SpellHandSlotDiceEntry struct {
	EffectID  string `json:"effect_id"`
	DiceCount int    `json:"dice_count"`
	DieSides  int    `json:"die_sides"`
}

// SpellHand retains typed SRD source facts without implementing them.
type SpellHand struct {
	Size           string                   `json:"size"`
	Placement      string                   `json:"placement"`
	Appearance     string                   `json:"appearance"`
	ObjectAC       int                      `json:"object_ac"`
	HitPoints      string                   `json:"hit_points"`
	ZeroHitPoints  string                   `json:"zero_hit_points"`
	OccupiesSpace  bool                     `json:"occupies_space"`
	CommandTimings []string                 `json:"command_timings"`
	MoveFeet       int                      `json:"move_feet"`
	ChooseCount    int                      `json:"choose_count"`
	Options        []SpellHandOptionsEntry  `json:"options"`
	SlotDice       []SpellHandSlotDiceEntry `json:"slot_dice"`
}

// SpellLock retains typed SRD source facts without implementing them.
type SpellLock struct {
	Objects               []string `json:"objects"`
	MustBeClosed          bool     `json:"must_be_closed"`
	Unlock                string   `json:"unlock"`
	Bypass                []string `json:"bypass"`
	PasswordOptional      bool     `json:"password_optional"`
	PasswordDistanceFeet  int      `json:"password_distance_feet"`
	PasswordUnlockSeconds int      `json:"password_unlock_seconds"`
}

// SpellSword retains typed SRD source facts without implementing them.
type SpellSword struct {
	Appearance      string   `json:"appearance"`
	Attack          string   `json:"attack"`
	AttackTimings   []string `json:"attack_timings"`
	ReachFeet       int      `json:"reach_feet"`
	MoveAction      string   `json:"move_action"`
	MoveFeet        int      `json:"move_feet"`
	MoveDestination string   `json:"move_destination"`
	RepeatTargets   []string `json:"repeat_targets"`
}

// SpellMagicAuraVariantsEntry retains typed SRD source facts without implementing them.
type SpellMagicAuraVariantsEntry struct {
	Choice  string   `json:"choice"`
	Target  string   `json:"target"`
	Effect  string   `json:"effect"`
	Options []string `json:"options"`
}

// SpellMagicAura retains typed SRD source facts without implementing them.
type SpellMagicAura struct {
	SourcePages       []int                         `json:"source_pages"`
	Variants          []SpellMagicAuraVariantsEntry `json:"variants"`
	DetectionExample  string                        `json:"detection_example"`
	RepeatDays        int                           `json:"repeat_days"`
	RepeatRequirement string                        `json:"repeat_requirement"`
	RepeatResult      string                        `json:"repeat_result"`
}

// SpellAstral retains typed SRD source facts without implementing them.
type SpellAstral struct {
	MaximumOtherTargets    int      `json:"maximum_other_targets"`
	IncludesCaster         bool     `json:"includes_caster"`
	Willing                bool     `json:"willing"`
	JacinthMinimumCostCP   int      `json:"jacinth_minimum_cost_cp"`
	SilverBarMinimumCostCP int      `json:"silver_bar_minimum_cost_cp"`
	MaterialsPerTarget     bool     `json:"materials_per_target"`
	DestinationPlane       string   `json:"destination_plane"`
	AlreadyOnPlane         string   `json:"already_on_plane"`
	BodyState              string   `json:"body_state"`
	BodyCondition          string   `json:"body_condition"`
	UnneededSustenance     []string `json:"unneeded_sustenance"`
	BodyAges               bool     `json:"body_ages"`
	AstralForm             string   `json:"astral_form"`
	Cord                   string   `json:"cord"`
	CordVisibleFeet        int      `json:"cord_visible_feet"`
	CordCutCondition       string   `json:"cord_cut_condition"`
	CordCutResult          string   `json:"cord_cut_result"`
	LeavingPlane           string   `json:"leaving_plane"`
	SeparateEffects        bool     `json:"separate_effects"`
	ZeroHitPoints          []string `json:"zero_hit_points"`
	ZeroHitPointsResult    string   `json:"zero_hit_points_result"`
	DismissAction          string   `json:"dismiss_action"`
	DismissResult          string   `json:"dismiss_result"`
	EndResult              string   `json:"end_result"`
}

// SpellAuguryOmensEntry is one row of the source Omens table.
type SpellAuguryOmensEntry struct {
	Omen    string `json:"omen"`
	Results string `json:"results"`
}

// SpellAugury retains typed SRD source facts without implementing them.
type SpellAugury struct {
	Subject                   string                  `json:"subject"`
	WindowMinutes             int                     `json:"window_minutes"`
	OmenSource                string                  `json:"omen_source"`
	ChosenBy                  string                  `json:"chosen_by"`
	Omens                     []SpellAuguryOmensEntry `json:"omens"`
	Limitation                string                  `json:"limitation"`
	RepeatScope               string                  `json:"repeat_scope"`
	CumulativeNoAnswerPercent int                     `json:"cumulative_no_answer_percent"`
}

// SpellLife retains typed SRD source facts without implementing them.
type SpellLife struct {
	EmanationFeet            int      `json:"emanation_feet"`
	Affects                  []string `json:"affects"`
	Resistance               string   `json:"resistance"`
	HitPointMaximumReduction string   `json:"hit_point_maximum_reduction"`
	ReviveTrigger            string   `json:"revive_trigger"`
	ReviveHitPoints          int      `json:"revive_hit_points"`
}

// SpellAwaken retains typed SRD source facts without implementing them.
type SpellAwaken struct {
	SourcePages         []int    `json:"source_pages"`
	TargetChoices       []string `json:"target_choices"`
	MaximumIntelligence int      `json:"maximum_intelligence"`
	GainedIntelligence  int      `json:"gained_intelligence"`
	Language            string   `json:"language"`
	PlantBecomes        string   `json:"plant_becomes"`
	PlantGains          []string `json:"plant_gains"`
	PlantStatistics     string   `json:"plant_statistics"`
	StatisticsExamples  []string `json:"statistics_examples"`
	Condition           string   `json:"condition"`
	ConditionDays       int      `json:"condition_days"`
	ConditionEnds       []string `json:"condition_ends"`
	AfterCondition      string   `json:"after_condition"`
}

// SpellBane retains typed SRD source facts without implementing them.
type SpellBane struct {
	BaseTargets              int      `json:"base_targets"`
	AdditionalTargetsPerSlot int      `json:"additional_targets_per_slot"`
	Save                     string   `json:"save"`
	SaveSuccess              string   `json:"save_success"`
	Rolls                    []string `json:"rolls"`
	DieSides                 int      `json:"die_sides"`
	RollModifier             string   `json:"roll_modifier"`
	TargetChoice             string   `json:"target_choice"`
}

// SpellBanishment retains typed SRD source facts without implementing them.
type SpellBanishment struct {
	BaseTargets              int      `json:"base_targets"`
	AdditionalTargetsPerSlot int      `json:"additional_targets_per_slot"`
	Save                     string   `json:"save"`
	SaveSuccess              string   `json:"save_success"`
	Destination              string   `json:"destination"`
	Condition                string   `json:"condition"`
	Return                   string   `json:"return"`
	NoReturnTypes            []string `json:"no_return_types"`
	NoReturnCondition        string   `json:"no_return_condition"`
	NoReturnSeconds          int      `json:"no_return_seconds"`
	NoReturnDestination      string   `json:"no_return_destination"`
}

func spell197Templates() map[string]any {
	return map[string]any{
		"antipathy":  SpellAntipathy{},
		"eye":        SpellEye{},
		"hand":       SpellHand{},
		"lock":       SpellLock{},
		"sword":      SpellSword{},
		"magic_aura": SpellMagicAura{},
		"astral":     SpellAstral{},
		"augury":     SpellAugury{},
		"life":       SpellLife{},
		"awaken":     SpellAwaken{},
		"bane":       SpellBane{},
		"banishment": SpellBanishment{},
	}
}

// spell197CastingSeconds maps the long casting times this batch introduces to
// their exact source length. Other casting times keep zero casting seconds.
var spell197CastingSeconds = map[string]int{"1-minute": 60, "1-hour": 3600, "8-hours": 28800}

// validSpell197Target lists the target shapes introduced by this batch. They
// are source descriptions only and select no runtime targeting behavior.
func validSpell197Target(t SpellTarget) bool {
	switch t.Shape {
	case "choice":
		return t.Count == 1 && slices.Contains([]string{"huge-or-smaller-creature-or-object-in-range", "willing-creature-or-unattended-object-touch", "low-intelligence-beast-or-plant-creature-or-natural-plant-touch"}, t.Eligibility)
	case "created-effect":
		return t.Count == 1 && slices.Contains([]string{"point-in-range", "visible-unoccupied-space-in-range"}, t.Eligibility)
	case "single":
		return t.Count == 1 && t.Eligibility == "closed-door-window-gate-container-or-hatch-touch"
	case "self-and-up-to-count":
		return t.Count == 8 && t.Eligibility == "willing-creatures-in-range"
	case "emanation":
		return t.Count == 0 && t.Eligibility == "self-and-allies-in-aura"
	case "up-to-count":
		return t.Count == 3 && t.Eligibility == "visible-creatures-in-range"
	default:
		return false
	}
}

// spell197TargetScaling returns the typed base and per-slot target counts that
// a slot-targets scaling entry references for Bane and Banishment.
func spell197TargetScaling(d *SpellSourceDetails) (base, additional int, ok bool) {
	switch {
	case d.Bane != nil:
		return d.Bane.BaseTargets, d.Bane.AdditionalTargetsPerSlot, true
	case d.Banishment != nil:
		return d.Banishment.BaseTargets, d.Banishment.AdditionalTargetsPerSlot, true
	}
	return 0, 0, false
}

// validateSpell197EffectDice checks slot-effect-dice scaling: each referenced
// effect has its own per-slot dice in the typed Arcane Hand payload, in order.
func (s SpellRecord) validateSpell197EffectDice(effects map[string]bool) error {
	sc := s.Scaling
	if sc.AboveLevel != s.Level || sc.DiceCount != 0 || sc.DieSides != 0 || len(sc.EffectIDs) == 0 || !uniqueNames(sc.EffectIDs) || s.SourceDetails == nil || s.SourceDetails.Hand == nil || !slices.Contains(s.RequiredMechanics, "scaling.per-slot-by-effect") {
		return fmt.Errorf("invalid per-effect slot scaling")
	}
	rows := s.SourceDetails.Hand.SlotDice
	if len(rows) != len(sc.EffectIDs) {
		return fmt.Errorf("per-effect slot scaling requires one typed row per effect")
	}
	for i, row := range rows {
		if row.EffectID != sc.EffectIDs[i] || !effects[row.EffectID] || row.DiceCount <= 0 || !slices.Contains([]int{4, 6, 8, 10, 12, 20}, row.DieSides) {
			return fmt.Errorf("invalid per-effect slot dice")
		}
	}
	return nil
}

func (s SpellRecord) validateSpell197() (bool, error) {
	names := map[string]string{"antipathy": "antipathy-sympathy", "eye": "arcane-eye", "hand": "arcane-hand", "lock": "arcane-lock", "sword": "arcane-sword", "magic_aura": "arcanists-magic-aura", "astral": "astral-projection", "augury": "augury", "life": "aura-of-life", "awaken": "awaken", "bane": "bane", "banishment": "banishment"}
	expected := ""
	for key, name := range names {
		if s.ID == prefix+"spells/"+name {
			expected = key
		}
	}
	d := s.SourceDetails
	if d == nil {
		if expected != "" {
			return true, fmt.Errorf("batch 197 details required")
		}
		return false, nil
	}
	// The long casting times added here need their exact length on any record.
	if seconds := spell197CastingSeconds[s.CastingTime]; seconds != 0 && d.CastingSeconds != seconds {
		return true, fmt.Errorf("casting seconds differ from casting time")
	}
	payloads := map[string]any{
		"antipathy":  d.Antipathy,
		"eye":        d.Eye,
		"hand":       d.Hand,
		"lock":       d.Lock,
		"sword":      d.Sword,
		"magic_aura": d.MagicAura,
		"astral":     d.Astral,
		"augury":     d.Augury,
		"life":       d.Life,
		"awaken":     d.Awaken,
		"bane":       d.Bane,
		"banishment": d.Banishment,
	}
	count := 0
	for key, value := range payloads {
		if reflect.ValueOf(value).IsNil() {
			continue
		}
		count++
		if expected != key {
			return true, fmt.Errorf("mismatched batch 197 payload")
		}
		raw, err := json.Marshal(value)
		if err != nil {
			return true, err
		}
		if err := requireSpell198Fields(raw, reflect.TypeOf(value).Elem()); err != nil {
			return true, err
		}
	}
	if expected == "" && count == 0 {
		return false, nil
	}
	// Exactly one typed payload of any batch may accompany a record.
	total := 0
	all := reflect.ValueOf(*d)
	for i := 0; i < all.NumField(); i++ {
		if f := all.Field(i); f.Kind() == reflect.Pointer && !f.IsNil() {
			total++
		}
	}
	if count != 1 || total != 1 || s.Scope != "complete-source-record" || d.CastingTriggers == nil || len(d.CastingTriggers) != 0 || strings.TrimSpace(d.EndTiming) == "" {
		return true, fmt.Errorf("invalid batch 197 source details")
	}
	if d.CastingSeconds != spell197CastingSeconds[s.CastingTime] {
		return true, fmt.Errorf("casting seconds differ from casting time")
	}
	if a := d.Antipathy; a != nil {
		if !slices.Equal(a.SourcePages, []int{109, 110}) || len(a.Modes) != 2 || a.Modes[0].Choice != "antipathy" || a.Modes[1].Choice != "sympathy" || a.ChooseCount != 1 || !slices.Equal(a.TargetChoices, []string{"creature", "object"}) || a.Save != "wisdom" || a.SaveDistanceFeet <= 0 || a.EndSaveDistanceFeet <= 0 || a.ImmunitySeconds <= 0 {
			return true, fmt.Errorf("invalid antipathy/sympathy variants")
		}
	}
	if h := d.Hand; h != nil {
		choices := []string{}
		for _, o := range h.Options {
			choices = append(choices, o.Choice)
		}
		if !slices.Equal(choices, []string{"clenched-fist", "forceful-hand", "grasping-hand", "interposing-hand"}) || h.ChooseCount != 1 || h.ObjectAC <= 0 || h.MoveFeet <= 0 || s.Scaling.Kind != "slot-effect-dice" {
			return true, fmt.Errorf("invalid hand options")
		}
	}
	if m := d.MagicAura; m != nil {
		if !slices.Equal(m.SourcePages, []int{110, 111}) || len(m.Variants) != 2 || m.Variants[0].Choice != "mask" || m.Variants[1].Choice != "false-aura" || len(m.Variants[1].Options) != 3 || m.RepeatDays <= 0 {
			return true, fmt.Errorf("invalid magic aura variants")
		}
	}
	if a := d.Astral; a != nil {
		m := s.Components.Material
		if a.MaximumOtherTargets != s.Target.Count || a.JacinthMinimumCostCP <= 0 || a.SilverBarMinimumCostCP <= 0 || m == nil || *m.CostCP != a.JacinthMinimumCostCP+a.SilverBarMinimumCostCP || !*m.Consumed || !a.MaterialsPerTarget {
			return true, fmt.Errorf("invalid astral targets or materials")
		}
	}
	if a := d.Augury; a != nil {
		omens := []string{}
		for _, o := range a.Omens {
			omens = append(omens, o.Omen)
		}
		if !slices.Equal(omens, []string{"Weal", "Woe", "Weal and woe", "Indifference"}) || a.WindowMinutes <= 0 || a.CumulativeNoAnswerPercent <= 0 || a.CumulativeNoAnswerPercent > 100 || !*s.Ritual {
			return true, fmt.Errorf("invalid augury omens or chance")
		}
	}
	if l := d.Life; l != nil && (l.EmanationFeet <= 0 || l.ReviveHitPoints <= 0) {
		return true, fmt.Errorf("invalid aura of life area or healing")
	}
	if a := d.Awaken; a != nil {
		if !slices.Equal(a.SourcePages, []int{111, 112}) || !slices.Equal(a.TargetChoices, []string{"beast-creature", "plant-creature", "natural-noncreature-plant"}) || a.GainedIntelligence <= a.MaximumIntelligence || a.ConditionDays <= 0 {
			return true, fmt.Errorf("invalid awaken targets")
		}
	}
	if b := d.Bane; b != nil && (b.Save != "charisma" || b.DieSides != 4 || s.Scaling.Kind != "slot-targets") {
		return true, fmt.Errorf("invalid bane save, die or scaling")
	}
	if b := d.Banishment; b != nil && (b.Save != "charisma" || !slices.Equal(b.NoReturnTypes, []string{"aberration", "celestial", "elemental", "fey", "fiend"}) || b.NoReturnSeconds != s.DurationSeconds || s.Scaling.Kind != "slot-targets") {
		return true, fmt.Errorf("invalid banishment save, types or scaling")
	}
	return true, nil
}
