# Item source batches

Item catalogue batches add descriptive SRD 5.2.1 item families under the
[item family contract](item-family-contract.md). A batch is data only. Every
imported record is unavailable: `MechanicsAvailable` always returns false, and
no resolver, inventory grant, tool or coverage entry reads a batch. Source pin,
license and attribution: [source provenance](srd-5.2.1-conformance.md).

## Potion of Healing

The Healing contract fixture (`item-families.json`, see the
[item family contract](item-family-contract.md)) owns Potion of Healing with
both source presentations: the adventuring-gear entry (p. 99) and the Potions
of Healing magic-item family (p. 236) with its four variants. It is the only
playable item record. A batch whose inventory lists either presentation, such as
#238 or #261, joins that entry to the fixture in its test and verifies the
source rows against it. It does not add a second record, which
`ComposeItemSourceIdentities` would reject as a shared identity.

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

## Magic items 04 (#228), pp. 215–216

Twelve families, each with the single variant `standard`, in planning inventory
order. Files: `item_details_228.go`, `testdata/items_228.json`,
`items_228_test.go` and `scripts/verify_srd_items_228.py`.

| Family | Page | Rarity | Attunement | Typed payload |
| --- | --- | --- | --- | --- |
| Censer of Controlling Air Elementals | 215 | Rare | none | `censer` |
| Chime of Opening | 215 | Rare | none | `chime`: ten uses, then useless |
| Circlet of Blasting | 215 | Uncommon | none | `circlet` |
| Cloak of Arachnida | 215 | Very Rare | required | `arachnida` |
| Cloak of Displacement | 215 | Rare | required | `displacement` |
| Cloak of Elvenkind | 215 | Uncommon | required | `elvenkind` |
| Cloak of Invisibility | 215 | Legendary | required | `invisibility`: continues onto p. 216 |
| Cloak of Protection | 216 | Uncommon | required | `protection` |
| Cloak of the Bat | 216 | Rare | required | `bat` |
| Cloak of the Manta Ray | 216 | Uncommon | required | `manta_ray` |
| Crystal Ball | 216 | Very Rare | required | `crystal_ball` |
| Crystal Ball of Mind Reading | 216 | Legendary | required | `mind_reading` |

The Cloak of Invisibility is the only family here with a charge pool: 3
charges, 1 per use, regaining 1d3 at dawn. A manifest family may now state
`Charges`; without it the decoder still requires no charges. The Chime of
Opening says "used 10 times", not charges, so its ten uses are a use limit in
the payload. Once-per-dawn properties are use limits, as in batch 01. None of
the twelve lists a price.

Composing the Healing contract fixture with batches #225 and #228 yields 25
families and 32 variants with no shared identity. The Healing fixture,
`ResolveHealingPotion` and the coverage register are unchanged.

