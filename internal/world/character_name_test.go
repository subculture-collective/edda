package world

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/llm"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// blankNames are names that leave a character with nothing to be called.
var blankNames = []string{"", "   ", "\t\n", " 　"}

// countingLLM counts the requests it receives and answers with its response.
type countingLLM struct {
	llm.Provider
	calls    int
	response llm.Response
}

func (p *countingLLM) Complete(context.Context, []llm.Message, []llm.Tool) (*llm.Response, error) {
	p.calls++
	response := p.response
	return &response, nil
}

// The TUI and the API both create characters through Orchestrator.Run.
func TestOrchestratorRejectsABlankCharacterName(t *testing.T) {
	for _, version := range []domain.RulesVersion{"", domain.RulesVersionLegacy, domain.RulesVersionSRD521} {
		for _, name := range blankNames {
			provider, queries := &countingLLM{}, &stubQuerier{}
			result, err := NewOrchestrator(provider, queries).Run(t.Context(), OrchestratorInput{
				Name: "Ashfall", Profile: testProfile(), UserID: uuid.New(), RulesVersion: version,
				CharacterProfile: &CharacterProfile{Name: name, Concept: "human fighter", Background: "Soldier"},
			}, nil)
			if result != nil || err == nil || err.Error() != "character name is required" {
				t.Errorf("rules %q, name %q: result = %v, error = %v, want \"character name is required\"", version, name, result, err)
			}
			if provider.calls != 0 || queries.createCampaignCalls != 0 || queries.createPlayerCharacterCalls != 0 {
				t.Errorf("rules %q, name %q: %d model requests, %d campaigns, %d characters; want none",
					version, name, provider.calls, queries.createCampaignCalls, queries.createPlayerCharacterCalls)
			}
		}
	}
}

// extract_character_profile is the only model tool that names a player character.
func TestCharacterInterviewRejectsABlankNameFromTheModel(t *testing.T) {
	extract := func(name string) llm.Response {
		return llm.Response{ToolCalls: []llm.ToolCall{{ID: "1", Name: extractCharacterToolName, Arguments: map[string]any{
			"name": name, "concept": "human fighter", "background": "Soldier", "personality": "stoic",
			"motivations": []any{"duty"}, "strengths": []any{"steady"}, "weaknesses": []any{"proud"},
			"rules_choices": map[string]any{"origin": "Human", "class": "Fighter", "background": "Soldier"},
		}}}}
	}
	for _, name := range blankNames {
		t.Run("legacy "+strings.TrimSpace(name), func(t *testing.T) {
			ci := NewCharacterInterviewer(&countingLLM{response: llm.Response{
				ToolCalls: []llm.ToolCall{{ID: "1", Name: extractCharacterToolName, Arguments: map[string]any{
					"name": name, "concept": "ranger", "background": "scout", "personality": "stoic",
					"motivations": []any{"duty"}, "strengths": []any{"steady"}, "weaknesses": []any{"proud"},
				}}},
			}}, nil)
			result, err := ci.Start(t.Context())
			if err == nil || ci.Done() || ci.Profile() != nil {
				t.Errorf("name %q: result = %+v, error = %v; want a refused profile", name, result, err)
			}
		})
		t.Run("srd "+strings.TrimSpace(name), func(t *testing.T) {
			ci := NewRulesCharacterInterviewer(&countingLLM{response: extract(name)}, nil, domain.RulesVersionSRD521)
			result, err := ci.Start(t.Context())
			if err != nil || result.Done || ci.Profile() != nil || !strings.Contains(result.Message, "name") {
				t.Errorf("name %q: result = %+v, error = %v; want a request for a name", name, result, err)
			}
		})
	}
}

// Surrounding white space is not part of a name.
func TestOrchestratorTrimsTheCharacterName(t *testing.T) {
	var saved []string
	startingLocationID := uuid.New()
	q := &stubQuerier{
		createCampaignFn: func(context.Context, statedb.CreateCampaignParams) (statedb.Campaign, error) {
			return statedb.Campaign{ID: pgUUID(uuid.New())}, nil
		},
		listLocationsByCampaignFn: func(_ context.Context, campaignID pgtype.UUID) ([]statedb.Location, error) {
			return []statedb.Location{{ID: pgUUID(startingLocationID), CampaignID: campaignID, Name: validResponse().StartingLocation}}, nil
		},
		createPlayerCharacterFn: func(_ context.Context, arg statedb.CreatePlayerCharacterParams) (statedb.PlayerCharacter, error) {
			saved = append(saved, arg.Name)
			// Stop here: only the saved name matters.
			return statedb.PlayerCharacter{}, errors.New("stop after the character")
		},
	}
	provider := &sequentialStubLLM{t: t, scripts: []stubSkeletonLLM{
		{content: mustJSON(t, validResponse())}, {content: mustSceneJSON(t, validSceneResponse())},
	}}
	character := testCharacterProfile()
	character.Name = " \t Mara of  Vell　\n"

	_, err := NewOrchestrator(provider, q).Run(t.Context(), OrchestratorInput{
		Name: "Ashfall", Profile: testProfile(), CharacterProfile: character, UserID: uuid.New(),
	}, nil)

	if len(saved) != 1 || saved[0] != "Mara of  Vell" {
		t.Fatalf("saved names = %q, want only %q (error: %v)", saved, "Mara of  Vell", err)
	}
	if character.Name != " \t Mara of  Vell　\n" {
		t.Errorf("the caller's profile changed to %q", character.Name)
	}
}
