# Stat-block catalogue batch 02 (#152)

This batch catalogues the twelve index entries the #82 inventory assigns to
#152: eleven stat blocks and one section heading. The stat blocks are
Allosaurus (PDF p. 344) and the Ancient Black, Blue, Brass, Bronze, Copper,
Gold, Green, Red, Silver and White Dragons (pp. 265, 267, 269, 271, 277, 292,
294–295, 319, 325 and 341). "Animals" (p. 344) is a section heading with no
stat block and is recorded as an index entry only. The fixture keeps the pinned
English SRD 5.2.1 source, CC-BY-4.0 attribution and adaptation notice.

IDs follow the inventory: `srd-5.2.1:content/monsters/<slug>` for the dragons
and `srd-5.2.1:content/animals/<slug>` for Allosaurus and the heading.
`Monster.Validate` accepts the `content/animals/` prefix for that reason. That
is the only change to the shared #143 contract.

## Layout

| File | Content |
| --- | --- |
| `internal/rules/srd521/testdata/monsters-statblocks-02.json` | The batch data |
| `internal/rules/srd521/monster_statblocks_02.go` | Types, strict decoding and validation |
| `internal/rules/srd521/monster_statblocks_02_test.go` | Literal source fixtures and guards |
| `scripts/verify_srd_statblocks_02.py` | Comparator against the pinned PDF |

Each stat block is a #143 `Monster` plus the typed fields that contract has no
member for: tags, skills, damage and condition immunities, lair XP, Legendary
Action and Legendary Resistance uses, and one detail record per action
(recharge, save, area, damage rolls, conditions, penalties, attack sequences,
Multiattack replacements and spell lists). Action text is kept verbatim,
including the Gold Dragon's "+17 to hit" and the Copper Dragon's "an 90-foot-long"
Line. The Green Dragon's legendary actions continue on p. 295.

## Checks

Validation recomputes ability modifiers, save and skill bonuses, passive
Perception, CR and lair XP, attack bonuses, spell save DCs and dice averages,
and requires a mechanic key for every speed, sense, skill, immunity, trait,
action and effect. The Go test compares each block with literal source values,
checks the twelve identities and pages against the inventory, rejects malformed
and loosely decoded input, and confirms the batch composes with the Bandit
catalogue without duplicate IDs. The Python comparator reads the PDF columns,
rebuilds every record and compares the whole file.

```sh
rtk proxy python3 -I scripts/verify_srd_statblocks_02.py /path/to/SRD_CC_v5.2.1.pdf
rtk proxy env GOMAXPROCS=2 go test -p 1 -short -race ./internal/rules/srd521/...
rtk proxy go vet ./internal/rules/srd521/...
```

This is offline catalogue data. Every block is unavailable: no mechanic key in
this batch is implemented, the Bandit fixture's keys do not unlock any block,
and nothing here adds NPC behavior, encounters, coverage or Full Rules support.
