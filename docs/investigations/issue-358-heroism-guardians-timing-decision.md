# Heroism start-turn and Spirit Guardians entry/end timing ownership (#358)

Heroism already has a bounded published execution owner, #912, which reaches both qualification gates through its cohort #848. It is reused unchanged and gets no duplicate. Spirit Guardians has no resolver anywhere in the tracker. Its residual is assigned to two bounded tasks: {{ISSUE_SG_SAVES}} for the cast, the designated exclusions, Emanation membership and the once-per-turn saves, and {{ISSUE_SG_SPEED}} for the Speed halving. Their exact bodies are committed beside this decision; native publication and readback are recorded in the companion receipt.

## Source and existing owners

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Source facts and fixtures retain the CC-BY-4.0 attribution and adaptation notice in [the source register](../srd-5.2.1-conformance.md). Every clause below was read from the pinned pages, and the scope file records the checked phrases per page.

| Named scope | PDF | Owner and boundary |
|---|---|---|
| Heroism record and list membership | 140 | Catalogue #209 holds the record. #856 owns exact Paladin membership and #503 the Bard list. A record grants no action. |
| Heroism cast and recipients | 140, 179 | #912: Action, V, S, one willing touched creature, Concentration up to one minute, one more recipient per slot level above one. #353 ends the exact cast on damage, incapacitation and the other Concentration endings. |
| Heroism Frightened immunity | 140 | #912, consuming the shared Frightened owner #488. |
| Heroism start-turn Temporary Hit Points | 18, 140, 190 | #912: at the start of each of the recipient's own turns it gains Temporary Hit Points equal to the caster's spellcasting ability modifier, with no invented minimum, no stacking and no grant after expiry. The grant belongs to the recipient's turn start, never to an optional activation. |
| Temporary Hit Points lasting to a Long Rest | 18 | Open splitter #272 names temporary-HP recovery boundaries and #428 owns the completed Long Rest. Not reassigned here. |
| Spirit Guardians record and Cleric eligibility | 164 | Catalogue #220 holds the record. Production eligibility stays with #503 list membership, #556 Cleric casting, #557 preparation, #558 slots and the level transitions through #597. None of them resolves the effect. |
| Spirit Guardians cast, designation, membership, triggers, save, damage, upcast, ending | 5, 16, 21, 104–106, 164–165, 177, 179, 181, 187 | No existing owner. New task {{ISSUE_SG_SAVES}}. |
| Spirit Guardians Speed halving | 5, 106, 164, 181, 188 | No existing owner. New task {{ISSUE_SG_SPEED}}. |
| Same-spell potency ordering and all other spell keys | 106 | Residual reconciliation #277. |

Whole-body search of all 1,892 issues, open and closed, found "Spirit Guardians" only in the roadmap #47, the tracking group #193, catalogue #220 and this spike. No issue names its prayer-scroll component, its designation or its triggers. The nearest analogues, #1044 Conjure Woodland Beings and #619 Aura of Life, own different spells and stay unchanged. The finite casting adapters #1447–#1450 cover the accepted 245-identity domain, which excludes this Cleric-only spell, so {{ISSUE_SG_SAVES}} implements its own stated casting modes or stays unavailable.

The tracker has no shared Temporary Hit Points primitive. Each consumer, including #912, #573, #807 and #1000, states the PDF18 rules in its own fixture. That is existing ownership, not a residual of this spike.

## Timing decision

The two spells do not share a trigger and get no shared scheduler.

**Heroism** is recipient-keyed: one grant per recipient at the start of that recipient's own turn while the cast is active. Its turn identity is the recipient's turn. The planning key `effect.periodic.recipient-start` stays with #912 and unregistered.

**Spirit Guardians** is not a start-of-turn effect. PDF164–165 name three triggers: the Emanation enters a creature's space, a creature enters the Emanation, and a creature ends its turn there. Each is a consequence of a committed movement or end-turn action. The decision fixes four facts for {{ISSUE_SG_SAVES}}:

