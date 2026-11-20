<!-- edda-275:coin-payment-sale -->
Decision #275; roadmap #47. One implementation PR; Task / Size L.

Source: English SRD 5.2.1, PDF 89 (Coins, Selling Equipment), 91 (Longsword), 95 (Candle) and 205–206 (Magic Item Values by Rarity); SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement one coin ledger, one listed-price payment and one equipment sale. Later lifestyle, lodging and service tasks consume this payment; they do not copy it. No shop inventory, haggling, banking, debt, treasure generation, carrying-capacity enforcement or item effect is added.

Coins follow PDF89 exactly: Copper Piece (CP) 1/100 GP, Silver Piece (SP) 1/10 GP, Electrum Piece (EP) 1/2 GP, Gold Piece (GP) 1 GP, Platinum Piece (PP) 10 GP. Store every amount as an integer number of CP: CP1, SP10, EP50, GP100, PP1000. An owner's holding is a non-negative integer count of each of the five denominations. Fifty coins of any denominations weigh one pound; expose that exact fraction as a derived fact and round nothing. The source's "about a third of an ounce" is descriptive, not a second weight rule. Preserve legacy reads of the existing fixed-creation `srd521_gold` Coin Pouch values and unrelated JSON; do not silently convert or duplicate them.

Payment: the price is an exact CP amount from an authoritative source row or an authoritative GM-set offer. The confirmed payload names the exact coins the payer hands over and any coins returned as change; coins paid minus coins returned must equal the price. The source states relative worth, not an exchange service, so no denomination is converted implicitly. Reject a payment the named coins cannot cover.

Purchase: one Equipment identity with a listed cost, offered by an authoritative GM or merchant offer fact for the current scene. The listed cost applies unless that offer records a different GM-set price. Ownership of the bought item is granted once; identity alone executes no property, mastery or effect. Named fixture: Longsword, PDF91, 15 GP, catalogue record from #248.

Sale (PDF89): Equipment fetches half its listed cost. Trade goods and valuables such as gems and art objects keep their full value; the Equipment chapter prints no price table for them, so full-value sale needs an authoritative GM record that classifies the object and states its value. Never infer that class from a name. If half the listed cost is not a whole number of CP, the source gives no rounding rule: that sale stays unavailable.

Magic items: PDF89 defers their prices to "Magic Items". PDF205–206 lets a GM who allows buying and selling use rarity to set prices; it does not say a sale fetches half. Buying or selling a magic item therefore needs an explicit campaign policy and an authoritative GM-set price. Under that policy show the rarity guide as guidance only and never apply it or the half-cost rule automatically: Common 100 GP, Uncommon 400 GP, Rare 4,000 GP, Very Rare 40,000 GP, Legendary 200,000 GP, Artifact priceless; halve the value for a consumable other than a Spell Scroll; a Spell Scroll's value is double its scribing cost; add the Equipment cost of an incorporated item. An identity printed in both chapters, such as Potion of Healing, takes the magic-item path. A seller asking for a service instead of coin is a GM fact, not a transaction here.

## Deterministic fixtures

- 3 EP + 5 SP is worth 200 CP, the same as 2 GP; 1 PP is worth 1,000 CP. 100 coins weigh 2 pounds; 49 coins weigh 49/50 pound.
- With an offer for a Longsword at its listed cost, paying 15 GP leaves a 20 GP holding at 5 GP. Paying 2 PP with 5 GP returned as change is also exact. Paying 14 GP, 2 PP with no change, coins the payer does not hold, or a forged price rejects without partial changes.
- Selling the Longsword fetches 750 CP, received as the named 7 GP and 5 SP. A Candle's listed 1 CP has no whole-copper half: unavailable, not rounded to 0 or 1.
- A GM-recorded gem valued at 50 GP sells for 50 GP. The same object without that record, or an ordinary Longsword claimed as a valuable, does not get full value.
- Guide witnesses under the enabled policy: +1 Armor (Plate Armor) shows 5,500 GP = 4,000 + 1,500; a Common consumable shows 50 GP; an Artifact shows no price. Without the policy or a GM-set price, a magic-item sale or purchase is unavailable; a Potion of Healing never sells for an automatic 25 GP.
- Reject negative, fractional, boolean or overflowing counts, an unknown denomination, a foreign owner or item, a missing or stale offer, a duplicate grant and a sale of an item the seller no longer owns.

## Durable implementation acceptance

Preview is inert: no dice, writes or disclosure of another owner's holdings. Confirm revalidates campaign, actor, ownership, offer, policy, price, source version and expected revisions. Commit coin deltas for both sides that are tracked, the item transfer, the saved operands and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without paying or granting again; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private holdings, GM offer internals and hidden prices.

Hard prerequisites: #195 item contract with listed cost, #248 Longsword catalogue record. Parent: #275. This task is required by item qualification #280 through the lifestyle and service tasks and the `edda-275:economy-qualification` task; it adds no direct #280, #279, #67 or #429 edge and changes no existing owner. Starting-equipment coin grants stay with their class and background owners, the optional higher-level money guide with #2053, crafting costs with #276 and ship or mount prices with #274.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
