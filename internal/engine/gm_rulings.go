package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/db"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GM rulings are server-owned facts that the campaign owner, acting as game
// master, declares through a trusted adapter. No model tool reaches them.

// ErrRulingNotFound hides whether a campaign, character or connection exists
// outside the caller's campaign.
var ErrRulingNotFound = errors.New("not found")

// RulingRefusal explains why a ruling was not recorded. Conflict marks a
// refusal caused by current state rather than by the declared facts.
type RulingRefusal struct {
	Reason   string
	Conflict bool
}

func (r RulingRefusal) Error() string { return r.Reason }

// RouteRuling is the route fact recorded on one connection.
type RouteRuling struct {
	FromID uuid.UUID          `json:"from_location_id"`
	ToID   uuid.UUID          `json:"to_location_id"`
	Route  srd521.TravelRoute `json:"route"`
}

// beginRuling takes the campaign turn lock and the campaign row, so a ruling
// never interleaves with a turn or confirmation. ownerID uuid.Nil skips the
// owner filter for the local trusted TUI.
func (e *Engine) beginRuling(ctx context.Context, campaignID, ownerID uuid.UUID) (pgx.Tx, func(), error) {
	if e.turnDB == nil {
		return nil, nil, ErrActionRefresh
	}
	release, err := e.acquireTurn(ctx, campaignID)
	if err != nil {
		return nil, nil, err
	}
	tx, err := e.turnDB.Begin(ctx)
	if err != nil {
		release()
		return nil, nil, err
	}
	done := func() { rollbackTurnTx(tx); release() }
	var owner any
	if ownerID != uuid.Nil {
		owner = db.ToPgUUID(ownerID)
	}
	var rules string
	err = tx.QueryRow(ctx, `SELECT rules_version FROM campaigns WHERE id=$1 AND ($2::uuid IS NULL OR created_by=$2) FOR UPDATE`, db.ToPgUUID(campaignID), owner).Scan(&rules)
	if errors.Is(err, pgx.ErrNoRows) {
		err = ErrRulingNotFound
	}
	if err == nil && rules != string(domain.RulesVersionSRD521) {
		err = RulingRefusal{Reason: "GM rulings apply only to SRD 5.2.1 campaigns."}
	}
	var pending bool
	if err == nil {
		err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM turn_attempts WHERE campaign_id=$1 AND status IN ('pending_mechanics','pending_narration'))`, db.ToPgUUID(campaignID)).Scan(&pending)
	}
	if err == nil && pending {
		err = ErrTurnRecoveryRequired
	}
	if err != nil {
		done()
		return nil, nil, err
	}
	return tx, done, nil
}

// recordRuling appends the ruling to the session log as a meta entry and
// commits it with the fact it records.
func recordRuling(ctx context.Context, tx pgx.Tx, campaignID uuid.UUID, tool string, args any, text string) error {
	encoded, err := json.Marshal(args)
	if err != nil {
		return err
	}
	calls, err := marshalAppliedToolCalls([]AppliedToolCall{{Tool: tool, Arguments: encoded}})
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO session_logs(campaign_id,turn_number,player_input,input_type,llm_response,tool_calls) SELECT $1,COALESCE(MAX(turn_number),0)+1,'GM ruling','meta',$2,$3 FROM session_logs WHERE campaign_id=$1`, db.ToPgUUID(campaignID), text, calls); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE campaigns SET updated_at=clock_timestamp() WHERE id=$1`, db.ToPgUUID(campaignID)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// DeclareRoute records travel route facts on the one connection that leads
// from one location to another, in either direction of a bidirectional
// connection. The facts pass the same strict parser the travel preview uses.
func (e *Engine) DeclareRoute(ctx context.Context, campaignID, ownerID, fromID, toID uuid.UUID, facts []byte) (*RouteRuling, error) {
	tx, done, err := e.beginRuling(ctx, campaignID, ownerID)
	if err != nil {
		return nil, err
	}
	defer done()
	route, err := game.ParseTravelRoute(facts)
	if err != nil {
		return nil, RulingRefusal{Reason: err.Error()}
	}
	rows, err := tx.Query(ctx, `SELECT lc.id,f.name,t.name FROM location_connections lc JOIN locations f ON f.campaign_id=lc.campaign_id AND f.id=$2 JOIN locations t ON t.campaign_id=lc.campaign_id AND t.id=$3 WHERE lc.campaign_id=$1 AND ((lc.from_location_id=$2 AND lc.to_location_id=$3) OR (lc.bidirectional AND lc.from_location_id=$3 AND lc.to_location_id=$2)) FOR UPDATE OF lc`, db.ToPgUUID(campaignID), db.ToPgUUID(fromID), db.ToPgUUID(toID))
	if err != nil {
		return nil, err
	}
	var ids []uuid.UUID
	var from, to string
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id, &from, &to); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	switch {
	case len(ids) == 0:
		return nil, ErrRulingNotFound
	case len(ids) > 1:
		return nil, RulingRefusal{Reason: "More than one connection leads there; the route is ambiguous.", Conflict: true}
	}
	encoded, err := json.Marshal(route)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE location_connections SET srd521_route=$1::jsonb WHERE id=$2 AND campaign_id=$3`, encoded, db.ToPgUUID(ids[0]), db.ToPgUUID(campaignID)); err != nil {
		return nil, err
	}
	road := "no road"
	if *route.GoodRoad {
		road = "a good road"
	}
	ruling := &RouteRuling{FromID: fromID, ToID: toID, Route: route}
	text := fmt.Sprintf("GM ruling: the route from %s to %s is %d feet of %s terrain with %s.", from, to, *route.DistanceFeet, route.Terrain, road)
	if err := recordRuling(ctx, tx, campaignID, "gm_declare_route", ruling, text); err != nil {
		return nil, err
	}
	return ruling, nil
}

