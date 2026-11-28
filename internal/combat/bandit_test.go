package combat

import (
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"testing"
)

func TestBanditGeometry(t *testing.T) {
	_, version, err := srd521.BanditRecord()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, action, cover string
		distance            int
		sight, legal        bool
	}{
		{"melee", "scimitar", "none", 5, true, true},
		{"melee far", "scimitar", "none", 10, true, false},
		{"normal ranged", "light-crossbow", "none", 80, true, true},
		{"ranged half cover", "light-crossbow", "half", 30, true, true},
		{"ranged near", "light-crossbow", "none", 5, true, false},
		{"long range", "light-crossbow", "none", 85, true, false},
		{"total cover", "light-crossbow", "total", 30, true, false},
		{"hidden", "scimitar", "none", 5, false, false},
		{"unknown cover", "scimitar", "", 5, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := srd521.BanditProfile{Version: version, ActionID: tc.action}
			s := Space{DistanceFeet: &tc.distance, LineOfSight: &tc.sight, Cover: tc.cover}
			if err := BanditAttack(&p, &s); (err == nil) != tc.legal {
				t.Fatalf("legal=%t err=%v", tc.legal, err)
			}
		})
	}
	if BanditAttack(nil, nil) == nil {
		t.Fatal("missing geometry accepted")
	}
	p := srd521.BanditProfile{Version: version, ActionID: "scimitar"}
	if BanditAttack(&p, &Space{}) == nil {
		t.Fatal("unknown distance accepted")
	}
}

func TestBanditTurn(t *testing.T) {
	_, version, err := srd521.BanditRecord()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, configured, cover, weapon string
		distance                        *int
		sight, unknown                  bool
		closed, after                   int
	}{
		{"scimitar in reach", "scimitar", "none", "scimitar", ptr(5), true, false, 0, 5},
		{"crossbow within 5 feet", "light-crossbow", "none", "scimitar", ptr(5), true, false, 0, 5},
		{"out of scimitar reach", "scimitar", "none", "light-crossbow", ptr(15), true, false, 0, 15},
		{"normal range edge", "scimitar", "none", "light-crossbow", ptr(80), true, false, 0, 80},
		{"walks into normal range", "light-crossbow", "none", "light-crossbow", ptr(110), true, false, 30, 80},
		{"long range after walking", "light-crossbow", "none", "", ptr(115), true, false, 30, 85},
		{"beyond long range", "light-crossbow", "none", "", ptr(400), true, false, 30, 370},
		{"blocked path", "light-crossbow", "half", "", ptr(100), true, false, 0, 100},
		{"half cover in range", "scimitar", "half", "light-crossbow", ptr(30), true, false, 0, 30},
		{"total cover", "scimitar", "total", "", ptr(30), true, false, 0, 30},
		{"unseen in reach", "scimitar", "none", "", ptr(5), false, false, 0, 5},
		{"distance unknown", "scimitar", "none", "", nil, true, true, 0, 0},
		{"cover unknown", "scimitar", "", "", ptr(5), true, true, 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			player, bandit := uuid.New(), uuid.New()
			e := &Encounter{Status: "active", Turn: 1, Bandit: &srd521.BanditProfile{Version: version, ActionID: tc.configured}, Actors: []EncounterActor{{ID: player, Player: true, HP: 12}, {ID: bandit, HP: 11, Speed: 30}}, Budget: TurnBudget{Action: true, Movement: 30}, Space: &Space{DistanceFeet: tc.distance, LineOfSight: ptr(tc.sight), Cover: tc.cover}}
			closed, profile, err := e.BanditTurn()
			weapon := ""
			if profile != nil {
				weapon = profile.ActionID
			}
			if (err != nil) != tc.unknown || closed != tc.closed || weapon != tc.weapon || e.Bandit.ActionID != tc.configured || (!tc.unknown && *e.Space.DistanceFeet != tc.after) {
				t.Fatalf("closed=%d weapon=%q err=%v space=%+v", closed, weapon, err, e.Space)
			}
		})
	}
}
