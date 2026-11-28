-- Read-only audit for the creation-form bug fixed in PR #2242.
-- See blank-character-names.md before acting on the results.

-- 1. Characters with a blank name. guided_form_defaults is true when the rest
--    of the row is what the broken TUI guided form saved.
SELECT
  p.id AS character_id,
  p.campaign_id,
  c.name AS campaign_name,
  u.name AS owner_name,
  c.rules_version,
  p.created_at,
  COALESCE(
    p.name = ''
    AND p.description = 'human barbarian. Acolyte'
    AND p.stats->>'concept' = 'human barbarian'
    AND p.stats->>'personality' = 'Lawful Good',
    false
  ) AS guided_form_defaults
FROM player_characters p
JOIN campaigns c ON c.id = p.campaign_id
LEFT JOIN users u ON u.id = c.created_by
WHERE p.name !~ '[^[:space:]]'
ORDER BY p.created_at, p.id;

-- 2. Candidates only: campaigns of the TUI's local user whose saved genre and
--    tone resemble the attributes form's defaults (Fantasy, Open Wilderness,
--    Gritty and Dark). The model wrote these values and the request was not
--    saved, so a match does not show that the defaults replaced a choice.
SELECT
  c.id AS campaign_id,
  c.name AS campaign_name,
  c.genre,
  c.tone,
  c.world_type,
  c.created_at,
  EXISTS (
    SELECT 1 FROM player_characters p
    WHERE p.campaign_id = c.id
      AND p.name = ''
      AND p.description = 'human barbarian. Acolyte'
      AND p.stats->>'personality' = 'Lawful Good'
  ) AS has_guided_form_character
FROM campaigns c
JOIN users u ON u.id = c.created_by
WHERE u.name = 'Player'
  AND c.genre ILIKE '%fantasy%'
  AND c.tone ~* '(grit|grim|dark)'
ORDER BY c.created_at, c.id;
