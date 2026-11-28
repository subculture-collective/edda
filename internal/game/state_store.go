package game

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/db"
	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// ReviseWorldFactCommand revises a canonical world fact inside a campaign.
type ReviseWorldFactCommand = domain.ReviseWorldFactCommand

// ReviseWorldFactResult reports the new fact and propagation details.
type ReviseWorldFactResult = domain.ReviseWorldFactResult
type RevealLocationCommand = domain.RevealLocationCommand
type RevealLocationResult = domain.RevealLocationResult
type MovePlayerCommand = domain.MovePlayerCommand
type MovePlayerResult = domain.MovePlayerResult

// StateStore owns campaign-scoped game state mutations.
type StateStore struct {
	queries stateQuerier
	timeDB  db.DBTX
}

type stateQuerier interface {
	GetFactByID(ctx context.Context, id pgtype.UUID) (statedb.WorldFact, error)
	SupersedeFact(ctx context.Context, arg statedb.SupersedeFactParams) (statedb.SupersedeFactRow, error)
	GetLocationByID(ctx context.Context, arg statedb.GetLocationByIDParams) (statedb.Location, error)
	GetConnectionsFromLocation(ctx context.Context, arg statedb.GetConnectionsFromLocationParams) ([]statedb.GetConnectionsFromLocationRow, error)
	GetPlayerCharacterByID(ctx context.Context, id pgtype.UUID) (statedb.PlayerCharacter, error)
	UpdatePlayerLocation(ctx context.Context, arg statedb.UpdatePlayerLocationParams) (statedb.PlayerCharacter, error)
	SetLocationPlayerKnown(ctx context.Context, id pgtype.UUID) error
	SetLocationPlayerVisited(ctx context.Context, id pgtype.UUID) error
}

// NewStateStore creates a persistence module for game state mutations.
func NewStateStore(q stateQuerier, timeDB ...db.DBTX) *StateStore {
	var t db.DBTX
	if len(timeDB) > 0 {
		t = timeDB[0]
	}
	return &StateStore{queries: q, timeDB: t}
}

// ReviseWorldFact supersedes a fact and propagates player-known visibility.
func (s *StateStore) ReviseWorldFact(ctx context.Context, cmd ReviseWorldFactCommand) (*ReviseWorldFactResult, error) {
	if s == nil || s.queries == nil {
		return nil, errors.New("state store is required")
	}
	if cmd.CampaignID == uuid.Nil {
		return nil, errors.New("campaign_id is required")
	}
	if cmd.FactID == uuid.Nil {
		return nil, errors.New("fact_id is required")
	}
	if cmd.NewFact == "" {
		return nil, errors.New("new_fact is required")
	}

	oldFact, err := s.queries.GetFactByID(ctx, dbutil.ToPgtype(cmd.FactID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrWorldFactNotFound
		}
		return nil, fmt.Errorf("get world fact: %w", err)
	}
	if dbutil.FromPgtype(oldFact.CampaignID) != cmd.CampaignID {
		return nil, domain.ErrWorldFactNotFound
	}
	if oldFact.SupersededBy.Valid {
		return nil, domain.ErrWorldFactSuperseded
	}

	newFact, err := s.queries.SupersedeFact(ctx, statedb.SupersedeFactParams{
		OldFactID:  dbutil.ToPgtype(cmd.FactID),
		CampaignID: dbutil.ToPgtype(cmd.CampaignID),
		Fact:       cmd.NewFact,
		Category:   oldFact.Category,
		Source:     oldFact.Source,
		Reveal:     cmd.RevealToPlayer,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrWorldFactSuperseded
		}
		return nil, fmt.Errorf("supersede world fact: %w", err)
	}

	return &ReviseWorldFactResult{
		OldFact:               factToDomain(oldFact),
		NewFact:               factToDomain(statedb.WorldFact(newFact)),
		PlayerKnownPropagated: newFact.PlayerKnown,
	}, nil
}

