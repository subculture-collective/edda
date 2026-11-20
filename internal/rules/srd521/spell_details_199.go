package srd521

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// SpellLightning retains Call Lightning's typed SRD source facts without implementing them.
type SpellLightning struct {
	CloudShape          string   `json:"cloud_shape"`
	CloudHeightFeet     int      `json:"cloud_height_feet"`
	CloudRadiusFeet     int      `json:"cloud_radius_feet"`
	CloudPlacement      string   `json:"cloud_placement"`
	BoltPoint           string   `json:"bolt_point"`
	BoltRadiusFeet      int      `json:"bolt_radius_feet"`
	Save                string   `json:"save"`
	SaveSuccess         string   `json:"save_success"`
	InitialBolt         string   `json:"initial_bolt"`
	RepeatAction        string   `json:"repeat_action"`
	RepeatPointChoices  []string `json:"repeat_point_choices"`
	StormCondition      string   `json:"storm_condition"`
	StormEffect         string   `json:"storm_effect"`
	StormBonusDiceCount int      `json:"storm_bonus_dice_count"`
	StormBonusDieSides  int      `json:"storm_bonus_die_sides"`
}

// SpellCalmChoicesEntry retains one per-creature Calm Emotions alternative.
type SpellCalmChoicesEntry struct {
	Choice     string   `json:"choice"`
	Conditions []string `json:"conditions"`
	Existing   string   `json:"existing"`
	Attitude   string   `json:"attitude"`
	Subjects   string   `json:"subjects"`
	EndsEarly  []string `json:"ends_early"`
	OnSpellEnd string   `json:"on_spell_end"`
}

// SpellCalm retains Calm Emotions' typed SRD source facts without implementing them.
type SpellCalm struct {
	RadiusFeet        int                     `json:"radius_feet"`
	Center            string                  `json:"center"`
	CreatureType      string                  `json:"creature_type"`
	Save              string                  `json:"save"`
	SaveSuccess       string                  `json:"save_success"`
	ChoicePerCreature bool                    `json:"choice_per_creature"`
	ChooseCount       int                     `json:"choose_count"`
	Choices           []SpellCalmChoicesEntry `json:"choices"`
}

// SpellChain retains Chain Lightning's typed SRD source facts without implementing them.
type SpellChain struct {
	FirstTarget            string   `json:"first_target"`
	SecondaryBolts         int      `json:"secondary_bolts"`
	SecondaryTargets       string   `json:"secondary_targets"`
	SecondaryRangeFeet     int      `json:"secondary_range_feet"`
	SecondaryRangeOrigin   string   `json:"secondary_range_origin"`
	TargetKinds            []string `json:"target_kinds"`
	MaxBoltsPerTarget      int      `json:"max_bolts_per_target"`
	Save                   string   `json:"save"`
	SaveSuccess            string   `json:"save_success"`
	AdditionalBoltsPerSlot int      `json:"additional_bolts_per_slot"`
}

// SpellCharm retains Charm Monster's or Charm Person's typed SRD source facts
// without implementing them. The two spells differ in creature type and level.
type SpellCharm struct {
	SourcePages              []int    `json:"source_pages"`
	CreatureType             string   `json:"creature_type"`
	Save                     string   `json:"save"`
	SaveAdvantageIf          string   `json:"save_advantage_if"`
	SaveSuccess              string   `json:"save_success"`
	Condition                string   `json:"condition"`
	Attitude                 string   `json:"attitude"`
	EndsEarly                []string `json:"ends_early"`
	AfterSpell               string   `json:"after_spell"`
	BaseTargets              int      `json:"base_targets"`
	AdditionalTargetsPerSlot int      `json:"additional_targets_per_slot"`
}

// SpellChillUpgradesEntry retains one Cantrip Upgrade row.
type SpellChillUpgradesEntry struct {
	CharacterLevel int `json:"character_level"`
	DiceCount      int `json:"dice_count"`
	DieSides       int `json:"die_sides"`
}

// SpellChill retains Chill Touch's typed SRD source facts without implementing them.
type SpellChill struct {
	Attack          string                    `json:"attack"`
	Reach           string                    `json:"reach"`
	HitRestriction  string                    `json:"hit_restriction"`
	RestrictionEnds string                    `json:"restriction_ends"`
	Upgrades        []SpellChillUpgradesEntry `json:"upgrades"`
}

