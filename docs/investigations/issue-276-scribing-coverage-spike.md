<!-- edda-276:scribing-coverage-spike -->
Decision #276; roadmap #47. Planning spike: one decision PR plus published bounded follow-ups; Task / Size M; label Planning/Spike.

Source: English SRD 5.2.1, PDF 103 and 244; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`.

## Named scope

"Scribing Spell Scrolls" on PDF103 beyond the #276 scroll scribing task (marker `edda-276:scroll-scribing`): the level 2 and level 4–9 rows of the Spell Scroll Costs table; scribes whose preparation comes from Bard, Druid, Paladin, Ranger, Sorcerer, Warlock or Wizard outputs; multiclass scribes, including which level "your level" means for a scribed cantrip; and cantrips known from a source other than the scribe's own class. PDF103 says "A spellcaster can transfer a spell to a scroll", so each class needs its actual preparation owner, as #557 is for the Cleric. Reading stays with #1142 and the #276 scribed scroll reading task.

## Acceptance criteria

- [ ] Inspect only the named scope: Spell Scroll scribing beyond a single-class Cleric's cantrip, level 1 and level 3 scrolls. Record exact source rules, existing owners and unknowns from the pinned PDF pages; do not guess production mechanics.
- [ ] Search open and closed issue titles and bodies for an existing owner of each clause. Publish a one-PR Task only for an unowned, source-backed residual, with a named fixture, source pages, costs and timing, durable/replay/disclosure tests, exclusions and hard prerequisites. Split by tool, class or category so each child is one reviewable PR.
- [ ] Attach each child under this spike. Before closing, add every child as a native blocker of the #276 crafting qualification task (marker `edda-276:crafting-qualification`), or of one aggregate task that blocks it if its 100-entry limit would be exceeded. Read back every child, parent and edge.
- [ ] Keep each source unknown unavailable and name the GM ruling it needs.

Hard prerequisites: none. The children extend the #276 procedure task named above and are blocked by it. Parent: #276. This spike directly blocks the #276 crafting qualification task, so item qualification #280 stays blocked until its children exist and qualify. Closing this spike is not implementation evidence. Add no #279, #67 or #429 edge and alter no existing owner. No runtime, migration, UI or deployment change.
