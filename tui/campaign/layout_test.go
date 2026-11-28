package campaign

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"git.subcult.tv/subculture-collective/edda/internal/world"
)

// longReply is a game-master reply far wider than any tested terminal.
var longReply = strings.Repeat("Tell me which debts the drowned kingdom still keeps. ", 8) + "ENDOFREPLY"

type sizedModel interface {
	tea.Model
	SetSize(width, height int)
}

func TestInterviewsWrapLongReplies(t *testing.T) {
	for _, tc := range []struct {
		name  string
		model func(t *testing.T) sizedModel
		reply tea.Msg
	}{
		{"campaign", func(t *testing.T) sizedModel {
			m := NewInterviewModel(t.Context(), nil)
			return &m
		}, interviewResponseMsg{message: longReply}},
		{"character", func(t *testing.T) sizedModel {
			m := NewCharInterviewModel(t.Context(), nil, &world.CampaignProfile{})
			return &m
		}, charInterviewResponseMsg{message: longReply}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := tc.model(t)
			m.SetSize(80, 24)
			updated, _ := m.Update(tc.reply)
			view := updated.View()
			assertShows(t, view, 80, 24, "GM:", "ENDOFREPLY")

			// A narrower terminal wraps the same history again.
			resized := asSized(t, updated)
			resized.SetSize(50, 24)
			assertShows(t, resized.View(), 50, 24, "GM:", "ENDOFREPLY")
		})

		t.Run(tc.name+"/typing leaves the history in place", func(t *testing.T) {
			m := tc.model(t)
			m.SetSize(80, 10)
			updated, _ := m.Update(tc.reply)
			updated, _ = updated.Update(tc.reply)
			updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyPgUp})
			scrolled := historyOf(updated)
			// The viewport's default bindings page on these letters.
			for _, key := range []string{"b", "u", "k", " ", "f", "d", "j"} {
				updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
				if got := historyOf(updated); got != scrolled {
					t.Fatalf("typing %q scrolled the history:\n%s", key, updated.View())
				}
			}
		})
	}
}

// asSized returns m as a model that can be resized.
func asSized(t *testing.T, m tea.Model) sizedModel {
	t.Helper()
	switch v := m.(type) {
	case InterviewModel:
		return &v
	case CharInterviewModel:
		return &v
	}
	t.Fatalf("unexpected model type %T", m)
	return nil
}

// historyOf returns the chat history an interview has on screen.
func historyOf(m tea.Model) string {
	switch v := m.(type) {
	case InterviewModel:
		return v.viewport.View()
	case CharInterviewModel:
		return v.viewport.View()
	}
	return ""
}

func assertShows(t *testing.T, view string, width, height int, want ...string) {
	t.Helper()
	if w, h := lipgloss.Width(view), lipgloss.Height(view); w > width || h > height {
		t.Errorf("view is %dx%d, terminal is %dx%d\n%s", w, h, width, height, view)
	}
	for _, text := range want {
		if !strings.Contains(view, text) {
			t.Errorf("%q is not on screen:\n%s", text, view)
		}
	}
}

func TestConfirmationKeys(t *testing.T) {
	newModel := func(background string) tea.Model {
		m := NewConfirmationModel("The Drowned Crown", "Short.", world.CampaignProfile{},
			world.CharacterProfile{Name: "Mara", Background: background})
		m.SetSize(80, 24)
		return m
	}
	press := func(m tea.Model, keys ...tea.KeyType) tea.Model {
		for _, key := range keys {
			m, _ = m.Update(tea.KeyMsg{Type: key})
		}
		return m
	}
	chosen := func(m tea.Model) tea.Msg {
		t.Helper()
		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		if cmd == nil {
			t.Fatal("Enter returned no command")
		}
		return cmd()
	}

	t.Run("up and down choose a button when the text fits", func(t *testing.T) {
		m := press(newModel("Short."), tea.KeyDown)
		if _, ok := chosen(m).(ChangeMsg); !ok {
			t.Fatalf("Down did not choose Go Back:\n%s", m.View())
		}
		if _, ok := chosen(press(m, tea.KeyUp)).(ConfirmedMsg); !ok {
			t.Fatal("Up did not choose Begin Adventure")
		}
	})

	t.Run("up and down scroll long text; left and right choose", func(t *testing.T) {
		start := newModel(strings.Repeat("A long history of debts. ", 120))
		down := press(start, tea.KeyDown)
		if down.View() == start.View() {
			t.Fatalf("Down did not scroll:\n%s", down.View())
		}
		if _, ok := chosen(down).(ConfirmedMsg); !ok {
			t.Fatal("scrolling changed the chosen button")
		}
		if got := press(down, tea.KeyUp).View(); got != start.View() {
			t.Fatalf("Up did not scroll back:\n%s", got)
		}
		if _, ok := chosen(press(start, tea.KeyRight)).(ChangeMsg); !ok {
			t.Fatal("Right did not choose Go Back")
		}
		if _, ok := chosen(press(start, tea.KeyRight, tea.KeyLeft)).(ConfirmedMsg); !ok {
			t.Fatal("Left did not choose Begin Adventure")
		}
	})
}