- Turn identity is the turn in progress (round, turn owner and sequence). The once-per-turn limit binds cast, creature and that turn, so the same creature can save on the caster's turn and again on its own.
- Phase is the trigger's own action: a movement step or the creature's turn end. There is no start phase and no background tick.
- Spatial revision is the existing exact pairwise distance from #67, changed only by committed five-foot movement. Inside means 15 feet or less with no Total Cover from the caster (PDF177); the caster is the origin and outside the area (PDF181). Edda stores no grid, so this distance test is a labelled adaptation.
- Deduplication is one record per cast, creature and turn, committed with the save.

The planning key `effect.area.entry-or-target-end.once-per-turn` moves to {{ISSUE_SG_SAVES}} and stays unregistered until that task qualifies.

## Recorded ambiguities

- The source lists no cast-time trigger. A creature already inside when the spell is cast does not save until it ends its turn there or re-enters.
- PDF187 lets the person whose turn it is order things that happen at the same time. {{ISSUE_SG_SAVES}} requires that explicit choice when another effect shares the moment and builds no universal order. #912 is unchanged.
- No page states how a Speed that changes during a move is accounted. {{ISSUE_SG_SPEED}} applies the PDF188 switching subtraction as a labelled Edda adaptation and reports it in the receipt.
- PDF140 states no removal of Heroism's Temporary Hit Points when the spell ends, so the PDF18 duration applies.
- PDF106 potency ordering for overlapping casts of the same spell stays with #277. The new task admits one active Spirit Guardians per encounter.
- Entries from forced movement and teleportation stay unavailable until their own owners emit trusted position changes.

## Native plan and caps

Both tasks are children of #193, type Task. {{ISSUE_SG_SAVES}} is Size L and blocked by #220 and #353. {{ISSUE_SG_SPEED}} is Size M and blocked by {{ISSUE_SG_SAVES}}; it is the terminal and directly blocks item qualification #280. Closed baselines #58, #61, #64, #66, #67, #68, #77 and #351 are consumed without new native edges; #67 already has 100 consumers. No existing body, parent or edge changes apart from adding the terminal to #280. The open spike is not made a prerequisite of its children.

Observed counts at drafting: #279 has 100 blockers, #280 has 96, #193 has 30 children, #220 blocks 3 and #353 blocks 76. After publication #280 has 97, #193 has 32, #220 blocks 4 and #353 blocks 77. The dependency graph with the planned edges and the parent hierarchy are acyclic across 1,892 issues.

Spell qualification #279 is at the native 100-blocker maximum, so the direct #279 edge that #358 asks for cannot be added. No accepted aggregate fits: #1477 is fixed to the 245-identity Sorcerer, Warlock and Wizard union, and #848–#850, #973 and #1096 are finite class rosters. This decision does not hang a Cleric-only spell on one of those contracts and does not remove an edge. The terminal's native mandatory path is #280, which #279 itself blocks, so the final coverage claim still waits for both tasks. #277 already blocks #279, is itself blocked by catalogue #220 and must wire every spell key's owner to the spell gate; it inherits these two owners. Replacing #358's own #279 edge with the terminal would restore a direct path but removes an existing edge, so it needs explicit owner approval.

Heroism's owner #912 keeps its existing path #912 → #848 → #279 and #280.

## Evidence and limits

The [scope file](issue-358-heroism-guardians-timing-scope.json) records the checked source phrases, the owner body hashes and gate paths, both task templates and hashes, the planned native edges and the cap counts. The two task bodies use the placeholders `{{ISSUE_SG_SAVES}}` and `{{ISSUE_SG_SPEED}}` until publication assigns numbers; the receipt then records the exact published bodies, parents, blockers, type and Size, and a fresh cap and acyclicity check.

Focused checks for this decision: the PDF hash, 48 recorded phrases found on their cited pages, every page cited by a task present in the checked set, body hashes against a full tracker read, and both graph checks. Runtime implementation remains paused. This change contains no runtime files, database migration, UI behavior, production mutation or Full Rules admission. The full Go suite was not run locally; hosted CI applies to this decision, and the implementation tasks keep their own runtime, PostgreSQL and device requirements.
