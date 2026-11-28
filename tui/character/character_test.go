package character

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"git.subcult.tv/subculture-collective/edda/tui/msgs"
)

// Compile-time checks: *Model must satisfy tea.Model and the sub-model
// interface (tea.Model + SetSize) expected by the Router. Importing the parent
// tui package here would create a cycle, so the interface is inlined.
var _ tea.Model = (*Model)(nil)
var _ interface {
	tea.Model
	SetSize(width, height int)
} = (*Model)(nil)

func TestNewReturnsModel(t *testing.T) {
	m := New()
	if m.width != 0 || m.height != 0 {
		t.Fatal("expected zero dimensions on a freshly created model")
	}
}

func TestSetSizeUpdatesDimensions(t *testing.T) {
	m := New()
	m.SetSize(80, 24)
	if m.width != 80 || m.height != 24 {
		t.Fatalf("expected 80x24, got %dx%d", m.width, m.height)
	}
}

func TestInitReturnsNil(t *testing.T) {
	m := New()
	if m.Init() != nil {
		t.Fatal("Init() should return nil")
	}
}

func TestEscEmitsNavigateBackMsg(t *testing.T) {
	m := New()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a command on Esc, got nil")
	}
	msg := cmd()
	if _, ok := msg.(msgs.NavigateBackMsg); !ok {
		t.Fatalf("expected NavigateBackMsg, got %T", msg)
	}
}

func TestOtherKeyReturnsNoCmd(t *testing.T) {
	m := New()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatal("expected no command for non-Esc keys")
	}
}

func TestViewContainsCharacterSheetHeader(t *testing.T) {
	m := New()
	m.SetSize(80, 24)
	got := m.View()
	if !strings.Contains(got, "Character Sheet") {
		t.Fatalf("expected View to contain %q, got:\n%s", "Character Sheet", got)
	}
}

func TestViewContainsEscHint(t *testing.T) {
	m := New()
	m.SetSize(80, 24)
	got := m.View()
	if !strings.Contains(got, "Esc") {
		t.Fatalf("expected View to contain Esc key hint, got:\n%s", got)
	}
}

// screen renders the sheet of c in a box tall enough for all of it.
func screen(c *game.CharacterRecord) string {
	m := New()
	m.SetSize(100, 400)
	m.Show(c)
	return m.View()
}

func assertShows(t *testing.T, view string, texts ...string) {
	t.Helper()
	for _, text := range texts {
		if !strings.Contains(view, text) {
			t.Errorf("the sheet does not show %q:\n%s", text, view)
		}
	}
}

func TestViewShowsASavedCharacterWithoutSRDRules(t *testing.T) {
	view := screen(&game.CharacterRecord{
		PlayerCharacter: domain.PlayerCharacter{
			Name: "Mara Venn", Description: "Ledger keeper. Owes the drowned.", Status: "active",
			Level: 2, HP: 7, MaxHP: 12, Experience: 1200,
			Stats:     []byte(`{"personality":"wry","bonds":["the ledger","a debt"],"luck_points":3,"secret":""}`),
			Abilities: []byte(`[{"name":"Second Sight","description":"Sees old debts"},{"name":"Plain"}]`),
		},
		RulesVersion: "legacy", Location: "The Drowned Gate",
	})
	assertShows(t, view, "Mara Venn", "Owes the drowned.", "Level 2 · Active · HP 7/12 · XP 1200/3000",
		"Location: The Drowned Gate", "Personality: wry", "Bonds: the ledger, a debt", "Luck points: 3", "Secret: —",
		"Second Sight: Sees old debts", "◆ Plain")
	if strings.Contains(view, "Skills") || strings.Contains(view, "No character yet") {
		t.Errorf("unexpected text:\n%s", view)
	}
}

