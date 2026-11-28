# Characters and campaigns saved by the broken creation forms

Until PR #2242 (merged 2026-11-28 13:46 UTC), both TUI creation forms ignored
the player's choices. A TUI binary built before that change still does.

- The **guided character form** saved an empty name, race Human, class
  Barbarian, background Acolyte and alignment Lawful Good.
- The **campaign attributes form** always requested Fantasy, Open Wilderness and
  Gritty and Dark.

## Find affected rows

[blank-character-names.sql](blank-character-names.sql) holds two `SELECT`
statements and changes nothing. To run it in a read-only session, with
`DATABASE_URL` pointing at the database:

```sh
PGOPTIONS='-c default_transaction_read_only=on' \
  psql "$DATABASE_URL" -f docs/data-checks/blank-character-names.sql
```

Save the output before players rename their characters. The queries find the
guided form's rows by the blank name, so a renamed character drops out.

1. **Characters.** Every row has a blank name. `guided_form_defaults = true`
   means the description, concept and alignment also match what the broken form
   saved. A blank-named row without that match came from another path, such as
   an API request for a legacy campaign.
2. **Campaigns.** These are candidates only. A campaign row does not record
   whether the attributes form made it or what the form requested. The model
   chose the saved genre and tone, and setting style is not saved. A player who
   chose the defaults looks the same as one whose choice was dropped, and a
   world whose genre or tone the model worded differently is missed.

## Repair

Rows are flagged, not rewritten. Nothing in the database holds what the player
chose, so a bulk update could only guess.

- **Name:** when the TUI loads a campaign that has a blank-named character, it
  asks the player to name it before play. Escape returns to the campaign list.
  The web app has no such prompt.
- **Race, class, background and alignment:** these cannot be recovered. The
  prompt says that the guided form replaced them. A player who wants other
  choices has to start a new campaign. These are legacy-rules characters, so the
  choices are in the description, stats and abilities text.
- **Campaign settings:** leave them. The world was already generated from the
  request. Tell the owner if the output shows a campaign they ask about.

## No database constraint yet

The server now refuses a blank name before it creates anything. There is no
`CHECK` constraint on `player_characters.name`:

- A validated constraint fails to apply to a database that still has a blank
  name, which would stop the migration.
- A `NOT VALID` constraint applies, but PostgreSQL still checks it on every
  later `UPDATE` of a blank-named row. Hit point and stat writes for an
  unrenamed character would then fail.

Add the constraint in a new migration after query 1 returns no rows on every
database.
