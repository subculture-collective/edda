# Declared fall (#147)

One fall, declared by the server, resolves outside encounters with damage and
Prone the player can inspect. Source: SRD 5.2.1 Rules Glossary, "Falling
[Hazard]", p. 182 of the pinned PDF (see
[conformance](srd-5.2.1-conformance.md)). This is not Full Rules support.

## Declaration

The fall is a server-owned ruling on the actor,
`player_characters.stats.srd521_fall`:

```json
{"id": "fall-1", "height_feet": 35, "landing": "solid",
 "flying": false, "onto_creature": false, "fall_prevention": false}
```

Every field is required. `landing` is `solid` or `water`. The player-stat tool
only changes existing integer stats, so the model cannot write this object.

These fail closed with a reason and no roll:

- a missing or negative height, or a missing landing;
- any other landing, such as acid or lava;
- unknown flying, creature or magic facts, or any of them `true`: falls out of
  flight, onto a creature and with Feather Fall, Slow Fall or similar magic
  are excluded;
- any field the slice does not know;
- an active encounter, or an actor at 0 HP.

## Preview and confirmation

The intent is `fall-resolve` with `entry_id` set to the declaration's `id`.
On solid ground `skill` must be empty. For water it may be `athletics` or
`acrobatics` to use the Reaction for the DC 15 entry check, or empty to
decline it.

The preview states the dice: one d6 per full 10 feet, so a fall under 10 feet
deals no damage, capped at 20d6. The preview revision covers the declaration,
the actor snapshot and the saved conditions, so a changed height or a
replaced declaration refuses an open preview.

Confirmation rolls the optional check, then the falling dice. A success halves
the damage, rounding down. In one transaction it:

- removes the declaration, so a second confirmation finds nothing to resolve;
- lowers HP, stopping at 0;
- saves Prone when any damage was taken;
- makes a caster's Concentration save;
- records the `srd_fall` receipt (height, landing, dice, rolled total, halving,
  damage, HP before and after, Prone before and after, the water check) and the
  session log.

A replay returns the saved result without rolling. At 0 HP the narrative says
that falling Unconscious and death saving throws require adjudication.

## Where Prone lives

Outside an encounter Prone is `stats.srd521_conditions.prone`. `combat-stand`
outside an encounter ends it without a movement cost, because no movement
budget is tracked there.

When an encounter starts, the player actor takes the saved Prone and the
actor's saved value is cleared, so only the encounter holds it. Crawling,
standing for half speed and Disadvantage on attacks then follow the
[Beast encounter](beast-encounter.md) rules. When the encounter completes, a
Prone player is saved back to `stats.srd521_conditions`.

## Player controls

Clients cannot declare a fall or change one. While a declaration has an `id`,
the character sheet carries it as `sheet.fall`. `height_feet` and `landing` are
included only when the declaration passes the checks above.

- **Web.** The Character tab shows the waiting fall, and the action menu offers
  to resolve it. For water it also offers each Reaction check. The preview,
  confirm and retry flow is the one other actions use, and the `srd_fall`
  receipt shows in the action panel and the narrative feed.
- **TUI.** `/fall` previews the declared fall; `/fall athletics` and
  `/fall acrobatics` preview it with the water Reaction. `/confirm-fall`
  commits and retries with the same action ID. `/cancel-fall` drops the preview
  and leaves the fall declared. The commands act for the campaign owner's first
  active character, and the confirmation narrative lists the dice.
