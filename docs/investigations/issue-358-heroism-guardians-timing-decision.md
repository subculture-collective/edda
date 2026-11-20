# Heroism start-turn and Spirit Guardians entry/end timing ownership (#358)

Heroism already has a bounded published execution owner, #912, which reaches both qualification gates through its cohort #848. It is reused unchanged and gets no duplicate. Spirit Guardians has no resolver anywhere in the tracker. Its residual is assigned to two bounded tasks: #2058 for the designated exclusions, Emanation membership and the once-per-turn saves, and #2059 for the Speed halving. Both are published as native Task issues; their exact bodies are committed beside this decision, and the readback is in the [link receipt](issue-358-link-verification.json).

## Source and existing owners

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Source facts and fixtures retain the CC-BY-4.0 attribution and adaptation notice in [the source register](../srd-5.2.1-conformance.md). Every clause below was read from the pinned pages, and the scope file records the checked phrases per page.

| Named scope | PDF | Owner and boundary |
|---|---|---|
| Heroism record and list membership | 140 | Catalogue #209 holds the record. #856 owns exact Paladin membership and #503 the Bard list. A record grants no action. |
| Heroism cast and recipients | 140, 179 | #912: Action, V, S, one willing touched creature, Concentration up to one minute, one more recipient per slot level above one. #353 ends the exact cast on a failed Constitution save after damage, on incapacitation and on the other Concentration endings. |
| Heroism Frightened immunity | 140 | #912, consuming the shared Frightened owner #488. |
| Heroism start-turn grant | 140 | #912: at the start of each of the recipient's own turns it gains Temporary Hit Points equal to the caster's spellcasting ability modifier, with no invented minimum, no permanent stacking and no grant after expiry. The grant belongs to the recipient's turn start, never to an optional activation. |
| Temporary Hit Points pool rules | 18, 190 | Named residual of open splitter #272, which already names temporary-HP recovery boundaries; #428 owns the completed Long Rest. See "Temporary Hit Points" below. |
| Spirit Guardians record and list membership | 39, 164 | Catalogue #220 holds the record. #503 owns Cleric list membership. A record grants no action. |
| Spirit Guardians Cleric casting modes | 36, 105–106, 164 | #556: slot of the spell's level or higher, one slot cast per turn, V, S, M, hands, DC 8 + PB + WIS. The unpriced, unconsumed prayer scroll may be replaced by a Component Pouch or by a Spellcasting Focus a feature allows, with a free hand (PDF105–106). #557 owns preparation, #558 slots and the level transitions run through #597. None of them resolves the effect. |
| Spirit Guardians designation, membership, triggers, save, damage, upcast, ending | 5, 16, 21–22, 36, 105–106, 164–165, 177, 179, 181, 187 | No existing owner. New task #2058. |
| Spirit Guardians Speed halving | 5, 106, 164, 181, 188 | No existing owner. New task #2059. |
| Priest casts Spirit Guardians 1/Day, spell save DC 13 | 316 | Catalogue #170 holds the Priest record. Monster spellcasting fixtures belong to the open monster splitter #182, which feeds world qualification #183. Out of this spike's stated targets and unavailable; that cast consumes #2058 and #2059 and does not copy them. |
| Same-spell potency ordering and all other spell keys | 106 | Residual reconciliation #277. |

Whole-body search of all 1,892 issues, open and closed, found "Spirit Guardians" only in the roadmap #47, the tracking group #193, catalogue #220 and this spike. No issue names its prayer-scroll component, its designation or its triggers. The nearest analogues, #1044 Conjure Woodland Beings and #619 Aura of Life, own different spells and stay unchanged. The spell is on the Cleric list only (PDF39, 164), and the Priest stat block casts it without a slot (PDF316). The finite casting adapters #1447–#1450 cover the accepted 245-identity domain, which excludes it.

#2058 does not re-implement Cleric casting. Like #619, it takes its committed cast operands (slot level, bound DC, alignment, cast ID) from an explicit isolated entitlement and leaves component, focus, slot-per-turn and DC validation with #556. #556 says "Holy Symbol replaces only noncostly nonconsumed M" and does not name the Component Pouch; the pouch path is part of its cited PDF104–106 obligation and is recorded here without editing #556. #556 cannot be a native prerequisite of the new tasks: #556 is blocked by #503, which is blocked by #277, and the terminal blocks #277 (see "Native plan and caps"), so that edge would close a cycle.

### Temporary Hit Points

#912's body states the grant amount, no invented minimum, a "permanent tempHP stacking" negative and an "expired" negative. It does not state four PDF18 and PDF190 clauses that a per-turn refresh needs:

- Temporary Hit Points are lost before real Hit Points when damage is taken.
- They cannot be added together. A creature that has some and receives more decides whether to keep the old amount or take the new one. Heroism raises that choice at every turn start.
- They last until depleted or until a finished Long Rest. PDF140 states no removal when Heroism ends, so they persist after the spell.
- They are not Hit Points and receiving them is not healing.

