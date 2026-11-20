# Item source batches

Item catalogue batches add descriptive SRD 5.2.1 item families under the
[item family contract](item-family-contract.md). A batch is data only. Every
imported record is unavailable: `MechanicsAvailable` always returns false, and
no resolver, inventory grant, tool or coverage entry reads a batch. Source pin,
license and attribution: [source provenance](srd-5.2.1-conformance.md).

## Layout

`internal/rules/srd521/item_source.go` is shared by all batches and should not
need editing for a new one. It defines the generic `ItemSourceBatch[D]` record,
strict decoding, manifest validation and `ComposeItemSourceIdentities`.

Each batch owns four files named after its catalogue issue, here `225`:

| File | Content |
| --- | --- |
| `internal/rules/srd521/item_details_225.go` | Typed payloads, `ItemDetails225`, the reviewed manifest, `LoadItems225` and `DecodeItems225` |
| `internal/rules/srd521/testdata/items_225.json` | The batch data |
| `internal/rules/srd521/items_225_test.go` | Identity, literal source-value, composition and rejection tests |
| `scripts/verify_srd_items_225.py` | Comparison with the pinned PDF |

A later batch adds its own four files and a section below. It reuses the shared
leaf types (`ItemDice`, `ItemCheck`) and names its payload types after its items.

## Record

A batch root records `schema_version: 1`, `rules_version: srd-5.2.1`, the
planning `batch` key, the `catalogue_issue` and the pinned `source`. Each family
has a canonical ID (`srd-5.2.1:content/items/<slug>`), name, kind, category, the
base items its category line allows or excludes, and one presentation that links
the planning inventory ID and page.

Each variant has its own stable ID (`<family-id>/variant/<slug>`), name, rarity,
attunement, charges, equip and use restrictions, cost, required mechanic keys
and a `details` object holding exactly one typed payload. A family with one
form uses the variant `standard`. A family whose heading lists several rarities,
such as "+1, +2, or +3", has one variant per rarity. Rows of a random table are
data inside the payload, not variants.

The reviewed manifest lives in code. It pins every family and variant identity,
page, category, base restriction, rarity, attunement, equip and use mode and
mechanic key list, in order. A batch that omits, adds, reorders or widens any of
them fails. Every field is required: the decoded value must re-encode to exactly
the supplied JSON, so an omitted zero, `false` or empty list fails too. Null,
unknown fields, duplicate or mixed-case members and trailing data fail. Payload
strings are lower-case keys, never prose; numbers are non-negative; tables must
cover their die exactly.

The manifest does not pin payload values. Inside `details` the decoder checks
shape and internal consistency only: a table covers its die, a rarity variant
carries its own bonus, a DC or die count is positive. It accepts a
well-formed payload with a different DC, die, distance or key, so
`DecodeItems225` does not prove that caller-supplied JSON matches the source.
The embedded data is pinned by the literal fixtures in
`TestItems225SourceValues`, and the Python comparison checks its values, table
rows and keys against the printed text.

No price is inferred from rarity. `unspecified` with zero amount and `none`
currency means no listed price, never free. A once-per-dawn property is a use
limit in the payload, not a charge pool.

## Magic items 01 (#225), pp. 209–211

Twelve families and sixteen variants, in planning inventory order:

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Adamantine Armor | 209 | `standard` (Uncommon) | `adamantine` |
| Ammunition of Slaying | 209 | `standard` (Very Rare) | `slaying`: 14-row 1d100 creature type table, save and extra damage |
| Ammunition, +1, +2, or +3 | 209 | `plus-1` (Uncommon), `plus-2` (Rare), `plus-3` (Very Rare) | `ammunition_bonus` |
| Amulet of Health | 209 | `standard` (Rare) | `health` |
| Amulet of Proof against Detection and Location | 209 | `standard` (Uncommon) | `proof` |
| Amulet of the Planes | 209 | `standard` (Very Rare) | `planes`: check and 5-row 1d100 destination table with its 1d6 and 1d8 sub-rolls |
| Animated Shield | 209 | `standard` (Very Rare) | `animated` |
| Apparatus of the Crab | 210 | `standard` (Legendary) | `apparatus`: statistics, vehicle rules and the 10-row lever table |
| Armor of Invulnerability | 210 | `standard` (Legendary) | `invulnerability`, including Metal Shell |
| Armor of Resistance | 210 | `standard` (Rare) | `resistance`: 10-row 1d10 damage type table |
| Armor of Vulnerability | 211 | `standard` (Rare) | `vulnerability`, including the curse |
| Armor, +1, +2, or +3 | 210 | `plus-1` (Rare), `plus-2` (Very Rare), `plus-3` (Legendary) | `armor_bonus` |

