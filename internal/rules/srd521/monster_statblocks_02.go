package srd521

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

//go:embed testdata/monsters-statblocks-02.json
var monsterBatch02JSON []byte

// MonsterBatch02 is stat-block batch 02 (#152): Allosaurus to Animals. It is
// source-backed catalogue data and grants no NPC behavior.
type MonsterBatch02 struct {
	SchemaVersion int                       `json:"schema_version"`
	Source        Source                    `json:"source"`
	Headings      []MonsterBatch02Heading   `json:"headings"`
	StatBlocks    []MonsterBatch02StatBlock `json:"stat_blocks"`
}

// MonsterBatch02Heading is an indexed source entry that prints no stat block.
type MonsterBatch02Heading struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	SourcePages []int  `json:"source_pages"`
}

// MonsterBatch02StatBlock pairs a contract Monster with the typed source fields
// the contract has no member for. ActionDetails follows Monster.Actions in order.
type MonsterBatch02StatBlock struct {
	Monster               Monster                      `json:"monster"`
	Tags                  []string                     `json:"tags"`
	Skills                map[string]int               `json:"skills"`
	DamageVulnerabilities []string                     `json:"damage_vulnerabilities"`
	DamageResistances     []string                     `json:"damage_resistances"`
	DamageImmunities      []string                     `json:"damage_immunities"`
	ConditionImmunities   []string                     `json:"condition_immunities"`
	LairXP                int                          `json:"lair_xp"`
	Legendary             *MonsterBatch02Legendary     `json:"legendary"`
	ActionDetails         []MonsterBatch02ActionDetail `json:"action_details"`
}

// MonsterBatch02Legendary retains Legendary Action and Resistance uses and the
// printed rule for spending them.
type MonsterBatch02Legendary struct {
	Uses               int    `json:"uses"`
	LairUses           int    `json:"lair_uses"`
	ResistanceUses     int    `json:"resistance_uses"`
	ResistanceLairUses int    `json:"resistance_lair_uses"`
	Text               string `json:"text"`
}

// MonsterBatch02ActionDetail types the numbers inside one action's source text.
// Zero, empty and null members mean the source prints no such value.
type MonsterBatch02ActionDetail struct {
	ActionID         string                      `json:"action_id"`
	Section          string                      `json:"section"`
	Recharge         *MonsterBatch02Recharge     `json:"recharge"`
	Save             *MonsterBatch02Save         `json:"save"`
	Area             *MonsterBatch02Area         `json:"area"`
	TargetRangeFeet  int                         `json:"target_range_feet"`
	Damage           []MonsterBatch02Damage      `json:"damage"`
	Conditions       []string                    `json:"conditions"`
	PushFeet         int                         `json:"push_feet"`
	DurationMinutes  int                         `json:"duration_minutes"`
	Penalties        []MonsterBatch02Penalty     `json:"penalties"`
	OncePerTurnCycle bool                        `json:"once_per_turn_cycle"`
	Movement         string                      `json:"movement"`
	Attacks          []MonsterBatch02AttackUse   `json:"attacks"`
	Replacement      *MonsterBatch02Replacement  `json:"replacement"`
	Casts            []MonsterBatch02Cast        `json:"casts"`
	Spellcasting     *MonsterBatch02Spellcasting `json:"spellcasting"`
}

// MonsterBatch02Recharge is the inclusive d6 range that restores an action.
type MonsterBatch02Recharge struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// MonsterBatch02Save is the saving throw an action calls for.
type MonsterBatch02Save struct {
	Ability       string `json:"ability"`
	DC            int    `json:"dc"`
	HalfOnSuccess bool   `json:"half_on_success"`
}

// MonsterBatch02Area is a Line (length and width), Cone (length) or Sphere
// (radius, centered within RangeFeet).
type MonsterBatch02Area struct {
	Shape     string `json:"shape"`
	Feet      int    `json:"feet"`
	WidthFeet int    `json:"width_feet"`
	RangeFeet int    `json:"range_feet"`
}

// MonsterBatch02Damage is one printed damage roll and what triggers it.
type MonsterBatch02Damage struct {
	Average  int    `json:"average"`
	Dice     int    `json:"dice"`
	Sides    int    `json:"sides"`
	Modifier int    `json:"modifier"`
	Type     string `json:"type"`
	Trigger  string `json:"trigger"`
}

