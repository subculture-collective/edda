# Bodily condition ownership (#270)

Four of the six named conditions already have bounded published owners: Grappled #778, Restrained and Poisoned #727, Petrified #1444, and the Prone instance with its attack consequences #488. Reuse them without changing their bodies. Two source-backed residuals have no owner and become two new Task/M children of #270: one for the Prone movement clauses and one for the Exhaustion level ledger. The child bodies are [the Prone movement task](issue-270-prone-movement-task.md) and [the Exhaustion task](issue-270-exhaustion-task.md). They are not published yet; the companion [scope](issue-270-bodily-conditions-scope.json) records the planned native edges, and the link receipt is written at publication.

## Source and existing owners

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. The named pages PDF176–191 and the pages they refer to (PDF6–7, 13–14, 17–18) were read from that file. Source facts and fixtures retain the CC-BY-4.0 attribution and adaptation notice in [the source register](../srd-5.2.1-conformance.md). #270's only prerequisite, #80, is closed.

| Named scope | PDF | Source clauses | Owner and boundary |
|---|---|---|---|
| Grappled | 182, 190 | Speed 0 that can't increase; Disadvantage on attacks against anyone but the grappler; drag at 1 extra foot unless the target is Tiny or two sizes smaller; one grapple per hand; escape action against the escape DC; ends on grappler Incapacitated, range exceeded or free release | #778 owns all of these with creator-bound instances. #779 owns the Unarmed Strike Grapple save, DC, size limit and free hand. Grappler feat stays #390. |
| Restrained | 187 | Speed 0 that can't increase; attacks against have Advantage, own attacks Disadvantage; Disadvantage on Dexterity saves | #727. The glossary gives no general ending; each source owns its own. |
| Poisoned | 186 | Disadvantage on attack rolls and ability checks | #727 projection; #504 and #616 trusted endings. Saves are not affected. |
| Petrified | 186 | Inanimate with nonmagical gear, weight ×10, no aging; Incapacitated; Speed 0; attacks against have Advantage; automatic Strength and Dexterity save failure; Resistance to all damage; Immunity to Poisoned | #1444 owns every clause and consumes #487 and #727. #624 owns the Greater Restoration ending. |
| Prone, attacks | 186 | Disadvantage on own attacks; attacks against have Advantage within 5 feet, otherwise Disadvantage | #488. Prone from Unconscious is #487, from Shove #779, from a fall #147. |
| Prone, movement | 14, 179, 186 | Dropping Prone for free unless Speed is 0; righting costs half Speed rounded down and is impossible at Speed 0; crawling costs 1 extra foot, 2 in Difficult Terrain | **New Prone movement task.** #488 names "crawl or spend half Speed to stand" but binds no rounding, no Speed-0 refusal, no crawl cost and no voluntary drop; no other body mentions them. |
| Exhaustion | 179–181, 183, 188 | Cumulative levels, the one condition that stacks; death at level 6; every D20 Test reduced by 2 × level; Speed reduced by 5 × level; condition ends at level 0; revival returns with 1 fewer level | **New Exhaustion task.** No issue owns the ledger or its penalties. |
| Exhaustion, reduction | 185 | A finished Long Rest lowers the level by 1 | #428 Long Rest, #623 Greater Restoration, #1001 Tireless. They keep their triggers and use the new ledger's single reduce-one output. |

The search covered the titles and whole bodies of all 1,892 issues, open and closed. Exhaustion appears in 36 bodies: three reducers, immunity lines in stat blocks and spells, and roadmap lists. None admits a level or applies a penalty. "Crawl" appears only in #488, and no body mentions dropping Prone or the Speed-0 righting rule. Owner body hashes and gate paths are in the scope file.

Dehydration, malnutrition and suffocation (PDF181, 185, 189) are sources of Exhaustion, not the condition. They stay with the open hazard split #150, which must publish them. Their rules that some levels can't be removed, or that all levels from one source are removed, are why the new ledger keeps provenance per level. Travel, contagion, temperature, item and stat-block sources outside the named pages stay with #150, #278 and their stat-block owners.

## The two new tasks

**Resolve Prone dropping righting and crawling against Speed 0.** It consumes the Prone instance from #488 and the Speed-0 projections from #727 and #778. Fixtures: Speed 30 rights for 15, Speed 25 for 12; a Grappled or Restrained Prone creature can neither right itself nor drop Prone, even with a Speed bonus; crawling 5 feet costs 10, or 15 in Difficult Terrain. Prerequisites: #488, #727, #778.

**Admit shared Exhaustion levels with D20 Test Speed and death consequences.** One ledger per creature, levels 1–5, death at 6, source provenance per level, Immunity honoured. The roll reduction is a numeric term in the #61 receipt, not Disadvantage; it covers Initiative and Death Saving Throws, and natural-20 and natural-1 results still follow the die face (PDF7, 18). Speed and every special speed fall by 5 feet per level (PDF188). Fixtures include level 2 turning 17 into 13 and Speed 30 into 20, level 3 righting from Prone for 7, a sixth level killing once, and revival from level 6 at level 5. Prerequisites: #61, #64, #66 and the Prone movement task.

Two choices in the Exhaustion body are Edda adaptations and are labelled as such there: Speed is clamped at 0 when the reduction exceeds it, and reduce-one takes the most recently gained removable level. The source states neither.

## Native edges and caps

Both tasks are children of #270, which goes from 4 to 6 sub-issues. The Prone movement task blocks the Exhaustion task, whose Speed fixture uses righting.

The spell gate #279 has 100 blockers, the native cap, so neither task can block it directly. The Exhaustion task instead blocks #623. #623 is an existing direct blocker of both #279 and #280, and its body removes one Exhaustion level and recomputes the derived penalties, which are exactly this ledger's outputs. Today #623's only prerequisite is #270 itself; without the new edge, closing this decision would leave it blocked by nothing concrete. This is the one addition to an existing non-gate issue: #623 goes from 1 to 2 blockers. Its body, parent and existing edge are unchanged. A direct #280 edge alone was rejected because #280 does not block #279.

Resulting mandatory paths: #279 → #623 → Exhaustion → Prone movement; #280 → #623 → the same; #337 → #551 → #279 → #623 and #337 → #970 → #971 → #1001 → #623. None passes through #270, so closing the decision removes no obligation.

Other counts after publication: #280 stays at 96 blockers and #337 at 62; consumers of #61 85, #64 50, #66 63, #488 98, #727 34 and #778 15, all under 100. #67 and #429 are at 100 consumers and gain none. No edge is added to #428, #1001, #390 or any other existing issue. #390 Grappler still names #270 for its shared grapple rules and has no native edge to #778 or #779; edges between two existing issues are outside this decision, and all three stay mandatory for #337. The simulated blocked-by graph and parent hierarchy are acyclic with the planned edges.

## Evidence and limits

Current code has only legacy string flags in `internal/combat/state.go`: `prone` and `poisoned` give attack Disadvantage, `exhausted` has no level or effect, and Grappled, Restrained and Petrified do not exist. `apply_condition` takes a model-supplied name and round count. None of this is a source-bound admission, and neither new task converts a legacy flag into the new contract.

Checked for this draft: the PDF hash, 47 exact source phrases on 17 pages against the extracted page text, every owner body, the cap counts and both acyclicity simulations, from a full issue readback on 2026-11-20. Not done: publication, native readback of the new issues and the link receipt. Runtime implementation remains paused. This change contains no runtime file, migration, UI behavior, production mutation or Full Rules admission; the Go suite was not run locally and hosted CI covers the decision head. The implementation tasks keep their own PostgreSQL, web/TUI and device requirements.
