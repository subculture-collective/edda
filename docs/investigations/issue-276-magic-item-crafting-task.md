<!-- edda-276:magic-item-crafting -->
Decision #276; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 92–94, 103, 123, 206–207, 210, 236 and 251; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement "Crafting Magic Items" on PDF206–207 for three named products: Armor, +1 (Plate Armor) from #225 and #251; Potion of Resistance from #238; Wand of Magic Detection from #245. Creating the item grants no use: wearing, drinking, charges, attunement and casting remain with their own owners. PDF206 sends Potion of Healing brewing and Spell Scroll scribing to the Equipment rules; route both to their #276 tasks and never through the rarity table.

Arcana (PDF206): the crafter and any assistants must have proficiency in the Arcana skill.

Tools (PDF206–207): use the tool the Magic Item Tools table requires for the item's category, and have proficiency with it; assistants need that proficiency too. Armor: Leatherworker's, Smith's or Weaver's Tools by kind of armor as the tools' descriptions note. Potion: Alchemist's Supplies or Herbalism Kit. Ring: Jeweler's Tools. Rod, Staff, Wand: Woodcarver's Tools. Scroll: Calligrapher's Supplies. Weapon: Leatherworker's, Smith's or Woodcarver's Tools by kind of weapon. Wondrous Item: Tinker's Tools or the tool required to make the nonmagical item on which the magic item is based. Retain all nine rows as data; execute only the three fixtures' rows.

Spells (PDF207): if the item lets its user cast any spells from it, the crafter must have all those spells prepared every day spent crafting. Use the actual Cleric preparation output #557 for Detect Magic (PDF123).

Time and cost (PDF207, all five rows retained): Common 5 days 50 GP; Uncommon 10 days 200 GP; Rare 50 days 2,000 GP; Very Rare 125 days 20,000 GP; Legendary 250 days 100,000 GP. Halve time and cost for a consumable item other than a Spell Scroll. Artifacts have no row and stay unavailable. Each day is 8 hours of work; days need not be consecutive. Assistants divide the time by the number of characters working; normally one other character, more only with a recorded GM allowance. The source gives no rounding for a result that is not a whole day (a Common or Very Rare consumable, or an odd count shared between two workers); such a project stays unavailable until a GM ruling records the day count.

Raw materials (PDF207): the table cost is the raw materials. The GM determines whether they are available. In a city the chance is 75 percent, in any other settlement 25 percent; if they are not available, the crafter must wait at least 7 days before checking again. Record the GM's settlement fact and outcome. When the GM chooses to roll, use an injected percentile (Edda adaptation: 1–75 succeeds in a city, 1–25 in another settlement) and save it once. No settlement fact means no check.

Incorporated item (PDF207): if the magic item incorporates an item that has a purchase cost, also pay that entire cost or craft that item under the Equipment rules. For the craft path, consume one actual item produced by the #276 nonmagical crafting task (marker `edda-276:nonmagical-crafting`); reuse its project ledger for the magic-item days.

Sentient item traits on PDF207 describe a GM-created persona and are not a crafting procedure here.

## Deterministic fixtures

- Armor, +1 (Plate Armor), Rare, Smith's Tools: 50 days and 2,000 GP, plus either 1,500 GP for the armor (3,500 GP in total, the source's own example) or one Plate Armor already crafted under the Equipment rules. One Arcana- and Smith's-proficient assistant makes it 25 days.
- Potion of Resistance, Uncommon consumable, Herbalism Kit or Alchemist's Supplies: 5 days and 100 GP. It casts no spell, so no preparation is required.
- Wand of Magic Detection, Uncommon, Woodcarver's Tools: 10 days and 200 GP, Detect Magic prepared on each recorded day. A day without it prepared rejects; later days resume.
- City availability: injected 75 succeeds and 76 fails; other settlement: 25 succeeds and 26 fails. After a failure, a recheck 6 days later rejects and one 7 days later is allowed.
- Reject independently: no Arcana; no tool proficiency; wrong category tool (Smith's Tools for the wand); assistant without Arcana or without the tool proficiency; second assistant without GM allowance; 1,999 GP; neither paying for nor supplying the Plate Armor; halving the wand; not halving the potion; a Common consumable (2.5 days) without a day-count ruling; an Artifact; Potion of Healing or Spell Scroll through this table; a 7-hour day; an extra day.
- Reject an invented crafting check, formula or recipe requirement, exotic ingredient quest, failure or curse roll, attunement grant or sale value.

## Durable implementation acceptance

Preview is inert: it never rolls, writes, spends or reveals private facts. Confirm revalidates campaign, actor, source version, tool and material ownership, proficiency entitlement, authoritative campaign clock and expected revisions. Commit each project start, workday, payment, completion and created item with public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without a second payment, workday or item; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, explicit unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private inventory, GM rulings not yet disclosed and hidden facts.

Hard prerequisites: #195 item contract; #225 Armor, +1 record; #238 Potion of Resistance record; #245 Wand of Magic Detection record; #251 Plate Armor record; #253 Alchemist's Supplies and Herbalism Kit records; #254 Smith's Tools record; #255 Woodcarver's Tools record; #557 Cleric preparation output; the #276 nonmagical crafting task. Parent: #276. This task reaches item qualification #280 through the #276 crafting qualification task (marker `edda-276:crafting-qualification`). Add no #279, #67 or #429 edge and alter no existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
