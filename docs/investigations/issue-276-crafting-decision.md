# Nonmagical crafting, potion brewing, scroll scribing and magic-item crafting (#276)

No published issue executes any of the four named crafting procedures. A title-and-body search of all 1,892 issues found catalogue data owners for the tools and products and adjacent item owners, but no task that spends raw materials, counts workdays or creates an item. This decision therefore plans five bounded tasks under #276: one per procedure and one qualification task that carries the single new edge to item qualification #280. The tasks are drafted here and are not yet published; the native numbers, readback and receipt are added at publication.

## Source

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Source facts and fixtures retain the CC-BY-4.0 attribution and adaptation notice in [the source register](../srd-5.2.1-conformance.md). Every clause below was compared with the extracted page text; the companion scope file lists the exact phrases and both cost tables.

| Named scope | PDF | Source rule | Existing owners | Planned task |
|---|---|---|---|---|
| Nonmagical crafting | 93–94, 103 | Required tool and proficiency, helpers proficient too; raw materials worth half the purchase cost, rounded down, availability decided by the GM; cost in GP / 10 days of 8 hours, fraction rounded up, days need not be consecutive; time divided by the number of workers, normally one assistant | Data only: tools #253–#255, products #247–#252 and #256–#262 | `nonmagical-crafting` |
| Brewing Potions of Healing | 94, 99, 103 | Herbalism Kit proficiency, that kit, 25 GP of raw material, 1 day of 8 hours | Data only: #253, #261. Drinking and administering: #266 (closed) | `potion-brewing` |
| Scribing Spell Scrolls | 93, 103 | Ten-row time and cost table; 8 hours per day, days need not be consecutive; Arcana or Calligrapher's Supplies proficiency; spell prepared each day; Material components at hand, consumed only on completion; scroll uses the scribe's save DC and attack bonus; a cantrip works at the scribe's level | Data only: #217, #253, #262. Reading and consumption: #1142. Spellbook copying: #1770. Cleric preparation: #557 | `scroll-scribing` |
| Crafting magic items | 206–207 | Arcana for crafter and assistants; nine-row category tool table; all castable spells prepared every day; five-row rarity time and cost table, halved for a consumable other than a Spell Scroll; 8-hour days; assistants; GM-decided materials with 75 and 25 percent chances and a 7-day wait; incorporated item paid in full or crafted | Data only: #225–#246, #251, #253–#255. Item casting examples: #1140, #2051 | `magic-item-crafting` |

Adjacent issues that mention crafting words do not own these rules. #1726 resolves Ink, Ink Pen and Parchment use and states that it infers no scribing. #1870–#1872 resolve the Fabricate spell. #1447 validates components for a cast inside its 245-spell boundary, not materials held during inscription. #275 will own prices, services and selling; #278 remains the planning reconciliation for item keys. None of their bodies, parents or edges change.

## Planned tasks

| Key | Title | Size | Existing prerequisites | Planned prerequisites | Blocks |
|---|---|---|---|---|---|
| `nonmagical-crafting` | Craft one nonmagical item with a required tool, raw materials and workdays | M | #195, #248, #251, #254, #255 | — | — |
| `potion-brewing` | Brew one Potion of Healing with a Herbalism Kit | S | #195, #253, #261 | `nonmagical-crafting` | — |
| `scroll-scribing` | Scribe one Spell Scroll from a prepared spell | M | #195, #217, #253, #262, #557, #1142 | `nonmagical-crafting` | — |
| `magic-item-crafting` | Craft one magic item by rarity with Arcana, a category tool and prepared spells | M | #195, #225, #238, #245, #251, #253, #254, #255, #557 | `nonmagical-crafting` | — |
| `crafting-qualification` | Qualify the finite crafting, brewing and scribing outputs | M | #247–#262 | the four tasks above | #280 |

All five are Task issues and native children of #276. Each body is the exact text of its `issue-276-*-task.md` file; the scope file records each SHA-256. The bodies name their siblings by marker, so they need no edit after numbers are assigned.

The nonmagical task also owns the one crafting-project ledger (project start, 8-hour workdays, completion). The other three procedures reuse it, which is why they wait for it. Their fixtures are named: Light Hammer, Longbow and Plate Armor; the basic Potion of Healing; Sacred Flame, Sanctuary and Revivify scrolls scribed by a Cleric; Armor, +1 (Plate Armor), Potion of Resistance and Wand of Magic Detection. Wider coverage is not implied. The qualification task reconciles every tool Craft entry against the catalogue, proves the routing between the four procedures and the two-path +1 armor composition, and reports the rest as unavailable.

Scribing is the one place where a new output meets an existing owner. PDF99 and PDF244 fix a Spell Scroll's save DC and attack bonus by spell level, and #1142 implements that. PDF103 says a scribed scroll uses the scribe's values. The scribing task saves those values on the created scroll and presents them to the unchanged #1142 protocol; found scrolls keep the fixed table.

## Unknowns kept unavailable

- PDF103 does not state the denomination for rounding half the purchase cost down. Only prices that halve to whole GP are executable; others need a recorded GM amount.
- Neither PDF103 nor PDF207 rounds a day count that is not whole after dividing among workers or halving for a consumable. Such projects need a recorded GM day count.
- The source does not say when raw materials are paid, and defines no refund, abandonment, check or failure. The single debit at project start is labelled an Edda accounting adaptation; nothing else is invented.
- PDF103 asks for a prepared spell, while cantrips are known. The Cantrips clause is read as allowing a known cantrip; other cantrip sources stay unavailable.
- The brewing rule names the basic Potion of Healing. Brewing the stronger PDF236 variants is not stated and stays unavailable.

Two findings fall outside the four named scopes and are reported without a new task. Heavy Crossbow, the crafting example on PDF103, is on PDF91 but is missing from the #80 weapon inventory (37 names) and from #247–#250; the qualification task fails closed on it. The Magic Item Rarities and Values table on PDF206 has no executable owner; only its Spell Scroll clause is used, as a consistency witness against the PDF95 prices.

## Native edges and limits

Observed 2026-11-20: #280 has 96 of 100 blocked-by entries, and sibling decisions #274 and #275 still have to add theirs. The four procedure tasks therefore reach #280 through the one qualification task, the accepted shared-aggregate pattern, which takes #280 to 97. #276 goes from 0 to 5 sub-issues. The largest planned blocked-by list is 20. #195 goes from 61 to 65 dependents and #557 from 24 to 26; no prerequisite approaches 100. Spell qualification #279 is full at 100 and receives no edge: scribing creates an item and adds no spell effect. No task adds a #67 or #429 consumer.

No planned prerequisite is downstream of #276 or #280, so the dependency graph stays acyclic. #280 blocks the Cleric qualification #551, which already waits for #557; the new path #557 → scribing → qualification → #280 → #551 runs the same way. The five tasks are leaves of the parent hierarchy. Caps and both graphs must be read again at publication, because the sibling decisions may publish first.

Closing #276 supplies no crafting runtime. Removing the qualification task's #280 edge would lose the mandatory path for all five tasks. This change contains no runtime file, migration, UI behavior, production mutation or Full Rules admission. Focused source checks and hosted CI apply to this decision; the implementation tasks keep their own PostgreSQL, web, TUI and device requirements.
