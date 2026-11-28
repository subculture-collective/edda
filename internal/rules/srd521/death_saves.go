package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// DeathSave is the explicit action that rolls a player's Death Saving Throw
// outside encounters, where Edda keeps no turn order (SRD 5.2.1 PDF pp. 17-18;
// adapted, see Attribution). Inside an encounter the save is rolled at the
// start of the player's turn instead.
const DeathSave = "death-save"

// Zero-HP receipt events.
const (
	ZeroHPDropped        = "dropped-to-0"
	ZeroHPMassiveDamage  = "massive-damage"
	ZeroHPDamageAtZero   = "damage-at-0"
	ZeroHPDeathSave      = "death-save"
	ZeroHPStableRecovery = "stable-recovery"
)

// DeathSaves is a player character's durable 0-HP state, stats.srd521_death_saves.
// Successes and Failures count toward three; both reset when the character
// regains Hit Points or becomes Stable. Dead is final: no supported rule
// revives the character. Missing state is a living character with no saves.
type DeathSaves struct {
	Successes int  `json:"successes,omitempty"`
	Failures  int  `json:"failures,omitempty"`
	Stable    bool `json:"stable,omitempty"`
	Dead      bool `json:"dead,omitempty"`
}

// Validate refuses counts and combinations the rules cannot produce.
func (s DeathSaves) Validate() error {
	bad := s.Successes < 0 || s.Failures < 0 || s.Successes > 2
	if s.Dead {
		bad = bad || s.Failures > 3 || s.Stable
	} else {
		bad = bad || s.Failures > 2 || (s.Stable && (s.Successes > 0 || s.Failures > 0))
	}
	if !bad {
		return nil
	}
	return fmt.Errorf("this character's death saving throw state is not one the rules produce; it requires adjudication")
}

// DecodeDeathSaves reads stats.srd521_death_saves strictly, like DecodeConditions.
func DecodeDeathSaves(stats []byte) (DeathSaves, error) {
	var envelope struct {
		Saves json.RawMessage `json:"srd521_death_saves"`
	}
	var s DeathSaves
	if len(stats) == 0 {
		return s, nil
	}
	if err := json.Unmarshal(stats, &envelope); err != nil {
		return s, fmt.Errorf("character stats are unreadable")
	}
	if len(envelope.Saves) == 0 || string(envelope.Saves) == "null" {
		return s, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(envelope.Saves))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&s); err != nil {
		return DeathSaves{}, fmt.Errorf("this character's death saving throw state is not one the rules produce; it requires adjudication")
	}
	return s, s.Validate()
}

// ZeroHPReceipt records one 0-HP event for the player character: dropping to
// 0 HP, death by massive damage, damage taken at 0 HP, a Death Saving Throw or
// a Stable character's recovery. Die is the d20 (or the d4 hours of recovery).
type ZeroHPReceipt struct {
	RulesVersion  string            `json:"rules_version"`
	Event         string            `json:"event"`
	Reason        string            `json:"reason"`
	Damage        int               `json:"damage,omitempty"`
	Critical      bool              `json:"critical,omitempty"`
	Remaining     int               `json:"remaining,omitempty"`
	Die           int               `json:"die,omitempty"`
	Outcome       string            `json:"outcome,omitempty"`
	HPBefore      int               `json:"hp_before"`
	HPAfter       int               `json:"hp_after"`
	MaxHP         int               `json:"max_hp"`
	Before        DeathSaves        `json:"before"`
	After         DeathSaves        `json:"after"`
	TimeBefore    int               `json:"time_before,omitempty"`
	TimeAfter     int               `json:"time_after,omitempty"`
	Condition     *ConditionReceipt `json:"condition,omitempty"`
	SourceIDs     []string          `json:"source_ids"`
	RecoveryHours int               `json:"recovery_hours,omitempty"`
}

// Regained reports a change that brings the character back from 0 HP: it
// ends Unconscious and leaves Prone.
func (r *ZeroHPReceipt) Regained() bool { return r.HPBefore == 0 && r.HPAfter > 0 }

func (s DeathSaves) counts() string {
	return fmt.Sprintf("%d successes, %d failures", s.Successes, s.Failures)
}

