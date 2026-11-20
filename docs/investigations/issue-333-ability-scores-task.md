<!-- edda-333:ability-scores -->
Decision #333; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1 PDF20–21 (Step 3: Ability Scores) and the background Ability Scores part on83; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Retain CC-BY-4.0 attribution and adaptation notice in `docs/srd-5.2.1-conformance.md`. Fixture and ownership decisions are Edda additions.

Implement the one creation step that generates, assigns and adjusts the six ability scores. Exactly one of three methods is recorded, and a campaign policy set by the GM may restrict which are offered:

- Standard Array: exactly the multiset 15, 14, 13, 12, 10, 8, each used once.
- Random Generation: six times, roll four d6 and record the total of the highest three. Save all 24 injected dice; each total is 3–18.
- Point Cost: 27 points; cost by score 8:0, 9:1, 10:2, 11:3, 12:4, 13:5, 14:7, 15:9. Scores below 8 or above 15 have no cost row and reject. A total above 27 rejects. The source states no minimum spend, so a valid underspent set is accepted and unspent points grant nothing.

Assign the six generated numbers to Strength, Dexterity, Constitution, Intelligence, Wisdom and Charisma, each number once. The Standard Array by Class table on PDF21 is a suggestion only and is never enforced. Then apply the selected background's adjustment from the #306 records: among its three listed abilities, +2 to one and +1 to a different one, or +1 to all three; no increase may raise a score above 20, and an ability the background does not list cannot be increased. Derive modifiers from the PDF21 table exactly: 3:−4, 4–5:−3, 6–7:−2, 8–9:−1, 10–11:+0, 12–13:+1, 14–15:+2, 16–17:+3, 18–19:+4, 20:+5. That table has no row below 3 or above 20; do not invent one here.

This task owns generation, assignment, the background adjustment and the recorded method. Background identity, feat, proficiencies and equipment belong to the #333 background task; later feat increases stay with #384/#389; custom backgrounds reuse this adjustment through their own task. The provenance contract stays with #430.

Positive fixtures: Standard Array assigned STR15/DEX14/CON13/INT8/WIS10/CHA12 with Soldier +2 STR/+1 CON gives STR17 (+3), CON14 (+2). Soldier +1 STR/DEX/CON gives 16/15/14. Point Cost 15,15,15,8,8,8 costs27; 15,14,13,12,10,8 costs27. Injected dice 6,5,4,1 give15; 1,1,1,1 give3; 6,6,6,2 give18. Reject: array value reused or replaced; point cost28; purchased 16 or 7; five or seven scores; client-supplied random total that does not match saved dice; +2 and +1 on the same ability; +2/+1/+1; an increase to an unlisted ability (Soldier Wisdom); a method the campaign policy disallows; regeneration after confirmation; a second roll set kept alongside the first.

Preview is inert: no dice, writes or history. Confirm revalidates campaign, actor, player/GM authority, source version, every choice and the expected revision, then persists choices, injected dice, grants, resource deltas, derived sheet and public history atomically. Fail at each stage and prove full isolated PostgreSQL rollback. Same-ID retry/reload/reconnect returns the saved receipt without reroll, regrant or recharge; a conflicting payload rejects and concurrent confirmations commit once. Shared web and trusted TUI expose only qualified choices and explicit unavailable reasons, with accessible cancel/confirm/reconnect and desktop/phone checks for changed web. Keep private GM facts, hidden actors and private choices out of discovery, preview, errors and public history. Tests use explicit isolated entitlements and grant no production admission.

Native parent #333; directly block #337. Hard prerequisites: #284, #306 and #430. Existing bodies, parents and edges remain unchanged. No new #283 child.

Implementation verification: deterministic source/numeric/choice tests with injectable dice, isolated PostgreSQL atomic rollback/replay/concurrency/reload, shared web/TUI disclosure/confirmation checks, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests do not qualify PostgreSQL or devices. No production grant, deployment or Full Rules promotion.
