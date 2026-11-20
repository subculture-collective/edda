# Bodily condition ownership (#270)

Four of the six named conditions already have bounded published owners: Grappled #778, Restrained and Poisoned #727, Petrified #1444, and the Prone instance with its attack consequences and movement options #488. Reuse them without changing their bodies. Four source-backed residuals have no owner and become four new Task children of #270:

| Issue | Key | Title | Size | Body |
|---|---|---|---|---|
| #2061 | `prone-movement` | Resolve dropping Prone, Speed-0 righting refusal and crawl cost | M | [task](issue-270-prone-movement-task.md) |
| #2062 | `flying-fall` | Trigger a flying creature fall from Incapacitated, Prone or Fly Speed 0 unless it hovers | S | [task](issue-270-flying-fall-task.md) |
| #2063 | `exhaustion-ledger` | Admit the shared Exhaustion level ledger with death at six and revival reduction | M | [task](issue-270-exhaustion-ledger-task.md) |
| #2064 | `exhaustion-penalties` | Project Exhaustion D20 Test and Speed reductions from the level ledger | M | [task](issue-270-exhaustion-penalties-task.md) |

All four are published as native Task issues on the edda Project. The bodies refer to each other by key (`edda-270:<key>`), which the table above resolves. The companion [scope](issue-270-bodily-conditions-scope.json) records the exact bodies, hashes and native edges, and the readback is in the [link receipt](issue-270-link-verification.json).

## Source and existing owners

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. The named pages PDF176–191 and the pages they refer to (PDF6–7, 13–14, 16–18, 194, 236) were read from that file. Source facts and fixtures retain the CC-BY-4.0 attribution and adaptation notice in [the source register](../srd-5.2.1-conformance.md). #270's only prerequisite, #80, is closed.

| Named scope | PDF | Source clauses | Owner and boundary |
|---|---|---|---|
| Grappled | 182, 190 | Speed 0 that can't increase; Disadvantage on attacks against anyone but the grappler; drag at 1 extra foot unless the target is Tiny or two sizes smaller; one grapple per hand; escape action against the escape DC; ends on grappler Incapacitated, range exceeded or free release | #778 owns all of these with creator-bound instances. #779 owns the Unarmed Strike Grapple save, DC, size limit and free hand. Grappler feat stays #390. |
| Restrained | 187 | Speed 0 that can't increase; attacks against have Advantage, own attacks Disadvantage; Disadvantage on Dexterity saves | #727. The glossary gives no general ending; each source owns its own. |
| Poisoned | 186 | Disadvantage on attack rolls and ability checks | #727 projection; #504 and #616 trusted endings. Saves are not affected. |
| Petrified | 186 | Inanimate with nonmagical gear, weight ×10, no aging; Incapacitated; Speed 0; attacks against have Advantage; automatic Strength and Dexterity save failure; Resistance to all damage; Immunity to Poisoned | #1444 owns every clause and consumes #487 and #727. #624 owns the Greater Restoration ending. |
| Prone, instance and attacks | 186 | Disadvantage on own attacks; attacks against have Advantage within 5 feet, otherwise Disadvantage; the only movement options are to crawl or to right yourself | #488, whose body states "Prone crawl or spend half Speed to stand". |
| Prone, movement numbers | 14, 179, 186 | Dropping Prone without an action or Speed unless Speed is 0; righting cost rounded down; no righting at Speed 0; crawling costs 1 extra foot, 2 in Difficult Terrain | **`prone-movement`.** It consumes #488's options and adds only what #488 does not state. "Crawl" and "half Speed to stand" occur in no other body, and no body mentions dropping Prone or the Speed-0 rule. |
| Prone or Fly Speed 0 while flying | 182, 183, 188 | A flying creature falls if Incapacitated, Prone or at Fly Speed 0, unless it can hover; Speed reduced to 0 reduces every special speed to 0 | **`flying-fall`.** No body mentions the rule; #147 says "Exclude flying" and keeps the fall itself; #1336 keeps the Fly spell ending. |
| Exhaustion, levels | 179–181, 183 | Cumulative levels, the one condition that stacks; death at level 6; condition ends at level 0; Immunity; revival returns with 1 fewer level | **`exhaustion-ledger`.** No issue stores a level. |
| Exhaustion, penalties | 180–181, 188 | Every D20 Test reduced by 2 × level; Speed and each special speed reduced by 5 × level | **`exhaustion-penalties`.** No issue applies either penalty. |
| Exhaustion, reduction | 185 | A finished Long Rest lowers the level by 1 | #428 keeps the Long Rest trigger. #623 keeps the shared one-level reduction event and its recomputation call; #1001 consumes it through #623, as its body says. The ledger supplies only the stored level and the removable-level rule. |

