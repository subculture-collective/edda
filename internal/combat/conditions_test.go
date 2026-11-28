package combat

import (
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
)

type queueEncounterDice []int

func (d *queueEncounterDice) Intn(int) int { v := (*d)[0]; *d = (*d)[1:]; return v }

func TestEncounterIncapacitatedInitiativeDisadvantage(t *testing.T) {
	// Player rolls 18 and 4 (raw 17 and 3); it keeps the lower. The foe rolls 11.
	dice := queueEncounterDice{17, 3, 10}
	e, err := NewEncounter(EncounterActor{ID: uuid.New(), HP: 10, Conditions: srd521.Conditions{Incapacitated: true}}, EncounterActor{ID: uuid.New(), HP: 10}, 0, &dice)
	if err != nil {
		t.Fatal(err)
	}
	if len(dice) != 0 {
		t.Fatalf("unused dice %v", dice)
	}
	if e.Actors[1].InitiativeDie != 4 || !e.Actors[1].Player || e.Actors[0].InitiativeDie != 11 {
		t.Fatalf("actors %+v", e.Actors)
	}
}

func TestEncounterConditionsGateTurnBudget(t *testing.T) {
	player := uuid.New()
	newEncounter := func(c srd521.Conditions) *Encounter {
		a := EncounterActor{ID: player, Player: true, HP: 10, Conditions: c}
		return &Encounter{Status: "active", Round: 1, Actors: []EncounterActor{a, {ID: uuid.New(), HP: 10}}, Budget: TurnBudget{Action: true, BonusAction: true, Movement: actorSpeed(a)}}
	}
	actions := []string{srd521.MartialAttack, srd521.SpellCast, srd521.ConsumableUse, srd521.ActionSurge, EncounterDash}
	inc := newEncounter(srd521.Conditions{Incapacitated: true})
	for _, kind := range actions {
		if err := inc.Spend(player, kind); err == nil {
			t.Errorf("Incapacitated %s accepted", kind)
		}
	}
	if inc.Budget.Movement != 30 || inc.Spend(player, EncounterMove) != nil {
		t.Fatalf("Incapacitated keeps its speed: %+v", inc.Budget)
	}
	unc := newEncounter(srd521.Conditions{Prone: true, Unconscious: true})
	if unc.Budget.Movement != 0 {
		t.Fatalf("Unconscious speed %d", unc.Budget.Movement)
	}
	for _, kind := range append(actions, EncounterMove, EncounterStand) {
		if err := unc.Spend(player, kind); err == nil {
			t.Errorf("Unconscious %s accepted", kind)
		}
	}
	if err := unc.Spend(player, EncounterEndTurn); err != nil {
		t.Fatal(err)
	}
	if !unc.Actors[0].Prone {
		t.Fatal("standing must not end Prone while Unconscious")
	}
}
