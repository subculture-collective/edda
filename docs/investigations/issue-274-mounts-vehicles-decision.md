# Mount, vehicle and cargo ownership (#274)

The three equipment tables on PDF100–101 and the rules printed beside them have no published execution owner. The #80 inventory lists these pages as an unindexed section assigned to #274, and none of its 753 indexed families covers them. This decision assigns them to six bounded tasks and one cohort qualification, all children of #274. The cohort task is the single new blocker of item qualification #280 and of exploration qualification #183. Stat blocks, travel-pace variants and the economy keep their existing owners. The seven tasks are published as #2074–#2080; `issue-274-link-verification.json` records the readback.

## Source and existing owners

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Source facts and fixtures retain the CC-BY-4.0 attribution and adaptation notice in [the source register](../srd-5.2.1-conformance.md). Every table row and quoted clause in the scope file was compared with the PDF page text.

| Named scope | PDF | Owner and boundary |
|---|---|---|
| Mounts and Other Animals: eight rows of carrying capacity and cost | 100 | New records task. Data only. |
| Stat blocks of those eight animals | 346–348, 357, 359–360, 364 | Reused: #156 Camel, #158 Draft Horse, #159 Elephant, #167 Mastiff, #168 Mule, #170 Pony, #171 Riding Horse, #176 Warhorse. Their mechanic keys stay with #182. These are catalogue records; #143 says catalogue presence is not playable support. |
| Executable Riding Horse statistics | 360 | Reused: #1999 owns the ordinary Riding Horse (Large Beast, Speed 60). The riding task takes it as a native prerequisite. No published task makes another table animal an acting creature, so only the Riding Horse can be ridden. |
| Carried load, Mounts and Cargo pull limit, Barding cost and weight | 92, 100 | New load task. Uses table capacities as written and the #251/#252 Armor-table records. |
| Mounted Combat: eligibility, mounting, dismounting, control | 15–16 | New riding task. A whole-body search for "Mounted Combat", "Controlling a Mount" and "dismount" found only #920 and #1931. #920 specifies shared Initiative and a controlled mount for the Find Steed steed, and #1931 the dismount at the end of Phantom Steed; neither names a general mounted-state provider. |
| Mounted Combat: falling off; Saddles | 16, 100 | New saddle-and-fall task on top of the riding task. The two saddle rules need a mounted state. |
| Tack, Harness, and Drawn Vehicles: ten priced rows | 100 | New records task. |
| Airborne and Waterborne Vehicles: seven rows; usual passenger prices | 100–101 | New records task. |
| Large-vehicle speed, crew, passengers, cargo | 5, 12, 100–101, 192 | New vessel-hour task on top of the #144 travel segment. Carve-out: for the seven table vehicles it owns the PDF12 sentence on waterborne vessels and the PDF192 Vehicles paragraph, both inside #150's page range. |
| Damage threshold and Ship Repair | 5, 101, 180 | New hull task. It also owns the vessel instance. |
| Mounted doubled distance; land vehicles choose a pace as normal | 12 | Reused by inference: #150's body names neither mounts nor vehicles; its ownership is read from "remaining travel pace variants" on PDF12–13. #144 excludes vehicles. |
| Hiring and paying hirelings; coins, buying and selling | 89, 102 | Reused: the #275 economy split. The vessel-hour and hull tasks wait on #275 and consume the crew and coin outputs it assigns. |

No current runtime code resolves a mount, vehicle, carried weight, damage threshold or coin payment. Starting gold exists only as a Coin Pouch equipment entry.

## Published tasks

Seven issues, each type Task under parent #274. Bodies are the companion task files, byte for byte. The reviewed files referred to each other with `{{274:<key>}}` placeholders; publication replaced them with the issue numbers below, and the scope file keeps both the template hash and the published hash.

