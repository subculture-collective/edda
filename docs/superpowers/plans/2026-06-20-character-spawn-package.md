# Character Spawn Package Implementation Plan

> **For agentic workers:** Execute this plan task-by-task. Recommended path:
> dispatch a fresh subagent per task, review each result with `review-quality`,
> then continue. For complex multi-agent splits, use
> `parallel-feature-development`, `team-composition-patterns`, and
> `team-communication-protocols`. Steps use checkbox (`- [ ]`) syntax for
> tracking.

**Goal:** Let a newly created Player Character start with inventory, player-known lore facts, and player-aware NPC relationships.

**Architecture:** Add a narrow `CharacterSpawnPackage` API/input type and apply it after `PersistCharacterProfile` but before opening-scene generation in `internal/world/orchestrator.go`. The implementation should use existing persistence seams (`items`, `world_facts`, `entity_relationships`) instead of new schema or a generic knowledge blob.

**Tech Stack:** Go, sqlc-generated `statedb.Querier`, existing startup HTTP API, existing `internal/world` tests.

---

## File Structure

- Modify `pkg/api/types.go`
  - Add request DTOs: `CharacterSpawnPackage`, `StarterItem`, `StarterKnownFact`, `StarterRelationship`.
  - Add `SpawnPackage *CharacterSpawnPackage` to `WorldBuildRequest`.
- Create `internal/world/spawn.go`
  - Define world-layer spawn types and `ApplySpawnPackage(ctx, q, campaignID, playerID, pkg)`.
  - Map starter inventory to `CreateItem`, known facts to `CreateFact`, relationships to `CreateRelationship`.
- Modify `internal/world/orchestrator.go`
  - Add `SpawnPackage *CharacterSpawnPackage` to `OrchestratorInput`.
  - Apply it immediately after character persistence and before scene generation.
- Modify `internal/handlers/startup.go`
  - Convert `api.CharacterSpawnPackage` into `world.CharacterSpawnPackage` and pass it to the orchestrator.
- Modify `internal/world/persist_adapter_orchestrator_test.go`
  - Extend `stubQuerier` with `CreateItem` and `CreateRelationship` support.
  - Add unit tests for spawn application.
  - Add orchestrator ordering/propagation test if practical.
- Modify `internal/handlers/startup_test.go`
  - Extend `startupStubQuerier` with `CreateItem` and `CreateRelationship` support.
  - Add handler test proving JSON spawn payload reaches persistence.

## Spawn Contract

Use this JSON shape in `WorldBuildRequest`:

```json
{
  "spawn_package": {
    "items": [
      {
        "name": "Rusty Lantern",
        "description": "A dented lantern that still burns cleanly.",
        "item_type": "tool",
        "rarity": "common",
        "quantity": 1,
        "equipped": false
      }
    ],
    "known_facts": [
      {
        "fact": "Ash storms swallow the roads at dusk.",
        "category": "environment"
      }
    ],
    "relationships": [
      {
        "target_entity_type": "npc",
        "target_entity_id": "11111111-1111-1111-1111-111111111111",
        "relationship_type": "owes_debt",
        "description": "Marshal Vey once saved Kael's life.",
        "strength": 30
      }
    ]
  }
}
```

Defaults:

- Item `item_type`: `misc`
- Item `rarity`: `common`
- Item `quantity`: `1` when omitted or less than `1`
- Item `properties`: `{}`
- Known fact `category`: `lore`
- Known fact `source`: `character_spawn`
- Relationship source is always the newly persisted player character: `source_entity_type="player_character"`, `source_entity_id=<playerID>`.

Do not implement codex entity creation in this pass. Codex tables are already player-known but require selecting or creating concrete language/culture/belief/economic records; that should be a follow-up once the UX decides whether spawn packages reference generated world entities by ID/name or create new codex entries.

---

### Task 1: Add API and world spawn types

**Files:**
- Modify: `pkg/api/types.go:262-269`
- Create: `internal/world/spawn.go`

- [ ] **Step 1: Add API DTOs**

In `pkg/api/types.go`, replace the `WorldBuildRequest` block with:

