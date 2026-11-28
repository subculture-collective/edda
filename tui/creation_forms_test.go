package tui

import (
	"context"
	"errors"
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"git.subcult.tv/subculture-collective/edda/internal/llm"
	"git.subcult.tv/subculture-collective/edda/internal/prompt"
	"git.subcult.tv/subculture-collective/edda/internal/world"
	"git.subcult.tv/subculture-collective/edda/tui/campaign"
)

// promptRecorder keeps each prompt it is sent and fails the request, which
// brings the player back to the attributes form.
type promptRecorder struct {
	stubLLMProvider
	prompts *[]string
}

func (p promptRecorder) Complete(_ context.Context, messages []llm.Message, _ []llm.Tool) (*llm.Response, error) {
	*p.prompts = append(*p.prompts, messages[0].Content)
	return nil, errors.New("no model")
}

// choose moves the cursor of a form's list by steps, then confirms the field.
func choose(t *testing.T, m tea.Model, steps int) tea.Model {
	t.Helper()
	move := key(tea.KeyDown)
	if steps < 0 {
		move, steps = key(tea.KeyUp), -steps
	}
	for range steps {
		m = send(t, m, move)
	}
	return send(t, m, key(tea.KeyEnter))
}

func TestAttributesFormSendsTheChosenValues(t *testing.T) {
	// A round moves from the values the form shows: the first options at the
	// start, and the previous choices after the failed request.
	type round struct {
		steps                     [3]int
		genre, settingStyle, tone string
	}
	for _, tc := range []struct {
		name   string
		rounds []round
	}{
		{"genre", []round{{[3]int{2, 0, 0}, "Horror", "Open Wilderness", "Gritty and Dark"}}},
		{"setting style", []round{{[3]int{0, 3, 0}, "Fantasy", "Underground", "Gritty and Dark"}}},
		{"tone", []round{{[3]int{0, 0, 4}, "Fantasy", "Open Wilderness", "Humorous"}}},
		{"all three, then again after a failure", []round{
			{[3]int{6, 5, 5}, "Steampunk", "Peaceful Realm", "Tense and Suspenseful"},
			{[3]int{0, 0, 0}, "Steampunk", "Peaceful Realm", "Tense and Suspenseful"},
			{[3]int{-1, -2, -3}, "Post-Apocalyptic", "Underground", "Epic and Grand"},
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var prompts []string
			var m tea.Model = resize(t, newTestLauncherWithProvider(promptRecorder{prompts: &prompts}), 80, 24)
			m = send(t, m, campaign.MethodChosenMsg{Method: campaign.MethodAttributes})

			for i, r := range tc.rounds {
				for _, steps := range r.steps {
					m = choose(t, m, steps)
				}
				if len(prompts) != i+1 {
					t.Fatalf("round %d: %d proposal requests, want %d\n%s", i, len(prompts), i+1, m.View())
				}
				if want := prompt.BuildProposalsPrompt(r.genre, r.settingStyle, r.tone); prompts[i] != want {
					sent := m.(Launcher).proposalAttrs
					t.Fatalf("round %d: requested %+v, want %q, %q, %q", i, sent, r.genre, r.settingStyle, r.tone)
				}
			}
		})
	}
}

func TestCharacterFormSendsTheChosenValues(t *testing.T) {
	launcher := newTestLauncherWithProvider(stubLLMProvider{})
	profile := testCampaignProfile()
	launcher.campaignProfile = &profile
	var m tea.Model = resize(t, launcher, 80, 24)
	m = send(t, m, campaign.CharMethodChosenMsg{Method: campaign.MethodAttributes})

	m = send(t, m, typed("Mara"))
	m = send(t, m, key(tea.KeyEnter))
	for _, steps := range []int{1, 2, 3, 4} { // race, class, background, alignment
		m = choose(t, m, steps)
	}

	want := world.BuildCharacterFromAttributes(world.CharacterAttributes{
		Name:       "Mara",
		Race:       world.D5ERaces[1],
		Class:      world.D5EClasses[2],
		Background: world.D5EBackgrounds[3],
		Alignment:  world.D5EAlignments[4],
	})
	got := m.(Launcher).characterProfile
	if got == nil {
		t.Fatalf("the form did not finish:\n%s", m.View())
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("character = %+v, want %+v", *got, *want)
	}
}
