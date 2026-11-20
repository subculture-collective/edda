<!-- edda-270:prone-movement -->
Decision #270; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 14, 179, 181–182, 186–188; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only the Prone movement clauses that no published task binds: voluntarily dropping Prone (PDF14), the exact cost and Speed-0 refusal of righting yourself (PDF186) and the crawling movement cost (PDF179). Consume the actual Prone instance and attack consequences from #488, the Speed-0 projections of Restrained from #727 and Grappled from #778, and their existing movement accounting. Do not copy a condition store, attack modifier, grapple, shove, fall or Unconscious rule. #487 keeps "remain Prone when Unconscious ends"; #779 keeps Shove; #147 keeps falling.

Dropping Prone: on its own turn a creature may give itself the Prone condition without an action, Bonus Action, Reaction or any movement, unless its Speed is 0. The self-applied instance is an ordinary source-bound Prone instance owned by that creature.

Righting: while Prone, the only movement options are crawling or spending movement equal to half the creature's current Speed, rounded down, to right itself, which ends the Prone condition. Use the authoritative current Speed after every qualified increase and reduction, not a printed base value. The cost is deducted from this turn's remaining movement; insufficient remaining movement rejects. If the current Speed is 0, righting is unavailable. A Speed that "is 0 and can't increase" from an actual Grappled or Restrained instance stays 0 for this check whatever bonus is present; the same holds for any other actual Speed-0 instance supplied by its own owner (#487, #1444). Righting ends the creature's Prone condition, meaning every Prone contribution on it, because PDF186 ends the condition itself; it never ends the Speed-0 condition or any other condition.

Crawling: each foot crawled costs 1 extra foot of movement, or 2 extra feet in Difficult Terrain (PDF179, 181). Crawling does not end Prone. This task owns only the crawl multiplier; which spaces are Difficult Terrain stays an authoritative scene fact supplied by its own owner. A different movement mode, special speed or drag cost keeps its own rule and is not combined here unless both facts are authoritative.

## Deterministic fixtures

- Speed 30, Prone, 30 feet of movement left: righting costs 15 and leaves 15. Speed 25 costs 12, not 13. Speed 5 costs 2. After spending 20 of Speed 30, righting (15) rejects with no mutation.
- Prone and Grappled through actual #778, or Prone and Restrained through actual #727: Speed is 0, righting and dropping Prone are both unavailable, and a +10 Speed bonus does not change that. When the Speed-0 instance ends, righting becomes available again at the then-current Speed.
- Crawling 5 feet costs 10 feet of movement; in Difficult Terrain it costs 15. With 10 feet left, a second 5-foot crawl in Difficult Terrain rejects.
- A standing creature with Speed 30 on its own turn drops Prone at no cost: 30 feet remain and its action, Bonus Action and Reaction are untouched. Off-turn dropping, dropping at Speed 0, dropping another creature, or dropping twice rejects.
- Reject righting for free, righting with an action instead of movement, a rounded-up or base-Speed cost, crawling at normal cost, automatic righting at turn start, clearing Restrained or Grappled by righting, and a client-declared Speed, terrain or remaining-movement value. Preserve unrelated conditions and the attack consequences of #488 before and after.

## Durable implementation acceptance

Preview is inert: no dice, writes or hidden facts. Confirm revalidates campaign, actor, turn ownership, current Prone and Speed-0 instances, current Speed, remaining movement, terrain fact, source version and expected revision. Commit the movement spend, the Prone instance creation or ending, derived sheet and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. Same action ID and payload replay after retry, reload or reconnect returns the saved receipt without spending again; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Private NPC Speed, position and conditions never enter previews, receipts, errors or public history. Preserve legacy reads and unrelated JSON; the legacy `prone` combat flag is not silently upgraded to this contract.

Hard prerequisites: #488 actual Prone instance and consequences, #727 Restrained Speed-0 projection, #778 Grappled Speed-0 projection and movement cost accounting. Parent: #270. Directly block the Exhaustion task published with this decision (`edda-270:exhaustion-levels`), which composes its Speed reduction with righting; through it this task is mandatory for #279, #280 and #337. #279 is at its 100-blocker cap: add no direct #279 edge. Add no #67 or #429 consumer and alter no existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
