package srd521

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

// SpellGlobe retains Globe of Invulnerability's typed SRD source facts without implementing them.
type SpellGlobe struct {
	EmanationFeet             int    `json:"emanation_feet"`
	BarrierMobility           string `json:"barrier_mobility"`
	BarrierAppearance         string `json:"barrier_appearance"`
	BlockedSpellMaxLevel      int    `json:"blocked_spell_max_level"`
	BlockedSpellOrigin        string `json:"blocked_spell_origin"`
	SpellTargeting            string `json:"spell_targeting"`
	AreaExclusion             string `json:"area_exclusion"`
	SlotBlockedLevelIncrement int    `json:"slot_blocked_level_increment"`
}

// SpellGlyph retains Glyph of Warding's inscription, trigger, rune and stored-spell rules.
type SpellGlyph struct {
	SourcePages         []int    `json:"source_pages"`
	InscribeLocations   []string `json:"inscribe_locations"`
	MaxDiameterFeet     int      `json:"max_diameter_feet"`
	BreakDistanceFeet   int      `json:"break_distance_feet"`
	BreakResult         string   `json:"break_result"`
	DetectionCheck      string   `json:"detection_check"`
	DetectionDC         string   `json:"detection_dc"`
	TriggerSet          string   `json:"trigger_set"`
	TriggerTypes        []string `json:"trigger_types"`
	SurfaceTriggers     []string `json:"surface_triggers"`
	ObjectTriggers      []string `json:"object_triggers"`
	TriggerEndsSpell    bool     `json:"trigger_ends_spell"`
	TriggerRefinements  []string `json:"trigger_refinements"`
	RefinementExamples  []string `json:"refinement_examples"`
	RuneShape           string   `json:"rune_shape"`
	RuneRadiusFeet      int      `json:"rune_radius_feet"`
	RuneSave            string   `json:"rune_save"`
	RuneDiceCount       int      `json:"rune_dice_count"`
	RuneDieSides        int      `json:"rune_die_sides"`
	RuneDamageTypes     []string `json:"rune_damage_types"`
	RuneDamageChoice    string   `json:"rune_damage_choice"`
	RuneSuccess         string   `json:"rune_success"`
	SpellMaxLevel       int      `json:"spell_max_level"`
	SpellPrepared       bool     `json:"spell_prepared"`
	SpellTargets        []string `json:"spell_targets"`
	SpellStoredEffect   string   `json:"spell_stored_effect"`
	SpellTarget         string   `json:"spell_target"`
	SpellAreaCenter     string   `json:"spell_area_center"`
	SpellHostileSummons string   `json:"spell_hostile_summons"`
	SpellConcentration  string   `json:"spell_concentration"`
	SlotRuneDiceCount   int      `json:"slot_rune_dice_count"`
	SlotRuneDieSides    int      `json:"slot_rune_die_sides"`
	SlotSpellLevel      string   `json:"slot_spell_level"`
}

// SpellGoodberry retains Goodberry's berries and nourishment without implementing them.
type SpellGoodberry struct {
	BerryCount      int    `json:"berry_count"`
	BerryLocation   string `json:"berry_location"`
	EatAction       string `json:"eat_action"`
	BerryHitPoints  int    `json:"berry_hit_points"`
	NourishmentDays int    `json:"nourishment_days"`
	UneatenBerries  string `json:"uneaten_berries"`
}

// SpellGrease retains Grease's area, terrain and Dexterity save without implementing them.
type SpellGrease struct {
	SquareFeet        int      `json:"square_feet"`
	Center            string   `json:"center"`
	Flammability      string   `json:"flammability"`
	Terrain           string   `json:"terrain"`
	Save              string   `json:"save"`
	SaveTriggers      []string `json:"save_triggers"`
	AppearFailure     string   `json:"appear_failure"`
	EnterOrEndFailure string   `json:"enter_or_end_failure"`
}

// SpellGreaterInvisibility retains the touched creature's Invisible condition.
type SpellGreaterInvisibility struct {
	Target       string `json:"target"`
	Condition    string `json:"condition"`
	ConditionEnd string `json:"condition_ends"`
}

// SpellGreaterRestoration retains the five removable effects without implementing them.
type SpellGreaterRestoration struct {
	Target                string   `json:"target"`
	ChoiceCount           int      `json:"choice_count"`
	ExhaustionLevels      int      `json:"exhaustion_levels_removed"`
	ConditionsRemoved     []string `json:"conditions_removed"`
	CurseIncludesAttune   bool     `json:"curse_includes_attunement"`
	AbilityScoreReduction string   `json:"ability_score_reduction"`
	HitPointMaxReduction  string   `json:"hit_point_maximum_reduction"`
}

