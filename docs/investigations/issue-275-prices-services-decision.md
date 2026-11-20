# Coins, selling, lifestyle and priced services (#275)

No published issue owns coin denominations, payment, purchase, sale, lifestyle, lodging, hirelings or spellcasting services. Thirteen creation owners own a once-only Gold Piece grant and spend nothing. A search of all 1,892 issue titles and bodies, open and closed, found nothing else beyond #275 itself, roadmap lists that name it, costly spell components and a few boundary mentions. This decision therefore assigns the five named scopes to four bounded implementation tasks, names one bounded follow-up for the magic-item value table that Selling Equipment points to, and adds one qualification task. All six become native children of #275. The qualification task is the single new blocker of item qualification #280. Their exact bodies are the six companion task files; issue numbers, publication and readback are recorded later in the link-verification receipt. Until then nothing here is published, and closing #275 supplies no runtime.

## Source

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Source facts and fixtures retain the CC-BY-4.0 attribution and adaptation notice in [the source register](../srd-5.2.1-conformance.md). The issue names PDF 89 and 101–102. Cross-references were read for fixtures and boundaries: Barbarian starting gold (28), Longsword (91), Plate Armor (92), Candle (95–96), Potion of Healing (95, 99, 236), crafted-item sale (103), Identify (142), Revivify (159), Magic Item Values by Rarity (205–206) and Armor, +1 (210).