// MonsterBatch02Penalty is a flat or rolled penalty the action imposes.
type MonsterBatch02Penalty struct {
	AppliesTo string `json:"applies_to"`
	Flat      int    `json:"flat"`
	Average   int    `json:"average"`
	Dice      int    `json:"dice"`
	Sides     int    `json:"sides"`
}

// MonsterBatch02AttackUse names another action of the block and how often it is made.
type MonsterBatch02AttackUse struct {
	ActionID string `json:"action_id"`
	Count    int    `json:"count"`
}

// MonsterBatch02Replacement lists what Multiattack may substitute for attacks.
type MonsterBatch02Replacement struct {
	Attacks int                           `json:"attacks"`
	Options []MonsterBatch02ReplaceOption `json:"options"`
}

// MonsterBatch02ReplaceOption is one substitute action; Spell is set for Spellcasting.
type MonsterBatch02ReplaceOption struct {
	ActionID string `json:"action_id"`
	Spell    string `json:"spell"`
	Level    int    `json:"level"`
}

// MonsterBatch02Cast is a spell an action casts. Level 0 means no stated version.
type MonsterBatch02Cast struct {
	Spell           string `json:"spell"`
	Level           int    `json:"level"`
	ViaSpellcasting bool   `json:"via_spellcasting"`
}

// MonsterBatch02Spellcasting retains the casting ability, DC, optional spell
// attack bonus and the printed spell list.
type MonsterBatch02Spellcasting struct {
	Ability     string                `json:"ability"`
	SaveDC      int                   `json:"save_dc"`
	AttackBonus *int                  `json:"attack_bonus"`
	Spells      []MonsterBatch02Spell `json:"spells"`
}

// MonsterBatch02Spell is one listed spell. Restriction is source text.
type MonsterBatch02Spell struct {
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Frequency   string `json:"frequency"`
	Restriction string `json:"restriction"`
}

var (
	monsterBatch02Damage     = strings.Fields("acid bludgeoning cold fire force lightning necrotic piercing poison psychic radiant slashing thunder")
	monsterBatch02Conditions = strings.Fields("blinded charmed deafened exhaustion frightened grappled incapacitated invisible paralyzed petrified poisoned prone restrained stunned unconscious")
	monsterBatch02Abilities  = strings.Fields("str dex con int wis cha")
	monsterBatch02Skills     = map[string]string{"deception": "cha", "history": "int", "insight": "wis", "perception": "wis", "persuasion": "cha", "stealth": "dex"}
)

// LoadMonsterBatch02 returns a fresh validated copy of the embedded batch.
func LoadMonsterBatch02() (MonsterBatch02, error) {
	return DecodeMonsterBatch02(bytes.NewReader(monsterBatch02JSON))
}

// DecodeMonsterBatch02 strictly decodes one batch document: no duplicate or
// unknown members, no trailing data, and every top-level member present.
func DecodeMonsterBatch02(r io.Reader) (MonsterBatch02, error) {
	var b MonsterBatch02
	data, err := io.ReadAll(r)
	if err != nil {
		return b, err
	}
	if err := uniqueItemJSON(json.NewDecoder(bytes.NewReader(data))); err != nil {
		return b, err
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&b); err != nil {
		return MonsterBatch02{}, fmt.Errorf("decode stat-block batch 02: %w", err)
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return MonsterBatch02{}, fmt.Errorf("stat-block batch 02 must contain one JSON object")
	}
	if err := b.Validate(); err != nil {
		return MonsterBatch02{}, err
	}
	return b, nil
}

// UnmarshalJSON rejects an omitted member, so a typed source field cannot
// disappear into a zero value.
func (s *MonsterBatch02StatBlock) UnmarshalJSON(data []byte) error {
	type plain MonsterBatch02StatBlock
	return decodeMonsterBatch02(data, (*plain)(s), "monster", "tags", "skills", "damage_vulnerabilities", "damage_resistances",
		"damage_immunities", "condition_immunities", "lair_xp", "legendary", "action_details")
}

