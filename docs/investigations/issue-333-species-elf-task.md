<!-- edda-333:species-elf -->
Decision #333; roadmap #47. One implementation PR; Task / Size L.

Source: English SRD 5.2.1 PDF84–85 (Elf and the Elven Lineages table); SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Retain CC-BY-4.0 attribution and adaptation notice in `docs/srd-5.2.1-conformance.md`. Fixture and ownership decisions are Edda additions.

This task owns only the species identity, its nested choices and the traits named here. Class creation and level applications consume its output and apply each grant once; it grants no class, background or feat. The creation provenance contract stays with #430. Elf is Humanoid, Medium, Speed30.

- Darkvision 60 feet.
- Elven Lineage: choose Drow, High Elf or Wood Elf at creation, and at the same time choose Intelligence, Wisdom or Charisma as the spellcasting ability for this trait's spells.
  - Drow: Darkvision becomes 120 feet; knows Dancing Lights. Level3 Faerie Fire; level5 Darkness.
  - High Elf: knows Prestidigitation, and after each completed Long Rest may replace that cantrip with a different cantrip from the Wizard spell list. Level3 Detect Magic; level5 Misty Step.
  - Wood Elf: Speed35; knows Druidcraft. Level3 Longstrider; level5 Pass without Trace.
  The level3 and level5 spells are gained at those character levels and are always prepared. Each can be cast once without a spell slot, regained on a completed Long Rest, and can also be cast with any slot of the appropriate level the character has.
- Fey Ancestry: Advantage on saving throws made to avoid or end the Charmed condition.
- Keen Senses: proficiency in exactly one of Insight, Perception or Survival, chosen at creation.
- Trance: no need to sleep; magic cannot put the Elf to sleep; a Long Rest can be finished in 4 hours of trancelike meditation with consciousness retained.

This task owns eligibility, the per-spell free use and its recovery. Slot versus free-use routing stays with #1448. Each spell executes only through its own effect owner: Dancing Lights #1281, Darkness #1282, Prestidigitation #1384–#1386, Detect Magic #265/#906, Misty Step #688, Druidcraft #1049/#1050, Longstrider #1066, Pass without Trace #1072. Faerie Fire has no published effect owner; its key is still assigned by #277, so the Drow level3 grant is recorded and its cast stays unavailable until that owner is published and qualified. A High Elf replacement cantrip is offered only when its own Wizard cantrip effect is qualified. Consume #61 and #488 for the Charmed save Advantage, #145 for Darkvision, #430 for Keen Senses provenance, #428 for Long Rest completion and the 4-hour Trance, and #968 for the level3/5 thresholds. Trance shortens only the Elf's own rest; other rest rules stay with #428.

Positive fixtures: Drow level1 has Darkvision120 and Dancing Lights with the chosen ability. Wood Elf has Speed35. High Elf level3 casts Detect Magic once without a slot, then only with a level1 or higher slot until a completed Long Rest; at level5 Misty Step has its own separate free use and needs a level2 or higher slot afterwards. Keen Senses Perception duplicating a class skill follows #430. A sleep effect on the Elf is refused by its own owner through this trait. Reject: no lineage or two; ability chosen outside the three or changed later; level3 spell at level2; one shared free use across two spells; free use restored by a Short Rest; a lineage spell cast with a slot below its level; a replacement cantrip from another list, a second replacement in one rest or replacement by a non-High Elf; two Keen Senses skills; Advantage on saves unrelated to Charmed.

Preview is inert: no dice, writes or history. Confirm revalidates campaign, actor, player/GM authority, source version, every choice and the expected revision, then persists choices, injected dice, grants, resource deltas, derived sheet and public history atomically. Fail at each stage and prove full isolated PostgreSQL rollback. Same-ID retry/reload/reconnect returns the saved receipt without reroll, regrant or recharge; a conflicting payload rejects and concurrent confirmations commit once. Shared web and trusted TUI expose only qualified choices and explicit unavailable reasons, with accessible cancel/confirm/reconnect and desktop/phone checks for changed web. Keep private GM facts, hidden actors and private choices out of discovery, preview, errors and public history. Tests use explicit isolated entitlements and grant no production admission.

Native parent #333; directly block #337. Hard prerequisites: #61, #145, #265, #277, #307, #428, #430, #488, #688, #906, #968, #1049, #1050, #1066, #1072, #1281, #1282, #1384, #1385, #1386 and #1448. Existing bodies, parents and edges remain unchanged. No new #283 child.

Implementation verification: deterministic source/numeric/choice tests with injectable dice, isolated PostgreSQL atomic rollback/replay/concurrency/reload, shared web/TUI disclosure/confirmation checks, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests do not qualify PostgreSQL or devices. No production grant, deployment or Full Rules promotion.
