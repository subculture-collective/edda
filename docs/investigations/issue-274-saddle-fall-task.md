<!-- edda-274:saddle-fall -->
Decision #274; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 15–16, 100; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

On top of the mounted state from {{274:ride}}, implement the PDF16 Falling Off rule and the two PDF100 saddle rules for the same Riding Horse fixture. This task adds no mounting, control, stat block, purchase or barding rule and does not change #920 or #1931.

State ownership: this task adds one fact to the mounted relation that {{274:ride}} owns: which saddle variant from the {{274:records}} records, if any, is fitted on the mount. The fitted saddle is an explicit isolated or GM-declared fact. Buying a saddle stays unavailable until the #275 split publishes a payment owner, and a saddle that is only carried is not fitted.

Falling off: when an effect is about to move the mount against its will while the rider is on it, the rider makes a DC 10 Dexterity saving throw through #61. On a failure the rider falls off, the mounted state ends, and the rider lands with the Prone condition (#488) in an unoccupied space within 5 feet of the mount. The rider makes the same save when knocked Prone while mounted and when the mount is knocked Prone. The landing space must be an authoritative unoccupied space; with none known, the outcome needs GM adjudication and is not auto-placed.

Saddles: a Military Saddle gives Advantage on any ability check the rider makes to remain mounted. The falling-off roll is a saving throw, so the Military Saddle does not affect it. An Exotic Saddle is required to ride an aquatic or flying mount; whether a mount is aquatic or flying is an authoritative fact. The source requires no saddle for any other mount, so do not require a Riding Saddle to ride a horse.

## Deterministic fixtures

- Falling off with injected dice: an effect moving the mount against its will with a save total of 9 leaves the rider Prone and dismounted in a known unoccupied space within 5 feet; a total of 10 keeps the rider mounted. The same save is called when the rider is knocked Prone while mounted and when the mount is knocked Prone. Willing movement of the mount calls no save. With no known unoccupied space within 5 feet the failure is not auto-placed.
- Military Saddle: a GM-called ability check to remain mounted rolls with Advantage; the same check without that saddle, with a Riding Saddle or with the saddle carried but not fitted rolls normally. A Military Saddle adds no Advantage to the falling-off save.
- Exotic Saddle, with an isolated fixture mount: a mount recorded as flying or aquatic cannot be ridden with no saddle, a Riding Saddle or a Military Saddle and can with an Exotic Saddle. A Riding Horse can be ridden bareback. A mount whose aquatic or flying status is unknown is unavailable.
- Negative, each violated independently: a save for a rider who is not mounted; a DC other than 10; a Strength or ability-check roll in place of the Dexterity save; the rider left mounted after a failure; a landing more than 5 feet away or in an occupied space; a purchased or free saddle appearing without the declared fact.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, encounter, rider, mount, mounted state, fitted saddle, positions, source version and expected revisions. Commit injected dice, the end of the mounted state, the condition, position and public history atomically. A forced save belongs to its triggering event, with no invented optional activation. Inject failures after each stage and prove full isolated PostgreSQL rollback. Same action ID and payload replay after retry, reload or reconnect returns the same receipt without rerolling; a conflicting payload rejects; racing confirmations commit once. The fitted saddle survives reload. Web and trusted TUI share availability, unavailable reasons and confirmation/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private checks, hidden creatures and private mount facts.

Hard prerequisites: {{274:ride}} mounted state, {{274:records}} saddle records, #61 checks and saves, #488 Prone. It consumes the merged #67 scene facts for occupancy and adds no #67 edge, because #67 is at the 100-consumer cap. Parent: #274. The cohort qualification {{274:cohort}} requires this task and is the direct blocker of #280 and #183; add no gate edge here.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