| Issue | Key | Title | Size | Blocked by |
|---|---|---|---|---|
| #2074 | records | Catalogue SRD mount, tack, drawn-vehicle and large-vehicle table records | M | #195 |
| #2075 | load | Load pack animals and hitched drawn vehicles within table carrying capacity | M | #2074, #251, #252 |
| #2076 | ride | Mount, dismount and control one trained Riding Horse | M | #2074, #66, #171, #780, #781, #782, #1999 |
| #2077 | unseat | Resolve falling off a mount and the Military and Exotic Saddle rules | M | #2074, #2076, #61, #488 |
| #2078 | hull | Apply ship damage thresholds and berthed repair to one vessel instance | M | #2074, #275 |
| #2079 | voyage | Sail one large-vehicle hour with crew, wind, current, passenger and cargo limits | L | #2074, #2078, #144, #275 |
| #2080 | cohort | Qualify the mount, drawn-vehicle and ship fixture cohort | M | #2074, #2075, #2076, #2077, #2078, #2079 |

Deviation from #274's acceptance criteria. #274 says to inspect only its named PDF100–101 scopes and to keep wider scope as named spikes. The ride and unseat tasks implement PDF15–16 Mounted Combat as implementation children. This is deliberate: the Military and Exotic Saddle rules cannot be tested without mounting, control and the falling-off save; the search above found no owner for that text; and a spike would only restate one column. The text is split over two tasks so each stays one reviewable PR.

Persisted state. The records are data only, so each rule task names the instance it owns. The hull task owns the vessel instance (identity, owner, current Hit Points), and the vessel-hour task is blocked by it and reads it. The load task owns the owned-animal and drawn-vehicle instance. The ride task owns the mounted relation between two encounter participants, and the unseat task adds the fitted saddle to it. A mount, vehicle, saddle or ship enters a campaign only as an explicit isolated or GM-declared fact.

Find Steed and Phantom Steed. #920 and #1931 specify mounted rules for their own steeds, and their bodies list exact native prerequisites that cannot be edited, so this decision adds no edge to them. The ride body requires extending any mounted or controlled-mount state they have already merged, and the scope file records the same duty for whichever merges second.

Carve-out from #150. #150's page range (PDF12–13, 192–201) contains the PDF12 "Vehicles" paragraph and the PDF192 "Vehicles" paragraph. For the seven table vehicles the vessel-hour task owns the waterborne-vessel sentence on PDF12 and the PDF192 paragraph, because they only point to the PDF101 table speed that #274 names. #150's split must not publish a second owner for them. It keeps the doubled mounted distance and land-vehicle pace. #150's body cannot be edited, so the carve-out is recorded here and in the scope file's follow-up obligations.

The #275 edges. Blocking on the open #275 split follows the precedent of #386 on #272, #388 on #333 and #390 on #270, but such an edge is satisfied when the planning issue closes. Before #275 closes, its split must add its skilled-hireling child as a native blocker of the vessel-hour task and its coin-payment child as a native blocker of the hull task. Both task bodies say so, and the cohort task fails qualification otherwise.

Five readings are fixed in the task bodies because they are easy to get wrong.

1. The Military Saddle gives Advantage on ability checks to remain mounted; the falling-off roll is a DC 10 Dexterity saving throw, so the saddle does not affect it.
2. A damage instance equal to the threshold is taken in full: the PDF180 rule says "equal to or greater than", and its example only shows 9 and 11 against a threshold of 10.
3. The PDF5 Round Down rule is applied where the source counts whole things or says so: the mounting cost (PDF15 says "round down") and private-accommodation passengers. It is not applied to a halved vessel speed or a halved repair time. The PDF101 table prints 1½ and 2½ miles per hour, so vessel speeds are not whole numbers and a rounded value would depend on the unit; and rounding half a day down would make a 1 Hit Point shipyard repair take no time. A Warship under sail against a strong wind moves 1¼ miles, and 15 Hit Points in a shipyard take 7½ days. The rounded results (1 mile, 7 days, 0 days) are negative fixtures. This is an Edda reading of an unstated interaction.
4. The strong-wind rule applies to "a ship sailing", and the source does not say which rows sail. Under sail or rowed is a declared fact for an hour with a strong opposing wind or a dead calm. The source treats the Keelboat and Rowboat as rowed vessels, so a rowed hour for them uses the table speed. No other row has a rowed speed.
5. PDF16 says "Domesticated horses, mules, and similar creatures have such training." Of the table animals, Horse, Draft; Horse, Riding; Warhorse and Mule are read as horses and mules. Camel, Elephant, Mastiff and Pony are similar creatures only by authoritative fact.

