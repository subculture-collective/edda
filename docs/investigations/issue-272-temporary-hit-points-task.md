<!-- edda-272:temporary-hit-point-pool -->
Decision #272; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 17–18, 129, 185, 190; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement the cross-source Temporary Hit Point clauses of PDF18 that no published grant owner states, with one actual grant as the witness: False Life through its existing resolver #1318. Do not copy that spell's dice, slot scaling or components. No new grant source, healing effect or class feature is enabled.

Every grant owner keeps its own amount, trigger, eligibility, per-source non-stacking and duration fixture unchanged: Potent Spellcasting #573, Wild Shape #645, Polymorph #698, #724 and #725, Patient Defense #807, Heroism #912, Tireless #1000, False Life #1318, Dark One's Blessing #1497, True Polymorph #1638–#1642 and Shapechange #1969. This task is a downstream consolidation: it sits after #1318 because existing edges cannot change, and it replaces none of those fixtures. It newly owns only these clauses:

- Damage removes the pool first and only the leftover carries over to Hit Points.
- A creature that already has a pool and receives more, from the same or a different source, decides whether to keep the old amount or take the new one. The decision is an explicit recorded choice, never a sum and never a silent default.
- Healing cannot restore the pool, and receiving Temporary Hit Points does not count as healing. A creature at full Hit Points can receive them. A creature at 0 Hit Points that receives them stays at 0 with its existing Unconscious or dying state.
- The pool lasts until depleted or until the creature finishes a Long Rest. Consume the canonical completed Long Rest #428: its completion removes any remaining pool, whatever its source, in the same atomic receipt. An interrupted Long Rest, including one that grants Short Rest benefits after at least 1 hour, and a completed Short Rest do not remove it. Temporary Hit Points granted after the completion remain.

Two consequences are Edda adaptations of "Temporary Hit Points aren't Hit Points", not source sentences: the Hit Point maximum does not cap the pool, and a grant fires no trigger that depends on healing or on regaining Hit Points. Healing stays capped at the maximum (PDF17) regardless of the pool.

Future consumer: Heroism #912, the published owner that decision #358 reused unchanged, raises the keep-or-replace choice at every turn start and must consume this task; it must neither auto-keep nor auto-replace. #912 has no native edge to this task, so this task reaches spell qualification #279 only indirectly; see the decision.

## Deterministic fixtures

- False Life with injected 2d4 = 5 gives 9. Taking 7 damage leaves 2 and unchanged Hit Points; taking 5 more leaves 0 and reduces Hit Points by 3. PDF18 example: 5 Temporary Hit Points and 7 damage lose the pool and 2 Hit Points.
- Holding 10 and receiving 12: the recorded choice yields 12 or 10; reject 22 and reject an unrecorded automatic replacement.
- At 14 of 20 Hit Points with a pool of 2, the closed potion fixture #266 with injected 2d4 = 6 heals 8 on a conscious ally: 20 Hit Points and the pool stays 2; healing never refills the pool.
- A completed Long Rest with a pool of 2 ends with a pool of 0 in the #428 receipt. A Long Rest interrupted in hour 3 yields Short Rest benefits and keeps the pool; a completed Short Rest keeps it.
- At 0 Hit Points a grant of 9 leaves the creature at 0 and Unconscious through #487, with unchanged Death Saving Throw counters.
- Reject pool removal at Long Rest start, at dawn or on a client clock; a duplicate removal on receipt replay; conversion of the pool into Hit Points; a pool capped at the Hit Point maximum; a healing trigger fired by a grant; removing another creature's pool; and an ending this task invents for a source-owned pool. Depletion and the completed Long Rest always apply; any earlier ending is only the one the grant's own source states.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, actor, source version, clock, targets, the current pool, the grant source, the recorded keep/replace choice and expected revisions. Commit the pool change, Hit Point change, recorded choice or rest-completion removal, and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the same receipt with no reroll; a conflicting payload rejects; racing confirmations commit once. Depletion and the Long Rest removal belong to their triggering confirmed event (the damage or the #428 completion), with no optional activation. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires a focused build, lint and desktop/phone checks. Public history omits private HP, hidden actors, private checks and undisclosed item facts.

Hard prerequisites: #65 committed attack damage, #266 healing witness, #428 completed Long Rest output, #487 Unconscious at 0 Hit Points, #1318 actual False Life grant. Parent: #272. Block the rest/recovery qualification task published with this decision; that task blocks #280. Spell qualification #279 already requires the unchanged #1318 through the accepted finite shared aggregate; do not add a direct #279, #280, #67 or #429 edge and do not alter any existing grant owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this task and its callers qualify. No deployment, source inference or Full Rules promotion.
