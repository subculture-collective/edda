<!-- edda-333:species-halfling -->
Decision #333; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1 PDF86 (Halfling); consumed general rules: Interactions with Rerolls8, Moving around Other Creatures14, Hide183; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Retain CC-BY-4.0 attribution and adaptation notice in `docs/srd-5.2.1-conformance.md`. Fixture and ownership decisions are Edda additions.

This task owns only the species identity, its nested choices and the traits named here. Class creation and level applications consume its output and apply each grant once; it grants no class, background or feat. The creation provenance contract stays with #430. Halfling is Humanoid, Small, Speed30. It has no nested choice and no limited-use resource.

- Brave: Advantage on saving throws made to avoid or end the Frightened condition.
- Halfling Nimbleness: can move through the space of any creature that is a size larger than the Halfling, but cannot stop in the same space.
- Luck: on rolling a 1 on the d20 of a D20 Test, may reroll the die and must use the new roll.
- Naturally Stealthy: can take the Hide action even when obscured only by a creature that is at least one size larger.

Consume #61 for the D20 Test receipt and save Advantage, #488 for Frightened, #67 for scene positions and sizes, and #1137 for the Hide action itself; this task adds only the extra permitted obscurement. Luck is offered immediately on the Halfling's own d20 showing 1, before the outcome is applied, for ability checks, attack rolls and saving throws. With Advantage or Disadvantage, PDF8 allows only one die to be rerolled, and only a die that shows 1 qualifies. The new roll stands even if it is another 1; no second Luck reroll applies to that new roll. Luck never applies to damage or other dice and is distinct from Heroic Inspiration #774. Sizes, occupied spaces and who obscures whom are authoritative scene facts; unknown facts keep Nimbleness and Naturally Stealthy unavailable.

Positive fixtures: a save to end Frightened rolls with Advantage. Injected d20 1 then14 uses14; 1 then1 stays1. With Advantage, dice1 and9 may reroll the1; injected17 gives17. A Small Halfling moves through a declared Medium creature's space and ends beyond it. Behind a declared Medium ally with no other cover, the Hide action is permitted and resolves through #1137. Reject: Luck on a 2, on a damage die, on another creature's roll, after the outcome is committed or twice on one die; keeping the original 1 after choosing to reroll; moving through a Small creature's space by this trait or ending the move inside a larger creature's space; Naturally Stealthy behind a Small creature; Advantage on saves unrelated to Frightened.

Preview is inert: no dice, writes or history. Confirm revalidates campaign, actor, player/GM authority, source version, every choice and the expected revision, then persists choices, injected dice, grants, resource deltas, derived sheet and public history atomically. Fail at each stage and prove full isolated PostgreSQL rollback. Same-ID retry/reload/reconnect returns the saved receipt without reroll, regrant or recharge; a conflicting payload rejects and concurrent confirmations commit once. Shared web and trusted TUI expose only qualified choices and explicit unavailable reasons, with accessible cancel/confirm/reconnect and desktop/phone checks for changed web. Keep private GM facts, hidden actors and private choices out of discovery, preview, errors and public history. Tests use explicit isolated entitlements and grant no production admission.

Native parent #333; directly block #337. Hard prerequisites: #61, #67, #307, #430, #488 and #1137. Existing bodies, parents and edges remain unchanged. No new #283 child.

Implementation verification: deterministic source/numeric/choice tests with injectable dice, isolated PostgreSQL atomic rollback/replay/concurrency/reload, shared web/TUI disclosure/confirmation checks, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests do not qualify PostgreSQL or devices. No production grant, deployment or Full Rules promotion.
