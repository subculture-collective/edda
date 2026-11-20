package srd521

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

// SpellCure retains Cure Wounds' typed SRD source facts without implementing them.
type SpellCure struct {
	Target           string `json:"target"`
	HealingDiceCount int    `json:"healing_dice_count"`
	HealingDieSides  int    `json:"healing_die_sides"`
	HealingModifier  string `json:"healing_modifier"`
	SlotDiceCount    int    `json:"slot_dice_count"`
	SlotDieSides     int    `json:"slot_die_sides"`
}

// SpellLights retains Dancing Lights' typed SRD source facts without implementing them.
type SpellLights struct {
	SourcePages       []int    `json:"source_pages"`
	LightCount        int      `json:"light_count"`
	LightSize         string   `json:"light_size"`
	Forms             []string `json:"forms"`
	CombinedForm      string   `json:"combined_form"`
	CombinedSize      string   `json:"combined_size"`
	LightLevel        string   `json:"light_level"`
	LightRadiusFeet   int      `json:"light_radius_feet"`
	MoveAction        string   `json:"move_action"`
	MoveFeet          int      `json:"move_feet"`
	MoveDestination   string   `json:"move_destination"`
	MaxSeparationFeet int      `json:"max_separation_feet"`
	SeparationFrom    string   `json:"separation_from"`
	VanishesWhen      string   `json:"vanishes_when"`
}

// SpellDarkness retains Darkness' typed SRD source facts without implementing them.
type SpellDarkness struct {
	DarknessKind      string   `json:"darkness_kind"`
	SphereRadiusFeet  int      `json:"sphere_radius_feet"`
	SphereOrigin      string   `json:"sphere_origin"`
	Darkvision        string   `json:"darkvision"`
	NonmagicalLight   string   `json:"nonmagical_light"`
	ObjectAlternative string   `json:"object_alternative"`
	EmanationFeet     int      `json:"emanation_feet"`
	ObjectBlocker     string   `json:"object_blocker"`
	BlockerExamples   []string `json:"blocker_examples"`
	DispelsLights     []string `json:"dispels_lights"`
	DispelsUpToLevel  int      `json:"dispels_up_to_level"`
}

// SpellDarkvision retains Darkvision's typed SRD source facts without implementing them.
type SpellDarkvision struct {
	Target         string `json:"target"`
	Sense          string `json:"sense"`
	SenseRangeFeet int    `json:"sense_range_feet"`
}

// SpellDaylight retains Daylight's typed SRD source facts without implementing them.
type SpellDaylight struct {
	LightKind         string   `json:"light_kind"`
	SphereRadiusFeet  int      `json:"sphere_radius_feet"`
	SphereOrigin      string   `json:"sphere_origin"`
	SphereLight       string   `json:"sphere_light"`
	DimLightExtraFeet int      `json:"dim_light_extra_feet"`
	ObjectAlternative string   `json:"object_alternative"`
	EmanationFeet     int      `json:"emanation_feet"`
	ObjectBlocker     string   `json:"object_blocker"`
	BlockerExamples   []string `json:"blocker_examples"`
	DispelsDarkness   string   `json:"dispels_darkness"`
	DispelsUpToLevel  int      `json:"dispels_up_to_level"`
}

// SpellDeathWard retains Death Ward's typed SRD source facts without implementing them.
type SpellDeathWard struct {
	Target                string `json:"target"`
	DropTrigger           string `json:"drop_trigger"`
	ReplacementHitPoints  int    `json:"replacement_hit_points"`
	DropEndsSpell         bool   `json:"drop_ends_spell"`
	InstantDeathTrigger   string `json:"instant_death_trigger"`
	InstantDeathResult    string `json:"instant_death_result"`
	InstantDeathEndsSpell bool   `json:"instant_death_ends_spell"`
}

