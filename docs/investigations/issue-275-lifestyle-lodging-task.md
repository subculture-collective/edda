<!-- edda-275:lifestyle-lodging -->
Decision #275; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 101–102 (Lifestyle Expenses; Food, Drink, and Lodging); SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement the periodic lifestyle choice and the seventeen Food, Drink, and Lodging purchases. Every charge uses the coin payment from the `edda-275:coin-payment-sale` task; implement no second ledger or payment.

Lifestyle (PDF101): Wretched is free; Squalid costs 1 SP per day, Poor 2 SP, Modest 1 GP, Comfortable 2 GP, Wealthy 4 GP and Aristocratic 10 GP. At the start of each week or month, whichever the GM chose for the campaign, the character chooses one lifestyle and pays its price for that period. The source defines neither the length of a week nor of a month: the period's day count must be an authoritative campaign calendar fact, and the charge is the daily price times that count. The expense covers lodging, food, equipment maintenance and other necessities. Lifestyles have no inherent consequences: add no modifier, condition, rest change, risk roll, reputation score or equipment decay. The GM may take the recorded lifestyle into account; that stays a GM judgment. The source states no mid-period change, refund, arrears or automatic downgrade, so none exists here. A character who cannot pay the chosen price must choose a lifestyle they can pay for.

Food, Drink, and Lodging (PDF101–102), each a single listed-price purchase: Ale (mug) 4 CP; Bread (loaf) 2 CP; Cheese (wedge) 1 SP; Inn Stay per Day at Squalid 7 CP, Poor 1 SP, Modest 5 SP, Comfortable 8 SP, Wealthy 2 GP, Aristocratic 4 GP; Meal at Squalid 1 CP, Poor 2 CP, Modest 1 SP, Comfortable 2 SP, Wealthy 3 SP, Aristocratic 6 SP; Wine (bottle) Common 2 SP, Fine 10 GP. There is no Wretched row. An inn stay is one night's lodging and grants no rest, healing or safety by itself: Short and Long Rests stay with #427 and #428. The table gives no weight for a meal or loaf, so a purchase does not by itself satisfy or change the PDF185 Malnutrition or PDF181 Dehydration hazards. Rations and other gear keep their own Equipment rows.

Inclusion (PDF101): prices for daily lodging and meals are included in the lifestyle's expenses. During a paid period, record that lifestyle's own Inn Stay per Day and Meal rows as covered at no further charge. Any other row, a different grade, or any row when no period is paid, is bought at its table price. Wretched covers nothing because the table has no Wretched row.

## Deterministic fixtures

- With a campaign week of 7 days, Modest costs 7 GP and Squalid 7 SP. With a campaign month of 30 days, Poor costs 60 SP (6 GP) and Aristocratic 300 GP. Wretched records the choice and charges nothing.
- A second choice inside the same paid period, a choice with no GM week-or-month setting, an unknown period length, a refund request, a lifestyle name outside the seven, or a holding that cannot cover the price rejects without partial changes. No debt or downgrade is created.
- Under a paid Modest period, a Modest inn stay and a Modest meal are covered and charge nothing; a Wealthy inn stay costs 2 GP; a mug of Ale costs 4 CP. With no paid period the Modest inn stay costs 5 SP and the Modest meal 1 SP.
- Fine Wine costs 10 GP and Common Wine 2 SP. Reject a Wretched inn stay or meal, an invented row, a negative or fractional quantity and a covered row charged twice or granted twice on replay.
- A recorded Aristocratic or Wretched lifestyle leaves every d20 Test, condition, rest and Hit Point unchanged.

## Durable implementation acceptance

Preview is inert. Confirm revalidates campaign, actor, GM period setting, calendar day count, current period, lifestyle or row identity, quantity, price, source version and expected revisions. Commit the choice or purchase, the coin payment, the coverage record and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without charging again; a conflicting payload rejects; racing confirmations commit once. A period boundary creates no automatic charge: the next choice is an explicit confirmed action. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private holdings and other characters' lifestyle payments unless the GM discloses them.

Hard prerequisite: the `edda-275:coin-payment-sale` task. Parent: #275. Required by item qualification #280 through the `edda-275:economy-qualification` task; add no direct #280, #279, #67 or #429 edge and change no existing owner. Ship passenger fares on PDF101 stay with #274.

Validation for the implementation PR: focused deterministic rules tests comparing all seven lifestyle prices and all seventeen rows with the source; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
