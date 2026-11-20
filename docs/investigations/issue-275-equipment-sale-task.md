<!-- edda-275:equipment-purchase-sale -->
Decision #275; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 89 (Selling Equipment), 91 (Longsword), 95–96 (Candle), 95 and 99 (Potion of Healing, as a boundary) and 103 (sale of a crafted item, as a boundary); SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement one nonmagical Equipment purchase at its listed cost and one Equipment sale. Every coin movement uses the holding and payment from the `edda-275:coin-holdings-payment` task; implement no second ledger or payment. No shop inventory, stock count, haggling or item effect is added.

Offer. The source prices items and does not say who is selling or buying. This PR adds one GM-recorded offer: campaign, scene, item identity, direction (the other party sells, or the other party buys), opened and withdrawn by explicit GM action. A purchase or sale needs a current offer for that identity and direction. No other issue owns such a record.

Purchase. One Equipment identity with a listed cost is bought for exactly that cost. Ownership of the bought item is granted once and its acquisition record names the purchase receipt. Identity alone executes no property, mastery or effect. Named fixture: Longsword, PDF91, 15 GP, catalogue record from #248.

Sale (PDF89). "Equipment fetches half its cost when sold." The seller receives half the listed cost as named coins. If half the listed cost is not a whole number of CP, the source gives no rounding rule and that sale stays unavailable. Trade goods and valuables such as gems and art objects keep their full value; the Equipment chapter prints no price table for them, so this PR adds one GM-recorded valuable: campaign, owner, object, class (trade good, gem or art object) and value in CP. Full-value sale needs that record. Never infer the class from a name.

Crafted items. PDF103 says a crafted item can be sold "at its normal price", which conflicts with the half-cost rule for that item. Crafting and the sale of a crafted item stay with the crafting decision #276. Here the half-cost sale applies only to an item whose acquisition record names no crafting receipt; an item whose record names one is unavailable for sale in this task.

Magic items. PDF89 defers their prices to "Magic Items". Buying or selling a magic item is unavailable in this task and belongs to `edda-275:magic-item-value-guide`. Potion of Healing has an Equipment row at 50 GP (PDF95) and its entry on PDF99 says "This potion is a magic item", so it takes the magic-item path and is unavailable here too.

## Deterministic fixtures

- With a selling offer for a Longsword, paying 15 GP leaves a 20 GP holding at 5 GP and grants one Longsword. Paying 2 PP with 5 GP returned as change is also exact. Paying 14 GP, or any payment with no current offer, rejects without partial changes.
- With a buying offer, selling the Longsword credits 750 CP, received as the named 7 GP and 5 SP, and the item leaves the seller once.
- A Candle's listed 1 CP (PDF95) has no whole-copper half: its sale is unavailable, not rounded to 0 or 1. Its purchase for 1 CP succeeds.
- A GM-recorded gem valued at 50 GP sells for 50 GP. The same object without that record, or an ordinary Longsword claimed as a valuable, does not get full value.
- A Longsword whose acquisition record names a crafting receipt is unavailable for sale: it fetches neither 750 CP nor 1,500 CP here.
- A Potion of Healing is unavailable for purchase and sale here; it never sells for an automatic 25 GP.
- Reject a foreign owner or item, a missing, stale or withdrawn offer, an offer in the wrong direction, a forged price, a duplicate grant and a sale of an item the seller no longer owns.

## Durable implementation acceptance

Preview is inert: no dice, writes or disclosure of another owner's holdings. Confirm revalidates campaign, actor, ownership, offer, acquisition record, valuable record, listed cost, source version and expected revisions. Commit the coin deltas, the item transfer, the saved operands and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without paying or granting again; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private holdings, GM offer internals and recorded values the GM has not disclosed.

Hard prerequisites: the `edda-275:coin-holdings-payment` task; #195 item contract with listed cost; catalogue records #248 (Longsword), #257 (Candle) and #261 (Potion of Healing). Parent: #275. Required by item qualification #280 through the `edda-275:economy-qualification` task; add no direct #280, #279, #67 or #429 edge and change no existing owner. Crafting costs and crafted-item sale stay with #276, ship and mount prices with #274.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