// SpellGuardian retains Guardian of Faith's guardian, trigger, damage and ending.
type SpellGuardian struct {
	GuardianSize        string   `json:"guardian_size"`
	GuardianNature      string   `json:"guardian_nature"`
	GuardianPosture     string   `json:"guardian_posture"`
	GuardianSpace       string   `json:"guardian_space"`
	OccupiesSpace       bool     `json:"occupies_space"`
	Invulnerable        bool     `json:"invulnerable"`
	Appearance          string   `json:"appearance"`
	TriggerRadiusFeet   int      `json:"trigger_radius_feet"`
	TriggerTargets      string   `json:"trigger_targets"`
	Triggers            []string `json:"triggers"`
	Save                string   `json:"save"`
	DamageFlat          int      `json:"damage_flat"`
	DamageType          string   `json:"damage_type"`
	SaveSuccess         string   `json:"save_success"`
	VanishesTotalDamage int      `json:"vanishes_total_damage"`
}

// SpellGuards retains Guards and Wards' area, password, four effects and extension.
type SpellGuards struct {
	SourcePages               []int    `json:"source_pages"`
	FloorAreaSquareFeet       int      `json:"floor_area_square_feet"`
	HeightFeet                int      `json:"height_feet"`
	Shapes                    []string `json:"shapes"`
	Unaffected                string   `json:"unaffected"`
	PasswordDistanceFeet      int      `json:"password_distance_feet"`
	PasswordResult            string   `json:"password_result"`
	DispelMagicOnSpell        string   `json:"dispel_magic_on_spell"`
	DispellableEffects        int      `json:"dispellable_effects"`
	AllDispelled              string   `json:"all_dispelled"`
	ExtensionDays             int      `json:"extension_days"`
	ExtensionResult           string   `json:"extension_result"`
	CorridorsObscurement      string   `json:"corridors_obscurement"`
	CorridorsMisdirectPercent int      `json:"corridors_misdirection_percent"`
	CorridorsMisdirection     string   `json:"corridors_misdirection"`
	DoorsLock                 string   `json:"doors_lock"`
	DoorsIllusionCount        int      `json:"doors_illusion_count"`
	DoorsIllusion             string   `json:"doors_illusion"`
	StairsWebs                string   `json:"stairs_webs"`
	StairsRegrowMinutes       int      `json:"stairs_regrow_minutes"`
	OtherEffectCount          int      `json:"other_effect_count"`
	OtherEffects              []string `json:"other_effects"`
}

// SpellGuidance retains Guidance's d4 skill bonus without implementing it.
type SpellGuidance struct {
	Target         string `json:"target"`
	Choice         string `json:"choice"`
	BonusDiceCount int    `json:"bonus_dice_count"`
	BonusDieSides  int    `json:"bonus_die_sides"`
	AppliesTo      string `json:"applies_to"`
}

// SpellGuidingBolt retains Guiding Bolt's attack, damage and Advantage window.
type SpellGuidingBolt struct {
	Attack          string `json:"attack"`
	Target          string `json:"target"`
	DamageDiceCount int    `json:"damage_dice_count"`
	DamageDieSides  int    `json:"damage_die_sides"`
	DamageType      string `json:"damage_type"`
	AdvantageOn     string `json:"advantage_on"`
	AdvantageWindow string `json:"advantage_window"`
	SlotDiceCount   int    `json:"slot_dice_count"`
	SlotDieSides    int    `json:"slot_die_sides"`
}

// SpellGust retains Gust of Wind's line, push, movement and flame rules.
type SpellGust struct {
	LineLengthFeet           int    `json:"line_length_feet"`
	LineWidthFeet            int    `json:"line_width_feet"`
	Direction                string `json:"direction"`
	Save                     string `json:"save"`
	PushFeet                 int    `json:"push_feet"`
	PushDirection            string `json:"push_direction"`
	EndTurnSave              string `json:"end_turn_save"`
	ApproachMovementCostFeet int    `json:"approach_movement_cost_feet"`
	ApproachMovementPerFeet  int    `json:"approach_movement_per_feet"`
	Disperses                string `json:"disperses"`
	Extinguishes             string `json:"extinguishes"`
	ProtectedFlamesExample   string `json:"protected_flames_example"`
	ProtectedFlamesEffect    string `json:"protected_flames_effect"`
	ProtectedFlamesPercent   int    `json:"protected_flames_extinguish_percent"`
	RedirectAction           string `json:"redirect_action"`
}

