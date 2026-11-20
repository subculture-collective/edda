<!-- edda-272:wand-of-magic-missiles-dawn -->
Decision #272; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 206, 251; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only the dawn recharge of one source identity, `srd-5.2.1:content/magic-items/wand-of-magic-missiles`: the wand regains 1d6 + 1 expended charges daily at dawn, never above its 7 charges. The cast, the charge spend and the last-charge destruction are the separate wand task published with this decision and are consumed here, not copied. No other item or recharge formula is enabled.

#1148 is a build-order prerequisite, not a published generic seam: it stays limited to the Wand of Magic Detection and its 1d3 formula. This adapter extends the authoritative dawn event that task introduces instead of creating a second dawn clock. Neither #1148's body nor its item domain is widened. On a world or plane where dawn does not occur, the GM determines when the item recharges (PDF206): require a trusted explicit GM recharge event and never fabricate time.

A wand that survived its last-charge roll sits at 0 charges and recharges at the next dawn like any other count. A destroyed wand no longer exists as an item and never recharges. A rest is not a recharge: neither a Short Rest nor a Long Rest changes the count, whether or not a dawn falls inside the rest.

## Deterministic fixtures

- Dawn at 0 charges with injected d6 = 3 gives 4. Dawn at 5 with d6 = 6 gives 7, not 12. Dawn at 7 stays 7.
- A wand at 0 charges after a last-charge roll of 2 recharges at dawn; a wand destroyed by a roll of 1 does not.
- A completed Long Rest with no dawn event leaves 2 charges at 2. A dawn event during that rest applies once.
- On a plane without dawn, a trusted GM recharge event with injected d6 = 1 turns 3 charges into 5; without that event nothing recharges.
- Reject recharge on rest start or rest completion, a client clock, a duplicate dawn for the same day, a reroll, a foreign wand, a recharge of a destroyed wand, the 1d3 formula of #1148, a total above 7, and a fabricated dawn on a plane without one.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, actor, source version, clock, item identity and existence, charges, the dawn or GM recharge event and expected revisions. Commit the event, die, recharge and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same event ID and payload after retry, reload or reconnect returns the same receipt with no reroll; a conflicting payload rejects; racing events commit once. The recharge belongs to its triggering authoritative event, with no optional activation. Dice are injectable. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires a focused build, lint and desktop/phone checks. Public history omits private HP, hidden actors, private checks and undisclosed item facts, including an undisclosed charge count.

Hard prerequisites: #1148 authoritative dawn event (build order; it brings its own prerequisites #195, #1140 and the Rogue execution cohort #1097) and the Wand of Magic Missiles cast task published with this decision. Parent: #272. Block the rest/recovery qualification task published with this decision; that task blocks #280. Do not add a direct #279, #280, #67 or #429 edge and do not alter any existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this task and its callers qualify. No deployment, source inference or Full Rules promotion.