The tracker has no shared Temporary Hit Points pool. #807 states the non-stacking choice and the duration for its own feature; #573, #1000 and the Polymorph owners state their own grants. This decision assigns the four clauses by name to the open splitter #272 as an explicit residual and publishes no pool task of its own, as #358 directs ("reusing #272/#277 ownership"). #912 is not edited. Its fixture cannot qualify until the pool owner from #272 exists, and its implementer must neither auto-replace nor auto-keep.

## Timing decision

The two spells do not share a trigger and get no shared scheduler.

**Heroism** is recipient-keyed: one grant per recipient at the start of that recipient's own turn while the cast is active. Its turn identity is the recipient's turn. The planning key `effect.periodic.recipient-start` stays with #912 and unregistered.

**Spirit Guardians** is not a start-of-turn effect. PDF164–165 name three triggers: the Emanation enters a creature's space, a creature enters the Emanation, and a creature ends its turn there. Each is a consequence of a committed movement or end-turn action. The decision fixes four facts for #2058:

- Turn identity is the turn in progress (round, turn owner and sequence). The once-per-turn limit binds cast, creature and that turn, so the same creature can save on the caster's turn and again on its own.
- Phase is the trigger's own action: a movement step or the creature's turn end. There is no start phase and no background tick.
- Spatial revision is the existing exact pairwise distance from #67, changed only by committed five-foot movement. Inside means 15 feet or less with no Total Cover from the caster (PDF177). Edda stores no grid, so this distance test is a labelled adaptation. The origin of an Emanation is outside its area unless its creator decides otherwise (PDF181); the task keeps that choice unavailable as a labelled restriction.
- Deduplication is one record per cast, creature and turn, committed with the save.

The planning key `effect.area.entry-or-target-end.once-per-turn` moves to #2058 and stays unregistered until that task qualifies.

## Recorded ambiguities

- The source lists no cast-time trigger. A creature already inside when the spell is cast does not save until it ends its turn there or re-enters.
- PDF187 lets the person whose turn it is order things that happen at the same time. No ordering owner exists and this spike invents no scheduler. In #2058, an action that would resolve the save together with another same-moment effect on the same creature stays unavailable with an explicit reason. #912 is unchanged.
- PDF181 lets the creator include the Emanation's origin. #2058 does not offer it: including the caster is unavailable as a labelled Edda restriction, not a source prohibition.
- PDF164 sets no limit on which creatures may be designated. Limiting the choice to creatures the caster's player already knows is a labelled Edda disclosure adaptation.
- PDF164–165 give a damage type for a good, neutral or evil caster and none for an unaligned one (PDF22). That cast stays unavailable. No published issue owns a player character's alignment fact, so the fixture supplies it and a production cast stays unavailable until an owner exists.
- No page states how a Speed that changes during a move is accounted. #2059 applies the PDF188 switching subtraction as a labelled Edda adaptation and reports it in the receipt.
- Heroism's keep-or-replace choice at each turn start and the persistence of its Temporary Hit Points after the spell ends are PDF18 rules that #912's body does not spell out. They belong to the #272 residual above.
- PDF106 potency ordering for overlapping casts of the same spell stays with #277. The new task admits one active Spirit Guardians per encounter.
- Entries from forced movement and teleportation stay unavailable until their own owners emit trusted position changes.

## Native plan and caps

Both tasks are children of #193, type Task, Size M. #2058 is blocked by #220 and #353. #2059 is blocked by #2058 and is the terminal. Closed baselines #58, #61, #64, #66, #67, #68, #77 and #351 are consumed without new native edges; #67 already has 100 consumers. The open spike is not made a prerequisite of its children.

The terminal blocks item qualification #280 directly. Spell qualification #279 holds 100 native blockers, the platform maximum, so the direct #279 edge that #358 asks for cannot be added. A path through #280 alone is not enough: #279 blocks #280 and not the reverse, and seven open issues wait on #279 without reaching #280 (#571, #574, #577, #578, #579, #582, #656). No finite aggregate fits: #1477 is fixed to the 245-identity Sorcerer, Warlock and Wizard union, and #848–#850, #973 and #1096 are finite class rosters.

The terminal therefore also blocks #277. #277 is an open direct blocker of #279, it is already blocked by catalogue #220, and its acceptance requires every remaining spell key to have an owner wired to the spell gate; Spirit Guardians is one of those keys. The native path is #279 → #277 → #2059 → #2058. This is the one addition to an existing non-gate issue: #277 goes from 30 to 31 blockers, and its body, parent and existing edges are unchanged. The consequence is that #277 cannot close before Spirit Guardians qualifies. #358's "directly to both" criterion is met directly for #280 and through one native hop for #279.

