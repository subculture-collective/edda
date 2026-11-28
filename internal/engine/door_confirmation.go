package engine

import (
	"context"

	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PreviewDoor lists the doors the campaign owner's character knows and, for a
// 1-based number from that list, previews opening that door for the TUI. The
// preview is nil when number names no listed door.
func (e *Engine) PreviewDoor(ctx context.Context, campaignID uuid.UUID, number int) ([]srd521.SheetDoor, *game.ActionPreview, error) {
	if e.baseDB == nil {
		return nil, nil, ErrActionRefresh
	}
	c, pc, err := e.ownerSnapshot(ctx, campaignID)
	if err != nil {
		return nil, nil, err
	}
	doors, err := game.KnownDoorChoices(ctx, e.baseDB, c, pc)
	if err != nil || number < 1 || number > len(doors) {
		return doors, nil, err
	}
	prepared, err := game.PrepareStoredAction(ctx, e.baseDB, c, pc, game.ActionIntent{ActorID: pc.ID, Kind: srd521.OpenDoor, EntryID: doors[number-1].ID}, false)
	if err != nil {
		return nil, nil, err
	}
	return doors, prepared.Preview, nil
}

// confirmOpenDoor opens the prepared door under the caller's locks and
// transaction. It changes only srd521_door.state, and only from closed, so a
// concurrent change refreshes instead of writing twice. The receipt and
// history are written by the caller in the same transaction.
func confirmOpenDoor(ctx context.Context, tx pgx.Tx, campaign uuid.UUID, prepared *game.StoredAction) (*srd521.DoorReceipt, error) {
	door := prepared.Door
	if door == nil || door.Err != nil {
		return nil, ErrActionRefresh
	}
	_, receipt, err := srd521.OpenDoorFacts(door.Facts, door.ID.String(), door.Name)
	if err != nil {
		return nil, ErrActionRefresh
	}
	tag, err := tx.Exec(ctx, `UPDATE npcs SET properties=jsonb_set(properties,'{srd521_door,state}','"open"'),updated_at=clock_timestamp() WHERE id=$1 AND campaign_id=$2 AND properties->'srd521_door'->>'state'='closed'`, door.ID, campaign)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() != 1 {
		return nil, ErrActionRefresh
	}
	return receipt, nil
}
