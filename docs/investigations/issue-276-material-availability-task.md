<!-- edda-276:magic-item-material-availability -->
Decision #276; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 207; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The percentile mapping and the fixtures below are Edda adaptations.

## Bounded outcome

Implement only the raw-material availability check in "Crafting Magic Items" on PDF207: "The GM determines whether appropriate raw materials are available. In a city, there is a 75 percent chance that the materials are available, and in any other settlement, that chance is 25 percent. If materials aren't available, you must wait at least 7 days before checking on the availability again."

The check produces the availability ruling that the #276 magic-item crafting task (marker `edda-276:magic-item-crafting`) requires before a project starts. It owns no cost, tool, proficiency, workday or product rule.

The GM records a settlement fact for the crafter's location: city, or other settlement. The source names no chance outside a settlement, so no settlement fact means no check; the GM can still record a direct ruling under the crafting task. The GM may decide the outcome directly, or roll. For a roll, use one injected percentile and save it once (Edda adaptation: 1–75 succeeds in a city, 1–25 in another settlement). After a failed check for one crafter and item, a new check is allowed only when the authoritative campaign clock is at least 7 days later. The source does not say whether the wait is per item, per settlement or per crafter; the per-crafter-and-item scope is an Edda conservative bound, and a recorded GM ruling may replace it.

Nonmagical crafting on PDF103 has no percentage or wait; do not apply this check there.

## Deterministic fixtures

- City: injected 75 succeeds and 76 fails. Other settlement: 25 succeeds and 26 fails.
- After a failure on day 10, a recheck on day 16 rejects and one on day 17 is allowed.
- A success for the Wand of Magic Detection lets that project start; it grants nothing for another item.
- A GM direct ruling of "available" with no roll is saved as a ruling and consumes no dice.
- Reject independently: no settlement fact; a settlement fact from a non-GM actor; a second roll replacing a saved one; a client-supplied roll; a recheck inside 7 days; the check applied to a nonmagical item, a Potion of Healing or a Spell Scroll.
- Reject an invented price increase, bribe, travel shortcut, partial availability or Arcana check.

## Durable implementation acceptance

Preview is inert: it never rolls, writes, spends or reveals private facts. Confirm revalidates campaign, actor, source version, tool and material ownership, proficiency entitlement, authoritative campaign clock and expected revisions. Commit each project start, workday, payment, completion and created item with public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without a second payment, workday or item; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, explicit unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private inventory, GM rulings not yet disclosed and hidden facts.

Hard prerequisites: the #276 magic-item crafting task; the #276 crafting project ledger task (marker `edda-276:crafting-project-ledger`) for the saved ruling. Parent: #276. This task reaches item qualification #280 through the #276 crafting qualification task (marker `edda-276:crafting-qualification`), which holds the one direct #280 edge. Add no #279, #67 or #429 edge and alter no existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