```go
// WorldBuildRequest finalizes startup choices and creates the campaign world.
type WorldBuildRequest struct {
	Name             string                 `json:"name"`
	Summary          string                 `json:"summary"`
	Profile          *CampaignProfile       `json:"profile"`
	CharacterProfile *CharacterProfile      `json:"character_profile"`
	RulesMode        string                 `json:"rules_mode,omitempty"`
	SpawnPackage     *CharacterSpawnPackage `json:"spawn_package,omitempty"`
}

// CharacterSpawnPackage describes starting campaign state granted to the player character.
type CharacterSpawnPackage struct {
	Items         []StarterItem         `json:"items,omitempty"`
	KnownFacts    []StarterKnownFact    `json:"known_facts,omitempty"`
	Relationships []StarterRelationship `json:"relationships,omitempty"`
}

type StarterItem struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	ItemType    string         `json:"item_type,omitempty"`
	Rarity      string         `json:"rarity,omitempty"`
	Properties  map[string]any `json:"properties,omitempty"`
	Equipped    bool           `json:"equipped,omitempty"`
	Quantity    int32          `json:"quantity,omitempty"`
}

type StarterKnownFact struct {
	Fact     string `json:"fact"`
	Category string `json:"category,omitempty"`
}

type StarterRelationship struct {
	TargetEntityType string `json:"target_entity_type"`
	TargetEntityID   string `json:"target_entity_id"`
	RelationshipType string `json:"relationship_type"`
	Description      string `json:"description,omitempty"`
	Strength         *int32 `json:"strength,omitempty"`
}
```

- [ ] **Step 2: Create world spawn types and mapper helpers**

Create `internal/world/spawn.go` with package `world` and types mirroring the API package without importing `pkg/api`:

```go
package world

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

type CharacterSpawnPackage struct {
	Items         []StarterItem
	KnownFacts    []StarterKnownFact
	Relationships []StarterRelationship
}

type StarterItem struct {
	Name        string
	Description string
	ItemType    string
	Rarity      string
	Properties  map[string]any
	Equipped    bool
	Quantity    int32
}

type StarterKnownFact struct {
	Fact     string
	Category string
}

type StarterRelationship struct {
	TargetEntityType string
	TargetEntityID   uuid.UUID
	RelationshipType string
	Description      string
	Strength         *int32
}
```

- [ ] **Step 3: Run typecheck/test for expected failures**

Run: `go test ./internal/world ./pkg/api`

Expected: compile fails until `ApplySpawnPackage` exists and/or unused imports are resolved in the next task.

---

### Task 2: Implement spawn package application

**Files:**
- Modify: `internal/world/spawn.go`
- Test: `internal/world/persist_adapter_orchestrator_test.go`

- [ ] **Step 1: Extend `stubQuerier` in `internal/world/persist_adapter_orchestrator_test.go`**

Add fields to `stubQuerier`:

```go
	createItemFn         func(context.Context, statedb.CreateItemParams) (statedb.Item, error)
	createRelationshipFn func(context.Context, statedb.CreateRelationshipParams) (statedb.EntityRelationship, error)

	createItemCalls         int
	createRelationshipCalls int

	lastCreateItemParams         statedb.CreateItemParams
	lastCreateRelationshipParams statedb.CreateRelationshipParams
```

Add methods:

```go
func (s *stubQuerier) CreateItem(ctx context.Context, arg statedb.CreateItemParams) (statedb.Item, error) {
	s.createItemCalls++
	s.lastCreateItemParams = arg
	if s.createItemFn != nil {
		return s.createItemFn(ctx, arg)
	}
	return statedb.Item{}, nil
}

func (s *stubQuerier) CreateRelationship(ctx context.Context, arg statedb.CreateRelationshipParams) (statedb.EntityRelationship, error) {
	s.createRelationshipCalls++
	s.lastCreateRelationshipParams = arg
	if s.createRelationshipFn != nil {
		return s.createRelationshipFn(ctx, arg)
	}
	return statedb.EntityRelationship{}, nil
}
```

- [ ] **Step 2: Add failing spawn application test**

Add this test to `internal/world/persist_adapter_orchestrator_test.go`:

