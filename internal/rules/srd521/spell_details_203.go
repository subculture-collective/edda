package srd521

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

// SpellDimensionDoor retains Dimension Door's typed SRD source facts without implementing them.
type SpellDimensionDoor struct {
	SourcePages         []int    `json:"source_pages"`
	DestinationChoices  []string `json:"destination_choices"`
	DirectionExamples   []string `json:"direction_examples"`
	Arrival             string   `json:"arrival"`
	CompanionCount      int      `json:"companion_count"`
	CompanionWilling    bool     `json:"companion_willing"`
	CompanionWithinFeet int      `json:"companion_within_feet"`
	CompanionArriveFeet int      `json:"companion_arrival_feet"`
	CompanionArrival    string   `json:"companion_arrival"`
	CollisionTrigger    string   `json:"collision_trigger"`
	CollisionDiceCount  int      `json:"collision_dice_count"`
	CollisionDieSides   int      `json:"collision_die_sides"`
	CollisionDamageType string   `json:"collision_damage_type"`
	CollisionDamaged    string   `json:"collision_damaged"`
	CollisionResult     string   `json:"collision_result"`
}

// SpellDisguiseSelf retains Disguise Self's typed SRD source facts without implementing them.
type SpellDisguiseSelf struct {
	Affects            []string `json:"affects"`
	HeightChangeFeet   int      `json:"height_change_feet"`
	HeightChanges      []string `json:"height_changes"`
	WeightChanges      []string `json:"weight_changes"`
	FormRequirement    string   `json:"form_requirement"`
	OtherExtent        string   `json:"other_extent"`
	PhysicalInspection string   `json:"physical_inspection"`
	InspectionExample  string   `json:"inspection_example"`
	DiscoveryAction    string   `json:"discovery_action"`
	DiscoveryCheck     string   `json:"discovery_check"`
	DiscoveryDC        string   `json:"discovery_dc"`
	DiscoveryResult    string   `json:"discovery_result"`
}

// SpellDisintegrate retains Disintegrate's typed SRD source facts without implementing them.
type SpellDisintegrate struct {
	TargetChoices       []string `json:"target_choices"`
	ForceCreationSample string   `json:"force_creation_example"`
	Save                string   `json:"save"`
	FailDiceCount       int      `json:"fail_dice_count"`
	FailDieSides        int      `json:"fail_die_sides"`
	FailBonus           int      `json:"fail_bonus"`
	DamageType          string   `json:"damage_type"`
	SuccessDamage       string   `json:"success_damage"`
	ZeroHitPointsResult string   `json:"zero_hit_points_result"`
	DustIncludes        string   `json:"dust_includes"`
	RevivedOnlyBy       []string `json:"revived_only_by"`
	AutomaticSize       string   `json:"automatic_size"`
	AutomaticTargets    []string `json:"automatic_targets"`
	OversizeSize        string   `json:"oversize_size"`
	OversizeResult      string   `json:"oversize_result"`
	SlotDiceCount       int      `json:"slot_dice_count"`
	SlotDieSides        int      `json:"slot_die_sides"`
}

// SpellDispelEvil retains Dispel Evil and Good's typed SRD source facts without implementing them.
type SpellDispelEvil struct {
	AffectedTypes    []string `json:"affected_types"`
	AttackEffect     string   `json:"attack_effect"`
	EarlyEnd         string   `json:"early_end"`
	BreakAction      string   `json:"break_action"`
	BreakTarget      string   `json:"break_target"`
	BreakResults     []string `json:"break_results"`
	DismissalAction  string   `json:"dismissal_action"`
	DismissalRange   int      `json:"dismissal_range_feet"`
	DismissalTarget  string   `json:"dismissal_target"`
	DismissalSave    string   `json:"dismissal_save"`
	DismissalFailure string   `json:"dismissal_failure"`
	UndeadElsewhere  string   `json:"undead_elsewhere"`
	FeyElsewhere     string   `json:"fey_elsewhere"`
}

// SpellDispelMagic retains Dispel Magic's typed SRD source facts without implementing them.
type SpellDispelMagic struct {
	TargetChoices  []string `json:"target_choices"`
	AutoEndUpToLvl int      `json:"auto_end_up_to_level"`
	CheckFromLevel int      `json:"check_from_level"`
	CheckAbility   string   `json:"check_ability"`
	CheckDCBase    int      `json:"check_dc_base"`
	CheckDCAdds    string   `json:"check_dc_adds"`
	CheckSuccess   string   `json:"check_success"`
	SlotAutoEnd    string   `json:"slot_auto_end"`
}