// SpellHallow retains Hallow's ward and all ten extra-effect choices.
type SpellHallow struct {
	SourcePages                  []int    `json:"source_pages"`
	TouchTarget                  string   `json:"touch_target"`
	Power                        string   `json:"power"`
	MaxRadiusFeet                int      `json:"max_radius_feet"`
	OverlapResult                string   `json:"overlap_result"`
	WardTypes                    []string `json:"ward_types"`
	WardChoice                   string   `json:"ward_choice"`
	WardEntry                    string   `json:"ward_entry"`
	WardPossessionAndConditions  []string `json:"ward_possession_and_conditions"`
	WardPossessionResult         string   `json:"ward_possession_result"`
	ExtraEffectCount             int      `json:"extra_effect_count"`
	ExtraEffects                 []string `json:"extra_effects"`
	Courage                      string   `json:"courage"`
	Darkness                     string   `json:"darkness"`
	Daylight                     string   `json:"daylight"`
	PeacefulRest                 string   `json:"peaceful_rest"`
	ExtradimensionalInterference string   `json:"extradimensional_interference"`
	Fear                         string   `json:"fear"`
	Resistance                   string   `json:"resistance"`
	Silence                      string   `json:"silence"`
	Tongues                      string   `json:"tongues"`
	Vulnerability                string   `json:"vulnerability"`
}

func spell208Templates() map[string]any {
	return map[string]any{
		"globe":                SpellGlobe{},
		"glyph":                SpellGlyph{},
		"goodberry":            SpellGoodberry{},
		"grease":               SpellGrease{},
		"greater_invisibility": SpellGreaterInvisibility{},
		"greater_restoration":  SpellGreaterRestoration{},
		"guardian":             SpellGuardian{},
		"guards":               SpellGuards{},
		"guidance":             SpellGuidance{},
		"guiding_bolt":         SpellGuidingBolt{},
		"gust":                 SpellGust{},
		"hallow":               SpellHallow{},
	}
}

// validSpell208Target lists the target eligibilities this batch introduces.
func validSpell208Target(t SpellTarget) bool {
	switch t.Shape {
	case "emanation":
		return t.Count == 0 && t.Eligibility == "area-within-barrier"
	case "choice":
		return t.Count == 1 && t.Eligibility == "surface-or-closable-object"
	case "square":
		return t.Count == 0 && t.Eligibility == "creatures-in-area-centered-in-range"
	case "created-effect":
		return (t.Count == 1 && (t.Eligibility == "unoccupied-visible-space-in-range" || t.Eligibility == "warded-floor-area")) ||
			(t.Count == 0 && t.Eligibility == "area-around-touched-point")
	case "line":
		return t.Count == 0 && t.Eligibility == "creatures-in-line"
	}
	return false
}

// validSpell208Scaling checks the one scaling kind this batch introduces: a
// spell-level limit that rises by a fixed amount per slot level above the base.
func (s SpellRecord) validSpell208Scaling(effects map[string]bool) bool {
	sc, d := s.Scaling, s.SourceDetails
	if d == nil || d.Globe == nil || sc.Kind != "slot-limit" || s.Level == 0 || sc.AboveLevel != s.Level || sc.DiceCount != 0 || sc.DieSides != 0 ||
		len(sc.EffectIDs) != 1 || !effects[sc.EffectIDs[0]] || d.Globe.SlotBlockedLevelIncrement <= 0 {
		return false
	}
	return slices.Contains(s.RequiredMechanics, "scaling.barrier-level-per-slot")
}

type spell208Effect struct {
	id, mechanic string
	dice, sides  int
	damage       string
}

// spell208Row pins every header value, scaling row, effect and dependency a
// batch record must carry, by equality with the printed SRD page.
type spell208Row struct {
	slug, name      string
	page, level     int
	school          string
	lists           []string
	casting         string
	castingSeconds  int
	concentration   bool
	rangeKind       string
	rangeFeet       int
	target          SpellTarget
	duration        string
	seconds         int
	endTiming       string
	verbal, somatic bool
	material        string
	costCP          int
	consumed        bool
	scaling         SpellScaling
	effects         []spell208Effect
	required        []string
}

