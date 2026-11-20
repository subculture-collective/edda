<!-- edda-274:vessel-hour -->
Decision #274; roadmap #47. One implementation PR; Task / Size L.

Source: English SRD 5.2.1, PDF 5, 12, 100–102, 192; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement one committed hour of travel aboard one of the seven large vehicles recorded by {{274:records}}, on top of the #144 one-hour travel segment. This task owns the PDF100–101 Large Vehicles notes on speed, crew, passengers and cargo and the PDF12 and PDF192 sentences that put travelers aboard a vessel at the vessel's speed with no pace choice. Land-vehicle pace and the doubled mounted distance stay with #150. Damage and repair belong to {{274:hull}}. Hiring and paying a crew belong to the #275 economy split.

Speed: travelers aboard use the vessel's table speed in miles per hour and do not choose a travel pace, so no pace modifier applies. Distances keep exact halves. A ship sailing against a strong wind moves at half speed. In a dead calm a waterborne ship cannot move under sail and must be rowed; the source gives no rowed speed for any ship, so a dead-calm hour is unavailable. Keelboats and Rowboats are used on lakes and rivers. For those two, going downstream adds the current's speed; the current must be a declared authoritative fact, and the typical 3 miles per hour is offered only as an explicit GM choice. They cannot be rowed against any significant current. The source attaches the current rule to those two vessels, so a current on any other vessel is unavailable. Being pulled upstream by draft animals has no source speed and stays unavailable. Whether the Airship counts as sailing is not stated, so an Airship hour against a strong wind is unavailable; with no declared wind it moves 8 miles.

Crew: a ship larger than a Keelboat or Rowboat needs at least its table number of skilled hirelings to function. Count only authoritative crew members recorded as skilled hirelings through the output that the #275 split assigns; until that output is qualified, the five crewed ships stay unavailable. The Keelboat and Rowboat need their crew of one and no hireling.

Passengers: the table number is Small and Medium passengers in hammocks. A ship outfitted with private accommodations carries one-fifth as many, rounded down (PDF5). The outfitting is an authoritative ship fact. The usual 5 SP hammock and 2 GP cabin prices are records only; no fare is charged here.

Cargo: the table gives cargo in tons. Declared cargo must be an authoritative amount in tons. The source does not define a ton in pounds, so converting item weights to tons is unavailable.

Hours: PDF12 says ships might be able to travel up to 24 hours per day depending on the vessel and crew. Hours beyond the campaign's ordinary travel day need an explicit GM-declared continuous-sailing fact, and a 25th hour in one day always rejects.

## Deterministic fixtures

- Sailing Ship, 20 skilled hirelings, no adverse wind: one hour moves 2 miles and offers no pace. Against a strong wind it moves 1 mile. In a dead calm the hour is unavailable, with rowing named as the reason. With 19 skilled hirelings, or 20 crew of whom one is untrained, it cannot sail.
- Warship: 2½ miles, and 1¼ miles against a strong wind. Galley: 80 crew required; 79 rejects; any passenger rejects because its passenger entry is absent.
- Keelboat, crew 1, declared 3 mph current downstream: 4 miles. With no declared current the downstream hour is unavailable and never assumes 3. Upstream against a significant current rejects. Rowboat, crew 1, still water: 1½ miles.
- Passengers: Longship 150 in hammocks, 151 rejects; with private accommodations 30, and 31 rejects. Sailing Ship 20 and 4. Keelboat 6 and 1. Rowboat 3 and 0.
- Cargo: Sailing Ship 100 tons commits and 100½ rejects. Keelboat ½ ton commits. Rowboat and any declared cargo tonnage rejects because its cargo entry is absent. Cargo given only in pounds is unavailable.
- Negative, each violated independently: a Fast or Slow pace aboard; a pace Disadvantage or Advantage aboard; halving for a wind that is not strong or not against; a mounted or land-vehicle segment routed through this task; an unknown route distance; a foreign vessel; a stale clock or route revision.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, vessel identity and record, crew roster, passenger and cargo counts, declared wind and current, route, clock, source version and expected revisions. Commit position, clock, distance and public history once in the #144 travel transaction. Inject failures after each stage and prove full isolated PostgreSQL rollback. Same action ID and payload replay after retry, reload or reconnect returns the same receipt; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, unavailable reasons and confirmation/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits hidden route facts and private crew or cargo facts.

Hard prerequisites: {{274:records}} vessel records, #144 travel segment, #275 for the skilled-hireling crew output its split assigns. Parent: #274. The cohort qualification {{274:cohort}} requires this task and is the direct blocker of #280; add no #280 edge here. Do not alter #144 or #150.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