```sh
rtk proxy python3 -I scripts/verify_srd_items_228.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 05 (#229), pp. 216–218

Twelve families, each with the single variant `standard`, in planning inventory
order. Files: `item_details_229.go`, `testdata/items_229.json`,
`items_229_test.go` and `scripts/verify_srd_items_229.py`.

| Family | Page | Rarity | Attunement | Typed payload |
| --- | --- | --- | --- | --- |
| Crystal Ball of Telepathy | 216 | Legendary | required | `telepathy` |
| Crystal Ball of True Seeing | 216 | Legendary | required | `true_seeing` |
| Cube of Force | 216 | Rare | required | `cube_of_force`: six-face charge table |
| Cubic Gate | 216 | Legendary | none | `cubic_gate` |
| Dagger of Venom | 216 | Rare | none | `venom`: base `dagger` |
| Dancing Sword | 217 | Very Rare | required | `dancing`: five base swords |
| Decanter of Endless Water | 217 | Uncommon | none | `decanter` |
| Deck of Illusions | 217 | Uncommon | none | `illusions`: 33-row 1d100 table, p. 217 into p. 218 |
| Defender | 218 | Legendary | required | `defender`: base any melee weapon |
| Demon Armor | 218 | Very Rare | required | `demon_armor`: base any light, medium or heavy armor |
| Dimensional Shackles | 218 | Rare | none | `shackles` |
| Dragon Orb | 218 | Artifact | required | `dragon_orb`: five-row spell table |

Three families have a charge pool. The Cubic Gate has 3 charges, 1 per use,
regaining 1d3 at dawn. The Cube of Force (10 charges, regains 1d6) and the
Dragon Orb (7 charges, regains 1d4 + 3) charge a different amount for each
spell, so the variant states `per_use: 0` and each cost is a row in the
payload table. Zero in a payload spell row means the source gives no spell
level or save DC. The Deck of Illusions table lists 32 creature rows of three
numbers each and the card drawer at 97–00, recorded as 97–100. The deck
subtracts one from 1d20 for missing cards, stored as `missing_cards_subtract`
because payload numbers are never negative. None of the twelve lists a price.

Composing the Healing contract fixture with batches #225 to #229 yields 61
families and 74 variants with no shared identity. The Healing fixture,
`ResolveHealingPotion` and the coverage register are unchanged.

```sh
rtk proxy python3 -I scripts/verify_srd_items_229.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 06 (#230), pp. 218–220

Twelve families in planning inventory order. Files: `item_details_230.go`,
`testdata/items_230.json`, `items_230_test.go` and
`scripts/verify_srd_items_230.py`. Dragon Scale Mail has ten variants, one per row
of the source's dragon table, and Elemental Gem has four, one per row of its gem
table; every other family has the single variant `standard`. Presentation pages
follow the page of the heading, as in batch 04.

| Family | Page | Rarity | Attunement | Typed payload |
| --- | --- | --- | --- | --- |
| Dragon Scale Mail | 218 | Very Rare | required | `scale_mail`: ten dragon variants; continues onto p. 219 |
| Dragon Slayer | 219 | Rare | none | `slayer` |
| Dust of Disappearance | 219 | Uncommon | none | `disappearance` |
| Dust of Dryness | 219 | Uncommon | none | `dryness`: continues across the column break |
| Dust of Sneezing and Choking | 219 | Uncommon | none | `sneezing` |
| Dwarven Plate | 219 | Very Rare | none | `plate` |
| Dwarven Thrower | 219 | Very Rare | required | `thrower`: the Dwarf or Belt of Dwarvenkind qualifier is payload data |
| Efficient Quiver | 219 | Uncommon | none | `quiver` |
| Efreeti Bottle | 220 | Very Rare | none | `bottle`: the 1d10 table, rolled on the first opening |
| Elemental Gem | 220 | Uncommon | none | `gem`: four gem variants |
| Elixir of Health | 220 | Rare | none | `elixir` |
| Elven Chain | 220 | Rare | none | `chain` |

None of the twelve lists a price or a charge pool. The Dust of Dryness pinches
are a quantity (1d6 + 4), not charges. Attunement restrictions stay empty in the
shared fields, so the Dwarven Thrower's qualifier is recorded in its payload. The
shared manifest and validator are unchanged.

The shared manifest stores one presentation page per family, so Dragon Scale
Mail records p. 218 although its text and table continue onto p. 219. In the
Efreeti Bottle table, `understands_owner_languages` and `obeys_owner_commands`
are `false` where the source does not state them (rows 1 and 10); `false` there
means not stated, not a source claim. The validator checks structure and the
dragon and gem tables; the literal fixture and the PDF verifier pin the numeric
values.

Composing the Healing contract fixture with batches #225 to #229 and #230 yields
73 families and 98 variants with no shared identity. The Healing fixture,
`ResolveHealingPotion` and the coverage register are unchanged.

```sh
rtk proxy python3 -I scripts/verify_srd_items_230.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 07 (#231), pp. 220–223

Twelve families and 25 variants, in planning inventory order. Files:
`item_details_231.go`, `testdata/items_231.json`, `items_231_test.go` and
`scripts/verify_srd_items_231.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Energy Bow | 220 | `standard` (Very Rare, attunement) | `energy_bow`: three properties |
| Eversmoking Bottle | 220 | `standard` (Uncommon) | `eversmoking_bottle`: continues onto p. 221 |
| Eyes of Charming | 221 | `standard` (Uncommon, attunement) | `eyes_of_charming`: 3 charges |
| Eyes of Minute Seeing | 221 | `standard` (Uncommon) | `eyes_of_minute_seeing` |
| Eyes of the Eagle | 221 | `standard` (Uncommon) | `eyes_of_the_eagle` |
| Feather Token | 221 | `anchor`, `bird`, `fan`, `swan-boat`, `tree`, `whip` | `feather_token`: its 1d100 row and numeric facts |
| Figurine of Wondrous Power | 221 | nine printed figurines, `bronze-griffon` to `silver-raven` | `figurine`: shared rules and its forms; continues onto p. 222 |
| Flame Tongue | 223 | `standard` (Rare, attunement) | `flame_tongue` |
| Folding Boat | 223 | `standard` (Rare) | `folding_boat`: three command words |
| Frost Brand | 223 | `standard` (Very Rare, attunement) | `frost_brand` |
| Gauntlets of Ogre Power | 223 | `standard` (Uncommon, attunement) | `gauntlets_of_ogre_power` |
| Gem of Brightness | 223 | `standard` (Uncommon) | `gem_of_brightness`: 50 charges, three command words |

Variant rules used in this batch:

- Feather Token and Figurine of Wondrous Power print "Rarity Varies". Each printed
  kind is a variant named as printed, with the rarity printed after its name.
  The six rows of the Feather Tokens table are payload data (`table_min`,
  `table_max`; a printed `00` is stored as 100), not separate records.
- All variants of a family share one payload type, so per-kind numbers are
  `ItemFact231` entries: a key, a non-negative amount and a unit. The unit
  `minus` marks an amount printed with a minus sign (the Giant Fly's modifiers
  and saves for Intelligence and Charisma).
- Ivory Goats is one variant with three forms (Goat of Terror, Goat of
  Traveling, Goat of Travail), because the source prints one heading and one
  rarity. The Goat of Traveling has no hour limit; its `duration_hours` is 0 and
  its 24 charges are facts of that form, not the variant's charge pool.
- The Ebony Fly's accompanying Giant Fly stat block is recorded as facts of that
  form. Creature names are keys, not links to monster records.
- The Energy Bow's, Flame Tongue's and Frost Brand's base weapons are base item
  options, not variants.

Eyes of Charming (3 charges, regains all at dawn) and the Gem of Brightness (50
charges, no recharge) are the only families with a charge pool. Their
`per_use` is 1, the least a use costs; the Eyes accept more charges and the
gem's command words cost 0, 1 and 5, all recorded in the payloads. The 50 GP
value of the depleted gem is a payload fact, not a price: every cost stays
unspecified. The Frost Brand's once-per-hour property is a use limit, not a
charge pool.

Composing the Healing contract fixture with batches #225 to #230 and this batch yields 85
families and 123 variants with no shared identity. The Healing fixture,
`ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads pp. 220–223 as one text stream, so the Eversmoking
Bottle and the figurine text continue across their page breaks.

```sh
rtk proxy python3 -I scripts/verify_srd_items_231.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 08 (#232), pp. 223–225

Twelve families and 12 variants, in planning inventory order. Files:
`item_details_232.go`, `testdata/items_232.json`, `items_232_test.go` and
`scripts/verify_srd_items_232.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Gem of Seeing | 223 | `standard` (Rare, attunement) | `gem_of_seeing`: 3 charges, regains 1d3 at dawn |
| Giant Slayer | 223 | `standard` (Rare) | `giant_slayer` |
| Glamoured Studded Leather | 223 | `standard` (Rare) | `glamoured_studded_leather` |
| Gloves of Missile Snaring | 224 | `standard` (Uncommon, attunement) | `gloves_of_missile_snaring` |
| Gloves of Swimming and Climbing | 224 | `standard` (Uncommon, attunement) | `gloves_of_swimming_and_climbing` |
| Gloves of Thievery | 224 | `standard` (Uncommon) | `gloves_of_thievery` |
| Goggles of Night | 224 | `standard` (Uncommon) | `goggles_of_night` |
| Hammer of Thunderbolts | 224 | `standard` (Legendary, attunement) | `hammer_of_thunderbolts`: 5 charges, regains 1d4 + 1 at dawn |
| Handy Haversack | 224 | `standard` (Rare) | `handy_haversack`: pouches and the nested-space gate |
| Hat of Disguise | 224 | `standard` (Uncommon, attunement) | `hat_of_disguise` |
| Hat of Many Spells | 224 | `standard` (Very Rare, attunement by a Wizard) | `hat_of_many_spells`: continues onto p. 225 |
| Headband of Intellect | 225 | `standard` (Uncommon, attunement) | `headband_of_intellect` |

Variant rules used in this batch:

- No family in this batch prints several kinds or "Rarity Varies", so each has one
  `standard` variant named as printed.
- The Hat of Many Spells is attuned by a Wizard. The shared attunement field has
  no class slot and stays `required` with no restrictions; the class is the
  payload fact `attunement_class` and the mechanic key `attunement.wizard`.
- The Hat's failure table is payload data: eleven rows (`min`, `max`; a printed
  `00` is stored as 100) whose further rolls (d10, d4 and d6) are `outcomes`
  that must cover their die exactly once. The Hat is listed on p. 224; its text
  and table continue on p. 225, recorded as `source_pages` [224, 225].
- Gem of Seeing (3 charges) and Hammer of Thunderbolts (5 charges) are the only
  families with a charge pool. Their `per_use` is 1. The Hammer's Giant's Bane
  benefits are payload data of the Hammer, not separate records. The 50 GP gem and
  1,000 GP component limit in the Hat's rows are payload facts, not prices: every
  cost stays unspecified.
- Belt of Giant Strength, Gauntlets of Ogre Power, Bag of Holding and Portable Hole
  appear only as keys; they are not links to item records.

Composing the Healing contract fixture with batches #225 to #231 and this batch yields 97
families and 135 variants with no shared identity. The Healing fixture,
`ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads pp. 223–225 as one text stream, so the Hat of Many
Spells text and table continue across the p. 224 to p. 225 break.

```sh
rtk proxy python3 -I scripts/verify_srd_items_232.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 09 (#233), pp. 225–227

Twelve families and 28 variants, in planning inventory order. Files:
`item_details_233.go`, `testdata/items_233.json`, `items_233_test.go` and
`scripts/verify_srd_items_233.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Helm of Brilliance | 225 | `standard` (Very Rare, attunement) | `helm_of_brilliance`: gems, four properties and the fire-damage mishap; continues onto p. 226 |
| Helm of Comprehending Languages | 226 | `standard` (Uncommon) | `helm_of_comprehending_languages` |
| Helm of Telepathy | 226 | `standard` (Uncommon, attunement) | `helm_of_telepathy` |
| Helm of Teleportation | 226 | `standard` (Rare, attunement) | `helm_of_teleportation`: 3 charges |
| Holy Avenger | 226 | `standard` (Legendary, attunement) | `holy_avenger`: class attunement and the aura |
| Horn of Blasting | 226 | `standard` (Rare) | `horn_of_blasting` |
| Horn of Valhalla | 226 | `silver`, `brass`, `bronze`, `iron` | `horn_of_valhalla`: shared rules and its 1d100 row |
| Horseshoes of Speed | 227 | `standard` (Rare) | `horseshoes_of_speed` |
| Horseshoes of a Zephyr | 226 | `standard` (Very Rare) | `horseshoes_of_a_zephyr` |
| Immovable Rod | 227 | `standard` (Uncommon) | `immovable_rod` |
| Instant Fortress | 227 | `standard` (Rare, attunement) | `instant_fortress`: statuette and tower |
| Ioun Stone | 227 | fourteen printed stones, `absorption` to `sustenance` | `ioun_stone`: shared rules and one stone; continues onto p. 228 |

Variant rules used in this batch:

- The Horn of Valhalla prints several rarities and the Ioun Stone prints
  "Rarity Varies". Each printed metal or stone is a variant named as printed,
  with the rarity printed with its heading: Silver and Brass are Rare, Bronze
  is Very Rare and Iron is Legendary. The four rows of the horn table are payload
  data (`table_min`, `table_max`; a printed `00` is stored as 100).
- The Ioun Stone text continues from p. 227 onto p. 228. The presentation keeps
  the inventory page 227; each stone's `source_pages` is the page that prints it
  (Mastery, Protection, Regeneration, Reserve, Strength and Sustenance are on
  p. 228). Shared rules repeat in every stone's payload. Per-stone numbers are
  `ItemFact233` entries: a key, a non-negative amount and a unit.
- The Holy Avenger's "Requires Attunement by a Paladin" is a payload fact
  (`attunement_by_any_of`); the shared validator keeps attunement restrictions
  empty. Its base weapons are any Simple or Martial weapon.
- The Helm of Brilliance text continues from p. 225 onto p. 226, so its
  `source_pages` is `[225, 226]`. Gem counts are dice, not quantities.
- The Helm of Teleportation is the only family with a charge pool (3 charges,
  1 per use, regains 1d3 at dawn). Other limits, such as the Helm of Telepathy's
  once-until-dawn spells and the Horn of Valhalla's 7 days, are use limits.
- Horseshoes of Speed is printed after Horseshoes of a Zephyr; the catalogue
  keeps inventory order and the comparison script handles the print order.

Composing the Healing contract fixture with batches #225 to #231 and this batch
yields 97 families and 151 variants with no shared identity. The Healing
fixture, `ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads pp. 225–228 as one text stream, so the Helm of
Brilliance and the Ioun Stone text continue across their page breaks.

```sh
rtk proxy python3 -I scripts/verify_srd_items_233.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 10 (#234), pp. 228–229

Twelve families and 12 variants, in planning inventory order. Files:
`item_details_234.go`, `testdata/items_234.json`, `items_234_test.go` and
`scripts/verify_srd_items_234.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Iron Bands | 228 | `standard` (Rare) | `iron_bands` |
| Iron Flask | 228 | `standard` (Legendary) | `iron_flask` |
| Javelin of Lightning | 228 | `standard` (Uncommon) | `javelin_of_lightning`: Lightning Bolt property |
| Lantern of Revealing | 228 | `standard` (Uncommon) | `lantern_of_revealing` |
| Luck Blade | 229 | `standard` (Legendary, attunement) | `luck_blade`: Luck and Wish properties |
| Mace of Disruption | 229 | `standard` (Rare, attunement) | `mace_of_disruption` |
| Mace of Smiting | 229 | `standard` (Rare) | `mace_of_smiting` |
| Mace of Terror | 229 | `standard` (Rare, attunement) | `mace_of_terror`: 3 charges |
| Mantle of Spell Resistance | 229 | `standard` (Rare, attunement) | `mantle_of_spell_resistance` |
| Manual of Bodily Health | 229 | `standard` (Very Rare) | `manual_of_bodily_health` |
| Manual of Gainful Exercise | 229 | `standard` (Very Rare) | `manual_of_gainful_exercise` |
| Manual of Golems | 229 | `standard` (Very Rare) | `manual_of_golems`: 1d20 table on p. 230 |

Variant rules used in this batch:

- No family prints more than one variant. The four golems of the Manual of
  Golems are rows of its 1d20 table, kept as payload data (`golems`) with their
  roll range, days and supplies cost; the book is one item, one rarity.
- The Manual of Golems table is printed at the top of p. 230, so its payload
  records `source_pages` 229 and 230 while its presentation page stays 229. The
  GP amounts are supplies the creator pays, not a price: every cost stays
  unspecified.
- The Luck Blade, Mace of Disruption, Mace of Smiting and Mace of Terror base
  weapons (Glaive to Shortsword, Mace) are base item options, not variants.
- The two ability manuals share one payload type, `ItemManualStudy`, under
  their own payload names. "A century" is stored as 100 years.

The Mace of Terror (3 charges, regains 1d3 at dawn) and the Luck Blade are the
only families with a charge pool. The Luck Blade "has 1d3 charges" with no
printed maximum and no recharge, so its pool records maximum 3, the top of the
starting roll, per use 1 and recharge none; the roll itself is the payload's
`initial_charges`. Once-until-dawn limits (Iron Bands, Javelin of Lightning,
Luck's reroll, Wish) are use limits, not charge pools.

Composing the Healing contract fixture with batches #225 to #231 and this batch
yields 97 families and 135 variants with no shared identity. The Healing fixture,
`ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads pp. 228–229 and the left column of p. 230 as one
text stream, so the Lantern of Revealing, the Mace of Terror and the Manual of
Golems continue across their column and page breaks.

```sh
rtk proxy python3 -I scripts/verify_srd_items_234.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 11 (#235), pp. 230–233

Twelve families and twelve variants, in planning inventory order. Files:
`item_details_235.go`, `testdata/items_235.json`, `items_235_test.go` and
`scripts/verify_srd_items_235.py`. `item_source.go` is unchanged.

| Family | Page | Rarity | Attunement | Typed payload |
| --- | --- | --- | --- | --- |
| Manual of Quickness of Action | 230 | Very Rare | none | `manual_of_quickness_of_action` |
| Marvelous Pigments | 230 | Very Rare | none | `marvelous_pigments` |
| Medallion of Thoughts | 230 | Uncommon | required | `medallion_of_thoughts`: 5 charges |
| Mirror of Life Trapping | 230 | Very Rare | none | `mirror_of_life_trapping` |
| Mithral Armor | 231 | Uncommon | none | `mithral_armor`: base any medium or heavy armor, excluding hide |
| Mysterious Deck | 231 | Legendary | none | `mysterious_deck`: both 1d100 columns, 22 card effects, Avatar of Death; continues onto p. 232 |
| Necklace of Adaptation | 232 | Uncommon | required | `necklace_of_adaptation` |
| Necklace of Fireballs | 233 | Rare | none | `necklace_of_fireballs` |
| Necklace of Prayer Beads | 233 | Rare | required | `necklace_of_prayer_beads`: 6-row 1d20 bead table |
| Nine Lives Stealer | 233 | Very Rare | required | `nine_lives_stealer`: base any simple or martial weapon |
| Oathbow | 233 | Very Rare | required | `oathbow`: base longbow or shortbow |
| Oil of Etherealness | 233 | Rare | none | `oil_of_etherealness` |

Every family prints one form, so each has the single variant `standard`. The
Mysterious Deck is one variant: its 13-card and 22-card decks are two columns of
one table, recorded as two row lists (`deck_13_cards`, `deck_22_cards`) that each
cover the 1d100 exactly. A printed `00` is stored as 100, and a card with a dash
in a column has no row in it. The 22 card effects are keys and numeric facts;
the Euryale penalty is the amount 2 with the unit `minus`, and dice such as the
Moon card's 1d3 are an amount with a unit such as `d3`. The Avatar of Death stat
block on p. 232 is a payload of the deck, not a monster record.

The Necklace of Prayer Beads restricts attunement to a Cleric, Druid or Paladin.
The shared decoder requires an empty attunement restriction list, so the classes
are a payload field (`attunement_by_any_of`), as in batch 06 (#230). Its table gives a
spell level only for the Bead of Curing (Cure Wounds, level 2); zero means none
printed.

Only two families have a charge pool. The Medallion of Thoughts has 5 charges,
1 per use, regaining 1d4 at dawn. The Nine Lives Stealer has 1d8 + 1 charges,
one lost per creature slain and never recharged; its variant records 9, the
highest result of the printed dice, and the payload records the dice. The beads
of the two necklaces and the pots of Marvelous Pigments are quantities rolled
when found, not charges. The 25 GP, 500 GP, 2,000 GP and 1,000 GP values are
payload facts, not prices; every cost stays unspecified.

Composing the Healing contract fixture with batches #225 to #231 and this batch
yields 97 families and 135 variants with no shared identity. The Healing
fixture, `ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads pp. 230–233 as one text stream, so the Mysterious
Deck continues across its page break.

```sh
rtk proxy python3 -I scripts/verify_srd_items_235.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 12 (#236), pp. 233–235

Twelve families and 12 variants, in planning inventory order. Files:
`item_details_236.go`, `testdata/items_236.json`, `items_236_test.go` and
`scripts/verify_srd_items_236.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Oil of Sharpness | 233 | `standard` (Very Rare) | `oil_of_sharpness`: +3 weapon and ammunition |
| Oil of Slipperiness | 233 | `standard` (Uncommon) | `oil_of_slipperiness`: Freedom of Movement and Grease |
| Pearl of Power | 234 | `standard` (Uncommon, attunement) | `pearl_of_power` |
| Periapt of Health | 234 | `standard` (Uncommon, attunement) | `periapt_of_health` |
| Periapt of Proof against Poison | 234 | `standard` (Rare, attunement) | `periapt_of_proof_against_poison` |
| Periapt of Wound Closure | 234 | `standard` (Uncommon, attunement) | `periapt_of_wound_closure`: Life Preservation and Natural Healing Boost |
| Philter of Love | 234 | `standard` (Uncommon) | `philter_of_love` |
| Pipes of Haunting | 234 | `standard` (Uncommon) | `pipes_of_haunting`: 3 charges |
| Pipes of the Sewers | 234 | `standard` (Uncommon, attunement) | `pipes_of_the_sewers`: 3 charges |
| Plate Armor of Etherealness | 234 | `standard` (Legendary, attunement) | `plate_armor_of_etherealness` |
| Portable Hole | 234 | `standard` (Rare) | `portable_hole` |
| Potion of Animal Friendship | 235 | `standard` (Uncommon) | `potion_of_animal_friendship`: level 3, save DC 13 |

Variant rules used in this batch:

- No family prints more than one variant, so each has the single variant
  `standard`. Potion of Giant Strength, which has a variant table, is outside
  this batch.
- The Oil of Slipperiness starts on p. 233 and ends on p. 234; its payload
  records `source_pages` 233 and 234. The Portable Hole starts on p. 234 and
  ends on p. 235; its payload records `source_pages` 234 and 235. Both keep
  their heading page as the presentation page.
- The Plate Armor of Etherealness base armor options (Half Plate Armor, Plate
  Armor) are base item options, not variants.
- The Pearl of Power's "Requires Attunement by a Spellcaster" qualifier is the
  payload's `attunement_by_any_of`; the shared attunement restriction list stays
  empty, as in earlier batches.

Pipes of Haunting and Pipes of the Sewers are the only families with a charge
pool: 3 charges, regaining 1d3 at dawn. Their `per_use` is 1, the least a use
costs; the Pipes of the Sewers accept 1 to 3 charges, recorded in the payload.
Once-until-dawn limits (Pearl of Power, Periapt of Health, Plate Armor of
Etherealness) are use limits, not charge pools. Potion and oil quantities are
not charge pools. No family lists a price, so every cost stays unspecified.

Composing the Healing contract fixture with batches #225 to #231 and this batch
yields 97 families and 135 variants with no shared identity. The Healing fixture,
`ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads pp. 233–235 as one text stream, column by column, so
the Oil of Slipperiness, the Pipes of Haunting and the Portable Hole continue
across their column and page breaks.

```sh
rtk proxy python3 -I scripts/verify_srd_items_236.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 13 (#237), pp. 235–236

Twelve families and 16 variants, in planning inventory order. Files:
`item_details_237.go`, `testdata/items_237.json`, `items_237_test.go` and
`scripts/verify_srd_items_237.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Potion of Clairvoyance | 235 | `standard` (Rare) | `clairvoyance` |
| Potion of Climbing | 235 | `standard` (Common) | `climbing` |
| Potion of Diminution | 235 | `standard` (Rare) | `diminution` |
| Potion of Flying | 235 | `standard` (Very Rare) | `flying` |
| Potion of Gaseous Form | 235 | `standard` (Rare) | `gaseous_form` |
| Potion of Giant Strength | 235 | `hill`, `frost-or-stone`, `fire`, `cloud`, `storm` | `giant_strength`: its table row |
| Potion of Growth | 235 | `standard` (Uncommon) | `growth` |
| Potion of Heroism | 236 | `standard` (Rare) | `heroism` |
| Potion of Invisibility | 236 | `standard` (Rare) | `invisibility` |
| Potion of Invulnerability | 236 | `standard` (Rare) | `invulnerability` |
| Potion of Longevity | 236 | `standard` (Very Rare) | `longevity` |
| Potion of Mind Reading | 236 | `standard` (Rare) | `mind_reading` |

Variant rules used in this batch:

- Potion of Giant Strength prints "Rarity Varies" and a five-row table. Each row is
  a variant named as printed, with the rarity printed in its row (Uncommon, Rare,
  Rare, Very Rare, Legendary). The row printed "frost or stone" lists two giants
  and is one variant; its payload holds both in `giants`.
- Every potion is category `potion`, uses the `consume` mode with the `drink-one`
  restriction (as the Elixir of Health in batch 06), and has no attunement, equip
  restriction, charges or price. The spells and conditions a potion names
  (Clairvoyance, Enlarge/Reduce, Gaseous Form, Bless, Detect Thoughts, Invisible)
  are keys, not links to spell or condition records.
- Heroism, Invisibility, Invulnerability and Longevity are on p. 236, after the
  Potions of Healing, which belong to the Healing contract fixture and are not
  part of this batch.

Composing the Healing contract fixture with batches #225 to #231 and this batch yields 97
families and 139 variants with no shared identity. The Healing fixture,
`ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads pp. 235–236 column by column as one text stream.

```sh
rtk proxy python3 -I scripts/verify_srd_items_237.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 14 (#238), pp. 236–238

Eleven families and 14 variants, in planning inventory order. Files:
`item_details_238.go`, `testdata/items_238.json`, `items_238_test.go` and
`scripts/verify_srd_items_238.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Potion of Poison | 236 | `standard` (Uncommon) | `poison` |
| Potion of Resistance | 236 | `standard` (Uncommon) | `resistance_potion`: the 10-row 1d10 damage type table |
| Potion of Speed | 236 | `standard` (Very Rare) | `speed` |
| Potion of Vitality | 236 | `standard` (Very Rare) | `vitality` |
| Potion of Water Breathing | 236 | `standard` (Uncommon) | `water_breathing` |
| Quarterstaff of the Acrobat | 237 | `standard` (Very Rare, attunement) | `acrobat` |
| Ring of Animal Influence | 237 | `standard` (Rare) | `animal_influence` |
| Ring of Djinni Summoning | 237 | `standard` (Legendary, attunement) | `djinni_summoning` |
| Ring of Elemental Command | 237 | `air`, `earth`, `fire`, `water` (all Legendary, attunement) | `elemental_command`: the shared properties and one linked plane's Elemental Focus and spells |
| Ring of Evasion | 238 | `standard` (Rare, attunement) | `evasion` |
| Ring of Feather Falling | 238 | `standard` (Rare, attunement) | `feather_falling` |

Variant rules used in this batch:

- **Potions of Healing is not a family of this batch.** It is the twelfth source
  family of #238, and the [Healing contract fixture](item-family-contract.md)
  already holds it as the canonical `potion-of-healing` family: the
  `magic-items/potions-of-healing` presentation and all four table variants.
  A second record would repeat that presentation, which composition rejects, so
  this batch holds the other eleven families. `TestItems238SourceIdentities`
  joins all twelve inventory records, and `verify_srd_items_238.py` compares the
  fixture's four rows (rarity and healing dice) with the p. 236 table. The
  fixture is unchanged.
- The Ring of Elemental Command is linked to one of four planes and the source
  names a linked ring that way, "Ring of Elemental Command (air)". `air`,
  `earth`, `fire` and `water` are variants, each Legendary, with the shared
  properties repeated in each payload. A plane with no printed value stores
  `none`, such as the Earth ring's speed.
- The Potion of Resistance table is a random table, so it is payload data in one
  `standard` variant, as for the Armor of Resistance in batch 01.
- The Potion of Water Breathing (pp. 236–237) and the Ring of Elemental Command
  (pp. 237–238) record `source_pages` in their payloads.
- The "until the start your next turn" wording of Elemental Compulsion is
  printed that way; the key records the intended phrase.

The Ring of Animal Influence (3 charges, 1d3 at dawn), the Ring of Evasion (3
charges, 1d3 at dawn) and the Ring of Elemental Command (5 charges, 1d4 + 1 at
dawn; Feather Fall costs 0) have charge pools. Their `per_use` is 1, the least
a charged use costs. The Quarterstaff of the Acrobat's deflection and the
Djinni's 24-hour wait are use limits. No family lists a price.

Composing the Healing contract fixture with batches #225 to #237 and this batch
yields 168 families and 229 variants with no shared identity. The Healing
fixture, `ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads the pages column by column as one text stream. It
lists every payload leaf with its printed phrase, fails on a leaf it does not
list, requires a number to appear in its phrase, and derives each family's
pages from the layout. It supplements the literal fixtures in
`TestItems238SourceValues`; it is not an automatic conformance proof.

```sh
rtk proxy python3 -I scripts/verify_srd_items_238.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 15 (#239), pp. 238–239

Twelve families and 12 variants, in planning inventory order. Files:
`item_details_239.go`, `testdata/items_239.json`, `items_239_test.go` and
`scripts/verify_srd_items_239.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Ring of Free Action | 238 | `standard` (Rare, attunement) | `free_action` |
| Ring of Invisibility | 238 | `standard` (Legendary, attunement) | `invisibility` |
| Ring of Jumping | 238 | `standard` (Uncommon, attunement) | `jumping` |
| Ring of Mind Shielding | 238 | `standard` (Uncommon, attunement) | `mind_shielding` |
| Ring of Protection | 238 | `standard` (Rare, attunement) | `protection` |
| Ring of Regeneration | 238 | `standard` (Very Rare, attunement) | `regeneration` |
| Ring of Resistance | 238 | `standard` (Rare) | `resistance_ring`: the 10-row 1d10 damage type and gemstone table |
| Ring of Shooting Stars | 238 | `standard` (Very Rare, attunement) | `shooting_stars`: free spells, three charged properties and the 4-row sphere damage table |
| Ring of Spell Storing | 239 | `standard` (Rare, attunement) | `spell_storing` |
| Ring of Spell Turning | 239 | `standard` (Legendary, attunement) | `spell_turning` |
| Ring of Swimming | 239 | `standard` (Uncommon) | `swimming` |
| Ring of Telekinesis | 239 | `standard` (Very Rare, attunement) | `telekinesis` |

Variant rules used in this batch:

- Every family has the single variant `standard`.
- The Ring of Resistance table is a random table, so its ten damage types and
  gemstones are payload data, as for the Armor of Resistance in batch 01.
- The Ring of Shooting Stars starts on p. 238 and ends on p. 239; its payload
  records `source_pages` 238 and 239. Its sphere damage table maps the number of
  spheres to dice.
- The Ring of Spell Storing's found contents, `1d6 − 1` levels, are the dice and
  a `found_levels_penalty` of 1. Stored levels are not charges.

The Ring of Shooting Stars is the only family with a charge pool: 6 charges,
regaining 1d6 at dawn. Its `per_use` is 1, the least a charged use costs; each
property records its own cost. No family lists a price.

Composing the Healing contract fixture with batches #225 to #238 and this batch
yields 180 families and 241 variants with no shared identity. The Healing
fixture, `ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads the pages column by column as one text stream. It
lists every payload leaf with its printed phrase, fails on a leaf it does not
list, requires a number to appear in its phrase, and derives each family's
pages from the layout. It supplements the literal fixtures in
`TestItems239SourceValues`; it is not an automatic conformance proof.

```sh
rtk proxy python3 -I scripts/verify_srd_items_239.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 16 (#240), pp. 239–241

Twelve families and 12 variants, in planning inventory order. Files:
`item_details_240.go`, `testdata/items_240.json`, `items_240_test.go` and
`scripts/verify_srd_items_240.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Ring of Three Wishes | 239 | `standard` (Legendary) | `three_wishes` |
| Ring of Warmth | 239 | `standard` (Uncommon, attunement) | `warmth` |
| Ring of Water Walking | 239 | `standard` (Uncommon) | `water_walking` |
| Ring of X-ray Vision | 239 | `standard` (Rare, attunement) | `x_ray` |
| Ring of the Ram | 239 | `standard` (Rare, attunement) | `ram` |
| Robe of Eyes | 240 | `standard` (Rare, attunement) | `eyes` |
| Robe of Scintillating Colors | 240 | `standard` (Very Rare, attunement) | `scintillating` |
| Robe of Stars | 240 | `standard` (Very Rare, attunement) | `stars` |
| Robe of Useful Items | 240 | `standard` (Uncommon) | `useful_items`: six fixed patches, the 13-row 1d100 table and its numeric facts |
| Robe of the Archmagi | 240 | `standard` (Legendary, attunement) | `archmagi` |
| Rod of Absorption | 241 | `standard` (Very Rare, attunement) | `absorption` |
| Rod of Alertness | 241 | `standard` (Very Rare, attunement) | `alertness` |

Variant rules used in this batch:

- Every family has the single variant `standard`.
- The Ring of X-ray Vision (pp. 239–240) and the Robe of Useful Items (pp.
  240–241) record `source_pages` in their payloads. The Ring of the Ram is
  printed before the other four rings; the catalogue keeps inventory order, which
  also places the Robe of Useful Items before the Robe of the Archmagi.
- The Robe of Useful Items patches are rows of a random table, so they are
  payload data. A printed `00` is stored as 100. Each row records its patch; only
  printed counts, sizes and values are `patch_facts`, and the iron door's and
  window's placement rules are `patch_rules`. Names such as Riding Horse,
  Mastiff and Spell Scroll are keys, not links to other records.
- The Robe of the Archmagi's attunement classes are the payload's
  `attunement_by_any_of`; the shared restriction list stays empty.

The Ring of Three Wishes (3 charges, never recharged), the Ring of the Ram (3
charges, 1d3 at dawn) and the Robe of Scintillating Colors (3 charges, 1d3 at
dawn) have charge pools. The Robe of Stars' stars, the Rod of Absorption's
stored levels and the patches are not charges. No family lists a price; the
patch values are payload facts.

Composing the Healing contract fixture with batches #225 to #239 and this batch
yields 192 families and 253 variants with no shared identity. The Healing
fixture, `ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads the pages column by column as one text stream. It
lists every payload leaf with its printed phrase, fails on a leaf it does not
list, requires a number to appear in its phrase, and derives each family's
pages from the layout. It supplements the literal fixtures in
`TestItems240SourceValues`; it is not an automatic conformance proof.

```sh
rtk proxy python3 -I scripts/verify_srd_items_240.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 17 (#241), pp. 241–243

Twelve families and 12 variants, in planning inventory order. Files:
`item_details_241.go`, `testdata/items_241.json`, `items_241_test.go` and
`scripts/verify_srd_items_241.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Rod of Lordly Might | 241 | `standard` (Legendary, attunement) | `lordly_might`: the mace, six button forms and the Drain Life, Paralyze and Terrify properties |
| Rod of Resurrection | 242 | `standard` (Legendary, attunement) | `resurrection` |
| Rod of Rulership | 242 | `standard` (Rare, attunement) | `rulership` |
| Rod of Security | 242 | `standard` (Very Rare) | `security` |
| Rope of Climbing | 242 | `standard` (Uncommon) | `climbing_rope` |
| Rope of Entanglement | 242 | `standard` (Rare) | `entanglement` |
| Scarab of Protection | 243 | `standard` (Legendary, attunement) | `scarab` |
| Scimitar of Speed | 243 | `standard` (Very Rare, attunement) | `speed` |
| Sending Stones | 243 | `standard` (Uncommon) | `sending_stones` |
| Sentinel Shield | 243 | `standard` (Uncommon) | `sentinel` |
| Shield of Missile Attraction | 243 | `standard` (Rare, attunement) | `missile_attraction` |
| Shield of the Cavalier | 243 | `standard` (Very Rare, attunement) | `cavalier` |

Variant rules used in this batch:

- Every family has the single variant `standard`. The Rod of Lordly Might's
  button forms are properties of one item, not variants.
- The Rod of Lordly Might (pp. 241–242) and the Rope of Entanglement (pp.
  242–243) record `source_pages` in their payloads.
- The shields use the base item option `shield`. The Scimitar of Speed's base is
  `scimitar`.
- Sending Stones are one family; the pair is payload data, not a quantity.

The Rod of Resurrection (5 charges, regaining 1 at dawn) and the Scarab of
Protection (12 charges, never recharged) have charge pools. Their `per_use` is
1; Resurrection's 5-charge cost is in the payload. Once-until-dawn properties
and the Rod of Security's 10-day wait are use limits. No family lists a price.

Composing the Healing contract fixture with batches #225 to #240 and this batch
yields 204 families and 265 variants with no shared identity. The Healing
fixture, `ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads the pages column by column as one text stream. It
lists every payload leaf with its printed phrase, fails on a leaf it does not
list, requires a number to appear in its phrase, and derives each family's
pages from the layout. It supplements the literal fixtures in
`TestItems241SourceValues`; it is not an automatic conformance proof.

```sh
rtk proxy python3 -I scripts/verify_srd_items_241.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 18 (#242), pp. 243–246

Twelve families and 23 variants, in planning inventory order. Files:
`item_details_242.go`, `testdata/items_242.json`, `items_242_test.go` and
`scripts/verify_srd_items_242.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Shield, +1, +2, or +3 | 243 | `plus-1`, `plus-2`, `plus-3` (Uncommon, Rare, Very Rare) | `shield_bonus`: that rarity's Armor Class bonus |
| Slippers of Spider Climbing | 244 | `standard` (Uncommon, attunement) | `spider_slippers` |
| Sovereign Glue | 244 | `standard` (Legendary) | `glue` |
| Spell Scroll | 244 | `cantrip`, `level-1`, `level-2`, `level-3`, `level-4`, `level-5`, `level-6`, `level-7`, `level-8`, `level-9` (Common, Common, Uncommon, Uncommon, Rare, Rare, Very Rare, Very Rare, Very Rare, Legendary) | `scroll`: that spell level's table row and the shared reading and copying rules |
| Spellguard Shield | 244 | `standard` (Very Rare, attunement) | `spellguard` |
| Sphere of Annihilation | 244 | `standard` (Legendary) | `annihilation`: control, movement and the 3-row 1d100 interaction table |
| Staff of Charming | 245 | `standard` (Rare, attunement) | `charming` |
| Staff of Fire | 245 | `standard` (Very Rare, attunement) | `fire` |
| Staff of Frost | 245 | `standard` (Very Rare, attunement) | `frost` |
| Staff of Healing | 245 | `standard` (Rare, attunement) | `healing` |
| Staff of Power | 245 | `standard` (Very Rare, attunement) | `power`: quarterstaff bonuses, nine spells, the last-charge roll and Retributive Strike |
| Staff of Striking | 246 | `standard` (Very Rare, attunement) | `striking` |

Variant rules used in this batch:

- "Shield, +1, +2, or +3" has one variant per rarity, `plus-1` to `plus-3`, as
  the armor, ammunition and weapon families of earlier batches do.
- The Spell Scroll prints "Rarity Varies" and a table by spell level. Each of the
  ten rows is a variant, `cantrip` and `level-1` to `level-9`, with the rarity,
  save DC and attack bonus of its row; the cantrip row stores `spell_level` 0.
  The shared reading and copying rules repeat in each payload. A variant names no
  spell. The mundane-gear Spell Scroll presentation belongs to a later gear batch,
  which must join this family's identity rather than add a second one.
- The Sphere of Annihilation (pp. 244–245) and the Staff of Power (pp. 245–246)
  record `source_pages` in their payloads. The Sphere's interaction table is a
  random table, so it is payload data; a printed `00` is stored as 100.
- Staff attunement classes are the payload's `attunement_by_any_of`; the shared
  restriction list stays empty.
- In the Staff of Power spell table, `printed_level` is `level-5-version` where
  the source says so and `none` where it prints no level. The Staff of Healing's Cure
  Wounds cost of 1 charge per spell level is its own `cure_wounds` field.
- The shields use the base item option `shield`.

The six staffs have charge pools: Charming (10, 1d8 + 2 at dawn), Fire, Frost,
Healing and Striking (10, 1d6 + 4 at dawn) and Power (20, 2d8 + 4 at dawn).
Their `per_use` is 1, the least a charged use costs. Once-until-dawn properties
are use limits. No family lists a price.

Composing the Healing contract fixture with batches #225 to #241 and this batch
yields 216 families and 288 variants with no shared identity. The Healing
fixture, `ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads the pages column by column as one text stream. It
lists every payload leaf with its printed phrase, fails on a leaf it does not
list, requires a number to appear in its phrase, and derives each family's
pages from the layout. It supplements the literal fixtures in
`TestItems242SourceValues`; it is not an automatic conformance proof.

```sh
rtk proxy python3 -I scripts/verify_srd_items_242.py /path/to/SRD_CC_v5.2.1.pdf
```

## Magic items 19 (#243), pp. 246–248

Twelve families and 12 variants, in planning inventory order. Files:
`item_details_243.go`, `testdata/items_243.json`, `items_243_test.go` and
`scripts/verify_srd_items_243.py`. `item_source.go` is unchanged.

| Family | Page | Variants | Typed payload |
| --- | --- | --- | --- |
| Staff of Swarming Insects | 246 | `standard` (Rare, attunement) | `swarming_insects` |
| Staff of Thunder and Lightning | 247 | `standard` (Very Rare, attunement) | `thunder_lightning` |
| Staff of Withering | 248 | `standard` (Rare, attunement) | `withering` |
| Staff of the Magi | 246 | `standard` (Legendary, attunement) | `magi`: quarterstaff bonuses, Spell Absorption, nineteen spells, the last-charge roll and Retributive Strike |
| Staff of the Python | 247 | `standard` (Uncommon, attunement) | `python` |
| Staff of the Woodlands | 247 | `standard` (Rare, attunement) | `woodlands` |
| Stone of Controlling Earth Elementals | 248 | `standard` (Rare) | `earth_stone` |
| Stone of Good Luck (Luckstone) | 248 | `standard` (Uncommon, attunement) | `luckstone` |
| Sun Blade | 248 | `standard` (Rare, attunement) | `sun_blade` |
| Sword of Life Stealing | 248 | `standard` (Rare, attunement) | `life_stealing` |
| Sword of Sharpness | 248 | `standard` (Very Rare, attunement) | `sharpness` |
| Sword of Wounding | 248 | `standard` (Rare, attunement) | `wounding` |

Variant rules used in this batch:

- Every family has the single variant `standard`. The sword families' weapon
  lists are base item options, not variants.
- The Staff of the Magi (pp. 246–247) and the Staff of Thunder and Lightning
  (pp. 247–248) record `source_pages` in their payloads. Inventory order is not
  print order: Staff of the Magi, the Python and the Woodlands are printed before
  Staff of Thunder and Lightning and Staff of Withering. The catalogue keeps
  inventory order.
- Staff attunement classes are the payload's `attunement_by_any_of`; the shared
  restriction list stays empty.
- In the Staff of the Magi spell table, `printed_level` is `level-7-version` where
  the source says so and `none` where it prints no level. A cost of 0 is printed.
- The Sword of Life Stealing's 15 and the Sword of Sharpness's 14 extra damage
  are flat printed amounts, not dice.

Swarming Insects (10, 1d6 + 4 at dawn), Withering (3, 1d3 at dawn), the Magi
(50, 4d6 + 2 at dawn) and the Woodlands (6, 1d6 at dawn) have charge pools.
Their `per_use` is 1, the least a charged use costs. Once-until-dawn properties
and the Staff of the Python's 1-hour wait are use limits. No family lists a
price.

Composing the Healing contract fixture with batches #225 to #242 and this batch
yields 228 families and 300 variants with no shared identity. The Healing
fixture, `ResolveHealingPotion` and the coverage register are unchanged.

The comparison script reads the pages column by column as one text stream. It
lists every payload leaf with its printed phrase, fails on a leaf it does not
list, requires a number to appear in its phrase, and derives each family's
pages from the layout. It supplements the literal fixtures in
`TestItems243SourceValues`; it is not an automatic conformance proof.

```sh
rtk proxy python3 -I scripts/verify_srd_items_243.py /path/to/SRD_CC_v5.2.1.pdf
```

This work includes material from the System Reference Document 5.2.1 (“SRD
5.2.1”) by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd.
The SRD 5.2.1 is licensed under the Creative Commons Attribution 4.0 International
License, available at https://creativecommons.org/licenses/by/4.0/legalcode.
Names, numeric values, table rows and requirements are adapted into Edda's
structured item records; the software license does not replace this source license.
