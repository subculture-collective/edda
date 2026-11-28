-- name: GetCampaignTime :one
SELECT campaign_id, day, hour, minute, updated_at, second
FROM campaign_time
WHERE campaign_id = $1;

-- name: UpsertCampaignTime :one
-- Sets the time to the start of the given minute, so a stale second from the
-- SRD rules clock does not shift it.
INSERT INTO campaign_time (campaign_id, day, hour, minute, second)
VALUES ($1, $2, $3, $4, 0)
ON CONFLICT (campaign_id) DO UPDATE
SET day = EXCLUDED.day, hour = EXCLUDED.hour, minute = EXCLUDED.minute, second = 0, updated_at = now()
RETURNING campaign_id, day, hour, minute, updated_at, second;