```go
func TestApplySpawnPackage_CreatesInventoryKnownFactsAndRelationships(t *testing.T) {
	campaignID := uuid.New()
	playerID := uuid.New()
	targetNPCID := uuid.New()
	strength := int32(30)
	q := &stubQuerier{
		createItemFn: func(_ context.Context, arg statedb.CreateItemParams) (statedb.Item, error) {
			return statedb.Item{ID: pgUUID(uuid.New()), Name: arg.Name}, nil
		},
		createFactFn: func(_ context.Context, arg statedb.CreateFactParams) (statedb.WorldFact, error) {
			return statedb.WorldFact{ID: pgUUID(uuid.New()), Fact: arg.Fact}, nil
		},
		createRelationshipFn: func(_ context.Context, arg statedb.CreateRelationshipParams) (statedb.EntityRelationship, error) {
			return statedb.EntityRelationship{ID: pgUUID(uuid.New()), RelationshipType: arg.RelationshipType}, nil
		},
	}

	err := ApplySpawnPackage(context.Background(), q, campaignID, playerID, &CharacterSpawnPackage{
		Items: []StarterItem{{Name: "Rusty Lantern", Description: "Still burns.", ItemType: "tool", Rarity: "common", Quantity: 2, Equipped: true, Properties: map[string]any{"light": true}}},
		KnownFacts: []StarterKnownFact{{Fact: "Ash storms swallow the roads at dusk.", Category: "environment"}},
		Relationships: []StarterRelationship{{TargetEntityType: "npc", TargetEntityID: targetNPCID, RelationshipType: "owes_debt", Description: "Marshal Vey once saved Kael.", Strength: &strength}},
	})
	if err != nil {
		t.Fatalf("ApplySpawnPackage returned error: %v", err)
	}

	if q.createItemCalls != 1 {
		t.Fatalf("CreateItem calls = %d, want 1", q.createItemCalls)
	}
	if got := q.lastCreateItemParams.CampaignID; got != pgUUID(campaignID) {
		t.Errorf("item campaign id = %#v, want %#v", got, pgUUID(campaignID))
	}
	if got := q.lastCreateItemParams.PlayerCharacterID; got != pgUUID(playerID) {
		t.Errorf("item player id = %#v, want %#v", got, pgUUID(playerID))
	}
	if got := q.lastCreateItemParams.Name; got != "Rusty Lantern" {
		t.Errorf("item name = %q, want Rusty Lantern", got)
	}
	if got := q.lastCreateItemParams.Description; !got.Valid || got.String != "Still burns." {
		t.Errorf("item description = %#v, want valid Still burns.", got)
	}
	if got := q.lastCreateItemParams.ItemType; got != "tool" {
		t.Errorf("item type = %q, want tool", got)
	}
	if got := q.lastCreateItemParams.Rarity; got != "common" {
		t.Errorf("item rarity = %q, want common", got)
	}
	if !q.lastCreateItemParams.Equipped {
		t.Error("item equipped = false, want true")
	}
	if got := q.lastCreateItemParams.Quantity; got != 2 {
		t.Errorf("item quantity = %d, want 2", got)
	}

	var properties map[string]bool
	if err := json.Unmarshal(q.lastCreateItemParams.Properties, &properties); err != nil {
		t.Fatalf("unmarshal item properties: %v", err)
	}
	if !properties["light"] {
		t.Errorf("item properties = %v, want light=true", properties)
	}

	if q.createFactCalls != 1 {
		t.Fatalf("CreateFact calls = %d, want 1", q.createFactCalls)
	}
	if got := q.lastCreateFactParams.Fact; got != "Ash storms swallow the roads at dusk." {
		t.Errorf("fact = %q, want ash storm fact", got)
	}
	if got := q.lastCreateFactParams.Category; got != "environment" {
		t.Errorf("fact category = %q, want environment", got)
	}
	if got := q.lastCreateFactParams.Source; got != "character_spawn" {
		t.Errorf("fact source = %q, want character_spawn", got)
	}
	if !q.lastCreateFactParams.PlayerKnown {
		t.Error("fact player_known = false, want true")
	}

	if q.createRelationshipCalls != 1 {
		t.Fatalf("CreateRelationship calls = %d, want 1", q.createRelationshipCalls)
	}
	rel := q.lastCreateRelationshipParams
	if rel.SourceEntityType != "player_character" {
		t.Errorf("relationship source type = %q, want player_character", rel.SourceEntityType)
	}
	if rel.SourceEntityID != pgUUID(playerID) {
		t.Errorf("relationship source id = %#v, want %#v", rel.SourceEntityID, pgUUID(playerID))
	}
	if rel.TargetEntityType != "npc" || rel.TargetEntityID != pgUUID(targetNPCID) {
		t.Errorf("relationship target = %s/%#v, want npc/%#v", rel.TargetEntityType, rel.TargetEntityID, pgUUID(targetNPCID))
	}
	if rel.RelationshipType != "owes_debt" {
		t.Errorf("relationship type = %q, want owes_debt", rel.RelationshipType)
	}
	if !rel.Description.Valid || rel.Description.String != "Marshal Vey once saved Kael." {
		t.Errorf("relationship description = %#v, want valid description", rel.Description)
	}
	if !rel.Strength.Valid || rel.Strength.Int32 != 30 {
		t.Errorf("relationship strength = %#v, want 30", rel.Strength)
	}
}
```

