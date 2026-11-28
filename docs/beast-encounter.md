# Explicit Wolf and Brown Bear encounter fixtures (#179, #180)

These fixtures run the Wolf (PDF p. 364) and Brown Bear (PDF p. 346) blocks of
the [pinned SRD 5.2.1](srd-5.2.1-conformance.md) in the existing two-actor
encounter, on the path the [Bandit fixture](bandit-encounter.md) uses. The
records are the existing typed transcriptions in `testdata/monsters_177.json`
and `testdata/monsters_155.json`, which keep the source attribution and
CC-BY-4.0 notice. This runtime adapts their mechanical fields into attack
resolution.

## Explicit instance contract

An NPC takes part only through `properties.srd521_target` with profile
`wolf-v1` or `brown-bear-v1`, the controlling player's `actor_id`, exact
`distance_feet`, `line_of_sight` and `cover`, and a `beast` object:

- `version`: the revision `srd521.BeastRecord(profile)` returns. It hashes the
  whole validated block, including required mechanics.
- `action_id`: `bite` for the Wolf, `multiattack` for the Brown Bear.
- `target_size`: the declared size of the player the NPC attacks.

A Wolf also needs `pack_ally`: `{"present": false}`, or
`{"present": true, "distance_feet": N, "incapacitated": bool}` for the wolf's
ally nearest the player. The ally is a declared fact, not an encounter actor.
A missing size, ally or spatial fact makes the preview ask for clarification
and rolls nothing. A changed `beast` object during an active encounter
requires adjudication. A changed `pack_ally` or spatial fact is a new ruling
that applies to the next attack.

Each profile has its own capability set. The block loads only when every
required mechanic key is in that set and its actions, traits and size-gate
text match what the runtime resolves; any other key or shape refuses it.
Skills, Darkvision and Climb are in the sets because the encounter derives no
check, sight ruling or path from them. `Monster.Available()` stays false for
every record.

## Rules

- Both attacks need visible reach within 5 feet and no total cover.
- Pack Tactics: the Wolf has Advantage when its ally is present, within 5 feet
  of the player and not Incapacitated.
- A hit from the Wolf's Bite (Medium or smaller) or the Bear's Claw (Large or
  smaller) gives the player the Prone condition, with no saving throw. A
  larger target is exempt; an already Prone target is not affected twice.
- Prone is saved on the player's encounter actor. Melee attacks against the
  Prone player have Advantage. Several sources still give one Advantage.
- While Prone, the player's attack rolls have Disadvantage: the longsword
  attack and the opportunity attack roll two d20s and keep the lower. The
  preview names the Disadvantage before anything is rolled. Advantage from
  another source would cancel it to one d20, but the player's attacks have no
  Advantage source in these fixtures.
- While Prone, the player crawls: each 5-foot move, approach or retreat costs
  10 feet of the turn's movement. The fixture declares no Difficult Terrain.
  Standing up (`combat-stand`) costs half of 30 feet and ends the condition.
- Every attack against the Prone player is a melee attack within 5 feet, so
  none has the beyond-5-feet Disadvantage.
- A Beast beyond 5 feet first uses its walk speed (40 feet for both blocks) to
  close the declared distance in 5-foot steps on a clear path, then attacks if
  it ends within 5 feet. If its speed runs out, or the path is not clear
  (cover or no line of sight), it takes its turn without attacking and the
  encounter continues. The End turn preview names the approach.
- Multiattack spends one Action. Ruling: the Bite resolves before the Claw.
  Reach, sight, size and revision are checked before any roll, so an invalid
  target rolls neither attack. If the Bite reduces the player to 0 HP, the
  Claw is not rolled.

## Persistence

The NPC turn runs inside the existing confirmation transaction. All attacks,
HP, Prone, budgets, cursor, clock, receipt and history commit together. A
failed die on any attack or a failed history write commits none of them;
retrying the same identity can roll again. Replaying a completed identity
returns its saved result without rolling or applying damage or Prone again.

Receipts are the shared AttackReceipt with `advantage` and `disadvantage`
(reasons), `d20s` (both dice when only one of them applies; empty when they
cancel) and `condition` (`prone-applied`, `already-prone` or `size-exempt`).

## Limits

No other monster, general animal behaviour, ally actors, pet UI, other
conditions that impose Disadvantage, Difficult Terrain, ranged attacks against
a Prone target, legendary actions, monster AI or catalogue promotion. The web
receipt shows the kept die; both d20s appear in the narrative.
