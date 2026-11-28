# GM rulings

Some SRD 5.2.1 facts are rulings the game master makes, not facts the model
may invent: the distance, terrain and road of a travel route, and a fall. The
campaign owner acts as game master and declares them through two HTTP routes.
No model tool writes them.

| Route | Body | Writes |
| --- | --- | --- |
| `PUT /campaigns/{id}/locations/{from}/routes/{to}` | the [route facts](travel-segments.md#route-facts) object | `location_connections.srd521_route` on the one connection from `from` to `to`, or a bidirectional connection from `to` to `from` |
| `POST /campaigns/{id}/character/fall` | `actor_id` and the [declared fall](declared-fall.md#declaration) facts without `id` | `player_characters.stats.srd521_fall`, with an `id` the server assigns |

## Checks

- Both routes sit behind `RequireOwnership`: another user's campaign returns
  404 and writes nothing. A location, connection or character outside the
  campaign also returns 404.
- Rulings apply only to SRD 5.2.1 campaigns; a legacy campaign returns 422.
- Facts are decoded strictly and checked with the same code the previews use
  (`game.ParseTravelRoute` and `game.DeclaredFall.Validate`). A refusal returns
  422 with the reason and writes nothing.
- Current state returns 409: two connections between the same locations, an
  unresolved fall already declared, an active encounter, a turn in progress or
  a turn that needs recovery.

## Audit

Each ruling takes the campaign turn lock and, in the same transaction as the
fact, adds a `meta` entry to the session log. Its text states the ruling, for
example "GM ruling: the route from Gate to Wood is 5280 feet of forest terrain
with no road.", and its tool call `gm_declare_route` or `gm_declare_fall`
carries the stored facts.

## Not supported

There is no web or TUI form yet; the routes are the interface. World creation
does not declare route facts, and the model cannot propose a ruling.
