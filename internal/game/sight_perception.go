package game

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/google/uuid"
)

// SightPerceptionEntryPrefix selects a sight-based Wisdom (Perception) check
// on one known target: EntryID is this prefix followed by the target's ID.
const SightPerceptionEntryPrefix = "sight-perception:"

// SightPerceptionTarget returns the target ID named by a sight check entry.
func SightPerceptionTarget(entryID string) (uuid.UUID, bool) {
	raw, ok := strings.CutPrefix(entryID, SightPerceptionEntryPrefix)
	if !ok {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(raw)
	return id, err == nil && id != uuid.Nil
}

// ReadSenses returns the actor's declared, server-owned senses from
// stats.srd521_senses. A missing or malformed declaration is unresolved (nil);
// it never defaults to "no Darkvision".
func ReadSenses(stats []byte) *combat.Senses {
	var envelope struct {
		Senses *combat.Senses `json:"srd521_senses"`
	}
	if json.Unmarshal(stats, &envelope) != nil || envelope.Senses == nil || envelope.Senses.DarkvisionFeet < 0 {
		return nil
	}
	return envelope.Senses
}

// cannotSee states a sight ruling without naming anything the player does
// not already know.
const cannotSee = "You can't see this target; a check that relies on sight automatically fails."

func prepareSightPerception(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, intent ActionIntent, lock bool) (*StoredAction, error) {
	p, err := PrepareCharacterAction(ctx, campaign, pc, ActionIntent{ActorID: intent.ActorID, Kind: srd521.AbilityCheck, Ability: "wisdom", Skill: "perception"})
	if err != nil {
		return nil, err
	}
	result := &StoredAction{Preview: p}
	p.Intent = intent
	if !p.Available {
		return result, nil
	}
	p.Available = false
	// The same reason covers unknown, hidden, absent and malformed targets so
	// the preview never confirms that an unknown actor exists.
	p.Reason = "Choose one known target for a sight-based Wisdom (Perception) check outside encounters."
	targetID, ok := SightPerceptionTarget(intent.EntryID)
	if intent.Kind != srd521.AbilityCheck || !ok || intent.Ability != "" || intent.Skill != "" || intent.SlotLevel != 0 {
		return result, nil
	}
	encounter, err := ReadEncounter(ctx, database, campaign.ID)
	if err != nil {
		return nil, err
	}
	if encounter != nil && encounter.Status == "active" {
		return result, nil
	}
	targets, err := readMartialTargets(ctx, database, campaign, pc, lock)
	if err != nil {
		return nil, err
	}
	var target *MartialTarget
	for i := range targets {
		if targets[i].ID == targetID {
			target = &targets[i]
		}
	}
	if target == nil {
		return result, nil
	}
	senses := ReadSenses(pc.Stats)
	clock, err := ReadSpellClock(ctx, database, campaign.ID)
	if err != nil {
		return nil, err
	}
	emitter, err := ReadLightEmitter(ctx, database, campaign.ID, pc, clock, lock)
	if err != nil {
		return nil, err
	}
	identity := []any{p.ExpectedRevision, target.ID, target.Name, target.Space, senses}
	if emitter != nil {
		identity = append(identity, emitter)
	}
	encoded, err := json.Marshal(identity)
	if err != nil {
		return nil, err
	}
	p.ExpectedRevision = fmt.Sprintf("sight-v1:%x", sha256.Sum256(encoded))
	p.Targets = []ActionTarget{{ID: target.ID, Label: target.Name, DistanceFeet: target.Space.DistanceFeet}}
	space, lit := illuminate(target.Space, emitter)
	sight, err := space.Sight(senses)
	if err != nil {
		p.Reason = err.Error()
		return result, nil
	}
	if !sight.Seen {
		p.Reason = cannotSee
		return result, nil
	}
	mode := "one ordinary d20"
	if sight.Circumstance == combat.SightDisadvantage {
		mode = "Disadvantage (lightly obscured)"
		sight.SourceIDs = append(sight.SourceIDs, "srd-5.2.1:content/d20-tests/advantage-disadvantage")
	}
	darkvision := ""
	if sight.Darkvision {
		darkvision = ", seen with Darkvision"
	}
	if lit {
		darkvision += fmt.Sprintf(", lit by your Light on %s", emitter.Item)
		sight.SourceIDs = append(sight.SourceIDs, "srd-5.2.1:content/spells/light")
	}
	result.Sight = &sight
	p.Available = true
	p.Reason = ""
	p.Options = []string{fmt.Sprintf("Sight-based Wisdom (Perception): declared %s light, %s obscurement%s; %s. The outcome requires adjudication; no scene effect or time cost.", sight.Light, sight.Obscurement, darkvision, mode)}
	p.SourceIDs = append(p.SourceIDs, sight.SourceIDs...)
	return result, nil
}
