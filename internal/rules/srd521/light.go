package srd521

import (
	"fmt"
	"slices"
)

// Light (SRD 5.2.1 p. 144) on one object the caster carries. The values below
// are pinned to the complete source record by light_test.go.
const (
	LightSpell           = "light"
	DismissLight         = "light-dismiss"
	LightDurationSeconds = 3600
	LightBrightFeet      = 20
	LightDimFeet         = 20
)

// LightColors is the explicit color vocabulary. The source lets the caster pick
// any color; a fixed list keeps free text out of saved rulings.
var LightColors = []string{"white", "yellow", "orange", "red", "green", "blue", "violet"}

// LightEffect is the caster's one Light cast. Absolute campaign seconds survive
// reloads and time advancement; the action ID identifies exactly which cast a
// later recast or dismissal removed.
type LightEffect struct {
	ID         string `json:"id"`
	ItemID     string `json:"item_id"`
	Color      string `json:"color"`
	StartedAt  int    `json:"started_at"`
	ExpiresAt  int    `json:"expires_at"`
	BrightFeet int    `json:"bright_feet"`
	DimFeet    int    `json:"dim_feet"`
}

// Validate rejects any saved effect whose identity, duration or radii differ
// from the supported slice.
func (l *LightEffect) Validate() error {
	if l == nil {
		return nil
	}
	if l.ID == "" || l.ItemID == "" || !slices.Contains(LightColors, l.Color) || l.StartedAt < 0 || l.ExpiresAt-l.StartedAt != LightDurationSeconds || l.BrightFeet != LightBrightFeet || l.DimFeet != LightDimFeet {
		return fmt.Errorf("invalid saved Light effect")
	}
	return nil
}

// Active reports whether the effect lasts at campaign second now.
func (l *LightEffect) Active(now int) bool {
	return l != nil && l.StartedAt <= now && now < l.ExpiresAt
}

// LightReceipt is the public, saved account of one cast or dismissal.
type LightReceipt struct {
	RulesVersion string       `json:"rules_version"`
	Kind         string       `json:"kind"`
	Reason       string       `json:"reason"`
	Item         string       `json:"item,omitempty"`
	Before       *LightEffect `json:"before,omitempty"`
	After        *LightEffect `json:"after,omitempty"`
	Time         int          `json:"time"`
	SourceIDs    []string     `json:"source_ids"`
}

// LightComponentsAvailable is the verbal and material ruling for the Magic
// Initiate cantrip: the recorded Holy Symbol focus stands in for the noncostly
// material component.
func LightComponentsAvailable(c *CastingBuild) bool {
	return c != nil && c.Verbal && c.Focus
}

func lightSpell(c *CastingBuild) Spell {
	return Spell{ID: LightSpell, Label: "Light", Known: true, Prepared: true, Components: []string{"verbal", "material (Holy Symbol focus)"}, ComponentsAvailable: LightComponentsAvailable(c), SlotLevels: []int{}, SourceID: prefix + "spells/light"}
}

// CastLight replaces the caster's previous Light, if any, with a new cast on
// one carried object. It leaves concentration, slots and Hit Dice unchanged.
func CastLight(b *Build, sheet *CharacterSheet, itemID, item, color, actionID string, now int) (*LightReceipt, error) {
	if b == nil || b.Class != "cleric" || b.Casting == nil || b.Casting.Choices == nil || !slices.Contains(b.Casting.Choices.FeatCantrips, "Light") || sheet == nil || sheet.HP.Value <= 0 {
		return nil, fmt.Errorf("conscious caster who knows Light required")
	}
	if !LightComponentsAvailable(b.Casting) || itemID == "" || actionID == "" || now < 0 || !slices.Contains(LightColors, color) {
		return nil, fmt.Errorf("Light components, object, color or clock unavailable")
	}
	r := &LightReceipt{RulesVersion: "srd-5.2.1", Kind: SpellCast, Item: item, Time: now, SourceIDs: []string{prefix + "spells/light", prefix + "feats/magic-initiate"}}
	if b.Casting.Light.Active(now) {
		r.Before = b.Casting.Light
	}
	b.Casting.Light = &LightEffect{ID: actionID, ItemID: itemID, Color: color, StartedAt: now, ExpiresAt: now + LightDurationSeconds, BrightFeet: LightBrightFeet, DimFeet: LightDimFeet}
	r.After = b.Casting.Light
	r.Reason = fmt.Sprintf("Light on %s (%s): bright light 20 feet and dim light 20 feet further for one hour; Magic action, no slot.", item, color)
	if r.Before != nil {
		r.Reason += " The previous Light ended."
	}
	return r, nil
}

// EndLight voluntarily dismisses the caster's active Light at no cost.
func EndLight(b *Build, item string, now int) (*LightReceipt, error) {
	if b == nil || b.Casting == nil || !b.Casting.Light.Active(now) {
		return nil, fmt.Errorf("no active Light to dismiss")
	}
	r := &LightReceipt{RulesVersion: "srd-5.2.1", Kind: DismissLight, Item: item, Before: b.Casting.Light, Time: now, SourceIDs: []string{prefix + "spells/light"}}
	b.Casting.Light = nil
	r.Reason = fmt.Sprintf("Light on %s dismissed; no action or slot cost.", item)
	return r, nil
}

// The bounded Light fixture: one object the caster carries, cast as the Magic
// Initiate cantrip. Objects carried by anyone else stay unavailable.
const lightFixtureScope = "bounded-carried-object-cantrip"

var lightFixtureTarget = SpellTarget{Shape: "single", Count: 1, Eligibility: "caster-carried-object-touch"}

func (s SpellRecord) lightFixtureAvailable() bool {
	if s.validate() != nil || s.SourceDetails != nil || len(s.Effects) != 1 {
		return false
	}
	m, e := s.Components.Material, s.Effects[0]
	if s.ID != prefix+"spells/light" || s.SourceID != s.ID || s.Name != "Light" || s.SourcePage != 144 || s.School != "evocation" || !slices.Equal(s.Lists, []string{"bard", "cleric", "sorcerer", "wizard"}) || s.CastingTime != "action" || s.Range != (SpellRange{"touch", 0}) {
		return false
	}
	if !*s.Components.Verbal || *s.Components.Somatic || m == nil || m.Description != "a firefly or phosphorescent moss" || *m.CostCP != 0 || *m.Consumed || *s.Ritual || *s.Concentration || s.Duration != "timed" || s.DurationSeconds != LightDurationSeconds {
		return false
	}
	if e.ID != "light" || e.Mechanic != "light.carried-object-bright-20-dim-20" || e.Timing != "for-duration" || e.DiceCount != 0 || e.DieSides != 0 || e.DamageType != "" || e.Bonus != 0 || *e.AddCastingModifier || *e.HalfInitialOnly {
		return false
	}
	return slices.Equal(s.RequiredMechanics, []string{"cast.magic-initiate-cantrip", "light.carried-object-bright-20-dim-20", "ending.light-recast-or-dismiss", "components.focus-for-material"})
}
