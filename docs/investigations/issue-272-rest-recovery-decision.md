# Interrupted rest, recovery and recharge ownership (#272)

Most named rest and recharge rules already have bounded published owners: #427 for the Short Rest, #428 for the Long Rest and #1148 for one wand's dawn recharge. Reuse them without changing their bodies. Whole-body inspection of all 1,892 issues found four source clauses with no owner. This decision assigns each to one implementation task and adds one qualification task that carries all of them to the gates. The five exact bodies are the companion task files; native publication, numbers and readback are recorded later in the link-verification receipt. Closing this decision supplies no runtime.

## Source and existing owners

The English SRD 5.2.1 PDF matches SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Source facts and fixtures retain the CC-BY-4.0 attribution and adaptation notice in [the source register](../srd-5.2.1-conformance.md). The issue cites PDF176–191 and 204–208; the clauses it names also sit on PDF17–18 (damage, healing and Temporary Hit Points) and PDF251 (the wand text), which were read as well.

| Named scope | PDF | Reused implementation contracts and boundary |
|---|---|---|
| Short Rest interruption and resumption | 17, 183, 187 | #427: a 1-hour rest, at least 1 Hit Point to start, Hit Point Dice one at a time with the Constitution modifier, minimum 1 each and capped at the maximum, exactly one recharge event. Initiative, a non-cantrip spell or any damage stops it with no benefit. The source gives a Short Rest no resumption rule; only a Long Rest can resume. A later Short Rest is a new rest. |
| Short Rest consumers | 102, 187 | Attunement #1141 and #2051, Prayer of Healing #926, Tireless #1001 and the class recoveries consume the completed #427 rest. An interrupted rest fails the attunement. None of them owns the rest. |
| Long Rest interruption and eligibility | 181, 185 | #428: at least 8 hours, 6 asleep, at most 2 of light activity, at least 1 Hit Point, 16 hours since the last one. Initiative, a non-cantrip spell, any damage or 1 hour of exertion interrupts; at least 1 hour rested gives Short Rest benefits; resuming adds 1 hour per interruption. Completion restores Hit Points, Hit Point Dice, a reduced maximum and reduced ability scores and removes 1 Exhaustion level, with one recovery event. |
| Per Day and X/Day uses | 186, 257 | #428 through each consumer: #1681, #1688 and #1717 recover on a completed Long Rest, not at dawn. Monster Recharge X–Y stays with #182. |
| Dawn recharge | 206, 251 | #1148: the Wand of Magic Detection regains 1d3 expended charges at each authoritative dawn, capped at 3, with a trusted GM event on a plane without dawn. #1140 owns its charge spend and #1132 the Thief's charge-saving roll. |
| Charge knowledge and consumption | 206 | #1613 reveals charges through Identify. Consumables that lose their magic are the closed #79 and #266 potion fixtures and the scroll #1142. |
| Maximum Hit Points | 17, 185 | Healing is capped at the maximum in #427 and each healing owner; #428 restores a reduced maximum; #627 owns Greater Restoration's removal. Reductions stay with their sources. |

## Unowned residuals

1. **Knocking Out a Creature (PDF17, 184).** A melee attack that would reduce a creature to 0 Hit Points can leave it at 1, Unconscious and starting a Short Rest. No issue owns it; Rogue #1120 is a different Cunning Strike effect. The task consumes #427 and #487 and names the three source endings. After an interrupted rest the source gives no further ending, so none is invented.
2. **Temporary Hit Points (PDF18, 190).** Six grant owners (#645, #725, #807, #1000, #1318, #1969) each restate part of the pool rule and one refers to a shared rule that no issue supplies. No Long Rest contract removes the pool on completion. The task owns the pool once, with False Life #1318 as its actual witness, and consumes #428 for the ending.
3. **Sleep during a Long Rest (PDF185, 191).** #428 requires 6 hours of sleep but does not state that the sleeper is Unconscious. The task binds each recorded sleep interval to one #487 instance. The source names no waking trigger, so waking requires an explicit trusted fact.
4. **Last-charge destruction (PDF206, 251).** #1140 and #1148 are limited to a wand with no destruction clause, and no other issue rolls for a last charge. The Wand of Magic Missiles is the smallest witness: no Attunement, 1 to 3 charges for levels 1 to 3, 1d6 + 1 at dawn, and a d20 on the last charge. The actual darts stay with #1368. This is also the first owner of a second dawn formula and of the rule that possession alone reveals no charge count.

Wider scope keeps its existing planning owner and is not executable: all other charged or dawn items, item resilience, Artifact destruction, sentient-item limits and attuned charge knowledge with #278; Elf Trance with #333; malnutrition and other hazards with #150; Exhaustion effects with #270.

## Planned native publication

| Task file | Title | Size | Blocked by |
|---|---|---|---|
| `issue-272-knock-out-task.md` | Knock out a creature with a melee attack and start its Short Rest | M | #61, #65, #427, #487 |
| `issue-272-temporary-hit-points-task.md` | Resolve the shared Temporary Hit Point pool and its Long Rest ending | M | #428, #1318 |
| `issue-272-long-rest-sleep-task.md` | Apply the Unconscious condition during Long Rest sleep | S | #428, #487 |
| `issue-272-wand-charges-task.md` | Spend, recharge and destroy a Wand of Magic Missiles by its charges | M | #195, #1140, #1148, #1368 |
| `issue-272-recovery-qualification-task.md` | Qualify the rest, recovery and recharge exception outputs | M | #427, #428 and the four tasks above |

All five are Task children of #272, beside #427 and #428. Only the qualification task touches a gate: it becomes a blocker of item qualification #280 and character qualification #337. Spell qualification #279 already has 100 native blockers, and four direct edges would fill #280. None of the four residuals is a spell effect; #279 blocks #280, so no full spell and item claim can pass #280 without them, and #279 already requires the unchanged #427, #428, #1318 and #1368 through existing paths. No existing body, parent or edge changes apart from those two gate additions. No new direct #67 or #429 consumer is introduced.

## Evidence and limits

The companion scope records the checked source clauses, every owner's body hash and gate paths, the five exact bodies and the planned edges. The graph was read at 2026-11-20T07:26:45Z: #272's only blocker #80 is closed. Counts now and after publication are #279 100 and 100, #280 96 and 97, #337 62 and 63, children of #272 2 and 7, consumers of #428 81 and 84, of #487 87 and 89, of #427 34 and 36. The planned dependency graph and the parent chain #272, #193, #54, #47 stay acyclic, and with #272 closed every new task and both rest tasks still reach #280 and #337. Publication must recheck these counts, because other open splitters may also add #280 blockers.

Current gameplay evidence is the uninterrupted Cleric Short and Long Rest in `internal/rules/srd521/lifecycle.go`; interruptions, Temporary Hit Points, item charges and dawn events have no runtime. This change contains no runtime file, migration, UI behavior, production mutation or Full Rules admission. Runtime implementation remains paused. Focused source and graph checks and hosted CI apply to this decision; the tasks keep their own PostgreSQL and device requirements.
