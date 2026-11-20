package srd521

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// SpellCommandOptionsEntry retains one source-listed Command option.
type SpellCommandOptionsEntry struct {
	Choice    string `json:"choice"`
	Effect    string `json:"effect"`
	Condition string `json:"condition"`
	EndsTurn  string `json:"ends_turn"`
}

// SpellCommand retains Command's typed SRD source facts without implementing them.
type SpellCommand struct {
	Words                    int                        `json:"words"`
	Save                     string                     `json:"save"`
	SaveSuccess              string                     `json:"save_success"`
	FollowTiming             string                     `json:"follow_timing"`
	Options                  []SpellCommandOptionsEntry `json:"options"`
	BaseTargets              int                        `json:"base_targets"`
	AdditionalTargetsPerSlot int                        `json:"additional_targets_per_slot"`
}

// SpellCommune retains Commune's typed SRD source facts without implementing them.
type SpellCommune struct {
	Contact               string `json:"contact"`
	QuestionLimit         int    `json:"question_limit"`
	QuestionKind          string `json:"question_kind"`
	AskTiming             string `json:"ask_timing"`
	Answer                string `json:"answer"`
	UnclearAnswer         string `json:"unclear_answer"`
	ShortPhraseAnswer     string `json:"short_phrase_answer"`
	RepeatCumulativePct   int    `json:"repeat_cumulative_percent"`
	RepeatAppliesFromCast int    `json:"repeat_applies_from_cast"`
	RepeatResetsOn        string `json:"repeat_resets_on"`
	RepeatEffect          string `json:"repeat_effect"`
	RitualAlternative     bool   `json:"ritual_alternative"`
}

// SpellNature retains Commune with Nature's typed SRD source facts without
// implementing them. Its text continues from PDF p.116 to p.117.
type SpellNature struct {
	SourcePages         []int    `json:"source_pages"`
	OutdoorRadiusMiles  int      `json:"outdoor_radius_miles"`
	UndergroundRadiusFt int      `json:"underground_radius_feet"`
	UndergroundSetting  string   `json:"underground_setting"`
	FailsWhere          string   `json:"fails_where"`
	FailureExamples     []string `json:"failure_examples"`
	ChooseCount         int      `json:"choose_count"`
	Facts               []string `json:"facts"`
	CreatureMinimumCR   int      `json:"creature_minimum_cr"`
	CreatureTypes       []string `json:"creature_types"`
	CreatureChoice      string   `json:"creature_choice"`
	PrevalentKinds      []string `json:"prevalent_kinds"`
	PrevalentChoice     string   `json:"prevalent_choice"`
	RitualAlternative   bool     `json:"ritual_alternative"`
}

// SpellComprehend retains Comprehend Languages' typed SRD source facts without implementing them.
type SpellComprehend struct {
	Understands        []string `json:"understands"`
	WrittenRequires    string   `json:"written_requires"`
	ReadingMinutesPage int      `json:"reading_minutes_per_page"`
	DecodesSymbols     bool     `json:"decodes_symbols"`
	DecodesSecrets     bool     `json:"decodes_secret_messages"`
	RitualAlternative  bool     `json:"ritual_alternative"`
}

// SpellCompulsion retains Compulsion's typed SRD source facts without implementing them.
type SpellCompulsion struct {
	TargetChoice    string `json:"target_choice"`
	Save            string `json:"save"`
	Condition       string `json:"condition"`
	DirectionAction string `json:"direction_action"`
	Direction       string `json:"direction"`
	Movement        string `json:"movement"`
	Route           string `json:"route"`
	RepeatSave      string `json:"repeat_save"`
	RepeatSuccess   string `json:"repeat_success"`
}

// SpellCone retains Cone of Cold's typed SRD source facts without implementing them.
type SpellCone struct {
	ConeFeet    int    `json:"cone_feet"`
	Origin      string `json:"origin"`
	Save        string `json:"save"`
	SaveSuccess string `json:"save_success"`
	KilledState string `json:"killed_state"`
}

// SpellConfusionBehaviorsEntry retains one row of the Confusion behavior table.
type SpellConfusionBehaviorsEntry struct {
	RollMin  int    `json:"roll_min"`
	RollMax  int    `json:"roll_max"`
	Behavior string `json:"behavior"`
}