- [ ] **Step 3: Implement `ApplySpawnPackage`**

In `internal/world/spawn.go`, add:

```go
func ApplySpawnPackage(ctx context.Context, q statedb.Querier, campaignID uuid.UUID, playerID uuid.UUID, pkg *CharacterSpawnPackage) error {
	if q == nil || pkg == nil {
		return nil
	}
	campaignPgID := dbutil.ToPgtype(campaignID)
	playerPgID := dbutil.ToPgtype(playerID)

	for _, item := range pkg.Items {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		properties := []byte(`{}`)
		if len(item.Properties) > 0 {
			encoded, err := json.Marshal(item.Properties)
			if err != nil {
				return fmt.Errorf("spawn package: encode item %q properties: %w", name, err)
			}
			properties = encoded
		}
		quantity := item.Quantity
		if quantity < 1 {
			quantity = 1
		}
		itemType := strings.TrimSpace(item.ItemType)
		if itemType == "" {
			itemType = "misc"
		}
		rarity := strings.TrimSpace(item.Rarity)
		if rarity == "" {
			rarity = "common"
		}
		_, err := q.CreateItem(ctx, statedb.CreateItemParams{
			CampaignID:        campaignPgID,
			PlayerCharacterID: playerPgID,
			Name:              name,
			Description:       pgtype.Text{String: strings.TrimSpace(item.Description), Valid: strings.TrimSpace(item.Description) != ""},
			ItemType:          itemType,
			Rarity:            rarity,
			Properties:        properties,
			Equipped:          item.Equipped,
			Quantity:          quantity,
		})
		if err != nil {
			return fmt.Errorf("spawn package: create item %q: %w", name, err)
		}
	}

	for _, fact := range pkg.KnownFacts {
		text := strings.TrimSpace(fact.Fact)
		if text == "" {
			continue
		}
		category := strings.TrimSpace(fact.Category)
		if category == "" {
			category = "lore"
		}
		_, err := q.CreateFact(ctx, statedb.CreateFactParams{
			CampaignID:  campaignPgID,
			Fact:        text,
			Category:    category,
			Source:      "character_spawn",
			PlayerKnown: true,
		})
		if err != nil {
			return fmt.Errorf("spawn package: create known fact %q: %w", text, err)
		}
	}

	for _, rel := range pkg.Relationships {
		targetType := strings.TrimSpace(rel.TargetEntityType)
		relType := strings.TrimSpace(rel.RelationshipType)
		if targetType == "" || rel.TargetEntityID == uuid.Nil || relType == "" {
			continue
		}
		strength := pgtype.Int4{}
		if rel.Strength != nil {
			strength = pgtype.Int4{Int32: *rel.Strength, Valid: true}
		}
		_, err := q.CreateRelationship(ctx, statedb.CreateRelationshipParams{
			CampaignID:       campaignPgID,
			SourceEntityType: "player_character",
			SourceEntityID:   playerPgID,
			TargetEntityType: targetType,
			TargetEntityID:   dbutil.ToPgtype(rel.TargetEntityID),
			RelationshipType: relType,
			Description:      pgtype.Text{String: strings.TrimSpace(rel.Description), Valid: strings.TrimSpace(rel.Description) != ""},
			Strength:         strength,
		})
		if err != nil {
			return fmt.Errorf("spawn package: create relationship %q: %w", relType, err)
		}
	}

	return nil
}
```

- [ ] **Step 4: Run test**

Run: `go test ./internal/world -run TestApplySpawnPackage_CreatesInventoryKnownFactsAndRelationships -v`

Expected: PASS.

---

### Task 3: Wire spawn package into world creation

**Files:**
- Modify: `internal/world/orchestrator.go`
- Modify: `internal/handlers/startup.go`
- Test: `internal/world/persist_adapter_orchestrator_test.go`

- [ ] **Step 1: Add orchestrator input field**

In `internal/world/orchestrator.go`, add to `OrchestratorInput`:

```go
	SpawnPackage     *CharacterSpawnPackage
```

- [ ] **Step 2: Apply package after player persistence**