## Kept unavailable

The source does not state these, so no task implements them: an overload penalty for a table mount or drawn vehicle; a pack load combined with pulling; Armor Class, training or donning for barding worn by a mount, and Shield barding; a rowed speed for the Galley, Longship, Sailing Ship or Warship, and so any dead-calm hour for them; towed speeds; a strong opposing wind with no declared propulsion; whether an Airship sails, and so an Airship hour against a strong wind or in a dead calm; a current on vessels other than the Keelboat and Rowboat; pounds per ton; Poison or Psychic damage to a vessel, since PDF12 says a vehicle is not an object; any consequence at 0 Hit Points; feeding or stabling requirements beyond the listed prices; which creatures are "similar creatures" trained to accept a rider; how a campaign acquires a mount, vehicle, saddle or ship other than by a declared fact; and using an owned table animal as an encounter mount. Purchase of any record and passenger fares stay unavailable until the #275 split publishes a payment owner.

Three items lie outside the named scopes and are recorded for the item and world reconcilers, not assigned here. The general Carrying Capacity rule on PDF178 has no published owner. The Lance record on PDF91 is owned by #248, and its "unless mounted" key is for the #278 reconciler; no execution task exists yet, and the ride task's mounted state is the fact that key will consume. The PDF201 rolling stone's damage threshold lies inside #150's unsplit range.

## Native graph and caps

Counts read at 2026-11-20T07:49:40Z across 1,892 issues. #280 has 96 of 100 blocked-by entries, and the other open splitters still need edges there, so the six implementation tasks reach #280 through the cohort task on one edge, 97 after publication. The cohort task also blocks #183, 13 of 100 before and 14 after: #183's range (PDF11–18, 176–203) contains the Mounted Combat rules and the vessel-speed sentences, and its body requires implementation children to be native blockers. #183 is already a transitive blocker of #280, so the edge adds no cycle. #274 has no sub-issues and gains seven. The busiest reused prerequisite is #488 at 97 of 100 consumers, 98 after the unseat task. #1999 goes from 4 to 5 consumers and #275 from 1 to 3. #67 and #429 are at 100 consumers and #279 at 100 blockers; none gains an edge. The ride and unseat tasks use #67's merged scene facts without a native edge. No task blocks spell qualification #279, because nothing here is a spell effect.

Those counts were a snapshot. #270 and #272 published against the same gates before this decision did, so publication re-read both ends of every edge immediately before adding it and would have stopped at 100. At publication #280 went from 98 to 99 blocked-by entries and #183 from 13 to 14; #488 went from 98 to 99 consumers, #61 from 86 to 87, #66 from 63 to 64, #195 from 62 to 63 and #1999 from 4 to 5. #280 and #488 each have one free slot left. The re-read values are in `issue-274-link-verification.json`.

Read back at 2026-11-20T15:11:54Z across 1,912 issues, the dependency graph (10,281 edges) and the parent hierarchy (1,908 edges) are acyclic and the parent chain is #274, #193, #54, #47. Spell qualification #279 gained no edge, but it already depended on #183 through Wall of Stone #681, so the #183 edge also places the cohort task behind #279 on that existing path. Closing #274 removes no obligation: every implementation task is a prerequisite of the cohort task, and the cohort task blocks #280 and #183.

## Evidence and limits

The scope file holds the checked clauses and table rows, the fixed readings, the follow-up obligations for #150, #275, #278, #920 and #1931, the owner body hashes, the planned bodies' template hashes, edges and cap counts. No existing issue body, parent or edge changes except the new #280 and #183 blocker and the consumer edges on the listed prerequisites. Runtime implementation remains paused. This change contains no runtime file, migration, UI behavior, production mutation or Full Rules admission. Focused source checks and hosted CI apply to this decision; the implementation tasks keep their own runtime, PostgreSQL and device requirements.
