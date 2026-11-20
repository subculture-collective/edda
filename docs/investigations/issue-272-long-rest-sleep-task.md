<!-- edda-272:long-rest-sleep-unconscious -->
Decision #272; roadmap #47. One implementation PR; Task / Size S.

Source: English SRD 5.2.1, PDF 107, 185, 191, 253; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. CC-BY-4.0 attribution and adaptation: `docs/srd-5.2.1-conformance.md`. The deterministic fixtures below are Edda adaptations.

## Bounded outcome

Implement only this Long Rest clause: during sleep, the resting creature has the Unconscious condition. The canonical Long Rest #428 already owns the 8-hour minimum, at least 6 hours of sleep, at most 2 hours of light activity, eligibility, interruption, resumption and recovery. This task binds each recorded sleep interval of that rest to one source-bound Unconscious instance from the shared output #487. It copies neither contract.

While a sleep interval of an active Long Rest is in effect, the creature has one Unconscious instance whose source is that rest. Light-activity intervals, such as standing watch, carry none. The consequences are exactly #487's output, including Incapacitated, Prone, Speed 0, attack and saving-throw effects, critical hits from within 5 feet, being unaware, and the Concentration ending that Incapacitated already causes. Ending the instance leaves the creature Prone.

The instance ends when the authoritative rest record ends that sleep interval: at the recorded end of sleep, at completion of the rest, or when a confirmed #428 interruption or stop event records that the sleeper woke. The Long Rest (PDF185) and Unconscious (PDF191) entries name no waking trigger. Noise, an attack or damage therefore do not by themselves remove the instance; the attack resolves against an Unconscious creature, and #428 decides whether the rest is interrupted. A waking fact must be explicit and trusted; otherwise it is unavailable. Source-specific wakers keep their owners and supply that trusted fact when they are qualified: the Alarm spell's mental ping, which awakens its caster (PDF107, #1035), and the Weapon of Warning, which awakens each naturally sleeping subject when combat begins (PDF253, retained with #278). This task accepts such a fact and implements neither effect. A Short Rest applies no sleep instance. A feature that replaces sleep, such as Elf Trance (PDF85), and magical sleep keep their own owners and instances.

## Deterministic fixtures

- A Long Rest from T0 with sleep recorded T0 to T0 + 6 hours and watch for the final 2 hours: at T0 + 3 hours the creature has exactly one rest-sourced Unconscious instance; at T0 + 7 hours it has none and is Prone until it stands; at T0 + 8 hours #428 completes unchanged.
- At T0 + 3 hours a melee hit from within 5 feet uses #487's Advantage and critical-hit output; the damage is a #428 interruption. The instance ends only with the recorded waking fact of that confirmed interruption.
- A trusted waking fact at T0 + 2 hours, such as the one an Alarm ping from #1035 supplies, ends the instance at that time and leaves the creature Prone; whether the rest continues is #428's decision.
- A sleeper concentrating on a spell loses Concentration through the existing #487/#353 output when the sleep interval begins; no separate rule is added here.
- Reject an Unconscious instance during a watch interval or a Short Rest; a rest that completes with less recorded sleep than #428 requires; waking from noise, damage or a client message without a trusted waking fact; a Weapon of Warning or Alarm waking that no qualified owner produced; a Perception penalty or other effect that #487 does not output; removal of a Sleep-spell, knock-out or 0 Hit Point Unconscious instance when the sleeper wakes; a duplicate instance on replay; and treating Trance or another sleep replacement as sleep.

## Durable implementation acceptance

Preview is inert and never rolls, writes or reveals private facts. Confirm revalidates campaign, actor, source version, clock, targets, the rest record, its sleep intervals, the waking fact and expected revisions. Commit the sleep interval, the condition instance or its ending, derived state and public history atomically. Inject a failure after each stage and prove full isolated PostgreSQL rollback. The same action ID and payload after retry, reload or reconnect returns the same receipt with no reroll; a conflicting payload rejects; racing confirmations commit once. Applying and ending the instance belong to the triggering confirmed rest event, with no optional activation. Web and trusted TUI share availability, unavailable reasons and confirm/cancel/reconnect behavior. Changed web requires a focused build, lint and desktop/phone checks. Public history omits private HP, hidden actors, private checks and undisclosed item facts.

Hard prerequisites: #428 canonical Long Rest, #487 shared Unconscious output. Parent: #272. Block the rest/recovery qualification task published with this decision; that task blocks #280. Do not add a direct #279, #280, #67 or #429 edge and do not alter any existing owner.

Validation for the implementation PR: focused deterministic rules tests; isolated PostgreSQL rollback/concurrency/replay/reload tests; shared web/TUI confirmation and disclosure fixtures; `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests alone do not qualify PostgreSQL or devices. Runtime remains unavailable until this task and its callers qualify. No deployment, source inference or Full Rules promotion.