// SpellFireball retains Delayed Blast Fireball's typed SRD source facts without implementing them.
type SpellFireball struct {
	BeadOrigin        string `json:"bead_origin"`
	BeadPoint         string `json:"bead_point"`
	ExplosionTrigger  string `json:"explosion_trigger"`
	RadiusFeet        int    `json:"radius_feet"`
	Save              string `json:"save"`
	DamageType        string `json:"damage_type"`
	FailDamage        string `json:"fail_damage"`
	SuccessDamage     string `json:"success_damage"`
	BaseDiceCount     int    `json:"base_dice_count"`
	BaseDieSides      int    `json:"base_die_sides"`
	GrowthDiceCount   int    `json:"growth_dice_count"`
	GrowthDieSides    int    `json:"growth_die_sides"`
	GrowthTiming      string `json:"growth_timing"`
	TouchSave         string `json:"touch_save"`
	TouchFailure      string `json:"touch_failure"`
	TouchSuccess      string `json:"touch_success"`
	ThrowFeet         int    `json:"throw_feet"`
	ThrownEnd         string `json:"thrown_end"`
	ThrownResult      string `json:"thrown_result"`
	Ignites           string `json:"ignites"`
	SlotBaseDiceCount int    `json:"slot_base_dice_count"`
	SlotBaseDieSides  int    `json:"slot_base_die_sides"`
}

// SpellDemiplane retains Demiplane's typed SRD source facts without implementing them.
type SpellDemiplane struct {
	SourcePages         []int    `json:"source_pages"`
	Surface             string   `json:"surface"`
	DoorSize            string   `json:"door_size"`
	DoorAppearance      string   `json:"door_appearance"`
	DoorOperation       string   `json:"door_operation"`
	RoomKind            string   `json:"room_kind"`
	RoomDimensionFeet   int      `json:"room_dimension_feet"`
	RoomMaterialChoices []string `json:"room_material_choices"`
	SpellEndDoor        string   `json:"spell_end_door"`
	ObjectsInside       string   `json:"objects_inside"`
	CreaturesInside     string   `json:"creatures_inside"`
	ShuntedCondition    string   `json:"shunted_condition"`
	ShuntedLanding      string   `json:"shunted_landing"`
	ConnectionChoices   []string `json:"connection_choices"`
	OtherCreatureNeeds  string   `json:"other_creature_needs"`
}

// SpellDetectBlockersEntry retains one source-listed barrier and its thickness.
type SpellDetectBlockersEntry struct {
	Material  string `json:"material"`
	Thickness string `json:"thickness"`
}

// SpellDetectEvil retains Detect Evil and Good's typed SRD source facts without implementing them.
type SpellDetectEvil struct {
	SenseRadiusFeet int                        `json:"sense_radius_feet"`
	DetectedTypes   []string                   `json:"detected_types"`
	SensedInfo      string                     `json:"sensed_info"`
	HallowDetection string                     `json:"hallow_detection"`
	Blockers        []SpellDetectBlockersEntry `json:"blockers"`
}

// SpellDetectMagic retains Detect Magic's typed SRD source facts without implementing them.
type SpellDetectMagic struct {
	SenseRadiusFeet int                        `json:"sense_radius_feet"`
	Senses          string                     `json:"senses"`
	FollowUpAction  string                     `json:"follow_up_action"`
	AuraSubjects    string                     `json:"aura_subjects"`
	AuraAppearance  string                     `json:"aura_appearance"`
	SchoolLearned   string                     `json:"school_learned"`
	Blockers        []SpellDetectBlockersEntry `json:"blockers"`
}

// SpellDetectPoison retains Detect Poison and Disease's typed SRD source facts without implementing them.
type SpellDetectPoison struct {
	SenseRadiusFeet int                        `json:"sense_radius_feet"`
	Detected        []string                   `json:"detected"`
	SensedInfo      string                     `json:"sensed_info"`
	KindSensed      string                     `json:"kind_sensed"`
	Blockers        []SpellDetectBlockersEntry `json:"blockers"`
}

