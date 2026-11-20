<!-- edda-275:magic-item-value-guide -->
Decision #275; roadmap #47. One implementation PR; Task / Size S.

Source: English SRD 5.2.1, PDF 89 (Selling Equipment defers magic-item prices), 205–206 (Magic Item Values by Rarity), 92 (Plate Armor), 210 (Armor, +1, +2, or +3), 95 and 99 (Potion of Healing) and 236 (Potions of Healing); SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

#275 names coins and selling; PDF89 sends magic-item prices to "Magic Items", and no issue owns that page's value table. This task is the bounded follow-up for it. It shows the rarity values to the GM and lets a magic item change hands only at a price the GM sets. It reuses the offer, purchase and sale of `edda-275:equipment-purchase-sale` and the payment of `edda-275:coin-holdings-payment`; it adds no ledger, no item effect, no attunement and no crafting.

Policy. PDF205: "If you allow characters to buy and sell magic items in your campaign, rarity can help you set prices for those items." This PR adds one GM-recorded campaign policy, off by default, set and cleared by explicit GM action. With the policy off, buying and selling a magic item stay unavailable.

Price. With the policy on, a magic-item offer must carry a GM-set price in CP; this PR adds that one field to the offer for magic-item identities only. There is no default price. The source does not say a sold magic item fetches half, so the half-cost rule is never applied to a magic item.

Guide (PDF206), shown to the GM beside the offer and never applied automatically: Common 100 GP, Uncommon 400 GP, Rare 4,000 GP, Very Rare 40,000 GP, Legendary 200,000 GP, Artifact priceless. Halve the value for a consumable other than a Spell Scroll. If the item incorporates an item with an Equipment cost, add that cost. The source says a Spell Scroll is worth double its scribing cost; the scribing table belongs to the crafting decision #276, so a Spell Scroll shows no guide value until that owner exists.

Availability (PDF205), shown as guide text only: "Common magic items can often be bought in a town or city. Uncommon and Rare magic items are usually found only in cities, and rarer magic items might be sold only in wondrous locations, such as a city on another plane of existence." "Often", "usually" and "might" leave this to the GM, so no location check is enforced. A seller who asks for a service instead of coin is a GM fact, not a transaction here.

## Deterministic fixtures

- Armor, +1 (Plate Armor) is Rare (PDF210) and Plate Armor costs 1,500 GP (PDF92): the guide shows 5,500 GP = 4,000 + 1,500, the source's own example. With the policy on and a GM-set price of 6,000 GP, the purchase charges 6,000 GP, not 5,500.
- Potion of Healing is Common (PDF236) and consumable: the guide shows 50 GP, half of 100. With the policy on and a GM-set price of 40 GP it sells for 40 GP. It never sells for an automatic 25 GP.
- An Artifact shows "priceless" and no number. A Spell Scroll shows no guide value.
- With the policy off, or with the policy on and no GM-set price, a magic-item purchase or sale is unavailable and charges nothing.
- Reject a guide value used as a price without a GM-set price, a half-cost sale of a magic item, a GM-set price on a nonmagical identity, a forged policy and a location refusal invented from the availability text.

## Durable implementation acceptance

Preview is inert. Confirm revalidates campaign, actor, policy, offer, GM-set price, item record and rarity, source version and expected revisions. The policy change, and each purchase or sale with its coin deltas, item transfer and public history, commits atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without paying or granting again; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior; the guide and availability text appear only to the GM. Changed web requires focused build/lint and desktop/phone checks. Public history omits private holdings, the policy's internals and undisclosed prices.

Hard prerequisites: the `edda-275:equipment-purchase-sale` task; #195 item contract with rarity and listed cost; catalogue records #225 (Armor, +1, +2, or +3), #251 (Plate Armor), #238 (Potions of Healing) and #261 (Potion of Healing). Parent: #275. Required by item qualification #280 through the `edda-275:economy-qualification` task; add no direct #280, #279, #67 or #429 edge and change no existing owner. Magic-item crafting and scroll scribing stay with #276.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
