<!-- edda-274:hull-repair -->
Decision #274; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 12, 16, 101, 180; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement two PDF101 rules for the seven large vehicles recorded by {{274:records}}: the Damage Threshold note, using the PDF180 glossary rule, and Ship Repair. A vessel instance stores current Hit Points up to its table maximum. This task resolves no attack against a vessel's AC, no travel, no crew rule and no purchase.

Damage threshold: this task consumes one already-resolved damage instance from a single attack or effect, with its amount and type. A vessel with a damage threshold takes the whole instance when the amount equals or exceeds the threshold and takes nothing when it is lower. Compare each instance separately; never add instances together to reach the threshold. A vessel whose threshold entry is absent takes every instance in full. Hit Points do not go below 0 (PDF16). PDF12 says a vehicle is not an object, so the Poison and Psychic immunity of objects is not applied by inference; an instance of either type against a vessel is unavailable without an authoritative ruling. PDF100–101 states no consequence for a vessel at 0 Hit Points, so record the value and invent no sinking, crew or cargo effect.

Ship repair: repairs can be made only while the vessel is berthed, an authoritative fact. Repairing 1 Hit Point takes 1 day and costs 20 GP for materials and labor. Where supplies and skilled labor are abundant, such as a city shipyard, also an authoritative fact, the time and cost are halved: half a day and 10 GP per Hit Point. Confirm an explicit whole number of Hit Points, no more than the missing amount. Commit only completed units: clock advance, GP cost and restored Hit Points together. The GP cost is paid through the coin output that the #275 economy split assigns; until that output is qualified, repair stays unavailable, and a free or unpaid repair never commits. No ability check, tool, crew or separate hireling wage is added.

## Deterministic fixtures

- Sailing Ship, 300 HP, threshold 15: an instance of 14 leaves 300; an instance of 15 leaves 285; two instances of 10 in one turn leave 300. Keelboat, threshold 10: 9 leaves 100 and 11 leaves 89. Rowboat, 50 HP, no threshold: an instance of 1 leaves 49. An Airship takes an instance of 1 in full.
- Negative, each violated independently: treating an instance equal to the threshold as superficial; subtracting the threshold from the damage; summing instances; a threshold on the Rowboat or Airship; HP below 0; an automatic sinking at 0; Poison or Psychic damage reduced to 0 by assumption.
- Repair: a berthed Sailing Ship at 285 of 300 repairs 15 HP for 15 days and 300 GP. In a shipyard the same repair takes 7½ days and 150 GP. Repairing 1 HP in a shipyard takes half a day and 10 GP. Reject a repair at sea or with an unknown berth, 16 HP when 15 are missing, a fractional Hit Point, a payer with 299 GP for the 300 GP repair, halving only the time or only the cost, and a shipyard halving without the authoritative abundance fact.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, vessel identity and record, current Hit Points, the damage instance or the berth, shipyard, payer and balance, clock, source version and expected revisions. Commit Hit Points, cost, clock and public history atomically. A damage application belongs to its triggering event, with no invented optional activation. Inject failures after each stage and prove full isolated PostgreSQL rollback. Same action ID and payload replay after retry, reload or reconnect returns the same receipt and never repeats damage, payment or repair; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, unavailable reasons and confirmation/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private balances and hidden attackers.

Hard prerequisites: {{274:records}} vessel records, #275 for the coin-payment output its split assigns. Parent: #274. The cohort qualification {{274:cohort}} requires this task and is the direct blocker of #280; add no #280 edge here. Other damage-threshold users, such as the PDF201 rolling stone, keep their own owners.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