// SpellWhispers retains Dissonant Whispers' typed SRD source facts without implementing them.
type SpellWhispers struct {
	Target              string `json:"target"`
	Perception          string `json:"perception"`
	Save                string `json:"save"`
	FailDiceCount       int    `json:"fail_dice_count"`
	FailDieSides        int    `json:"fail_die_sides"`
	DamageType          string `json:"damage_type"`
	FailMovementTrigger string `json:"fail_movement_trigger"`
	FailMovement        string `json:"fail_movement"`
	FailRoute           string `json:"fail_route"`
	SuccessDamage       string `json:"success_damage"`
	SlotDiceCount       int    `json:"slot_dice_count"`
	SlotDieSides        int    `json:"slot_die_sides"`
}

// SpellDivination retains Divination's typed SRD source facts without implementing them.
type SpellDivination struct {
	Contact                  string   `json:"contact"`
	QuestionCount            int      `json:"question_count"`
	QuestionSubject          string   `json:"question_subject"`
	QuestionWindowDays       int      `json:"question_window_days"`
	ReplyBy                  string   `json:"reply_by"`
	ReplyTruth               string   `json:"reply_truth"`
	ReplyForms               []string `json:"reply_forms"`
	AccountsForChange        bool     `json:"accounts_for_changed_circumstances"`
	ChangedCircumstanceModel string   `json:"changed_circumstances_example"`
	RepeatReset              string   `json:"repeat_reset"`
	NoAnswerChancePercent    int      `json:"no_answer_chance_percent"`
	NoAnswerCumulative       bool     `json:"no_answer_cumulative"`
	NoAnswerAfter            string   `json:"no_answer_after"`
}

// SpellDivineFavor retains Divine Favor's typed SRD source facts without implementing them.
type SpellDivineFavor struct {
	AttackKind    string `json:"attack_kind"`
	Trigger       string `json:"trigger"`
	ExtraDiceCnt  int    `json:"extra_dice_count"`
	ExtraDieSides int    `json:"extra_die_sides"`
	DamageType    string `json:"damage_type"`
	Applies       string `json:"applies"`
}

// SpellDivineSmite retains Divine Smite's typed SRD source facts without implementing them.
type SpellDivineSmite struct {
	TriggerHits       []string `json:"trigger_hits"`
	TriggerTiming     string   `json:"trigger_timing"`
	ExtraDiceCount    int      `json:"extra_dice_count"`
	ExtraDieSides     int      `json:"extra_die_sides"`
	DamageType        string   `json:"damage_type"`
	TypeBonusDice     int      `json:"type_bonus_dice_count"`
	TypeBonusDieSides int      `json:"type_bonus_die_sides"`
	TypeBonusTypes    []string `json:"type_bonus_types"`
	SlotDiceCount     int      `json:"slot_dice_count"`
	SlotDieSides      int      `json:"slot_die_sides"`
}

// SpellDivineWordRow retains one printed row of the Divine Word Effects table.
type SpellDivineWordRow struct {
	MinHitPoints    int  `json:"min_hit_points"`
	MaxHitPoints    int  `json:"max_hit_points"`
	Dies            bool `json:"dies"`
	Blinded         bool `json:"blinded"`
	Deafened        bool `json:"deafened"`
	Stunned         bool `json:"stunned"`
	DurationSeconds int  `json:"duration_seconds"`
}

// SpellDivineWord retains Divine Word's typed SRD source facts without implementing them.
type SpellDivineWord struct {
	SourcePages      []int                `json:"source_pages"`
	Affects          string               `json:"affects"`
	Save             string               `json:"save"`
	HitPointLimit    int                  `json:"hit_point_limit"`
	LimitBasis       string               `json:"limit_basis"`
	Rows             []SpellDivineWordRow `json:"rows"`
	BanishTypes      []string             `json:"banish_types"`
	BanishRegardless bool                 `json:"banish_regardless_of_hit_points"`
	BanishResult     string               `json:"banish_result"`
	BanishReturnHrs  int                  `json:"banish_return_hours"`
	BanishException  string               `json:"banish_return_exception"`
}

// SpellDominateSlotDuration retains one printed concentration length for a slot level.
type SpellDominateSlotDuration struct {
	SlotLevel       int  `json:"slot_level"`
	DurationSeconds int  `json:"duration_seconds"`
	OrHigher        bool `json:"or_higher"`
}

