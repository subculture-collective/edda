<!-- edda-333:species-human -->
Decision #333; roadmap #47. One implementation PR; Task / Size S.

Source: English SRD 5.2.1 PDF86 (Human); SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Retain CC-BY-4.0 attribution and adaptation notice in `docs/srd-5.2.1-conformance.md`. Fixture and ownership decisions are Edda additions.

This task owns only the species identity, its nested choices and the traits named here. Class creation and level applications consume its output and apply each grant once; it grants no class, background or feat. The creation provenance contract stays with #430. Human is Humanoid, Speed30. The existing fixed Human Soldier Fighter and Human Acolyte Cleric builds keep working and gain no wider admission from this task.

- Size: Medium or Small, chosen when the species is selected and fixed afterwards.
- Resourceful: gains Heroic Inspiration whenever the Human finishes a Long Rest.
- Skillful: proficiency in one skill of the player's choice, from all eighteen skills.
- Versatile: one Origin feat of the player's choice. Skilled is the source's recommendation, never a default or requirement.

Resourceful is the gain only: #774 owns the single Heroic Inspiration pool, its maximum of one and its spending, and #428 owns Long Rest completion. A Human who already holds Heroic Inspiration at the end of the rest gains none; PDF8 lets that player give it to another player character who lacks it, and this task offers that as an explicit confirmed choice through #774's pool. Skillful's provenance and duplicate handling follow #430. Versatile grants through #384 and executes only through the chosen feat's owner, #385 Alert, #386 Magic Initiate, #387 Savage Attacker or #388 Skilled; it is a second, separate Origin feat beside the background's. A feat that is not repeatable cannot be taken twice across the two grants; repeatability follows each feat owner.

Positive fixtures: Small Human records size Small with Speed30. A completed Long Rest without Heroic Inspiration gives one; with one already held it gives none to the Human and may pass one to a declared party member who lacks it. Skillful Insight records species provenance. Versatile Skilled with three proficiencies and a Soldier's Savage Attacker coexist. Reject: no size or a later size change; Heroic Inspiration from a Short Rest or an interrupted Long Rest; a second held Inspiration; two Skillful skills; a General, Fighting Style or Epic Boon feat for Versatile; Skilled applied by default without a choice; a non-repeatable feat duplicated across background and Versatile; regrant on replay.

Preview is inert: no dice, writes or history. Confirm revalidates campaign, actor, player/GM authority, source version, every choice and the expected revision, then persists choices, injected dice, grants, resource deltas, derived sheet and public history atomically. Fail at each stage and prove full isolated PostgreSQL rollback. Same-ID retry/reload/reconnect returns the saved receipt without reroll, regrant or recharge; a conflicting payload rejects and concurrent confirmations commit once. Shared web and trusted TUI expose only qualified choices and explicit unavailable reasons, with accessible cancel/confirm/reconnect and desktop/phone checks for changed web. Keep private GM facts, hidden actors and private choices out of discovery, preview, errors and public history. Tests use explicit isolated entitlements and grant no production admission.

Native parent #333; directly block #337. Hard prerequisites: #307, #384, #385, #386, #387, #388, #428, #430 and #774. Existing bodies, parents and edges remain unchanged. No new #283 child.

Implementation verification: deterministic source/numeric/choice tests with injectable dice, isolated PostgreSQL atomic rollback/replay/concurrency/reload, shared web/TUI disclosure/confirmation checks, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests do not qualify PostgreSQL or devices. No production grant, deployment or Full Rules promotion.