func TestViewShowsTheServerDerivedSRDSheet(t *testing.T) {
	stats, err := os.ReadFile("../../internal/rules/srd521/testdata/fighter1.json")
	if err != nil {
		t.Fatal(err)
	}
	sheet, err := srd521.Project(stats, 1, 9, 12, 0)
	if err != nil {
		t.Fatal(err)
	}
	height := 30
	sheet.Fall = &srd521.PendingFall{HeightFeet: &height, Landing: "water"}
	sheet.DeathSaves = &srd521.DeathSaves{Successes: 1, Failures: 2}
	sheet.Prone = true
	sheet.Doors = []srd521.SheetDoor{{Name: "Iron door", State: "closed"}}
	sheet.TravelPaces = srd521.TravelPaces

	view := screen(&game.CharacterRecord{
		PlayerCharacter: domain.PlayerCharacter{Name: "Fixture fighter", Level: 1, HP: 9, MaxHP: 12, Status: "active", Stats: stats},
		RulesVersion:    "srd-5.2.1", Sheet: sheet,
	})
	strength := sheet.Abilities["strength"]
	assertShows(t, view, "Fixture fighter", "HP 9/12 · XP 0/300",
		fmt.Sprintf("Armor class %d", sheet.Defenses["armor-class"].Value),
		fmt.Sprintf("Initiative %+d", sheet.Defenses["initiative"].Value),
		fmt.Sprintf("Strength %d (%+d), save %+d", strength.Score, strength.Modifier.Value, sheet.Saves["strength"].Value),
		fmt.Sprintf("Athletics %+d", sheet.Skills["athletics"].Value),
		"Second Wind", "Weapon Mastery (read only)",
		"At 0 HP: Unconscious", "1 successes, 2 failures", "A fall of 30 feet onto water", "Prone",
		"Iron door (closed)", "known route", "does not claim Full Rules support")
	// The raw build is not shown beside the derived values.
	if strings.Contains(view, "fighting_style") || strings.Contains(view, "Srd521") {
		t.Errorf("the raw SRD build is on screen:\n%s", view)
	}
}

func TestViewShowsSpellcastingEffects(t *testing.T) {
	sheet := &srd521.CharacterSheet{
		Spellcasting: &srd521.Spellcasting{Ability: "wisdom", Modifier: srd521.Value{Value: 3}, Spells: []srd521.Spell{
			{ID: "light", Label: "Light", Known: true, Prepared: true, Playable: true},
			{ID: "guidance", Label: "Guidance", Known: true},
		}},
		Concentration: &srd521.Concentration{SpellLabel: "Shield of Faith", ExpiresAt: 29400},
		Light:         &srd521.LightEffect{Color: "blue", ExpiresAt: 32400},
		LightItem:     "Mace",
	}
	assertShows(t, screen(&game.CharacterRecord{RulesVersion: "srd-5.2.1", Sheet: sheet}),
		"Casting ability: Wisdom (+3)", "Light · Known · Prepared · Supported",
		"Guidance · Known · Not prepared · Effect unavailable",
		"Concentration: Shield of Faith, expires Day 1, 08:10:00", "Light on Mace (blue), expires Day 1, 09:00:00.")

	sheet.Light, sheet.Concentration = nil, nil
	assertShows(t, screen(&game.CharacterRecord{RulesVersion: "srd-5.2.1", Sheet: sheet}),
		"No concentration effect active.", "No Light active.")
}

// As on the web, what is wrong with a saved build is shown, and a failure
// inside the server is not.
func TestViewSaysWhenTheRulesSheetIsUnavailable(t *testing.T) {
	for name, tc := range map[string]struct {
		err    error
		shown  []string
		hidden []string
	}{
		"unsupported build": {
			err:   errors.New("six explicit ability scores required"),
			shown: []string{"Rules sheet unavailable: six explicit ability scores required"},
		},
		"failed load": {
			err:    fmt.Errorf("%w: %w", game.ErrSheetLoadFailed, errors.New(`relation "secret_table" does not exist`)),
			shown:  []string{"The character sheet could not be loaded."},
			hidden: []string{"secret_table", "unavailable:"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			view := screen(&game.CharacterRecord{
				PlayerCharacter: domain.PlayerCharacter{Name: "Odd build", Level: 1, Stats: []byte(`{"srd521":{"class":"wizard"}}`)},
				RulesVersion:    "srd-5.2.1", SheetError: tc.err,
			})
			assertShows(t, view, append(tc.shown, "Odd build")...)
			for _, unwanted := range append(tc.hidden, "wizard", "XP 0/") {
				if strings.Contains(view, unwanted) {
					t.Errorf("the sheet shows %q:\n%s", unwanted, view)
				}
			}
		})
	}
}