// SpellOrb retains Chromatic Orb's typed SRD source facts without implementing them.
type SpellOrb struct {
	DamageTypes      []string `json:"damage_types"`
	Attack           string   `json:"attack"`
	LeapTrigger      string   `json:"leap_trigger"`
	LeapDieSides     int      `json:"leap_die_sides"`
	LeapRangeFeet    int      `json:"leap_range_feet"`
	LeapTarget       string   `json:"leap_target"`
	LeapAttack       string   `json:"leap_attack"`
	LeapDamage       string   `json:"leap_damage"`
	BaseSlotLeaps    int      `json:"base_slot_leaps"`
	LeapAgainMinSlot int      `json:"leap_again_min_slot"`
	MaxLeaps         string   `json:"max_leaps"`
	TargetingLimit   string   `json:"targeting_limit"`
}

// SpellDeath retains Circle of Death's typed SRD source facts without implementing them.
type SpellDeath struct {
	RadiusFeet  int    `json:"radius_feet"`
	Center      string `json:"center"`
	Save        string `json:"save"`
	SaveSuccess string `json:"save_success"`
	Appearance  string `json:"appearance"`
}

// SpellClairvoyanceFocusChoicesEntry pairs a source-listed focus with its sense.
type SpellClairvoyanceFocusChoicesEntry struct {
	Focus string `json:"focus"`
	Sense string `json:"sense"`
}

// SpellClairvoyance retains Clairvoyance's typed SRD source facts without implementing them.
type SpellClairvoyance struct {
	RangeMiles           int                                  `json:"range_miles"`
	LocationChoices      []string                             `json:"location_choices"`
	FamiliarDefinition   string                               `json:"familiar_definition"`
	ObviousExamples      []string                             `json:"obvious_examples"`
	SensorInvisible      bool                                 `json:"sensor_invisible"`
	SensorIntangible     bool                                 `json:"sensor_intangible"`
	SensorInvulnerable   bool                                 `json:"sensor_invulnerable"`
	SensorMoves          bool                                 `json:"sensor_moves"`
	SenseChoices         []string                             `json:"sense_choices"`
	SenseUse             string                               `json:"sense_use"`
	SwitchAction         string                               `json:"switch_action"`
	FocusChoices         []SpellClairvoyanceFocusChoicesEntry `json:"focus_choices"`
	FocusMinimumCostCP   int                                  `json:"focus_minimum_cost_cp"`
	SensorSeenByExamples []string                             `json:"sensor_seen_by_examples"`
	VisibleSensor        string                               `json:"visible_sensor"`
}

// SpellClone retains Clone's typed SRD source facts without implementing them.
// Only the diamond is consumed; the vessel is a separate required object.
// The aggregate material.consumed flag means "includes a consumed component"; the split lives here.
type SpellClone struct {
	SourcePages             []int    `json:"source_pages"`
	TargetChoices           []string `json:"target_choices"`
	FleshMinimumCubicInches int      `json:"flesh_minimum_cubic_inches"`
	FormsIn                 string   `json:"forms_in"`
	GrowthDays              int      `json:"growth_days"`
	AgeChoices              []string `json:"age_choices"`
	InertState              string   `json:"inert_state"`
	TransferTrigger         string   `json:"transfer_trigger"`
	SoulConditions          []string `json:"soul_conditions"`
	Identical               []string `json:"identical"`
	Excluded                []string `json:"excluded"`
	OriginalRemains         string   `json:"original_remains"`
	DiamondMinimumCostCP    int      `json:"diamond_minimum_cost_cp"`
	DiamondConsumed         bool     `json:"diamond_consumed"`
	VesselMinimumCostCP     int      `json:"vessel_minimum_cost_cp"`
	VesselConsumed          bool     `json:"vessel_consumed"`
	VesselRequirements      []string `json:"vessel_requirements"`
}

// SpellCloudkill retains Cloudkill's typed SRD source facts without implementing them.
type SpellCloudkill struct {
	RadiusFeet       int      `json:"radius_feet"`
	Center           string   `json:"center"`
	Appearance       string   `json:"appearance"`
	Obscurement      string   `json:"obscurement"`
	DispersedBy      string   `json:"dispersed_by"`
	DispersalExample string   `json:"dispersal_example"`
	DispersalEffect  string   `json:"dispersal_effect"`
	Save             string   `json:"save"`
	SaveSuccess      string   `json:"save_success"`
	SaveTriggers     []string `json:"save_triggers"`
	SavesPerTurn     int      `json:"saves_per_turn"`
	MoveFeet         int      `json:"move_feet"`
	MoveDirection    string   `json:"move_direction"`
	MoveTiming       string   `json:"move_timing"`
}

