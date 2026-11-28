package game

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/google/uuid"
)

// TravelOption is one known destination reached by one connection from the
// actor's location, with the plan derived from its declared route facts.
type TravelOption struct {
	RouteID uuid.UUID
	ToID    uuid.UUID
	To      string
	Plan    *srd521.TravelPlan
	Reason  string
}

// TravelOption returns only an available destination in this snapshot.
func (s *StoredAction) TravelOption(id uuid.UUID) *TravelOption {
	for i := range s.Travel {
		if s.Travel[i].ToID == id && s.Travel[i].Reason == "" && s.Travel[i].Plan != nil {
			return &s.Travel[i]
		}
	}
	return nil
}

// ReadTravelDay reads stats.srd521_travel. A missing record is no travel yet;
// a malformed one fails closed.
func ReadTravelDay(stats []byte) (*srd521.TravelDay, error) {
	var envelope struct {
		Travel *srd521.TravelDay `json:"srd521_travel"`
	}
	if err := json.Unmarshal(stats, &envelope); err != nil {
		return nil, err
	}
	if t := envelope.Travel; t != nil && (t.Day < 1 || t.Seconds < 0 || t.Seconds > srd521.TravelDaySeconds) {
		return nil, errors.New("invalid travel record")
	}
	return envelope.Travel, nil
}

type travelRow struct {
	RouteID uuid.UUID
	ToID    uuid.UUID
	To      string
	Facts   []byte
}

