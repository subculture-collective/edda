package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

const (
	// DetectMagicID is the only feat spell entry this slice can cast.
	DetectMagicID = "detect-magic"
	// DetectMagicStudy is the Magic action that reveals auras while concentrating.
	DetectMagicStudy = "detect-magic-study"
	// DetectMagicFactsProfile names the trusted scene-fact contract.
	DetectMagicFactsProfile = "detect-magic-v1"
	// DetectMagicRadiusFeet is the source sensing radius (SRD 5.2.1 p. 123).
	DetectMagicRadiusFeet = 30
	detectMagicSeconds    = 600
)

// detectMagicBarriers are the three source blockers. "none" declares an open
// path; any other value is unknown and fails closed.
var detectMagicBarriers = []string{"stone-dirt-or-wood-1-foot", "metal-1-inch", "lead-thin-sheet"}

// DetectMagicFacts is an explicitly adjudicated, player-known magic bearer.
// Every field is required; nothing is inferred from prose or item names.
// Bearer says whether the npcs row is an "object" or a "creature"; sensing
// treats both alike, and the NPC listings leave objects out.
type DetectMagicFacts struct {
	Profile      string   `json:"profile"`
	Bearer       string   `json:"bearer"`
	ActorIDs     []string `json:"actor_ids"`
	LocationID   string   `json:"location_id"`
	Hidden       *bool    `json:"hidden"`
	DistanceFeet *int     `json:"distance_feet"`
	Barrier      string   `json:"barrier"`
	Visible      *bool    `json:"visible"`
	Origin       string   `json:"origin"`
	School       string   `json:"school"`
}

// spellNameMinorWords stay lowercase after the first word of a spell name.
var spellNameMinorWords = []string{"a", "an", "and", "for", "from", "into", "of", "on", "or", "the", "to", "via", "with", "without"}

var detectMagicFactKeys = []string{"profile", "bearer", "actor_ids", "location_id", "hidden", "distance_feet", "barrier", "visible", "origin", "school"}

// ParseDetectMagicFacts accepts only the exact vocabulary. Hidden lore, unknown
// bearers, unknown barriers, unknown origins and unknown schools are rejected, never guessed.
func ParseDetectMagicFacts(raw []byte) (DetectMagicFacts, bool) {
	var facts DetectMagicFacts
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil || len(fields) != len(detectMagicFactKeys) {
		return facts, false
	}
	for _, key := range detectMagicFactKeys {
		value := bytes.TrimSpace(fields[key])
		if len(value) == 0 || bytes.Equal(value, []byte("null")) {
			return facts, false
		}
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&facts) != nil || facts.Profile != DetectMagicFactsProfile || len(facts.ActorIDs) == 0 || facts.LocationID == "" || *facts.Hidden || *facts.DistanceFeet < 0 {
		return facts, false
	}
	if facts.Bearer != "object" && facts.Bearer != "creature" {
		return facts, false
	}
	if facts.Barrier != "none" && !slices.Contains(detectMagicBarriers, facts.Barrier) {
		return facts, false
	}
	schools := strings.Fields("abjuration conjuration divination enchantment evocation illusion necromancy transmutation")
	if (facts.Origin != "spell" || !slices.Contains(schools, facts.School)) && (facts.Origin != "other" || facts.School != "") {
		return facts, false
	}
	return facts, true
}

// DetectMagicSubject is one validated bearer, keyed by its scene identity.
type DetectMagicSubject struct {
	ID    string
	Label string
	Facts DetectMagicFacts
}

// DetectMagicAura is what a study discloses for one visible sensed subject.
type DetectMagicAura struct {
	SubjectID    string `json:"subject_id"`
	Subject      string `json:"subject"`
	DistanceFeet int    `json:"distance_feet"`
	School       string `json:"school,omitempty"`
	Origin       string `json:"origin"`
}

// DetectMagicReceipt is the saved, public account of a cast or a study.
type DetectMagicReceipt struct {
	RulesVersion   string            `json:"rules_version"`
	Kind           string            `json:"kind"`
	Spell          string            `json:"spell"`
	Use            string            `json:"use"`
	FeatUsesBefore int               `json:"magic_initiate_uses_before"`
	FeatUsesAfter  int               `json:"magic_initiate_uses_after"`
	SlotsBefore    int               `json:"slots_before"`
	SlotsAfter     int               `json:"slots_after"`
	Before         *Concentration    `json:"before,omitempty"`
	After          *Concentration    `json:"after,omitempty"`
	TimeBefore     int               `json:"time_before"`
	TimeAfter      int               `json:"time_after"`
	Presence       bool              `json:"presence"`
	Auras          []DetectMagicAura `json:"auras"`
	Reason         string            `json:"reason"`
	SourceIDs      []string          `json:"source_ids"`
}

// SensedDetectMagic keeps subjects within 30 feet with a declared open path.
func SensedDetectMagic(subjects []DetectMagicSubject) []DetectMagicSubject {
	out := []DetectMagicSubject{}
	for _, s := range subjects {
		if s.Facts.Barrier == "none" && *s.Facts.DistanceFeet <= DetectMagicRadiusFeet {
			out = append(out, s)
		}
	}
	return out
}

