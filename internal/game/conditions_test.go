package game

import (
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
)

func TestConditionReason(t *testing.T) {
	saved := []byte(`{"srd521_conditions":{"prone":true,"unconscious":true}}`)
	encounter := func(c srd521.Conditions, status string) *combat.Encounter {
		return &combat.Encounter{Status: status, Actors: []combat.EncounterActor{{ID: uuid.New(), Player: true, Conditions: c}, {ID: uuid.New()}}}
	}
	for name, tc := range map[string]struct {
		stats     []byte
		encounter *combat.Encounter
		intent    ActionIntent
		want      string
	}{
		"no conditions":                   {[]byte(`{}`), nil, ActionIntent{Kind: srd521.SpellCast}, ""},
		"unconscious check":               {saved, nil, ActionIntent{Kind: srd521.AbilityCheck}, "Unconscious"},
		"unconscious stand":               {saved, nil, ActionIntent{Kind: combat.EncounterStand}, "Unconscious"},
		"unconscious fall, no reaction":   {saved, nil, ActionIntent{Kind: srd521.FallResolve}, ""},
		"unconscious fall with reaction":  {saved, nil, ActionIntent{Kind: srd521.FallResolve, Skill: "athletics"}, "Unconscious"},
		"unconscious end turn":            {saved, nil, ActionIntent{Kind: combat.EncounterEndTurn}, ""},
		"encounter owns state":            {saved, encounter(srd521.Conditions{}, "active"), ActionIntent{Kind: srd521.MartialAttack}, ""},
		"encounter incapacitated attack":  {[]byte(`{}`), encounter(srd521.Conditions{Incapacitated: true}, "active"), ActionIntent{Kind: srd521.MartialAttack}, "Incapacitated"},
		"encounter incapacitated react":   {[]byte(`{}`), encounter(srd521.Conditions{Incapacitated: true}, "active"), ActionIntent{Kind: combat.EncounterReact}, "Incapacitated"},
		"encounter incapacitated move":    {[]byte(`{}`), encounter(srd521.Conditions{Incapacitated: true}, "active"), ActionIntent{Kind: combat.EncounterMove}, ""},
		"encounter incapacitated decline": {[]byte(`{}`), encounter(srd521.Conditions{Incapacitated: true}, "active"), ActionIntent{Kind: combat.EncounterDecline}, ""},
		"completed encounter uses stats":  {saved, encounter(srd521.Conditions{}, "completed"), ActionIntent{Kind: srd521.AbilityCheck}, "Unconscious"},
		"invalid encounter state":         {[]byte(`{}`), encounter(srd521.Conditions{Unconscious: true}, "active"), ActionIntent{Kind: combat.EncounterEndTurn}, "adjudication"},
		"unknown saved state":             {[]byte(`{"srd521_conditions":{"stunned":true}}`), nil, ActionIntent{Kind: combat.EncounterEndTurn}, "adjudication"},
	} {
		reason, err := conditionReason(tc.stats, tc.encounter, tc.intent)
		if err != nil {
			reason = err.Error()
		}
		if (tc.want == "") != (reason == "") || !strings.Contains(reason, tc.want) {
			t.Errorf("%s: reason %q, want %q", name, reason, tc.want)
		}
	}
}
