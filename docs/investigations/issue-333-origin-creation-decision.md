# Origin and creation choices (#333)

Class choice, class equipment, level advancement, higher-level starts, multiclassing and the four Origin feats already have published owners, and this decision reuses them unchanged. No published issue executes a background other than the two fixed builds, any species trait, ability-score generation, the language and alignment choices, trinkets, GM-authored backgrounds, the player's Passive Perception or typed damage Resistance. Those residuals are assigned to 22 bounded tasks, each one PR, each a native child of #333 and a direct blocker of character qualification #337. Twenty implement a rule and two are qualification-only outputs that keep #430 and #488 under GitHub's limit of 100 outgoing blocking relations. One of the twenty, `drow-faerie-fire`, stays open until the open spell splitter #277 publishes the Faerie Fire effect. Closing this decision supplies no runtime.

The tasks are drafted here as exact body files. Their issue numbers, native links and readback are recorded in the link-verification receipt when they are published; until then the scope status is `SOURCE_RECONCILED_PUBLICATION_PENDING`.

## Source reconciliation

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Retain CC-BY-4.0 attribution and fixture adaptation in [the source register](../srd-5.2.1-conformance.md). Both blockers of #333, the background catalogue #306 and the species catalogue #307, are closed.

Three page corrections against the issue's cited range (pp. 18–27, 83–86, 192):

- PDF18 holds only the end of Playing the Game (Temporary Hit Points). Character Creation starts on PDF19.
- Creating a Background starts on PDF192; steps 2 to 5 are on PDF193.
- Trinkets occupy PDF26–27, rows 01 to 00.