// readTravelRoutes lists connections from the actor's location to locations
// the player knows or has visited. Unknown destinations are never listed.
func readTravelRoutes(ctx context.Context, database statedb.DBTX, campaign, from uuid.UUID, lock bool) ([]travelRow, error) {
	locking := ""
	if lock {
		locking = " FOR UPDATE OF lc"
	}
	rows, err := database.Query(ctx, `SELECT lc.id,l.id,l.name,lc.srd521_route FROM location_connections lc JOIN locations l ON l.campaign_id=lc.campaign_id AND l.id=CASE WHEN lc.from_location_id=$2 THEN lc.to_location_id ELSE lc.from_location_id END WHERE lc.campaign_id=$1 AND (lc.from_location_id=$2 OR (lc.bidirectional AND lc.to_location_id=$2)) AND l.id<>$2 AND (l.player_known OR l.player_visited) ORDER BY l.id,lc.id`+locking, dbutil.ToPgtype(campaign), dbutil.ToPgtype(from))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []travelRow{}
	for rows.Next() {
		var row travelRow
		if err := rows.Scan(&row.RouteID, &row.ToID, &row.To, &row.Facts); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

func travelLabel(option TravelOption) string {
	plan := option.Plan
	if plan == nil {
		return option.To
	}
	road := "no road"
	if plan.GoodRoad {
		road = "good road"
	}
	label := fmt.Sprintf("%s: %s (%d feet), %s, %s; maximum pace %s", option.To, travelMiles(plan.DistanceFeet), plan.DistanceFeet, plan.Terrain, road, plan.MaximumPace)
	if option.Reason == "" {
		label += fmt.Sprintf("; %s at a %s pace", travelDuration(plan.Seconds), plan.Pace)
	}
	return label
}

func travelMiles(feet int) string {
	return fmt.Sprintf("%.2f miles", float64(feet)/5280)
}

func travelDuration(seconds int) string {
	return fmt.Sprintf("%d min %02d s", seconds/60, seconds%60)
}

func prepareTravel(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, intent ActionIntent, lock bool) (*StoredAction, error) {
	p, err := PrepareCharacterAction(ctx, campaign, pc, ActionIntent{ActorID: intent.ActorID, Kind: srd521.AbilityCheck, Ability: "wisdom"})
	if err != nil {
		return nil, err
	}
	p.Intent = intent
	result := &StoredAction{Preview: p}
	if !p.Available {
		return result, nil
	}
	p.Available, p.Modifier, p.Options = false, nil, []string{}
	p.Reason = "Choose a fast, normal or slow pace for one hour of travel."
	if intent.Ability != "" || intent.Skill != "" || intent.SlotLevel != 0 {
		return result, nil
	}
	if pc.CurrentLocationID == nil {
		p.Reason = "Travel needs a current location."
		return result, nil
	}
	if result.Clock, err = ReadSpellClock(ctx, database, campaign.ID); err != nil {
		return nil, err
	}
	encounter, err := ReadEncounter(ctx, database, campaign.ID)
	if err != nil {
		return nil, err
	}
	rows, err := readTravelRoutes(ctx, database, campaign.ID, *pc.CurrentLocationID, lock)
	if err != nil {
		return nil, err
	}
	encoded, err := json.Marshal([]any{p.ExpectedRevision, intent, result.Clock, encounter, rows})
	if err != nil {
		return nil, err
	}
	p.ExpectedRevision = fmt.Sprintf("travel-v1:%x", sha256.Sum256(encoded))
	if encounter != nil && encounter.Status == "active" {
		p.Reason = "Travel is unavailable during an encounter."
		return result, nil
	}
	if _, pending, _ := readDeclaredFall(pc.Stats); pending {
		p.Reason = "Resolve the declared fall before travelling."
		return result, nil
	}
	if ReadConditions(pc.Stats).Prone {
		p.Reason = "Stand up from Prone before travelling."
		return result, nil
	}
	pace := intent.EntryID
	if pace != srd521.PaceFast && pace != srd521.PaceNormal && pace != srd521.PaceSlow {
		return result, nil
	}
	day, err := ReadTravelDay(pc.Stats)
	if err != nil {
		p.Reason = "This character's travel record is unreadable."
		return result, nil
	}
	travelled := day.TravelledOn(result.Clock)
	count := map[uuid.UUID]int{}
	for _, row := range rows {
		count[row.ToID]++
	}
	for _, row := range rows {
		option := TravelOption{RouteID: row.RouteID, ToID: row.ToID, To: row.To}
		var route srd521.TravelRoute
		switch {
		case count[row.ToID] > 1:
			option.Reason = "More than one route leads there; the route is ambiguous."
		case len(row.Facts) == 0:
			option.Reason = "Route distance, terrain and road are not declared; travel time prose is not an SRD distance."
		case decodeTravelRoute(row.Facts, &route) != nil:
			option.Reason = "Route facts include unknown or malformed fields, such as mounts or vehicles, which are not supported."
		default:
			plan, err := srd521.PlanTravel(route, pace)
			option.Plan = plan
			if err != nil {
				option.Reason = err.Error()
				var unresolved srd521.TravelUnresolved
				if errors.As(err, &unresolved) {
					option.Plan = nil
				}
			} else if travelled+plan.Seconds > srd521.TravelDaySeconds {
				option.Reason = fmt.Sprintf("Travel beyond 8 hours a day requires Extended Travel saving throws, which are not supported; %s travelled today.", travelDuration(travelled))
			}
		}
		if containsTravel(result.Travel, row.ToID) {
			continue // An ambiguous destination is listed once.
		}
		result.Travel = append(result.Travel, option)
		p.Targets = append(p.Targets, ActionTarget{ID: option.ToID, Label: travelLabel(option), UnavailableReason: option.Reason})
		if option.Reason == "" {
			p.Available = true
		}
	}
	if len(rows) == 0 {
		p.Reason = "No known connected route leads from here."
		return result, nil
	}
	if !p.Available {
		p.Reason = "No known route can be travelled within one hour at a " + pace + " pace."
		return result, nil
	}
	p.Reason = ""
	p.Options = []string{fmt.Sprintf("Travel one segment of at most one hour at a %s pace (%d miles an hour). Pace effect while travelling: %s. No checks are made during the segment; navigation, weather, food and water, forced march and random encounters are not modelled.", pace, srd521.TravelMilesPerHour(pace), strings.Join(srd521.TravelPaceEffects(pace), "; "))}
	p.CostStatus = "fixture-only"
	p.SourceIDs = srd521.TravelSourceIDs()
	return result, nil
}

// decodeTravelRoute rejects unknown fields so an undeclared mechanic, such as
// a mount, vehicle or flying party, is never silently ignored.
func decodeTravelRoute(data []byte, route *srd521.TravelRoute) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode(route)
}

func containsTravel(options []TravelOption, id uuid.UUID) bool {
	for _, option := range options {
		if option.ToID == id {
			return true
		}
	}
	return false
}
