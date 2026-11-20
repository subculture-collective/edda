<!-- edda-274:riding -->
Decision #274; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 15–16, 100, 360; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement the mounted state for one fixture mount, the Riding Horse, under three PDF15–16 Mounted Combat rules: eligibility, mounting and dismounting, and controlling a mount. The Riding Horse acts with the executable statistics owned by #1999 (PDF360: Large Beast, Speed 60 ft., Initiative +1); #171 owns its catalogue record. This task adds no stat block and makes no other creature an acting mount. Falling off and the two saddle rules belong to #2077. Purchase, travel pace, carried load, barding and the Lance stay outside.

State ownership: this task owns the mounted relation between a rider and a mount that are both participants in one encounter. The mount enters the encounter as an explicit isolated or GM-declared creature. This task does not create an owned-animal instance and does not link one to the encounter mount; #2075 owns the owned-animal instance, and joining the two is not stated by any published task and stays unavailable. #920 (Find Steed) and #1931 (Phantom Steed) specify control and dismounting for their own steeds and are not changed. No native edge joins them to this task because their bodies list exact prerequisites and cannot be edited. If either has merged a mounted or controlled-mount state first, extend that state here and do not add a second one.

Eligibility: a willing creature at least one size larger than the rider, with an appropriate anatomy, can serve as a mount. Size comes from qualified statistics; willingness and anatomy are authoritative facts, and an unknown fact leaves the action unavailable.

Mounting and dismounting: during the rider's own move, mount a creature within 5 feet or dismount. Either costs movement equal to half the rider's Speed, rounded down as PDF15 states. Use the rider's current authoritative Speed. No action, Bonus Action or check is charged.

Control: a mount can be controlled only if it has been trained to accept a rider. PDF16 says "Domesticated horses, mules, and similar creatures have such training." The source does not list the similar creatures. Of the PDF100 table animals, this task reads Horse, Draft; Horse, Riding; Warhorse; and Mule as domesticated horses and mules. Camel, Elephant, Mastiff, Pony and every other creature are trained only by an authoritative fact. Only the Riding Horse has executable statistics here, so the reading enables no other mount. When a rider mounts a controlled mount, its Initiative changes to match the rider's. It moves on the rider's turn as directed and has only three action options on that turn: Dash, Disengage and Dodge, resolved by #780, #781 and #782. It can move and act on the turn it is mounted. An independent mount keeps its own place in the Initiative order and acts through its own existing turn owner; rider-directed actions for it reject.

## Deterministic fixtures

- A Medium rider with Speed 30 mounts an adjacent willing Riding Horse for 15 feet of movement. Speed 25 costs 12. With 14 feet of movement left at Speed 30 the mount is unavailable. Dismounting costs the same. A mount 10 feet away, an unwilling creature, an unknown anatomy, an action or Bonus Action charge, and mounting on another creature's turn are rejected independently.
- Size, with isolated fixture creatures: a Medium rider cannot mount a Medium creature. A Small rider can when the other facts hold. A willing creature with no recorded training can be ridden only as an independent mount.
- Control: on mounting, the horse's Initiative becomes the rider's. On that turn it may Dash, Disengage or Dodge and move 60 feet as directed. Reject an attack or any fourth action option for a controlled mount, rider-directed movement of an independent mount, and a Speed other than the 60 feet from #1999.
- Negative: a second mounted state beside one already merged for a Find Steed or Phantom Steed; a catalogue record alone used as an acting mount; an owned animal from #2075 treated as the encounter mount.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, encounter, actor, mount identity, size, willingness, training, positions, remaining movement, source version and expected revisions. Commit the mounted state, movement cost, Initiative change and public history atomically. Inject failures after each stage and prove full isolated PostgreSQL rollback. Same action ID and payload replay after retry, reload or reconnect returns the same receipt; a conflicting payload rejects; racing confirmations commit once. The mounted state survives reload. Web and trusted TUI share availability, unavailable reasons and confirmation/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits hidden creatures and private mount facts.

Hard prerequisites: #2074 mount-table row for the Riding Horse, #66 initiative and movement economy, #171 Riding Horse record, #1999 executable Riding Horse statistics, #780 Dash, #781 Disengage, #782 Dodge. It consumes the merged #67 scene facts for distance and adds no #67 edge, because #67 is at the 100-consumer cap. Parent: #274. The cohort qualification #2080 requires this task and is the direct blocker of #280 and #183; add no gate edge here.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
