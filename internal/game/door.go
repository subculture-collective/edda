package game

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/google/uuid"
)

// KnownDoor is one door the actor knows at their location, read from the
// server-owned npcs.properties.srd521_door facts. Err is a player-safe reason
// when the facts are not an ordinary supported door.
type KnownDoor struct {
	ID    uuid.UUID
	Name  string
	Raw   []byte
	Facts srd521.DoorFacts
	Err   error
}

// readKnownDoors returns the doors at the actor's location that the campaign
// history has involved and whose facts name this actor and are not hidden.
// Hidden and foreign doors are skipped, so callers cannot tell them from none.
func readKnownDoors(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, lock bool) ([]KnownDoor, error) {
	doors := []KnownDoor{}
	if pc.CurrentLocationID == nil {
		return doors, nil
	}
	locking := ""
	if lock {
		locking = " FOR UPDATE OF n"
	}
	rows, err := database.Query(ctx, `SELECT n.id,n.name,n.properties->'srd521_door' FROM npcs n WHERE n.campaign_id=$1 AND n.location_id=$2 AND n.properties ? 'srd521_door' AND EXISTS(SELECT 1 FROM session_logs l WHERE l.campaign_id=n.campaign_id AND n.id=ANY(l.npcs_involved)) ORDER BY n.id`+locking, dbutil.ToPgtype(campaign.ID), *pc.CurrentLocationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var door KnownDoor
		if err := rows.Scan(&door.ID, &door.Name, &door.Raw); err != nil {
			return nil, err
		}
		if !srd521.DoorFor(door.Raw, pc.ID.String()) {
			continue
		}
		door.Facts, door.Err = srd521.ParseDoorFacts(door.Raw)
		doors = append(doors, door)
	}
	return doors, rows.Err()
}

// applyDoors resolves each target's declared door. A door that is unknown,
// hidden, foreign or unsupported leaves the target's sight unresolved.
func applyDoors(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, targets []MartialTarget, lock bool) error {
	var doors []KnownDoor
	for i := range targets {
		space := &targets[i].Space
		if space.Door == "" {
			continue
		}
		if doors == nil {
			var err error
			if doors, err = readKnownDoors(ctx, database, campaign, pc, lock); err != nil {
				return err
			}
		}
		state, raw := "", []byte(nil)
		for _, door := range doors {
			if door.ID.String() == space.Door && door.Err == nil {
				state, raw = door.Facts.State, door.Raw
			}
		}
		revision := sha256.Sum256(append([]byte(space.Revision+"\x00"), raw...))
		*space = space.ThroughDoor(state)
		space.Revision = fmt.Sprintf("%x", revision)
	}
	return nil
}

// unknownDoor is one reason for unknown, hidden, foreign and absent doors so a
// preview never confirms that a door exists.
const unknownDoor = "Choose one known door at your location."

func prepareOpenDoor(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, intent ActionIntent, lock bool) (*StoredAction, error) {
	p, err := PrepareCharacterAction(ctx, campaign, pc, ActionIntent{ActorID: intent.ActorID, Kind: srd521.AbilityCheck, Ability: "strength"})
	if err != nil {
		return nil, err
	}
	p.Intent = intent
	result := &StoredAction{Preview: p}
	if !p.Available {
		return result, nil
	}
	p.Available, p.Modifier, p.Reason = false, nil, unknownDoor
	id, err := uuid.Parse(intent.EntryID)
	if err != nil || id == uuid.Nil || intent.Ability != "" || intent.Skill != "" || intent.SlotLevel != 0 {
		return result, nil
	}
	doors, err := readKnownDoors(ctx, database, campaign, pc, lock)
	if err != nil {
		return nil, err
	}
	var door *KnownDoor
	for i := range doors {
		if doors[i].ID == id {
			door = &doors[i]
		}
	}
	if door == nil {
		return result, nil
	}
	encounter, err := ReadEncounter(ctx, database, campaign.ID)
	if err != nil {
		return nil, err
	}
	encoded, err := json.Marshal([]any{p.ExpectedRevision, door.ID, door.Name, door.Raw, encounter})
	if err != nil {
		return nil, err
	}
	p.ExpectedRevision = fmt.Sprintf("door-v1:%x", sha256.Sum256(encoded))
	target := ActionTarget{ID: door.ID, Label: door.Name, ReachFeet: srd521.DoorReachFeet}
	if door.Err == nil {
		target.DistanceFeet = door.Facts.DistanceFeet
	}
	p.Targets = []ActionTarget{target}
	p.SourceIDs = append([]string(nil), srd521.DoorSourceIDs...)
	if door.Err != nil {
		p.Reason = door.Err.Error()
		return result, nil
	}
	if encounter != nil && encounter.Status == "active" {
		p.Reason = "Opening a door during an encounter is not supported yet; the one free object interaction per turn is not tracked."
		if encounter.Reaction != nil && encounter.Reaction.Pending != nil {
			p.Reason = "Choose or decline the pending reaction before continuing."
		}
		return result, nil
	}
	if _, _, err := srd521.OpenDoorFacts(door.Facts, door.ID.String(), door.Name); err != nil {
		p.Reason = err.Error()
		return result, nil
	}
	result.Door = door
	p.Available, p.Reason = true, ""
	p.CostStatus = "fixture-only"
	p.Options = []string{fmt.Sprintf("Open %s (closed, unlocked, %d feet away, within your %d-foot reach). Outside combat this takes no action and no campaign time. Targets behind the door are re-evaluated with it open.", door.Name, *door.Facts.DistanceFeet, srd521.DoorReachFeet)}
	return result, nil
}