// SpellDetectThoughts retains Detect Thoughts' two options and probe without implementing them.
type SpellDetectThoughts struct {
	ActivationAction     string                     `json:"activation_action"`
	Options              []string                   `json:"options"`
	SenseRadiusFeet      int                        `json:"sense_radius_feet"`
	SenseSubjects        string                     `json:"sense_subjects"`
	SenseReadsContent    bool                       `json:"sense_reads_content"`
	SenseBlockers        []SpellDetectBlockersEntry `json:"sense_blockers"`
	ReadRangeFeet        int                        `json:"read_range_feet"`
	ReadTargets          []string                   `json:"read_targets"`
	ReadLearns           string                     `json:"read_learns"`
	ReadNoLanguageResult string                     `json:"read_no_language_result"`
	ProbeAction          string                     `json:"probe_action"`
	ProbeSave            string                     `json:"probe_save"`
	ProbeFailure         []string                   `json:"probe_failure"`
	ProbeFailureExamples []string                   `json:"probe_failure_examples"`
	ProbeSuccess         string                     `json:"probe_success"`
	TargetKnowsProbed    bool                       `json:"target_knows_probed"`
	TargetCheckAction    string                     `json:"target_check_action"`
	TargetCheck          string                     `json:"target_check"`
	TargetCheckDC        string                     `json:"target_check_dc"`
	TargetCheckSuccess   string                     `json:"target_check_success"`
	TargetCheckWindow    string                     `json:"target_check_window"`
}

func spell202Templates() map[string]any {
	return map[string]any{
		"cure":            SpellCure{},
		"lights":          SpellLights{},
		"darkness":        SpellDarkness{},
		"darkvision":      SpellDarkvision{},
		"daylight":        SpellDaylight{},
		"death_ward":      SpellDeathWard{},
		"fireball":        SpellFireball{},
		"demiplane":       SpellDemiplane{},
		"detect_evil":     SpellDetectEvil{},
		"detect_magic":    SpellDetectMagic{},
		"detect_poison":   SpellDetectPoison{},
		"detect_thoughts": SpellDetectThoughts{},
	}
}

func validSpell202Target(t SpellTarget) bool {
	switch t.Shape {
	case "created-effect":
		return (t.Count == 4 && t.Eligibility == "torch-size-lights-in-range") || (t.Count == 1 && t.Eligibility == "visible-flat-solid-surface-in-range")
	case "choice":
		return t.Count == 1 && t.Eligibility == "point-in-range-or-unattended-object"
	}
	return false
}

type spell202Effect struct {
	id, mechanic string
	dice, sides  int
	damage       string
	modifier     bool
}

// spell202Row pins every header value, scaling row, effect and dependency a
// batch record must carry, by equality with the printed SRD page.
type spell202Row struct {
	slug, name            string
	page, level           int
	school                string
	lists                 []string
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
	scaling               SpellScaling
	effects               []spell202Effect
	required              []string
}