var spell208Rows = map[string]spell208Row{
	"globe": {"globe-of-invulnerability", "Globe of Invulnerability", 136, 6, "abjuration", []string{"sorcerer", "wizard"}, "action", 0, true, "self", 0, SpellTarget{"emanation", 0, "area-within-barrier"}, "timed", 60, "duration-or-concentration-ended", true, true, "a glass bead", 0, false,
		SpellScaling{"slot-limit", 6, 0, 0, []string{"barrier"}}, []spell208Effect{{"barrier", "barrier.globe-of-invulnerability", 0, 0, ""}},
		[]string{"barrier.globe-of-invulnerability", "scaling.barrier-level-per-slot", "concentration.source-duration", "components.material"}},
	"glyph": {"glyph-of-warding", "Glyph of Warding", 136, 3, "abjuration", []string{"bard", "cleric", "wizard"}, "1-hour", 3600, false, "touch", 0, SpellTarget{"choice", 1, "surface-or-closable-object"}, "until-dispelled", 0, "dispelled-or-triggered", true, true, "powdered diamond worth 200+ GP, which the spell consumes", 20000, true,
		SpellScaling{"slot", 3, 1, 8, []string{"explosion"}}, []spell208Effect{{"explosion", "damage.glyph-explosive-rune", 5, 8, ""}, {"stored-spell", "storage.glyph-spell", 0, 0, ""}},
		[]string{"damage.glyph-explosive-rune", "storage.glyph-spell", "scaling.per-slot", "components.material"}},
	"goodberry": {"goodberry", "Goodberry", 137, 1, "conjuration", []string{"druid", "ranger"}, "action", 0, false, "self", 0, SpellTarget{"single", 1, "self"}, "timed", 86400, "duration-ended", true, true, "a sprig of mistletoe", 0, false,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell208Effect{{"berries", "create.goodberry-berries", 0, 0, ""}, {"eat", "healing.goodberry-eat", 0, 0, ""}},
		[]string{"create.goodberry-berries", "healing.goodberry-eat", "components.material"}},
	"grease": {"grease", "Grease", 137, 1, "conjuration", []string{"sorcerer", "wizard"}, "action", 0, false, "feet", 60, SpellTarget{"square", 0, "creatures-in-area-centered-in-range"}, "timed", 60, "duration-ended", true, true, "a bit of pork rind or butter", 0, false,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell208Effect{{"terrain", "terrain.grease-difficult", 0, 0, ""}, {"save", "save.grease-dexterity-prone", 0, 0, ""}},
		[]string{"terrain.grease-difficult", "save.grease-dexterity-prone", "components.material"}},
	"greater_invisibility": {"greater-invisibility", "Greater Invisibility", 137, 4, "illusion", []string{"bard", "sorcerer", "wizard"}, "action", 0, true, "touch", 0, SpellTarget{"single", 1, "creature-touch"}, "timed", 60, "duration-or-concentration-ended", true, true, "", 0, false,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell208Effect{{"invisible", "condition.greater-invisibility-invisible", 0, 0, ""}},
		[]string{"condition.greater-invisibility-invisible", "concentration.source-duration"}},
	"greater_restoration": {"greater-restoration", "Greater Restoration", 137, 5, "abjuration", []string{"bard", "cleric", "druid", "paladin", "ranger"}, "action", 0, false, "touch", 0, SpellTarget{"single", 1, "creature-touch"}, "instantaneous", 0, "instantaneous", true, true, "diamond dust worth 100+ GP, which the spell consumes", 10000, true,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell208Effect{{"remove", "remove.greater-restoration", 0, 0, ""}},
		[]string{"remove.greater-restoration", "components.material"}},
	"guardian": {"guardian-of-faith", "Guardian of Faith", 137, 4, "conjuration", []string{"cleric"}, "action", 0, false, "feet", 30, SpellTarget{"created-effect", 1, "unoccupied-visible-space-in-range"}, "timed", 28800, "duration-ended-or-total-damage-dealt", true, false, "", 0, false,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell208Effect{{"guardian", "summon.guardian-of-faith-guardian", 0, 0, ""}, {"ward", "damage.guardian-of-faith-ward", 0, 0, "radiant"}},
		[]string{"summon.guardian-of-faith-guardian", "damage.guardian-of-faith-ward", "save.dexterity-half", "ending.guardian-total-damage"}},
	"guards": {"guards-and-wards", "Guards and Wards", 137, 6, "abjuration", []string{"bard", "wizard"}, "1-hour", 3600, false, "touch", 0, SpellTarget{"created-effect", 1, "warded-floor-area"}, "timed", 86400, "duration-ended-or-all-effects-dispelled", true, true, "a silver rod worth 10+ GP", 1000, false,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell208Effect{{"ward", "ward.guards-and-wards-area", 0, 0, ""}, {"corridors", "area.guards-and-wards-corridors", 0, 0, ""}, {"doors", "lock.guards-and-wards-doors", 0, 0, ""}, {"stairs", "web.guards-and-wards-stairs", 0, 0, ""}, {"other", "spell.guards-and-wards-other-effect", 0, 0, ""}},
		[]string{"ward.guards-and-wards-area", "area.guards-and-wards-corridors", "lock.guards-and-wards-doors", "web.guards-and-wards-stairs", "spell.guards-and-wards-other-effect", "dispel.guards-and-wards-effects", "extension.guards-and-wards-365-days", "components.material"}},
	"guidance": {"guidance", "Guidance", 138, 0, "divination", []string{"cleric", "druid"}, "action", 0, true, "touch", 0, SpellTarget{"single", 1, "willing-creature-touch"}, "timed", 60, "duration-or-concentration-ended", true, true, "", 0, false,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell208Effect{{"check", "bonus.guidance-ability-check", 1, 4, ""}},
		[]string{"bonus.guidance-ability-check", "concentration.source-duration"}},
	"guiding_bolt": {"guiding-bolt", "Guiding Bolt", 138, 1, "evocation", []string{"cleric"}, "action", 0, false, "feet", 120, SpellTarget{"single", 1, "creature-in-range"}, "timed", 6, "end-of-caster-next-turn", true, true, "", 0, false,
		SpellScaling{"slot", 1, 1, 6, []string{"damage"}}, []spell208Effect{{"attack", "attack.guiding-bolt-ranged-spell", 0, 0, ""}, {"damage", "damage.guiding-bolt", 4, 6, "radiant"}, {"advantage", "advantage.guiding-bolt-next-attack", 0, 0, ""}},
		[]string{"attack.guiding-bolt-ranged-spell", "damage.guiding-bolt", "advantage.guiding-bolt-next-attack", "scaling.per-slot"}},
	"gust": {"gust-of-wind", "Gust of Wind", 138, 2, "evocation", []string{"druid", "ranger", "sorcerer", "wizard"}, "action", 0, true, "self", 0, SpellTarget{"line", 0, "creatures-in-line"}, "timed", 60, "duration-or-concentration-ended", true, true, "a legume seed", 0, false,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell208Effect{{"gust", "push.gust-of-wind", 0, 0, ""}, {"movement", "movement.gust-closer-double-cost", 0, 0, ""}, {"environment", "environment.gust-of-wind-gas-and-flames", 0, 0, ""}, {"redirect", "direction.gust-bonus-action", 0, 0, ""}},
		[]string{"push.gust-of-wind", "movement.gust-closer-double-cost", "environment.gust-of-wind-gas-and-flames", "direction.gust-bonus-action", "concentration.source-duration", "components.material"}},
	"hallow": {"hallow", "Hallow", 138, 5, "abjuration", []string{"cleric"}, "24-hours", 86400, false, "touch", 0, SpellTarget{"created-effect", 0, "area-around-touched-point"}, "until-dispelled", 0, "dispelled", true, true, "incense worth 1,000+ GP, which the spell consumes", 100000, true,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell208Effect{{"area", "area.hallow", 0, 0, ""}, {"ward", "ward.hallow-hallowed-ward", 0, 0, ""}, {"extra", "effect.hallow-extra", 0, 0, ""}},
		[]string{"area.hallow", "ward.hallow-hallowed-ward", "effect.hallow-extra", "components.material"}},
}