The only way to a direct #279 edge is to replace #358's own #279 slot with the terminal at publication. That removes an existing edge, so it needs the owner's explicit approval and is not part of this plan. If the owner approves it, the #277 edge is not added.

Counts read immediately before publication (2026-11-20T08:03:49Z, 1,892 issues): #279 had 100 blockers, #280 had 96, #277 had 30, #193 had 30 children, #220 blocked 3 and #353 blocked 76. Counts read back after publication (2026-11-20T08:06:27Z, 1,894 issues): #279 has 100, #280 has 97, #277 has 31, #193 has 32 children, #220 blocks 4 and #353 blocks 77. No issue exceeds 100 blockers, 100 dependents or 100 children. The dependency graph (10,206 edges) and the parent hierarchy (1,890 edges) are acyclic, and #279 and #280 each reach #2059 and #2058. Compared with the read before publication, the only differences are the two new issues and one added blocker each on #280 and #277; no existing title, state, body, type, parent or edge changed.

Heroism's owner #912 keeps its existing path #912 → #848 → #279 and #280.

## Evidence and limits

The [scope file](issue-358-heroism-guardians-timing-scope.json) records the checked source phrases, the owner body hashes and gate paths, both task templates and hashes, the planned native edges and the cap counts. The [link receipt](issue-358-link-verification.json) records the exact published bodies of #2058 and #2059, which match the committed task files byte for byte, with their parents, blockers, type and Size, the cap and acyclicity checks over every issue, and that #279 and #280 each reach the terminal.

Focused checks for this decision: the PDF hash, the recorded phrases found on their cited pages, every page cited by a task or by this decision present in the checked set, body hashes against a full tracker read, both graph checks and gate reachability. Runtime implementation remains paused. This change contains no runtime files, database migration, UI behavior, production mutation or Full Rules admission. The full Go suite was not run locally; hosted CI applies to this decision, and the implementation tasks keep their own runtime, PostgreSQL and device requirements.

## Addendum, 2026-11-20: the terminal moves from #277 to #279

This addendum corrects one edge of "Native plan and caps". The rest of the decision, both task bodies and the first [link receipt](issue-358-link-verification.json) stand as dated evidence.

The original plan made the terminal #2059 a blocker of the open splitter #277, so that the full gate #279 reached it through #277. That left a splitter blocked by a runtime implementation task: #277 could not close until #2059 was implemented, although assigning keys to owners does not depend on any owner's runtime work. The orchestrator has since set a standing gate policy that forbids this and approves the replacement the original plan left to the owner:

- #279 and #280 are at the native cap of 100 blocked-by entries.
- A closed splitter's own edge on a gate carries no remaining obligation once its terminal qualification child exists.
- When a gate is full, that splitter's own gate edge is replaced by its terminal qualification child, so the count stays at 100.
- No other edge is removed, no issue exceeds 100, no open splitter is blocked by a runtime implementation task, and both graphs stay acyclic.

#358 is closed and its terminal #2059 exists, so its #279 slot was replaced. Three native changes were applied in this order, so that #279 kept a path to #2059 throughout:

1. Removed "#279 blocked-by #358".
2. Added "#279 blocked-by #2059".
3. Removed "#277 blocked-by #2059".

Read before the change (2026-11-20T16:17:09Z, 1,951 issues, 10,562 dependency edges): #279 had 100 blockers including closed #358, #277 had 31 including #2059, and #280 had 100 including #2059. Read back after it (2026-11-20T16:20:44Z, 1,951 issues, 10,561 dependency edges): #279 has 100 blockers, #277 has 30 and #280 has 100. The only differences in the whole tracker are those three edges. Every issue body, state and parent is unchanged, and the bodies of #2058 and #2059 still match the committed task files.

The native paths are now #279 → #2059 → #2058 and, as before, #280 → #2059 → #2058. #277 is blocked only by #80 and the catalogue issues #196–#224, its 30 entries from before this decision, and no longer reaches #2059. The seven issues that wait on #279 without reaching #280 (#571, #574, #577, #578, #579, #582, #656) still reach the terminal through #279. The dependency graph and the parent hierarchy are acyclic, and the largest counts are 100 blockers, 100 dependents and 95 children.

Two consequences are recorded and not acted on:

- 237 open issues that waited on the terminal only through #277 (#277 itself, #399, #503, #2111 and their dependents) no longer do. That is the intended effect: they depend on spell-key assignment, not on Spirit Guardians runtime work. 53 issues still reach the terminal.
- The cycle that kept #556 from being a native prerequisite of #2058 ran through the #277 edge and no longer exists. This change adds no #556 edge; #2058's body and blockers are unchanged.

#358 still holds a blocked-by slot on #280, which also blocks directly on #2059. That edge was outside this correction and is untouched. The [rewire receipt](issue-358-gate-rewire-verification.json) records the edges, counts, paths and graph checks. No runtime file changed and the full Go suite was not run.
