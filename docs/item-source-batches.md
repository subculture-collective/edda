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

This work includes material from the System Reference Document 5.2.1 (“SRD
5.2.1”) by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd.
The SRD 5.2.1 is licensed under the Creative Commons Attribution 4.0 International
License, available at https://creativecommons.org/licenses/by/4.0/legalcode.
Names, numeric values, table rows and requirements are adapted into Edda's
structured item records; the software license does not replace this source license.
