package game

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"slices"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/google/uuid"
)

// prepareDetectMagic previews the Magic Initiate free-use cast or the follow-up
// study. Preparation never writes; scene facts are fingerprinted so a changed
// ruling requires a fresh preview.
func prepareDetectMagic(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, intent ActionIntent, lock bool) (*StoredAction, error) {
	preview, err := PrepareCharacterAction(ctx, campaign, pc, ActionIntent{ActorID: intent.ActorID, Kind: srd521.AbilityCheck, Ability: "wisdom"})
	if err != nil {
		return nil, err
	}
	preview.Intent = intent
	result := &StoredAction{Preview: preview}
	if !preview.Available {
		return result, nil
	}
	preview.Available, preview.Modifier = false, nil
	study := intent.Kind == srd521.DetectMagicStudy
	preview.Reason = "Choose Detect Magic with the recorded Magic Initiate (Cleric) free use."
	if intent.Ability != "" || intent.Skill != "" || (study && intent.EntryID != "") {
		return result, nil
	}
	if intent.SlotLevel != 0 {
		preview.Reason = "Casting Detect Magic with a spell slot or as a Ritual is unavailable; only the Magic Initiate free use is supported."
		return result, nil
	}
	sheet, err := srd521.Project(pc.Stats, pc.Level, pc.HP, pc.MaxHP, pc.Experience)
	if err != nil || sheet.Spellcasting == nil {
		return result, nil
	}
	subjects, revision, err := readDetectMagicSubjects(ctx, database, campaign, pc, lock)
	if err != nil {
		return nil, err
	}
	result.Magic = subjects
	encoded, err := json.Marshal([]any{preview.ExpectedRevision, revision})
	if err != nil {
		return nil, err
	}
	preview.ExpectedRevision = fmt.Sprintf("detect-magic-v1:%x", sha256.Sum256(encoded))
	preview.CostStatus = "fixture-only"
	if study {
		preview.Options = []string{"Study with Detect Magic: Magic action; see a faint aura around each visible adjudicated creature or object in 30 feet that bears sensed magic, and learn the school of spell-created effects. No slot or feat use."}
		preview.Costs = []ActionCost{{ResourceID: "magic-action", Amount: 1}}
		preview.SourceIDs = []string{"srd-5.2.1:content/spells/detect-magic", "srd-5.2.1:content/rules-glossary/concentration"}
		now, err := ReadSpellClock(ctx, database, campaign.ID)
		if err != nil {
			return nil, err
		}
		c := sheet.Concentration
		if c != nil && c.StartedAt <= now && c.ExpiresAt > now {
			result.Concentration = c
		}
		switch {
		case c == nil || c.Spell != srd521.DetectMagicID || c.StartedAt > now || c.ExpiresAt <= now:
			preview.Reason = "Detect Magic concentration is required."
		case len(srd521.SensedDetectMagic(subjects)) == 0:
			preview.Reason = "No adjudicated magic is sensed within 30 feet."
		default:
			preview.Available, preview.Reason = true, ""
		}
		return result, nil
	}
	preview.Targets = []ActionTarget{{ID: pc.ID, Label: pc.Name + " (self; 30-foot sense)"}}
	preview.Options = []string{"Detect Magic; Magic Initiate (Cleric) free use, no spell slot; Magic action; verbal and somatic components; sense adjudicated magic within 30 feet through open paths (1 foot of stone, dirt or wood, 1 inch of metal or a thin sheet of lead blocks it); concentration up to 10 minutes. Replaces current concentration. Slot and Ritual casting unavailable."}
	preview.Costs = []ActionCost{{ResourceID: "magic-initiate", Amount: 1}, {ResourceID: "magic-action", Amount: 1}}
	preview.SourceIDs = []string{"srd-5.2.1:content/spells/detect-magic", "srd-5.2.1:content/feats/magic-initiate", "srd-5.2.1:content/rules-glossary/concentration"}
	choices := sheet.Spellcasting.Choices
	var spell *srd521.Spell
	for i := range sheet.Spellcasting.Spells {
		if sheet.Spellcasting.Spells[i].ID == srd521.DetectMagicID && sheet.Spellcasting.Spells[i].Playable {
			spell = &sheet.Spellcasting.Spells[i]
		}
	}
	switch {
	case choices == nil || spell == nil:
		preview.Reason = "No recorded Magic Initiate (Cleric) Detect Magic choice."
	case choices.FeatUses < 1:
		preview.Reason = "The Magic Initiate free use is spent; it returns after a Long Rest."
	case !spell.ComponentsAvailable:
		preview.Reason = "Detect Magic requires verbal and somatic components."
	default:
		preview.Available, preview.Reason = true, ""
	}
	return result, nil
}

// readDetectMagicSubjects reads only trusted srd521_magic facts on player-known
// NPCs at the caster's location that name this caster. Invalid or hidden facts
// are skipped and never disclosed.
func readDetectMagicSubjects(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, lock bool) ([]srd521.DetectMagicSubject, string, error) {
	subjects := []srd521.DetectMagicSubject{}
	if pc.CurrentLocationID == nil {
		return subjects, "no-location", nil
	}
	locking := ""
	if lock {
		locking = " FOR UPDATE OF n"
	}
	rows, err := database.Query(ctx, `SELECT n.id,n.name,n.properties->'srd521_magic' FROM npcs n WHERE n.campaign_id=$1 AND n.location_id=$2 AND n.properties ? 'srd521_magic' AND EXISTS(SELECT 1 FROM session_logs l WHERE l.campaign_id=n.campaign_id AND n.id=ANY(l.npcs_involved)) ORDER BY n.id`+locking, campaign.ID, *pc.CurrentLocationID)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	seen := []string{}
	for rows.Next() {
		var id uuid.UUID
		var name string
		var raw []byte
		if err := rows.Scan(&id, &name, &raw); err != nil {
			return nil, "", err
		}
		seen = append(seen, id.String(), name, string(raw))
		facts, ok := srd521.ParseDetectMagicFacts(raw)
		if !ok || facts.LocationID != pc.CurrentLocationID.String() || !slices.Contains(facts.ActorIDs, pc.ID.String()) {
			continue
		}
		subjects = append(subjects, srd521.DetectMagicSubject{ID: id.String(), Label: name, Facts: facts})
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}
	data, err := json.Marshal(seen)
	if err != nil {
		return nil, "", err
	}
	return subjects, fmt.Sprintf("%x", sha256.Sum256(data)), nil
}