func (s SpellRecord) validateSpell208() (bool, error) {
	expected := ""
	for key, row := range spell208Rows {
		if s.ID == prefix+"spells/"+row.slug {
			expected = key
		}
	}
	if s.SourceDetails == nil {
		if expected != "" {
			return true, fmt.Errorf("batch 208 details required")
		}
		return false, nil
	}
	d := s.SourceDetails
	payloads := map[string]any{
		"globe": d.Globe, "glyph": d.Glyph, "goodberry": d.Goodberry, "grease": d.Grease,
		"greater_invisibility": d.GreaterInvisibility, "greater_restoration": d.GreaterRestoration,
		"guardian": d.Guardian, "guards": d.Guards, "guidance": d.Guidance, "guiding_bolt": d.GuidingBolt,
		"gust": d.Gust, "hallow": d.Hallow,
	}
	count := 0
	for key, value := range payloads {
		if reflect.ValueOf(value).IsNil() {
			continue
		}
		count++
		if expected != key {
			return true, fmt.Errorf("mismatched batch 208 payload")
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
	row := spell208Rows[expected]
	if count != 1 || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" || d.CastingSeconds != row.castingSeconds || d.CastingTriggers == nil || len(d.CastingTriggers) != 0 {
		return true, fmt.Errorf("invalid batch 208 source details")
	}
	if err := s.validateSpell208Header(row); err != nil {
		return true, err
	}
	if err := d.validateSpell208Payload(expected); err != nil {
		return true, err
	}
	return true, nil
}

func (s SpellRecord) validateSpell208Header(row spell208Row) error {
	flags := func(p *bool) bool { return p != nil && *p }
	c := s.Components
	if s.Name != row.name || s.SourcePage != row.page || s.Level != row.level || s.School != row.school || !slices.Equal(s.Lists, row.lists) ||
		s.CastingTime != row.casting || flags(s.Ritual) || flags(s.Concentration) != row.concentration ||
		s.Range != (SpellRange{row.rangeKind, row.rangeFeet}) || s.Target != row.target || s.Duration != row.duration || s.DurationSeconds != row.seconds ||
		s.SourceDetails.EndTiming != row.endTiming || flags(c.Verbal) != row.verbal || flags(c.Somatic) != row.somatic {
		return fmt.Errorf("batch 208 header differs from source")
	}
	if m := c.Material; row.material == "" {
		if m != nil {
			return fmt.Errorf("batch 208 material differs from source")
		}
	} else if m == nil || m.Description != row.material || m.CostCP == nil || *m.CostCP != row.costCP || m.Consumed == nil || *m.Consumed != row.consumed {
		return fmt.Errorf("batch 208 material differs from source")
	}
	sc := s.Scaling
	if sc.Kind != row.scaling.Kind || sc.AboveLevel != row.scaling.AboveLevel || sc.DiceCount != row.scaling.DiceCount || sc.DieSides != row.scaling.DieSides || !slices.Equal(sc.EffectIDs, row.scaling.EffectIDs) {
		return fmt.Errorf("batch 208 scaling differs from source")
	}
	if len(s.Effects) != len(row.effects) {
		return fmt.Errorf("batch 208 effects differ from source")
	}
	for i, want := range row.effects {
		e := s.Effects[i]
		if e.ID != want.id || e.Mechanic != want.mechanic || e.DiceCount != want.dice || e.DieSides != want.sides || e.DamageType != want.damage || e.Bonus != 0 || e.AddCastingModifier == nil || *e.AddCastingModifier || e.HalfInitialOnly == nil || *e.HalfInitialOnly {
			return fmt.Errorf("batch 208 effect %s differs from source", want.id)
		}
	}
	if !slices.Equal(s.RequiredMechanics, row.required) {
		return fmt.Errorf("batch 208 required mechanics differ from source")
	}
	slots := []int{}
	for level := s.Level; level <= 9 && s.Level > 0; level++ {
		slots = append(slots, level)
	}
	if !slices.Equal(s.SlotLevels, slots) {
		return fmt.Errorf("batch 208 slot levels differ from source")
	}
	return nil
}

func (d *SpellSourceDetails) validateSpell208Payload(key string) error {
	bad := func() error { return fmt.Errorf("invalid batch 208 %s payload", key) }
	switch key {
	case "globe":
		x := d.Globe
		if x.EmanationFeet != 10 || x.BarrierMobility != "immobile" || x.BarrierAppearance != "shimmering" || x.BlockedSpellMaxLevel != 5 || x.BlockedSpellOrigin != "cast-from-outside-barrier" || x.SpellTargeting != "may-target-creatures-and-objects-within-but-no-effect" || x.AreaExclusion != "area-within-barrier-excluded-from-areas-of-effect" || x.SlotBlockedLevelIncrement != 1 {
			return bad()
		}
	case "glyph":
		x := d.Glyph
		if !slices.Equal(x.SourcePages, []int{136, 137}) || !slices.Equal(x.InscribeLocations, []string{"surface", "within-closable-object"}) || x.MaxDiameterFeet != 10 || x.BreakDistanceFeet != 10 || x.BreakResult != "glyph-broken-and-spell-ends-without-triggering" ||
			x.DetectionCheck != "wisdom-perception" || x.DetectionDC != "caster-spell-save-dc" || x.TriggerSet != "when-glyph-inscribed" || !slices.Equal(x.TriggerTypes, []string{"explosive-rune", "spell-glyph"}) ||
			!slices.Equal(x.SurfaceTriggers, []string{"touching-or-stepping-on-glyph", "removing-object-covering-glyph", "approaching-within-distance"}) || !slices.Equal(x.ObjectTriggers, []string{"opening-object", "seeing-glyph"}) || !x.TriggerEndsSpell ||
			!slices.Equal(x.TriggerRefinements, []string{"only-chosen-creature-types", "exempt-creatures-meeting-conditions"}) || !slices.Equal(x.RefinementExamples, []string{"aberrations", "password"}) ||
			x.RuneShape != "sphere" || x.RuneRadiusFeet != 20 || x.RuneSave != "dexterity" || x.RuneDiceCount != 5 || x.RuneDieSides != 8 || !slices.Equal(x.RuneDamageTypes, []string{"acid", "cold", "fire", "lightning", "thunder"}) || x.RuneDamageChoice != "when-glyph-created" || x.RuneSuccess != "half-damage" ||
			x.SpellMaxLevel != 3 || !x.SpellPrepared || !slices.Equal(x.SpellTargets, []string{"single-creature", "area"}) || x.SpellStoredEffect != "no-immediate-effect" || x.SpellTarget != "creature-that-triggered-glyph" || x.SpellAreaCenter != "creature-that-triggered-glyph" ||
			x.SpellHostileSummons != "appear-closest-to-intruder-and-attack-it" || x.SpellConcentration != "lasts-until-end-of-full-duration" || x.SlotRuneDiceCount != 1 || x.SlotRuneDieSides != 8 || x.SlotSpellLevel != "up-to-slot-level" {
			return bad()
		}
	case "goodberry":
		x := d.Goodberry
		if x.BerryCount != 10 || x.BerryLocation != "caster-hand" || x.EatAction != "bonus-action" || x.BerryHitPoints != 1 || x.NourishmentDays != 1 || x.UneatenBerries != "disappear-when-spell-ends" {
			return bad()
		}
	case "grease":
		x := d.Grease
		if x.SquareFeet != 10 || x.Center != "point-within-range" || x.Flammability != "nonflammable" || x.Terrain != "difficult-terrain" || x.Save != "dexterity" ||
			!slices.Equal(x.SaveTriggers, []string{"grease-appears-for-creature-standing-in-area", "creature-enters-area", "creature-ends-turn-in-area"}) || x.AppearFailure != "prone-condition" || x.EnterOrEndFailure != "fall-prone" {
			return bad()
		}
	case "greater_invisibility":
		x := d.GreaterInvisibility
		if x.Target != "creature-touched" || x.Condition != "invisible" || x.ConditionEnd != "spell-ends" {
			return bad()
		}
	case "greater_restoration":
		x := d.GreaterRestoration
		if x.Target != "creature-touched" || x.ChoiceCount != 1 || x.ExhaustionLevels != 1 || !slices.Equal(x.ConditionsRemoved, []string{"charmed", "petrified"}) || !x.CurseIncludesAttune || x.AbilityScoreReduction != "any-reduction-to-one-ability-score" || x.HitPointMaxReduction != "any-reduction" {
			return bad()
		}
	case "guardian":
		x := d.Guardian
		if x.GuardianSize != "large" || x.GuardianNature != "spectral" || x.GuardianPosture != "hovers" || x.GuardianSpace != "unoccupied-space-visible-in-range" || !x.OccupiesSpace || !x.Invulnerable || x.Appearance != "form-appropriate-for-deity-or-pantheon" ||
			x.TriggerRadiusFeet != 10 || x.TriggerTargets != "enemy" || !slices.Equal(x.Triggers, []string{"moves-within-radius-first-time-on-turn", "starts-turn-within-radius"}) || x.Save != "dexterity" || x.DamageFlat != 20 || x.DamageType != "radiant" || x.SaveSuccess != "half-damage" || x.VanishesTotalDamage != 60 {
			return bad()
		}
	case "guards":
		x := d.Guards
		if !slices.Equal(x.SourcePages, []int{137, 138}) || x.FloorAreaSquareFeet != 2500 || x.HeightFeet != 20 ||
			!slices.Equal(x.Shapes, []string{"one-50-foot-square", "one-hundred-contiguous-5-foot-squares", "twenty-five-contiguous-10-foot-squares"}) ||
			x.Unaffected != "individuals-specified-when-cast" || x.PasswordDistanceFeet != 5 || x.PasswordResult != "speaker-immune-to-effects" || x.DispelMagicOnSpell != "no-effect" || x.DispellableEffects != 4 || x.AllDispelled != "spell-ends" ||
			x.ExtensionDays != 365 || x.ExtensionResult != "lasts-until-all-effects-dispelled" || x.CorridorsObscurement != "heavily-obscured" || x.CorridorsMisdirectPercent != 50 || x.CorridorsMisdirection != "creature-other-than-caster-believes-opposite-direction" ||
			x.DoorsLock != "arcane-lock" || x.DoorsIllusionCount != 10 || x.DoorsIllusion != "appear-as-plain-sections-of-wall" || x.StairsWebs != "web-spell-top-to-bottom" || x.StairsRegrowMinutes != 10 || x.OtherEffectCount != 1 ||
			!slices.Equal(x.OtherEffects, []string{"dancing-lights-four-corridors-repeating-program", "magic-mouth-two-locations", "stinking-cloud-two-locations-vapors-return-within-10-minutes", "gust-of-wind-one-corridor-or-room-continuous", "suggestion-one-5-foot-square-entering-creature-mentally"}) {
			return bad()
		}
	case "guidance":
		x := d.Guidance
		if x.Target != "willing-creature-touched" || x.Choice != "skill" || x.BonusDiceCount != 1 || x.BonusDieSides != 4 || x.AppliesTo != "ability-check-using-chosen-skill" {
			return bad()
		}
	case "guiding_bolt":
		x := d.GuidingBolt
		if x.Attack != "ranged-spell-attack" || x.Target != "creature-in-range" || x.DamageDiceCount != 4 || x.DamageDieSides != 6 || x.DamageType != "radiant" || x.AdvantageOn != "next-attack-roll-against-target" || x.AdvantageWindow != "before-end-of-caster-next-turn" || x.SlotDiceCount != 1 || x.SlotDieSides != 6 {
			return bad()
		}
	case "gust":
		x := d.Gust
		if x.LineLengthFeet != 60 || x.LineWidthFeet != 10 || x.Direction != "chosen-by-caster" || x.Save != "strength" || x.PushFeet != 15 || x.PushDirection != "away-from-caster-along-line" || x.EndTurnSave != "same-save-when-creature-ends-turn-in-line" ||
			x.ApproachMovementCostFeet != 2 || x.ApproachMovementPerFeet != 1 || x.Disperses != "gas-or-vapor" || x.Extinguishes != "candles-and-similar-unprotected-flames" || x.ProtectedFlamesExample != "lanterns" || x.ProtectedFlamesEffect != "dance-wildly" || x.ProtectedFlamesPercent != 50 || x.RedirectAction != "bonus-action-on-later-turns" {
			return bad()
		}
	case "hallow":
		x := d.Hallow
		if !slices.Equal(x.SourcePages, []int{138, 139}) || x.TouchTarget != "point" || x.Power != "holy-or-unholy" || x.MaxRadiusFeet != 60 || x.OverlapResult != "spell-fails-if-radius-includes-area-under-hallow" ||
			!slices.Equal(x.WardTypes, []string{"aberration", "celestial", "elemental", "fey", "fiend", "undead"}) || x.WardChoice != "any-of-listed-types" || x.WardEntry != "chosen-types-cannot-willingly-enter" ||
			!slices.Equal(x.WardPossessionAndConditions, []string{"possessed", "charmed", "frightened"}) || x.WardPossessionResult != "not-possessed-charmed-or-frightened-by-chosen-types-while-in-area" || x.ExtraEffectCount != 1 ||
			!slices.Equal(x.ExtraEffects, []string{"courage", "darkness", "daylight", "peaceful-rest", "extradimensional-interference", "fear", "resistance", "silence", "tongues", "vulnerability"}) ||
			x.Courage != "chosen-types-cannot-gain-frightened-condition-in-area" ||
			x.Darkness != "darkness-fills-area-normal-light-and-magical-light-from-lower-level-spells-cannot-illuminate" ||
			x.Daylight != "bright-light-fills-area-magical-darkness-from-lower-level-spells-cannot-extinguish" ||
			x.PeacefulRest != "interred-dead-bodies-cannot-be-turned-into-undead" ||
			x.ExtradimensionalInterference != "chosen-types-cannot-enter-or-exit-by-teleportation-or-interplanar-travel" ||
			x.Fear != "chosen-types-have-frightened-condition-in-area" ||
			x.Resistance != "chosen-types-have-resistance-to-one-chosen-damage-type-in-area" ||
			x.Silence != "no-sound-emanates-from-or-reaches-into-area" ||
			x.Tongues != "chosen-types-communicate-with-any-creature-in-area-without-shared-language" ||
			x.Vulnerability != "chosen-types-have-vulnerability-to-one-chosen-damage-type-in-area" {
			return bad()
		}
	}
	return nil
}