Amulet of Health, both other amulets, Animated Shield and the Armors of
Invulnerability, Resistance and Vulnerability require attunement. None of the
twelve lists a price or has charges. The source says ten pieces of magic
ammunition are worth one potion of the same rarity; that is recorded as a
relation, not a price.

Composing the Healing contract fixture with this batch yields 13 families and 20
variants with no shared identity. The Healing fixture, `ResolveHealingPotion`
and the coverage register are unchanged.

To recheck the PDF pin, inventory identities and pages, printed category and
rarity lines, every table row and the source values, run:

```sh
rtk proxy python3 -I scripts/verify_srd_items_225.py /path/to/SRD_CC_v5.2.1.pdf
```

The comparison supplements the offline literal fixtures; it is not an automatic
conformance proof. The PDF and extracted text are not vendored.

## Magic items 02 (#226), pp. 211–213

Twelve families and eighteen variants, in planning inventory order:

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Arrow-Catching Shield | 211 | `standard` (Rare) | `arrow_catching` |
| Bag of Beans | 211 | `standard` (Rare) | `beans`: bag, explosion save and damage, and the 12-row 1d100 effect table |
| Bag of Devouring | 212 | `standard` (Very Rare) | `devouring`: pull chance, both checks, capacity and destruction |
| Bag of Holding | 212 | `standard` (Uncommon) | `holding`: dimensions, capacity, air and the Astral gate |
| Bag of Tricks | 212 | `gray`, `rust`, `tan` (all Uncommon) | `tricks`: that colour's 8-row 1d8 creature table |
| Bead of Force | 212 | `standard` (Rare) | `force`: save, damage and the force sphere |
| Bead of Nourishment | 213 | `standard` (Common) | `nourishment` |
| Belt of Dwarvenkind | 213 | `standard` (Rare) | `dwarvenkind`: three benefits, the beard chance and two further benefits |
| Belt of Giant Strength | 213 | `hill` (Rare), `frost-or-stone` (Very Rare), `fire` (Very Rare), `cloud` (Legendary), `storm` (Legendary) | `giant_strength`: giant types and Strength score |
| Berserker Axe | 213 | `standard` (Rare) | `berserker`, including the curse |
| Boots of Elvenkind | 213 | `standard` (Uncommon) | `elvenkind` |
| Boots of Levitation | 213 | `standard` (Rare) | `levitation` |

Variant rules used in this batch:

- Bag of Tricks prints one rarity and three tables, "Gray Bag of Tricks", "Rust
  Bag of Tricks" and "Tan Bag of Tricks". Each is a variant with its own table.
- Belt of Giant Strength prints "Rarity Varies" and a five-row table. Each row
  is a variant, named as printed. The "frost or stone" row stays one variant
  that lists both giant types; splitting it would invent two names.
- Berserker Axe's Battleaxe, Greataxe and Halberd are base item options, like
  the armor types of batch #225, not variants.

Arrow-Catching Shield, Belt of Dwarvenkind, Belt of Giant Strength, Berserker
Axe and Boots of Levitation require attunement. None of the twelve lists a price
or has charges. Bag of Tricks' three uses before the next dawn are a use limit
in the payload, not a charge pool.

Bead of Force starts on p. 212 and continues on p. 213. Its presentation keeps
the planning inventory page, 212, and its payload records `source_pages`
`[212, 213]`. The Bag of Beans table row printed as 96–00 is recorded as 96–100.
Its row numbers are typed quantities: dice, a flat bonus and a unit. Creature
names in a row are keys, not links to monster records.

