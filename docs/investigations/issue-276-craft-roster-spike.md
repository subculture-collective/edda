<!-- edda-276:craft-roster-spike -->
Decision #276; roadmap #47. Planning spike: one decision PR plus published bounded follow-ups; Task / Size M; label Planning/Spike.

Source: English SRD 5.2.1, PDF 91–100 and 103; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`.

## Named scope

Every tool Craft entry on PDF93–94 other than the three products of the #276 nonmagical crafting task (marker `edda-276:nonmagical-crafting`: Light Hammer, Longbow, Plate Armor). That is the named products and class entries of Alchemist's, Brewer's, Calligrapher's (Ink only), Carpenter's, Cartographer's, Cobbler's, Cook's, Glassblower's, Jeweler's, Leatherworker's, Mason's, Painter's, Potter's, Smith's, Tinker's, Weaver's and Woodcarver's tools, the Disguise Kit, the Herbalism Kit (Antitoxin, Candle, Healer's Kit) and the Poisoner's Kit, resolved against catalogue records #247–#262.

Known open points to settle from the source, not by inference: whether the tool records #253–#255 carry structured Craft entries; products listed under two tools (Club, Greatclub, Quarterstaff; Antitoxin; Holy Symbol; Arcane Focus; Druidic Focus), for which PDF103 says only "You must use the required tool"; variable-price entries (Arcane Focus, Druidic Focus, Holy Symbol, Ammunition); prices whose half is not a whole number of Copper Pieces; GM-assigned tools for unlisted items; and Heavy Crossbow, which is on PDF91 but absent from the #80 inventory and #247–#250. Forgery Kit, Gaming Set, Musical Instrument, Navigator's Tools and Thieves' Tools have no Craft entry and grant nothing.

## Acceptance criteria

- [ ] Inspect only the named scope: the tool Craft entries on PDF93–94 that the #276 nonmagical crafting task does not execute. Record exact source rules, existing owners and unknowns from the pinned PDF pages; do not guess production mechanics.
- [ ] Search open and closed issue titles and bodies for an existing owner of each clause. Publish a one-PR Task only for an unowned, source-backed residual, with a named fixture, source pages, costs and timing, durable/replay/disclosure tests, exclusions and hard prerequisites. Split by tool, class or category so each child is one reviewable PR.
- [ ] Attach each child under this spike. Before closing, add every child as a native blocker of the #276 crafting qualification task (marker `edda-276:crafting-qualification`), or of one aggregate task that blocks it if its 100-entry limit would be exceeded. Read back every child, parent and edge.
- [ ] Keep each source unknown unavailable and name the GM ruling it needs.

Hard prerequisites: none. The children extend the #276 procedure task named above and are blocked by it. Parent: #276. This spike directly blocks the #276 crafting qualification task, so item qualification #280 stays blocked until its children exist and qualify. Closing this spike is not implementation evidence. Add no #279, #67 or #429 edge and alter no existing owner. No runtime, migration, UI or deployment change.
