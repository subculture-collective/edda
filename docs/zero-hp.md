# Dropping to 0 HP

Tracker #2164. SRD 5.2.1 PDF pp. 17–18 ("Dropping to 0 Hit Points", "Death
Saving Throws", "Stabilizing a Character"), p. 181 ("Death Saving Throw"),
p. 188 ("Stable") and p. 191 ("Unconscious"). This covers the player character only.
It is partial evidence, not Full Rules support.

## Damage

Falls and NPC attacks (Bandit, Wolf, Brown Bear and the fixture foe) pass each
damage source through one engine function in the confirmation transaction:

- Damage that drops the character from above 0 HP to 0 HP kills it when the
  remaining damage is at least its HP maximum. Otherwise the character gains
  the Unconscious condition through `ApplyConditionChange`, which also makes it
  Prone, drops held weapons and a Shield, and ends Concentration.
- Damage at 0 HP is one Death Saving Throw failure, or two from a Critical
  Hit. Damage of at least the HP maximum kills. Damage ends Stable.
- Each Multiattack attack is a separate damage source.

## Attacks against an Unconscious player

NPC attacks against an Unconscious player have Advantage. Within 5 feet every
hit is a Critical Hit (the receipt's `auto_critical` names why). A Prone
player gives Advantage within 5 feet and Disadvantage beyond it, so a Bandit's
Light Crossbow at range rolls one d20 when the two cancel.

## Death Saving Throws

`stats.srd521_death_saves` holds `successes`, `failures`, `stable` and `dead`.
It is decoded strictly: an unknown field or impossible count makes every action
unavailable with an adjudication reason.

- **In an encounter** the server rolls the save when the player's turn starts
  at 0 HP and the character is neither Stable nor dead. It runs in the same
  confirmation as the NPC turn that hands the turn back. The player can only
  end the turn while Unconscious.
- **Outside an encounter** Edda has no turn order, so the player confirms each
  save with the `death-save` action (web action menu, TUI `/death-save`). Each
  save stands for the start of one turn and advances campaign time 6 seconds.
  A Stable character uses the same action to wait 1d4 hours and regain 1 HP.

A d20 of 10 or more succeeds and below 10 fails. A 1 counts as two failures.
A 20 regains 1 HP, which ends Unconscious and leaves Prone. Three successes
make the character Stable and reset the counts. Three failures kill it.
Regaining HP clears the counts.

A dead character has no available actions; no rule revives it. An encounter
completes when the player dies and hands its conditions back to the character.

## Free-text turns

A character that is dead, at 0 HP or Unconscious cannot act, so a free-text
turn never reaches the model. The web story and the TUI show the reason and
point to the `death-save` action (or, in an encounter, to ending the turn).
The check reads the same `srd521_death_saves` and condition state as the
rules actions. SRD 5.2.1 campaigns already refuse every free-text turn with a
generic guard; this reason replaces it while the character cannot act. Edda
has no mode in which the model narrates the world while the player waits.

## Characters left at 0 HP by older builds

Before #2211, a fall or NPC attack could leave a character at 0 HP without
the Unconscious condition. Every other action is unavailable for that
character, and its reason says what to do next. Outside an encounter the only
choice is `death-save`. Inside an active encounter `death-save` is refused;
the character can end its turn, or decline a pending opportunity attack (an
older build could offer one, and the turn cannot end while it waits).
Confirming `death-save` once applies Unconscious through
`ApplyConditionChange`, which makes it Prone, drops held weapons and a Shield,
and ends Concentration. It rolls nothing and takes no time; its `srd_death_save` receipt has the event
`unconscious-repair`. The next `death-save` rolls the first save. A repair on
load would write state from a read, and a data migration could not end
Concentration or drop held items, so the repair is an explicit action.

Inside an active encounter the repair runs when the player's turn starts, in
the same confirmation, before that turn's Death Saving Throw. It applies
Unconscious to the encounter's player actor through `ApplyConditionChange`
and rolls nothing; the encounter result carries its receipt as
`unconscious_repair` beside `death_save`. The save that follows uses the
ordinary rules.

The API sends the encounter's `unconscious_repair` and `death_save`, and the
`receipt` of an `srd_death_save` call, as `srd_death_save` resolution events,
live and in history. The web shows each as a 0 HP receipt.

Resolution events are listed in the order the engine resolved them, and the web
shows them in that order. In an encounter turn, each NPC attack comes before the
Concentration saves it caused (`concentration_after_attack` records which attack
caused each one), then the 0 HP receipts, then any effect the clock expired.
Turns saved without that field list their saves after all the attacks.

Every roll is in a durable receipt (`zero_hp` on fall and NPC attack receipts,
`death_save` on the encounter result, `srd_death_save` for the explicit
action). Previews never roll or write. Replays return the saved result.

## Help from an ally

Outside an encounter a trusted ally NPC can help the Unconscious player with
the `ally-aid` action (web action menu, "At 0 HP"). The player takes no
action; the ally does. Sources: PDF p. 10 and p. 18 (the Help action can
stabilize a creature at 0 HP with a DC 10 Wisdom (Medicine) check), p. 99
(Potion of Healing: a Bonus Action to administer it to a creature within
5 feet) and p. 204 ("Using a Potion").

- **Who can help.** An NPC whose `properties.srd521_ally` holds the
  `aid-player-v1` profile: the shared ally facts from
  [the Cure Wounds ally slice](spell-actions.md#willing-adjacent-ally-263)
  plus `medicine_modifier`, an integer from -5 to +28. The ruling names the
  player (`actor_id`) and the current location, and the NPC is within 5 feet,
  in sight, without cover. Missing, extra, null or mistyped facts leave the NPC
  out. Model tools cannot write the `srd521_ally` namespace.
- **Medicine** (`entry_id: "medicine"`). The ally rolls d20 plus its
  modifier against DC 10. Success makes the player Stable and resets the
  counts. Failure changes nothing. A 20 or a 1 on this check has no special
  effect. It is not offered to a Stable player.
- **Potion** (`entry_id: <item id>`). The ally takes one of the player's own
  Potions of Healing and gives it to the player. Edda tracks no NPC
  inventory, so the potion comes from the player's pack; taking it is an
  Edda ruling under the ally's trusted fixture. The player regains 2d4 + 2 HP
  from 0. Unconscious ends through `ApplyConditionChange`, Prone remains, and
  the counts clear. One potion is used. A Stable player can also be healed.
- **One round.** Each confirmation stands for one 6-second round and advances
  campaign time by 6 seconds. The ally acts first. If the player is still
  dying afterwards (failed Medicine), the round's Death Saving Throw is rolled
  in the same transaction, so repeated attempts are not free.

The confirmation names the helper as `target_id`. Only an ally offered by the
locked, revalidated preview crosses it. The receipt is a `srd_death_save`
tool call with event `ally-stabilize` or `ally-healing`, the helper, the
check or the potion, and the following save as `death_save`. Previews never
roll or write; replays return the saved receipt; a history failure rolls back
HP, counts, conditions, the potion and the clock. Evidence:
`internal/rules/srd521/ally_aid_test.go`,
`internal/handlers/ally_aid_integration_test.go` and
`frontend/tests/browser/death-save.spec.ts`.

## Not supported

- Help inside an encounter. The two-actor encounter has no ally that can act.
- An ally casting Cure Wounds or Spare the Dying (#958). NPCs have no
  spellcasting state, and Spare the Dying has no spell implementation.
- A terminal command for ally aid. The TUI keeps `/death-save` only; ally
  selection is web-only, as with the ally potion slice.
- The Unconscious player still cannot drink a potion or cast Cure Wounds.
  `ResolveHealingPotion` refuses 0 HP; `ResolveAdministeredPotion` accepts it
  for a potion given by another creature.
- Falls while at 0 HP, monster death saves, Temporary Hit Points and HP
  maximum reduction.
