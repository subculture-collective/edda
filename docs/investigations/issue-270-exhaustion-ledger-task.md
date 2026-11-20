<!-- edda-270:exhaustion-ledger -->
Decision #270; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 17–18, 179–181, 183, 185, 189, 194, 236; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only the level ledger of the Exhaustion condition (PDF181): one authoritative per-creature level with per-level provenance, gain, Immunity, death at level 6, the ending at level 0, the predicate that says which level may be removed, and the revival reduction. No published task owns these. The two penalties the level causes belong to the Exhaustion penalty task (`edda-270:exhaustion-penalties`) and are not implemented here.

Reduction triggers stay where they are. #428 keeps the completed Long Rest trigger (PDF185). #623 keeps the trusted one-level reduction event and its recomputation call, and #1001 consumes that event through #623 as its body says. This task supplies the stored level they lower and the rule for which level is removable; it copies no trigger and keeps no second counter. Sources of Exhaustion (dehydration, malnutrition, suffocation, travel, temperature, contagions, items and stat blocks) stay with their own owners; this task admits a level only from a trusted, source-bound gain event and enables no such source in production.

Levels (PDF179, 181): Exhaustion is the one condition that stacks with itself. Each trusted receipt of the condition adds exactly 1 level and records its source, receipt and order. Valid levels are 0–6. A living creature holds 0–5, and 0 means the condition is absent. Reaching 6 commits the creature's death in the same transaction as a trusted source-bound death record, and the dead creature keeps level 6. No seventh level is ever added. No published task owns a general dead state, 0 Hit Points or Death Saving Throws (PDF17–18); this task adds none of them and does not call the legacy combat death code. The death record it writes is the authoritative death fact that #618 reads.

Immunity to the Exhaustion condition (PDF183) means a gain does not affect the creature at all: no level, no penalty, no history of a hidden level.

Removal (PDF181): when the level reaches 0 the condition ends. A reduction removes exactly one removable level, never below 0. Levels differ only in provenance; taking the most recently gained removable level first is an Edda ordering, not source text. A source owner may mark its own levels non-removable until its stated fact is met (PDF181 dehydration, PDF185 malnutrition) or remove all levels of its own source (PDF189 suffocation). This task supplies only that provenance and predicate, using isolated trusted facts, and never infers a meal, drink or breath.

Revival (PDF180): "If the creature died with any Exhaustion levels, it returns with 1 fewer level." Apply it to the actual #618 Revivify receipt: a creature that died at level 2 returns at 1, and one that died at level 6 returns at 5. Revival itself, its Hit Points and its other consequences stay with #618. Every other revival owner (#947 Raise Dead and later ones) is a caller of this same output. The ledger fails closed: a dead creature with stored levels can become alive only through a revival receipt that applied this reduction, so a revival owner that has not yet consumed it stays unavailable for such a creature.

Not in this task, and unavailable until the source owner supplies a rule, because the source gives no precedence:

- A level that "lasts until the contagion ends" (PDF194) against ordinary reduction.
- An effect that removes "any Exhaustion levels you have" (PDF236 Potion of Vitality) against levels a source made non-removable.
- The revival reduction when the only remaining levels are non-removable.

## Deterministic fixtures

- Two trusted gains from different sources give level 2 with two provenance rows, not a single flag. A third gain gives 3.
- At level 5 one more gain commits level 6 and one death record; replay does not kill or add again. A gain on the dead creature, a seventh level, a client-declared level, a multi-level gain without a qualified source, a negative level or a gain on an immune creature rejects without mutation.
- Ledger primitive, exercised by an isolated test event (the first real callers are #623 and #428): lowering level 1 ends the condition; at level 0 it is a no-op without history. With one non-removable level and one ordinary level it removes the ordinary level; with only the non-removable level it removes nothing. Remove-all-of-source removes only that source's levels.
- Died at level 2, revived by the actual #618 receipt: level 1. Died at level 6: level 5 and alive. A forged, foreign or replayed revival receipt changes nothing. Marking a dead creature with levels alive without that receipt rejects.
- Reject remove-all against a non-removable level, a contagion-duration level and revival against only non-removable levels as unavailable, each with a reason and no mutation.

## Durable implementation acceptance

Preview is inert: no dice, writes or hidden facts. Confirm revalidates campaign, actor, trusted source and receipt, immunity, current level, source version and expected revision. Commit the level change, provenance, any death record and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. Same action ID and payload replay after retry, reload or reconnect returns the saved receipt with no second level or death; a conflicting payload rejects; racing gains or reductions each commit once in a serial order. Web and trusted TUI show the player's own level, share unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. An NPC's level and sources never enter previews, receipts, errors or public history. Preserve legacy reads and unrelated JSON: the legacy `exhausted` combat flag and player status carry no level and are never converted into one.

Hard prerequisites: #64 confirmation and receipts; #618 actual Revivify receipt and its authoritative death fact. Parent: #270. Directly block the Exhaustion penalty task (`edda-270:exhaustion-penalties`); through it and the existing blocker #623 this task is mandatory for #279, #280 and #337. #279 is at its 100-blocker cap: add no direct #279 edge. Add no #67 or #429 consumer, no edge to #428, #947 or #1001, and change no existing body, parent or prerequisite.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
