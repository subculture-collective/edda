# Offline spell record contract (schema version 1)

Issue #194 defines `SpellRecords`, decoded by `DecodeSpellRecords(io.Reader)` in
`internal/rules/srd521/spell_record.go`. This is offline data validation. It has no
runtime consumers, action registration, database writes or coverage promotion.
`testdata/spell_records.json` contains one complete new source record, Acid Arrow,
and adaptations of the two existing bounded fixtures. Subsequent catalogue
imports are separately owned; the #218, #198, #196, #199, #200 and #201 batches are described below.

The envelope retains the exact `Source` pin, license, attribution and adaptation
notice checked by the conformance register. Each record has its own stable
`srd-5.2.1:content/spells/<slug>` ID, matching source ID and one-based PDF page.
The existing issue-80 inventory supplies these identities and pages. The source
range and class lists describe source facts; they do not grant playable targets
or casting abilities. See [the source pin](srd-5.2.1-conformance.md) and
[the expansion decision](investigations/issue-80-spell-item-expansion.md).

| Field | Contract |
| --- | --- |
| `scope` | `complete-source-record` or `bounded-self-level-1` |
| `level`, `school`, `lists` | Spell level 0–9, school and unique class lists |
| `casting_time`, `range` | Casting time category; feet, touch or self range |
| `target` | Single target shape, count and explicit eligibility |
| `components` | Explicit verbal/somatic flags; material description, cost in copper pieces and explicit consumption flag, or null material |
| `ritual`, `concentration`, `duration`, `duration_seconds` | Explicit flags and instantaneous or positive timed duration |
| `scaling` | None, or per-slot dice increments with base level and named effect references |
| `effects` | Named mechanic dependency, timing, dice, damage type, bonus, casting modifier and miss behavior |
| `required_mechanics` | Nonempty unique dependencies; names never register resolvers |
| `slot_levels` | Ordered unique data slot levels; bounded fixtures require exactly `[1]` |

The initial version covers the shapes needed by these three fixtures. The #218
extensions below add reaction triggers, long casting times and area shapes; #196
adds cantrip upgrades, one-minute casts and Emanations. Other new spell-specific
payloads still need explicit contract extensions and offline fixtures before import.
A structurally valid record is not proof of complete source transcription.
Reviewers must compare the actual source and preserve every required mechanic.
Unknown JSON fields, invalid shapes, duplicate IDs, missing source links and
invalid dependency references fail validation. Unknown mechanic names can be
retained for future planning, but remain unavailable.

Acid Arrow records the Wizard level-two Evocation, Action, 90-foot range,
V/S/powdered rhubarb leaf, instantaneous duration and single target. It retains
4d4 initial Acid damage, half initial damage only on a miss, and 2d4 damage at
the end of the target's next turn on a hit. Both damage components gain 1d4 per
slot above level two. Its dependencies remain unavailable, including ranged
spell attacks, delayed damage, miss damage, material handling and scaling.
No Acid Arrow action is added.

Cure Wounds retains source page 121, Touch range and source class lists while
restricting the fixture to conscious self and one level-one slot. Its 2d8 plus
casting modifier healing is checked against `ResolveCureWounds`. Higher-slot
healing and other targets belong to later work. Shield of Faith retains page
162, 60-foot source range, V/S/prayer scroll, Bonus Action, +2 AC and ten-minute
concentration; its fixture also stays self-only and level-one-slot-only. The
existing Cleric resolver supplies the Holy Symbol focus substitution. These
bounded records do not claim complete catalogue coverage for either spell.

`MechanicsAvailable()` checks only the explicit dependencies of those two
existing bounded fixtures. Complete source records always return false. Missing,
unknown or unsupported required keys also return false. This method is a data
inspection result; the existing resolvers still enforce preparation, components,
slots, target eligibility and lifecycle rules. No production path calls it.