In `internal/world/orchestrator.go`, after the `orchestrator character persisted` log and before `report("Setting the scene…")`, add:

```go
	if err := ApplySpawnPackage(ctx, o.queries, campaignUUID, dbCharacterID(characterRow), input.SpawnPackage); err != nil {
		logger().Error("orchestrator spawn package application failed", "campaign_id", campaignUUID, "character_id", dbCharacterID(characterRow), "error", err)
		return nil, fmt.Errorf("orchestrator: apply spawn package: %w", err)
	}
```

- [ ] **Step 3: Add API-to-world conversion helper**

In `internal/handlers/startup.go`, add imports for `github.com/google/uuid` if missing, then add below `BuildWorld`:

```go
func startupSpawnPackageToWorld(pkg *api.CharacterSpawnPackage) (*world.CharacterSpawnPackage, error) {
	if pkg == nil {
		return nil, nil
	}
	out := &world.CharacterSpawnPackage{
		Items:      make([]world.StarterItem, 0, len(pkg.Items)),
		KnownFacts: make([]world.StarterKnownFact, 0, len(pkg.KnownFacts)),
	}
	for _, item := range pkg.Items {
		out.Items = append(out.Items, world.StarterItem{
			Name:        item.Name,
			Description: item.Description,
			ItemType:    item.ItemType,
			Rarity:      item.Rarity,
			Properties:  item.Properties,
			Equipped:    item.Equipped,
			Quantity:    item.Quantity,
		})
	}
	for _, fact := range pkg.KnownFacts {
		out.KnownFacts = append(out.KnownFacts, world.StarterKnownFact{Fact: fact.Fact, Category: fact.Category})
	}
	for _, rel := range pkg.Relationships {
		targetID, err := uuid.Parse(strings.TrimSpace(rel.TargetEntityID))
		if err != nil {
			return nil, fmt.Errorf("invalid spawn relationship target_entity_id %q", rel.TargetEntityID)
		}
		out.Relationships = append(out.Relationships, world.StarterRelationship{
			TargetEntityType: rel.TargetEntityType,
			TargetEntityID:   targetID,
			RelationshipType: rel.RelationshipType,
			Description:      rel.Description,
			Strength:         rel.Strength,
		})
	}
	return out, nil
}
```

- [ ] **Step 4: Pass converted spawn package to orchestrator**

In `BuildWorld`, after validating `CharacterProfile`, add:

```go
	spawnPackage, err := startupSpawnPackageToWorld(req.SpawnPackage)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
```

Then include `SpawnPackage: spawnPackage,` in `world.OrchestratorInput`.

- [ ] **Step 5: Run package tests**

Run: `go test ./internal/world ./internal/handlers ./pkg/api`

Expected: PASS or only failures unrelated to this change. Fix compile errors before continuing.

---

### Task 4: Add handler coverage for startup spawn payload

**Files:**
- Modify: `internal/handlers/startup_test.go`

- [ ] **Step 1: Extend startup stub**

Add fields to `startupStubQuerier`:

```go
	createItemFn         func(context.Context, statedb.CreateItemParams) (statedb.Item, error)
	createRelationshipFn func(context.Context, statedb.CreateRelationshipParams) (statedb.EntityRelationship, error)
	lastCreateItemParams         statedb.CreateItemParams
	lastCreateRelationshipParams statedb.CreateRelationshipParams
	createItemCalls         int
	createRelationshipCalls int
```

Add methods:

```go
func (s *startupStubQuerier) CreateItem(ctx context.Context, arg statedb.CreateItemParams) (statedb.Item, error) {
	s.createItemCalls++
	s.lastCreateItemParams = arg
	if s.createItemFn != nil {
		return s.createItemFn(ctx, arg)
	}
	return statedb.Item{}, nil
}

func (s *startupStubQuerier) CreateRelationship(ctx context.Context, arg statedb.CreateRelationshipParams) (statedb.EntityRelationship, error) {
	s.createRelationshipCalls++
	s.lastCreateRelationshipParams = arg
	if s.createRelationshipFn != nil {
		return s.createRelationshipFn(ctx, arg)
	}
	return statedb.EntityRelationship{}, nil
}
```

- [ ] **Step 2: Add spawn payload assertions to `TestBuildWorld_Success`**

In `TestBuildWorld_Success`, create `playerID := uuid.New()` and `npcID := uuid.New()` near the existing IDs. Change `createPlayerCharacterFn` to return `playerID`.