var spell202Rows = map[string]spell202Row{
	"cure": {"cure-wounds", "Cure Wounds", 121, 1, "abjuration", []string{"bard", "cleric", "druid", "paladin", "ranger"}, false, false, "touch", 0, SpellTarget{"single", 1, "creature-touch"}, "instantaneous", 0, "instantaneous", true, true, "", 0,
		SpellScaling{"slot", 1, 2, 8, []string{"healing"}}, []spell202Effect{{"healing", "healing.cure-wounds", 2, 8, "", true}}, []string{"healing.cure-wounds", "scaling.per-slot"}},
	"lights": {"dancing-lights", "Dancing Lights", 121, 0, "illusion", []string{"bard", "sorcerer", "wizard"}, false, true, "feet", 120, SpellTarget{"created-effect", 4, "torch-size-lights-in-range"}, "timed", 60, "duration-or-concentration-ended", true, true, "a bit of phosphorus", 0,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell202Effect{{"lights", "light.dancing-lights", 0, 0, "", false}, {"move", "movement.dancing-lights-bonus-action", 0, 0, "", false}},
		[]string{"light.dancing-lights", "movement.dancing-lights-bonus-action", "light.dim-10-foot-radius", "concentration.source-duration", "components.material"}},
	"darkness": {"darkness", "Darkness", 122, 2, "evocation", []string{"sorcerer", "warlock", "wizard"}, false, true, "feet", 60, SpellTarget{"choice", 1, "point-in-range-or-unattended-object"}, "timed", 600, "duration-or-concentration-ended", true, false, "bat fur and a piece of coal", 0,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell202Effect{{"darkness", "area.magical-darkness-sphere-15", 0, 0, "", false}, {"object-darkness", "area.magical-darkness-emanation-15", 0, 0, "", false}, {"dispel", "dispel.overlapping-light-spell-level-2-or-lower", 0, 0, "", false}},
		[]string{"area.magical-darkness-sphere-15", "area.magical-darkness-emanation-15", "dispel.overlapping-light-spell-level-2-or-lower", "concentration.source-duration", "components.material"}},
	"darkvision": {"darkvision", "Darkvision", 122, 2, "transmutation", []string{"druid", "ranger", "sorcerer", "wizard"}, false, false, "touch", 0, SpellTarget{"single", 1, "willing-creature-touch"}, "timed", 28800, "duration-ended", true, true, "a dried carrot", 0,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell202Effect{{"darkvision", "sense.darkvision-150", 0, 0, "", false}},
		[]string{"sense.darkvision-150", "components.material"}},
	"daylight": {"daylight", "Daylight", 122, 3, "evocation", []string{"cleric", "druid", "paladin", "ranger", "sorcerer"}, false, false, "feet", 60, SpellTarget{"choice", 1, "point-in-range-or-unattended-object"}, "timed", 3600, "duration-ended", true, true, "", 0,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell202Effect{{"daylight", "area.sunlight-sphere-60", 0, 0, "", false}, {"object-daylight", "area.sunlight-emanation-60", 0, 0, "", false}, {"dispel", "dispel.overlapping-darkness-spell-level-3-or-lower", 0, 0, "", false}},
		[]string{"area.sunlight-sphere-60", "area.sunlight-emanation-60", "dispel.overlapping-darkness-spell-level-3-or-lower"}},
	"death_ward": {"death-ward", "Death Ward", 122, 4, "abjuration", []string{"cleric", "paladin"}, false, false, "touch", 0, SpellTarget{"single", 1, "creature-touch"}, "timed", 28800, "duration-ended", true, true, "", 0,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell202Effect{{"drop-to-one", "protection.death-ward-drop-to-1", 0, 0, "", false}, {"instant-death", "protection.death-ward-negate-instant-death", 0, 0, "", false}},
		[]string{"protection.death-ward-drop-to-1", "protection.death-ward-negate-instant-death"}},
	"fireball": {"delayed-blast-fireball", "Delayed Blast Fireball", 122, 7, "evocation", []string{"sorcerer", "wizard"}, false, true, "feet", 150, SpellTarget{"sphere", 0, "creatures-in-sphere"}, "timed", 60, "duration-or-concentration-ended", true, true, "a ball of bat guano and sulfur", 0,
		SpellScaling{"slot", 7, 1, 6, []string{"damage"}}, []spell202Effect{{"damage", "damage.delayed-blast-fireball", 12, 6, "fire", false}, {"growth", "damage.delayed-blast-fireball-growth", 1, 6, "fire", false}, {"touch", "interaction.delayed-blast-bead-touch", 0, 0, "", false}, {"ignite", "ignite.flammable-objects", 0, 0, "", false}},
		[]string{"damage.delayed-blast-fireball", "damage.delayed-blast-fireball-growth", "interaction.delayed-blast-bead-touch", "ignite.flammable-objects", "area.sphere-20-radius", "save.dexterity-half", "scaling.per-slot", "concentration.source-duration", "components.material"}},
	"demiplane": {"demiplane", "Demiplane", 122, 8, "conjuration", []string{"sorcerer", "warlock", "wizard"}, false, false, "feet", 60, SpellTarget{"created-effect", 1, "visible-flat-solid-surface-in-range"}, "timed", 3600, "duration-ended", false, true, "", 0,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell202Effect{{"door", "portal.demiplane-door", 0, 0, "", false}, {"demiplane", "space.demiplane-room", 0, 0, "", false}, {"ending", "ending.demiplane-door-vanishes", 0, 0, "", false}},
		[]string{"portal.demiplane-door", "space.demiplane-room", "ending.demiplane-door-vanishes", "condition.prone-on-shunt"}},
	"detect_evil": {"detect-evil-and-good", "Detect Evil and Good", 123, 1, "divination", []string{"cleric", "paladin"}, false, true, "self", 0, SpellTarget{"single", 1, "self"}, "timed", 600, "duration-or-concentration-ended", true, true, "", 0,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell202Effect{{"sense", "sense.detect-evil-and-good", 0, 0, "", false}},
		[]string{"sense.detect-evil-and-good", "sense.hallow-active", "blocker.detection-barriers", "concentration.source-duration"}},
	"detect_magic": {"detect-magic", "Detect Magic", 123, 1, "divination", []string{"bard", "cleric", "druid", "paladin", "ranger", "sorcerer", "warlock", "wizard"}, true, true, "self", 0, SpellTarget{"single", 1, "self"}, "timed", 600, "duration-or-concentration-ended", true, true, "", 0,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell202Effect{{"sense", "sense.detect-magic", 0, 0, "", false}, {"aura", "sense.detect-magic-aura-school", 0, 0, "", false}},
		[]string{"sense.detect-magic", "sense.detect-magic-aura-school", "blocker.detection-barriers", "cast.ritual", "concentration.source-duration"}},
	"detect_poison": {"detect-poison-and-disease", "Detect Poison and Disease", 123, 1, "divination", []string{"cleric", "druid", "paladin", "ranger"}, true, true, "self", 0, SpellTarget{"single", 1, "self"}, "timed", 600, "duration-or-concentration-ended", true, true, "a yew leaf", 0,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell202Effect{{"sense", "sense.detect-poison-and-disease", 0, 0, "", false}},
		[]string{"sense.detect-poison-and-disease", "blocker.detection-barriers", "cast.ritual", "concentration.source-duration", "components.material"}},
	"detect_thoughts": {"detect-thoughts", "Detect Thoughts", 123, 2, "divination", []string{"bard", "sorcerer", "wizard"}, false, true, "self", 0, SpellTarget{"single", 1, "self"}, "timed", 60, "duration-or-concentration-ended", true, true, "1 Copper Piece", 1,
		SpellScaling{"none", 0, 0, 0, []string{}}, []spell202Effect{{"sense", "sense.thoughts", 0, 0, "", false}, {"read", "read.thoughts", 0, 0, "", false}, {"probe", "probe.thoughts-wisdom-save", 0, 0, "", false}},
		[]string{"sense.thoughts", "read.thoughts", "probe.thoughts-wisdom-save", "check.arcana-ends-probe", "blocker.detection-barriers", "concentration.source-duration", "components.material"}},
}

