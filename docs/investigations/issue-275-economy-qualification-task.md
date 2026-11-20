<!-- edda-275:economy-qualification -->
Decision #275; roadmap #47. One qualification PR; Task / Size M.

Source: English SRD 5.2.1, PDF 89 and 101–102, with the consumed rows on 28, 91, 92, 95–96, 99, 103, 142, 159, 205–206, 210 and 236; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`.

Qualify the finite composition of the five #275 implementation outputs: `edda-275:coin-holdings-payment`, `edda-275:equipment-purchase-sale`, `edda-275:magic-item-value-guide`, `edda-275:lifestyle-lodging` and `edda-275:hired-services`. Own no ledger, price row, payment, purchase, sale, policy, lifestyle, service or spell effect. A record's presence, a merged task or the closed decision #275 is not execution evidence. This task exists so the five outputs reach item qualification #280 through one mandatory edge.

Required deterministic composition controls, run through the actual qualified rules, game, engine and handler callers with isolated explicit entitlements:

- One ledger: a single holding of 2 PP, 3 GP, 4 EP, 5 SP and 6 CP (2,556 CP) pays, in order, a Longsword at 15 GP, a 7-day Modest lifestyle at 7 GP, one mug of Ale at 4 CP and a skilled hireling for 1 day at 2 GP. The named coins and change in each receipt account for exactly 2,404 CP spent and 152 CP left. No task keeps a second balance.
- One store: a character whose creation receipt granted 75 GP pays 15 GP and holds 60 GP in the record its creation owner wrote; a legacy Coin Pouch of 29 GP is counted once.
- Sale after purchase: the bought Longsword sells for 750 CP, half of its listed 1,500 CP, and the item leaves the seller once. A Candle's sale, a crafted Longsword's sale and a Potion of Healing's half-cost sale stay unavailable.
- Inclusion: under the paid Modest period a Modest inn stay the GM marked as included charges nothing; the same stay without a mark charges 5 SP and a Wealthy inn stay 2 GP through the same ledger. Nothing is included without a GM mark.
- Magic items: with the campaign policy off, Armor, +1 (Plate Armor) cannot be bought. With it on, the guide shows 5,500 GP and the purchase charges the GM-set 6,000 GP through the same ledger.
- Service total: Revivify in a town charges 600 GP once, the 300 GP table row plus the 300 GP component, and its effect is the single committed result of #618, used unchanged. A village request and an unpayable request charge nothing and cast nothing. If #618 does not accept the GM-supplied caster, the control is the refusal: nothing charged, nothing cast.
- Source comparison: all five coin values, the fifty-coins-per-pound rule, seven lifestyle prices, seventeen Food, Drink, and Lodging rows, three Hirelings rows, seven Spellcasting Services rows and the six rarity values match the pinned PDF text.
- Fail-closed boundaries hold together: no implicit denomination exchange, no rounding of an odd half-cost, no automatic magic-item price, no enforced magic-item location, no automatic inclusion, no lifestyle consequence, no rest or hazard change from a meal or inn stay, no hireling stat block and no spell effect without its own qualified resolver.

Preview is inert. Confirmation revalidates actor, campaign, source, offer, policy, period, inclusion mark, settlement class and expected revisions. Inject each failure for full isolated PostgreSQL rollback across the composed sequence; a same-ID retry, reload or reconnect returns each saved receipt without paying, granting or casting again; conflicts reject and racing confirmations commit once. Web and trusted TUI agree on supported choices, unavailable reasons and cancel/confirm/reconnect. Changed web requires focused build/lint and desktop/phone fixtures. Do not expose private holdings, offers or hidden state.

Native parent #275; directly block #280. Exact prerequisites: the five #275 implementation tasks named above. Preserve every existing body, parent and prerequisite. No new #279, #67 or #429 edge. Mount, vehicle, cargo and crew prices stay with #274, crafting and scribing with #276 and residual item keys with #278.

Qualification must run focused source and contrary tests, isolated PostgreSQL replay/rollback/concurrency, shared web/TUI fixtures, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`; short tests alone do not qualify database or devices. Runtime remains unqualified until the five outputs and this integration proof pass. No deployment, production grant or Full Rules promotion.