// SpellConfusion retains Confusion's typed SRD source facts without implementing them.
type SpellConfusion struct {
	RadiusFeet               int                            `json:"radius_feet"`
	Center                   string                         `json:"center"`
	Save                     string                         `json:"save"`
	SaveSuccess              string                         `json:"save_success"`
	Restrictions             []string                       `json:"restrictions"`
	BehaviorDieSides         int                            `json:"behavior_die_sides"`
	BehaviorTiming           string                         `json:"behavior_timing"`
	Behaviors                []SpellConfusionBehaviorsEntry `json:"behaviors"`
	DirectionDieSides        int                            `json:"direction_die_sides"`
	Directions               []string                       `json:"directions"`
	RepeatSave               string                         `json:"repeat_save"`
	RepeatSuccess            string                         `json:"repeat_success"`
	AdditionalRadiusFeetSlot int                            `json:"additional_radius_feet_per_slot"`
}

// SpellAnimals retains Conjure Animals' typed SRD source facts without
// implementing them. Its text continues from PDF p.117 to p.118.
type SpellAnimals struct {
	SourcePages      []int    `json:"source_pages"`
	Size             string   `json:"size"`
	Form             string   `json:"form"`
	Placement        string   `json:"placement"`
	FormChoice       string   `json:"form_choice"`
	FormExamples     []string `json:"form_examples"`
	SaveAdvantage    string   `json:"save_advantage"`
	AdvantageRangeFt int      `json:"advantage_range_feet"`
	MoveFeet         int      `json:"move_feet"`
	MoveTiming       string   `json:"move_timing"`
	MoveDestination  string   `json:"move_destination"`
	TriggerRangeFeet int      `json:"trigger_range_feet"`
	Triggers         []string `json:"triggers"`
	Save             string   `json:"save"`
	SaveSuccess      string   `json:"save_success"`
	SavesPerTurn     int      `json:"saves_per_turn"`
}

// SpellCelestialLightsEntry retains one Conjure Celestial light.
type SpellCelestialLightsEntry struct {
	Choice      string `json:"choice"`
	Effect      string `json:"effect"`
	Save        string `json:"save"`
	SaveSuccess string `json:"save_success"`
}

// SpellCelestial retains Conjure Celestial's typed SRD source facts without implementing them.
type SpellCelestial struct {
	Appearance      string                      `json:"appearance"`
	RadiusFeet      int                         `json:"radius_feet"`
	HeightFeet      int                         `json:"height_feet"`
	Center          string                      `json:"center"`
	ChoicePerTarget string                      `json:"choice_per_target"`
	Lights          []SpellCelestialLightsEntry `json:"lights"`
	Illumination    string                      `json:"illumination"`
	MoveFeet        int                         `json:"move_feet"`
	MoveTiming      string                      `json:"move_timing"`
	Triggers        []string                    `json:"triggers"`
	AffectedPerTurn int                         `json:"affected_per_turn"`
}

// SpellElementalElementsEntry pairs a source-listed element with its damage type.
type SpellElementalElementsEntry struct {
	Element    string `json:"element"`
	DamageType string `json:"damage_type"`
}

// SpellElemental retains Conjure Elemental's typed SRD source facts without implementing them.
type SpellElemental struct {
	Size             string                        `json:"size"`
	Form             string                        `json:"form"`
	Placement        string                        `json:"placement"`
	Elements         []SpellElementalElementsEntry `json:"elements"`
	Triggers         []string                      `json:"triggers"`
	TriggerRangeFeet int                           `json:"trigger_range_feet"`
	TriggerCondition string                        `json:"trigger_condition"`
	Save             string                        `json:"save"`
	Condition        string                        `json:"condition"`
	ConditionEnds    string                        `json:"condition_ends"`
	RepeatSave       string                        `json:"repeat_save"`
	RepeatSuccess    string                        `json:"repeat_success"`
}

// SpellFey retains Conjure Fey's typed SRD source facts without implementing them.
type SpellFey struct {
	Size            string   `json:"size"`
	Origin          string   `json:"origin"`
	Form            string   `json:"form"`
	Placement       string   `json:"placement"`
	Attack          string   `json:"attack"`
	AttackRangeFeet int      `json:"attack_range_feet"`
	Condition       string   `json:"condition"`
	ConditionEnds   string   `json:"condition_ends"`
	FearSources     []string `json:"fear_sources"`
	LaterAction     string   `json:"later_action"`
	TeleportFeet    int      `json:"teleport_feet"`
	TeleportTo      string   `json:"teleport_to"`
	LaterAttack     string   `json:"later_attack"`
}

