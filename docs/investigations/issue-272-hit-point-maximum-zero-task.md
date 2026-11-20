<!-- edda-272:hit-point-maximum-zero -->
Decision #272; roadmap #47. One implementation PR; Task / Size S.

Source: English SRD 5.2.1, PDF 17, 185, 187, 341; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only this general rule (PDF17): a creature dies if its Hit Point maximum reaches 0. The actual witness is Wight Life Drain #1705, a published reducer whose failed save lowers the target's Hit Point maximum by the damage taken. Do not copy its save, dice or Resistance handling. Every reducer keeps its own amount, trigger and restoration; this task adds only the death predicate that reads the committed maximum. No new reducer, restoration or healing effect is enabled.

When a committed reduction leaves the creature's Hit Point maximum at 0, the creature dies in that same transaction. The rule applies to characters and monsters alike and bypasses the 0 Hit Point flow: no Death Saving Throw, no Unconscious instance and no Stable state is created for that event. A maximum of 1 or more never triggers it, even when current Hit Points are 0. Use the existing authoritative dead status that the closed attack slice #65 already commits in `internal/combat`; do not create a second death model. Edda adaptation: a reduction larger than the remaining maximum stores 0, never a negative maximum.

Recovery boundary: death is not a reduced maximum waiting for a rest. A dead creature cannot start a Short or Long Rest (each needs at least 1 Hit Point), so #428's restoration of a reduced maximum never applies to it, and neither #428 nor Greater Restoration #627 returns it to life. Returning to life stays with its own owners, such as Revivify #618, and is otherwise unavailable.

## Deterministic fixtures

- Maximum 7, Hit Points 7, failed Life Drain save with injected d8 = 5 for 7 Necrotic damage: Hit Points 0, maximum 0, dead in that receipt with no Death Saving Throw counters and no Unconscious instance.
- Maximum 8, Hit Points 8, the same 7 damage: Hit Points 1, maximum 1, alive.
- Maximum 20, Hit Points 3, the same 7 damage: Hit Points 0, maximum 13. This rule does not fire; the ordinary 0 Hit Point flow applies unchanged.
- Maximum 5, Hit Points 5, injected d8 = 8 for 10 damage halved by Resistance to 5, as #1705 computes: maximum 0, dead.
- A successful save deals no damage, lowers nothing and kills nothing.
- Reject death at a maximum of 1; a death deferred to a later turn, rest or dawn; a Death Saving Throw, Stable state or 1d4-hour recovery for this death; a negative stored maximum; a client-declared maximum or reduction; a second death on replay; a rest started by the dead creature; and revival by a Long Rest or by restoring the maximum.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, actor, source version, clock, targets, the reducer's committed result, the current maximum and expected revisions. Commit the reduction, Hit Points, the death and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the same receipt with no reroll; a conflicting payload rejects; racing confirmations commit once. The death belongs to the triggering confirmed reduction, with no optional activation. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires a focused build, lint and desktop/phone checks. Public history omits private HP and maximums, hidden actors, private checks and undisclosed item facts.

Hard prerequisites: #65 committed damage and dead status, #1705 actual Hit Point maximum reduction. Parent: #272. Block the rest/recovery qualification task published with this decision; that task blocks #280. Do not add a direct #279, #280, #67 or #429 edge and do not alter any existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this task and its callers qualify. No deployment, source inference or Full Rules promotion.
