package handlers

import (
	"slices"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
)

// Clients show resolution events in list order, so the list follows the order
// the engine resolved them: each Concentration save after the attack that
// caused it, then the new turn's 0 HP receipts, then an expiry from the clock.
func TestResolutionEventsFollowResolutionOrder(t *testing.T) {
	attack := func(weapon string) *srd521.AttackReceipt {
		return &srd521.AttackReceipt{RulesVersion: "srd-5.2.1", Weapon: weapon, Outcome: "hit"}
	}
	effect := func(reason string) *srd521.EffectReceipt {
		return &srd521.EffectReceipt{RulesVersion: "srd-5.2.1", Reason: reason}
	}
	save := &srd521.ZeroHPReceipt{RulesVersion: "srd-5.2.1", Event: "death-save", Reason: "death save"}
	for name, test := range map[string]struct {
		tool    string
		payload map[string]any
		want    []string
	}{
		"saves follow their attacks": {"srd_encounter", map[string]any{
			"npc_receipts":               []*srd521.AttackReceipt{attack("Bite"), attack("Claw")},
			"concentration_receipts":     []*srd521.EffectReceipt{effect("after Bite"), effect("after Claw"), effect("expired")},
			"concentration_after_attack": []int{1, 2},
			"death_save":                 save,
		}, []string{"NPC Bite hit", "after Bite", "NPC Claw hit", "after Claw", "0 HP death save", "expired"}},
		// A record saved before the order was kept lists saves after attacks.
		"unordered record": {"srd_encounter", map[string]any{
			"npc_receipts":           []*srd521.AttackReceipt{attack("Bite")},
			"concentration_receipts": []*srd521.EffectReceipt{effect("saved")},
		}, []string{"NPC Bite hit", "saved"}},
		// A reaction can end the round: the expiry follows the player's attack.
		"expiry after the action": {"srd_attack", map[string]any{
			"receipt":                &srd521.AttackReceipt{RulesVersion: "srd-5.2.1", Kind: srd521.MartialAttack, Outcome: "miss"},
			"concentration_receipts": []*srd521.EffectReceipt{effect("expired")},
		}, []string{"Longsword attack miss", "expired"}},
		"lifecycle receipt first": {"srd_effect", map[string]any{
			"receipt":                effect("rested"),
			"concentration_receipts": []*srd521.EffectReceipt{effect("expired")},
		}, []string{"rested", "expired"}},
	} {
		t.Run(name, func(t *testing.T) {
			var got []string
			for _, event := range engineResolutionEventsToAPI([]engine.AppliedToolCall{{Tool: test.tool, Result: mustJSON(t, test.payload)}}) {
				if event.Type == "srd_effect" {
					got = append(got, event.Outcome)
				} else {
					got = append(got, event.Label+" "+event.Outcome)
				}
			}
			if !slices.Equal(got, test.want) {
				t.Fatalf("events %q, want %q", got, test.want)
			}
		})
	}
}
