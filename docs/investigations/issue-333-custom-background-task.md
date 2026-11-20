<!-- edda-333:custom-background -->
Decision #333; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1 PDF192–193 (Creating a Background, steps1–5), 19 (the GM might offer additional backgrounds) and83 (Parts of a Background); SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Retain CC-BY-4.0 attribution and adaptation notice in `docs/srd-5.2.1-conformance.md`. Fixture and ownership decisions are Edda additions.

Implement creation or customization of a background by the GM as campaign content, and its selection at character creation. Only an authoritative GM may author one; a player selects it only after the GM publishes it to the campaign. The five steps are exact:

1. Choose three abilities. The source's notes about which abilities suit which occupation are advice and are not enforced. The adjustment itself (+2/+1 or +1/+1/+1, cap20) is the #333 ability-score task's rule applied to these three.
2. Choose one feat from the Origin category. In SRD 5.2.1 that is Alert, Magic Initiate, Savage Attacker or Skilled. The grant goes through #384 and executes only through #385–#388; Magic Initiate's list and Skilled's three proficiencies are fixed by the GM or left to the feat owner's own choices, never invented here.
3. Choose two skills. They need not relate to the abilities.
4. Choose one tool.
5. Assemble an equipment package worth 50GP including unspent gold, with no Martial weapons and no armor. Value uses only catalogued item costs; an item without a catalogued cost makes the package unavailable. PDF83 says each background offers a choice between its package and 50GP, so the player selecting an authored background chooses that package or 50GP.

Customizing one of the four printed backgrounds creates a new campaign-scoped record under these same steps; the printed record is never edited. A selected authored background grants exactly what the #333 background task grants for printed ones: one feat, two skills, one tool and the package once, with provenance under #430. The record is versioned; a character keeps the version it was created with.

Positive fixtures: a GM-authored Deckhand record with abilities STR/DEX/WIS, feat Skilled, skills Athletics and Perception, tool Navigator's Tools and a package of catalogued items totalling exactly 50GP with remaining coins selects and grants once. A customized Soldier that swaps Intimidation for Survival is a distinct record. Reject: player-authored record; two or four abilities or a repeated ability; a General, Fighting Style or Epic Boon feat; one or three skills; no tool or two tools; a package worth 49 or 51GP; a Martial weapon; any armor or Shield; an uncatalogued item; editing a printed background; editing a version already used by a created character; selection before publication or from another campaign.

Preview is inert: no dice, writes or history. Confirm revalidates campaign, actor, player/GM authority, source version, every choice and the expected revision, then persists choices, injected dice, grants, resource deltas, derived sheet and public history atomically. Fail at each stage and prove full isolated PostgreSQL rollback. Same-ID retry/reload/reconnect returns the saved receipt without reroll, regrant or recharge; a conflicting payload rejects and concurrent confirmations commit once. Shared web and trusted TUI expose only qualified choices and explicit unavailable reasons, with accessible cancel/confirm/reconnect and desktop/phone checks for changed web. Keep private GM facts, hidden actors and private choices out of discovery, preview, errors and public history. Tests use explicit isolated entitlements and grant no production admission.

Native parent #333; directly block #337. Hard prerequisites: the #333 ability-score and background tasks published by this decision, and #301, #306, #384, #385, #386, #387, #388 and #430. Existing bodies, parents and edges remain unchanged. No new #283 child.

Implementation verification: deterministic source/numeric/choice tests with injectable dice, isolated PostgreSQL atomic rollback/replay/concurrency/reload, shared web/TUI disclosure/confirmation checks, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests do not qualify PostgreSQL or devices. No production grant, deployment or Full Rules promotion.
