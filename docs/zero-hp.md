# Dropping to 0 HP

Tracker #2164. SRD 5.2.1 pp. 16–17 ("Dropping to 0 Hit Points", "Death Saving
Throws", "Stabilizing a Character"), p. 180 ("Death Saving Throw"), p. 187
("Stable") and p. 191 ("Unconscious"). This covers the player character only.
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

Every roll is in a durable receipt (`zero_hp` on fall and NPC attack receipts,
`death_save` on the encounter result, `srd_death_save` for the explicit
action). Previews never roll or write. Replays return the saved result.

## Not supported

- Stabilizing by another creature (Help action with DC 10 Wisdom (Medicine),
  Spare the Dying #958). No ally can act in the two-actor encounter.
- Healing from 0 HP by Cure Wounds or a Potion of Healing. Both are the
  Unconscious character's own actions and the rules refuse them at 0 HP; no
  other creature can heal the player yet.
- Falls while at 0 HP, monster death saves, Temporary Hit Points and HP
  maximum reduction.
- Narrative (LLM) turns are not blocked while the player is Unconscious or dead.
