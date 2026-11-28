package db

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
)

// advanceCampaignTimeSQL adds $2 minutes to the campaign clock in one
// statement. A campaign without a clock row starts at day 1, 08:00.
//
// The SRD rules clock reads the time and later writes an absolute value while
// it holds the campaign row FOR UPDATE. Adding in SQL alone could land between
// the two and be overwritten, so this takes a conflicting campaign row lock
// before it touches campaign_time. That order also avoids a deadlock when the
// clock row does not exist yet. FOR NO KEY UPDATE still lets unrelated rows
// that reference the campaign be inserted.
//
// The update leaves campaign_time.second alone on purpose: adding whole
// minutes keeps the SRD rules clock's absolute-second expiries exact.
const advanceCampaignTimeSQL = `
WITH locked AS (
  SELECT id FROM campaigns WHERE id = $1 FOR NO KEY UPDATE
)
INSERT INTO campaign_time AS t (campaign_id, day, hour, minute)
SELECT id, 1 + (480 + $2::bigint) / 1440, (480 + $2::bigint) / 60 % 24, $2::bigint % 60 FROM locked
ON CONFLICT (campaign_id) DO UPDATE
SET day = t.day + (t.hour * 60 + t.minute + $2::bigint) / 1440,
    hour = (t.hour * 60 + t.minute + $2::bigint) / 60 % 24,
    minute = (t.minute + $2::bigint) % 60,
    updated_at = now()
RETURNING day, hour, minute
`

// AdvanceCampaignTime atomically adds non-negative hours and minutes to the
// campaign clock and returns the new time.
func AdvanceCampaignTime(ctx context.Context, database DBTX, campaignID uuid.UUID, hours, minutes int) (day, hour, minute int, err error) {
	// Larger values could wrap the minute total; the day column rejects far
	// smaller ones.
	if hours > math.MaxInt32 || minutes > math.MaxInt32 {
		return 0, 0, 0, fmt.Errorf("advance of %d hours %d minutes is too large", hours, minutes)
	}
	delta := int64(hours)*60 + int64(minutes)
	err = database.QueryRow(ctx, advanceCampaignTimeSQL, ToPgUUID(campaignID), delta).Scan(&day, &hour, &minute)
	return day, hour, minute, err
}

// LockCampaignForClock takes the campaign row lock that AdvanceCampaignTime
// takes. A legacy writer that changes a character row and then advances the
// clock in one transaction calls it first, so it locks the campaign before the
// character, the order the SRD rules actions use.
func LockCampaignForClock(ctx context.Context, database DBTX, campaignID uuid.UUID) error {
	var locked bool
	if err := database.QueryRow(ctx, `SELECT true FROM campaigns WHERE id = $1 FOR NO KEY UPDATE`, ToPgUUID(campaignID)).Scan(&locked); err != nil {
		return fmt.Errorf("lock campaign: %w", err)
	}
	return nil
}
