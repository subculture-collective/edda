# Offline spell record contract (schema version 1)

Issue #194 defines `SpellRecords`, decoded by `DecodeSpellRecords(io.Reader)` in
`internal/rules/srd521/spell_record.go`. This is offline data validation. It has no
runtime consumers, action registration, database writes or coverage promotion.
`testdata/spell_records.json` contains one complete new source record, Acid Arrow,
and adaptations of the two existing bounded fixtures. Subsequent catalogue
imports are separately owned; the #218, #198, #196, #199, #200, #201, #202, #203, #204, #205, #206, #207, #208,
#209, #210, #211, #212, #213, #214, #215, #216, #217, #219, #220, #221, #222, #223 and #224 batches are described below.

The envelope retains the exact `Source` pin, license, attribution and adaptation
notice checked by the conformance register. Each record has its own stable
`srd-5.2.1:content/spells/<slug>` ID, matching source ID and one-based PDF page.
The existing issue-80 inventory supplies these identities and pages. The source
range and class lists describe source facts; they do not grant playable targets
or casting abilities. See [the source pin](srd-5.2.1-conformance.md) and
[the expansion decision](investigations/issue-80-spell-item-expansion.md).

| Field | Contract |
| --- | --- |
| `scope` | `complete-source-record`, `bounded-self-level-1` or `bounded-carried-object-cantrip` (Light, #264) |
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

## Issue #202: Cure Wounds through Detect Thoughts

`testdata/spells_202.json` imports the twelve families the issue-80 inventory
assigns to #202, on PDF pp.121-123: Cure Wounds, Dancing Lights, Darkness,
Darkvision, Daylight, Death Ward, Delayed Blast Fireball, Demiplane, Detect Evil
and Good, Detect Magic, Detect Poison and Disease and Detect Thoughts. Dancing
Lights keeps its p.121 to p.122 continuation and Demiplane its p.122 to p.123
continuation in explicit `source_pages` lists. `spells_202_test.go` supplies
literal source-value fixtures for every header, variant, effect and scaling
value, and checks the family list and pages against the decision inventory.

The typed payloads retain Cure Wounds' healing and per-slot dice; Dancing
Lights' four lights, forms and movement; Darkness' and Daylight's sphere,
object emanation, blocker and dispelled light; Darkvision's 150-foot sense;
Death Ward's two protections; Delayed Blast Fireball's bead, growth, touch and
throw rules; Demiplane's door, room and connection choices; the shared barrier
list of the three detection spells; and Detect Thoughts' two options and probe.

Schema version 1 gains only shapes present in this batch, on top of #200:

- Target shapes `created-effect` (four torch-size lights, or a visible flat
  solid surface) and `choice` (a point in range or an unattended object).
- Detect Magic and Detect Poison and Disease use `casting_time = action` with
  `ritual = true`, as Animal Messenger does, and require `cast.ritual`.
- The complete Cure Wounds record is separate from the bounded `cure-wounds`
  compatibility fixture that shares its ID; composition skips the fixture, so
  composing the fixtures with #202 yields no duplicate.

Validation pins each record's header, scaling row, effects, material,
dependency list and typed values by equality with the printed page, so a
positive but wrong value is rejected. Composition of the bounded fixtures with
#196, #218, #198, #199, #200, #201 and #202 yields 96 unique records. Every imported record
remains unavailable: no resolver, action, entitlement, tool, inventory,
coverage or Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_202.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #203: Dimension Door through Dominate Monster

`testdata/spells_203.json` imports the twelve families the issue-80 inventory
assigns to #203, on PDF pp.123-125: Dimension Door, Disguise Self, Disintegrate,
Dispel Evil and Good, Dispel Magic, Dissonant Whispers, Divination, Divine Favor,
Divine Smite, Divine Word, Dominate Beast and Dominate Monster. Dimension Door
keeps its p.123 to p.124 continuation and Dominate Monster its p.125 to p.126
continuation in explicit `source_pages` lists. `spells_203_test.go` supplies
literal source-value fixtures for every header, variant, effect and scaling value,
and checks the family list and pages against the decision inventory.

The typed payloads retain Dimension Door's destination choices, companion rule and
collision damage; Disguise Self's appearance limits and Study check; Disintegrate's
failed-save damage, dust and object rules; Dispel Evil and Good's two special
functions; Dispel Magic's automatic and checked endings; Dissonant Whispers' damage
and forced move; Divination's question and cumulative no-answer chance; Divine
Favor's and Divine Smite's extra damage and triggers; Divine Word's hit-point table
and banishment; and the shared Dominate shape with each spell's slot-duration table.

Schema version 1 gains only shapes present in this batch, on top of #202:

- Target shape `choice` with the eligibilities `destination-in-range`,
  `visible-creature-nonmagical-object-or-force-creation-in-range` and
  `creature-object-or-magical-effect-in-range`.
- Scaling kind `slot-auto-end` (Dispel Magic), which requires the typed Dispel Magic
  payload and `scaling.auto-end-by-slot-level`; and `slot-duration-table`
  (Dominate Beast and Dominate Monster), which requires a typed slot-duration table
  and `scaling.duration-by-slot-table`.
- Disintegrate's printed +40 is the effect `bonus`; Divine Smite keeps its printed
  hit trigger in `casting_triggers`; Divination uses `casting_time = action` with
  `ritual = true` and a consumed 25 GP material.

Validation pins each record's header, scaling row, effects, material, trigger list,
dependency list and typed values by equality with the printed page, so a positive
but wrong value is rejected. Composition of the bounded fixtures with #196, #218,
#198, #199, #200, #201, #202 and #203 through #208 yield 168 unique records. Every imported
record remains unavailable: no resolver, action, entitlement, tool, inventory,
coverage or Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_203.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #204: Dominate Person through Enthrall

`testdata/spells_204.json` imports the twelve families the issue-80 inventory
assigns to #204, on PDF pp.126-128: Dominate Person, Dragon's Breath, Dream,
Druidcraft, Earthquake, Eldritch Blast, Elementalism, Enhance Ability,
Enlarge/Reduce, Ensnaring Strike, Entangle and Enthrall. Enlarge/Reduce keeps
its p.127 to p.128 continuation in an explicit `source_pages` list; every other
family sits on one page. `spells_204_test.go` supplies literal source-value
fixtures for every header, option, effect and scaling value, and checks the
family list and pages against the decision inventory.

The typed payloads retain Dominate Person's link, commands and per-slot
duration table; Dragon's Breath's five damage types and cone; Dream's trance,
dream and terror options; the option lists of Druidcraft and Elementalism;
Earthquake's terrain, fissures, structure damage and collapse; Eldritch Blast's
beam count by character level; Enhance Ability's extra targets; both
Enlarge/Reduce effects; and the vines, area and save of Ensnaring Strike,
Entangle and Enthrall.

Schema version 1 gains only shapes present in this batch, on top of #202:

- Range kind `special` (Dream). It carries `feet = 0` and requires
  `range.special`.
- Target eligibilities `visible-humanoid-in-range`, `known-creature-same-plane`,
  `creature-or-object-in-range` and `visible-creature-or-object-in-range` on a
  `single` target, and `visible-ground-point-in-range` on a `created-effect`.
- Scaling kind `cantrip-beams` (Eldritch Blast) uses zero dice fields, names the
  beam effect and requires `scaling.cantrip-beams`; the three beam steps are in
  the typed payload. `slot-duration` also accepts Dominate Person's typed table,
  and `slot-targets` accepts Enhance Ability's typed target increment.
- Earthquake's flat 50 structure damage is stored as `bonus` on a zero-dice
  effect.
- Enthrall's `-10` Wisdom (Perception) penalty is a signed payload value.
- Ensnaring Strike keeps its Bonus Action trigger in `casting_triggers`.
- Dream is a one-minute cast with `casting_seconds = 60`.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong value
is rejected. The #204 verifier script checks identities, pages, printed headers and
anchor phrases only; typed payload values (slot-duration tables, fissure numbers,
option lists) are verified by the literal fixtures and pins, not by the script. A batch payload on a record of another batch is rejected too.
Composition of the bounded fixtures with #196, #218, #198, #199, #200, #201,
#202 and #203 through #208 yield 168 unique records. Every imported record remains
unavailable: no resolver, action, entitlement, tool, inventory, coverage or Full
Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_204.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #205: Etherealness through Find Traps

`testdata/spells_205.json` imports the twelve families the issue-80 inventory
assigns to #205, on PDF pp.128-131: Etherealness, Expeditious Retreat, Eyebite,
Fabricate, Faerie Fire, Faithful Hound, False Life, Fear, Feather Fall, Find
Familiar, Find Steed and Find Traps. Each family sits wholly on one page, so no
record needs a `source_pages` list. The Otherworldly Steed stat block and the
tail of Find the Path on p.131 belong to other families and are not imported;
Find Steed names the stat block as the unavailable dependency
`stat-block.otherworldly-steed`. `spells_205_test.go` supplies literal
source-value fixtures for every header, variant, effect and scaling value, and
checks the family list and pages against the decision inventory.

The typed payloads retain Etherealness' border-plane rules, shunt damage and
per-slot extra targets; Expeditious Retreat's repeat Dash; Eyebite's three
chosen effects and retargeting; Fabricate's size limits and exclusions; Faerie
Fire's cube, colours and light; Faithful Hound's bark, bite and movement; False
Life's dice and per-slot bonus; Fear's cone and end-of-turn save; Feather Fall's
trigger and descent; Find Familiar's forms, telepathy and dismissal rules; Find
Steed's summoning, combat and disappearance rules; and Find Traps' sensing and
exclusions.

Schema version 1 gains only shapes present in this batch, on top of #202:

- Target shapes `choice` (visible raw materials), `cube` (objects and creatures
  in a cube), `up-to-count` (five falling creatures) and `selected-objects`
  (traps in line of sight in range).
- Etherealness uses `slot-targets` scaling with a typed base and per-slot target
  count; False Life uses `slot-bonus` scaling with a typed per-slot temporary
  hit point increment, as Aid does for hit points.
- Feather Fall is a `reaction` and records its trigger in `casting_triggers`.
  Find Familiar is a one-hour cast with `ritual = true`, as Augury is for one
  minute, and requires `cast.ritual`.
- Find Steed keeps `scaling.kind = none`: its source text uses the slot level as
  the stat block's spell level, which the payload and the dependency
  `stat-block.slot-level-as-spell-level` record without implementing it.

Validation pins each record's header, scaling row, effects (including timing and
bonus), material, dependency list and every typed payload value by equality
with the printed page, so a positive but wrong value is rejected. Composition of
the bounded fixtures with #196, #218, #198, #199, #200, #201, #202 and #203 through #208
yields 168 unique records. Every imported record remains unavailable: no
resolver, action, entitlement, tool, inventory, coverage or Full Rules behavior
is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_205.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #206: Find the Path through Fly

