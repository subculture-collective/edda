<!-- edda-276:magic-item-coverage-spike -->
Decision #276; roadmap #47. Planning spike: one decision PR plus published bounded follow-ups; Task / Size M; label Planning/Spike.

Source: English SRD 5.2.1, PDF 206–207; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`.

## Named scope

"Crafting Magic Items" on PDF206–207 beyond the #276 magic-item crafting task (marker `edda-276:magic-item-crafting`), over the catalogue records #225–#246: the Ring, Rod, Staff, Weapon and Wondrous Item rows of the Magic Item Tools table and the Scroll row for scrolls other than Spell Scrolls; armor made with Leatherworker's or Weaver's Tools; the Common, Very Rare and Legendary rows; items that cast spells prepared through a class other than Cleric; incorporated items other than Plate Armor; and "the tool required to make the nonmagical item on which the magic item is based". Artifacts have no row and stay unavailable. Halved or shared times that are not whole days stay unavailable without a GM day-count ruling.

## Acceptance criteria

- [ ] Inspect only the named scope: magic-item crafting beyond Armor, +1 (Plate Armor), Potion of Resistance and Wand of Magic Detection. Record exact source rules, existing owners and unknowns from the pinned PDF pages; do not guess production mechanics.
- [ ] Search open and closed issue titles and bodies for an existing owner of each clause. Publish a one-PR Task only for an unowned, source-backed residual, with a named fixture, source pages, costs and timing, durable/replay/disclosure tests, exclusions and hard prerequisites. Split by tool, class or category so each child is one reviewable PR.
- [ ] Attach each child under this spike. Before closing, add every child as a native blocker of the #276 crafting qualification task (marker `edda-276:crafting-qualification`), or of one aggregate task that blocks it if its 100-entry limit would be exceeded. Read back every child, parent and edge.
- [ ] Keep each source unknown unavailable and name the GM ruling it needs.

Hard prerequisites: none. The children extend the #276 procedure task named above and are blocked by it. Parent: #276. This spike directly blocks the #276 crafting qualification task, so item qualification #280 stays blocked until its children exist and qualify. Closing this spike is not implementation evidence. Add no #279, #67 or #429 edge and alter no existing owner. No runtime, migration, UI or deployment change.
