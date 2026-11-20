# Stat-block batch 05 (#155)

This batch catalogues the twelve stat blocks the
[#82 inventory](investigations/issue-82-world-inventory.json) assigns to #155.
It is data only. Every record is unavailable: `Monster.Available` implements no
mechanic key, the embedded runtime catalogue still holds only the Bandit, and no
resolver, encounter, tool or coverage entry reads the batch. Source pin, license
and attribution: [source provenance](srd-5.2.1-conformance.md).

## Records

| Record | PDF page | Inventory family |
| --- | --- | --- |
| Black Bear | 345 | animals |
| Black Dragon Wyrmling | 263–264 | monsters |
| Black Pudding | 265–266 | monsters |
| Blink Dog | 266 | monsters |
| Blood Hawk | 345–346 | animals |
| Blue Dragon Wyrmling | 266 | monsters |
| Boar | 346 | animals |
| Bone Devil | 267–268 | monsters |
| Brass Dragon Wyrmling | 268 | monsters |
| Bronze Dragon Wyrmling | 269–270 | monsters |
| Brown Bear | 346 | animals |
| Bugbear Stalker | 271 | monsters |

IDs follow the inventory, so the four animal blocks use
`srd-5.2.1:content/animals/<slug>`. Blocks that continue on a second page keep
their text rejoined across the page break.

## Layout

| File | Content |
| --- | --- |
| `internal/rules/srd521/testdata/monsters_155.json` | The batch data |
| `internal/rules/srd521/monsters_155_test.go` | Literal source values, inventory, dependency-key, rejection and composition tests |
| `scripts/verify_srd_monsters_155.py` | Comparison with the pinned PDF |

The layout and record conventions follow batches #151 to #154. No shared source
file changed: the existing stat-block contract already types every value here.

## Source values the contract cannot type

These stay in the entry's retained source text and are pinned by the tests and
the PDF comparison.

- **Bugbear Stalker, Javelin** is a "Melee or Ranged Attack Roll" with reach
  10 ft. or range 30/120 ft. Its key is `action/melee-or-ranged-attack`, as for
  the #153 Archmage, and its damage dice are typed on its effect.
- **Bugbear Stalker, Morningstar** prints "(with Advantage if the target is
  Grappled by the bugbear)" in its attack roll. It is typed as the key
  `attack/advantage-against-grappled-target`, as for the #153 Ankheg.
- **Blood Hawk, Beak** prints a second damage expression for a Bloodied target.
  It is typed on its own effect, `effect/bloodied-piercing-damage`.
- **Black Pudding** Corrosive Form, Dissolving Pseudopod and Split print
  penalties, a 2-foot eating rate and a trigger with size and Hit Point limits.
  Each stays in text under its own key; only the 1d8 Acid damage is typed.
- **Blink Dog** languages are one printed phrase, kept as a single entry.

## Verification

```sh
GOMAXPROCS=2 go test -p 1 -short -race ./internal/rules/srd521/...
go vet ./internal/rules/srd521/...
python3 -I scripts/verify_srd_monsters_155.py /path/to/SRD_CC_v5.2.1.pdf
```

The verifier checks the PDF's SHA-256, prints each record from its typed fields
and compares it with the PDF text, compares entry headings against the PDF's
bold-italic runs, and ties every typed number and dependency key to a phrase its
own entry prints.