var spell202Blockers = []SpellDetectBlockersEntry{{"stone-dirt-or-wood", "1-foot"}, {"metal", "1-inch"}, {"lead", "thin-sheet"}}

func (s SpellRecord) validateSpell202() (bool, error) {
	expected := ""
	for key, row := range spell202Rows {
		if s.ID == prefix+"spells/"+row.slug {
			expected = key
		}
	}
	if s.SourceDetails == nil {
		// The bounded Cure Wounds compatibility fixture shares this ID.
		if expected != "" && s.Scope == "complete-source-record" {
			return true, fmt.Errorf("batch 202 details required")
		}
		return false, nil
	}
	d := s.SourceDetails
	payloads := map[string]any{
		"cure": d.Cure, "lights": d.Lights, "darkness": d.Darkness, "darkvision": d.Darkvision,
		"daylight": d.Daylight, "death_ward": d.DeathWard, "fireball": d.Fireball, "demiplane": d.Demiplane,
		"detect_evil": d.DetectEvil, "detect_magic": d.DetectMagic, "detect_poison": d.DetectPoison, "detect_thoughts": d.DetectThoughts,
	}
	count := 0
	for key, value := range payloads {
		if reflect.ValueOf(value).IsNil() {
			continue
		}
		count++
		if expected != key {
			return true, fmt.Errorf("mismatched batch 202 payload")
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
	if count != 1 || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" || d.CastingSeconds != 0 || d.CastingTriggers == nil || len(d.CastingTriggers) != 0 {
		return true, fmt.Errorf("invalid batch 202 source details")
	}
	row := spell202Rows[expected]
	if err := s.validateSpell202Header(row); err != nil {
		return true, err
	}
	if err := d.validateSpell202Payload(expected); err != nil {
		return true, err
	}
	return true, nil
}

func (s SpellRecord) validateSpell202Header(row spell202Row) error {
	flags := func(p *bool) bool { return p != nil && *p }
	c := s.Components
	if s.Name != row.name || s.SourcePage != row.page || s.Level != row.level || s.School != row.school || !slices.Equal(s.Lists, row.lists) ||
		s.CastingTime != "action" || flags(s.Ritual) != row.ritual || flags(s.Concentration) != row.concentration ||
		s.Range != (SpellRange{row.rangeKind, row.rangeFeet}) || s.Target != row.target || s.Duration != row.duration || s.DurationSeconds != row.seconds ||
		s.SourceDetails.EndTiming != row.endTiming || flags(c.Verbal) != row.verbal || flags(c.Somatic) != row.somatic {
		return fmt.Errorf("batch 202 header differs from source")
	}
	if m := c.Material; row.material == "" {
		if m != nil {
			return fmt.Errorf("batch 202 material differs from source")
		}
	} else if m == nil || m.Description != row.material || m.CostCP == nil || *m.CostCP != row.costCP || m.Consumed == nil || *m.Consumed {
		return fmt.Errorf("batch 202 material differs from source")
	}
	sc := s.Scaling
	if sc.Kind != row.scaling.Kind || sc.AboveLevel != row.scaling.AboveLevel || sc.DiceCount != row.scaling.DiceCount || sc.DieSides != row.scaling.DieSides || !slices.Equal(sc.EffectIDs, row.scaling.EffectIDs) {
		return fmt.Errorf("batch 202 scaling differs from source")
	}
	if len(s.Effects) != len(row.effects) {
		return fmt.Errorf("batch 202 effects differ from source")
	}
	for i, want := range row.effects {
		e := s.Effects[i]
		if e.ID != want.id || e.Mechanic != want.mechanic || e.DiceCount != want.dice || e.DieSides != want.sides || e.DamageType != want.damage || e.Bonus != 0 || e.AddCastingModifier == nil || *e.AddCastingModifier != want.modifier || e.HalfInitialOnly == nil || *e.HalfInitialOnly {
			return fmt.Errorf("batch 202 effect %s differs from source", want.id)
		}
	}
	if !slices.Equal(s.RequiredMechanics, row.required) {
		return fmt.Errorf("batch 202 required mechanics differ from source")
	}
	slots := []int{}
	for level := s.Level; level <= 9 && s.Level > 0; level++ {
		slots = append(slots, level)
	}
	if !slices.Equal(s.SlotLevels, slots) {
		return fmt.Errorf("batch 202 slot levels differ from source")
	}
	return nil
}

func (d *SpellSourceDetails) validateSpell202Payload(key string) error {
	bad := func() error { return fmt.Errorf("invalid batch 202 %s payload", key) }
	switch key {
	case "cure":
		c := d.Cure
		if c.Target != "creature-touched" || c.HealingDiceCount != 2 || c.HealingDieSides != 8 || c.HealingModifier != "spellcasting-ability-modifier" || c.SlotDiceCount != 2 || c.SlotDieSides != 8 {
			return bad()
		}
	case "lights":
		l := d.Lights
		if !slices.Equal(l.SourcePages, []int{121, 122}) || l.LightCount != 4 || l.LightSize != "torch" || !slices.Equal(l.Forms, []string{"torches", "lanterns", "glowing-orbs"}) || l.CombinedSize != "medium" || l.CombinedForm != "vaguely-humanlike" || l.LightLevel != "dim-light" || l.LightRadiusFeet != 10 || l.MoveAction != "bonus-action" || l.MoveFeet != 60 || l.MaxSeparationFeet != 20 || l.SeparationFrom != "another-light-created-by-spell" || l.MoveDestination != "space-within-range" || l.VanishesWhen != "exceeds-spell-range" {
			return bad()
		}
	case "darkness":
		x := d.Darkness
		if x.DarknessKind != "magical-darkness" || x.SphereRadiusFeet != 15 || x.EmanationFeet != 15 || x.SphereOrigin != "point-in-range" || x.Darkvision != "cannot-see-through" || x.NonmagicalLight != "cannot-illuminate" || x.ObjectAlternative != "object-not-worn-or-carried" || x.ObjectBlocker != "covered-by-opaque-item" || !slices.Equal(x.BlockerExamples, []string{"bowl", "helm"}) || !slices.Equal(x.DispelsLights, []string{"bright-light", "dim-light"}) || x.DispelsUpToLevel != 2 {
			return bad()
		}
	case "darkvision":
		x := d.Darkvision
		if x.Target != "willing-creature-touched" || x.Sense != "darkvision" || x.SenseRangeFeet != 150 {
			return bad()
		}
	case "daylight":
		x := d.Daylight
		if x.LightKind != "sunlight" || x.SphereRadiusFeet != 60 || x.EmanationFeet != 60 || x.SphereOrigin != "point-in-range" || x.SphereLight != "bright-light" || x.DimLightExtraFeet != 60 || x.ObjectAlternative != "object-not-worn-or-carried" || x.ObjectBlocker != "covered-by-opaque-item" || !slices.Equal(x.BlockerExamples, []string{"bowl", "helm"}) || x.DispelsDarkness != "darkness-from-spell" || x.DispelsUpToLevel != 3 {
			return bad()
		}
	case "death_ward":
		x := d.DeathWard
		if x.Target != "creature-touched" || x.DropTrigger != "would-drop-to-0-hit-points-first-time" || x.ReplacementHitPoints != 1 || !x.DropEndsSpell || x.InstantDeathTrigger != "effect-that-would-kill-instantly-without-dealing-damage" || x.InstantDeathResult != "effect-negated" || !x.InstantDeathEndsSpell {
			return bad()
		}
	case "fireball":
		x := d.Fireball
		if x.BeadOrigin != "beam-of-yellow-light-from-caster" || x.BeadPoint != "chosen-point-in-range" || x.ExplosionTrigger != "spell-ends" || x.RadiusFeet != 20 || x.Save != "dexterity" || x.DamageType != "fire" || x.FailDamage != "total-accumulated-damage" || x.SuccessDamage != "half-accumulated-damage" || x.BaseDiceCount != 12 || x.BaseDieSides != 6 || x.GrowthDiceCount != 1 || x.GrowthDieSides != 6 || x.GrowthTiming != "whenever-caster-turn-ends-and-spell-has-not-ended" || x.TouchSave != "dexterity" || x.TouchFailure != "spell-ends-and-bead-explodes" || x.TouchSuccess != "creature-can-throw-bead" || x.ThrowFeet != 40 || x.ThrownEnd != "enters-creature-space-or-collides-with-solid-object" || x.ThrownResult != "spell-ends-and-bead-explodes" || x.Ignites != "flammable-objects-not-worn-or-carried-in-explosion" || x.SlotBaseDiceCount != 1 || x.SlotBaseDieSides != 6 {
			return bad()
		}
	case "demiplane":
		x := d.Demiplane
		if !slices.Equal(x.SourcePages, []int{122, 123}) || x.Surface != "flat-solid-surface-visible-in-range" || x.DoorSize != "medium" || x.DoorAppearance != "shadowy" || x.DoorOperation != "opened-and-closed" || x.RoomKind != "empty-room" || x.RoomDimensionFeet != 30 || !slices.Equal(x.RoomMaterialChoices, []string{"wood", "stone"}) || x.SpellEndDoor != "vanishes" || x.ObjectsInside != "remain" || x.CreaturesInside != "remain-unless-they-opt-to-be-shunted-through-door" || x.ShuntedCondition != "prone" || x.ShuntedLanding != "unoccupied-spaces-closest-to-former-door-space" || !slices.Equal(x.ConnectionChoices, []string{"new-demiplane", "demiplane-from-previous-casting", "demiplane-by-another-creature"}) || x.OtherCreatureNeeds != "know-nature-and-contents" {
			return bad()
		}
	case "detect_evil":
		x := d.DetectEvil
		if x.SenseRadiusFeet != 30 || !slices.Equal(x.DetectedTypes, []string{"aberration", "celestial", "elemental", "fey", "fiend", "undead"}) || x.SensedInfo != "location" || x.HallowDetection != "whether-active-and-where" || !slices.Equal(x.Blockers, spell202Blockers) {
			return bad()
		}
	case "detect_magic":
		x := d.DetectMagic
		if x.SenseRadiusFeet != 30 || x.Senses != "presence-of-magical-effects" || x.FollowUpAction != "magic-action" || x.AuraSubjects != "visible-creature-or-object-bearing-the-magic" || x.AuraAppearance != "faint-aura" || x.SchoolLearned != "school-of-magic-if-effect-created-by-spell" || !slices.Equal(x.Blockers, spell202Blockers) {
			return bad()
		}
	case "detect_poison":
		x := d.DetectPoison
		if x.SenseRadiusFeet != 30 || !slices.Equal(x.Detected, []string{"poisons", "poisonous-or-venomous-creatures", "magical-contagions"}) || x.SensedInfo != "location" || x.KindSensed != "kind-of-poison-creature-or-contagion" || !slices.Equal(x.Blockers, spell202Blockers) {
			return bad()
		}
	case "detect_thoughts":
		x := d.DetectThoughts
		if x.ActivationAction != "magic-action-on-later-turns" || !slices.Equal(x.Options, []string{"sense-thoughts", "read-thoughts"}) || x.SenseRadiusFeet != 30 || x.SenseSubjects != "thoughts-of-creatures-that-know-languages-or-are-telepathic" || x.SenseReadsContent || !slices.Equal(x.SenseBlockers, spell202Blockers) ||
			x.ReadRangeFeet != 30 || !slices.Equal(x.ReadTargets, []string{"visible-creature", "creature-detected-with-sense-thoughts"}) || x.ReadLearns != "what-is-most-on-targets-mind-right-now" || x.ReadNoLanguageResult != "learn-nothing" ||
			x.ProbeAction != "magic-action-on-next-turn" || x.ProbeSave != "wisdom" || !slices.Equal(x.ProbeFailure, []string{"reasoning", "emotions", "something-that-looms-large"}) || !slices.Equal(x.ProbeFailureExamples, []string{"worry", "love", "hate"}) || x.ProbeSuccess != "spell-ends" || !x.TargetKnowsProbed ||
			x.TargetCheckAction != "action-on-its-turn" || x.TargetCheck != "intelligence-arcana" || x.TargetCheckDC != "caster-spell-save-dc" || x.TargetCheckSuccess != "spell-ends" || x.TargetCheckWindow != "until-caster-shifts-attention-from-target-mind" {
			return bad()
		}
	}
	return nil
}