The search covered the titles and whole bodies of all 1,892 issues, open and closed. Exhaustion appears in 38 bodies: three reducers, immunity lines in stat blocks and spells, and roadmap lists. None admits a level or applies a penalty. "While flying", "fall if you have" and "Fly Speed is reduced to 0" appear in none. Owner body hashes and gate paths are in the scope file.

### Sources of Prone

| Source | PDF | Owner |
|---|---|---|
| Unconscious | 191 | #487 |
| Unarmed Strike Shove | 190 | #779 |
| Landing after a fall | 182 | #147 |
| Long Jump landing in Difficult Terrain, DC 10 Acrobatics | 185 | #150, the open hazard split, which names jump and must publish it |
| Burning: an action to give yourself Prone and put out the fire | 178 | Unowned. No issue names Burning. Nearest planning owner #150, whose list does not name it. `prone-movement` refuses it: it costs an action and is not the free drop. |
| Ending a turn in another creature's space | 14 | Unowned. No issue names it and no task owns positions or space occupancy. Nearest planning owner #150. |
| Falling off a mount, and the save when rider or mount is knocked Prone | 16 | Unowned. Nearest planning owner the mount split #274, whose named pages are PDF100–101. |

The three unowned sources lie outside the six named conditions, so they get no child here. They stay unavailable and are listed as follow-ups below.

Dehydration, malnutrition and suffocation (PDF181, 185, 189) are sources of Exhaustion, not the condition. They stay with #150, which must publish them. Their rules that some levels can't be removed, or that all levels from one source are removed, are why the ledger keeps provenance per level. Travel, contagion, temperature, item and stat-block sources outside the named pages stay with #150, #278 and their stat-block owners.

## The four new tasks

**#2061 `prone-movement`.** Voluntary drop (PDF14), righting cost from the current Speed rounded down, the Speed-0 refusal against actual #727 and #778 instances, and the two crawl figures. The cost comes from Speed, never a special speed. Because no task owns a per-turn movement record, it adds the minimal one for its own spends. Fixtures: Speed 30 rights for 15, Speed 25 for 12, Speed 30 with Fly Speed 60 still 15; a standing Grappled or Restrained creature cannot drop Prone and a Prone one cannot right itself, even with a Speed bonus; crawling 5 feet costs 10, or 15 with a trusted Difficult Terrain fact. Prerequisites: #488, #727, #778, #780.

**#2062 `flying-fall`.** Decides only whether a fall starts, then hands one declared fall with the trusted height to #147. Triggers are an actual Prone or Incapacitated instance, or a Fly Speed of 0, which PDF188 makes follow from Grappled and Restrained. Hover prevents the fall. Prerequisites: #147, #487, #727, #778, `prone-movement`.

**#2063 `exhaustion-ledger`.** Levels 0–6 with provenance: a living creature holds 0–5, reaching 6 commits a death record and the dead creature keeps level 6. Immunity, the removable-level rule and the revival reduction on the actual #618 receipt. It fails closed: a dead creature with levels becomes alive only through a receipt that applied the reduction, so #947 and later revival owners stay unavailable for such creatures until they call it. Prerequisites: #64, #618.

**#2064 `exhaustion-penalties`.** A projection of the ledger level into the #61 check and save receipts, the #65 attack roll, the #66 Initiative roll and the Speed #780 reads. Not Disadvantage; die-face results still follow the die. Composition fixtures: level 3 rights from Prone for 7, and level 4 brings Fly Speed 20 to 0 and a fall. Prerequisites: `exhaustion-ledger`, `prone-movement`, `flying-fall`, #61, #65, #66, #780.

Two choices are Edda adaptations and are labelled as such in the bodies: Speed is clamped at 0 when the reduction exceeds it, and a reduction takes the most recently gained removable level. The source states neither.

## Native edges and caps

All four tasks are children of #270, which goes from 4 to 8 sub-issues.

The spell gate #279 has 100 blockers, the native cap, so no task can block it directly. `exhaustion-penalties` blocks #623 instead. #623 is an existing direct blocker of both #279 and #280, and its body removes one Exhaustion level and "recompute[s] all derived penalties", which is this projection. Today #623's only prerequisite is #270 itself; without the new edge, closing this decision would leave it blocked by nothing concrete. #623 goes from 1 to 2 blockers. Its body, parent and existing edge are unchanged.

