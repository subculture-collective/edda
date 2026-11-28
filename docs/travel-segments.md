# One-hour travel segments (#144)

A character can travel one segment of at most one hour along a known
connected route at an explicit Fast, Normal or Slow pace. Confirmation moves
the character and advances the campaign clock once, in the same transaction
as the receipt and the session log. Source: SRD 5.2.1 p. 12 (Travel Pace) and
p. 192 (Travel Pace, Travel Terrain).

## Intent

| Intent | Target | Effect |
| --- | --- | --- |
| `travel-hour`, `entry_id: fast \| normal \| slow` | one offered destination location | Moves the actor, marks the destination visited, advances the clock by the exact travel time and records `srd_travel`. |

The preview lists each connection from the actor's location to a location the
player knows or has visited. Unknown destinations are never listed. Each
target shows distance, terrain, road, maximum pace and travel time, or the
reason it is unavailable. The preview revision binds the pace and the route
facts, so a changed pace, a changed route or an earlier move returns
`ErrActionRefresh` without changing anything. The web API and the TUI share
`game.PrepareStoredAction` and `engine.ConfirmAction`; there is no web or TUI
control yet.

## Route facts

Route facts live in `location_connections.srd521_route` (migration 00045). No
tool writes this column. The legacy `travel_time` text is never read as a
distance. Every key is required:

```json
{"profile":"travel-route-v1","distance_feet":15840,"terrain":"forest","good_road":false}
```

- `terrain` is one Travel Terrain entry: `arctic`, `coastal`, `desert`,
  `forest`, `grassland`, `hill`, `mountain`, `swamp`, `underdark` or `urban`.
  `waterborne` needs a vehicle and is unavailable.
- Unknown keys, a missing key, unknown terrain or a distance of zero or less
  make the route unavailable.
- Two connections to the same destination make it ambiguous and unavailable.
- A bidirectional connection uses the same facts in both directions.

## Rules

- One hour covers 4 miles at Fast, 3 at Normal and 2 at Slow. A route longer
  than that at the chosen pace is unavailable; travel time is the distance at
  that rate, rounded up to whole seconds.
- Terrain sets the maximum pace. A good road raises it one step. A Fast pace
  in arctic terrain needs equipment such as skis, which is not modelled, so it
  is unavailable.
- The pace effect is shown in the preview and receipt: Fast gives
  Disadvantage on Wisdom (Perception or Survival) and Dexterity (Stealth)
  checks, Normal on Dexterity (Stealth) checks, and Slow gives Advantage on
  Wisdom (Perception or Survival) checks. No checks are made during a segment.
- `stats.srd521_travel = {"day": N, "seconds": S}` counts travel on the
  current campaign day. A segment that would pass 8 hours is unavailable,
  because Extended Travel saving throws are not modelled. A malformed record
  makes travel unavailable.
- Travel is unavailable during an active encounter, while the actor has a
  saved Prone condition (#147), or while a declared fall is unresolved.

## Durable effects

Travel advances the clock through the same `campaign_time` write as waits and
rests. Light and concentration keep their absolute `expires_at`. A
concentration effect that has ended by the arrival second is removed and named
in the receipt, as a wait or rest does.

## Not supported

Forced march and Extended Travel, mounts, vehicles, flying or special
movement, slower travelers in a group, navigation, foraging, weather, food and
water, and random encounters. The legacy `move_player` tool is unchanged.
