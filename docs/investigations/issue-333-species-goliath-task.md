<!-- edda-333:species-goliath -->
Decision #333; roadmap #47. One implementation PR; Task / Size L.

Source: English SRD 5.2.1 PDF85–86 (Goliath); consumed general rules: Carrying Capacity178, Teleportation190; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Retain CC-BY-4.0 attribution and adaptation notice in `docs/srd-5.2.1-conformance.md`. Fixture and ownership decisions are Edda additions.

This task owns only the species identity, its nested choices and the traits named here. Class creation and level applications consume its output and apply each grant once; it grants no class, background or feat. The creation provenance contract stays with #430. Goliath is Humanoid, Medium, Speed35.

- Giant Ancestry: choose exactly one benefit at creation. Its uses equal the Proficiency Bonus and all return on a completed Long Rest.
  - Cloud's Jaunt: Bonus Action, magically teleport up to 30 feet to an unoccupied space the Goliath can see.
  - Fire's Burn: on hitting a target with an attack roll and dealing damage, also deal 1d10 Fire damage to it.
  - Frost's Chill: on such a hit, also deal 1d6 Cold damage and reduce the target's Speed by 10 feet until the start of the Goliath's next turn.
  - Hill's Tumble: on such a hit against a Large or smaller creature, give it the Prone condition.
  - Stone's Endurance: on taking damage, a Reaction rolls 1d12, adds the Constitution modifier and reduces the damage by that total.
  - Storm's Thunder: on taking damage from a creature within 60 feet, a Reaction deals 1d8 Thunder damage to that creature.
- Large Form, from character level5: a Bonus Action changes size to Large if the space is big enough, for 10 minutes or until ended (no action). During it: Advantage on Strength checks and Speed +10 feet. Once used, unavailable until a completed Long Rest.
- Powerful Build: Advantage on any ability check made to end the Grappled condition, and the Goliath counts as one size larger when determining carrying capacity.

Consume #65 for the hit and damage receipt the on-hit options attach to, #66 for Bonus Actions, #81 for the Reaction windows, #61 for check Advantage, #488 for Prone, #778 for the Grappled escape check, #428 for recovery and #968 for the level5 threshold. On-hit options are optional, spend one use each and are declared on the confirmed hit, never retroactively. Carrying capacity follows the PDF178 table for the Goliath's own sheet: a Medium Goliath uses the Large row (Strength×30 lb. carry, Strength×60 lb. drag/lift/push), and in Large Form the Huge row. Visible unoccupied destinations, target size, distance and room for Large Form are authoritative scene facts; unknown facts keep the benefit unavailable.

Positive fixtures: PB2 gives two uses. Fire's Burn injected6 adds6 Fire to a confirmed hit. Frost's Chill injected3 adds3 Cold and Speed30 becomes20 until the Goliath's next turn starts. Hill's Tumble knocks a Large target Prone and is refused against Huge. Stone's Endurance injected7 with CON+2 reduces9 damage to0, never below. Storm's Thunder injected5 deals5 to an attacker at60 feet and is refused at65. Level5 Large Form gives Speed45 and Strength-check Advantage. Strength16 carries480 lb. and drags960. Reject: two ancestry benefits or a change after creation; an on-hit option on a miss or on a hit that dealt no damage; a third use at PB2; a Reaction when none is available; teleport to an occupied or unseen space or beyond30 feet; Large Form at level4, in a declared too-small space or twice before a Long Rest; uses restored by a Short Rest; Advantage on checks other than to end Grappled.

Preview is inert: no dice, writes or history. Confirm revalidates campaign, actor, player/GM authority, source version, every choice and the expected revision, then persists choices, injected dice, grants, resource deltas, derived sheet and public history atomically. Fail at each stage and prove full isolated PostgreSQL rollback. Same-ID retry/reload/reconnect returns the saved receipt without reroll, regrant or recharge; a conflicting payload rejects and concurrent confirmations commit once. Shared web and trusted TUI expose only qualified choices and explicit unavailable reasons, with accessible cancel/confirm/reconnect and desktop/phone checks for changed web. Keep private GM facts, hidden actors and private choices out of discovery, preview, errors and public history. Tests use explicit isolated entitlements and grant no production admission.

Native parent #333; directly block #337. Hard prerequisites: #61, #65, #66, #81, #307, #428, #430, #488, #778 and #968. Existing bodies, parents and edges remain unchanged. No new #283 child.

Implementation verification: deterministic source/numeric/choice tests with injectable dice, isolated PostgreSQL atomic rollback/replay/concurrency/reload, shared web/TUI disclosure/confirmation checks, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests do not qualify PostgreSQL or devices. No production grant, deployment or Full Rules promotion.
