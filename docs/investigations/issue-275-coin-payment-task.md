<!-- edda-275:coin-holdings-payment -->
Decision #275; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 89 (Coins) and 28 (the Barbarian 75 GP starting option used as a fixture); SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement one coin holding per character and one exact coin payment. The purchase, sale, magic-item, lifestyle, lodging and service tasks of #275 consume this holding and payment; none copies them. This task buys and sells nothing: no item, shop inventory, haggling, banking, debt, treasure generation, carrying-capacity enforcement or item effect is added.

Coins follow PDF89 exactly: Copper Piece (CP) 1/100 GP, Silver Piece (SP) 1/10 GP, Electrum Piece (EP) 1/2 GP, Gold Piece (GP) 1 GP, Platinum Piece (PP) 10 GP. Store every amount as an integer number of CP: CP1, SP10, EP50, GP100, PP1000. A holding is a non-negative integer count of each of the five denominations. Fifty coins of any denominations weigh one pound; expose that exact fraction as a derived fact and round nothing. The source's "about a third of an ounce" is descriptive, not a second weight rule.

One store. Creation owners already grant Gold Pieces once and spend none: #433, #508, #554, #635, #730, #785, #853, #982, #1101, #1177, #1480, #1762 and #2053. They record the amount under source key `srd-5.2.1:mechanic/currency/gp` or as "inventory/currency" or "inventory/gold", and the two fixed creation builds store it as the integer `srd521_gold` in a Coin Pouch item's JSON. That recorded GP amount is the GP count of this holding. The ledger reads it where the merged creation owners persist it, or from `srd521_gold` while none has merged, and writes every GP change back to that same record. The CP, SP, EP and PP counts are new integer fields beside it on that record. Keep no second GP balance, convert or duplicate no existing value, change no creation owner, and preserve legacy reads and unrelated JSON.

Payment. The price is an exact CP amount supplied by an authoritative caller. In this PR the only caller is one GM-recorded charge that the PR adds: campaign, payer, exact CP price and a label, opened and withdrawn by explicit GM action. The later #275 tasks supply source-row prices through the same payment. The confirmed payload names the exact coins the payer hands over and any coins returned as change; coins paid minus coins returned must equal the price. A credit names the exact coins received. The source states relative worth, not an exchange service, so no denomination is converted implicitly. Reject a payment the named coins cannot cover.

## Deterministic fixtures

- 3 EP + 5 SP is worth 200 CP, the same as 2 GP; 1 PP is worth 1,000 CP. 100 coins weigh 2 pounds; 49 coins weigh 49/50 pound.
- With a GM charge of 15 GP, paying 15 GP leaves a 20 GP holding at 5 GP. Paying 2 PP with 5 GP returned as change is also exact. Paying 14 GP, 2 PP with no change, coins the payer does not hold, or a forged or withdrawn charge rejects without partial changes.
- A character whose creation receipt granted 75 GP (the Barbarian option B on PDF28, seeded as its owner #433 records it, under an isolated test entitlement) pays 15 GP and holds 60 GP in that same record.
- A fixed-build character with the legacy Coin Pouch `{"srd521_gold":29}` holds 29 GP once. Paying 15 GP leaves 14 GP in that record. The holding is never 58 GP, no second balance appears and unrelated JSON keys are unchanged.
- A credit of 7 GP and 5 SP raises the holding by exactly those coins.
- Reject negative, fractional, boolean or overflowing counts, an unknown denomination, a foreign payer, a missing or stale charge and an implicit exchange of one denomination for another.

## Durable implementation acceptance

Preview is inert: no dice, writes or disclosure of another character's holdings. Confirm revalidates campaign, actor, payer, charge, price, held coins, source version and expected revisions. Commit the coin deltas, the saved operands and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without paying again; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private holdings and GM charge internals.

Hard prerequisite: #195 item contract, which the Coin Pouch record follows. Parent: #275. This task is required by item qualification #280 through the other #275 tasks and the `edda-275:economy-qualification` task; it adds no direct #280, #279, #67 or #429 edge, no edge to a creation owner, and changes no existing owner. Starting-equipment coin grants stay with their class owners named above and the optional higher-level money guide with #2053. Purchase and sale belong to `edda-275:equipment-purchase-sale`.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
