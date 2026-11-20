# Coins, selling, lifestyle and priced services (#275)

No published issue owns any of the five named scopes. A search of all 1,892 issue titles and bodies, open and closed, found only #275 itself, roadmap lists that name it, unrelated spells and a few boundary mentions. This decision therefore assigns the scopes to three bounded implementation tasks and one qualification task. All four become native children of #275. The qualification task is the single new blocker of item qualification #280. Their exact bodies are the four companion task files; issue numbers, publication and readback are recorded later in the link-verification receipt. Until then nothing here is published, and closing #275 supplies no runtime.

## Source

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Source facts and fixtures retain the CC-BY-4.0 attribution and adaptation notice in [the source register](../srd-5.2.1-conformance.md). The issue names PDF 89 and 101–102. Four cross-references were read as well: Longsword (91) and Candle (95) for sale fixtures, Revivify (159) for a costly-component service, and Magic Item Values by Rarity (205–206) because Selling Equipment defers magic-item prices to it.

| Named scope | PDF | Source rule | Owner |
|---|---|---|---|
| Coins | 89 | CP 1/100, SP 1/10, EP 1/2, GP 1 and PP 10 GP; fifty coins weigh a pound | coin task |
| Selling | 89, 205–206 | Equipment fetches half its cost; trade goods and valuables keep full value; magic-item prices are set by the GM with a rarity guide | coin task |
| Lifestyle | 101 | Seven lifestyles from free to 10 GP per day, chosen and paid at the start of each week or month (GM's choice); no inherent consequences | lifestyle task |
| Food, drink, lodging | 101–102 | Seventeen priced rows; daily lodging and meals are included in the lifestyle's expenses | lifestyle task |
| Hirelings | 102 | Skilled 2 GP per day, untrained 2 SP per day, messenger 2 CP per mile; the table's pay is a minimum | service task |
| Spellcasting services | 102, 159 | Seven rows from 30 GP to 100,000 GP with settlement availability; add the cost of expensive components | service task, effect from the spell's own owner |

The companion scope file holds every checked clause, the complete tables in copper pieces, the search patterns with their matches, and the body hash of each adjacent issue.

## Existing evidence and unknowns

The runtime has two fixed creation builds that store a `srd521_gold` number on a Coin Pouch item, and the item contract #195 records a listed cost. There is no coin ledger, purchase, sale, lifestyle or service transaction. Starting-equipment coin grants belong to their class owners, for example #433 and #730, and the optional higher-level money guide to #2053. Those owners grant coins once; none spends them.

The source leaves several things open. The tasks keep each one unavailable or bind it to an authoritative campaign fact, and guess nothing:

- The length of a week or month. The lifestyle charge needs a campaign calendar day count.
- Rounding when half a listed cost is not a whole copper piece, as with a 1 CP Candle. That sale is unavailable.
- Prices for trade goods, gems and art objects. The Equipment chapter prints no table, so full-value sale needs a GM record of class and value.
- Whether a sold magic item fetches half. PDF205–206 only lets the GM set prices by rarity, so the guide is shown and never applied automatically. Potion of Healing appears in both chapters and takes the magic-item path.
- Exchange between denominations. The source gives relative worth only, so each payment names the exact coins paid and returned.
- Weight of a bought meal or loaf. A purchase does not satisfy the Malnutrition (PDF185) or Dehydration (PDF181) hazards, which have no owner here.
- What a hireling does once hired, and whether a willing caster is present. Both are GM facts.

## The four tasks

| Key | Title | Size | Blocked by | Blocks |
|---|---|---|---|---|
| `edda-275:coin-payment-sale` | Confirm coin holdings, listed-price payment and equipment sale | L | #195, #248 | lifestyle, service and qualification tasks |
| `edda-275:lifestyle-lodging` | Confirm lifestyle expenses and Food, Drink, and Lodging purchases | M | coin task | qualification task |
| `edda-275:hired-services` | Confirm hireling, messenger and spellcasting service payments | M | coin task, #618 | qualification task |
| `edda-275:economy-qualification` | Qualify coin, lifestyle and service payments for item coverage | M | the three tasks above | #280 |

The coin task owns the only ledger and payment. Its named fixture is the Longsword: bought for 15 GP and sold for 750 CP. It is size L because it carries holdings, payment, purchase, sale and the guide-only magic-item values. Lifestyle and lodging share one task because the inclusion rule joins them: under a paid Modest period, a Modest inn stay and meal are covered, and a Wealthy inn stay still costs 2 GP. Hirelings and spellcasting services share one task because both are a priced agreement with an authoritative GM offer. Its Revivify fixture costs 600 GP in a town, the 300 GP row plus the 300 GP diamond, and is refused in a village. The spell's effect stays with #618; the service task charges nothing when no qualified resolver can run.

The qualification task composes the three outputs and adds no rule. It runs one holding of 2,556 CP through a purchase, a lifestyle period, a mug of Ale and a hireling, and checks that 2,404 CP are spent and 152 CP remain.

## Native routing and caps

Children attach under #275, as #2051 did under #273. #275 has no sub-issues and will have four. #280 has 96 of 100 blockers. The sibling decisions #274 and #276 also still owe their children to #280, so three open decisions share four free slots. This decision takes one: only the qualification task blocks #280, and the three implementation tasks reach #280 through it. Removing that edge would lose the mandatory path for all three. #279 is at its 100-blocker cap and gets no edge; these are item and equipment rules, and the one spell effect involved is already required by #279 through #618.

New prerequisite edges land on #195 (61 dependents, becomes 62), #248 (12, becomes 13) and #618 (9, becomes 10). No edge touches #67 or #429. With the planned edges the dependency graph stays acyclic, and the parent chain #275 → #193 → #54 → #47 is unchanged. No existing body, parent or edge changes apart from the one new #280 blocker.

Mount, vehicle, passenger and crew prices stay with #274. Crafting, brewing and scribing costs stay with #276. Residual item keys stay with #278. Rests stay with #427 and #428, and a player caster's own costly materials with #1447.

## Evidence and limits

A focused check compared 54 clauses and every table row with the pinned PDF text, confirmed the hash, and confirmed that no existing title matches a new one and that #275's only prerequisite, #80, is closed. It did not publish or read back any native issue; that receipt follows publication. Runtime implementation remains paused. This change contains no runtime file, migration, UI behavior, production mutation or Full Rules admission. The implementation tasks keep their own PostgreSQL, web, TUI and device requirements.