// UnmarshalJSON requires every detail member, including explicit nulls and zeros.
func (a *MonsterBatch02ActionDetail) UnmarshalJSON(data []byte) error {
	type plain MonsterBatch02ActionDetail
	return decodeMonsterBatch02(data, (*plain)(a), "action_id", "section", "recharge", "save", "area", "target_range_feet", "damage",
		"conditions", "push_feet", "duration_minutes", "penalties", "once_per_turn_cycle", "movement", "attacks", "replacement",
		"casts", "spellcasting")
}

func decodeMonsterBatch02(data []byte, dst any, required ...string) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	for _, key := range required {
		if _, ok := fields[key]; !ok {
			return fmt.Errorf("explicit stat-block member %s required", key)
		}
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	return d.Decode(dst)
}

// Validate checks provenance, identities and every typed field. It never
// consults or grants runtime support.
func (b MonsterBatch02) Validate() error {
	if b.SchemaVersion != 1 {
		return fmt.Errorf("unsupported stat-block batch schema %d", b.SchemaVersion)
	}
	if err := validateSource(b.Source); err != nil {
		return err
	}
	if len(b.StatBlocks) == 0 || b.Headings == nil {
		return fmt.Errorf("stat blocks and explicit headings required")
	}
	seen := map[string]bool{}
	for _, h := range b.Headings {
		if seen[h.ID] || !contentID.MatchString(h.ID) || !strings.HasPrefix(h.ID, "srd-5.2.1:content/") || strings.TrimSpace(h.Name) == "" || h.Kind != "section-heading" ||
			len(h.SourcePages) != 2 || h.SourcePages[0] < 1 || h.SourcePages[0] > h.SourcePages[1] || h.SourcePages[1] > 364 {
			return fmt.Errorf("invalid or duplicate heading %q", h.ID)
		}
		seen[h.ID] = true
	}
	for _, s := range b.StatBlocks {
		if seen[s.Monster.ID] {
			return fmt.Errorf("duplicate stat-block ID %q", s.Monster.ID)
		}
		seen[s.Monster.ID] = true
		if err := s.Validate(); err != nil {
			return fmt.Errorf("%s: %w", s.Monster.ID, err)
		}
	}
	return nil
}

