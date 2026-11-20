# Stat-block batch 03 (#153)

This batch catalogues the twelve stat blocks the
[#82 inventory](investigations/issue-82-world-inventory.json) assigns to #153.
It is data only. Every record is unavailable: `Monster.Available` implements no
mechanic key, the embedded runtime catalogue still holds only the Bandit, and no
resolver, encounter, tool or coverage entry reads the batch. Source pin, license
and attribution: [source provenance](srd-5.2.1-conformance.md).

## Records

| Record | PDF page | Inventory family |
| --- | --- | --- |
| Animated Armor | 259 | monsters |
| Animated Flying Sword | 259 | monsters |
| Animated Rug of Smothering | 259 | monsters |
| Ankheg | 259 | monsters |
| Ankylosaurus | 344 | animals |
| Ape | 344 | animals |
| Archelon | 344–345 | animals |
| Archmage | 305 | monsters |
| Assassin | 260 | monsters |
| Awakened Shrub | 260 | monsters |
| Awakened Tree | 260 | monsters |
| Axe Beak | 260 | monsters |

The Archelon's abilities, trait and actions continue on p. 345. IDs follow the
inventory, so the three p. 344 blocks use `srd-5.2.1:content/animals/<slug>`.

## Layout

| File | Content |
| --- | --- |
| `internal/rules/srd521/testdata/monsters_153.json` | The batch data |
| `internal/rules/srd521/monsters_153_test.go` | Literal source values, inventory, dependency-key, rejection and composition tests |
| `scripts/verify_srd_monsters_153.py` | Comparison with the pinned PDF |

The layout and record conventions follow the first stat-block batch (#151):
the #143 `Monster` record with the optional stat-block fields in
`monster_stat_block.go` (tags, skills, damage and condition defenses, hover,
bonus actions, reactions, limited uses, saves, targets, effect dice and spell
lists). Each present field requires a dependency key, and removing any listed
key invalidates the record. This is the first batch with bonus actions,
reactions, gear and vulnerabilities.

Two shared edits were needed. `Monster.Validate` accepts the inventory's
`content/animals/` ID prefix. `monster_source_exceptions.go` adds
`monsterPrintedXP`, described below.

## Source values the contract cannot type

These stay in the entry's retained source text and are pinned by the tests and
the PDF comparison. None is turned into a typed attack.

- **Awakened Shrub, Rake** deals a fixed 1 Slashing damage with no dice. Its key
  is `action/fixed-damage-melee-attack`.
- **Archmage, Arcane Burst** is a "Melee or Ranged Attack Roll" with reach 5 ft.
  or range 150 ft. Its key is `action/melee-or-ranged-attack`; the damage dice
  are typed on its effect.
- **Escape DCs** (Animated Rug of Smothering and Ankheg, DC 13), target size
  limits and the Ankheg's Advantage against a creature it grapples.
- **Archmage immunities** print "Psychic; Charmed (with Mind Blank)". The record
  lists both and requires `defense/immunity-with-mind-blank` through the Mind
  Blank spell entry.
- **Languages** such as "Common plus five other languages" are one list entry.
  "Languages None" is an empty list.

## Source inconsistencies kept as printed

- The Archmage prints `CR 12 (XP 8,000; PB +4)`. The table on p. 256 gives
  8,400 XP for CR 12. The record keeps 8,000. `monsterPrintedXP` accepts only
  that exact pair for that record, and the encounter budget still rejects it.
- The Archmage's Misty Step entry says "The mage casts Misty Step".
- The Ankheg's Acid Spray and the Ape's Rock print "Recharge 6", stored as
  recharge 6.

## Checks

```sh
rtk proxy python3 -I scripts/verify_srd_monsters_153.py /path/to/SRD_CC_v5.2.1.pdf
rtk proxy env GOMAXPROCS=2 go test -p 1 -short -race ./internal/rules/srd521/...
rtk proxy go vet ./internal/rules/srd521/...
```

The comparison supplements the offline literal fixtures; it is not an automatic
conformance proof. The PDF and extracted text are not vendored.

This work includes material from the System Reference Document 5.2.1 (“SRD
5.2.1”) by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd.
The SRD 5.2.1 is licensed under the Creative Commons Attribution 4.0 International
License, available at https://creativecommons.org/licenses/by/4.0/legalcode.
Names, numeric values and rules text are adapted into Edda's structured monster
records; the software license does not replace this source license.
