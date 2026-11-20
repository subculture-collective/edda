<!-- edda-276:potion-brewing -->
Decision #276; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 94, 99, 103, 206–207 and 236; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only "Brewing Potions of Healing" on PDF103. A character who has proficiency with the Herbalism Kit creates one Potion of Healing by using that kit and 25 GP of raw material over the course of 1 day (8 hours of work). The product is the basic Potion of Healing on PDF99 (50 GP, 2d4 + 2). Its drinking and administering stay with the existing potion fixtures, including #266; this task adds no healing resolver.

Reuse the #276 crafting project ledger task (marker `edda-276:crafting-project-ledger`): one project, one payment of 25 GP, one 8-hour workday, one created potion. The brewing clause is specific and replaces the general formula. Do not derive 5 days from 50 GP / 10. The brewing clause names no assistant rule, no availability ruling and no check; add none. Require the actual Herbalism Kit in the brewer's authoritative inventory and the brewer's own proficiency with it.

The Herbalism Kit's other Craft entries (Antitoxin, Candle, Healer's Kit) are nonmagical items under the general procedure and belong to the #276 Craft roster spike (marker `edda-276:craft-roster-spike`). The Magic Item Tools table on PDF207 lists Alchemist's Supplies or Herbalism Kit for Potions made under the magic-item rules; PDF206 sends Potion of Healing brewing to the Equipment rule, which names only the Herbalism Kit. Greater, superior and supreme Potions of Healing (PDF236) are not this product: the brewing rule does not name them, and they stay unavailable.

## Deterministic fixtures

- Proficient brewer with a Herbalism Kit and 25 GP: after one confirmed 8-hour workday, the brewer has paid 25 GP once and inventory gains exactly one basic Potion of Healing.
- Two separate projects on two days cost 50 GP and yield two potions; one project never yields two.
- Reject independently: no Herbalism Kit proficiency; kit absent; 24 GP; Alchemist's Supplies substituted for the kit; the 5-day general formula; an assistant halving the day; a greater Potion of Healing; a foreign actor or inventory.
- Reject an invented check, failure chance, ingredient search roll or batch brewing.

## Durable implementation acceptance

Preview is inert: it never rolls, writes, spends or reveals private facts. Confirm revalidates campaign, actor, source version, tool and material ownership, proficiency entitlement, authoritative campaign clock and expected revisions. Commit each project start, workday, payment, completion and created item with public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without a second payment, workday or item; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, explicit unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private inventory, GM rulings not yet disclosed and hidden facts.

Hard prerequisites: #195 item contract; #253 Herbalism Kit and Alchemist's Supplies records; #261 Potion of Healing record; the #276 crafting project ledger task (marker `edda-276:crafting-project-ledger`). Parent: #276. This task reaches item qualification #280 through the #276 crafting qualification task (marker `edda-276:crafting-qualification`), which holds the one direct #280 edge. Add no #279, #67 or #429 edge and alter no existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
