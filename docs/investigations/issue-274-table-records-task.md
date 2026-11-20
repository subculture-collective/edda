<!-- edda-274:table-records -->
Decision #274; roadmap #47. One data-only PR; Task / Size M.

Source: English SRD 5.2.1, PDF 100–101; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Add source-linked records for the three equipment tables that the #80 inventory left outside its 753 indexed families. Use the #195 versioned item-record contract: stable ID, per-record page, pin, attribution and explicit required mechanic keys. This PR adds data and offline validation only. It adds no purchase flow, inventory grant, creature record, travel rule or resolver, and a record's presence makes nothing playable.

Mounts and Other Animals, PDF100, eight rows as carrying capacity / cost: Camel 450 lb. / 50 GP; Elephant 1,320 lb. / 200 GP; Horse, Draft 540 lb. / 50 GP; Horse, Riding 480 lb. / 75 GP; Mastiff 195 lb. / 25 GP; Mule 420 lb. / 8 GP; Pony 225 lb. / 30 GP; Warhorse 540 lb. / 400 GP. Each row carries an opaque reference to its stat-block source key. The stat blocks stay with the monster catalogue batches #156 Camel, #158 Draft Horse, #159 Elephant, #167 Mastiff, #168 Mule, #170 Pony, #171 Riding Horse and #176 Warhorse. Do not copy ability scores, size, Speed or traits here and do not recompute a capacity from them.

Tack, Harness, and Drawn Vehicles, PDF100, ten priced rows as weight / cost: Carriage 600 lb. / 100 GP; Cart 200 lb. / 15 GP; Chariot 100 lb. / 250 GP; Feed per day 10 lb. / 5 CP; Saddle, Exotic 40 lb. / 60 GP; Saddle, Military 30 lb. / 20 GP; Saddle, Riding 25 lb. / 10 GP; Sled 300 lb. / 20 GP; Stabling per day no weight / 5 SP; Wagon 400 lb. / 35 GP. Saddle is one family with exactly three explicit variants. A saddle includes its bit, bridle and reins; do not create those as separate records.

Airborne and Waterborne Vehicles, PDF101, seven rows as speed / crew / passengers / cargo tons / AC / HP / damage threshold / cost: Airship 8 mph / 10 / 20 / 1 / 13 / 300 / none / 40,000 GP; Galley 4 mph / 80 / none / 150 / 15 / 500 / 20 / 30,000 GP; Keelboat 1 mph / 1 / 6 / 1/2 / 15 / 100 / 10 / 3,000 GP; Longship 3 mph / 40 / 150 / 10 / 15 / 300 / 15 / 10,000 GP; Rowboat 1½ mph / 1 / 3 / none / 11 / 50 / none / 50 GP; Sailing Ship 2 mph / 20 / 20 / 100 / 15 / 300 / 15 / 10,000 GP; Warship 2½ mph / 60 / 60 / 200 / 15 / 500 / 20 / 25,000 GP. Store speeds and cargo as exact rationals. A source dash is an explicit absent value, never zero. PDF100 also gives the Rowboat a carried weight of 100 pounds; record it on that row only.

Record the two usual passenger prices on PDF101 as source notes marked as usual and variable: a hammock 5 SP per day and a small private cabin 2 GP per day. Record the required mechanic keys that later tasks own: table carrying capacity, drawn-vehicle pull, barding derivation, saddle benefit and requirement, vessel speed, crew, passengers, cargo, damage threshold and ship repair. Every key stays unavailable in this PR.

## Deterministic fixtures

- Positive: exactly 8 mount rows, 10 tack and vehicle rows in 8 families with 3 saddle variants, and 7 large-vehicle rows; every value above matches the page it cites; IDs are unique and stable; each mount row names one stat-block key.
- Negative, each violated independently: a missing or extra row; a fourth saddle variant; a transposed value such as Draft Horse 480 lb. or Riding Horse 540 lb.; 5 SP for Feed or 5 CP for Stabling; a dash stored as 0; Rowboat speed stored as 1 or 2 instead of 1½; Keelboat cargo stored as 0 or 1; Galley passengers stored as 0; a Rowboat or Airship damage threshold; a weight for Stabling; a carried weight for any ship other than the Rowboat; a record without page, pin or attribution; a duplicate ID.
- A capacity derived from a stat block, a ton converted to pounds, a price converted between coins, a default wind or current, and a passenger price stored as a fixed fare are all rejected. The source gives none of them as a table value.

## Acceptance

Validation runs offline against the committed data with no model call, network or production data. The catalogue loader fails closed on unknown fields, units and mechanic keys. No runtime entitlement, shop, starting-equipment grant, coverage promotion or UI affordance changes; web and trusted TUI continue to show these families as unavailable. Compare every row with its cited PDF page before closing.

Hard prerequisite: #195 item-record contract. Parent: #274. The cohort qualification #2080 requires this task and is the direct blocker of item qualification #280 and exploration qualification #183; add no gate edge here. Buying, selling and coin handling stay with the #275 economy split, and creature statistics stay with the monster catalogue and #182.

Validation for the implementation PR: offline source, schema, name, variant and value fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. No deployment, source inference or Full Rules promotion.
