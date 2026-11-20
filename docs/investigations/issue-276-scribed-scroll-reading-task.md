<!-- edda-276:scribed-scroll-reading -->
Decision #276; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 99, 103, 159 and 244; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement one step for reading a Spell Scroll created by the #276 scroll scribing task (marker `edda-276:scroll-scribing`). PDF103 says a scribed scroll's spell uses the scribe's spell save DC and spell attack bonus, and a scribed cantrip works as if the caster were the scribe's level. PDF99 and PDF244 fix those values by spell level, and #1142 implements that table for every scroll.

This task owns the choice between the two. For a scroll instance that carries saved scribe values, use the saved DC, attack bonus and cantrip caster level in place of #1142's fixed table. For every other scroll, #1142's table applies unchanged. This lands after #1142 and calls its read path; it does not edit #1142's issue body, parent or edges, and it changes nothing else #1142 owns: identity, class-list permission, the Thief override, the higher-level ability check (DC 10 + spell level), casting time, interruption and consumption. It adds no spell effect. Effect dispatch still needs a qualified owner; a spell without one stays unavailable before any cost.

The saved values are read only from the scroll instance. Never take them from the client payload, the reader's sheet or the scribe's current sheet.

## Deterministic fixtures

- Sacred Flame scribed by a level-5 Cleric (WIS +3, PB +3) and read by a level-1 Cleric: Dexterity save DC 14 and 2d8 Radiant through the #957 resolver.
- A found Sacred Flame scroll without saved values, same reader: the #1142 fixed DC 13 and the unchanged #1142 and #957 behavior.
- A scribed Sanctuary scroll presents DC 14 and attack bonus +6 to dispatch. No Sanctuary effect owner is published, so the read is unavailable before consumption; the saved values stay on the scroll.
- The scribe later reaches level 11 and Wisdom 20: the scroll still reads at DC 14 and caster level 5.
- Reject independently: a scribed scroll falling back to DC 13 or +5; a found scroll using any scribe value; the reader's own DC, attack bonus or level; values supplied by the client; changed saved values; saved values applied to the DC 10 + level check; a foreign inventory.
- Reject an invented upcast, slot use, ritual read or permission that the saved values grant.

## Durable implementation acceptance

Preview is inert: it never rolls, writes, spends or reveals private facts. Confirm revalidates campaign, actor, source version, tool and material ownership, proficiency entitlement, authoritative campaign clock and expected revisions. Commit each project start, workday, payment, completion and created item with public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without a second payment, workday or item; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, explicit unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private inventory, GM rulings not yet disclosed and hidden facts.

Hard prerequisites: #1142 scroll read and consumption protocol; #957 Sacred Flame resolver; the #276 scroll scribing task for scroll instances with saved values. Parent: #276. This task reaches item qualification #280 through the #276 crafting qualification task (marker `edda-276:crafting-qualification`), which holds the one direct #280 edge. Add no #279, #67 or #429 edge and alter no existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
