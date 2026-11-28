package engine

import (
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
)

// Multiattack saves once per damaging attack, so the narrative names each save
// after the attack that caused it, as the resolution events do.
func TestNPCTurnNarrativeNamesEachSaveAfterItsAttack(t *testing.T) {
	attack := func(weapon string) *srd521.AttackReceipt {
		return &srd521.AttackReceipt{Weapon: weapon, Outcome: "hit", Die: 12, Total: 14, Damage: 3, DamageType: "slashing"}
	}
	attacks := []*srd521.AttackReceipt{attack("Bite"), attack("Claw"), attack("Claw")}
	saves := []*srd521.EffectReceipt{{Reason: "Concentration maintained after damage."}, {Reason: "Concentration lost after damage."}}
	// The second attack caused no save: the second save belongs to the third.
	got := npcTurnNarrative(attacks, saves, []int{1, 3})
	hit := ": hit (d20 12 + 0 = 14), 3 slashing damage."
	want := " Bite" + hit + " Concentration maintained after damage." +
		" Claw" + hit +
		" Claw" + hit + " Concentration lost after damage."
	if got != want {
		t.Fatalf("narrative %q, want %q", got, want)
	}
}