Offline tests cover the exact Acid Arrow payload, bounded existing fixtures,
unavailable dependencies and unchanged partial coverage. The negative fixture
matrix in `testdata/spell_record_invalid.json` exercises malformed fields,
source/attribution/pin loss, duplicate mechanic keys and attempted target/slot
widening; tests also reject duplicate records and trailing JSON. Run:

```sh
rtk proxy go test -short -race -count=1 ./internal/rules/srd521
rtk proxy go test -short -race ./...
rtk proxy go vet ./...
```

Source verification used the existing local PDF with SHA-256
`8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`:
spell entries on pages 107, 121 and 162, with component rules on pages 105–106.
CI and offline tests require no PDF or network. The PDF and extracted text are
not vendored.

This work includes material from the System Reference Document 5.2.1 (“SRD
5.2.1”) by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd.
The SRD 5.2.1 is licensed under the Creative Commons Attribution 4.0 International
License, available at https://creativecommons.org/licenses/by/4.0/legalcode.
Spell facts are adapted into structured Edda metadata; Cure Wounds and Shield of
Faith are deliberately restricted to existing fixtures. The software license
does not replace CC-BY-4.0 for the source-derived material.

## Issue #218: Scorching Ray through Shield of Faith

`testdata/spells_218.json` contains exactly twelve complete source families from
pinned PDF pages 159–162. `spells_218_test.go` checks their names, initial source
pages, headers, materials, effects, scaling, mechanic dependencies and nested
choices. Seeming continues onto page 161; Shield continues onto page 162. The
source PDF hash was checked against the pin above. No PDF or extracted text is
needed by the offline tests.

This batch extends schema version 1 with optional, typed `source_details` on
complete source records. It adds only shapes present in this batch: a ten-minute
cast, conditional Bonus Action and Reaction triggers, unlimited range, selected
creatures, a Sphere, rays, alternative targets, indefinite duration and per-slot
ray count scaling. `scaling.kind = slot-count` uses zero dice fields and the typed
`rays.additional_per_slot` count; its effect reference names the additional ray.
The existing `slot` shape still represents damage dice additions. A timed round
is recorded as six seconds, with Shield's exact start-of-next-turn boundary
retained separately in `end_timing`.

Each selected detail payload requires every field, including explicit false,
zero and list values. There is no catch-all prose or arbitrary parameter map.
The typed payloads retain:

- Scorching Ray's three attacks, one/several-target choice and extra ray per slot.
- Scrying's creature and seen-location modes, all three knowledge and three
  physical-connection modifiers, success lockout and sensor behavior.
- Searing Smite's melee weapon/Unarmed Strike triggers, initial and recurring
  damage, after-damage Constitution save and scaling of both damage effects.
- Secret Chest's separate 5,000 GP and 50 GP component minima, dimensions,
  capacity, Ethereal storage, Magic actions, daily cumulative decay and endings.
  The common material field records the combined 5,050 GP minimum; the typed
  fields retain both independently required objects. Neither is consumed.
- See Invisibility's creature/object perception and ghostly Ethereal view.
- Seeming's per-target appearance choices, height/weight/limb constraints,
  unwilling-target save, physical inspection and Study discovery check.
- Sending's two addressing choices, 25-word message and immediate like reply,
  cross-plane delivery failure, sender feedback and recipient blocking.
- Sequester's object/willing-creature choices, concealment, creature suspension,
  waived sustenance and aging, chosen early-end condition and damage ending.
- Shapechange's seen-form and CR restrictions, first-form temporary Hit Points,
  full retained-statistics list, retained Spellcasting and equipment choices.
- Shatter's ten-foot Sphere, Constitution save/half damage, Construct
  Disadvantage, eligible nonmagical objects and per-slot damage scaling.
- Shield's two reaction triggers, +5 AC including the triggering attack,
  Magic Missile protection and exact next-turn expiration.
- Shield of Faith's chosen creature within 60 feet, +2 AC and ten-minute
  concentration, with all source slot levels retained as unavailable metadata.

