<!-- edda-270:prone-movement -->
Decision #270; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 14, 16, 178–179, 181–182, 185–188; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only the Prone movement clauses that #488 does not state. #488 keeps the Prone instance, its attack consequences and the two movement options as its body names them ("crawl or spend half Speed to stand"); this task consumes that output and never keeps a second Prone store or a second righting action. It adds four things: voluntarily dropping Prone (PDF14), the round-down of the righting cost from the current Speed (PDF186), the refusal to right yourself or drop Prone at Speed 0 against actual Speed-0 instances (PDF14, 186), and the crawling cost numbers (PDF179). #487 keeps "remain Prone when Unconscious ends"; #779 keeps Shove; #147 keeps falling.

Dropping Prone (PDF14): on its own turn a creature can give itself the Prone condition "without using an action or any of your Speed, but you can’t do so if your Speed is 0." The self-applied instance is an ordinary source-bound Prone instance owned by that creature. The source names only the action and Speed; the fixture assertion that the Bonus Action and Reaction are also untouched is an Edda fixture control.

Righting (PDF186): the cost is half the creature's Speed, rounded down. Use the authoritative current Speed that #780 reads, after every qualified increase and reduction, not a printed base value and never a Fly, Climb, Swim or Burrow Speed. Insufficient movement left this turn rejects. If the current Speed is 0, righting is unavailable. A Speed that "is 0 and can't increase" from an actual Restrained (#727) or Grappled (#778) instance stays 0 for this check whatever bonus is present; the same holds for any other actual Speed-0 instance supplied by its own owner (#487, #1444). Righting ends the creature's Prone condition, meaning every Prone contribution on it, because PDF186 ends the condition itself; it never ends the Speed-0 condition or any other condition.

Crawling (PDF179): each foot crawled costs 1 extra foot of movement, or 2 extra feet in Difficult Terrain. Crawling does not end Prone. This task owns only those two figures. The general Difficult Terrain cost and the facts that make a space Difficult Terrain (PDF14, 181) belong to the open hazard split #150; until its child exists the terrain fact is an isolated trusted fact and production crawling through Difficult Terrain stays unavailable. A drag cost keeps its own #778 rule and is not combined here unless both facts are authoritative.

Movement spent: no published task owns a general per-turn movement record. This task adds the minimal one its own spends need: the turn's allowance from the authoritative Speed of #780, feet spent by righting and crawling, and replay-safe receipts. It adds no path, position, Opportunity Attack or other movement mode.

Not in this task, and unavailable until an owner exists:

- A flying creature that becomes Prone falls unless it can hover (PDF182). The flying fall task published with this decision (`edda-270:flying-fall`) owns that and consumes this task's drop.
- Burning (PDF178) lets a creature give itself Prone as an action to put out fire. That is an action-cost self-Prone from a hazard, not the free drop; it is refused here.
- Prone from a Long Jump landing (PDF185), from ending a turn in another creature's space (PDF14) and from falling off a mount (PDF16) are sources with their own owners.
- How movement spent on righting counts against a special speed used in the same turn (PDF188 says only to subtract "the distance already moved" when switching speeds). Righting in a turn in which the creature also uses a special speed stays unavailable.

## Deterministic fixtures

- Speed 30, Prone, 30 feet of movement left: righting costs 15 and leaves 15. Speed 25 costs 12, not 13. Speed 5 costs 2. After spending 20 of Speed 30, righting (15) rejects with no mutation. Speed 30 with Fly Speed 60 and no movement yet this turn: righting costs 15, not 30.
- A standing creature Grappled through actual #778, or Restrained through actual #727, has Speed 0 and cannot drop Prone, even with a +10 Speed bonus. A Prone creature with either instance cannot right itself, with or without the bonus. When the Speed-0 instance ends, righting becomes available again at the then-current Speed.
- Crawling 5 feet costs 10 feet of movement; with a trusted Difficult Terrain fact it costs 15. With 10 feet left, a second 5-foot crawl in Difficult Terrain rejects.
- A standing creature with Speed 30 on its own turn drops Prone at no cost: 30 feet remain and its action is unspent (Edda control: Bonus Action and Reaction also unspent). Off-turn dropping, dropping another creature, dropping while already Prone, or a Burning-style drop that claims the free path rejects.
- Reject righting for free, righting with an action instead of movement, a rounded-up, base-Speed or special-speed cost, crawling at normal cost, automatic righting at turn start, clearing Restrained or Grappled by righting, and a client-declared Speed, terrain or movement-left value. Preserve unrelated conditions and the attack consequences of #488 before and after.

## Durable implementation acceptance

Preview is inert: no dice, writes or hidden facts. Confirm revalidates campaign, actor, turn ownership, current Prone and Speed-0 instances, current Speed, movement spent this turn, terrain fact, source version and expected revision. Commit the movement spend, the Prone instance creation or ending, derived sheet and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. Same action ID and payload replay after retry, reload or reconnect returns the saved receipt without spending again; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Private NPC Speed, position and conditions never enter previews, receipts, errors or public history. Preserve legacy reads and unrelated JSON; the legacy `prone` combat flag is not silently upgraded to this contract.

Hard prerequisites: #488 actual Prone instance, attack consequences and movement options; #727 Restrained Speed-0 projection; #778 Grappled Speed-0 projection; #780 authoritative modified Speed. Parent: #270. Directly block the flying fall task (`edda-270:flying-fall`), which consumes the drop, and the Exhaustion penalty task (`edda-270:exhaustion-penalties`). #279 is at its 100-blocker cap, so this task takes no direct #279 edge: the edge to the penalty task is its route to #279, #280 and #337 through the existing blocker #623, and that task's composition fixture (righting at a reduced Speed) is the reason the edge is valid. Add no #67 or #429 consumer and alter no existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