// SpellMinor retains Conjure Minor Elementals' typed SRD source facts without implementing them.
type SpellMinor struct {
	Origin        string   `json:"origin"`
	EmanationFeet int      `json:"emanation_feet"`
	Trigger       string   `json:"trigger"`
	DamageTypes   []string `json:"damage_types"`
	TypeChoice    string   `json:"type_choice"`
	Terrain       string   `json:"terrain"`
	TerrainFor    string   `json:"terrain_for"`
}

func spell200Templates() map[string]any {
	return map[string]any{
		"command":    SpellCommand{},
		"commune":    SpellCommune{},
		"nature":     SpellNature{},
		"comprehend": SpellComprehend{},
		"compulsion": SpellCompulsion{},
		"cone":       SpellCone{},
		"confusion":  SpellConfusion{},
		"animals":    SpellAnimals{},
		"celestial":  SpellCelestial{},
		"elemental":  SpellElemental{},
		"fey":        SpellFey{},
		"minor":      SpellMinor{},
	}
}

// validSpell200Target lists the target eligibilities this batch introduces.
func validSpell200Target(t SpellTarget) bool {
	switch t.Shape {
	case "created-effect":
		return t.Count == 1 && t.Eligibility == "unoccupied-space-in-range"
	case "emanation":
		return t.Count == 0 && t.Eligibility == "creatures-in-emanation"
	}
	return false
}

// spell200Identity pins each record's printed level, school, casting, range,
// target, ritual, concentration, duration, scaling and effect dice, so a record
// cannot adopt another spell's values. Effects read id:dice:sides:type:modifier.
var spell200Identity = map[string]struct {
	Name           string
	Level          int
	School         string
	CastingTime    string
	Range          SpellRange
	Target         SpellTarget
	Ritual         bool
	Concentration  bool
	Seconds        int
	Scaling        string
	ScalingDice    [2]int
	ScalingEffects []string
	Effects        []string
}{
	"command":    {"command", 1, "enchantment", "action", SpellRange{"feet", 60}, SpellTarget{"up-to-count", 1, "visible-creatures-in-range"}, false, false, 0, "slot-targets", [2]int{0, 0}, []string{"command"}, []string{"command:0:0::false"}},
	"commune":    {"commune", 5, "divination", "1-minute", SpellRange{"self", 0}, SpellTarget{"single", 1, "self"}, true, false, 60, "none", [2]int{0, 0}, []string{}, []string{"answers:0:0::false", "no-answer:0:0::false"}},
	"nature":     {"commune-with-nature", 5, "divination", "1-minute", SpellRange{"self", 0}, SpellTarget{"single", 1, "self"}, true, false, 0, "none", [2]int{0, 0}, []string{}, []string{"facts:0:0::false"}},
	"comprehend": {"comprehend-languages", 1, "divination", "action", SpellRange{"self", 0}, SpellTarget{"single", 1, "self"}, true, false, 3600, "none", [2]int{0, 0}, []string{}, []string{"understand:0:0::false"}},
	"compulsion": {"compulsion", 4, "enchantment", "action", SpellRange{"feet", 30}, SpellTarget{"selected-creatures", 0, "visible-creatures-in-range"}, false, true, 60, "none", [2]int{0, 0}, []string{}, []string{"charmed:0:0::false", "direction:0:0::false"}},
	"cone":       {"cone-of-cold", 5, "evocation", "action", SpellRange{"self", 0}, SpellTarget{"cone", 0, "creatures-in-cone"}, false, false, 0, "slot", [2]int{1, 8}, []string{"damage"}, []string{"damage:8:8:cold:false", "statue:0:0::false"}},
	"confusion":  {"confusion", 4, "enchantment", "action", SpellRange{"feet", 90}, SpellTarget{"sphere", 0, "creatures-in-sphere"}, false, true, 60, "slot-radius", [2]int{0, 0}, []string{"behavior"}, []string{"behavior:0:0::false"}},
	"animals":    {"conjure-animals", 3, "conjuration", "action", SpellRange{"feet", 60}, SpellTarget{"created-effect", 1, "visible-unoccupied-space-in-range"}, false, true, 600, "slot", [2]int{1, 10}, []string{"damage"}, []string{"damage:3:10:slashing:false", "pack-move:0:0::false", "advantage:0:0::false"}},
	"celestial":  {"conjure-celestial", 7, "conjuration", "action", SpellRange{"feet", 90}, SpellTarget{"created-effect", 1, "point-in-range"}, false, true, 600, "slot", [2]int{1, 12}, []string{"healing", "damage"}, []string{"healing:4:12::true", "damage:6:12:radiant:false", "light:0:0::false"}},
	"elemental":  {"conjure-elemental", 5, "conjuration", "action", SpellRange{"feet", 60}, SpellTarget{"created-effect", 1, "unoccupied-space-in-range"}, false, true, 600, "slot", [2]int{1, 8}, []string{"damage", "repeat-damage"}, []string{"damage:8:8:spirit-element:false", "restrained:0:0::false", "repeat-damage:4:8:spirit-element:false"}},
	"fey":        {"conjure-fey", 6, "conjuration", "action", SpellRange{"feet", 60}, SpellTarget{"created-effect", 1, "visible-unoccupied-space-in-range"}, false, true, 600, "slot", [2]int{1, 12}, []string{"damage"}, []string{"damage:3:12:psychic:true", "fear:0:0::false", "teleport:0:0::false"}},
	"minor":      {"conjure-minor-elementals", 4, "conjuration", "action", SpellRange{"self", 0}, SpellTarget{"emanation", 0, "creatures-in-emanation"}, false, true, 600, "slot", [2]int{1, 8}, []string{"extra-damage"}, []string{"extra-damage:2:8:chosen-type:false", "terrain:0:0::false"}},
}