// ResolvePlayerDamage applies damage to a player character's Hit Points and
// returns the HP after it and, when the damage reaches or lands at 0 HP, the
// receipt. Reaching 0 HP with remaining damage of at least maxHP is death
// (massive damage); otherwise the character drops to 0 HP and the caller
// applies Unconscious. Damage at 0 HP is one failure, two from a Critical
// Hit, death at three failures, and death when it is at least maxHP. Damage
// also ends Stable. A dead character takes no further effect.
func ResolvePlayerDamage(hp, maxHP, damage int, critical bool, saves DeathSaves) (int, *ZeroHPReceipt, error) {
	if hp < 0 || maxHP < 1 || hp > maxHP || damage < 0 {
		return 0, nil, fmt.Errorf("invalid damage state")
	}
	if err := saves.Validate(); err != nil {
		return 0, nil, err
	}
	if saves.Dead {
		return 0, nil, fmt.Errorf("the character is dead")
	}
	if damage == 0 || hp-damage > 0 {
		return max(0, hp-damage), nil, nil
	}
	r := &ZeroHPReceipt{RulesVersion: "srd-5.2.1", Damage: damage, Critical: critical, HPBefore: hp, MaxHP: maxHP, Before: saves}
	if hp > 0 {
		r.Remaining = damage - hp
		r.SourceIDs = []string{prefix + "damage/dropping-to-0-hit-points"}
		if r.Remaining >= maxHP {
			r.Event, r.After = ZeroHPMassiveDamage, DeathSaves{Dead: true}
			r.SourceIDs = append(r.SourceIDs, prefix+"damage/instant-death")
			r.Reason = fmt.Sprintf("Dropped to 0 HP with %d damage remaining, at least the %d HP maximum: the character dies.", r.Remaining, maxHP)
			return 0, r, nil
		}
		r.Event, r.After = ZeroHPDropped, DeathSaves{}
		r.SourceIDs = append(r.SourceIDs, prefix+"damage/falling-unconscious", prefix+"rules-glossary/unconscious")
		r.Reason = "Dropped to 0 HP: Unconscious and Prone, making Death Saving Throws at the start of each turn."
		return 0, r, nil
	}
	r.Event = ZeroHPDamageAtZero
	r.SourceIDs = []string{prefix + "damage/death-saving-throws"}
	if saves.Stable {
		r.SourceIDs = append(r.SourceIDs, prefix+"damage/stabilizing-a-character")
	}
	if damage >= maxHP {
		r.After = DeathSaves{Successes: saves.Successes, Failures: saves.Failures, Dead: true}
		r.Outcome = "dead"
		r.Reason = fmt.Sprintf("%d damage at 0 HP is at least the %d HP maximum: the character dies.", damage, maxHP)
		return 0, r, nil
	}
	failures := 1
	if critical {
		failures = 2
	}
	r.After = DeathSaves{Successes: saves.Successes, Failures: saves.Failures + failures}
	r.Outcome = "failure"
	if failures == 2 {
		r.Outcome = "two-failures"
	}
	r.Reason = fmt.Sprintf("Damage at 0 HP: %d Death Saving Throw failure", failures)
	if failures == 2 {
		r.Reason = "Damage from a Critical Hit at 0 HP: 2 Death Saving Throw failures"
	}
	if saves.Stable {
		r.Reason += "; no longer Stable"
	}
	if r.After.Failures >= 3 {
		r.After.Failures, r.After.Dead, r.Outcome = 3, true, "dead"
		r.Reason += ". Three failures: the character dies."
		return 0, r, nil
	}
	r.Reason += " (" + r.After.counts() + ")."
	return 0, r, nil
}

// ResolveDeathSave rolls one Death Saving Throw for a character at 0 HP that
// is neither Stable nor dead: a d20 of 10 or higher succeeds, a 1 is two
// failures and a 20 regains 1 Hit Point. Three successes make the character
// Stable and three failures kill it.
func ResolveDeathSave(saves DeathSaves, maxHP int, dice Dice) (*ZeroHPReceipt, error) {
	if err := saves.Validate(); err != nil {
		return nil, err
	}
	if saves.Dead || saves.Stable || maxHP < 1 || dice == nil {
		return nil, fmt.Errorf("no Death Saving Throw is due")
	}
	roll := dice.Intn(20)
	if roll < 0 || roll >= 20 {
		return nil, fmt.Errorf("invalid d20")
	}
	r := &ZeroHPReceipt{RulesVersion: "srd-5.2.1", Event: ZeroHPDeathSave, Die: roll + 1, MaxHP: maxHP, Before: saves, After: saves, SourceIDs: []string{prefix + "damage/death-saving-throws", prefix + "rules-glossary/death-saving-throw"}}
	switch {
	case r.Die == 20:
		r.Outcome, r.HPAfter, r.After = "regain-1-hp", 1, DeathSaves{}
		r.Reason = "Death Saving Throw: natural 20. Regains 1 HP; Unconscious ends and Prone remains."
		return r, nil
	case r.Die == 1:
		r.Outcome = "two-failures"
		r.After.Failures += 2
	case r.Die < 10:
		r.Outcome = "failure"
		r.After.Failures++
	default:
		r.Outcome = "success"
		r.After.Successes++
	}
	r.Reason = fmt.Sprintf("Death Saving Throw: d20 %d, %s", r.Die, map[string]string{"two-failures": "a 1 counts as two failures", "failure": "failure", "success": "success"}[r.Outcome])
	switch {
	case r.After.Failures >= 3:
		r.After = DeathSaves{Successes: saves.Successes, Failures: 3, Dead: true}
		r.Outcome = "dead"
		r.Reason += ". Three failures: the character dies."
	case r.After.Successes >= 3:
		r.After = DeathSaves{Stable: true}
		r.Outcome = "stable"
		r.SourceIDs = append(r.SourceIDs, prefix+"damage/stabilizing-a-character", prefix+"rules-glossary/stable")
		r.Reason += ". Three successes: Stable at 0 HP, still Unconscious; no more Death Saving Throws."
	default:
		r.Reason += " (" + r.After.counts() + ")."
	}
	return r, nil
}

// ResolveStableRecovery is a Stable character that is not healed regaining
// 1 Hit Point after 1d4 hours. The caller advances the campaign clock by
// RecoveryHours and ends Unconscious.
func ResolveStableRecovery(saves DeathSaves, maxHP int, dice Dice) (*ZeroHPReceipt, error) {
	if err := saves.Validate(); err != nil {
		return nil, err
	}
	if !saves.Stable || saves.Dead || maxHP < 1 || dice == nil {
		return nil, fmt.Errorf("the character is not Stable")
	}
	roll := dice.Intn(4)
	if roll < 0 || roll >= 4 {
		return nil, fmt.Errorf("invalid d4")
	}
	return &ZeroHPReceipt{RulesVersion: "srd-5.2.1", Event: ZeroHPStableRecovery, Die: roll + 1, RecoveryHours: roll + 1, Outcome: "regain-1-hp", HPAfter: 1, MaxHP: maxHP, Before: saves, SourceIDs: []string{prefix + "damage/stabilizing-a-character"}, Reason: fmt.Sprintf("Stable and not healed: after 1d4 = %d hours, regains 1 HP. Unconscious ends and Prone remains.", roll+1)}, nil
}
