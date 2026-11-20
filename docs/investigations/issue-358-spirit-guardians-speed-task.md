<!-- edda-358:spirit-guardians-speed -->
Decision #358; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 5, 106, 164, 181, 188; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only one clause of Spirit Guardians (PDF164): the Speed of every creature that the caster did not designate is halved while that creature is in the Emanation. Reuse the committed cast, designation set and membership revision of {{ISSUE_SG_SAVES}}; do not copy its cast, triggers, saves or damage. No other Speed modifier, Difficult Terrain rule, condition, spell or movement mode is added.

Projection: while a non-designated creature is inside the active Emanation, its Speed is half its otherwise authoritative Speed, rounded down (PDF5). Every special speed it has is halved for the same time (PDF188). The halving is one source-bound modifier tied to the cast ID. It ends the moment the creature is outside, the caster's Concentration ends or the duration expires, and it never changes the stored base Speed. The effects of the same spell cast more than once do not combine (PDF106), so the halving never applies twice; {{ISSUE_SG_SAVES}} admits only one active cast per encounter. The caster and every designated creature keep their Speed.

Movement budget: the source states no accounting for a Speed that changes during a move. As a labelled Edda adaptation of the PDF188 switching rule, each five-foot step is checked against the Speed in force in the space the creature is leaving, minus the distance it has already moved this turn. A result of zero or less ends its movement for that turn. Report this adaptation in the receipt's source derivation. Do not present it as source text.

Dash and other consumers read the projected Speed from their own owners; this task changes none of them. Unknown base Speed, membership or designation facts keep the movement unavailable with an explicit reason.

## Deterministic fixtures

- A non-designated creature with Speed 30 that starts its turn inside has Speed 15. With Speed 25 it has 12. With Speed 30 and a Fly Speed 40 it has 15 and 20.
- A designated creature and the caster keep Speed 30 inside the area.
- Starting 25 feet from the caster with Speed 30, a creature moves two steps to 15 feet. The second step is paid at Speed 30 and triggers the entry save owned by {{ISSUE_SG_SAVES}}. Inside, 15 − 10 leaves 5 feet. A creature that has already moved 20 feet when it enters has 15 − 20, so it cannot move again that turn.
- A creature that starts inside with Speed 30 and steps from 15 feet to 20 feet pays that step at Speed 15. Outside, 30 − 5 leaves 25 feet.
- When the caster's Concentration ends during the creature's turn, the next step uses Speed 30 minus the distance already moved.
- Reject a client-supplied Speed, a halving for a designated creature, a halving after the effect ended, a second halving and any change to the stored base Speed. Reload mid-turn restores the same projected Speed and distance moved.

## Durable implementation acceptance

Preview is inert and shows the projected Speed and remaining movement without writing. Confirm of each movement revalidates campaign, actor, cast ID, designation, membership revision, base Speed, distance moved, clock and turn cursor, then commits the step, the membership change and public history in one transaction. Inject a history-write failure and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the committed receipt; a conflicting payload rejects; racing confirmations commit once. Public receipts omit hidden creatures and private NPC facts. Web and trusted TUI share availability, disabled reasons and confirmation, cancel and reconnect behavior; changed web requires focused build, lint and desktop and phone checks.

Hard prerequisite: {{ISSUE_SG_SAVES}}, which supplies the cast, designation and membership and already requires #220 and #353. Consumed closed baselines without new native edges: #58 sheet projection, #64 confirmation, #66 movement budget and turn identity, #67 trusted distance. Parent: #193. Directly block item qualification #280; this terminal carries {{ISSUE_SG_SAVES}} to that gate. Spell qualification #279 already holds 100 native blockers, the platform maximum; do not add, remove or repurpose an edge there in this PR.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback, concurrency, replay and reload tests with `-tags=integration`; shared web and TUI movement fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. The fixture uses explicit isolated entitlement and grants no production eligibility. No deployment, source inference or Full Rules promotion.