Set `createItemFn` and `createRelationshipFn` in `queries`:

```go
		createItemFn: func(_ context.Context, arg statedb.CreateItemParams) (statedb.Item, error) {
			return statedb.Item{ID: dbutil.ToPgtype(uuid.New()), Name: arg.Name}, nil
		},
		createRelationshipFn: func(_ context.Context, arg statedb.CreateRelationshipParams) (statedb.EntityRelationship, error) {
			return statedb.EntityRelationship{ID: dbutil.ToPgtype(uuid.New()), RelationshipType: arg.RelationshipType}, nil
		},
```

Add `SpawnPackage` to the `api.WorldBuildRequest` body:

```go
		SpawnPackage: &api.CharacterSpawnPackage{
			Items: []api.StarterItem{{Name: "Rusty Lantern", Description: "Still burns.", ItemType: "tool", Quantity: 1}},
			KnownFacts: []api.StarterKnownFact{{Fact: "Kael knows the ash roads.", Category: "background"}},
			Relationships: []api.StarterRelationship{{TargetEntityType: "npc", TargetEntityID: npcID.String(), RelationshipType: "knows", Description: "Old patrol contact"}},
		},
```

After existing response assertions, add:

```go
	if queries.createItemCalls != 1 {
		t.Fatalf("CreateItem calls = %d, want 1", queries.createItemCalls)
	}
	if got := queries.lastCreateItemParams.PlayerCharacterID; got != dbutil.ToPgtype(playerID) {
		t.Errorf("starter item player id = %#v, want %#v", got, dbutil.ToPgtype(playerID))
	}
	if got := queries.lastCreateItemParams.Name; got != "Rusty Lantern" {
		t.Errorf("starter item name = %q, want Rusty Lantern", got)
	}
	if queries.createRelationshipCalls != 1 {
		t.Fatalf("CreateRelationship calls = %d, want 1", queries.createRelationshipCalls)
	}
	if got := queries.lastCreateRelationshipParams.TargetEntityID; got != dbutil.ToPgtype(npcID) {
		t.Errorf("starter relationship target id = %#v, want %#v", got, dbutil.ToPgtype(npcID))
	}
```

- [ ] **Step 3: Add invalid relationship target test**

Add a new test `TestBuildWorld_InvalidSpawnRelationshipTargetReturnsBadRequest` that posts a valid world build request with:

```go
SpawnPackage: &api.CharacterSpawnPackage{
	Relationships: []api.StarterRelationship{{TargetEntityType: "npc", TargetEntityID: "not-a-uuid", RelationshipType: "knows"}},
},
```

Expected assertion:

```go
if rec.Code != http.StatusBadRequest {
	t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
}
```

- [ ] **Step 4: Run handler tests**

Run: `go test ./internal/handlers -run 'TestBuildWorld' -v`

Expected: PASS.

---

### Task 5: Full validation and follow-up notes

**Files:**
- No required code edits unless tests expose issues.

- [ ] **Step 1: Run focused tests**

Run:

```bash
go test ./internal/world ./internal/handlers ./pkg/api
```

Expected: PASS.

- [ ] **Step 2: Run broader tests if focused suite passes**

Run:

```bash
go test ./...
```

Expected: PASS or documented unrelated failures.

- [ ] **Step 3: Manual API smoke payload**

Use the JSON contract above against `/api/v1/campaigns/start/world` in a dev run if the app server is already running. Verify:

- `/api/v1/campaigns/{id}/characters/{playerID}/inventory` includes starter items.
- `/api/v1/campaigns/{id}/facts` includes starter known facts.
- `/api/v1/campaigns/{id}/relationships` includes starter relationships if relationship listing filters include player-aware/default rows.

- [ ] **Step 4: Record follow-up**

If codex spawning is still desired, create a follow-up issue/plan: "Codex spawn unlocks for generated languages/cultures/beliefs/economies." That design needs a decision about referencing existing generated codex records versus creating new records during spawn.

---

## Self-Review

- Spec coverage: covers starter inventory, starting knowledge as player-known facts, NPC/relationship hooks, and leaves codex unlocks as a follow-up because current codex tables need concrete generated entity references.
- Placeholder scan: no TBD/TODO placeholders are required for implementation.
- Type consistency: API `CharacterSpawnPackage` converts to world `CharacterSpawnPackage`; world layer depends on `statedb.Querier`; relationships use UUIDs internally and strings at JSON boundary.