// DeclareFall records a pending fall on the owner's character. The server
// assigns its ID; the facts must pass DeclaredFall.Validate, and an
// unresolved declaration is never replaced.
func (e *Engine) DeclareFall(ctx context.Context, campaignID, ownerID, actorID uuid.UUID, fall game.DeclaredFall) (*game.DeclaredFall, error) {
	tx, done, err := e.beginRuling(ctx, campaignID, ownerID)
	if err != nil {
		return nil, err
	}
	defer done()
	_, pc, err := game.ReadActionSnapshot(ctx, tx, campaignID, actorID, ownerID, true)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRulingNotFound
	}
	if err != nil {
		return nil, err
	}
	if fall.ID != "" {
		return nil, RulingRefusal{Reason: "The server assigns the fall id."}
	}
	if err := fall.Validate(); err != nil {
		return nil, RulingRefusal{Reason: err.Error()}
	}
	if pc.Status != "active" {
		return nil, RulingRefusal{Reason: "Only an active character can fall.", Conflict: true}
	}
	if game.HasDeclaredFall(pc.Stats) {
		return nil, RulingRefusal{Reason: "Resolve the declared fall before declaring another.", Conflict: true}
	}
	encounter, err := game.ReadEncounter(ctx, tx, campaignID)
	if err != nil {
		return nil, err
	}
	if encounter != nil && encounter.Status == "active" {
		return nil, RulingRefusal{Reason: "Falls during an encounter are not supported.", Conflict: true}
	}
	fall.ID = "fall-" + uuid.NewString()
	encoded, err := json.Marshal(fall)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521_fall}',$1::jsonb),updated_at=clock_timestamp() WHERE id=$2 AND campaign_id=$3`, encoded, db.ToPgUUID(actorID), db.ToPgUUID(campaignID)); err != nil {
		return nil, err
	}
	landing := map[string]string{srd521.FallSolid: "solid ground", srd521.FallWater: "water"}[fall.Landing]
	text := fmt.Sprintf("GM ruling: %s falls %d feet onto %s. Resolve the fall to roll its damage.", pc.Name, *fall.HeightFeet, landing)
	if err := recordRuling(ctx, tx, campaignID, "gm_declare_fall", fall, text); err != nil {
		return nil, err
	}
	return &fall, nil
}
