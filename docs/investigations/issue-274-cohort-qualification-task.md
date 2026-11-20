<!-- edda-274:cohort-qualification -->
Decision #274; roadmap #47. One qualification PR; Task / Size M.

Source: English SRD 5.2.1, PDF 5, 12, 15–16, 92, 100–102, 180, 192; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

Qualify the five #274 outputs as one cohort for the item gate. This task owns qualification only. It adds no record, rule, resolver, purchase flow or production grant, and it copies nothing from its inputs. Item qualification #280 had 96 of 100 blocked-by entries when #274 was decided, with other splitters still to publish, so this single task carries the whole cohort to #280 on one edge.

Exact roster: {{274:records}} table records; {{274:load}} carried load, pull limit and barding derivation; {{274:ride}} riding, control, falling off and saddles; {{274:voyage}} large-vehicle travel hour; {{274:hull}} damage threshold and ship repair. Each input must have merged and passed its own acceptance. A closed planning issue, a catalogue record or a published task is not evidence.

Required witnesses, run through the actual rules, game, engine and handler callers with isolated explicit entitlements:

- Records: the 8 mount rows, 10 tack and vehicle rows and 7 large-vehicle rows still match PDF100–101, and every mechanic key they declare has a qualified resolver from this roster or is reported as unavailable with its named reason.
- Land: two Draft Horses pull a Wagon with 5,000 lb. of cargo and reject 5,001 lb.; a Warhorse carrying Chain Mail barding (300 GP, 110 lb.) has 430 lb. of capacity left.
- Riding: a Medium rider with Speed 30 mounts a Riding Horse for 15 feet, directs a Dash, fails a DC 10 Dexterity save with a total of 9 and lands Prone; a Military Saddle does not change that save.
- Water: a Sailing Ship with 20 skilled hirelings moves 2 miles in an hour and 1 mile against a strong wind; it takes nothing from a 14-damage instance and 15 from a 15-damage instance; berthed in a shipyard it repairs 15 HP for 7½ days and 150 GP.
- Limits stay unavailable and are reported as such: overload penalties, barding worn by a mount, Shield barding, rowed and towed speeds, an Airship against a strong wind, pounds-to-tons conversion, Poison or Psychic damage to a vessel, any consequence at 0 Hit Points, fares, and purchase of any record.

Report data completeness separately from mechanical support for the mounts and vehicles section that #280 names. Preview is inert. Confirmation revalidates actor, campaign, source version, ownership, authoritative facts and expected revisions. Inject each failure for full isolated PostgreSQL rollback; same-ID retry, reload and reconnect return the saved receipt without rerolling or charging again; conflicts reject; racing confirmations commit once. Web and trusted TUI must agree on supported choices, unavailable reasons and cancel/confirm/reconnect. Changed web requires focused build/lint and desktop/phone fixtures. Do not expose hidden or private state.

Native parent #274; directly block #280. Exact prerequisites: {{274:records}}, {{274:load}}, {{274:ride}}, {{274:voyage}}, {{274:hull}}. Preserve every existing body, parent and prerequisite. Add no #279, #183, #67 or #429 edge. Closing #274 removes none of these obligations. Qualification must run focused source and contrary tests, isolated PostgreSQL replay/rollback/concurrency, shared web/TUI fixtures, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`; short tests alone do not qualify database or devices. No deployment, production grant or Full Rules promotion.