`testdata/spells_206.json` imports the twelve families the issue-80 inventory
assigns to #206, on PDF pp.131-133: Find the Path, Finger of Death, Fire Bolt,
Fire Shield, Fire Storm, Fireball, Flame Blade, Flame Strike, Flaming Sphere,
Flesh to Stone, Floating Disk and Fly. Find Traps, Fog Cloud and Forbiddance
share these pages but belong to other batches and are not imported here. Every
family sits on a single page, so no `source_pages` continuation lists are used.
`spells_206_test.go` supplies literal source-value fixtures for every header,
variant, effect and scaling value, and checks the family list and pages against
the decision inventory.

The typed payloads retain Find the Path's route, invalid destinations and
known information; Finger of Death's damage, flat bonus and Zombie; Fire Bolt's
attack, ignition and three cantrip tiers; Fire Shield's two shield choices,
resistances and retaliation; Fire Storm's ten contiguous cubes; Fireball's
sphere and per-slot dice; Flame Blade's blade, attack and light; Flame Strike's
cylinder with separate Fire and Radiant dice; Flaming Sphere's damage,
movement, barriers and light; Flesh to Stone's repeating save and its
three-of-a-kind counts; Floating Disk's size, capacity, following and ending
distances; and Fly's speed, hover and per-slot target increase. The Fireball
payload is `fireball_burst` because `fireball` already names the Delayed Blast
Fireball payload.

Schema version 1 gains only shapes present in this batch, on top of #202:

- Target shapes `choice` (a creature or object in range), `cubes` (up to ten
  contiguous 10-foot Cubes), `cylinder` and `created-effect` (an unoccupied
  ground space in range).
- Find the Path uses `casting_time = 1-minute` with `casting_seconds = 60` and
  a one-day (86,400 second) concentration duration; Flame Blade uses
  `bonus-action`; Floating Disk uses `action` with `ritual = true`.
- Fire Bolt uses `cantrip-upgrade` scaling and Fly uses `slot-targets` scaling,
  each tied to its own typed payload.

Validation pins each record's header, scaling row, effects, material,
dependency list and typed values by equality with the printed page, so a
positive but wrong value is rejected. Composition of the bounded fixtures with
#196, #218, #198, #199, #200, #201, #202 and #203 through #208 yield 168 unique records.
Every imported record remains unavailable: no resolver, action, entitlement,
tool, inventory, coverage or Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_206.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #207: Fog Cloud through Glibness

`testdata/spells_207.json` imports the twelve families the issue-80 inventory
assigns to #207, on PDF pp.133-136: Fog Cloud, Forbiddance, Forcecage, Foresight,
Freedom of Movement, Freezing Sphere, Gaseous Form, Gate, Geas, Gentle Repose,
Giant Insect and Glibness. Gaseous Form keeps its p.134 to p.135 continuation and
Giant Insect its p.135 to p.136 continuation, including the Giant Insect stat
block, in explicit `source_pages` lists. `spells_207_test.go` supplies literal
source-value fixtures for every header, effect, table row and scaling value, and
checks the family list and pages against the decision inventory.

Giant Insect's generic effect rows (`jab`, `jab-poison`, `web`) carry base dice
only with bonus 0. The source adds +3 and the spell's level to Poison Jab and
Web Bolt; the typed stat block (`piercing_bonus`, `bludgeoning_bonus` and the
`*_adds_spell_level` flags) is authoritative for those values.

