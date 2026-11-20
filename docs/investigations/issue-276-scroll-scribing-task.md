<!-- edda-276:scroll-scribing -->
Decision #276; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 93, 99, 103, 159, 206–207 and 244; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only "Scribing Spell Scrolls" on PDF103 for three named Cleric spells from #217: Sacred Flame (cantrip), Sanctuary (level 1) and Revivify (level 3). The created item is the Spell Scroll identity of #262 and #1142. This task owns the scribing project and the scribe-bound operands saved on the created scroll. It adds no spell effect, no class grant and no new scroll-reading rule; reading and consumption remain #1142.

Time and cost (PDF103, all ten rows retained as data): cantrip 1 day 15 GP; level 1 1 day 25 GP; level 2 3 days 100 GP; level 3 5 days 150 GP; level 4 10 days 1,000 GP; level 5 25 days 1,500 GP; level 6 40 days 10,000 GP; level 7 50 days 12,500 GP; level 8 60 days 15,000 GP; level 9 120 days 50,000 GP. Each day of inscription is 8 hours of work; days need not be consecutive. PDF206–207 keep Spell Scrolls out of the magic-item rarity table and its consumable halving. The scribing rule names no assistants and no raw-material availability ruling; add neither. Reuse the crafting-project ledger from the #276 nonmagical crafting task (marker `edda-276:nonmagical-crafting`).

Prerequisites for the scribe (PDF103): proficiency in the Arcana skill or with Calligrapher's Supplies, and the spell prepared on each day of the inscription. Use the actual Cleric preparation output #557. A day on which the spell is not prepared cannot be recorded. The rule requires proficiency, not possession of the supplies; do not require the tool item. A cantrip is scribed under the Cantrips clause: accept it only when it is among the scribe's own qualified known cantrips on that day. The scribe must also have at hand any Material components the spell requires on each recorded day. A consumed Material component is consumed only when the scroll is completed, in the completing transaction.

Scribe-bound operands (PDF103): the scroll's spell uses the scribe's spell save DC and spell attack bonus. For a cantrip, the version on the scroll works as if the caster were the scribe's level. Save these values on the created scroll instance at completion, with the class that supplied the preparation. They never change afterwards. PDF99 and PDF244 give fixed values by spell level for a Spell Scroll; #1142 keeps those for every scroll without saved scribe operands. For a scribed instance, present the saved operands to the unchanged #1142 protocol; do not edit, copy or widen #1142.

The scroll bears the spell at its own level. The source gives no rule for scribing at a higher level, from an item, from a ritual book entry or from another caster's preparation; all stay unavailable. Copying a scroll into a spellbook remains #1770.

## Deterministic fixtures

- A level-5 Cleric with Arcana proficiency, Sanctuary prepared and its mirror shard at hand: 1 day, 25 GP, one level-1 scroll saving the scribe's own DC and attack bonus (for WIS +3 and PB +3: DC 14, +6). The shard is not consumed.
- The same Cleric with Calligrapher's Supplies proficiency instead of Arcana scribes Sacred Flame: 1 day, 15 GP, saved caster level 5.
- Revivify: 5 days, 150 GP, a diamond worth 300+ GP at hand on each day. Record days 1–3, unprepare the spell, and day 4 rejects; prepare it again and days 4–5 complete. The diamond leaves inventory only in the completing transaction.
- A found scroll without saved operands still reads at the fixed #1142 values; a scribed one never falls back to DC 13 / +5.
- Reject independently: neither proficiency; spell not prepared; another class's or another character's preparation; Material component absent on a day; diamond consumed at start; 24 GP for level 1; a 7-hour day; a sixth day; an assistant dividing the time; a level-2 scroll of Sanctuary; an unknown spell identity; rewriting saved operands after the scribe levels up.
- Reject an invented Arcana check, failure or mishap roll, ink or parchment consumption, slot expenditure or sale.

## Durable implementation acceptance

Preview is inert: it never rolls, writes, spends or reveals private facts. Confirm revalidates campaign, actor, source version, tool and material ownership, proficiency entitlement, authoritative campaign clock and expected revisions. Commit each project start, workday, payment, completion and created item with public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without a second payment, workday or item; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, explicit unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private inventory, GM rulings not yet disclosed and hidden facts.

Hard prerequisites: #195 item contract; #217 Sacred Flame, Sanctuary and Revivify records; #253 Calligrapher's Supplies record; #262 Spell Scroll record; #557 Cleric preparation output; #1142 scroll identity and consumption protocol; the #276 nonmagical crafting task for the shared project ledger. Parent: #276. This task reaches item qualification #280 through the #276 crafting qualification task (marker `edda-276:crafting-qualification`). It adds an item-creation mode, not a spell effect; add no #279, #67 or #429 edge and alter no existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
