package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ErrConcentrationCaster reports an event for a caster outside the campaign.
var ErrConcentrationCaster = errors.New("concentration caster is not in this campaign")

// ApplyConcentrationEvents is the trusted adapter that ends or tests the exact
// concentration cast after server-resolved damage or condition events. Callers
// pass events built from results they already committed in tx (NPC attacks,
// falls, and later #487's Incapacitated/Unconscious transitions); nothing
// here reads client input. hpBefore is the caster's HP before the first
// damage event. It locks the caster row and the caster's cast rows, saves the
// caster's casting state and HP, and removes only that caster's recipient
// links for an ended cast, all inside tx. Unknown events fail closed with no
// change and no roll.
func ApplyConcentrationEvents(ctx context.Context, tx pgx.Tx, campaignID, actorID uuid.UUID, hpBefore int, events []srd521.ConcentrationEvent, dice srd521.Dice) ([]*srd521.EffectReceipt, error) {
	var stats []byte
	var level, maxHP, xp int
	err := tx.QueryRow(ctx, `SELECT stats,level,max_hp,experience FROM player_characters WHERE id=$1 AND campaign_id=$2 FOR UPDATE`, actorID, campaignID).Scan(&stats, &level, &maxHP, &xp)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrConcentrationCaster
	}
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `SELECT 1 FROM concentration_casts WHERE campaign_id=$1 AND caster_id=$2 FOR UPDATE`, campaignID, actorID); err != nil {
		return nil, err
	}
	var envelope struct {
		Build srd521.Build `json:"srd521"`
	}
	if err := json.Unmarshal(stats, &envelope); err != nil {
		return nil, err
	}
	sheet, err := srd521.Project(stats, level, hpBefore, maxHP, xp)
	if err != nil {
		return nil, fmt.Errorf("concentration caster: %w", err)
	}
	now, err := game.ReadSpellClock(ctx, tx, campaignID)
	if err != nil {
		return nil, err
	}
	receipts, err := srd521.ResolveConcentrationEvents(&envelope.Build, sheet, events, now, dice)
	if err != nil || len(receipts) == 0 {
		return nil, err
	}
	if err := saveCaster(ctx, tx, campaignID, actorID, &envelope.Build, receipts[len(receipts)-1].HPAfter); err != nil {
		return nil, err
	}
	return receipts, nil
}