The typed payloads retain Fog Cloud's sphere and radius increment; Forbiddance's
ward, planar proofing, damage, password and daily recasting; Forcecage's cage and
box forms, trapping and Charisma save; Foresight's two advantages; Freedom of
Movement's protections, swim speed and escape; Freezing Sphere's blast, ice,
break-free check and held globe; Gaseous Form's cloud form, defenses and limits;
Gate's portal, named-creature option and planar rulers; Geas' command, damage,
enders and per-slot durations; Gentle Repose's protections; Giant Insect's forms,
commands and its full stat block with the Poison Jab, Web Bolt and Venomous Spew
actions; and Glibness' Charisma check and truth rule.

Schema version 1 gains only shapes present in this batch, on top of #202:

- Target shapes `created-effect` (a floor area of up to 40,000 square feet, or an
  area in range) and `single` (a corpse or remains by touch).
- Casting times `10-minutes` with `ritual = true` (Forbiddance) and `1-minute`
  (Foresight, Geas) use the typed `casting_seconds` of 600 and 60; Gentle Repose
  uses `action` with `ritual = true`.
- `slot-radius` also accepts Fog Cloud's `additional_radius_feet_per_slot`;
  `slot-targets` also accepts the per-slot extra creature of Freedom of Movement
  and Gaseous Form; `slot-duration` also accepts Geas' typed slot table and
  requires `scaling.geas-slot-duration`.
- Scaling kind `slot-spell-level` (Giant Insect) records that the slot's level is
  the spell's level in the stat block, and requires
  `scaling.stat-block-spell-level`. Its effect carries no dice.

Validation pins each record's header, scaling row, effects, material,
dependency list and typed values, including every stat block row, by equality
with the printed pages, so a positive but wrong value is rejected. Composition of
the bounded fixtures with #196, #218, #198, #199, #200, #201, #202 and #203 through #208
yields 168 unique records. Every imported record remains unavailable: no
resolver, action, entitlement, tool, inventory, coverage or Full Rules behavior
is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_207.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #208: Globe of Invulnerability through Hallow

`testdata/spells_208.json` imports the twelve families the issue-80 inventory
assigns to #208, on PDF pp.136-138: Globe of Invulnerability, Glyph of Warding,
Goodberry, Grease, Greater Invisibility, Greater Restoration, Guardian of Faith,
Guards and Wards, Guidance, Guiding Bolt, Gust of Wind and Hallow. Glyph of
Warding keeps its p.136 to p.137 continuation, Guards and Wards its p.137 to
p.138 continuation and Hallow its p.138 to p.139 continuation in explicit
`source_pages` lists. `spells_208_test.go` supplies literal source-value
fixtures for every header, variant, effect and scaling value, and checks the
family list and pages against the decision inventory.

The typed payloads retain Globe of Invulnerability's barrier and per-slot spell
level limit; Glyph of Warding's inscription, triggers, explosive rune and stored
spell; Goodberry's berries; Grease's square and save; Greater Invisibility's
condition; Greater Restoration's five removable effects; Guardian of Faith's
guardian, trigger and 60-damage ending; Guards and Wards' area, password, four
effects and 365-day extension; Guidance's d4 bonus; Guiding Bolt's attack,
damage and Advantage window; Gust of Wind's line, push and flame rules; and
Hallow's ward and all ten extra-effect choices.

Schema version 1 gains only shapes present in this batch, on top of #202:

- Casting time `24-hours` (Hallow), with `casting_seconds = 86400`.
- Scaling kind `slot-limit`: a spell-level limit that rises by a typed amount per
  slot level above the base (Globe of Invulnerability), alongside the existing
  `slot` dice scaling used by Glyph of Warding and Guiding Bolt.
- Target shapes `emanation`, `choice`, `square`, `created-effect` and `line`
  with the eligibilities this batch lists. Hallow's "area around" a touched
  point is `created-effect` with count 0, not `sphere`: the page never names a
  sphere. Glyph of Warding's
  `until-dispelled` duration records `dispelled-or-triggered` as its end timing.
- Flat quantities that are not dice, such as Guardian of Faith's 20 Radiant
  damage and Goodberry's 1 Hit Point, live in the typed payload; their effects
  carry no dice and no bonus.

Validation pins each record's header, scaling row, effects, material,
dependency list and typed values by equality with the printed page, so a
positive but wrong value is rejected. Composition of the bounded fixtures with
#196, #218, #198, #199, #200, #201, #202 and #203 through #208 yield 168 unique records.
Every imported record remains unavailable: no resolver, action, entitlement,
tool, inventory, coverage or Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_208.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #209: Hallucinatory Terrain through Hold Monster

`testdata/spells_209.json` imports the twelve families the issue-80 inventory
assigns to #209, on PDF pp.139-140: Hallucinatory Terrain, Harm, Haste, Heal,
Healing Word, Heat Metal, Hellish Rebuke, Heroes' Feast, Heroism, Hex, Hideous
Laughter and Hold Monster. Heat Metal keeps its p.139 to p.140 continuation and
Hold Monster its p.140 to p.141 continuation in explicit `source_pages` lists.
`spells_209_test.go` supplies literal source-value fixtures for every header,
option, effect and scaling value, and checks the family list and pages against
the decision inventory.

The typed payloads retain Hallucinatory Terrain's disguise and disbelief check;
Harm's Hit Point maximum reduction; Haste's benefits, allowed actions and
lethargy; Heal's flat healing and slot increment; Healing Word's dice; Heat
Metal's repeat damage and held-object save; Hellish Rebuke's trigger; the
benefits of Heroes' Feast; Heroism's Temporary Hit Points; Hex's chosen
ability, retargeting and slot duration table; and the conditions and repeat
saves of Hideous Laughter and Hold Monster.

Schema version 1 gains only shapes present in this batch:

- Target eligibilities `natural-terrain-in-cube` on a `cube`;
  `visible-willing-creature-in-range`,
  `visible-manufactured-metal-object-in-range` and
  `creature-that-damaged-caster` on a `single` target; and
  `surface-in-unoccupied-cube-next-to-caster` on a `created-effect`.
- `slot-bonus` also accepts Heal's typed flat healing increment, `slot-targets`
  the typed increments of Heroism, Hideous Laughter and Hold Monster, and
  `slot-duration` Hex's typed table.
- Hellish Rebuke keeps its Reaction trigger in `casting_triggers`. Heal's 70
  Hit Points and Haste's +2 Armor Class are effect `bonus` values.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong value
is rejected. The #209 verifier script checks identities, pages, printed headers,
material costs and anchor phrases only; typed payload values are covered by the
literal fixtures and pins, not by the script. A batch payload on a record of
another batch is rejected too. Composition of the bounded fixtures with #196,
#218, #198, #199, #200, #201, #202 and #203 through #214 yields 240 unique
records. Every imported record remains unavailable: no resolver, action,
entitlement, tool, inventory, coverage or Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_209.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #210: Hold Person through Insect Plague

