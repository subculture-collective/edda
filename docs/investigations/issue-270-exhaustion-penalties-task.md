<!-- edda-270:exhaustion-penalties -->
Decision #270; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 6–7, 13, 17–18, 180–182, 188; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only the two derived penalties of the Exhaustion condition (PDF181), read from the actual level of the Exhaustion ledger task (`edda-270:exhaustion-ledger`). This task stores no level and admits none. #623 "recompute all derived penalties" calls this projection after its reduction; #428 and #1001 see the same recomputation through the ledger level.

D20 Tests (PDF6–7, 180–181): "When you make a D20 Test, the roll is reduced by 2 times your Exhaustion level." D20 Tests are ability checks, attack rolls and saving throws. The reduction is a numeric term recorded beside the natural die in the #61 receipt for checks and saves and in the #65 receipt for attack rolls; it is not Disadvantage, and it combines with Advantage, Disadvantage and other modifiers without replacing them. Initiative is a Dexterity check (PDF13) and takes the reduction in the #66 Initiative roll. A hit or miss that the source ties to the number on the d20 (PDF7) keeps using the die face. Damage rolls, DCs the creature sets, AC and passive scores are not D20 Tests and get no invented reduction.

Death Saving Throws are saving throws (PDF17), so the reduction is defined for them, and PDF18 ties their 1 and 20 results to the die face. No published task owns Death Saving Throws. This task writes no fixture against the legacy combat death-save code; the reduction on a Death Saving Throw stays unavailable until a source-bound owner exists and consumes this projection.

Speed (PDF181, 188): "Your Speed is reduced by a number of feet equal to 5 times your Exhaustion level", and each special speed is reduced by the same amount. The reduction feeds the one authoritative current Speed that #780 reads. A Speed that is 0 and can't increase stays 0. A reduction larger than the Speed yields 0, never a negative stored Speed; this clamp is an Edda adaptation, not source text. When the reduction brings a flying creature's Fly Speed to 0, the flying fall task (`edda-270:flying-fall`) decides the fall; this task only supplies the reduced value.

## Deterministic fixtures

- Level 2: natural 14 with +3 becomes 17 − 4 = 13 on an ability check and on a saving throw, and Speed 30 becomes 20 with Climb Speed 30 also 20. Level 5: −10 and Speed 5. Level 0: no term in the receipt and no reduction.
- Level 1, attack roll natural 20 through #65: the reduction is recorded and the attack still hits; natural 1 still misses. Level 2 Initiative is reduced by 4. With Advantage the reduction applies once to the kept die.
- Level 3, Speed 30: Speed 15, Dash through #780 adds 15, and righting from Prone through the Prone movement task costs 7.
- Level 4, flying with Fly Speed 20 and no hover: Fly Speed 0 and one fall through the flying fall task; with hover it stays aloft. Fly Speed 30 at level 4 is 10 and no fall.
- Speed 20 at level 5 is 0 and cannot right itself. Speed 0 from actual Grappled #778 stays 0 at every level and returns to the reduced value when the grapple ends.
- After an actual ledger reduction from level 1 to 0, the roll term and Speed reduction are gone in the same transaction.
- Reject treating the penalty as Disadvantage, applying it to a damage roll, save DC, AC or passive score, applying it twice, a client-declared level or penalty, and a Death Saving Throw reduction through legacy code.

## Durable implementation acceptance

The penalty is a projection: it has no activation of its own and is applied inside another owner's roll or movement, in that owner's receipt and transaction. Preview shows the reduced modifier and Speed from facts the player may know. A level change and both recomputed penalties commit atomically with the derived sheet and public history. Inject a failure after each stage and prove full isolated PostgreSQL rollback. Same action ID and payload replay after retry, reload or reconnect returns the saved receipt with the same term and no reroll; a conflicting payload rejects; a roll racing a level change uses one committed level. Web and trusted TUI show the player's own level and both penalties, share unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. An NPC's level, penalties and private rolls never enter previews, receipts, errors or public history. Preserve legacy reads and unrelated JSON.

Hard prerequisites: the Exhaustion ledger task (`edda-270:exhaustion-ledger`); #61 check and save receipts; #65 attack roll; #66 Initiative and turn economy; #780 authoritative modified Speed; the Prone movement task (`edda-270:prone-movement`) and the flying fall task (`edda-270:flying-fall`), whose outputs the composition fixtures above use. Those last two edges also route both tasks to the capped #279. Parent: #270. Directly block #623, the existing direct #279 and #280 blocker that removes one level and recomputes these penalties; that single edge makes all four tasks of this decision mandatory for #279, #280 and, through #551 and #1001, #337. Add no direct #279 edge, no #67 or #429 consumer, no edge to #428 or #1001, and change no existing body, parent or other prerequisite.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