// SpellDominate retains Dominate Beast's and Dominate Monster's typed SRD source facts.
type SpellDominate struct {
	SourcePages     []int                       `json:"source_pages"`
	TargetKind      string                      `json:"target_kind"`
	Save            string                      `json:"save"`
	FailCondition   string                      `json:"failure_condition"`
	ConditionLength string                      `json:"condition_duration"`
	SaveAdvantage   string                      `json:"save_advantage"`
	RepeatSave      string                      `json:"repeat_save"`
	RepeatSuccess   string                      `json:"repeat_success"`
	Link            string                      `json:"link"`
	CommandTiming   string                      `json:"command_timing"`
	CommandAction   string                      `json:"command_action"`
	CommandExamples []string                    `json:"command_examples"`
	Obey            string                      `json:"obey"`
	Idle            string                      `json:"idle"`
	ReactionCommand string                      `json:"reaction_command"`
	ReactionCost    string                      `json:"reaction_cost"`
	SlotDurations   []SpellDominateSlotDuration `json:"slot_durations"`
}

func spell203Templates() map[string]any {
	return map[string]any{
		"dimension_door":     SpellDimensionDoor{},
		"disguise_self":      SpellDisguiseSelf{},
		"disintegrate":       SpellDisintegrate{},
		"dispel_evil":        SpellDispelEvil{},
		"dispel_magic":       SpellDispelMagic{},
		"dissonant_whispers": SpellWhispers{},
		"divination":         SpellDivination{},
		"divine_favor":       SpellDivineFavor{},
		"divine_smite":       SpellDivineSmite{},
		"divine_word":        SpellDivineWord{},
		"dominate_beast":     SpellDominate{},
		"dominate_monster":   SpellDominate{},
	}
}

func validSpell203Target(t SpellTarget) bool {
	return t.Shape == "choice" && t.Count == 1 && slices.Contains([]string{"destination-in-range", "visible-creature-nonmagical-object-or-force-creation-in-range", "creature-object-or-magical-effect-in-range"}, t.Eligibility)
}

// validSpell203Scaling checks the two scaling kinds introduced by this batch.
// Each references a typed payload and a named, unavailable dependency.
func (s SpellRecord) validSpell203Scaling(effects map[string]bool) bool {
	sc, d := s.Scaling, s.SourceDetails
	if d == nil || s.Level == 0 || sc.AboveLevel != s.Level || sc.DiceCount != 0 || sc.DieSides != 0 || len(sc.EffectIDs) == 0 || !uniqueNames(sc.EffectIDs) {
		return false
	}
	for _, id := range sc.EffectIDs {
		if !effects[id] {
			return false
		}
	}
	switch sc.Kind {
	case "slot-auto-end":
		return d.DispelMagic != nil && d.DispelMagic.SlotAutoEnd != "" && slices.Contains(s.RequiredMechanics, "scaling.auto-end-by-slot-level")
	case "slot-duration-table":
		table := func(p *SpellDominate) bool { return p != nil && len(p.SlotDurations) > 0 }
		return (table(d.DominateBeast) || table(d.DominateMonster)) && slices.Contains(s.RequiredMechanics, "scaling.duration-by-slot-table")
	}
	return false
}

type spell203Effect struct {
	id, mechanic string
	dice, sides  int
	damage       string
	bonus        int
}

// spell203Row pins every header value, scaling row, effect and dependency a
// batch record must carry, by equality with the printed SRD page.
type spell203Row struct {
	slug, name            string
	page, level           int
	school                string
	lists                 []string
	casting               string
	ritual, concentration bool
	rangeKind             string
	rangeFeet             int
	target                SpellTarget
	duration              string
	seconds               int
	endTiming             string
	verbal, somatic       bool
	material              string
	costCP                int
	consumed              bool
	triggers              []string
	scaling               SpellScaling
	effects               []spell203Effect
	required              []string
}

var spell203None = SpellScaling{"none", 0, 0, 0, []string{}}

