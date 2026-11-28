package srd521

import "fmt"

// Mechanic keys for server-resolved events that can end Concentration
// (SRD 5.2.1 Rules Glossary, "Concentration"; adapted as fixture keys, not
// registered spell support). Any other key fails closed.
const (
	ConcentrationDamageSource  = "concentration.damage.per-source"
	ConcentrationIncapacitated = "concentration.end.incapacitated"
	ConcentrationUnconscious   = "concentration.end.unconscious"
	ConcentrationDead          = "concentration.end.dead"
)

// ConcentrationEvent is one trusted, already-resolved event on the caster.
// Callers build it from committed damage or condition results inside their
// own transaction; it is never decoded from client input.
type ConcentrationEvent struct {
	Kind string
	// SourceID names one damage source, such as one attack or one fall.
	// Consecutive damage events with the same SourceID are one instance.
	SourceID string
	// Damage is the actual damage taken; HPAfter is the caster's HP after it.
	Damage  int
	HPAfter int
}

// ResolveConcentrationEvents applies trusted events in order to the caster's
// detached build. Each positive damage source rolls one Constitution save at
// DC min(30,max(10,floor(damage/2))); zero HP, incapacitation, Unconscious and
// death end Concentration without a roll. Events after Concentration ends
// produce no receipt and no roll. Every event is validated before any change.
func ResolveConcentrationEvents(b *Build, sheet *CharacterSheet, events []ConcentrationEvent, now int, dice Dice) ([]*EffectReceipt, error) {
	if b == nil || sheet == nil || now < 0 {
		return nil, fmt.Errorf("caster and campaign clock required")
	}
	merged := make([]ConcentrationEvent, 0, len(events))
	seen := map[string]bool{}
	hp := sheet.HP.Value
	for _, e := range events {
		switch e.Kind {
		case ConcentrationDamageSource:
			if e.SourceID == "" || e.Damage < 0 || e.HPAfter < 0 {
				return nil, fmt.Errorf("damage event needs a source, actual damage and resulting HP")
			}
			if e.HPAfter > hp {
				return nil, fmt.Errorf("damage cannot raise HP")
			}
			hp = e.HPAfter
			if n := len(merged); n > 0 && merged[n-1].Kind == e.Kind && merged[n-1].SourceID == e.SourceID {
				merged[n-1].Damage += e.Damage
				merged[n-1].HPAfter = e.HPAfter
				continue
			}
			if seen[e.SourceID] {
				return nil, fmt.Errorf("damage source %q repeated after another event", e.SourceID)
			}
			seen[e.SourceID] = true
		case ConcentrationIncapacitated, ConcentrationUnconscious, ConcentrationDead:
			if e.Damage != 0 {
				return nil, fmt.Errorf("condition event carries no damage")
			}
		default:
			return nil, fmt.Errorf("unsupported concentration event %q", e.Kind)
		}
		merged = append(merged, e)
	}
	receipts := []*EffectReceipt{}
	hp = sheet.HP.Value
	for _, e := range merged {
		before := hp
		if e.Kind == ConcentrationDamageSource {
			hp = e.HPAfter
		}
		if b.Casting == nil || b.Casting.Concentration == nil || (e.Kind == ConcentrationDamageSource && e.Damage == 0) {
			continue
		}
		c := b.Casting
		r := &EffectReceipt{RulesVersion: "srd-5.2.1", Kind: "concentration-damage", Mechanic: e.Kind, Before: c.Concentration, TimeBefore: now, TimeAfter: now, HPBefore: before, HPAfter: hp, SlotsBefore: c.Slots, SlotsAfter: c.Slots, HitDiceBefore: b.HitDice, HitDiceAfter: b.HitDice, SourceIDs: []string{prefix + "rules-glossary/concentration"}}
		switch {
		case e.Kind == ConcentrationDamageSource && hp <= 0:
			c.Concentration = nil
			r.Reason = "Concentration ended at zero HP."
		case e.Kind == ConcentrationDamageSource:
			if sheet.Spellcasting == nil || dice == nil {
				return nil, fmt.Errorf("Concentration save needs a supported caster and dice")
			}
			saveSheet := *sheet
			saveSheet.HP.Value = before
			check, err := ResolveCheck(&saveSheet, Check{Kind: SavingThrow, Ability: "constitution"}, Difficulty{Value: min(30, max(10, e.Damage/2)), PlayerKnown: true, OutcomeKnown: true}, dice)
			if err != nil {
				return nil, err
			}
			receipt := check.PlayerReceipt()
			r.Check = &receipt
			r.Reason = "Concentration maintained after damage."
			if !check.Success() {
				c.Concentration = nil
				r.Reason = "Concentration lost after damage."
			}
		default:
			c.Concentration = nil
			r.Kind = "concentration-ended"
			r.Reason = map[string]string{ConcentrationIncapacitated: "Concentration ended: Incapacitated.", ConcentrationUnconscious: "Concentration ended: Unconscious.", ConcentrationDead: "Concentration ended at death."}[e.Kind]
			r.SourceIDs = append(r.SourceIDs, prefix+"rules-glossary/incapacitated")
			if e.Kind == ConcentrationUnconscious {
				r.SourceIDs = append(r.SourceIDs, prefix+"rules-glossary/unconscious")
			}
		}
		r.After = c.Concentration
		receipts = append(receipts, r)
	}
	return receipts, nil
}
