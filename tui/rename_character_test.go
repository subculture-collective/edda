package tui

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/config"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/tui/campaign"
)

// characterStore holds one campaign's characters and records renames.
type characterStore struct {
	noopQuerier
	characters []statedb.PlayerCharacter
	renames    []statedb.RenamePlayerCharacterParams
	listErr    error
	renameErr  error
}

func (s *characterStore) GetPlayerCharacterByCampaign(context.Context, pgtype.UUID) ([]statedb.PlayerCharacter, error) {
	return slices.Clone(s.characters), s.listErr
}

func (s *characterStore) RenamePlayerCharacter(_ context.Context, arg statedb.RenamePlayerCharacterParams) (int64, error) {
	if s.renameErr != nil {
		return 0, s.renameErr
	}
	s.renames = append(s.renames, arg)
	for i := range s.characters {
		if s.characters[i].ID == arg.ID {
			s.characters[i].Name = arg.Name
		}
	}
	return 1, nil
}

func namedCharacter(id byte, name string) statedb.PlayerCharacter {
	return statedb.PlayerCharacter{
		ID: pgtype.UUID{Bytes: [16]byte{id}, Valid: true}, Name: name,
		Description: pgtype.Text{String: "human barbarian. Acolyte", Valid: true},
	}
}

// sendUntilApp gives msg to m and runs what follows until the game opens. The
// game's own commands would wait on its refresh timer, so they are not run.
func sendUntilApp(t *testing.T, m tea.Model, msg tea.Msg) tea.Model {
	t.Helper()
	queue := []tea.Msg{msg}
	for steps := 0; len(queue) > 0; steps++ {
		if steps > 50 {
			t.Fatal("commands did not settle")
		}
		next, cmd := m.Update(queue[0])
		m, queue = next, queue[1:]
		if _, ok := m.(App); ok {
			return m
		}
		if cmd == nil {
			continue
		}
		queue = append(queue, results(cmd)...)
	}
	return m
}

// results runs cmd and returns the messages that move the launcher on.
func results(cmd tea.Cmd) []tea.Msg {
	switch msg := cmd().(type) {
	case tea.BatchMsg:
		var all []tea.Msg
		for _, c := range msg {
			if c != nil {
				all = append(all, results(c)...)
			}
		}
		return all
	case campaignLoadedMsg, campaign.CharacterRenamedMsg, characterRenameDoneMsg, campaign.BackMsg:
		return []tea.Msg{msg}
	default:
		return nil
	}
}

func renameLauncher(t *testing.T, store *characterStore) tea.Model {
	t.Helper()
	l := NewLauncherWithEngine(config.Config{}, t.Context(), store, nil)
	return resize(t, l, 80, 24)
}

func TestLoadingACampaignAsksToNameAnUnnamedCharacter(t *testing.T) {
	saved := makeTestCampaign(1, "The Drowned Crown")
	store := &characterStore{characters: []statedb.PlayerCharacter{namedCharacter(7, "")}}
	m := sendUntilApp(t, renameLauncher(t, store), campaign.SelectedMsg{Campaign: saved})

	view := m.View()
	for _, want := range []string{"Name Your Character", "human barbarian. Acolyte"} {
		if !strings.Contains(view, want) {
			t.Fatalf("%q is not on screen:\n%s", want, view)
		}
	}

	// A blank name stays on the prompt and saves nothing.
	m, _ = m.Update(typed("   "))
	m = sendUntilApp(t, m, key(tea.KeyEnter))
	if view := m.View(); !strings.Contains(view, "character name is required") || len(store.renames) != 0 {
		t.Fatalf("blank name: renames = %v, view:\n%s", store.renames, view)
	}

	m, _ = m.Update(typed("Mara "))
	m = sendUntilApp(t, m, key(tea.KeyEnter))
	want := statedb.RenamePlayerCharacterParams{ID: namedCharacter(7, "").ID, CampaignID: saved.ID, Name: "Mara"}
	if len(store.renames) != 1 || store.renames[0] != want {
		t.Fatalf("renames = %+v, want %+v", store.renames, want)
	}
	if _, ok := m.(App); !ok {
		t.Fatalf("the game did not open:\n%s", m.View())
	}
}

func TestLoadingANamedCampaignOpensTheGame(t *testing.T) {
	for name, store := range map[string]*characterStore{
		"named":              {characters: []statedb.PlayerCharacter{namedCharacter(7, "Mara")}},
		"no characters":      {},
		"failed name lookup": {listErr: errors.New("no connection")},
	} {
		t.Run(name, func(t *testing.T) {
			m := sendUntilApp(t, renameLauncher(t, store), campaign.SelectedMsg{Campaign: makeTestCampaign(1, "Saved")})
			if _, ok := m.(App); !ok {
				t.Fatalf("the game did not open:\n%s", m.View())
			}
		})
	}
}

func TestRenamePromptCanBeLeftAndReportsAFailedSave(t *testing.T) {
	open := func(t *testing.T) (*characterStore, tea.Model) {
		store := &characterStore{characters: []statedb.PlayerCharacter{namedCharacter(7, "  ")}}
		return store, sendUntilApp(t, renameLauncher(t, store), campaign.SelectedMsg{Campaign: makeTestCampaign(1, "Saved")})
	}

	t.Run("escape returns to the campaign list", func(t *testing.T) {
		store, m := open(t)
		m, cmd := m.Update(key(tea.KeyEsc))
		m, _ = m.Update(cmd())
		if got := m.(Launcher).state; got != launcherLoading || len(store.renames) != 0 {
			t.Fatalf("state = %v, renames = %v", got, store.renames)
		}
	})

	t.Run("failed save", func(t *testing.T) {
		store, m := open(t)
		store.renameErr = errors.New("write refused")
		m, _ = m.Update(typed("Mara"))
		m = sendUntilApp(t, m, key(tea.KeyEnter))
		l, ok := m.(Launcher)
		if !ok || !strings.Contains(l.errMsg, "write refused") {
			t.Fatalf("want the save error on the campaign list:\n%s", m.View())
		}
	})
}