// RevealLocation marks a campaign location as known to the player.
func (s *StateStore) RevealLocation(ctx context.Context, cmd RevealLocationCommand) (*RevealLocationResult, error) {
	if s == nil || s.queries == nil {
		return nil, errors.New("state store is required")
	}
	if cmd.CampaignID == uuid.Nil {
		return nil, errors.New("campaign_id is required")
	}
	if cmd.LocationID == uuid.Nil {
		return nil, errors.New("location_id is required")
	}

	location, err := s.queries.GetLocationByID(ctx, statedb.GetLocationByIDParams{ID: dbutil.ToPgtype(cmd.LocationID), CampaignID: dbutil.ToPgtype(cmd.CampaignID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("location not found in campaign")
		}
		return nil, fmt.Errorf("get location: %w", err)
	}
	if dbutil.FromPgtype(location.CampaignID) != cmd.CampaignID {
		return nil, domain.ErrWorldFactNotFound
	}
	if err := s.queries.SetLocationPlayerKnown(ctx, dbutil.ToPgtype(cmd.LocationID)); err != nil {
		return nil, fmt.Errorf("reveal location: %w", err)
	}
	return &RevealLocationResult{LocationID: cmd.LocationID, LocationName: location.Name}, nil
}

// MovePlayer moves a player character between connected locations.
func (s *StateStore) MovePlayer(ctx context.Context, cmd MovePlayerCommand) (*MovePlayerResult, error) {
	if s == nil || s.queries == nil {
		return nil, errors.New("state store is required")
	}
	if cmd.CampaignID == uuid.Nil {
		return nil, errors.New("campaign_id is required")
	}
	if cmd.PlayerCharacterID == uuid.Nil {
		return nil, errors.New("player_character_id is required")
	}
	if cmd.CurrentLocationID == uuid.Nil {
		return nil, errors.New("current_location_id is required")
	}
	if cmd.TargetLocationID == uuid.Nil {
		return nil, errors.New("target_location_id is required")
	}

	player, err := s.queries.GetPlayerCharacterByID(ctx, dbutil.ToPgtype(cmd.PlayerCharacterID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("player character not found in campaign")
		}
		return nil, fmt.Errorf("get player character: %w", err)
	}
	if dbutil.FromPgtype(player.CampaignID) != cmd.CampaignID {
		return nil, errors.New("player character not found in campaign")
	}

	currentLocation, err := s.queries.GetLocationByID(ctx, statedb.GetLocationByIDParams{ID: dbutil.ToPgtype(cmd.CurrentLocationID), CampaignID: db.ToPgUUID(cmd.CampaignID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("current location not found in campaign")
		}
		return nil, fmt.Errorf("get current location: %w", err)
	}
	targetLocation, err := s.queries.GetLocationByID(ctx, statedb.GetLocationByIDParams{ID: dbutil.ToPgtype(cmd.TargetLocationID), CampaignID: db.ToPgUUID(cmd.CampaignID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("target location not found in campaign")
		}
		return nil, fmt.Errorf("get target location: %w", err)
	}
	if dbutil.FromPgtype(currentLocation.CampaignID) != cmd.CampaignID || dbutil.FromPgtype(targetLocation.CampaignID) != cmd.CampaignID {
		return nil, errors.New("location does not belong to campaign")
	}
	connections, err := s.queries.GetConnectionsFromLocation(ctx, statedb.GetConnectionsFromLocationParams{CampaignID: db.ToPgUUID(cmd.CampaignID), LocationID: dbutil.ToPgtype(cmd.CurrentLocationID)})
	if err != nil {
		return nil, fmt.Errorf("get connections: %w", err)
	}
	connected := false
	var travelTime string
	for _, connection := range connections {
		if dbutil.FromPgtype(connection.ConnectedLocationID) == cmd.TargetLocationID {
			connected = true
			travelTime = connection.TravelTime.String
			break
		}
	}
	if !connected {
		return nil, errors.New("target location is not connected to current location")
	}
	if _, err := s.queries.UpdatePlayerLocation(ctx, statedb.UpdatePlayerLocationParams{CurrentLocationID: dbutil.ToPgtype(cmd.TargetLocationID), ID: dbutil.ToPgtype(cmd.PlayerCharacterID)}); err != nil {
		return nil, fmt.Errorf("update player location: %w", err)
	}

	result := &MovePlayerResult{PlayerCharacterID: cmd.PlayerCharacterID, FromLocationID: cmd.CurrentLocationID, ToLocationID: cmd.TargetLocationID, ToLocationName: targetLocation.Name, ToLocationDescription: targetLocation.Description.String, ToLocationType: targetLocation.LocationType.String, TravelTime: travelTime, VisitedMarked: true}
	if err := s.queries.SetLocationPlayerVisited(ctx, dbutil.ToPgtype(cmd.TargetLocationID)); err != nil {
		result.VisitedMarked = false
		result.VisitedWarning = err.Error()
	}
	if travelTime != "" {
		if dbTX := s.timeDB; dbTX != nil {
			advHours, advMinutes := parseTravelTime(travelTime)
			result.Day, result.Hour, result.Minute, err = db.AdvanceCampaignTime(ctx, dbTX, cmd.CampaignID, advHours, advMinutes)
			if err != nil {
				result.TimeWarning = err.Error()
			}
		}
	}
	return result, nil
}

func factToDomain(f statedb.WorldFact) domain.WorldFact {
	var supersededBy *uuid.UUID
	if f.SupersededBy.Valid {
		id := dbutil.FromPgtype(f.SupersededBy)
		supersededBy = &id
	}
	return domain.WorldFact{
		ID:           dbutil.FromPgtype(f.ID),
		CampaignID:   dbutil.FromPgtype(f.CampaignID),
		Fact:         f.Fact,
		Category:     f.Category,
		Source:       f.Source,
		SupersededBy: supersededBy,
		CreatedAt:    f.CreatedAt.Time,
	}
}

func parseTravelTime(s string) (hours, minutes int) {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return 1, 0
	}
	if strings.Contains(s, "hour") {
		parts := strings.FieldsFunc(s, func(r rune) bool { return r < '0' || r > '9' })
		if len(parts) > 0 {
			hours, _ = strconv.Atoi(parts[0])
		}
		if idx := strings.Index(s, "minute"); idx >= 0 {
			for _, p := range strings.FieldsFunc(s[idx-5:], func(r rune) bool { return r < '0' || r > '9' }) {
				if p != "" {
					minutes, _ = strconv.Atoi(p)
					break
				}
			}
		}
		if hours == 0 {
			hours = 1
		}
		return hours, minutes
	}
	if strings.Contains(s, "minute") || strings.Contains(s, "min") {
		parts := strings.FieldsFunc(s, func(r rune) bool { return r < '0' || r > '9' })
		if len(parts) > 0 {
			minutes, _ = strconv.Atoi(parts[0])
			return 0, minutes
		}
		return 0, 30
	}
	return 1, 0
}
