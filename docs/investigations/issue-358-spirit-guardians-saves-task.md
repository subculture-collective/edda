<!-- edda-358:spirit-guardians-emanation-saves -->
Decision #358; roadmap #47. One implementation PR; Task / Size L.

Source: English SRD 5.2.1, PDF 5, 16, 21, 104–106, 164–165, 177, 179, 181, 187; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only one source identity: Spirit Guardians, cast by one isolated, explicitly entitled Cleric fixture with a level-three or higher slot. This task owns the cast, the creatures designated as unaffected, membership of the caster's 15-foot Emanation, the three source triggers, the once-per-turn Wisdom save and its damage. It does not own the Speed halving ({{ISSUE_SG_SPEED}}), any other spell, a generic area or start-of-turn scheduler, NPC spellcasting or production Cleric eligibility.

Cast (PDF164, 104–106): Action, range Self, V, S, M (a prayer scroll; unpriced and unconsumed, so only a source-allowed focus may replace it), Concentration up to 10 minutes, one slot of level three or higher and at most one slot spent to cast on the turn. Outside the finite shared casting adapters, this task implements exactly these modes for its own fixture or stays unavailable. Starting the cast ends the caster's earlier Concentration through the existing owner. The cosmetic angelic, fey or fiendish form has no mechanical effect.

Designation (PDF164): only when the spell is cast, the caster may designate any creatures known to the caster as unaffected. Persist that set with the cast. It cannot be extended, reduced or re-chosen later. A creature not designated at cast time is affected, including one that arrives later.

Membership (PDF177, 181): the Emanation extends 15 feet from the caster and moves with the caster. The caster is its origin and is not in the area; do not offer self-inclusion. A creature is inside when its trusted exact distance to the caster is 15 feet or less and no Total Cover blocks every straight line from the caster to it. Unknown distance, sight or cover facts keep the affected action or turn transition unavailable with an explicit reason; never guess a position.

Triggers (PDF164–165): a non-designated creature makes a Wisdom saving throw against the caster's bound spell save DC whenever (1) the caster's movement brings the Emanation onto the creature, (2) the creature enters the Emanation, or (3) the creature ends its own turn inside it. Nothing else triggers the save: not the cast itself for a creature already inside, not the start of any turn, not the caster's turn end, and not staying inside during another creature's turn. On a failed save the creature takes 3d8 damage; on a success, half as much, rounded down (PDF5, 16). The damage type is Radiant when the caster's authoritative alignment is good or neutral and Necrotic when it is evil (PDF21); an unknown or unaligned caster keeps the cast unavailable. Each slot level above three adds 1d8. When one movement step forces two or more creatures to save at the same time, roll the damage once for all of them (PDF16).

Once per turn (PDF165): a creature makes this save at most once per turn. Bind the save to the cast ID, the creature ID and the identity of the turn in progress (round, turn owner and turn sequence), not to a round or to wall time. The same creature may save again on a later turn in the same round.

Timing identity: every trigger is a consequence of a committed movement or end-turn action. Preview of that action lists the saves it will force without rolling. Confirm commits the movement or turn end, the membership revision, the save, the damage and the dedup record together. No optional activation is invented for the caster, and no background tick runs between actions.

Simultaneous events (PDF187): when the turn-end save coincides with another effect at the same moment for the same creature, the controller of the creature whose turn it is chooses the order. Require that explicit, authoritative choice before the turn advances. Do not default an order and do not build a universal scheduler.

Overlap (PDF106): the effects of the same spell cast more than once do not combine, and the most potent or most recent cast applies. This PR supports one active Spirit Guardians in an encounter. A second concurrent Spirit Guardians cast in the same encounter stays unavailable with an explicit reason; general potency ordering remains with residual reconciliation #277.

Ending (PDF179): the Emanation, membership records and pending triggers end with the caster's Concentration through #353 (damage, Incapacitated, death, a new Concentration cast, voluntary end) or at ten minutes. Ending never refunds the slot and never reverses committed damage.

