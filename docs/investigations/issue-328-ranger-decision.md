# Ranger / Hunter source and implementation reconciliation (#328)

This is a decision and publication change. Ranger and Hunter remain unavailable.
No runtime, grants, coverage, fixed Fighter/Cleric fixtures or production state
are changed. Final native publication and immutable parent seal are pending.

Prepared base: `3042172c935bdefa606c008c8ddda277b7548fd8` (accepted Monk #843).
The pinned SRD5.2.1 PDF SHA256 is
`8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`.
PDF57-59 Ranger,60 full spell list,61 Hunter were extracted directly. Eight
independent PDF comparators cover Ranger, Hunter, advancement, multiclass and
all four feat catalogues. CC-BY-4.0 attribution/adaptations follow
[the source register](../srd-5.2.1-conformance.md).

Canonical counts checked: **72 Ranger +21 Hunter =93 records**. The JSON matrix
retains each exact serialized source record and all216 contextual mechanic and
option key occurrences, with rules/game/engine/combat/database/handler/web/TUI
findings. Repeated shared keys retain their source contexts. Every named feature
is owned by its execution leaf, with a native prerequisite path from its granting
row; ASI and selectors do not substitute for behavior.

`creation.go` admits only fixed Human Soldier Fighter / Human Acolyte Cleric.
`advancement.go` and `advancement_preview.go` implement Fighter1→2, not Ranger.
Database storage and the web sheet cannot admit Ranger. The TUI character view
is a placeholder; its launcher enters the shared world flow, without Ranger key
consumption.
Source catalogue and label presence therefore remain metadata evidence only.

Class tables use **Ranger level** for features, preparations, slots and Favored
Enemy counts. XP/livePB independently select the COMPLETE canonical advancement
table by sum of authoritative admitted class levels AFTER transition. Table PB
is single-class context only. Fighter4+Ranger1 reaches total5/6500XP/PB3 with
later-class D10/fixed6 HP. Fighter4/Ranger1→Ranger2 reaches total6/14000XP/PB3
with class2 features. Unsupported combinations and totals>20 stay unavailable.

The published prerequisite gateway [#968](https://github.com/subculture-collective/edda/issues/968)
is Task/M, parent336, sole blocker429. Readback preserves429's full body,
prerequisites and99old outgoing edges;968 fills slot100. It qualifies the
existing protocol and owns no transaction, die or class eligibility. Own Ranger
HP adapter consumes968; every level/later-entry contract consumes ownHP.
Monk783 and Paladin851 retain their paths. No more direct429consumers are legal.

The exact source spell roster is48Ranger spells plus11Druid cantrips. Existing
bounded spell owners are reused; FOURCLASS membership503 is not Ranger membership.
New spell contracts split lifecycle, movement, damage and world disclosures.
Hunter's Mark uses attack-roll hits/Force damage, not the old weapon-only variant;
transfer requires target0HP before expiry and spends a BonusAction, not a new cast.
All feature fixtures require isolated entitlements and grant no production access.

## All twenty class rows

| Ranger level | Class-context PB | Prepared | Favored Enemy | Slots 1/2/3/4/5 | Features |
|---|---|---|---|---|---|
| 1 | +2 | 2 | 2 | 2/0/0/0/0 | Spellcasting, Favored Enemy, Weapon Mastery |
| 2 | +2 | 3 | 2 | 2/0/0/0/0 | Deft Explorer, Fighting Style |
| 3 | +2 | 4 | 2 | 3/0/0/0/0 | Ranger Subclass |
| 4 | +2 | 5 | 2 | 3/0/0/0/0 | Ability Score Improvement |
| 5 | +3 | 6 | 3 | 4/2/0/0/0 | Extra Attack |
| 6 | +3 | 6 | 3 | 4/2/0/0/0 | Roving |
| 7 | +3 | 7 | 3 | 4/3/0/0/0 | Hunter subclass feature |
| 8 | +3 | 7 | 3 | 4/3/0/0/0 | Ability Score Improvement |
| 9 | +4 | 9 | 4 | 4/3/2/0/0 | Expertise |
| 10 | +4 | 9 | 4 | 4/3/2/0/0 | Tireless |
| 11 | +4 | 10 | 4 | 4/3/3/0/0 | Hunter subclass feature |
| 12 | +4 | 10 | 4 | 4/3/3/0/0 | Ability Score Improvement |
| 13 | +5 | 11 | 5 | 4/3/3/1/0 | Relentless Hunter |
| 14 | +5 | 11 | 5 | 4/3/3/1/0 | Nature’s Veil |
| 15 | +5 | 12 | 5 | 4/3/3/2/0 | Hunter subclass feature |
| 16 | +5 | 12 | 5 | 4/3/3/2/0 | Ability Score Improvement |
| 17 | +6 | 14 | 6 | 4/3/3/3/1 | Precise Hunter |
| 18 | +6 | 14 | 6 | 4/3/3/3/1 | Feral Senses |
| 19 | +6 | 15 | 6 | 4/3/3/3/2 | Epic Boon |
| 20 | +6 | 15 | 6 | 4/3/3/3/2 | Foe Slayer |

## Exact bounded owners and native routing

The following is a proposed manifest, not live publication proof. Native IDs,
Task/Size readback and full graph receipt will replace this state after reservation.

- `class-gate`: Qualify finite Ranger class gate. Parent `328`; needs class-cohort, levels-cohort, shared-gate, #328.
- `class-cohort`: Qualify finite Ranger class cohort. Parent `328`; needs hp, core, skills, equipment, later-entry, spellcasting, favored-enemy, weapon-mastery, deft-explorer, fighting-style, ranger-subclass, ability-score-improvement, extra-attack, roving, expertise, relentless-hunter, natures-veil, precise-hunter, feral-senses, epic-boon, foe-slayer, tireless-hp, tireless-rest, slots, membership, druidic-cantrips, hunter-identity, hunters-lore, colossus-slayer, horde-breaker, escape-the-horde, multiattack-defense, superior-hunters-prey, superior-hunters-defense, hunters-prey, defensive-tactics.
- `levels-cohort`: Qualify finite Ranger levels cohort. Parent `328`; needs level-1, level-2, level-3, level-4, level-5, level-6, level-7, level-8, level-9, level-10, level-11, level-12, level-13, level-14, level-15, level-16, level-17, level-18, level-19, level-20.
- `shared-gate`: Qualify finite Ranger shared gate. Parent `328`; needs spell-cohort-1, spell-cohort-2, spell-cohort-3, spell-cohort-4, #263, #265, #355, #614, #615, #616, #618, #621, #623, #624, #625, #626, #627, #682, #692, #693, #700, #906, #909, #923, #924, #927, #934, #941, #954, #955, #956, #958, #965, #966.
- `spell-cohort-1`: Qualify finite Ranger spell cohort 1. Parent `277`; needs spell-alarm, spell-animal-friendship, spell-animal-messenger-admission, spell-animal-messenger-journey, spell-barkskin, spell-commune-with-nature-facts, spell-conjure-animals-aura, spell-conjure-animals-movement, spell-conjure-animals-damage, spell-conjure-woodland-beings-aura, spell-conjure-woodland-beings-disengage, spell-darkvision, spell-dominate-beast-charm, spell-dominate-beast-control, spell-druidcraft-weather.
- `spell-cohort-2`: Qualify finite Ranger spell cohort 2. Parent `277`; needs spell-druidcraft-nature, spell-elementalism-air-earth, spell-elementalism-fire-water, spell-elementalism-sculpt, spell-enhance-ability, spell-ensnaring-strike-hit, spell-ensnaring-strike-tick-escape, spell-entangle, spell-find-traps, spell-goodberry-create, spell-goodberry-eat, spell-gust-of-wind, spell-hunters-mark-cast, spell-hunters-mark-transfer, spell-jump.
- `spell-cohort-3`: Qualify finite Ranger spell cohort 3. Parent `277`; needs spell-locate-animals-or-plants, spell-longstrider, spell-meld-into-stone-entry, spell-meld-into-stone-inside, spell-meld-into-stone-expulsion, spell-message, spell-nondetection, spell-pass-without-trace, spell-plant-growth-overgrowth, spell-plant-growth-enrichment, spell-poison-spray, spell-produce-flame-light, spell-produce-flame-attack, spell-protection-from-energy, spell-shillelagh.
- `spell-cohort-4`: Qualify finite Ranger spell cohort 4. Parent `277`; needs spell-silence, spell-speak-with-plants-facts, spell-speak-with-plants-terrain, spell-spike-growth-terrain, spell-spike-growth-movement, spell-starry-wisp, spell-stoneskin, spell-water-breathing, spell-water-walk, spell-wind-wall-placement, spell-wind-wall-damage, spell-wind-wall-barrier.
- `execution-dependencies`: Qualify finite Ranger shared execution dependency contracts. Parent `328`; needs #61, #64, #65, #66, #67, #273, #353.
- `hp`: Adapt Ranger own D10 fixed6 HP to shared transaction. Parent `class-cohort`; needs #968.
- `core`: Confirm Ranger first-character core proficiencies. Parent `class-cohort`; needs #430.
- `skills`: Select Ranger three first-character skills. Parent `class-cohort`; needs #430.
- `equipment`: Choose Ranger exact starting equipment. Parent `class-cohort`; needs .
- `later-entry`: Admit exact later-class Ranger entry. Parent `class-cohort`; needs hp, #430.
- `spellcasting`: Prepare Ranger spells with Wisdom and Druidic Focus. Parent `class-cohort`; needs execution-dependencies, membership, slots, spell-alarm, spell-animal-friendship, spell-animal-messenger-admission, spell-animal-messenger-journey, spell-barkskin, spell-commune-with-nature-facts, spell-conjure-animals-aura, spell-conjure-animals-damage, spell-conjure-animals-movement, spell-conjure-woodland-beings-aura, spell-conjure-woodland-beings-disengage, spell-darkvision, spell-dominate-beast-charm, spell-dominate-beast-control, spell-enhance-ability, spell-ensnaring-strike-hit, spell-ensnaring-strike-tick-escape, spell-entangle, spell-find-traps, spell-goodberry-create, spell-goodberry-eat, spell-gust-of-wind, spell-hunters-mark-cast, spell-hunters-mark-transfer, spell-jump, spell-locate-animals-or-plants, spell-longstrider, spell-meld-into-stone-entry, spell-meld-into-stone-expulsion, spell-meld-into-stone-inside, spell-nondetection, spell-pass-without-trace, spell-plant-growth-enrichment, spell-plant-growth-overgrowth, spell-protection-from-energy, spell-silence, spell-speak-with-plants-facts, spell-speak-with-plants-terrain, spell-spike-growth-movement, spell-spike-growth-terrain, spell-stoneskin, spell-water-breathing, spell-water-walk, spell-wind-wall-barrier, spell-wind-wall-damage, spell-wind-wall-placement, #263, #265, #355, #428, #614, #615, #616, #618, #621, #623, #624, #625, #626, #627, #682, #692, #693, #700, #906, #909, #923, #924, #927, #934, #941, #965, #966.
- `favored-enemy`: Cast Favored Enemy exact 5.2.1 Hunter’s Mark. Parent `class-cohort`; needs spell-hunters-mark-cast, spell-hunters-mark-transfer, #428.
- `weapon-mastery`: Select and replace two Ranger mastery kinds. Parent `class-cohort`; needs #428, #489, #490, #491, #492, #493, #494, #495, #496.
- `deft-explorer`: Apply Deft Explorer expertise and languages. Parent `class-cohort`; needs #430.
- `fighting-style`: Choose Ranger Fighting Style or Druidic Warrior. Parent `class-cohort`; needs druidic-cantrips, #384, #385, #386, #387, #388, #389, #390, #391, #392, #393, #394, #395, #396, #397, #398, #399, #400, #401.
- `ranger-subclass`: Select Hunter and exact Ranger subclass timing. Parent `class-cohort`; needs hunter-identity.
- `ability-score-improvement`: Select Ranger qualified ASI at four class levels. Parent `class-cohort`; needs #384, #385, #386, #387, #388, #389, #390, #391, #392, #393, #394, #395, #396, #397, #398, #399, #400, #401.
- `extra-attack`: Resolve Ranger two own-turn Attack attacks. Parent `class-cohort`; needs execution-dependencies.
- `roving`: Apply Roving nonheavy Speed climb and swim. Parent `class-cohort`; needs execution-dependencies.
- `expertise`: Select two Ranger Expertise skills. Parent `class-cohort`; needs #430.
- `relentless-hunter`: Protect Hunter’s Mark concentration against damage only. Parent `class-cohort`; needs execution-dependencies, spell-hunters-mark-cast, #487.
- `natures-veil`: Spend Nature’s Veil uses for exact Invisible duration. Parent `class-cohort`; needs #428, #962.
- `precise-hunter`: Apply Precise Hunter marked-target attack advantage. Parent `class-cohort`; needs execution-dependencies, spell-hunters-mark-transfer.
- `feral-senses`: Resolve Feral Senses exact thirty-foot Blindsight. Parent `class-cohort`; needs execution-dependencies, #962.
- `epic-boon`: Select qualified Ranger Epic Boon at class nineteen. Parent `class-cohort`; needs #384, #385, #386, #387, #388, #389, #390, #391, #392, #393, #394, #395, #396, #397, #398, #399, #400, #401.
- `foe-slayer`: Replace own Hunter’s Mark damage die with D10. Parent `class-cohort`; needs execution-dependencies, spell-hunters-mark-cast.
- `tireless-hp`: Resolve Tireless Magic action temporary HP. Parent `class-cohort`; needs execution-dependencies, #428.
- `tireless-rest`: Resolve Tireless Short Rest Exhaustion reduction. Parent `class-cohort`; needs #427, #623.
- `slots`: Apply Ranger class slot counts and combined multiclass slot table. Parent `class-cohort`; needs execution-dependencies, #428.
- `membership`: Validate exact Ranger source spell list membership. Parent `class-cohort`; needs #293, #315.
- `druidic-cantrips`: Prepare and replace two Druidic Warrior Ranger cantrips. Parent `class-cohort`; needs membership, spell-druidcraft-nature, spell-druidcraft-weather, spell-elementalism-air-earth, spell-elementalism-fire-water, spell-elementalism-sculpt, spell-message, spell-poison-spray, spell-produce-flame-attack, spell-produce-flame-light, spell-shillelagh, spell-starry-wisp, #954, #955, #956, #958.
- `hunter-identity`: Validate chosen Hunter class-three identity. Parent `class-cohort`; needs .
- `hunters-lore`: Disclose Hunter’s Lore exact own marked defenses. Parent `class-cohort`; needs hunter-identity, spell-hunters-mark-transfer.
- `colossus-slayer`: Resolve Colossus Slayer once per turn extra weapon damage. Parent `class-cohort`; needs execution-dependencies, hunter-identity.
- `horde-breaker`: Resolve Horde Breaker same-weapon distinct nearby target. Parent `class-cohort`; needs execution-dependencies, hunter-identity.
- `escape-the-horde`: Apply Escape the Horde Opportunity Attack disadvantage. Parent `class-cohort`; needs execution-dependencies, hunter-identity.
- `multiattack-defense`: Apply Multiattack Defense same-creature remainder-turn disadvantage. Parent `class-cohort`; needs execution-dependencies, hunter-identity.
- `superior-hunters-prey`: Resolve Superior Hunter’s Prey current mark damage to nearby visible creature. Parent `class-cohort`; needs execution-dependencies, foe-slayer, hunter-identity, spell-hunters-mark-cast.
- `superior-hunters-defense`: Resolve Superior Hunter’s Defense reaction resistance through turn end. Parent `class-cohort`; needs execution-dependencies, hunter-identity.
- `hunters-prey`: Choose or rest-replace Hunter’s Prey option. Parent `class-cohort`; needs colossus-slayer, horde-breaker, #427, #428.
- `defensive-tactics`: Choose or rest-replace Defensive Tactics option. Parent `class-cohort`; needs escape-the-horde, multiattack-defense, #427, #428.
- `level-1`: Apply Ranger class level 1 first-character or later-class entry. Parent `levels-cohort`; needs core, equipment, favored-enemy, hp, later-entry, skills, slots, spellcasting, weapon-mastery.
- `level-2`: Apply Ranger class level 2 from class 1. Parent `levels-cohort`; needs deft-explorer, favored-enemy, fighting-style, hp, later-entry, level-1, slots, spellcasting.
- `level-3`: Apply Ranger class level 3 from class 2. Parent `levels-cohort`; needs favored-enemy, hp, hunters-lore, hunters-prey, later-entry, level-2, ranger-subclass, slots, spellcasting.
- `level-4`: Apply Ranger class level 4 from class 3. Parent `levels-cohort`; needs ability-score-improvement, favored-enemy, hp, later-entry, level-3, slots, spellcasting.
- `level-5`: Apply Ranger class level 5 from class 4. Parent `levels-cohort`; needs extra-attack, favored-enemy, hp, later-entry, level-4, slots, spellcasting.
- `level-6`: Apply Ranger class level 6 from class 5. Parent `levels-cohort`; needs favored-enemy, hp, later-entry, level-5, roving, slots, spellcasting.
- `level-7`: Apply Ranger class level 7 from class 6. Parent `levels-cohort`; needs defensive-tactics, favored-enemy, hp, later-entry, level-6, slots, spellcasting.
- `level-8`: Apply Ranger class level 8 from class 7. Parent `levels-cohort`; needs ability-score-improvement, favored-enemy, hp, later-entry, level-7, slots, spellcasting.
- `level-9`: Apply Ranger class level 9 from class 8. Parent `levels-cohort`; needs expertise, favored-enemy, hp, later-entry, level-8, slots, spellcasting.
- `level-10`: Apply Ranger class level 10 from class 9. Parent `levels-cohort`; needs favored-enemy, hp, later-entry, level-9, slots, spellcasting, tireless-hp, tireless-rest.
- `level-11`: Apply Ranger class level 11 from class 10. Parent `levels-cohort`; needs favored-enemy, hp, later-entry, level-10, slots, spellcasting, superior-hunters-prey.
- `level-12`: Apply Ranger class level 12 from class 11. Parent `levels-cohort`; needs ability-score-improvement, favored-enemy, hp, later-entry, level-11, slots, spellcasting.
- `level-13`: Apply Ranger class level 13 from class 12. Parent `levels-cohort`; needs favored-enemy, hp, later-entry, level-12, relentless-hunter, slots, spellcasting.
- `level-14`: Apply Ranger class level 14 from class 13. Parent `levels-cohort`; needs favored-enemy, hp, later-entry, level-13, natures-veil, slots, spellcasting.
- `level-15`: Apply Ranger class level 15 from class 14. Parent `levels-cohort`; needs favored-enemy, hp, later-entry, level-14, slots, spellcasting, superior-hunters-defense.
- `level-16`: Apply Ranger class level 16 from class 15. Parent `levels-cohort`; needs ability-score-improvement, favored-enemy, hp, later-entry, level-15, slots, spellcasting.
- `level-17`: Apply Ranger class level 17 from class 16. Parent `levels-cohort`; needs favored-enemy, hp, later-entry, level-16, precise-hunter, slots, spellcasting.
- `level-18`: Apply Ranger class level 18 from class 17. Parent `levels-cohort`; needs favored-enemy, feral-senses, hp, later-entry, level-17, slots, spellcasting.
- `level-19`: Apply Ranger class level 19 from class 18. Parent `levels-cohort`; needs epic-boon, favored-enemy, hp, later-entry, level-18, slots, spellcasting.
- `level-20`: Apply Ranger class level 20 from class 19. Parent `levels-cohort`; needs favored-enemy, foe-slayer, hp, later-entry, level-19, slots, spellcasting.
- `spell-alarm`: Resolve Alarm effect. Parent `spell-cohort-1`; needs execution-dependencies.
- `spell-animal-friendship`: Resolve Animal Friendship effect. Parent `spell-cohort-1`; needs execution-dependencies, #488.
- `spell-animal-messenger-admission`: Resolve Animal Messenger admission. Parent `spell-cohort-1`; needs execution-dependencies.
- `spell-animal-messenger-journey`: Resolve Animal Messenger journey. Parent `spell-cohort-1`; needs execution-dependencies, spell-animal-messenger-admission.
- `spell-barkskin`: Resolve Barkskin effect. Parent `spell-cohort-1`; needs execution-dependencies.
- `spell-commune-with-nature-facts`: Resolve Commune with Nature facts. Parent `spell-cohort-1`; needs execution-dependencies.
- `spell-conjure-animals-aura`: Resolve Conjure Animals aura. Parent `spell-cohort-1`; needs execution-dependencies.
- `spell-conjure-animals-movement`: Resolve Conjure Animals movement. Parent `spell-cohort-1`; needs execution-dependencies, spell-conjure-animals-aura.
- `spell-conjure-animals-damage`: Resolve Conjure Animals damage. Parent `spell-cohort-1`; needs execution-dependencies, spell-conjure-animals-aura.
- `spell-conjure-woodland-beings-aura`: Resolve Conjure Woodland Beings aura. Parent `spell-cohort-1`; needs execution-dependencies.
- `spell-conjure-woodland-beings-disengage`: Resolve Conjure Woodland Beings disengage. Parent `spell-cohort-1`; needs execution-dependencies, spell-conjure-woodland-beings-aura, #781.
- `spell-darkvision`: Resolve Darkvision effect. Parent `spell-cohort-1`; needs execution-dependencies.
- `spell-dominate-beast-charm`: Resolve Dominate Beast charm. Parent `spell-cohort-1`; needs execution-dependencies, #488.
- `spell-dominate-beast-control`: Resolve Dominate Beast control. Parent `spell-cohort-1`; needs execution-dependencies, spell-dominate-beast-charm, #488.
- `spell-druidcraft-weather`: Resolve Druidcraft weather. Parent `spell-cohort-1`; needs execution-dependencies.
- `spell-druidcraft-nature`: Resolve Druidcraft nature. Parent `spell-cohort-2`; needs execution-dependencies, spell-druidcraft-weather.
- `spell-elementalism-air-earth`: Resolve Elementalism air-earth. Parent `spell-cohort-2`; needs execution-dependencies.
- `spell-elementalism-fire-water`: Resolve Elementalism fire-water. Parent `spell-cohort-2`; needs execution-dependencies, spell-elementalism-air-earth.
- `spell-elementalism-sculpt`: Resolve Elementalism sculpt. Parent `spell-cohort-2`; needs execution-dependencies, spell-elementalism-air-earth.
- `spell-enhance-ability`: Resolve Enhance Ability effect. Parent `spell-cohort-2`; needs execution-dependencies.
- `spell-ensnaring-strike-hit`: Resolve Ensnaring Strike hit. Parent `spell-cohort-2`; needs execution-dependencies, #727.
- `spell-ensnaring-strike-tick-escape`: Resolve Ensnaring Strike tick-escape. Parent `spell-cohort-2`; needs execution-dependencies, spell-ensnaring-strike-hit, #727.
- `spell-entangle`: Resolve Entangle effect. Parent `spell-cohort-2`; needs execution-dependencies, #727.
- `spell-find-traps`: Resolve Find Traps effect. Parent `spell-cohort-2`; needs execution-dependencies.
- `spell-goodberry-create`: Resolve Goodberry create. Parent `spell-cohort-2`; needs execution-dependencies.
- `spell-goodberry-eat`: Resolve Goodberry eat. Parent `spell-cohort-2`; needs execution-dependencies, spell-goodberry-create.
- `spell-gust-of-wind`: Resolve Gust of Wind effect. Parent `spell-cohort-2`; needs execution-dependencies.
- `spell-hunters-mark-cast`: Resolve Hunter’s Mark cast. Parent `spell-cohort-2`; needs execution-dependencies.
- `spell-hunters-mark-transfer`: Resolve Hunter’s Mark transfer. Parent `spell-cohort-2`; needs execution-dependencies, spell-hunters-mark-cast.
- `spell-jump`: Resolve Jump effect. Parent `spell-cohort-2`; needs execution-dependencies.
- `spell-locate-animals-or-plants`: Resolve Locate Animals or Plants effect. Parent `spell-cohort-3`; needs execution-dependencies.
- `spell-longstrider`: Resolve Longstrider effect. Parent `spell-cohort-3`; needs execution-dependencies.
- `spell-meld-into-stone-entry`: Resolve Meld into Stone entry. Parent `spell-cohort-3`; needs execution-dependencies.
- `spell-meld-into-stone-inside`: Resolve Meld into Stone inside. Parent `spell-cohort-3`; needs execution-dependencies, spell-meld-into-stone-entry.
- `spell-meld-into-stone-expulsion`: Resolve Meld into Stone expulsion. Parent `spell-cohort-3`; needs execution-dependencies, spell-meld-into-stone-entry, #488.
- `spell-message`: Resolve Message effect. Parent `spell-cohort-3`; needs execution-dependencies.
- `spell-nondetection`: Resolve Nondetection effect. Parent `spell-cohort-3`; needs execution-dependencies.
- `spell-pass-without-trace`: Resolve Pass without Trace effect. Parent `spell-cohort-3`; needs execution-dependencies.
- `spell-plant-growth-overgrowth`: Resolve Plant Growth overgrowth. Parent `spell-cohort-3`; needs execution-dependencies.
- `spell-plant-growth-enrichment`: Resolve Plant Growth enrichment. Parent `spell-cohort-3`; needs execution-dependencies, spell-plant-growth-overgrowth.
- `spell-poison-spray`: Resolve Poison Spray effect. Parent `spell-cohort-3`; needs execution-dependencies.
- `spell-produce-flame-light`: Resolve Produce Flame light. Parent `spell-cohort-3`; needs execution-dependencies.
- `spell-produce-flame-attack`: Resolve Produce Flame attack. Parent `spell-cohort-3`; needs execution-dependencies, spell-produce-flame-light.
- `spell-protection-from-energy`: Resolve Protection from Energy effect. Parent `spell-cohort-3`; needs execution-dependencies.
- `spell-shillelagh`: Resolve Shillelagh effect. Parent `spell-cohort-3`; needs execution-dependencies.
- `spell-silence`: Resolve Silence effect. Parent `spell-cohort-4`; needs execution-dependencies, #961.
- `spell-speak-with-plants-facts`: Resolve Speak with Plants facts. Parent `spell-cohort-4`; needs execution-dependencies.
- `spell-speak-with-plants-terrain`: Resolve Speak with Plants terrain. Parent `spell-cohort-4`; needs execution-dependencies, spell-speak-with-plants-facts.
- `spell-spike-growth-terrain`: Resolve Spike Growth terrain. Parent `spell-cohort-4`; needs execution-dependencies.
- `spell-spike-growth-movement`: Resolve Spike Growth movement. Parent `spell-cohort-4`; needs execution-dependencies, spell-spike-growth-terrain.
- `spell-starry-wisp`: Resolve Starry Wisp effect. Parent `spell-cohort-4`; needs execution-dependencies, #962.
- `spell-stoneskin`: Resolve Stoneskin effect. Parent `spell-cohort-4`; needs execution-dependencies.
- `spell-water-breathing`: Resolve Water Breathing effect. Parent `spell-cohort-4`; needs execution-dependencies.
- `spell-water-walk`: Resolve Water Walk effect. Parent `spell-cohort-4`; needs execution-dependencies.
- `spell-wind-wall-placement`: Resolve Wind Wall placement. Parent `spell-cohort-4`; needs execution-dependencies.
- `spell-wind-wall-damage`: Resolve Wind Wall damage. Parent `spell-cohort-4`; needs execution-dependencies, spell-wind-wall-placement.
- `spell-wind-wall-barrier`: Resolve Wind Wall barrier. Parent `spell-cohort-4`; needs execution-dependencies, spell-wind-wall-placement.

The dependency DAG and containment DAG are separate. Each shared terminal must
feed a finite spell cohort→shared gate→279/280, and shared gate→class gate→337.
Class/level cohorts feed class gate→337. No new283children or direct shared leaf
root edges. A separately proposed ninth finite execution dependency gate routes
consumers to exact existing primitives to respect outgoing caps; it must never
depend on consuming leaves. Gate closure supplies no substitute implementation.

## Verification and remaining release work

The one released `GOMAXPROCS=2 go test -p 2 -short -race ./...` and
`GOMAXPROCS=2 go vet -p 2 ./...` passed on this base. All eight pinned comparators
passed. Source comparators are rerunnable; the full suite will not be repeated
for narrow publication changes. Go1.27.3 lint shim lacks its binary; no install
or host changes were made. No database, browser/device or production journey
was run; runtime behavior has not changed.

Parent owns serial47/54/283/337 supplements, fresh independent review, final CI
and merge. Final accepted post-Paladin prefix release is required before any
immutable final seal. Keep328 OPEN until its decision PR merges. Source/old
owner bodies remain untouched; no watch, polling, merge, deploy or closure.

Capacity projection from initial live graph plus proposed exact edges:279incoming98,
280incoming93,337incoming55,429outgoing100,283children95,328children5.
Existing primitive outgoing projections61=75,64=41,65=66,66=55,67=92,273=32,
353=68. These are projections, not a published native graph receipt.
