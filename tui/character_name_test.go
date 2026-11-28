package tui

import (
	"context"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"git.subcult.tv/subculture-collective/edda/internal/llm"
	"git.subcult.tv/subculture-collective/edda/internal/world"
	"git.subcult.tv/subculture-collective/edda/tui/campaign"
)

type countingProvider struct {
	stubLLMProvider
	calls *int
}

func (p countingProvider) Complete(context.Context, []llm.Message, []llm.Tool) (*llm.Response, error) {
	*p.calls++
	return &llm.Response{}, nil
}

// worldBuildError runs the commands that follow a confirmation, side by side as
// Bubble Tea does, until the world build reports an error.
func worldBuildError(t *testing.T, m tea.Model, cmd tea.Cmd) error {
	t.Helper()
	messages := make(chan tea.Msg, 16)
	var run func(tea.Cmd)
	run = func(cmd tea.Cmd) {
		if cmd != nil {
			go func() { messages <- cmd() }()
		}
	}
	run(cmd)
	deadline := time.After(5 * time.Second)
	for {
		select {
		case msg := <-messages:
			switch msg := msg.(type) {
			case tea.BatchMsg:
				for _, c := range msg {
					run(c)
				}
			case campaign.WorldErrorMsg:
				return msg.Err
			case campaign.WorldReadyMsg:
				t.Fatal("the world was built")
			case nil:
			default:
				next, follow := m.Update(msg)
				m = next
				run(follow)
			}
		case <-deadline:
			t.Fatal("the world build reported no error")
		}
	}
}

func TestConfirmingABlankCharacterNameBuildsNothing(t *testing.T) {
	for _, name := range []string{"", "   "} {
		calls := 0
		var m tea.Model = resize(t, newTestLauncherWithProvider(countingProvider{calls: &calls}), 80, 24)
		next, cmd := m.Update(campaign.ConfirmedMsg{
			Name: "Ashfall", Profile: testCampaignProfile(),
			CharacterProfile: world.CharacterProfile{Name: name, Concept: "human barbarian", Background: "Acolyte"},
		})
		err := worldBuildError(t, next, cmd)
		if err == nil || err.Error() != "character name is required" {
			t.Errorf("name %q: error = %v, want \"character name is required\"", name, err)
		}
		if calls != 0 {
			t.Errorf("name %q: %d model requests, want none", name, calls)
		}
	}
}

func TestCharacterFormAsksForAName(t *testing.T) {
	for _, name := range []string{"", "   "} {
		var m tea.Model = resize(t, newTestLauncherWithProvider(stubLLMProvider{}), 80, 24)
		m = send(t, m, campaign.CharMethodChosenMsg{Method: campaign.MethodAttributes})
		if name != "" {
			m = send(t, m, typed(name))
		}
		m = send(t, m, key(tea.KeyEnter))

		view := m.View()
		if !strings.Contains(view, "name is required") || strings.Contains(view, "timeout") {
			t.Errorf("name %q: want \"name is required\" and no \"timeout\":\n%s", name, view)
		}
		if !strings.Contains(view, "Character Name") {
			t.Errorf("name %q: the form left the name field:\n%s", name, view)
		}
	}
}
