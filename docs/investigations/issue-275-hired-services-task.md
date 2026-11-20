<!-- edda-275:hired-services -->
Decision #275; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 102 (Hirelings; Spellcasting) and 159 (Revivify); SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement the price, eligibility and payment of the three Hirelings rows and the seven Spellcasting Services rows. Every charge uses the coin payment from the `edda-275:coin-payment-sale` task. This task records an agreement and its payment. It creates no controllable creature, stat block, combat ally, skill check, crew automation or spell effect.

Hirelings (PDF102): a skilled hireling performs a service that involves a proficiency (weapon, tool or skill), such as a mercenary, artisan or scribe, and costs 2 GP per day. An untrained hireling does work that needs no particular proficiency, such as a laborer or porter, and costs 2 SP per day. A messenger costs 2 CP per mile. The table's pay is a minimum; some expert hirelings require more. An authoritative GM offer may therefore set a higher price for a row and never a lower one. Charge whole days and whole miles; the source prices no fraction of either. Whether a willing hireling exists, what the hireling then does, and any loyalty or risk are GM facts. Ship crews of skilled hirelings stay with #274; a hired creature that acts in play needs its own qualified world and creature owners.

Spellcasting Services (PDF102): Cantrip 30 GP and levels 1 and 2 at 50 GP and 200 GP in a village, town or city; level 3 at 300 GP and levels 4–5 at 2,000 GP in a town or city only; levels 6–8 at 20,000 GP and level 9 at 100,000 GP in a city only. If the spell has expensive components, add their cost. Use the cost printed on the spell's actual source record; for a "worth N+ GP" component add N unless the authoritative offer names a higher actual component value. The source prices the spell's own level; a request for a higher-level slot has no price and stays unavailable. The table says where a caster can be found, not that one is present: require an authoritative settlement class (village, town or city) and an authoritative GM offer of a willing caster for that spell. An unknown settlement class is unavailable.

The service receipt binds payer, caster offer, spell identity and level, table row, component costs and total. The effect itself runs only through that spell's own qualified resolver with authoritative caster operands; without them the service is unavailable and nothing is charged. Catalogue presence, this price table and the payment grant no spell execution, learned spell or slot.

## Deterministic fixtures

- Skilled hireling for 3 days: 6 GP. An expert offered at 5 GP per day for 3 days: 15 GP. A skilled hireling offered at 1 GP per day rejects as below the minimum. Untrained hireling for 5 days: 10 SP. Messenger for 25 miles: 50 CP.
- Zero, negative or fractional days or miles, a missing offer, an untrained price for proficient work, and a holding that cannot cover the total reject without partial changes.
- Revivify (PDF159, level 3, a diamond worth 300+ GP, which the spell consumes) in a town: 300 GP + 300 GP = 600 GP. The same request in a village rejects: level 3 is town or city only. With the offer naming a 500 GP diamond the total is 800 GP. Its effect uses the existing Revivify resolver #618 and that owner's own timing and target rules.
- A level-1 spell with no expensive component in a village: 50 GP. A cantrip: 30 GP. A level-5 spell in a town: 2,000 GP plus components. A level-6 spell in a town and a level-9 spell anywhere but a city reject; a level-9 spell in a city costs 100,000 GP plus components.
- Reject an invented level-10 row, a component cost omitted or counted twice, an upcast request, a spell with no qualified resolver, a forged settlement class and a payment taken when the effect cannot resolve.

## Durable implementation acceptance

Preview is inert and reveals no private caster, hireling or GM facts. Confirm revalidates campaign, actor, offer, settlement class, row, days or miles, spell record and component cost, resolver availability, price, source version and expected revisions. Commit the agreement or service receipt, the coin payment, the delegated effect's own committed result where one runs, and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the saved receipt without charging or casting again; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private holdings, offer internals and hidden targets.

Hard prerequisites: the `edda-275:coin-payment-sale` task and #618 actual Revivify resolver. Parent: #275. Required by item qualification #280 through the `edda-275:economy-qualification` task; add no direct #280, #279, #67 or #429 edge and change no existing owner. #1447 keeps a player caster's own costly and consumed materials; scroll scribing and crafting labor stay with #276.

Validation for the implementation PR: focused deterministic rules tests comparing all three hireling rows and all seven service rows with the source; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