`DecodeItems226` pins a few literals beyond shape: the five belt rows, the three
Bag of Tricks creature tables, the Bag of Beans DC 15 Dexterity save and the Bag
of Devouring DC 15 and DC 20 checks. Other payload values are pinned only by
`TestItems226SourceValues` and the Python comparison, as in batch #225.

Composing the Healing contract fixture, batch #225 and this batch yields 25
families and 38 variants with no shared identity. The Healing fixture, batch
#225, `ResolveHealingPotion` and the coverage register are unchanged.

To recheck the PDF pin, inventory identities and pages, printed category and
rarity lines, every table row and the source values, run:

```sh
rtk proxy python3 -I scripts/verify_srd_items_226.py /path/to/SRD_CC_v5.2.1.pdf
```

The comparison ties every payload key to its printed phrase and fails when a
payload gains a key it does not map. It supplements the offline literal fixtures
in `TestItems226SourceValues`; it is not an automatic conformance proof.

## Magic items 03 (#227), pp. 214–215

Files: `item_details_227.go` (`ItemDetails227`, `LoadItems227`,
`DecodeItems227`), `testdata/items_227.json`, `items_227_test.go` and
`scripts/verify_srd_items_227.py`. `item_source.go` is unchanged.

Twelve families and twelve variants, in planning inventory order. Every family
is a Wondrous Item with one printed rarity, so each has the single variant
`standard`:

| Family | Page | Rarity | Attunement | Typed payload |
| --- | --- | --- | --- | --- |
| Boots of Speed | 214 | Rare | Required | `boots_speed`: ten-minute total until a Long Rest |
| Boots of Striding and Springing | 214 | Uncommon | Required | `boots_striding` |
| Boots of the Winterlands | 214 | Uncommon | Required | `boots_winterlands`: Cold Resistance and Winter Strider |
| Bowl of Commanding Water Elementals | 214 | Rare | None | `bowl`: summon and the bowl's size |
| Bracers of Archery | 214 | Uncommon | Required | `bracers_archery` |
| Bracers of Defense | 214 | Rare | Required | `bracers_defense` |
| Brazier of Commanding Fire Elementals | 214 | Rare | None | `brazier`: summon |
| Brooch of Shielding | 214 | Uncommon | Required | `brooch` |
| Broom of Flying | 214 | Uncommon | Required | `broom`: riding, sending and recalling |
| Candle of Invocation | 214 | Very Rare | Required | `candle`: burn time, light and the 16-row 1d100 Outer Plane table |
| Cape of the Mountebank | 215 | Rare | None | `cape` |
| Carpet of Flying | 215 | Very Rare | None | `carpet`: 4-row 1d100 size table |

Candle of Invocation starts on p. 214 and its Outer Plane table is on p. 215.
The presentation keeps the inventory page 214; the payload records
`table_page: 215`. The four Carpet of Flying sizes are rows of its random table,
so they are payload data, not variants. A printed `00` is stored as 100.

None of the twelve lists a price or has charges. The Boots of Speed ten-minute
total, the candle's four-hour burn time and the once-until-next-dawn properties
of the bowl, brazier and cape are use limits in the payload, not charge pools.
The Bowl and Brazier share the `ItemBowlBrazierSummon` type because the source
prints the same summoning property for both.

Composing the Healing contract fixture, magic items 01 and 02 and this batch
yields 37 families and 50 variants with no shared identity. The Healing fixture,
`ResolveHealingPotion` and the coverage register are unchanged.

To recheck the PDF pin, inventory identities and pages, printed category and
rarity lines, every table row and the source values, run:

```sh
rtk proxy python3 -I scripts/verify_srd_items_227.py /path/to/SRD_CC_v5.2.1.pdf
```

This work includes material from the System Reference Document 5.2.1 (“SRD
5.2.1”) by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd.
The SRD 5.2.1 is licensed under the Creative Commons Attribution 4.0 International
License, available at https://creativecommons.org/licenses/by/4.0/legalcode.
Names, numeric values, table rows and requirements are adapted into Edda's
structured item records; the software license does not replace this source license.
