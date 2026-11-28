//go:build integration

package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/handlers"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
)

// The TUI's data tabs read through the engine. They must show what the web
// client's character, inventory and quest routes answer.
func TestIntegrationPlayerRecordsMatchTheWebRoutes(t *testing.T) {
	ctx := context.Background()
	q := statedb.New(testPool)
	fighter, err := os.ReadFile("../rules/srd521/testdata/fighter1.json")
	if err != nil {
		t.Fatal(err)
	}
	gameEngine, err := engine.New(testPool, nil, config.LLMConfig{})
	if err != nil {
		t.Fatal(err)
	}
	characters := &handlers.CharacterHandlers{Queries: q, Pool: testPool, Logger: log.New(io.Discard)}
	world := &handlers.WorldHandlers{Queries: q, Logger: log.New(io.Discard)}
	router := chi.NewRouter()
	router.Get("/campaigns/{id}/character", characters.GetCharacter)
	router.Get("/campaigns/{id}/character/inventory", characters.GetCharacterInventory)
	router.Get("/campaigns/{id}/quests", world.ListQuests)
	router.Get("/campaigns/{id}/quests/{qid}", world.GetQuest)

	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	newCampaign := func(name, rules string) statedb.Campaign {
		t.Helper()
		c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: name, RulesVersion: rules, Status: "active", CreatedBy: dbutil.ToPgtype(testUserID)})
		must(err)
		return c
	}

	// Another campaign's quest and item must not appear.
	other := newCampaign("Other", "legacy")
	_, err = q.CreateQuest(ctx, statedb.CreateQuestParams{CampaignID: other.ID, Title: "Foreign quest", QuestType: "short_term", Status: "active"})
	must(err)

	for _, rules := range []string{"legacy", "srd-5.2.1"} {
		t.Run(rules, func(t *testing.T) {
			c := newCampaign("Records "+rules, rules)
			base := "/campaigns/" + dbutil.FromPgtype(c.ID).String()
			get := func(path string, into any) {
				t.Helper()
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, base+path, nil))
				if rec.Code != http.StatusOK {
					t.Fatalf("GET %s: %d %s", path, rec.Code, rec.Body.String())
				}
				must(json.Unmarshal(rec.Body.Bytes(), into))
			}
			location, err := q.CreateLocation(ctx, statedb.CreateLocationParams{CampaignID: c.ID, Name: "The Drowned Gate"})
			must(err)
			var owners []pgtype.UUID
			for _, name := range []string{"Older", "Newer"} {
				pc, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: name, Stats: fighter, Hp: 9, MaxHp: 12, Level: 1, Status: "active", Abilities: []byte(`[]`), CurrentLocationID: location.ID})
				must(err)
				owners = append(owners, pc.ID)
				_, err = q.CreateItem(ctx, statedb.CreateItemParams{CampaignID: c.ID, PlayerCharacterID: pc.ID, Name: name + "'s lantern", ItemType: "misc", Rarity: "common", Quantity: 1})
				must(err)
			}
			for _, status := range []string{"completed", "active"} {
				quest, err := q.CreateQuest(ctx, statedb.CreateQuestParams{CampaignID: c.ID, Title: status + " debt", QuestType: "short_term", Status: status})
				must(err)
				for i, description := range []string{"second", "first"} {
					_, err = q.CreateObjective(ctx, statedb.CreateObjectiveParams{QuestID: quest.ID, Description: description, OrderIndex: int32(1 - i)})
					must(err)
				}
			}

			records, err := gameEngine.PlayerRecords(ctx, dbutil.FromPgtype(c.ID))
			must(err)

			var character api.CharacterResponse
			get("/character", &character)
			got := records.Character
			if got == nil || got.Name != "Newer" || got.ID.String() != character.ID || got.HP != character.HP || got.MaxHP != character.MaxHP || got.RulesVersion != character.RulesVersion {
				t.Fatalf("character = %+v, web %+v", got, character)
			}
			if got.Location != "The Drowned Gate" || got.SheetError != nil {
				t.Errorf("location %q, sheet error %v", got.Location, got.SheetError)
			}
			if (got.Sheet != nil) != (rules == "srd-5.2.1") {
				t.Fatalf("sheet = %+v for %s", got.Sheet, rules)
			}
			if engineSheet, webSheet := encode(t, got.Sheet), encode(t, character.Sheet); engineSheet != webSheet {
				t.Errorf("sheet differs from the web's:\n%s\n%s", engineSheet, webSheet)
			}

			var items []api.ItemResponse
			get("/character/inventory", &items)
			if len(items) != 1 || len(records.Inventory) != 1 || records.Inventory[0].Name != "Newer's lantern" || records.Inventory[0].ID.String() != items[0].ID {
				t.Errorf("inventory = %+v, web %+v", records.Inventory, items)
			}

			var quests []api.QuestResponse
			get("/quests", &quests)
			if len(quests) != 2 || len(records.Quests) != 2 {
				t.Fatalf("quests = %+v, web %+v", records.Quests, quests)
			}
			for i, listed := range quests {
				var full api.QuestResponse
				get("/quests/"+listed.ID, &full)
				quest := records.Quests[i]
				if quest.ID.String() != full.ID || quest.Title != full.Title || string(quest.Status) != full.Status || len(quest.Objectives) != 2 || len(full.Objectives) != 2 {
					t.Fatalf("quest %d = %+v, web %+v", i, quest, full)
				}
				for j, objective := range quest.Objectives {
					if objective.Description != full.Objectives[j].Description {
						t.Errorf("quest %d objective %d = %q, web %q", i, j, objective.Description, full.Objectives[j].Description)
					}
				}
			}
		})
	}

	t.Run("a campaign with nothing yet", func(t *testing.T) {
		records, err := gameEngine.PlayerRecords(ctx, dbutil.FromPgtype(newCampaign("Empty", "legacy").ID))
		must(err)
		if records.Character != nil || len(records.Inventory) != 0 || len(records.Quests) != 0 {
			t.Errorf("records = %+v, want none", records)
		}
	})

	t.Run("an unsupported SRD build", func(t *testing.T) {
		c := newCampaign("Unsupported", "srd-5.2.1")
		_, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: "Odd", Stats: []byte(`{"srd521":{"class":"wizard"}}`), Hp: 1, MaxHp: 1, Level: 1, Status: "active", Abilities: []byte(`[]`)})
		must(err)
		records, err := gameEngine.PlayerRecords(ctx, dbutil.FromPgtype(c.ID))
		must(err)
		if records.Character == nil || records.Character.Sheet != nil || records.Character.SheetError == nil {
			t.Errorf("character = %+v, want a sheet error", records.Character)
		}
	})
}

func encode(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