The full Shield of Faith source record is separate from the existing bounded
fixture. `ComposeSpellSourceCatalogue` validates its inputs, excludes bounded
compatibility fixtures and rejects duplicate complete source IDs. Composing the
#194 and #218 fixtures yields thirteen distinct source records: Acid Arrow plus
this batch. The bounded fixtures keep their existing IDs and exact availability
checks; added detail payloads cannot make those checks pass.

All twelve source records remain unavailable regardless of their dependency
keys. This batch supplies data for later effect-instance work (#351); it adds
no resolver, action, entitlement, tool, inventory or coverage claim. Negative
fixtures cover omitted detail fields, missing dependencies, invalid variants,
costs, chances, triggers and scaling references, unknown dependencies and
attempts to attach source metadata to bounded compatibility records.

The source-derived JSON retains the pinned source, CC-BY-4.0 attribution and an
adaptation notice. It is an adaptation into structured metadata under the source
license described above; the repository software license does not replace it.

## Issue #198: Barkskin through Burning Hands

`testdata/spells_198.json` imports the exact twelve families on PDF pp.112–114.
Bestow Curse retains its p.113 continuation and Blink its p.114 continuation in
explicit `source_pages` lists. All records retain the pinned source attribution
and adaptation notice. `spells_198_test.go` supplies literal source-value fixtures
for every header, choice, effect, dependency and scaling value.

The typed extensions retain Barkskin's AC 17 and willing target; Beacon of Hope's
Wisdom/Death save Advantage and maximum healing; Befuddlement's damage, Magic
restrictions, 30-day repeat save and three restoration alternatives; all four
Bestow Curse choices, six ability options and five slot-duration bands; Black
Tentacles' Strength saves, once-per-turn triggers and Athletics escape; both Blade
Barrier forms, dimensions, cover and Force damage; Bless's unconsumed Holy Symbol
worth at least 5 GP and target scaling; both Blight target modes; both
Blindness/Deafness conditions and repeat saves; Blink's d6, Ethereal interaction
and return cases; Blur's perception exceptions; and Burning Hands' Cone damage
and eligible object ignition.

`slot-targets` references the typed additional-target count for Bless and
Blindness/Deafness. `slot-duration` references Bestow Curse's slot table, including
concentration removal at slot 5 and indefinite duration at slot 9. Zero seconds
with `until_dispelled=true` is explicit source data. Curse choices use `none` or
`not-applicable` for fields outside the selected choice rather than omitting
required values. These fields are metadata; no effect resolver reads them.

Every new payload field is required recursively. Null, omitted values, blank
labels, empty required lists (including `[ ]`), mixed entries, duplicate keys,
case aliases and unknown fields fail offline decoding. Required common numeric
fields also distinguish omission/null from explicit zero. Negative Scrying
modifiers in the earlier batch remain supported. Composition of #194, #218 and
#198 yields 25 unique complete source records and excludes both bounded fixtures.
Their exact availability checks are unchanged; every imported record remains
unavailable. No runtime, inventory, tool, action, coverage or Full Rules behavior
is added.

To recheck the PDF pin, inventory identities/pages and edition-specific source
anchors, run:

```sh
rtk proxy python3 scripts/verify_srd_spells_198.py /tmp/edda-wave1/SRD_CC_v5.2.1.pdf
```

The source check supplements the full offline literal fixtures; it is not an
automatic conformance proof. The PDF and extracted text are not vendored. No
network, inference, UI or database tests are needed for this data-only batch.

## Issue #196: Acid Arrow through Antimagic Field

`testdata/spells_196.json` imports the exact twelve families the issue-80
inventory assigns to #196, on PDF pp.107–109. Animal Friendship retains its p.108
continuation and Animate Objects its p.109 continuation, including the Animated
Object stat block printed with that spell, in explicit `source_pages` lists.
`spells_196_test.go` supplies literal source-value fixtures for every header,
choice, variant, effect, dependency and scaling value, and checks the batch's IDs,
names and pages against the inventory.

Acid Arrow is restated unchanged from the #194 fixture and has no typed payload:
the common fields already hold its whole source entry. The two files therefore
cannot be composed together; a later composition uses this batch for Acid Arrow
and `spell_records.json` only for its two bounded compatibility fixtures.

The typed extensions retain Acid Splash's five-foot Sphere, Dexterity save and
three cantrip upgrade tiers; Aid's three targets and equal maximum/current Hit
Point increase; Alarm's door/window/area choice, exemptions and both audible and
mental variants; all three Alter Self options, including the four natural weapons
and their damage types; Animal Friendship's Wisdom save and damage ending; Animal
Messenger's Charisma save, Challenge Rating exception, message limit and both
travel rates; Animal Shapes' form limits, retained statistics and Temporary Hit
Points; Animate Dead's bones/corpse choice, commands and 24-hour control with
reassertion of up to four creatures; Animate Objects' eligibility, size costs,
commands and the full Animated Object stat block with its per-size Hit Points and
Slam damage; Antilife Shell's exempt creature types and forced-passage ending;
and each Antimagic Field prohibition, suppression rule and exception.

Schema version 1 gains only shapes present in this batch:

- A cantrip has level 0 and an explicitly empty `slot_levels` list; every other
  spell still needs a nonempty one. `slot_levels` must be present and non-null.
  `cantrip-upgrade` scaling carries the per-tier die and references the typed tiers.
- `casting_time = 1-minute` requires `casting_seconds = 60`, and the reverse.
  `ritual = true` records the source's "or Ritual" alternative.
- `slot-bonus` (Aid), `slot-creatures` (Animate Dead) and `slot-dice-by-size`
  (Animate Objects) use zero common dice and reference a typed increment.
  `slot-targets` and `slot-duration` now also accept Animal Friendship's target
  increment and Animal Messenger's 48-hour increment.
- New target shapes: `emanation`, `selected-objects`, and further eligibility
  labels for the existing `single`, `choice`, `sphere` and `selected-creatures`.
- Stat block ability modifiers and saves are the only signed payload values.
  Validation pins the whole Animated Object ability table, so a changed sign or
  value fails rather than passing as a valid negative.

Each payload is tied to its own spell ID and its source scaling kind. Each of
the twelve IDs is also bound to its printed name, page, level, range, duration
and concentration, and validation compares the printed numbers, saves and
immunity lists by equality (for example the 5-foot Sphere, Aid's 5 Hit Points,
the 20-foot Cube, twenty-five words, Challenge Rating 4, four reasserted
creatures, AC 15 and both 10-foot Emanations), so a structurally valid record
with a different source value is rejected. Acid Arrow has no payload; its
4d4, 2d4 and 1d4 dice are pinned with its header. The
recursive presence, null, blank, empty-list, duplicate-key, alias and unknown-field
guards from #198 apply to every new field. Composition of #196, #218 and #198
yields 36 unique complete source records. Every imported record remains
unavailable: no resolver, action, entitlement, tool, inventory, coverage or Full
Rules behavior is added, and the bounded fixtures' availability checks are
unchanged.

To recheck the PDF pin, inventory identities/pages, all twelve source headers and
edition-specific source anchors, run:

```sh
rtk proxy python3 scripts/verify_srd_spells_196.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #197: Antipathy/Sympathy through Banishment

`testdata/spells_197.json` imports the exact twelve families the issue-80
inventory assigns to catalogue issue #197, on PDF pp.109–112: Antipathy/Sympathy,
Arcane Eye, Arcane Hand, Arcane Lock, Arcane Sword, Arcanist’s Magic Aura, Astral
Projection, Augury, Aura of Life, Awaken, Bane and Banishment. Antipathy/Sympathy
continues onto p.110, Arcanist’s Magic Aura onto p.111 and Awaken onto p.112;
their payloads keep explicit `source_pages` lists. `spells_197_test.go` holds
literal source-value fixtures for every header, choice, effect, dependency and
scaling value, and compares the identities and first pages with the inventory.

The typed payloads retain:

- Antipathy/Sympathy's two modes with their conditions and movement rules, the
  creature-or-object target, caster-specified creature kind, 120-foot Wisdom
  saves, the Sympathy 5-foot restriction and damage save, and the one-minute
  immunity after a successful save.
- Arcane Eye's Invisible, invulnerable sensor, 30-foot Darkvision, 30-foot Bonus
  Action movement and one-inch opening.
- Arcane Hand's Large AC 20 object, Hit Points, 60-foot movement and all four
  options: Clenched Fist, Forceful Hand, Grasping Hand and Interposing Hand.
- Arcane Lock's five eligible closed objects, bypass list and password.
- Arcane Sword's melee spell attack, 4d12 Force damage plus modifier and 30-foot
  Bonus Action movement.
- Arcanist's Magic Aura's Mask (creature) and False Aura (object) variants,
  the three False Aura choices and the 30-day until-dispelled rule.
- Astral Projection's caster plus up to eight willing targets, both per-target
  consumed components, suspended body, silver cord, plane exit and endings.
- Augury's four-row Omens table, 30-minute window and cumulative 25 percent
  chance of no answer.
- Aura of Life's 30-foot Emanation, Necrotic Resistance, protected Hit Point
  maximum and 1 Hit Point recovery.
- Awaken's three target kinds, Intelligence 3 limit and Intelligence 10 result,
  language, natural-plant changes and 30-day Charmed condition.
- Bane's three visible targets, Charisma save, 1d4 subtraction and target scaling.
- Banishment's Charisma save, demiplane, Incapacitated condition, return space,
  five creature types that do not return and target scaling.

This batch extends schema version 1 with three casting times, `1-minute`,
`1-hour` and `8-hours`. Their `casting_seconds` must be 60, 3,600 and 28,800;
other casting times in this batch keep zero. Augury is a one-minute cast with
`ritual=true`. New target shapes are `created-effect`, `emanation` and
`self-and-up-to-count`, plus further `choice`, `single` and `up-to-count`
eligibilities. `slot-targets` now also reads Bane's and Banishment's typed
counts. The new `slot-effect-dice` scaling kind uses zero common dice fields and
one typed `hand.slot_dice` row per referenced effect, in order: 2d8 for
Clenched Fist and 2d6 for the Grasping Hand crush.

Priced materials are in copper pieces: Arcane Lock 2,500 (consumed), Arcane
Sword 25,000, Augury 2,500 and Awaken 100,000 (consumed). Astral Projection's
common field records the 110,000 per-target minimum; the typed fields keep the
100,000 jacinth and 10,000 silver bar separately. A record accepts exactly one
typed payload across all batches, and it must be the payload for that spell.

Composition of #194, #218, #198 and #197 yields 37 unique complete source
records and excludes both bounded fixtures, whose availability checks are
unchanged. Every imported record remains unavailable. No runtime, inventory,
tool, action, coverage or Full Rules behavior is added. The JSON retains the
pinned source, CC-BY-4.0 attribution and an adaptation notice.

To recheck the PDF pin, inventory identities and pages, verbatim source headers
and edition-specific anchors, run:

```sh
rtk proxy python3 scripts/verify_srd_spells_197.py /path/to/SRD_CC_v5.2.1.pdf
```

The comparator rebuilds each source header from the structured record and
requires it to equal the PDF text. It supplements the literal offline fixtures
and is not an automatic conformance proof. The PDF and extracted text are not
vendored.
## Issue #199: Call Lightning through Color Spray

`testdata/spells_199.json` imports the twelve families the issue-80 inventory
assigns to #199, on PDF pp.114-116: Call Lightning, Calm Emotions, Chain
Lightning, Charm Monster, Charm Person, Chill Touch, Chromatic Orb, Circle of
Death, Clairvoyance, Clone, Cloudkill and Color Spray. Charm Monster keeps its
p.114 to p.115 continuation and Clone its p.115 to p.116 continuation in explicit
`source_pages` lists. `spells_199_test.go` supplies literal source-value fixtures
for every header, choice, variant, effect, dependency and scaling value.

The typed payloads retain Call Lightning's cylinder cloud, bolt and storm bonus;
Calm Emotions' two per-creature choices; Chain Lightning's three secondary bolts
and per-slot bolt increment; the Monster and Person charm targets and endings;
Chill Touch's three cantrip tiers; Chromatic Orb's damage types and leaps;
Circle of Death's sphere; Clairvoyance's sensor, senses and focuses; Clone's
choices, 120-day growth, diamond and vessel; Cloudkill's moving sphere; and Color
Spray's cone.

Schema version 1 gains only shapes present in this batch, on top of #196:

- `casting_time = 1-hour` (Clone), with `casting_seconds = 3600`.
- Target shapes `cylinder`, `chain` and `location`, plus labels for `cone`,
  `up-to-count`, `single` and `sphere` (for example the humanoid and creature
  sphere labels).
- `slot-count` scaling also covers Chain Lightning's extra bolt per slot, and
  `slot-targets` also covers the Charm Monster and Charm Person extra target.
- Chill Touch reuses #196's `cantrip-upgrade` scaling; its tiers live in the
  `chill` payload, which that scaling check now accepts beside Acid Splash's.

Validation pins the printed measurements by equality (Call Lightning's 10-foot
height, 60-foot radius and 5-foot bolt, the 20-foot Calm Emotions and Cloudkill
radii, Chain Lightning's and the orb's 30-foot range, Circle of Death's 60-foot
radius, Cloudkill's 10-foot move, the 15-foot Color Spray cone and Clone's 120
days), so a positive but wrong value is rejected.

#196 restates Acid Arrow, so the #194 fixture cannot be composed beside it.
Composition of the two bounded compatibility fixtures with #196, #218, #198
and #199 yields 48 unique records. Every imported record remains unavailable: no
resolver, action, entitlement, tool, inventory, coverage or Full Rules behavior
is added.

To recheck the PDF pin, identities, pages and source anchors, run:

```sh
rtk proxy python3 scripts/verify_srd_spells_199.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #200: Command through Conjure Minor Elementals

`testdata/spells_200.json` imports the twelve families the issue-80 inventory
assigns to #200, on PDF pp.116-118: Command, Commune, Commune with Nature,
Comprehend Languages, Compulsion, Cone of Cold, Confusion, Conjure Animals,
Conjure Celestial, Conjure Elemental, Conjure Fey and Conjure Minor Elementals.
Commune with Nature keeps its p.116 to p.117 continuation and Conjure Animals its
p.117 to p.118 continuation in explicit `source_pages` lists. Conjure Woodland
Beings, which follows on p.118, belongs to a later batch and is not imported.
`spells_200_test.go` supplies literal source-value fixtures for every header,
choice, table row, effect and scaling value, and checks the family list and
pages against the inventory.

The typed payloads retain Command's five options; Commune's three yes-or-no
questions and cumulative no-answer chance; Commune with Nature's two radii and
five facts; Comprehend Languages' reading limits; Compulsion's directed
movement; Cone of Cold's cone; Confusion's behavior table and direction die;
Conjure Animals' pack, movement and triggers; Conjure Celestial's two lights;
Conjure Elemental's four elements and restraint; Conjure Fey's attack, fear and
teleport; and Conjure Minor Elementals' emanation and damage types.

Schema version 1 gains only shapes present in this batch, on top of #199:

- Commune, Commune with Nature and Comprehend Languages set `ritual = true`
  with their normal casting time (`1-minute` or `action`) and a
  `ritual_alternative` flag in the payload, as #196 does for Alarm.
- Target `created-effect` gains `unoccupied-space-in-range`, and `emanation`
  gains `creatures-in-emanation`. The other targets reuse existing shapes.
- `slot-targets` also covers Command's extra creature per slot.
- `slot-radius` scaling (Confusion) references `additional_radius_feet_per_slot`
  in the `confusion` payload and the `scaling.radius-per-slot` dependency.
- `slot` scaling also lists several effects: Conjure Celestial scales its healing
  and damage together, and Conjure Elemental scales both its 8d8 and 4d8 damage
  because the source says only "the damage increases".

Validation pins each record's printed level, school, casting time, range, target,
ritual flag, concentration, duration, scaling dice and effect dice by equality,
and the printed measurements of each payload (for example the 60-foot Cone of
Cold, the 10-foot Confusion sphere, the 40-foot Celestial cylinder and the 30-foot
Fey teleport), so a positive but wrong value is rejected.

Composition of the two bounded compatibility fixtures with #196, #218, #198,
#197, #199 and #200 yields 72 unique records. Every imported record remains
unavailable: no resolver, action, entitlement, tool, inventory, coverage or Full
Rules behavior is added.

To recheck the PDF pin, identities, pages and source anchors, run:

```sh
rtk proxy python3 scripts/verify_srd_spells_200.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #201: Conjure Woodland Beings through Creation

`testdata/spells_201.json` imports the twelve families the issue-80 inventory
assigns to #201, on PDF pp.118-121: Conjure Woodland Beings, Contact Other
Plane, Contagion, Contingency, Continual Flame, Control Water, Control Weather,
Counterspell, Create Food and Water, Create Undead, Create or Destroy Water and
Creation. Conjure Woodland Beings keeps its p.118 to p.119 continuation and
Control Water its p.119 to p.120 continuation in explicit `source_pages` lists.
`spells_201_test.go` checks the family list and pages against the inventory and
supplies literal source-value fixtures for every header, choice, table row,
effect, dependency and scaling value.

The typed payloads retain the woodland emanation and its triggers; Contact Other
Plane's save, five questions, one-word answers and failure; Contagion's
repeated saves and seven-day duration; Contingency's contingent-spell limits and
statuette; Continual Flame's light; Control Water's four effects, including the
flood and whirlpool measurements; Control Weather's delay and three stage
tables; Counterspell's trigger and wasted costs; the food and water amounts;
Create Undead's corpses, commands, 24-hour control and per-slot creature
options; Create or Destroy Water's two choices; and Creation's Materials table.

Schema version 1 gains only shapes present in this batch, on top of #199:

- Duration `special` (Creation). It has `duration_seconds = 0` and is never
  concentration; the real durations live in the `creation.materials` table and
  the record requires `duration.special-by-material`. The shared validator rejects `special` unless typed Creation details are present.
- Casting time `reaction` with exactly one entry in `casting_triggers`
  (Counterspell). Every other record in this batch keeps that list empty.
- Target shape `cube`, plus labels for `emanation`, `up-to-count`, `single` and
  `choice` (the water and corpse targets).
- Scaling kind `slot-increment`, which reads the typed per-slot water and Cube
  increments, and `slot-creature-table`, which reads Create Undead's rows for
  slots 7, 8 and 9. Both use zero common dice fields.
- Conjure Woodland Beings reuses the `slot` scaling kind for its extra 1d8.

Validation pins the printed measurements by equality (the 10-foot emanation,
DC 15, three saves, seven days, the 1,500 GP statuette, the 20-foot light radii,
the 100-foot Cube, the whirlpool sizes, the 5-mile weather radius, the 45 pounds
and 30 gallons, the 150 GP onyx, the 24-hour control and every Creation duration
row), so a positive but wrong value is rejected.

Composition of the two bounded compatibility fixtures with #196, #218, #198,
#199 and #201 yields 72 unique records. Every imported record remains
unavailable: no resolver, action, entitlement, tool, inventory, coverage or Full
Rules behavior is added.

To recheck the PDF pin, identities, pages, rebuilt headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_201.py /path/to/SRD_CC_v5.2.1.pdf
```

The comparator rebuilds each source header from the structured record and
requires it to equal the PDF text. It supplements the literal offline fixtures
and is not an automatic conformance proof. The PDF and extracted text are not
vendored.
