<!-- edda-270:flying-fall -->
Decision #270; roadmap #47. One implementation PR; Task / Size S.

Source: English SRD 5.2.1, PDF 182–184, 186–188; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only the Flying rule of PDF182: "While flying, you fall if you have the Incapacitated or Prone condition or your Fly Speed is reduced to 0. You can stay aloft in those circumstances if you can hover." No published task owns it, and #147 excludes flying. This task decides only whether a fall starts. The fall itself, its damage and the landing Prone stay with #147, which receives one declared fall with the trusted height.

Triggers, each read from its actual owner and never from a client flag:

- Prone: an actual instance from #488, including a voluntary drop through the Prone movement task (`edda-270:prone-movement`).
- Incapacitated: an actual instance from #487, including the Incapacitated part of Unconscious, Paralyzed, Stunned or Petrified when their owners supply it.
- Fly Speed reduced to 0: PDF188 says a change to Speed changes every special speed equally and gives the example that Speed reduced to 0 makes a Climb Speed 0. So an actual Grappled (#778) or Restrained (#727) Speed-0 instance makes the Fly Speed 0, and so does any authoritative reduction that reaches 0.

Hover (PDF183) is an authoritative creature or effect fact from its own owner, such as a stat block or #1335. With it the creature stays aloft in all three circumstances. Without a known hover fact the creature does not hover.

No published task owns an authoritative "currently flying at this height" state. This task reads it only as an isolated trusted fact and enables no production flight; a fall trigger for a creature whose flying state or height is unknown stays unavailable. Ending of the Fly spell stays #1336. Falling onto another creature and magical fall prevention keep #147's exclusions.

## Deterministic fixtures

- Trusted flying at 30 feet, Fly Speed 60, no hover: gaining Prone from actual #488, or dropping Prone through the Prone movement task, starts one fall of 30 feet handed to #147. The same creature with a trusted hover fact stays aloft and keeps Prone.
- Flying, no hover, becomes Incapacitated through actual #487: one fall. With hover: none.
- Flying, no hover, Grappled through actual #778 or Restrained through actual #727: Fly Speed is 0 and it falls once; a +10 Speed bonus changes nothing. With hover: stays aloft with Fly Speed 0.
- An isolated trusted Speed reduction of 20 on Fly Speed 20 gives 0 and a fall; a reduction of 15 gives Fly Speed 5 and no fall.
- A creature on the ground, a creature with a Fly Speed that is not flying, or one with unknown height: no fall and an unavailable reason. Two triggers in the same transaction start one fall, not two. Replay starts none.
- Reject a client-declared flying, hover, height or Fly Speed value, a fall for a hovering creature, a fall that skips #147's damage and landing rules, and ending the triggering condition because the creature landed.

## Durable implementation acceptance

The trigger is evaluated inside the transaction that commits the triggering condition or Speed change; it adds no optional activation of its own. Preview of the triggering action shows the pending fall only from facts the player may know. Commit the condition or Speed change, the fall start handed to #147, derived sheet and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. Same action ID and payload replay after retry, reload or reconnect returns the saved receipt with no second fall; a conflicting payload rejects; racing triggers commit one fall. Web and trusted TUI show the same fall notice and unavailable reasons and share confirm/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. A hidden NPC's height, hover fact, Speed and conditions never enter previews, receipts, errors or public history. Preserve legacy reads and unrelated JSON.

Hard prerequisites: #147 declared fall; #487 actual Incapacitated instance; #727 Restrained and #778 Grappled Speed-0 projections; the Prone movement task (`edda-270:prone-movement`), which also carries the #488 Prone instance. Parent: #270. Directly block the Exhaustion penalty task (`edda-270:exhaustion-penalties`). #279 is at its 100-blocker cap, so this task takes no direct #279 edge: the edge to the penalty task is its route to #279, #280 and #337 through the existing blocker #623, and that task's composition fixture (a Fly Speed reduced to 0 by Exhaustion) is the reason the edge is valid. Add no #67 or #429 consumer and alter no existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