`testdata/spells_210.json` imports the twelve families the issue-80 inventory
assigns to #210, on PDF pp.141-143: Hold Person, Holy Aura, Hunter's Mark,
Hypnotic Pattern, Ice Knife, Ice Storm, Identify, Illusory Script,
Imprisonment, Incendiary Cloud, Inflict Wounds and Insect Plague. Ice Storm
keeps its p.141 to p.142 continuation and Incendiary Cloud its p.142 to p.143
continuation in explicit `source_pages` lists. `spells_210_test.go` supplies
literal source-value fixtures for every header, option, effect and scaling
value, and checks the family list and pages against the decision inventory.

The typed payloads retain Hold Person's Humanoid-only targets; Holy Aura's
protections and retaliation; Hunter's Mark's tracking, retargeting and slot
duration table; Hypnotic Pattern's endings; Ice Knife's two damage steps; Ice
Storm's area and terrain; what Identify reveals; Illusory Script's readers and
alternatives; all five Imprisonment effects, its ending trigger and the level 9
Dispel Magic requirement; Incendiary Cloud's save triggers and drift; and the
damage areas of Inflict Wounds and Insect Plague.

Schema version 1 gains only shapes present in this batch:

- Target eligibilities `chosen-creatures-in-aura` on an `emanation`,
  `creatures-seeing-pattern-in-cube` on a `cube`,
  `creatures-in-cylinder-centered-in-range` on a `cylinder`,
  `object-or-creature-touched-throughout-casting` on a `choice` and
  `writing-material-touched` on a `single` target.
- `slot-targets` accepts Hold Person's typed increment and `slot-duration`
  accepts Hunter's Mark's typed table.
- Identify and Illusory Script are one-minute casts with `ritual = true`.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong value
is rejected. The #210 verifier script checks identities, pages, printed headers,
material costs and anchor phrases only; typed payload values are covered by the
literal fixtures and pins, not by the script. A batch payload on a record of
another batch is rejected too. Composition of the bounded fixtures with #196,
#218, #198, #199, #200, #201, #202 and #203 through #214 yields 240 unique
records. Every imported record remains unavailable: no resolver, action,
entitlement, tool, inventory, coverage or Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_210.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #211: Instant Summons through Locate Creature

`testdata/spells_211.json` imports the twelve families the issue-80 inventory
assigns to #211, on PDF pp.143-145: Instant Summons, Invisibility, Irresistible
Dance, Jump, Knock, Legend Lore, Lesser Restoration, Levitate, Light, Lightning
Bolt, Locate Animals or Plants and Locate Creature. Knock keeps its p.143 to
p.144 continuation in an explicit `source_pages` list. `spells_211_test.go`
supplies literal source-value fixtures for every header, option, effect and
scaling value, and checks the family list and pages against the decision
inventory.

The typed payloads retain Instant Summons' object limits and recall;
Invisibility's early endings; both Irresistible Dance outcomes; Jump's movement
exchange; Knock's lock rules; Legend Lore's results; the conditions Lesser
Restoration ends; Levitate's weight, movement and altitude rules; Light's
radii; Lightning Bolt's line; and the ranges and blockers of both locating
spells.

Schema version 1 gains only shapes present in this batch:

- Target eligibilities `sapphire-and-eligible-object-touch`,
  `visible-access-prevented-object-in-range`,
  `visible-creature-or-loose-object-in-range` and
  `unattended-large-or-smaller-object-touch` on a `single` target.
- `slot-targets` accepts the typed increments of Invisibility and Jump.
- Legend Lore's `cost_cp = 45000` is its whole component cost: 250 GP of
  consumed incense plus four 50 GP ivory strips that are not consumed. As with
  Clone, `consumed` is true because part is consumed; the payload keeps the
  split.
- Jump prints the singular `Component:` label and no Concentration.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong value
is rejected. The #211 verifier script checks identities, pages, printed headers,
material costs and anchor phrases only; typed payload values are covered by the
literal fixtures and pins, not by the script. A batch payload on a record of
another batch is rejected too. Composition of the bounded fixtures with #196,
#218, #198, #199, #200, #201, #202 and #203 through #214 yields 240 unique
records. Every imported record remains unavailable: no resolver, action,
entitlement, tool, inventory, coverage or Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_211.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #212: Locate Object through Mass Cure Wounds

`testdata/spells_212.json` imports the twelve families the issue-80 inventory
assigns to #212, on PDF pp.145-147: Locate Object, Longstrider, Mage Armor,
Mage Hand, Magic Circle, Magic Jar, Magic Missile, Magic Mouth, Magic Weapon,
Magnificent Mansion, Major Image and Mass Cure Wounds. Magic Circle keeps its
p.145 to p.146 continuation and Magic Weapon its p.146 to p.147 continuation in
explicit `source_pages` lists. `spells_212_test.go` supplies literal
source-value fixtures for every header, option, effect and scaling value, and
checks the family list and pages against the decision inventory.

The typed payloads retain Locate Object's range and blocker; Longstrider's
Speed increase; Mage Armor's base Armor Class; Mage Hand's uses and limits;
Magic Circle's creature types, protections, reverse direction and slot
duration; Magic Jar's possession, return and death rules; Magic Missile's
darts; Magic Mouth's message and trigger limits; Magic Weapon's bonus table;
the Magnificent Mansion and its servants; Major Image's illusion and the slot
level at which it lasts until dispelled; and Mass Cure Wounds' area and target
count.

Schema version 1 gains only shapes present in this batch:

- Target shape `darts` (Magic Missile, three darts), plus the eligibilities
  this batch lists on `single`, `cylinder`, `created-effect` and `up-to-count`
  targets.
- Scaling kind `slot-bonus-table` (Magic Weapon) requires its two typed rows
  and `scaling.bonus-by-slot`; `slot-until-dispelled` (Major Image) requires
  the typed minimum slot and `scaling.until-dispelled-by-slot`.
- `slot-count` also accepts Magic Missile's per-slot dart, `slot-duration`
  Magic Circle's typed increment and `slot-targets` Longstrider's increment.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong value
is rejected. The #212 verifier script checks identities, pages, printed headers,
material costs and anchor phrases only; typed payload values are covered by the
literal fixtures and pins, not by the script. A batch payload on a record of
another batch is rejected too. Composition of the bounded fixtures with #196,
#218, #198, #199, #200, #201, #202 and #203 through #214 yields 240 unique
records. Every imported record remains unavailable: no resolver, action,
entitlement, tool, inventory, coverage or Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_212.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #213: Mass Heal through Mirage Arcane

`testdata/spells_213.json` imports the twelve families the issue-80 inventory
assigns to #213, on PDF pp.147-150: Mass Heal, Mass Healing Word, Mass
Suggestion, Maze, Meld into Stone, Mending, Message, Meteor Swarm, Mind Blank,
Mind Spike, Minor Illusion and Mirage Arcane. Message keeps its p.148 to p.149
continuation and Mirage Arcane its p.149 to p.150 continuation in explicit
`source_pages` lists. `spells_213_test.go` supplies literal source-value
fixtures for every header, option, effect and scaling value, and checks the
family list and pages against the decision inventory.

The typed payloads retain Mass Heal's shared pool; Mass Healing Word's targets;
Mass Suggestion's limits and slot duration table; Maze's escape check; both
Meld into Stone expulsion cases; Mending's size limit; Message's barriers;
Meteor Swarm's four points and objects; Mind Blank's protections; Mind Spike's
tracking; both Minor Illusion options; and Mirage Arcane's terrain rules.