var spell203Rows = map[string]spell203Row{
	"dimension_door": {"dimension-door", "Dimension Door", 123, 4, "conjuration", []string{"bard", "sorcerer", "warlock", "wizard"}, "action", false, false, "feet", 500, SpellTarget{"choice", 1, "destination-in-range"}, "instantaneous", 0, "instantaneous", true, false, "", 0, false, []string{},
		spell203None, []spell203Effect{{"teleport", "teleport.dimension-door", 0, 0, "", 0}, {"companion", "teleport.dimension-door-companion", 0, 0, "", 0}, {"collision", "damage.dimension-door-collision", 4, 6, "force", 0}},
		[]string{"teleport.dimension-door", "teleport.dimension-door-companion", "damage.dimension-door-collision"}},
	"disguise_self": {"disguise-self", "Disguise Self", 124, 1, "illusion", []string{"bard", "sorcerer", "wizard"}, "action", false, false, "self", 0, SpellTarget{"single", 1, "self"}, "timed", 3600, "duration-ended", true, true, "", 0, false, []string{},
		spell203None, []spell203Effect{{"disguise", "illusion.disguise-self", 0, 0, "", 0}, {"inspection", "check.disguise-self-investigation", 0, 0, "", 0}},
		[]string{"illusion.disguise-self", "check.disguise-self-investigation"}},
	"disintegrate": {"disintegrate", "Disintegrate", 124, 6, "transmutation", []string{"sorcerer", "wizard"}, "action", false, false, "feet", 60, SpellTarget{"choice", 1, "visible-creature-nonmagical-object-or-force-creation-in-range"}, "instantaneous", 0, "instantaneous", true, true, "a lodestone and dust", 0, false, []string{},
		SpellScaling{"slot", 6, 3, 6, []string{"damage"}}, []spell203Effect{{"damage", "damage.disintegrate", 10, 6, "force", 40}, {"dust", "disintegrate.zero-hit-points-gray-dust", 0, 0, "", 0}, {"object", "disintegrate.automatic-object-or-force-creation", 0, 0, "", 0}},
		[]string{"damage.disintegrate", "disintegrate.zero-hit-points-gray-dust", "disintegrate.automatic-object-or-force-creation", "save.dexterity", "scaling.per-slot", "components.material"}},
	"dispel_evil": {"dispel-evil-and-good", "Dispel Evil and Good", 124, 5, "abjuration", []string{"cleric", "paladin"}, "action", false, true, "self", 0, SpellTarget{"single", 1, "self"}, "timed", 60, "duration-or-concentration-ended", true, true, "powdered silver and iron", 0, false, []string{},
		spell203None, []spell203Effect{{"disadvantage", "disadvantage.attacks-against-caster-by-listed-types", 0, 0, "", 0}, {"break", "special.break-enchantment", 0, 0, "", 0}, {"dismissal", "special.dismissal", 0, 0, "", 0}},
		[]string{"disadvantage.attacks-against-caster-by-listed-types", "special.break-enchantment", "special.dismissal", "save.charisma", "concentration.source-duration", "components.material"}},
	"dispel_magic": {"dispel-magic", "Dispel Magic", 124, 3, "abjuration", []string{"bard", "cleric", "druid", "paladin", "ranger", "sorcerer", "warlock", "wizard"}, "action", false, false, "feet", 120, SpellTarget{"choice", 1, "creature-object-or-magical-effect-in-range"}, "instantaneous", 0, "instantaneous", true, true, "", 0, false, []string{},
		SpellScaling{"slot-auto-end", 3, 0, 0, []string{"check"}}, []spell203Effect{{"end", "dispel.ongoing-spell-level-3-or-lower", 0, 0, "", 0}, {"check", "check.dispel-ongoing-spell-level-4-or-higher", 0, 0, "", 0}},
		[]string{"dispel.ongoing-spell-level-3-or-lower", "check.dispel-ongoing-spell-level-4-or-higher", "scaling.auto-end-by-slot-level"}},
	"dissonant_whispers": {"dissonant-whispers", "Dissonant Whispers", 124, 1, "enchantment", []string{"bard"}, "action", false, false, "feet", 60, SpellTarget{"single", 1, "visible-creature-in-range"}, "instantaneous", 0, "instantaneous", true, false, "", 0, false, []string{},
		SpellScaling{"slot", 1, 1, 6, []string{"damage"}}, []spell203Effect{{"damage", "damage.dissonant-whispers", 3, 6, "psychic", 0}, {"flee", "movement.dissonant-whispers-reaction-flee", 0, 0, "", 0}},
		[]string{"damage.dissonant-whispers", "movement.dissonant-whispers-reaction-flee", "save.wisdom-half", "scaling.per-slot"}},
	"divination": {"divination", "Divination", 125, 4, "divination", []string{"cleric", "druid", "wizard"}, "action", true, false, "self", 0, SpellTarget{"single", 1, "self"}, "instantaneous", 0, "instantaneous", true, true, "incense worth 25+ GP, which the spell consumes", 2500, true, []string{},
		spell203None, []spell203Effect{{"contact", "divination.contact-god", 0, 0, "", 0}, {"answer", "divination.cumulative-no-answer-chance", 0, 0, "", 0}},
		[]string{"divination.contact-god", "divination.cumulative-no-answer-chance", "cast.ritual", "components.material"}},
	"divine_favor": {"divine-favor", "Divine Favor", 125, 1, "transmutation", []string{"paladin"}, "bonus-action", false, false, "self", 0, SpellTarget{"single", 1, "self"}, "timed", 60, "duration-ended", true, true, "", 0, false, []string{},
		spell203None, []spell203Effect{{"damage", "damage.divine-favor-weapon-hit", 1, 4, "radiant", 0}},
		[]string{"damage.divine-favor-weapon-hit"}},
	"divine_smite": {"divine-smite", "Divine Smite", 125, 1, "evocation", []string{"paladin"}, "bonus-action", false, false, "self", 0, SpellTarget{"single", 1, "hit-target"}, "instantaneous", 0, "instantaneous", true, false, "", 0, false, []string{"hit-with-melee-weapon", "hit-with-unarmed-strike"},
		SpellScaling{"slot", 1, 1, 8, []string{"damage"}}, []spell203Effect{{"damage", "damage.divine-smite", 2, 8, "radiant", 0}, {"bonus", "damage.divine-smite-fiend-undead", 1, 8, "radiant", 0}},
		[]string{"damage.divine-smite", "damage.divine-smite-fiend-undead", "trigger.bonus-action-after-melee-or-unarmed-hit", "scaling.per-slot"}},
	"divine_word": {"divine-word", "Divine Word", 125, 7, "evocation", []string{"cleric"}, "bonus-action", false, false, "feet", 30, SpellTarget{"selected-creatures", 0, "creatures-in-range"}, "instantaneous", 0, "instantaneous", true, false, "", 0, false, []string{},
		spell203None, []spell203Effect{{"table", "effect.divine-word-hit-point-table", 0, 0, "", 0}, {"banish", "banish.divine-word-planar", 0, 0, "", 0}},
		[]string{"effect.divine-word-hit-point-table", "banish.divine-word-planar", "save.charisma"}},
	"dominate_beast": {"dominate-beast", "Dominate Beast", 125, 4, "enchantment", []string{"druid", "ranger", "sorcerer"}, "action", false, true, "feet", 60, SpellTarget{"single", 1, "visible-beast-in-range"}, "timed", 60, "duration-or-concentration-ended", true, true, "", 0, false, []string{},
		SpellScaling{"slot-duration-table", 4, 0, 0, []string{"charm"}}, spell203DominateEffects, spell203DominateRequired},
	"dominate_monster": {"dominate-monster", "Dominate Monster", 125, 8, "enchantment", []string{"bard", "sorcerer", "warlock", "wizard"}, "action", false, true, "feet", 60, SpellTarget{"single", 1, "visible-creature-in-range"}, "timed", 3600, "duration-or-concentration-ended", true, true, "", 0, false, []string{},
		SpellScaling{"slot-duration-table", 8, 0, 0, []string{"charm"}}, spell203DominateEffects, spell203DominateRequired},
}

