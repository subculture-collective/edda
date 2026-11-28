package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Condition names for the shared Incapacitated and Unconscious state (SRD
// 5.2.1 Rules Glossary pp. 184 and 191; adapted, see Attribution). Any other
// name fails closed.
const (
	ConditionIncapacitated = "incapacitated"
	ConditionUnconscious   = "unconscious"
)

// Conditions is an actor's durable condition state. It is saved in
// stats.srd521_conditions outside encounters and on the player actor inside
// one. Each condition has exactly one flag: Unconscious implies Incapacitated
// and keeps Prone set, so neither is copied into another field.
// Incapacitated records only an Incapacitated source other than Unconscious.
// Class features such as Rage and Danger Sense read IsIncapacitated and
// Unconscious; they never keep their own copies.
type Conditions struct {
	Prone         bool `json:"prone,omitempty"`
	Incapacitated bool `json:"incapacitated,omitempty"`
	Unconscious   bool `json:"unconscious,omitempty"`
}

// IsIncapacitated reports the Incapacitated condition from any supported source.
func (c Conditions) IsIncapacitated() bool { return c.Incapacitated || c.Unconscious }

// Validate refuses a state the supported rules cannot produce.
func (c Conditions) Validate() error {
	if c.Unconscious && !c.Prone {
		return fmt.Errorf("an Unconscious actor is Prone; this condition state requires adjudication")
	}
	return nil
}

// DecodeConditions reads stats.srd521_conditions strictly: an unknown
// condition or a state Validate refuses returns an error, so callers fail
// closed instead of ignoring it. Missing state is no condition.
func DecodeConditions(stats []byte) (Conditions, error) {
	var envelope struct {
		Conditions json.RawMessage `json:"srd521_conditions"`
	}
	var c Conditions
	if len(stats) == 0 {
		return c, nil
	}
	if err := json.Unmarshal(stats, &envelope); err != nil {
		return c, fmt.Errorf("character stats are unreadable")
	}
	if len(envelope.Conditions) == 0 || string(envelope.Conditions) == "null" {
		return c, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(envelope.Conditions))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&c); err != nil {
		return Conditions{}, fmt.Errorf("this character has conditions the supported rules do not cover; actions require adjudication")
	}
	return c, c.Validate()
}

// ConditionEffects lists what the current conditions do. It is the one place
// those rules are written down; previews, encounters and later features read it.
type ConditionEffects struct {
	// Incapacitated: no actions, Bonus Actions or Reactions, no speech,
	// Disadvantage on Initiative, and Concentration ends.
	NoActions              bool `json:"no_actions,omitempty"`
	NoSpeech               bool `json:"no_speech,omitempty"`
	InitiativeDisadvantage bool `json:"initiative_disadvantage,omitempty"`
	// Unconscious also: held objects drop, Speed 0, attacks against it have
	// Advantage, Strength and Dexterity saves fail, a hit from within 5 feet
	// is a Critical Hit, and it is unaware of its surroundings.
	DropsHeldObjects   bool     `json:"drops_held_objects,omitempty"`
	SpeedZero          bool     `json:"speed_zero,omitempty"`
	AttackedAdvantage  bool     `json:"attacked_with_advantage,omitempty"`
	AutoFailSaves      []string `json:"auto_fail_saves,omitempty"`
	CriticalWithinFeet int      `json:"critical_within_feet,omitempty"`
	Unaware            bool     `json:"unaware,omitempty"`
	SourceIDs          []string `json:"source_ids,omitempty"`
}

// Effects returns the effects of c.
func (c Conditions) Effects() ConditionEffects {
	var e ConditionEffects
	if c.IsIncapacitated() {
		e.NoActions, e.NoSpeech, e.InitiativeDisadvantage = true, true, true
		e.SourceIDs = append(e.SourceIDs, prefix+"rules-glossary/incapacitated")
	}
	if c.Unconscious {
		e.DropsHeldObjects, e.SpeedZero, e.AttackedAdvantage, e.Unaware = true, true, true, true
		e.AutoFailSaves = []string{"strength", "dexterity"}
		e.CriticalWithinFeet = 5
		e.SourceIDs = append(e.SourceIDs, prefix+"rules-glossary/unconscious", prefix+"rules-glossary/prone")
	}
	return e
}

// AutoFailsSave reports whether a saving throw with ability fails automatically.
func (c Conditions) AutoFailsSave(ability string) bool {
	for _, a := range c.Effects().AutoFailSaves {
		if a == ability {
			return true
		}
	}
	return false
}

// UnavailableReason names why an incapacitated actor cannot take a choice, or
// returns "" when c does not restrict it. move reports a choice that only
// spends movement; it stays open while Incapacitated but not at Speed 0.
func (c Conditions) UnavailableReason(move bool) string {
	switch {
	case c.Unconscious:
		return "You are Unconscious: you can't act, move or speak."
	case c.Incapacitated && !move:
		return "You are Incapacitated: you can't take actions, Bonus Actions or Reactions, or speak."
	}
	return ""
}

// ConditionChange applies or ends one condition. Callers build it from a
// server-resolved result; it is never decoded from client input.
type ConditionChange struct {
	Condition string
	Apply     bool
}

// ConditionReceipt records one condition change. Dropped lists the item IDs
// the actor dropped; Concentration lists any cast the change ended.
type ConditionReceipt struct {
	RulesVersion  string           `json:"rules_version"`
	Condition     string           `json:"condition"`
	Applied       bool             `json:"applied"`
	Changed       bool             `json:"changed"`
	Before        Conditions       `json:"before"`
	After         Conditions       `json:"after"`
	Dropped       []string         `json:"dropped,omitempty"`
	Concentration []*EffectReceipt `json:"concentration,omitempty"`
	SourceIDs     []string         `json:"source_ids"`
}

// ResolveConditionChange returns the receipt for change and the Concentration
// events it raises. Applying Unconscious also applies Prone; ending it leaves
// Prone. Ending Unconscious leaves an Incapacitated flag from another source.
// A change that is already true is a no-op receipt with no events. Unknown
// conditions and invalid states fail closed.
func ResolveConditionChange(before Conditions, change ConditionChange) (*ConditionReceipt, []ConcentrationEvent, error) {
	if err := before.Validate(); err != nil {
		return nil, nil, err
	}
	after := before
	var source string
	switch change.Condition {
	case ConditionIncapacitated:
		after.Incapacitated = change.Apply
		source = prefix + "rules-glossary/incapacitated"
	case ConditionUnconscious:
		after.Unconscious = change.Apply
		if change.Apply {
			after.Prone = true
		}
		source = prefix + "rules-glossary/unconscious"
	default:
		return nil, nil, fmt.Errorf("unsupported condition %q", change.Condition)
	}
	r := &ConditionReceipt{RulesVersion: "srd-5.2.1", Condition: change.Condition, Applied: change.Apply, Changed: after != before, Before: before, After: after, SourceIDs: []string{source}}
	if change.Condition == ConditionUnconscious && change.Apply && !before.Prone {
		r.SourceIDs = append(r.SourceIDs, prefix+"rules-glossary/prone")
	}
	var events []ConcentrationEvent
	if !before.IsIncapacitated() && after.IsIncapacitated() {
		kind := ConcentrationIncapacitated
		if after.Unconscious {
			kind = ConcentrationUnconscious
		}
		events = []ConcentrationEvent{{Kind: kind}}
		r.SourceIDs = append(r.SourceIDs, prefix+"rules-glossary/concentration")
	}
	return r, events, nil
}