// SpellSpray retains Color Spray's typed SRD source facts without implementing them.
type SpellSpray struct {
	ConeFeet      int    `json:"cone_feet"`
	Origin        string `json:"origin"`
	Save          string `json:"save"`
	SaveSuccess   string `json:"save_success"`
	Condition     string `json:"condition"`
	ConditionEnds string `json:"condition_ends"`
	Appearance    string `json:"appearance"`
}

func spell199Templates() map[string]any {
	return map[string]any{
		"lightning":    SpellLightning{},
		"calm":         SpellCalm{},
		"chain":        SpellChain{},
		"monster":      SpellCharm{},
		"person":       SpellCharm{},
		"chill":        SpellChill{},
		"orb":          SpellOrb{},
		"death":        SpellDeath{},
		"clairvoyance": SpellClairvoyance{},
		"clone":        SpellClone{},
		"cloudkill":    SpellCloudkill{},
		"spray":        SpellSpray{},
	}
}

// sourcePayloadCount counts every selected typed payload across all batches, so
// a record can never carry payloads belonging to two spells.
func (d *SpellSourceDetails) sourcePayloadCount() int {
	count := 0
	v := reflect.ValueOf(*d)
	for i := 0; i < v.NumField(); i++ {
		if f := v.Field(i); f.Kind() == reflect.Pointer && !f.IsNil() {
			count++
		}
	}
	return count
}