Schema version 1 gains only shapes present in this batch:

- Range kind `sight` (Mirage Arcane). It carries `feet = 0`, requires source
  details and lists `range.sight`.
- Target shape `spheres` (Meteor Swarm, four points), plus the eligibilities
  this batch lists on `up-to-count`, `single`, `choice` and `square` targets.
- `slot-duration` accepts Mass Suggestion's typed table. Its 366 days are
  31622400 seconds.
- Meteor Swarm's 1 mile is stored as 5280 feet. Mass Heal's 700-point pool
  lives in the payload because the caster divides it; its effects carry no
  bonus.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong value
is rejected. The #213 verifier script checks identities, pages, printed headers,
material costs and anchor phrases only; typed payload values are covered by the
literal fixtures and pins, not by the script. A batch payload on a record of
another batch is rejected too. Composition of the bounded fixtures with #196,
#218, #198, #199, #200, #201, #202 and #203 through #214 yields 240 unique
records. Every imported record remains unavailable: no resolver, action,
entitlement, tool, inventory, coverage or Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_213.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #214: Mirror Image through Phantom Steed

`testdata/spells_214.json` imports the twelve families the issue-80 inventory
assigns to #214, on PDF pp.150-152: Mirror Image, Mislead, Misty Step, Modify
Memory, Moonbeam, Move Earth, Nondetection, Pass without Trace, Passwall,
Phantasmal Force, Phantasmal Killer and Phantom Steed. Moonbeam keeps its p.150
to p.151 continuation and Phantasmal Force its p.151 to p.152 continuation in
explicit `source_pages` lists. The PDF prints Passwall before Pass without
Trace; the records follow the inventory order. `spells_214_test.go` supplies
literal source-value fixtures for every header, option, effect and scaling
value, and checks the family list and pages against the decision inventory.

The typed payloads retain Mirror Image's duplicate roll; Mislead's double;
Misty Step's distance; Modify Memory's options, interruptions and slot windows;
Moonbeam's area, save triggers and shape-shift rule; Move Earth's limits;
Nondetection's targets; Pass without Trace's bonus; Passwall's dimensions;
Phantasmal Force's belief and damage; Phantasmal Killer's repeat save; and
Phantom Steed's stat block changes.

Schema version 1 gains only shapes present in this batch:

- Scaling kind `slot-memory-window` (Modify Memory) requires its three typed
  rows, the any-time slot and `scaling.memory-window-by-slot`.
- Target eligibilities `terrain-area-in-range` on a `square`,
  `willing-creature-place-or-object-touch` on a `choice`,
  `caster-and-chosen-creatures-in-aura` on an `emanation` and
  `visible-point-on-wood-plaster-or-stone-surface` on a `created-effect`.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong value
is rejected. The #214 verifier script checks identities, pages, printed headers,
material costs and anchor phrases only; typed payload values are covered by the
literal fixtures and pins, not by the script. A batch payload on a record of
another batch is rejected too. Composition of the bounded fixtures with #196,
#218, #198, #199, #200, #201, #202 and #203 through #214 yields 240 unique
records. Every imported record remains unavailable: no resolver, action,
entitlement, tool, inventory, coverage or Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_214.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #215: Planar Ally through Prismatic Spray

`testdata/spells_215.json` imports the twelve families the issue-80 inventory
assigns to #215, on PDF pp.152-155: Planar Ally, Planar Binding, Plane Shift,
Plant Growth, Poison Spray, Polymorph, Power Word Heal, Power Word Kill, Power
Word Stun, Prayer of Healing, Prestidigitation and Prismatic Spray. Planar
Binding keeps its p.152 to p.153 continuation, Polymorph keeps its p.153 to
p.154 continuation and Prismatic Spray keeps its p.154 to p.155 continuation
in explicit `source_pages` lists. `spells_215_test.go` supplies literal
source-value fixtures for every header, option, effect and scaling value, and
checks the family list and pages against the decision inventory. Agents
transcribed the values from the pinned PDF's text; the fixtures, literals and
pins come from that one transcription, so their agreement is not independent
proof of it.

The typed payloads retain Planar Ally's payment rates and bargaining rules;
Planar Binding's four slot durations; Plane Shift's travelers and both
destination forms; Plant Growth's two casting times and their effects; Poison
Spray's cantrip rows; Polymorph's retained statistics and ending; the three
Power Words' thresholds and outcomes; Prayer of Healing's rest and healing;
Prestidigitation's six options; and all eight Prismatic Spray rows.

Plant Growth prints two casting times. The record's `casting_time` is `action`
and the payload keeps both the Overgrowth and the 8-hour Enrichment rows.
Enrichment's half-mile radius is stored as 2640 feet. Prestidigitation prints
"Up to 1 hour" without Concentration; it is a timed record whose `end_timing`
is `up-to-duration`.

Schema version 1 gains only shapes present in this batch:

- Existing scaling kind `slot-duration` also accepts Planar Binding's four typed
  rows with `scaling.duration-by-slot`; `cantrip-upgrade` also accepts Poison
  Spray's three typed rows.
- Target eligibilities `unoccupied-space-in-range-for-loyal-creature` on
  `created-effect`, `celestial-elemental-fey-or-fiend-in-range` on `single`,
  `caster-and-willing-creatures-linked-in-circle` on `selected-creatures`,
  `overgrowth-or-enrichment-point-in-range` and
  `one-chosen-minor-effect-in-range` on `choice`,
  `chosen-creatures-in-range-for-entire-casting` on `up-to-count` and
  `creatures-in-self-origin-cone` on `cone`.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong
value is rejected. The #215 verifier script checks identities, pages, printed
headers, material costs and anchor phrases only; typed payload values are
covered by the literal fixtures and pins, not by the script. A batch payload
on a record of another batch is rejected too. Composition of the bounded
fixtures with #196, #197, #218, #198, #199, #200, #201, #202, #203, #204,
#205, #206, #207, #208, #209, #210, #211, #212, #213, #214 and #215 yields 252
unique records. Every imported record remains unavailable: no resolver,
action, entitlement, tool, inventory, coverage or Full Rules behavior is
added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_215.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #216: Prismatic Wall through Ray of Frost

`testdata/spells_216.json` imports the twelve families the issue-80 inventory
assigns to #216, on PDF pp.155-157: Prismatic Wall, Private Sanctum, Produce
Flame, Programmed Illusion, Project Image, Protection from Energy, Protection
from Evil and Good, Protection from Poison, Purify Food and Drink, Raise Dead,
Ray of Enfeeblement and Ray of Frost. Project Image keeps its p.156 to p.157
continuation in an explicit `source_pages` list. `spells_216_test.go` supplies
literal source-value fixtures for every header, option, effect and scaling
value, and checks the family list and pages against the decision inventory.
Agents transcribed the values from the pinned PDF's text; the fixtures,
literals and pins come from that one transcription, so their agreement is not
independent proof of it.

