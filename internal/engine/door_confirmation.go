package engine

import (
	"context"

	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

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
