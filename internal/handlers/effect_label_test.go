package handlers

import (
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
)

// Every path that shows a caster effect receipt names its concentration spells,
// including receipts saved before the names were sent.
func TestResolutionEventsNameConcentrationSpells(t *testing.T) {
	receipt := &srd521.EffectReceipt{
		RulesVersion: "srd-5.2.1",
		Before:       &srd521.Concentration{ID: "faith-id", Spell: "shield-of-faith"},
		After:        &srd521.Concentration{ID: "detect-id", Spell: srd521.DetectMagicID},
	}
	for name, call := range map[string]struct {
		tool    string
		payload map[string]any
	}{
		"lifecycle": {"srd_effect", map[string]any{"receipt": receipt}},
		"expiry":    {"srd_attack", map[string]any{"concentration_receipts": []*srd521.EffectReceipt{receipt}}},
		"encounter": {"srd_encounter", map[string]any{"concentration_receipts": []*srd521.EffectReceipt{receipt}}},
		"fall":      {"srd_fall", map[string]any{"receipt": &srd521.FallReceipt{RulesVersion: "srd-5.2.1", Concentration: receipt}}},
	} {
		t.Run(name, func(t *testing.T) {
			var shown []*srd521.EffectReceipt
			for _, event := range engineResolutionEventsToAPI([]engine.AppliedToolCall{{Tool: call.tool, Result: mustJSON(t, call.payload)}}) {
				switch r := event.Details["receipt"].(type) {
				case *srd521.EffectReceipt:
					shown = append(shown, r)
				case *srd521.FallReceipt:
					shown = append(shown, r.Concentration)
				}
			}
			if len(shown) == 0 {
				t.Fatal("no caster effect receipt shown")
			}
			for _, r := range shown {
				if r.Before.SpellLabel != "Shield of Faith" || r.After.SpellLabel != "Detect Magic" {
					t.Fatalf("labels %q and %q", r.Before.SpellLabel, r.After.SpellLabel)
				}
			}
		})
	}
}