var spell203DominateEffects = []spell203Effect{{"charm", "condition.charmed", 0, 0, "", 0}, {"link", "link.dominate-telepathic-command", 0, 0, "", 0}, {"reaction", "command.dominate-reaction", 0, 0, "", 0}}

var spell203DominateRequired = []string{"condition.charmed", "link.dominate-telepathic-command", "command.dominate-reaction", "save.wisdom", "save.advantage-when-fighting", "save.repeat-on-damage", "scaling.duration-by-slot-table", "concentration.source-duration"}

func (s SpellRecord) validateSpell203() (bool, error) {
	expected := ""
	for key, row := range spell203Rows {
		if s.ID == prefix+"spells/"+row.slug {
			expected = key
		}
	}
	if s.SourceDetails == nil {
		if expected != "" {
			return true, fmt.Errorf("batch 203 details required")
		}
		return false, nil
	}
	d := s.SourceDetails
	payloads := map[string]any{
		"dimension_door": d.DimensionDoor, "disguise_self": d.DisguiseSelf, "disintegrate": d.Disintegrate, "dispel_evil": d.DispelEvil,
		"dispel_magic": d.DispelMagic, "dissonant_whispers": d.DissonantWhispers, "divination": d.Divination, "divine_favor": d.DivineFavor,
		"divine_smite": d.DivineSmite, "divine_word": d.DivineWord, "dominate_beast": d.DominateBeast, "dominate_monster": d.DominateMonster,
	}
	count := 0
	for key, value := range payloads {
		if reflect.ValueOf(value).IsNil() {
			continue
		}
		count++
		if expected != key {
			return true, fmt.Errorf("mismatched batch 203 payload")
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
	row := spell203Rows[expected]
	if count != 1 || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" || d.CastingSeconds != 0 || d.CastingTriggers == nil || !slices.Equal(d.CastingTriggers, row.triggers) {
		return true, fmt.Errorf("invalid batch 203 source details")
	}
	if err := s.validateSpell203Header(row); err != nil {
		return true, err
	}
	if err := d.validateSpell203Payload(expected); err != nil {
		return true, err
	}
	return true, nil
}

func (s SpellRecord) validateSpell203Header(row spell203Row) error {
	flags := func(p *bool) bool { return p != nil && *p }
	c := s.Components
	if s.Name != row.name || s.SourcePage != row.page || s.Level != row.level || s.School != row.school || !slices.Equal(s.Lists, row.lists) ||
		s.CastingTime != row.casting || flags(s.Ritual) != row.ritual || flags(s.Concentration) != row.concentration ||
		s.Range != (SpellRange{row.rangeKind, row.rangeFeet}) || s.Target != row.target || s.Duration != row.duration || s.DurationSeconds != row.seconds ||
		s.SourceDetails.EndTiming != row.endTiming || flags(c.Verbal) != row.verbal || flags(c.Somatic) != row.somatic {
		return fmt.Errorf("batch 203 header differs from source")
	}
	if m := c.Material; row.material == "" {
		if m != nil {
			return fmt.Errorf("batch 203 material differs from source")
		}
	} else if m == nil || m.Description != row.material || m.CostCP == nil || *m.CostCP != row.costCP || m.Consumed == nil || *m.Consumed != row.consumed {
		return fmt.Errorf("batch 203 material differs from source")
	}
	sc := s.Scaling
	if sc.Kind != row.scaling.Kind || sc.AboveLevel != row.scaling.AboveLevel || sc.DiceCount != row.scaling.DiceCount || sc.DieSides != row.scaling.DieSides || !slices.Equal(sc.EffectIDs, row.scaling.EffectIDs) {
		return fmt.Errorf("batch 203 scaling differs from source")
	}
	if len(s.Effects) != len(row.effects) {
		return fmt.Errorf("batch 203 effects differ from source")
	}
	for i, want := range row.effects {
		e := s.Effects[i]
		if e.ID != want.id || e.Mechanic != want.mechanic || e.DiceCount != want.dice || e.DieSides != want.sides || e.DamageType != want.damage || e.Bonus != want.bonus || e.AddCastingModifier == nil || *e.AddCastingModifier || e.HalfInitialOnly == nil || *e.HalfInitialOnly {
			return fmt.Errorf("batch 203 effect %s differs from source", want.id)
		}
	}
	if !slices.Equal(s.RequiredMechanics, row.required) {
		return fmt.Errorf("batch 203 required mechanics differ from source")
	}
	slots := []int{}
	for level := s.Level; level <= 9 && s.Level > 0; level++ {
		slots = append(slots, level)
	}
	if !slices.Equal(s.SlotLevels, slots) {
		return fmt.Errorf("batch 203 slot levels differ from source")
	}
	return nil
}

func (d *SpellSourceDetails) validateSpell203Payload(key string) error {
	bad := func() error { return fmt.Errorf("invalid batch 203 %s payload", key) }
	switch key {
	case "dimension_door":
		x := d.DimensionDoor
		if !slices.Equal(x.SourcePages, []int{123, 124}) || !slices.Equal(x.DestinationChoices, []string{"visible-place", "visualized-place", "described-by-distance-and-direction"}) ||
			!slices.Equal(x.DirectionExamples, []string{"200-feet-straight-downward", "300-feet-upward-to-the-northwest-at-45-degree-angle"}) || x.Arrival != "exactly-the-spot-desired" ||
			x.CompanionCount != 1 || !x.CompanionWilling || x.CompanionWithinFeet != 5 || x.CompanionArriveFeet != 5 || x.CompanionArrival != "space-within-5-feet-of-destination-space" ||
			x.CollisionTrigger != "arrival-space-occupied-by-creature-or-completely-filled-by-objects" || x.CollisionDiceCount != 4 || x.CollisionDieSides != 6 || x.CollisionDamageType != "force" ||
			x.CollisionDamaged != "caster-and-each-traveling-creature" || x.CollisionResult != "teleportation-fails" {
			return bad()
		}
	case "disguise_self":
		x := d.DisguiseSelf
		if !slices.Equal(x.Affects, []string{"caster", "clothing", "armor", "weapons", "other-belongings-on-person"}) || x.HeightChangeFeet != 1 || !slices.Equal(x.HeightChanges, []string{"shorter", "taller"}) ||
			!slices.Equal(x.WeightChanges, []string{"heavier", "lighter"}) || x.FormRequirement != "same-basic-arrangement-of-limbs" || x.OtherExtent != "up-to-caster" ||
			x.PhysicalInspection != "changes-fail-to-hold-up" || x.InspectionExample != "added-hat-objects-pass-through-and-touch-feels-nothing" || x.DiscoveryAction != "study-action" ||
			x.DiscoveryCheck != "intelligence-investigation" || x.DiscoveryDC != "caster-spell-save-dc" || x.DiscoveryResult != "creature-discerns-disguise" {
			return bad()
		}
	case "disintegrate":
		x := d.Disintegrate
		if !slices.Equal(x.TargetChoices, []string{"creature", "nonmagical-object", "creation-of-magical-force"}) || x.ForceCreationSample != "wall-of-force" || x.Save != "dexterity" || x.FailDiceCount != 10 || x.FailDieSides != 6 || x.FailBonus != 40 ||
			x.DamageType != "force" || x.SuccessDamage != "none-stated" || x.ZeroHitPointsResult != "disintegrated-into-gray-dust" || x.DustIncludes != "target-and-everything-nonmagical-worn-and-carried" ||
			!slices.Equal(x.RevivedOnlyBy, []string{"true-resurrection", "wish"}) || x.AutomaticSize != "large-or-smaller" || !slices.Equal(x.AutomaticTargets, []string{"nonmagical-object", "creation-of-magical-force"}) ||
			x.OversizeSize != "huge-or-larger" || x.OversizeResult != "disintegrates-10-foot-cube-portion" || x.SlotDiceCount != 3 || x.SlotDieSides != 6 {
			return bad()
		}
	case "dispel_evil":
		x := d.DispelEvil
		if !slices.Equal(x.AffectedTypes, []string{"celestial", "elemental", "fey", "fiend", "undead"}) || x.AttackEffect != "disadvantage-on-attack-rolls-against-caster" || x.EarlyEnd != "either-special-function" ||
			x.BreakAction != "magic-action" || x.BreakTarget != "touched-creature-possessed-by-or-charmed-or-frightened-by-listed-types" || !slices.Equal(x.BreakResults, []string{"no-longer-possessed", "no-longer-charmed", "no-longer-frightened"}) ||
			x.DismissalAction != "magic-action" || x.DismissalRange != 5 || x.DismissalTarget != "visible-creature-of-listed-type" || x.DismissalSave != "charisma" || x.DismissalFailure != "sent-back-to-home-plane-if-not-there" ||
			x.UndeadElsewhere != "shadowfell" || x.FeyElsewhere != "feywild" {
			return bad()
		}
	case "dispel_magic":
		x := d.DispelMagic
		if !slices.Equal(x.TargetChoices, []string{"creature", "object", "magical-effect"}) || x.AutoEndUpToLvl != 3 || x.CheckFromLevel != 4 || x.CheckAbility != "spellcasting-ability" || x.CheckDCBase != 10 ||
			x.CheckDCAdds != "spell-level" || x.CheckSuccess != "spell-ends" || x.SlotAutoEnd != "spell-level-equal-to-or-less-than-slot-level" {
			return bad()
		}
	case "dissonant_whispers":
		x := d.DissonantWhispers
		if x.Target != "one-creature-of-choice-that-caster-can-see-in-range" || x.Perception != "hears-discordant-melody-in-its-mind" || x.Save != "wisdom" || x.FailDiceCount != 3 || x.FailDieSides != 6 || x.DamageType != "psychic" ||
			x.FailMovementTrigger != "must-immediately-use-reaction-if-available" || x.FailMovement != "move-as-far-away-from-caster-as-it-can" || x.FailRoute != "safest-route" || x.SuccessDamage != "half-damage-only" ||
			x.SlotDiceCount != 1 || x.SlotDieSides != 6 {
			return bad()
		}
	case "divination":
		x := d.Divination
		if x.Contact != "god-or-gods-servants" || x.QuestionCount != 1 || x.QuestionSubject != "specific-goal-event-or-activity" || x.QuestionWindowDays != 7 || x.ReplyBy != "gm" || x.ReplyTruth != "truthful" ||
			!slices.Equal(x.ReplyForms, []string{"short-phrase", "cryptic-rhyme"}) || x.AccountsForChange || x.ChangedCircumstanceModel != "casting-of-other-spells" || x.RepeatReset != "long-rest" ||
			x.NoAnswerChancePercent != 25 || !x.NoAnswerCumulative || x.NoAnswerAfter != "first-casting" {
			return bad()
		}
	case "divine_favor":
		x := d.DivineFavor
		if x.AttackKind != "attacks-with-weapons" || x.Trigger != "hit" || x.ExtraDiceCnt != 1 || x.ExtraDieSides != 4 || x.DamageType != "radiant" || x.Applies != "until-spell-ends" {
			return bad()
		}
	case "divine_smite":
		x := d.DivineSmite
		if !slices.Equal(x.TriggerHits, []string{"melee-weapon", "unarmed-strike"}) || x.TriggerTiming != "immediately-after-hitting-target" || x.ExtraDiceCount != 2 || x.ExtraDieSides != 8 || x.DamageType != "radiant" ||
			x.TypeBonusDice != 1 || x.TypeBonusDieSides != 8 || !slices.Equal(x.TypeBonusTypes, []string{"fiend", "undead"}) || x.SlotDiceCount != 1 || x.SlotDieSides != 8 {
			return bad()
		}
	case "divine_word":
		x := d.DivineWord
		want := []SpellDivineWordRow{{0, 20, true, false, false, false, 0}, {21, 30, false, true, true, true, 3600}, {31, 40, false, true, true, false, 600}, {41, 50, false, false, true, false, 60}}
		if !slices.Equal(x.SourcePages, []int{125}) || x.Affects != "each-creature-of-choice-in-range" || x.Save != "charisma" || x.HitPointLimit != 50 || x.LimitBasis != "current-hit-points" || !slices.Equal(x.Rows, want) ||
			!slices.Equal(x.BanishTypes, []string{"celestial", "elemental", "fey", "fiend"}) || !x.BanishRegardless || x.BanishResult != "forced-back-to-plane-of-origin-if-not-there" || x.BanishReturnHrs != 24 || x.BanishException != "wish-spell" {
			return bad()
		}
	case "dominate_beast":
		if !validSpell203Dominate(d.DominateBeast, []int{125}, "beast", []SpellDominateSlotDuration{{5, 600, false}, {6, 3600, false}, {7, 28800, true}}) {
			return bad()
		}
	case "dominate_monster":
		if !validSpell203Dominate(d.DominateMonster, []int{125, 126}, "creature", []SpellDominateSlotDuration{{9, 28800, false}}) {
			return bad()
		}
	}
	return nil
}

func validSpell203Dominate(x *SpellDominate, pages []int, kind string, table []SpellDominateSlotDuration) bool {
	return slices.Equal(x.SourcePages, pages) && x.TargetKind == kind && x.Save == "wisdom" && x.FailCondition == "charmed" && x.ConditionLength == "spell-duration" &&
		x.SaveAdvantage == "caster-or-allies-fighting-target" && x.RepeatSave == "whenever-target-takes-damage" && x.RepeatSuccess == "spell-ends-on-target" && x.Link == "telepathic-while-same-plane" &&
		x.CommandTiming == "caster-turn" && x.CommandAction == "no-action-required" && slices.Equal(x.CommandExamples, []string{"Attack that creature", "Move over there", "Fetch that object"}) &&
		x.Obey == "target-does-its-best-to-obey-on-its-turn" && x.Idle == "acts-and-moves-as-it-likes-focusing-on-protecting-itself" && x.ReactionCommand == "caster-can-command-target-to-take-reaction" &&
		x.ReactionCost == "caster-must-take-own-reaction" && slices.Equal(x.SlotDurations, table)
}
