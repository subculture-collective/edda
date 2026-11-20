# Origin and creation choices (#333)

Class choice, class equipment, level advancement, higher-level starts, multiclassing and the four Origin feats already have published owners, and this decision reuses them unchanged. No published issue executes a background other than the two fixed builds, any species trait, ability-score generation, the language and alignment choices, trinkets, GM-authored backgrounds or the player's Passive Perception. Those residuals are assigned to fifteen bounded tasks, each one implementation PR, each a native child of #333 and a direct blocker of character qualification #337. Closing this decision supplies no runtime.

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
| Languages | 20 | New `languages-alignment`. Rare languages only through #640, #987 and #1105. |
| Ability scores | 20–21, 83 | New `ability-scores`: three methods, assignment, background adjustment, modifier table. Later increases stay with #384/#389. |
| Alignment | 21–22 | New `languages-alignment`, as a descriptive record. |
| Saves, skills, HP, Initiative, AC, attacks, spell numbers | 22–23 | Existing derived sheet (`internal/rules/srd521/character.go`, #61/#65/#68) and the class tasks. |
| Passive Perception | 22, 186 | New `passive-perception`. The derived sheet does not carry it and no issue names it for player characters. |
| Advancement, bonus feats, higher-level starts | 19, 23–24 | Decision #336: #285, #429, #968, #2053, #2054. |
| Multiclassing | 24–26 | Decision #335: #2056. |
| Trinkets | 26–27 | New `trinket`. #2053 already names #333 as the owner. |
| Nine species with nested choices | 83–86 | Nine new tasks, one per species. |
| Origin feat effects | 87 | Decision #334: #384–#388. |
| Creating a Background | 19, 83, 192–193 | New `custom-background`, built on `ability-scores` and `backgrounds`. |
| Creation provenance and duplicate proficiencies | 19–23 | Existing #430, already a child of #333, unchanged. |

The character sheet advice, the Class Overview's likes and complexity columns, Imagine Your Past and Present and the Tiers of Play carry no rule to execute and need no owner. The Standard Array by Class table is a suggestion and is recorded as such in `ability-scores`.

Current runtime evidence: `internal/rules/srd521/creation.go` admits exactly Human Soldier Fighter and Human Acolyte Cleric. A search of all 1,892 issue titles and bodies, open and closed, found no owner for any species trait name, for Criminal or Sage, for the three generation methods, the language tables, trinkets or authored backgrounds. The companion scope records each clause, its checked source phrases and each reused owner with its body hash.

## New tasks

| Key | Title | Size | Hard prerequisites |
|---|---|---|---|
| `ability-scores` | Confirm creation ability scores by Standard Array, 4d6 or 27-point cost | M | #284, #306, #430 |
| `backgrounds` | Confirm Acolyte, Criminal, Sage or Soldier grants and equipment choice | M | #195, #247, #249, #253, #254, #256–#261, #306, #384–#387, #430 |
| `languages-alignment` | Record creation languages and alignment from the exact source tables | S | #284, #430 |
| `trinket` | Confirm one optional creation trinket roll from the exact d100 table | S | #195, #284 |
| `passive-perception` | Project player Passive Perception on the derived sheet | S | #61, #284 |
| `custom-background` | Confirm one GM-authored background under the five source steps | M | `ability-scores`, `backgrounds`, #301, #306, #384–#388, #430 |
| `dragonborn` | Confirm Dragonborn ancestry, Breath Weapon, resistance and level5 flight | M | #61, #65, #66, #145, #147, #307, #428, #430, #487, #968 |
| `dwarf` | Confirm Dwarf resilience, per-level toughness and Stonecunning | M | #61, #66, #145, #307, #428, #430, #504, #727, #968 |
| `elf` | Confirm Elf lineage spells, Keen Senses, Fey Ancestry and Trance | L | #61, #145, #265, #277, #307, #428, #430, #488, #688, #906, #968, #1049, #1050, #1066, #1072, #1281, #1282, #1384–#1386, #1448 |
| `gnome` | Confirm Gnome cunning, lineage spells and Rock Gnome clockwork devices | M | #61, #66, #145, #307, #428, #430, #700, #955, #1375, #1376, #1384–#1386, #1448 |
| `goliath` | Confirm Goliath Giant Ancestry options, Large Form and Powerful Build | L | #61, #65, #66, #81, #307, #428, #430, #488, #778, #968 |
| `halfling` | Confirm Halfling Brave, Luck, Nimbleness and Naturally Stealthy | M | #61, #67, #307, #430, #488, #1137 |
| `human` | Confirm Human size, Resourceful, Skillful and Versatile choices | S | #307, #384–#388, #428, #430, #774 |
| `orc` | Confirm Orc Adrenaline Rush and Relentless Endurance | M | #66, #145, #307, #427, #428, #430, #780 |
| `tiefling` | Confirm Tiefling size, Fiendish Legacy spells and Otherworldly Presence | M | #145, #307, #428, #430, #678, #683, #695, #959, #968, #1075, #1264, #1282, #1318, #1448, #1607, #1631, #1632 |

One species per task is the smallest boundary a reviewer can check against a single source entry. Elf and Goliath are Size L because each carries three or six nested options with separate resources. Species tasks own eligibility, nested choices, free uses and recovery; every spell they grant executes only through its existing effect owner, and slot versus free-use routing stays with #1448, as decision #273 requires.

Three ownership points a reviewer should check:

- Faerie Fire (Drow, level 3) has no published effect owner. Its key is still assigned by the open spell splitter #277, so `elf` is blocked by #277, records the grant and keeps the cast unavailable until that owner exists.
- #774 owns the Heroic Inspiration pool and says the Human Resourceful gain remains with #333. `human` owns only that gain.
- Some general rules have no published shared owner and are consumed in a bounded way inside the task that needs them: Tremorsense (PDF190) in `dwarf`, Carrying Capacity for the character's own sheet (PDF178) in `goliath`, Temporary Hit Points and Massive Damage (PDF17–18) in `orc`, a Fly Speed for the flight's duration (PDF182) in `dragonborn`, and typed Resistance (PDF187) in `dragonborn`, `dwarf` and `tiefling`. The three Resistance consumers must share one projection in the damage pipeline: whichever lands first adds it and the others reuse it. None is widened into a general framework.

Where the source is silent, the task says so and refuses rather than inventing: a repeated language roll waits for an explicit new roll or choice, a fourth clockwork device is rejected while three exist, and an underspent point-cost set is accepted without converting the unspent points.

## Native links and caps

Each task has parent #333 and blocks #337 directly, which is what the issue's third criterion asks. No shared aggregate is needed:

| Native list | Before | After | Cap |
|---|---:|---:|---:|
| #337 blocked-by | 62 | 77 | 100 |
| #333 sub-issues | 1 | 16 | 100 |
| #283 sub-issues | 95 | 95 | 100 |
| Largest new blocked-by list (`elf`) | — | 21 | 100 |

No new task is a child of #283. `custom-background` is the only task blocked by other new tasks, so `ability-scores` and `backgrounds` are published first. The dependency graph was read in full (1,892 issues, 10,201 blocked-by edges): no planned prerequisite is transitively blocked by #337, so the fifteen new edges into #337 close no cycle. The parent chain stays #333 → #283 → #54 → #47. No existing body, parent or edge changes; #430 stays a child of #333 and keeps its own edges.

## Deterministic witnesses

Soldier +2 STR/+1 CON on an assigned 15/13 gives 17 and 14. Point costs 15,15,15,8,8,8 and 15,14,13,12,10,8 both total 27. A Dwarf Fighter 1 with CON +2 has 13 Hit Points and 22 after the fixed level 2 application. A level 5 Dragonborn with CON +2 has breath DC 13 and 2d10. A Goliath with Strength 16 carries 480 lb. An Orc with maximum 12 and current 6 Hit Points may drop to 1 against 17 damage and dies against 18. Wisdom 15 with Perception proficiency at level 1 gives Passive Perception 14, the source's own example.

## Verification and limits

Eighty source phrases across PDF8, 17–24, 26–27, 83–87, 94, 178, 186, 190 and 192–193 were compared with the pinned PDF text, and the spell levels the lineage tasks rely on (Detect Magic, Hellish Rebuke 1; Misty Step, Darkness 2) were read from the spell descriptions. Task body hashes, owner body hashes, cap counts and the cycle check are in the scope file.

No runtime file, catalogue, migration, UI or existing issue is changed. Runtime implementation remains paused. Source checks and hosted CI qualify this decision head; PostgreSQL, browser/device, production and Full Rules qualification belong to the later implementation PRs.