func (s SpellRecord) validateSpell199() (bool, error) {
	names := map[string]string{"lightning": "call-lightning", "calm": "calm-emotions", "chain": "chain-lightning", "monster": "charm-monster", "person": "charm-person", "chill": "chill-touch", "orb": "chromatic-orb", "death": "circle-of-death", "clairvoyance": "clairvoyance", "clone": "clone", "cloudkill": "cloudkill", "spray": "color-spray"}
	expected := ""
	for key, name := range names {
		if s.ID == prefix+"spells/"+name {
			expected = key
		}
	}
	if s.SourceDetails == nil {
		if expected != "" {
			return true, fmt.Errorf("batch 199 details required")
		}
		return false, nil
	}
	d := s.SourceDetails
	payloads := map[string]any{
		"lightning":    d.Lightning,
		"calm":         d.Calm,
		"chain":        d.Chain,
		"monster":      d.Monster,
		"person":       d.Person,
		"chill":        d.Chill,
		"orb":          d.Orb,
		"death":        d.Death,
		"clairvoyance": d.Clairvoyance,
		"clone":        d.Clone,
		"cloudkill":    d.Cloudkill,
		"spray":        d.Spray,
	}
	count := 0
	for key, value := range payloads {
		if reflect.ValueOf(value).IsNil() {
			continue
		}
		count++
		if expected != key {
			return true, fmt.Errorf("mismatched batch 199 payload")
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
		return true, fmt.Errorf("exactly one batch 199 payload required")
	}
	castingSeconds := map[string]int{"action": 0, "10-minutes": 600, "1-hour": 3600}
	seconds, ok := castingSeconds[s.CastingTime]
	if !ok || s.Scope != "complete-source-record" || d.CastingSeconds != seconds || len(d.CastingTriggers) != 0 || d.CastingTriggers == nil || strings.TrimSpace(d.EndTiming) == "" {
		return true, fmt.Errorf("invalid batch 199 source details")
	}
	if d.Lightning != nil {
		l := d.Lightning
		if l.CloudShape != "cylinder" || l.CloudHeightFeet != 10 || l.CloudRadiusFeet != 60 || l.BoltRadiusFeet != 5 || l.Save != "dexterity" || !slices.Equal(l.RepeatPointChoices, []string{"same-point", "different-point"}) || l.StormBonusDiceCount <= 0 || l.StormBonusDieSides <= 0 {
			return true, fmt.Errorf("invalid call lightning area, save or storm bonus")
		}
	}
	if d.Calm != nil {
		c := d.Calm
		choices := []string{}
		for _, v := range c.Choices {
			choices = append(choices, v.Choice)
		}
		if !slices.Equal(choices, []string{"condition-immunity", "indifference"}) || c.ChooseCount != 1 || !c.ChoicePerCreature || c.RadiusFeet != 20 || c.Save != "charisma" || c.CreatureType != "humanoid" {
			return true, fmt.Errorf("invalid calm emotions choices")
		}
	}
	if d.Chain != nil {
		c := d.Chain
		if c.SecondaryBolts != 3 || c.SecondaryRangeFeet != 30 || c.MaxBoltsPerTarget != 1 || c.Save != "dexterity" || !slices.Equal(c.TargetKinds, []string{"creature", "object"}) {
			return true, fmt.Errorf("invalid chain lightning bolts or targets")
		}
	}
	for _, charm := range []struct {
		payload *SpellCharm
		kind    string
		pages   []int
	}{{d.Monster, "any-creature", []int{114, 115}}, {d.Person, "humanoid", []int{115}}} {
		if c := charm.payload; c != nil && (c.CreatureType != charm.kind || !slices.Equal(c.SourcePages, charm.pages) || c.Save != "wisdom" || c.Condition != "charmed" || !slices.Equal(c.EndsEarly, []string{"caster-damages-target", "caster-allies-damage-target"})) {
			return true, fmt.Errorf("invalid charm target, pages or ending")
		}
	}
	if d.Chill != nil {
		want := []SpellChillUpgradesEntry{{5, 2, 10}, {11, 3, 10}, {17, 4, 10}}
		if s.Level != 0 || !slices.Equal(d.Chill.Upgrades, want) {
			return true, fmt.Errorf("invalid cantrip upgrades")
		}
	}
	if d.Orb != nil {
		o := d.Orb
		if !slices.Equal(o.DamageTypes, []string{"acid", "cold", "fire", "lightning", "poison", "thunder"}) || o.LeapDieSides != 8 || o.LeapRangeFeet != 30 || o.BaseSlotLeaps != 1 || o.LeapAgainMinSlot != 2 {
			return true, fmt.Errorf("invalid orb damage types or leaps")
		}
	}
	if d.Death != nil && (d.Death.RadiusFeet != 60 || d.Death.Save != "constitution") {
		return true, fmt.Errorf("invalid circle of death area or save")
	}
	if d.Clairvoyance != nil {
		c := d.Clairvoyance
		if c.RangeMiles != 1 || s.Range != (SpellRange{"feet", 5280}) || !slices.Equal(c.LocationChoices, []string{"familiar", "obvious-unfamiliar"}) || !slices.Equal(c.SenseChoices, []string{"seeing", "hearing"}) || !slices.Equal(c.FocusChoices, []SpellClairvoyanceFocusChoicesEntry{{"jeweled-horn", "hearing"}, {"glass-eye", "seeing"}}) || c.FocusMinimumCostCP <= 0 || !c.SensorInvisible || !c.SensorIntangible || !c.SensorInvulnerable || c.SensorMoves {
			return true, fmt.Errorf("invalid clairvoyance range, senses or focus")
		}
	}
	if d.Clone != nil {
		c := d.Clone
		m := s.Components.Material
		if !slices.Equal(c.SourcePages, []int{115, 116}) || !slices.Equal(c.TargetChoices, []string{"creature", "at-least-1-cubic-inch-of-its-flesh"}) || !slices.Equal(c.AgeChoices, []string{"same-age", "younger"}) || c.GrowthDays != 120 || c.FleshMinimumCubicInches <= 0 || c.DiamondMinimumCostCP <= 0 || c.VesselMinimumCostCP <= 0 || !c.DiamondConsumed || c.VesselConsumed || m == nil || *m.CostCP != c.DiamondMinimumCostCP+c.VesselMinimumCostCP || !*m.Consumed {
			return true, fmt.Errorf("invalid clone choices, growth or materials")
		}
	}
	if d.Cloudkill != nil {
		c := d.Cloudkill
		if c.RadiusFeet != 20 || c.SavesPerTurn != 1 || c.MoveFeet != 10 || c.Save != "constitution" || !slices.Equal(c.SaveTriggers, []string{"on-cast-in-sphere", "sphere-moves-into-space", "enters-sphere", "ends-turn-in-sphere"}) {
			return true, fmt.Errorf("invalid cloudkill area, movement or frequency")
		}
	}
	if d.Spray != nil && (d.Spray.ConeFeet != 15 || d.Spray.Save != "constitution" || d.Spray.Condition != "blinded") {
		return true, fmt.Errorf("invalid color spray area, save or condition")
	}
	return true, nil
}
