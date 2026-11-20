# Mount, vehicle and cargo ownership (#274)

The three equipment tables on PDF100–101 and the rules printed beside them have no published execution owner. The #80 inventory lists these pages as an unindexed section assigned to #274, and none of its 753 indexed families covers them. This decision assigns them to five bounded tasks and one cohort qualification, all children of #274. The cohort task is the single new blocker of item qualification #280. Stat blocks, travel-pace variants and the economy keep their existing owners. Nothing is published yet: the companion scope records the planned issues and edges, and the link-verification receipt is written at publication.

## Source and existing owners

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Source facts and fixtures retain the CC-BY-4.0 attribution and adaptation notice in [the source register](../srd-5.2.1-conformance.md). Every table row and quoted clause in the scope file was compared with the PDF page text.

| Named scope | PDF | Owner and boundary |
|---|---|---|
| Mounts and Other Animals: eight rows of carrying capacity and cost | 100 | New records task. Data only. |
| Stat blocks of those eight animals | 346–348, 357, 359–360, 364 | Reused: #156 Camel, #158 Draft Horse, #159 Elephant, #167 Mastiff, #168 Mule, #170 Pony, #171 Riding Horse, #176 Warhorse. Their mechanic keys stay with #182. |
| Carried load, Mounts and Cargo pull limit, Barding cost and weight | 92, 100 | New load task. Uses table capacities as written and the #251/#252 Armor-table records. |
| Saddles | 100 | New riding task. The two saddle rules need a mounted state. |
| Mounted Combat | 15–16 | New riding task. Whole-body search found no owner; #920 and #1931 only consume a mounted state for their own steeds. |
| Tack, Harness, and Drawn Vehicles: ten priced rows | 100 | New records task. |
| Airborne and Waterborne Vehicles: seven rows; usual passenger prices | 100–101 | New records task. |
| Large-vehicle speed, crew, passengers, cargo | 12, 100–101, 192 | New vessel-hour task on top of the #144 travel segment. |
| Damage threshold and Ship Repair | 101, 180 | New hull task. |
| Mounted doubled distance; land vehicles choose a pace as normal | 12 | Reused: #150 owns the remaining travel-pace variants on PDF12–13. #144 excludes vehicles. |
| Hiring and paying hirelings; coins, buying and selling | 89, 102 | Reused: the #275 economy split. The vessel-hour and hull tasks wait on #275 and consume the crew and coin outputs it assigns. |

No current runtime code resolves a mount, vehicle, carried weight, damage threshold or coin payment. Starting gold exists only as a Coin Pouch equipment entry.

## Planned tasks

Six issues, each type Task under parent #274. Bodies are the companion task files. They refer to each other with `{{274:<key>}}` placeholders that publication replaces with issue numbers.

| Key | Title | Size | Blocked by |
|---|---|---|---|
| records | Catalogue SRD mount, tack, drawn-vehicle and large-vehicle table records | M | #195 |
| load | Load pack animals and hitched drawn vehicles within table carrying capacity | M | records, #251, #252 |
| ride | Ride one trained mount with mounting cost, control, falling off and saddle rules | L | records, #61, #66, #171, #488, #780, #781, #782 |
| voyage | Sail one large-vehicle hour with crew, wind, current, passenger and cargo limits | L | records, #144, #275 |
| hull | Apply ship damage thresholds and berthed repair by the day | M | records, #275 |
| cohort | Qualify the mount, drawn-vehicle and ship fixture cohort | M | records, load, ride, voyage, hull |

The riding task includes PDF15–16 although #274 cites PDF100–101. The Military and Exotic Saddle rules cannot be tested without mounting, control and the falling-off save, the text is one column, and a separate spike would only restate it. The task is bounded to one Riding Horse fixture and changes neither Find Steed nor Phantom Steed.

Two readings are fixed in the task bodies because they are easy to get wrong. The Military Saddle gives Advantage on ability checks to remain mounted; the falling-off roll is a DC 10 Dexterity saving throw, so the saddle does not affect it. A damage instance equal to the threshold is taken in full: the PDF180 rule says "equal to or greater than", and its example only shows 9 and 11 against a threshold of 10.

## Kept unavailable

The source does not state these, so no task implements them: an overload penalty for a table mount or drawn vehicle; a pack load combined with pulling; Armor Class, training or donning for barding worn by a mount, and Shield barding; rowed and towed vessel speeds; whether an Airship sails for the strong-wind rule; a current on vessels other than the Keelboat and Rowboat; pounds per ton; Poison or Psychic damage to a vessel, since PDF12 says a vehicle is not an object; any consequence at 0 Hit Points; feeding or stabling requirements beyond the listed prices. Purchase of any record and passenger fares stay unavailable until the #275 split publishes a payment owner. The general Carrying Capacity rule on PDF178 and the Lance's mounted exception on PDF91 have no published owner. They lie outside the named scopes and are recorded for the item and world reconcilers, not assigned here.

## Native graph and caps

Counts read at 2026-11-20T07:34:19Z across 1,892 issues. #280 has 96 of 100 blocked-by entries, and the other open splitters still need edges there, so the five implementation tasks reach #280 through the cohort task on one edge, 97 after publication. #274 has no sub-issues and gains six. The busiest reused prerequisite is #488 at 97 of 100 consumers, 98 after the riding task. #67 and #429 are at 100 consumers and #279 at 100 blockers; none gains an edge. The riding task uses #67's merged scene facts without a native edge. No task blocks spell qualification #279, because nothing here is a spell effect, and #183 is unchanged.

With the planned edges added to the full graph, the dependency graph stays acyclic and the parent chain is #274, #193, #54, #47. Closing #274 removes no obligation: every implementation task is a prerequisite of the cohort task, and the cohort task blocks #280.

## Evidence and limits

The scope file holds the checked clauses and table rows, the owner body hashes, the planned bodies' template hashes, edges and cap counts. No existing issue body, parent or edge changes except the one new #280 blocker. Runtime implementation remains paused. This change contains no runtime file, migration, UI behavior, production mutation or Full Rules admission. Focused source checks and hosted CI apply to this decision; the implementation tasks keep their own runtime, PostgreSQL and device requirements.
