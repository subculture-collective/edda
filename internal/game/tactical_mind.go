package game

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// TacticalWindow is server-only adjudication. Only the preview is public.
type TacticalWindow struct {
	CheckID uuid.UUID
	Base    srd521.CheckReceipt
	DC      int
}

// OpenTacticalMindWindow is for trusted adjudication in the transaction that
// records the original check/history. No HTTP/model endpoint accepts a DC.
// Call after writing history and use a locked, current snapshot. Unknown
// outcomes from ordinary sheet checks never open this window automatically.
func OpenTacticalMindWindow(ctx context.Context, tx pgx.Tx, campaign domain.Campaign, pc domain.PlayerCharacter, checkID uuid.UUID, resolution *srd521.CheckResolution, dc int) error {
	if resolution == nil || resolution.Success() || checkID == uuid.Nil || pc.UserID == uuid.Nil || pc.UserID != campaign.CreatedBy {
		return fmt.Errorf("owned failed check required")
	}
	base := resolution.PlayerReceipt()
	if _, err := srd521.ResolveTacticalMind(base, dc, 1, "decline", nil); err != nil {
		return err
	}
	p, err := PrepareCharacterAction(ctx, campaign, pc, ActionIntent{ActorID: pc.ID, Kind: srd521.AbilityCheck, Ability: "strength"})
	if err != nil {
		return err
	}
	if !p.Available || pc.Level != 2 {
		return fmt.Errorf("supported Fighter 2 required")
	}
	data, err := json.Marshal(base)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO tactical_mind_windows(actor_id,owner_id,campaign_id,check_id,base_receipt,difficulty,revision,turn_number) VALUES($1,$2,$3,$4,$5,$6,$7,(SELECT COALESCE(MAX(turn_number),0) FROM session_logs WHERE campaign_id=$3)) ON CONFLICT(actor_id) DO UPDATE SET owner_id=EXCLUDED.owner_id,campaign_id=EXCLUDED.campaign_id,check_id=EXCLUDED.check_id,base_receipt=EXCLUDED.base_receipt,difficulty=EXCLUDED.difficulty,revision=EXCLUDED.revision,turn_number=EXCLUDED.turn_number,resolved=false`, pc.ID, pc.UserID, campaign.ID, checkID, data, dc, p.ExpectedRevision)
	return err
}

func prepareTacticalMind(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, intent ActionIntent, lock bool) (*StoredAction, error) {
	p, err := PrepareCharacterAction(ctx, campaign, pc, ActionIntent{ActorID: intent.ActorID, Kind: srd521.AbilityCheck, Ability: "strength"})
	if err != nil {
		return nil, err
	}
	result := &StoredAction{Preview: p}
	p.Intent = intent
	if !p.Available {
		return result, nil
	}
	p.Available = false
	p.Modifier = nil
	p.Reason = "Tactical Mind requires an owned, saved ordinary failed ability-check window and available Second Wind."
	if pc.Level != 2 || pc.UserID == uuid.Nil || pc.UserID != campaign.CreatedBy || intent.Ability != "" || intent.Skill != "" || intent.SlotLevel != 0 || (intent.EntryID != "use" && intent.EntryID != "decline") {
		return result, nil
	}
	var b struct {
		Build srd521.Build `json:"srd521"`
	}
	if err := json.Unmarshal(pc.Stats, &b); err != nil {
		return nil, err
	}
	if b.Build.Class != "fighter" || b.Build.Advancement == nil || b.Build.SecondWind < 1 {
		return result, nil
	}
	locking := ""
	if lock {
		locking = " FOR UPDATE"
	}
	var w TacticalWindow
	var data []byte
	err = database.QueryRow(ctx, `SELECT check_id,base_receipt,difficulty FROM tactical_mind_windows WHERE actor_id=$1 AND campaign_id=$2 AND owner_id=$3 AND NOT resolved AND revision=$4 AND turn_number=(SELECT COALESCE(MAX(turn_number),0) FROM session_logs WHERE campaign_id=$2)`+locking, pc.ID, campaign.ID, pc.UserID, p.ExpectedRevision).Scan(&w.CheckID, &data, &w.DC)
	if err == pgx.ErrNoRows {
		return result, nil
	}
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &w.Base); err != nil {
		return nil, err
	}
	if _, err = srd521.ResolveTacticalMind(w.Base, w.DC, b.Build.SecondWind, "decline", nil); err != nil {
		return result, nil
	}
	encoded, err := json.Marshal([]any{p.ExpectedRevision, w, intent})
	if err != nil {
		return nil, err
	}
	p.ExpectedRevision = fmt.Sprintf("tactical-v1:%x", sha256.Sum256(encoded))
	p.Available = true
	p.Reason = ""
	p.Options = []string{"Apply one d10 to the saved failed check; spend one Second Wind only on success. Decline closes the window without rolling or spending."}
	p.SourceIDs = []string{"srd-5.2.1:content/fighter/tactical-mind"}
	p.CostStatus = "conditional-on-success"
	result.Tactical = &w
	result.Build = &b.Build
	return result, nil
}