| Scope | PDF | Actual ownership |
|---|---|---|
| Class choice, level 1 class features, class proficiencies, class HP and Hit Point Dice | 19–20, 22–23 | The twelve class qualifications #498/#501/#551/#630/#728/#776/#845/#970/#1093/#1475/#1758/#2036 and their own first-class tasks. |
| Class starting equipment | 20 | One task per class: #433/#508/#554/#635/#730/#785/#853/#982/#1101/#1177/#1480/#1762. |
| Background choice and grants | 19–20, 83, 94 | New `backgrounds`. Feat selection #384; effects #385 Alert, #386 Magic Initiate, #387 Savage Attacker; provenance #430. |
| Spending creation coins | 20 | Price owners under #275. The equipment tasks record coins only. |
| Languages | 20 | New `languages-alignment`, which owns the one record of the Standard and Rare Languages tables. #987 and #1105 consume that record; #640 grants Druidic. |
| Ability scores | 20–21, 83 | New `ability-scores`: three methods, assignment, background adjustment, modifier table. Later increases stay with #384/#389. |
| Alignment | 21–22 | New `languages-alignment`, as a descriptive record. |
| Saves, skills, HP, Initiative, AC, attacks, spell numbers | 22–23 | Existing derived sheet (`internal/rules/srd521/character.go`, #61/#65/#68) and the class tasks. |
| Passive Perception | 22, 186 | New `passive-perception`. The derived sheet does not carry it and no issue names it for player characters. |
| Advancement, bonus feats, higher-level starts | 19, 23–24 | Decision #336: #285, #429, #968, #2053, #2054. |
| Multiclassing | 24–26 | Decision #335: #2056. |
| Trinkets | 26–27 | New `trinket`. #2053 already names #333 as the owner. |
| Nine species with nested choices | 83–86 | Thirteen new tasks: one each for seven species, three for Elf and three for Goliath. |
| Typed damage Resistance | 5, 17, 187 | New `typed-resistance`, the one shared owner. Dragonborn, Dwarf and Tiefling register sources and consume it. |
| Origin feat effects | 87 | Decision #334: #384–#388. |
| Creating a Background | 19, 83, 192–193 | New `custom-background`, built on `ability-scores` and `backgrounds`. |
| Creation provenance and duplicate proficiencies | 19–23 | Existing #430, already a child of #333, unchanged; consumed through new `creation-output`. |
| Charmed, Frightened and Prone | 178, 182, 186 | Existing #488, unchanged; consumed through new `condition-output`. |

These passages carry no rule to execute and need no owner: the character sheet advice, the Class Overview's likes and complexity columns, Imagine Your Past and Present and the Tiers of Play. Two more are recorded as dispositions so that nobody validates against them. Species age (PDF83) is descriptive and unvalidated: the text allows any age within a species' normal life span, and the SRD species entries give no life spans. The Ability Scores and Backgrounds table (PDF19–20) lists only Acolyte and Soldier; it is a suggestion and never restricts a background or ability choice, which follow the four background records on PDF83. The Standard Array by Class table is likewise a suggestion and is recorded as such in `ability-scores`.

Current runtime evidence: `internal/rules/srd521/creation.go` admits exactly Human Soldier Fighter and Human Acolyte Cleric, and the only mention of Resistance in non-test runtime code is a comment in `attack.go`. A search of all 1,892 issue titles and bodies, open and closed, found no owner for any species trait name, for Criminal or Sage, for the three generation methods, the language tables, trinkets, authored backgrounds or a typed Resistance projection. The companion scope records each clause, its checked source phrases and each reused owner with its body hash.

## New tasks

Backticked prerequisites are other tasks of this decision.

| Key | Title | Size | Hard prerequisites |
|---|---|---|---|
| `creation-output` | Qualify creation provenance output for origin consumers | S | #430 |
| `condition-output` | Qualify Charmed, Frightened and Prone outputs for origin consumers | S | #488 |
| `typed-resistance` | Apply typed damage Resistance once per damage instance | S | #65 |
| `trinket` | Confirm one optional creation trinket roll from the exact d100 table | S | #195, #284 |
| `passive-perception` | Project player Passive Perception on the derived sheet | S | #61, #284 |
| `ability-scores` | Confirm creation ability scores by Standard Array, 4d6 or 27-point cost | M | `creation-output`, #284, #306 |
| `backgrounds` | Confirm Acolyte, Criminal, Sage or Soldier grants and equipment choice | M | `creation-output`, #195, #247, #249, #253, #254, #256–#261, #306, #384–#387 |
| `languages-alignment` | Record creation languages and alignment from the exact source tables | S | `creation-output`, #284 |
| `dragonborn` | Confirm Dragonborn ancestry, Breath Weapon, resistance and level5 flight | M | `creation-output`, `typed-resistance`, #61, #65, #66, #145, #307, #428, #487, #968 |
| `dwarf` | Confirm Dwarf resilience, per-level toughness and Stonecunning | M | `creation-output`, `typed-resistance`, #61, #66, #145, #307, #428, #504, #727, #968 |
| `elf` | Confirm Elf lineage choice, cantrips, Fey Ancestry, Keen Senses and Trance | M | `creation-output`, `condition-output`, #61, #145, #307, #428, #1049, #1050, #1281, #1384–#1386 |
| `gnome` | Confirm Gnome cunning, lineage spells and Rock Gnome clockwork devices | M | `creation-output`, #61, #66, #145, #307, #428, #700, #955, #1375, #1376, #1384–#1386, #1448 |
| `goliath` | Confirm Goliath Powerful Build, Giant Ancestry choice and on-hit boons | M | `creation-output`, `condition-output`, #61, #65, #66, #307, #428, #778 |
| `halfling` | Confirm Halfling Brave, Luck, Nimbleness and Naturally Stealthy | M | `creation-output`, `condition-output`, #61, #307, #1137 |
| `human` | Confirm Human size, Resourceful, Skillful and Versatile choices | S | `creation-output`, #307, #384–#388, #428, #774 |
| `orc` | Confirm Orc Adrenaline Rush and Relentless Endurance | M | `creation-output`, #65, #66, #145, #307, #427, #428, #487, #780 |
| `tiefling` | Confirm Tiefling size, Fiendish Legacy spells and Otherworldly Presence | M | `creation-output`, `typed-resistance`, #145, #307, #428, #678, #683, #695, #959, #968, #1075, #1264, #1282, #1318, #1448, #1607, #1631, #1632 |
| `custom-background` | Confirm one GM-authored background under the five source steps | M | `creation-output`, `ability-scores`, `backgrounds`, #247–#262, #301, #306, #384–#388 |
| `elf-lineage-spells` | Confirm Elf level 3 and 5 lineage spell uses and High Elf cantrip replacement | M | `elf`, #265, #428, #688, #906, #968, #1066, #1072, #1282, #1448 |
| `goliath-reactions` | Resolve Goliath Cloud's Jaunt, Stone's Endurance and Storm's Thunder | M | `goliath`, #66, #81 |
| `goliath-large-form` | Confirm Goliath Large Form size change at character level 5 | M | `goliath`, #61, #66, #968 |
| `drow-faerie-fire` | Enable the Drow level 3 Faerie Fire lineage cast | S | `elf-lineage-spells`, #277; pending: the Faerie Fire owner from #277 |

One species per task is the default boundary, because a reviewer can check it against a single source entry. Elf and Goliath are each split in three so that no task carries more than one nested option family with its own resources: `elf` holds the identity, the lineage choice and the passive traits, `elf-lineage-spells` the six level 3 and 5 spells and the High Elf cantrip replacement, and `drow-faerie-fire` the one spell with no effect owner; `goliath` holds the ancestry choice, the three on-hit boons and Powerful Build, `goliath-reactions` the teleport and the two Reactions, and `goliath-large-form` the size change. Species tasks own eligibility, nested choices, free uses and recovery; every spell they grant executes only through its existing effect owner, and slot versus free-use routing stays with #1448, as decision #273 requires.

Ownership points a reviewer should check:

- Faerie Fire (Drow, level 3; Level 1 spell, PDF129) has no published effect owner. It appears only in the catalogue issue #205, and its key is still assigned by the open spell splitter #277. `drow-faerie-fire` therefore owns that one grant. It is blocked by #277 and by `elf-lineage-spells`, it blocks #337 directly, and its body keeps it open until the effect owner that #277 publishes exists, is linked as its blocker and is merged, and the cast executes through it. That future edge is recorded in the scope file as a pending blocker for the publisher of #277. No other Elf task waits on #277. #277 is a planning issue and is itself blocked by 30 open catalogue issues, so only this one task is held.
- The High Elf may replace its cantrip with any of the fifteen Wizard cantrips on PDF79. `elf-lineage-spells` owns the replacement choice and offers a cantrip only while that cantrip's own effect owner is qualified. It cannot close on an empty offer, because Dancing Lights #1281 is required through `elf`.
- Typed Resistance (PDF17, 187) is a shared damage rule with no published owner: #435 and other features only say they consume it. `typed-resistance` is its single owner and blocks #337; `dragonborn`, `dwarf` and `tiefling` are blocked by it and register a named source without halving anything themselves. If an existing consumer such as #435 lands first, `typed-resistance` adopts that one projection.
- The Standard and Rare Languages tables have one record, owned by `languages-alignment`. #987 and #1105 consume it but no existing edge orders them after this task, and existing edges are not changed here. If either lands first, `languages-alignment` adopts that record and adds no second table.
- #774 owns the Heroic Inspiration pool and says the Human Resourceful gain remains with #333. `human` owns only that gain.
- A fall when Draconic Flight ends in the air has no owner: #147 resolves one declared fall and excludes flying. `dragonborn` ends the Fly Speed, records the fact and computes no fall.
- Some general rules have no published shared owner and are consumed in a bounded way inside the single task that needs them: Tremorsense (PDF190) in `dwarf`, Carrying Capacity for the character's own sheet (PDF178) and a timed Speed reduction in `goliath`, teleportation (PDF190) in `goliath-reactions`, Temporary Hit Points and Massive Damage (PDF17–18) in `orc`, and a Fly Speed for the flight's duration (PDF182) in `dragonborn`. None is widened into a general framework.

Where the source is silent, the task says so and either refuses or records an Edda decision as one: a repeated language roll waits for an explicit new roll or choice; a fourth clockwork device is rejected while three exist; an underspent point-cost set is accepted without converting the unspent points; a Halfling Luck reroll that is also a 1 stands and is not rerolled again; the High Elf holds one replaceable cantrip at a time; and Large Form changes no equipment and gains no forced ending beyond the three the source names.

## Native links and caps

Each task has parent #333 and blocks #337 directly, which is what the issue's third criterion asks. GitHub limits an issue to 100 blocked-by relations, 100 outgoing blocking relations and a parent to 100 sub-issues. All three were checked for every issue this plan touches, from a live read of all 1,892 issues.

Two reused owners are too close to the outgoing limit to take one edge per consumer, so each is consumed through one qualification-only output, the pattern #1738 uses for Warlock:

- #430 has 98 outgoing relations and thirteen origin tasks need it. `creation-output` is blocked by #430 and blocks those thirteen and #337.
- #488 has 97 and three tasks need it. `condition-output` is blocked by #488 and blocks `elf`, `goliath` and `halfling` and #337. Three direct edges would fit but would fill the cap exactly.

#67 already has 100 outgoing relations, so no task takes an edge to it. It is closed, and every task that reads its scene facts requires it transitively through another prerequisite.

| Native list | Before | After | Cap |
|---|---:|---:|---:|
| #337 blocked-by | 62 | 84 | 100 |
| #333 sub-issues | 1 | 23 | 100 |
| #283 sub-issues | 95 | 95 | 100 |
| Largest new blocked-by list (`custom-background`) | — | 26 | 100 |
| Largest new outgoing list (`creation-output`) | — | 14 | 100 |
| #430 outgoing | 98 | 99 | 100 |
| #488 outgoing | 97 | 98 | 100 |
| #67 outgoing | 100 | 100 | 100 |
| #61 outgoing | 84 | 92 | 100 |
| #428 outgoing | 81 | 90 | 100 |
| #487 outgoing | 87 | 89 | 100 |
| #65 outgoing | 73 | 77 | 100 |
| #384 outgoing | 67 | 70 | 100 |
| #66 outgoing | 62 | 69 | 100 |

The scope file lists the before and after outgoing count of all 71 reused prerequisites; apart from #430 and #488 none passes 92. Other splitters publish in parallel, so the publisher rereads these counts immediately before creating edges.

No new task is a child of #283. Publication order follows the new-to-new edges: the two outputs and `typed-resistance` first, then the tasks that depend only on them, then `custom-background`, `elf-lineage-spells`, the two Goliath follow-ups and last `drow-faerie-fire`. The dependency graph was read in full (1,892 issues, 10,201 blocked-by edges): no planned prerequisite is transitively blocked by #337, the new-to-new edges are acyclic, and so the 22 new edges into #337 close no cycle. The parent chain stays #333 → #283 → #54 → #47. No existing body, parent or edge changes; #430 stays a child of #333 and keeps its own edges.

## Deterministic witnesses

Soldier +2 STR/+1 CON on an assigned 15/13 gives 17 and 14. Point costs 15,15,15,8,8,8 and 15,14,13,12,10,8 both total 27. A Dwarf Fighter 1 with CON +2 has 13 Hit Points and 22 after the fixed level 2 application. A level 5 Dragonborn with CON +2 has breath DC 13 and 2d10; injected 7 and 4 deal 11 or 5, and 5 or 2 to a resistant target. 28 Fire damage reduced by 5 and then halved is 11, the source's own example. A Goliath with Strength 16 carries 480 lb., and 960 lb. in Large Form. Stone's Endurance rolling 7 with CON +2 turns 12 damage into 3. An Orc with maximum 12 and current 6 Hit Points may drop to 1 against 17 damage and dies against 18. Wisdom 15 with Perception proficiency at level 1 gives Passive Perception 14, the source's own example.

## Verification and limits

118 source phrases across PDF5, 8, 14, 17–27, 79, 83–87, 94, 129, 178–179, 182–184, 186–187, 190, 192–193 were compared with the pinned PDF text. The spell levels the lineage tasks rely on were read from the spell descriptions: Detect Magic, Longstrider, Faerie Fire and Hellish Rebuke are level 1; Misty Step, Darkness and Pass without Trace are level 2. Task body hashes, owner titles, states and body hashes read live, cap counts and the cycle check are in the scope file.

No runtime file, catalogue, migration, UI or existing issue is changed. Runtime implementation remains paused. Source checks and hosted CI qualify this decision head; PostgreSQL, browser/device, production and Full Rules qualification belong to the later implementation PRs.