The typed payloads retain Prismatic Wall's dimensions, light, seven layer rows
and layer destruction; Private Sanctum's cube sizes, six security properties
and permanence; Produce Flame's light and thrown attack; Programmed
Illusion's trigger and reset; Project Image's control and discovery; the three
Protection spells' benefits; Purify Food and Drink's area; Raise Dead's limits
and penalty; Ray of Enfeeblement's two save results; and Ray of Frost's slow
and cantrip rows.

Project Image's range of 500 miles is stored as 2640000 feet, and its
"up to 1 day" as 86400 seconds. Prismatic Wall rows with no damage or no
damage threshold store `none` and 0. Protection from Evil and Good prints
"Concentration up to 10 minutes" without the usual comma.

Schema version 1 gains only shapes present in this batch:

- Scaling kind `slot-cube-size` (Private Sanctum) requires its typed cube
  increment and `scaling.cube-size-per-slot`.
- Existing scaling kind `cantrip-upgrade` also accepts the typed rows of Produce
  Flame and Ray of Frost.
- Target eligibilities `wall-or-globe-centered-in-range` on `choice`,
  `warded-area-in-range` on `cube`, `triggered-illusion-in-range` and
  `previously-seen-location-in-range` on `created-effect`,
  `nonmagical-food-and-drink-in-sphere` on `sphere` and `dead-creature-touch`
  on `single`.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong
value is rejected. The #216 verifier script checks identities, pages, printed
headers, material costs and anchor phrases only; typed payload values are
covered by the literal fixtures and pins, not by the script. A batch payload
on a record of another batch is rejected too. Composition of the bounded
fixtures with #196, #197, #218, #198, #199, #200, #201, #202, #203, #204,
#205, #206, #207, #208, #209, #210, #211, #212, #213, #214, #215 and #216
yields 264 unique records. Every imported record remains unavailable: no
resolver, action, entitlement, tool, inventory, coverage or Full Rules
behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_216.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #217: Ray of Sickness through Sanctuary

`testdata/spells_217.json` imports the twelve families the issue-80 inventory
assigns to #217, on PDF pp.158-159: Ray of Sickness, Regenerate, Reincarnate,
Remove Curse, Resilient Sphere, Resistance, Resurrection, Reverse Gravity,
Revivify, Rope Trick, Sacred Flame and Sanctuary. Resurrection keeps its p.158
to p.159 continuation in an explicit `source_pages` list. `spells_217_test.go`
supplies literal source-value fixtures for every header, option, effect and
scaling value, and checks the family list and pages against the decision
inventory. Agents transcribed the values from the pinned PDF's text; the
fixtures, literals and pins come from that one transcription, so their
agreement is not independent proof of it.

The typed payloads retain Ray of Sickness's poison and slot die; Regenerate's
healing and regrowth; Reincarnate's ten-row species table; Remove Curse's
cursed-item rule; Resilient Sphere's barrier and movement; Resistance's eleven
damage types and reduction die; Resurrection's limits, penalty and long-dead
cost; Reverse Gravity's cylinder; Revivify's one-minute window; Rope Trick's
space; Sacred Flame's cover rule and cantrip rows; and Sanctuary's save and
endings.

The PDF prints Regenerate before Ray of Sickness; the records follow the
inventory order. "A century" is stored as 100 years and "the last minute" as
60 seconds. Resistance prints the singular "Component:" label.

Schema version 1 gains only shapes present in this batch:

- Existing scaling kind `cantrip-upgrade` also accepts Sacred Flame's typed rows.
- Target eligibilities `dead-humanoid-or-piece-of-one-touch`,
  `creature-or-object-touch`, `large-or-smaller-creature-or-object-in-range`
  and `rope-touch` on `single` and `creatures-and-objects-in-reversed-gravity`
  on `cylinder`.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong
value is rejected. The #217 verifier script checks identities, pages, printed
headers, material costs and anchor phrases only; typed payload values are
covered by the literal fixtures and pins, not by the script. A batch payload
on a record of another batch is rejected too. Composition of the bounded
fixtures with #196, #197, #218, #198, #199, #200, #201, #202, #203, #204,
#205, #206, #207, #208, #209, #210, #211, #212, #213, #214, #215, #216 and
#217 yields 276 unique records. Every imported record remains unavailable: no
resolver, action, entitlement, tool, inventory, coverage or Full Rules
behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_217.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #219: Shillelagh through Speak with Animals

`testdata/spells_219.json` imports the twelve families the issue-80 inventory
assigns to #219, on PDF pp.162-164: Shillelagh, Shining Smite, Shocking Grasp,
Silence, Silent Image, Simulacrum, Sleep, Sleet Storm, Slow, Sorcerous Burst,
Spare the Dying and Speak with Animals. Simulacrum keeps its p.162 to p.163
continuation in an explicit `source_pages` list. `spells_219_test.go` supplies
literal source-value fixtures for every header, option, effect and scaling
value, and checks the family list and pages against the decision inventory.
Agents transcribed the values from the pinned PDF's text; the fixtures,
literals and pins come from that one transcription, so their agreement is not
independent proof of it.

The typed payloads retain Shillelagh's weapon die table; Shining Smite's extra
damage, light and reveal; Shocking Grasp's restriction; Silence's sphere;
Silent Image's control and discovery; Simulacrum's creation, control and
repair; Sleep's two saves; Sleet Storm's cylinder; Slow's six penalties;
Sorcerous Burst's damage choices and extra dice; Spare the Dying's range
table; and Speak with Animals' information limits.

Simulacrum's 12-hour casting adds casting time `12-hours` with 43200 casting
seconds. Shining Smite keeps its two printed triggers and prints the singular
"Component:" label, as does Sorcerous Burst. Sorcerous Burst's damage effect
stores the damage type `chosen`; the seven choices are in the payload.

Schema version 1 gains only shapes present in this batch:

- Casting time `12-hours` (Simulacrum) requires typed source details.
- Cantrip scaling kind `cantrip-die-table` (Shillelagh) requires its three die
  rows and `scaling.cantrip-die-upgrade`.
- Cantrip scaling kind `cantrip-range` (Spare the Dying) requires its three
  range rows and `scaling.cantrip-range-upgrade`.
- Existing scaling kind `cantrip-upgrade` also accepts the typed rows of
  Shocking Grasp and Sorcerous Burst.
- Target eligibilities `held-club-or-quarterstaff`,
  `beast-or-humanoid-within-10-feet-for-entire-casting` and
  `creature-with-0-hit-points-not-dead` on `single`,
  `area-centered-on-chosen-point` and `chosen-creatures-in-sphere` on
  `sphere`, `visible-phenomenon-in-range` on `created-effect`,
  `area-centered-on-chosen-point` on `cylinder` and
  `chosen-creatures-in-cube-in-range` on `up-to-count`.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong
value is rejected. The #219 verifier script checks identities, pages, printed
headers, material costs and anchor phrases only; typed payload values are
covered by the literal fixtures and pins, not by the script. A batch payload
on a record of another batch is rejected too. Composition of the bounded
fixtures with #196, #197, #218, #198, #199, #200, #201, #202, #203, #204,
#205, #206, #207, #208, #209, #210, #211, #212, #213, #214, #215, #216, #217
and #219 yields 288 unique records. Every imported record remains unavailable:
no resolver, action, entitlement, tool, inventory, coverage or Full Rules
behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_219.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #220: Speak with Dead through Suggestion

