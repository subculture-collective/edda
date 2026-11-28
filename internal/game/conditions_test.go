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
		"unconscious death save":          {saved, nil, ActionIntent{Kind: srd521.DeathSave}, ""},
		"dead end turn":                   {[]byte(`{"srd521_death_saves":{"failures":3,"dead":true}}`), nil, ActionIntent{Kind: combat.EncounterEndTurn}, "dead"},
		"dead death save":                 {[]byte(`{"srd521_death_saves":{"dead":true}}`), nil, ActionIntent{Kind: srd521.DeathSave}, "dead"},
		"unknown death save state":        {[]byte(`{"srd521_death_saves":{"failures":4}}`), nil, ActionIntent{Kind: combat.EncounterEndTurn}, "adjudication"},
	} {
		reason, err := conditionReason(tc.stats, 1, tc.encounter, tc.intent)
		if err != nil {
			reason = err.Error()
		}
		if (tc.want == "") != (reason == "") || !strings.Contains(reason, tc.want) {
			t.Errorf("%s: reason %q, want %q", name, reason, tc.want)
		}
	}
}

// An older build could leave 0 HP without Unconscious: outside an active
// encounter only the repairing Death Saving Throw stays open, and inside one
// only ending the turn and declining a Reaction.
func TestConditionReasonZeroHPWithoutUnconscious(t *testing.T) {
	encounter := func(c srd521.Conditions) *combat.Encounter {
		return &combat.Encounter{Status: "active", Actors: []combat.EncounterActor{{ID: uuid.New(), Player: true, Conditions: c}, {ID: uuid.New()}}}
	}
	active := encounter(srd521.Conditions{})
	for name, tc := range map[string]struct {
		encounter *combat.Encounter
		kind      string
		want      string
	}{
		"check":                {nil, srd521.AbilityCheck, "death-save"},
		"stand":                {nil, combat.EncounterStand, "death-save"},
		"death save":           {nil, srd521.DeathSave, ""},
		"encounter end turn":   {active, combat.EncounterEndTurn, ""},
		"encounter decline":    {active, combat.EncounterDecline, ""},
		"encounter attack":     {active, srd521.MartialAttack, "End your turn"},
		"encounter react":      {active, combat.EncounterReact, "End your turn"},
		"encounter move":       {active, combat.EncounterMove, "End your turn"},
		"encounter owns state": {encounter(srd521.Conditions{Prone: true, Incapacitated: true, Unconscious: true}), srd521.MartialAttack, "You are Unconscious"},
	} {
		reason, err := conditionReason([]byte(`{}`), 0, tc.encounter, ActionIntent{Kind: tc.kind})
		if err != nil || (tc.want == "") != (reason == "") || !strings.Contains(reason, tc.want) {
			t.Errorf("%s: reason %q %v, want %q", name, reason, err, tc.want)
		}
	}
}

func TestZeroHPTurnReason(t *testing.T) {
	unconscious := []byte(`{"srd521_conditions":{"prone":true,"unconscious":true}}`)
	active := func(c srd521.Conditions) *combat.Encounter {
		return &combat.Encounter{Status: "active", Actors: []combat.EncounterActor{{ID: uuid.New(), Player: true, Conditions: c}, {ID: uuid.New()}}}
	}
	for name, tc := range map[string]struct {
		stats     []byte
		hp        int
		encounter *combat.Encounter
		want      string
	}{
		"has HP":                {[]byte(`{}`), 5, nil, ""},
		"unconscious":           {unconscious, 0, nil, "/death-save"},
		"stable":                {[]byte(`{"srd521_conditions":{"prone":true,"unconscious":true},"srd521_death_saves":{"stable":true}}`), 0, nil, "Stable"},
		"dead":                  {[]byte(`{"srd521_death_saves":{"failures":3,"dead":true}}`), 0, nil, "dead"},
		"missing unconscious":   {[]byte(`{}`), 0, nil, "without the Unconscious condition"},
		"in an encounter":       {[]byte(`{}`), 0, active(srd521.Conditions{Unconscious: true, Incapacitated: true, Prone: true}), "End your turn"},
		"unconscious with HP":   {unconscious, 3, nil, "Unconscious and cannot act"},
		"unknown death saves":   {[]byte(`{"srd521_death_saves":{"failures":4}}`), 0, nil, "adjudication"},
		"completed encounter":   {unconscious, 0, &combat.Encounter{Status: "completed"}, "/death-save"},
		"encounter has HP, foe": {[]byte(`{}`), 8, active(srd521.Conditions{}), ""},
	} {
		reason := ZeroHPTurnReason(tc.stats, tc.hp, tc.encounter)
		if (tc.want == "") != (reason == "") || !strings.Contains(reason, tc.want) {
			t.Errorf("%s: reason %q, want %q", name, reason, tc.want)
		}
	}
}
