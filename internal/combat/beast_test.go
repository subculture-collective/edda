package combat

import (
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
)

func ptr[T any](v T) *T { return &v }

func beastProfile(t *testing.T, profile, action, size string) *srd521.BeastProfile {
	t.Helper()
	_, version, err := srd521.BeastRecord(profile)
	if err != nil {
		t.Fatal(err)
	}
	return &srd521.BeastProfile{Version: version, ActionID: action, TargetSize: size}
}

func TestPackAllyEligibility(t *testing.T) {
	for _, tc := range []struct {
		name     string
		ally     *PackAlly
		eligible bool
		unknown  bool
	}{
		{"undeclared", nil, false, true},
		{"present unknown", &PackAlly{}, false, true},
		{"no ally", &PackAlly{Present: ptr(false)}, false, false},
		{"adjacent", &PackAlly{Present: ptr(true), DistanceFeet: ptr(5), Incapacitated: ptr(false)}, true, false},
		{"incapacitated", &PackAlly{Present: ptr(true), DistanceFeet: ptr(5), Incapacitated: ptr(true)}, false, false},
		{"too far", &PackAlly{Present: ptr(true), DistanceFeet: ptr(10), Incapacitated: ptr(false)}, false, false},
		{"distance unknown", &PackAlly{Present: ptr(true), Incapacitated: ptr(false)}, false, true},
		{"condition unknown", &PackAlly{Present: ptr(true), DistanceFeet: ptr(5)}, false, true},
	} {
		ok, err := tc.ally.Eligible()
		if ok != tc.eligible || (err != nil) != tc.unknown {
			t.Fatalf("%s: %t %v", tc.name, ok, err)
		}
	}
}

func TestBeastGeometryAndAdvantage(t *testing.T) {
	adjacent := &PackAlly{Present: ptr(true), DistanceFeet: ptr(5), Incapacitated: ptr(false)}
	none := &PackAlly{Present: ptr(false)}
	for _, tc := range []struct {
		name, profile, action, size string
		distance                    int
		sight                       bool
		ally                        *PackAlly
		prone, legal                bool
		advantage                   int
	}{
		{"wolf pack", "wolf-v1", "bite", "medium", 5, true, adjacent, false, true, 1},
		{"wolf alone", "wolf-v1", "bite", "medium", 5, true, none, false, true, 0},
		{"wolf pack prone", "wolf-v1", "bite", "medium", 5, true, adjacent, true, true, 2},
		{"wolf ally unknown", "wolf-v1", "bite", "medium", 5, true, nil, false, false, 0},
		{"wolf size unknown", "wolf-v1", "bite", "", 5, true, none, false, false, 0},
		{"wolf out of reach", "wolf-v1", "bite", "medium", 10, true, none, false, false, 0},
		{"wolf unseen", "wolf-v1", "bite", "medium", 5, false, none, false, false, 0},
		{"bear ignores ally", "brown-bear-v1", "multiattack", "medium", 5, true, nil, false, true, 0},
		{"bear prone target", "brown-bear-v1", "multiattack", "medium", 5, true, nil, true, true, 1},
		{"bear single attack", "brown-bear-v1", "claw", "medium", 5, true, nil, false, false, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := beastProfile(t, tc.profile, tc.action, tc.size)
			space := &Space{DistanceFeet: &tc.distance, LineOfSight: &tc.sight, Cover: "none", PackAlly: tc.ally}
			e := &Encounter{Beast: p, Space: space, Status: "active", Actors: []EncounterActor{{ID: uuid.New(), Player: true, HP: 10, Prone: tc.prone}, {ID: uuid.New(), HP: 10}}}
			reasons, err := e.BeastAdvantage()
			if (err == nil) != tc.legal || (BeastAttack(p, space) == nil) != tc.legal || len(reasons) != tc.advantage {
				t.Fatalf("legal=%t err=%v reasons=%v", tc.legal, err, reasons)
			}
		})
	}
	if BeastAttack(nil, &Space{}) == nil {
		t.Fatal("missing profile accepted")
	}
}

func TestProneStandAndRefusals(t *testing.T) {
	player := uuid.New()
	e := &Encounter{Status: "active", Round: 1, Actors: []EncounterActor{{ID: player, Player: true, HP: 10, Prone: true}, {ID: uuid.New(), HP: 10}}, Budget: TurnBudget{Action: true, BonusAction: true, Movement: 30}, Space: &Space{DistanceFeet: ptr(5), LineOfSight: ptr(true), Cover: "none"}}
	for _, kind := range []string{srd521.MartialAttack, EncounterMove, EncounterApproach, EncounterRetreat} {
		if err := e.Spend(player, kind); err == nil {
			t.Fatalf("%s while Prone accepted", kind)
		}
	}
	if !e.Budget.Action || e.Budget.Movement != 30 {
		t.Fatal("refusal spent budget")
	}
	e.Budget.Movement = 10
	if err := e.Spend(player, EncounterStand); err == nil || !e.Actors[0].Prone {
		t.Fatal("stood without half speed")
	}
	e.Budget.Movement = 30
	if err := e.Spend(player, EncounterStand); err != nil || e.Actors[0].Prone || e.Budget.Movement != 15 {
		t.Fatalf("stand: %v %+v", err, e.Budget)
	}
	if err := e.Spend(player, EncounterStand); err == nil {
		t.Fatal("stood twice")
	}
	if err := e.Spend(player, srd521.MartialAttack); err != nil {
		t.Fatal(err)
	}
	// A pending opportunity attack is refused while the player is Prone.
	e.Actors[0].Prone = true
	e.Turn = 1
	e.Reaction = &ReactionState{Available: true, Pending: &ReactionTrigger{TargetID: e.Actors[1].ID}}
	if err := e.SpendReaction(); err == nil || !e.Reaction.Available {
		t.Fatal("Prone opportunity attack accepted")
	}
	if !e.Actors[0].Prone {
		t.Fatal("refusal changed condition")
	}
}
