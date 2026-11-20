<!-- edda-333:species-elf -->
Decision #333; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1 PDF84–85 (Elf; level1 row of the Elven Lineages table); SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Retain CC-BY-4.0 attribution and adaptation notice in `docs/srd-5.2.1-conformance.md`. Fixture and ownership decisions are Edda additions.

This task owns only the species identity, its nested choices and the traits named here. Class creation and level applications consume its output and apply each grant once; it grants no class, background or feat. The creation provenance contract stays with #430 and is consumed through the #333 creation provenance output. Elf is Humanoid, Medium, Speed30.

- Darkvision 60 feet.
- Elven Lineage: choose Drow, High Elf or Wood Elf at creation, and at the same time choose Intelligence, Wisdom or Charisma as the spellcasting ability for this trait's spells. Level1 benefits:
  - Drow: Darkvision becomes 120 feet; knows Dancing Lights.
  - High Elf: knows Prestidigitation.
  - Wood Elf: Speed35; knows Druidcraft.
- Fey Ancestry: Advantage on saving throws made to avoid or end the Charmed condition.
- Keen Senses: proficiency in exactly one of Insight, Perception or Survival, chosen at creation.
- Trance: no need to sleep; magic cannot put the Elf to sleep; a Long Rest can be finished in 4 hours of trancelike meditation with consciousness retained.

The level3 and level5 lineage spells and the High Elf cantrip replacement belong to the #333 Elf lineage spell task, and the Drow level3 Faerie Fire cast to the #333 Faerie Fire task; this task records the lineage and ability they read and offers none of them. Each cantrip executes only through its own effect owner: Dancing Lights #1281, Prestidigitation #1384–#1386, Druidcraft #1049/#1050. Consume #61 for the save Advantage, #488 for Charmed through the #333 condition output, #145 for Darkvision, #430 for Keen Senses provenance through the #333 creation provenance output, and #428 for Long Rest completion and the 4-hour Trance. Trance shortens only the Elf's own rest; other rest rules stay with #428.

Positive fixtures: Drow level1 has Darkvision120 and Dancing Lights with the chosen ability. High Elf casts Prestidigitation with the chosen ability, not a class ability. Wood Elf has Speed35 and Druidcraft. A save to end Charmed rolls with Advantage. Keen Senses Perception duplicating a class skill follows #430. An Elf completes a Long Rest after 4 hours of Trance. A sleep effect on the Elf is refused by its own owner through this trait. Reject: no lineage or two; ability chosen outside the three or changed later; a cantrip from another lineage; Wood Elf Speed or Drow Darkvision on another lineage; two Keen Senses skills or one outside the three; Advantage on saves unrelated to Charmed; a 4-hour Long Rest for a non-Elf; any level3 or level5 lineage spell offered by this task.

Preview is inert: no dice, writes or history. Confirm revalidates campaign, actor, player/GM authority, source version, every choice and the expected revision, then persists choices, injected dice, grants, resource deltas, derived sheet and public history atomically. Fail at each stage and prove full isolated PostgreSQL rollback. Same-ID retry/reload/reconnect returns the saved receipt without reroll, regrant or recharge; a conflicting payload rejects and concurrent confirmations commit once. Shared web and trusted TUI expose only qualified choices and explicit unavailable reasons, with accessible cancel/confirm/reconnect and desktop/phone checks for changed web. Keep private GM facts, hidden actors and private choices out of discovery, preview, errors and public history. Tests use explicit isolated entitlements and grant no production admission.

Native parent #333; directly block #337. Hard prerequisites: the #333 creation provenance output task published by this decision, the #333 condition output task published by this decision, #61, #145, #307, #428, #1049, #1050, #1281, #1384, #1385 and #1386. Existing bodies, parents and edges remain unchanged. No new #283 child.

Implementation verification: deterministic source/numeric/choice tests with injectable dice, isolated PostgreSQL atomic rollback/replay/concurrency/reload, shared web/TUI disclosure/confirmation checks, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests do not qualify PostgreSQL or devices. No production grant, deployment or Full Rules promotion.
