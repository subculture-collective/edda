# Nonmagical crafting, potion brewing, scroll scribing and magic-item crafting (#276)

No published issue executes any of the four named crafting procedures. A title-and-body search of all 1,892 issues found catalogue data owners for the tools and products and adjacent item owners, but no task that spends raw materials, counts workdays or creates an item. This decision publishes eleven children under #276, #2114–#2124: seven bounded implementation tasks, three named coverage spikes for the wider scope, and one qualification task, #2124, that carries the edge to item qualification #280. #280 was full at 100 blockers, so #2124 took the place of #276's own edge there; see Native edges and limits. The readback is in [the link verification](issue-276-link-verification.json).

## Source

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Source facts and fixtures retain the CC-BY-4.0 attribution and adaptation notice in [the source register](../srd-5.2.1-conformance.md). Every clause below was compared with the extracted page text; the companion scope file lists the exact phrases and both cost tables. The PDF column is the union of the Source lines of the planned tasks in that row.

| Named scope | PDF | Source rule | Existing owners | Planned children |
|---|---|---|---|---|
| Shared project record | 103, 206–207 | A working day is 8 hours and days need not be consecutive, for all four procedures | None. Item contract #195; the runtime campaign clock initialized under #34 | `crafting-project-ledger` |
| Nonmagical crafting | 91–94, 103 | Required tool and proficiency, helpers proficient too; raw materials worth half the purchase cost, rounded down, availability decided by the GM; cost in GP / 10 days of 8 hours, fraction rounded up; time divided by the number of workers, normally one assistant | Data only: tools #253–#255, products #247–#252 and #256–#262 | `nonmagical-crafting`, `craft-roster-spike` |
| Brewing Potions of Healing | 94, 99, 103, 206–207, 236 | Herbalism Kit proficiency, that kit, 25 GP of raw material, 1 day of 8 hours | Data only: #253, #261. Drinking and administering: #266 (closed) | `potion-brewing` |
| Scribing Spell Scrolls | 93, 99, 103, 159, 206–207, 244 | Ten-row time and cost table; Arcana or Calligrapher's Supplies proficiency; spell prepared each day; Material components at hand, consumed only on completion; scroll uses the scribe's save DC and attack bonus; a cantrip works at the scribe's level | Data only: #217, #242, #253, #262. Fixed-table reading and consumption: #1142. Owned materials: #1447, #618. Cleric preparation: #557. Spellbook copying: #1770 | `scroll-scribing`, `scribed-scroll-reading`, `scribing-coverage-spike` |
| Crafting magic items | 92–94, 103, 123, 206–207, 210, 236, 251 | Arcana for crafter and assistants; nine-row category tool table; all castable spells prepared every day; five-row rarity time and cost table, halved for a consumable other than a Spell Scroll; GM-decided materials with 75 and 25 percent chances and a 7-day wait; incorporated item paid in full or crafted | Data only: #202, #225–#246, #251, #253–#255. Item casting examples: #1140, #2051 | `magic-item-crafting`, `magic-item-material-availability`, `magic-item-coverage-spike` |
| Roster, routing and composition | 91–95, 99, 103, 206–207, 244 | Tool Craft entries; Spell Scroll value is double its scribing cost; PDF206 routes brewing and scribing to the Equipment rules | None | `crafting-qualification` |

Adjacent issues that mention crafting words do not own these rules. #1726 resolves Ink, Ink Pen and Parchment use and states that it infers no scribing. #1870–#1872 resolve the Fabricate spell. #275 owns prices, coins, services and selling; #278 remains the planning reconciliation for item keys. None of their bodies, parents or edges change.

## Clauses that needed an explicit owner

| Clause or need | Owner | Why |
|---|---|---|
| Paying a crafting cost | #2083, Confirm coin holdings and exact listed-price payment (decision #275, marker `edda-275:coin-holdings-payment`), called once by `crafting-project-ledger` | The runtime stores `srd521_gold` on a Coin Pouch row and nothing spends it. #2083 is the only published coin ledger, so #276 adds no second debit |
| Recorded GM rulings | `crafting-project-ledger` | A search for "GM ruling", "recorded GM" and "GM allowance" over all issue bodies found no owner |
| Tool Craft entries as data | `nonmagical-crafting` for Smith's and Woodcarver's Tools; `craft-roster-spike` for the rest | The #253–#255 bodies do not say their records carry Craft entries. The qualification task reports a missing entry as absent data |
| Spell Scroll record for levels 2–9 | #242 (PDF244 family) | #262 transcribes only the PDF99 entry, which the source defines as Spell Scroll (Cantrip) or (Level 1) |
| Owned priced and consumed components | #1447 for the ownership model, #618 for the Revivify diamond | #1447's cast boundary does not list Sanctuary or Revivify, so scribing reuses its ownership validation and no cast admission |
| Reading a scribed scroll with the scribe's values | `scribed-scroll-reading` | #1142 fixes DC and attack bonus by spell level and has no input for saved values. Its body, parent and edges stay untouched; the new task selects the saved values before that table and lands after #1142 |
| Detect Magic record for the wand | #202 | #245 and #557 do not carry the spell record |
| Wider coverage | The three spikes | #276 asks to keep wider scope as explicitly named bounded spikes |

