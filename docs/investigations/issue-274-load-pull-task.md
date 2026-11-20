<!-- edda-274:load-pull -->
Decision #274; roadmap #47. One implementation PR; Task / Size M.

Source: English SRD 5.2.1, PDF 92, 100; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement three PDF100 rules for the eight animals in the Mounts and Other Animals table, using the records from {{274:records}}: carried load against table carrying capacity, the Mounts and Cargo pull limit for a carriage, cart, chariot, sled or wagon, and the Barding cost and weight derivation. No other creature, vehicle or load rule is enabled.

Carried load: an owned animal identified by one of the eight table rows can carry gear whose total authoritative weight does not exceed that row's carrying capacity. Use the table value as written. Do not derive a capacity from Strength, size or the Mule's Beast of Burden trait; the general Carrying Capacity rule on PDF178 is not owned here. An item, rider or creature without an authoritative weight cannot be loaded. A load above capacity is unavailable; the source gives this table no overload penalty, so invent none.

Pull: an animal pulling one of the five drawn vehicles can move weight up to five times its table carrying capacity, and that weight includes the vehicle's own weight. When several animals pull the same vehicle, add their carrying capacities first and then multiply by five. An animal is either carrying or hitched in one committed arrangement; the source does not say how a pack load combines with pulling, so that combination is unavailable. Team size limits, harness items and terrain are not in the source and are not invented. The PDF12 rule that travelers in land vehicles choose a pace as normal, and the doubled mounted distance, are travel-pace variants that stay with #150.

Barding: for a qualified Light, Medium or Heavy armor record from the PDF92 Armor table (#251, #252), the barding variant costs four times the armor's cost and weighs twice its weight. This task derives the priced, weighted variant and counts its weight as carried load. It does not apply an Armor Class, Strength or Stealth entry, training requirement or don/doff time to a mount; the source states none for mounts, so wearing barding stays unavailable. Shield barding is not adjudicated and stays unavailable. Purchase is not committed here.

## Deterministic fixtures

- Mule, capacity 420 lb.: 420 lb. of known-weight gear commits; 421 lb. is unavailable; an item of unknown weight is unavailable. The Mule's 420 lb. is used as written and is not recomputed from its stat block or its Beast of Burden trait.
- One Draft Horse (540 lb.) hitched to a Wagon (400 lb.): limit 2,700 lb. including the wagon, so 2,300 lb. of cargo commits and 2,301 lb. is unavailable. Two Draft Horses: (540 + 540) × 5 = 5,400 lb., so 5,000 lb. of cargo commits and 5,001 lb. is unavailable. A Mule and a Pony on one Cart (200 lb.): (420 + 225) × 5 = 3,225 lb. including the cart.
- Negative, each violated independently: multiplying each animal by five and also adding a bonus; excluding the vehicle's weight; applying the pull limit to a pack load or to a Rowboat; hitching an animal that is not one of the eight rows; hitching one animal to two vehicles; a foreign or unowned animal or vehicle; a stale load revision.
- Chain Mail (75 GP, 55 lb.) derives barding at 300 GP and 110 lb. Leather Armor (10 GP, 10 lb.) derives 40 GP and 20 lb. A Warhorse (540 lb.) carrying Chain Mail barding has 430 lb. of capacity left. Reject barding at the base cost or weight, a doubled cost, a quadrupled weight, Shield barding, an AC change on the mount and barding from an armor record that is not qualified.
- A Rowboat counts as 100 lb. of carried load from its record.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, actor, ownership of each animal, vehicle and item, source version, every weight and expected revisions. Commit the load or hitch arrangement, derived totals and public history atomically. Inject failures after each stage and prove full isolated PostgreSQL rollback. Same action ID and payload replay after retry, reload or reconnect returns the same receipt; a conflicting payload rejects; racing confirmations commit once. Web and trusted TUI share availability, unavailable reasons and confirmation/cancel/reconnect behavior. Changed web requires focused build/lint and desktop/phone checks. Public history omits private item facts.

Hard prerequisites: {{274:records}} table records, #251 and #252 Armor-table records. It adds no #67 or #429 edge. Parent: #274. The cohort qualification {{274:cohort}} requires this task and is the direct blocker of #280; add no #280 edge here. Do not alter #1725 Backpack capacity or any existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until implementation and its callers qualify. No deployment, source inference or Full Rules promotion.
