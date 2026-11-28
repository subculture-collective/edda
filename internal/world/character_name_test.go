package world

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/llm"
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
