<!-- edda-276:nonmagical-crafting -->
Decision #276; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 91–94 and 103; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures and the choices marked as Edda adaptations below are not SRD text.

## Bounded outcome

Implement the "Crafting Nonmagical Items" procedure on PDF103 for three named products: Light Hammer and Longbow from #248 and Plate Armor from #251, with Smith's Tools from #254 and Woodcarver's Tools from #255. A catalogue record, a tool's Craft entry or this task's closure never makes another item craftable. Every other Craft entry belongs to the #276 Craft roster spike (marker `edda-276:craft-roster-spike`) and stays unavailable with that reason.

Tools (PDF93–94, 103): the crafter must use the required tool and have proficiency with it. The required tool comes from the tool's own Craft entry. Smith's Tools: any Melee weapon except Club, Greatclub, Quarterstaff and Whip; Medium armor except Hide; Heavy armor. Woodcarver's Tools: Ranged weapons except Pistol, Musket and Sling. The #254 and #255 issue bodies do not say their records carry the Craft entry. If the records lack it, this task adds those two entries as procedure data with their PDF94 text, limited to the class entries above; it adds no catalogue record and no other tool's entry. Require the actual tool in the crafter's authoritative inventory. Anyone who helps must also have proficiency with that tool; the source does not say a helper needs a second copy, so do not require or invent one.

For an item no Craft entry lists, PDF103 says the GM assigns the tool. Without a recorded GM assignment the item is unavailable. The source does not say what happens when two tools list one product; none of the three products is listed twice, and such products stay with the roster spike.

Raw materials (PDF103): worth half the purchase cost, rounded down. The GM determines whether appropriate raw materials are available; require that recorded ruling and never infer it. The three fixtures halve exactly (1 GP, 25 GP, 750 GP). The source does not name the unit the rounding uses. A half that is an exact number of Copper Pieces needs no rounding. A price whose half is not a whole number of Copper Pieces stays unavailable until a GM ruling records the amount.

Time (PDF103): days of 8 working hours equal the purchase cost in GP divided by 10, with a fraction rounded up to a day. Days need not be consecutive. Characters can combine efforts: divide the time by the number of characters working. Normally only one other character can assist; more assistants need a recorded GM allowance. The source gives no rounding for a shared time that is not a whole day; such a project stays unavailable until a GM ruling records the day count.

Completion (PDF103): when every requirement is met the crafter makes the item and can use it. Selling it at the normal price belongs to decision #275.

Project, workdays, payment and rulings: reuse the #276 crafting project ledger task (marker `edda-276:crafting-project-ledger`). This task supplies the validated descriptor (product, tool, participants, day count, raw-material amount) and adds no second project record or debit.

## Deterministic fixtures

- Light Hammer, 2 GP, Smith's Tools: raw materials 1 GP; 2 / 10 = 0.2 rounds up to 1 day.
- Longbow, 50 GP, Woodcarver's Tools: raw materials 25 GP; 5 days. Record days 1–2, an unrelated gap, then days 3–5; the item appears only on day 5.
- Plate Armor, 1,500 GP, Smith's Tools: raw materials 750 GP; 150 days alone, 75 days with one proficient assistant.
- The source's own Heavy Crossbow example (50 GP, 5 days) is checked as arithmetic only. Heavy Crossbow has no catalogue identity in the #80 inventory, so it is not craftable here.
- Reject independently: no proficiency; tool absent; a tool whose Craft entry excludes the item (Hide Armor with Smith's Tools, Musket with Woodcarver's Tools); Shield, which no Craft entry lists, without a GM tool assignment; an unproficient assistant; a second assistant without GM allowance; 24 GP for the Longbow; missing GM material ruling; Longbow with one assistant (2.5 days) without a day-count ruling; any other product, with the roster spike named as the reason.
- Reject an invented ability check, quality roll, masterwork bonus, Fabricate shortcut or tool granted by the project.

## Durable implementation acceptance

Preview is inert: it never rolls, writes, spends or reveals private facts. Confirm revalidates campaign, actor, source version, tool and material ownership, proficiency entitlement, authoritative campaign clock and expected revisions. Commit each project start, workday, payment, completion and created item with public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without a second payment, workday or item; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, explicit unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private inventory, GM rulings not yet disclosed and hidden facts.

Hard prerequisites: #195 item contract; #248 Light Hammer, Longbow and Musket records; #251 Plate Armor, Hide Armor and Shield records; #254 Smith's Tools record; #255 Woodcarver's Tools record; the #276 crafting project ledger task (marker `edda-276:crafting-project-ledger`). Parent: #276. This task reaches item qualification #280 through the #276 crafting qualification task (marker `edda-276:crafting-qualification`), which holds the one direct #280 edge. Add no #279, #67 or #429 edge and alter no existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