// Validate checks one block against the Monster contract and then every
// derived number the source lets us recompute.
func (s MonsterBatch02StatBlock) Validate() error {
	m := s.Monster
	if err := m.Validate(); err != nil {
		return err
	}
	needs := func(key string) error {
		if !validMonsterKey(key) || !slices.Contains(m.RequiredMechanics, key) {
			return fmt.Errorf("missing mechanic dependency %q", key)
		}
		return nil
	}
	if !encounterValidXP(m.CR, int64(m.XP)) {
		return fmt.Errorf("XP does not match CR")
	}
	for _, key := range monsterBatch02Abilities {
		a := m.Abilities[key]
		if a.Save != a.Modifier && a.Save != a.Modifier+m.ProficiencyBonus {
			return fmt.Errorf("save %s is neither unproficient nor proficient", key)
		}
	}
	for kind := range m.Speeds {
		if kind != "walk" {
			if err := needs("speed/" + kind); err != nil {
				return err
			}
		}
	}
	for sense := range m.Senses.Special {
		if err := needs("sense/" + sense); err != nil {
			return err
		}
	}
	if s.Tags == nil || s.Skills == nil {
		return fmt.Errorf("explicit tags and skills required")
	}
	for _, tag := range s.Tags {
		if !familyName.MatchString(tag) {
			return fmt.Errorf("invalid tag %q", tag)
		}
	}
	for skill, bonus := range s.Skills {
		ability, ok := monsterBatch02Skills[skill]
		mod := m.Abilities[ability].Modifier
		if !ok || (bonus != mod+m.ProficiencyBonus && bonus != mod+2*m.ProficiencyBonus) {
			return fmt.Errorf("skill %s is neither proficient nor expert", skill)
		}
		if err := needs("skill/" + skill); err != nil {
			return err
		}
	}
	passive := 10 + m.Abilities["wis"].Modifier
	if bonus, ok := s.Skills["perception"]; ok {
		passive = 10 + bonus
	}
	if m.Senses.PassivePerception != passive {
		return fmt.Errorf("passive perception does not follow Perception")
	}
	for _, list := range []struct {
		values []string
		valid  []string
		prefix string
		suffix string
	}{
		{s.DamageVulnerabilities, monsterBatch02Damage, "vulnerability/", "-damage"},
		{s.DamageResistances, monsterBatch02Damage, "resistance/", "-damage"},
		{s.DamageImmunities, monsterBatch02Damage, "immunity/", "-damage"},
		{s.ConditionImmunities, monsterBatch02Conditions, "immunity/", "-condition"},
	} {
		if list.values == nil {
			return fmt.Errorf("explicit vulnerabilities, resistances and immunities required")
		}
		for i, v := range list.values {
			if !slices.Contains(list.valid, v) || slices.Contains(list.values[:i], v) {
				return fmt.Errorf("invalid or duplicate %s%s", list.prefix, v)
			}
			if err := needs(list.prefix + v + list.suffix); err != nil {
				return err
			}
		}
	}
	if s.LairXP != 0 {
		cr, err := strconv.Atoi(m.CR)
		if err != nil || !encounterValidXP(strconv.Itoa(cr+1), int64(s.LairXP)) {
			return fmt.Errorf("lair XP must be the next challenge rating's XP")
		}
		if err := needs("monster/lair-experience"); err != nil {
			return err
		}
	}
	if len(s.ActionDetails) != len(m.Actions) {
		return fmt.Errorf("one action detail per action required")
	}
	legendaryActions, resistance := 0, false
	for _, t := range m.Traits {
		resistance = resistance || t.Key == "trait/legendary-resistance"
	}
	for i, d := range s.ActionDetails {
		a := m.Actions[i]
		if d.ActionID != a.ID {
			return fmt.Errorf("action detail %d does not follow action %q", i, a.ID)
		}
		if err := s.validateDetail(a, d); err != nil {
			return fmt.Errorf("%s: %w", a.ID, err)
		}
		if d.Section == "legendary-action" {
			legendaryActions++
		} else if legendaryActions > 0 {
			return fmt.Errorf("legendary actions must follow ordinary actions")
		}
	}
	if l := s.Legendary; l == nil {
		if legendaryActions != 0 || resistance {
			return fmt.Errorf("legendary actions or resistance without legendary uses")
		}
	} else {
		if legendaryActions == 0 || !resistance || l.Uses < 1 || l.LairUses < l.Uses || l.ResistanceUses < 1 || l.ResistanceLairUses < l.ResistanceUses || strings.TrimSpace(l.Text) == "" {
			return fmt.Errorf("invalid legendary uses")
		}
		printed := fmt.Sprintf("(%d/Day, or %d/Day in Lair)", l.ResistanceUses, l.ResistanceLairUses)
		if !slices.ContainsFunc(m.Traits, func(t MonsterMechanic) bool {
			return t.Key == "trait/legendary-resistance" && strings.Contains(t.Text, printed)
		}) || !strings.HasPrefix(l.Text, fmt.Sprintf("Legendary Action Uses: %d (%d in Lair). ", l.Uses, l.LairUses)) {
			return fmt.Errorf("legendary uses do not match the retained source text")
		}
		for _, key := range []string{"monster/legendary-actions", "monster/lair-legendary-uses"} {
			if err := needs(key); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s MonsterBatch02StatBlock) validateDetail(a MonsterAction, d MonsterBatch02ActionDetail) error {
	m := s.Monster
	has := func(key string) bool {
		return slices.ContainsFunc(a.Effects, func(e MonsterMechanic) bool { return e.Key == key })
	}
	action := func(id string) bool {
		return slices.ContainsFunc(m.Actions, func(o MonsterAction) bool { return o.ID == id })
	}
	if d.Damage == nil || d.Conditions == nil || d.Penalties == nil || d.Attacks == nil || d.Casts == nil {
		return fmt.Errorf("explicit damage, conditions, penalties, attacks and casts required")
	}
	switch d.Section {
	case "action":
		if strings.HasPrefix(a.Key, "legendary-action/") {
			return fmt.Errorf("section does not match key")
		}
	case "legendary-action":
		if a.Key != "legendary-action/"+a.ID {
			return fmt.Errorf("section does not match key")
		}
	default:
		return fmt.Errorf("invalid section %q", d.Section)
	}
	for _, e := range a.Effects {
		if !strings.Contains(a.Text, e.Text) {
			return fmt.Errorf("effect %q is not a clause of the action text", e.Key)
		}
	}
	if (d.Recharge != nil) != has("effect/recharge") || (d.Recharge != nil && (d.Recharge.Min < 2 || d.Recharge.Min > d.Recharge.Max || d.Recharge.Max != 6)) {
		return fmt.Errorf("invalid recharge")
	}
	if (d.Save != nil) != has("effect/saving-throw") {
		return fmt.Errorf("save and saving-throw effect must agree")
	}
	if d.Save != nil && (!slices.Contains(monsterBatch02Abilities, d.Save.Ability) || d.Save.DC < 1 || d.Save.DC > 30 || d.Save.HalfOnSuccess != has("effect/half-damage-on-success")) {
		return fmt.Errorf("invalid save")
	}
	if d.Area != nil {
		ar := d.Area
		ok := ar.Feet > 0 && has("effect/"+ar.Shape+"-area")
		switch ar.Shape {
		case "line":
			ok = ok && ar.WidthFeet > 0 && ar.RangeFeet == 0
		case "cone":
			ok = ok && ar.WidthFeet == 0 && ar.RangeFeet == 0
		case "sphere":
			ok = ok && ar.WidthFeet == 0 && ar.RangeFeet > 0
		default:
			ok = false
		}
		if !ok || d.Save == nil || d.TargetRangeFeet != 0 {
			return fmt.Errorf("invalid area")
		}
	} else if has("effect/line-area") || has("effect/cone-area") || has("effect/sphere-area") {
		return fmt.Errorf("area effect without typed area")
	}
	if d.TargetRangeFeet < 0 || (d.TargetRangeFeet > 0) != has("effect/single-target-in-range") {
		return fmt.Errorf("invalid target range")
	}
	if (a.Attack != nil) != (len(d.Damage) > 0 && d.Damage[0].Trigger == "hit") {
		return fmt.Errorf("attack and hit damage must agree")
	}
	damageEffects := 0
	for _, e := range a.Effects {
		if strings.HasSuffix(e.Key, "-damage") {
			damageEffects++
		}
	}
	if damageEffects != len(d.Damage) {
		return fmt.Errorf("every damage effect needs a typed roll")
	}
	for i, roll := range d.Damage {
		if !validMonsterDice(MonsterHP{roll.Average, roll.Dice, roll.Sides, roll.Modifier}) || !slices.Contains(monsterBatch02Damage, roll.Type) || !has("effect/"+roll.Type+"-damage") {
			return fmt.Errorf("invalid damage roll %d", i)
		}
		if a.Attack != nil {
			at := a.Attack
			str := m.Abilities["str"].Modifier
			if roll.Trigger != "hit" || (i == 0 && (at.Damage != MonsterHP{roll.Average, roll.Dice, roll.Sides, roll.Modifier} || at.DamageType != roll.Type || at.Bonus != str+m.ProficiencyBonus || roll.Modifier != str)) || (i > 0 && roll.Modifier != 0) {
				return fmt.Errorf("attack damage %d does not follow the typed attack", i)
			}
		} else if roll.Trigger != "failure" || d.Save == nil || roll.Modifier != 0 {
			return fmt.Errorf("save damage %d needs a failed save", i)
		}
	}
	for i, c := range d.Conditions {
		if !slices.Contains(monsterBatch02Conditions, c) || slices.Contains(d.Conditions[:i], c) || !has("effect/"+c+"-condition") {
			return fmt.Errorf("invalid condition %q", c)
		}
	}
	if d.PushFeet < 0 || (d.PushFeet > 0) != has("effect/forced-movement") || d.DurationMinutes < 0 {
		return fmt.Errorf("invalid push or duration")
	}
	for _, p := range d.Penalties {
		rolled := p.Dice > 0 && slices.Contains([]int{4, 6, 8, 10, 12, 20}, p.Sides) && p.Flat == 0 && (p.Average == 0 || p.Average == p.Dice*(p.Sides+1)/2)
		flat := p.Flat > 0 && p.Dice == 0 && p.Sides == 0 && p.Average == 0
		if !slices.Contains([]string{"armor-class", "d20-test", "damage-roll"}, p.AppliesTo) || !(rolled || flat) || !has("effect/"+p.AppliesTo+"-penalty") {
			return fmt.Errorf("invalid penalty")
		}
	}
	if d.OncePerTurnCycle != has("effect/once-per-turn-cycle") || (d.OncePerTurnCycle && d.Section != "legendary-action") {
		return fmt.Errorf("invalid once-per-turn-cycle limit")
	}
	if !slices.Contains([]string{"", "half-speed", "half-fly-speed"}, d.Movement) || (d.Movement != "" && !has("effect/"+d.Movement+"-movement")) {
		return fmt.Errorf("invalid movement")
	}
	for _, use := range d.Attacks {
		if use.Count < 1 || use.ActionID == a.ID || !action(use.ActionID) {
			return fmt.Errorf("attack use names no other action")
		}
	}
	if (a.Key == "action/multiattack") != (a.ID == "multiattack") || (a.ID == "multiattack" && len(d.Attacks) == 0) {
		return fmt.Errorf("Multiattack needs its attack sequence")
	}
	if (d.Replacement != nil) != has("effect/multiattack-replacement") {
		return fmt.Errorf("replacement and effect must agree")
	}
	if r := d.Replacement; r != nil {
		if a.ID != "multiattack" || r.Attacks < 1 || len(r.Options) == 0 {
			return fmt.Errorf("invalid replacement")
		}
		for _, o := range r.Options {
			if !action(o.ActionID) || (o.ActionID == "spellcasting") != (o.Spell != "") || o.Level < 0 || o.Level > 9 || (o.Spell != "" && !s.lists(o.Spell, o.Level)) {
				return fmt.Errorf("replacement option %q is not in the block", o.ActionID)
			}
		}
	}
	for _, c := range d.Casts {
		if c.Level < 0 || c.Level > 9 || !has("effect/spell-"+monsterBatch02Slug(c.Spell)) || c.ViaSpellcasting != has("effect/spellcasting-use") ||
			(c.ViaSpellcasting && !s.lists(c.Spell, c.Level)) || (!c.ViaSpellcasting && d.Spellcasting == nil) {
			return fmt.Errorf("cast %q is not backed by the block", c.Spell)
		}
	}
	if sc := d.Spellcasting; sc != nil {
		if sc.Ability != "cha" || !has("effect/charisma-spellcasting-ability") || !has("effect/spell-save-dc") || !has("effect/no-material-components") ||
			sc.SaveDC != 8+m.Abilities["cha"].Modifier+m.ProficiencyBonus || sc.Spells == nil ||
			(sc.AttackBonus != nil) != has("effect/spell-attack-bonus") || (sc.AttackBonus != nil && *sc.AttackBonus != sc.SaveDC-8) ||
			(a.ID == "spellcasting") != (len(sc.Spells) > 0) {
			return fmt.Errorf("invalid spellcasting")
		}
		for i, sp := range sc.Spells {
			if strings.TrimSpace(sp.Name) == "" || sp.Level < 0 || sp.Level > 9 || (sp.Frequency != "at-will" && sp.Frequency != "1/day") ||
				!has("effect/spell-"+monsterBatch02Slug(sp.Name)) || (sp.Restriction != "" && !strings.Contains(a.Text, sp.Restriction)) ||
				slices.ContainsFunc(sc.Spells[:i], func(o MonsterBatch02Spell) bool { return o.Name == sp.Name }) {
				return fmt.Errorf("invalid spell %q", sp.Name)
			}
		}
	} else if a.ID == "spellcasting" {
		return fmt.Errorf("Spellcasting needs its typed list")
	}
	return s.validateDetailText(a, d)
}

// validateDetailText requires every typed source number of a detail to be
// printed in the retained action text, so a typed value cannot drift from it.
func (s MonsterBatch02StatBlock) validateDetailText(a MonsterAction, d MonsterBatch02ActionDetail) error {
	var want []string
	if at := a.Attack; at != nil {
		want = append(want, fmt.Sprintf("Attack Roll: +%d", at.Bonus))
		if at.Reach > 0 {
			want = append(want, fmt.Sprintf("reach %d ft.", at.Reach))
		}
		if at.Range > 0 {
			want = append(want, fmt.Sprintf("range %d", at.Range))
		}
	}
	if r := d.Recharge; r != nil {
		if r.Min == r.Max {
			want = append(want, fmt.Sprintf("(Recharge %d)", r.Min))
		} else {
			want = append(want, fmt.Sprintf("(Recharge %d–%d)", r.Min, r.Max))
		}
	}
	if d.Save != nil {
		want = append(want, fmt.Sprintf("Saving Throw: DC %d", d.Save.DC))
	}
	if ar := d.Area; ar != nil {
		switch ar.Shape {
		case "line":
			want = append(want, fmt.Sprintf("%d-foot-long, %d-foot-wide Line", ar.Feet, ar.WidthFeet))
		case "cone":
			want = append(want, fmt.Sprintf("%d-foot Cone", ar.Feet))
		case "sphere":
			want = append(want, fmt.Sprintf("%d-foot-radius Sphere", ar.Feet), fmt.Sprintf("within %d feet", ar.RangeFeet))
		}
	}
	if d.TargetRangeFeet > 0 {
		want = append(want, fmt.Sprintf("within %d feet", d.TargetRangeFeet))
	}
	for _, roll := range d.Damage {
		dice := fmt.Sprintf("%dd%d", roll.Dice, roll.Sides)
		if roll.Modifier != 0 {
			dice += fmt.Sprintf(" + %d", roll.Modifier)
		}
		want = append(want, fmt.Sprintf("%d (%s) %s%s damage", roll.Average, dice, strings.ToUpper(roll.Type[:1]), roll.Type[1:]))
	}
	if d.PushFeet > 0 {
		want = append(want, fmt.Sprintf("pushed up to %d feet", d.PushFeet))
	}
	if d.DurationMinutes > 0 {
		want = append(want, fmt.Sprintf("%d minute", d.DurationMinutes))
	}
	for _, p := range d.Penalties {
		switch {
		case p.Flat > 0:
			want = append(want, fmt.Sprintf("−%d penalty", p.Flat))
		case p.Average > 0:
			want = append(want, fmt.Sprintf("%d (%dd%d)", p.Average, p.Dice, p.Sides))
		default:
			want = append(want, fmt.Sprintf("%dd%d", p.Dice, p.Sides))
		}
	}
	if sc := d.Spellcasting; sc != nil {
		want = append(want, fmt.Sprintf("spell save DC %d", sc.SaveDC))
		if sc.AttackBonus != nil {
			want = append(want, fmt.Sprintf("+%d to hit with spell attacks", *sc.AttackBonus))
		}
	}
	for _, phrase := range want {
		if !monsterBatch02Prints(a.Text, phrase) {
			return fmt.Errorf("typed value %q is not printed in the action text", phrase)
		}
	}
	return nil
}

// monsterBatch02Prints reports whether text contains phrase as whole numbers:
// no digit may adjoin the match, so "DC 2" is not found in "DC 22".
func monsterBatch02Prints(text, phrase string) bool {
	digit := func(b byte) bool { return b >= '0' && b <= '9' }
	for from := 0; ; {
		i := strings.Index(text[from:], phrase)
		if i < 0 {
			return false
		}
		start, end := from+i, from+i+len(phrase)
		if (start == 0 || !digit(text[start-1]) || !digit(phrase[0])) && (end == len(text) || !digit(text[end]) || !digit(phrase[len(phrase)-1])) {
			return true
		}
		from = start + 1
	}
}

// lists reports whether the block's Spellcasting action prints the spell at that level.
func (s MonsterBatch02StatBlock) lists(spell string, level int) bool {
	for _, d := range s.ActionDetails {
		if d.ActionID == "spellcasting" && d.Spellcasting != nil {
			return slices.ContainsFunc(d.Spellcasting.Spells, func(sp MonsterBatch02Spell) bool { return sp.Name == spell && sp.Level == level })
		}
	}
	return false
}

func monsterBatch02Slug(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "-")
}

// Available reports whole-block runtime support. No batch 02 mechanic is
// implemented, so every block is unavailable; the typed details grant nothing.
func (s MonsterBatch02StatBlock) Available() bool {
	return s.Validate() == nil && s.Monster.Available()
}