func detectMagicCaster(b *Build, sheet *CharacterSheet) error {
	if b == nil || b.Class != "cleric" || b.Casting == nil || sheet == nil || sheet.Spellcasting == nil || sheet.HP.Value <= 0 {
		return fmt.Errorf("conscious supported caster required")
	}
	return nil
}

// CastDetectMagic spends the recorded Magic Initiate (Cleric) free use, never a
// slot. It starts ten-minute concentration, replacing any current effect, and
// reports only whether adjudicated magic is sensed within range.
func CastDetectMagic(b *Build, sheet *CharacterSheet, subjects []DetectMagicSubject, actionID string, now int) (*DetectMagicReceipt, error) {
	if err := detectMagicCaster(b, sheet); err != nil {
		return nil, err
	}
	c := b.Casting
	ch := c.Choices
	if ch == nil || ch.FeatSpell != "Detect Magic" || ch.FeatAbility != "wisdom" {
		return nil, fmt.Errorf("no recorded Magic Initiate (Cleric) Detect Magic choice")
	}
	if ch.FeatUses != 1 {
		return nil, fmt.Errorf("Magic Initiate free use spent")
	}
	if !c.Verbal || !c.Somatic || actionID == "" || now < 0 {
		return nil, fmt.Errorf("verbal and somatic components and an action identity required")
	}
	r := &DetectMagicReceipt{RulesVersion: "srd-5.2.1", Kind: SpellCast, Spell: DetectMagicID, Use: "magic-initiate-free", FeatUsesBefore: 1, SlotsBefore: c.Slots, SlotsAfter: c.Slots, Before: c.Concentration, TimeBefore: now, TimeAfter: now, Auras: []DetectMagicAura{}}
	ch.FeatUses = 0
	c.Concentration = &Concentration{ID: actionID, Spell: DetectMagicID, StartedAt: now, ExpiresAt: now + detectMagicSeconds}
	r.After = c.Concentration
	r.Presence = len(SensedDetectMagic(subjects)) > 0
	r.Reason = "Detect Magic cast with the Magic Initiate (Cleric) free use; Magic action; no spell slot; concentration up to 10 minutes."
	if r.Before != nil {
		r.Reason += " Replaced " + ConcentrationLabel(r.Before.Spell) + " concentration."
	}
	if r.Presence {
		r.Reason += " Adjudicated magic is sensed within 30 feet; take the Magic action to study visible auras."
	} else {
		r.Reason += " No adjudicated magic is sensed within 30 feet; this does not establish that none is present."
	}
	r.SourceIDs = []string{prefix + "spells/detect-magic", prefix + "feats/magic-initiate", prefix + "rules-glossary/concentration"}
	return r, nil
}

// StudyDetectMagic is the source follow-up Magic action. It requires live
// Detect Magic concentration and sensed magic, and discloses only visible
// adjudicated subjects; schools appear only for spell-created effects.
func StudyDetectMagic(b *Build, sheet *CharacterSheet, subjects []DetectMagicSubject, now int) (*DetectMagicReceipt, error) {
	if err := detectMagicCaster(b, sheet); err != nil {
		return nil, err
	}
	c := b.Casting
	live := c.Concentration
	if live == nil || live.Spell != DetectMagicID || live.StartedAt > now || live.ExpiresAt <= now {
		return nil, fmt.Errorf("Detect Magic concentration required")
	}
	present := SensedDetectMagic(subjects)
	if len(present) == 0 {
		return nil, fmt.Errorf("no adjudicated magic sensed")
	}
	r := &DetectMagicReceipt{RulesVersion: "srd-5.2.1", Kind: DetectMagicStudy, Spell: DetectMagicID, SlotsBefore: c.Slots, SlotsAfter: c.Slots, Before: live, After: live, TimeBefore: now, TimeAfter: now, Presence: true, Auras: []DetectMagicAura{}}
	if c.Choices != nil {
		r.FeatUsesBefore, r.FeatUsesAfter = c.Choices.FeatUses, c.Choices.FeatUses
	}
	for _, s := range present {
		if *s.Facts.Visible {
			r.Auras = append(r.Auras, DetectMagicAura{SubjectID: s.ID, Subject: s.Label, DistanceFeet: *s.Facts.DistanceFeet, School: s.Facts.School, Origin: s.Facts.Origin})
		}
	}
	r.Reason = "Detect Magic study; Magic action; no slot or feat use."
	if len(r.Auras) == 0 {
		r.Reason += " No sensed magic is on a visible creature or object."
	}
	for _, a := range r.Auras {
		r.Reason += fmt.Sprintf(" Faint aura around %s (%d feet)", a.Subject, a.DistanceFeet)
		if a.School != "" {
			r.Reason += ": " + a.School + " spell."
		} else {
			r.Reason += "; not created by a spell."
		}
	}
	r.SourceIDs = []string{prefix + "spells/detect-magic", prefix + "rules-glossary/concentration"}
	return r, nil
}

// ConcentrationLabel names a concentration spell by turning its ID into title
// case. The spell records that hold SRD names are not built into the binary, so
// a name with punctuation, such as "Hunter’s Mark", loses it here.
func ConcentrationLabel(spell string) string {
	words := strings.FieldsFunc(spell, func(r rune) bool { return r == '-' })
	if len(words) == 0 {
		return "Unnamed spell"
	}
	for i, word := range words {
		if i == 0 || !slices.Contains(spellNameMinorWords, word) {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}