The edges from `prone-movement` and `flying-fall` to `exhaustion-penalties` are cap routing. The penalty task needs their outputs only for its two composition fixtures; the edges exist because neither task has another route to the full #279, following the shared-aggregate precedent of the #273 decision. They place #623, #1001 and the Exhaustion path behind #147, #487, #488, #727, #778 and #780.

Resulting mandatory paths: #279 → #623 → `exhaustion-penalties` → each of the other three; #280 → #623 → the same; #337 → #551 → #279 → #623 and #337 → #970 → #971 → #1001 → #623. None passes through #270, so closing the decision removes no obligation.

Counts read back after publication, all at or under 100: #279 100 blockers, #280 97 and #337 62, none of them touched (#280 was 96 when this was drafted; #358 added #2059 first); #623 2 blockers; consumers of #61 85, #64 50, #65 74, #66 63, #147 5, #487 88, #488 98, #618 10, #727 35, #778 16, #780 12. #67 and #429 are at 100 consumers and gain none. The blocked-by graph and parent hierarchy are acyclic in the native readback.

No edge is added to #428. It would be acyclic, but it would put Long Rest and its 81 consumers behind Revivify #618, which the ledger needs for the revival clause. #428's "reduce Exhaustion by 1" stays unavailable until the ledger exists, and both gates stay blocked by the ledger through #623. No edge is added to #1001, which already reaches the ledger through #623, or to #947 or #390.

## Publisher confirmations

The reviewed package left three points for the publisher. All three were applied as planned and are recorded in the link receipt.

- The reading of #488. Its phrase "crawl or spend half Speed to stand" stays with #488, and `prone-movement` adds only the unstated clauses. #488's body cannot be edited to say so.
- The #623 edge. It adds a blocker to an existing non-gate issue, which the #273, #335 and #336 decisions did not do. #623's body still says "Native prerequisites: #270." while its native blockers are #270 and #2064. The link receipt records that divergence with #623's unchanged body hash.
- The ledger is ordered after Revivify #618 so that the revival reduction is tested on a real receipt.

## Follow-ups for the owner

These are outside this decision and should be repeated in the PR body.

- #390 Grappler names #270 for its shared grapple rules and has no native edge to #778 or #779. When #270 closes it becomes startable without them. All three stay mandatory for #337.
- #150 must publish the Long Jump landing rule, and should take Burning (PDF178) and the shared-space Prone rule (PDF14). #274 should take the mounted fall (PDF16).
- No published task owns Death Saving Throws, 0 Hit Points or a general dead state (PDF17–18). The Exhaustion reduction on a Death Saving Throw stays unavailable until one exists.
- No published task owns a flying state with height, or general movement along a path with positions.
- #428 has no native edge to the ledger.

## Evidence and limits

Current code has only legacy string flags in `internal/combat/state.go`: `prone` and `poisoned` give attack Disadvantage, `exhausted` has no level or effect, and Grappled, Restrained and Petrified do not exist. Its death-save code is legacy too. `apply_condition` takes a model-supplied name and round count. None of this is a source-bound admission, and no new task converts a legacy flag into the new contract.

Unknowns the source does not settle are listed in the scope file and stay unavailable: a contagion level that "lasts until the contagion ends" (PDF194), removing "any Exhaustion levels" (PDF236) against non-removable levels, revival against only non-removable levels, and how righting counts against a special speed used in the same turn.

Checked for the draft: the PDF hash, 62 exact source phrases on 22 pages against the extracted page text, every owner body, the cap counts and both acyclicity simulations, from a full issue readback on 2026-11-20.

Checked at publication, read back at 2026-11-20T08:14:40Z: each of #2061–#2064 matches its committed task file byte for byte, with the planned title, OPEN state, type Task, Size, parent #270 and blocker set. A full page-through of 1,898 issues shows 10,225 dependency edges and 1,894 parent edges, both graphs acyclic, and no issue above 100 blockers, dependents or children. Against the read taken just before publication, the only differences in the 1,894 earlier issues are the added blocker #2064 on #623 and #358 closing, which this publication did not cause. Every reused owner keeps its body hash, parent and blockers. Runtime implementation remains paused. This change contains no runtime file, migration, UI behavior, production mutation or Full Rules admission; the Go suite was not run locally and hosted CI covers the decision head. The implementation tasks keep their own PostgreSQL, web/TUI and device requirements.