`testdata/spells_220.json` imports the twelve families the issue-80 inventory
assigns to #220, on PDF pp.164-166: Speak with Dead, Speak with Plants, Spider
Climb, Spike Growth, Spirit Guardians, Spiritual Weapon, Starry Wisp, Stinking
Cloud, Stone Shape, Stoneskin, Storm of Vengeance and Suggestion. Spirit
Guardians keeps its p.164 to p.165 continuation and Storm of Vengeance keeps
its p.165 to p.166 continuation in explicit `source_pages` lists.
`spells_220_test.go` supplies literal source-value fixtures for every header,
option, effect and scaling value, and checks the family list and pages against
the decision inventory. Agents transcribed the values from the pinned PDF's
text; the fixtures, literals and pins come from that one transcription, so
their agreement is not independent proof of it.

The typed payloads retain Speak with Dead's question limits; Speak with Plants'
area and terrain choices; Spider Climb's movement and extra targets; Spike
Growth's damage interval and recognition check; Spirit Guardians' three
alignment rows; Spiritual Weapon's attack and movement; Starry Wisp's reveal;
Stinking Cloud's poison; Stone Shape's limits; Stoneskin's resistances; Storm
of Vengeance's initial strike and turn rows; and Suggestion's wording limit and
endings.

Spirit Guardians stores one Radiant and one Necrotic damage effect; the
payload's alignment rows say which applies, and the slot scaling names both.
Storm of Vengeance's range of 1 mile is stored as 5280 feet. "The past day" is
stored as 1 day.

Schema version 1 gains only shapes present in this batch:

- Existing scaling kind `slot-targets` also accepts Spider Climb; `cantrip-upgrade`
  also accepts Starry Wisp's typed rows.
- Target eligibilities `corpse-in-range`, `stone-object-or-section-touch` and
  `visible-hearing-understanding-creature-in-range` on `single`,
  `plants-in-immobile-emanation` and `undesignated-creatures-in-emanation` on
  `emanation`, `ground-in-sphere-on-point-in-range`,
  `area-centered-on-point-in-range` and
  `creatures-and-objects-under-storm-cloud` on `sphere` and
  `chosen-space-in-range` on `created-effect`.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong
value is rejected. The #220 verifier script checks identities, pages, printed
headers, material costs and anchor phrases only; typed payload values are
covered by the literal fixtures and pins, not by the script. A batch payload
on a record of another batch is rejected too. Composition of the bounded
fixtures with #196, #197, #218, #198, #199, #200, #201, #202, #203, #204,
#205, #206, #207, #208, #209, #210, #211, #212, #213, #214, #215, #216, #217,
#219 and #220 yields 300 unique records. Every imported record remains
unavailable: no resolver, action, entitlement, tool, inventory, coverage or
Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_220.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #221: Summon Dragon through Tiny Hut

`testdata/spells_221.json` imports the twelve families the issue-80 inventory
assigns to #221, on PDF pp.166-169: Summon Dragon, Sunbeam, Sunburst, Symbol,
Telekinesis, Telepathic Bond, Teleport, Teleportation Circle, Thaumaturgy,
Thunderwave, Time Stop and Tiny Hut. Telekinesis keeps its p.167 to p.168
continuation and Teleportation Circle keeps its p.168 to p.169 continuation in
explicit `source_pages` lists. `spells_221_test.go` supplies literal
source-value fixtures for every header, option, effect and scaling value, and
checks the family list and pages against the decision inventory. Agents
transcribed the values from the pinned PDF's text; the fixtures, literals and
pins come from that one transcription, so their agreement is not independent
proof of it.

The typed payloads retain Summon Dragon's Draconic Spirit stat block; Sunbeam's
line and light; Sunburst's blindness; Symbol's triggers and six glyph effects;
Telekinesis's creature and object options; Telepathic Bond's link; Teleport's
outcome table, familiarity meanings and results; Teleportation Circle's sigil
and permanence rules; Thaumaturgy's six options; Thunderwave's push; Time
Stop's turns and endings; and Tiny Hut's barrier.

The pinned PDF ends Telekinesis's fine-control sentence on p.168 at "such as
manipulating a simple tool," with no further text. The record keeps that one
example and notes the printed ending; nothing is supplied from another source.
In Teleport's table, printed `00` is stored as 100 and a dash is an absent
outcome row. Teleportation Circle's "1 round" is stored as 6 seconds.
Thaumaturgy prints "Up to 1 minute" without Concentration; it is a timed record
whose `end_timing` is `up-to-duration`. Symbol prints "Until dispelled or
triggered"; the duration is `until-dispelled` and `end_timing` names the
trigger. The Draconic Spirit's ability modifiers and saves reuse
`SpellStatAbility`.

Schema version 1 gains only shapes present in this batch:

- Scaling kind `slot-stat-block` (Summon Dragon) requires the typed slot-level
  rule and `scaling.stat-block-by-slot`.
- Target eligibilities `creatures-in-self-origin-line` on `line`,
  `creatures-in-sphere-on-point-in-range` on `sphere`,
  `surface-or-closable-object-touch`, `visible-creature-or-object-in-range`,
  `willing-creatures-or-single-object-in-range` and
  `one-chosen-minor-wonder-in-range` on `choice`,
  `chosen-willing-creatures-in-range` on `up-to-count`,
  `circle-on-ground-in-range` on `created-effect`,
  `creatures-and-unsecured-objects-in-self-origin-cube` on `cube` and
  `stationary-emanation-around-caster` on `emanation`.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong
value is rejected. The #221 verifier script checks identities, pages, printed
headers, material costs and anchor phrases only; typed payload values are
covered by the literal fixtures and pins, not by the script. A batch payload
on a record of another batch is rejected too. Composition of the bounded
fixtures with #196, #197, #218, #198, #199, #200, #201, #202, #203, #204,
#205, #206, #207, #208, #209, #210, #211, #212, #213, #214, #215, #216, #217,
#219, #220 and #221 yields 312 unique records. Every imported record remains
unavailable: no resolver, action, entitlement, tool, inventory, coverage or
Full Rules behavior is added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_221.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline literal
fixtures and is not an automatic conformance proof.

## Issue #222: Tongues through Vitriolic Sphere

`testdata/spells_222.json` imports the twelve families the issue-80 inventory
assigns to #222, on PDF pp.169-172: Tongues, Transport via Plants, Tree
Stride, True Polymorph, True Resurrection, True Seeing, True Strike, Tsunami,
Unseen Servant, Vampiric Touch, Vicious Mockery and Vitriolic Sphere. Tongues
keeps its p.169 to p.170 continuation, True Resurrection keeps its p.170 to
p.171 continuation and Vicious Mockery keeps its p.171 to p.172 continuation
in explicit `source_pages` lists. `spells_222_test.go` supplies literal
source-value fixtures for every header, option, effect and scaling value, and
checks the family list and pages against the decision inventory. Agents
transcribed the values from the pinned PDF's text; the fixtures, literals and
pins come from that one transcription, so their agreement is not independent
proof of it.