## Deterministic fixtures

- Cleric fixture with WIS +3, proficiency bonus +3 and a good alignment: DC 8 + 3 + 3 = 14, level-three slot, injected 3d8 = 4, 5, 6. A failed save takes 15 Radiant; a successful save takes 7. An evil caster deals Necrotic. A level-four slot rolls 4d8; a level-two slot, a second slot on the same turn, a missing V or S, an absent material and focus, or a caster without the explicit entitlement is rejected without mutation.
- Caster moves from 20 feet to 15 feet of an opponent on the caster's turn: one save on that turn. Stepping away and back on the same turn forces no second save. On the opponent's own turn it stays and ends its turn inside: a second save, on a different turn.
- An opponent 20 feet away moves to 15 feet on its own turn and ends its turn inside: exactly one save, at entry, and none at turn end. An opponent that starts its turn inside and never leaves saves only at its turn end, not at its turn start.
- An opponent already at 10 feet when the spell is cast does not save at the cast. At exactly 15 feet a creature is inside; at 20 feet it is outside. Total Cover from the caster keeps a creature within 15 feet outside the area.
- A creature designated at the cast never saves and takes no damage at any trigger. A request to designate a creature after the cast, to add the caster to the area, or to trigger a save at the start of a turn is rejected.
- Two non-designated creatures brought inside by one caster step share one damage roll and make separate saves. Each has its own dedup record.
- With a second end-of-turn effect pending for the same creature, the turn does not advance until its controller supplies the order; both orders produce the stated results.
- A second Spirit Guardians cast while one is active in the same encounter is rejected without mutation.
- Concentration ending between two triggers leaves the committed first result intact and forces no further save. Reload between any two steps restores the cast, designation, membership, dedup records and remaining duration exactly.

## Durable implementation acceptance

Preview is inert: no dice, no writes and no private facts. Confirm locks and revalidates campaign, caster, every affected creature, entitlement, slot, components, alignment, effect revision, spatial revision, clock and turn cursor. Commit cost, cast, designation, movement or turn end, membership, saves, damage, Concentration consequences for damaged creatures, dedup records and public history in one transaction. Inject a failure after each stage, including the history write, and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the committed receipt with no extra roll, charge or save; a conflicting payload rejects; racing confirmations commit once. Reject foreign campaign or actor IDs, hidden creatures, forged costs, forged distances and forged turn identities. Public receipts omit private NPC HP, save modifiers the player has not learned and hidden creatures. Preserve unrelated JSON and legacy reads. Web and trusted TUI share availability, disabled reasons and confirmation, cancel and reconnect behavior; changed web requires focused build, lint and desktop and phone checks.

Entries in this PR come only from the existing five-foot movement of the caster or of the creature. Forced movement and teleportation stay unavailable as entry sources until their own owners emit trusted position changes; do not infer them. Fixtures with more creatures than the current encounter supports run in the rules layer and in isolated handler state; player availability covers only the actors the encounter can represent.

The fixture demonstrates its level-five Cleric entitlement explicitly. Imported spell data, catalogue #220, closing #323 or closing #358 grant nothing. Production discovery stays unavailable until the existing Cleric owners supply list membership (#503), Wisdom and Holy Symbol casting (#556), preparation (#557), slots (#558) and the level transitions through #597.

Hard prerequisites: #220 Spirit Guardians source record, #353 Concentration ending events. Consumed closed baselines without new native edges: #61 saves, #64 confirmation, #66 turn identity and action economy, #67 trusted distance, sight and cover facts, #68 slot casting, #77 Concentration and duration, #351 independent cast identities. Parent: #193. This task is a hard prerequisite of {{ISSUE_SG_SPEED}}, which carries the direct edge to item qualification #280. Spell qualification #279 already holds 100 native blockers, the platform maximum; do not add, remove or repurpose an edge there in this PR.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback, concurrency, replay and reload tests with `-tags=integration`; shared web and TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. No deployment, source inference or Full Rules promotion.
