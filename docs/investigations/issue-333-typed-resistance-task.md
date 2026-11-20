<!-- edda-333:typed-resistance -->
Decision #333; roadmap #47. One implementation PR; Task / Size S.

Source: English SRD 5.2.1 PDF17 (Resistance and Vulnerability; Order of Application), the Rules Glossary entry Resistance on187 and Round Down on5; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Retain CC-BY-4.0 attribution and adaptation notice in `docs/srd-5.2.1-conformance.md`. Fixture and ownership decisions are Edda additions.

Implement the one shared projection of Resistance to a damage type in the damage pipeline. No published issue owns it: existing features that grant Resistance, such as Rage in #435, name it as something they consume. This task is the single owner; a feature registers a named source of Resistance to a type and never halves damage itself.

- A creature holds zero or more named sources of Resistance per damage type, each with its own origin and ending.
- When that creature takes damage of a type it resists, the damage of that type is halved and rounded down.
- Resistance is applied only once to an instance of damage, however many sources grant it.
- Order follows PDF17: bonuses, penalties and multipliers first, then Resistance. A save for half is part of the first step.
- An instance with several damage types halves only the resisted part.

Vulnerability and Immunity are the third step and a separate rule on the same page; they have no published owner and are not added here, so a creature with either keeps that part unavailable. This task grants no Resistance: the Dragonborn, Dwarf and Tiefling tasks of decision #333 and every class or spell owner register their own sources. If an existing consumer such as #435 lands first and has already added a Resistance projection, this task adopts that one projection, completes it to the rules above and adds no second store. Consume #65 for the committed damage event.

Positive fixtures: 9 Fire damage against Resistance to Fire is4. Two sources of Resistance to Fire still give4. 1 becomes0. 28 Fire damage reduced by5 to23 is then halved to11, the source's own example without its Vulnerability step. An instance of6 Fire and5 Slashing against Resistance to Fire is3 and5. A failed save for11 and a successful save for5 become5 and2. Removing one of two sources leaves Resistance; removing the last ends it. Reject: halving a type the creature does not resist; halving twice for two sources; rounding up; halving before a reduction or before the save result; halving healing or Temporary Hit Points gained; a client-supplied halved amount; a source with no origin; a second resistance store; a Vulnerability or Immunity result invented here.

The projection has no command of its own. It is applied inside the confirmed damage event, in the same atomic write and rollback as that event: an injected failure leaves Hit Points and history unchanged, and same-ID replay returns the saved receipt without halving again. The receipt records the amount before and after and the source used. Shared web and trusted TUI show a character's own Resistances with their sources on the sheet and the applied halving on the public damage receipt, with desktop/phone checks for changed web; a hidden actor's Resistances are never disclosed before they apply. Tests use explicit isolated entitlements and grant no production admission.

Native parent #333; directly block #337. Hard prerequisites: #65. Existing bodies, parents and edges remain unchanged. No new #283 child.

Implementation verification: deterministic source/numeric/choice tests with injectable dice, isolated PostgreSQL atomic rollback/replay/concurrency/reload, shared web/TUI disclosure/confirmation checks, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests do not qualify PostgreSQL or devices. No production grant, deployment or Full Rules promotion.
