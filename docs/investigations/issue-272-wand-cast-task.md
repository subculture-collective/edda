<!-- edda-272:wand-of-magic-missiles-cast -->
Decision #272; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 146, 206, 251; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only one source identity, `srd-5.2.1:content/magic-items/wand-of-magic-missiles`: an Uncommon wand that needs no Attunement and has 7 charges. While holding it, the user can expend no more than 3 charges to cast Magic Missile from it: 1 charge casts the level 1 version and each additional charge raises the spell's level by 1. Reuse the actual Magic Missile output #1368 for the darts, targets and damage; do not copy it. This task owns the wand's charge cost for that cast and its last-charge destruction. The dawn recharge is the separate task published with this decision. No other item, learned spell, class admission or unrestricted spell dispatch is enabled.

PDF206 supplies the item-casting channel: no user spell slot, no components, the spell's normal casting time, range and duration. Magic Missile has no attack roll, saving throw or Concentration, so no spellcasting ability, save DC or Concentration is invented. #1140 is a build-order prerequisite, not a published generic seam: it stays limited to the Wand of Magic Detection, and this adapter extends the item charge record that task introduces instead of creating a second charge model. Neither #1140's body nor its item domain is widened.

If an expenditure removes the wand's last charge, roll 1d20 once after the charges are spent. On a 1 the wand crumbles into ashes and is destroyed; the source does not cancel the cast that spent the charge, so it still resolves. On 2–20 the wand remains with 0 charges. A destroyed wand is no longer a usable or rechargeable item. An expenditure that leaves at least 1 charge never rolls. If a separately qualified effect prevents the expenditure, as Thief #1132 does for its own entitled user, no last charge was expended and no roll occurs; this task does not implement that hook.

Charge knowledge stays with its owners. PDF206 states two reveal paths: Identify (#1613) and Attunement, which this wand does not have (attuned knowledge for other items stays with #278). The source states no rule for a holder who has neither. Edda disclosure adaptation: this task adds no other reveal path, so the count is shown only from an existing authoritative disclosure fact, and an unaffordable confirmation rejects without revealing it.

## Deterministic fixtures

- At 7 charges, spending 3 casts the level 3 version: #1368 resolves 5 darts, and 4 charges remain. Spending 1 casts the level 1 version with 3 darts.
- At 2 charges, spending 2 leaves 0 and rolls one d20: an injected 1 destroys the wand after the 4-dart cast resolves; an injected 2 keeps the wand at 0 charges.
- At 3 charges, spending 2 leaves 1 and rolls no d20.
- Reject 4 charges in one cast, a cast at 0 charges, a spend above the remaining count, a wand not held, a foreign actor or wand, a slot or component cost, an Attunement requirement, an attack roll or save DC, a d20 when charges remain, two d20 rolls for one expenditure, a cancelled cast on a roll of 1, a cast from a destroyed wand, a reroll, and a charge count disclosed by an error or unavailable reason.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, actor, source version, clock, targets, item ownership, holding, charges, the chosen charge count, source-bound #1368 availability and expected revisions. Commit the charge spend, the cast and its effect, the last-charge roll and destruction, and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the same receipt with no reroll; a conflicting payload rejects; racing confirmations commit once. The last-charge roll belongs to its triggering confirmed expenditure, with no optional activation. Dice are injectable. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires a focused build, lint and desktop/phone checks. Public history omits private HP, hidden actors, private checks and undisclosed item facts.

Hard prerequisites: #195 item contract, #1140 item charge record (build order; it brings its own prerequisites #906 and the Rogue execution cohort #1097), #1368 actual Magic Missile resolver. Parent: #272. Block the wand dawn-recharge task and the rest/recovery qualification task published with this decision; the qualification task blocks #280. Spell qualification #279 already requires the unchanged #1368 through the accepted finite shared aggregate; do not add a direct #279, #280, #67 or #429 edge and do not alter any existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this task and its callers qualify. No deployment, source inference or Full Rules promotion.
