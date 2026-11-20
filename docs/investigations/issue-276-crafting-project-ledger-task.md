<!-- edda-276:crafting-project-ledger -->
Decision #276; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 103 and 206–207; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The project record, the payment point and the fixtures below are Edda adaptations, not SRD text.

## Bounded outcome

Implement the one durable crafting-project record that the four #276 procedures reuse: nonmagical crafting (marker `edda-276:nonmagical-crafting`), potion brewing (`edda-276:potion-brewing`), scroll scribing (`edda-276:scroll-scribing`) and magic-item crafting (`edda-276:magic-item-crafting`). This task executes no SRD procedure. It owns no formula, table, tool, proficiency, preparation or product rule; each procedure supplies an already validated project descriptor. Until a procedure task qualifies, no player can start a project.

Project record: one start, a fixed required day count, recorded workdays, one completion. The start saves the procedure kind, product identity, participants, required day count and cost. PDF103 and PDF207 both set a working day at 8 hours and say the days need not be consecutive. A workday records 8 hours for every participant on a campaign day the existing authoritative campaign clock has already reached. This task reads that clock and never advances it. The last required workday creates exactly one item instance through the #195 item contract, and creating it grants no use.

Payment: the source gives each procedure a cost but does not say when it is paid. Edda adaptation: the cost is paid once, in the transaction that starts the project. The payment is the single coin payment that decision #275 assigns (marker `edda-275:coin-payment-sale` in that decision). This task owns no coin ledger, denomination, conversion or change rule and reads no `srd521_gold` value itself. Amounts are exact Copper Piece counts supplied by the procedure. When that payment output is not published and qualified, every crafting project is unavailable. Do not add a second debit path.

Recorded GM rulings: no published issue owns one. This task stores a typed ruling on a project: its kind, value, the GM actor who made it, and whether it is disclosed. The procedures name the kinds they accept (raw-material availability, assigned tool, raw-material amount, day count, extra assistants, workday without components). The ledger stores and replays rulings and infers none. A ruling made by a non-GM actor rejects.

Edda conservative bounds, not SRD text: a participant records at most one 8-hour workday per campaign day across all projects, because the source defines a working day as 8 hours and says nothing about two in one day. The source defines no refund, abandonment, partial product, check or failure, so none exists: a project stays open until completed.

## Deterministic fixtures

Use a test-only procedure descriptor; it grants no production availability.

- A 1-day, 25 GP project: start pays 2,500 CP once; one workday completes it and creates one item.
- A 5-day, 150 GP project: record days 1–2, let the clock pass unrelated days, then record days 3–5. The item appears only with day 5.
- A 2-participant project records 8 hours for both participants on each workday.
- A saved GM day-count ruling is returned unchanged after reload and is hidden from public history until disclosed.
- Reject independently: coins that cannot cover the cost (no project is opened); a 7-hour day; a workday beyond the required count; a second workday for one participant on the same campaign day; a workday on a campaign day the clock has not reached; a workday by a non-participant; a second completion; a ruling from a non-GM actor; a foreign campaign or actor; a descriptor from an unqualified procedure.
- Reject an invented refund, partial item, progress percentage, check, failure roll or time advance.

## Durable implementation acceptance

Preview is inert: it never rolls, writes, spends or reveals private facts. Confirm revalidates campaign, actor, source version, tool and material ownership, proficiency entitlement, authoritative campaign clock and expected revisions. Commit each project start, workday, payment, completion and created item with public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without a second payment, workday or item; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, explicit unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private inventory, GM rulings not yet disclosed and hidden facts.

Hard prerequisites: #195 item contract; decision #275 for the coin payment. The `edda-275:coin-payment-sale` task has no native number yet. At publication, find it by marker and add it as a native blocker of this task; if it is not yet published, record that edge as pending in the link verification and add it when the number exists. Parent: #276. This task reaches item qualification #280 through the #276 crafting qualification task (marker `edda-276:crafting-qualification`), which holds the one direct #280 edge. Add no #279, #67 or #429 edge and alter no existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