func (s SpellRecord) validateSpell200() (bool, error) {
	expected := ""
	for key, id := range spell200Identity {
		if s.ID == prefix+"spells/"+id.Name {
			expected = key
		}
	}
	if s.SourceDetails == nil {
		if expected != "" {
			return true, fmt.Errorf("batch 200 details required")
		}
		return false, nil
	}
	d := s.SourceDetails
	payloads := map[string]any{
		"command":    d.Command,
		"commune":    d.Commune,
		"nature":     d.Nature,
		"comprehend": d.Comprehend,
		"compulsion": d.Compulsion,
		"cone":       d.Cone,
		"confusion":  d.Confusion,
		"animals":    d.Animals,
		"celestial":  d.Celestial,
		"elemental":  d.Elemental,
		"fey":        d.Fey,
		"minor":      d.Minor,
	}
	count := 0
	for key, value := range payloads {
		if reflect.ValueOf(value).IsNil() {
			continue
		}
		count++
		if expected != key {
			return true, fmt.Errorf("mismatched batch 200 payload")
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
	if count != 1 || d.sourcePayloadCount() != 1 {
		return true, fmt.Errorf("exactly one batch 200 payload required")
	}
	if s.Ritual == nil || s.Concentration == nil {
		return true, fmt.Errorf("explicit ritual and concentration flags required")
	}
	for _, e := range s.Effects {
		if e.AddCastingModifier == nil {
			return true, fmt.Errorf("explicit casting modifier flag required")
		}
	}
	id := spell200Identity[expected]
	seconds := map[string]int{"action": 0, "1-minute": 60}[s.CastingTime]
	if s.Scope != "complete-source-record" || s.Level != id.Level || s.School != id.School || s.CastingTime != id.CastingTime || *s.Ritual != id.Ritual || *s.Concentration != id.Concentration || s.DurationSeconds != id.Seconds || s.Scaling.Kind != id.Scaling || s.Range != id.Range || s.Target != id.Target || s.Scaling.DiceCount != id.ScalingDice[0] || s.Scaling.DieSides != id.ScalingDice[1] {
		return true, fmt.Errorf("batch 200 record differs from its printed identity")
	}
	effects := []string{}
	for _, e := range s.Effects {
		effects = append(effects, fmt.Sprintf("%s:%d:%d:%s:%t", e.ID, e.DiceCount, e.DieSides, e.DamageType, *e.AddCastingModifier))
	}
	if !slices.Equal(effects, id.Effects) || !slices.Equal(append([]string{}, s.Scaling.EffectIDs...), append([]string{}, id.ScalingEffects...)) {
		return true, fmt.Errorf("batch 200 effects differ from the printed dice")
	}
	if d.CastingSeconds != seconds || len(d.CastingTriggers) != 0 || d.CastingTriggers == nil || strings.TrimSpace(d.EndTiming) == "" {
		return true, fmt.Errorf("invalid batch 200 source details")
	}
	if (id.Seconds == 0) != (s.Duration == "instantaneous") {
		return true, fmt.Errorf("invalid batch 200 duration")
	}
	if c := d.Command; c != nil {
		choices := []string{}
		for _, o := range c.Options {
			choices = append(choices, o.Choice)
		}
		if !slices.Equal(choices, []string{"approach", "drop", "flee", "grovel", "halt"}) || c.Words != 1 || c.Save != "wisdom" || c.BaseTargets != 1 || c.AdditionalTargetsPerSlot != 1 {
			return true, fmt.Errorf("invalid command options, save or targets")
		}
	}
	if c := d.Commune; c != nil && (c.QuestionLimit != 3 || c.RepeatCumulativePct != 25 || c.RepeatAppliesFromCast != 2 || !c.RitualAlternative) {
		return true, fmt.Errorf("invalid commune questions or repeat chance")
	}
	if n := d.Nature; n != nil {
		if !slices.Equal(n.SourcePages, []int{116, 117}) || n.OutdoorRadiusMiles != 3 || n.UndergroundRadiusFt != 300 || n.ChooseCount != 3 || len(n.Facts) != 5 || n.CreatureMinimumCR != 10 || !slices.Equal(n.CreatureTypes, []string{"celestial", "elemental", "fey", "fiend", "undead"}) || !slices.Equal(n.PrevalentKinds, []string{"plant", "mineral", "beast"}) || !n.RitualAlternative {
			return true, fmt.Errorf("invalid commune with nature area, facts or pages")
		}
	}
	if c := d.Comprehend; c != nil && (len(c.Understands) != 2 || c.ReadingMinutesPage != 1 || c.DecodesSymbols || c.DecodesSecrets || !c.RitualAlternative) {
		return true, fmt.Errorf("invalid comprehend languages reading limits")
	}
	if c := d.Compulsion; c != nil && (c.Save != "wisdom" || c.Condition != "charmed") {
		return true, fmt.Errorf("invalid compulsion save or condition")
	}
	if c := d.Cone; c != nil && (c.ConeFeet != 60 || c.Save != "constitution") {
		return true, fmt.Errorf("invalid cone of cold area or save")
	}
	if c := d.Confusion; c != nil {
		rows := [][2]int{}
		for _, b := range c.Behaviors {
			rows = append(rows, [2]int{b.RollMin, b.RollMax})
		}
		if c.RadiusFeet != 10 || c.Save != "wisdom" || c.BehaviorDieSides != 10 || c.DirectionDieSides != 4 || c.AdditionalRadiusFeetSlot != 5 || !slices.Equal(rows, [][2]int{{1, 1}, {2, 6}, {7, 8}, {9, 10}}) || !slices.Equal(c.Directions, []string{"north", "east", "south", "west"}) {
			return true, fmt.Errorf("invalid confusion area, table or scaling")
		}
	}
	if a := d.Animals; a != nil && (!slices.Equal(a.SourcePages, []int{117, 118}) || a.Size != "large" || a.MoveFeet != 30 || a.TriggerRangeFeet != 10 || a.AdvantageRangeFt != 5 || a.Save != "dexterity" || a.SavesPerTurn != 1 || len(a.Triggers) != 3) {
		return true, fmt.Errorf("invalid conjure animals pack, triggers or pages")
	}
	if c := d.Celestial; c != nil {
		lights := []string{}
		for _, l := range c.Lights {
			lights = append(lights, l.Choice)
		}
		if c.RadiusFeet != 10 || c.HeightFeet != 40 || c.MoveFeet != 30 || c.AffectedPerTurn != 1 || !slices.Equal(lights, []string{"healing-light", "searing-light"}) || len(c.Triggers) != 3 {
			return true, fmt.Errorf("invalid conjure celestial cylinder or lights")
		}
	}
	if e := d.Elemental; e != nil {
		want := []SpellElementalElementsEntry{{"air", "lightning"}, {"earth", "thunder"}, {"fire", "fire"}, {"water", "cold"}}
		if e.Size != "large" || e.TriggerRangeFeet != 5 || e.Save != "dexterity" || e.Condition != "restrained" || !slices.Equal(e.Elements, want) || len(e.Triggers) != 2 {
			return true, fmt.Errorf("invalid conjure elemental elements or restraint")
		}
	}
	if f := d.Fey; f != nil && (f.Size != "medium" || f.AttackRangeFeet != 5 || f.TeleportFeet != 30 || f.Condition != "frightened" || !slices.Equal(f.FearSources, []string{"caster", "spirit"})) {
		return true, fmt.Errorf("invalid conjure fey attack, fear or teleport")
	}
	if m := d.Minor; m != nil && (m.EmanationFeet != 15 || !slices.Equal(m.DamageTypes, []string{"acid", "cold", "fire", "lightning"})) {
		return true, fmt.Errorf("invalid conjure minor elementals area or damage types")
	}
	return true, nil
}