| Named scope | PDF | Source rule | Owner |
|---|---|---|---|
| Coins | 89 | CP 1/100, SP 1/10, EP 1/2, GP 1 and PP 10 GP; fifty coins weigh a pound | coin task |
| Selling | 89 | Equipment fetches half its cost; trade goods and valuables keep full value; magic-item prices are in "Magic Items" | sale task |
| Lifestyle | 101 | Seven lifestyles from free to 10 GP per day, chosen and paid at the start of each week or month (GM's choice); no inherent consequences | lifestyle task |
| Food, drink, lodging | 101–102 | Seventeen priced rows; daily lodging and meals are included in the lifestyle's expenses | lifestyle task |
| Hirelings | 102 | Skilled 2 GP per day, untrained 2 SP per day, messenger 2 CP per mile; the table's pay is a minimum | service task |
| Spellcasting services | 102, 159 | Seven rows from 30 GP to 100,000 GP with settlement availability; add the cost of expensive components | service task, effect from the spell's own owner |
| Follow-up: magic-item values | 205–206 | If the GM allows buying and selling, rarity helps set prices; six values, halved for most consumables, plus an incorporated item's cost | value-guide task |

The magic-item value table is not one of the five named scopes. #275 asks that wider scope stay an explicitly named bounded follow-up, and PDF89 points to it. The #80 inventory pins PDF206 to the crafting decision #276 for "Crafting Magic Items"; neither #276's body nor any other issue cites the value table. It is therefore published here as its own small child and not folded into the coin or sale task.

The companion scope file holds every checked clause, the complete tables in copper pieces, the search patterns with their matches, and the body hash of each adjacent issue.

## Existing evidence and unknowns

The runtime has two fixed creation builds that store a `srd521_gold` number in a Coin Pouch item's JSON, and the item contract #195 records a listed cost. There is no coin ledger, purchase, sale, lifestyle or service transaction.

Thirteen open issues grant Gold Pieces once at creation: #433, #508, #554, #635, #730, #785, #853, #982, #1101, #1177, #1480, #1762 and the higher-level money guide #2053. Ten of them list the source key `srd-5.2.1:mechanic/currency/gp`, and #508, #554, #635 and #785 name it among their owned keys, each for its creation grant only. None spends coins. The coin task does not take that key from them. It treats the GP amount they record, and the legacy `srd521_gold` value, as the GP count of the single holding, and adds the other four denominations beside it. No edge to those issues is added and none is edited.

The source leaves several things open. The tasks keep each one unavailable, bind it to a GM-recorded fact, or mark the choice as an Edda adaptation. Nothing is presented as a source rule that the PDF does not print:

- The length of a week or month, and campaign time generally. No issue owns a calendar, so the lifestyle task adds one GM-recorded period with a day count.
- What being unable to pay for a lifestyle means. Adaptation: the choice is rejected and nothing is recorded.
- What "included in your lifestyle's expenses" covers: which row, how many meals, and the price of another grade. Adaptation: nothing is included automatically; the GM marks a specific inn stay or meal as included.
- Rounding when half a listed cost is not a whole copper piece, as with a 1 CP Candle. That sale is unavailable.
- Prices for trade goods, gems and art objects. The Equipment chapter prints no table, so full-value sale needs a GM record of class and value.
- Sale of a crafted item. PDF103 says it sells "at its normal price", against the half-cost rule. That sale stays with #276 and is unavailable in the sale task.
- Whether a sold magic item fetches half, and where one can be bought. PDF205–206 lets the GM set prices by rarity and describes availability with "often", "usually" and "might". The guide and the availability text are shown and never applied or enforced. Potion of Healing has an Equipment price, but PDF99 calls it a magic item, so it takes the magic-item path.
- A Spell Scroll's value. It is double the scribing cost, and the scribing table belongs to #276, so no guide value is shown until that owner exists.
- Exchange between denominations. The source gives relative worth only, so each payment names the exact coins paid and returned.
- Who sells, who buys, whether a willing hireling or caster is present, and the settlement class. No issue owns an offer or a settlement record; the sale and service tasks each add the GM record they need.
- For a hired spell: the cost of a "worth N+ GP" component and a higher-level slot. Adaptations: add N unless the GM names a higher value; a higher slot has no price. The source adds the component cost whether or not the spell consumes it.
- A hired caster's operands. No issue owns a non-player caster. The service task uses #618 unchanged and refuses the service if that resolver does not accept the GM-supplied caster.
- Weight of a bought meal or loaf. A purchase does not satisfy the Malnutrition (PDF185) or Dehydration (PDF181) hazards, which have no owner here.
- What a hireling does once hired. That is a GM fact.

## The six tasks

| Key | Title | Size | Blocked by | Blocks |
|---|---|---|---|---|
| `edda-275:coin-holdings-payment` | Confirm coin holdings and exact listed-price payment | M | #195 | sale, lifestyle, service and qualification tasks |
| `edda-275:equipment-purchase-sale` | Confirm equipment purchase and half-cost or full-value sale | M | coin task, #195, #248, #257, #261 | value-guide and qualification tasks |
| `edda-275:magic-item-value-guide` | Show magic-item values by rarity for GM-priced purchase and sale | S | sale task, #195, #225, #238, #251, #261 | qualification task |
| `edda-275:lifestyle-lodging` | Confirm lifestyle expenses and Food, Drink, and Lodging purchases | M | coin task | qualification task |
| `edda-275:hired-services` | Confirm hireling, messenger and spellcasting service payments | M | coin task, #618 | qualification task |
| `edda-275:economy-qualification` | Qualify coin, lifestyle and service payments for item coverage | M | the five tasks above | #280 |

The coin task owns the only holding and payment and reads the GP amount the creation owners write. Its fixture is a 75 GP creation grant that pays 15 GP and leaves 60 GP. The sale task adds the GM offer, the listed-cost purchase and the sale; its named fixture is the Longsword, bought for 15 GP and sold for 750 CP. Each fixture that needs a catalogue record has that record's issue as a blocker: Longsword #248, Candle #257, Potion of Healing #261 and #238, Armor, +1 #225 and Plate Armor #251. The value-guide task carries the source's own example, Armor, +1 (Plate Armor) at 5,500 GP, and charges only a GM-set price.

Lifestyle and lodging share one task because the inclusion sentence joins them. Hirelings and spellcasting services share one task because both are a priced agreement with a GM offer. Its Revivify fixture costs 600 GP in a town, the 300 GP row plus the 300 GP diamond, and is refused in a village. The spell's effect stays with #618. The other service rows are price previews that commit nothing.

The qualification task composes the five outputs and adds no rule. It runs one holding of 2,556 CP through a purchase, a lifestyle period, a mug of Ale and a hireling, and checks that 2,404 CP are spent and 152 CP remain.

## Native routing and caps

Children attach under #275, as #2051 did under #273. #275 has no sub-issues and will have six. #280 has 96 of 100 blockers. This decision takes one, which makes 97: only the qualification task blocks #280, and the five implementation tasks reach #280 through it. Removing that edge would lose the mandatory path for all five. Three slots then remain for the sibling decisions #274 and #276 together, so each of them must also reach #280 through one aggregate blocker. #279 is at its 100-blocker cap and gets no edge; these are item and equipment rules, and the one spell effect involved is already required by #279 through #618. Every count is read again immediately before publication and recorded in the link-verification receipt.

New prerequisite edges land on #195 (61 dependents, becomes 64), #248 (12, becomes 13), #257 (1, becomes 2), #261 (4, becomes 6), #225 (1, becomes 2), #238 (1, becomes 2), #251 (2, becomes 3) and #618 (9, becomes 10). No edge touches #67, #429 or any creation owner, and none points at an unpublished issue. With the planned edges the dependency graph stays acyclic, and the parent chain #275 → #193 → #54 → #47 is unchanged. No existing body, parent or edge changes apart from the one new #280 blocker.

Mount, vehicle, passenger and crew prices stay with #274. Crafting, brewing and scribing costs, and the sale of a crafted item at its normal price, stay with #276. Residual item keys stay with #278. Rests stay with #427 and #428, and a player caster's own costly materials with #1447.

## Evidence and limits

A focused check compared every listed clause and table row with the pinned PDF text, confirmed the hash, and confirmed that no existing title matches a new one and that #275's only prerequisite, #80, is closed. It did not publish or read back any native issue; that receipt follows publication. Runtime implementation remains paused. This change contains no runtime file, migration, UI behavior, production mutation or Full Rules admission. The implementation tasks keep their own PostgreSQL, web, TUI and device requirements.
