package handlers

import (
	"encoding/json"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
)

// The explicit death-save and ally-aid actions save one receipt, which
// reaches clients as the same event.
func TestDeathSaveActionReceiptEvent(t *testing.T) {
	receipt := &srd521.ZeroHPReceipt{RulesVersion: "srd-5.2.1", Event: "death-save", Reason: "d20 9, failure", Die: 9, Outcome: "failure"}
	events := engineResolutionEventsToAPI([]engine.AppliedToolCall{{Tool: "srd_death_save", Result: mustJSON(t, map[string]any{"receipt": receipt})}})
	if len(events) != 1 || events[0].Type != "srd_death_save" || events[0].Outcome != receipt.Reason || string(mustJSON(t, events[0].Details["receipt"])) != string(mustJSON(t, receipt)) {
		t.Fatalf("events: %+v", events)
	}
}

// The receipts of a player turn that starts at 0 HP reach live responses and
// history as 0 HP events, the repair before the Death Saving Throw.
func TestEncounterZeroHPReceiptEvents(t *testing.T) {
	repair := &srd521.ZeroHPReceipt{RulesVersion: "srd-5.2.1", Event: srd521.ZeroHPUnconsciousRepair, Reason: "now Unconscious"}
	save := &srd521.ZeroHPReceipt{RulesVersion: "srd-5.2.1", Event: "death-save", Reason: "d20 12, success", Die: 12, Outcome: "success"}
	for name, test := range map[string]struct {
		payload map[string]any
		want    []*srd521.ZeroHPReceipt
	}{
		"repair and save": {map[string]any{"unconscious_repair": repair, "death_save": save}, []*srd521.ZeroHPReceipt{repair, save}},
		"save only":       {map[string]any{"death_save": save}, []*srd521.ZeroHPReceipt{save}},
		"neither":         {map[string]any{}, nil},
		"unknown rules":   {map[string]any{"unconscious_repair": &srd521.ZeroHPReceipt{Reason: "unversioned"}}, nil},
	} {
		t.Run(name, func(t *testing.T) {
			call := engine.AppliedToolCall{Tool: "srd_encounter", Result: mustJSON(t, test.payload)}
			live := engineResolutionEventsToAPI([]engine.AppliedToolCall{call})
			history := engineResolutionEventsToAPI(appliedToolCallsFromSessionLog(mustJSON(t, []engine.AppliedToolCall{call})))
			if len(live) != len(test.want) || string(mustJSON(t, live)) != string(mustJSON(t, history)) {
				t.Fatalf("live=%+v history=%+v", live, history)
			}
			for i, want := range test.want {
				got, err := json.Marshal(live[i].Details["receipt"])
				if err != nil {
					t.Fatal(err)
				}
				if live[i].Type != "srd_death_save" || live[i].Outcome != want.Reason || string(got) != string(mustJSON(t, want)) {
					t.Fatalf("event %d: %+v", i, live[i])
				}
			}
		})
	}
}