The typed payloads retain True Polymorph's three transformation forms, its
retained statistics and its full-Concentration result; True Resurrection's
limits, restoration and replacement body; True Strike's weapon attack and
extra Radiant dice; Tsunami's wall, movement, damage reduction and swimming;
Unseen Servant's statistics, commands, example tasks and endings; Vampiric
Touch's healing and repeat attack; and Vitriolic Sphere's initial and delayed
damage.

Tsunami's 1-mile range is stored as 5280 feet and its 6 rounds as 36 seconds.
Maintaining Concentration on True Polymorph for the full hour makes it last
until dispelled: the record is a timed Concentration record and the payload
and `end_timing` keep that result. True Strike's weapon cost of 1+ CP is
stored as 1 copper piece.

Schema version 1 gains only shapes present in this batch:

- Existing scaling kind `cantrip-upgrade` also accepts the three typed rows of
  True Strike and Vicious Mockery. True Strike's rows are its extra Radiant
  dice and its scaled effect is `extra`, which has no dice below level 5.
- Target eligibilities `large-or-larger-inanimate-plant-in-range`,
  `creature-dead-no-longer-than-200-years-touch` and
  `visible-or-heard-creature-in-range` on `single`,
  `visible-creature-or-nonmagical-object-in-range` on `choice`,
  `creatures-in-water-wall-on-point-in-range` on `wall`,
  `unoccupied-ground-space-in-range` on `created-effect` and
  `creatures-in-acid-sphere-on-point-in-range` on `sphere`.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong
value is rejected. The #222 verifier script checks identities, pages, printed
headers, material costs and anchor phrases only; typed payload values are
covered by the literal fixtures and pins, not by the script. A batch payload
on a record of another batch is rejected too. Composition of the bounded
fixtures with #196, #197 and every batch from #198 through #222 yields 324
unique records. Every imported record remains unavailable: no resolver,
action, entitlement, tool, inventory, coverage or Full Rules behavior is
added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_222.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline
literal fixtures and is not an automatic conformance proof.

## Issue #223: Wall of Fire through Wind Wall

`testdata/spells_223.json` imports the twelve families the issue-80 inventory
assigns to #223, on PDF pp.172-174: Wall of Fire, Wall of Force, Wall of Ice,
Wall of Stone, Wall of Thorns, Warding Bond, Water Breathing, Water Walk, Web,
Weird, Wind Walk and Wind Wall. Wall of Ice keeps its p.172 to p.173
continuation and Water Walk keeps its p.173 to p.174 continuation in explicit
`source_pages` lists. `spells_223_test.go` supplies literal source-value
fixtures for every header, option, effect and scaling value, and checks the
family list and pages against the decision inventory. Agents transcribed the
values from the pinned PDF's text; the fixtures, literals and pins come from
that one transcription, so their agreement is not independent proof of it.

The typed payloads retain each wall's forms, dimensions, object statistics,
saves and damage; Warding Bond's rings, benefits, shared damage and endings;
Water Walk's example surfaces and transitions; Web's anchoring, escape check
and burning; Weird's repeated save; Wind Walk's cloud form, reversion and
descent; and Wind Wall's blocked and unaffected projectiles.

Wall of Force's 1/4-inch thickness is stored as one quarter-inch. Wall of
Stone becomes permanent after full Concentration: the record is a timed
Concentration record and the payload and `end_timing` keep that result.
Warding Bond's rings cost 50+ GP each, so the material cost is 10000 copper
pieces for the pair.

Schema version 1 gains only shapes present in this batch:

- Existing scaling kind `slot-effect-dice` also accepts Wall of Ice's two
  typed rows in `source_details.wall_of_ice.slot_dice`, for the appearance
  damage and the frigid air. Wall of Fire and Wall of Thorns use `slot` on
  both of their damage effects.
- Target eligibilities `creatures-in-wall-of-fire`,
  `force-wall-on-point-in-range`, `creatures-in-ice-wall`,
  `stone-wall-on-point-in-range`, `creatures-in-wall-of-thorns` and
  `creatures-in-wind-wall` on `wall`, `other-willing-creature-touch` on
  `single`, `chosen-willing-creatures-in-range` for ten creatures on
  `up-to-count` and `self-and-up-to-count`, `webs-in-cube-on-point-in-range`
  on `cube` and `chosen-creatures-in-sphere-on-point-in-range` on
  `selected-creatures`.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong
value is rejected. The #223 verifier script checks identities, pages, printed
headers, material costs and anchor phrases only; typed payload values are
covered by the literal fixtures and pins, not by the script. A batch payload
on a record of another batch is rejected too. Composition of the bounded
fixtures with #196, #197 and every batch from #198 through #223 yields 336
unique records. Every imported record remains unavailable: no resolver,
action, entitlement, tool, inventory, coverage or Full Rules behavior is
added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_223.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline
literal fixtures and is not an automatic conformance proof.

## Issue #224: Wish through Zone of Truth

`testdata/spells_224.json` imports the three families the issue-80 inventory
assigns to #224, on PDF p.175: Wish, Word of Recall and Zone of Truth.
`spells_224_test.go` supplies literal source-value fixtures for every header,
option, effect and scaling value, and checks the family list and pages against
the decision inventory. Agents transcribed the values from the pinned PDF's
text; the fixtures, literals and pins come from that one transcription, so
their agreement is not independent proof of it.

The typed payloads retain Wish's spell duplication, all seven listed
alternative effects (Object Creation, Instant Health, Resistance, Spell
Immunity, Sudden Learning, Roll Redo and Reshape Reality) and its casting
stress; Word of Recall's sanctuary rules; and Zone of Truth's save and limits.
Wish's 25,000 GP object limit is stored as 2500000 copper pieces. Its stress
effect records 1d10 Necrotic damage; the payload says that this is per level
of each spell cast.

Schema version 1 gains only target eligibilities: `one-chosen-wish-use` on
`choice`, `willing-creatures-within-range` for five creatures on
`self-and-up-to-count` and `creatures-in-truth-zone-on-point-in-range` on
`sphere`.

This is the last spell batch. `spells_224_test.go` also checks that the
composed source catalogue holds exactly the 339 spell families, names and
pages in the issue-80 inventory. It is a check of identities and pages, not of
conformance.

Validation pins each record's whole printed value, including header, scaling
row, effects, required mechanics and typed payload, so a positive but wrong
value is rejected. The #224 verifier script checks identities, pages, printed
headers, material costs and anchor phrases only; typed payload values are
covered by the literal fixtures and pins, not by the script. A batch payload
on a record of another batch is rejected too. Composition of the bounded
fixtures with #196, #197 and every batch from #198 through #224 yields 339
unique records. Every imported record remains unavailable: no resolver,
action, entitlement, tool, inventory, coverage or Full Rules behavior is
added.

To recheck the PDF pin, identities, pages, printed headers and source anchors,
run:

```sh
rtk proxy python3 scripts/verify_srd_spells_224.py /path/to/SRD_CC_v5.2.1.pdf
```

The same limits as #198 apply: the source check supplements the offline
literal fixtures and is not an automatic conformance proof.
