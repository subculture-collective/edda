<!-- edda-274:riding -->
Decision #274; roadmap #47. One implementation PR; Task / Size L.

Source: English SRD 5.2.1, PDF 5, 15–16, 100, 360; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement riding for one fixture mount, the Riding Horse (Large, Speed 60 ft., stat block from #171), with the PDF15–16 Mounted Combat rules and the two PDF100 saddle rules. The saddle rules cannot be evaluated without a mounted state, and no published task owns PDF15–16, so this task owns both. It adds no mount stat block, purchase, travel pace, carrying rule, barding effect or Lance handling. #920 Find Steed control and #1931 Phantom Steed keep their own contracts and are not changed.

Eligibility: a willing creature at least one size larger than the rider, with an appropriate anatomy, can serve as a mount. Size comes from qualified records; willingness and anatomy are authoritative facts, and an unknown fact leaves the action unavailable.

Mounting and dismounting: during the rider's own move, mount a creature within 5 feet or dismount. Either costs movement equal to half the rider's Speed, rounded down (PDF5). Use the rider's current authoritative Speed. No action, Bonus Action or check is charged.

Control: a mount can be controlled only if it has been trained to accept a rider. Domesticated horses and mules have that training; for any other creature it must be an authoritative fact. When a rider mounts a controlled mount, its Initiative changes to match the rider's. It moves on the rider's turn as directed and has only three action options on that turn: Dash, Disengage and Dodge, resolved by #780, #781 and #782. It can move and act on the turn it is mounted. An independent mount keeps its own place in the Initiative order and acts through its own existing turn owner; rider-directed actions for it reject.

Falling off: when an effect is about to move the mount against its will while the rider is on it, the rider makes a DC 10 Dexterity saving throw through #61. On a failure the rider falls off and lands with the Prone condition (#488) in an unoccupied space within 5 feet of the mount. The rider makes the same save when knocked Prone while mounted and when the mount is knocked Prone. The landing space must be an authoritative unoccupied space; with none known, the outcome needs GM adjudication and is not auto-placed.

Saddles, from the {{274:records}} records: a Military Saddle gives Advantage on any ability check the rider makes to remain mounted. The falling-off roll is a saving throw, so the Military Saddle does not affect it. An Exotic Saddle is required to ride an aquatic or flying mount; whether a mount is aquatic or flying is an authoritative record fact. The source requires no saddle for any other mount, so do not require a Riding Saddle to ride a horse.

## Deterministic fixtures

- A Medium rider with Speed 30 mounts an adjacent willing Riding Horse for 15 feet of movement. Speed 25 costs 12. With 14 feet of movement left at Speed 30 the mount is unavailable. Dismounting costs the same. A mount 10 feet away, an unwilling creature, an unknown anatomy, an action or Bonus Action charge, and mounting on another creature's turn are rejected independently.
- Size, with isolated fixture creatures: a Medium rider cannot mount a Medium creature. A Small rider can when the other facts hold. A willing creature with no recorded training can be ridden only as an independent mount.
- Control: on mounting, the horse's Initiative becomes the rider's. On that turn it may Dash, Disengage or Dodge and move 60 feet as directed. Reject an attack or any fourth action option for a controlled mount, and rider-directed movement of an independent mount.
- Falling off with injected dice: an effect moving the mount against its will with a save total of 9 leaves the rider Prone in a known unoccupied space within 5 feet; a total of 10 keeps the rider mounted. The same save is called when the rider is knocked Prone while mounted and when the mount is knocked Prone. Willing movement of the mount calls no save. A Military Saddle adds no Advantage to this save.
- Military Saddle: a GM-called ability check to remain mounted rolls with Advantage; the same check without that saddle, with a Riding Saddle or with the saddle carried but not on the mount rolls normally.
- Exotic Saddle: a mount recorded as flying or aquatic cannot be ridden with no saddle, a Riding Saddle or a Military Saddle and can with an Exotic Saddle. A Riding Horse can be ridden bareback. A mount whose aquatic or flying status is unknown is unavailable.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, encounter, actor, mount identity, size, willingness, training, saddle, positions, remaining movement, source version and expected revisions. Commit the mounted state, movement cost, Initiative change, injected dice, condition and public history atomically. A forced save belongs to its triggering event, with no invented optional activation. Inject failures after each stage and prove full isolated PostgreSQL rollback. Same action ID and payload replay after retry, reload or reconnect returns the same receipt without rerolling; a conflicting payload rejects; racing confirmations commit once. The mounted state survives reload. Web and trusted TUI share availability, unavailable reasons and confirmation/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private checks, hidden creatures and private mount facts.

Hard prerequisites: {{274:records}} saddle records, #61 checks and saves, #66 initiative and movement economy, #171 Riding Horse record, #488 Prone, #780 Dash, #781 Disengage, #782 Dodge. It consumes the merged #67 scene facts for distance and occupancy and adds no #67 edge, because #67 is at the 100-consumer cap. Parent: #274. The cohort qualification {{274:cohort}} requires this task and is the direct blocker of #280; add no #280 edge here.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