## Published children

| Issue | Key | Title | Size | Existing prerequisites | Planned prerequisites | Blocks |
|---|---|---|---|---|---|---|
| #2114 | `crafting-project-ledger` | Record one crafting project: start, 8-hour workdays, payment and completion | M | #195, #275, #2083 | — | — |
| #2115 | `nonmagical-crafting` | Craft one nonmagical item with a required tool, raw materials and workdays | M | #195, #248, #251, #254, #255 | `crafting-project-ledger` | — |
| #2116 | `potion-brewing` | Brew one Potion of Healing with a Herbalism Kit | M | #195, #253, #261 | `crafting-project-ledger` | — |
| #2117 | `scroll-scribing` | Scribe one Spell Scroll from a prepared Cleric spell | M | #195, #217, #242, #253, #262, #557, #618, #1142, #1447 | `crafting-project-ledger` | — |
| #2118 | `scribed-scroll-reading` | Read one scribed Spell Scroll with its saved scribe values | M | #957, #1142 | `scroll-scribing` | — |
| #2119 | `magic-item-crafting` | Craft one magic item by rarity with Arcana, a category tool and prepared spells | M | #195, #202, #225, #238, #245, #251, #253, #254, #255, #557 | `crafting-project-ledger`, `nonmagical-crafting` | — |
| #2120 | `magic-item-material-availability` | Check magic-item raw material availability by settlement | M | — | `crafting-project-ledger`, `magic-item-crafting` | — |
| #2121 | `craft-roster-spike` | Split the remaining tool Craft entries into bounded crafting PRs | M | — | — | — |
| #2122 | `scribing-coverage-spike` | Split Spell Scroll scribing for other classes and spell levels | M | — | — | — |
| #2123 | `magic-item-coverage-spike` | Split magic-item crafting for the remaining categories and rarities | M | — | — | — |
| #2124 | `crafting-qualification` | Qualify the finite crafting, brewing and scribing outputs | M | #242, #247–#262 | the ten children above | #280 |

All eleven are Task issues and native children of #276. The three spikes also take the Planning/Spike label, as #276 has. Each body is the exact text of its `issue-276-*` file; the scope file records each SHA-256. The bodies name their siblings by marker, so they need no edit after numbers are assigned. Every child is Size M, the size all published children and #276 itself carry.

The ledger is its own task because four procedures share it and no one of them should be the hidden owner of project records, payment timing or GM rulings. The availability check is split from magic-item crafting, and scribed-scroll reading from scribing, so that each implementation stays one reviewable PR.

The implementation tasks execute named fixtures only: Light Hammer, Longbow and Plate Armor; the basic Potion of Healing; Sacred Flame, Sanctuary and Revivify scrolls scribed by a single-class Cleric; Armor, +1 (Plate Armor), Potion of Resistance and Wand of Magic Detection. Nothing else becomes craftable when they close. The remaining Craft entries, the other classes and spell levels, and the other magic-item categories and rarities each have a named spike. Each spike blocks the qualification task and must add its own children as blockers of that task before it closes, so #280 stays blocked until the wider scope has owners that qualify.

## Unknowns kept unavailable

- PDF103 does not name the unit for rounding half the purchase cost down. A half that is an exact number of Copper Pieces needs no rounding; any other price needs a recorded GM amount.
- Neither PDF103 nor PDF207 rounds a day count that is not whole after dividing among workers or halving for a consumable. Such projects need a recorded GM day count.
- PDF103 says "You must use the required tool" and does not say what applies when two tools list one product. No executed fixture is listed twice; those products stay with the Craft roster spike.
- The source does not say when a cost is paid, and defines no refund, abandonment, check or failure. Paying once at project start is an Edda adaptation; nothing else is invented.
- PDF103 asks for a prepared spell, while cantrips are known. A scribe's own qualified Cleric cantrip is accepted as an Edda reading; other cantrip sources stay unavailable.
- PDF103 ties "on each day of the inscription" to preparation only and gives no timing for Material components. Completion requires them, because consumption is fixed there. An earlier workday without them needs a recorded GM ruling; that is an Edda conservative bound.
- The source does not say a Component Pouch or Spellcasting Focus can replace a component during scribing, so neither does.
- PDF93 lists Spell Scroll under Calligrapher's Supplies, and PDF103 names only proficiency. A scribe who qualifies only through the supplies must hold them or have a GM ruling; a scribe with Arcana needs no tool. This is an Edda conservative bound.
- PDF103 says a scribed cantrip works as if the caster were "your level". The scribe's total character level is saved, and only single-class scribes execute; multiclass scribes stay with the scribing coverage spike.
- The brewing rule names the basic Potion of Healing. Brewing the stronger PDF236 variants is not stated and stays unavailable.
- PDF207 does not say whether the 7-day wait is per item, settlement or crafter, and names no chance outside a settlement. The check is tracked per crafter and item as an Edda bound; no settlement fact means no check.
- The source does not say whether a character can work two 8-hour days in one day. One workday per participant per campaign day is an Edda bound.

