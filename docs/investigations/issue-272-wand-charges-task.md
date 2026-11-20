<!-- edda-272:wand-of-magic-missiles-charges -->
Decision #272; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 146, 206, 251; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only one source identity: Wand of Magic Missiles, an Uncommon wand that needs no Attunement. It has 7 charges. While holding it, the user can expend no more than 3 charges to cast Magic Missile from it: 1 charge casts the level 1 version and each additional charge raises the spell's level by 1. Reuse the actual Magic Missile output #1368 for the darts, targets and damage; do not copy it. This adapter owns the wand's charge costs, its dawn recharge and its last-charge destruction. No other item, learned spell, class admission or unrestricted spell dispatch is enabled.

PDF206 supplies the item-casting channel: no user spell slot, no components, the spell's normal casting time, range and duration. Magic Missile has no attack roll, saving throw or Concentration, so no spellcasting ability, save DC or Concentration is invented. Reuse the charge-spend seam of #1140 and the authoritative dawn event of #1148; those tasks stay limited to the Wand of Magic Detection and are not widened. The wand regains 1d6 + 1 expended charges daily at dawn, never above 7. On a world or plane where dawn does not occur, the GM determines the recharge: require a trusted explicit GM recharge event and never fabricate time.

If an expenditure removes the wand's last charge, roll 1d20 once after the charges are spent. On a 1 the wand crumbles into ashes and is destroyed; the source does not cancel the cast that spent the charge, so it still resolves. On 2–20 the wand remains with 0 charges and recharges at a later dawn. A destroyed wand is no longer a usable or rechargeable item. An expenditure that leaves at least 1 charge never rolls. If a separately qualified effect prevents the expenditure, as Thief #1132 does for its own entitled user, no last charge was expended and no roll occurs; this task does not implement that hook.

PDF206 also governs what a creature knows. This wand has no Attunement, so the attuned-creature knowledge clause never applies to it. Disclose the remaining count only from an authoritative identification or disclosure fact, such as the Identify output #1613 or a trusted GM disclosure. Possession alone reveals nothing, and an unaffordable confirmation rejects without revealing the count.

## Deterministic fixtures

- At 7 charges, spending 3 casts the level 3 version: #1368 resolves 5 darts, and 4 charges remain. Spending 1 casts the level 1 version with 3 darts.
- At 2 charges, spending 2 leaves 0 and rolls one d20: an injected 1 destroys the wand after the 4-dart cast resolves; an injected 2 keeps the wand at 0 charges.
- Dawn at 0 charges with injected d6 = 3 gives 4. Dawn at 5 with d6 = 6 gives 7, not 12. Dawn at 7 stays 7.
- Reject 4 charges in one cast, a cast at 0 charges, a spend above the remaining count, a wand not held, a foreign actor or wand, a slot or component cost, an Attunement requirement, an attack roll or save DC, a d20 when charges remain, two d20 rolls for one expenditure, a cast or recharge from a destroyed wand, recharge on a rest, a client clock, a duplicate dawn, a reroll, and a fabricated dawn on a plane without one.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, actor, source version, clock, targets, item ownership, holding, charges, the chosen charge count, source-bound #1368 availability and expected revisions. Commit the charge spend, the cast and its effect, the last-charge roll and destruction, or the dawn event, die and recharge, and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the same receipt with no reroll; a conflicting payload rejects; racing confirmations commit once. The last-charge roll and the dawn recharge belong to their triggering confirmed event, with no optional activation. Dice are injectable. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires a focused build, lint and desktop/phone checks. Public history omits private HP, hidden actors, private checks and undisclosed item facts.

Hard prerequisites: #195 item contract, #1140 charge-spend seam, #1148 authoritative dawn recharge event, #1368 actual Magic Missile resolver. Parent: #272. Block the rest/recovery qualification task published with this decision; that task blocks #280. Spell qualification #279 already requires the unchanged #1368 through the accepted finite shared aggregate; do not add a direct #279, #280, #67 or #429 edge and do not alter any existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this task and its callers qualify. No deployment, source inference or Full Rules promotion.