Two findings fall outside the four named scopes and are reported without a new task. Heavy Crossbow, the crafting example on PDF103, is on PDF91 but is missing from the #80 weapon inventory (37 names) and from #247–#250; the qualification task fails closed on it and the Craft roster spike names it. The Magic Item Rarities and Values table on PDF206 has no executable owner; only its Spell Scroll clause is used, as a consistency witness against the PDF95 prices.

## Native edges and limits

**Standing gate policy.** The qualification gates #279 and #280 are at the native cap of 100 blocked-by entries. The orchestrator decided the following for every splitter. A closed splitter's own edge on a gate, or the edge of one that closes in the publishing PR, carries no remaining obligation once its terminal qualification child exists. When a gate is full, that edge is replaced by the terminal qualification child: remove "gate blocked-by splitter", add "gate blocked-by terminal", and the count stays at 100. No other edge is removed, the count never exceeds 100, no open splitter is made blocked by a runtime implementation task, and both graphs stay acyclic.

**Applied to #280.** The review read had #280 at 96 blockers, and the draft kept the #276 edge and added the qualification task as a 97th. Decisions #270, #272, #274 and #275 published first and filled #280. This publication therefore removed "#280 blocked-by #276" and added "#280 blocked-by #2124" in consecutive calls, after every other write had succeeded. #280 was at 99 only between those two calls and is at 100 again; its other 99 blockers are unchanged. #276 closes with this PR, and its acceptance criteria ask for exact implementation blockers on item qualification: #2124 is that blocker, and the ten other children reach #280 only through it. The #280 body still lists #276 in its own text and was not edited.

#276 goes from 0 to 11 sub-issues and no longer blocks anything. The largest blocked-by list is 27, on #2124. #195 goes from 68 to 73 dependents, #557 from 24 to 26 and #1142 from 2 to 4; no prerequisite approaches 100. Spell qualification #279 is full at 100 and receives no edge: scribing creates an item and adds no spell effect. No child adds a #67 or #429 consumer.

**The coin payment edge.** The reviewed draft named a marker, `edda-275:coin-payment-sale`, for a coin task with no native number, and treated the edge as pending. Decision #275 published that task as #2083, Confirm coin holdings and exact listed-price payment, with the marker `edda-275:coin-holdings-payment`. The ledger body and the plan were corrected before publication, and #2114 is blocked by #2083 as an ordinary edge. #2114 also keeps its #275 edge, which is satisfied because #275 is closed. The ledger body says every project is unavailable until the payment output qualifies.

No prerequisite of the eleven children is downstream of #280, whose live downstream set is #337, #551, #630, #1093, #2053, #2054 and #2056, so the dependency graph stays acyclic. #618 already blocks #280 and #557 already blocks #551; the new paths through the qualification task run the same way. The planned edges among the eleven children follow one order, ledger first and qualification last. The children are leaves of the parent hierarchy until a spike publishes its own. Both graphs were read in full before and after publication.

Closing #276 supplies no crafting runtime. Removing #2124's #280 edge would lose the mandatory path for all eleven children. This change contains no runtime file, migration, UI behavior, production mutation or Full Rules admission. Focused source checks and hosted CI apply to this decision; the implementation tasks keep their own PostgreSQL, web, TUI and device requirements.

## Publication

Published 2026-11-20, read back at 16:12:54Z. The eleven issues are #2114–#2124 in the order of the table above. Each body matches its task file byte for byte, and each is an open Task of Size M on the edda Project with #276 as parent; #2121, #2122 and #2123 carry the Planning/Spike label. The publication added 11 parent links and 68 blocked-by edges (67 planned, 1 on the gate) and removed one edge, "#280 blocked-by #276". Both ends of every added edge were re-read immediately before it was added.

A full page-through of 1,951 issues after publication shows the dependency graph (10,562 edges) and the parent hierarchy (1,947 edges) acyclic, with no issue above 100 blockers, dependents or children. The read taken just before publication showed 1,940 issues and 10,495 edges, also acyclic. Between the two reads no earlier issue changed title, type, parent or body, and the only edge changes are the ones listed here. #225 and #275 closed between the review read and publication; both stay native blockers and are satisfied.

The search figure of 1,892 issues is the review read. Source clauses were checked against the pinned PDF during review, not again at publication. The receipt is [issue-276-link-verification.json](issue-276-link-verification.json).
